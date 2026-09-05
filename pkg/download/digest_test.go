package download

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	getter "github.com/hashicorp/go-getter"
)

// The same asset name is published by both repositories with different bytes, so a
// manifest keyed only by name would be ambiguous. See llama-cpp-builder b10783.
const collidingName = "llama-b10783-bin-ubuntu-vulkan-arm64.tar.gz"

// testManifest builds a manifest that holds one asset for each source.
func testManifest() *manifest {
	return &manifest{
		Version:     1,
		Tag:         "b10783",
		UpstreamTag: "b10783",
		Sources: map[string]manifestSource{
			"hybridgroup/llama-cpp-builder": {
				Tag: "b10783",
				Assets: map[string]manifestAsset{
					collidingName: {SHA256: "bbbb"},
				},
			},
			"ggml-org/llama.cpp": {
				Tag: "b10783",
				Assets: map[string]manifestAsset{
					collidingName: {SHA256: "gggg"},
				},
			},
		},
	}
}

func TestManifestDigestForUsesTheSource(t *testing.T) {
	m := testManifest()

	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "the builder asset",
			url:  "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b10783/" + collidingName,
			want: "bbbb",
		},
		{
			name: "the upstream asset of the same name",
			url:  "https://github.com/ggml-org/llama.cpp/releases/download/b10783/" + collidingName,
			want: "gggg",
		},
		{
			name: "a repository the manifest does not name",
			url:  "https://github.com/someone/mirror/releases/download/b10783/" + collidingName,
			want: "",
		},
		{
			name: "a tag that does not agree",
			url:  "https://github.com/ggml-org/llama.cpp/releases/download/b10600/" + collidingName,
			want: "",
		},
		{
			name: "an asset the manifest does not name",
			url:  "https://github.com/ggml-org/llama.cpp/releases/download/b10783/llama-b10783-bin-macos-arm64.tar.gz",
			want: "",
		},
		{
			name: "a URL that is not a release asset",
			url:  "https://example.com/llama.tar.gz",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.digestFor(tt.url); got != tt.want {
				t.Errorf("digestFor(%s) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}

func TestVerifyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "asset.tar.gz")
	content := []byte("the bytes of an asset")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	sum := sha256.Sum256(content)
	digest := hex.EncodeToString(sum[:])

	if err := verifyFile(path, digest); err != nil {
		t.Errorf("verifyFile() with the right digest failed: %v", err)
	}

	// A digest is hexadecimal, so the case must not matter.
	if err := verifyFile(path, strings.ToUpper(digest)); err != nil {
		t.Errorf("verifyFile() with an upper case digest failed: %v", err)
	}

	// An empty digest checks nothing, which VerifyIfAvailable permits.
	if err := verifyFile(path, ""); err != nil {
		t.Errorf("verifyFile() with no digest failed: %v", err)
	}

	err := verifyFile(path, strings.Repeat("0", 64))
	if !errors.Is(err, ErrDigestMismatch) {
		t.Errorf("verifyFile() with a wrong digest = %v, want ErrDigestMismatch", err)
	}
}

func TestFetchManifest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/digests/b10783.json" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprint(w, `{"version":1,"tag":"b10783","upstream_tag":"b10783","sources":{
			"hybridgroup/llama-cpp-builder":{"tag":"b10783","assets":{"a.tar.gz":{"sha256":"abcd"}}}}}`)
	}))
	defer server.Close()

	original := digestsURL
	digestsURL = server.URL + "/digests/%s.json"
	defer func() { digestsURL = original }()

	m, err := fetchManifest(context.Background(), "b10783", "")
	if err != nil {
		t.Fatalf("fetchManifest() failed: %v", err)
	}
	if m.Tag != "b10783" {
		t.Errorf("Tag = %q, want b10783", m.Tag)
	}
	want := "abcd"
	url := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b10783/a.tar.gz"
	if got := m.digestFor(url); got != want {
		t.Errorf("digestFor() = %q, want %q", got, want)
	}

	if _, err := fetchManifest(context.Background(), "b00000", ""); err == nil {
		t.Error("fetchManifest() for a tag with no manifest returned no error")
	}
}

// serveManifest starts a server that gives a manifest naming the assets in digests,
// and points both manifest URLs at it for the length of the test, as a release that
// publishes both copies does. It gives back the SHA-256 of the manifest bytes, in
// hexadecimal, which a test pins with [Target.ManifestSHA256].
func serveManifest(t *testing.T, tag string, digests map[string]string) string {
	t.Helper()

	assets := make([]string, 0, len(digests))
	for name, digest := range digests {
		assets = append(assets, fmt.Sprintf("%q:{\"sha256\":%q}", name, digest))
	}
	body := fmt.Sprintf(`{"version":1,"tag":%q,"upstream_tag":%q,"sources":{
		"hybridgroup/llama-cpp-builder":{"tag":%q,"assets":{%s}},
		"ggml-org/llama.cpp":{"tag":%q,"assets":{%s}}}}`,
		tag, tag, tag, strings.Join(assets, ","), tag, strings.Join(assets, ","))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)

	originalAsset, originalPages := manifestAssetURL, digestsURL
	manifestAssetURL = server.URL + "/releases/download/%[1]s/%[1]s.json"
	digestsURL = server.URL + "/digests/%s.json"
	t.Cleanup(func() { manifestAssetURL, digestsURL = originalAsset, originalPages })

	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
}

func TestDefaultResolverReportsDigests(t *testing.T) {
	serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})

	assets, err := defaultResolver{}.ResolveAssets(Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	})
	if err != nil {
		t.Fatalf("ResolveAssets() failed: %v", err)
	}
	if len(assets) != 1 {
		t.Fatalf("got %d assets, want 1", len(assets))
	}
	if assets[0].SHA256 != "abcd" {
		t.Errorf("SHA256 = %q, want abcd", assets[0].SHA256)
	}
}

func TestDefaultResolverWithNoManifest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	original := digestsURL
	digestsURL = server.URL + "/digests/%s.json"
	defer func() { digestsURL = original }()

	// A manifest that is not there must not stop the resolver. The policy decides
	// what happens to an asset that has no digest.
	assets, err := defaultResolver{}.ResolveAssets(Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	})
	if err != nil {
		t.Fatalf("ResolveAssets() failed: %v", err)
	}
	if len(assets) != 1 || assets[0].SHA256 != "" {
		t.Errorf("got %v, want one asset with no digest", assets)
	}
}

func TestResolveAssetsFromAPlainResolver(t *testing.T) {
	resolver := ResolverFunc(func(Target) ([]string, error) {
		return []string{"https://example.com/a.tar.gz"}, nil
	})

	assets, _, err := resolveAssets(context.Background(), Target{Version: "b10783"}, resolver, VerifyIfAvailable)
	if err != nil {
		t.Fatalf("resolveAssets() failed: %v", err)
	}
	if len(assets) != 1 || assets[0].URL != "https://example.com/a.tar.gz" || assets[0].SHA256 != "" {
		t.Errorf("got %v, want one asset with no digest", assets)
	}
}

func TestInstallStopsOnADigestThatDoesNotAgree(t *testing.T) {
	body := createMockTarGz(t, "b10783")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(body)
	}))
	defer server.Close()

	dest := t.TempDir()
	asset := Asset{
		URL:    server.URL + "/llama-b10783-bin-ubuntu-cpu-arm64.tar.gz",
		SHA256: strings.Repeat("0", 64),
	}

	err := get(context.Background(), asset, dest, nil)
	if !errors.Is(err, ErrDigestMismatch) {
		t.Fatalf("get() = %v, want ErrDigestMismatch", err)
	}

	// Nothing may be extracted, and the archive must not stay behind.
	entries, err := os.ReadDir(dest)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("the destination holds %v, want it empty", names)
	}
}

// A destination whose name holds "404" must not make a digest that does not agree
// look like a file that is not there.
func TestInstallStopsOnADigestThatDoesNotAgreeUnderA404Path(t *testing.T) {
	body := createMockTarGz(t, "b10783")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(body)
	}))
	defer server.Close()

	dest := filepath.Join(t.TempDir(), "build-404")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}

	asset := Asset{
		URL:    server.URL + "/llama-b10783-bin-ubuntu-cpu-arm64.tar.gz",
		SHA256: strings.Repeat("0", 64),
	}

	err := get(context.Background(), asset, dest, nil)
	if !errors.Is(err, ErrDigestMismatch) {
		t.Fatalf("get() = %v, want ErrDigestMismatch", err)
	}
	if errors.Is(err, ErrFileNotFound) {
		t.Errorf("get() = %v, want no ErrFileNotFound", err)
	}
}

func TestInstallAcceptsADigestThatAgrees(t *testing.T) {
	body := createMockTarGz(t, "b10783")
	sum := sha256.Sum256(body)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(body)
	}))
	defer server.Close()

	dest := t.TempDir()
	asset := Asset{
		URL:    server.URL + "/llama-b10783-bin-ubuntu-cpu-arm64.tar.gz",
		SHA256: hex.EncodeToString(sum[:]),
	}

	if err := get(context.Background(), asset, dest, nil); err != nil {
		t.Fatalf("get() failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "libllama.so")); err != nil {
		t.Errorf("the library was not extracted: %v", err)
	}
}

func TestInstallRequiresADigest(t *testing.T) {
	originalGet := getFunc
	getFunc = func(context.Context, Asset, string, getter.ProgressTracker) error {
		t.Error("nothing may be downloaded when a digest is required and missing")
		return nil
	}
	defer func() { getFunc = originalGet }()

	resolver := ResolverFunc(func(Target) ([]string, error) {
		return []string{"https://example.com/a.tar.gz"}, nil
	})

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	}, t.TempDir(), nil, resolver, WithVerify(VerifyRequired))

	if !errors.Is(err, ErrDigestMissing) {
		t.Errorf("Install() = %v, want ErrDigestMissing", err)
	}
}

func TestInstallWarnsWhenAnAssetHasNoDigest(t *testing.T) {
	originalGet := getFunc
	getFunc = func(context.Context, Asset, string, getter.ProgressTracker) error { return nil }
	defer func() { getFunc = originalGet }()

	var warned []string
	originalWarning := VerifyWarning
	VerifyWarning = func(url string) { warned = append(warned, url) }
	defer func() { VerifyWarning = originalWarning }()

	resolver := ResolverFunc(func(Target) ([]string, error) {
		return []string{"https://example.com/a.tar.gz"}, nil
	})

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	}, t.TempDir(), nil, resolver)
	if err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	if len(warned) != 1 || warned[0] != "https://example.com/a.tar.gz" {
		t.Errorf("warned about %v, want the one asset", warned)
	}
}

func TestVerifyOffSkipsTheCheck(t *testing.T) {
	var got []Asset
	originalGet := getFunc
	getFunc = func(_ context.Context, asset Asset, _ string, _ getter.ProgressTracker) error {
		got = append(got, asset)
		return nil
	}
	defer func() { getFunc = originalGet }()

	serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	}, t.TempDir(), nil, nil, WithVerify(VerifyOff))
	if err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	if len(got) != 1 || got[0].SHA256 != "" {
		t.Errorf("got %v, want the asset with no digest to check", got)
	}
}

func TestVerifyOffDoesNotFetchAManifest(t *testing.T) {
	originalGet := getFunc
	getFunc = func(context.Context, Asset, string, getter.ProgressTracker) error { return nil }
	defer func() { getFunc = originalGet }()

	var fetched int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetched++
		fmt.Fprint(w, `{"version":1,"sources":{}}`)
	}))
	defer server.Close()

	original := digestsURL
	digestsURL = server.URL + "/digests/%s.json"
	defer func() { digestsURL = original }()

	err := Install(context.Background(), Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783",
	}, t.TempDir(), nil, nil, WithVerify(VerifyOff))
	if err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// An air-gapped install that checks nothing must not wait on a manifest.
	if fetched != 0 {
		t.Errorf("fetched the manifest %d times, want 0", fetched)
	}
}

func TestFetchManifestPrefersTheReleaseAsset(t *testing.T) {
	// The two copies hold the same bytes in a real release. They differ here only so
	// that the manifest says which one answered.
	asset := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"version":1,"tag":"b10783","upstream_tag":"asset"}`)
	}))
	defer asset.Close()
	pages := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"version":1,"tag":"b10783","upstream_tag":"pages"}`)
	}))
	defer pages.Close()

	originalAsset, originalPages := manifestAssetURL, digestsURL
	manifestAssetURL = asset.URL + "/releases/download/%[1]s/%[1]s.json"
	digestsURL = pages.URL + "/digests/%s.json"
	defer func() { manifestAssetURL, digestsURL = originalAsset, originalPages }()

	m, err := fetchManifest(context.Background(), "b10783", "")
	if err != nil {
		t.Fatalf("fetchManifest() failed: %v", err)
	}
	if m.UpstreamTag != "asset" {
		t.Errorf("the manifest came from %q, want the release asset", m.UpstreamTag)
	}
}

func TestFetchManifestFallsBackToTheVersionSite(t *testing.T) {
	// A release published before the manifest was an asset has only the site copy.
	asset := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer asset.Close()
	pages := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"version":1,"tag":"b10783","upstream_tag":"pages"}`)
	}))
	defer pages.Close()

	originalAsset, originalPages := manifestAssetURL, digestsURL
	manifestAssetURL = asset.URL + "/releases/download/%[1]s/%[1]s.json"
	digestsURL = pages.URL + "/digests/%s.json"
	defer func() { manifestAssetURL, digestsURL = originalAsset, originalPages }()

	m, err := fetchManifest(context.Background(), "b10783", "")
	if err != nil {
		t.Fatalf("fetchManifest() failed: %v", err)
	}
	if m.UpstreamTag != "pages" {
		t.Errorf("the manifest came from %q, want the version site", m.UpstreamTag)
	}

	// A manifest in neither place is an error, and only then.
	digestsURL = asset.URL + "/digests/%s.json"
	if _, err := fetchManifest(context.Background(), "b10783", ""); err == nil {
		t.Error("fetchManifest() with no manifest anywhere returned no error")
	}
}

