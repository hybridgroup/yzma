package download

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	// ErrNoInstallRecord means a library directory has no install record, so
	// nothing says which release should be there.
	ErrNoInstallRecord = errors.New("no install record")

	// ErrNoFileDigests means the release manifest has a digest for each archive
	// but not for the files inside, so the installation cannot be checked.
	ErrNoFileDigests = errors.New("the digests of this release do not cover files")

	// ErrRecordMismatch means the install record does not match the release it
	// is checked against.
	ErrRecordMismatch = errors.New("the install record does not agree")
)

// FileState is what [VerifyInstall] found for one name.
type FileState int

const (
	// FileVerified means the file on disk has the bytes the publisher recorded.
	FileVerified FileState = iota

	// FileChanged means the file exists but its bytes are different.
	FileChanged

	// FileMissing means the install should have created the file but it is gone.
	FileMissing

	// FileUnexpected means the file is in the directory but no asset of this
	// install contains it. Another install in the same directory causes these.
	FileUnexpected
)

// String returns the name of a state.
func (s FileState) String() string {
	switch s {
	case FileVerified:
		return "verified"
	case FileChanged:
		return "changed"
	case FileMissing:
		return "missing"
	case FileUnexpected:
		return "unexpected"
	default:
		return "unknown"
	}
}

// FileReport is what [VerifyInstall] found for one name.
type FileReport struct {
	Name  string    `json:"name"`
	State FileState `json:"state"`
}

// MarshalJSON writes the state as its name.
func (s FileState) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

// VerifyReport is what [VerifyInstall] found in a library directory.
type VerifyReport struct {
	// Tag is the release the files were checked against.
	Tag string `json:"tag"`

	// LibPath is the directory that was checked.
	LibPath string `json:"lib_path"`

	// Files holds one entry for each name, sorted by name.
	Files []FileReport `json:"files"`

	// Counts of each state.
	Verified   int `json:"verified"`
	Changed    int `json:"changed"`
	Missing    int `json:"missing"`
	Unexpected int `json:"unexpected"`
}

// OK reports whether every installed file still matches what the publisher recorded.
// Files that belong to something else do not make it false.
func (r *VerifyReport) OK() bool {
	return r.Changed == 0 && r.Missing == 0
}

