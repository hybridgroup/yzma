package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
)

var (
	modelFile *string
	prompt    *string
	libPath   *string
	verbose   *bool

	contextSize *int
	predictSize *int
	batchSize   *int

	pooling     *string
	poolingType llama.PoolingType
)

func showUsage() {
	fmt.Println(`
Usage:
embeddings -model [model file path] -lib [llama.cpp .so file path] -prompt [omit this flag for a chat session] -v`)
}

func handleFlags() error {
	modelFile = flag.String("model", "", "model file to use")
	prompt = flag.String("p", "", "prompt")
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files")
	verbose = flag.Bool("v", false, "verbose logging")

	contextSize = flag.Int("c", 4096, "context size for model")
	predictSize = flag.Int("n", -1, "predict size for model")
	batchSize = flag.Int("b", 2048, "max batch size for model")

	pooling = flag.String("pooling", "mean", "pooling type for embeddings (mean, cls, rank, last, none)")

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

	if *predictSize < 0 {
		*predictSize = *contextSize
	}

	switch *pooling {
	case "mean":
		poolingType = llama.PoolingTypeMean
	case "cls":
		poolingType = llama.PoolingTypeCLS
	case "rank":
		poolingType = llama.PoolingTypeRank
	case "last":
		poolingType = llama.PoolingTypeLast
	case "none":
		poolingType = llama.PoolingTypeNone
	default:
		poolingType = llama.PoolingTypeUnspecified
	}

	return nil
}
