package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
)

var (
	vocab   llama.Vocab
	model   llama.Model
	lctx    llama.Context
	sampler llama.Sampler

	messages []llama.ChatMessage
)

func main() {
	if err := handleFlags(); err != nil {
		showUsage()
		os.Exit(0)
	}

	if err := llama.Load(*libPath); err != nil {
		fmt.Println("unable to load library", err.Error())
		os.Exit(1)
	}
	if err := mtmd.Load(*libPath); err != nil {
		fmt.Println("unable to load library", err.Error())
		os.Exit(1)
	}

	mctxParams := mtmd.ContextParamsDefault()
	if !*verbose {
		llama.LogSet(llama.LogSilent())
		mtmd.LogSet(llama.LogSilent())
	}

	llama.Init()
	defer llama.Close()

	fmt.Println("Loading model", *modelFile)
	model, err := llama.ModelLoadFromFile(*modelFile, llama.ModelDefaultParams())
	if err != nil {
		fmt.Println("unable to load model from file", err.Error())
		os.Exit(1)
	}
	if model == 0 {
		fmt.Println("unable to load model from file", *modelFile)
		os.Exit(1)
	}
	defer llama.ModelFree(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 4096
	ctxParams.NBatch = 2048

	lctx, err := llama.InitFromModel(model, ctxParams)
	if err != nil {
		fmt.Println("unable to initialize context from model", err.Error())
		os.Exit(1)
	}
	defer llama.Free(lctx)

	vocab := llama.ModelGetVocab(model)

	// pass in flags as params to samplers
	sp := llama.DefaultSamplerParams()
	sp.Temp = float32(*temperature)
	sp.TopK = int32(*topK)
	sp.TopP = float32(*topP)
	sp.MinP = float32(*minP)

	sampler := llama.NewSampler(model, llama.DefaultSamplers, sp)
	mtmdCtx, err := mtmd.InitFromFile(*projFile, model, mctxParams)
	if err != nil {
		fmt.Println("unable to initialize context from file", err.Error())
		os.Exit(1)
	}
	defer mtmd.Free(mtmdCtx)

	if *template == "" {
		*template = llama.ModelChatTemplate(model, "")
	}

	messages = make([]llama.ChatMessage, 0)
	if *systemPrompt != "" {
		messages = append(messages, llama.NewChatMessage("system", *systemPrompt))
	}
	messages = append(messages, llama.NewChatMessage("user", *prompt+mtmd.DefaultMarker()))

	output := mtmd.InputChunksInit()
	input := mtmd.NewInputText(chatTemplate(true), true, true)
	bitmap := mtmd.BitmapInitFromFile(mtmdCtx, *imageFile)
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
}

func chatTemplate(add bool) string {
	buf := make([]byte, 1024)
	len := llama.ChatApplyTemplate(*template, messages, add, buf)
	result := string(buf[:len])
	return result
}
