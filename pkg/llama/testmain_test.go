package llama

import (
	"flag"
	"os"
	"runtime"
	"strings"
	"testing"
)

var (
	benchModel    Model
	benchCtx      Context
	benchTemplate string
	benchReady    bool
)

// benchThreads is the thread count of the text benchmark on each machine. The
// model is too small to use more threads, thus more threads make it slower.
const benchThreads = 4

var (
	nCtx       int
	nThreads   int
	device     string
	threadpool bool

	benchThreadpool Threadpool
)

func init() {
	flag.IntVar(&nCtx, "nctx", 8192, "number of context tokens for llama.Context")
	flag.IntVar(&nThreads, "threads", benchThreads, "number of CPU threads, 0 for the value of llama.Threads")
	flag.StringVar(&device, "device", "", "comma-separated list of devices to use for benchmarking (e.g. 'CUDA0')")
	flag.BoolVar(&threadpool, "threadpool", false, "hold the CPU threads to the performance cores")
}

func TestMain(m *testing.M) {
	flag.Parse() // Parse flags before running tests

	code := m.Run()

	if benchReady {
		benchmarkTeardown()
	}

	os.Exit(code)
}

func benchmarkSetupOnce(b *testing.B) {
	if benchReady {
		return
	}

	modelFile := benchmarkModelFileName(b)

	benchmarkSetup(b)

	mparams := ModelDefaultParams()
	mparams.LoadMode = LoadModeNone

	if strings.EqualFold(device, "CPU") {
		mparams.SetCPUOnly()
	} else if device != "" {
		devs := []GGMLBackendDevice{}
		devices := strings.SplitSeq(device, ",")
		for d := range devices {
			dev := GGMLBackendDeviceByName(d)
			if dev == 0 {
				b.Fatalf("unknown device: %s", d)
			}
			devs = append(devs, dev)
		}

		devs = append(devs, 0) // NULL terminator required by llama.cpp
		if err := mparams.SetDevices(devs); err != nil {
			b.Fatalf("SetDevices failed: %v", err)
		}
		defer runtime.KeepAlive(devs)
	}

	model, err := ModelLoadFromFile(modelFile, mparams)
	if err != nil {
		b.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	benchModel = model

	params := ContextDefaultParams()
	params.NBatch = 1024
	params.NCtx = uint32(nCtx)
	if nThreads > 0 {
		params.NThreads = int32(nThreads)
		params.NThreadsBatch = int32(nThreads)
	}

	ctx, err := InitFromModel(model, params)
	if err != nil {
		b.Fatalf("InitFromModel failed: %v", err)
	}
	benchCtx = ctx

	if threadpool {
		tp, err := newBenchThreadpool(params.NThreads)
		if err != nil {
			b.Fatalf("newBenchThreadpool failed: %v", err)
		}
		AttachThreadpool(ctx, uintptr(tp), uintptr(tp))
		benchThreadpool = tp
	}

	benchTemplate = ModelChatTemplate(model, "")

	benchReady = true
}

// newBenchThreadpool holds n threads to the first n performance CPUs, thus the
// pool has the thread count of the context.
func newBenchThreadpool(n int32) (Threadpool, error) {
	cpus := PerformanceCPUs()
	if len(cpus) == 0 {
		return 0, ErrNoPerformanceCPUs
	}
	if int(n) < len(cpus) {
		cpus = cpus[:n]
	}

	params := ThreadpoolParamsDefault(int32(len(cpus)))
	params.SetCPUs(cpus)
	return ThreadpoolNew(&params)
}

func benchmarkTeardown() {
	if benchThreadpool != 0 {
		DetachThreadpool(benchCtx)
	}
	Free(benchCtx)
	if benchThreadpool != 0 {
		ThreadpoolFree(benchThreadpool)
	}
	ModelFree(benchModel)

	LogSet(LogNormal)
	Close()
}
