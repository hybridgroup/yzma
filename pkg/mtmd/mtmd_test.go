package mtmd

import (
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestDefaultMarker(t *testing.T) {
	marker := DefaultMarker()
	if marker == "" {
		t.Fatal("DefaultMarker returned an empty string")
	}
	t.Logf("DefaultMarker returned: %s", marker)
}

func TestContextParamsDefault(t *testing.T) {
	params := ContextParamsDefault()
	if params.Threads <= 0 {
		t.Fatal("ContextParamsDefault returned invalid thread count")
	}
	if params.ProgressCallback != 0 {
		t.Fatal("ContextParamsDefault returned non-nil ProgressCallback")
	}
	if params.ProgressCallbackUserData != 0 {
		t.Fatal("ContextParamsDefault returned non-nil ProgressCallbackUserData")
	}
	t.Logf("ContextParamsDefault returned: %+v", params)
}

func TestSetProgressCallback(t *testing.T) {
	var params ContextParamsType

	// nil callback should clear the pointer
	params.SetProgressCallback(nil)
	if params.ProgressCallback != 0 {
		t.Fatal("SetProgressCallback(nil) did not clear ProgressCallback")
	}

	// non-nil callback should set a non-zero pointer
	params.SetProgressCallback(func(progress float32, userData uintptr) bool {
		return true
	})
	if params.ProgressCallback == 0 {
		t.Fatal("SetProgressCallback did not set ProgressCallback")
	}
}

func TestInitFromFileAndFree(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	if ctx == 0 {
		t.Fatal("InitFromFile returned an invalid context")
	}

	t.Log("InitFromFile successfully initialized the context")

	Free(ctx)
	t.Log("Free successfully freed the context")
}

func TestSupportVision(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	supportsVision := SupportVision(ctx)
	if !supportsVision {
		t.Fatal("SupportVision expected true for multimodal model")
	}
}

func TestTokenize(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	text := NewInputText("Here is an image: <__media__>", true, true)

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open image file")
	}

	bitmap := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap)

	bitmaps := []Bitmap{bitmap} // Replace with actual bitmap data if available

	result := Tokenize(ctx, chunks, text, bitmaps)
	if result != 0 {
		t.Fatalf("Tokenize failed with result: %d", result)
	}

	t.Log("Tokenize successfully tokenized the input text")
}

func TestHelperEvalChunks(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	lctx, err := llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer llama.Free(lctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	var nPast llama.Pos = 0
	var seqID llama.SeqId = 1
	var nBatch int32 = 1
	var logitsLast bool = true
	var newNPast llama.Pos

	result := HelperEvalChunks(ctx, lctx, chunks, nPast, seqID, nBatch, logitsLast, &newNPast)
	if result != 0 {
		t.Fatalf("HelperEvalChunks failed with result: %d", result)
	}

	t.Log("HelperEvalChunks successfully evaluated the chunks")
}

func TestEncodeChunk(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)
	idx := uint64(1)
	chunk := InputChunksGet(chunks, idx)

	err = EncodeChunk(ctx, chunk)
	if err != nil {
		t.Fatalf("EncodeChunk failed: %v", err)
	}

	t.Log("EncodeChunk successfully encoded the chunk")
}

func TestGetOutputEmbd(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)
	idx := uint64(1)
	chunk := InputChunksGet(chunks, idx)

	err = EncodeChunk(ctx, chunk)
	if err != nil {
		t.Fatalf("EncodeChunk failed: %v", err)
	}

	sz := llama.ModelNEmbdInp(model) * int32(InputChunkGetNTokens(chunk))
	if sz <= 0 {
		t.Fatal("Calculated embedding size is invalid")
	}
	embd, err := GetOutputEmbd(ctx, sz)
	if err != nil {
		t.Fatalf("GetOutputEmbd failed: %v", err)
	}
	if embd == nil {
		t.Fatal("GetOutputEmbd returned nil")
	}

	t.Logf("GetOutputEmbd successfully retrieved embeddings of length: %d", len(embd))
}

func TestDecodeUseNonCausal(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	// Test with nil chunk (default image chunk behavior)
	useNonCausal := DecodeUseNonCausal(ctx, 0)
	t.Logf("DecodeUseNonCausal with nil chunk returned: %v", useNonCausal)

	// Test with zero context returns false
	if DecodeUseNonCausal(0, 0) {
		t.Fatal("DecodeUseNonCausal expected false for zero context")
	}

	// Test with an actual image chunk
	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)
	for i := uint64(0); i < InputChunksSize(chunks); i++ {
		chunk := InputChunksGet(chunks, i)
		result := DecodeUseNonCausal(ctx, chunk)
		chunkType := InputChunkGetType(chunk)
		t.Logf("DecodeUseNonCausal with chunk[%d] (type=%d) returned: %v", i, chunkType, result)

		// For image chunks, the result should match the nil-chunk default
		if chunkType == InputChunkTypeImage && result != useNonCausal {
			t.Fatalf("DecodeUseNonCausal mismatch: nil chunk returned %v, image chunk returned %v", useNonCausal, result)
		}
	}
}

func TestDecodeUseMRope(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	useMRope := DecodeUseMRope(ctx)
	t.Logf("DecodeUseMRope returned: %v", useMRope)

	// Test with zero context returns false
	if DecodeUseMRope(0) {
		t.Fatal("DecodeUseMRope expected false for zero context")
	}
}

func TestSupportAudio(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	supportsAudio := SupportAudio(ctx)
	t.Logf("SupportAudio returned: %v", supportsAudio)

	// Test with zero context returns false
	if SupportAudio(0) {
		t.Fatal("SupportAudio expected false for zero context")
	}
}

func TestGetAudioSampleRate(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	sampleRate := GetAudioSampleRate(ctx)
	t.Logf("GetAudioSampleRate returned: %d", sampleRate)
	if sampleRate != -1 && sampleRate <= 0 {
		t.Fatalf("GetAudioSampleRate returned an invalid sample rate: %d", sampleRate)
	}

	// Test with zero context returns -1
	if GetAudioSampleRate(0) != -1 {
		t.Fatal("GetAudioSampleRate expected -1 for zero context")
	}
}
