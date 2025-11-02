package llama

// Common types matching llama.cpp
type (
	Token  int32
	Pos    int32
	SeqId  int32
	Memory uintptr
)

// Constants from llama.h
const (
	DefaultSeed = 0xFFFFFFFF
	TokenNull   = -1

	// File magic numbers
	FileMagicGGLA = 0x67676c61
	FileMagicGGSN = 0x6767736e
	FileMagicGGSQ = 0x67677371

	// Session constants
	SessionMagic   = FileMagicGGSN
	SessionVersion = 9

	StateSeqMagic   = FileMagicGGSQ
	StateSeqVersion = 2

	// maximum token value
	MaxToken = 0x7fffffff
)

// Enums
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

type RoPEType int32

const (
	RoPETypeNone   RoPEType = -1
	RoPETypeNorm   RoPEType = 0
	RoPETypeNEOX   RoPEType = 2
	RoPETypeMROPE  RoPEType = 8
	RoPETypeIMROPE RoPEType = 40
	RoPETypeVision RoPEType = 24
)

type TokenType int32

const (
	TokenTypeUndefined TokenType = iota
	TokenTypeNormal
	TokenTypeUnknown
	TokenTypeControl
	TokenTypeUserDefined
	TokenTypeUnused
	TokenTypeByte
)

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
	FtypeGUESSED         Ftype = 1024
)

type RopeScalingType int32

const (
	RopeScalingTypeUnspecified RopeScalingType = -1
	RopeScalingTypeNone        RopeScalingType = 0
	RopeScalingTypeLinear      RopeScalingType = 1
	RopeScalingTypeYARN        RopeScalingType = 2
	RopeScalingTypeLongROPE    RopeScalingType = 4
)

type PoolingType int32

const (
	PoolingTypeUnspecified PoolingType = -1
	PoolingTypeNone        PoolingType = 0
	PoolingTypeMean        PoolingType = 1
	PoolingTypeCLS         PoolingType = 2
	PoolingTypeLast        PoolingType = 3
	PoolingTypeRank        PoolingType = 4
)

type AttentionType int32

const (
	AttentionTypeCausal    AttentionType = 0
	AttentionTypeNonCausal AttentionType = 1
)

type FlashAttentionType int32

const (
	FlashAttentionTypeAuto     FlashAttentionType = -1
	FlashAttentionTypeDisabled FlashAttentionType = 0
	FlashAttentionTypeEnabled  FlashAttentionType = 1
)

type SplitMode int32

const (
	SplitModeNone  SplitMode = 0
	SplitModeLayer SplitMode = 1
	SplitModeRow   SplitMode = 2
)

type GpuBackend int32

const (
	GpuBackendNone   GpuBackend = 0
	GpuBackendCPU    GpuBackend = 1
	GpuBackendCUDA   GpuBackend = 2
	GpuBackendMetal  GpuBackend = 3
	GpuBackendHIP    GpuBackend = 4
	GpuBackendVulkan GpuBackend = 5
	GpuBackendOpenCL GpuBackend = 6
	GpuBackendSYCL   GpuBackend = 7
)

// String returns the string representation of the GPU backend
func (b GpuBackend) String() string {
	switch b {
	case GpuBackendNone:
		return "None"
	case GpuBackendCPU:
		return "CPU"
	case GpuBackendCUDA:
		return "CUDA"
	case GpuBackendMetal:
		return "Metal"
	case GpuBackendHIP:
		return "HIP"
	case GpuBackendVulkan:
		return "Vulkan"
	case GpuBackendOpenCL:
		return "OpenCL"
	case GpuBackendSYCL:
		return "SYCL"
	default:
		return "Unknown"
	}
}

type NumaStrategy int32

const (
	NumaStrategyDisabled   NumaStrategy = 0
	NumaStrategyDistribute NumaStrategy = 1
	NumaStrategyIsolate    NumaStrategy = 2
	NumaStrategyNumactl    NumaStrategy = 3
	NumaStrategyMirror     NumaStrategy = 4
	NumaStrategyCount      NumaStrategy = 5
)

type LogLevel int32

const (
	LogLevelNone     LogLevel = 0
	LogLevelDebug    LogLevel = 1
	LogLevelInfo     LogLevel = 2
	LogLevelWarn     LogLevel = 3
	LogLevelError    LogLevel = 4
	LogLevelContinue LogLevel = 5
)

// Opaque types (represented as pointers)
type (
	Model       uintptr
	Context     uintptr
	Vocab       uintptr
	AdapterLora uintptr
)

// Structs
type TokenData struct {
	Id    Token   // token id
	Logit float32 // log-odds of the token
	P     float32 // probability of the token
}

type TokenDataArray struct {
	Data     *TokenData // pointer to token data array
	Size     uint64     // number of tokens
	Selected int64      // index of selected token (-1 if none)
	Sorted   uint8      // whether the array is sorted by probability (bool as uint8)
}

