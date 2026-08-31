//go:build js && wasm

package llamawasm

// SamplerChainInit makes a chain of samplers.
func SamplerChainInit(params SamplerChainParams) Sampler {
	if !Loaded() {
		return 0
	}
	return Sampler(call("_yzma_sampler_chain_new", int(params.NoPerf)))
}

// SamplerChainAdd puts a sampler at the end of a chain.
//
// The chain then owns the sampler. Thus the handle of the sampler is no longer
// valid and SamplerFree of that handle has no result. When you free the chain,
// you free each sampler in it.
func SamplerChainAdd(chain Sampler, smpl Sampler) {
	if !Loaded() {
		return
	}
	callVoid("_yzma_sampler_chain_add", int(chain), int(smpl))
}

// SamplerInitGreedy makes a sampler that always takes the token with the
// highest probability.
func SamplerInitGreedy() Sampler {
	return newSampler("_yzma_sampler_greedy")
}

// SamplerInitDist makes a sampler that takes a token at random, following the
// probability of each one. A seed of 0xFFFFFFFF takes a random seed.
func SamplerInitDist(seed uint32) Sampler {
	return newSampler("_yzma_sampler_dist", int(seed))
}

// SamplerInitTemp makes a sampler that changes the shape of the distribution.
// A value below 1.0 makes the output more sure, and a value above 1.0 makes it
// more varied.
func SamplerInitTemp(t float32) Sampler {
	return newSampler("_yzma_sampler_temp", float64(t))
}

// SamplerInitTopK makes a sampler that keeps only the k most probable tokens.
func SamplerInitTopK(k int32) Sampler {
	return newSampler("_yzma_sampler_top_k", int(k))
}

// SamplerInitTopP makes a sampler that keeps the most probable tokens up to a
// total probability of p.
func SamplerInitTopP(p float32, keep uint32) Sampler {
	return newSampler("_yzma_sampler_top_p", float64(p), int(keep))
}

// SamplerInitMinP makes a sampler that removes every token whose probability
// is below p of the probability of the most likely token.
func SamplerInitMinP(p float32, keep uint32) Sampler {
	return newSampler("_yzma_sampler_min_p", float64(p), int(keep))
}

// SamplerInitPenalties makes a sampler that lowers the probability of the
// tokens that came before.
func SamplerInitPenalties(nVocab int32, lastN int32, repeat float32, freq float32, present float32) Sampler {
	return newSampler("_yzma_sampler_penalties", int(nVocab), int(lastN),
		float64(repeat), float64(freq), float64(present))
}

// SamplerSample takes the next token. An idx of -1 uses the logits of the last
// token of the batch.
func SamplerSample(smpl Sampler, ctx Context, idx int32) Token {
	if !Loaded() {
		return -1
	}
	return Token(call("_yzma_sampler_sample", int(smpl), int(ctx), int(idx)))
}

// SamplerAccept gives the sampler the token that the program selected. The
// samplers that examine the previous tokens need this.
func SamplerAccept(smpl Sampler, token Token) {
	if !Loaded() {
		return
	}
	callVoid("_yzma_sampler_accept", int(smpl), int(token))
}

// SamplerReset puts a sampler back to its starting state.
func SamplerReset(smpl Sampler) {
	if !Loaded() {
		return
	}
	callVoid("_yzma_sampler_reset", int(smpl))
}

// SamplerFree frees a sampler or a chain of samplers.
func SamplerFree(smpl Sampler) {
	if !Loaded() {
		return
	}
	callVoid("_yzma_sampler_free", int(smpl))
}

func newSampler(name string, args ...any) Sampler {
	if !Loaded() {
		return 0
	}
	h := call(name, args...)
	if h < 0 {
		return 0
	}
	return Sampler(h)
}
