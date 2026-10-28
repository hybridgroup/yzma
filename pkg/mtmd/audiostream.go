package mtmd

import (
	"errors"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/loader"
	"github.com/jupiterrider/ffi"
)

// ErrNoAudioStream is returned when libmtmd does not have the audio stream API.
var ErrNoAudioStream = errors.New("mtmd_audio_stream not found in library")

// AudioStream is an opaque handle to a realtime audio input stream.
// It turns PCM audio into decoder ready embedding rows, one row per decoder position.
type AudioStream uintptr

// AudioStreamParams mirrors the C struct mtmd_audio_stream_params.
type AudioStreamParams struct {
	// StepMs is the audio encoded per encoder call, rounded up to the model step.
	StepMs int32
	// Prompt holds tokens that open the stream, they go out with the first rows.
	// The memory must stay valid until AudioStreamInit returns.
	Prompt *llama.Token
	// NPrompt is the number of tokens in Prompt.
	NPrompt uint64
}

// Return values of AudioStreamPush and AudioStreamFlush.
const (
	AudioStreamOK    int32 = 0
	AudioStreamEnded int32 = 1
)

// Return values of AudioStreamStep. A positive value is the number of rows ready.
const (
	AudioStreamStepNeedAudio int32 = 0
	AudioStreamStepEnd       int32 = -1
	AudioStreamStepError     int32 = -2
)

var (
	// ffiTypeAudioStreamParams mirrors struct mtmd_audio_stream_params
	ffiTypeAudioStreamParams = ffi.NewType(&ffi.TypeSint32, &ffi.TypePointer, &ffiTypeSize)

	// MTMD_API bool mtmd_support_audio_stream(const mtmd_context * ctx);
	supportAudioStreamFunc ffi.Fun

	// MTMD_API struct mtmd_audio_stream_params mtmd_audio_stream_params_default(void);
	audioStreamParamsDefaultFunc ffi.Fun

	// MTMD_API mtmd_audio_stream * mtmd_audio_stream_init(mtmd_context * ctx, struct mtmd_audio_stream_params params);
	audioStreamInitFunc ffi.Fun

	// MTMD_API void mtmd_audio_stream_free(mtmd_audio_stream * stream);
	audioStreamFreeFunc ffi.Fun

	// MTMD_API int32_t mtmd_audio_stream_push(mtmd_audio_stream * stream, const float * samples, size_t n_samples);
	audioStreamPushFunc ffi.Fun

	// MTMD_API int32_t mtmd_audio_stream_flush(mtmd_audio_stream * stream, bool end);
	audioStreamFlushFunc ffi.Fun

	// MTMD_API int32_t mtmd_audio_stream_step(mtmd_audio_stream * stream, llama_token last_token);
	audioStreamStepFunc ffi.Fun

	// MTMD_API const float * mtmd_audio_stream_get_output_embd(const mtmd_audio_stream * stream);
	audioStreamGetOutputEmbdFunc ffi.Fun

	// MTMD_API int64_t mtmd_audio_stream_get_n_rows(const mtmd_audio_stream * stream);
	audioStreamGetNRowsFunc ffi.Fun

	// MTMD_API int64_t mtmd_audio_stream_get_n_rows_flushed(const mtmd_audio_stream * stream);
	audioStreamGetNRowsFlushedFunc ffi.Fun

	hasAudioStream bool
)

func loadAudioStreamFuncs(lib loader.Lib) error {
	var err error

	// Older libraries do not have the audio stream API, so it is optional.
	hasAudioStream = false

	if supportAudioStreamFunc, err = lib.Prep("mtmd_support_audio_stream", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return nil
	}

	if audioStreamParamsDefaultFunc, err = lib.Prep("mtmd_audio_stream_params_default", &ffiTypeAudioStreamParams); err != nil {
		return nil
	}

	if audioStreamInitFunc, err = lib.Prep("mtmd_audio_stream_init", &ffi.TypePointer, &ffi.TypePointer, &ffiTypeAudioStreamParams); err != nil {
		return nil
	}

	if audioStreamFreeFunc, err = lib.Prep("mtmd_audio_stream_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return nil
	}

	if audioStreamPushFunc, err = lib.Prep("mtmd_audio_stream_push", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffiTypeSize); err != nil {
		return nil
	}

	if audioStreamFlushFunc, err = lib.Prep("mtmd_audio_stream_flush", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeUint8); err != nil {
		return nil
	}

	if audioStreamStepFunc, err = lib.Prep("mtmd_audio_stream_step", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return nil
	}

	if audioStreamGetOutputEmbdFunc, err = lib.Prep("mtmd_audio_stream_get_output_embd", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return nil
	}

	if audioStreamGetNRowsFunc, err = lib.Prep("mtmd_audio_stream_get_n_rows", &ffi.TypeSint64, &ffi.TypePointer); err != nil {
		return nil
	}

	if audioStreamGetNRowsFlushedFunc, err = lib.Prep("mtmd_audio_stream_get_n_rows_flushed", &ffi.TypeSint64, &ffi.TypePointer); err != nil {
		return nil
	}

	hasAudioStream = true

	return nil
}

