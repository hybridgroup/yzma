package download

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	getter "github.com/hashicorp/go-getter"
)

func TestInstallUsesResolver(t *testing.T) {
	var got []string
	originalGet := getFunc
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		got = append(got, url)
		return nil
	}
	defer func() { getFunc = originalGet }()

	target := Target{Arch: AMD64, OS: Linux, Processor: CUDA, Version: "b7974"}
	resolver := ResolverFunc(func(resolved Target) ([]string, error) {
		if resolved != target {
			t.Errorf("resolver got target %+v, want %+v", resolved, target)
		}
		return []string{"https://example.com/runtime.zip", "https://example.com/llama.tar.gz"}, nil
	})

	if err := Install(context.Background(), target, t.TempDir(), nil, resolver); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}
	// Order matters: an auxiliary runtime lands before the libraries linking it.
	wantURLs := []string{"https://example.com/runtime.zip", "https://example.com/llama.tar.gz"}
	if len(got) != len(wantURLs) {
		t.Fatalf("downloaded %v, want %v", got, wantURLs)
	}
	for i := range wantURLs {
		if got[i] != wantURLs[i] {
			t.Fatalf("downloaded %v, want %v", got, wantURLs)
		}
	}
}

func TestInstallResolverError(t *testing.T) {
	originalGet := getFunc
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		t.Fatal("Install() downloaded despite a resolver error")
		return nil
	}
	defer func() { getFunc = originalGet }()

	sentinel := errors.New("no build for this target")
	resolver := ResolverFunc(func(Target) ([]string, error) { return nil, sentinel })

	err := Install(context.Background(), Target{Arch: AMD64, OS: Linux, Processor: CPU, Version: "b7974"},
		t.TempDir(), nil, resolver)
	if !errors.Is(err, sentinel) {
		t.Fatalf("Install() error = %v, want %v", err, sentinel)
	}
}

func TestInstallNilResolverUsesDefault(t *testing.T) {
	var got []string
	originalGet := getFunc
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		got = append(got, url)
		return nil
	}
	defer func() { getFunc = originalGet }()

	target := Target{Arch: AMD64, OS: Linux, Processor: CPU, Version: "b7974"}
	if err := Install(context.Background(), target, t.TempDir(), nil, nil); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}
	want := "https://github.com/ggml-org/llama.cpp/releases/download/b7974/llama-b7974-bin-ubuntu-x64.tar.gz"
	if len(got) != 1 || got[0] != want {
		t.Fatalf("downloaded %v, want [%s]", got, want)
	}
}

func TestInstallInvalidVersion(t *testing.T) {
	err := Install(context.Background(), Target{Arch: AMD64, OS: Linux, Processor: CPU, Version: "nogood"},
		t.TempDir(), nil, nil)
	if err != ErrInvalidVersion {
		t.Fatalf("Install() error = %v, want %v", err, ErrInvalidVersion)
	}
}

// The Windows CUDA build needs the CUDA runtime archive alongside it.
func TestDefaultResolverWindowsCUDA(t *testing.T) {
	urls, err := DefaultResolver.Resolve(Target{Arch: AMD64, OS: Windows, Processor: CUDA, Version: "b7974"})
	if err != nil {
		t.Fatalf("Resolve() failed: %v", err)
	}
	if len(urls) != 2 {
		t.Fatalf("Resolve() returned %d urls, want 2: %v", len(urls), urls)
	}
	if !strings.Contains(urls[0], "cudart-llama-bin-win-cuda") {
		t.Errorf("first asset = %q, want the CUDA runtime archive", urls[0])
	}
	if !strings.Contains(urls[1], "llama-b7974-bin-win-cuda") {
		t.Errorf("second asset = %q, want the llama.cpp build", urls[1])
	}
}

func TestDefaultResolverUnknownProcessor(t *testing.T) {
	_, err := DefaultResolver.Resolve(Target{Arch: ARM64, OS: Windows, Processor: CUDA, Version: "b7974"})
	if err == nil {
		t.Fatal("Resolve() accepted a target with no published build")
	}
}

// A tagged release has no binaries of its own, so the llama.cpp assets come from the
// nightly build tag while the llama-cpp-builder assets keep the release tag.
func TestDefaultResolverTaggedRelease(t *testing.T) {
	target := Target{Arch: AMD64, OS: Linux, Processor: CPU, Version: "v0.3.0", UpstreamVersion: "b10621"}
	urls, err := DefaultResolver.Resolve(target)
	if err != nil {
		t.Fatalf("Resolve() failed: %v", err)
	}
	want := "https://github.com/ggml-org/llama.cpp/releases/download/b10621/llama-b10621-bin-ubuntu-x64.tar.gz"
	if urls[0] != want {
		t.Errorf("Resolve() = %q, want %q", urls[0], want)
	}

	target.Arch = ARM64
	urls, err = DefaultResolver.Resolve(target)
	if err != nil {
		t.Fatalf("Resolve() failed: %v", err)
	}
	want = "https://github.com/hybridgroup/llama-cpp-builder/releases/download/v0.3.0/llama-v0.3.0-bin-ubuntu-cpu-arm64.tar.gz"
	if urls[0] != want {
		t.Errorf("Resolve() = %q, want %q", urls[0], want)
	}
}

