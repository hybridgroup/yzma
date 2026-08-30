//go:build js && wasm

package llamawasm

// InitFromModel makes an inference context for a model.
func InitFromModel(model Model, params ContextParams) (Context, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}

	handle, err := callErr("_yzma_context_new",
		int(model),
		int(params.NCtx),
		int(params.NBatch),
		int(params.NUbatch),
		int(params.NThreads),
		int(params.Embeddings),
		int(params.PoolingType),
	)
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

// Decode runs a batch of tokens through the model. The positions of the tokens
// continue from the state of the context, the same as llama.BatchGetOne with
// llama.Decode.
func Decode(ctx Context, batch Batch) (int32, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if batch.NTokens == 0 {
		return 0, nil
	}

	ptr, err := tokenScratch.reserve(len(batch.tokens) * 4)
	if err != nil {
		return 0, err
	}
	writeTokens(ptr, batch.tokens)

	return callErr("_yzma_decode", int(ctx), ptr, len(batch.tokens))
}

// Encode runs a batch of tokens through an encoder model.
func Encode(ctx Context, batch Batch) (int32, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if batch.NTokens == 0 {
		return 0, nil
	}

	ptr, err := tokenScratch.reserve(len(batch.tokens) * 4)
	if err != nil {
		return 0, err
	}
	writeTokens(ptr, batch.tokens)

	return callErr("_yzma_encode", int(ctx), ptr, len(batch.tokens))
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

// MemoryClear removes everything from the memory of the context, which starts
// the generation again from an empty state.
//
// On a native platform this takes a llama.Memory that comes from
// llama.GetMemory. Here it takes the context, because the shim holds the
// memory itself.
func MemoryClear(ctx Context, data bool) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	clear := 0
	if data {
		clear = 1
	}
	_, err := callErr("_yzma_memory_clear", int(ctx), clear)
	return err
}
