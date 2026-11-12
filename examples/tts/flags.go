package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
)

var (
	modelFile  *string
	projFile   *string
	prompt     *string
	outputFile *string
	voice      *string
	lang       *string
	libPath    *string
	verbose    *bool

	temperature *float64
	topK        *int
	topP        *float64
	minP        *float64
	seed        *uint
	contextSize *int
	batchSize   *int
	maxFrames   *int
)

func showUsage() {
	fmt.Println(`
Usage:
tts -model [model file path] -mmproj [projector file path] -lib [llama.cpp .so file path] -p [text to speak] -o [output WAV file]`)
}

func handleFlags() error {
	modelFile = flag.String("model", "", "model file to use")
	projFile = flag.String("mmproj", "", "projector file to use")
	prompt = flag.String("p", "", "text to speak")
	outputFile = flag.String("o", "output.wav", "output WAV file")
	voice = flag.String("voice", "", "preset voice, the first one of the model if empty")
	lang = flag.String("lang", "", "language of the text, if the model supports it")
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files")
	verbose = flag.Bool("v", false, "verbose logging")

	temperature = flag.Float64("temp", 0.8, "temperature for the backbone")
	topK = flag.Int("top-k", 40, "top-k for the backbone")
	minP = flag.Float64("min-p", 0.05, "min-p for the backbone")
	topP = flag.Float64("top-p", 0.95, "top-p for the backbone")
	seed = flag.Uint("seed", math.MaxUint32, "random seed, the default gives a random one")

	contextSize = flag.Int("c", 4096, "context size for model")
	batchSize = flag.Int("b", 2048, "max batch size for model")
	maxFrames = flag.Int("n", 512, "max number of audio frames")

	flag.Parse()

	if len(*modelFile) == 0 ||
		len(*projFile) == 0 ||
		len(*prompt) == 0 {

		return errors.New("missing a flag")
	}

	if len(*libPath) == 0 && os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	return nil
}
