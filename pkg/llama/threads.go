package llama

import (
	"errors"
	"runtime"
)

// Threads gives a good number of threads for inference on the CPU of this
// machine. It counts the cores that do the arithmetic well, which is one
// thread for each physical core, and only the performance cores of a machine
// that has two kinds. The count comes from the operating system, and a machine
// that tells nothing gets half of its logical CPUs.
//
// llama.cpp asks for four threads unless a caller changes it, which is slow on
// a machine with many cores. Thus [ContextDefaultParams] sends this value.
func Threads() int32 {
	if n := mathCores(); n > 0 {
		return int32(n)
	}
	return int32(defaultThreads())
}

// ErrNoPerformanceCPUs says that the system does not tell which CPUs belong
// to the performance cores.
var ErrNoPerformanceCPUs = errors.New("the system does not name the performance CPUs")

// PerformanceCPUs gives one CPU for each core that does the arithmetic well,
// which is one CPU of each performance core. It gives nothing when the system
// says nothing. Use it to hold the threads of a pool to a core, see
// [NewPerformanceThreadpool].
func PerformanceCPUs() []int32 {
	return mathCPUs()
}

// defaultThreads is the count that llama.cpp uses when it can read nothing
// about the cores. Half of the logical CPUs is one thread for each physical
// core of a machine with SMT.
func defaultThreads() int {
	n := runtime.NumCPU()
	if n > 4 {
		n /= 2
	}
	if n < 1 {
		n = 1
	}
	return n
}
