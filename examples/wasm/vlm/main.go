//go:build js && wasm

// Vlm answers a question about an image in a browser.
//
// It is the same shape as examples/wasm/chat, with the multimodal calls of
// pkg/llamawasm added: the page decodes the image and sends the pixels, and this
// program turns them into a bitmap, puts them in the prompt beside the text, and
// then runs the same loop of sampling and decoding.
//
// Build it with TinyGo:
//
//	tinygo build -target wasm -o build/wasm/yzma-vlm.wasm ./examples/wasm/vlm
//
// See wasm/README.md for how to serve the result.
package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/hybridgroup/yzma/pkg/llamawasm"
)

const (
	modelPath   = "/models/model.gguf"
	projectPath = "/models/mmproj.gguf"
)

var (
	model   llamawasm.Model
	ctx     llamawasm.Context
	mctx    llamawasm.MtmdContext
	vocab   llamawasm.Vocab
	sampler llamawasm.Sampler
	nBatch  int32 = 2048
)

func main() {
	if err := llamawasm.Load(""); err != nil {
		post("error", err.Error())
		return
	}

	llamawasm.LogSet(llamawasm.LogSilent())
	llamawasm.Init()

	if !llamawasm.MtmdSupported() {
		post("error", "this build of llama.cpp has no multimodal calls")
		return
	}

	js.Global().Set("yzmaLoadModel", js.FuncOf(loadModel))
	js.Global().Set("yzmaOpenModel", js.FuncOf(openModel))
	js.Global().Set("yzmaDescribe", js.FuncOf(describe))

	post("ready", backendReport())

	<-make(chan struct{})
}

// loadModel(modelURL, projectorURL) gets both files and makes the contexts.
func loadModel(this js.Value, args []js.Value) any {
	if len(args) < 2 || !args[0].Truthy() || !args[1].Truthy() {
		post("error", "loadModel needs the URL of a model and of a projector")
		return nil
	}
	modelURL, projectorURL := args[0].String(), args[1].String()

	go func() {
		for _, file := range []struct {
			name, url, path string
		}{
			{"model", modelURL, modelPath},
			{"projector", projectorURL, projectPath},
		} {
			post("status", "downloading the "+file.name)

			err := llamawasm.FetchModelFile(file.path, file.url, func(done, total int64) {
				if total > 0 {
					post("progress", fmt.Sprintf("%s %d%%", file.name, done*100/total))
				}
			})
			if err != nil {
				post("error", err.Error())
				return
			}
		}

		open()
	}()

	return nil
}

// openModel() loads the files that are already in the filesystem of the module.
// A test that has them puts them there itself.
func openModel(this js.Value, args []js.Value) any {
	go open()
	return nil
}

// open loads the model, the context, and the projector.
func open() {
	post("status", "loading the model")

	params := llamawasm.ModelDefaultParams()
	onGPU := llamawasm.GPUDevice() != ""
	if onGPU {
		params.NGpuLayers = 999
	}

	var err error
	if model, err = llamawasm.ModelLoadFromFile(modelPath, params); err != nil {
		post("error", err.Error())
		return
	}

	// A model with eyes needs room for the tokens of the image as well as the
	// text, so the context is larger than the one the chat example uses.
	ctxParams := llamawasm.ContextDefaultParams()
	ctxParams.NCtx = 4096
	ctxParams.NBatch = uint32(nBatch)

	if ctx, err = llamawasm.InitFromModel(model, ctxParams); err != nil {
		post("error", err.Error())
		return
	}

	vocab = llamawasm.ModelGetVocab(model)

	post("status", "loading the projector")

	if mctx, err = llamawasm.MtmdInitFromFile(projectPath, model, 0, onGPU); err != nil {
		post("error", err.Error())
		return
	}

	if !llamawasm.MtmdSupportVision(mctx) {
		post("error", "this projector does not take images")
		return
	}

	post("loaded", llamawasm.ModelDesc(model)+", "+backendReport())
}

// describe(prompt, width, height, rgba, maxTokens) answers a question about an
// image. The pixels come from a canvas of the page, so they are RGBA.
func describe(this js.Value, args []js.Value) any {
	if len(args) < 4 {
		post("error", "describe needs a prompt, a size, and the pixels")
		return nil
	}

	prompt := args[0].String()
	width := int32(args[1].Int())
	height := int32(args[2].Int())

	rgba := make([]byte, args[3].Get("length").Int())
	js.CopyBytesToGo(rgba, args[3])

	maxTokens := int32(128)
	if len(args) > 4 && args[4].Truthy() {
		maxTokens = int32(args[4].Int())
	}

	go run(prompt, width, height, rgba, maxTokens)

	return nil
}

