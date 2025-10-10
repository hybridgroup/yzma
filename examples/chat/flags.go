package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
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
)

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
