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
