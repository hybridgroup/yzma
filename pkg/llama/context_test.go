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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	params := ContextDefaultParams()
	ctx, err := InitFromModel(model, params)
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}

	if err := Free(ctx); err != nil {
		t.Fatalf("Free failed: %v", err)
	}
}

func TestWarmup(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	params := ContextDefaultParams()
	ctx, err := InitFromModel(model, params)
	if ctx == 0 || err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	SetWarmup(ctx, true)

	Warmup(ctx, model)
	t.Log("Warmup completed successfully")

	SetWarmup(ctx, false)
}

func TestNCtx(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	// tokenize prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, prompt, true, true)

	// create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	nVocab := int(VocabNTokens(vocab))
	logits, err := GetLogitsIth(ctx, -1, nVocab)
	if err != nil {
		t.Fatalf("GetLogitsIth returned an error: %v", err)
	}
	if logits == nil {
		t.Fatal("GetLogitsIth returned nil")
	}
	t.Logf("GetLogitsIth returned %d logits", len(logits))
}

func TestGetEmbeddingsIth(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	params := ContextDefaultParams()
	params.PoolingType = PoolingTypeMean
	params.Embeddings = 1

	ctx, err := InitFromModel(model, params)
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	// Tokenize a prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, prompt, true, true)

	// Create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	nEmbeddings := VocabNTokens(vocab)
	embeddings, err := GetEmbeddingsIth(ctx, -1, nEmbeddings) // Get embeddings for the last token
	if err != nil {
		t.Fatalf("GetEmbeddingsIth returned an error: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	params := ContextDefaultParams()
	params.PoolingType = PoolingTypeMean
	params.Embeddings = 1

	ctx, err := InitFromModel(model, params)
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	// Tokenize a prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, prompt, true, true)

	// Create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	seqID := SeqId(0)
	nEmbeddings := int32(VocabNTokens(vocab))
	embeddings, err := GetEmbeddingsSeq(ctx, seqID, nEmbeddings)
	if err != nil {
		t.Fatalf("GetEmbeddingsSeq returned an error: %v", err)
	}
	if embeddings == nil {
		t.Fatal("GetEmbeddingsSeq returned nil")
	}
	if len(embeddings) != int(nEmbeddings) {
		t.Fatalf("GetEmbeddingsSeq returned %d embeddings, expected %d", len(embeddings), nEmbeddings)
	}
	t.Logf("GetEmbeddingsSeq returned %d embeddings", len(embeddings))
}

func TestGetEmbeddings(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	params := ContextDefaultParams()
	params.PoolingType = PoolingTypeNone
	params.Embeddings = 1

	ctx, err := InitFromModel(model, params)
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	// Tokenize a prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, prompt, true, true)

	// Create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	nOutputs := len(tokens)
	nEmbeddings := int(ModelNEmbd(model))
	embeddings, err := GetEmbeddings(ctx, nOutputs, nEmbeddings)
	if err != nil {
		t.Fatalf("GetEmbeddings returned an error: %v", err)
	}
	if embeddings == nil {
		t.Fatal("GetEmbeddings returned nil")
	}
	if len(embeddings) != nOutputs*nEmbeddings {
		t.Fatalf("GetEmbeddings returned %d values, expected %d", len(embeddings), nOutputs*nEmbeddings)
	}
	t.Logf("GetEmbeddings returned %d values", len(embeddings))
}

func TestSetEmbeddings(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	// Enable causal attention
	SetCausalAttn(ctx, true)
	t.Log("SetCausalAttn successfully set causal attention to true")

	// Disable causal attention
	SetCausalAttn(ctx, false)
	t.Log("SetCausalAttn successfully set causal attention to false")
}

func TestGetLogits(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	// Tokenize a prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, prompt, true, true)

	// Create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	nTokens := len(tokens)
	nVocab := int(VocabNTokens(vocab))
	logits, err := GetLogits(ctx, nTokens, nVocab)
	if err != nil {
		t.Fatalf("GetLogits returned an error: %v", err)
	}
	if logits == nil {
		t.Fatal("GetLogits returned nil")
	}
	if len(logits) != nTokens*nVocab {
		t.Fatalf("GetLogits returned %d values, expected %d", len(logits), nTokens*nVocab)
	}
	t.Logf("GetLogits returned %d values", len(logits))
}

func TestEncode(t *testing.T) {
	modelFile := testEncoderModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	if !ModelHasEncoder(model) {
		t.Skip("Model does not have an encoder; skipping Encode test")
	}

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	// Tokenize a prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, prompt, true, true)

	// Create batch
	batch := BatchGetOne(tokens)

	// Call Encode
	result, err := Encode(ctx, batch)
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}
	if result != 0 {
		t.Fatalf("Encode returned non-zero result: %d", result)
	}
	t.Logf("Encode succeeded with result: %d", result)
}

func TestSynchronize(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	// Call Synchronize and ensure no error occurs
	if err := Synchronize(ctx); err != nil {
		t.Fatalf("Synchronize returned error: %v", err)
	}
	t.Log("Synchronize completed successfully")
}

func TestGetPoolingType(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	poolingType := GetPoolingType(ctx)
	t.Logf("GetPoolingType returned: %d", poolingType)
	// Optionally, check for valid enum values if known
}
