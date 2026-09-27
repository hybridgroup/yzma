//go:build js && wasm

package llamawasm

// SamplerChainInit creates a sampler chain.
func SamplerChainInit(params SamplerChainParams) Sampler {
	if !Loaded() {
		return 0
	}
	return Sampler(call("_yzma_sampler_chain_new", int(params.NoPerf)))
}

// SamplerChainAdd appends a sampler to a chain.
//
// The chain then owns the sampler. So the sampler handle is no longer valid
// and SamplerFree on it does nothing. Freeing the chain frees every sampler
// in it.
func SamplerChainAdd(chain Sampler, smpl Sampler) {
	if !Loaded() {
		return
	}
	callVoid("_yzma_sampler_chain_add", int(chain), int(smpl))
}

// SamplerInitGreedy creates a sampler that always picks the most probable
// token.
func SamplerInitGreedy() Sampler {
	return newSampler("_yzma_sampler_greedy")
}

// SamplerInitDist creates a sampler that picks a token at random, weighted by
// probability. A seed of 0xFFFFFFFF uses a random seed.
func SamplerInitDist(seed uint32) Sampler {
	return newSampler("_yzma_sampler_dist", int(seed))
}

// SamplerInitTemp creates a sampler that reshapes the distribution. A value
// below 1.0 makes the output more confident, and a value above 1.0 makes it
// more varied.
func SamplerInitTemp(t float32) Sampler {
	return newSampler("_yzma_sampler_temp", float64(t))
}

// SamplerInitTopK creates a sampler that keeps only the k most probable tokens.
func SamplerInitTopK(k int32) Sampler {
	return newSampler("_yzma_sampler_top_k", int(k))
}

// SamplerInitTopP creates a sampler that keeps the most probable tokens up to a
// total probability of p.
func SamplerInitTopP(p float32, keep uint32) Sampler {
	return newSampler("_yzma_sampler_top_p", float64(p), int(keep))
}

// SamplerInitMinP creates a sampler that removes every token whose probability
// is below p times the probability of the most likely token.
func SamplerInitMinP(p float32, keep uint32) Sampler {
	return newSampler("_yzma_sampler_min_p", float64(p), int(keep))
}

// SamplerInitPenalties creates a sampler that lowers the probability of
// tokens that already appeared.
func SamplerInitPenalties(nVocab int32, lastN int32, repeat float32, freq float32, present float32) Sampler {
	return newSampler("_yzma_sampler_penalties", int(nVocab), int(lastN),
		float64(repeat), float64(freq), float64(present))
}

// SamplerInitTypical creates a sampler that keeps the tokens whose surprise is
// near the average.
func SamplerInitTypical(p float32, keep uint32) Sampler {
	return newSampler("_yzma_sampler_typical", float64(p), int(keep))
}

// SamplerInitXTC creates a sampler that sometimes removes a probable token,
// which produces more varied text.
func SamplerInitXTC(p float32, t float32, minKeep uint32, seed uint32) Sampler {
	return newSampler("_yzma_sampler_xtc", float64(p), float64(t), int(minKeep), int(seed))
}

// SamplerInitTopNSigma creates a sampler that keeps the tokens whose logit is
// within n standard deviations of the largest one.
func SamplerInitTopNSigma(n float32) Sampler {
	return newSampler("_yzma_sampler_top_n_sigma", float64(n))
}

// SamplerInitTempExt creates a sampler that adjusts the temperature based on
// the entropy of the distribution.
func SamplerInitTempExt(t float32, delta float32, exponent float32) Sampler {
	return newSampler("_yzma_sampler_temp_ext", float64(t), float64(delta), float64(exponent))
}

// SamplerInitMirostat creates a sampler that keeps the surprise of the text
// near tau.
func SamplerInitMirostat(nVocab int32, seed uint32, tau, eta float32, m int32) Sampler {
	return newSampler("_yzma_sampler_mirostat", int(nVocab), int(seed), float64(tau),
		float64(eta), int(m))
}

// SamplerInitMirostatV2 creates a sampler that keeps the surprise of the text
// near tau, and needs no vocabulary.
func SamplerInitMirostatV2(seed uint32, tau, eta float32) Sampler {
	return newSampler("_yzma_sampler_mirostat_v2", int(seed), float64(tau), float64(eta))
}

// SamplerInitAdaptiveP creates a sampler that picks tokens whose probability
// is near a target.
func SamplerInitAdaptiveP(target float32, decay float32, seed uint32) Sampler {
	return newSampler("_yzma_sampler_adaptive_p", float64(target), float64(decay), int(seed))
}

// SamplerInitInfill creates a sampler for a fill in the middle prompt. Put it
// after the top-k and top-p samplers.
func SamplerInitInfill(vocab Vocab) Sampler {
	return newSampler("_yzma_sampler_infill", int(vocab))
}

// SamplerInitGrammar creates a sampler that only allows text matching a GBNF
// grammar. Pass the name of the start rule in root, which is usually "root".
func SamplerInitGrammar(vocab Vocab, grammar, root string) Sampler {
	if !has("_yzma_sampler_grammar") {
		return 0
	}

	grammarPtr, freeGrammar, err := allocString(grammar)
	if err != nil {
		return 0
	}
	defer freeGrammar()

	rootPtr, freeRoot, err := allocString(root)
	if err != nil {
		return 0
	}
	defer freeRoot()

	return newSampler("_yzma_sampler_grammar", int(vocab), grammarPtr, rootPtr)
}

