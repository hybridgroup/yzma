package speculative

import (
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestNextNZeroContext(t *testing.T) {
	// The guards must hold before any library is loaded.
	SetEmbeddingsNextN(0, true, false)
	SetNextNLayerOffset(0, 0)

	if got := GetEmbeddingsNextN(0, 1, 4); got != nil {
		t.Fatalf("GetEmbeddingsNextN = %v, want nil", got)
	}

	if got := GetEmbeddingsNextNIth(0, 0, 4); got != nil {
		t.Fatalf("GetEmbeddingsNextNIth = %v, want nil", got)
	}
}

func TestAvailable(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// An older llama.cpp build without the NextN functions is a valid state.
	t.Logf("Available returned: %v", Available())
}

func TestGetEmbeddingsNextN(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if !Available() {
		t.Skip("NextN functions not in this llama.cpp build")
	}

	modelFile := testModelFileName(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	ctx, err := llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer llama.Free(ctx)

	tokens := llama.Tokenize(llama.ModelGetVocab(model), "hello world", true, true)
	if len(tokens) == 0 {
		t.Fatal("Tokenize returned no tokens")
	}

	batch := llama.BatchInit(int32(len(tokens)), 0, 1)
	defer llama.BatchFree(batch)

	for i, token := range tokens {
		if err := batch.Add(token, llama.Pos(i), []llama.SeqId{0}, true); err != nil {
			t.Fatalf("batch.Add failed: %v", err)
		}
	}

	SetEmbeddingsNextN(ctx, true, false)

	if _, err := llama.Decode(ctx, batch); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	nEmbd := int(llama.ModelNEmbd(model))

	embd := GetEmbeddingsNextN(ctx, len(tokens), nEmbd)
	if embd == nil {
		t.Skip("model has no NextN layers")
	}

	if len(embd) != len(tokens)*nEmbd {
		t.Fatalf("len(GetEmbeddingsNextN) = %d, want %d", len(embd), len(tokens)*nEmbd)
	}

	row := GetEmbeddingsNextNIth(ctx, 0, nEmbd)
	if len(row) != nEmbd {
		t.Fatalf("len(GetEmbeddingsNextNIth) = %d, want %d", len(row), nEmbd)
	}

	for i := range row {
		if row[i] != embd[i] {
			t.Fatalf("row[%d] = %v, want %v", i, row[i], embd[i])
		}
	}
}
