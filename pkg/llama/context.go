package llama

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var FFITypeContextParams = ffi.NewType(
	&ffi.TypeUint32, &ffi.TypeUint32,
	&ffi.TypeUint32, &ffi.TypeUint32,
	&ffi.TypeSint32, &ffi.TypeSint32,
	&ffi.TypeSint32, &ffi.TypeSint32,
	&ffi.TypeSint32, &ffi.TypeSint32,
	&ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeFloat, &ffi.TypeFloat,
	&ffi.TypeUint32, &ffi.TypeFloat,
	&ffi.TypePointer, &ffi.TypePointer,
	&ffi.TypeSint32, &ffi.TypeSint32,
	&ffi.TypePointer, &ffi.TypePointer,
	&ffi.TypeUint8, &ffi.TypeUint8,
	&ffi.TypeUint8, &ffi.TypeUint8,
	&ffi.TypeUint8, &ffi.TypeUint8)

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

	// LLAMA_API llama_memory_t llama_get_memory (const struct llama_context * ctx);
	getMemoryFunc ffi.Fun

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

	// Get all output token embeddings.
	// when pooling_type == LLAMA_POOLING_TYPE_NONE or when using a generative model,
	// the embeddings for which llama_batch.logits[i] != 0 are stored contiguously
	// in the order they have appeared in the batch.
	// shape: [n_outputs*n_embd]
	// Otherwise, returns NULL.
	// TODO: deprecate in favor of llama_get_embeddings_ith() (ref: https://github.com/ggml-org/llama.cpp/pull/14853#issuecomment-3113143522)
	// LLAMA_API float * llama_get_embeddings(struct llama_context * ctx);
	getEmbeddingsFunc ffi.Fun

	// LLAMA_API float * llama_get_logits_ith(struct llama_context * ctx, int32_t i);
	getLogitsIthFunc ffi.Fun

	// Token logits obtained from the last call to llama_decode()
	// The logits for which llama_batch.logits[i] != 0 are stored contiguously
	// in the order they have appeared in the batch.
	// Rows: number of tokens for which llama_batch.logits[i] != 0
	// Cols: n_vocab
	// TODO: deprecate in favor of llama_get_logits_ith() (ref: https://github.com/ggml-org/llama.cpp/pull/14853#issuecomment-3113143522)
	// LLAMA_API float * llama_get_logits(struct llama_context * ctx);
	getLogitsFunc ffi.Fun

	// LLAMA_API uint32_t llama_n_ctx(const struct llama_context * ctx);
	nCtxFunc ffi.Fun

	// LLAMA_API uint32_t llama_n_batch(const struct llama_context * ctx);
	nBatchFunc ffi.Fun

	// LLAMA_API uint32_t llama_n_ubatch(const struct llama_context * ctx);
	nUBatchFunc ffi.Fun

	// LLAMA_API uint32_t llama_n_seq_max(const struct llama_context * ctx);
	nSeqMaxFunc ffi.Fun

	// LLAMA_API const struct llama_model * llama_get_model(const struct llama_context * ctx);
	getModelFunc ffi.Fun

	// LLAMA_API void llama_set_embeddings(struct llama_context * ctx, bool embeddings);
	setEmbeddingsFunc ffi.Fun

	// LLAMA_API void llama_set_causal_attn(struct llama_context * ctx, bool causal_attn);
	setCausalAttnFunc ffi.Fun
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

	if getMemoryFunc, err = lib.Prep("llama_get_memory", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_get_memory", err)
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

	if getEmbeddingsFunc, err = lib.Prep("llama_get_embeddings", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_get_embeddings", err)
	}

	if getLogitsIthFunc, err = lib.Prep("llama_get_logits_ith", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("llama_get_logits_ith", err)
	}

	if getLogitsFunc, err = lib.Prep("llama_get_logits", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_get_logits", err)
	}

	if nCtxFunc, err = lib.Prep("llama_n_ctx", &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("llama_n_ctx", err)
	}

	if nBatchFunc, err = lib.Prep("llama_n_batch", &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("llama_n_batch", err)
	}

	if nUBatchFunc, err = lib.Prep("llama_n_ubatch", &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("llama_n_ubatch", err)
	}

	if nSeqMaxFunc, err = lib.Prep("llama_n_seq_max", &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("llama_n_seq_max", err)
	}

	if getModelFunc, err = lib.Prep("llama_get_model", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_get_model", err)
	}

	if setEmbeddingsFunc, err = lib.Prep("llama_set_embeddings", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeUint8); err != nil {
		return loadError("llama_set_embeddings", err)
	}

	if setCausalAttnFunc, err = lib.Prep("llama_set_causal_attn", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeUint8); err != nil {
		return loadError("llama_set_causal_attn", err)
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
func GetEmbeddingsIth(ctx Context, i int32, nVocab int32) []float32 {
	var result *float32
	getEmbeddingsIthFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &i)

	if result == nil {
		return nil
	}

	return unsafe.Slice(result, nVocab)
}

// GetEmbeddingsSeq gets the embeddings for this sequence ID.
func GetEmbeddingsSeq(ctx Context, seqID SeqId, nVocab int32) []float32 {
	var result *float32
	getEmbeddingsSeqFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &seqID)

	if result == nil {
		return nil
	}

	return unsafe.Slice(result, nVocab)
}

