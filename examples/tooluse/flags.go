package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	modelFile    *string
	libPath      *string
	verbose      *bool
	temperature  *float64
	contextSize  *int
	predictSize  *int
	userQuestion *string
)

func showUsage() {
	fmt.Println(`
Usage:
  tools -model [model file path] -lib [llama.cpp .so file path]

This example demonstrates tool/function calling with LLMs.
The model will be asked to calculate "What is 15 + 27?" and will use
the available calculator tools to compute the answer.`)
}

func handleFlags() error {
	modelFile = flag.String("model", "", "model file to use")
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files")
	verbose = flag.Bool("v", false, "verbose logging")
	temperature = flag.Float64("temp", 0.1, "temperature for model")
	contextSize = flag.Int("c", 4096, "context size for model")
	predictSize = flag.Int("n", 512, "number of tokens to predict")
	userQuestion = flag.String("question", "What is 15 + 27?", "question to ask the model")

	flag.Parse()

	if len(*modelFile) == 0 {
		return errors.New("missing model flag")
	}

	if len(*libPath) == 0 && os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	return nil
}
