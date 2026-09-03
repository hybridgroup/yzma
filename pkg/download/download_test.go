package download

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	getter "github.com/hashicorp/go-getter"
)

func TestLlamaLatestVersion(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/llama-cpp-builder/version.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Return a mock response
		response := struct {
			TagName string `json:"tag_name"`
		}{
			TagName: "b7974",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Override the version URL for testing
	originalURL := currentVersionURL
	currentVersionURL = server.URL + "/llama-cpp-builder/version.json"
	defer func() { currentVersionURL = originalURL }()

	version, err := LlamaLatestVersion()
	if err != nil {
		t.Fatal("could not get latest version", err)
	}

	if !strings.HasPrefix(version, "b") {
		t.Fatalf("Expected version should start with 'b', got '%s'", version)
	}

	if version != "b7974" {
		t.Fatalf("Expected version 'b7974', got '%s'", version)
	}

	t.Logf("LlamaLatestVersion returned: %s", version)
}

func TestLlamaLatestVersion_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message": "API rate limit exceeded"}`))
	}))
	defer server.Close()

	// Override the version URL for testing
	originalURL := currentVersionURL
	currentVersionURL = server.URL + "/llama-cpp-builder/version.json"
	defer func() { currentVersionURL = originalURL }()

	// Reduce retry count for faster test
	originalRetryCount := RetryCount
	originalRetryDelay := RetryDelay
	RetryCount = 1
	RetryDelay = 0
	defer func() {
		RetryCount = originalRetryCount
		RetryDelay = originalRetryDelay
	}()

	_, err := LlamaLatestVersion()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// createMockTarGz creates a mock .tar.gz file containing a fake libllama.so
// with a top-level directory prefix (e.g., "llama-b7974/")
func createMockTarGz(t *testing.T, version string) []byte {
	t.Helper()

	var buf strings.Builder
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	prefix := "llama-" + version + "/"

	// Add the top-level directory
	hdr := &tar.Header{
		Name:     prefix,
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}

	// Add a fake libllama.so file
	content := []byte("fake library content")
	hdr = &tar.Header{
		Name: prefix + "libllama.so",
		Mode: 0755,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("failed to write tar content: %v", err)
	}

	// Add a subdirectory with another file
	hdr = &tar.Header{
		Name:     prefix + "lib/",
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}

	content2 := []byte("another file")
	hdr = &tar.Header{
		Name: prefix + "lib/libggml.so",
		Mode: 0755,
		Size: int64(len(content2)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write(content2); err != nil {
		t.Fatalf("failed to write tar content: %v", err)
	}

	tw.Close()
	gzw.Close()

	return []byte(buf.String())
}

func TestGetLinuxCPU(t *testing.T) {
	version := "b7974"

	// Create mock tar.gz content with version prefix
	mockTarGz := createMockTarGz(t, version)

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/b7974/llama-b7974-bin-ubuntu-x64.tar.gz"
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s, want %s", r.URL.Path, expectedPath)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/gzip")
		w.Write(mockTarGz)
	}))
	defer server.Close()

	arch := "amd64"
	osVer := "linux"
	processor := "cpu"
	dest := t.TempDir()

	// Override the get function to use our mock server
	originalGet := getFunc
	getFunc = func(ctx context.Context, asset Asset, dest string, progress getter.ProgressTracker) error {
		// Replace the real URL with our mock server URL
		mockURL := server.URL + "/b7974/llama-b7974-bin-ubuntu-x64.tar.gz"
		return downloadAndExtractTarGz(Asset{URL: mockURL}, dest, nil)
	}
	defer func() { getFunc = originalGet }()

	err := Get(arch, osVer, processor, version, dest)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// Check that files were extracted without the prefix directory
	expectedFile := filepath.Join(dest, "libllama.so")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", expectedFile)
	}

	// Also check the subdirectory file
	expectedFile2 := filepath.Join(dest, "lib", "libggml.so")
	if _, err := os.Stat(expectedFile2); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", expectedFile2)
	}

	// Verify the prefix directory was NOT created
	prefixDir := filepath.Join(dest, "llama-"+version)
	if _, err := os.Stat(prefixDir); !os.IsNotExist(err) {
		t.Fatalf("Prefix directory should not exist: %s", prefixDir)
	}

	t.Logf("Get() successfully downloaded and extracted files to: %s", dest)
}

