package mtmd

import (
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestInputChunksInitAndFree(t *testing.T) {
	chunks := InputChunksInit()
	if chunks == InputChunks(0) {
		t.Fatal("InputChunksInit returned an invalid InputChunks")
	}

	t.Log("InputChunksInit successfully initialized InputChunks")

	InputChunksFree(chunks)
	t.Log("InputChunksFree successfully freed InputChunks")
}

func TestInputChunksSize(t *testing.T) {
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

	model := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx := InitFromFile(mmprojFile, model, params)
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	size := InputChunksSize(chunks)
	t.Logf("InputChunksSize returned: %d", size)
	if size == 0 {
		t.Fatalf("invalid chunk size: %d", size)
	}

	idx := uint32(1)
	chunk := InputChunksGet(chunks, idx)
	if chunk == InputChunk(0) {
		t.Fatalf("InputChunksGet returned an invalid chunk for index %d", idx)
	}

	t.Logf("InputChunksGet successfully retrieved chunk at index %d", idx)

	chunkType := InputChunkGetType(chunk)
	t.Logf("InputChunkGetType returned: %d", chunkType)

	switch chunkType {
	case InputChunkTypeText:
		tokens := InputChunkGetTokensText(chunk)
		t.Logf("InputChunkGetTokensText returned %d tokens", len(tokens))
	}
}

func TestInputChunkGetNTokens(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx := InitFromFile(mmprojFile, model, params)
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint32(1)
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

	model := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx := InitFromFile(mmprojFile, model, params)
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint32(1)
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

	model := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx := InitFromFile(mmprojFile, model, params)
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint32(1)
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

	model := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	params := ContextParamsDefault()
	ctx := InitFromFile(mmprojFile, model, params)
	defer Free(ctx)

	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	testSetupChunks(t, ctx, chunks)

	idx := uint32(1)
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
