// Package compare measures yzma against a model server. yzma calls llama.cpp
// in the same process. The servers answer over an OpenAI compatible REST
// interface. That interface is the only permitted difference. The model, the
// prompt, the image and the sampler must be the same everywhere.
package compare

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"time"

	_ "image/png" // the test image can be a PNG
)

// EmbedPrompt is the text that the embeddings suite turns into a vector. An
// embedding gives no tokens back, thus the cost of a request is almost all
// transport, and that is what this suite measures.
const EmbedPrompt = "A llama farm sits high on the altiplano, where the air is thin and the grass is short."

// The prompt and the image that every engine gets.
const (
	// The text prompt asks for a long answer. A short answer stops at the end
	// of generation, and the engines do not agree on that last token, thus one
	// token of difference would move tokens a second by a tenth or more.
	TextPrompt  = "Write a long description of a llama farm."
	ImagePrompt = "What is in this image?"
	ImageFile   = "../../images/domestic_llama.jpg"
)

// Request is one generation. Every engine gets the same one.
type Request struct {
	Prompt    string
	Image     []byte // the bytes of an image, empty for the text suite
	MaxTokens int
	Seed      uint32

	// Embeddings says that this request wants a vector and not tokens.
	Embeddings bool
}

// Result is what one generation gives.
type Result struct {
	Text         string
	Tokens       int           // tokens that the engine made, without the prompt
	PromptTokens int           // tokens of the prompt, with the image
	Dimensions   int           // size of the vector of the embeddings suite
	FirstToken   time.Duration // from the call to the first token
	Total        time.Duration // the whole request
}

// The count of the prompt tokens says if the engines do the same work. An
// image gives most of them, and a projector that splits the image into tiles
// gives three times the tokens of one that does not.

// Engine makes text with a model.
type Engine interface {
	Name() string

	// Load makes the engine ready and puts the model in memory. A timed run
	// must never pay for the load.
	Load(Request) error

	Generate(Request) (Result, error)

	// Embed turns the prompt into a vector. Result holds no made tokens, thus
	// FirstToken and Total are the same, because there is no stream.
	Embed(Request) (Result, error)

	Close()
}

// A server keeps the prompt of the last request, and the image with it. A
// second request with the same bytes then costs almost nothing, while yzma
// empties its cache after each generation. To compare the engines, each run
// must therefore get an image and a prompt that no engine has seen.

// EmbedVariant gives the text of one run of the embeddings suite. Each one is
// different, thus no engine answers from a cache.
func EmbedVariant(n int) string {
	return fmt.Sprintf("Request %d. %s", n, EmbedPrompt)
}

// TextVariant gives the prompt of one run. Each one is different, thus no
// engine answers from a cache.
// The number goes first. A server keeps a prompt that begins as the last one
// did and reuses those tokens, thus a common start would give it most of the
// prompt for free while yzma reads all of it again.
func TextVariant(n int) string {
	return fmt.Sprintf("Request %d. %s", n, TextPrompt)
}

// ImageVariants makes count images from one file. Each one has the size and
// the content of the first, with a small block of pixels that only it has.
// The vision model therefore does the same work for each, and no cache of a
// server can answer for another.
func ImageVariants(path string, count int) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	variants := make([][]byte, 0, count)
	for n := range count {
		bounds := src.Bounds()
		canvas := image.NewRGBA(bounds)
		draw.Draw(canvas, bounds, src, bounds.Min, draw.Src)

		// One block of 8 by 8 pixels in a corner carries the number of the
		// run. It is too small to change what the model sees.
		shade := color.RGBA{R: uint8(n * 7), G: uint8(n * 11), B: uint8(n * 13), A: 255}
		for y := bounds.Min.Y; y < bounds.Min.Y+8 && y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Min.X+8 && x < bounds.Max.X; x++ {
				canvas.Set(x, y, shade)
			}
		}

		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, canvas, &jpeg.Options{Quality: 95}); err != nil {
			return nil, err
		}
		variants = append(variants, buf.Bytes())
	}

	return variants, nil
}
