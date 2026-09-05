package download

import (
	"context"
	"errors"
	"strings"
	"testing"

	getter "github.com/hashicorp/go-getter"
)

// aDigest is a well formed SHA-256 value that no test asset actually has.
const aDigest = "0000000000000000000000000000000000000000000000000000000000000000"

func TestParsePinnedVersion(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		wantTag    string
		wantDigest string
		wantErr    error
	}{
		{
			name:    "a version with no digest passes through",
			version: "b10783",
			wantTag: "b10783",
		},
		{
			name:    "an empty version passes through",
			version: "",
			wantTag: "",
		},
		{
			name:    "latest with no digest passes through",
			version: "latest",
			wantTag: "latest",
		},
		{
			name:       "a nightly build takes a digest",
			version:    "b10783@sha256:" + aDigest,
			wantTag:    "b10783",
			wantDigest: aDigest,
		},
		{
			name:       "a tagged release takes a digest",
			version:    "v0.3.0@sha256:" + aDigest,
			wantTag:    "v0.3.0",
			wantDigest: aDigest,
		},
		{
			name:       "a digest in capitals becomes lower case",
			version:    "b10783@SHA256:" + strings.ToUpper("abc") + aDigest[3:],
			wantTag:    "b10783",
			wantDigest: "abc" + aDigest[3:],
		},
		{
			name:    "a digest with no algorithm is refused",
			version: "b10783@" + aDigest,
			wantErr: ErrInvalidDigest,
		},
		{
			name:    "an algorithm that is not sha256 is refused",
			version: "b10783@sha512:" + aDigest,
			wantErr: ErrInvalidDigest,
		},
		{
			name:    "a digest that is too short is refused",
			version: "b10783@sha256:abcd",
			wantErr: ErrInvalidDigest,
		},
		{
			name:    "a digest that is not hexadecimal is refused",
			version: "b10783@sha256:" + strings.Repeat("z", 64),
			wantErr: ErrInvalidDigest,
		},
		{
			name:    "an empty digest is refused",
			version: "b10783@sha256:",
			wantErr: ErrInvalidDigest,
		},
		{
			name:    "latest cannot take a digest",
			version: "latest@sha256:" + aDigest,
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "an empty version cannot take a digest",
			version: "@sha256:" + aDigest,
			wantErr: ErrInvalidVersion,
		},
		{
			name:    "a version that is not a release cannot take a digest",
			version: "nightly@sha256:" + aDigest,
			wantErr: ErrInvalidVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, digest, err := ParsePinnedVersion(tt.version)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ParsePinnedVersion(%q) = %v, want %v", tt.version, err, tt.wantErr)
				}
				if tag != "" || digest != "" {
					t.Errorf("got tag %q and digest %q, want both empty on an error", tag, digest)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePinnedVersion(%q) failed: %v", tt.version, err)
			}
			if tag != tt.wantTag {
				t.Errorf("tag = %q, want %q", tag, tt.wantTag)
			}
			if digest != tt.wantDigest {
				t.Errorf("digest = %q, want %q", digest, tt.wantDigest)
			}
		})
	}
}

func TestFetchManifestChecksTheRawBytes(t *testing.T) {
	manifestDigest := serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})

	m, err := fetchManifest(context.Background(), "b10783", manifestDigest)
	if err != nil {
		t.Fatalf("fetchManifest() with the right digest failed: %v", err)
	}
	if m.Tag != "b10783" {
		t.Errorf("Tag = %q, want b10783", m.Tag)
	}

	_, err = fetchManifest(context.Background(), "b10783", aDigest)
	if !errors.Is(err, ErrDigestMismatch) {
		t.Errorf("fetchManifest() with a digest that does not agree = %v, want ErrDigestMismatch", err)
	}
}

// recordAssets makes getFunc collect the assets it is given rather than download
// them, and gives back the collected slice.
func recordAssets(t *testing.T) *[]Asset {
	t.Helper()

	var got []Asset
	original := getFunc
	getFunc = func(_ context.Context, asset Asset, _ string, _ getter.ProgressTracker) error {
		got = append(got, asset)
		return nil
	}
	t.Cleanup(func() { getFunc = original })

	return &got
}

func TestInstallAcceptsAPinnedManifest(t *testing.T) {
	manifestDigest := serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})
	got := recordAssets(t)

	dest := t.TempDir()
	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
		Version: "b10783@sha256:" + manifestDigest,
	}, dest, nil, nil)
	if err != nil {
		t.Fatalf("Install() with a pin that agrees failed: %v", err)
	}

	if len(*got) != 1 {
		t.Fatalf("got %d assets, want 1", len(*got))
	}
	if (*got)[0].SHA256 != "abcd" {
		t.Errorf("SHA256 = %q, want abcd, so the pinned manifest gave the digest", (*got)[0].SHA256)
	}

	// Only the tag reaches the URLs and the record. The pin is not part of either.
	if strings.Contains((*got)[0].URL, "@") {
		t.Errorf("the URL %q holds the pin, want only the tag", (*got)[0].URL)
	}
	record, err := ReadInstallRecord(dest)
	if err != nil {
		t.Fatalf("ReadInstallRecord() failed: %v", err)
	}
	if record.Tag != "b10783" {
		t.Errorf("record Tag = %q, want b10783", record.Tag)
	}
}

