package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jupiterrider/ffi"
)

// LoadLibrary The path can be an empty string to use the location as set by the YZMA_LIB env variable.
// The lib should be the "short name" for the library, for example:
// gguf, llama, mtmd
func LoadLibrary(path, lib string) (ffi.Lib, error) {
	if os.Getenv("YZMA_LIB") != "" {
		path = os.Getenv("YZMA_LIB")
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
