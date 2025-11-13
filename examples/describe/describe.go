package main

import (
	"fmt"
	"path"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
)

func describe(tmpFile string) error {
	llama.Init()
	defer llama.BackendFree()

	mtmdCtxParams := mtmd.ContextParamsDefault()

	switch {
	case *verbose:
		fmt.Println("Using model", path.Join(*modelsDir, *modelFile))
	default:
		llama.LogSet(llama.LogSilent())
		mtmdCtxParams.Verbosity = llama.LogLevelContinue
	}

	model, err := llama.ModelLoadFromFile(path.Join(*modelsDir, *modelFile), llama.ModelDefaultParams())
	if err != nil {
		return fmt.Errorf("unable to load model from file %s: %v", *modelFile, err)
	}
	if model == 0 {
		return fmt.Errorf("unable to load model from file %s", *modelFile)
	}

	defer llama.ModelFree(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 4096
	ctxParams.NBatch = 2048

	lctx, err := llama.InitFromModel(model, ctxParams)
	if err != nil {
		return fmt.Errorf("unable to initialize context from model: %v", err)
	}
	defer llama.Free(lctx)

	vocab := llama.ModelGetVocab(model)
	sampler := llama.NewSampler(model, llama.DefaultSamplers, llama.DefaultSamplerParams())

	mtmdCtx, err := mtmd.InitFromFile(path.Join(*modelsDir, *projFile), model, mtmdCtxParams)
	if err != nil {
		return fmt.Errorf("unable to initialize context from file: %v", err)
	}
	defer mtmd.Free(mtmdCtx)

	template = llama.ModelChatTemplate(model, "")
	messages = []llama.ChatMessage{llama.NewChatMessage("user", *prompt+mtmd.DefaultMarker())}
	output := mtmd.InputChunksInit()
	input := mtmd.NewInputText(chatTemplate(true), true, true)

	bitmap := mtmd.BitmapInitFromFile(mtmdCtx, tmpFile)
	defer mtmd.BitmapFree(bitmap)

	mtmd.Tokenize(mtmdCtx, output, input, []mtmd.Bitmap{bitmap})

	var n llama.Pos
	mtmd.HelperEvalChunks(mtmdCtx, lctx, output, 0, 0, int32(ctxParams.NBatch), true, &n)

	var sz int32 = 1
	batch := llama.BatchInit(1, 0, 1)
	batch.NSeqId = &sz
	batch.NTokens = 1
	seqs := unsafe.SliceData([]llama.SeqId{0})
	batch.SeqId = &seqs

	fmt.Println()

	for i := 0; i < llama.MaxToken; i++ {
		token := llama.SamplerSample(sampler, lctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			fmt.Println()
			break
		}

		buf := make([]byte, 128)
		l := llama.TokenToPiece(vocab, token, buf, 0, true)

		fmt.Print(string(buf[:l]))

		batch.Token = &token
		batch.Pos = &n

		llama.Decode(lctx, batch)
		n++
	}

	return nil
}

func chatTemplate(add bool) string {
	buf := make([]byte, 1024)
	len := llama.ChatApplyTemplate(template, messages, add, buf)
	result := string(buf[:len])
	return result
}
