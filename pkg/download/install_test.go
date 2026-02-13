package download

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	getter "github.com/hashicorp/go-getter"
)

func TestInstall(t *testing.T) {
	dest := t.TempDir()
	version := "b7974"

	// Create mock tar.gz content
	mockTarGz := createMockTarGz(t, version)

	// Create mock servers
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/ggml-org/llama.cpp/releases/latest" {
			t.Errorf("unexpected API path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response := struct {
			TagName string `json:"tag_name"`
		}{
			TagName: version,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer apiServer.Close()

	downloadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(mockTarGz)
	}))
	defer downloadServer.Close()

	// Override URLs and functions for testing
	originalAPIURL := apiURL
	originalLlamaCppVersionDocURL := llamaCppVersionDocURL
	originalGetFunc := getFunc

	apiURL = apiServer.URL + "/repos/ggml-org/llama.cpp/releases/latest"
	llamaCppVersionDocURL = apiServer.URL + "/repos/ggml-org/llama.cpp/releases/latest"

	getFunc = func(ctx context.Context, url string, dest string, progress getter.ProgressTracker) error {
		// Use mock download server
		mockURL := downloadServer.URL + "/mock.tar.gz"
		return downloadAndExtractTarGz(mockURL, dest, nil)
	}

	defer func() {
		apiURL = originalAPIURL
		llamaCppVersionDocURL = originalLlamaCppVersionDocURL
		getFunc = originalGetFunc
	}()

	// Test: Check not installed initially
	exists := alreadyInstalled(dest)
	if exists {
		t.Fatalf("should NOT see libraries are installed")
	}

	// Test: Initial install
	if err := initialInstall(dest, CPU); err != nil {
		t.Fatalf("should be able to install libraries: %v", err)
	}

	exists = alreadyInstalled(dest)
	if !exists {
		t.Fatalf("should see libraries are installed")
	}

	// Verify version file was created
	versionFilePath := filepath.Join(dest, versionFile)
	if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
		t.Fatalf("version file should exist")
	}

	// Test: Create old version file
	if err := createVersionFile(dest, "b1000"); err != nil {
		t.Fatalf("should be able to update the version file to a different version: %v", err)
	}

	// Test: Check version is not latest
	isLatest, version1, err := alreadyLatestVersion(dest)
	if err != nil {
		t.Fatalf("should be able to check if the version is latest: %v", err)
	}

	if isLatest {
		t.Fatalf("should NOT see that the version is latest")
	}

	if version1 != version {
		t.Fatalf("expected version %s, got %s", version, version1)
	}

	// Test: Upgrade install
	if err := upgradeInstall(dest, CPU, version1); err != nil {
		t.Fatalf("should be able to upgrade the libraries: %v", err)
	}

	// Test: Check version is now latest
	isLatest, version2, err := alreadyLatestVersion(dest)
	if err != nil {
		t.Fatalf("should be able to check if the version is latest: %v", err)
	}

	if !isLatest {
		t.Fatalf("should see that the version is latest")
	}

	if version1 != version2 {
		t.Fatalf("expected versions to match: %s != %s", version1, version2)
	}
}

func TestAlreadyInstalled(t *testing.T) {
	dest := t.TempDir()

	// Test: Not installed
	if alreadyInstalled(dest) {
		t.Fatal("should not be installed in empty directory")
	}

	// Test: Create version file
	versionFilePath := filepath.Join(dest, versionFile)
	if err := os.WriteFile(versionFilePath, []byte(`{"tag_name":"b7974"}`), 0644); err != nil {
		t.Fatalf("failed to create version file: %v", err)
	}

	// Test: Now installed
	if !alreadyInstalled(dest) {
		t.Fatal("should be installed after version file created")
	}
}

func TestAlreadyLatestVersion(t *testing.T) {
	dest := t.TempDir()
	version := "b7974"

	// Create mock API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			TagName string `json:"tag_name"`
		}{
			TagName: version,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	originalAPIURL := apiURL
	apiURL = server.URL + "/repos/ggml-org/llama.cpp/releases/latest"
	defer func() { apiURL = originalAPIURL }()

	// Test: Same version
	if err := createVersionFile(dest, version); err != nil {
		t.Fatalf("failed to create version file: %v", err)
	}

	isLatest, returnedVersion, err := alreadyLatestVersion(dest)
	if err != nil {
		t.Fatalf("failed to check version: %v", err)
	}

	if !isLatest {
		t.Fatal("should be latest version")
	}

	if returnedVersion != version {
		t.Fatalf("expected version %s, got %s", version, returnedVersion)
	}

	// Test: Old version
	if err := createVersionFile(dest, "b1000"); err != nil {
		t.Fatalf("failed to create version file: %v", err)
	}

	isLatest, returnedVersion, err = alreadyLatestVersion(dest)
	if err != nil {
		t.Fatalf("failed to check version: %v", err)
	}

	if isLatest {
		t.Fatal("should not be latest version")
	}

	if returnedVersion != version {
		t.Fatalf("expected version %s, got %s", version, returnedVersion)
	}
}

func TestCreateVersionFile(t *testing.T) {
	dest := t.TempDir()
	version := "b7974"

	if err := createVersionFile(dest, version); err != nil {
		t.Fatalf("failed to create version file: %v", err)
	}

	versionFilePath := filepath.Join(dest, versionFile)
	data, err := os.ReadFile(versionFilePath)
	if err != nil {
		t.Fatalf("failed to read version file: %v", err)
	}

	var tag tag
	if err := json.Unmarshal(data, &tag); err != nil {
		t.Fatalf("failed to unmarshal version file: %v", err)
	}

	if tag.TagName != version {
		t.Fatalf("expected version %s, got %s", version, tag.TagName)
	}
}

func TestDownloadVersionFile(t *testing.T) {
	version := "b7974"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			TagName string `json:"tag_name"`
		}{
			TagName: version,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	returnedVersion, err := downloadVersionFile(server.URL)
	if err != nil {
		t.Fatalf("failed to download version file: %v", err)
	}

	if returnedVersion != version {
		t.Fatalf("expected version %s, got %s", version, returnedVersion)
	}
}

func TestDownloadVersionFile_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := downloadVersionFile(server.URL)
	if err == nil {
		t.Fatal("expected error for failed request")
	}
}

func TestInstallLibraries_NoUpgrade(t *testing.T) {
	dest := t.TempDir()
	version := "b7974"

	// Create version file to simulate already installed
	if err := createVersionFile(dest, version); err != nil {
		t.Fatalf("failed to create version file: %v", err)
	}

	// Run install with allowUpgrade=false
	err := InstallLibraries(dest, CPU, false)
	if err != nil {
		t.Fatalf("InstallLibraries should succeed when already installed and upgrade not allowed: %v", err)
	}
}
