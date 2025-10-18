package llama

import (
	"testing"
)

func TestContextDefaultParams(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ContextDefaultParams()
	if params == (ContextParams{}) {
		t.Fatal("ContextDefaultParams returned empty parameters")
	}
}

func TestInitModel(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	params := ContextDefaultParams()
	ctx := InitFromModel(model, params)
	if ctx == (Context(0)) {
		t.Fatal("Failed to initialize context")
	}

	Free(ctx)
	// No direct way to verify, but ensure no panic or error occurs
}

func TestWarmup(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	params := ContextDefaultParams()
	ctx := InitFromModel(model, params)
	if ctx == (Context(0)) {
		t.Fatal("Failed to initialize context")
	}
	defer Free(ctx)

	SetWarmup(ctx, true)
	SetWarmup(ctx, false)
	// No direct way to verify, but ensure no panic or error occurs
}

func TestNCtx(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	nCtx := NCtx(ctx)
	if nCtx == 0 {
		t.Fatal("NCtx returned 0, which is invalid")
	}
	t.Logf("NCtx returned: %d", nCtx)
}

func TestNBatch(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	nBatch := NBatch(ctx)
	if nBatch == 0 {
		t.Fatal("NBatch returned 0, which is invalid")
	}
	t.Logf("NBatch returned: %d", nBatch)
}

func TestNUBatch(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	nUBatch := NUBatch(ctx)
	if nUBatch == 0 {
		t.Fatal("NUBatch returned 0, which is invalid")
	}
	t.Logf("NUBatch returned: %d", nUBatch)
}

func TestNSeqMax(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	nSeqMax := NSeqMax(ctx)
	if nSeqMax == 0 {
		t.Fatal("NSeqMax returned 0, which is invalid")
	}
	t.Logf("NSeqMax returned: %d", nSeqMax)
}

func TestGetModel(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	retrievedModel := GetModel(ctx)
	if retrievedModel == (Model(0)) {
		t.Fatal("GetModel returned an empty model")
	}

	t.Logf("GetModel successfully retrieved the model: %v", retrievedModel)
}

func TestGetLogitsIth(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// tokenize prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	count := Tokenize(vocab, prompt, nil, true, true)
	tokens := make([]Token, count)
	Tokenize(vocab, prompt, tokens, true, true)

	// create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	nVocab := int(VocabNTokens(vocab))
	logits := GetLogitsIth(ctx, -1, nVocab)
	if logits == nil {
		t.Fatal("GetLogitsIth returned nil")
	}
	t.Logf("GetLogitsIth returned %d logits", len(logits))
}

func TestGetEmbeddingsIth(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	params := ContextDefaultParams()
	params.PoolingType = PoolingTypeMean
	params.Embeddings = 1

	ctx := InitFromModel(model, params)
	defer Free(ctx)

	// Tokenize a prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	count := Tokenize(vocab, prompt, nil, true, true)
	tokens := make([]Token, count)
	Tokenize(vocab, prompt, tokens, true, true)

	// Create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	nEmbeddings := VocabNTokens(vocab)
	embeddings := GetEmbeddingsIth(ctx, -1, nEmbeddings) // Get embeddings for the last token
	if embeddings == nil {
		t.Fatal("GetEmbeddingsIth returned nil")
	}
	if len(embeddings) != int(nEmbeddings) {
		t.Fatalf("GetEmbeddingsIth returned %d embeddings, expected %d", len(embeddings), nEmbeddings)
	}
	t.Logf("GetEmbeddingsIth returned %d embeddings", len(embeddings))
}

func TestGetEmbeddingsSeq(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	params := ContextDefaultParams()
	params.PoolingType = PoolingTypeMean
	params.Embeddings = 1

	ctx := InitFromModel(model, params)
	defer Free(ctx)

	// Tokenize a prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	count := Tokenize(vocab, prompt, nil, true, true)
	tokens := make([]Token, count)
	Tokenize(vocab, prompt, tokens, true, true)

	// Create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	seqID := SeqId(0)
	nEmbeddings := int32(VocabNTokens(vocab))
	embeddings := GetEmbeddingsSeq(ctx, seqID, nEmbeddings)
	if embeddings == nil {
		t.Fatal("GetEmbeddingsSeq returned nil")
	}
	if len(embeddings) != int(nEmbeddings) {
		t.Fatalf("GetEmbeddingsSeq returned %d embeddings, expected %d", len(embeddings), nEmbeddings)
	}
	t.Logf("GetEmbeddingsSeq returned %d embeddings", len(embeddings))
}

func TestSetEmbeddings(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// Enable embeddings
	SetEmbeddings(ctx, true)
	t.Log("SetEmbeddings successfully set embeddings to true")

	// Disable embeddings
	SetEmbeddings(ctx, false)
	t.Log("SetEmbeddings successfully set embeddings to false")
}

func TestSetCausalAttn(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// Enable causal attention
	SetCausalAttn(ctx, true)
	t.Log("SetCausalAttn successfully set causal attention to true")

	// Disable causal attention
	SetCausalAttn(ctx, false)
	t.Log("SetCausalAttn successfully set causal attention to false")
}
