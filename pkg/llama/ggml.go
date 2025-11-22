package llama

import (
	"errors"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

type (
	GGMLBackendDeviceType int32
	GGMLBackendDevice     uintptr
	GGMLBackendReg        uintptr
)

const (
	// CPU device using system memory
	GGMLBackendDeviceTypeCPU GGMLBackendDeviceType = iota
	// GPU device using dedicated memory
	GGMLBackendDeviceTypeGPU
	// integrated GPU device using host memory
	GGMLBackendDeviceTypeIGPU
	// accelerator devices intended to be used together with the CPU backend (e.g. BLAS or AMX)
	GGMLBackendDeviceTypeACCEL
)

var (
	// GGML_API void ggml_backend_load_all(void);
	ggmlBackendLoadAllFunc ffi.Fun

	// GGML_API void ggml_backend_load_all(void);
	ggmlBackendLoadAllFromPath ffi.Fun

	// Unload a backend if loaded dynamically and unregister it
	// GGML_API void               ggml_backend_unload(ggml_backend_reg_t reg);
	ggmlBackendUnloadFunc ffi.Fun

	// Device enumeration
	// GGML_API size_t             ggml_backend_dev_count(void);
	ggmlBackendDevCountFunc ffi.Fun

	// GGML_API ggml_backend_dev_t ggml_backend_dev_get(size_t index);
	ggmlBackendDevGetFunc ffi.Fun

	// GGML_API ggml_backend_dev_t ggml_backend_dev_by_name(const char * name);
	ggmlBackendDevByNameFunc ffi.Fun

	// GGML_API ggml_backend_dev_t ggml_backend_dev_by_type(enum ggml_backend_dev_type type);
	ggmlBackendDevByTypeFunc ffi.Fun

	// GGML_API size_t             ggml_backend_reg_count(void);
	ggmlBackendRegCountFunc ffi.Fun

	// GGML_API ggml_backend_reg_t ggml_backend_reg_get(size_t index);
	ggmlBackendRegGetFunc ffi.Fun

	// GGML_API ggml_backend_reg_t ggml_backend_reg_by_name(const char * name);
	ggmlBackendRegByNameFunc ffi.Fun
)

func loadGGML(lib ffi.Lib) error {
	var err error

	if ggmlBackendLoadAllFunc, err = lib.Prep("ggml_backend_load_all", &ffi.TypeVoid); err != nil {
		return loadError("ggml_backend_load_all", err)
	}

	if ggmlBackendLoadAllFromPath, err = lib.Prep("ggml_backend_load_all_from_path", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_load_all_from_path", err)
	}

	if ggmlBackendUnloadFunc, err = lib.Prep("ggml_backend_unload", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_unload", err)
	}

	if ggmlBackendDevCountFunc, err = lib.Prep("ggml_backend_dev_count", &ffi.TypeUint64); err != nil {
		return loadError("ggml_backend_dev_count", err)
	}

	if ggmlBackendDevGetFunc, err = lib.Prep("ggml_backend_dev_get", &ffi.TypePointer, &ffi.TypeUint64); err != nil {
		return loadError("ggml_backend_dev_get", err)
	}

	if ggmlBackendDevByNameFunc, err = lib.Prep("ggml_backend_dev_by_name", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_dev_by_name", err)
	}

	if ggmlBackendDevByTypeFunc, err = lib.Prep("ggml_backend_dev_by_type", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("ggml_backend_dev_by_type", err)
	}

	if ggmlBackendRegCountFunc, err = lib.Prep("ggml_backend_reg_count", &ffi.TypeUint64); err != nil {
		return loadError("ggml_backend_reg_count", err)
	}

	if ggmlBackendRegGetFunc, err = lib.Prep("ggml_backend_reg_get", &ffi.TypePointer, &ffi.TypeUint64); err != nil {
		return loadError("ggml_backend_reg_get", err)
	}

	if ggmlBackendRegByNameFunc, err = lib.Prep("ggml_backend_reg_by_name", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("ggml_backend_reg_by_name", err)
	}

	return nil
}

// GGMLBackendLoadAll loads all backends using the default search paths.
func GGMLBackendLoadAll() {
	ggmlBackendLoadAllFunc.Call(nil)
}

// GGMLBackendLoadAllFromPath loads all backends from a specific path.
func GGMLBackendLoadAllFromPath(path string) error {
	if path == "" {
		return errors.New("invalid path")
	}

	p := &[]byte(path + "\x00")[0]
	ggmlBackendLoadAllFromPath.Call(nil, unsafe.Pointer(&p))

	return nil
}

// GGMLBackendUnload unloads a backend if loaded dynamically and unregisters it.
func GGMLBackendUnload(reg GGMLBackendReg) {
	if reg == 0 {
		return
	}

	ggmlBackendUnloadFunc.Call(nil, unsafe.Pointer(&reg))
}

// GGMLBackendDeviceCount returns the number of backend devices.
func GGMLBackendDeviceCount() uint64 {
	var ret ffi.Arg
	ggmlBackendDevCountFunc.Call(unsafe.Pointer(&ret))
	return uint64(ret)
}

// GGMLBackendDeviceGet returns the backend device at the given index.
func GGMLBackendDeviceGet(index uint64) GGMLBackendDevice {
	var ret GGMLBackendDevice
	ggmlBackendDevGetFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&index))
	return ret
}

// GGMLBackendDeviceByName returns the backend device by its name.
func GGMLBackendDeviceByName(name string) GGMLBackendDevice {
	namePtr, _ := utils.BytePtrFromString(name)
	var ret GGMLBackendDevice
	ggmlBackendDevByNameFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&namePtr))
	return ret
}

// GGMLBackendDeviceByType returns the backend device by its type.
func GGMLBackendDeviceByType(devType GGMLBackendDeviceType) GGMLBackendDevice {
	var ret GGMLBackendDevice
	ggmlBackendDevByTypeFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&devType))
	return ret
}

// GGMLBackendRegCount returns the number of backend registrations.
func GGMLBackendRegCount() uint64 {
	var ret ffi.Arg
	ggmlBackendRegCountFunc.Call(unsafe.Pointer(&ret))
	return uint64(ret)
}

// GGMLBackendRegGet returns the backend registration at the given index.
func GGMLBackendRegGet(index uint64) GGMLBackendReg {
	var ret GGMLBackendReg
	ggmlBackendRegGetFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&index))
	return ret
}

// GGMLBackendRegByName returns the backend registration by its name.
func GGMLBackendRegByName(name string) GGMLBackendReg {
	namePtr, _ := utils.BytePtrFromString(name)
	var ret GGMLBackendReg
	ggmlBackendRegByNameFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&namePtr))
	return ret
}
