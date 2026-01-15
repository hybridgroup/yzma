package mtmd

import (
	"sync/atomic"
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestNewEvalCallbackSimple(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	var callCount atomic.Int32

	cb := NewEvalCallbackSimple(func(tensor uintptr, ask bool) bool {
		callCount.Add(1)
		return true
	})

	if cb == 0 {
		t.Fatal("NewEvalCallbackSimple returned a nil callback")
	}

	t.Logf("NewEvalCallbackSimple created callback at %#x", cb)
}

func TestContextParamsWithCallback(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	var evalCount atomic.Int32

	params := ContextParamsDefault()
	params.SetEvalCallbackSimple(func(tensor uintptr, ask bool) bool {
		evalCount.Add(1)
		return true
	})

	if params.CbEval == 0 {
		t.Fatal("CbEval was not set on context params")
	}

	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	t.Logf("Context initialized with eval callback, callback ptr: %#x", params.CbEval)
}

func TestNewEvalCallback(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	var callCount atomic.Int32

	cb := NewEvalCallback(func(tensor uintptr, ask bool, userData uintptr) bool {
		callCount.Add(1)
		return true
	})

	if cb == 0 {
		t.Fatal("NewEvalCallback returned a nil callback")
	}

	t.Logf("NewEvalCallback created callback at %#x", cb)
}
