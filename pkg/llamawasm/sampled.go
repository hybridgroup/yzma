//go:build js && wasm

package llamawasm

// llama.cpp can put the sampler in the compute graph, thus the backend gives
// the token and the values that made it. The calls here read that result and
// follow the ones in pkg/llama. This is experimental in llama.cpp.

// SetSampler attaches a sampler to a sequence of a context, which makes the
// backend sample while it decodes. It gives false when the backend does not
// take the sampler, and the program then samples itself with SamplerSample.
func SetSampler(ctx Context, seqID SeqId, sampler Sampler) (bool, error) {
	if !Loaded() {
		return false, ErrNotLoaded
	}
	if !has("_yzma_set_sampler") {
		return false, ErrNoBackendSampling
	}

	rc, err := callErr("_yzma_set_sampler", int(ctx), int(seqID), int(sampler))
	if err != nil {
		return false, err
	}
	return rc == 1, nil
}

// GetSampledTokenIth gives the token that the backend sampled for the output
// at i. It gives TokenNull when the backend sampled no token.
func GetSampledTokenIth(ctx Context, i int32) (Token, error) {
	if !Loaded() {
		return TokenNull, ErrNotLoaded
	}
	if !has("_yzma_get_sampled_token_ith") {
		return TokenNull, ErrNoBackendSampling
	}

	// A token can be -1, thus only the value of a bad handle is an error.
	rc := call("_yzma_get_sampled_token_ith", int(ctx), int(i))
	if rc <= errBadHandle {
		return TokenNull, shimError("_yzma_get_sampled_token_ith", rc)
	}
	return Token(rc), nil
}

// GetSampledProbsCountIth gives the number of probabilities that the backend
// sampler made for the output at i.
func GetSampledProbsCountIth(ctx Context, i int32) (int32, error) {
	return sampledCount("_yzma_get_sampled_probs_count_ith", ctx, i)
}

// GetSampledLogitsCountIth gives the number of logits that the backend sampler
// made for the output at i.
func GetSampledLogitsCountIth(ctx Context, i int32) (int32, error) {
	return sampledCount("_yzma_get_sampled_logits_count_ith", ctx, i)
}

// GetSampledCandidatesCountIth gives the number of candidates that the backend
// sampler made for the output at i.
func GetSampledCandidatesCountIth(ctx Context, i int32) (int32, error) {
	return sampledCount("_yzma_get_sampled_candidates_count_ith", ctx, i)
}

// GetSampledProbsIth gives the probabilities of the output at i. The n
// argument is GetSampledProbsCountIth of the same output.
func GetSampledProbsIth(ctx Context, i, n int32) ([]float32, error) {
	return sampledFloats("_yzma_get_sampled_probs_ith", ctx, i, n)
}

// GetSampledLogitsIth gives the logits of the output at i. The n argument is
// GetSampledLogitsCountIth of the same output.
func GetSampledLogitsIth(ctx Context, i, n int32) ([]float32, error) {
	return sampledFloats("_yzma_get_sampled_logits_ith", ctx, i, n)
}

// GetSampledCandidatesIth gives the tokens that the backend sampler kept for
// the output at i. The n argument is GetSampledCandidatesCountIth of the same
// output.
func GetSampledCandidatesIth(ctx Context, i, n int32) ([]Token, error) {
	if !Loaded() {
		return nil, ErrNotLoaded
	}
	if !has("_yzma_get_sampled_candidates_ith") {
		return nil, ErrNoBackendSampling
	}
	if n <= 0 {
		return nil, nil
	}

	ptr, err := outScratch.reserve(int(n) * 4)
	if err != nil {
		return nil, err
	}

	if _, err := callErr("_yzma_get_sampled_candidates_ith", int(ctx), int(i), ptr, int(n)); err != nil {
		return nil, err
	}
	return readTokens(ptr, int(n)), nil
}

// sampledCount reads one of the three counts of the backend sampler.
func sampledCount(name string, ctx Context, i int32) (int32, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if !has(name) {
		return 0, ErrNoBackendSampling
	}
	return callErr(name, int(ctx), int(i))
}

// sampledFloats reads one of the two float arrays of the backend sampler.
func sampledFloats(name string, ctx Context, i, n int32) ([]float32, error) {
	if !Loaded() {
		return nil, ErrNotLoaded
	}
	if !has(name) {
		return nil, ErrNoBackendSampling
	}
	return floats(name, n, int(ctx), int(i))
}
