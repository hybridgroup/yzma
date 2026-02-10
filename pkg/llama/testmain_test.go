package llama

import (
	"flag"
	"os"
	"strings"
	"testing"
)

var (
	benchModel    Model
	benchCtx      Context
	benchTemplate string
	benchReady    bool
)

var (
	nCtx   int
	device string
)

func init() {
	flag.IntVar(&nCtx, "nctx", 8192, "number of context tokens for llama.Context")
	flag.StringVar(&device, "device", "", "comma-separated list of devices to use for benchmarking (e.g. 'CUDA0')")
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
	mparams.UseMmap = 0

	if device != "" {
		devs := []GGMLBackendDevice{}
		devices := strings.Split(device, ",")
		for _, d := range devices {
			dev := GGMLBackendDeviceByName(d)
			if dev == 0 {
				b.Fatalf("unknown device: %s", d)
			}
			devs = append(devs, dev)
		}

		mparams.SetDevices(devs)
	}

	model, err := ModelLoadFromFile(modelFile, mparams)
	if err != nil {
		b.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	benchModel = model

	params := ContextDefaultParams()
	params.NBatch = 1024
	params.NCtx = uint32(nCtx)

	ctx, err := InitFromModel(model, params)
	if err != nil {
		b.Fatalf("InitFromModel failed: %v", err)
	}
	benchCtx = ctx

	benchTemplate = ModelChatTemplate(model, "")

	benchReady = true
}

func benchmarkTeardown() {
	Free(benchCtx)
	ModelFree(benchModel)

	LogSet(LogNormal)
	Close()
}
