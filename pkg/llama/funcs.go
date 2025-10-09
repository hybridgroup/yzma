package llama

import (
	"github.com/jupiterrider/ffi"
)

var (
	backendInitFunc ffi.Fun

	backendFreeFunc ffi.Fun
)

func loadFuncs(lib ffi.Lib) error {
	var err error
	if backendInitFunc, err = lib.Prep("llama_backend_init", &ffi.TypeVoid); err != nil {
		return err
	}

	if backendFreeFunc, err = lib.Prep("llama_backend_free", &ffi.TypeVoid); err != nil {
		return err
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
