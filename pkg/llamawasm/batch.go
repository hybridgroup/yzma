//go:build js && wasm

package llamawasm

import (
	"errors"
	"fmt"
)

// Errors from [Batch.Add] and [Batch.SetLogit]. They match the llama package
// errors, so the same code builds for a native platform and for a browser.
var (
	// ErrBatchFull means the batch already holds as many tokens as
	// BatchInit gave it room for.
	ErrBatchFull = errors.New("batch is full")

	// ErrTooManySeqIDs means the call passed more sequence IDs than the
	// batch nSeqMax.
	ErrTooManySeqIDs = errors.New("too many sequence IDs for batch")

	// ErrBatchNotWritable means the batch has no arrays to write. A zero
	// value batch and a batch from BatchGetOne have none.
	ErrBatchNotWritable = errors.New("batch owns no writable token arrays")

	// ErrBatchIndexRange means the index does not point to a token in the
	// batch.
	ErrBatchIndexRange = errors.New("batch index out of range")

	// ErrNoBatch means the module predates the calls that take a batch with
	// positions, which arrived in ABI version 6.
	ErrNoBatch = errors.New("llamawasm: this llama.cpp module has no batch calls, install a newer build")
)

// BatchGetOne creates a batch that holds the given tokens. Token positions
// continue from the state of the context that decodes the batch, and every
// token goes to sequence 0.
//
// The batch has no arrays of its own, so [Batch.Add] and [Batch.SetLogit]
// return [ErrBatchNotWritable] on it. Use [BatchInit] for a batch to fill in.
func BatchGetOne(tokens []Token) Batch {
	return Batch{
		NTokens: int32(len(tokens)),
		tokens:  tokens,
	}
}

// BatchInit creates a batch with room for nTokens tokens, each with up to
// nSeqMax sequence IDs. Fill it with [Batch.Add].
//
// The embd argument is always 0 here. A WebAssembly module cannot take an
// embedding as batch input, because the shim has no call for it.
//
// The batch is an ordinary Go value, so [BatchFree] is not needed. It
// exists so the same code builds for a native platform.
func BatchInit(nTokens int32, embd int32, nSeqMax int32) Batch {
	if nTokens < 1 || nSeqMax < 1 || embd != 0 {
		return Batch{}
	}

	n := int(nTokens)
	return Batch{
		NTokens:   0,
		tokens:    make([]Token, n),
		pos:       make([]Pos, n),
		nSeqID:    make([]int32, n),
		seqIDs:    make([]SeqId, n*int(nSeqMax)),
		logits:    make([]int8, n),
		capTokens: nTokens,
		capSeq:    nSeqMax,
	}
}

// BatchFree does nothing. Batch memory belongs to Go here, so the garbage
// collector reclaims it.
func BatchFree(batch Batch) error {
	return nil
}

// Tokens returns the tokens in the batch.
func (b Batch) Tokens() []Token {
	return b.tokens[:b.NTokens]
}

// Clear sets the batch token count to zero.
func (b *Batch) Clear() error {
	b.NTokens = 0

	return nil
}

// writable reports whether the batch has the arrays that Add and SetLogit write.
func (b *Batch) writable() bool {
	return b.capTokens > 0 && b.pos != nil && b.nSeqID != nil && b.seqIDs != nil && b.logits != nil
}

// SetLogit sets whether the model computes logits for the token at index idx.
//
// The index must point to a token the batch already holds, so it must be
// below NTokens. llama.cpp only reads flags inside that range.
func (b *Batch) SetLogit(idx int32, logits bool) error {
	if !b.writable() {
		return ErrBatchNotWritable
	}
	if idx < 0 || idx >= b.NTokens {
		return fmt.Errorf("%w: index %d not in [0,%d)", ErrBatchIndexRange, idx, b.NTokens)
	}

	if logits {
		b.logits[idx] = 1
	} else {
		b.logits[idx] = 0
	}

	return nil
}

// Add appends a token to the batch with its position, sequences, and logit
// flag.
//
// It writes nothing and returns an error if the batch is full ([ErrBatchFull]),
// if seqIDs has more IDs than the batch nSeqMax ([ErrTooManySeqIDs]), or if
// the batch has no arrays to write ([ErrBatchNotWritable]).
func (b *Batch) Add(token Token, pos Pos, seqIDs []SeqId, logits bool) error {
	if !b.writable() {
		return ErrBatchNotWritable
	}

	i := b.NTokens

	if i < 0 || i >= b.capTokens {
		return fmt.Errorf("%w: index %d, capacity %d", ErrBatchFull, i, b.capTokens)
	}
	if int32(len(seqIDs)) > b.capSeq {
		return fmt.Errorf("%w: %d sequence IDs for a batch with n_seq_max %d", ErrTooManySeqIDs, len(seqIDs), b.capSeq)
	}

	b.tokens[i] = token
	b.pos[i] = pos
	b.nSeqID[i] = int32(len(seqIDs))

	start := int(i) * int(b.capSeq)
	copy(b.seqIDs[start:start+int(b.capSeq)], seqIDs)

	// SetLogit checks the index against NTokens, so the count must include
	// the new token before its flag can change.
	b.NTokens++

	return b.SetLogit(i, logits)
}
