// Compares the two batch paths of Decode, to find if the empty fields of
// BatchGetOne cause the macOS trap. The directory name starts with an
// underscore, so the Go tools keep it out of ./... builds.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
)

var (
	mode      = flag.String("mode", "getone", "batch mode, getone or full")
	modelFile = flag.String("model", "", "model file to use")
	libPath   = flag.String("lib", "", "path to llama.cpp compiled library files")
	verbose   = flag.Bool("v", false, "verbose logging")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
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
		return fmt.Errorf("unable to load model from file %s: %w", *modelFile, err)
	}
	defer llama.ModelFree(model)

	ctx, err := llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil {
		return fmt.Errorf("unable to initialize context from model: %w", err)
	}
	defer llama.Free(ctx)

	vocab := llama.ModelGetVocab(model)
	tokens := llama.Tokenize(vocab, "Hello world", true, true)

	batch, err := makeBatch(*mode, tokens)
	if err != nil {
		return err
	}

	ret, err := llama.Decode(ctx, batch)
	if err != nil {
		return fmt.Errorf("decode failed: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("decode returned non-zero: %d", ret)
	}

	fmt.Printf("mode %s decoded %d tokens\n", *mode, len(tokens))

	return nil
}

// makeBatch builds the batch two ways. BatchGetOne leaves Pos, NSeqId, SeqId
// and Logits empty, while Add sets all of them.
func makeBatch(mode string, tokens []llama.Token) (llama.Batch, error) {
	switch mode {
	case "getone":
		return llama.BatchGetOne(tokens), nil
	case "full":
		batch := llama.BatchInit(int32(len(tokens)), 0, 1)
		for i, token := range tokens {
			last := i == len(tokens)-1
			if err := batch.Add(token, llama.Pos(i), []llama.SeqId{0}, last); err != nil {
				return batch, fmt.Errorf("unable to add token %d: %w", i, err)
			}
		}
		return batch, nil
	default:
		return llama.Batch{}, fmt.Errorf("unknown mode %s", mode)
	}
}
