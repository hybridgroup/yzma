package llama

import (
	"runtime"
	"testing"
)

func TestThreads(t *testing.T) {
	n := Threads()
	if n < 1 {
		t.Errorf("Threads gave %d, want at least 1", n)
	}
	if n > int32(runtime.NumCPU()) {
		t.Errorf("Threads gave %d, want no more than the %d logical CPUs", n, runtime.NumCPU())
	}
	t.Logf("Threads gave %d of %d logical CPUs", n, runtime.NumCPU())
}

func TestDefaultThreads(t *testing.T) {
	if n := defaultThreads(); n < 1 {
		t.Errorf("defaultThreads gave %d, want at least 1", n)
	}
}
