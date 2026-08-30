//go:build js && wasm

// Chat runs llama.cpp inference in a browser.
//
// The program is the same shape as examples/hello, but it uses the llamawasm
// package instead of the llama package, and it sends each piece of text to the
// page instead of printing it.
//
// Build it with TinyGo:
//
//	tinygo build -target wasm -o build/wasm/yzma.wasm ./examples/wasm/chat
//
// or with the standard toolchain:
//
//	GOOS=js GOARCH=wasm go build -o build/wasm/yzma.wasm ./examples/wasm/chat
//
// See wasm/README.md for how to serve the result.
package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/hybridgroup/yzma/pkg/llamawasm"
)

const modelPath = "/models/model.gguf"

var (
	model   llamawasm.Model
	ctx     llamawasm.Context
	vocab   llamawasm.Vocab
	sampler llamawasm.Sampler
)

func main() {
	if err := llamawasm.Load(""); err != nil {
		post("error", err.Error())
		return
	}

	llamawasm.LogSet(llamawasm.LogSilent())
	llamawasm.Init()

	// The page calls these.
	js.Global().Set("yzmaLoadModel", js.FuncOf(loadModel))
	js.Global().Set("yzmaOpenModel", js.FuncOf(openModel))
	js.Global().Set("yzmaGenerate", js.FuncOf(generate))

	post("ready", backendReport())

	// Keep the program alive so that the page can call into it.
	<-make(chan struct{})
}

// loadModel(url) gets a model over the network and makes a context for it.
func loadModel(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		post("error", "loadModel needs a URL")
		return nil
	}
	url := args[0].String()

	go func() {
		post("status", "downloading the model")

		err := llamawasm.FetchModelFile(modelPath, url, func(done, total int64) {
			if total > 0 {
				post("progress", fmt.Sprintf("%d%%", done*100/total))
			}
		})
		if err != nil {
			post("error", err.Error())
			return
		}

		open(modelPath)
	}()

	return nil
}

// openModel(path) loads a model that is already in the filesystem of the
// llama.cpp module. A test that has the file puts it there itself.
func openModel(this js.Value, args []js.Value) any {
	path := modelPath
	if len(args) > 0 && args[0].Truthy() {
		path = args[0].String()
	}

	go open(path)

	return nil
}

// open loads the model at path and makes a context for it.
func open(path string) {
	post("status", "loading the model")

	params := llamawasm.ModelDefaultParams()

	// A build with WebGPU has a device, so put every layer on it. A build on the
	// CPU has none, and the value does nothing there.
	if llamawasm.GPUDevice() != "" {
		params.NGpuLayers = 999
	}

	var err error
	if model, err = llamawasm.ModelLoadFromFile(path, params); err != nil {
		post("error", err.Error())
		return
	}

	ctxParams := llamawasm.ContextDefaultParams()
	ctxParams.NCtx = 2048
	ctxParams.NBatch = 512

	if ctx, err = llamawasm.InitFromModel(model, ctxParams); err != nil {
		post("error", err.Error())
		return
	}

	vocab = llamawasm.ModelGetVocab(model)

	post("loaded", llamawasm.ModelDesc(model)+", "+backendReport())
}

// backendReport says what does the computation.
func backendReport() string {
	if device := llamawasm.GPUDevice(); device != "" {
		return fmt.Sprintf("backend: %s (%s)", llamawasm.Backend(), device)
	}
	return fmt.Sprintf("backend: %s", llamawasm.Backend())
}

// generate(prompt, maxTokens) makes text and sends each piece to the page.
func generate(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		post("error", "generate needs a prompt")
		return nil
	}
	prompt := args[0].String()

	maxTokens := int32(128)
	if len(args) > 1 && args[1].Truthy() {
		maxTokens = int32(args[1].Int())
	}

	go func() {
		if model == 0 || ctx == 0 {
			post("error", "load a model first")
			return
		}

		// Each generation starts from an empty state.
		if err := llamawasm.MemoryClear(ctx, true); err != nil {
			post("error", err.Error())
			return
		}

		if sampler != 0 {
			llamawasm.SamplerFree(sampler)
		}
		sampler = llamawasm.SamplerChainInit(llamawasm.SamplerChainDefaultParams())
		llamawasm.SamplerChainAdd(sampler, llamawasm.SamplerInitGreedy())

		tokens := llamawasm.Tokenize(vocab, prompt, true, false)
		if len(tokens) == 0 {
			post("error", "the prompt has no tokens")
			return
		}

		batch := llamawasm.BatchGetOne(tokens)
		buf := make([]byte, 64)
		start := time.Now()

		var count int32
		for pos := int32(0); pos < maxTokens; pos += batch.NTokens {
			if _, err := llamawasm.Decode(ctx, batch); err != nil {
				post("error", err.Error())
				return
			}

			token := llamawasm.SamplerSample(sampler, ctx, -1)
			if llamawasm.VocabIsEOG(vocab, token) {
				break
			}
			llamawasm.SamplerAccept(sampler, token)

			if n := llamawasm.TokenToPiece(vocab, token, buf, 0, true); n > 0 {
				post("token", string(buf[:n]))
			}

			count++
			batch = llamawasm.BatchGetOne([]llamawasm.Token{token})
		}

		elapsed := time.Since(start).Seconds()
		if elapsed > 0 {
			post("done", fmt.Sprintf("%d tokens, %.2f tokens/s", count, float64(count)/elapsed))
			return
		}
		post("done", fmt.Sprintf("%d tokens", count))
	}()

	return nil
}

// post sends a message to whatever holds this module. In a worker that is the
// page, and in Node it is the console.
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
