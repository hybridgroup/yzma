package llama

import (
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

func TestMain(m *testing.M) {
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
	bd := os.Getenv("YZMA_BENCHMARK_DEVICE")
	if bd != "" {
		devs := []GGMLBackendDevice{}
		devices := strings.Split(bd, ",")
		for _, d := range devices {
			dev := GGMLBackendDeviceByName(d)
			if dev == 0 {
				b.Fatalf("unknown device: %s", d)
			}
			devs = append(devs, dev)
		}

		mparams.SetDevices(devs)
	}

	params := ContextDefaultParams()

	switch {
	case strings.Contains(bd, "CUDA"), strings.Contains(bd, "VULKAN"):
		mparams.UseDirectIO = 1
		params.NCtx = 32000

	case runtime.GOOS == "darwin":
		params.NCtx = 16000

	default:
		params.NCtx = 8192
	}

	params.NBatch = 1024
	mparams.UseMmap = 0

	model, err := ModelLoadFromFile(modelFile, mparams)
	if err != nil {
		b.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	benchModel = model

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
