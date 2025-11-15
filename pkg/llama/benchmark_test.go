package llama

import (
	"testing"
)

func BenchmarkInference(b *testing.B) {
	modelFile := benchmarkModelFileName(b)

	benchmarkSetup(b)
	defer benchmarkCleanup(b)
	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
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

	for i := 0; i < b.N; i++ {
		benchmarkInference(b, ctx, model, "This is a test")
	}
}

func benchmarkInference(b *testing.B, ctx Context, model Model, text string) {
	vocab := ModelGetVocab(model)

	// get tokens from the prompt
	tokens := Tokenize(vocab, text, true, false)

	batch := BatchGetOne(tokens)

	sampler := SamplerChainInit(SamplerChainDefaultParams())
	SamplerChainAdd(sampler, SamplerInitGreedy())

	for pos := int32(0); pos < 24; pos += batch.NTokens {
		Decode(ctx, batch)
		token := SamplerSample(sampler, ctx, -1)
		if VocabIsEOG(vocab, token) {
			break
		}

		buf := make([]byte, 36)
		TokenToPiece(vocab, token, buf, 0, true)

		batch = BatchGetOne([]Token{token})
	}

	Synchronize(ctx)
	mem, err := GetMemory(ctx)
	if err != nil {
		b.Fatalf("GetMemory failed: %v", err)
	}
	MemoryClear(mem, true)
}
