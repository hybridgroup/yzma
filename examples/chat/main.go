package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

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
	var tensorBuftBuf []llama.TensorBuftOverride
	switch {
	case *cmoe:
		tensorBuftBuf = []llama.TensorBuftOverride{llama.NewTensorBuftAllFFNExprsOverride(), {}} // sentinel-terminated
		if err := mParams.SetTensorBufOverrides(tensorBuftBuf); err != nil {
			fmt.Println("SetTensorBufOverrides failed:", err)
			os.Exit(1)
		}
	case *ncmoe > 0:
		tensorBuftBuf = make([]llama.TensorBuftOverride, 0, *ncmoe+1)
		for i := 0; i < *ncmoe; i++ {
			tensorBuftBuf = append(tensorBuftBuf, llama.NewTensorBuftBlockOverride(i))
		}
		tensorBuftBuf = append(tensorBuftBuf, llama.TensorBuftOverride{}) // sentinel
		if err := mParams.SetTensorBufOverrides(tensorBuftBuf); err != nil {
			fmt.Println("SetTensorBufOverrides failed:", err)
			os.Exit(1)
		}
	}

	var err error
	model, err = llama.ModelLoadFromFile(*modelFile, mParams)
	runtime.KeepAlive(tensorBuftBuf)
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
	if *threads > 0 {
		ctxParams.NThreads = int32(*threads)
		ctxParams.NThreadsBatch = int32(*threads)
	}

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
		turn()

		return
	}

	// chat session
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("USER> ")
		pmpt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("unable to read user input", err.Error())
			os.Exit(1)
		}

		messages = append(messages, llama.NewChatMessage("user", pmpt))
		turn()
	}
}

// decoded is the formatted history already in the context.
var decoded string

// turn decodes the part of the history that is not in the context yet, then
// generates the reply and adds it to the history.
func turn() {
	formatted, ok := chatTemplate(true)
	if !ok {
		fmt.Println("unable to apply chat template", *template)
		os.Exit(1)
	}

	// Start again from an empty context when the new history does not extend the old one.
	text, ok := strings.CutPrefix(formatted, decoded)
	if !ok || decoded == "" || llama.ModelHasEncoder(model) {
		clearMemory()
		text = formatted
		decoded = ""
	}

	response := chat(text, decoded == "")
	messages = append(messages, llama.NewChatMessage("assistant", response))

	if decoded, ok = chatTemplate(false); !ok {
		decoded = ""
	}
}

func clearMemory() {
	mem, err := llama.GetMemory(lctx)
	if err != nil {
		fmt.Println("unable to get memory", err.Error())
		os.Exit(1)
	}
	if err := llama.MemoryClear(mem, true); err != nil {
		fmt.Println("unable to clear memory", err.Error())
		os.Exit(1)
	}
}

func chat(text string, first bool) string {
	tokens := llama.Tokenize(vocab, text, first, true)

	batch := llama.BatchGetOne(tokens)

	if llama.ModelHasEncoder(model) {
		if _, err := llama.Encode(lctx, batch); err != nil {
			fmt.Println("unable to encode", err.Error())
			os.Exit(1)
		}

		start := llama.ModelDecoderStartToken(model)
		if start == llama.TokenNull {
			start = llama.VocabBOS(vocab)
		}

		batch = llama.BatchGetOne([]llama.Token{start})
	}

	fmt.Println()

	response := ""
	buf := make([]byte, 256)
	for pos := int32(0); pos < int32(*predictSize); pos += batch.NTokens {
		if _, err := llama.Decode(lctx, batch); err != nil {
			fmt.Println("unable to decode, the context may be full:", err.Error())
			os.Exit(1)
		}
		token := llama.SamplerSample(sampler, lctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			fmt.Println()
			break
		}

		l := llama.TokenToPiece(vocab, token, buf, 0, false)
		if l < 0 {
			// A negative result is the size the piece needs.
			buf = make([]byte, -l)
			l = llama.TokenToPiece(vocab, token, buf, 0, false)
		}
		next := ""
		if l > 0 && int(l) <= len(buf) {
			next = string(buf[:l])
		}

		batch = llama.BatchGetOne([]llama.Token{token})

		fmt.Print(next)
		response += next
	}

	fmt.Println()

	return response
}

// chatTemplate formats messages and reports false when llama.cpp does not know the template.
func chatTemplate(add bool) (string, bool) {
	buf := make([]byte, 4096)
	n := llama.ChatApplyTemplate(*template, messages, add, buf)
	if int(n) > len(buf) {
		// The result is the full length even when buf is too short.
		buf = make([]byte, n)
		n = llama.ChatApplyTemplate(*template, messages, add, buf)
	}
	if n < 0 || int(n) > len(buf) {
		return "", false
	}
	return string(buf[:n]), true
}
