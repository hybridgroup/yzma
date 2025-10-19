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
	osVer := "linux"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(osVer, processor, version, dest)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	expectedFile := filepath.Join(dest, "llama-b6795-bin-ubuntu-x64.zip")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Downloaded file not found: %s", expectedFile)
	}

	t.Logf("Get() successfully downloaded the file to: %s", expectedFile)
}

func TestGetInvalidOS(t *testing.T) {
	version := "b6795"
	osVer := "cpm"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(osVer, processor, version, dest)
	if err != errUnknownOS {
		t.Fatalf("Get() should have failed: %v", err)
	}
}

func TestGetInvalidProcessor(t *testing.T) {
	version := "b6795"
	osVer := "windows"
	processor := "flux"
	dest := t.TempDir()

	err := Get(osVer, processor, version, dest)
	if err != errUnknownProcessor {
		t.Fatalf("Get() should have failed: %v", err)
	}
}

func TestGetInvalidVersion(t *testing.T) {
	version := "nogood"
	osVer := "linux"
	processor := "cpu"
	dest := t.TempDir()

	err := Get(osVer, processor, version, dest)
	if err != errInvalidVersion {
		t.Fatalf("Get() should have failed: %v", err)
	}
}
