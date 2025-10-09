package llama

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	backendInitFunc ffi.Fun

	backendFreeFunc ffi.Fun

	// GGML_API void ggml_backend_load_all(void);
	ggmlBackendLoadAllFunc ffi.Fun

	// GGML_API void ggml_backend_load_all(void);
	ggmlBackendLoadAllFromPath ffi.Fun
)

func loadFuncs(lib ffi.Lib) error {
	var err error
	if backendInitFunc, err = lib.Prep("llama_backend_init", &ffi.TypeVoid); err != nil {
		return fmt.Errorf("llama_backend_init: %w", err)
	}

	if backendFreeFunc, err = lib.Prep("llama_backend_free", &ffi.TypeVoid); err != nil {
		return fmt.Errorf("llama_backend_free: %w", err)
	}

	if runtime.GOOS == "windows" {
		path := os.Getenv("YZMA_LIB")
		filename := filepath.Join(path, "ggml.dll")
		lib, err = ffi.Load(filename)
		if err != nil {
			return fmt.Errorf("load ggml.dll: %w", err)
		}
	}

	if ggmlBackendLoadAllFunc, err = lib.Prep("ggml_backend_load_all", &ffi.TypeVoid); err != nil {
		return fmt.Errorf("ggml_backend_load_all: %w", err)
	}

	if ggmlBackendLoadAllFromPath, err = lib.Prep("ggml_backend_load_all_from_path", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ggml_backend_load_all_from_path: %w", err)
	}

	return nil
}

// BackendInit initializes the llama.cpp back-end.
func BackendInit() {
	backendInitFunc.Call(nil)
}

// BackendFree frees the llama.cpp back-end.
func BackendFree() {
	backendFreeFunc.Call(nil)
}

// GGMLBackendLoadAll loads all backends using the default search paths.
func GGMLBackendLoadAll() {
	ggmlBackendLoadAllFunc.Call(nil)
}

// GGMLBackendLoadAllFromPath loads all backends from a specific path.
func GGMLBackendLoadAllFromPath(path string) {
	p := &[]byte(path + "\x00")[0]
	ggmlBackendLoadAllFromPath.Call(nil, unsafe.Pointer(&p))
}
