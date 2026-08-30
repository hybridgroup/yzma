//go:build js && wasm

// Package llamawasm runs llama.cpp inference in a browser.
//
// The [llama] package loads the llama.cpp shared libraries with libffi. A
// WebAssembly module cannot do this: it has no dlopen and no libffi, and TinyGo
// cannot compile the C++ of llama.cpp. This package therefore uses a different
// way to reach the same library. llama.cpp is a second WebAssembly module, made
// by Emscripten, and this package calls it through JavaScript.
//
// The names and the order of the calls are the same as in the [llama] package
// for the part of the API that this package has, so a program moves from one to
// the other with a change of the import:
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
// finds out if the page can use more than one thread, takes the correct
// llama.cpp module, and puts the result in globalThis.yzmaReady. [Load] waits
// for that promise.
//
// # Run it in a worker
//
// Every call into llama.cpp is synchronous, and one token takes milliseconds.
// A call from the main thread of a page stops the page. Put this code and the
// llama.cpp module in a Web Worker, and send the result to the page with
// postMessage. One [Decode] does one batch, so the worker can send each token
// as it comes.
//
// # One call at a time
//
// The package keeps scratch memory in the llama.cpp module for the calls that
// pass tokens and text, so only one goroutine may call into it at a time. This
// is not a limit in practice: llama.cpp itself takes one call at a time, and
// the generation loop is one goroutine.
//
// # WebGPU
//
// There are three builds of llama.cpp, and the JavaScript glue takes the best
// one that the browser can run: WebGPU, the CPU with more than one thread, or
// the CPU with one thread. [Backend] says which one it took, and [GPUDevice]
// names the GPU that llama.cpp found.
//
// A page can have WebGPU while llama.cpp has no device, because the backend
// needs an adapter with f16 shaders. Ask [GPUDevice], which reports what
// llama.cpp really has.
//
// Set NGpuLayers in [ModelParams] to put layers on the GPU. It does nothing in
// a build for the CPU.
//
// # Limits
//
// A build for the CPU uses SIMD. WebGPU needs Chrome or Edge 137 and later,
// because the backend waits for the GPU inside a synchronous call and that
// needs JavaScript Promise Integration.
//
// A WebAssembly module can address 4 GB, and one JavaScript ArrayBuffer holds
// at most 2 GB, so a model of more than 2 GB must be in splits.
//
// The package has the calls that text generation and embeddings need. It does
// not have multimodal input, LoRA adapters, saved state, or quantization.
package llamawasm
