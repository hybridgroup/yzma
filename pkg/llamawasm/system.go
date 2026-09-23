//go:build js && wasm

package llamawasm

// The calls here describe llama.cpp and the machine. Each one follows the same
// call in pkg/llama. A module before ABI version 9 has none of them, thus each
// gives a zero value.

// PrintSystemInfo gives the features of the CPU and of the backends that
// llama.cpp uses.
func PrintSystemInfo() string {
	return callString("_yzma_print_system_info", 1024)
}

// FtypeName gives the name of a kind of quantization, as ModelFtype gives it.
func FtypeName(ftype Ftype) string {
	return callString("_yzma_ftype_name", 64, int(ftype))
}

// TimeUs gives the time of llama.cpp in microseconds.
func TimeUs() int64 {
	if !has("_yzma_time_us") {
		return 0
	}
	return int64(callValue("_yzma_time_us").Float())
}

// MaxDevices gives the largest number of devices that llama.cpp can use.
func MaxDevices() uint64 { return systemUint64("_yzma_max_devices") }

// MaxParallelSequences gives the largest number of sequences of a context.
func MaxParallelSequences() uint64 { return systemUint64("_yzma_max_parallel_sequences") }

// SupportsGpuOffload tells if the module has a backend that is not the CPU.
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
