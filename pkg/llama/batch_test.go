package llama

import (
	"testing"
	"unsafe"
)

func TestBatchInit(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	nTokens := int32(1)
	embd := int32(0)
	nSeqMax := int32(1)

	batch := BatchInit(nTokens, embd, nSeqMax)
	if batch == (Batch{}) {
		t.Fatal("BatchInit returned an empty batch")
	}

	BatchFree(batch)
}

func TestBatchGetOne(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	tokens := []Token{1, 2, 3, 4, 5}
	batch := BatchGetOne(tokens)
	if batch == (Batch{}) {
		t.Fatal("BatchGetOne returned an empty batch")
	}
}

func TestBatchClear(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	batch := BatchInit(2, 0, 1)
	batch.NTokens = 2
	err := batch.Clear()
	if err != nil {
		t.Fatalf("Clear returned error: %v", err)
	}
	if batch.NTokens != 0 {
		t.Errorf("Clear did not reset NTokens to 0, got %d", batch.NTokens)
	}
	BatchFree(batch)
}

func TestBatchSetLogit(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	batch := BatchInit(2, 0, 1)
	batch.NTokens = 2

	batch.SetLogit(0, true)
	batch.SetLogit(1, false)

	logits := unsafe.Slice((*int8)(batch.Logits), int(batch.NTokens))
	if logits[0] != 1 {
		t.Errorf("SetLogit did not set index 0 to 1, got %d", logits[0])
	}
	if logits[1] != 0 {
		t.Errorf("SetLogit did not set index 1 to 0, got %d", logits[1])
	}
	BatchFree(batch)
}

func TestBatchAdd(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	batch := BatchInit(2, 0, 2)
	defer BatchFree(batch)

	token := Token(42)
	pos := Pos(7)
	seqIDs := []SeqId{1, 2}
	logits := true

	batch.Add(token, pos, seqIDs, logits)

	if batch.NTokens != 1 {
		t.Errorf("Add did not increment NTokens, got %d", batch.NTokens)
	}

	tokens := unsafe.Slice((*Token)(batch.Token), int(batch.NTokens))
	if tokens[0] != token {
		t.Errorf("Add did not set token correctly, got %v", tokens[0])
	}

	poses := unsafe.Slice((*Pos)(batch.Pos), int(batch.NTokens))
	if poses[0] != pos {
		t.Errorf("Add did not set pos correctly, got %v", poses[0])
	}

	nSeqIds := unsafe.Slice((*int32)(batch.NSeqId), int(batch.NTokens))
	if nSeqIds[0] != int32(len(seqIDs)) {
		t.Errorf("Add did not set nSeqIds correctly, got %v", nSeqIds[0])
	}

	logitVals := unsafe.Slice((*int8)(batch.Logits), int(batch.NTokens))
	if logitVals[0] != 1 {
		t.Errorf("Add did not set logits correctly, got %v", logitVals[0])
	}
}
