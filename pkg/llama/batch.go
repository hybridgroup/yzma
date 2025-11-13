package llama

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	FFITypeBatch = ffi.NewType(&ffi.TypeSint32,
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

	if batchInitFunc, err = lib.Prep("llama_batch_init", &FFITypeBatch, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32); err != nil {
		return loadError("llama_batch_init", err)
	}

	if batchFreeFunc, err = lib.Prep("llama_batch_free", &ffi.TypeVoid, &FFITypeBatch); err != nil {
		return loadError("llama_batch_free", err)
	}

	if batchGetOneFunc, err = lib.Prep("llama_batch_get_one", &FFITypeBatch, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
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
