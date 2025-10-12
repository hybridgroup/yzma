package llama

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	FFITypeContextParams = ffi.NewType(&ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32, &ffi.TypeUint32,
		&ffi.TypeSint32, &ffi.TypeSint32,
		&ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32,
		&ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat,
		&ffi.TypeUint32, &ffi.TypeFloat,
		&ffi.TypePointer, &ffi.TypePointer,
		&ffi.TypeSint32, &ffi.TypeSint32,
		&ffi.TypePointer, &ffi.TypePointer,
		&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8)
)

var (
	// LLAMA_API struct llama_context_params        llama_context_default_params(void);
	contextDefaultParamsFunc ffi.Fun

	// LLAMA_API void llama_free(struct llama_context * ctx);
	freeFunc ffi.Fun

	// LLAMA_API void llama_set_warmup(struct llama_context * ctx, bool warmup);
	setWarmupFunc ffi.Fun

	// LLAMA_API int32_t llama_encode(
	//         		struct llama_context * ctx,
	//          	struct llama_batch   batch);
	encodeFunc ffi.Fun

	// LLAMA_API int32_t llama_decode(
	// 				struct llama_context * ctx,
	// 				struct llama_batch   batch);
	decodeFunc ffi.Fun

	// LLAMA_API void                           llama_perf_context_reset(      struct llama_context * ctx);
	perfContextResetFunc ffi.Fun

	// LLAMA_API void llama_memory_clear(
	// 				llama_memory_t mem,
	// 				bool data);
	memoryClearFunc ffi.Fun

	// LLAMA_API llama_memory_t llama_get_memory (const struct llama_context * ctx);
	getMemoryFunc ffi.Fun

	// LLAMA_API bool llama_memory_seq_rm(
	//         		llama_memory_t mem,
	//           	llama_seq_id seq_id,
	//              llama_pos p0,
	//              llama_pos p1);
	memorySeqRmFunc ffi.Fun

	// LLAMA_API void llama_synchronize(struct llama_context * ctx);
	synchronizeFunc ffi.Fun

	// LLAMA_API  enum llama_pooling_type   llama_pooling_type(const struct llama_context * ctx); // TODO: rename to llama_get_pooling_type
	poolingTypeFunc ffi.Fun

	// Get the embeddings for the ith token. For positive indices, Equivalent to:
	// llama_get_embeddings(ctx) + ctx->output_ids[i]*n_embd
	// Negative indicies can be used to access embeddings in reverse order, -1 is the last embedding.
	// shape: [n_embd] (1-dimensional)
	// returns NULL for invalid ids.
	// LLAMA_API float * llama_get_embeddings_ith(struct llama_context * ctx, int32_t i);
	getEmbeddingsIthFunc ffi.Fun

	// Get the embeddings for a sequence id
	// Returns NULL if pooling_type is LLAMA_POOLING_TYPE_NONE
	// when pooling_type == LLAMA_POOLING_TYPE_RANK, returns float[n_cls_out] with the rank(s) of the sequence
	// otherwise: float[n_embd] (1-dimensional)
	// LLAMA_API float * llama_get_embeddings_seq(struct llama_context * ctx, llama_seq_id seq_id);
	getEmbeddingsSeqFunc ffi.Fun
)

func loadContextFuncs(lib ffi.Lib) error {
	var err error

	if contextDefaultParamsFunc, err = lib.Prep("llama_context_default_params", &FFITypeContextParams); err != nil {
		return loadError("llama_context_default_params", err)
	}

	if freeFunc, err = lib.Prep("llama_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("llama_free", err)
	}

	if setWarmupFunc, err = lib.Prep("llama_set_warmup", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeUint8); err != nil {
		return loadError("llama_set_warmup", err)
	}

	if encodeFunc, err = lib.Prep("llama_encode", &ffi.TypeSint32, &ffi.TypePointer, &FFITypeBatch); err != nil {
		return loadError("llama_encode", err)
	}

	if decodeFunc, err = lib.Prep("llama_decode", &ffi.TypeSint32, &ffi.TypePointer, &FFITypeBatch); err != nil {
		return loadError("llama_decode", err)
	}

	if perfContextResetFunc, err = lib.Prep("llama_perf_context_reset", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("llama_perf_context_reset", err)
	}

	if memoryClearFunc, err = lib.Prep("llama_memory_clear", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeUint8); err != nil {
		return loadError("llama_memory_clear", err)
	}

	if getMemoryFunc, err = lib.Prep("llama_get_memory", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_get_memory", err)
	}

	if memorySeqRmFunc, err = lib.Prep("llama_memory_seq_rm", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32); err != nil {
		return loadError("llama_memory_seq_rm", err)
	}

	if synchronizeFunc, err = lib.Prep("llama_synchronize", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("llama_synchronize", err)
	}

	if poolingTypeFunc, err = lib.Prep("llama_pooling_type", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_pooling_type", err)
	}

	if getEmbeddingsIthFunc, err = lib.Prep("llama_get_embeddings_ith", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("llama_get_embeddings_ith", err)
	}

	if getEmbeddingsSeqFunc, err = lib.Prep("llama_get_embeddings_seq", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("llama_get_embeddings_seq", err)
	}

	return nil
}

