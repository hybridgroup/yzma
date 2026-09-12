package mtmd

import (
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestInputChunksInitAndFree(t *testing.T) {
	testSetup(t)
	chunks := InputChunksInit()
	if chunks == InputChunks(0) {
		t.Fatal("InputChunksInit returned an invalid InputChunks")
	}

	t.Log("InputChunksInit successfully initialized InputChunks")

	InputChunksFree(chunks)
	t.Log("InputChunksFree successfully freed InputChunks")
}

func TestInputChunksSize(t *testing.T) {
	testSetup(t)
	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	size := InputChunksSize(chunks)
	if size != 0 {
		t.Fatalf("InputChunksSize returned a non-zero size for an empty InputChunks: %d", size)
	}

	t.Logf("InputChunksSize returned: %d", size)
}

func TestInputChunksGetType(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	size := InputChunksSize(chunks)
	t.Logf("InputChunksSize returned: %d", size)
	if size == 0 {
		t.Fatalf("invalid chunk size: %d", size)
	}

	idx := uint64(size - 1) // Use the last chunk index to ensure we are testing a valid chunk
	chunk := InputChunksGet(chunks, idx)
	if chunk == InputChunk(0) {
		t.Fatalf("InputChunksGet returned an invalid chunk for index %d", idx)
	}

	t.Logf("InputChunksGet successfully retrieved chunk at index %d", idx)

	chunkType := InputChunkGetType(chunk)
	if chunkType != InputChunkTypeText && chunkType != InputChunkTypeImage {
		t.Fatalf("InputChunkGetType returned an unexpected type: %d", chunkType)
	}
}

func TestInputChunkGetTokensText(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	size := InputChunksSize(chunks)
	t.Logf("InputChunksSize returned: %d", size)
	if size == 0 {
		t.Fatalf("invalid chunk size: %d", size)
	}

	// find text chunk
	var textChunk InputChunk
	for i := uint64(0); i < size; i++ {
		chunk := InputChunksGet(chunks, i)
		if chunk == InputChunk(0) {
			t.Fatalf("InputChunksGet returned an invalid chunk for index %d", i)
		}

		chunkType := InputChunkGetType(chunk)
		if chunkType == InputChunkTypeText {
			textChunk = chunk
			break
		}
	}

	if textChunk == InputChunk(0) {
		t.Fatal("No text chunk found in InputChunks")
	}
	tokens := InputChunkGetTokensText(textChunk)
	if tokens == nil {
		t.Fatal("InputChunkGetTokensText returned nil")
	}
}

func TestInputChunkGetNTokens(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint64(1)
	chunk := InputChunksGet(chunks, idx)
	if chunk == InputChunk(0) {
		t.Fatalf("InputChunksGet returned an invalid chunk for index %d", idx)
	}

	t.Logf("InputChunksGet successfully retrieved chunk at index %d", idx)

	nTokens := InputChunkGetNTokens(chunk)
	t.Logf("InputChunkGetNTokens returned: %d", nTokens)
}

func TestInputChunkGetId(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint64(1)
	chunk := InputChunksGet(chunks, idx)
	if chunk == InputChunk(0) {
		t.Fatalf("InputChunksGet returned an invalid chunk for index %d", idx)
	}

	t.Logf("InputChunksGet successfully retrieved chunk at index %d", idx)

	id := InputChunkGetId(chunk)
	t.Logf("InputChunkGetId returned: %s", id)
}

func TestInputChunkGetNPos(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint64(1)
	chunk := InputChunksGet(chunks, idx)
	if chunk == InputChunk(0) {
		t.Fatalf("InputChunksGet returned an invalid chunk for index %d", idx)
	}

	t.Logf("InputChunksGet successfully retrieved chunk at index %d", idx)

	nPos := InputChunkGetNPos(chunk)
	t.Logf("InputChunkGetNPos returned: %d", nPos)
}

func TestInputChunkCopyAndFree(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint64(1)
	chunk := InputChunksGet(chunks, idx)
	if chunk == InputChunk(0) {
		t.Fatalf("InputChunksGet returned an invalid chunk for index %d", idx)
	}

	t.Logf("InputChunksGet successfully retrieved chunk at index %d", idx)

	copy := InputChunkCopy(chunk)
	if copy == InputChunk(0) {
		t.Fatal("InputChunkCopy returned an invalid chunk")
	}

	t.Log("InputChunkCopy successfully created a copy of the chunk")

	InputChunkFree(copy)
	t.Log("InputChunkFree successfully freed the copied chunk")
}

func TestInputChunkGetTokensImage(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint64(1)
	chunk := InputChunksGet(chunks, idx)
	if chunk == InputChunk(0) {
		t.Fatalf("InputChunksGet returned an invalid chunk for index %d", idx)
	}

	t.Logf("InputChunksGet successfully retrieved chunk at index %d", idx)

	tokens := InputChunkGetTokensImage(chunk)
	if tokens == ImageTokens(0) {
		t.Fatalf("InputChunkGetTokensImage returned a nil pointer")
	}

	t.Logf("InputChunkGetTokensImage returned image tokens")

	nTokens := ImageTokensGetNTokens(tokens)
	nx := ImageTokensGetNX(tokens)
	ny := ImageTokensGetNY(tokens)
	id := ImageTokensGetId(tokens)
	nPos := ImageTokensGetNPos(tokens)

	t.Logf("n_tokens: %d", nTokens)
	t.Logf("nx:      %d", nx)
	t.Logf("ny:      %d", ny)
	t.Logf("id:      %s", id)
	t.Logf("n_pos:   %d", nPos)
}

// TestDecoderPosLayout keeps the Go struct at the 16 byte layout of the C struct
// mtmd_decoder_pos. A reordered or resized field makes the FFI return buffer wrong.
func TestDecoderPosLayout(t *testing.T) {
	var pos DecoderPos

	if got := unsafe.Sizeof(pos); got != 16 {
		t.Fatalf("unsafe.Sizeof(DecoderPos) = %d, want 16", got)
	}

	offsets := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"T", unsafe.Offsetof(pos.T), 0},
		{"X", unsafe.Offsetof(pos.X), 4},
		{"Y", unsafe.Offsetof(pos.Y), 8},
		{"Z", unsafe.Offsetof(pos.Z), 12},
	}

	for _, o := range offsets {
		if o.got != o.want {
			t.Errorf("offset of %s = %d, want %d", o.name, o.got, o.want)
		}
	}
}

