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

func TestModelInvalidFile(t *testing.T) {
	modelFile := "invalid_model.gguf"

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	if model != 0 {
		t.Fatal("ModelLoadFromFile should have failed for invalid file")
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

func TestModelMetaCount(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	count := ModelMetaCount(model)
	t.Logf("ModelMetaCount returned: %d", count)
	if count < 0 {
		t.Fatal("ModelMetaCount returned negative value")
	}
}

func TestModelMetaKeyByIndex(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	count := ModelMetaCount(model)
	if count <= 0 {
		t.Skip("No metadata keys to test")
	}
	key, ok := ModelMetaKeyByIndex(model, 0)
	if !ok {
		t.Fatal("ModelMetaKeyByIndex failed for index 0")
	}
	t.Logf("ModelMetaKeyByIndex returned: %s", key)
}

func TestModelMetaValStrByIndex(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	count := ModelMetaCount(model)
	if count <= 0 {
		t.Skip("No metadata values to test")
	}
	val, ok := ModelMetaValStrByIndex(model, 0)
	if !ok {
		t.Fatal("ModelMetaValStrByIndex failed for index 0")
	}
	t.Logf("ModelMetaValStrByIndex returned: %s", val)
}

func TestModelMetaValStr(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	count := ModelMetaCount(model)
	if count <= 0 {
		t.Skip("No metadata to test")
	}
	key, ok := ModelMetaKeyByIndex(model, 0)
	if !ok {
		t.Skip("ModelMetaKeyByIndex failed for index 0")
	}
	val, ok := ModelMetaValStr(model, key)
	if !ok {
		t.Fatal("ModelMetaValStr failed for key:", key)
	}
	t.Logf("ModelMetaValStr returned: %s", val)
}
