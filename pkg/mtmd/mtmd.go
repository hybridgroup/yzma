package mtmd

import (
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

//	struct mtmd_input_text {
//	    const char * text;
//	    bool add_special;
//	    bool parse_special;
//	};
type InputText struct {
	Text         *byte
	AddSpecial   bool
	ParseSpecial bool
}

// Opaque types (represented as pointers)
type Context uintptr
type ImageTokens uintptr
type InputChunk uintptr
type InputChunks uintptr

//	struct mtmd_context_params {
//	    bool use_gpu;
//	    bool print_timings;
//	    int n_threads;
//	    enum ggml_log_level verbosity;
//	    const char * image_marker; // deprecated, use media_marker instead
//	    const char * media_marker;
//	};
type ContextParamsType struct {
	UseGPU       bool
	PrintTimings bool
	Threads      int32
	Verbosity    llama.LogLevel
	ImageMarker  *byte
	MediaMarker  *byte
}

var (
	FFITypeContextParams = ffi.NewType(&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer)
	FFITypeInputText     = ffi.NewType(&ffi.TypePointer, &ffi.TypeUint8, &ffi.TypeUint8)
)

var (
	// MTMD_API const char * mtmd_default_marker(void);
	defaultMarkerFunc ffi.Fun

	// MTMD_API struct mtmd_context_params mtmd_context_params_default(void);
	contextParamsDefaultFunc ffi.Fun

	// MTMD_API mtmd_context * mtmd_init_from_file(const char * mmproj_fname,
	//                                         const struct llama_model * text_model,
	//                                         const struct mtmd_context_params ctx_params);
	initFromFileFunc ffi.Fun

	// MTMD_API void mtmd_free(mtmd_context * ctx);
	freeFunc ffi.Fun

	// MTMD_API bool mtmd_support_vision(mtmd_context * ctx);
	supportVisionFunc ffi.Fun

	// MTMD_API int32_t mtmd_tokenize(mtmd_context * ctx,
	//                            mtmd_input_chunks * output,
	//                            const mtmd_input_text * text,
	//                            const mtmd_bitmap ** bitmaps,
	//                            size_t n_bitmaps);
	tokenizeFunc ffi.Fun

	// MTMD_API int32_t mtmd_helper_eval_chunks(mtmd_context * ctx,
	//                                          struct llama_context * lctx,
	//                                          const mtmd_input_chunks * chunks,
	//                                          llama_pos n_past,
	//                                          llama_seq_id seq_id,
	//                                          int32_t n_batch,
	//                                          bool logits_last,
	//                                          llama_pos * new_n_past);
	helperEvalChunksFunc ffi.Fun

	// MTMD_API bool mtmd_decode_use_non_causal(mtmd_context * ctx);
	decodeUseNonCausalFunc ffi.Fun

	// MTMD_API bool mtmd_decode_use_mrope(mtmd_context * ctx);
	decodeUseMRopeFunc ffi.Fun

	// MTMD_API bool mtmd_support_audio(mtmd_context * ctx);
	supportAudioFunc ffi.Fun
)

func loadFuncs(lib ffi.Lib) error {
	var err error

	if defaultMarkerFunc, err = lib.Prep("mtmd_default_marker", &ffi.TypePointer); err != nil {
		return loadError("mtmd_default_marker", err)
	}

	if contextParamsDefaultFunc, err = lib.Prep("mtmd_context_params_default", &FFITypeContextParams); err != nil {
		return loadError("mtmd_context_params_default", err)
	}

	if initFromFileFunc, err = lib.Prep("mtmd_init_from_file", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &FFITypeContextParams); err != nil {
		return loadError("mtmd_init_from_file", err)
	}

	if freeFunc, err = lib.Prep("mtmd_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("mtmd_free", err)
	}

	if supportVisionFunc, err = lib.Prep("mtmd_support_vision", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("mtmd_support_vision", err)
	}

	if tokenizeFunc, err = lib.Prep("mtmd_tokenize", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint64); err != nil {
		return loadError("mtmd_tokenize", err)
	}

	if helperEvalChunksFunc, err = lib.Prep("mtmd_helper_eval_chunks", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer,
		&ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeUint8, &ffi.TypePointer); err != nil {

		return loadError("mtmd_helper_eval_chunks", err)
	}

	if decodeUseNonCausalFunc, err = lib.Prep("mtmd_decode_use_non_causal", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("mtmd_decode_use_non_causal", err)
	}

	if decodeUseMRopeFunc, err = lib.Prep("mtmd_decode_use_mrope", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("mtmd_decode_use_mrope", err)
	}

	if supportAudioFunc, err = lib.Prep("mtmd_support_audio", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("mtmd_support_audio", err)
	}

	return nil
}

func DefaultMarker() string {
	var marker *byte
	defaultMarkerFunc.Call(unsafe.Pointer(&marker))
	return utils.BytePtrToString(marker)
}

func ContextParamsDefault() ContextParamsType {
	var ctx ContextParamsType
	contextParamsDefaultFunc.Call(unsafe.Pointer(&ctx))
	return ctx
}

// InitFromFile initializes the mtmd context. mmprojFname is a projector file. model is a model that has already been opened.
// ctxParams are the ContextParamsType for the new Context.
func InitFromFile(mmprojFname string, model llama.Model, ctxParams ContextParamsType) Context {
	var ctx Context
	file := &[]byte(mmprojFname + "\x00")[0]
	initFromFileFunc.Call(unsafe.Pointer(&ctx), unsafe.Pointer(&file), unsafe.Pointer(&model), unsafe.Pointer(&ctxParams))
	return ctx
}

// Free frees a Context that has already been created using InitFromFile.
func Free(ctx Context) {
	freeFunc.Call(nil, unsafe.Pointer(&ctx))
}

// SupportVision returns whether the current model supports vision input.
func SupportVision(ctx Context) bool {
	var result ffi.Arg
	supportVisionFunc.Call(&result, unsafe.Pointer(&ctx))

	return result.Bool()
}

// Tokenize an input text prompt and a list of bitmaps (images/audio)
// the prompt must have the input image marker (default: "<__media__>") in it
// the default marker is defined by mtmd_default_marker()
// the marker will be replaced with the image/audio chunk
// for example:
//
//	"here is an image: <__media__>\ndescribe it in detail."
//	this will gives 3 chunks:
//	1. "here is an image: <start_of_image>"
//	2. (image/audio tokens)
//	3. "<end_of_image>\ndescribe it in detail."
//
// number of bitmaps must be equal to the number of markers in the prompt
// this function is thread-safe (shared ctx)
// return values:
//
//	0 on success
//	1 on number of bitmaps not matching the number of markers
//	2 on image preprocessing error
func Tokenize(ctx Context, out InputChunks, text *InputText, bitmaps []Bitmap) int32 {
	bt := unsafe.SliceData(bitmaps)
	nBitmaps := uint64(len(bitmaps))

	var result ffi.Arg
	tokenizeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), unsafe.Pointer(&out), unsafe.Pointer(&text), unsafe.Pointer(&bt), unsafe.Pointer(&nBitmaps))

	return int32(result)
}

