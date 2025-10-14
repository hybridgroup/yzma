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
