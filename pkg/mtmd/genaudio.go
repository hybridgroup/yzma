package mtmd

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/loader"
	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

// ErrNoGenAudio is returned when libmtmd does not have the audio generation API.
var ErrNoGenAudio = errors.New("mtmd_helper_gen_audio not found in library")

// HelperGenAudio is an opaque handle to a stateful audio generation helper.
// Free it before its llama and mtmd contexts, and never call it at the same time as other users of them.
type HelperGenAudio uintptr

// HelperGenAudioOutType is the format of the generated audio.
type HelperGenAudioOutType int32

const (
	// HelperGenAudioOutTypePCM is raw PCM.
	HelperGenAudioOutTypePCM HelperGenAudioOutType = iota
	// HelperGenAudioOutTypeWAV is WAV PCM 16 bit LE mono.
	HelperGenAudioOutTypeWAV
)

// HelperGenAudioOutMode is what a response that continues a sequence outputs.
type HelperGenAudioOutMode int32

const (
	// HelperGenAudioOutAuto lets the model switch between text and audio itself.
	HelperGenAudioOutAuto HelperGenAudioOutMode = iota
	HelperGenAudioOutText
	HelperGenAudioOutAudio
	// HelperGenAudioOutMixed is interleaved text and audio.
	HelperGenAudioOutMixed
)

// HelperGenAudioInput mirrors the C struct mtmd_helper_gen_audio_inp of llama.cpp b12219 and later.
// Use SetPrompt and SetLang to fill the string fields.
type HelperGenAudioInput struct {
	// Stream means HelperGenAudioGetOutput must be called after each HelperGenAudioStepGen.
	Stream bool
	SeqID  llama.SeqId
	// Prompt is the text to speak. It must stay valid until HelperGenAudioSetInput returns.
	Prompt    *byte
	PromptLen uint64
	// NPast above 0 continues a sequence of NPast positions with raw text.
	NPast   int32
	OutMode HelperGenAudioOutMode
	// SpeakerRef is an optional audio clip to clone or a voice from GenAudioGetVoice.
	SpeakerRef Bitmap
	// Lang is optional. It must stay valid until HelperGenAudioSetInput returns.
	Lang *byte
	TopK int32
	TopP float32
	// Temp is the audio frame temperature. Below 0 is the model default and 0 is greedy.
	Temp float32
	// Seed is the random seed. math.MaxUint32 gives a random one.
	Seed    uint32
	OutType HelperGenAudioOutType
}

// SetPrompt copies text into the input as the prompt.
func (in *HelperGenAudioInput) SetPrompt(text string) {
	if text == "" {
		in.Prompt = nil
		in.PromptLen = 0
		return
	}
	buf := []byte(text)
	in.Prompt = &buf[0]
	in.PromptLen = uint64(len(buf))
}

// SetLang copies lang into the input. An empty lang clears it.
func (in *HelperGenAudioInput) SetLang(lang string) error {
	if lang == "" {
		in.Lang = nil
		return nil
	}
	p, err := utils.BytePtrFromString(lang)
	if err != nil {
		return err
	}
	in.Lang = p
	return nil
}