// VerifyInstall checks the files in libPath against the digests that the publisher
// recorded for the release installed there.
//
// An empty tag uses the tag from the install record. Pass a tag to name the release
// that must be installed, without trusting the record. The tag may carry the expected
// digest of the manifest, in the form "b10785@sha256:<digest>", so the site that
// serves the manifest is not trusted either.
//
// An install saves its release manifest next to the record, so a check needs no
// network. The saved manifest is checked against the same digest, and if it is
// missing or does not match it is fetched as before.
func VerifyInstall(ctx context.Context, libPath, tag string) (*VerifyReport, error) {
	tag, manifestDigest, err := ParsePinnedVersion(tag)
	if err != nil {
		return nil, err
	}

	record, err := ReadInstallRecord(libPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w in %s", ErrNoInstallRecord, libPath)
		}
		return nil, err
	}

	target, err := record.target()
	if err != nil {
		return nil, err
	}

	named := tag != ""
	if !named {
		tag = record.Tag
	}

	m, err := installManifest(ctx, libPath, record, tag, manifestDigest)
	if err != nil {
		return nil, err
	}

	// When the caller names a tag, resolve that tag's assets again. The recorded
	// URLs are not used, because a record with the wrong tag can just as easily
	// list the wrong assets.
	assets := record.Assets
	if named {
		target.Version = tag
		target.UpstreamVersion = ""
		if IsTaggedRelease(tag) {
			// The manifest names the nightly build that holds the binaries for a
			// tagged release, so the release page is not needed.
			target.UpstreamVersion = m.UpstreamTag
			if target.UpstreamVersion == "" {
				upstream, err := LlamaNightlyTag(tag)
				if err != nil {
					return nil, err
				}
				target.UpstreamVersion = upstream
			}
		}

		urls, err := DefaultResolver.Resolve(target)
		if err != nil {
			return nil, err
		}
		assets = make([]Asset, len(urls))
		for i, url := range urls {
			assets[i] = Asset{URL: url}
		}
	}

	// Collect the files that the assets of this install should have created.
	wantFiles := make(map[string]string)
	wantLinks := make(map[string]string)
	found := 0
	for _, asset := range assets {
		entry, ok := m.assetFor(asset.URL)
		if !ok {
			continue
		}
		found++
		for name, digest := range entry.Files {
			wantFiles[name] = digest
		}
		for name, destination := range entry.Links {
			wantLinks[name] = destination
		}
	}

	// A record that lists assets of another release cannot be checked against this
	// one. This is what a hand edited record looks like.
	if found == 0 {
		return nil, fmt.Errorf("%w: the install record names assets that %s does not publish",
			ErrRecordMismatch, tag)
	}

	if len(wantFiles) == 0 && len(wantLinks) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrNoFileDigests, tag)
	}

	report := &VerifyReport{Tag: tag, LibPath: libPath}
	seen := make(map[string]bool)

	err = filepath.WalkDir(libPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		name, err := filepath.Rel(libPath, path)
		if err != nil {
			return err
		}
		name = filepath.ToSlash(name)

		// The record and the manifest are not part of any asset.
		if name == InstallRecordName || name == InstallManifestName {
			return nil
		}

		seen[name] = true

		if entry.Type()&fs.ModeSymlink != 0 {
			want, ok := wantLinks[name]
			if !ok {
				report.add(name, FileUnexpected)
				return nil
			}
			got, err := os.Readlink(path)
			if err != nil {
				return err
			}
			if got != want {
				report.add(name, FileChanged)
				return nil
			}
			report.add(name, FileVerified)
			return nil
		}

		want, ok := wantFiles[name]
		if !ok {
			report.add(name, FileUnexpected)
			return nil
		}

		got, err := hashFile(path)
		if err != nil {
			return err
		}
		if !strings.EqualFold(got, want) {
			report.add(name, FileChanged)
			return nil
		}
		report.add(name, FileVerified)
		return nil
	})
	if err != nil {
		return nil, err
	}

	for name := range wantFiles {
		if !seen[name] {
			report.add(name, FileMissing)
		}
	}
	for name := range wantLinks {
		if !seen[name] {
			report.add(name, FileMissing)
		}
	}

	sort.Slice(report.Files, func(i, j int) bool {
		return report.Files[i].Name < report.Files[j].Name
	})

	return report, nil
}

// installManifest returns the digest manifest of a release. It tries the copy saved next
// to the record first, so a check needs no network. If that copy is missing or does not
// match the digest, it fetches the manifest and saves it for the next check.
func installManifest(ctx context.Context, libPath string, record *InstallRecord, tag, want string) (*manifest, error) {
	// Without a digest from the caller, use the one the install recorded, so a
	// manifest that changed on disk is not trusted.
	cached := want
	if cached == "" {
		cached = record.ManifestSHA256
	}
	if m, _, ok := loadCachedManifest(libPath, tag, cached); ok {
		return m, nil
	}

	m, body, err := fetchManifestBody(ctx, tag, want)
	if err != nil {
		return nil, err
	}

	cacheManifest(libPath, record, tag, body)

	return m, nil
}

// cacheManifest saves a manifest next to the install record. It only speeds up the next
// check, so an unwritable directory is not an error.
func cacheManifest(libPath string, record *InstallRecord, tag string, body []byte) {
	if record.Tag != tag {
		return
	}
	if err := WriteInstallManifest(libPath, body); err != nil {
		return
	}

	sum := sha256.Sum256(body)
	digest := hex.EncodeToString(sum[:])
	if strings.EqualFold(record.ManifestSHA256, digest) {
		return
	}

	record.ManifestSHA256 = digest
	_ = WriteInstallRecord(libPath, *record)
}

// add adds one result to the report and counts it.
func (r *VerifyReport) add(name string, state FileState) {
	r.Files = append(r.Files, FileReport{Name: name, State: state})

	switch state {
	case FileVerified:
		r.Verified++
	case FileChanged:
		r.Changed++
	case FileMissing:
		r.Missing++
	case FileUnexpected:
		r.Unexpected++
	}
}

// hashFile returns the SHA-256 of a file, in hexadecimal.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
