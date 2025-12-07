package loader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jupiterrider/ffi"
)

// LoadLibrary The path can be an empty string to use the location as set by the YZMA_LIB env variable.
// The lib should be the "short name" for the library, for example:
// gguf, llama, mtmd
func LoadLibrary(path, lib string) (ffi.Lib, error) {
	if os.Getenv("YZMA_LIB") != "" {
		path = os.Getenv("YZMA_LIB")
	}

	// Ensure the library path is in LD_LIBRARY_PATH (Linux) or equivalent
	if path != "" {
		if err := ensureLibraryPath(path); err != nil {
			return ffi.Lib{}, err
		}
	}

	var filename string
	switch runtime.GOOS {
	case "linux", "freebsd":
		filename = filepath.Join(path, fmt.Sprintf("lib%s.so", lib))
	case "windows":
		filename = filepath.Join(path, fmt.Sprintf("%s.dll", lib))
	case "darwin":
		filename = filepath.Join(path, fmt.Sprintf("lib%s.dylib", lib))
	}

	return ffi.Load(filename)
}

// ensureLibraryPath ensures the given path is in the library search path
func ensureLibraryPath(path string) error {
	envVar := getLibraryPathEnvVar()
	if envVar == "" {
		return errors.New("unsupported OS for setting library path")
	}

	currentPath := os.Getenv(envVar)

	// Check if path is already in the library path
	if currentPath != "" {
		separator := string(os.PathListSeparator)
		paths := strings.Split(currentPath, separator)
		for _, p := range paths {
			if p == path {
				return nil // Already in path
			}
		}
		// Prepend the new path
		return os.Setenv(envVar, path+separator+currentPath)
	}

	// Set the path if not already set
	return os.Setenv(envVar, path)
}

// getLibraryPathEnvVar returns the library path environment variable for the current OS
func getLibraryPathEnvVar() string {
	switch runtime.GOOS {
	case "linux", "freebsd":
		return "LD_LIBRARY_PATH"
	case "darwin":
		return "DYLD_LIBRARY_PATH"
	case "windows":
		return "PATH"
	default:
		return ""
	}
}