// serveVersionFile points currentVersionURL and previousVersionURL at a server that
// gives the two bodies for the length of the test.
func serveVersionFile(t *testing.T, current, previous string) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := current
		if strings.Contains(r.URL.Path, "previous") {
			body = previous
		}
		if body == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprint(w, body)
	}))
	t.Cleanup(server.Close)

	originalCurrent, originalPrevious := currentVersionURL, previousVersionURL
	currentVersionURL = server.URL + "/version.json"
	previousVersionURL = server.URL + "/previous.json"
	t.Cleanup(func() { currentVersionURL, previousVersionURL = originalCurrent, originalPrevious })
}

func TestManifestDigestFromTheVersionFiles(t *testing.T) {
	digest := strings.Repeat("ab", 32)

	tests := []struct {
		name     string
		current  string
		previous string
		tag      string
	}{
		{
			name:    "the digest field of the current version",
			current: fmt.Sprintf(`{"tag_name":"b10783","manifest_sha256":%q}`, digest),
			tag:     "b10783",
		},
		{
			name:    "the pin of the current version, for a file with no digest field",
			current: fmt.Sprintf(`{"tag_name":"b10783","pin":"b10783@sha256:%s"}`, digest),
			tag:     "b10783",
		},
		{
			name:     "the previous version",
			current:  `{"tag_name":"b10784"}`,
			previous: fmt.Sprintf(`{"tag_name":"b10783","manifest_sha256":%q}`, digest),
			tag:      "b10783",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serveVersionFile(t, tt.current, tt.previous)

			got, err := ManifestDigest(context.Background(), tt.tag)
			if err != nil {
				t.Fatalf("ManifestDigest() failed: %v", err)
			}
			if got != digest {
				t.Errorf("ManifestDigest() = %q, want %q", got, digest)
			}

			want := tt.tag + "@sha256:" + digest
			if pin, err := PinnedVersion(context.Background(), tt.tag); err != nil || pin != want {
				t.Errorf("PinnedVersion() = %q, %v, want %q", pin, err, want)
			}
		})
	}
}

