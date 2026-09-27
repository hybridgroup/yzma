//go:build js && wasm

// Package llamawasm runs llama.cpp inference in a browser.
//
// The [llama] package loads the llama.cpp shared libraries with libffi. A
// WebAssembly module cannot do this. It has no dlopen and no libffi, and TinyGo
// cannot compile the llama.cpp C++ code. So llama.cpp is a second WebAssembly
// module, built by Emscripten, and this package calls it through JavaScript.
//
// Call names and order match the [llama] package for the part of the API that
// this package covers. So a program can switch packages by changing the import.
//
//	llamawasm.Load("")
//	llamawasm.LogSet(llamawasm.LogSilent())
//	llamawasm.Init()
//
//	model, err := llamawasm.ModelLoadFromFile("model.gguf", llamawasm.ModelDefaultParams())
//	ctx, err := llamawasm.InitFromModel(model, llamawasm.ContextDefaultParams())
//	vocab := llamawasm.ModelGetVocab(model)
//
//	tokens := llamawasm.Tokenize(vocab, prompt, true, false)
//	batch := llamawasm.BatchGetOne(tokens)
//
//	sampler := llamawasm.SamplerChainInit(llamawasm.SamplerChainDefaultParams())
//	llamawasm.SamplerChainAdd(sampler, llamawasm.SamplerInitGreedy())
//
// # What the page must do first
//
// The JavaScript glue in the yzma wasm directory must run before [Load]. It
// checks whether the page can use more than one thread, selects the right
// llama.cpp module, and puts the result in globalThis.yzmaReady. [Load] waits
// for that promise.
//
// # Run it in a worker
//
// Each call into llama.cpp is synchronous and one token takes milliseconds.
// A call from the main thread blocks the page. Put this code and the llama.cpp
// module in a Web Worker, and send results to the page with postMessage. One
// [Decode] runs one batch, so the worker can send each token right away.
//
// # One call at a time
//
// The package keeps scratch memory in the llama.cpp module for calls that
// pass tokens and text. So only one goroutine can call into it at a time.
// This is not a limit in practice, because llama.cpp accepts one call at a time
// and the generation loop is one goroutine.
//
// # WebGPU
//
// There are three builds of llama.cpp. The JavaScript glue selects the best one
// that the browser can run, which is WebGPU, the CPU with more than one thread,
// or the CPU with one thread. [Backend] returns the selection and [GPUDevice]
// returns the name of the GPU that llama.cpp found.
//
// A page can have WebGPU while llama.cpp has no device, because the backend
// needs an adapter with f16 shaders. Use [GPUDevice], which reports what
// llama.cpp actually has.
//
// Set NGpuLayers in [ModelParams] to put layers on the GPU. A CPU build ignores
// this value.
//
// # Images
//
// The llama.cpp multimodal library, mtmd, is in every build. [MtmdBitmapInit]
// takes the pixels of an image, [MtmdTokenize] puts them into a prompt with the
// text, and [MtmdHelperEvalChunks] runs both through the model. See mtmd.go for
// the call order.
//
// The pixels must be RGB. A page decodes the image with a canvas, so any
// format the browser reads works and the build needs no image library.
//
// Images only. Audio needs the page to decode and resample the samples, and
// video needs ffmpeg in a subprocess.
//
// # Chat templates and tool calling
//
// [ModelChatTemplate] returns the template stored in the GGUF. The yzma template
// and message packages are pure Go, so they render a multi turn conversation
// and parse the tool calls that come back. See examples/wasm/tools.
//
// [ChatApplyTemplate] takes only one message, which is enough for a question
// about an image.
//
// # Limits
//
// A CPU build uses SIMD. WebGPU needs Chrome or Edge 137 or later, or Firefox
// 153 or later, because the backend waits for the GPU in a synchronous call and
// that needs JavaScript Promise Integration. The loader picks the CPU in
// Firefox, which is faster there. Chrome on Linux needs switches for Vulkan,
// see wasm/README.md.
//
// A WebAssembly module can address 4 GB and one JavaScript ArrayBuffer holds at
// most 2 GB. So a model larger than 2 GB must be split.
//
// The package has the calls that text generation, embeddings, and images need,
// and it saves context state in memory. It does not have audio, video,
// LoRA adapters, state in a file, or quantization.
//
// The shim has no end of turn token and no grammar sampler. So the message
// package StopMarkers approximates, and a grammar cannot force a tool call.
package llamawasm
