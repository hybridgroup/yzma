package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
)

var (
	modelFile    *string
	projFile     *string
	imageFile    *string
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
	if err := mtmd.Load(*libPath); err != nil {
		fmt.Println("unable to load library", err.Error())
		os.Exit(1)
	}

	llama.Init()
	defer llama.BackendFree()

	fmt.Println("Loading model", *modelFile)
	model := llama.ModelLoadFromFile(*modelFile, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 4096
	ctxParams.NBatch = 2048

	lctx := llama.InitFromModel(model, ctxParams)
	defer llama.Free(lctx)

	vocab := llama.ModelGetVocab(model)
	// TODO: pass in flags as params to samplers
	sampler := llama.NewSampler(model, llama.DefaultSamplers)
	mtmdCtx := mtmd.InitFromFile(*projFile, model, mtmd.ContextParamsDefault())
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
		llama.TokenToPiece(vocab, token, buf, 0, true)

		fmt.Print(string(buf))

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

func showUsage() {
	fmt.Println(`
Usage:
vlm -model [model file path] -proj [projector file path] -lib [llama.cpp .so file path] -prompt [what you want to ask] -image [image file path]`)
}

func handleFlags() error {
	modelFile = flag.String("model", "", "model file to use")
	projFile = flag.String("mmproj", "", "projector file to use")
	imageFile = flag.String("image", "", "image file to use")
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

	if len(*modelFile) == 0 ||
		len(*projFile) == 0 ||
		len(*prompt) == 0 ||
		len(*imageFile) == 0 {

		return errors.New("missing a flag")
	}

	if os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	return nil
}
