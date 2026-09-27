// Package compare measures yzma against a model server. yzma calls llama.cpp
// in the same process. The servers answer over an OpenAI compatible REST
// interface. That interface is the only allowed difference. The model,
// prompt, image and sampler must be the same everywhere.
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
// embedding returns no tokens, so a request costs almost only transport, and
// that is what this suite measures.
const EmbedPrompt = "A llama farm sits high on the altiplano, where the air is thin and the grass is short."

// The prompt and image that every engine gets.
const (
	// The text prompt asks for a long answer. A short answer stops at end of
	// generation, and the engines disagree on that last token, so a one token
	// difference would move tokens a second by a tenth or more.
	TextPrompt  = "Write a long description of a llama farm."
	ImagePrompt = "What is in this image?"
	ImageFile   = "../../images/domestic_llama.jpg"
)

// Request is one generation. Every engine gets the same one.
type Request struct {
	Prompt    string
	Image     []byte // image bytes, empty for the text suite
	MaxTokens int
	Seed      uint32

	// Embeddings means this request wants a vector and not tokens.
	Embeddings bool
}

// Result is the outcome of one generation.
type Result struct {
	Text         string
	Tokens       int           // tokens the engine generated, without the prompt
	PromptTokens int           // prompt tokens, including the image
	Dimensions   int           // vector size for the embeddings suite
	FirstToken   time.Duration // from the call to the first token
	Total        time.Duration // the whole request
}

// The prompt token count shows whether the engines do the same work. The image
// makes up most of it, and a projector that splits the image into tiles
// produces three times the tokens of one that doesn't.

// Engine generates text with a model.
type Engine interface {
	Name() string

	// Load prepares the engine and loads the model into memory. A timed run
	// must never pay for the load.
	Load(Request) error

	Generate(Request) (Result, error)

	// Embed turns the prompt into a vector. Result holds no generated tokens,
	// and FirstToken equals Total because there is no stream.
	Embed(Request) (Result, error)

	Close()
}

// A server keeps the last request's prompt and image. A second request with
// the same bytes then costs almost nothing, while yzma clears its cache after
// each generation. So to compare the engines fairly, each run must get an
// image and a prompt that no engine has seen.

// EmbedVariant returns the text for one embeddings suite run. Each one is
// different, so no engine answers from a cache.
func EmbedVariant(n int) string {
	return fmt.Sprintf("Request %d. %s", n, EmbedPrompt)
}

// TextVariant returns the prompt for one run. Each one is different, so no
// engine answers from a cache.
// The number goes first. A server reuses the tokens a prompt shares with the
// last one, so a common prefix would give it most of the prompt for free while
// yzma reads all of it again.
func TextVariant(n int) string {
	return fmt.Sprintf("Request %d. %s", n, TextPrompt)
}

// Images makes a different image for each run from one file. Each has the
// file's size and content plus a small unique block of pixels. So the vision
// model does the same work for each, and no server cache can answer for another.
type Images struct {
	src image.Image
}

// NewImages reads the file. A size other than zero scales it first.
func NewImages(path string, size image.Point) (*Images, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	if size.X > 0 && size.Y > 0 {
		src = scale(src, size)
	}

	return &Images{src: src}, nil
}

// Variant returns the image for run n. No two runs get the same bytes.
func (im *Images) Variant(n int) ([]byte, error) {
	bounds := im.src.Bounds()
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, im.src, bounds.Min, draw.Src)

	// An 8 by 8 pixel block in a corner encodes the run number, with a white
	// pixel for each set bit. JPEG preserves black and white edges, and the
	// block is too small to change what the model sees.
	for bit := range 64 {
		shade := color.RGBA{A: 255}
		if uint64(n)>>bit&1 == 1 {
			shade = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		canvas.Set(bounds.Min.X+bit%8, bounds.Min.Y+bit/8, shade)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, canvas, &jpeg.Options{Quality: 95}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// scale returns the image at the given size, with bilinear sampling. Each
// engine scales images its own way, so an image none of them resizes makes
// them do the same work.
func scale(src image.Image, size image.Point) image.Image {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))

	for y := range size.Y {
		fy := (float64(y)+0.5)*float64(b.Dy())/float64(size.Y) - 0.5
		y0 := clamp(int(fy), b.Dy()-1)
		y1 := clamp(y0+1, b.Dy()-1)
		wy := max(fy-float64(y0), 0)

		for x := range size.X {
			fx := (float64(x)+0.5)*float64(b.Dx())/float64(size.X) - 0.5
			x0 := clamp(int(fx), b.Dx()-1)
			x1 := clamp(x0+1, b.Dx()-1)
			wx := max(fx-float64(x0), 0)

			var out [4]float64
			for _, p := range [4]struct {
				x, y int
				w    float64
			}{
				{x0, y0, (1 - wx) * (1 - wy)},
				{x1, y0, wx * (1 - wy)},
				{x0, y1, (1 - wx) * wy},
				{x1, y1, wx * wy},
			} {
				r, g, bl, a := src.At(b.Min.X+p.x, b.Min.Y+p.y).RGBA()
				out[0] += float64(r) * p.w
				out[1] += float64(g) * p.w
				out[2] += float64(bl) * p.w
				out[3] += float64(a) * p.w
			}
			dst.SetRGBA(x, y, color.RGBA{
				R: uint8(out[0] / 257), G: uint8(out[1] / 257),
				B: uint8(out[2] / 257), A: uint8(out[3] / 257),
			})
		}
	}

	return dst
}

func clamp(v, most int) int {
	return min(max(v, 0), most)
}
