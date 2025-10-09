package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
)

var (
	modelFile    *string
	prompt       *string
	systemPrompt *string
	template     *string
	libPath      *string
	verbose      *bool

	temperature *float64
	topK        *int
	topP        *float64
	minP        *float64
	contextSize *int
	predictSize *int
	batchSize   *int

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
		llama.LogSet(llama.LogSilent(), uintptr(0))
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
	llama.SamplerChainAdd(sampler, llama.SamplerInitDist(llama.DEFAULT_SEED))

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
		if start == llama.TOKEN_NULL {
			start = llama.VocabBOS(vocab)
		}

		batch = llama.BatchGetOne([]llama.Token{start})
	}

	fmt.Println()

	response := ""
	for pos := int32(0); pos+batch.NTokens < int32(*predictSize); pos += batch.NTokens {
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

func showUsage() {
	fmt.Println(`
Usage:
chat -model [model file path] -lib [llama.cpp .so file path] -prompt [omit this flag for a chat session] -v`)
}

func handleFlags() error {
	modelFile = flag.String("model", "", "model file to use")
	prompt = flag.String("p", "", "prompt")
	systemPrompt = flag.String("sys", "", "system prompt")
	template = flag.String("template", "", "template name")
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files")
	verbose = flag.Bool("v", false, "verbose logging")

	temperature = flag.Float64("temp", 0.8, "temperature for model")
	topK = flag.Int("top-k", 40, "top-k for model")
	minP = flag.Float64("min-p", 0.1, "min-p for model")
	topP = flag.Float64("top-p", 0.9, "top-p for model")

	contextSize = flag.Int("c", 4096, "context size for model")
	predictSize = flag.Int("n", -1, "predict size for model")
	batchSize = flag.Int("b", 2048, "max batch size for model")

	flag.Parse()

	if len(*modelFile) == 0 {
		return errors.New("missing model flag")
	}

	if os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	if *predictSize < 0 {
		*predictSize = *contextSize //llama.MaxToken
	}

	return nil
}
