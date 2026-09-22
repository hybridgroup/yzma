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

// SchedPriority is the priority that the threads of a pool run at.
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

// ThreadpoolParams says how to make a thread pool. CPUMask holds one entry for
// each CPU of the machine. A mask of all zeros lets the system put the threads
// where it wants, which costs speed on a machine with two kinds of core.
type ThreadpoolParams struct {
	CPUMask   [maxThreads]uint8
	NThreads  int32
	Prio      SchedPriority
	Poll      uint32
	StrictCPU uint8
	Paused    uint8
}

// SetCPUs puts one CPU of the given list in the mask and asks for strict
// placement, which gives each thread one CPU of its own.
func (p *ThreadpoolParams) SetCPUs(cpus []int32) {
	p.CPUMask = [maxThreads]uint8{}
	for _, cpu := range cpus {
		if cpu >= 0 && cpu < maxThreads {
			p.CPUMask[cpu] = 1
		}
	}
	p.StrictCPU = 1
}

// ThreadpoolParamsDefault gives the parameters of a pool of n threads with no
// CPU mask.
func ThreadpoolParamsDefault(nThreads int32) ThreadpoolParams {
	var p ThreadpoolParams
	pp := &p
	ggmlThreadpoolParamsInitFunc.Call(nil, unsafe.Pointer(&pp), unsafe.Pointer(&nThreads))
	return p
}

var (
	// ErrNoThreadpool says that this build of llama.cpp makes no thread pool.
	ErrNoThreadpool = errors.New("the CPU backend has no thread pool")

	threadpoolNewFn  ffi.Fun
	threadpoolFreeFn ffi.Fun
)

// cpuProcAddress gives the address of a call of the CPU backend. The calls of
// a backend are not in the shared library of llama.cpp. They come from the
// register of the backend, which is how llama.cpp finds them.
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

// loadThreadpoolFuncs finds the calls of the thread pool. Call it after the
// backend starts, because the register of the CPU backend exists only then.
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

// ThreadpoolNew makes a thread pool. Give it to a context with
// [AttachThreadpool] and free it with [ThreadpoolFree] after the context goes.
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

// NewPerformanceThreadpool makes a thread pool of one thread for each
// performance core, each held to a CPU of its own. It gives ErrNoThreadpool
// when the build of llama.cpp makes no pool and ErrNoPerformanceCPUs when the
// system does not say which CPUs to use.
//
// Give the pool to a context with [AttachThreadpool]. The model of that
// context must come from [ModelParams.SetCPUOnly], not from a device list that
// names the CPU, because such a list leaves the pool with no work.
func NewPerformanceThreadpool() (Threadpool, error) {
	cpus := PerformanceCPUs()
	if len(cpus) == 0 {
		return 0, ErrNoPerformanceCPUs
	}

	params := ThreadpoolParamsDefault(int32(len(cpus)))
	params.SetCPUs(cpus)

	return ThreadpoolNew(&params)
}