var (
	// ffiTypeHelperGenAudioInput mirrors struct mtmd_helper_gen_audio_inp
	ffiTypeHelperGenAudioInput = ffi.NewType(&ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypePointer, &ffiTypeSize,
		&ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer,
		&ffi.TypeSint32, &ffi.TypeFloat, &ffi.TypeFloat, &ffi.TypeUint32, &ffi.TypeSint32)

	// MTMD_API struct mtmd_helper_gen_audio_inp mtmd_helper_gen_audio_inp_default(void);
	helperGenAudioInpDefaultFunc ffi.Fun

	// MTMD_API mtmd_helper_gen_audio * mtmd_helper_gen_audio_init(struct llama_context * lctx, struct mtmd_context * mctx);
	helperGenAudioInitFunc ffi.Fun

	// MTMD_API void mtmd_helper_gen_audio_free(mtmd_helper_gen_audio * ctx);
	helperGenAudioFreeFunc ffi.Fun

	// MTMD_API void mtmd_helper_gen_audio_reset(mtmd_helper_gen_audio * ctx);
	helperGenAudioResetFunc ffi.Fun

	// MTMD_API int32_t mtmd_helper_gen_audio_set_input(mtmd_helper_gen_audio * ctx, const struct mtmd_helper_gen_audio_inp * inp);
	helperGenAudioSetInputFunc ffi.Fun

	// MTMD_API int32_t mtmd_helper_gen_audio_step_prompt(mtmd_helper_gen_audio * ctx, int32_t n_batch);
	helperGenAudioStepPromptFunc ffi.Fun

	// MTMD_API int32_t mtmd_helper_gen_audio_step_gen(mtmd_helper_gen_audio * ctx, llama_token sampled,
	//                  const float * h_state_in, const float ** h_state_out, bool * out_stop);
	helperGenAudioStepGenFunc ffi.Fun

	// MTMD_API int32_t mtmd_helper_gen_audio_get_output(mtmd_helper_gen_audio * ctx, int32_t * out_sample_rate,
	//                  const char ** out_data, size_t * out_data_len, int64_t * out_n_samples);
	helperGenAudioGetOutputFunc ffi.Fun

	// MTMD_API int32_t mtmd_helper_gen_audio_flush(mtmd_helper_gen_audio * ctx);
	helperGenAudioFlushFunc ffi.Fun

	// MTMD_API const char * mtmd_gen_audio_voice_name(const mtmd_context * ctx, int32_t idx);
	genAudioVoiceNameFunc ffi.Fun

	// MTMD_API mtmd_bitmap * mtmd_gen_audio_get_voice(const mtmd_context * ctx, const char * name);
	genAudioGetVoiceFunc ffi.Fun

	hasHelperGenAudio bool
	hasGenAudioVoice  bool
)

func loadGenAudioFuncs(lib loader.Lib) error {
	// Older libraries do not have the audio generation API, so it is optional.
	hasHelperGenAudio = loadHelperGenAudioFuncs(lib)
	hasGenAudioVoice = loadGenAudioVoiceFuncs(lib)
	return nil
}

func loadHelperGenAudioFuncs(lib loader.Lib) bool {
	var err error

	if helperGenAudioInpDefaultFunc, err = lib.Prep("mtmd_helper_gen_audio_inp_default", &ffiTypeHelperGenAudioInput); err != nil {
		return false
	}

	if helperGenAudioInitFunc, err = lib.Prep("mtmd_helper_gen_audio_init", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return false
	}

	if helperGenAudioFreeFunc, err = lib.Prep("mtmd_helper_gen_audio_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return false
	}

	if helperGenAudioResetFunc, err = lib.Prep("mtmd_helper_gen_audio_reset", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return false
	}

	if helperGenAudioSetInputFunc, err = lib.Prep("mtmd_helper_gen_audio_set_input", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return false
	}

	if helperGenAudioStepPromptFunc, err = lib.Prep("mtmd_helper_gen_audio_step_prompt", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return false
	}

	if helperGenAudioStepGenFunc, err = lib.Prep("mtmd_helper_gen_audio_step_gen", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32,
		&ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return false
	}

	if helperGenAudioGetOutputFunc, err = lib.Prep("mtmd_helper_gen_audio_get_output", &ffi.TypeSint32, &ffi.TypePointer,
		&ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return false
	}

	if helperGenAudioFlushFunc, err = lib.Prep("mtmd_helper_gen_audio_flush", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return false
	}

	return true
}

func loadGenAudioVoiceFuncs(lib loader.Lib) bool {
	var err error

	if genAudioVoiceNameFunc, err = lib.Prep("mtmd_gen_audio_voice_name", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return false
	}

	if genAudioGetVoiceFunc, err = lib.Prep("mtmd_gen_audio_get_voice", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return false
	}

	return true
}