// SamplerInitGrammarLazyPatterns creates a grammar sampler that only starts
// after a trigger pattern or token appears.
func SamplerInitGrammarLazyPatterns(vocab Vocab, grammar, root string,
	triggerPatterns []string, triggerTokens []Token) Sampler {
	if !has("_yzma_sampler_grammar_lazy") {
		return 0
	}

	grammarPtr, freeGrammar, err := allocString(grammar)
	if err != nil {
		return 0
	}
	defer freeGrammar()

	rootPtr, freeRoot, err := allocString(root)
	if err != nil {
		return 0
	}
	defer freeRoot()

	patternsPtr, freePatterns, err := allocStrings(triggerPatterns)
	if err != nil {
		return 0
	}
	defer freePatterns()

	tokensPtr, freeTokens, err := allocTokens(triggerTokens)
	if err != nil {
		return 0
	}
	defer freeTokens()

	return newSampler("_yzma_sampler_grammar_lazy", int(vocab), grammarPtr, rootPtr,
		patternsPtr, len(triggerPatterns), tokensPtr, len(triggerTokens))
}

// SamplerInitDry creates a DRY sampler, which lowers the probability of a
// repeated sequence. llama.cpp turns a negative penaltyLast into 0, so pass
// the context size to cover the whole history.
func SamplerInitDry(vocab Vocab, multiplier float32, base float32, allowedLength int32,
	penaltyLast int32, seqBreakers []string) Sampler {
	if !has("_yzma_sampler_dry") {
		return 0
	}

	breakersPtr, freeBreakers, err := allocStrings(seqBreakers)
	if err != nil {
		return 0
	}
	defer freeBreakers()

	return newSampler("_yzma_sampler_dry", int(vocab), float64(multiplier), float64(base),
		int(allowedLength), int(penaltyLast), breakersPtr, len(seqBreakers))
}

// SamplerInitLogitBias creates a sampler that shifts the logit of each token in
// tokens by the value at the same index in biases.
//
// The signature differs from llama.SamplerInitLogitBias, which takes a pointer
// to an array of llama.LogitBias. A struct cannot cross the module boundary,
// so this takes two slices of the same length.
func SamplerInitLogitBias(nVocab int32, tokens []Token, biases []float32) Sampler {
	if !has("_yzma_sampler_logit_bias") || len(tokens) != len(biases) {
		return 0
	}

	tokensPtr, freeTokens, err := allocTokens(tokens)
	if err != nil {
		return 0
	}
	defer freeTokens()

	biasesPtr, freeBiases, err := allocFloats(biases)
	if err != nil {
		return 0
	}
	defer freeBiases()

	return newSampler("_yzma_sampler_logit_bias", int(nVocab), tokensPtr, biasesPtr, len(tokens))
}

// SamplerName returns the name of a sampler.
func SamplerName(smpl Sampler) string {
	if !has("_yzma_sampler_name") {
		return ""
	}

	const size = 128
	ptr, err := pieceScratch.reserve(size)
	if err != nil {
		return ""
	}

	n := call("_yzma_sampler_name", int(smpl), ptr, size)
	if n <= 0 {
		return ""
	}
	return string(readBytes(ptr, int(n)))
}

// SamplerGetSeed returns the seed of a sampler, or 0xFFFFFFFF for a sampler
// that has no seed.
func SamplerGetSeed(smpl Sampler) uint32 {
	if !has("_yzma_sampler_get_seed") {
		return 0xFFFFFFFF
	}
	return uint32(call("_yzma_sampler_get_seed", int(smpl)))
}

// SamplerClone copies a sampler. The caller owns the copy and must
// free it or put it in a chain.
func SamplerClone(smpl Sampler) Sampler {
	if !has("_yzma_sampler_clone") {
		return 0
	}
	return newSampler("_yzma_sampler_clone", int(smpl))
}

// SamplerChainN returns the number of samplers in a chain.
func SamplerChainN(chain Sampler) int {
	if !has("_yzma_sampler_chain_n") {
		return 0
	}

	n := call("_yzma_sampler_chain_n", int(chain))
	if n < 0 {
		return 0
	}
	return int(n)
}

// SamplerChainGet returns the sampler at position i of a chain. An i of -1
// returns the chain itself.
//
// The chain keeps the sampler. So SamplerFree on the result frees nothing,
// and the result becomes stale when the chain is freed.
func SamplerChainGet(chain Sampler, i int32) Sampler {
	if !has("_yzma_sampler_chain_get") {
		return 0
	}
	return newSampler("_yzma_sampler_chain_get", int(chain), int(i))
}

// SamplerChainRemove removes the sampler at position i from a chain. The caller
// then owns the sampler and must free it.
func SamplerChainRemove(chain Sampler, i int32) Sampler {
	if !has("_yzma_sampler_chain_remove") {
		return 0
	}
	return newSampler("_yzma_sampler_chain_remove", int(chain), int(i))
}

// SamplerSample samples the next token. An idx of -1 uses the logits of the
// last token in the batch.
func SamplerSample(smpl Sampler, ctx Context, idx int32) Token {
	if !Loaded() {
		return -1
	}
	return Token(call("_yzma_sampler_sample", int(smpl), int(ctx), int(idx)))
}

// SamplerAccept passes the selected token to the sampler. Samplers that look
// at previous tokens need this.
func SamplerAccept(smpl Sampler, token Token) {
	if !Loaded() {
		return
	}
	callVoid("_yzma_sampler_accept", int(smpl), int(token))
}

// SamplerReset resets a sampler to its initial state.
func SamplerReset(smpl Sampler) {
	if !Loaded() {
		return
	}
	callVoid("_yzma_sampler_reset", int(smpl))
}

// SamplerFree frees a sampler or a sampler chain.
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
