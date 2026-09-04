//go:build js && wasm

package llamawasm

import (
	"errors"
	"testing"
)

func TestBatchGetOne(t *testing.T) {
	batch := BatchGetOne([]Token{1, 2, 3})

	if batch.NTokens != 3 {
		t.Errorf("NTokens = %d, want 3", batch.NTokens)
	}
	if len(batch.Tokens()) != 3 {
		t.Errorf("len(Tokens()) = %d, want 3", len(batch.Tokens()))
	}
	if batch.writable() {
		t.Error("a batch from BatchGetOne must not be writable")
	}
	if err := batch.Add(1, 0, []SeqId{0}, true); !errors.Is(err, ErrBatchNotWritable) {
		t.Errorf("Add gave %v, want ErrBatchNotWritable", err)
	}
	if err := batch.SetLogit(0, true); !errors.Is(err, ErrBatchNotWritable) {
		t.Errorf("SetLogit gave %v, want ErrBatchNotWritable", err)
	}
}

func TestBatchInit(t *testing.T) {
	batch := BatchInit(4, 0, 2)
	defer BatchFree(batch)

	if !batch.writable() {
		t.Fatal("a batch from BatchInit must be writable")
	}
	if batch.NTokens != 0 {
		t.Errorf("NTokens = %d, want 0", batch.NTokens)
	}
}

// An embedding cannot be the input of a batch here, and a size below one has no
// meaning. Each one gives a batch that holds nothing.
func TestBatchInitBadArguments(t *testing.T) {
	tests := []struct {
		name                   string
		nTokens, embd, nSeqMax int32
	}{
		{"no tokens", 0, 0, 1},
		{"no sequences", 4, 0, 0},
		{"an embedding", 4, 64, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch := BatchInit(tt.nTokens, tt.embd, tt.nSeqMax)
			if batch.writable() {
				t.Error("the batch must not be writable")
			}
		})
	}
}

func TestBatchAdd(t *testing.T) {
	batch := BatchInit(3, 0, 2)
	defer BatchFree(batch)

	if err := batch.Add(11, 0, []SeqId{0, 1}, false); err != nil {
		t.Fatalf("Add gave %v", err)
	}
	if err := batch.Add(22, 1, []SeqId{1}, true); err != nil {
		t.Fatalf("Add gave %v", err)
	}

	if batch.NTokens != 2 {
		t.Errorf("NTokens = %d, want 2", batch.NTokens)
	}

	want := []Token{11, 22}
	got := batch.Tokens()
	if len(got) != len(want) {
		t.Fatalf("len(Tokens()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Tokens()[%d] = %d, want %d", i, got[i], want[i])
		}
	}

	if batch.pos[1] != 1 {
		t.Errorf("pos[1] = %d, want 1", batch.pos[1])
	}
	if batch.nSeqID[0] != 2 || batch.nSeqID[1] != 1 {
		t.Errorf("nSeqID = %v, want [2 1]", batch.nSeqID[:2])
	}

	// The identifiers of each token take capSeq places, one token after the
	// other, which is the shape that the shim takes.
	if batch.seqIDs[0] != 0 || batch.seqIDs[1] != 1 || batch.seqIDs[2] != 1 {
		t.Errorf("seqIDs = %v, want [0 1 1 ...]", batch.seqIDs)
	}

	if batch.logits[0] != 0 || batch.logits[1] != 1 {
		t.Errorf("logits = %v, want [0 1]", batch.logits[:2])
	}
}

func TestBatchAddFull(t *testing.T) {
	batch := BatchInit(1, 0, 1)
	defer BatchFree(batch)

	if err := batch.Add(1, 0, []SeqId{0}, true); err != nil {
		t.Fatalf("Add gave %v", err)
	}
	if err := batch.Add(2, 1, []SeqId{0}, true); !errors.Is(err, ErrBatchFull) {
		t.Errorf("Add gave %v, want ErrBatchFull", err)
	}
	if batch.NTokens != 1 {
		t.Errorf("NTokens = %d, want 1", batch.NTokens)
	}
}

func TestBatchAddTooManySeqIDs(t *testing.T) {
	batch := BatchInit(2, 0, 1)
	defer BatchFree(batch)

	if err := batch.Add(1, 0, []SeqId{0, 1}, true); !errors.Is(err, ErrTooManySeqIDs) {
		t.Errorf("Add gave %v, want ErrTooManySeqIDs", err)
	}
	if batch.NTokens != 0 {
		t.Errorf("NTokens = %d, want 0", batch.NTokens)
	}
}

func TestBatchSetLogitOutOfRange(t *testing.T) {
	batch := BatchInit(4, 0, 1)
	defer BatchFree(batch)

	if err := batch.Add(1, 0, []SeqId{0}, false); err != nil {
		t.Fatalf("Add gave %v", err)
	}

	// Index 1 has room in the arrays, but the batch does not hold a token
	// there yet, and llama.cpp reads no flag beyond NTokens.
	for _, idx := range []int32{-1, 1, 4} {
		if err := batch.SetLogit(idx, true); !errors.Is(err, ErrBatchIndexRange) {
			t.Errorf("SetLogit(%d) gave %v, want ErrBatchIndexRange", idx, err)
		}
	}

	if err := batch.SetLogit(0, true); err != nil {
		t.Errorf("SetLogit(0) gave %v", err)
	}
	if batch.logits[0] != 1 {
		t.Errorf("logits[0] = %d, want 1", batch.logits[0])
	}
}

func TestBatchClear(t *testing.T) {
	batch := BatchInit(2, 0, 1)
	defer BatchFree(batch)

	if err := batch.Add(1, 0, []SeqId{0}, true); err != nil {
		t.Fatalf("Add gave %v", err)
	}
	if err := batch.Clear(); err != nil {
		t.Fatalf("Clear gave %v", err)
	}
	if batch.NTokens != 0 {
		t.Errorf("NTokens = %d, want 0", batch.NTokens)
	}

	// A batch that is clear takes tokens again from the start.
	if err := batch.Add(2, 0, []SeqId{0}, true); err != nil {
		t.Fatalf("Add gave %v", err)
	}
	if batch.tokens[0] != 2 {
		t.Errorf("tokens[0] = %d, want 2", batch.tokens[0])
	}
}
