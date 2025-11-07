package llama

import (
	"fmt"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

// Opaque types (represented as pointers)
type GGMLBackendBufferType uintptr

var (
	// GGML_API ggml_backend_buffer_type_t ggml_backend_cpu_buffer_type(void);
	ggmlBackendCpuBufferType ffi.Fun
)

func loadGGMLBase(lib ffi.Lib) error {
	var err error

	if ggmlBackendCpuBufferType, err = lib.Prep("ggml_backend_cpu_buffer_type", &ffi.TypeVoid); err != nil {
		return loadError("ggml_backend_cpu_buffer_type", err)
	}

	return nil
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

// NewTensorBuftBlockOverride creates a TensorBuftOverride for a specific FFN block index to execute in the CPU.
func NewTensorBuftBlockOverride(index int) TensorBuftOverride {
	return NewTensorBuftOverride(ffnExprBlockRegex(index))
}

// NewTensorBuftAllFFNExprsOverride creates a TensorBuftOverride for all FFN expression tensors to execute in the CPU.
func NewTensorBuftAllFFNExprsOverride() TensorBuftOverride {
	return NewTensorBuftOverride(ffnExprsRegex)
}

// NewTensorBuftOverride creates a TensorBuftOverride for a custom pattern to execute in the CPU.
func NewTensorBuftOverride(pattern string) TensorBuftOverride {
	data, err := utils.BytePtrFromString(pattern)
	if err != nil {
		return TensorBuftOverride{}
	}
	return TensorBuftOverride{
		Pattern: data,
		Type:    GGMLBackendCpuBufferType(),
	}
}