// GetEmbeddings retrieves all output token embeddings.
// Returns a slice of float32 of length nOutputs * nEmbeddings, or nil if not available.
func GetEmbeddings(ctx Context, nOutputs, nEmbeddings int) []float32 {
	var result *float32
	getEmbeddingsFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	if result == nil || nOutputs <= 0 || nEmbeddings <= 0 {
		return nil
	}
	return unsafe.Slice(result, nOutputs*nEmbeddings)
}

// GetLogitsIth retrieves the logits for the ith token.
func GetLogitsIth(ctx Context, i int32, nVocab int) []float32 {
	var logitsPtr *float32
	getLogitsIthFunc.Call(unsafe.Pointer(&logitsPtr), unsafe.Pointer(&ctx), unsafe.Pointer(&i))

	if logitsPtr == nil {
		return nil
	}

	return unsafe.Slice(logitsPtr, nVocab)
}

// GetLogits retrieves all token logits from the last call to llama_decode.
// Returns a slice of float32 of length nTokens * nVocab, or nil if not available.
func GetLogits(ctx Context, nTokens, nVocab int) []float32 {
	var result *float32
	getLogitsFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	if result == nil || nTokens <= 0 || nVocab <= 0 {
		return nil
	}
	return unsafe.Slice(result, nTokens*nVocab)
}

// NCtx returns the number of context tokens.
func NCtx(ctx Context) uint32 {
	var result ffi.Arg
	nCtxFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return uint32(result)
}

// NBatch returns the number of batch tokens.
func NBatch(ctx Context) uint32 {
	var result ffi.Arg
	nBatchFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return uint32(result)
}

// NUBatch returns the number of micro-batch tokens.
func NUBatch(ctx Context) uint32 {
	var result ffi.Arg
	nUBatchFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return uint32(result)
}

// NSeqMax returns the maximum number of sequences.
func NSeqMax(ctx Context) uint32 {
	var result ffi.Arg
	nSeqMaxFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return uint32(result)
}

// GetModel retrieves the model associated with the given context.
func GetModel(ctx Context) Model {
	var model Model
	getModelFunc.Call(unsafe.Pointer(&model), unsafe.Pointer(&ctx))

	return model
}

// SetEmbeddings sets whether the context outputs embeddings or not.
func SetEmbeddings(ctx Context, embeddings bool) {
	setEmbeddingsFunc.Call(nil, unsafe.Pointer(&ctx), &embeddings)
}

// SetCausalAttn sets whether to use causal attention or not.
func SetCausalAttn(ctx Context, causalAttn bool) {
	setCausalAttnFunc.Call(nil, unsafe.Pointer(&ctx), &causalAttn)
}
