//go:build js && wasm

package llamawasm

// InitFromModel makes an inference context for a model.
func InitFromModel(model Model, params ContextParams) (Context, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}

	// A module before ABI version 6 makes a context of one sequence only.
	if params.NSeqMax > 1 && !has("_yzma_context_new_seq") {
		return 0, ErrNoBatch
	}

	args := []any{
		int(model),
		int(params.NCtx),
		int(params.NBatch),
		int(params.NUbatch),
		int(params.NThreads),
		int(params.Embeddings),
		int(params.PoolingType),
	}

	name := "_yzma_context_new"
	if has("_yzma_context_new_seq") {
		name = "_yzma_context_new_seq"
		args = append(args, int(params.NSeqMax))
	}

	handle, err := callErr(name, args...)
	if err != nil {
		return 0, err
	}
	return Context(handle), nil
}

// Free frees a context.
func Free(ctx Context) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	callVoid("_yzma_context_free", int(ctx))
	return nil
}

// NCtx gives the size of the context.
func NCtx(ctx Context) uint32 {
	if !Loaded() {
		return 0
	}
	n := call("_yzma_context_n_ctx", int(ctx))
	if n < 0 {
		return 0
	}
	return uint32(n)
}

// NBatch gives the largest logical batch of the context. A module before ABI
// version 6 gives 0.
func NBatch(ctx Context) uint32 {
	return contextSize(ctx, "_yzma_context_n_batch")
}

// NUBatch gives the largest physical batch of the context.
func NUBatch(ctx Context) uint32 {
	return contextSize(ctx, "_yzma_context_n_ubatch")
}

// NSeqMax gives the largest number of sequences that the context holds.
func NSeqMax(ctx Context) uint32 {
	return contextSize(ctx, "_yzma_context_n_seq_max")
}

// NCtxSeq gives the size of the context of one sequence.
func NCtxSeq(ctx Context) uint32 {
	return contextSize(ctx, "_yzma_context_n_ctx_seq")
}

// contextSize reads a size of the context. It gives 0 if the module has no
// such call and if the call fails.
func contextSize(ctx Context, name string) uint32 {
	if !has(name) {
		return 0
	}
	n := call(name, int(ctx))
	if n < 0 {
		return 0
	}
	return uint32(n)
}

// Decode runs a batch of tokens through the model.
//
// A batch from [BatchGetOne] takes its positions from the state of the
// context, the same as llama.BatchGetOne with llama.Decode. A batch from
// [BatchInit] carries the position, the sequences, and the logit flag of each
// token, and the module must be of ABI version 6 or later for it.
func Decode(ctx Context, batch Batch) (int32, error) {
	return run(ctx, batch, "_yzma_decode")
}

// Encode runs a batch of tokens through an encoder model.
func Encode(ctx Context, batch Batch) (int32, error) {
	return run(ctx, batch, "_yzma_encode")
}

// run sends a batch to the shim. The name is the call for a batch of tokens
// only, and the call for a batch with positions is the same name with _batch.
func run(ctx Context, batch Batch, name string) (int32, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if batch.NTokens == 0 {
		return 0, nil
	}

	n := int(batch.NTokens)

	ptr, err := tokenScratch.reserve(n * 4)
	if err != nil {
		return 0, err
	}
	writeTokens(ptr, batch.tokens[:n])

	if !batch.writable() {
		return callErr(name, int(ctx), ptr, n)
	}

	if !has(name + "_batch") {
		return 0, ErrNoBatch
	}

	posPtr, err := posScratch.reserve(n * 4)
	if err != nil {
		return 0, err
	}
	writeInt32s(posPtr, batch.pos[:n])

	nSeqPtr, err := nSeqScratch.reserve(n * 4)
	if err != nil {
		return 0, err
	}
	writeInt32s(nSeqPtr, batch.nSeqID[:n])

	// The identifiers of every token go in one array, thus the shim makes the
	// array of pointers that llama_batch needs.
	seqCount := n * int(batch.capSeq)
	seqPtr, err := seqScratch.reserve(seqCount * 4)
	if err != nil {
		return 0, err
	}
	writeInt32s(seqPtr, batch.seqIDs[:seqCount])

	logitPtr, err := logitScratch.reserve(n)
	if err != nil {
		return 0, err
	}
	writeInt8s(logitPtr, batch.logits[:n])

	return callErr(name+"_batch", int(ctx), ptr, posPtr, nSeqPtr, seqPtr, int(batch.capSeq), logitPtr, n)
}

// GetEmbeddingsSeq gives the embedding of a sequence. The n argument is the
// number of values to read, which is ModelNEmbd of the model.
func GetEmbeddingsSeq(ctx Context, seqID SeqId, n int32) ([]float32, error) {
	if !Loaded() {
		return nil, ErrNotLoaded
	}
	if n <= 0 {
		return nil, nil
	}

	ptr, err := embdScratch.reserve(int(n) * 4)
	if err != nil {
		return nil, err
	}

	if _, err := callErr("_yzma_get_embeddings_seq", int(ctx), int(seqID), ptr, int(n)); err != nil {
		return nil, err
	}
	return readFloats(ptr, int(n)), nil
}
