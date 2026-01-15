package llama

import (
	"os"
	"path/filepath"
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
	model, err := ModelLoadFromFile(modelFile, params)
	if err == nil {
		t.Fatal("ModelLoadFromFile should have failed for invalid file")
	}
	if model != 0 {
		t.Fatal("ModelLoadFromFile should have failed for invalid file")
	}
}

func TestModelHasDecoder(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	hasDecoder := ModelHasDecoder(model)
	if !hasDecoder {
		t.Fatal("ModelHasDecoder returned false, but the model should have a decoder")
	}
}

func TestModelNEmbdInp(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nEmbdInp := ModelNEmbdInp(model)
	if nEmbdInp <= 0 {
		t.Fatal("ModelNEmbdInp returned an invalid value")
	}
}

func TestModelNEmbdOut(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nEmbdOut := ModelNEmbdOut(model)
	if nEmbdOut <= 0 {
		t.Fatal("ModelNEmbdOut returned an invalid value")
	}
	t.Logf("ModelNEmbdOut returned: %d", nEmbdOut)
}

func TestModelNLayer(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nLayer := ModelNLayer(model)
	if nLayer <= 0 {
		t.Fatal("ModelNLayer returned an invalid value")
	}
}

func TestModelNHead(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)
	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nHead := ModelNHead(model)
	if nHead <= 0 {
		t.Fatal("ModelNHead returned an invalid value")
	}
}

func TestModelNHeadKV(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nHeadKV := ModelNHeadKV(model)
	if nHeadKV <= 0 {
		t.Fatal("ModelNHeadKV returned an invalid value")
	}
}

func TestModelNSWA(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nSWA := ModelNSWA(model)
	if nSWA < 0 {
		t.Fatal("ModelNSWA returned an invalid value")
	}
}

func TestModelNCtxTrain(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nClsOut := ModelNClsOut(model)
	t.Logf("ModelNClsOut returned: %d", nClsOut)
}

func TestModelClsLabel(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	label := ModelClsLabel(model, 0)
	t.Logf("ModelClsLabel returned: %s", label)
}

func TestModelDesc(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	desc := ModelDesc(model)
	t.Logf("ModelDesc returned: %s", desc)
}

func TestModelSize(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	size := ModelSize(model)
	t.Logf("ModelSize returned: %d bytes", size)
}

func TestModelIsRecurrent(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	isRecurrent := ModelIsRecurrent(model)
	t.Logf("ModelIsRecurrent returned: %v", isRecurrent)
}

func TestModelIsHybrid(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	isHybrid := ModelIsHybrid(model)
	t.Logf("ModelIsHybrid returned: %v", isHybrid)
}

func TestModelIsDiffusion(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	isDiffusion := ModelIsDiffusion(model)
	t.Logf("ModelIsDiffusion returned: %v", isDiffusion)
}

func TestModelRopeFreqScaleTrain(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(testModelFileName(t), ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	freqScale := ModelRopeFreqScaleTrain(model)
	t.Logf("ModelRopeFreqScaleTrain returned: %f", freqScale)
}

func TestModelRopeType(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(testModelFileName(t), ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ropeType := ModelRopeType(model)
	t.Logf("ModelRopeType returned: %d", ropeType)
}

func TestModelMetaCount(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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

func TestModelMetaKeyStr(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// Try a few likely valid and invalid keys
	invalidKey := ModelMetaKey(-12345)

	s := ModelMetaKeyStr(ModelMetaKeySamplingTopK)
	if s == "" {
		t.Log("ModelMetaKeyStr returned empty string for valid key (may be expected if no keys defined at 0)")
	} else {
		t.Logf("ModelMetaKeyStr(%d) returned: %q", ModelMetaKeySamplingTopK, s)
	}

	s = ModelMetaKeyStr(invalidKey)
	if s != "" {
		t.Fatalf("ModelMetaKeyStr should return empty string for invalid key, got: %q", s)
	}
}

func TestModelLoadCallback(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	progressCalls := 0
	callback := func(progress float32, userData uintptr) uint8 {
		progressCalls++
		t.Logf("Model loading progress: %.2f%%", progress*100)
		return 1 // continue loading
	}

	params := ModelDefaultParams()
	params.SetProgressCallback(callback)
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	if model == 0 {
		t.Fatal("ModelLoadFromFile failed to load model")
	}
	if progressCalls == 0 {
		t.Fatal("Progress callback was not called during model loading")
	}
}

func TestModelQuantizeDefaultParams(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelQuantizeDefaultParams()
	if params == (ModelQuantizeParams{}) {
		t.Fatal("ModelQuantizeDefaultParams returned empty parameters")
	}
}

func TestModelQuantize(t *testing.T) {
	modelFile := os.Getenv("YZMA_TEST_QUANTIZE_MODEL")
	if modelFile == "" {
		t.Skip("YZMA_TEST_QUANTIZE_MODEL env var not set; skipping TestModelQuantize")
	}

	tmpDir := t.TempDir()
	quantizedModelFile := filepath.Join(tmpDir, "quantized_model.gguf")

	testSetup(t)
	defer testCleanup(t)

	params := ModelQuantizeDefaultParams()
	params.NThread = 8
	params.Ftype = FtypeMostlyQ4_K_M
	result := ModelQuantize(modelFile, quantizedModelFile, &params)
	if result != 0 {
		t.Fatalf("ModelQuantize failed with error code: %d", result)
	}

	// Load the quantized model to verify it was created correctly
	quantizedModel, err := ModelLoadFromFile(quantizedModelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(quantizedModel)

	if quantizedModel == 0 {
		t.Fatal("Failed to load the quantized model")
	}
}

func TestModelChatTemplate(t *testing.T) {
	modelFile := testMMMModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	template := ModelChatTemplate(model, "")
	if template == "" {
		t.Fatal("ModelChatTemplate returned an empty string")
	}
	t.Logf("ModelChatTemplate returned: %s", template)
}

func TestModelLoadFromSplitsInvalid(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	paths := []string{"invalid_split1.gguf", "invalid_split2.gguf"}
	model, err := ModelLoadFromSplits(paths, params)
	if err == nil {
		t.Fatal("ModelLoadFromSplits should have failed for invalid files")
	}
	if model != 0 {
		t.Fatal("ModelLoadFromSplits should have failed for invalid files")
	}
}

func TestModelLoadFromSplitsValid(t *testing.T) {
	testSplitModelFileNames := testSplitModelFileNames(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromSplits(testSplitModelFileNames, params)
	if err != nil {
		t.Fatalf("ModelLoadFromSplits failed: %v", err)
	}
	defer ModelFree(model)

	if model == 0 {
		t.Fatal("ModelLoadFromSplits failed to load model")
	}
}