func TestInstallAcceptsAPinOnTheTargetField(t *testing.T) {
	manifestDigest := serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})
	got := recordAssets(t)

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
		Version:        "b10783",
		ManifestSHA256: manifestDigest,
	}, t.TempDir(), nil, nil)
	if err != nil {
		t.Fatalf("Install() with ManifestSHA256 set failed: %v", err)
	}
	if len(*got) != 1 || (*got)[0].SHA256 != "abcd" {
		t.Errorf("got %v, want one asset with the digest abcd", *got)
	}
}

func TestInstallStopsOnAPinnedManifestThatDoesNotAgree(t *testing.T) {
	serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})

	original := getFunc
	getFunc = func(context.Context, Asset, string, getter.ProgressTracker) error {
		t.Error("nothing may be downloaded when the manifest does not agree with the pin")
		return nil
	}
	defer func() { getFunc = original }()

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
		Version: "b10783@sha256:" + aDigest,
	}, t.TempDir(), nil, nil)

	if !errors.Is(err, ErrDigestMismatch) {
		t.Errorf("Install() = %v, want ErrDigestMismatch", err)
	}
}

func TestInstallRefusesAPinWhenVerificationIsOff(t *testing.T) {
	original := getFunc
	getFunc = func(context.Context, Asset, string, getter.ProgressTracker) error {
		t.Error("nothing may be downloaded when a pin is given with VerifyOff")
		return nil
	}
	defer func() { getFunc = original }()

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
		Version: "b10783@sha256:" + aDigest,
	}, t.TempDir(), nil, nil, WithVerify(VerifyOff))

	if !errors.Is(err, ErrVerifyDisabled) {
		t.Errorf("Install() = %v, want ErrVerifyDisabled", err)
	}
}

func TestInstallRefusesAPinThatDisagreesWithTheTargetField(t *testing.T) {
	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
		Version:        "b10783@sha256:" + aDigest,
		ManifestSHA256: strings.Repeat("a", 64),
	}, t.TempDir(), nil, nil)

	if !errors.Is(err, ErrInvalidDigest) {
		t.Errorf("Install() = %v, want ErrInvalidDigest", err)
	}
}

func TestInstallRefusesAnAssetThePinnedManifestDoesNotName(t *testing.T) {
	// The manifest names another asset, so the one that resolves has no digest.
	manifestDigest := serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-x64.tar.gz": "abcd",
	})

	original := getFunc
	getFunc = func(context.Context, Asset, string, getter.ProgressTracker) error {
		t.Error("nothing may be downloaded when the pinned manifest omits an asset")
		return nil
	}
	defer func() { getFunc = original }()

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
		Version: "b10783@sha256:" + manifestDigest,
	}, t.TempDir(), nil, nil)

	if !errors.Is(err, ErrDigestMissing) {
		t.Errorf("Install() = %v, want ErrDigestMissing", err)
	}
}

func TestInstallPinsEveryAssetOfATargetThatTakesMoreThanOne(t *testing.T) {
	// A WebAssembly target takes all three browser builds.
	manifestDigest := serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-wasm-simd.tar.gz":    "aaaa",
		"llama-b10783-bin-wasm-simd-mt.tar.gz": "bbbb",
		"llama-b10783-bin-wasm-webgpu.tar.gz":  "cccc",
	})
	got := recordAssets(t)

	err := Install(context.Background(), Target{
		Arch: AMD64, OS: Wasm, Processor: CPU,
		Version: "b10783@sha256:" + manifestDigest,
	}, t.TempDir(), nil, nil)
	if err != nil {
		t.Fatalf("Install() for a WebAssembly target failed: %v", err)
	}

	if len(*got) != 3 {
		t.Fatalf("got %d assets, want 3", len(*got))
	}
	for _, asset := range *got {
		if asset.SHA256 == "" {
			t.Errorf("%s has no digest, want every asset of a pin checked", asset.URL)
		}
	}
}

func TestAPinMakesAManifestThatCannotBeReadFatal(t *testing.T) {
	// digestsURL points at a port that nothing listens on, per TestMain.
	original := getFunc
	getFunc = func(context.Context, Asset, string, getter.ProgressTracker) error {
		t.Error("nothing may be downloaded when a pinned manifest cannot be read")
		return nil
	}
	defer func() { getFunc = original }()

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
		Version: "b10783@sha256:" + aDigest,
	}, t.TempDir(), nil, nil)

	if err == nil {
		t.Error("Install() with a pin and no manifest returned no error")
	}
}

