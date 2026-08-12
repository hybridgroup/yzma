package llama

import (
	"errors"
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

func TestBatchAddFull(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	batch := BatchInit(2, 0, 1)
	defer BatchFree(batch)

	if err := batch.Add(1, 0, []SeqId{0}, true); err != nil {
		t.Fatalf("first Add returned error: %v", err)
	}
	if err := batch.Add(2, 1, []SeqId{0}, true); err != nil {
		t.Fatalf("second Add returned error: %v", err)
	}

	// Batch is now at capacity (2); a third Add must fail without writing.
	if err := batch.Add(3, 2, []SeqId{0}, true); err == nil {
		t.Fatal("Add past capacity did not return an error")
	}
	if batch.NTokens != 2 {
		t.Errorf("Add past capacity mutated NTokens, got %d want 2", batch.NTokens)
	}
}

func TestBatchAddTooManySeqIDs(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// n_seq_max == 1, so a token carrying 2 sequence IDs must be rejected.
	batch := BatchInit(2, 0, 1)
	defer BatchFree(batch)

	if err := batch.Add(1, 0, []SeqId{0, 1}, true); err == nil {
		t.Fatal("Add with seqIDs longer than n_seq_max did not return an error")
	}
	if batch.NTokens != 0 {
		t.Errorf("rejected Add mutated NTokens, got %d want 0", batch.NTokens)
	}
}

func TestBatchSetLogitOutOfRange(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	batch := BatchInit(4, 0, 1)
	defer BatchFree(batch)

	// An empty batch has no token at any index, so every index is rejected.
	if err := batch.SetLogit(0, true); !errors.Is(err, ErrBatchIndexRange) {
		t.Errorf("SetLogit(0) on an empty batch = %v, want ErrBatchIndexRange", err)
	}

	if err := batch.Add(1, 0, []SeqId{0}, true); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if err := batch.Add(2, 1, []SeqId{0}, true); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}

	if err := batch.SetLogit(-1, true); !errors.Is(err, ErrBatchIndexRange) {
		t.Errorf("SetLogit(-1) = %v, want ErrBatchIndexRange", err)
	}
	if err := batch.SetLogit(0, true); err != nil {
		t.Errorf("SetLogit(0) returned error: %v", err)
	}
	if err := batch.SetLogit(1, true); err != nil {
		t.Errorf("SetLogit(1) returned error: %v", err)
	}

	// Index 2 is inside the allocation but past the two tokens actually in the
	// batch. llama_decode only reads logit flags over [0, n_tokens), so a flag
	// set here would be dropped; it is a caller index bug, not a write to
	// honour silently.
	if err := batch.SetLogit(2, true); !errors.Is(err, ErrBatchIndexRange) {
		t.Errorf("SetLogit past NTokens = %v, want ErrBatchIndexRange", err)
	}
	if err := batch.SetLogit(4, true); !errors.Is(err, ErrBatchIndexRange) {
		t.Errorf("SetLogit at capacity = %v, want ErrBatchIndexRange", err)
	}
}

// TestBatchNotWritable pins that the batches which own no writable arrays are
// refused rather than dereferenced: llama_batch_get_one leaves pos, n_seq_id,
// seq_id and logits NULL, and llama_batch_init leaves token NULL when it is
// asked for an embedding batch.
func TestBatchNotWritable(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	t.Run("zero_value", func(t *testing.T) {
		var batch Batch
		if err := batch.Add(1, 0, []SeqId{0}, true); !errors.Is(err, ErrBatchNotWritable) {
			t.Errorf("Add = %v, want ErrBatchNotWritable", err)
		}
		if err := batch.SetLogit(0, true); !errors.Is(err, ErrBatchNotWritable) {
			t.Errorf("SetLogit = %v, want ErrBatchNotWritable", err)
		}
	})

	t.Run("get_one", func(t *testing.T) {
		tokens := []Token{1, 2, 3}
		batch := BatchGetOne(tokens)
		if batch.NTokens != 3 {
			t.Fatalf("BatchGetOne NTokens = %d, want 3", batch.NTokens)
		}
		if err := batch.Add(4, 3, []SeqId{0}, true); !errors.Is(err, ErrBatchNotWritable) {
			t.Errorf("Add = %v, want ErrBatchNotWritable", err)
		}
		if err := batch.SetLogit(2, true); !errors.Is(err, ErrBatchNotWritable) {
			t.Errorf("SetLogit = %v, want ErrBatchNotWritable", err)
		}
	})

	t.Run("embedding_batch", func(t *testing.T) {
		batch := BatchInit(4, 8, 1)
		defer BatchFree(batch)

		if batch.Token != nil {
			t.Skip("llama_batch_init allocated token for an embedding batch")
		}
		if err := batch.Add(1, 0, []SeqId{0}, true); !errors.Is(err, ErrBatchNotWritable) {
			t.Errorf("Add = %v, want ErrBatchNotWritable", err)
		}
	})
}

// TestBatchDataLayout pins the invariant the libffi descriptor depends on: the
// struct handed to C is BatchData, and it must stay exactly the seven C fields.
// Batch itself carries Go-only bookkeeping and is deliberately larger.
func TestBatchDataLayout(t *testing.T) {
	const cBatchSize = 56 // int32 + 6 pointers, LP64, with tail padding

	if got := unsafe.Sizeof(BatchData{}); got != cBatchSize {
		t.Errorf("sizeof(BatchData) = %d, want %d: it must mirror C llama_batch exactly", got, cBatchSize)
	}
	if got := unsafe.Offsetof(Batch{}.BatchData); got != 0 {
		t.Errorf("BatchData is at offset %d in Batch, want 0", got)
	}
}
