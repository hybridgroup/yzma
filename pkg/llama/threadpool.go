package llama

import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

// Threadpool is a set of threads that ggml uses to compute a graph.
type Threadpool uintptr

// SchedPriority is the priority that pool threads run at.
type SchedPriority int32

const (
	SchedPriorityLow SchedPriority = iota - 1
	SchedPriorityNormal
	SchedPriorityMedium
	SchedPriorityHigh
	SchedPriorityRealtime
)

// maxThreads is GGML_MAX_N_THREADS, the length of the CPU mask.
const maxThreads = 512

// ThreadpoolParams configures a thread pool. CPUMask holds one entry per CPU.
// An all zero mask lets the system place the threads, which costs speed on a
// hybrid machine.
type ThreadpoolParams struct {
	CPUMask   [maxThreads]uint8
	NThreads  int32
	Prio      SchedPriority
	Poll      uint32
	StrictCPU uint8
	Paused    uint8
}

// SetCPUs sets the given CPUs in the mask and turns on strict placement, so
// each thread gets its own CPU.
func (p *ThreadpoolParams) SetCPUs(cpus []int32) {
	p.CPUMask = [maxThreads]uint8{}
	for _, cpu := range cpus {
		if cpu >= 0 && cpu < maxThreads {
			p.CPUMask[cpu] = 1
		}
	}
	p.StrictCPU = 1
}

// ThreadpoolParamsDefault returns the parameters for a pool of n threads with
// no CPU mask.
func ThreadpoolParamsDefault(nThreads int32) ThreadpoolParams {
	var p ThreadpoolParams
	pp := &p
	ggmlThreadpoolParamsInitFunc.Call(nil, unsafe.Pointer(&pp), unsafe.Pointer(&nThreads))
	return p
}

var (
	// ErrNoThreadpool means this llama.cpp build cannot create a thread pool.
	ErrNoThreadpool = errors.New("the CPU backend has no thread pool")

	threadpoolNewFn  ffi.Fun
	threadpoolFreeFn ffi.Fun
)

// cpuProcAddress returns the address of a CPU backend function. Backend
// functions are not exported by the llama.cpp shared library. They come from
// the backend registry, which is how llama.cpp finds them.
func cpuProcAddress(name string) uintptr {
	dev := GGMLBackendDeviceByType(GGMLBackendDeviceTypeCPU)
	if dev == 0 {
		return 0
	}
	reg := GGMLBackendDeviceBackendReg(dev)
	if reg == 0 {
		return 0
	}

	cname, err := utils.BytePtrFromString(name)
	if err != nil {
		return 0
	}

	var addr uintptr
	ggmlBackendRegGetProcAddressFunc.Call(unsafe.Pointer(&addr), unsafe.Pointer(&reg), unsafe.Pointer(&cname))
	runtime.KeepAlive(cname)

	return addr
}

// loadThreadpoolFuncs binds the thread pool functions. Call it after the
// backend starts, because the CPU backend registry exists only then.
func loadThreadpoolFuncs() error {
	newAddr := cpuProcAddress("ggml_threadpool_new")
	freeAddr := cpuProcAddress("ggml_threadpool_free")
	if newAddr == 0 || freeAddr == 0 {
		return ErrNoThreadpool
	}

	threadpoolNewFn = ffi.Fun{Addr: newAddr, Cif: new(ffi.Cif)}
	if s := ffi.PrepCif(threadpoolNewFn.Cif, ffi.DefaultAbi, 1, &ffi.TypePointer, &ffi.TypePointer); s != ffi.OK {
		return loadError("ggml_threadpool_new", errors.New(s.String()))
	}

	threadpoolFreeFn = ffi.Fun{Addr: freeAddr, Cif: new(ffi.Cif)}
	if s := ffi.PrepCif(threadpoolFreeFn.Cif, ffi.DefaultAbi, 1, &ffi.TypeVoid, &ffi.TypePointer); s != ffi.OK {
		return loadError("ggml_threadpool_free", errors.New(s.String()))
	}

	return nil
}

// ThreadpoolNew creates a thread pool. Attach it to a context with
// [AttachThreadpool] and free it with [ThreadpoolFree] after the context is freed.
func ThreadpoolNew(params *ThreadpoolParams) (Threadpool, error) {
	if params == nil {
		return 0, errors.New("no thread pool parameters")
	}
	if threadpoolNewFn.Addr == 0 {
		if err := loadThreadpoolFuncs(); err != nil {
			return 0, err
		}
	}

	var tp Threadpool
	threadpoolNewFn.Call(unsafe.Pointer(&tp), unsafe.Pointer(&params))
	runtime.KeepAlive(params)
	if tp == 0 {
		return 0, ErrNoThreadpool
	}

	return tp, nil
}

// ThreadpoolFree frees a thread pool. Detach it from every context first.
func ThreadpoolFree(tp Threadpool) {
	if tp == 0 || threadpoolFreeFn.Addr == 0 {
		return
	}
	threadpoolFreeFn.Call(nil, unsafe.Pointer(&tp))
}

// NewPerformanceThreadpool creates a pool with one thread per performance
// core, each pinned to its own CPU. It returns ErrNoThreadpool when the
// llama.cpp build cannot create a pool, and ErrNoPerformanceCPUs when the
// system does not report which CPUs to use.
//
// Attach the pool with [AttachThreadpool]. The context's model must be loaded
// with [ModelParams.SetCPUOnly], not a device list that names the CPU, because
// such a list leaves the pool with no work.
func NewPerformanceThreadpool() (Threadpool, error) {
	cpus := PerformanceCPUs()
	if len(cpus) == 0 {
		return 0, ErrNoPerformanceCPUs
	}

	params := ThreadpoolParamsDefault(int32(len(cpus)))
	params.SetCPUs(cpus)

	return ThreadpoolNew(&params)
}