func run(prompt string, width, height int32, rgba []byte, maxTokens int32) {
	if model == 0 || ctx == 0 || mctx == 0 {
		post("error", "load a model first")
		return
	}

	// Every answer starts from an empty state.
	if err := llamawasm.MemoryClear(ctx, true); err != nil {
		post("error", err.Error())
		return
	}

	post("status", "looking at the image")

	bitmap, err := llamawasm.MtmdBitmapInit(width, height, dropAlpha(rgba))
	if err != nil {
		post("error", err.Error())
		return
	}
	defer llamawasm.MtmdBitmapFree(bitmap)

	chunks, err := llamawasm.MtmdInputChunksInit()
	if err != nil {
		post("error", err.Error())
		return
	}
	defer llamawasm.MtmdInputChunksFree(chunks)

	text := buildPrompt(prompt)

	if err := llamawasm.MtmdTokenize(mctx, chunks, text, true, true, []llamawasm.MtmdBitmap{bitmap}); err != nil {
		post("error", err.Error())
		return
	}

	// This is where the image goes through the projector and into the model. It
	// is the slow part, so it is timed on its own: putting it into the rate of
	// the tokens would say nothing about either.
	encodeStart := time.Now()

	nPast, err := llamawasm.MtmdHelperEvalChunks(mctx, ctx, chunks, 0, 0, nBatch, true)
	if err != nil {
		post("error", err.Error())
		return
	}

	encode := time.Since(encodeStart).Seconds()
	post("status", fmt.Sprintf("answering, %d tokens of image and text in %.1fs", nPast, encode))

	start := time.Now()

	if sampler != 0 {
		llamawasm.SamplerFree(sampler)
	}
	sampler = llamawasm.SamplerChainInit(llamawasm.SamplerChainDefaultParams())
	llamawasm.SamplerChainAdd(sampler, llamawasm.SamplerInitGreedy())

	buf := make([]byte, 64)

	var count int32
	for count < maxTokens {
		token := llamawasm.SamplerSample(sampler, ctx, -1)
		if llamawasm.VocabIsEOG(vocab, token) {
			break
		}
		llamawasm.SamplerAccept(sampler, token)

		if n := llamawasm.TokenToPiece(vocab, token, buf, 0, true); n > 0 {
			post("token", string(buf[:n]))
		}
		count++

		// The positions carry on from the chunks, which the module knows.
		if _, err := llamawasm.Decode(ctx, llamawasm.BatchGetOne([]llamawasm.Token{token})); err != nil {
			post("error", err.Error())
			return
		}
	}

	elapsed := time.Since(start).Seconds()
	if elapsed > 0 {
		post("done", fmt.Sprintf("%d tokens, %.2f tokens/s, and %.1fs for the image",
			count, float64(count)/elapsed, encode))
		return
	}
	post("done", fmt.Sprintf("%d tokens, and %.1fs for the image", count, encode))
}

// buildPrompt puts the marker of the model and the question into the chat format
// of the model. The marker is where the image goes.
func buildPrompt(question string) string {
	marker := llamawasm.MtmdMarker(mctx)
	if marker == "" {
		marker = "<__media__>"
	}

	text := marker + "\n" + question

	// A model with no chat template takes the text as it is.
	if formatted, err := llamawasm.ChatApplyTemplate(model, "user", text, true); err == nil && formatted != "" {
		return formatted
	}
	return text
}

// dropAlpha turns the RGBA of a canvas into the RGB that a bitmap needs.
func dropAlpha(rgba []byte) []byte {
	rgb := make([]byte, 0, len(rgba)/4*3)
	for i := 0; i+3 < len(rgba); i += 4 {
		rgb = append(rgb, rgba[i], rgba[i+1], rgba[i+2])
	}
	return rgb
}

func backendReport() string {
	if device := llamawasm.GPUDevice(); device != "" {
		return fmt.Sprintf("backend: %s (%s)", llamawasm.Backend(), device)
	}
	return fmt.Sprintf("backend: %s", llamawasm.Backend())
}

// post sends a message to whatever holds this module.
func post(kind, text string) {
	message := map[string]any{"kind": kind, "text": text}

	if fn := js.Global().Get("postMessage"); fn.Type() == js.TypeFunction {
		fn.Invoke(message)
		return
	}
	if fn := js.Global().Get("yzmaOnMessage"); fn.Type() == js.TypeFunction {
		fn.Invoke(message)
		return
	}
	fmt.Printf("%s: %s\n", kind, text)
}
