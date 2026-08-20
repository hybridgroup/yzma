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
	// GGML_API ggml_backend_buffer_type_t ggml_backend_cpu_buffer_type(void);
	ggmlBackendCpuBufferType ffi.Fun

	// GGML_API const char * ggml_backend_dev_name(ggml_backend_dev_t device);
	ggmlBackendDevNameFunc ffi.Fun

	// GGML_API void ggml_backend_dev_memory(ggml_backend_dev_t device, size_t * free, size_t * total);
	ggmlBackendDevMemoryFunc ffi.Fun

	// GGML_API const char * ggml_backend_dev_description(ggml_backend_dev_t device);
	ggmlBackendDevDescriptionFunc ffi.Fun

	// GGML_API enum ggml_backend_dev_type ggml_backend_dev_type(ggml_backend_dev_t device);
	ggmlBackendDevTypeFunc ffi.Fun

	// GGML_API ggml_backend_reg_t ggml_backend_dev_backend_reg(ggml_backend_dev_t device);
	ggmlBackendDevBackendRegFunc ffi.Fun

	// GGML_API const char * ggml_backend_reg_name(ggml_backend_reg_t reg);
	ggmlBackendRegNameFunc ffi.Fun

	// GGML_API int64_t ggml_blck_size(enum ggml_type type);
	ggmlBlckSizeFunc ffi.Fun

	// GGML_API size_t ggml_type_size(enum ggml_type type);
	ggmlTypeSizeFunc ffi.Fun

	// GGML_API size_t ggml_row_size(enum ggml_type type, int64_t ne);
	ggmlRowSizeFunc ffi.Fun

	// GGML_API const char * ggml_type_name(enum ggml_type type);
	ggmlTypeNameFunc ffi.Fun
)

func loadGGMLBase(lib ffi.Lib) error {
	var err error

	if ggmlBackendCpuBufferType, err = lib.Prep("ggml_backend_cpu_buffer_type", &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_cpu_buffer_type", err)
	}

	if ggmlBackendDevNameFunc, err = lib.Prep("ggml_backend_dev_name", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_dev_name", err)
	}

	if ggmlBackendDevMemoryFunc, err = lib.Prep("ggml_backend_dev_memory", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_dev_memory", err)
	}

	if ggmlBackendDevDescriptionFunc, err = lib.Prep("ggml_backend_dev_description", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_dev_description", err)
	}

	if ggmlBackendDevTypeFunc, err = lib.Prep("ggml_backend_dev_type", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_dev_type", err)
	}

	if ggmlBackendDevBackendRegFunc, err = lib.Prep("ggml_backend_dev_backend_reg", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_dev_backend_reg", err)
	}

	if ggmlBackendRegNameFunc, err = lib.Prep("ggml_backend_reg_name", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_reg_name", err)
	}

	if ggmlBlckSizeFunc, err = lib.Prep("ggml_blck_size", &ffi.TypeSint64, &ffi.TypeSint32); err != nil {
		return loadError("ggml_blck_size", err)
	}

	if ggmlTypeSizeFunc, err = lib.Prep("ggml_type_size", &ffi.TypeUint64, &ffi.TypeSint32); err != nil {
		return loadError("ggml_type_size", err)
	}

	if ggmlRowSizeFunc, err = lib.Prep("ggml_row_size", &ffi.TypeUint64, &ffi.TypeSint32, &ffi.TypeSint64); err != nil {
		return loadError("ggml_row_size", err)
	}

	if ggmlTypeNameFunc, err = lib.Prep("ggml_type_name", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("ggml_type_name", err)
	}

	return nil
}

// GGMLBackendCpuBufferType returns the buffer type used for CPU backends.
func GGMLBackendCpuBufferType() GGMLBackendBufferType {
	var ret uintptr
	ggmlBackendCpuBufferType.Call(unsafe.Pointer(&ret))
	return GGMLBackendBufferType(ret)
}

const ffnExprsRegex = `\.ffn_(up|down|gate)_(ch|)exps`

// MoEExpertTensorPattern is the canonical regex matching routed expert tensors.
// It matches ffn_(up|down|gate)_exps and ffn_(up|down|gate)_chexps tensor names.
const MoEExpertTensorPattern = ffnExprsRegex

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

// GGMLBackendDeviceName returns the name of the given backend device.
func GGMLBackendDeviceName(device GGMLBackendDevice) string {
	var ret *byte
	ggmlBackendDevNameFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&device))

	name := utils.BytePtrToString(ret)
	return name
}

// GGMLBackendDeviceDescription returns the description of the given backend device.
func GGMLBackendDeviceDescription(device GGMLBackendDevice) string {
	if device == 0 {
		return ""
	}

	var ret *byte
	ggmlBackendDevDescriptionFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&device))

	return utils.BytePtrToString(ret)
}

// GGMLBackendDevType returns the type of the given backend device.
// A device that is not valid gives GGMLBackendDeviceTypeCPU.
func GGMLBackendDevType(device GGMLBackendDevice) GGMLBackendDeviceType {
	if device == 0 {
		return GGMLBackendDeviceTypeCPU
	}

	// libffi always stores a full 8-byte ffi_arg for an integer return, so
	// the return buffer must be ffi.Arg-wide, not GGMLBackendDeviceType-wide (int32).
	var ret ffi.Arg
	ggmlBackendDevTypeFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&device))

	return GGMLBackendDeviceType(int32(ret))
}

// GGMLBackendDeviceBackendReg returns the backend registration that owns the given device.
func GGMLBackendDeviceBackendReg(device GGMLBackendDevice) GGMLBackendReg {
	if device == 0 {
		return 0
	}

	var ret GGMLBackendReg
	ggmlBackendDevBackendRegFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&device))

	return ret
}

// GGMLBackendRegName returns the name of the given backend registration.
func GGMLBackendRegName(reg GGMLBackendReg) string {
	if reg == 0 {
		return ""
	}

	var ret *byte
	ggmlBackendRegNameFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&reg))

	return utils.BytePtrToString(ret)
}

// GGMLBlockSize returns the number of elements in one block of the given type.
// The type must be one of the GGMLType constants; any other value reads past the
// end of the type table in C.
func GGMLBlockSize(t GGMLType) int64 {
	// libffi always stores a full 8-byte ffi_arg for an integer return, so the
	// return buffer must be ffi.Arg-wide.
	var ret ffi.Arg
	ggmlBlckSizeFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&t))

	return int64(ret)
}

// GGMLTypeSize returns the size in bytes of one block of the given type.
// The type must be one of the GGMLType constants; any other value reads past the
// end of the type table in C.
func GGMLTypeSize(t GGMLType) uint64 {
	var ret ffi.Arg
	ggmlTypeSizeFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&t))

	return uint64(ret)
}

// GGMLRowSize returns the size in bytes of a row of ne elements of the given type.
// The type must be one of the GGMLType constants; any other value reads past the
// end of the type table in C.
func GGMLRowSize(t GGMLType, ne int64) uint64 {
	var ret ffi.Arg
	ggmlRowSizeFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&t), unsafe.Pointer(&ne))

	return uint64(ret)
}

// GGMLTypeName returns the name of the given type, for example "q4_K".
// The type must be one of the GGMLType constants; any other value reads past the
// end of the type table in C.
func GGMLTypeName(t GGMLType) string {
	var ret *byte
	ggmlTypeNameFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&t))

	return utils.BytePtrToString(ret)
}
