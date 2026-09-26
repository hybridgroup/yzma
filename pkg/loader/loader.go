package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jupiterrider/ffi"
)

// Lib is a handle to a shared library.
type Lib struct {
	lib ffi.Lib
}

// Prep gets the address of a function and describes its signature.
func (l Lib) Prep(name string, ret *ffi.Type, args ...*ffi.Type) (ffi.Fun, error) {
	return l.lib.Prep(name, ret, args...)
}

// LoadLibrary The path can be an empty string to use the location as set by the YZMA_LIB env variable.
// The lib should be the "short name" for the library, for example:
// gguf, llama, mtmd
func LoadLibrary(path, lib string) (Lib, error) {
	if path == "" && os.Getenv("YZMA_LIB") != "" {
		path = os.Getenv("YZMA_LIB")
	}

	// Ensure the library path is set
	if path == "" {
		return Lib{}, fmt.Errorf("library path not specified and YZMA_LIB env variable not set")
	}

	filename := GetLibraryFilename(path, lib)

	l, err := ffi.Load(filename)
	if err != nil {
		return Lib{}, err
	}

	return Lib{lib: l}, nil
}

// GetLibraryFilename returns the full path to the library file for the given path and library name.
// The library name should be the "short name" (e.g., "llama", "gguf", "mtmd").
// The function returns the appropriate filename based on the current OS:
//   - Linux/FreeBSD: lib<name>.so
//   - Windows: <name>.dll
//   - Darwin: lib<name>.dylib
func GetLibraryFilename(path, lib string) string {
	switch runtime.GOOS {
	case "linux", "freebsd":
		return filepath.Join(path, fmt.Sprintf("lib%s.so", lib))
	case "windows":
		return filepath.Join(path, fmt.Sprintf("%s.dll", lib))
	case "darwin":
		return filepath.Join(path, fmt.Sprintf("lib%s.dylib", lib))
	default:
		return filepath.Join(path, lib)
	}
}
