package download

import "testing"

func TestVersionFile(t *testing.T) {
	dest := t.TempDir()

	exists := doesVersionFileExist(dest)
	if exists {
		t.Fatalf("should NOT see the version file existing")
	}

	tag, err := downloadVersionFile(llamaCppVersionDocURL)
	if err != nil {
		t.Fatalf("should be able to download version file: %v", err)
	}

	if err := createVersionFile(dest, tag); err != nil {
		t.Fatalf("should be able to create version file: %v", err)
	}

	exists = doesVersionFileExist(dest)
	if !exists {
		t.Fatalf("should see the version file exists")
	}

	exists, version, err := isLatestVersionInstalled(dest)
	if err != nil {
		t.Fatalf("should be able to check latest version with no errors: %v", err)
	}

	if !exists {
		t.Fatalf("should see the latest version is installed")
	}

	if version != tag.TagName {
		t.Fatalf("should see the latest version is what is expected: %s != %s", version, tag.TagName)
	}
}
