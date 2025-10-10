package llama

import (
	"testing"
)

func TestModelDefaultParams(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	if params == (ModelParams{}) {
		t.Fatal("ModelDefaultParams returned empty parameters")
	}
}

func TestModelHasDecoder(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	hasDecoder := ModelHasDecoder(model)
	if !hasDecoder {
		t.Fatal("ModelHasDecoder returned false, but the model should have a decoder")
	}
}

func TestModelNCtxTrain(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	nCtxTrain := ModelNCtxTrain(model)
	if nCtxTrain <= 0 {
		t.Fatal("ModelNCtxTrain returned an invalid value")
	}
}
