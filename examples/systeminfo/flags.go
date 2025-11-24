package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	libPath *string
)

func showUsage() {
	fmt.Println(`
Usage:
systeminfo -lib [llama.cpp .so file path]`)
}

func handleFlags() error {
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files (leave empty to use YZMA_LIB env var)")

	flag.Parse()

	if len(*libPath) == 0 && os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	return nil
}
