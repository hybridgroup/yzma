//go:build js && wasm

package llamawasm

// InitFromModel creates an inference context for a model.
func InitFromModel(model Model, params ContextParams) (Context, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}

	// A module before ABI version 6 creates a single sequence context.
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

	// Only ABI version 9 and later modules take NoPerf, and only ABI version
	// 10 and later take KVUnified and NOutputsMax.
	name := "_yzma_context_new"
	switch {
	case has("_yzma_context_new_kv"):
		name = "_yzma_context_new_kv"
		args = append(args, int(params.NSeqMax), int(params.NoPerf), int(params.KVUnified), int(params.NOutputsMax))
	case params.KVUnified != 0:
		return 0, ErrNoKVUnified
	case has("_yzma_context_new_ext"):
		name = "_yzma_context_new_ext"
		args = append(args, int(params.NSeqMax), int(params.NoPerf))
	case has("_yzma_context_new_seq"):
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

// NCtx returns the context size.
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

// NBatch returns the largest logical batch size of the context. A module
// before ABI version 6 returns 0.
func NBatch(ctx Context) uint32 {
	return contextSize(ctx, "_yzma_context_n_batch")
}

// NUBatch returns the largest physical batch size of the context.
func NUBatch(ctx Context) uint32 {
	return contextSize(ctx, "_yzma_context_n_ubatch")
}

// NSeqMax returns the maximum number of sequences the context holds.
func NSeqMax(ctx Context) uint32 {
	return contextSize(ctx, "_yzma_context_n_seq_max")
}

// NCtxSeq returns the context size of one sequence.
func NCtxSeq(ctx Context) uint32 {
	return contextSize(ctx, "_yzma_context_n_ctx_seq")
}

// contextSize reads a context size. It returns 0 if the module has no such
// call or if the call fails.
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

// GetPoolingType returns how the context pools the token embeddings of a
// sequence into one. It returns PoolingTypeUnspecified when the module has no
// such call.
func GetPoolingType(ctx Context) PoolingType {
	if !has("_yzma_context_pooling_type") {
		return PoolingTypeUnspecified
	}

	// A pooling type can be -1, so only the bad handle value is a failure.
	rc := call("_yzma_context_pooling_type", int(ctx))
	if rc <= errBadHandle {
		return PoolingTypeUnspecified
	}
	return PoolingType(rc)
}

// SetEmbeddings sets whether Decode produces embeddings instead of logits.
// ContextParams sets the same thing when the context is created.
func SetEmbeddings(ctx Context, embeddings bool) error {
	return setContextFlag("_yzma_set_embeddings", ctx, embeddings)
}

// SetCausalAttn sets whether attention is causal. An embedding model needs
// every token to attend to every other one, so it uses false.
func SetCausalAttn(ctx Context, causal bool) error {
	return setContextFlag("_yzma_set_causal_attn", ctx, causal)
}

// setContextFlag sends a boolean flag to the context.
func setContextFlag(name string, ctx Context, on bool) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	if !has(name) {
		return ErrNoContextFlags
	}
	_, err := callErr(name, int(ctx), boolToInt(on))
	return err
}

// Synchronize waits for the context computation to finish. An asynchronous
// backend such as WebGPU needs this before timing measurements.
func Synchronize(ctx Context) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	if !has("_yzma_synchronize") {
		return ErrNoContextFlags
	}
	_, err := callErr("_yzma_synchronize", int(ctx))
	return err
}

// Decode runs a batch of tokens through the model.
//
// A batch from [BatchGetOne] takes its positions from the context state, the
// same as llama.BatchGetOne with llama.Decode. A batch from [BatchInit] carries
// the position, sequences, and logit flag of each token, and needs a module of
// ABI version 6 or later.
func Decode(ctx Context, batch Batch) (int32, error) {
	return run(ctx, batch, "_yzma_decode")
}

// Encode runs a batch of tokens through an encoder model.
func Encode(ctx Context, batch Batch) (int32, error) {
	return run(ctx, batch, "_yzma_encode")
}

// run sends a batch to the shim. The name is the call for a tokens only batch,
// and the call for a batch with positions is the same name with _batch.
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

	// All token IDs go in one array, so the shim builds the pointer array
	// that llama_batch needs.
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

// GetEmbeddingsSeq returns the embedding of a sequence. The n argument is the
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

// GetEmbeddingsIth returns the embedding of one token in the last batch. An i
// of -1 selects the last token that has an embedding. The n argument is the
// number of values to read, which is ModelNEmbd of the model.
//
// The values belong to the context and the next Decode replaces them.
func GetEmbeddingsIth(ctx Context, i, n int32) ([]float32, error) {
	return floats("_yzma_get_embeddings_ith", n, int(ctx), int(i))
}

// GetEmbeddings returns the embeddings of every token in the last batch that
// has one. The values are stored one token after another, so the result holds
// nOutputs times nEmbd values.
//
// nEmbd is ModelNEmbd of the model. The values belong to the context and the
// next Decode replaces them.
func GetEmbeddings(ctx Context, nOutputs, nEmbd int32) ([]float32, error) {
	if nOutputs <= 0 {
		return nil, nil
	}
	return floats("_yzma_get_embeddings", nOutputs*nEmbd, int(ctx))
}

// GetLogitsIth returns the logits of one token in the last batch. An i of -1
// selects the last token that has logits. The nVocab argument is the number of
// values to read, which is VocabNTokens of the vocabulary.
//
// The values belong to the context and the next Decode replaces them.
func GetLogitsIth(ctx Context, i, nVocab int32) ([]float32, error) {
	return floats("_yzma_get_logits_ith", nVocab, int(ctx), int(i))
}

// GetLogits returns the logits of every token in the last batch that has them.
// The values are stored one token after another, so the result holds nTokens
// times nVocab values.
//
// nVocab is VocabNTokens of the vocabulary. The values belong to the context
// and the next Decode replaces them.
func GetLogits(ctx Context, nTokens, nVocab int32) ([]float32, error) {
	if nTokens <= 0 {
		return nil, nil
	}
	return floats("_yzma_get_logits", nTokens*nVocab, int(ctx))
}

// floats runs a call that returns a float array and reads the result. The
// count is the number of values, and args come before the pointer and count
// that every such call takes last.
func floats(name string, count int32, args ...any) ([]float32, error) {
	if !Loaded() {
		return nil, ErrNotLoaded
	}
	if !has(name) {
		return nil, ErrNoOutputs
	}
	if count <= 0 {
		return nil, nil
	}

	ptr, err := outScratch.reserve(int(count) * 4)
	if err != nil {
		return nil, err
	}

	if _, err := callErr(name, append(args, ptr, int(count))...); err != nil {
		return nil, err
	}
	return readFloats(ptr, int(count)), nil
}
