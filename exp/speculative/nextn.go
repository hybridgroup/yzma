package speculative

import (
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

// SetEmbeddingsNextN sets if the context keeps the NextN hidden states of the
// next decode.
//
// A masked value of false keeps all rows, indexed by batch position. Use it on
// the target context. A masked value of true keeps only the rows with a set
// logits flag, indexed through the output ids table. Use it on the draft
// context.
//
// The call has no effect when the library does not export the function.
func SetEmbeddingsNextN(ctx llama.Context, value, masked bool) {
	if ctx == 0 || setEmbeddingsNextNFunc.Cif == nil {
		return
	}
	setEmbeddingsNextNFunc.Call(nil, unsafe.Pointer(&ctx), &value, &masked)
}

// GetEmbeddingsNextN gets the NextN hidden states of the last decode. nRows is
// the number of rows that the caller expects, usually Batch.NTokens. nEmbd is
// the embedding width from [llama.ModelNEmbd].
//
// It returns nil when the library does not export the function, or when no
// hidden states are available. This usually means that SetEmbeddingsNextN was
// not enabled before the decode.
//
// The result points to memory that llama.cpp owns. Do not keep it after the
// next decode or synchronize call. Copy the rows that must stay.
func GetEmbeddingsNextN(ctx llama.Context, nRows, nEmbd int) []float32 {
	if ctx == 0 || getEmbeddingsNextNFunc.Cif == nil {
		return nil
	}

	var result *float32
	getEmbeddingsNextNFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))

	if result == nil {
		return nil
	}

	return unsafe.Slice(result, nRows*nEmbd)
}

// GetEmbeddingsNextNIth gets the NextN hidden state row for the ith output of
// the last decode. nEmbd is the embedding width from [llama.ModelNEmbd].
//
// On a masked context, i goes through the output ids table, so it must match a
// batch position with a set logits flag. On an unmasked context, i is the
// batch position.
//
// It returns nil when the library does not export the function, or when the
// row is not available. The result points to memory that llama.cpp owns. Do
// not keep it after the next decode or synchronize call.
func GetEmbeddingsNextNIth(ctx llama.Context, i int32, nEmbd int) []float32 {
	if ctx == 0 || getEmbeddingsNextNIthFunc.Cif == nil {
		return nil
	}

	var result *float32
	getEmbeddingsNextNIthFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &i)

	if result == nil {
		return nil
	}

	return unsafe.Slice(result, nEmbd)
}

// SetNextNLayerOffset selects the appended NextN block that the MTP decode
// graph runs. The offset counts from the end of the trunk layers.
//
// The call has no effect when the library does not export the function.
func SetNextNLayerOffset(ctx llama.Context, offset int32) {
	if ctx == 0 || setNextNLayerOffsetFunc.Cif == nil {
		return
	}
	setNextNLayerOffsetFunc.Call(nil, unsafe.Pointer(&ctx), &offset)
}
