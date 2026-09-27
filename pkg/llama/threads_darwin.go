package llama

import "golang.org/x/sys/unix"

// mathCores counts the compute cores on this machine. On Apple Silicon these
// are the perflevel0 cores. An Intel Mac has no perf levels, so it returns the
// physical cores. It returns 0 when the system reports nothing, and the caller
// then uses its own default.
func mathCores() int {
	for _, name := range []string{"hw.perflevel0.physicalcpu", "hw.physicalcpu"} {
		if n, err := unix.SysctlUint32(name); err == nil && n > 0 {
			return int(n)
		}
	}
	return 0
}

// mathCPUs returns nil on macOS. The system does not let a thread pick a CPU,
// so ggml does not pin threads there either.
func mathCPUs() []int32 { return nil }