func TestImageTokensGetDecoderPos(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	if got := ImageTokensGetDecoderPos(ImageTokens(0), 0, 0); got != (DecoderPos{}) {
		t.Fatalf("ImageTokensGetDecoderPos of nil image tokens = %+v, want the zero value", got)
	}

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx, err := InitFromFile(mmprojFile, model, params)
	if err != nil {
		t.Fatalf("Failed to initialize context from file: %v", err)
	}
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	chunk := InputChunksGet(chunks, 1)
	if chunk == InputChunk(0) {
		t.Fatal("InputChunksGet returned an invalid chunk for index 1")
	}

	tokens := InputChunkGetTokensImage(chunk)
	if tokens == ImageTokens(0) {
		t.Fatal("InputChunkGetTokensImage returned a nil pointer")
	}

	nTokens := ImageTokensGetNTokens(tokens)
	if nTokens == 0 {
		t.Fatal("ImageTokensGetNTokens returned 0")
	}

	// The deprecated grid gives the values that the M-RoPE layout must agree with.
	nx := ImageTokensGetNX(tokens)
	ny := ImageTokensGetNY(tokens)
	mrope := DecodeUseMRope(ctx)

	t.Logf("n_tokens: %d nx: %d ny: %d mrope: %v", nTokens, nx, ny, mrope)

	// A position other than 0 shows that the call uses pos_0.
	pos0 := llama.Pos(7)

	for i := uint64(0); i < nTokens; i++ {
		got := ImageTokensGetDecoderPos(tokens, pos0, i)

		var want DecoderPos
		if mrope {
			// M-RoPE keeps one temporal position. X is the column and Y is the row,
			// so a swap of the two planes makes this test fail.
			want = DecoderPos{
				T: uint32(pos0),
				X: uint32(pos0) + uint32(i%nx),
				Y: uint32(pos0) + uint32(i/nx),
				Z: 0,
			}
		} else {
			// Models without M-RoPE get the same sequential position in each plane.
			p := uint32(pos0) + uint32(i)
			want = DecoderPos{T: p, X: p, Y: p, Z: p}
		}

		if got != want {
			t.Fatalf("ImageTokensGetDecoderPos(%d) = %+v, want %+v", i, got, want)
		}

		if i < 4 {
			t.Logf("token %d: %+v", i, got)
		}
	}

	// An index out of range returns the zero value instead of an abort in C.
	if got := ImageTokensGetDecoderPos(tokens, pos0, nTokens); got != (DecoderPos{}) {
		t.Fatalf("ImageTokensGetDecoderPos of an out of range index = %+v, want the zero value", got)
	}
}
