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

// SamplerInitTypical makes a sampler that keeps the tokens whose surprise is
// near the average.
func SamplerInitTypical(p float32, keep uint32) Sampler {
	return newSampler("_yzma_sampler_typical", float64(p), int(keep))
}

// SamplerInitXTC makes a sampler that removes a probable token part of the
// time, which gives more varied text.
func SamplerInitXTC(p float32, t float32, minKeep uint32, seed uint32) Sampler {
	return newSampler("_yzma_sampler_xtc", float64(p), float64(t), int(minKeep), int(seed))
}

// SamplerInitTopNSigma makes a sampler that keeps the tokens whose logit is
// within n standard deviations of the largest one.
func SamplerInitTopNSigma(n float32) Sampler {
	return newSampler("_yzma_sampler_top_n_sigma", float64(n))
}

// SamplerInitTempExt makes a sampler that changes the temperature with the
// entropy of the distribution.
func SamplerInitTempExt(t float32, delta float32, exponent float32) Sampler {
	return newSampler("_yzma_sampler_temp_ext", float64(t), float64(delta), float64(exponent))
}

// SamplerInitMirostat makes a sampler that holds the surprise of the text near
// tau.
func SamplerInitMirostat(nVocab int32, seed uint32, tau, eta float32, m int32) Sampler {
	return newSampler("_yzma_sampler_mirostat", int(nVocab), int(seed), float64(tau),
		float64(eta), int(m))
}

// SamplerInitMirostatV2 makes a sampler that holds the surprise of the text
// near tau, and needs no vocabulary.
func SamplerInitMirostatV2(seed uint32, tau, eta float32) Sampler {
	return newSampler("_yzma_sampler_mirostat_v2", int(seed), float64(tau), float64(eta))
}

// SamplerInitAdaptiveP makes a sampler that takes the tokens whose probability
// is near a target.
func SamplerInitAdaptiveP(target float32, decay float32, seed uint32) Sampler {
	return newSampler("_yzma_sampler_adaptive_p", float64(target), float64(decay), int(seed))
}

// SamplerInitInfill makes a sampler for a fill in the middle prompt. Put it
// after the top-k and top-p samplers.
func SamplerInitInfill(vocab Vocab) Sampler {
	return newSampler("_yzma_sampler_infill", int(vocab))
}

// SamplerInitGrammar makes a sampler that permits only the text that a GBNF
// grammar describes. Give the name of the first rule in root, which is usually
// "root".
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

// SamplerInitGrammarLazyPatterns makes a grammar sampler that starts only after
// a pattern or a token of the trigger appears.
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

// SamplerInitDry makes a DRY sampler, which lowers the probability of a
// sequence that came before. A penaltyLast below 0 becomes 0 in llama.cpp, thus
// give the size of the context for the whole history.
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

// SamplerInitLogitBias makes a sampler that moves the logit of each token in
// tokens by the value at the same position in biases.
//
// The signature is not the same as llama.SamplerInitLogitBias, which takes a
// pointer to an array of llama.LogitBias. A struct cannot cross the boundary of
// the module, thus this takes two slices. They must have the same length.
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

// SamplerName gives the name of a sampler.
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

// SamplerGetSeed gives the seed of a sampler. The result is 0xFFFFFFFF for a
// sampler that has no seed.
func SamplerGetSeed(smpl Sampler) uint32 {
	if !has("_yzma_sampler_get_seed") {
		return 0xFFFFFFFF
	}
	return uint32(call("_yzma_sampler_get_seed", int(smpl)))
}

// SamplerClone makes a copy of a sampler. The caller owns the copy and must
// free it or put it in a chain.
func SamplerClone(smpl Sampler) Sampler {
	if !has("_yzma_sampler_clone") {
		return 0
	}
	return newSampler("_yzma_sampler_clone", int(smpl))
}

// SamplerChainN gives the number of samplers in a chain.
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

// SamplerChainGet gives the sampler at position i of a chain. An i of -1 gives
// the chain itself.
//
// The chain keeps the sampler. Thus SamplerFree of the result frees nothing,
// and the result becomes stale when the chain goes away.
func SamplerChainGet(chain Sampler, i int32) Sampler {
	if !has("_yzma_sampler_chain_get") {
		return 0
	}
	return newSampler("_yzma_sampler_chain_get", int(chain), int(i))
}

// SamplerChainRemove takes the sampler at position i out of a chain. The caller
// then owns the sampler and must free it.
func SamplerChainRemove(chain Sampler, i int32) Sampler {
	if !has("_yzma_sampler_chain_remove") {
		return 0
	}
	return newSampler("_yzma_sampler_chain_remove", int(chain), int(i))
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
