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

func TestModelNClsOut(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	nClsOut := ModelNClsOut(model)
	t.Logf("ModelNClsOut returned: %d", nClsOut)
}

func TestModelClsLabel(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	label := ModelClsLabel(model, 0)
	t.Logf("ModelClsLabel returned: %s", label)
}

func TestModelDesc(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	desc := ModelDesc(model)
	t.Logf("ModelDesc returned: %s", desc)
}

func TestModelSize(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	size := ModelSize(model)
	t.Logf("ModelSize returned: %d bytes", size)
}

func TestModelIsRecurrent(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	isRecurrent := ModelIsRecurrent(model)
	t.Logf("ModelIsRecurrent returned: %v", isRecurrent)
}

func TestModelIsHybrid(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	isHybrid := ModelIsHybrid(model)
	t.Logf("ModelIsHybrid returned: %v", isHybrid)
}

func TestModelIsDiffusion(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	isDiffusion := ModelIsDiffusion(model)
	t.Logf("ModelIsDiffusion returned: %v", isDiffusion)
}

func TestModelRopeFreqScaleTrain(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(testModelFileName(t), ModelDefaultParams())
	defer ModelFree(model)

	freqScale := ModelRopeFreqScaleTrain(model)
	t.Logf("ModelRopeFreqScaleTrain returned: %f", freqScale)
}

func TestModelRopeType(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(testModelFileName(t), ModelDefaultParams())
	defer ModelFree(model)

	ropeType := ModelRopeType(model)
	t.Logf("ModelRopeType returned: %d", ropeType)
}
