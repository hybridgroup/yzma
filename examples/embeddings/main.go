package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

const usage = `
  -model string
    	path to the GGUF model
  -prompt string
    	prompt to embed
  -pooling string
    	pooling type for embeddings (mean, cls, none)
`

func run() error {
	var modelFlag, prompt, pooling string
	flag.StringVar(&modelFlag, "model", "", "path to the GGUF model")
	flag.StringVar(&prompt, "prompt", "", "prompt to embed")
	flag.StringVar(&pooling, "pooling", "mean", "pooling type for embeddings (mean, cls, none)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: embeddings %s\n", usage)
	}
	flag.Parse()
	if modelFlag == "" || prompt == "" {
		flag.Usage()
		return nil
	}
	var poolingType llama.PoolingType
	switch pooling {
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

	// load library and init
	llama.Load(os.Getenv("YZMA_LIB"))
	llama.Init()
	defer llama.BackendFree()

	// load model
	model := llama.ModelLoadFromFile(modelFlag, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	// create context
	params := llama.ContextDefaultParams()
	params.PoolingType = poolingType
	params.Embeddings = 1
	lctx := llama.InitFromModel(model, params)
	defer llama.Free(lctx)

	// tokenize prompt
	vocab := llama.ModelGetVocab(model)
	count := llama.Tokenize(vocab, prompt, nil, true, true)
	tokens := make([]llama.Token, count)
	llama.Tokenize(vocab, prompt, tokens, true, true)

	// create batch and decode
	batch := llama.BatchGetOne(tokens)
	llama.Decode(lctx, batch)

	// get embeddings
	nEmbd := llama.ModelNEmbd(model)
	vec := llama.GetEmbeddingsSeq(lctx, 0, nEmbd)

	// normalize embeddings
	var sum float64
	for _, v := range vec {
		sum += float64(v * v)
	}
	sum = math.Sqrt(sum)
	norm := float32(1.0 / sum)

	var b strings.Builder
	for i, v := range vec {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%f", v*norm))
	}
	fmt.Println(b.String())

	return nil
}
