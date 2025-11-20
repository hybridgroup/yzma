package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	modelsDir *string
	modelFile *string
	projFile  *string
	imageFile string
	prompt    *string
	template  string
	libPath   string
	verbose   *bool
)

func showUsage() {
	fmt.Println(`
Usage:
describe [image file path]`)
}

func handleFlags() error {
	modelsDir = flag.String("models", "", "models directory to use")
	modelFile = flag.String("model", "moondream2-text-model-f16_ct-vicuna.gguf", "model file to use")
	projFile = flag.String("mmproj", "moondream2-mmproj-f16-20250414.gguf", "projector file to use")
	prompt = flag.String("p", "Describe what is in this image.", "prompt")
	verbose = flag.Bool("v", false, "verbose logging")

	flag.Parse()

	if len(*modelFile) == 0 ||
		len(*projFile) == 0 ||
		len(*prompt) == 0 {

		return errors.New("missing a flag")
	}

	if len(libPath) == 0 && os.Getenv("YZMA_LIB") != "" {
		libPath = os.Getenv("YZMA_LIB")
	}

	if len(libPath) == 0 {
		return errors.New("missing YZMA_LIB env var")
	}

	if len(flag.Args()) == 0 {
		return errors.New("missing file")
	}

	if len(*modelsDir) == 0 {
		homedir, _ := os.UserHomeDir()
		*modelsDir = path.Join(homedir, "models")
	}

	imageFile = flag.Args()[0]

	return nil
}
