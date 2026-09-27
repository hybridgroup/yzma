//go:build js && wasm

package llamawasm

// llama.cpp can put the sampler in the compute graph, so the backend returns
// the token and the values that produced it. The calls here read that result
// and follow the ones in pkg/llama. This is experimental in llama.cpp.

// SetSampler attaches a sampler to a context sequence, so the backend samples
// while it decodes. It returns false when the backend does not accept the
// sampler, and the program then samples itself with SamplerSample.
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

// GetSampledTokenIth returns the token that the backend sampled for output i.
// It returns TokenNull when the backend sampled no token.
func GetSampledTokenIth(ctx Context, i int32) (Token, error) {
	if !Loaded() {
		return TokenNull, ErrNotLoaded
	}
	if !has("_yzma_get_sampled_token_ith") {
		return TokenNull, ErrNoBackendSampling
	}

	// A token can be -1, so only the bad handle value is an error.
	rc := call("_yzma_get_sampled_token_ith", int(ctx), int(i))
	if rc <= errBadHandle {
		return TokenNull, shimError("_yzma_get_sampled_token_ith", rc)
	}
	return Token(rc), nil
}

// GetSampledProbsCountIth returns the number of probabilities the backend
// sampler produced for output i.
func GetSampledProbsCountIth(ctx Context, i int32) (int32, error) {
	return sampledCount("_yzma_get_sampled_probs_count_ith", ctx, i)
}

// GetSampledLogitsCountIth returns the number of logits the backend sampler
// produced for output i.
func GetSampledLogitsCountIth(ctx Context, i int32) (int32, error) {
	return sampledCount("_yzma_get_sampled_logits_count_ith", ctx, i)
}

// GetSampledCandidatesCountIth returns the number of candidates the backend
// sampler produced for output i.
func GetSampledCandidatesCountIth(ctx Context, i int32) (int32, error) {
	return sampledCount("_yzma_get_sampled_candidates_count_ith", ctx, i)
}

// GetSampledProbsIth returns the probabilities for output i. The n argument
// is GetSampledProbsCountIth of the same output.
func GetSampledProbsIth(ctx Context, i, n int32) ([]float32, error) {
	return sampledFloats("_yzma_get_sampled_probs_ith", ctx, i, n)
}

// GetSampledLogitsIth returns the logits for output i. The n argument is
// GetSampledLogitsCountIth of the same output.
func GetSampledLogitsIth(ctx Context, i, n int32) ([]float32, error) {
	return sampledFloats("_yzma_get_sampled_logits_ith", ctx, i, n)
}

// GetSampledCandidatesIth returns the tokens the backend sampler kept for
// output i. The n argument is GetSampledCandidatesCountIth of the same output.
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

// sampledCount reads one of the three backend sampler counts.
func sampledCount(name string, ctx Context, i int32) (int32, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if !has(name) {
		return 0, ErrNoBackendSampling
	}
	return callErr(name, int(ctx), int(i))
}

// sampledFloats reads one of the two backend sampler float arrays.
func sampledFloats(name string, ctx Context, i, n int32) ([]float32, error) {
	if !Loaded() {
		return nil, ErrNotLoaded
	}
	if !has(name) {
		return nil, ErrNoBackendSampling
	}
	return floats(name, n, int(ctx), int(i))
}
