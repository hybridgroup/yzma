package llama

import (
	"testing"
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
