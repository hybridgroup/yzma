//go:build js && wasm

package llamawasm

// The calls here describe llama.cpp and the machine. Each one follows the same
// call in pkg/llama. A module before ABI version 9 has none of them, so each
// returns a zero value.

// PrintSystemInfo returns the features of the CPU and the backends that
// llama.cpp uses.
func PrintSystemInfo() string {
	return callString("_yzma_print_system_info", 1024)
}

// FtypeName returns the name of a quantization type, as returned by ModelFtype.
func FtypeName(ftype Ftype) string {
	return callString("_yzma_ftype_name", 64, int(ftype))
}

// TimeUs returns the llama.cpp time in microseconds.
func TimeUs() int64 {
	if !has("_yzma_time_us") {
		return 0
	}
	return int64(callValue("_yzma_time_us").Float())
}

// MaxDevices returns the maximum number of devices llama.cpp can use.
func MaxDevices() uint64 { return systemUint64("_yzma_max_devices") }

// MaxParallelSequences returns the maximum number of sequences in a context.
func MaxParallelSequences() uint64 { return systemUint64("_yzma_max_parallel_sequences") }

// SupportsGpuOffload reports whether the module has a non CPU backend.
func SupportsGpuOffload() bool { return systemUint64("_yzma_supports_gpu_offload") == 1 }

func systemUint64(name string) uint64 {
	if !has(name) {
		return 0
	}
	n := call(name)
	if n < 0 {
		return 0
	}
	return uint64(n)
}
