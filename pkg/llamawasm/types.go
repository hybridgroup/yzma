//go:build js && wasm

package llamawasm

// The shim gives each object of llama.cpp a small int32 handle. A handle of 0
// is never valid. The Go code never holds an address inside the llama.cpp
// module, so the module is free to move or grow its memory.
type (
	// Token is one token of a vocabulary.
	Token int32

	// Pos is the position of a token in a sequence.
	Pos int32

	// SeqId is the identifier of a sequence.
	SeqId int32

	// Model is a handle to a loaded model.
	Model int32

	// Context is a handle to an inference context.
	Context int32

	// Vocab is a handle to the vocabulary of a model.
	Vocab int32

	// Sampler is a handle to a sampler or to a chain of samplers.
	Sampler int32
)

// PoolingType is how the embeddings of the tokens of a sequence become one
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

// ModelParams holds what the shim can set while it loads a model.
//
// The struct is much smaller than llama.ModelParams, because a WebAssembly
// build has no device list.
type ModelParams struct {
	// NGpuLayers is the number of layers to put on the GPU. A value larger than
	// the number of layers of the model puts them all there.
	//
	// It does nothing in a build for the CPU. Test [GPUDevice] first: it is
	// empty when llama.cpp has no GPU.
	NGpuLayers int32
}

// ModelDefaultParams gives the parameters that a model uses if the program
// changes nothing.
func ModelDefaultParams() ModelParams {
	return ModelParams{
		NGpuLayers: 0,
	}
}

// ContextParams holds what the shim can set while it makes a context.
type ContextParams struct {
	NCtx        uint32      // size of the text context, 0 = from the model
	NBatch      uint32      // largest logical batch, 0 = from llama.cpp
	NUbatch     uint32      // largest physical batch, 0 = from llama.cpp
	NThreads    int32       // number of threads, 0 = from the module
	PoolingType PoolingType // how to pool embeddings
	Embeddings  uint8       // 1 to compute embeddings
}

// ContextDefaultParams gives the parameters that a context uses if the program
// changes nothing.
//
// NThreads comes from [Threads], the number that the module can use. Leaving it
// at 0 would give the four threads that llama.cpp asks for whatever the machine
// has, which is slower on a machine with more.
func ContextDefaultParams() ContextParams {
	return ContextParams{
		NCtx:        0,
		NBatch:      0,
		NUbatch:     0,
		NThreads:    Threads(),
		PoolingType: PoolingTypeUnspecified,
		Embeddings:  0,
	}
}

// SamplerChainParams holds the parameters of a chain of samplers.
type SamplerChainParams struct {
	NoPerf uint8 // 1 to stop the measurement of the time of each sample
}

// SamplerChainDefaultParams gives the parameters that a chain uses if the
// program changes nothing.
func SamplerChainDefaultParams() SamplerChainParams {
	return SamplerChainParams{NoPerf: 1}
}

// Batch holds the tokens of one call to Decode or Encode.
//
// On a native platform llama.Batch is a C struct. Here it is an ordinary Go
// struct, because the shim makes the C batch itself.
type Batch struct {
	// NTokens is the number of tokens in the batch. A generation loop reads
	// it to move the position forward.
	NTokens int32

	tokens []Token
}
