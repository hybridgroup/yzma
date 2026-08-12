package llama

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"
)

func TestModelDefaultParams(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	if params == (ModelParams{}) {
		t.Fatal("ModelDefaultParams returned empty parameters")
	}
}

// TestModelDefaultParamsLayout guards against llama_model_params layout drift.
// ffiTypeModelParams is hand-maintained, so if llama.cpp adds, removes, or reorders
// a field the native defaults land in the wrong Go fields. Checking the values that
// llama.cpp is known to write catches that.
func TestModelDefaultParamsLayout(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()

	// The C bools. A field outside {0, 1} means libffi wrote struct padding
	// into it, which is what a missing or misplaced descriptor entry looks like.
	bools := []struct {
		name  string
		value uint8
	}{
		{"VocabOnly", params.VocabOnly},
		{"CheckTensors", params.CheckTensors},
		{"UseExtraBufts", params.UseExtraBufts},
		{"NoHost", params.NoHost},
		{"NoAlloc", params.NoAlloc},
		{"LoadMTP", params.LoadMTP},
	}
	for _, b := range bools {
		if b.value > 1 {
			t.Errorf("%s is %d, want a 0 or 1 C bool", b.name, b.value)
		}
	}

	// llama.cpp defaults every pointer field to NULL.
	if params.Devices != 0 {
		t.Errorf("Devices is %#x, want 0", params.Devices)
	}
	if params.TensorBuftOverrides != 0 {
		t.Errorf("TensorBuftOverrides is %#x, want 0", params.TensorBuftOverrides)
	}
	if params.TensorSplit != nil {
		t.Error("TensorSplit is not nil, want nil")
	}
	if params.ProgressCallback != 0 {
		t.Errorf("ProgressCallback is %#x, want 0", params.ProgressCallback)
	}
	if params.ProgressCallbackUserData != 0 {
		t.Errorf("ProgressCallbackUserData is %#x, want 0", params.ProgressCallbackUserData)
	}
	if params.KvOverrides != 0 {
		t.Errorf("KvOverrides is %#x, want 0", params.KvOverrides)
	}

	// The int32 fields. Ranges rather than exact values, so an upstream default
	// tweak does not fail the test but a shifted field does.
	// -1 means all layers, which is the current llama.cpp default.
	if params.NGpuLayers < -1 {
		t.Errorf("NGpuLayers is %d, want >= -1", params.NGpuLayers)
	}
	if params.SplitMode < SplitModeNone || params.SplitMode > SplitModeRow {
		t.Errorf("SplitMode is %d, want a valid SplitMode", params.SplitMode)
	}
	if params.LoadMode < LoadModeNone || params.LoadMode > LoadModeDirectIO {
		t.Errorf("LoadMode is %d, want a valid LoadMode", params.LoadMode)
	}
	if params.MainGpu < 0 {
		t.Errorf("MainGpu is %d, want >= 0", params.MainGpu)
	}
}

