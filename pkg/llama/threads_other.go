//go:build !linux && !darwin && !windows

package llama

// mathCores returns 0 on systems where yzma cannot query the cores, so the
// caller uses its own default.
func mathCores() int { return 0 }

// mathCPUs returns nil on systems where yzma cannot query the cores.
func mathCPUs() []int32 { return nil }
