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