func TestManifestDigestFromTheRelease(t *testing.T) {
	digest := strings.Repeat("cd", 32)

	// The version files name another tag, so the release has the answer.
	serveVersionFile(t, `{"tag_name":"b10784"}`, `{"tag_name":"b10782"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"assets":[
			{"name":"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz","digest":"sha256:%s"},
			{"name":"b10783.json","digest":"sha256:%s"}]}`, strings.Repeat("ef", 32), digest)
	}))
	defer server.Close()

	original := releaseAPIURL
	releaseAPIURL = server.URL + "/releases/tags/%s"
	defer func() { releaseAPIURL = original }()

	got, err := ManifestDigest(context.Background(), "b10783")
	if err != nil {
		t.Fatalf("ManifestDigest() failed: %v", err)
	}
	if got != digest {
		t.Errorf("ManifestDigest() = %q, want the digest of the manifest asset %q", got, digest)
	}
}

func TestManifestDigestWithNoManifestPublished(t *testing.T) {
	serveVersionFile(t, `{"tag_name":"b10784"}`, `{"tag_name":"b10782"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// An asset that is still uploading has no digest yet, and an older release
		// has no manifest asset at all.
		fmt.Fprint(w, `{"assets":[
			{"name":"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz","digest":"sha256:abcd"},
			{"name":"b10783.json","digest":null}]}`)
	}))
	defer server.Close()

	original := releaseAPIURL
	releaseAPIURL = server.URL + "/releases/tags/%s"
	defer func() { releaseAPIURL = original }()

	_, err := ManifestDigest(context.Background(), "b10783")
	if !errors.Is(err, ErrNoManifestDigest) {
		t.Errorf("ManifestDigest() = %v, want ErrNoManifestDigest", err)
	}

	if _, err := ManifestDigest(context.Background(), "not-a-version"); !errors.Is(err, ErrInvalidVersion) {
		t.Errorf("ManifestDigest() with a bad tag = %v, want ErrInvalidVersion", err)
	}
}
