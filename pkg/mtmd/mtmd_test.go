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
	t.Logf("ContextParamsDefault returned: %+v", params)
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
	t.Logf("SupportVision returned: %v", supportsVision)
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

	useNonCausal := DecodeUseNonCausal(ctx)
	t.Logf("DecodeUseNonCausal returned: %v", useNonCausal)
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
}

func TestGetAudioBitrate(t *testing.T) {
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

	bitrate := GetAudioBitrate(ctx)
	t.Logf("GetAudioBitrate returned: %d", bitrate)
	if bitrate != -1 && bitrate <= 0 {
		t.Fatal("GetAudioBitrate returned an invalid bitrate")
	}
}
