package mtmd

import (
	"os"
	"runtime"
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

func TestGetMarker(t *testing.T) {
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

	marker := GetMarker(ctx)
	if marker == "" {
		t.Fatal("GetMarker returned an empty string")
	}
	if marker != DefaultMarker() {
		t.Fatalf("GetMarker returned %q, want %q", marker, DefaultMarker())
	}
	t.Logf("GetMarker returned: %s", marker)

	// Test with zero context returns empty string
	if GetMarker(0) != "" {
		t.Fatal("GetMarker expected empty string for zero context")
	}
}

func TestSupportVideo(t *testing.T) {
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

	// Simply log the result; video support depends on the build-time MTMD_VIDEO flag.
	supportsVideo := SupportVideo(ctx)
	t.Logf("SupportVideo returned: %v", supportsVideo)

	// Test with zero context returns false
	if SupportVideo(0) {
		t.Fatal("SupportVideo expected false for zero context")
	}
}

func TestVideoInitParamsDefault(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := VideoInitParamsDefault()
	if params.FPSTarget <= 0 {
		t.Fatalf("VideoInitParamsDefault returned non-positive FPSTarget: %v", params.FPSTarget)
	}
	if params.TimestampIntervalMs <= 0 {
		t.Fatalf("VideoInitParamsDefault returned non-positive TimestampIntervalMs: %v", params.TimestampIntervalMs)
	}
	t.Logf("VideoInitParamsDefault: fps_target=%.2f timestamp_interval_ms=%d", params.FPSTarget, params.TimestampIntervalMs)
}

func TestVideoInitAndFree(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)
	videoFile := testVideoFileName(t)

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

	if !SupportVideo(ctx) {
		t.Skip("video support not available in this build")
	}

	videoParams := VideoInitParamsDefault()
	videoCtx := VideoInit(ctx, videoFile, videoParams)
	if videoCtx == 0 {
		t.Fatal("VideoInit returned an invalid VideoContext (is ffprobe installed?)")
	}
	defer VideoFree(videoCtx)

	info := VideoGetInfo(videoCtx)
	if info.Width == 0 || info.Height == 0 {
		t.Fatalf("VideoGetInfo returned zero dimensions: %+v", info)
	}
	if info.FPS <= 0 {
		t.Fatalf("VideoGetInfo returned non-positive FPS: %v", info.FPS)
	}
	t.Logf("VideoGetInfo: %dx%d @ %.2f fps, ~%d frames", info.Width, info.Height, info.FPS, info.NFrames)
}

func TestVideoInitFromBuf(t *testing.T) {
	if runtime.GOOS == "windows" {
		// Pipe-based video probing deadlocks on Windows when ffprobe exits early
		// after reading enough data; the C++ feeder thread blocks indefinitely
		// instead of receiving a broken-pipe error. See llama.cpp issue #24429.
		t.Skip("VideoInitFromBuf pipe mode is not reliable on Windows")
	}
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)
	videoFile := testVideoFileName(t)

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

	if !SupportVideo(ctx) {
		t.Skip("video support not available in this build")
	}

	buf, err := os.ReadFile(videoFile)
	if err != nil {
		t.Fatalf("could not read video file: %v", err)
	}

	videoParams := VideoInitParamsDefault()
	videoCtx := VideoInitFromBuf(ctx, buf, videoParams)
	if videoCtx == 0 {
		t.Fatal("VideoInitFromBuf returned an invalid VideoContext (is ffprobe installed?)")
	}
	defer VideoFree(videoCtx)

	info := VideoGetInfo(videoCtx)
	if info.Width == 0 || info.Height == 0 {
		t.Fatalf("VideoGetInfo returned zero dimensions after VideoInitFromBuf: %+v", info)
	}
	t.Logf("VideoGetInfo (from buf): %dx%d @ %.2f fps, ~%d frames", info.Width, info.Height, info.FPS, info.NFrames)
}

func TestVideoFreeNilIsNoop(t *testing.T) {
	// VideoFree(0) must not panic or crash.
	VideoFree(0)
}

func TestVideoGetInfoZeroContext(t *testing.T) {
	// VideoGetInfo(0) must return a zero-value struct without crashing.
	info := VideoGetInfo(0)
	if info.Width != 0 || info.Height != 0 {
		t.Fatalf("VideoGetInfo(0) returned non-zero info: %+v", info)
	}
}

func TestVideoInitZeroContext(t *testing.T) {
	// VideoInit with a zero context must return 0 without crashing.
	videoCtx := VideoInit(0, "nonexistent.mp4", VideoInitParams{})
	if videoCtx != 0 {
		t.Fatal("VideoInit with zero context should return 0")
	}
}
