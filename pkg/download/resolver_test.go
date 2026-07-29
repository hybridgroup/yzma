package download

import (
	"context"
	"errors"
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
	want := []string{"https://example.com/runtime.zip", "https://example.com/llama.tar.gz"}
	if len(got) != len(want) {
		t.Fatalf("downloaded %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("downloaded %v, want %v", got, want)
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
