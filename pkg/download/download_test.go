package download

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLlamaLatestVersion(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("skipping test since github API sends 403 error")
	}

	version, err := LlamaLatestVersion()
	if err != nil {
		t.Fatal("count not get latest version", err)
	}

	if !strings.HasPrefix(version, "b") {
		t.Fatalf("Expected version should start with 'b', got '%s'", version)
	}

	t.Logf("LlamaLatestVersion returned: %s", version)
}

func TestGetLinuxCPU(t *testing.T) {
	version := "b6795"
	arch := "amd64"
	osVer := "linux"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(arch, osVer, processor, version, dest)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	expectedFile := filepath.Join(dest, "libllama.so")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", expectedFile)
	}

	t.Logf("Get() successfully downloaded the file to: %s", expectedFile)
}

func TestGetInvalidOS(t *testing.T) {
	version := "b6795"
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
	version := "b6795"
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
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Linux, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-ubuntu-x64.zip//build/bin"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxCPU_ARM64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(ARM64, Linux, CPU, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Linux ARM64 CPU")
	}
}

func TestGetDownloadLocationAndFilename_LinuxCUDA_AMD64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Linux, CUDA, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-ubuntu-cuda-x64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxCUDA_ARM64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Linux, CUDA, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-ubuntu-cuda-arm64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxVulkan_AMD64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Linux, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-ubuntu-vulkan-x64.zip//build/bin"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_LinuxVulkan_ARM64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Linux, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-ubuntu-vulkan-arm64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_DarwinMetal_ARM64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Darwin, Metal, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-macos-arm64.zip//build/bin"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_DarwinMetal_AMD64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(AMD64, Darwin, Metal, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for macOS AMD64 Metal")
	}
}

func TestGetDownloadLocationAndFilename_DarwinCPU_ARM64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Darwin, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-macos-arm64.zip//build/bin"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_DarwinCPU_AMD64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Darwin, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-macos-x64.zip//build/bin"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsCPU_AMD64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Windows, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-win-cpu-x64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsCPU_ARM64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(ARM64, Windows, CPU, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-win-cpu-arm64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsCUDA_ARM64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(ARM64, Windows, CUDA, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Windows ARM64 CUDA")
	}
}

func TestGetDownloadLocationAndFilename_WindowsVulkan_AMD64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	location, filename, err := getDownloadLocationAndFilename(AMD64, Windows, Vulkan, version, dest)
	if err != nil {
		t.Fatalf("getDownloadLocationAndFilename() failed: %v", err)
	}

	expectedLocation := "https://github.com/ggml-org/llama.cpp/releases/download/b6795"
	expectedFilename := "llama-b6795-bin-win-vulkan-x64.zip"

	if location != expectedLocation {
		t.Errorf("location = %q, want %q", location, expectedLocation)
	}
	if filename != expectedFilename {
		t.Errorf("filename = %q, want %q", filename, expectedFilename)
	}
}

func TestGetDownloadLocationAndFilename_WindowsVulkan_ARM64(t *testing.T) {
	version := "b6795"
	dest := t.TempDir()

	_, _, err := getDownloadLocationAndFilename(ARM64, Windows, Vulkan, version, dest)
	if err == nil {
		t.Fatal("getDownloadLocationAndFilename() should have failed for Windows ARM64 Vulkan")
	}
}
