package mtmd

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func BenchmarkMultimodalInference(b *testing.B) {
	modelFile := benchmarkModelFileName(b)
	projectFile := benchmarkProjectorFileName(b)

	benchmarkSetup(b)
	defer benchmarkCleanup(b)

	mparams := llama.ModelDefaultParams()
	if os.Getenv("YZMA_BENCHMARK_DEVICE") != "" {
		devs := []llama.GGMLBackendDevice{}
		devices := strings.Split(os.Getenv("YZMA_BENCHMARK_DEVICE"), ",")
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
	defer llama.ModelFree(model)

	params := llama.ContextDefaultParams()
	params.NCtx = 4096
	params.NBatch = 2048

	ctx, err := llama.InitFromModel(model, params)
	if err != nil {
		b.Fatalf("InitFromModel failed: %v", err)
	}
	defer llama.Free(ctx)

	mtmdCtx, err := InitFromFile(projectFile, model, ContextParamsDefault())
	if err != nil {
		fmt.Println("unable to initialize context from file", err.Error())
		os.Exit(1)
	}
	defer Free(mtmdCtx)

	template := llama.ModelChatTemplate(model, "")

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
	if err != nil {
		b.Fatal("count not open file")
	}
	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	total := 0
	b.ResetTimer()
	for b.Loop() {
		total += benchmarkMultimodalInference(b, mtmdCtx, ctx, model, template, bitmap, "What is in this image?")
	}

	// Calculate tokens/second
	elapsedSeconds := b.Elapsed().Seconds()
	tokensPerSecond := float64(total) / elapsedSeconds
	b.ReportMetric(tokensPerSecond, "tokens/s")
}

func benchmarkMultimodalInference(b *testing.B, mctx Context, ctx llama.Context, model llama.Model, template string, bitmap Bitmap, text string) int {
	total := 0
	vocab := llama.ModelGetVocab(model)

	messages := make([]llama.ChatMessage, 0)
	messages = append(messages, llama.NewChatMessage("user", DefaultMarker()+text))
	input := NewInputText(chatTemplate(template, messages, true), true, true)

	output := InputChunksInit()
	defer InputChunksFree(output)

	result := Tokenize(mctx, output, input, []Bitmap{bitmap})
	if result != 0 {
		b.Fatalf("Tokenize failed with result: %d", result)
	}
	for i := uint64(0); i < uint64(InputChunksSize(output)); i++ {
		chunk := InputChunksGet(output, i)
		total += int(InputChunkGetNTokens(chunk))
	}

	var n llama.Pos
	nBatch := llama.NBatch(ctx)
	res := HelperEvalChunks(mctx, ctx, output, 0, 0, int32(nBatch), true, &n)
	if res != 0 {
		b.Fatalf("HelperEvalChunks failed with result: %d", res)
	}

	sp := llama.DefaultSamplerParams()
	sp.Temp = float32(0.7)

	sampler := llama.NewSampler(model, llama.DefaultSamplers, sp)
	defer llama.SamplerFree(sampler)

	for i := 0; i < int(nBatch); i++ {
		token := llama.SamplerSample(sampler, ctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}
		if token == llama.TokenNull {
			b.Fatalf("SamplerSample returned TokenNull")
		}

		total += 1
		buf := make([]byte, 128)
		llama.TokenToPiece(vocab, token, buf, 0, true)

		// Use BatchGetOne instead of manually constructing the batch
		batch := llama.BatchGetOne([]llama.Token{token})
		batch.Pos = &n

		llama.Decode(ctx, batch)
		n++
	}

	b.StopTimer()

	llama.Synchronize(ctx)
	mem, err := llama.GetMemory(ctx)
	if err != nil {
		b.Fatalf("GetMemory failed: %v", err)
	}
	llama.MemoryClear(mem, true)

	b.StartTimer()
	return total
}

func chatTemplate(template string, messages []llama.ChatMessage, add bool) string {
	buf := make([]byte, 1024)
	len := llama.ChatApplyTemplate(template, messages, add, buf)
	result := string(buf[:len])
	return result
}
