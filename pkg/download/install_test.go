package download

import (
	"os"
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
