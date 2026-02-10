package mtmd

import (
	"flag"
	"fmt"
	"os"
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
	projectFile := benchmarkProjectorFileName(b)

	benchmarkSetup(b)

	mparams := llama.ModelDefaultParams()
	mparams.UseMmap = 0

	if device != "" {
		devs := []llama.GGMLBackendDevice{}
		devices := strings.Split(device, ",")
		for _, d := range devices {
			dev := llama.GGMLBackendDeviceByName(d)
			if dev == 0 {
				b.Fatalf("unknown device: %s", d)
			}
			devs = append(devs, dev)
		}

		mparams.SetDevices(devs)
	}

	model, err := llama.ModelLoadFromFile(modelFile, mparams)
	if err != nil {
		b.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	benchModel = model

	params := llama.ContextDefaultParams()
	params.NCtx = uint32(nCtx)
	params.NBatch = 1024

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
