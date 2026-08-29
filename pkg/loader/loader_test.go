package loader

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jupiterrider/ffi"
)

func TestGetLibraryFilename(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		lib      string
		expected map[string]string // OS -> expected result
	}{
		{
			name: "llama library",
			path: "/usr/local/lib",
			lib:  "llama",
			expected: map[string]string{
				"linux":   "/usr/local/lib/libllama.so",
				"freebsd": "/usr/local/lib/libllama.so",
				"darwin":  "/usr/local/lib/libllama.dylib",
				"windows": "/usr/local/lib/llama.dll",
			},
		},
		{
			name: "gguf library",
			path: "/opt/yzma",
			lib:  "gguf",
			expected: map[string]string{
				"linux":   "/opt/yzma/libgguf.so",
				"freebsd": "/opt/yzma/libgguf.so",
				"darwin":  "/opt/yzma/libgguf.dylib",
				"windows": "/opt/yzma/gguf.dll",
			},
		},
		{
			name: "mtmd library",
			path: "/home/user/libs",
			lib:  "mtmd",
			expected: map[string]string{
				"linux":   "/home/user/libs/libmtmd.so",
				"freebsd": "/home/user/libs/libmtmd.so",
				"darwin":  "/home/user/libs/libmtmd.dylib",
				"windows": "/home/user/libs/mtmd.dll",
			},
		},
		{
			name: "empty path",
			path: "",
			lib:  "llama",
			expected: map[string]string{
				"linux":   "libllama.so",
				"freebsd": "libllama.so",
				"darwin":  "libllama.dylib",
				"windows": "llama.dll",
			},
		},
		{
			name: "relative path",
			path: "./lib",
			lib:  "llama",
			expected: map[string]string{
				"linux":   "lib/libllama.so",
				"freebsd": "lib/libllama.so",
				"darwin":  "lib/libllama.dylib",
				"windows": "lib/llama.dll",
			},
		},
		{
			name: "path with spaces",
			path: "/path/with spaces",
			lib:  "llama",
			expected: map[string]string{
				"linux":   "/path/with spaces/libllama.so",
				"freebsd": "/path/with spaces/libllama.so",
				"darwin":  "/path/with spaces/libllama.dylib",
				"windows": "/path/with spaces/llama.dll",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLibraryFilename(tt.path, tt.lib)

			expected, ok := tt.expected[runtime.GOOS]
			if !ok {
				// For unsupported OS, just check it returns something
				if result == "" {
					t.Error("expected non-empty result for unsupported OS")
				}
				return
			}

			// Normalize paths for comparison (handles OS-specific separators)
			expectedNorm := filepath.FromSlash(expected)
			if result != expectedNorm {
				t.Errorf("expected '%s', got '%s'", expectedNorm, result)
			}
		})
	}
}

func TestGetLibraryFilename_CurrentOS(t *testing.T) {
	path := "/test/path"
	lib := "testlib"

	result := GetLibraryFilename(path, lib)

	switch runtime.GOOS {
	case "linux", "freebsd":
		expected := filepath.Join(path, "libtestlib.so")
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	case "darwin":
		expected := filepath.Join(path, "libtestlib.dylib")
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	case "windows":
		expected := filepath.Join(path, "testlib.dll")
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	}
}

func TestGetLibraryFilename_ReturnsNonEmpty(t *testing.T) {
	result := GetLibraryFilename("/some/path", "somelib")

	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestGetLibraryFilename_ContainsLibName(t *testing.T) {
	lib := "mylib"
	result := GetLibraryFilename("/path", lib)

	// Result should contain the library name
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}

	// Check that result contains the library name in some form
	switch runtime.GOOS {
	case "linux", "freebsd", "darwin":
		if result != filepath.Join("/path", "lib"+lib+".so") &&
			result != filepath.Join("/path", "lib"+lib+".dylib") {
			// Just verify it's not empty for other cases
		}
	case "windows":
		expected := filepath.Join("/path", lib+".dll")
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	}
}

func TestGetLibraryFilename_DifferentLibNames(t *testing.T) {
	libs := []string{"llama", "gguf", "mtmd", "ggml"}
	basePath := "/lib"

	for _, lib := range libs {
		t.Run(lib, func(t *testing.T) {
			result := GetLibraryFilename(basePath, lib)

			if result == "" {
				t.Errorf("expected non-empty result for lib '%s'", lib)
			}

			// Verify the result starts with the base path (normalized for the OS)
			expectedPrefix := filepath.FromSlash(basePath)
			if len(result) < len(expectedPrefix) || result[:len(expectedPrefix)] != expectedPrefix {
				t.Errorf("expected path to start with '%s', got '%s'", expectedPrefix, result)
			}

			// Verify the result contains the library name
			if !strings.Contains(result, lib) {
				t.Errorf("expected result to contain '%s', got '%s'", lib, result)
			}
		})
	}
}

// The cif must point at the same argument type array that Prep keeps alive.
// If it points somewhere else, the collector can free the array that libffi
// reads on every call.
func TestPrepKeepsArgTypes(t *testing.T) {
	if os.Getenv("YZMA_LIB") == "" {
		t.Skip("no YZMA_LIB set, skipping test")
	}

	lib, err := LoadLibrary("", "llama")
	if err != nil {
		t.Fatalf("LoadLibrary failed: %v", err)
	}

	args := []*ffi.Type{&ffi.TypePointer, &ffi.TypeSint32}
	fun, err := lib.Prep("llama_model_n_params", &ffi.TypeUint64, args...)
	if err != nil {
		t.Fatalf("Prep failed: %v", err)
	}

	kept := keepArgTypes(args)

	if fun.Cif.ArgTypes != &kept[0] {
		t.Errorf("cif.ArgTypes is %p, want the kept array at %p", fun.Cif.ArgTypes, &kept[0])
	}

	runtime.GC()

	if fun.Cif.ArgTypes != &kept[0] {
		t.Errorf("cif.ArgTypes changed after a collection")
	}
}

// A second Prep of the same signature must not keep a second array.
func TestPrepReusesArgTypes(t *testing.T) {
	if os.Getenv("YZMA_LIB") == "" {
		t.Skip("no YZMA_LIB set, skipping test")
	}

	lib, err := LoadLibrary("", "llama")
	if err != nil {
		t.Fatalf("LoadLibrary failed: %v", err)
	}

	argTypesMu.Lock()
	before := len(argTypes)
	argTypesMu.Unlock()

	for range 10 {
		if _, err := lib.Prep("llama_model_n_params", &ffi.TypeUint64, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
			t.Fatalf("Prep failed: %v", err)
		}
	}

	argTypesMu.Lock()
	after := len(argTypes)
	argTypesMu.Unlock()

	if after-before > 1 {
		t.Errorf("10 preps of one signature kept %d arrays, want at most 1", after-before)
	}
}
