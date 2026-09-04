package download

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrDigestMismatch means a downloaded asset does not have the bytes that the
	// publisher recorded for it.
	ErrDigestMismatch = errors.New("digest does not match")

	// ErrDigestMissing means no digest was found for an asset and the policy is
	// [VerifyRequired].
	ErrDigestMissing = errors.New("no digest for asset")

	// ErrInvalidDigest means a pinned digest does not have the form
	// "sha256:" followed by 64 hexadecimal characters.
	ErrInvalidDigest = errors.New("invalid digest")

	// ErrVerifyDisabled means a pinned digest was given with [VerifyOff]. A pin asks
	// for a check, so the two do not agree.
	ErrVerifyDisabled = errors.New("a pinned digest needs verification, which is off")
)

// digestsURL is the URL of the digest manifest for a llama.cpp release tag. It is on
// the same site as the version files, so it does not use the GitHub API, which rate
// limits. See the llama-cpp-builder repo for how the manifests are made.
var digestsURL = "https://hybridgroup.github.io/llama-cpp-builder/digests/%s.json"

// VerifyPolicy says what [Install] does about the digest of an asset.
type VerifyPolicy int

const (
	// VerifyIfAvailable checks an asset that has a digest and permits one that does
	// not. This is the default.
	VerifyIfAvailable VerifyPolicy = iota

	// VerifyRequired makes an asset with no digest an error. Use it where the
	// libraries must be known, such as a deployment that is measured.
	VerifyRequired

	// VerifyOff checks nothing.
	VerifyOff
)

// VerifyWarning is called for each asset that installs with no digest under
// [VerifyIfAvailable]. Set it to nil to say nothing.
var VerifyWarning = func(url string) {
	fmt.Fprintf(os.Stderr, "yzma: no digest for %s, installed with no check\n", url)
}

// Asset is a release asset to install and, when it is known, the digest of its bytes.
type Asset struct {
	// URL is where the asset is downloaded from.
	URL string `json:"url"`

	// SHA256 is the expected digest of the asset, in hexadecimal. An empty value
	// means the digest is not known.
	SHA256 string `json:"sha256,omitempty"`
}

// AssetResolver reports the assets to install for a Target with their expected
// digests. [Install] takes it in place of [Resolver] when a resolver has both, so an
// existing [Resolver] keeps working and gives no digests.
type AssetResolver interface {
	ResolveAssets(target Target) (assets []Asset, err error)
}

// sha256Pattern matches the hexadecimal form of a SHA-256 digest.
var sha256Pattern = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)

// ParsePinnedVersion splits a version that carries the expected digest of its digest
// manifest. It takes either form.
//
//	b10785
//	b10785@sha256:<64 hexadecimal characters>
//
// A version with no digest gives an empty digest and no error, so a caller can pass
// any version through this. What a digest pins is the manifest that names the digest
// of every asset of the release. One version selects a different set of assets for
// each target, so no single archive digest covers them all.
//
// "latest" and an empty version name whichever release is newest at the time, so
// neither can carry a digest.
func ParsePinnedVersion(version string) (tag string, digest string, err error) {
	tag, rest, found := strings.Cut(version, "@")
	if !found {
		return version, "", nil
	}

	algorithm, value, found := strings.Cut(rest, ":")
	switch {
	case !found:
		return "", "", fmt.Errorf("%w: %q names no algorithm, want sha256:<digest>", ErrInvalidDigest, rest)
	case strings.ToLower(algorithm) != "sha256":
		return "", "", fmt.Errorf("%w: unknown algorithm %q, want sha256", ErrInvalidDigest, algorithm)
	case !sha256Pattern.MatchString(value):
		return "", "", fmt.Errorf("%w: %q is not 64 hexadecimal characters", ErrInvalidDigest, value)
	}

	if tag == "" || tag == "latest" {
		return "", "", fmt.Errorf("%w: a digest needs an exact version, not %q", ErrInvalidVersion, tag)
	}
	if err := VersionIsValid(tag); err != nil {
		return "", "", fmt.Errorf("%w: %s", err, tag)
	}

	return tag, strings.ToLower(value), nil
}

