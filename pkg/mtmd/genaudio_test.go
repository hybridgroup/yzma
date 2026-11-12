package mtmd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"os"
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestHelperGenAudioInputLayout(t *testing.T) {
	var in HelperGenAudioInput
	if size := unsafe.Sizeof(in); size != 72 {
		t.Fatalf("HelperGenAudioInput size %d does not match C size 72", size)
	}
	offsets := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"SeqID", unsafe.Offsetof(in.SeqID), 4},
		{"Prompt", unsafe.Offsetof(in.Prompt), 8},
		{"PromptLen", unsafe.Offsetof(in.PromptLen), 16},
		{"NPast", unsafe.Offsetof(in.NPast), 24},
		{"OutMode", unsafe.Offsetof(in.OutMode), 28},
		{"SpeakerRef", unsafe.Offsetof(in.SpeakerRef), 32},
		{"Lang", unsafe.Offsetof(in.Lang), 40},
		{"TopK", unsafe.Offsetof(in.TopK), 48},
		{"TopP", unsafe.Offsetof(in.TopP), 52},
		{"Temp", unsafe.Offsetof(in.Temp), 56},
		{"Seed", unsafe.Offsetof(in.Seed), 60},
		{"OutType", unsafe.Offsetof(in.OutType), 64},
	}
	for _, o := range offsets {
		if o.got != o.want {
			t.Errorf("%s offset %d does not match C offset %d", o.name, o.got, o.want)
		}
	}
}

func TestHelperGenAudioInputStrings(t *testing.T) {
	var in HelperGenAudioInput
	in.SetPrompt("hello")
	if in.PromptLen != 5 || unsafe.String(in.Prompt, in.PromptLen) != "hello" {
		t.Fatal("SetPrompt did not set the prompt")
	}
	in.SetPrompt("")
	if in.Prompt != nil || in.PromptLen != 0 {
		t.Fatal("SetPrompt with empty text did not clear the prompt")
	}

	if err := in.SetLang("english"); err != nil || in.Lang == nil {
		t.Fatalf("SetLang failed: %v", err)
	}
	if err := in.SetLang(""); err != nil || in.Lang != nil {
		t.Fatal("SetLang with empty lang did not clear it")
	}
	if err := in.SetLang("a\x00b"); err == nil {
		t.Fatal("SetLang with a NUL byte did not fail")
	}
}

func TestHelperGenAudioZeroHandles(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if !hasHelperGenAudio {
		if _, err := HelperGenAudioInputDefault(); !errors.Is(err, ErrNoGenAudio) {
			t.Fatalf("HelperGenAudioInputDefault without the symbols returned %v", err)
		}
		if _, err := HelperGenAudioInit(0, 0); !errors.Is(err, ErrNoGenAudio) {
			t.Fatalf("HelperGenAudioInit without the symbols returned %v", err)
		}
		t.Skip("mtmd_helper_gen_audio not in this library")
	}

	input, err := HelperGenAudioInputDefault()
	if err != nil {
		t.Fatalf("HelperGenAudioInputDefault failed: %v", err)
	}
	if input.TopK != 50 || input.TopP != 1 || input.Temp != -1 || input.Seed != math.MaxUint32 ||
		input.OutType != HelperGenAudioOutTypeWAV || input.Stream || input.Prompt != nil {
		t.Fatalf("unexpected default input %+v", input)
	}

	if _, err := HelperGenAudioInit(0, 0); err == nil {
		t.Fatal("HelperGenAudioInit with zero contexts did not fail")
	}
	HelperGenAudioFree(0)
	HelperGenAudioReset(0)
	if err := HelperGenAudioSetInput(0, &input); err == nil {
		t.Fatal("HelperGenAudioSetInput with zero helper did not fail")
	}
	if _, err := HelperGenAudioStepPrompt(0, 1); err == nil {
		t.Fatal("HelperGenAudioStepPrompt with zero helper did not fail")
	}
	if _, _, err := HelperGenAudioStepGen(0, llama.TokenNull, nil); err == nil {
		t.Fatal("HelperGenAudioStepGen with zero helper did not fail")
	}
	if _, err := HelperGenAudioGetOutput(0); err == nil {
		t.Fatal("HelperGenAudioGetOutput with zero helper did not fail")
	}
	if err := HelperGenAudioFlush(0); err == nil {
		t.Fatal("HelperGenAudioFlush with zero helper did not fail")
	}
	if hasGenAudioVoice {
		if _, err := GenAudioVoiceName(0, 0); err == nil {
			t.Fatal("GenAudioVoiceName with zero context did not fail")
		}
		if _, err := GenAudioGetVoice(0, "default"); err == nil {
			t.Fatal("GenAudioGetVoice with zero context did not fail")
		}
	}
}

