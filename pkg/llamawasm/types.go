//go:build js && wasm

package llamawasm

// The shim gives each llama.cpp object a small int32 handle. A handle of 0 is
// never valid. The Go code holds no address in the llama.cpp module, so the
// module can move or grow its memory.
type (
	// Token is one token in a vocabulary.
	Token int32

	// Pos is the position of a token in a sequence.
	Pos int32

	// SeqId is a sequence ID.
	SeqId int32

	// Model is a handle to a loaded model.
	Model int32

	// Context is a handle to an inference context.
	Context int32

	// Vocab is a handle to the vocabulary of a model.
	Vocab int32

	// Sampler is a handle to a sampler or a sampler chain.
	Sampler int32
)

// TokenNull is the value for a missing token. It matches llama.TokenNull.
const TokenNull Token = -1

// Ftype is the quantization type of a model. The values match llama.Ftype.
type Ftype int32

const (
	FtypeAllF32          Ftype = 0
	FtypeMostlyF16       Ftype = 1
	FtypeMostlyQ4_0      Ftype = 2
	FtypeMostlyQ4_1      Ftype = 3
	FtypeMostlyQ8_0      Ftype = 7
	FtypeMostlyQ5_0      Ftype = 8
	FtypeMostlyQ5_1      Ftype = 9
	FtypeMostlyQ2_K      Ftype = 10
	FtypeMostlyQ3_K_S    Ftype = 11
	FtypeMostlyQ3_K_M    Ftype = 12
	FtypeMostlyQ3_K_L    Ftype = 13
	FtypeMostlyQ4_K_S    Ftype = 14
	FtypeMostlyQ4_K_M    Ftype = 15
	FtypeMostlyQ5_K_S    Ftype = 16
	FtypeMostlyQ5_K_M    Ftype = 17
	FtypeMostlyQ6_K      Ftype = 18
	FtypeMostlyIQ2_XXS   Ftype = 19
	FtypeMostlyIQ2_XS    Ftype = 20
	FtypeMostlyQ2_K_S    Ftype = 21
	FtypeMostlyIQ3_XS    Ftype = 22
	FtypeMostlyIQ3_XXS   Ftype = 23
	FtypeMostlyIQ1_S     Ftype = 24
	FtypeMostlyIQ4_NL    Ftype = 25
	FtypeMostlyIQ3_S     Ftype = 26
	FtypeMostlyIQ3_M     Ftype = 27
	FtypeMostlyIQ2_S     Ftype = 28
	FtypeMostlyIQ2_M     Ftype = 29
	FtypeMostlyIQ4_XS    Ftype = 30
	FtypeMostlyIQ1_M     Ftype = 31
	FtypeMostlyBF16      Ftype = 32
	FtypeMostlyTQ1_0     Ftype = 36
	FtypeMostlyTQ2_0     Ftype = 37
	FtypeMostlyMXFP4_MOE Ftype = 38
	FtypeMostlyNVFP4     Ftype = 39
	FtypeMostlyQ1_0      Ftype = 40
	FtypeMostlyQ2_0      Ftype = 41
	FtypeGUESSED         Ftype = 1024
)

// ModelMetaKey is a model metadata key that llama.cpp knows. The values match
// enum llama_model_meta_key.
type ModelMetaKey int32

const (
	ModelMetaKeySamplingSequence ModelMetaKey = iota
	ModelMetaKeySamplingTopK
	ModelMetaKeySamplingTopP
	ModelMetaKeySamplingMinP
	ModelMetaKeySamplingXTCProb
	ModelMetaKeySamplingXTCThold
	ModelMetaKeySamplingTemp
	ModelMetaKeySamplingPenaltyLastN
	ModelMetaKeySamplingPenaltyRepeat
	ModelMetaKeySamplingMirostat
	ModelMetaKeySamplingMirostatTau
	ModelMetaKeySamplingMirostatEta
)

// RopeScalingType is how a model scales RoPE positions. The values match
// llama.RopeScalingType.
type RopeScalingType int32

const (
	RopeScalingTypeUnspecified RopeScalingType = -1
	RopeScalingTypeNone        RopeScalingType = 0
	RopeScalingTypeLinear      RopeScalingType = 1
	RopeScalingTypeYARN        RopeScalingType = 2
	RopeScalingTypeLongROPE    RopeScalingType = 3
	RopeScalingTypeMaxValue    RopeScalingType = RopeScalingTypeLongROPE
)

// VocabType is the tokenizer type of a vocabulary. The values match
// llama.VocabType.
type VocabType int32

