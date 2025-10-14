package llama

import (
	"testing"
)

func TestMemoryClear(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	// Clear memory with data = true
	MemoryClear(mem, true)
	// No direct way to verify, but ensure no panic or error occurs
	t.Log("MemoryClear executed successfully with data = true")

	// Clear memory with data = false
	MemoryClear(mem, false)
	t.Log("MemoryClear executed successfully with data = false")
}

func TestMemorySeqRm(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	// Remove tokens for a specific sequence and position range
	seqID := SeqId(1)
	p0 := Pos(0)
	p1 := Pos(10)

	success := MemorySeqRm(mem, seqID, p0, p1)
	if !success {
		t.Fatal("MemorySeqRm failed to remove tokens for the specified sequence and position range")
	}
	t.Logf("MemorySeqRm executed successfully for seqID: %d, p0: %d, p1: %d", seqID, p0, p1)

	// Test with seqID < 0 (match any sequence)
	success = MemorySeqRm(mem, -1, p0, p1)
	if !success {
		t.Fatal("MemorySeqRm failed to remove tokens for any sequence")
	}
	t.Log("MemorySeqRm executed successfully for seqID < 0 (match any sequence)")

	// Test with p1 < 0 (remove tokens from p0 to infinity)
	success = MemorySeqRm(mem, seqID, p0, -1)
	if !success {
		t.Fatal("MemorySeqRm failed to remove tokens from p0 to infinity")
	}
	t.Log("MemorySeqRm executed successfully for p1 < 0 (remove tokens from p0 to infinity)")
}

func TestMemorySeqCp(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	seqIDSrc := SeqId(1)
	seqIDDst := SeqId(2)
	p0 := Pos(0)
	p1 := Pos(10)

	MemorySeqCp(mem, seqIDSrc, seqIDDst, p0, p1)
	t.Logf("MemorySeqCp executed successfully for seqIDSrc: %d, seqIDDst: %d, p0: %d, p1: %d", seqIDSrc, seqIDDst, p0, p1)
}

func TestMemorySeqKeep(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	seqID := SeqId(1)
	MemorySeqKeep(mem, seqID)
	t.Logf("MemorySeqKeep executed successfully for seqID: %d", seqID)
}

func TestMemorySeqAdd(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	seqID := SeqId(1)
	p0 := Pos(0)
	p1 := Pos(10)
	delta := Pos(5)

	MemorySeqAdd(mem, seqID, p0, p1, delta)
	t.Logf("MemorySeqAdd executed successfully for seqID: %d, p0: %d, p1: %d, delta: %d", seqID, p0, p1, delta)
}

func TestMemorySeqDiv(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	seqID := SeqId(1)
	p0 := Pos(0)
	p1 := Pos(10)
	d := 2

	MemorySeqDiv(mem, seqID, p0, p1, d)
	t.Logf("MemorySeqDiv executed successfully for seqID: %d, p0: %d, p1: %d, d: %d", seqID, p0, p1, d)
}

func TestMemorySeqPosMin(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	seqID := SeqId(1)
	posMin := MemorySeqPosMin(mem, seqID)
	t.Logf("MemorySeqPosMin returned: %d for seqID: %d", posMin, seqID)
}

func TestMemorySeqPosMax(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	seqID := SeqId(1)
	posMax := MemorySeqPosMax(mem, seqID)
	t.Logf("MemorySeqPosMax returned: %d for seqID: %d", posMax, seqID)
}

func TestMemoryCanShift(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(testModelFileName(t), params)
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	mem := GetMemory(ctx)
	if mem == Memory(0) {
		t.Fatal("GetMemory returned an empty memory object")
	}

	canShift := MemoryCanShift(mem)
	t.Logf("MemoryCanShift returned: %v", canShift)
}
