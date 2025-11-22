package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
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

	if !*verbose {
		llama.LogSet(llama.LogSilent())
	}

	llama.Init()
	defer llama.Close()

	mParams := llama.ModelDefaultParams()

	// handle Mixture of Experts (MoE) options
	switch {
	case *cmoe:
		overrides := []llama.TensorBuftOverride{llama.NewTensorBuftAllFFNExprsOverride()}
		mParams.SetTensorBufOverrides(overrides)
	case *ncmoe > 0:
		overrides := make([]llama.TensorBuftOverride, 0)
		for i := 0; i < *ncmoe; i++ {
			overrides = append(overrides, llama.NewTensorBuftBlockOverride(i))
		}
		mParams.SetTensorBufOverrides(overrides)
	}

	var err error
	model, err = llama.ModelLoadFromFile(*modelFile, mParams)
	if err != nil {
		fmt.Println("unable to load model from file", err.Error())
		os.Exit(1)
	}
	if model == 0 {
		fmt.Println("unable to load model from file", *modelFile)
		os.Exit(1)
	}

	defer llama.ModelFree(model)

	vocab = llama.ModelGetVocab(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = uint32(*contextSize)
	ctxParams.NBatch = uint32(*batchSize)
	ctxParams.NUbatch = uint32(*uBatchSize)

	lctx, err = llama.InitFromModel(model, ctxParams)
	if err != nil {
		fmt.Println("unable to initialize context from model", err.Error())
		os.Exit(1)
	}
	defer llama.Free(lctx)

	// pass in flags as params to samplers
	sp := llama.DefaultSamplerParams()
	sp.Temp = float32(*temperature)
	sp.TopK = int32(*topK)
	sp.TopP = float32(*topP)
	sp.MinP = float32(*minP)

	samplers := []llama.SamplerType{llama.SamplerTypeTopK, llama.SamplerTypeTopP, llama.SamplerTypeMinP, llama.SamplerTypeTemperature}
	sampler = llama.NewSampler(model, samplers, sp)

	if *template == "" {
		*template = llama.ModelChatTemplate(model, "")
	}
	if *template == "" {
		*template = "chatml"
	}

	messages = make([]llama.ChatMessage, 0)
	if *systemPrompt != "" {
		messages = append(messages, llama.NewChatMessage("system", *systemPrompt))
	}

	// single message
	if len(*prompt) > 0 {
		messages = append(messages, llama.NewChatMessage("user", *prompt))
		chat(chatTemplate(true), true)

		return
	}

	// chat session
	first := true
	for {
		fmt.Print("USER> ")
		reader := bufio.NewReader(os.Stdin)
		pmpt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("unable to read user input", err.Error())
			os.Exit(1)
		}

		messages = append(messages, llama.NewChatMessage("user", pmpt))
		chat(chatTemplate(true), first)
		first = false
	}
}

func chat(text string, first bool) {
	tokens := llama.Tokenize(vocab, text, first, true)

	batch := llama.BatchGetOne(tokens)

	if llama.ModelHasEncoder(model) {
		llama.Encode(lctx, batch)

		start := llama.ModelDecoderStartToken(model)
		if start == llama.TokenNull {
			start = llama.VocabBOS(vocab)
		}

		batch = llama.BatchGetOne([]llama.Token{start})
	}

	fmt.Println()

	response := ""
	for pos := int32(0); pos < int32(*predictSize); pos += batch.NTokens {
		llama.Decode(lctx, batch)
		token := llama.SamplerSample(sampler, lctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			fmt.Println()
			break
		}

		buf := make([]byte, 256)
		l := llama.TokenToPiece(vocab, token, buf, 0, false)
		next := string(buf[:l])

		batch = llama.BatchGetOne([]llama.Token{token})

		fmt.Print(next)
		response += next
	}

	fmt.Println()
}

func chatTemplate(add bool) string {
	buf := make([]byte, 1024)
	len := llama.ChatApplyTemplate(*template, messages, add, buf)
	result := string(buf[:len])
	return result
}
