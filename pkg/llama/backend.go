package llama

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	backendInitFunc ffi.Fun

	backendFreeFunc ffi.Fun

	// LLAMA_API size_t llama_max_devices(void);
	maxDevicesFunc ffi.Fun

	// LLAMA_API size_t llama_max_parallel_sequences(void);
	maxParallelSequencesFunc ffi.Fun

	// LLAMA_API bool llama_supports_mmap(void);
	supportsMmapFunc ffi.Fun

	// LLAMA_API bool llama_supports_mlock(void);
	supportsMlockFunc ffi.Fun

	// LLAMA_API bool llama_supports_gpu_offload(void);
	supportsGpuOffloadFunc ffi.Fun

	// LLAMA_API bool llama_supports_rpc(void);
	supportsRpcFunc ffi.Fun

	// LLAMA_API int64_t llama_time_us(void);
	timeUsFunc ffi.Fun
)

func loadBackendFuncs(lib ffi.Lib) error {
	var err error
	if backendInitFunc, err = lib.Prep("llama_backend_init", &ffi.TypeVoid); err != nil {
		return loadError("llama_backend_init", err)
	}

	if backendFreeFunc, err = lib.Prep("llama_backend_free", &ffi.TypeVoid); err != nil {
		return loadError("llama_backend_free", err)
	}

	if maxDevicesFunc, err = lib.Prep("llama_max_devices", &ffi.TypeUint32); err != nil {
		return loadError("llama_max_devices", err)
	}

	if maxParallelSequencesFunc, err = lib.Prep("llama_max_parallel_sequences", &ffi.TypeUint32); err != nil {
		return loadError("llama_max_parallel_sequences", err)
	}

	if supportsMmapFunc, err = lib.Prep("llama_supports_mmap", &ffi.TypeUint8); err != nil {
		return loadError("llama_supports_mmap", err)
	}

	if supportsMlockFunc, err = lib.Prep("llama_supports_mlock", &ffi.TypeUint8); err != nil {
		return loadError("llama_supports_mlock", err)
	}

	if supportsGpuOffloadFunc, err = lib.Prep("llama_supports_gpu_offload", &ffi.TypeUint8); err != nil {
		return loadError("llama_supports_gpu_offload", err)
	}

	if supportsRpcFunc, err = lib.Prep("llama_supports_rpc", &ffi.TypeUint8); err != nil {
		return loadError("llama_supports_rpc", err)
	}

	if timeUsFunc, err = lib.Prep("llama_time_us", &ffi.TypeSint64); err != nil {
		return loadError("llama_time_us", err)
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

// MaxDevices returns the maximum number of devices supported.
func MaxDevices() uint64 {
	var result ffi.Arg
	maxDevicesFunc.Call(unsafe.Pointer(&result))
	return uint64(result)
}

// MaxParallelSequences returns the maximum number of parallel sequences supported.
func MaxParallelSequences() uint64 {
	var result ffi.Arg
	maxParallelSequencesFunc.Call(unsafe.Pointer(&result))
	return uint64(result)
}

// SupportsMmap checks if memory-mapped files are supported.
func SupportsMmap() bool {
	var result ffi.Arg
	supportsMmapFunc.Call(unsafe.Pointer(&result))
	return result.Bool()
}

// SupportsMlock checks if memory locking is supported.
func SupportsMlock() bool {
	var result ffi.Arg
	supportsMlockFunc.Call(unsafe.Pointer(&result))
	return result.Bool()
}

// SupportsGpuOffload checks if GPU offloading is supported.
func SupportsGpuOffload() bool {
	var result ffi.Arg
	supportsGpuOffloadFunc.Call(unsafe.Pointer(&result))
	return result.Bool()
}

// SupportsRpc checks if RPC is supported.
func SupportsRpc() bool {
	var result ffi.Arg
	supportsRpcFunc.Call(unsafe.Pointer(&result))
	return result.Bool()
}

// TimeUs returns the current time in microseconds.
func TimeUs() int64 {
	var result ffi.Arg
	timeUsFunc.Call(unsafe.Pointer(&result))
	return int64(result)
}
