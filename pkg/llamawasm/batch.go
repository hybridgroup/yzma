//go:build js && wasm

package llamawasm

import (
	"errors"
	"fmt"
)

// Errors of [Batch.Add] and [Batch.SetLogit]. They agree with the errors of the
// llama package, thus the same code builds for a native platform and for a
// browser.
var (
	// ErrBatchFull means that the batch already holds the tokens that
	// BatchInit gave it room for.
	ErrBatchFull = errors.New("batch is full")

	// ErrTooManySeqIDs means that the call gave more sequence identifiers
	// than the nSeqMax of the batch.
	ErrTooManySeqIDs = errors.New("too many sequence IDs for batch")

	// ErrBatchNotWritable means that the batch holds no arrays to write. A
	// batch of the zero value and a batch from BatchGetOne do not.
	ErrBatchNotWritable = errors.New("batch owns no writable token arrays")

	// ErrBatchIndexRange means that the index does not address a token of the
	// batch.
	ErrBatchIndexRange = errors.New("batch index out of range")

	// ErrNoBatch says that the module is from a release before the calls that
	// take a batch with positions, which are in ABI version 6 and later.
	ErrNoBatch = errors.New("llamawasm: this llama.cpp module has no batch calls, install a newer build")
)

// BatchGetOne makes a batch that holds the given tokens. The positions of the
// tokens continue from the state of the context that decodes the batch, and
// every token goes to sequence 0.
//
// The batch holds no arrays of its own, thus [Batch.Add] and [Batch.SetLogit]
// report [ErrBatchNotWritable] on it. Use [BatchInit] for a batch to fill in.
func BatchGetOne(tokens []Token) Batch {
	return Batch{
		NTokens: int32(len(tokens)),
		tokens:  tokens,
	}
}

// BatchInit makes a batch with room for nTokens tokens, each one with a
// maximum of nSeqMax sequence identifiers. Fill it with [Batch.Add].
//
// The embd argument is always 0 here. A WebAssembly module cannot take an
// embedding as the input of a batch, because the shim has no call for it.
//
// The batch is an ordinary Go value, thus [BatchFree] is not necessary. It
// exists so that the same code builds for a native platform.
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

// BatchFree does nothing. The memory of a batch belongs to Go here, thus the
// garbage collector takes it.
func BatchFree(batch Batch) error {
	return nil
}

// Tokens gives the tokens of the batch.
func (b Batch) Tokens() []Token {
	return b.tokens[:b.NTokens]
}

// Clear sets the number of tokens of the batch to zero.
func (b *Batch) Clear() error {
	b.NTokens = 0

	return nil
}

// writable tells if the batch holds the arrays that Add and SetLogit write.
func (b *Batch) writable() bool {
	return b.capTokens > 0 && b.pos != nil && b.nSeqID != nil && b.seqIDs != nil && b.logits != nil
}

// SetLogit sets if the model computes the logits of the token at index idx.
//
// The index must address a token that the batch already holds, thus it must be
// below NTokens. llama.cpp reads a flag only inside that range.
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

// Add puts a token in the batch with its position, its sequences, and its
// logit flag.
//
// It writes nothing and gives an error if the batch is full ([ErrBatchFull]),
// if seqIDs holds more identifiers than the nSeqMax of the batch
// ([ErrTooManySeqIDs]), or if the batch holds no arrays to write
// ([ErrBatchNotWritable]).
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

	// SetLogit measures the index against NTokens, thus the count must hold
	// the new token before the flag of it can change.
	b.NTokens++

	return b.SetLogit(i, logits)
}
