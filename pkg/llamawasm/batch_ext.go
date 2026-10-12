//go:build js && wasm

package llamawasm

import (
	"errors"
	"fmt"
)

// BatchExt is an extended batch. The entries stay in Go until [Process]
// sends them to the module, so a token embedding cannot be an entry here.
type BatchExt struct {
	b *batchExt
}

type batchExt struct {
	ctx     Context
	cap     int
	nSeqMax int
	tokens  []Token
	pos     []Pos
	posSet  []bool
	seqs    [][]SeqId
	outputs []int8
	order   []int32
	ordered bool
}

// ProcessType selects how [Process] runs a [BatchExt].
type ProcessType int32

const (
	ProcessTypeEncode ProcessType = 0
	ProcessTypeDecode ProcessType = 1
)

// DecisionOrder marks the entries that a decision head (clef, vllm-decision2) reads.
// A run of entries with the same value is one span, and spans must be separated by [DecisionOrderNone].
type DecisionOrder int32

const (
	DecisionOrderNone           DecisionOrder = 0
	DecisionOrderQuestionNoul   DecisionOrder = 1
	DecisionOrderQuestionChoice DecisionOrder = 2
	DecisionOrderQuestionScore  DecisionOrder = 3
	DecisionOrderOption         DecisionOrder = 4
)

// Errors reported by the [BatchExt] functions. They match the llama package errors.
var (
	// ErrBatchExtInvalidToken means the token ID was rejected.
	ErrBatchExtInvalidToken = errors.New("invalid token or embedding for batch")
	// ErrBatchExtInvalidSeqID means the sequence ID is not below the context n_seq_max.
	ErrBatchExtInvalidSeqID = errors.New("invalid sequence ID for batch")
	// ErrBatchExtRejected means the index or value was rejected.
	ErrBatchExtRejected = errors.New("batch rejected the value")

	// ErrNoBatchExt means the module predates the extended batch, which arrived in ABI version 11.
	ErrNoBatchExt = errors.New("llamawasm: this llama.cpp module has no extended batch, install a newer build")

	errInvalidBatchExt = errors.New("invalid extended batch")
)

// BatchExtInit creates an extended batch for the context.
// It holds up to [NBatch] entries and accepts sequence IDs below [NSeqMax].
func BatchExtInit(ctx Context) (BatchExt, error) {
	if !Loaded() {
		return BatchExt{}, ErrNotLoaded
	}
	if !has("_yzma_process") {
		return BatchExt{}, ErrNoBatchExt
	}
	n, seqs := NBatch(ctx), NSeqMax(ctx)
	if n == 0 || seqs == 0 {
		return BatchExt{}, fmt.Errorf("llamawasm: invalid context handle %d", ctx)
	}

	return BatchExt{&batchExt{ctx: ctx, cap: int(n), nSeqMax: int(seqs)}}, nil
}

// BatchExtFree does nothing. The batch belongs to Go here.
func BatchExtFree(batch BatchExt) error {
	if batch.b == nil {
		return errInvalidBatchExt
	}
	return nil
}

// BatchExtClear removes all entries from the batch.
func BatchExtClear(batch BatchExt) error {
	b := batch.b
	if b == nil {
		return errInvalidBatchExt
	}
	b.tokens, b.pos, b.posSet = b.tokens[:0], b.pos[:0], b.posSet[:0]
	b.seqs, b.outputs, b.order = b.seqs[:0], b.outputs[:0], b.order[:0]
	b.ordered = false

	return nil
}

// BatchExtAddToken adds a token and returns its index.
// Set its position with [BatchExtSetPos] before processing the batch.
func BatchExtAddToken(batch BatchExt, seqID SeqId, token Token) (int32, error) {
	b := batch.b
	if b == nil {
		return 0, errInvalidBatchExt
	}
	if len(b.tokens) >= b.cap {
		return -1, ErrBatchFull
	}
	if token < 0 {
		return -2, ErrBatchExtInvalidToken
	}
	if seqID < 0 || int(seqID) >= b.nSeqMax {
		return -3, ErrBatchExtInvalidSeqID
	}

	b.tokens = append(b.tokens, token)
	b.pos = append(b.pos, 0)
	b.posSet = append(b.posSet, false)
	b.seqs = append(b.seqs, []SeqId{seqID})
	b.outputs = append(b.outputs, 0)
	b.order = append(b.order, 0)

	return int32(len(b.tokens) - 1), nil
}

// BatchExtAddSeq adds the entry at idx to another sequence at the same position.
func BatchExtAddSeq(batch BatchExt, idx int32, seqID SeqId) error {
	b, err := batch.entry(idx)
	if err != nil {
		return err
	}
	if seqID < 0 || int(seqID) >= b.nSeqMax || len(b.seqs[idx]) >= b.nSeqMax {
		return fmt.Errorf("%w: add sequence %d to index %d", ErrBatchExtRejected, seqID, idx)
	}
	b.seqs[idx] = append(b.seqs[idx], seqID)

	return nil
}

