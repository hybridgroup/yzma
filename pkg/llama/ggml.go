package llama

import (
	"fmt"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

// Opaque types (represented as pointers)
type GGMLBackendBufferType uintptr

var (
	// GGML_API void ggml_backend_load_all(void);
	ggmlBackendLoadAllFunc ffi.Fun

	// GGML_API void ggml_backend_load_all(void);
	ggmlBackendLoadAllFromPath ffi.Fun

	// GGML_API ggml_backend_buffer_type_t ggml_backend_cpu_buffer_type(void);
	ggmlBackendCpuBufferType ffi.Fun
)

func loadGGML(lib ffi.Lib) error {
	var err error

	if ggmlBackendLoadAllFunc, err = lib.Prep("ggml_backend_load_all", &ffi.TypeVoid); err != nil {
		return loadError("ggml_backend_load_all", err)
	}

	if ggmlBackendLoadAllFromPath, err = lib.Prep("ggml_backend_load_all_from_path", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_load_all_from_path", err)
	}

	if ggmlBackendCpuBufferType, err = lib.Prep("ggml_backend_cpu_buffer_type", &ffi.TypeVoid); err != nil {
		return loadError("ggml_backend_cpu_buffer_type", err)
	}

	return nil
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

// GGMLBackendCpuBufferType returns the buffer type used for CPU backends.
func GGMLBackendCpuBufferType() GGMLBackendBufferType {
	var ret uintptr
	ggmlBackendCpuBufferType.Call(&ret)
	return GGMLBackendBufferType(ret)
}

const ffnExprsRegex = `\\.ffn_(up|down|gate)_(ch|)exps`

func ffnExprBlockRegex(index int) string {
	return fmt.Sprintf("blk\\.%d%s", index, ffnExprsRegex)
}

// TensorBuftBlockOverride creates a TensorBuftOverride for a specific block index to execute in the CPU.
func TensorBuftBlockOverride(index int) TensorBuftOverride {
	pattern := ffnExprBlockRegex(index)
	data, err := utils.BytePtrFromString(pattern)
	if err != nil {
		return TensorBuftOverride{}
	}
	return TensorBuftOverride{
		Pattern: data,
		Type:    GGMLBackendCpuBufferType(),
	}
}

// TensorBuftAllFFNExprsOverride creates a TensorBuftOverride for all FFN expression tensors to execute in the CPU.
func TensorBuftAllFFNExprsOverride() TensorBuftOverride {
	pattern := ffnExprsRegex
	data, err := utils.BytePtrFromString(pattern)
	if err != nil {
		return TensorBuftOverride{}
	}
	return TensorBuftOverride{
		Pattern: data,
		Type:    GGMLBackendCpuBufferType(),
	}
}
