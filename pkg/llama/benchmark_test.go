package llama

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func BenchmarkInference(b *testing.B) {
	modelFile := benchmarkModelFileName(b)

	benchmarkSetup(b)
	defer benchmarkCleanup(b)

	mparams := ModelDefaultParams()
	if os.Getenv("YZMA_BENCHMARK_DEVICE") != "" {
		devs := []GGMLBackendDevice{}
		devices := strings.Split(os.Getenv("YZMA_BENCHMARK_DEVICE"), ",")
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
	defer ModelFree(model)

	params := ContextDefaultParams()
	params.NCtx = 4096
	params.NBatch = 2048

	ctx, err := InitFromModel(model, params)
	if err != nil {
		b.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	total := 0
	b.ResetTimer()
	for b.Loop() {
		total += benchmarkInference(b, ctx, model, "Are you ready to go?")
	}

	// Calculate tokens/second
	elapsedSeconds := b.Elapsed().Seconds()
	tokensPerSecond := float64(total) / elapsedSeconds
	b.ReportMetric(tokensPerSecond, "tokens/s")

	// extra time to cleanup
	if runtime.GOOS == "darwin" {
		time.Sleep(time.Second)
	}
}

func benchmarkInference(b *testing.B, ctx Context, model Model, text string) int {
	vocab := ModelGetVocab(model)

	// get tokens from the prompt
	tokens := Tokenize(vocab, text, true, false)
	total := len(tokens)

	batch := BatchGetOne(tokens)

	sampler := SamplerChainInit(SamplerChainDefaultParams())
	SamplerChainAdd(sampler, SamplerInitGreedy())

	for pos := int32(0); pos < 24; pos += batch.NTokens {
		Decode(ctx, batch)
		token := SamplerSample(sampler, ctx, -1)
		if VocabIsEOG(vocab, token) {
			break
		}

		total++
		buf := make([]byte, 36)
		TokenToPiece(vocab, token, buf, 0, true)

		batch = BatchGetOne([]Token{token})
	}

	b.StopTimer()

	Synchronize(ctx)
	mem, err := GetMemory(ctx)
	if err != nil {
		b.Fatalf("GetMemory failed: %v", err)
	}
	MemoryClear(mem, true)

	b.StartTimer()

	return total
}
