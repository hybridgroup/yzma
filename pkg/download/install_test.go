package download

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAlreadyInstalled(t *testing.T) {
	dest := t.TempDir()

	// Should not be installed in a new temp dir
	if AlreadyInstalled(dest) {
		t.Fatal("AlreadyInstalled should return false for empty directory")
	}

	// Create a dummy library file to simulate installation
	libFile := filepath.Join(dest, LibraryName(runtime.GOOS))
	if err := os.WriteFile(libFile, []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to create dummy library file: %v", err)
	}

	if !AlreadyInstalled(dest) {
		t.Fatal("AlreadyInstalled should return true when library file exists")
	}
}

func TestHasCUDA_ParseVersion(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("CUDA is not available on macOS, skipping test")
	}

	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Use a shell command to echo fake output
	execCommand = func(name string, arg ...string) *exec.Cmd {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", "echo CUDA Version: 13.0")
		}
		return exec.Command("sh", "-c", "echo CUDA Version: 13.0")
	}

	ok, version := HasCUDA()
	if !ok {
		t.Fatal("expected CUDA to return true")
	}
	if version != "13.0" {
		t.Fatalf("expected version '13.0', got '%s'", version)
	}
}

func TestCUDA_NoCUDA(t *testing.T) {
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Simulate nvidia-smi not found or error
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false") // always fails
	}

	ok, version := HasCUDA()
	if ok {
		t.Fatal("expected CUDA to return false")
	}
	if version != "" {
		t.Fatalf("expected empty version, got '%s'", version)
	}
}

func TestHasROCm_ParseVersion(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("ROCm is not available on non-Linux, skipping test")
	}

	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Use a shell command to echo fake rocminfo output
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("sh", "-c", "echo 'Runtime Version:         1.1'")
	}

	ok, version := HasROCm()
	if !ok {
		t.Fatal("expected ROCm to return true")
	}
	if version != "1.1" {
		t.Fatalf("expected version '1.1', got '%s'", version)
	}
}

func TestHasROCm_NoROCm(t *testing.T) {
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Simulate rocminfo not found or error
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("false") // always fails
	}

	ok, version := HasROCm()
	if ok {
		t.Fatal("expected ROCm to return false")
	}
	if version != "" {
		t.Fatalf("expected empty version, got '%s'", version)
	}
}