// BatchExtSetOutputEmbd sets whether output embeddings are computed for the entry at idx.
// llama.cpp currently treats this the same as [BatchExtSetOutputLogits].
func BatchExtSetOutputEmbd(batch BatchExt, idx int32, value bool) error {
	return BatchExtSetOutputLogits(batch, idx, value)
}

// BatchExtSetOutputLogits sets whether logits are computed for the entry at idx.
// llama.cpp currently treats this the same as [BatchExtSetOutputEmbd].
func BatchExtSetOutputLogits(batch BatchExt, idx int32, value bool) error {
	b, err := batch.entry(idx)
	if err != nil {
		return err
	}
	b.outputs[idx] = 0
	if value {
		b.outputs[idx] = 1
	}

	return nil
}

// BatchExtSetPos sets the position of the entry at idx. A token uses one position.
func BatchExtSetPos(batch BatchExt, idx int32, pos ...Pos) error {
	b, err := batch.entry(idx)
	if err != nil {
		return err
	}
	if len(pos) != 1 {
		return fmt.Errorf("%w: %d positions, want 1", ErrBatchExtRejected, len(pos))
	}
	b.pos[idx], b.posSet[idx] = pos[0], true

	return nil
}

// BatchExtSetDecisionOrder sets which part of a decision prompt the entry at idx belongs to.
func BatchExtSetDecisionOrder(batch BatchExt, idx int32, order DecisionOrder) error {
	b, err := batch.entry(idx)
	if err != nil {
		return err
	}
	if order < DecisionOrderNone || order > DecisionOrderOption {
		return fmt.Errorf("%w: set decision order at index %d", ErrBatchExtRejected, idx)
	}
	b.order[idx] = int32(order)
	b.ordered = true

	return nil
}

// Process encodes or decodes the batch. The return values are the same as [Decode].
func Process(ctx Context, typ ProcessType, batch BatchExt) (int32, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	b := batch.b
	if b == nil {
		return 0, errInvalidBatchExt
	}
	if !has("_yzma_process") {
		return 0, ErrNoBatchExt
	}
	n := len(b.tokens)
	if n == 0 {
		return 0, nil
	}
	if i := indexFalse(b.posSet); i >= 0 {
		return 0, fmt.Errorf("%w: index %d has no position", ErrBatchExtRejected, i)
	}

	nSeqMax := 1
	for _, s := range b.seqs {
		nSeqMax = max(nSeqMax, len(s))
	}
	nSeq := make([]int32, n)
	seqIDs := make([]SeqId, n*nSeqMax)
	for i, s := range b.seqs {
		nSeq[i] = int32(len(s))
		copy(seqIDs[i*nSeqMax:], s)
	}

	ptr, err := tokenScratch.reserve(n * 4)
	if err != nil {
		return 0, err
	}
	writeTokens(ptr, b.tokens)

	posPtr, err := posScratch.reserve(n * 4)
	if err != nil {
		return 0, err
	}
	writeInt32s(posPtr, b.pos)

	nSeqPtr, err := nSeqScratch.reserve(n * 4)
	if err != nil {
		return 0, err
	}
	writeInt32s(nSeqPtr, nSeq)

	seqPtr, err := seqScratch.reserve(len(seqIDs) * 4)
	if err != nil {
		return 0, err
	}
	writeInt32s(seqPtr, seqIDs)

	outPtr, err := logitScratch.reserve(n)
	if err != nil {
		return 0, err
	}
	writeInt8s(outPtr, b.outputs)

	orderPtr := 0
	if b.ordered {
		if orderPtr, err = orderScratch.reserve(n * 4); err != nil {
			return 0, err
		}
		writeInt32s(orderPtr, b.order)
	}

	return callErr("_yzma_process", int(ctx), int(typ), ptr, posPtr, nSeqPtr, seqPtr, nSeqMax, outPtr, orderPtr, n)
}

// entry returns the batch data when idx is an entry of the batch.
func (batch BatchExt) entry(idx int32) (*batchExt, error) {
	b := batch.b
	if b == nil {
		return nil, errInvalidBatchExt
	}
	if idx < 0 || int(idx) >= len(b.tokens) {
		return nil, fmt.Errorf("%w: index %d not in [0,%d)", ErrBatchExtRejected, idx, len(b.tokens))
	}
	return b, nil
}

func indexFalse(v []bool) int {
	for i, ok := range v {
		if !ok {
			return i
		}
	}
	return -1
}