var errInvalidGenAudio = errors.New("invalid audio generation helper handle")

// HelperGenAudioInputDefault returns the default audio generation input.
func HelperGenAudioInputDefault() (HelperGenAudioInput, error) {
	var input HelperGenAudioInput
	if !hasHelperGenAudio {
		return input, ErrNoGenAudio
	}
	helperGenAudioInpDefaultFunc.Call(unsafe.Pointer(&input))
	return input, nil
}

// HelperGenAudioInit creates an audio generation helper for the llama context lctx and the mtmd context ctx.
// The llama context must have embeddings enabled. The caller must call HelperGenAudioFree when done.
func HelperGenAudioInit(lctx llama.Context, ctx Context) (HelperGenAudio, error) {
	var gen HelperGenAudio
	if !hasHelperGenAudio {
		return gen, ErrNoGenAudio
	}
	if lctx == 0 {
		return gen, errors.New("invalid llama context handle")
	}
	if ctx == 0 {
		return gen, errors.New("invalid mtmd context handle")
	}
	helperGenAudioInitFunc.Call(unsafe.Pointer(&gen), unsafe.Pointer(&lctx), unsafe.Pointer(&ctx))
	if gen == 0 {
		return gen, errors.New("mtmd_helper_gen_audio_init failed")
	}
	return gen, nil
}

// HelperGenAudioFree frees an audio generation helper.
func HelperGenAudioFree(gen HelperGenAudio) {
	if !hasHelperGenAudio || gen == 0 {
		return
	}
	helperGenAudioFreeFunc.Call(nil, unsafe.Pointer(&gen))
}

// HelperGenAudioReset clears the helper state. It does not clear the llama sequence.
func HelperGenAudioReset(gen HelperGenAudio) {
	if !hasHelperGenAudio || gen == 0 {
		return
	}
	helperGenAudioResetFunc.Call(nil, unsafe.Pointer(&gen))
}

// HelperGenAudioSetInput starts a new generation from input.
// It fails when the model does not support audio generation.
func HelperGenAudioSetInput(gen HelperGenAudio, input *HelperGenAudioInput) error {
	if !hasHelperGenAudio {
		return ErrNoGenAudio
	}
	if gen == 0 {
		return errInvalidGenAudio
	}
	if input == nil {
		return errors.New("invalid audio generation input")
	}

	var result ffi.Arg
	helperGenAudioSetInputFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&gen), unsafe.Pointer(&input))
	if status := int32(result); status != 0 {
		return fmt.Errorf("mtmd_helper_gen_audio_set_input failed with status %d", status)
	}
	return nil
}

// HelperGenAudioStepPrompt processes at most nBatch prompt tokens.
// It returns the number of prompt tokens left, 0 when the prompt is done.
func HelperGenAudioStepPrompt(gen HelperGenAudio, nBatch int32) (int32, error) {
	if !hasHelperGenAudio {
		return 0, ErrNoGenAudio
	}
	if gen == 0 {
		return 0, errInvalidGenAudio
	}

	var result ffi.Arg
	helperGenAudioStepPromptFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&gen), &nBatch)
	remaining := int32(result)
	if remaining < 0 {
		return 0, fmt.Errorf("mtmd_helper_gen_audio_step_prompt failed with status %d", remaining)
	}
	return remaining, nil
}

// HelperGenAudioStepGen generates one frame from the sampled token and the backbone hidden state.
// The returned state is native memory valid until the next step or reset, and nil when no frame was made.
func HelperGenAudioStepGen(gen HelperGenAudio, sampled llama.Token, hStateIn []float32) (hStateOut []float32, stop bool, err error) {
	if !hasHelperGenAudio {
		return nil, false, ErrNoGenAudio
	}
	if gen == 0 {
		return nil, false, errInvalidGenAudio
	}

	hIn := unsafe.SliceData(hStateIn)
	var hOut *float32
	hOutPtr := &hOut
	stopPtr := &stop

	var result ffi.Arg
	helperGenAudioStepGenFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&gen), &sampled,
		unsafe.Pointer(&hIn), unsafe.Pointer(&hOutPtr), unsafe.Pointer(&stopPtr))
	if status := int32(result); status != 0 {
		return nil, stop, fmt.Errorf("mtmd_helper_gen_audio_step_gen failed with status %d", status)
	}
	if hOut == nil {
		return nil, stop, nil
	}
	return unsafe.Slice(hOut, len(hStateIn)), stop, nil
}

