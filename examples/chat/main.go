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
	defer llama.BackendFree()

	model = llama.ModelLoadFromFile(*modelFile, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	vocab = llama.ModelGetVocab(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = uint32(*contextSize)
	ctxParams.NBatch = uint32(*batchSize)

	lctx = llama.InitFromModel(model, ctxParams)
	defer llama.Free(lctx)

	sampler = llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	if *topK != 0 {
		llama.SamplerChainAdd(sampler, llama.SamplerInitTopK(int32(*topK)))
	}
	if *topP < 1.0 {
		llama.SamplerChainAdd(sampler, llama.SamplerInitTopP(float32(*topP), 1))
	}
	if *minP > 0 {
		llama.SamplerChainAdd(sampler, llama.SamplerInitMinP(float32(*minP), 1))
	}
	llama.SamplerChainAdd(sampler, llama.SamplerInitTempExt(float32(*temperature), 0, 1.0))
	llama.SamplerChainAdd(sampler, llama.SamplerInitDist(llama.DefaultSeed))

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
	// call once to get the size
	count := llama.Tokenize(vocab, text, nil, first, true)

	// now get the actual tokens
	tokens := make([]llama.Token, count)
	llama.Tokenize(vocab, text, tokens, first, true)

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
