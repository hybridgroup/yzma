package mtmd

import (
	"errors"
	"math"
	"os"
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestAudioStreamParamsLayout(t *testing.T) {
	if size := unsafe.Sizeof(AudioStreamParams{}); size != 4 {
		t.Fatalf("AudioStreamParams size %d does not match C size 4", size)
	}
}

func TestAudioStreamZeroHandles(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if SupportAudioStream(0) {
		t.Fatal("SupportAudioStream with zero context returned true")
	}
	if !hasAudioStream {
		if _, err := AudioStreamParamsDefault(); !errors.Is(err, ErrNoAudioStream) {
			t.Fatalf("AudioStreamParamsDefault without the symbols returned %v", err)
		}
		if _, err := AudioStreamInit(0, AudioStreamParams{}); !errors.Is(err, ErrNoAudioStream) {
			t.Fatalf("AudioStreamInit without the symbols returned %v", err)
		}
		t.Skip("mtmd_audio_stream not in this library")
	}

	params, err := AudioStreamParamsDefault()
	if err != nil {
		t.Fatalf("AudioStreamParamsDefault failed: %v", err)
	}
	if params.StepMs != 320 {
		t.Fatalf("default StepMs %d, expected 320", params.StepMs)
	}

	if _, err := AudioStreamInit(0, params); err == nil {
		t.Fatal("AudioStreamInit with zero context did not fail")
	}
	AudioStreamFree(0)
	if r := AudioStreamPush(0, []float32{0}); r != AudioStreamEnded {
		t.Fatalf("AudioStreamPush with zero stream returned %d", r)
	}
	if r := AudioStreamFlush(0, true); r != AudioStreamEnded {
		t.Fatalf("AudioStreamFlush with zero stream returned %d", r)
	}
	if r := AudioStreamStep(0, llama.TokenNull); r != AudioStreamStepError {
		t.Fatalf("AudioStreamStep with zero stream returned %d", r)
	}
	if _, err := AudioStreamGetOutputEmbd(0, 1, 1); err == nil {
		t.Fatal("AudioStreamGetOutputEmbd with zero stream did not fail")
	}
	if AudioStreamGetNRows(0) != 0 || AudioStreamGetNRowsFlushed(0) != 0 {
		t.Fatal("row counters with zero stream are not 0")
	}
}

func TestAudioStreamUnsupportedModel(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	if !hasAudioStream {
		t.Skip("mtmd_audio_stream not in this library")
	}

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	ctx, err := InitFromFile(mmprojFile, model, ContextParamsDefault())
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	if SupportAudioStream(ctx) {
		t.Skip("test model supports the audio stream")
	}

	params, _ := AudioStreamParamsDefault()
	stream, err := AudioStreamInit(ctx, params)
	if err == nil {
		AudioStreamFree(stream)
		t.Fatal("AudioStreamInit on an unsupported model did not fail")
	}
}

func TestAudioStream(t *testing.T) {
	modelFile := os.Getenv("YZMA_TEST_AUDIOSTREAM_MODEL")
	mmprojFile := os.Getenv("YZMA_TEST_AUDIOSTREAM_MMPROJ")
	if modelFile == "" || mmprojFile == "" {
		t.Skip("no YZMA_TEST_AUDIOSTREAM_MODEL or YZMA_TEST_AUDIOSTREAM_MMPROJ skipping test")
	}

	testSetup(t)
	defer testCleanup(t)

	if !hasAudioStream {
		t.Skip("mtmd_audio_stream not in this library")
	}

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	ctx, err := InitFromFile(mmprojFile, model, ContextParamsDefault())
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	if !SupportAudioStream(ctx) {
		t.Fatal("audio stream model does not support the audio stream")
	}

	params, _ := AudioStreamParamsDefault()
	stream, err := AudioStreamInit(ctx, params)
	if err != nil {
		t.Fatalf("AudioStreamInit failed: %v", err)
	}
	defer AudioStreamFree(stream)

	nEmbd := llama.ModelNEmbdInp(model)
	sampleRate := GetAudioSampleRate(ctx)
	if sampleRate <= 0 {
		t.Fatalf("invalid audio sample rate %d", sampleRate)
	}

	// One second of a 440 Hz tone.
	samples := make([]float32, sampleRate)
	for i := range samples {
		samples[i] = 0.1 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sampleRate)))
	}

	lctxParams := llama.ContextDefaultParams()
	lctxParams.NCtx = 2048
	lctx, err := llama.InitFromModel(model, lctxParams)
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer llama.Free(lctx)

	sampler := llama.SamplerChainInit(llama.SamplerChainParams{})
	llama.SamplerChainAdd(sampler, llama.SamplerInitGreedy())
	defer llama.SamplerFree(sampler)

	// Each row is decoded and the sampled token is passed to the next step.
	var lastToken llama.Token = llama.TokenNull
	drain := func() int64 {
		var produced int64
		for {
			n := AudioStreamStep(stream, lastToken)
			switch {
			case n > 0:
				embd, err := AudioStreamGetOutputEmbd(stream, n, nEmbd)
				if err != nil {
					t.Fatalf("AudioStreamGetOutputEmbd failed: %v", err)
				}
				if len(embd) != int(n)*int(nEmbd) {
					t.Fatalf("embedding length %d, expected %d", len(embd), int(n)*int(nEmbd))
				}
				batch := llama.Batch{BatchData: llama.BatchData{NTokens: n, Embd: unsafe.SliceData(embd)}}
				if r, _ := llama.Decode(lctx, batch); r != 0 {
					t.Fatalf("Decode returned %d", r)
				}
				lastToken = llama.SamplerSample(sampler, lctx, -1)
				produced += int64(n)
			case n == AudioStreamStepNeedAudio, n == AudioStreamStepEnd:
				return produced
			default:
				t.Fatalf("AudioStreamStep returned %d", n)
			}
		}
	}

	if r := AudioStreamPush(stream, samples); r != AudioStreamOK {
		t.Fatalf("AudioStreamPush returned %d", r)
	}
	produced := drain()

	if r := AudioStreamFlush(stream, false); r != AudioStreamOK {
		t.Fatalf("AudioStreamFlush(false) returned %d", r)
	}
	produced += drain()

	flushed := AudioStreamGetNRowsFlushed(stream)
	rows := AudioStreamGetNRows(stream)
	t.Logf("rows %d, flushed %d, produced by step %d", rows, flushed, produced)
	if rows < flushed {
		t.Fatalf("rows %d did not reach the flush endpoint %d", rows, flushed)
	}
	if produced > rows {
		t.Fatalf("step produced %d rows, more than the counter %d", produced, rows)
	}

	if r := AudioStreamPush(stream, samples); r != AudioStreamOK {
		t.Fatalf("AudioStreamPush after a non ending flush returned %d", r)
	}
	if r := AudioStreamFlush(stream, true); r != AudioStreamOK {
		t.Fatalf("AudioStreamFlush(true) returned %d", r)
	}
	drain()

	if r := AudioStreamPush(stream, samples); r != AudioStreamEnded {
		t.Fatalf("AudioStreamPush after the end returned %d", r)
	}
	if n := AudioStreamStep(stream, lastToken); n != AudioStreamStepEnd {
		t.Fatalf("AudioStreamStep after the end returned %d", n)
	}
}