// A nightly tag must not go to the llama.cpp release page, which sometimes rate limits.
func TestInstallNightlyTagSkipsNightlyTagLookup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request for %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalURL := nightlyTagURL
	nightlyTagURL = server.URL + "/%s/nightly-tag.txt"
	defer func() { nightlyTagURL = originalURL }()

	originalGet := getFunc
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		return nil
	}
	defer func() { getFunc = originalGet }()

	target := Target{Arch: AMD64, OS: Linux, Processor: CPU, Version: "b7974"}
	if err := Install(context.Background(), target, t.TempDir(), nil, nil); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}
}

// A pinned version must not ask the version server which build is the newest.
func TestInstallUsesDefaultVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request for %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalCurrent := currentVersionURL
	currentVersionURL = server.URL + "/llama-cpp-builder/version.json"
	defer func() { currentVersionURL = originalCurrent }()

	originalDefault := DefaultVersion
	DefaultVersion = "b7974"
	defer func() { DefaultVersion = originalDefault }()

	var got []string
	originalGet := getFunc
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		got = append(got, url)
		return nil
	}
	defer func() { getFunc = originalGet }()

	target := Target{Arch: AMD64, OS: Linux, Processor: CPU}
	if err := Install(context.Background(), target, t.TempDir(), nil, nil); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	want := "https://github.com/ggml-org/llama.cpp/releases/download/b7974/llama-b7974-bin-ubuntu-x64.tar.gz"
	if len(got) != 1 || got[0] != want {
		t.Fatalf("downloaded %v, want [%s]", got, want)
	}
}

// An explicit version wins over the pinned one.
func TestInstallExplicitVersionOverridesDefault(t *testing.T) {
	originalDefault := DefaultVersion
	DefaultVersion = "v0.3.0"
	defer func() { DefaultVersion = originalDefault }()

	var got []string
	originalGet := getFunc
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		got = append(got, url)
		return nil
	}
	defer func() { getFunc = originalGet }()

	target := Target{Arch: AMD64, OS: Linux, Processor: CPU, Version: "b7974"}
	if err := Install(context.Background(), target, t.TempDir(), nil, nil); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	want := "https://github.com/ggml-org/llama.cpp/releases/download/b7974/llama-b7974-bin-ubuntu-x64.tar.gz"
	if len(got) != 1 || got[0] != want {
		t.Fatalf("downloaded %v, want [%s]", got, want)
	}
}

// "latest" asks for the newest build even when a version is pinned.
func TestInstallLatestSkipsDefaultVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"tag_name": "b8000"}`))
	}))
	defer server.Close()

	originalCurrent := currentVersionURL
	currentVersionURL = server.URL + "/llama-cpp-builder/version.json"
	defer func() { currentVersionURL = originalCurrent }()

	originalDefault := DefaultVersion
	DefaultVersion = "b7974"
	defer func() { DefaultVersion = originalDefault }()

	var got []string
	originalGet := getFunc
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		got = append(got, url)
		return nil
	}
	defer func() { getFunc = originalGet }()

	target := Target{Arch: AMD64, OS: Linux, Processor: CPU, Version: "latest"}
	if err := Install(context.Background(), target, t.TempDir(), nil, nil); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	want := "https://github.com/ggml-org/llama.cpp/releases/download/b8000/llama-b8000-bin-ubuntu-x64.tar.gz"
	if len(got) != 1 || got[0] != want {
		t.Fatalf("downloaded %v, want [%s]", got, want)
	}
}

// A pinned version is already published, so a missing file is a real error and must
// not fall back to the previous build.
func TestInstallDefaultVersionDoesNotFallBack(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request for %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalCurrent := currentVersionURL
	originalPrev := previousVersionURL
	currentVersionURL = server.URL + "/llama-cpp-builder/version.json"
	previousVersionURL = server.URL + "/llama-cpp-builder/previous.json"
	defer func() {
		currentVersionURL = originalCurrent
		previousVersionURL = originalPrev
	}()

	originalDefault := DefaultVersion
	DefaultVersion = "b7974"
	defer func() { DefaultVersion = originalDefault }()

	originalGet := getFunc
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		return ErrFileNotFound
	}
	defer func() { getFunc = originalGet }()

	target := Target{Arch: AMD64, OS: Linux, Processor: CPU}
	err := Install(context.Background(), target, t.TempDir(), nil, nil)
	if !errors.Is(err, ErrFileNotFound) {
		t.Fatalf("Install() error = %v, want %v", err, ErrFileNotFound)
	}
}
