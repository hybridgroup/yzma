package loader

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/jupiterrider/ffi"
)

func LoadLibrary(path string) (ffi.Lib, error) {
	if os.Getenv("YZMA_LIB") != "" {
		path = os.Getenv("YZMA_LIB")
	}

	var filename string
	switch runtime.GOOS {
	case "linux", "freebsd":
		filename = filepath.Join(path, "libmtmd.so")
	case "windows":
		filename = filepath.Join(path, "llama.dll")
	case "darwin":
		filename = filepath.Join(path, "libmtmd.dylib")
	}

	return ffi.Load(filename)
}