func TestGetInvalidOS(t *testing.T) {
	version := "b7974"
	arch := "amd64"
	osVer := "cpm"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(arch, osVer, processor, version, dest)
	if err != ErrUnknownOS {
		t.Fatalf("Get() should have failed: %v", err)
	}
}

func TestGetInvalidProcessor(t *testing.T) {
	version := "b7974"
	arch := "amd64"
	osVer := "windows"
	processor := "flux"
	dest := t.TempDir()

	err := Get(arch, osVer, processor, version, dest)
	if err != ErrUnknownProcessor {
		t.Fatalf("Get() should have failed: %v", err)
	}
}

func TestGetInvalidVersion(t *testing.T) {
	version := "nogood"
	arch := "amd64"
	osVer := "linux"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(arch, osVer, processor, version, dest)
	if err != ErrInvalidVersion {
		t.Fatalf("Get() should have failed: %v", err)
	}
}

func TestGetDownloadLocationAndFilename_LinuxCPU(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Linux, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-x64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxCPU_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Linux, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-cpu-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxCUDA_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Linux, CUDA, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-cuda-13-x64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxCUDA_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Linux, CUDA, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-cuda-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxVulkan_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Linux, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-vulkan-x64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxVulkan_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Linux, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-vulkan-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_BookwormCPU_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Bookworm, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-cpu-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_BookwormCPU_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(AMD64, Bookworm, CPU, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Bookworm AMD64 CPU")
	}
}

func TestGetDownloadLocationAndFilename_BookwormCUDA_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Bookworm, CUDA, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-cuda-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_BookwormCUDA_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(AMD64, Bookworm, CUDA, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Bookworm AMD64 CUDA")
	}
}

func TestGetDownloadLocationAndFilename_BookwormVulkan_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Bookworm, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-vulkan-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_BookwormVulkan_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(AMD64, Bookworm, Vulkan, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Bookworm AMD64 Vulkan")
	}
}

func TestGetDownloadLocationAndFilename_TrixieCPU_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Trixie, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-trixie-cpu-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_TrixieCPU_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Trixie, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-x64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_TrixieCUDA_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(ARM64, Trixie, CUDA, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Trixie ARM64 CUDA")
	}
}

func TestGetDownloadLocationAndFilename_TrixieCUDA_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Trixie, CUDA, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-cuda-13-x64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_TrixieVulkan_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Trixie, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-trixie-vulkan-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_TrixieVulkan_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Trixie, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-vulkan-x64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_DarwinMetal_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Darwin, Metal, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-macos-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_DarwinMetal_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(AMD64, Darwin, Metal, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for macOS AMD64 Metal")
	}
}

func TestGetDownloadLocationAndFilename_DarwinCPU_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Darwin, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-macos-arm64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_DarwinCPU_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Darwin, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-macos-x64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsCPU_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Windows, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-win-cpu-x64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsCPU_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Windows, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-win-cpu-arm64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsCUDA_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(ARM64, Windows, CUDA, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Windows ARM64 CUDA")
	}
}

func TestGetDownloadLocationAndFilename_WindowsVulkan_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Windows, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-win-vulkan-x64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsVulkan_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(ARM64, Windows, Vulkan, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Windows ARM64 Vulkan")
	}
}

func TestGetDownloadLocationAndFilename_LinuxROCm_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Linux, ROCm, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-ubuntu-rocm-7.2-x64.tar.gz"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxROCm_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(ARM64, Linux, ROCm, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Linux ARM64 ROCm")
	}
}

