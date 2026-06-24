package mtmd

import (
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestBatchInitAndFree(t *testing.T) {
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

	batch, err := BatchInit(ctx)
	if err != nil {
		t.Fatalf("BatchInit failed: %v", err)
	}
	if batch == 0 {
		t.Fatal("BatchInit returned an invalid batch handle")
	}

	t.Log("BatchInit successfully created a batch")

	BatchFree(batch)
	t.Log("BatchFree successfully freed the batch")
}

func TestBatchInitInvalidContext(t *testing.T) {
	_, err := BatchInit(0)
	if err == nil {
		t.Fatal("BatchInit with zero context should return an error")
	}
}

func TestBatchFreeZeroBatch(t *testing.T) {
	// BatchFree with a zero handle should not panic
	BatchFree(0)
}

func TestBatchAddChunk(t *testing.T) {
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

	batch, err := BatchInit(ctx)
	if err != nil {
		t.Fatalf("BatchInit failed: %v", err)
	}
	defer BatchFree(batch)

	// Find the image chunk (index 1 is the media chunk from testSetupChunks)
	var imageChunk InputChunk
	for i := uint64(0); i < InputChunksSize(chunks); i++ {
		c := InputChunksGet(chunks, i)
		if InputChunkGetType(c) == InputChunkTypeImage {
			imageChunk = c
			break
		}
	}
	if imageChunk == 0 {
		t.Fatal("no image chunk found in tokenized output")
	}

	result := BatchAddChunk(batch, imageChunk)
	if result != BatchAddSuccess {
		t.Fatalf("BatchAddChunk failed with result: %d", result)
	}

	t.Log("BatchAddChunk successfully added the image chunk")
}

func TestBatchAddChunkInvalidHandles(t *testing.T) {
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

	batch, err := BatchInit(ctx)
	if err != nil {
		t.Fatalf("BatchInit failed: %v", err)
	}
	defer BatchFree(batch)

	// Zero chunk should return error
	if BatchAddChunk(batch, 0) != BatchAddError {
		t.Fatal("BatchAddChunk with zero chunk should return BatchAddError")
	}

	// Zero batch should return error
	chunks := InputChunksInit()
	defer InputChunksFree(chunks)
	testSetupChunks(t, ctx, chunks)

	for i := uint64(0); i < InputChunksSize(chunks); i++ {
		c := InputChunksGet(chunks, i)
		if InputChunkGetType(c) == InputChunkTypeImage {
			if BatchAddChunk(0, c) != BatchAddError {
				t.Fatal("BatchAddChunk with zero batch should return BatchAddError")
			}
			break
		}
	}
}

func TestBatchAddChunkTextChunkRejected(t *testing.T) {
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

	batch, err := BatchInit(ctx)
	if err != nil {
		t.Fatalf("BatchInit failed: %v", err)
	}
	defer BatchFree(batch)

	// Text chunks must be rejected (result != BatchAddSuccess)
	for i := uint64(0); i < InputChunksSize(chunks); i++ {
		c := InputChunksGet(chunks, i)
		if InputChunkGetType(c) == InputChunkTypeText {
			result := BatchAddChunk(batch, c)
			if result == BatchAddSuccess {
				t.Fatal("BatchAddChunk should reject text chunks")
			}
			t.Logf("BatchAddChunk correctly rejected text chunk with result: %d", result)
			return
		}
	}
	t.Skip("no text chunk found to test rejection")
}

func TestBatchEncode(t *testing.T) {
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

	batch, err := BatchInit(ctx)
	if err != nil {
		t.Fatalf("BatchInit failed: %v", err)
	}
	defer BatchFree(batch)

	var imageChunk InputChunk
	for i := uint64(0); i < InputChunksSize(chunks); i++ {
		c := InputChunksGet(chunks, i)
		if InputChunkGetType(c) == InputChunkTypeImage {
			imageChunk = c
			break
		}
	}
	if imageChunk == 0 {
		t.Fatal("no image chunk found in tokenized output")
	}

	if result := BatchAddChunk(batch, imageChunk); result != BatchAddSuccess {
		t.Fatalf("BatchAddChunk failed with result: %d", result)
	}

	if err := BatchEncode(batch); err != nil {
		t.Fatalf("BatchEncode failed: %v", err)
	}

	t.Log("BatchEncode successfully encoded the batch")
}

func TestBatchEncodeInvalidHandle(t *testing.T) {
	if err := BatchEncode(0); err == nil {
		t.Fatal("BatchEncode with zero batch should return an error")
	}
}

func TestBatchGetOutputEmbd(t *testing.T) {
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

	batch, err := BatchInit(ctx)
	if err != nil {
		t.Fatalf("BatchInit failed: %v", err)
	}
	defer BatchFree(batch)

	var imageChunk InputChunk
	for i := uint64(0); i < InputChunksSize(chunks); i++ {
		c := InputChunksGet(chunks, i)
		if InputChunkGetType(c) == InputChunkTypeImage {
			imageChunk = c
			break
		}
	}
	if imageChunk == 0 {
		t.Fatal("no image chunk found in tokenized output")
	}

	if result := BatchAddChunk(batch, imageChunk); result != BatchAddSuccess {
		t.Fatalf("BatchAddChunk failed with result: %d", result)
	}

	if err := BatchEncode(batch); err != nil {
		t.Fatalf("BatchEncode failed: %v", err)
	}

	sz := llama.ModelNEmbdInp(model) * int32(InputChunkGetNTokens(imageChunk))
	if sz <= 0 {
		t.Fatal("calculated embedding size is invalid")
	}

	embd, err := BatchGetOutputEmbd(batch, imageChunk, sz)
	if err != nil {
		t.Fatalf("BatchGetOutputEmbd failed: %v", err)
	}
	if len(embd) == 0 {
		t.Fatal("BatchGetOutputEmbd returned an empty slice")
	}

	t.Logf("BatchGetOutputEmbd returned embeddings of length: %d", len(embd))
}

func TestBatchGetOutputEmbdInvalidHandles(t *testing.T) {
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

	var imageChunk InputChunk
	for i := uint64(0); i < InputChunksSize(chunks); i++ {
		c := InputChunksGet(chunks, i)
		if InputChunkGetType(c) == InputChunkTypeImage {
			imageChunk = c
			break
		}
	}
	if imageChunk == 0 {
		t.Fatal("no image chunk found in tokenized output")
	}

	batch, err := BatchInit(ctx)
	if err != nil {
		t.Fatalf("BatchInit failed: %v", err)
	}
	defer BatchFree(batch)

	// Zero batch handle
	if _, err := BatchGetOutputEmbd(0, imageChunk, 1); err == nil {
		t.Fatal("BatchGetOutputEmbd with zero batch should return an error")
	}

	// Zero chunk handle
	if _, err := BatchGetOutputEmbd(batch, 0, 1); err == nil {
		t.Fatal("BatchGetOutputEmbd with zero chunk should return an error")
	}
}

func TestBatchAddResultConstants(t *testing.T) {
	// Verify the constant values match the C API contract
	if BatchAddSuccess != 0 {
		t.Fatalf("BatchAddSuccess expected 0, got %d", BatchAddSuccess)
	}
	if BatchAddError != 1 {
		t.Fatalf("BatchAddError expected 1, got %d", BatchAddError)
	}
	if BatchAddTooLarge != 2 {
		t.Fatalf("BatchAddTooLarge expected 2, got %d", BatchAddTooLarge)
	}
	if BatchAddIncompatible != 3 {
		t.Fatalf("BatchAddIncompatible expected 3, got %d", BatchAddIncompatible)
	}
}

func TestBatchEncodeMultipleChunks(t *testing.T) {
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

	// Tokenize the same image twice to get two identical image chunks
	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open image file")
	}

	bitmap1 := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap1)
	bitmap2 := BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer BitmapFree(bitmap2)

	chunks1 := InputChunksInit()
	defer InputChunksFree(chunks1)
	text1 := NewInputText("First image: <__media__>", true, true)
	if Tokenize(ctx, chunks1, text1, []Bitmap{bitmap1}) != 0 {
		t.Fatal("Tokenize failed for first image")
	}

	chunks2 := InputChunksInit()
	defer InputChunksFree(chunks2)
	text2 := NewInputText("Second image: <__media__>", false, false)
	if Tokenize(ctx, chunks2, text2, []Bitmap{bitmap2}) != 0 {
		t.Fatal("Tokenize failed for second image")
	}

	var imgChunk1, imgChunk2 InputChunk
	for i := uint64(0); i < InputChunksSize(chunks1); i++ {
		c := InputChunksGet(chunks1, i)
		if InputChunkGetType(c) == InputChunkTypeImage {
			imgChunk1 = c
			break
		}
	}
	for i := uint64(0); i < InputChunksSize(chunks2); i++ {
		c := InputChunksGet(chunks2, i)
		if InputChunkGetType(c) == InputChunkTypeImage {
			imgChunk2 = c
			break
		}
	}
	if imgChunk1 == 0 || imgChunk2 == 0 {
		t.Fatal("could not find image chunks")
	}

	batch, err := BatchInit(ctx)
	if err != nil {
		t.Fatalf("BatchInit failed: %v", err)
	}
	defer BatchFree(batch)

	// Add both same-size image chunks; same-size images should be batchable
	r1 := BatchAddChunk(batch, imgChunk1)
	if r1 != BatchAddSuccess {
		t.Fatalf("BatchAddChunk for first chunk failed with result: %d", r1)
	}

	r2 := BatchAddChunk(batch, imgChunk2)
	switch r2 {
	case BatchAddSuccess:
		t.Log("both chunks added to batch; encoding together")
		if err := BatchEncode(batch); err != nil {
			t.Fatalf("BatchEncode failed: %v", err)
		}
		sz1 := llama.ModelNEmbdInp(model) * int32(InputChunkGetNTokens(imgChunk1))
		embd1, err := BatchGetOutputEmbd(batch, imgChunk1, sz1)
		if err != nil {
			t.Fatalf("BatchGetOutputEmbd for chunk1 failed: %v", err)
		}
		sz2 := llama.ModelNEmbdInp(model) * int32(InputChunkGetNTokens(imgChunk2))
		embd2, err := BatchGetOutputEmbd(batch, imgChunk2, sz2)
		if err != nil {
			t.Fatalf("BatchGetOutputEmbd for chunk2 failed: %v", err)
		}
		t.Logf("encoded 2 chunks in one batch; embd1 len=%d embd2 len=%d", len(embd1), len(embd2))

	case BatchAddTooLarge:
		t.Log("second chunk not added: batch too large (batch_max_tokens limit); encoding first chunk only")
		if err := BatchEncode(batch); err != nil {
			t.Fatalf("BatchEncode failed: %v", err)
		}

	case BatchAddIncompatible:
		t.Log("second chunk not added: incompatible with existing chunks; encoding first chunk only")
		if err := BatchEncode(batch); err != nil {
			t.Fatalf("BatchEncode failed: %v", err)
		}

	default:
		t.Fatalf("BatchAddChunk for second chunk returned unexpected result: %d", r2)
	}
}