// HelperGenAudioOutput is the generated audio.
type HelperGenAudioOutput struct {
	SampleRate int32
	// Data is a copy owned by the caller.
	Data []byte
	// NSamples is the number of PCM samples.
	NSamples int64
}

// HelperGenAudioGetOutput returns a copy of the generated audio.
// In stream mode it returns only the audio since the previous call.
func HelperGenAudioGetOutput(gen HelperGenAudio) (HelperGenAudioOutput, error) {
	var output HelperGenAudioOutput
	if !hasHelperGenAudio {
		return output, ErrNoGenAudio
	}
	if gen == 0 {
		return output, errInvalidGenAudio
	}

	var data *byte
	var dataLen uint64
	sampleRatePtr := &output.SampleRate
	dataPtr := &data
	dataLenPtr := &dataLen
	nSamplesPtr := &output.NSamples

	var result ffi.Arg
	helperGenAudioGetOutputFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&gen), unsafe.Pointer(&sampleRatePtr),
		unsafe.Pointer(&dataPtr), unsafe.Pointer(&dataLenPtr), unsafe.Pointer(&nSamplesPtr))
	if status := int32(result); status != 0 {
		return HelperGenAudioOutput{}, fmt.Errorf("mtmd_helper_gen_audio_get_output failed with status %d", status)
	}
	if data != nil && dataLen > 0 {
		output.Data = append([]byte(nil), unsafe.Slice(data, dataLen)...)
	}
	return output, nil
}

// HelperGenAudioFlush forces any buffered frames through the decoder.
// In stream mode call it once before the last HelperGenAudioGetOutput.
func HelperGenAudioFlush(gen HelperGenAudio) error {
	if !hasHelperGenAudio {
		return ErrNoGenAudio
	}
	if gen == 0 {
		return errInvalidGenAudio
	}

	var result ffi.Arg
	helperGenAudioFlushFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&gen))
	if status := int32(result); status != 0 {
		return fmt.Errorf("mtmd_helper_gen_audio_flush failed with status %d", status)
	}
	return nil
}

// GenAudioVoiceName returns the name of the idx-th preset voice of the model.
// It returns an empty string when idx is out of range.
func GenAudioVoiceName(ctx Context, idx int32) (string, error) {
	if !hasGenAudioVoice {
		return "", ErrNoGenAudio
	}
	if ctx == 0 {
		return "", errors.New("invalid mtmd context handle")
	}

	var name *byte
	genAudioVoiceNameFunc.Call(unsafe.Pointer(&name), unsafe.Pointer(&ctx), &idx)
	if name == nil {
		return "", nil
	}
	return utils.BytePtrToString(name), nil
}

// GenAudioGetVoice returns the preset voice with this name as a bitmap to use as SpeakerRef.
// It returns 0 when the model has no such voice. The caller must call BitmapFree when done.
func GenAudioGetVoice(ctx Context, name string) (Bitmap, error) {
	var bitmap Bitmap
	if !hasGenAudioVoice {
		return bitmap, ErrNoGenAudio
	}
	if ctx == 0 {
		return bitmap, errors.New("invalid mtmd context handle")
	}

	namePtr, err := utils.BytePtrFromString(name)
	if err != nil {
		return bitmap, err
	}
	genAudioGetVoiceFunc.Call(unsafe.Pointer(&bitmap), unsafe.Pointer(&ctx), unsafe.Pointer(&namePtr))
	return bitmap, nil
}