// TestModelLoadFromFileLoadMTP checks that a caller-set LoadMTP survives the trip
// through libffi. The test model has no MTP layers, so llama.cpp has nothing extra
// to load and the request is a no-op, but a bad struct layout would corrupt the
// neighboring fields and break the load.
func TestModelLoadFromFileLoadMTP(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	params.LoadMTP = 1

	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile with LoadMTP failed: %v", err)
	}
	defer ModelFree(model)

	if params.LoadMTP != 1 {
		t.Errorf("LoadMTP is %d after loading, want 1", params.LoadMTP)
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

func TestModelNLayerNextN(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nLayerNextN := ModelNLayerNextN(model)
	t.Logf("ModelNLayerNextN returned: %d", nLayerNextN)
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

// TestMetaStr drives the buffer growth directly, standing in for a model whose
// metadata is large enough to overflow the initial buffer. The stub mimics
// snprintf: it writes what fits, always NUL-terminates, and returns the length
// it would have written.
func TestMetaStr(t *testing.T) {
	snprintf := func(value string) func(b *byte, bLen uint64) int32 {
		return func(b *byte, bLen uint64) int32 {
			buf := unsafe.Slice(b, bLen)
			n := copy(buf[:len(buf)-1], value)
			buf[n] = 0
			return int32(len(value))
		}
	}

	tests := []struct {
		name  string
		size  int
		value string
	}{
		{"fits", 128, "llama 135M Q2_K"},
		{"exactly one byte short", 16, "123456789012345"},
		{"needs the whole buffer", 16, "1234567890123456"},
		{"grows once", 8, strings.Repeat("x", 100)},
		{"grows a long way", 4, strings.Repeat("chat template ", 5000)},
		{"empty", 8, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := metaStr(tt.size, snprintf(tt.value))
			if !ok {
				t.Fatalf("metaStr(%d) failed for a %d byte value", tt.size, len(tt.value))
			}
			if got != tt.value {
				t.Errorf("metaStr(%d) returned %d bytes, want %d", tt.size, len(got), len(tt.value))
			}
		})
	}

	if _, ok := metaStr(128, func(*byte, uint64) int32 { return -1 }); ok {
		t.Error("metaStr reported success for a negative result")
	}
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
	if len(desc) == 0 {
		t.Fatal("ModelDesc returned an empty string")
	}
	t.Logf("ModelDesc returned: %s", desc)
}

func TestModelFtype(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ftype := ModelFtype(model)
	if ftype == FtypeGUESSED {
		t.Fatal("ModelFtype returned FtypeGUESSED, which is invalid")
	}
	t.Logf("ModelFtype returned: %d", ftype)
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
	if size == 0 {
		t.Fatal("ModelSize returned 0, which is invalid")
	}
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

	// A float return read through an ffi.Arg yields ~1.2e19 rather than a
	// scaling factor, so pin the magnitude and not just the call. llama.cpp
	// stores 1.0f/ropescale, which exceeds 1 for a GGUF with a scaling
	// factor below 1, so the ceiling only has to exclude the garbage value.
	if freqScale <= 0 || freqScale > 1e6 {
		t.Fatalf("ModelRopeFreqScaleTrain returned %g, want a plausible scaling factor", freqScale)
	}
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

	// Enumerate every entry, checking that each round trip is well formed and
	// that no value comes back carrying its NUL terminator.
	//
	// This does not exercise the grow-and-retry path: reaching it needs a value
	// over 32 KiB, which the small test model has none of, so these assertions
	// hold against a truncating implementation too. TestMetaStr drives that path
	// directly and is what pins it.
	for i := range count {
		key, ok := ModelMetaKeyByIndex(model, i)
		if !ok {
			t.Errorf("ModelMetaKeyByIndex failed for index %d", i)
			continue
		}
		val, ok := ModelMetaValStrByIndex(model, i)
		if !ok {
			t.Errorf("ModelMetaValStrByIndex failed for index %d", i)
			continue
		}
		if strings.ContainsRune(key, 0) || strings.ContainsRune(val, 0) {
			t.Errorf("metadata entry %d (%q) contains a NUL terminator", i, key)
		}
	}
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

func TestModelParamsSetTensorBufOverrides(t *testing.T) {
	t.Run("empty_clears_pointer", func(t *testing.T) {
		p := ModelDefaultParams()
		p.TensorBuftOverrides = 0xDEAD
		if err := p.SetTensorBufOverrides(nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.TensorBuftOverrides != 0 {
			t.Fatalf("TensorBuftOverrides = %#x, want 0", p.TensorBuftOverrides)
		}
	})

	t.Run("valid_sentinel_sets_pointer", func(t *testing.T) {
		b := byte('x')
		overrides := []TensorBuftOverride{{Pattern: &b}, {}}
		p := ModelDefaultParams()
		if err := p.SetTensorBufOverrides(overrides); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := uintptr(unsafe.Pointer(&overrides[0]))
		if p.TensorBuftOverrides != want {
			t.Fatalf("TensorBuftOverrides = %#x, want %#x", p.TensorBuftOverrides, want)
		}
	})

	t.Run("missing_sentinel_returns_error", func(t *testing.T) {
		b := byte('x')
		overrides := []TensorBuftOverride{{Pattern: &b}}
		p := ModelDefaultParams()
		p.TensorBuftOverrides = 0xDEAD
		if err := p.SetTensorBufOverrides(overrides); err == nil {
			t.Fatal("expected error for missing sentinel")
		}
		if p.TensorBuftOverrides != 0xDEAD {
			t.Fatal("TensorBuftOverrides was modified on error")
		}
	})
}

func TestModelParamsSetDevices(t *testing.T) {
	t.Run("empty_clears_pointer", func(t *testing.T) {
		p := ModelDefaultParams()
		p.Devices = 0xDEAD
		if err := p.SetDevices(nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Devices != 0 {
			t.Fatalf("Devices = %#x, want 0", p.Devices)
		}
	})

	t.Run("valid_null_terminated_sets_pointer", func(t *testing.T) {
		devices := []GGMLBackendDevice{GGMLBackendDevice(1), 0}
		p := ModelDefaultParams()
		if err := p.SetDevices(devices); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := uintptr(unsafe.Pointer(&devices[0]))
		if p.Devices != want {
			t.Fatalf("Devices = %#x, want %#x", p.Devices, want)
		}
	})

	t.Run("missing_terminator_returns_error", func(t *testing.T) {
		devices := []GGMLBackendDevice{GGMLBackendDevice(1)}
		p := ModelDefaultParams()
		p.Devices = 0xDEAD
		if err := p.SetDevices(devices); err == nil {
			t.Fatal("expected error for missing NULL terminator")
		}
		if p.Devices != 0xDEAD {
			t.Fatal("Devices was modified on error")
		}
	})
}

func TestModelNParams(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	nParams := ModelNParams(model)
	if nParams == 0 {
		t.Fatal("ModelNParams returned 0")
	}
	t.Logf("ModelNParams returned: %d", nParams)
}

func TestModelSaveToFile(t *testing.T) {
	t.Skip("skipped: crashes on Linux CI due to C-level backend state after repeated init/free cycles")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "saved_model.gguf")
	ModelSaveToFile(model, outPath)

	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("ModelSaveToFile did not create file: %v", err)
	}
	t.Log("ModelSaveToFile completed successfully")
}

func TestSplitPath(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	path, length := SplitPath("/models/ggml-model-q4_0", 2, 4)
	if length < 0 {
		t.Fatal("SplitPath returned negative length")
	}
	expected := "/models/ggml-model-q4_0-00003-of-00004.gguf"
	if path != expected {
		t.Fatalf("SplitPath = %q, want %q", path, expected)
	}
	t.Logf("SplitPath returned: %s", path)
}

func TestSplitPrefix(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	prefix, length := SplitPrefix("/models/ggml-model-q4_0-00003-of-00004.gguf", 2, 4)
	if length < 0 {
		t.Fatal("SplitPrefix returned negative length")
	}
	expected := "/models/ggml-model-q4_0"
	if prefix != expected {
		t.Fatalf("SplitPrefix = %q, want %q", prefix, expected)
	}
	t.Logf("SplitPrefix returned: %s", prefix)
}
