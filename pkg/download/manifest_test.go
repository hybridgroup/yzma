//go:build manifest

// This file holds the check that the built-in platform table still names assets that
// llama.cpp publishes. It talks to the GitHub API, so it stays behind the "manifest"
// build tag and out of the ordinary test run: run it with `make test-manifest`, or
// with `go test -tags manifest -run TestDefaultResolverMatchesRelease ./pkg/download/`.
//
// Set YZMA_TEST_LLAMA_TAG to check a build other than the newest one.

package download

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// llamaReleaseAssets lists the asset names published for a llama.cpp release tag.
func llamaReleaseAssets(tag string) (map[string]bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/ggml-org/llama.cpp/releases/tags/%s", tag), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("received status code %d for release %s: %s", resp.StatusCode, tag, body)
	}

	var result struct {
		Assets []struct {
			Name string `json:"name"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	names := make(map[string]bool, len(result.Assets))
	for _, asset := range result.Assets {
		names[asset.Name] = true
	}
	return names, nil
}

// TestDefaultResolverMatchesRelease catches upstream filename drift, which tests that
// only assert hard-coded strings cannot see. It checks every target the built-in table
// sends to the llama.cpp release page; the targets served by llama-cpp-builder are
// built by this project and are not part of that release.
func TestDefaultResolverMatchesRelease(t *testing.T) {
	tag := os.Getenv("YZMA_TEST_LLAMA_TAG")
	if tag == "" {
		latest, err := LlamaLatestVersion()
		if err != nil {
			t.Fatalf("LlamaLatestVersion() failed: %v", err)
		}
		tag = latest
	}
	if IsTaggedRelease(tag) {
		upstream, err := LlamaNightlyTag(tag)
		if err != nil {
			t.Fatalf("LlamaNightlyTag(%s) failed: %v", tag, err)
		}
		tag = upstream
	}

	assets, err := llamaReleaseAssets(tag)
	if err != nil {
		t.Fatalf("listing the assets of %s failed: %v", tag, err)
	}
	if len(assets) == 0 {
		t.Fatalf("release %s publishes no assets", tag)
	}

	targets := []Target{
		{Arch: AMD64, OS: Linux, Processor: CPU},
		{Arch: AMD64, OS: Linux, Processor: Vulkan},
		{Arch: AMD64, OS: Linux, Processor: ROCm},
		{Arch: AMD64, OS: Trixie, Processor: CPU},
		{Arch: AMD64, OS: Trixie, Processor: Vulkan},
		{Arch: ARM64, OS: Darwin, Processor: Metal},
		{Arch: ARM64, OS: Darwin, Processor: CPU},
		{Arch: AMD64, OS: Darwin, Processor: CPU},
		{Arch: AMD64, OS: Windows, Processor: CPU},
		{Arch: ARM64, OS: Windows, Processor: CPU},
		{Arch: AMD64, OS: Windows, Processor: CUDA},
		{Arch: AMD64, OS: Windows, Processor: Vulkan},
		{Arch: AMD64, OS: Windows, Processor: ROCm},
	}

	const releasePrefix = "https://github.com/ggml-org/llama.cpp/releases/download/"
	for _, target := range targets {
		target.Version = tag
		name := fmt.Sprintf("%s/%s/%s", target.OS, target.Arch, target.Processor)
		t.Run(name, func(t *testing.T) {
			urls, err := DefaultResolver.Resolve(target)
			if err != nil {
				t.Fatalf("Resolve(%+v) failed: %v", target, err)
			}
			for _, url := range urls {
				if !strings.HasPrefix(url, releasePrefix) {
					continue
				}
				asset := url[strings.LastIndex(url, "/")+1:]
				if !assets[asset] {
					t.Errorf("release %s does not publish %q", tag, asset)
				}
			}
		})
	}
}

// TestDigestManifestCoversResolvedAssets checks that the published digest manifest
// names every asset the built-in table can ask for. An asset the manifest misses
// installs with no check under the default policy, so this catches a manifest that
// falls behind the table.
func TestDigestManifestCoversResolvedAssets(t *testing.T) {
	tag := os.Getenv("YZMA_TEST_LLAMA_TAG")
	if tag == "" {
		latest, err := LlamaLatestVersion()
		if err != nil {
			t.Fatalf("LlamaLatestVersion() failed: %v", err)
		}
		tag = latest
	}

	target := Target{Version: tag}
	if IsTaggedRelease(tag) {
		upstream, err := LlamaNightlyTag(tag)
		if err != nil {
			t.Fatalf("LlamaNightlyTag(%s) failed: %v", tag, err)
		}
		target.UpstreamVersion = upstream
	}

	m, err := fetchManifest(context.Background(), tag, "")
	if err != nil {
		t.Fatalf("fetching the digests of %s failed: %v", tag, err)
	}

	// Every combination the table accepts, so a new platform is covered as soon as
	// the table names it.
	oses := []OS{Linux, Bookworm, Trixie, Darwin, Windows, Wasm}
	arches := []Arch{AMD64, ARM64}
	processors := []Processor{CPU, CUDA, Metal, ROCm, Vulkan, WebGPU}

	seen := make(map[string]bool)
	for _, os := range oses {
		for _, arch := range arches {
			for _, processor := range processors {
				candidate := target
				candidate.Arch, candidate.OS, candidate.Processor = arch, os, processor

				urls, err := DefaultResolver.Resolve(candidate)
				if err != nil {
					// The table does not build this combination.
					continue
				}

				for _, url := range urls {
					if seen[url] {
						continue
					}
					seen[url] = true

					if m.digestFor(url) == "" {
						t.Errorf("the digests of %s do not name %s (%s/%s/%s)",
							tag, url, os, arch, processor)
					}
				}
			}
		}
	}

	if len(seen) == 0 {
		t.Fatal("no assets resolved")
	}
	t.Logf("checked %d assets for %s", len(seen), tag)
}