func TestGetDownloadLocationAndFilename_WindowsROCm_AMD64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Windows, ROCm, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b7974"
	expectedFilename := "llama-b7974-bin-win-hip-radeon-x64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsROCm_ARM64(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(ARM64, Windows, ROCm, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Windows ARM64 ROCm")
	}
}

func TestGetDownloadLocationAndFilename_DarwinROCm(t *testing.T) {
	version := "b7974"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(AMD64, Darwin, ROCm, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Darwin ROCm")
	}
}

func TestGet404Error(t *testing.T) {
	version := "b7974"

	// Create a mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	arch := "amd64"
	osVer := "linux"
	processor := "cpu"
	dest := t.TempDir()

	// Override the get function to use our mock server
	originalGet := getFunc
	getFunc = func(ctx context.Context, asset Asset, dest string, progress getter.ProgressTracker) error {
		mockURL := server.URL + "/mock.tar.gz"
		return get(ctx, Asset{URL: mockURL}, dest, nil)
	}
	defer func() { getFunc = originalGet }()

	err := Get(arch, osVer, processor, version, dest)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}

	if !errors.Is(err, ErrFileNotFound) {
		t.Fatalf("expected ErrFileNotFound, got: %v", err)
	}

	t.Logf("Got expected error: %v", err)
}

