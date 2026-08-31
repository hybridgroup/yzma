# yzma in WebAssembly

This directory holds the parts that a browser needs to run yzma. These are the
JavaScript glue that selects the correct WebAssembly build of llama.cpp, a Web
Worker that holds the program, a page, a static server, and a test for Node.

The Go code is in [`pkg/llamawasm`](../pkg/llamawasm) and the example is in
[`examples/wasm/chat`](../examples/wasm/chat).

## How it works

There are two WebAssembly modules.

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
TinyGo cannot compile the C++ of llama.cpp. Thus llama.cpp becomes a second
WebAssembly module, made by Emscripten with a small C shim, and the Go code
calls that module through JavaScript. The shim is in the `wasm` directory of the
[llama-cpp-builder](https://github.com/hybridgroup/llama-cpp-builder) repo.

The generation loop stays in Go. It does one `Decode` and one `SamplerSample`
for each token, as in the `examples/hello` program.

## Files

| File | What it does |
| --- | --- |
| `yzma-loader.js` | Finds the best build that the browser can run, which is WebGPU, more than one thread, or one thread. It loads that build of llama.cpp and puts it in `globalThis.yzmaReady`. |
| `worker.js` | Runs llama.cpp and the Go program in a Web Worker and sends each piece of text to the page. |
| `index.html` | A page that loads a model and makes text. |
| `vlm.html` | A page that asks a question about an image. |
| `serve/main.go` | A static server that sets the headers for a build with more than one thread. |
| `node/run.js` | Runs the same build in Node with no browser. CI uses this test. |
| `node/vlm.js` | The same test for an image model. It makes its own pixels, because Node has no canvas. |

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

Add `?mode=cpu` or `?mode=webgpu` to the URL of a page to select the backend
yourself.

`make wasm-example-go` builds the same program with the standard Go toolchain.
The binary is larger, which is an aid if TinyGo cannot build a dependency.

## Run it without a browser

```
make wasm-example
make test-wasm
```

`node/run.js` loads a small model, makes tokens with the greedy sampler, and
prints them. The greedy sampler always takes the most probable token, thus the
output does not change between runs and a test can compare it.

`make test-wasm-mt` does the same with the build that uses more than one thread.
Node gives `SharedArrayBuffer` without the headers that a browser needs, thus
this tests that build outside a browser.

`make test-wasm-webgpu` tests the fallback from WebGPU. Node has no WebGPU, thus
the loader must select a CPU build and the program must make text. Only a
browser can run the WebGPU build.

## Threads

llama.cpp asks for **four** threads unless a caller changes it. On a machine
with more cores this loses much speed.

| Tokens a second, in Chrome | Four threads | Every thread |
| --- | --- | --- |
| SmolLM-135M Q2_K | 55.9 | 63.3 |
| Gemma 3 1B Q2_K | 8.7 | 18.5 |
| SmolVLM-256M Q8_0, the answer | 56.3 | 96.9 |

Thus `ContextDefaultParams` and `MtmdContextParamsDefault` send
`llamawasm.Threads()`, which the JavaScript glue reads from the machine. The
glue makes the thread pool of the module the same size. A pool that is too small
is worse than a small number of threads, because llama.cpp then waits for threads
that cannot start. The thread that starts them is busy with computation.

The WebGPU build also gets an advantage, although it has no threads. A request
for one thread in place of four removes the wait at the barriers that four
threads make and only one thread reaches. SmolLM-135M changed from 47.6 tokens a
second to 63.3.

The threads have **no effect on the image**. The projector used 30.4 seconds on
four threads and 33.3 seconds on sixteen, and 38.3 seconds against 42.7 seconds
in a browser. Thus a large number of threads is a little worse. The GPU makes an
image fast, not the CPU.

## Speed

Measured in Chrome on one machine, an RTX 4070 with an Intel integrated GPU,
with the greedy sampler.

| Model | Backend | Tokens a second |
| --- | --- | --- |
| SmolLM-135M Q2_K | one thread | 10.8 |
| SmolLM-135M Q2_K | more threads | 63.3 |
| SmolLM-135M Q2_K | WebGPU | 63.3 |
| Gemma 3 1B Q2_K | more threads | 18.5 |
| Gemma 3 1B Q2_K | WebGPU | 38.7 |

The GPU is faster on the larger model. On the smaller model the two results
agree, because each operation is too small to justify the transfer to the GPU.
Test both with `?mode=cpu` and `?mode=webgpu`.

An image gives a different result. This is the same photo of 960 by 720 through
the projector of SmolVLM-256M Q8_0, and then 32 tokens of answer.

| Backend | Time for the image | Tokens a second |
| --- | --- | --- |
| more threads, in Chrome | 42.7 s | 96.9 |
| WebGPU, in Chrome | 1.6 s | 64.4 |
| one thread, in Node | 80 s | 17.4 |

A projector computes many numbers at the same time, which is the function of a
GPU. Thus the GPU is 25 times faster. On the CPU the reader waits for the image
and not for the answer. The threads make the answer faster but do not change the
time of the image. Thus a page with images needs WebGPU more than a page with
only text.

The size of the image has almost no effect. A model has its own resolution and
changes the size of the image. Thus 224 by 224 and 448 by 448 used almost the
same time and both gave 148 tokens of image.

The CPU builds give the same text each time. The GPU gives the same text for the
first tokens and then gives different text, because the shaders do the
calculations in a different order.

## Images

`vlm.html` and `examples/wasm/vlm` answer a question about an image. The
multimodal library of llama.cpp, mtmd, is in each build, thus you install
nothing more.

**The page decodes the image, not llama.cpp.** It draws the file on a canvas and
sends the pixels to the program.

```js
const bitmap = await createImageBitmap(file);
context.drawImage(bitmap, 0, 0, width, height);
const { data } = context.getImageData(0, 0, width, height); // RGBA
worker.postMessage({ kind: "describe", prompt, width, height, rgba: data.buffer }, [data.buffer]);
```

Thus each format that the browser reads is usable and the WebAssembly build
needs no image library. The Go side removes the alpha byte and sends the RGB to
mtmd.

The calls follow the mtmd package with an Mtmd prefix, because one package also
holds the llama calls.

```go
mctx, err := llamawasm.MtmdInitFromFile("/models/mmproj.gguf", model, 0, onGPU)
bitmap, err := llamawasm.MtmdBitmapInit(width, height, rgb)
chunks, err := llamawasm.MtmdInputChunksInit()
llamawasm.MtmdTokenize(mctx, chunks, prompt, true, true, []llamawasm.MtmdBitmap{bitmap})
nPast, err := llamawasm.MtmdHelperEvalChunks(mctx, ctx, chunks, 0, 0, nBatch, true)
// then the usual loop of SamplerSample and Decode, as for text
```

The prompt must hold one marker for each image. `MtmdMarker` gives the marker of
the model. `ChatApplyTemplate` puts the marker and the question into the correct
format for the model.

Two files come down for this, the model and its projector, the mmproj file.

Images only. Audio needs the page to decode and resample the samples, and video
needs ffmpeg in a subprocess, which a browser does not have.

## Why a worker

Each call into llama.cpp is synchronous and one token takes milliseconds. A call
from the main thread stops the page. The worker also lets the page show each
token immediately, because one `Decode` does one batch.

## WebGPU

There are three builds of llama.cpp. `yzma-loader.js` selects the best one that
the browser can run.

| Build | What it needs |
| --- | --- |
| `yzma_wasm_webgpu` | WebGPU with f16 shaders, and JSPI. Chrome and Edge 137 or later. |
| `yzma_wasm_mt` | `SharedArrayBuffer`, thus a page with the COOP and COEP headers. |
| `yzma_wasm` | Nothing. It operates in all browsers. |

A page can set the choice with `globalThis.yzmaMode`, which accepts `auto` (the
default), `webgpu`, or `cpu`. With `webgpu` the loader still falls back to the
CPU if the browser cannot run that build, because a slow page is better than a
page that does not operate.

`llamawasm.Backend()` gives the name of the part that computes and
`llamawasm.GPUDevice()` gives the name of the GPU that llama.cpp found. Ask
llama.cpp and not the browser, because a page can have WebGPU while llama.cpp
has no device.

### f16 shaders and NVIDIA

The backend of llama.cpp needs `shader-f16` and reports no device without it. In
a browser the backend uses the adapter of the browser and sets no options. This
gives two results.

- An Intel integrated GPU gives f16 and the WebGPU build operates.
- A discrete NVIDIA card does **not** give f16 in a browser. Dawn has the
  `vulkan_enable_f16_on_nvidia` option and llama.cpp sets it outside a browser,
  but not in one. A page also cannot set it, because it is a flag of the
  browser. Thus Chrome uses the CPU on such a machine unless the user starts it
  with this command.

  ```
  google-chrome --enable-dawn-features=vulkan_enable_f16_on_nvidia
  ```

The fallback makes the page slow, but the page operates.

## More than one thread

The faster build of llama.cpp needs `SharedArrayBuffer`. A browser gives that
only to a page with these headers.

```
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Embedder-Policy: require-corp
```

`wasm/serve` sets them. If a host does not set them, `yzma-loader.js` selects
the build with one thread, which is slower but operates in all browsers.
`llamawasm.Threaded()` gives the selection.

An isolated page can get a model from another origin only if that origin sends
the CORS headers. Hugging Face sends them, thus the model in `index.html` comes
down without a change. A model on a host with no CORS headers needs a copy on
the origin of the page.

## Limits

- WebGPU needs Chrome or Edge 137 or later, and an adapter with f16 shaders. All
  other browsers use the CPU with SIMD.
- A browser does not give the matrix instructions of a subgroup, which llama.cpp
  uses only outside a browser. Thus the GPU is slower in a page than the same
  backend on a desktop.
- An operation larger than `maxStorageBufferBindingSize` returns to the CPU.
- One JavaScript ArrayBuffer holds a maximum of 2 GB, thus a larger model must be
  in splits.
- `pkg/llamawasm` has the calls for text generation, embeddings, and images. It
  does not have audio, video, LoRA adapters, saved state, or quantization.
