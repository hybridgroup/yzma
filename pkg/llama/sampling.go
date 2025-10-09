package llama

import (
	"math"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

type SamplerType int32

const (
	SamplerTypeNone        SamplerType = iota
	SamplerTypeDry                     = 1
	SamplerTypeTopK                    = 2
	SamplerTypeTopP                    = 3
	SamplerTypeMinP                    = 4
	SamplerTypeTypicalP                = 6
	SamplerTypeTemperature             = 7
	SamplerTypeXTC                     = 8
	SamplerTypeInfill                  = 9
	SamplerTypePenalties               = 10
	SamplerTypeTopNSigma               = 11
)

type Sampler uintptr

var (
	FFISamplerChainParams = ffi.NewType(&ffi.TypePointer)
)

var (
	// LLAMA_API struct llama_sampler_chain_params  llama_sampler_chain_default_params(void);
	samplerChainDefaultParamsFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_chain_init(struct llama_sampler_chain_params params);
	samplerChainInitFunc ffi.Fun

	// LLAMA_API void llama_sampler_chain_add(struct llama_sampler * chain, struct llama_sampler * smpl);
	samplerChainAddFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_greedy(void);
	samplerInitGreedyFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_dist  (uint32_t seed);
	samplerInitDistFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_logit_bias(
	//                  int32_t   n_vocab,
	//                  int32_t   n_logit_bias,
	//   				const llama_logit_bias * logit_bias);
	samplerInitLogitBiasFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_penalties(
	// 						int32_t   penalty_last_n,   // last n tokens to penalize (0 = disable penalty, -1 = context size)
	// 						float   penalty_repeat,   // 1.0 = disabled
	// 						float   penalty_freq,     // 0.0 = disabled
	// 						float   penalty_present); // 0.0 = disabled
	samplerInitPenaltiesFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_dry(
	// 	const struct llama_vocab *  vocab,
	// 						int32_t    n_ctx_train,
	// 						float    dry_multiplier,
	// 						float    dry_base,
	// 						int32_t    dry_allowed_length,
	// 						int32_t    dry_penalty_last_n,
	// 					const char ** seq_breakers,
	// 						size_t    num_breakers);
	samplerInitDryFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_top_n_sigma(float   n);
	samplerInitTopNSigmaFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_top_k      (int32_t k);
	samplerInitTopKFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_typical    (float   p, size_t min_keep);
	samplerInitTypicalFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_top_p      (float   p, size_t min_keep);
	samplerInitTopPFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_min_p      (float   p, size_t min_keep);
	samplerInitMinPFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_xtc        (float   p, float   t,     size_t min_keep, uint32_t seed);
	samplerInitXTCFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_temp_ext   (float   t, float   delta, float exponent);
	samplerInitTempExtFunc ffi.Fun

	// LLAMA_API struct llama_sampler * llama_sampler_init_grammar(
	// 					const struct llama_vocab * vocab,
	//               	const char * grammar_str,
	//               	const char * grammar_root);
	samplerInitGrammarFunc ffi.Fun

	// LLAMA_API llama_token llama_sampler_sample(struct llama_sampler * smpl, struct llama_context * ctx, int32_t idx);
	samplerSampleFunc ffi.Fun

	// LLAMA_API void  llama_sampler_accept(struct llama_sampler * smpl, llama_token token);
	samplerAcceptFunc ffi.Fun

	// LLAMA_API void lama_sampler_free  (struct llama_sampler * smpl);
	samplerFreeFunc ffi.Fun
)

func loadSamplingFuncs(lib ffi.Lib) error {
	var err error
	if samplerChainDefaultParamsFunc, err = lib.Prep("llama_sampler_chain_default_params", &FFISamplerChainParams); err != nil {
		return err
	}

	if samplerChainInitFunc, err = lib.Prep("llama_sampler_chain_init", &ffi.TypePointer, &FFISamplerChainParams); err != nil {
		return err
	}

	if samplerChainAddFunc, err = lib.Prep("llama_sampler_chain_add", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if samplerInitGreedyFunc, err = lib.Prep("llama_sampler_init_greedy", &ffi.TypePointer); err != nil {
		return err
	}

	if samplerInitDistFunc, err = lib.Prep("llama_sampler_init_dist", &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return err
	}

	if samplerInitLogitBiasFunc, err = lib.Prep("llama_sampler_init_logit_bias", &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return err
	}

	if samplerInitPenaltiesFunc, err = lib.Prep("llama_sampler_init_penalties", &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat); err != nil {
		return err
	}

	if samplerInitDryFunc, err = lib.Prep("llama_sampler_init_dry", &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat,
		&ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return err
	}

	if samplerInitTopNSigmaFunc, err = lib.Prep("llama_sampler_init_top_n_sigma", &ffi.TypePointer, &ffi.TypeFloat); err != nil {
		return err
	}

	if samplerInitTopKFunc, err = lib.Prep("llama_sampler_init_top_k", &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return err
	}

	if samplerInitTypicalFunc, err = lib.Prep("llama_sampler_init_typical", &ffi.TypePointer, &ffi.TypeFloat, &ffi.TypeUint32); err != nil {
		return err
	}

	if samplerInitTopPFunc, err = lib.Prep("llama_sampler_init_top_p", &ffi.TypePointer, &ffi.TypeFloat, &ffi.TypeUint32); err != nil {
		return err
	}

	if samplerInitMinPFunc, err = lib.Prep("llama_sampler_init_top_p", &ffi.TypePointer, &ffi.TypeFloat, &ffi.TypeUint32); err != nil {
		return err
	}

	if samplerInitXTCFunc, err = lib.Prep("llama_sampler_init_xtc", &ffi.TypePointer, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeUint32, &ffi.TypeUint32); err != nil {
		return err
	}

	if samplerInitTempExtFunc, err = lib.Prep("llama_sampler_init_temp_ext", &ffi.TypePointer, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeFloat); err != nil {
		return err
	}

	if samplerInitGrammarFunc, err = lib.Prep("llama_sampler_init_grammar", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if samplerSampleFunc, err = lib.Prep("llama_sampler_sample", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return err
	}

	if samplerAcceptFunc, err = lib.Prep("llama_sampler_accept", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return err
	}

	if samplerFreeFunc, err = lib.Prep("llama_sampler_free", &ffi.TypePointer); err != nil {
		return err
	}

	return nil
}

// SamplerChainDefaultParams returns the default parameters to create a new sampling chain.
func SamplerChainDefaultParams() SamplerChainParams {
	var p SamplerChainParams
	samplerChainDefaultParamsFunc.Call(unsafe.Pointer(&p))

	return p
}

// SamplerChainInit initializes a new sampling chain.
func SamplerChainInit(params SamplerChainParams) Sampler {
	var p Sampler
	samplerChainInitFunc.Call(unsafe.Pointer(&p), unsafe.Pointer(&params))

	return p
}

// SamplerChainAdd adds a sampler to a sampling chain.
func SamplerChainAdd(chain Sampler, smpl Sampler) {
	samplerChainAddFunc.Call(nil, unsafe.Pointer(&chain), unsafe.Pointer(&smpl))
}

func SamplerInitGreedy() Sampler {
	var p Sampler
	samplerInitGreedyFunc.Call(unsafe.Pointer(&p))

	return p
}

func SamplerInitDist(seed uint32) Sampler {
	var p Sampler
	samplerInitDistFunc.Call(unsafe.Pointer(&p), unsafe.Pointer(&seed))

	return p
}

func SamplerInitLogitBias(nVocab int32, nLogitBias int32, logitBias *LogitBias) Sampler {
	var p Sampler
	samplerInitLogitBiasFunc.Call(unsafe.Pointer(&p), unsafe.Pointer(&nVocab), unsafe.Pointer(&nLogitBias), unsafe.Pointer(&logitBias))

	return p
}

func SamplerInitPenalties(lastN int32, repeat float32, freq float32, present float32) Sampler {
	var p Sampler
	samplerInitPenaltiesFunc.Call(unsafe.Pointer(&p), unsafe.Pointer(&lastN), unsafe.Pointer(&repeat), unsafe.Pointer(&freq), unsafe.Pointer(&present))

	return p
}

func SamplerInitDry(vocab Vocab, nCtxTrain int32, multiplier float32, base float32, allowedLength int32, penaltyLast int32,
	seqBreakers **byte, numBreakers uint32) Sampler {
	var p Sampler
	samplerInitDryFunc.Call(unsafe.Pointer(&p), unsafe.Pointer(&nCtxTrain), unsafe.Pointer(&multiplier), unsafe.Pointer(&base), unsafe.Pointer(&allowedLength), unsafe.Pointer(&penaltyLast),
		unsafe.Pointer(seqBreakers), unsafe.Pointer(&numBreakers))

	return p
}

func SamplerInitTopNSigma(n float32) Sampler {
	var p Sampler
	samplerInitTopNSigmaFunc.Call(unsafe.Pointer(&p), unsafe.Pointer(&n))

	return p
}

func SamplerInitTopK(k int32) Sampler {
	var p Sampler
	samplerInitTopKFunc.Call(unsafe.Pointer(&p), unsafe.Pointer(&k))

	return p
}

func SamplerInitTypical(p float32, keep uint32) Sampler {
	var s Sampler
	samplerInitTypicalFunc.Call(unsafe.Pointer(&s), unsafe.Pointer(&p), unsafe.Pointer(&keep))

	return s
}

func SamplerInitTopP(p float32, keep uint32) Sampler {
	var s Sampler
	samplerInitTopPFunc.Call(unsafe.Pointer(&s), unsafe.Pointer(&p), unsafe.Pointer(&keep))

	return s
}

func SamplerInitMinP(p float32, keep uint32) Sampler {
	var s Sampler
	samplerInitMinPFunc.Call(unsafe.Pointer(&s), unsafe.Pointer(&p), unsafe.Pointer(&keep))

	return s
}

func SamplerInitXTC(p float32, t float32, minKeep uint32, seed uint32) Sampler {
	var s Sampler
	samplerInitXTCFunc.Call(unsafe.Pointer(&s), unsafe.Pointer(&p), unsafe.Pointer(&t), unsafe.Pointer(&minKeep), unsafe.Pointer(&seed))

	return s
}

func SamplerInitTempExt(t float32, delta float32, exponent float32) Sampler {
	var s Sampler
	samplerInitTempExtFunc.Call(unsafe.Pointer(&s), unsafe.Pointer(&t), unsafe.Pointer(&delta), unsafe.Pointer(&exponent))

	return s
}

func SamplerInitGrammar(vocab Vocab, grammar, root string) Sampler {
	grmr, _ := utils.BytePtrFromString(grammar)
	r, _ := utils.BytePtrFromString(root)

	var s Sampler
	samplerInitGrammarFunc.Call(unsafe.Pointer(&s), unsafe.Pointer(&vocab), unsafe.Pointer(&grmr), unsafe.Pointer(&r))

	return s
}

func SamplerSample(smpl Sampler, ctx Context, idx int32) Token {
	var result ffi.Arg
	samplerSampleFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&smpl), unsafe.Pointer(&ctx), unsafe.Pointer(&idx))

	return Token(result)
}

func SamplerAccept(smpl Sampler, token Token) {
	samplerAcceptFunc.Call(nil, unsafe.Pointer(&smpl), unsafe.Pointer(&token))
}

func SamplerFree(smpl Sampler) {
	samplerFreeFunc.Call(nil, unsafe.Pointer(&smpl))
}

var (
	DefaultSamplers = []SamplerType{
		SamplerTypePenalties,
		SamplerTypeDry,
		SamplerTypeTopNSigma,
		SamplerTypeTopK,
		SamplerTypeTypicalP,
		SamplerTypeTopP,
		SamplerTypeMinP,
		SamplerTypeXTC,
		SamplerTypeTemperature,
	}
)

// NewSampler creates a new sampling chain.
func NewSampler(model Model, samplers []SamplerType) Sampler {
	vocab := ModelGetVocab(model)
	nTokens := VocabNTokens(vocab)

	params := SamplerChainDefaultParams()
	sampler := SamplerChainInit(params)

	logitBiasEOG := make([]LogitBias, 0)

	for i := int32(0); i < nTokens; i++ {
		token := Token(i)
		if VocabIsEOG(vocab, token) {
			logitBiasEOG = append(logitBiasEOG, LogitBias{Token: token, Bias: math.SmallestNonzeroFloat32})
		}
	}

	bias := SamplerInitLogitBias(nTokens, int32(len(logitBiasEOG)), unsafe.SliceData(logitBiasEOG))
	SamplerChainAdd(sampler, bias)

	for samplerType := range samplers {
		switch samplerType {
		case SamplerTypeDry:
			seqBreakers := []string{"\n", ":", "\"", "*"}
			var combined []*byte
			for _, s := range seqBreakers {
				ptr, err := utils.BytePtrFromString(s)
				if err != nil {
					panic(err)
				}
				combined = append(combined, ptr)
			}
			seqBreakersPtr := unsafe.SliceData(combined)

			dry := SamplerInitDry(vocab, ModelNCtxTrain(model), 0, 1.75, 2, 4096, seqBreakersPtr, uint32(len(seqBreakers)))
			SamplerChainAdd(sampler, dry)

		case SamplerTypeTopK:
			topK := SamplerInitTopK(40)
			SamplerChainAdd(sampler, topK)

		case SamplerTypeTopP:
			topP := SamplerInitTopP(0.95, 0)
			SamplerChainAdd(sampler, topP)

		case SamplerTypeMinP:
			minP := SamplerInitMinP(0.05, 0)
			SamplerChainAdd(sampler, minP)

		case SamplerTypeTypicalP:
			typical := SamplerInitTypical(1.0, 0)
			SamplerChainAdd(sampler, typical)

		case SamplerTypeTemperature:
			temp := SamplerInitTempExt(0.2, 0, 1.0)
			SamplerChainAdd(sampler, temp)

		case SamplerTypeXTC:
			xtc := SamplerInitXTC(0, 0.1, 0, DefaultSeed)
			SamplerChainAdd(sampler, xtc)

		case SamplerTypeInfill:
			// TODO: add implementation

		case SamplerTypePenalties:
			penalties := SamplerInitPenalties(64, 1.0, 0, 0)
			SamplerChainAdd(sampler, penalties)

		case SamplerTypeTopNSigma:
			topNSigma := SamplerInitTopNSigma(-1.0)
			SamplerChainAdd(sampler, topNSigma)
		}
	}

	// always add this last
	dist := SamplerInitDist(DefaultSeed)
	SamplerChainAdd(sampler, dist)

	return sampler
}
