package llama

import "golang.org/x/sys/unix"

// mathCores counts the cores of this machine that do the arithmetic well. The
// first level of performance holds the performance cores of an Apple Silicon
// machine. An Intel Mac has no level and gives the physical cores. The count
// is 0 when the system says nothing, and then the caller uses its own default.
func mathCores() int {
	for _, name := range []string{"hw.perflevel0.physicalcpu", "hw.physicalcpu"} {
		if n, err := unix.SysctlUint32(name); err == nil && n > 0 {
			return int(n)
		}
	}
	return 0
}

// mathCPUs gives nothing on macOS. The system does not let a thread choose a
// CPU, thus ggml holds no thread to a core there either.
func mathCPUs() []int32 { return nil }
