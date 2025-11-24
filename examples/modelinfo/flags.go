package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	libPath   *string
	modelFile *string
)

func showUsage() {
	fmt.Println(`
Usage:
modelinfo -lib [llama.cpp .so file path] -model [model file]`)
}

func handleFlags() error {
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files (leave empty to use YZMA_LIB env var)")
	modelFile = flag.String("model", "", "model file to use")

	flag.Parse()

	if len(*libPath) == 0 && os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	if len(*modelFile) == 0 {
		return errors.New("missing model flag")
	}

	return nil
}
