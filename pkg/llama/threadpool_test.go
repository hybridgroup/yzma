package llama

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestThreadpoolParamsLayout(t *testing.T) {
	// The struct crosses to C by address, thus a wrong size or a wrong place
	// of a field gives a pool that ignores the mask without saying so.
	var p ThreadpoolParams
	if got := unsafe.Sizeof(p); got != 528 {
		t.Errorf("ThreadpoolParams is %d bytes, want 528", got)
	}
	if got := unsafe.Offsetof(p.NThreads); got != 512 {
		t.Errorf("NThreads is at %d, want 512", got)
	}
	if got := unsafe.Offsetof(p.StrictCPU); got != 524 {
		t.Errorf("StrictCPU is at %d, want 524", got)
	}
}

func TestThreadpoolParamsSetCPUs(t *testing.T) {
	var p ThreadpoolParams
	p.SetCPUs([]int32{0, 2, 4, maxThreads, -1})

	for _, cpu := range []int{0, 2, 4} {
		if p.CPUMask[cpu] != 1 {
			t.Errorf("CPUMask[%d] is 0, want 1", cpu)
		}
	}
	if p.CPUMask[1] != 0 {
		t.Error("CPUMask[1] is 1, want 0")
	}
	if p.StrictCPU != 1 {
		t.Error("SetCPUs left StrictCPU at 0, want 1")
	}

	// A second call replaces the mask rather than adding to it.
	p.SetCPUs([]int32{1})
	if p.CPUMask[0] != 0 || p.CPUMask[1] != 1 {
		t.Error("SetCPUs added to the old mask, want a new one")
	}
}

func TestPerformanceCPUs(t *testing.T) {
	cpus := PerformanceCPUs()
	if len(cpus) > runtime.NumCPU() {
		t.Errorf("PerformanceCPUs gave %d CPUs, want no more than %d", len(cpus), runtime.NumCPU())
	}
	seen := map[int32]bool{}
	for _, cpu := range cpus {
		if cpu < 0 || int(cpu) >= runtime.NumCPU() {
			t.Errorf("PerformanceCPUs gave CPU %d, want one of this machine", cpu)
		}
		if seen[cpu] {
			t.Errorf("PerformanceCPUs gave CPU %d twice", cpu)
		}
		seen[cpu] = true
	}
	if n := len(cpus); n > 0 && int32(n) != Threads() {
		t.Errorf("PerformanceCPUs gave %d CPUs, want the %d of Threads", n, Threads())
	}
}

func TestThreadpoolNewNoParams(t *testing.T) {
	if _, err := ThreadpoolNew(nil); err == nil {
		t.Error("ThreadpoolNew(nil) gave no error, want one")
	}
}

func TestThreadpoolFreeZero(t *testing.T) {
	ThreadpoolFree(0) // must not panic
}

func TestSetCPUOnly(t *testing.T) {
	p := ModelParams{NGpuLayers: 99, Devices: 1}
	p.SetCPUOnly()

	if p.NGpuLayers != 0 {
		t.Errorf("NGpuLayers is %d, want 0", p.NGpuLayers)
	}
	if p.Devices != 0 {
		t.Errorf("Devices is %d, want 0", p.Devices)
	}
}
