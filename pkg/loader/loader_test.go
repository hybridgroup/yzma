package loader

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGetLibraryPathEnvVar(t *testing.T) {
	result := getLibraryPathEnvVar()

	switch runtime.GOOS {
	case "linux", "freebsd":
		if result != "LD_LIBRARY_PATH" {
			t.Errorf("expected 'LD_LIBRARY_PATH', got '%s'", result)
		}
	case "darwin":
		if result != "DYLD_LIBRARY_PATH" {
			t.Errorf("expected 'DYLD_LIBRARY_PATH', got '%s'", result)
		}
	case "windows":
		if result != "PATH" {
			t.Errorf("expected 'PATH', got '%s'", result)
		}
	default:
		if result != "" {
			t.Errorf("expected empty string for unsupported OS, got '%s'", result)
		}
	}
}

func TestGetLibraryPathEnvVar_ReturnsNonEmpty(t *testing.T) {
	result := getLibraryPathEnvVar()

	// On supported platforms, should return a non-empty string
	switch runtime.GOOS {
	case "linux", "freebsd", "darwin", "windows":
		if result == "" {
			t.Errorf("expected non-empty string for %s", runtime.GOOS)
		}
	}
}

func TestEnsureLibraryPath_AddsNewPath(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	// Clear the env var
	os.Setenv(envVar, "")

	testPath := "/test/lib/path"
	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	if result != testPath {
		t.Errorf("expected '%s', got '%s'", testPath, result)
	}
}

func TestEnsureLibraryPath_PrependsToExistingPath(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	existingPath := "/existing/path"
	os.Setenv(envVar, existingPath)

	testPath := "/test/lib/path"
	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	expected := testPath + string(os.PathListSeparator) + existingPath
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestEnsureLibraryPath_DoesNotDuplicatePath(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	testPath := "/test/lib/path"
	os.Setenv(envVar, testPath)

	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	if result != testPath {
		t.Errorf("expected '%s', got '%s'", testPath, result)
	}

	// Verify it wasn't duplicated
	count := strings.Count(result, testPath)
	if count != 1 {
		t.Errorf("path was duplicated, count: %d", count)
	}
}

func TestEnsureLibraryPath_DoesNotDuplicateInMiddle(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	testPath := "/test/lib/path"
	separator := string(os.PathListSeparator)
	existingPaths := "/first/path" + separator + testPath + separator + "/last/path"
	os.Setenv(envVar, existingPaths)

	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	if result != existingPaths {
		t.Errorf("expected '%s', got '%s'", existingPaths, result)
	}
}

func TestEnsureLibraryPath_DoesNotDuplicateAtEnd(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	testPath := "/test/lib/path"
	separator := string(os.PathListSeparator)
	existingPaths := "/first/path" + separator + "/second/path" + separator + testPath
	os.Setenv(envVar, existingPaths)

	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	if result != existingPaths {
		t.Errorf("expected '%s', got '%s'", existingPaths, result)
	}
}

func TestEnsureLibraryPath_MultiplePaths(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	os.Setenv(envVar, "")

	paths := []string{"/path/one", "/path/two", "/path/three"}
	for _, p := range paths {
		err := ensureLibraryPath(p)
		if err != nil {
			t.Fatalf("ensureLibraryPath failed for '%s': %v", p, err)
		}
	}

	result := os.Getenv(envVar)
	separator := string(os.PathListSeparator)

	// All paths should be present
	for _, p := range paths {
		if !strings.Contains(result, p) {
			t.Errorf("path '%s' not found in '%s'", p, result)
		}
	}

	// Verify order (last added should be first)
	expected := "/path/three" + separator + "/path/two" + separator + "/path/one"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestEnsureLibraryPath_CalledTwiceWithSamePath(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	os.Setenv(envVar, "")

	testPath := "/test/lib/path"

	// Call twice with the same path
	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("first ensureLibraryPath failed: %v", err)
	}

	err = ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("second ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	if result != testPath {
		t.Errorf("expected '%s', got '%s'", testPath, result)
	}

	// Verify path appears only once
	separator := string(os.PathListSeparator)
	parts := strings.Split(result, separator)
	count := 0
	for _, p := range parts {
		if p == testPath {
			count++
		}
	}
	if count != 1 {
		t.Errorf("path should appear exactly once, appeared %d times", count)
	}
}

func TestEnsureLibraryPath_WithExistingMultiplePaths(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	separator := string(os.PathListSeparator)
	existingPaths := "/usr/lib" + separator + "/usr/local/lib" + separator + "/opt/lib"
	os.Setenv(envVar, existingPaths)

	testPath := "/new/lib/path"
	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	expected := testPath + separator + existingPaths
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestEnsureLibraryPath_PathWithSpaces(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	os.Setenv(envVar, "")

	testPath := "/path/with spaces/lib"
	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	if result != testPath {
		t.Errorf("expected '%s', got '%s'", testPath, result)
	}
}

func TestEnsureLibraryPath_PathWithSpecialChars(t *testing.T) {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		t.Skip("unsupported OS")
	}

	// Save original value
	original := os.Getenv(envVar)
	defer os.Setenv(envVar, original)

	os.Setenv(envVar, "")

	testPath := "/path/with-dashes_and_underscores/lib"
	err := ensureLibraryPath(testPath)
	if err != nil {
		t.Fatalf("ensureLibraryPath failed: %v", err)
	}

	result := os.Getenv(envVar)
	if result != testPath {
		t.Errorf("expected '%s', got '%s'", testPath, result)
	}
}
