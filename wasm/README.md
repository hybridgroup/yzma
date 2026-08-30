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
| `yzma-loader.js` | Finds what the browser can run — WebGPU, more than one thread, or one thread — then loads that build of llama.cpp and puts it in `globalThis.yzmaReady`. |
| `worker.js` | Runs llama.cpp and the Go program in a Web Worker, and sends each piece of text to the page. |
| `index.html` | A page that loads a model and makes text. |
| `vlm.html` | A page that asks a question about an image. |
| `serve/main.go` | A static server that sets the headers that a build with more than one thread needs. |
| `node/run.js` | Runs the same build in Node, with no browser. This is the test that CI uses. |
| `node/vlm.js` | The same for a model with eyes. It makes its own pixels, because Node has no canvas. |

## Build and run

```
# get the WebAssembly build of llama.cpp
make download-llama.cpp-wasm

# build the programs with TinyGo
make wasm-example
make wasm-vlm-example

# serve them
make serve-wasm
```

Then <http://localhost:8080> is the chat page and
<http://localhost:8080/vlm.html> is the page that takes an image.

Add `?mode=cpu` or `?mode=webgpu` to the URL of either page to choose the backend
instead of letting the loader choose.

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

`make test-wasm-webgpu` tests the other half of the WebGPU story, which is what
happens where there is none: Node has no WebGPU, so the loader must take a build
on the CPU and the program must still make text. Only a browser can run the
WebGPU build itself.

## Speed

Measured in Chrome on one machine, an RTX 4070 with an Intel integrated GPU,
with the greedy sampler:

| Model | Backend | Tokens a second |
| --- | --- | --- |
| SmolLM-135M Q2_K | one thread | 10.8 |
| SmolLM-135M Q2_K | more threads | 55.9 |
| SmolLM-135M Q2_K | WebGPU | 47.6 |
| Gemma 3 1B Q2_K | more threads | 8.7 |
| Gemma 3 1B Q2_K | WebGPU | 32.3 |

The GPU wins on the larger model and loses on the smaller one, where the work of
each operation is too small to pay for the trip to the GPU. Try both: the page
takes `?mode=cpu` and `?mode=webgpu`.

An image is a different story. This is the same photo of 960 by 720 through the
projector of SmolVLM-256M Q8_0, and then 32 tokens of answer:

| Backend | Time for the image | Tokens a second |
| --- | --- | --- |
| more threads, in Chrome | 38.3 s | 56.3 |
| WebGPU, in Chrome | 1.5 s | 61.7 |
| one thread, in Node | 80 s | 17.4 |

Putting an image through a projector is a wall of numbers all at once, which is
what a GPU is for, so the GPU is twenty five times faster at it. On the CPU it is
the image, and not the answer, that a reader waits for.

The size of the image hardly matters: a model has a resolution of its own and
resizes what it gets, so 224 by 224 and 448 by 448 both took about the same time
and both came to 148 tokens of image.

The builds on the CPU give the same text every time. The GPU gives the same text
for the first few tokens and then goes its own way, because the shaders do the
arithmetic in a different order than the CPU does.

## A model with eyes

`vlm.html` and `examples/wasm/vlm` answer a question about an image. The
multimodal library of llama.cpp, mtmd, is in every build, so nothing more needs
installing.

**The page decodes the image, not llama.cpp.** It draws the file on a canvas and
sends the pixels to the program:

```js
const bitmap = await createImageBitmap(file);
context.drawImage(bitmap, 0, 0, width, height);
const { data } = context.getImageData(0, 0, width, height); // RGBA
worker.postMessage({ kind: "describe", prompt, width, height, rgba: data.buffer }, [data.buffer]);
```

So any format the browser reads works, and no image library goes into the
WebAssembly build. The Go side drops the alpha byte and hands the RGB to mtmd.

The calls follow the mtmd package, with an Mtmd prefix because one package holds
the llama calls as well:

```go
mctx, err := llamawasm.MtmdInitFromFile("/models/mmproj.gguf", model, 0, onGPU)
bitmap, err := llamawasm.MtmdBitmapInit(width, height, rgb)
chunks, err := llamawasm.MtmdInputChunksInit()
llamawasm.MtmdTokenize(mctx, chunks, prompt, true, true, []llamawasm.MtmdBitmap{bitmap})
nPast, err := llamawasm.MtmdHelperEvalChunks(mctx, ctx, chunks, 0, 0, nBatch, true)
// then the same loop of SamplerSample and Decode as for text alone
```

The prompt must hold one marker for each image, and `MtmdMarker` gives the
marker of the model. `ChatApplyTemplate` puts the marker and the question into
the format the model expects.

Two models come down for this: the model itself and its projector, the mmproj
file.

Images only. Audio and video stay out: audio needs the page to decode and
resample the samples itself, and video needs ffmpeg in a subprocess, which a
browser has not got.

## Why a worker

Every call into llama.cpp is synchronous, and one token takes milliseconds. A
call from the main thread stops the page. The worker also lets the page show
each token as it comes, because one `Decode` handles one batch.

## WebGPU

There are three builds of llama.cpp, and `yzma-loader.js` takes the best one
that the browser can run:

| Build | What it needs |
| --- | --- |
| `yzma_wasm_webgpu` | WebGPU with f16 shaders, and JSPI. Chrome and Edge 137 and later. |
| `yzma_wasm_mt` | `SharedArrayBuffer`, so a page with the COOP and COEP headers. |
| `yzma_wasm` | Nothing. It works everywhere. |

A page can force the choice with `globalThis.yzmaMode`, which takes `auto` (the
default), `webgpu`, or `cpu`. `webgpu` still falls back to the CPU if the
browser cannot run that build, because a page that does not work at all is
worse than a page that is slower.

`llamawasm.Backend()` says what does the computation, and
`llamawasm.GPUDevice()` names the GPU that llama.cpp found. Ask llama.cpp and
not the browser: a page can have WebGPU while llama.cpp still has no device.

### f16 shaders, and NVIDIA

The backend of llama.cpp needs `shader-f16`, and it reports no device at all
without it. Two things follow from how the backend asks for an adapter in a
browser, where it takes the adapter of the browser with no options of its own:

- An Intel integrated GPU gives f16, and the WebGPU build runs.
- A discrete NVIDIA card does **not** give f16 in a browser. Dawn has a toggle
  for it, `vulkan_enable_f16_on_nvidia`, and llama.cpp turns it on outside a
  browser but cannot in one. A page cannot turn it on either: it is a flag of
  the browser. So Chrome takes the CPU on such a machine unless somebody starts
  it with:

  ```
  google-chrome --enable-dawn-features=vulkan_enable_f16_on_nvidia
  ```

The fallback makes this a slower page and not a broken one.

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

- WebGPU needs Chrome or Edge 137 and later, and an adapter with f16 shaders.
  Everything else runs on the CPU with SIMD.
- A browser does not get the matrix instructions of a subgroup, which llama.cpp
  uses only outside a browser, so the GPU is slower in a page than the same
  backend is on a desktop.
- An operation larger than `maxStorageBufferBindingSize` goes back to the CPU.
- One JavaScript ArrayBuffer holds at most 2 GB, so a larger model must be in
  splits.
- `pkg/llamawasm` has the calls that text generation, embeddings, and images
  need. It does not have audio, video, LoRA adapters, saved state, or
  quantization.