func TestHelperGenAudioUnsupportedModel(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	if !hasHelperGenAudio || !hasGenAudioVoice {
		t.Skip("mtmd_helper_gen_audio not in this library")
	}

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	lctx, err := llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer llama.Free(lctx)

	ctx, err := InitFromFile(mmprojFile, model, ContextParamsDefault())
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	name, err := GenAudioVoiceName(ctx, 0)
	if err != nil || name != "" {
		t.Fatalf("GenAudioVoiceName on a model without voices returned %q, %v", name, err)
	}
	voice, err := GenAudioGetVoice(ctx, "default")
	if err != nil || voice != 0 {
		BitmapFree(voice)
		t.Fatalf("GenAudioGetVoice on a model without voices returned %v, %v", voice, err)
	}

	gen, err := HelperGenAudioInit(lctx, ctx)
	if err != nil {
		t.Fatalf("HelperGenAudioInit failed: %v", err)
	}
	defer HelperGenAudioFree(gen)

	input, _ := HelperGenAudioInputDefault()
	input.SetPrompt("hello")
	if err := HelperGenAudioSetInput(gen, &input); err == nil {
		t.Fatal("HelperGenAudioSetInput on an unsupported model did not fail")
	}
}

func TestHelperGenAudio(t *testing.T) {
	modelFile := os.Getenv("YZMA_TEST_GENAUDIO_MODEL")
	mmprojFile := os.Getenv("YZMA_TEST_GENAUDIO_MMPROJ")
	if modelFile == "" || mmprojFile == "" {
		t.Skip("no YZMA_TEST_GENAUDIO_MODEL or YZMA_TEST_GENAUDIO_MMPROJ skipping test")
	}

	testSetup(t)
	defer testCleanup(t)

	if !hasHelperGenAudio || !hasGenAudioVoice {
		t.Skip("mtmd_helper_gen_audio not in this library")
	}

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	params := llama.ContextDefaultParams()
	params.NCtx = 4096
	params.Embeddings = 1
	lctx, err := llama.InitFromModel(model, params)
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer llama.Free(lctx)

	ctx, err := InitFromFile(mmprojFile, model, ContextParamsDefault())
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	var speaker Bitmap
	name, err := GenAudioVoiceName(ctx, 0)
	if err != nil {
		t.Fatalf("GenAudioVoiceName failed: %v", err)
	}
	if name != "" {
		if speaker, err = GenAudioGetVoice(ctx, name); err != nil || speaker == 0 {
			t.Fatalf("GenAudioGetVoice(%q) returned %v, %v", name, speaker, err)
		}
		defer BitmapFree(speaker)
	}
	if unknown, err := GenAudioGetVoice(ctx, "no such voice"); err != nil || unknown != 0 {
		BitmapFree(unknown)
		t.Fatalf("GenAudioGetVoice for an unknown voice returned %v, %v", unknown, err)
	}

	gen, err := HelperGenAudioInit(lctx, ctx)
	if err != nil {
		t.Fatalf("HelperGenAudioInit failed: %v", err)
	}
	defer HelperGenAudioFree(gen)

	vocab := llama.ModelGetVocab(model)
	var bias []llama.LogitBias
	for _, tok := range llama.VocabGetSuppressTokens(vocab) {
		bias = append(bias, llama.LogitBias{Token: tok, Bias: float32(math.Inf(-1))})
	}

	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	if len(bias) > 0 {
		llama.SamplerChainAdd(sampler, llama.SamplerInitLogitBias(llama.VocabNTokens(vocab), int32(len(bias)), &bias[0]))
	}
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopK(40))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopP(0.95, 1))
	llama.SamplerChainAdd(sampler, llama.SamplerInitMinP(0.05, 1))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTemp(0.8))
	llama.SamplerChainAdd(sampler, llama.SamplerInitDist(42))
	defer llama.SamplerFree(sampler)

	mem, err := llama.GetMemory(lctx)
	if err != nil {
		t.Fatalf("GetMemory failed: %v", err)
	}
	nEmbd := llama.ModelNEmbd(model)

	speak := func() HelperGenAudioOutput {
		input, err := HelperGenAudioInputDefault()
		if err != nil {
			t.Fatalf("HelperGenAudioInputDefault failed: %v", err)
		}
		input.SetPrompt("Hello from yzma.")
		input.SpeakerRef = speaker
		input.Seed = 42

		if err := HelperGenAudioSetInput(gen, &input); err != nil {
			t.Fatalf("HelperGenAudioSetInput failed: %v", err)
		}
		for {
			remaining, err := HelperGenAudioStepPrompt(gen, int32(params.NBatch))
			if err != nil {
				t.Fatalf("HelperGenAudioStepPrompt failed: %v", err)
			}
			if remaining == 0 {
				break
			}
		}

		sampled := llama.SamplerSample(sampler, lctx, -1)
		hState, err := llama.GetEmbeddingsIth(lctx, -1, nEmbd)
		if err != nil || hState == nil {
			t.Fatalf("GetEmbeddingsIth returned %v", err)
		}

		frames := 0
		for stop := false; !stop && frames < 512; {
			var next []float32
			next, stop, err = HelperGenAudioStepGen(gen, sampled, hState)
			if err != nil {
				t.Fatalf("HelperGenAudioStepGen failed at frame %d: %v", frames, err)
			}
			if next == nil {
				break
			}
			frames++
			hState = next
			sampled = llama.SamplerSample(sampler, lctx, -1)
		}
		if frames == 0 {
			t.Fatal("no audio frames generated")
		}

		output, err := HelperGenAudioGetOutput(gen)
		if err != nil {
			t.Fatalf("HelperGenAudioGetOutput failed: %v", err)
		}
		return output
	}

	first := speak()
	if first.SampleRate <= 0 || first.NSamples <= 0 {
		t.Fatalf("invalid output sample rate %d samples %d", first.SampleRate, first.NSamples)
	}
	if len(first.Data) < 44 || string(first.Data[0:4]) != "RIFF" || string(first.Data[8:12]) != "WAVE" {
		t.Fatal("output is not a WAV file")
	}
	if rate := int32(binary.LittleEndian.Uint32(first.Data[24:28])); rate != first.SampleRate {
		t.Fatalf("WAV sample rate %d does not match %d", rate, first.SampleRate)
	}
	if want := 44 + 2*first.NSamples; int64(len(first.Data)) != want {
		t.Fatalf("WAV size %d does not match %d samples", len(first.Data), first.NSamples)
	}

	kept := bytes.Clone(first.Data)
	HelperGenAudioReset(gen)
	if _, err := llama.MemorySeqRm(mem, 0, -1, -1); err != nil {
		t.Fatalf("MemorySeqRm failed: %v", err)
	}
	llama.SamplerReset(sampler)

	second := speak()
	if !bytes.Equal(first.Data, kept) {
		t.Fatal("output of the first run changed after the helper was reused")
	}
	if !bytes.Equal(first.Data, second.Data) {
		t.Fatal("the same seed gave different audio")
	}
}
