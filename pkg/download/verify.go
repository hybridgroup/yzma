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
	// ErrNoInstallRecord means a library directory has no record of an install, so
	// there is nothing to say which release should be there.
	ErrNoInstallRecord = errors.New("no install record")

	// ErrNoFileDigests means the manifest for a release gives the digest of each
	// archive but not of the files in them, so an installation cannot be checked.
	ErrNoFileDigests = errors.New("the digests of this release do not cover files")

	// ErrRecordMismatch means the install record does not agree with the release it
	// is checked against.
	ErrRecordMismatch = errors.New("the install record does not agree")
)

// FileState is what [VerifyInstall] found for one name.
type FileState int

const (
	// FileVerified means the file on disk has the bytes the publisher recorded.
	FileVerified FileState = iota

	// FileChanged means the file is there and the bytes are not the same.
	FileChanged

	// FileMissing means the install should have put the file there and it is gone.
	FileMissing

	// FileUnexpected means the file is in the directory and no asset of this
	// install holds it. Another install in the same directory makes these.
	FileUnexpected
)

// String gives the name of a state.
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

// OK reports whether every file the install put there is still what the publisher
// recorded. A file that belongs to something else does not make it false.
func (r *VerifyReport) OK() bool {
	return r.Changed == 0 && r.Missing == 0
}

// VerifyInstall checks the files in libPath against the digests that the publisher
// recorded for the release installed there.
//
// An empty tag takes the tag from the install record. Give a tag to name the release
// that must be there, which does not trust the record. The tag may carry the expected
// digest of the digest manifest, in the form "b10785@sha256:<digest>", which does not
// trust the site that serves the manifest either.
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

	// A caller that names a tag gets the assets of that tag, resolved again. The
	// recorded URLs are never used then, because a record that says the wrong tag
	// can name the wrong assets as easily.
	assets := record.Assets
	if tag != "" {
		target.Version = tag
		target.UpstreamVersion = ""
		if IsTaggedRelease(tag) {
			upstream, err := LlamaNightlyTag(tag)
			if err != nil {
				return nil, err
			}
			target.UpstreamVersion = upstream
		}

		urls, err := DefaultResolver.Resolve(target)
		if err != nil {
			return nil, err
		}
		assets = make([]Asset, len(urls))
		for i, url := range urls {
			assets[i] = Asset{URL: url}
		}
	} else {
		tag = record.Tag
	}

	m, err := fetchManifest(ctx, tag, manifestDigest)
	if err != nil {
		return nil, err
	}

	// Gather what the assets of this install should have put in the directory.
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

	// A record that names assets of another release cannot be checked against this
	// one. A record that was changed by hand looks like this.
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

		// The record is not part of any asset.
		if name == InstallRecordName {
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

// add puts one result in the report and counts it.
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

// hashFile gives the SHA-256 of a file, in hexadecimal.
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
