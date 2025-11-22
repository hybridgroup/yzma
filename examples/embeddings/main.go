package main

import (
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

func run() error {
	if err := handleFlags(); err != nil {
		showUsage()
		return err
	}

	if err := llama.Load(*libPath); err != nil {
		return fmt.Errorf("unable to load library: %w", err)
	}

	if !*verbose {
		llama.LogSet(llama.LogSilent())
	}

	llama.Init()
	defer llama.Close()

	model, err := llama.ModelLoadFromFile(*modelFile, llama.ModelDefaultParams())
	if err != nil {
		return fmt.Errorf("unable to load model from file %s: %v", *modelFile, err)
	}
	if model == 0 {
		return fmt.Errorf("unable to load model from file %s", *modelFile)
	}
	defer llama.ModelFree(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = uint32(*contextSize)
	ctxParams.NBatch = uint32(*batchSize)
	ctxParams.PoolingType = poolingType
	ctxParams.Embeddings = 1

	lctx, err := llama.InitFromModel(model, ctxParams)
	if err != nil {
		return fmt.Errorf("unable to initialize context from model: %v", err)
	}
	defer llama.Free(lctx)

	// tokenize prompt
	vocab := llama.ModelGetVocab(model)
	tokens := llama.Tokenize(vocab, *prompt, true, true)

	// create batch and decode
	batch := llama.BatchGetOne(tokens)
	llama.Decode(lctx, batch)

	// get embeddings
	nEmbd := llama.ModelNEmbd(model)
	vec, err := llama.GetEmbeddingsSeq(lctx, 0, nEmbd)
	if err != nil {
		return fmt.Errorf("unable to get embeddings: %v", err)
	}

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
