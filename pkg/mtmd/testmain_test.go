package mtmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

var (
	benchModel    llama.Model
	benchCtx      llama.Context
	benchMtmdCtx  Context
	benchTemplate string
	benchBitmap   Bitmap
	benchImgData  []byte
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
	projectFile := benchmarkProjectorFileName(b)

	benchmarkSetup(b)

	mparams := llama.ModelDefaultParams()
	bd := os.Getenv("YZMA_BENCHMARK_DEVICE")
	if bd != "" {
		devs := []llama.GGMLBackendDevice{}
		devices := strings.Split(bd, ",")
		for _, d := range devices {
			dev := llama.GGMLBackendDeviceByName(d)
			if dev == 0 {
				b.Fatalf("unknown device: %s", d)
			}
			devs = append(devs, dev)
		}

		mparams.SetDevices(devs)
	}

	params := llama.ContextDefaultParams()

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

	model, err := llama.ModelLoadFromFile(modelFile, mparams)
	if err != nil {
		b.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	benchModel = model

	ctx, err := llama.InitFromModel(model, params)
	if err != nil {
		b.Fatalf("InitFromModel failed: %v", err)
	}
	benchCtx = ctx

	mprms := ContextParamsDefault()
	mprms.ImageMinTokens = 1024

	mtmdCtx, err := InitFromFile(projectFile, model, mprms)
	if err != nil {
		fmt.Println("unable to initialize context from file", err.Error())
		os.Exit(1)
	}
	benchMtmdCtx = mtmdCtx

	benchTemplate = llama.ModelChatTemplate(model, "")

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
	if err != nil {
		b.Fatal("could not open file")
	}
	benchImgData = data
	benchBitmap = BitmapInit(x, y, uintptr(unsafe.Pointer(&benchImgData[0])))

	benchReady = true
}

func benchmarkTeardown() {
	BitmapFree(benchBitmap)
	Free(benchMtmdCtx)
	llama.Free(benchCtx)
	llama.ModelFree(benchModel)

	llama.LogSet(llama.LogNormal)
	LogSet(llama.LogNormal)
	llama.Close()
}
