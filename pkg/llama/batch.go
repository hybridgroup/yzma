package llama

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

// Errors reported by [Batch.Add] and [Batch.SetLogit]. They are distinguished so
// a caller can tell an expected condition (a full batch, which just means it is
// time to decode) from a programming error, without matching on message text.
var (
	// ErrBatchFull means the batch already holds its allocated n_tokens.
	ErrBatchFull = errors.New("batch is full")
	// ErrTooManySeqIDs means more sequence IDs were given than the n_seq_max
	// the batch was allocated with.
	ErrTooManySeqIDs = errors.New("too many sequence IDs for batch")
	// ErrBatchNotWritable means the batch owns no writable arrays: it is the
	// zero value, came from [BatchGetOne], or is an embedding batch whose
	// token array was never allocated.
	ErrBatchNotWritable = errors.New("batch owns no writable token arrays")
	// ErrBatchIndexRange means the index does not address a token in the batch.
	ErrBatchIndexRange = errors.New("batch index out of range")
)

var (
	// ffiTypeBatch represents the C struct llama_batch
	ffiTypeBatch = ffi.NewType(&ffi.TypeSint32,
		&ffi.TypePointer, &ffi.TypePointer,
		&ffi.TypePointer, &ffi.TypePointer,
		&ffi.TypePointer, &ffi.TypePointer)
)

var (
	// LLAMA_API struct llama_batch llama_batch_init(
	//         int32_t n_tokens,
	batchInitFunc ffi.Fun

	// LLAMA_API void llama_batch_free(struct llama_batch batch);
	batchFreeFunc ffi.Fun

	// LLAMA_API struct llama_batch llama_batch_get_one(
	//               llama_token * tokens,
	//                   int32_t   n_tokens);
	batchGetOneFunc ffi.Fun
)

func loadBatchFuncs(lib ffi.Lib) error {
	var err error

	if batchInitFunc, err = lib.Prep("llama_batch_init", &ffiTypeBatch, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32); err != nil {
		return loadError("llama_batch_init", err)
	}

	if batchFreeFunc, err = lib.Prep("llama_batch_free", &ffi.TypeVoid, &ffiTypeBatch); err != nil {
		return loadError("llama_batch_free", err)
	}

	if batchGetOneFunc, err = lib.Prep("llama_batch_get_one", &ffiTypeBatch, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("llama_batch_get_one", err)
	}

	return nil
}

// BatchInit allocates a batch of tokens on the heap that can hold a maximum of nTokens.
// Each token can be assigned up to nSeqMax sequence ids
// The batch has to be freed with [BatchFree].
// If embd != 0, Batch.embd will be allocated with size of nTokens * embd * sizeof(float)
// Otherwise, Batch.token will be allocated to store nTokens [Token]
// The rest of the Batch members are allocated with size n_tokens
// All members are left uninitialized.
func BatchInit(nTokens int32, embd int32, nSeqMax int32) Batch {
	var batch Batch
	batchInitFunc.Call(unsafe.Pointer(&batch.BatchData), &nTokens, &embd, &nSeqMax)
	batch.capTokens = nTokens
	batch.capSeq = nSeqMax

	// llama.cpp zeroes n_tokens here, but llama.h documents the members as
	// left uninitialized. Add and SetLogit bound themselves against NTokens,
	// so it is set explicitly rather than relying on that.
	batch.NTokens = 0

	return batch
}

// BatchFree frees a Batch of tokens allocated with BatchInit.
func BatchFree(batch Batch) error {
	batchFreeFunc.Call(nil, unsafe.Pointer(&batch.BatchData))

	return nil
}

// BatchGetOne returns Batch for single sequence of tokens.
// The sequence ID will be fixed to 0.
// The position of the tokens will be tracked automatically by [Decode].
//
// The returned batch borrows the caller's tokens and owns no writable arrays of
// its own: llama_batch_get_one leaves pos, n_seq_id, seq_id and logits NULL. It
// can be passed to [Decode] or [Encode], but [Batch.Add] and [Batch.SetLogit]
// report [ErrBatchNotWritable] on it. Use [BatchInit] for a batch you intend to
// fill in yourself.
func BatchGetOne(tokens []Token) Batch {
	var batch Batch
	if len(tokens) == 0 {
		return batch
	}
	toks := unsafe.SliceData(tokens)
	nTokens := int32(len(tokens))

	batchGetOneFunc.Call(unsafe.Pointer(&batch.BatchData), unsafe.Pointer(&toks), &nTokens)

	return batch
}

// Clear resets the token count of the batch to zero.
func (b *Batch) Clear() error {
	b.NTokens = 0

	return nil
}

// writable reports whether the batch owns the arrays Add and SetLogit write to.
// A zero Batch and one from [BatchGetOne] do not, and llama_batch_init leaves
// token NULL when it was asked for an embedding batch, so each pointer has to be
// checked rather than inferred from capTokens alone.
func (b *Batch) writable() bool {
	return b.capTokens > 0 &&
		b.Token != nil && b.Pos != nil && b.NSeqId != nil && b.SeqId != nil && b.Logits != nil
}

// SetLogit sets whether to compute logits for the token at index idx in the batch.
//
// idx must address a token already in the batch, that is, it must be less than
// NTokens. llama_decode only reads logit flags over [0, n_tokens), so a flag set
// beyond that would be silently dropped rather than honoured; the tighter bound
// reports the caller's index error instead. Writes are also confined to the
// C-allocated array, which is what makes the bound a safety property and not
// just a diagnostic.
func (b *Batch) SetLogit(idx int32, logits bool) error {
	if !b.writable() {
		return ErrBatchNotWritable
	}
	if idx < 0 || idx >= b.NTokens {
		return fmt.Errorf("%w: index %d not in [0,%d)", ErrBatchIndexRange, idx, b.NTokens)
	}

	logitPtr := &unsafe.Slice((*int8)(b.Logits), int(b.capTokens))[idx]
	if logits {
		*logitPtr = 1
	} else {
		*logitPtr = 0
	}

	return nil
}

// Add adds a token to the batch with the given position, sequence IDs, and logits flag.
//
// It writes nothing and returns an error if the batch is already full
// ([ErrBatchFull]), if seqIDs is longer than the n_seq_max the batch was
// allocated with ([ErrTooManySeqIDs]), or if the batch owns no writable arrays
// ([ErrBatchNotWritable]), each of which would otherwise corrupt memory.
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

	// Set token and position
	unsafe.Slice((*Token)(b.Token), int(b.capTokens))[i] = token
	unsafe.Slice((*Pos)(b.Pos), int(b.capTokens))[i] = pos

	// Set number of sequence IDs
	unsafe.Slice((*int32)(b.NSeqId), int(b.capTokens))[i] = int32(len(seqIDs))

	// Set sequence IDs if present
	seqIDPtrs := unsafe.Slice((**SeqId)(b.SeqId), int(b.capTokens))
	if seqIDPtrs[i] != nil && len(seqIDs) > 0 {
		seqSlice := unsafe.Slice((*SeqId)(seqIDPtrs[i]), len(seqIDs))
		for j, sid := range seqIDs {
			seqSlice[j] = sid
		}
	}

	// SetLogit bounds idx against NTokens, so the count has to reflect the
	// token just written before the flag for it can be set.
	b.NTokens++

	return b.SetLogit(i, logits)
}