type Batch struct {
	NTokens int32    // number of tokens
	Token   *Token   // tokens
	Embd    *float32 // embeddings (if using embeddings instead of tokens)
	Pos     *Pos     // positions
	NSeqId  *int32   // number of sequence IDs per token
	SeqId   **SeqId  // sequence IDs
	Logits  *int8    // whether to compute logits for each token
}

// Model parameters
type ModelParams struct {
	Devices                  uintptr   // ggml_backend_dev_t * - NULL-terminated list of devices
	TensorBuftOverrides      uintptr   // const struct llama_model_tensor_buft_override *
	NGpuLayers               int32     // number of layers to store in VRAM
	SplitMode                SplitMode // how to split the model across multiple GPUs
	MainGpu                  int32     // the GPU that is used for the entire model
	TensorSplit              *float32  // proportion of the model to offload to each GPU
	ProgressCallback         uintptr   // llama_progress_callback function pointer
	ProgressCallbackUserData uintptr   // context pointer passed to the progress callback
	KvOverrides              uintptr   // const struct llama_model_kv_override *
	VocabOnly                uint8     // only load the vocabulary, no weights (bool as uint8)
	UseMmap                  uint8     // use mmap if possible (bool as uint8)
	UseMlock                 uint8     // force system to keep model in RAM (bool as uint8)
	CheckTensors             uint8     // validate model tensor data (bool as uint8)
	UseExtraBufts            uint8     // use extra buffer types (bool as uint8)
}

// Context parameters
type ContextParams struct {
	NCtx               uint32             // text context, 0 = from model
	NBatch             uint32             // logical maximum batch size
	NUbatch            uint32             // physical maximum batch size
	NSeqMax            uint32             // max number of sequences
	NThreads           int32              // number of threads to use for generation
	NThreadsBatch      int32              // number of threads to use for batch processing
	RopeScalingType    RopeScalingType    // RoPE scaling type
	PoolingType        PoolingType        // pooling type for embeddings
	AttentionType      AttentionType      // attention type
	FlashAttentionType FlashAttentionType // when to enable Flash Attention
	RopeFreqBase       float32            // RoPE base frequency
	RopeFreqScale      float32            // RoPE frequency scaling factor
	YarnExtFactor      float32            // YaRN extrapolation mix factor
	YarnAttnFactor     float32            // YaRN magnitude scaling factor
	YarnBetaFast       float32            // YaRN low correction dim
	YarnBetaSlow       float32            // YaRN high correction dim
	YarnOrigCtx        uint32             // YaRN original context size
	DefragThold        float32            // defragment the KV cache if holes/size > thold
	CbEval             uintptr            // evaluation callback
	CbEvalUserData     uintptr            // user data for evaluation callback
	TypeK              int32              // data type for K cache
	TypeV              int32              // data type for V cache
	AbortCallback      uintptr            // abort callback
	AbortCallbackData  uintptr            // user data for abort callback
	Embeddings         uint8              // whether to compute and return embeddings (bool as uint8)
	Offload_kqv        uint8              // whether to offload K, Q, V to GPU (bool as uint8)
	NoPerf             uint8              // whether to measure performance (bool as uint8)
	OpOffload          uint8              // offload host tensor operations to device
	SwaFull            uint8              // use full-size SWA cache (https://github.com/ggml-org/llama.cpp/pull/13194#issuecomment-2868343055)
	KVUnified          uint8              // use a unified buffer across the input sequences when computing the attentions
}

// Model quantize parameters
type ModelQuantizeParams struct {
	NThread              int32 // number of threads to use for quantizing
	Ftype                Ftype // quantize to this llama_ftype
	OutputTensorType     int32 // output tensor type
	TokenEmbeddingType   int32 // itoken embeddings tensor type
	AllowRequantize      uint8 // allow quantizing non-f32/f16 tensors (bool as uint8)
	QuantizeOutputTensor uint8 // quantize output.weight (bool as uint8)
	OnlyCopy             uint8 // only copy tensors - ftype, allow_requantize and quantize_output_tensor are ignored
	Pure                 uint8 // quantize all tensors to the default type
	KeepSplit            uint8 // keep split tensors (bool as uint8)
	IMatrix              *byte // pointer to importance matrix data
	KvOverrides          *byte // pointer to vector containing overrides
	TensorTypes          *byte // pointer to vector containing tensor types
	PruneLayers          *byte // pointer to vector containing layer indices to prune
}

// Chat message
type ChatMessage struct {
	Role    *byte // role string
	Content *byte // content string
}

// Sampler chain parameters
type SamplerChainParams struct {
	NoPerf uint8 // whether to measure performance timings (bool as uint8)
}

// Logit bias
type LogitBias struct {
	Token Token
	Bias  float32
}