// SupportAudioStream reports whether the model supports the realtime audio stream.
// It returns false when libmtmd does not have the audio stream API.
func SupportAudioStream(ctx Context) bool {
	if !hasAudioStream || ctx == 0 {
		return false
	}
	var result ffi.Arg
	supportAudioStreamFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return result.Bool()
}

// AudioStreamParamsDefault returns the default audio stream parameters.
func AudioStreamParamsDefault() (AudioStreamParams, error) {
	var params AudioStreamParams
	if !hasAudioStream {
		return params, ErrNoAudioStream
	}
	audioStreamParamsDefaultFunc.Call(unsafe.Pointer(&params))
	return params, nil
}

// AudioStreamInit creates an audio stream. It returns 0 when the model does not support it.
// The ctx must have been created with the text model passed to InitFromFile.
// The caller must call AudioStreamFree when done.
func AudioStreamInit(ctx Context, params AudioStreamParams) (AudioStream, error) {
	var stream AudioStream
	if !hasAudioStream {
		return stream, ErrNoAudioStream
	}
	if ctx == 0 {
		return stream, errors.New("invalid mtmd context handle")
	}
	audioStreamInitFunc.Call(unsafe.Pointer(&stream), unsafe.Pointer(&ctx), unsafe.Pointer(&params))
	if stream == 0 {
		return stream, errors.New("mtmd_audio_stream_init failed")
	}
	return stream, nil
}

// AudioStreamFree frees an audio stream.
func AudioStreamFree(stream AudioStream) {
	if !hasAudioStream || stream == 0 {
		return
	}
	audioStreamFreeFunc.Call(nil, unsafe.Pointer(&stream))
}

// AudioStreamPush feeds PCM F32 samples at GetAudioSampleRate to the stream.
// It returns AudioStreamOK, or AudioStreamEnded after the stream has ended.
func AudioStreamPush(stream AudioStream, samples []float32) int32 {
	if !hasAudioStream || stream == 0 {
		return AudioStreamEnded
	}
	samplesPtr := unsafe.SliceData(samples)
	nSamples := uint64(len(samples))

	var result ffi.Arg
	audioStreamPushFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&stream), unsafe.Pointer(&samplesPtr), &nSamples)
	return int32(result)
}

// AudioStreamFlush pads the stream so that all pending rows can be produced.
// With end false the stream stays open, with end true no more audio is accepted.
// It returns AudioStreamOK, or AudioStreamEnded after the stream has ended.
func AudioStreamFlush(stream AudioStream, end bool) int32 {
	if !hasAudioStream || stream == 0 {
		return AudioStreamEnded
	}
	var endArg uint8
	if end {
		endArg = 1
	}

	var result ffi.Arg
	audioStreamFlushFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&stream), &endArg)
	return int32(result)
}

// AudioStreamStep produces the next decoder rows.
// lastToken is the token sampled from the previous row, or llama.TokenNull before the first sample.
// It returns the number of rows ready, AudioStreamStepNeedAudio, AudioStreamStepEnd or AudioStreamStepError.
func AudioStreamStep(stream AudioStream, lastToken llama.Token) int32 {
	if !hasAudioStream || stream == 0 {
		return AudioStreamStepError
	}
	var result ffi.Arg
	audioStreamStepFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&stream), &lastToken)
	return int32(result)
}

// AudioStreamGetOutputEmbd returns the rows of the last step.
// nRows is the AudioStreamStep result and nEmbd is llama.ModelNEmbdInp of the text model.
// The slice points to native memory owned by the stream, valid until the next step or free.
func AudioStreamGetOutputEmbd(stream AudioStream, nRows, nEmbd int32) ([]float32, error) {
	if !hasAudioStream {
		return nil, ErrNoAudioStream
	}
	if stream == 0 {
		return nil, errors.New("invalid audio stream handle")
	}
	if nRows <= 0 || nEmbd <= 0 {
		return nil, nil
	}

	var embdPtr unsafe.Pointer
	audioStreamGetOutputEmbdFunc.Call(unsafe.Pointer(&embdPtr), unsafe.Pointer(&stream))
	if embdPtr == nil {
		return nil, errors.New("mtmd_audio_stream_get_output_embd returned null pointer")
	}

	return unsafe.Slice((*float32)(embdPtr), int(nRows)*int(nEmbd)), nil
}

// AudioStreamGetNRows returns the rows produced so far, the prefix rows included.
func AudioStreamGetNRows(stream AudioStream) int64 {
	if !hasAudioStream || stream == 0 {
		return 0
	}
	var result ffi.Arg
	audioStreamGetNRowsFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&stream))
	return int64(result)
}

// AudioStreamGetNRowsFlushed returns the row count the stream reaches at the end of the last flush.
// Reaching it means the flushed audio rows were produced, not that they were decoded.
func AudioStreamGetNRowsFlushed(stream AudioStream) int64 {
	if !hasAudioStream || stream == 0 {
		return 0
	}
	var result ffi.Arg
	audioStreamGetNRowsFlushedFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&stream))
	return int64(result)
}
