package download

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	getter "github.com/hashicorp/go-getter"
)

func TestLlamaLatestVersion(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/repos/ggml-org/llama.cpp/releases/latest" {
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

	// Override the API URL for testing
	originalURL := apiURL
	apiURL = server.URL + "/repos/ggml-org/llama.cpp/releases/latest"
	defer func() { apiURL = originalURL }()

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

	// Override the API URL for testing
	originalURL := apiURL
	apiURL = server.URL + "/repos/ggml-org/llama.cpp/releases/latest"
	defer func() { apiURL = originalURL }()

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
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		// Replace the real URL with our mock server URL
		mockURL := server.URL + "/b7974/llama-b7974-bin-ubuntu-x64.tar.gz"
		return downloadAndExtractTarGz(mockURL, dest, nil)
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
	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		mockURL := server.URL + "/mock.tar.gz"
		return get(ctx, mockURL, dest, nil)
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
