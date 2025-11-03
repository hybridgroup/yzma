package main

import (
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
)

var (
	libPath = os.Getenv("YZMA_LIB")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <modelFile>\n", os.Args[0])
		os.Exit(1)
	}
	modelFile := os.Args[1]

	llama.Load(libPath)
	llama.LogSet(llama.LogSilent())

	llama.Init()

	model := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	defer llama.ModelFree(model)

	desc := llama.ModelDesc(model)
	fmt.Printf("Model Description: %s\n", desc)

	size := llama.ModelSize(model)
	fmt.Printf("Model Size: %d tensors\n", size)

	encoder := llama.ModelHasEncoder(model)
	fmt.Printf("Model Has Encoder: %v\n", encoder)

	decoder := llama.ModelHasDecoder(model)
	fmt.Printf("Model Has Decoder: %v\n", decoder)

	recurrent := llama.ModelIsRecurrent(model)
	fmt.Printf("Model Is Recurrent: %v\n", recurrent)

	hybrid := llama.ModelIsHybrid(model)
	fmt.Printf("Model Is Hybrid: %v\n", hybrid)

	count := llama.ModelMetaCount(model)
	fmt.Printf("Model Metadata (%d entries):\n", count)
	for i := int32(0); i < count; i++ {
		key, ok := llama.ModelMetaKeyByIndex(model, i)
		if !ok {
			fmt.Printf("Error getting key for index %d\n", i)
			continue
		}
		value, ok := llama.ModelMetaValStrByIndex(model, i)
		if !ok {
			fmt.Printf("Error getting value for index %d\n", i)
			continue
		}
		fmt.Printf("  %s: %s\n", key, value)
	}
}
