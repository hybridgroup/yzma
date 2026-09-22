//go:build !linux && !darwin

package llama

// mathCores gives 0 on a system that yzma cannot ask about the cores, and then
// the caller uses its own default.
func mathCores() int { return 0 }