func TestAPlainVersionStillInstallsWithoutAManifest(t *testing.T) {
	// The same conditions as the test above, but with no pin, so the install goes
	// ahead with a warning. This is what keeps an air-gapped mirror working.
	originalWarning := VerifyWarning
	VerifyWarning = nil
	defer func() { VerifyWarning = originalWarning }()

	got := recordAssets(t)

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	}, t.TempDir(), nil, nil)
	if err != nil {
		t.Fatalf("Install() with a plain version failed: %v", err)
	}
	if len(*got) != 1 || (*got)[0].SHA256 != "" {
		t.Errorf("got %v, want one asset with no digest", *got)
	}
}

func TestAPinWorksWithACustomResolver(t *testing.T) {
	manifestDigest := serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})
	got := recordAssets(t)

	// A resolver that names an asset of the release still gets the pinned digest,
	// because the manifest is read by Install and not by the resolver.
	resolver := ResolverFunc(func(target Target) ([]string, error) {
		return []string{
			"https://github.com/hybridgroup/llama-cpp-builder/releases/download/" +
				target.Version + "/llama-" + target.Version + "-bin-ubuntu-cpu-arm64.tar.gz",
		}, nil
	})

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
		Version: "b10783@sha256:" + manifestDigest,
	}, t.TempDir(), nil, resolver)
	if err != nil {
		t.Fatalf("Install() with a custom resolver and a pin failed: %v", err)
	}
	if len(*got) != 1 || (*got)[0].SHA256 != "abcd" {
		t.Errorf("got %v, want one asset with the digest abcd", *got)
	}
}

func TestVerifyInstallRefusesABadPin(t *testing.T) {
	_, err := VerifyInstall(context.Background(), t.TempDir(), "b10783@sha512:"+aDigest)
	if !errors.Is(err, ErrInvalidDigest) {
		t.Errorf("VerifyInstall() = %v, want ErrInvalidDigest", err)
	}
}

func TestInstallTakesThePinOfTheDefaultVersion(t *testing.T) {
	manifestDigest := serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})
	got := recordAssets(t)

	original := DefaultVersion
	DefaultVersion = "b10783@sha256:" + manifestDigest
	defer func() { DefaultVersion = original }()

	dest := t.TempDir()
	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU,
	}, dest, nil, nil)
	if err != nil {
		t.Fatalf("Install() with no version and a pinned default failed: %v", err)
	}

	if len(*got) != 1 || (*got)[0].SHA256 != "abcd" {
		t.Fatalf("got %v, want one asset with the digest from the pinned manifest", *got)
	}

	// Only the tag reaches the URLs and the record.
	if strings.Contains((*got)[0].URL, "@") {
		t.Errorf("the URL %q holds the pin, want only the tag", (*got)[0].URL)
	}
	record, err := ReadInstallRecord(dest)
	if err != nil {
		t.Fatalf("ReadInstallRecord() failed: %v", err)
	}
	if record.Tag != "b10783" {
		t.Errorf("record Tag = %q, want b10783", record.Tag)
	}
}

func TestInstallWithNoPinAndNoManifest(t *testing.T) {
	// Neither manifest URL answers, per TestMain. A version with no digest still
	// installs, because only a pin makes the manifest mandatory.
	got := recordAssets(t)

	warned := 0
	originalWarning := VerifyWarning
	VerifyWarning = func(string) { warned++ }
	defer func() { VerifyWarning = originalWarning }()

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	}, t.TempDir(), nil, nil)
	if err != nil {
		t.Fatalf("Install() of a version with no digest failed: %v", err)
	}

	if len(*got) != 1 || (*got)[0].SHA256 != "" {
		t.Errorf("got %v, want one asset with no digest", *got)
	}
	if warned != 1 {
		t.Errorf("VerifyWarning was called %d times, want 1", warned)
	}
}

func TestInstallWithNoPinTakesTheDigestsItFinds(t *testing.T) {
	serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})
	got := recordAssets(t)

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	}, t.TempDir(), nil, nil)
	if err != nil {
		t.Fatalf("Install() of a version with no digest failed: %v", err)
	}

	if len(*got) != 1 || (*got)[0].SHA256 != "abcd" {
		t.Errorf("got %v, want one asset with the digest from the manifest", *got)
	}
}

func TestDefaultTag(t *testing.T) {
	original := DefaultVersion
	defer func() { DefaultVersion = original }()

	tests := []struct {
		version string
		want    string
	}{
		{version: "", want: ""},
		{version: "b10783", want: "b10783"},
		{version: "b10783@sha256:" + aDigest, want: "b10783"},
		{version: "b10783@sha256:nonsense", want: ""},
	}

	for _, tt := range tests {
		DefaultVersion = tt.version
		if got := DefaultTag(); got != tt.want {
			t.Errorf("DefaultTag() with DefaultVersion %q = %q, want %q", tt.version, got, tt.want)
		}
	}
}