// NewInputText create a new InputText to be used for calling Tokenize.
func NewInputText(text string, addSpecial, parseSpecial bool) *InputText {
	text += "\x00"
	p := unsafe.StringData(text)
	return &InputText{
		Text:         p,
		AddSpecial:   addSpecial,
		ParseSpecial: parseSpecial,
	}
}

// HelperEvalChunks is a helper function that automatically:
// 1. run llama.Decode() on text chunks
// 2. run mtmd.Encode() on image chunks, then mtmd.GetOutputEmbd() and then llama.Decode()
// if any of the mtmd.Encode() or llama.Decode() calls return non-zero, stop and forward the error
// otherwise, returns 0 on success
// this function is NOT thread-safe
func HelperEvalChunks(ctx Context, lctx llama.Context, chunks InputChunks, nPast llama.Pos, seqID llama.SeqId, nBatch int32, logitsLast bool, newNPast *llama.Pos) int32 {
	muHelperEvalChunks.Lock()
	defer muHelperEvalChunks.Unlock()

	var result ffi.Arg
	helperEvalChunksFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), unsafe.Pointer(&lctx), unsafe.Pointer(&chunks), unsafe.Pointer(&nPast), unsafe.Pointer(&seqID),
		unsafe.Pointer(&nBatch), unsafe.Pointer(&logitsLast), unsafe.Pointer(&newNPast))

	return int32(result)
}

// DecodeUseNonCausal checks if the non-causal mask needs to be set before llama_decode.
func DecodeUseNonCausal(ctx Context) bool {
	var result ffi.Arg
	decodeUseNonCausalFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return result.Bool()
}

// DecodeUseMRope checks if the current model uses M-RoPE for llama_decode.
func DecodeUseMRope(ctx Context) bool {
	var result ffi.Arg
	decodeUseMRopeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return result.Bool()
}

// SupportAudio checks if the current model supports audio input.
func SupportAudio(ctx Context) bool {
	var result ffi.Arg
	supportAudioFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return result.Bool()
}
