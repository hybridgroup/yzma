//go:build js && wasm

// Package llamawasm runs llama.cpp inference in a browser.
//
// The [llama] package loads the llama.cpp shared libraries with libffi. A
// WebAssembly module cannot do this. It has no dlopen and no libffi, and TinyGo
// cannot compile the C++ of llama.cpp. Thus llama.cpp is a second WebAssembly
// module, made by Emscripten, and this package calls it through JavaScript.
//
// The names and the order of the calls agree with the [llama] package for the
// part of the API that this package has. Thus a program moves from one package
// to the other with a change of the import.
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
// The JavaScript glue in the wasm directory of yzma must run before [Load]. It
// finds out if the page can use more than one thread, selects the correct
// llama.cpp module, and puts the result in globalThis.yzmaReady. [Load] waits
// for that promise.
//
// # Run it in a worker
//
// Each call into llama.cpp is synchronous and one token takes milliseconds.
// A call from the main thread stops the page. Put this code and the llama.cpp
// module in a Web Worker, and send the result to the page with postMessage. One
// [Decode] does one batch, thus the worker can send each token immediately.
//
// # One call at a time
//
// The package keeps scratch memory in the llama.cpp module for the calls that
// pass tokens and text. Thus only one goroutine can call into it at a time.
// This is not a limit in practice, because llama.cpp accepts one call at a time
// and the generation loop is one goroutine.
//
// # WebGPU
//
// There are three builds of llama.cpp. The JavaScript glue selects the best one
// that the browser can run, which is WebGPU, the CPU with more than one thread,
// or the CPU with one thread. [Backend] gives the selection and [GPUDevice]
// gives the name of the GPU that llama.cpp found.
//
// A page can have WebGPU while llama.cpp has no device, because the backend
// needs an adapter with f16 shaders. Use [GPUDevice], which gives the true
// condition of llama.cpp.
//
// Set NGpuLayers in [ModelParams] to put layers on the GPU. A CPU build ignores
// this value.
//
// # Images
//
// The multimodal library of llama.cpp, mtmd, is in each build. [MtmdBitmapInit]
// takes the pixels of an image, [MtmdTokenize] puts them into a prompt with the
// text, and [MtmdHelperEvalChunks] runs both through the model. See mtmd.go for
// the order of the calls.
//
// The pixels must be RGB. A page decodes the image with a canvas, thus each
// format that the browser reads is usable and the build needs no image library.
//
// Images only. Audio needs the page to decode and resample the samples, and
// video needs ffmpeg in a subprocess.
//
// # Chat templates and tool calling
//
// [ModelChatTemplate] gives the template that the GGUF holds. The template and
// message packages of yzma are pure Go, thus they render a conversation with
// turns and parse the tool calls that come back. See examples/wasm/tools.
//
// [ChatApplyTemplate] takes one message only, which is sufficient for a
// question about an image.
//
// # Limits
//
// A CPU build uses SIMD. WebGPU needs Chrome or Edge 137 or later, because the
// backend waits for the GPU in a synchronous call and that needs JavaScript
// Promise Integration.
//
// A WebAssembly module can address 4 GB and one JavaScript ArrayBuffer holds a
// maximum of 2 GB. Thus a model of more than 2 GB must be in splits.
//
// The package has the calls that text generation, embeddings, and images need.
// It does not have audio, video, LoRA adapters, saved state, or quantization.
//
// The shim gives no end of turn token and no grammar sampler. Thus StopMarkers
// of the message package approximates, and a grammar cannot force a tool call.
package llamawasm