// ContextDefaultParams returns the default params to initialize a model context.
func ContextDefaultParams() ContextParams {
	var p ContextParams
	contextDefaultParamsFunc.Call(unsafe.Pointer(&p))
	return p
}

// Free frees the resources for a model context.
func Free(ctx Context) {
	freeFunc.Call(nil, unsafe.Pointer(&ctx))
}

// SetWarmup sets the model context warmup mode on or off.
func SetWarmup(ctx Context, warmup bool) {
	setWarmupFunc.Call(nil, unsafe.Pointer(&ctx), &warmup)
}

// Encode encodes a batch of Token.
func Encode(ctx Context, batch Batch) int32 {
	var result ffi.Arg
	encodeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), unsafe.Pointer(&batch))

	return int32(result)
}

// Decode decodes a batch of Token.
func Decode(ctx Context, batch Batch) int32 {
	var result ffi.Arg
	decodeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), unsafe.Pointer(&batch))

	return int32(result)
}

// PerfContextReset resets the performance metrics for the model context.
func PerfContextReset(ctx Context) {
	perfContextResetFunc.Call(nil, unsafe.Pointer(&ctx))
}

// MemoryClear clears the memory contents.
// If data == true, the data buffers will also be cleared together with the metadata.
func MemoryClear(mem Memory, data bool) {
	memoryClearFunc.Call(nil, unsafe.Pointer(&mem), unsafe.Pointer(&data))
}

// MemorySeqRm removes all tokens that belong to the specified sequence and have positions in [p0, p1).
// Returns false if a partial sequence cannot be removed. Removing a whole sequence never fails.
// seqID < 0 : match any sequence
// p0 < 0     : [0,  p1]
// p1 < 0     : [p0, inf)
func MemorySeqRm(mem Memory, seqID SeqId, p0, p1 Pos) bool {
	var result ffi.Arg
	memorySeqRmFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&mem), &seqID, &p0, &p1)

	return result.Bool()
}

// GetMemory returns the current Memory for the Context.
func GetMemory(ctx Context) Memory {
	var mem Memory
	getMemoryFunc.Call(unsafe.Pointer(&mem), unsafe.Pointer(&ctx))

	return mem
}

// Synchronize waits until all computations are finished.
// This is automatically done when using one of the functions that obtains computation results
// and is not necessary to call it explicitly in most cases.
func Synchronize(ctx Context) {
	synchronizeFunc.Call(nil, unsafe.Pointer(&ctx))
}

// GetPoolingType returns the PoolingType for this context.
func GetPoolingType(ctx Context) PoolingType {
	var result ffi.Arg
	poolingTypeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))

	return PoolingType(result)
}

// GetEmbeddingsIth gets the embeddings for the ith token.
func GetEmbeddingsIth(ctx Context, i int32) []float32 {
	var result ffi.Arg
	getEmbeddingsIthFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &i)

	return unsafe.Slice(((*float32)(unsafe.Pointer(&result))), i)
}

// GetEmbeddingsSeq gets the embeddings for this sequence ID.
func GetEmbeddingsSeq(ctx Context, seqID SeqId, i int32) []float32 {
	var result ffi.Arg
	getEmbeddingsSeqFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &seqID)

	return unsafe.Slice(((*float32)(unsafe.Pointer(&result))), i)
}
