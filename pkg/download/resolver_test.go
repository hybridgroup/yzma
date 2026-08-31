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

// A wasm target takes every build from the llama-cpp-builder release page,
// because the JavaScript glue chooses at run time: WebGPU where the browser has
// it, more than one thread where the page is isolated, one thread otherwise.
func TestDefaultResolverWasm(t *testing.T) {
	for _, prcssr := range []Processor{CPU, WebGPU} {
		urls, err := DefaultResolver.Resolve(Target{Arch: AMD64, OS: Wasm, Processor: prcssr, Version: "b7974"})
		if err != nil {
			t.Fatalf("Resolve(%s) failed: %v", prcssr, err)
		}
		if len(urls) != 3 {
			t.Fatalf("Resolve(%s) returned %d urls, want 3: %v", prcssr, len(urls), urls)
		}
		for _, want := range []string{
			"llama-cpp-builder/releases/download/b7974/llama-b7974-bin-wasm-simd-mt.tar.gz",
			"llama-cpp-builder/releases/download/b7974/llama-b7974-bin-wasm-webgpu.tar.gz",
			"llama-cpp-builder/releases/download/b7974/llama-b7974-bin-wasm-simd.tar.gz",
		} {
			found := false
			for _, url := range urls {
				if strings.HasSuffix(url, want) {
					found = true
				}
			}
			if !found {
				t.Errorf("Resolve(%s) = %v, want an asset ending in %q", prcssr, urls, want)
			}
		}
	}
}

// CUDA, Metal, ROCm and Vulkan have no meaning in a browser.
func TestDefaultResolverWasmUnknownProcessor(t *testing.T) {
	for _, prcssr := range []Processor{CUDA, Metal, ROCm, Vulkan} {
		if _, err := DefaultResolver.Resolve(Target{Arch: AMD64, OS: Wasm, Processor: prcssr, Version: "b7974"}); err == nil {
			t.Errorf("Resolve() accepted a wasm target with the %s processor", prcssr)
		}
	}
}

func TestDefaultResolverLibraryNameWasm(t *testing.T) {
	if got := LibraryName(Wasm.String()); got != "yzma_wasm.wasm" {
		t.Errorf("LibraryName(wasm) = %q, want %q", got, "yzma_wasm.wasm")
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

// llama.cpp renamed its ROCm assets at build b10356, so resolution has to follow the
// naming that the requested build actually published.
func TestDefaultResolverROCmNaming(t *testing.T) {
	tests := []struct {
		name    string
		os      OS
		version string
		want    string
	}{
		{"linux before the rename", Linux, "b10355", "llama-b10355-bin-ubuntu-rocm-7.2-x64.tar.gz"},
		{"linux at the rename", Linux, "b10356", "llama-b10356-bin-ubuntu-rocm-7.14-x64.tar.gz"},
		{"linux after the rename", Linux, "b10705", "llama-b10705-bin-ubuntu-rocm-7.14-x64.tar.gz"},
		{"windows before the rename", Windows, "b10355", "llama-b10355-bin-win-hip-radeon-x64.zip"},
		{"windows at the rename", Windows, "b10356", "llama-b10356-bin-win-rocm-7.14-x64.zip"},
		{"windows after the rename", Windows, "b10705", "llama-b10705-bin-win-rocm-7.14-x64.zip"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, err := DefaultResolver.Resolve(Target{
				Arch: AMD64, OS: tt.os, Processor: ROCm, Version: tt.version,
			})
			if err != nil {
				t.Fatalf("Resolve() failed: %v", err)
			}
			want := "https://github.com/ggml-org/llama.cpp/releases/download/" + tt.version + "/" + tt.want
			if len(urls) != 1 || urls[0] != want {
				t.Fatalf("Resolve() = %v, want [%s]", urls, want)
			}
		})
	}
}

// A tagged release carries its binaries under a nightly build tag, so that tag, not
// the release name, decides the ROCm asset names.
func TestDefaultResolverROCmTaggedRelease(t *testing.T) {
	urls, err := DefaultResolver.Resolve(Target{
		Arch: AMD64, OS: Linux, Processor: ROCm, Version: "v0.3.0", UpstreamVersion: "b10355",
	})
	if err != nil {
		t.Fatalf("Resolve() failed: %v", err)
	}
	want := "https://github.com/ggml-org/llama.cpp/releases/download/b10355/llama-b10355-bin-ubuntu-rocm-7.2-x64.tar.gz"
	if len(urls) != 1 || urls[0] != want {
		t.Fatalf("Resolve() = %v, want [%s]", urls, want)
	}
}
