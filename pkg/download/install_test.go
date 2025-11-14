package download

import (
	"runtime"
	"testing"
)

func TestInstall(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("skipping test since github API sends 403 error")
	}

	dest := t.TempDir()

	exists := alreadyInstalled(dest)
	if exists {
		t.Fatalf("should NOT see libraries are installed")
	}

	if err := initialInstall(dest); err != nil {
		t.Fatalf("should be able to install libraries: %v", err)
	}

	exists = alreadyInstalled(dest)
	if !exists {
		t.Fatalf("should see libraries are installed")
	}

	if err := createVersionFile(dest, "1.0.0"); err != nil {
		t.Fatalf("should be able update the version file is a different version: %v", err)
	}

	isLatest, version1, err := alreadyLatestVersion(dest)
	if err != nil {
		t.Fatalf("should be able to check if the version is latest: %v", err)
	}

	if isLatest {
		t.Fatalf("should NOT see that the version is latest")
	}

	if err := upgradeInstall(dest, version1); err != nil {
		t.Fatalf("should be able to upgrade the libraries: %v", err)
	}

	isLatest, version2, err := alreadyLatestVersion(dest)
	if err != nil {
		t.Fatalf("should be able to check if the version is latest: %v", err)
	}

	if !isLatest {
		t.Fatalf("should see that the version is latest")
	}

	if version1 != version2 {
		t.Fatalf("should see that the version is updated")
	}
}
