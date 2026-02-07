package llama

import (
	"testing"
)

func BenchmarkInference(b *testing.B) {
	benchmarkSetupOnce(b)

	total := 0
	b.ResetTimer()
	for b.Loop() {
		total += benchmarkInference(b, benchCtx, benchModel, "Are you ready to go?")
	}

	elapsedSeconds := b.Elapsed().Seconds()
	tokensPerSecond := float64(total) / elapsedSeconds
	b.ReportMetric(tokensPerSecond, "tokens/s")
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
