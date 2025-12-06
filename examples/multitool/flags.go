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
  multitool -model [model file path] -lib [llama.cpp .so file path]

This example demonstrates multi-step tool/function calling with LLMs.
The model will be asked a question that requires multiple calculation steps,
and will use the available calculator tools to compute the answer.

Example questions that require multiple steps:
  - "What is (15 + 27) * 3?"
  - "Calculate 100 - 25, then multiply by 4"
  - "Add 10 and 20, then subtract 5"`)
}

func handleFlags() error {
	modelFile = flag.String("model", "", "model file to use")
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files")
	verbose = flag.Bool("v", false, "verbose logging")
	temperature = flag.Float64("temp", 0.1, "temperature for model")
	contextSize = flag.Int("c", 4096, "context size for model")
	predictSize = flag.Int("n", 512, "number of tokens to predict")
	userQuestion = flag.String("question", "What is (15 + 27) * 3?", "question to ask the model (use a multi-step calculation)")

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