// manifest holds the digests that a llama.cpp release publishes. The assets come from
// more than one repository, and an asset name can occur in more than one of them with
// different bytes, so the assets are grouped by the repository that published them.
type manifest struct {
	Version     int                       `json:"version"`
	Tag         string                    `json:"tag"`
	UpstreamTag string                    `json:"upstream_tag"`
	Sources     map[string]manifestSource `json:"sources"`
}

// manifestSource holds the assets of one repository.
type manifestSource struct {
	Tag    string                   `json:"tag"`
	Assets map[string]manifestAsset `json:"assets"`
}

// manifestAsset holds the digest of one asset. Files and Links describe what the
// asset holds, which only the repository that built it can report.
type manifestAsset struct {
	SHA256 string            `json:"sha256"`
	Files  map[string]string `json:"files,omitempty"`
	Links  map[string]string `json:"links,omitempty"`
}

// assetURLPattern matches a GitHub release asset URL, and gives the repository, the
// release tag, and the asset name.
var assetURLPattern = regexp.MustCompile(`^https://github\.com/([^/]+/[^/]+)/releases/download/([^/]+)/([^?]+)`)

// digestFor gives the expected digest of an asset URL, or "" when the manifest does
// not name it. The repository and the tag must both agree, because the same name can
// belong to a different asset in another repository.
func (m *manifest) digestFor(url string) string {
	asset, ok := m.assetFor(url)
	if !ok {
		return ""
	}
	return asset.SHA256
}

// assetFor finds the manifest entry for an asset URL.
func (m *manifest) assetFor(url string) (manifestAsset, bool) {
	parts := assetURLPattern.FindStringSubmatch(url)
	if parts == nil {
		return manifestAsset{}, false
	}

	source, ok := m.Sources[parts[1]]
	if !ok || source.Tag != parts[2] {
		return manifestAsset{}, false
	}

	asset, ok := source.Assets[parts[3]]
	return asset, ok
}

// maxManifestSize is the largest digest manifest that fetchManifest reads. A release
// names about twenty assets, so this is far more than enough.
const maxManifestSize = 8 << 20

// fetchManifest gets the digest manifest for a llama.cpp release tag. A want that is
// not empty is the expected SHA-256 of the raw manifest bytes, in hexadecimal. The
// bytes are checked before they are decoded.
func fetchManifest(ctx context.Context, tag string, want string) (*manifest, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(digestsURL, tag), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d for the digests of %s", resp.StatusCode, tag)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxManifestSize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read the digests of %s: %w", tag, err)
	}
	if len(body) > maxManifestSize {
		return nil, fmt.Errorf("the digests of %s are larger than %d bytes", tag, maxManifestSize)
	}

	if want != "" {
		sum := sha256.Sum256(body)
		got := hex.EncodeToString(sum[:])
		if !strings.EqualFold(got, want) {
			return nil, fmt.Errorf("%w for the digests of %s: expected %s, got %s", ErrDigestMismatch, tag, want, got)
		}
	}

	var m manifest
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

// verifyFile checks that a file has the expected digest. An empty digest checks
// nothing, because [VerifyIfAvailable] permits an asset that has none.
func verifyFile(path, want string) error {
	if want == "" {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open the downloaded file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("failed to read the downloaded file: %w", err)
	}

	got := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(got, want) {
		return fmt.Errorf("%w for %s: expected %s, got %s", ErrDigestMismatch, path, want, got)
	}

	return nil
}

// ParseVerifyPolicy reads a verify policy name. The names are "available" for
// [VerifyIfAvailable], "require" for [VerifyRequired], and "off" for [VerifyOff]. An
// empty name gives the default.
func ParseVerifyPolicy(name string) (VerifyPolicy, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "available", "if-available":
		return VerifyIfAvailable, nil
	case "require", "required":
		return VerifyRequired, nil
	case "off", "none":
		return VerifyOff, nil
	default:
		return VerifyIfAvailable, fmt.Errorf("unknown verify policy: %s", name)
	}
}

// String gives the name of a policy.
func (p VerifyPolicy) String() string {
	switch p {
	case VerifyRequired:
		return "require"
	case VerifyOff:
		return "off"
	default:
		return "available"
	}
}
