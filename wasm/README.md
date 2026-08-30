# yzma in WebAssembly

This directory holds what a browser needs to run yzma: the JavaScript glue that
takes the correct WebAssembly build of llama.cpp, a Web Worker that holds the
program, a page, a static server, and a test that runs in Node.

The Go code is in [`pkg/llamawasm`](../pkg/llamawasm), and the example is in
[`examples/wasm/chat`](../examples/wasm/chat).

## How it works

There are two WebAssembly modules:

```
   page (index.html)
         |  postMessage
   Web Worker
     |            \
   Go program      llama.cpp module
   (TinyGo)   -->  (Emscripten)
       through JavaScript
```

On a native platform yzma calls llama.cpp with libffi and loads the shared
libraries at run time. A WebAssembly module has no `dlopen` and no libffi, and
TinyGo cannot compile the C++ of llama.cpp. So llama.cpp becomes a second
WebAssembly module, made by Emscripten with a small C shim, and the Go code
calls that module through JavaScript. The shim is in the
[llama-cpp-builder](https://github.com/hybridgroup/llama-cpp-builder) repo, in
its `wasm` directory.

The generation loop stays in Go. One `Decode` and one `SamplerSample` for each
token, the same as the `examples/hello` program.

## Files

| File | What it does |
| --- | --- |
| `yzma-loader.js` | Tests if the page can use more than one thread, then loads the correct llama.cpp module and puts it in `globalThis.yzmaReady`. |
| `worker.js` | Runs llama.cpp and the Go program in a Web Worker, and sends each piece of text to the page. |
| `index.html` | A page that loads a model and makes text. |
| `serve/main.go` | A static server that sets the headers that a build with more than one thread needs. |
| `node/run.js` | Runs the same build in Node, with no browser. This is the test that CI uses. |

## Build and run

```
# get the WebAssembly build of llama.cpp
make download-llama.cpp-wasm

# build the program with TinyGo
make wasm-example

# serve it
make serve-wasm
```

Then open <http://localhost:8080>.

`make wasm-example-go` builds the same program with the standard Go toolchain.
The binary is larger, which is useful if TinyGo cannot build a dependency.

## Run it without a browser

```
make wasm-example
make test-wasm
```

`node/run.js` loads a small model, makes tokens with the greedy sampler, and
prints them. The greedy sampler always takes the most probable token, so the
output does not change from one run to the next and a test can compare it.

`make test-wasm-mt` does the same with the build that uses more than one thread.
Node gives `SharedArrayBuffer` without the headers that a browser needs, so this
tests that build outside a browser.

With SmolLM-135M Q2_K on one machine, the build with one thread made 10.8 tokens
a second, the build with more threads made 36.1 in Node, and the same build in
Chrome made 55.9. All of them gave the same text, because the greedy sampler
does not change from one run to the next.

## Why a worker

Every call into llama.cpp is synchronous, and one token takes milliseconds. A
call from the main thread stops the page. The worker also lets the page show
each token as it comes, because one `Decode` handles one batch.

## More than one thread

The faster build of llama.cpp needs `SharedArrayBuffer`, and a browser gives
that only to a page that comes with these headers:

```
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Embedder-Policy: require-corp
```

`wasm/serve` sets them. If a host does not, `yzma-loader.js` takes the build
with one thread, which is slower but works everywhere. `llamawasm.Threaded()`
says which build the page took.

An isolated page can only get a model from another origin if that origin sends
the CORS headers. Hugging Face does, so the model in `index.html` comes down as
it is, but a model on a host that sends no CORS headers needs a copy on the
origin of the page.

## Limits

- The computation uses the CPU and SIMD. There is no GPU.
- One JavaScript ArrayBuffer holds at most 2 GB, so a larger model must be in
  splits.
- `pkg/llamawasm` has the calls that text generation and embeddings need. It
  does not have multimodal input, LoRA adapters, saved state, or quantization.
