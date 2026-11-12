package main

import (
	"fmt"
	"math"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
)

func main() {
	if err := handleFlags(); err != nil {
		showUsage()
		os.Exit(0)
	}

	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	if err := llama.Load(*libPath); err != nil {
		return fmt.Errorf("unable to load library: %w", err)
	}
	if err := mtmd.Load(*libPath); err != nil {
		return fmt.Errorf("unable to load library: %w", err)
	}

	if !*verbose {
		llama.LogSet(llama.LogSilent())
		mtmd.LogSet(llama.LogSilent())
	}

	llama.Init()
	defer llama.Close()

	model, err := llama.ModelLoadFromFile(*modelFile, llama.ModelDefaultParams())
	if err != nil {
		return fmt.Errorf("unable to load model from file: %w", err)
	}
	defer llama.ModelFree(model)

	// The helper reads the backbone hidden state, so the context needs embeddings.
	params := llama.ContextDefaultParams()
	params.NCtx = uint32(*contextSize)
	params.NBatch = uint32(*batchSize)
	params.Embeddings = 1
	lctx, err := llama.InitFromModel(model, params)
	if err != nil {
		return fmt.Errorf("unable to initialize context: %w", err)
	}
	defer llama.Free(lctx)

	mctx, err := mtmd.InitFromFile(*projFile, model, mtmd.ContextParamsDefault())
	if err != nil {
		return fmt.Errorf("unable to load projector: %w", err)
	}
	defer mtmd.Free(mctx)

	speaker, err := loadVoice(mctx)
	if err != nil {
		return err
	}
	if speaker != 0 {
		defer mtmd.BitmapFree(speaker)
	}

	gen, err := mtmd.HelperGenAudioInit(lctx, mctx)
	if err != nil {
		return err
	}
	defer mtmd.HelperGenAudioFree(gen)

	input, err := mtmd.HelperGenAudioInputDefault()
	if err != nil {
		return err
	}
	input.SetPrompt(*prompt)
	if err := input.SetLang(*lang); err != nil {
		return err
	}
	input.SpeakerRef = speaker
	input.TopK = int32(*topK)
	input.TopP = float32(*topP)
	input.Seed = uint32(*seed)

	if err := mtmd.HelperGenAudioSetInput(gen, &input); err != nil {
		return fmt.Errorf("model does not support audio generation: %w", err)
	}

	for {
		remaining, err := mtmd.HelperGenAudioStepPrompt(gen, int32(*batchSize))
		if err != nil {
			return err
		}
		if remaining == 0 {
			break
		}
	}

	sampler := newSampler(model)
	defer llama.SamplerFree(sampler)

	sampled := llama.SamplerSample(sampler, lctx, -1)
	hState, err := llama.GetEmbeddingsIth(lctx, -1, llama.ModelNEmbd(model))
	if err != nil {
		return err
	}

	frames := 0
	for stop := false; !stop && frames < *maxFrames; {
		var next []float32
		next, stop, err = mtmd.HelperGenAudioStepGen(gen, sampled, hState)
		if err != nil {
			return err
		}
		if next == nil {
			break
		}
		frames++
		hState = next
		sampled = llama.SamplerSample(sampler, lctx, -1)
	}

	output, err := mtmd.HelperGenAudioGetOutput(gen)
	if err != nil {
		return err
	}
	if err := os.WriteFile(*outputFile, output.Data, 0o644); err != nil {
		return err
	}

	fmt.Printf("wrote %s, %.2f seconds of audio from %d frames\n",
		*outputFile, float64(output.NSamples)/float64(output.SampleRate), frames)
	return nil
}

// loadVoice returns the preset voice to speak with, or 0 when the model has none.
func loadVoice(mctx mtmd.Context) (mtmd.Bitmap, error) {
	name := *voice
	if name == "" {
		first, err := mtmd.GenAudioVoiceName(mctx, 0)
		if err != nil {
			return 0, err
		}
		if first == "" {
			return 0, nil
		}
		name = first
	}

	speaker, err := mtmd.GenAudioGetVoice(mctx, name)
	if err != nil {
		return 0, err
	}
	if speaker == 0 {
		return 0, fmt.Errorf("unknown voice %q", name)
	}
	if *verbose {
		fmt.Println("speaking with voice", name)
	}
	return speaker, nil
}

// newSampler samples the backbone tokens, without the tokens the model suppresses.
func newSampler(model llama.Model) llama.Sampler {
	vocab := llama.ModelGetVocab(model)
	var bias []llama.LogitBias
	for _, tok := range llama.VocabGetSuppressTokens(vocab) {
		bias = append(bias, llama.LogitBias{Token: tok, Bias: float32(math.Inf(-1))})
	}

	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	if len(bias) > 0 {
		llama.SamplerChainAdd(sampler, llama.SamplerInitLogitBias(llama.VocabNTokens(vocab), int32(len(bias)), &bias[0]))
	}
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopK(int32(*topK)))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopP(float32(*topP), 1))
	llama.SamplerChainAdd(sampler, llama.SamplerInitMinP(float32(*minP), 1))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTemp(float32(*temperature)))
	llama.SamplerChainAdd(sampler, llama.SamplerInitDist(uint32(*seed)))
	return sampler
}