const (
	VocabTypeNone VocabType = iota
	VocabTypeSPM
	VocabTypeBPE
	VocabTypeWPM
	VocabTypeUGM
	VocabTypeRWKV
	VocabTypePLAMO2
)

// TokenAttr holds the attributes of one token. The values match
// llama.TokenAttr.
type TokenAttr int32

const (
	TokenAttrUndefined  TokenAttr = 0
	TokenAttrUnknown    TokenAttr = 1 << 0
	TokenAttrUnused     TokenAttr = 1 << 1
	TokenAttrNormal     TokenAttr = 1 << 2
	TokenAttrControl    TokenAttr = 1 << 3
	TokenAttrUserDef    TokenAttr = 1 << 4
	TokenAttrByte       TokenAttr = 1 << 5
	TokenAttrNormalized TokenAttr = 1 << 6
	TokenAttrLstrip     TokenAttr = 1 << 7
	TokenAttrRstrip     TokenAttr = 1 << 8
	TokenAttrSingleWord TokenAttr = 1 << 9
)

// PoolingType is how the token embeddings of a sequence combine into one
// embedding.
type PoolingType int32

const (
	PoolingTypeUnspecified PoolingType = -1
	PoolingTypeNone        PoolingType = 0
	PoolingTypeMean        PoolingType = 1
	PoolingTypeCLS         PoolingType = 2
	PoolingTypeLast        PoolingType = 3
	PoolingTypeRank        PoolingType = 4
)

// ModelParams holds the settings the shim can apply when it loads a model.
//
// The struct is much smaller than llama.ModelParams, because a WebAssembly
// build has no device list.
type ModelParams struct {
	// NGpuLayers is the number of layers to put on the GPU. A value larger than
	// the model layer count puts them all there.
	//
	// A CPU build ignores this value. Check [GPUDevice] first, because it is
	// empty when llama.cpp has no GPU.
	NGpuLayers int32
}

// ModelDefaultParams returns the default model parameters.
func ModelDefaultParams() ModelParams {
	return ModelParams{
		NGpuLayers: 0,
	}
}

// ContextParams holds the settings the shim can apply when it creates a context.
type ContextParams struct {
	NCtx        uint32      // size of the text context, 0 = from the model
	NBatch      uint32      // largest logical batch, 0 = from llama.cpp
	NUbatch     uint32      // largest physical batch, 0 = from llama.cpp
	NSeqMax     uint32      // largest number of sequences, 0 = one
	NThreads    int32       // number of threads, 0 = from the module
	PoolingType PoolingType // how to pool embeddings
	Embeddings  uint8       // 1 to compute embeddings
	NoPerf      uint8       // 1 to skip timing each batch
	KVUnified   uint8       // 1 for sequences to share one cache, ABI 10 and later
	NOutputsMax uint32      // largest number of outputs of a batch, 0 = NBatch, ABI 10 and later
}

// ContextDefaultParams returns the default context parameters.
//
// NThreads comes from [Threads], which is the number the module can use. A
// value of 0 gives the llama.cpp default of four threads, which is slow on a
// machine with more cores.
func ContextDefaultParams() ContextParams {
	return ContextParams{
		NCtx:        0,
		NBatch:      0,
		NUbatch:     0,
		NSeqMax:     0,
		NThreads:    Threads(),
		PoolingType: PoolingTypeUnspecified,
		Embeddings:  0,
		NoPerf:      1,
	}
}

// SamplerChainParams holds the parameters of a sampler chain.
type SamplerChainParams struct {
	NoPerf uint8 // 1 to skip timing each sample
}

// SamplerChainDefaultParams returns the default sampler chain parameters.
func SamplerChainDefaultParams() SamplerChainParams {
	return SamplerChainParams{NoPerf: 1}
}

// Batch holds the tokens of one call to Decode or Encode.
//
// On a native platform llama.Batch is a C struct. Here it is an ordinary Go
// struct, because the shim builds the C batch itself.
//
// [BatchGetOne] returns a batch that holds only tokens, and the context assigns
// their positions. [BatchInit] returns a batch that also holds the position,
// sequences, and logit flag of each token, which a program needs to keep more
// than one sequence in one context.
type Batch struct {
	// NTokens is the number of tokens in the batch. A generation loop reads
	// it to move the position forward.
	NTokens int32

	tokens []Token

	// These are empty in a batch from BatchGetOne. seqIDs holds capSeq IDs
	// for each token, one after the other, because that is the layout the
	// shim expects.
	pos       []Pos
	nSeqID    []int32
	seqIDs    []SeqId
	logits    []int8
	capTokens int32
	capSeq    int32
}
