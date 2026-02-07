package mtmd

import (
	"runtime"
	"testing"
	"time"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func BenchmarkMultimodalInference(b *testing.B) {
	benchmarkSetupOnce(b)

	total := 0
	b.ResetTimer()
	for b.Loop() {
		total += benchmarkMultimodalInference(b, benchMtmdCtx, benchCtx, benchModel, benchTemplate, benchBitmap, "What is in this image?")
	}

	elapsedSeconds := b.Elapsed().Seconds()
	tokensPerSecond := float64(total) / elapsedSeconds
	b.ReportMetric(tokensPerSecond, "tokens/s")

	if runtime.GOOS == "darwin" {
		time.Sleep(time.Second)
	}
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
	for i := uint64(0); i < InputChunksSize(output); i++ {
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
	sp.Temp = float32(0.8)
	sp.TopK = int32(40)
	sp.TopP = float32(0.9)
	sp.MinP = float32(0.1)

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