func TestGetWithContext_FallbackToPreviousVersion(t *testing.T) {
	// Create a mock server to serve versions and mock files
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/llama-cpp-builder/version.json":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"tag_name": "b9000"}`))
		case "/llama-cpp-builder/previous.json":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"tag_name": "b8999"}`))
		case "/b9000/llama-b9000-bin-ubuntu-x64.tar.gz":
			// Simulate a 404 for the latest version
			w.WriteHeader(http.StatusNotFound)
		case "/b8999/llama-b8999-bin-ubuntu-x64.tar.gz":
			// Successfully return the previous version
			w.Header().Set("Content-Type", "application/gzip")
			w.Write(createMockTarGz(t, "b8999"))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Override URLs
	originalCurrent := currentVersionURL
	originalPrev := previousVersionURL
	currentVersionURL = server.URL + "/llama-cpp-builder/version.json"
	previousVersionURL = server.URL + "/llama-cpp-builder/previous.json"
	defer func() {
		currentVersionURL = originalCurrent
		previousVersionURL = originalPrev
	}()

	// Override getFunc to use our test server
	originalGet := getFunc
	getFunc = func(ctx context.Context, asset Asset, dest string, progress getter.ProgressTracker) error {
		url := asset.URL
		// Mock url rewriting to point to our test server
		parts := strings.Split(url, "/")
		filename := parts[len(parts)-1]
		versionPart := parts[len(parts)-2]

		mockURL := server.URL + "/" + versionPart + "/" + filename
		err := downloadAndExtractTarGz(Asset{URL: mockURL}, dest, nil)
		if err != nil && strings.Contains(err.Error(), "404") {
			return fmt.Errorf("%w: %s", ErrFileNotFound, url)
		}
		return err
	}
	defer func() { getFunc = originalGet }()

	dest := t.TempDir()

	// Call GetWithContext with auto version resolve
	err := GetWithContext(context.Background(), "amd64", "linux", "cpu", "", dest, nil)
	if err != nil {
		t.Fatalf("GetWithContext() failed: %v", err)
	}

	// Verify the downloaded file exist
	expectedFile := filepath.Join(dest, "libllama.so")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", expectedFile)
	}
}

// createMockVersionedTarGz creates a mock .tar.gz laid out like the llama.cpp
// release archives: a versioned library plus the unversioned symlink used by the
// loader, both under a top-level directory prefix.
func createMockVersionedTarGz(t *testing.T, version, soVersion string) []byte {
	t.Helper()

	var buf strings.Builder
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	prefix := "llama-" + version + "/"

	if err := tw.WriteHeader(&tar.Header{
		Name:     prefix,
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}

	versioned := "libllama." + soVersion + ".so"
	content := []byte("fake library " + version)
	if err := tw.WriteHeader(&tar.Header{
		Name: prefix + versioned,
		Mode: 0755,
		Size: int64(len(content)),
	}); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("failed to write tar content: %v", err)
	}

	if err := tw.WriteHeader(&tar.Header{
		Name:     prefix + "libllama.so",
		Mode:     0777,
		Typeflag: tar.TypeSymlink,
		Linkname: versioned,
	}); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}

	tw.Close()
	gzw.Close()

	return []byte(buf.String())
}

// TestExtractTarGzUpgradeRepointsSymlink verifies that installing over an older
// install repoints the unversioned symlink at the newly installed library
// instead of leaving it aimed at the superseded one.
func TestExtractTarGzUpgradeRepointsSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows, skipping test")
	}

	dest := t.TempDir()

	// Simulate an existing older install: a versioned library and the
	// unversioned symlink pointing at it.
	old := filepath.Join(dest, "libllama.0.0.10219.so")
	if err := os.WriteFile(old, []byte("old library"), 0755); err != nil {
		t.Fatalf("failed to create old library: %v", err)
	}
	link := filepath.Join(dest, "libllama.so")
	if err := os.Symlink("libllama.0.0.10219.so", link); err != nil {
		t.Fatalf("failed to create old symlink: %v", err)
	}

	mockTarGz := createMockVersionedTarGz(t, "b10375", "0.0.10375")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(mockTarGz)
	}))
	defer server.Close()

	if err := downloadAndExtractTarGz(Asset{URL: server.URL + "/b10375/llama-b10375-bin-ubuntu-x64.tar.gz"}, dest, nil); err != nil {
		t.Fatalf("downloadAndExtractTarGz() failed: %v", err)
	}

	target, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	if target != "libllama.0.0.10375.so" {
		t.Fatalf("symlink points at %q, want the newly installed libllama.0.0.10375.so", target)
	}

	// The superseded versioned file is intentionally left in place, and must not
	// have been overwritten through the stale symlink.
	content, err := os.ReadFile(old)
	if err != nil {
		t.Fatalf("failed to read superseded library: %v", err)
	}
	if string(content) != "old library" {
		t.Fatalf("superseded library was modified: %q", content)
	}
}

func TestVersionIsValid(t *testing.T) {
	valid := []string{"b7974", "v0.3.0", "v1.2.3-rc1"}
	for _, version := range valid {
		if err := VersionIsValid(version); err != nil {
			t.Errorf("VersionIsValid(%q) = %v, want nil", version, err)
		}
	}

	invalid := []string{"", "nogood", "b", "banana", "0.3.0", "latest"}
	for _, version := range invalid {
		if err := VersionIsValid(version); err == nil {
			t.Errorf("VersionIsValid(%q) accepted an invalid version", version)
		}
	}
}

func TestLlamaNightlyTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v0.3.0/nightly-tag.txt" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte("b10621\n"))
	}))
	defer server.Close()

	originalURL := nightlyTagURL
	nightlyTagURL = server.URL + "/%s/nightly-tag.txt"
	defer func() { nightlyTagURL = originalURL }()

	tag, err := LlamaNightlyTag("v0.3.0")
	if err != nil {
		t.Fatalf("LlamaNightlyTag() failed: %v", err)
	}
	if tag != "b10621" {
		t.Fatalf("LlamaNightlyTag() = %q, want b10621", tag)
	}

	// A nightly tag needs no lookup.
	if tag, err = LlamaNightlyTag("b7974"); err != nil || tag != "b7974" {
		t.Fatalf("LlamaNightlyTag() = %q, %v, want b7974, nil", tag, err)
	}

	if _, err = LlamaNightlyTag("v9.9.9"); err == nil {
		t.Fatal("LlamaNightlyTag() accepted a release with no nightly tag asset")
	}
}
