package llama

import (
	"fmt"
	"unsafe"

	"github.com/jupiterrider/ffi"
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
	batchInitFunc.Call(unsafe.Pointer(&batch), &nTokens, &embd, &nSeqMax)
	batch.capTokens = nTokens
	batch.capSeq = nSeqMax

	return batch
}

// BatchFree frees a Batch of tokens allocated with BatchInit.
func BatchFree(batch Batch) error {
	batchFreeFunc.Call(nil, unsafe.Pointer(&batch))

	return nil
}

// BatchGetOne returns Batch for single sequence of tokens.
// The sequence ID will be fixed to 0.
// The position of the tokens will be tracked automatically by [Decode].
func BatchGetOne(tokens []Token) Batch {
	var batch Batch
	if len(tokens) == 0 {
		return batch
	}
	toks := unsafe.SliceData(tokens)
	nTokens := int32(len(tokens))

	batchGetOneFunc.Call(unsafe.Pointer(&batch), unsafe.Pointer(&toks), &nTokens)

	return batch
}

// Clear resets the token count of the batch to zero.
func (b *Batch) Clear() error {
	b.NTokens = 0

	return nil
}

// SetLogit sets whether to compute logits for the token at index idx in the batch.
// It returns an error if idx is outside the batch capacity to avoid writing past
// the C-allocated array.
func (b *Batch) SetLogit(idx int32, logits bool) error {
	if idx < 0 || idx >= b.capTokens {
		return fmt.Errorf("llama: SetLogit index %d out of range [0,%d)", idx, b.capTokens)
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
// It returns an error (without writing) if the batch is already full or if seqIDs is
// longer than the n_seq_max the batch was allocated with, to avoid heap corruption.
func (b *Batch) Add(token Token, pos Pos, seqIDs []SeqId, logits bool) error {
	i := b.NTokens

	if i < 0 || i >= b.capTokens {
		return fmt.Errorf("llama: batch full: cannot add token at index %d (capacity %d)", i, b.capTokens)
	}
	if int32(len(seqIDs)) > b.capSeq {
		return fmt.Errorf("llama: too many sequence IDs %d for token (n_seq_max %d)", len(seqIDs), b.capSeq)
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

	b.NTokens++

	return b.SetLogit(i, logits)
}
