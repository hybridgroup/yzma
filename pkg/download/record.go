package download

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// InstallRecordName is the file that [Install] writes in the library directory to
// record what it installed. [VerifyInstall] reads it.
const InstallRecordName = "yzma-install.json"

// InstallManifestName is the file that [Install] writes in the library directory to save
// the release's digest manifest. [VerifyInstall] reads it instead of fetching one.
const InstallManifestName = "yzma-manifest.json"

// InstallRecord describes which llama.cpp release is installed in a library directory,
// and which assets it came from.
//
// The record sits next to the libraries, so anything that can change them can change
// it too. Give [VerifyInstall] a tag to check against instead.
type InstallRecord struct {
	Version int `json:"version"`

	// Tag is the llama.cpp release that was installed.
	Tag string `json:"tag"`

	// UpstreamTag is the nightly build tag that held the binaries for a tagged
	// release. It is empty for a nightly tag, which names its own assets.
	UpstreamTag string `json:"upstream_tag,omitempty"`

	// Arch, OS and Processor name the build, so a check can resolve the assets
	// again for a different tag.
	Arch      string `json:"arch"`
	OS        string `json:"os"`
	Processor string `json:"processor"`

	// ManifestSHA256 is the digest of the manifest kept in [InstallManifestName], in
	// hexadecimal. It is the pin when the install was pinned.
	ManifestSHA256 string `json:"manifest_sha256,omitempty"`

	Installed time.Time `json:"installed"`

	// Assets are the downloaded assets, in install order.
	Assets []Asset `json:"assets"`
}

// recordVersion is the format of the records this release writes. Version 2 added the
// manifest digest. A version 1 record has no digest and no saved manifest, so
// [VerifyInstall] fetches the manifest as before.
const recordVersion = 2

// WriteInstallRecord writes the record of an install into libPath.
func WriteInstallRecord(libPath string, record InstallRecord) error {
	record.Version = recordVersion

	body, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')

	path := filepath.Join(libPath, InstallRecordName)
	if err := os.WriteFile(path, body, 0644); err != nil {
		return fmt.Errorf("failed to write the install record: %w", err)
	}

	return nil
}

// ReadInstallRecord reads the record that [Install] left in libPath.
func ReadInstallRecord(libPath string) (*InstallRecord, error) {
	body, err := os.ReadFile(filepath.Join(libPath, InstallRecordName))
	if err != nil {
		return nil, err
	}

	var record InstallRecord
	if err := json.Unmarshal(body, &record); err != nil {
		return nil, fmt.Errorf("failed to read the install record: %w", err)
	}

	return &record, nil
}

// WriteInstallManifest saves the raw digest manifest of a release in libPath. The bytes
// are saved unchanged, because a pin is the digest of those exact bytes.
func WriteInstallManifest(libPath string, body []byte) error {
	path := filepath.Join(libPath, InstallManifestName)
	if err := os.WriteFile(path, body, 0644); err != nil {
		return fmt.Errorf("failed to write the install manifest: %w", err)
	}

	return nil
}

// ReadInstallManifest reads the manifest that [Install] left in libPath.
func ReadInstallManifest(libPath string) ([]byte, error) {
	return os.ReadFile(filepath.Join(libPath, InstallManifestName))
}

// target rebuilds the [Target] that made an install.
func (r *InstallRecord) target() (Target, error) {
	arch, err := ParseArch(r.Arch)
	if err != nil {
		return Target{}, ErrUnknownArch
	}

	operatingSystem, err := ParseOS(r.OS)
	if err != nil {
		return Target{}, ErrUnknownOS
	}

	processor, err := ParseProcessor(r.Processor)
	if err != nil {
		return Target{}, ErrUnknownProcessor
	}

	return Target{
		Arch:            arch,
		OS:              operatingSystem,
		Processor:       processor,
		Version:         r.Tag,
		UpstreamVersion: r.UpstreamTag,
	}, nil
}
