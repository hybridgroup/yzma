# yzma in WebAssembly

This directory holds what a browser needs to run yzma: the JavaScript glue that
picks the right WebAssembly build of llama.cpp, a Web Worker that runs the
program, a page, a static server, and a Node test.

The Go code is in [`pkg/llamawasm`](../pkg/llamawasm) and the example is in
[`examples/wasm/chat`](../examples/wasm/chat).

This page is for people who work on this code. If you want to use yzma in a
browser, read [Run yzma in a browser](https://yzma.ai/docs/tutorials/browser/)
and [Build for a browser](https://yzma.ai/docs/guides/browser/).

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
TinyGo cannot compile the llama.cpp C++ code. So llama.cpp becomes a second
WebAssembly module, built by Emscripten with a small C shim, and the Go code
calls that module through JavaScript. The shim is in the `wasm` directory of the
[llama-cpp-builder](https://github.com/hybridgroup/llama-cpp-builder) repo.

The generation loop stays in Go. It does one `Decode` and one `SamplerSample`
per token, like the `examples/hello` program.

## Files

| File | What it does |
| --- | --- |
| `yzma-loader.js` | Finds the best build the browser can run, which is WebGPU, multithreaded, or single thread. It loads that llama.cpp build and puts it in `globalThis.yzmaReady`. |
| `worker.js` | Runs llama.cpp and the Go program in a Web Worker and sends each piece of text to the page. |
| `index.html` | A page that loads a model and generates text. |
| `vlm.html` | A page that asks a question about an image. |
| `tools.html` | A page where the model calls tools. |
| `decide.html` | A page that asks a System One model typed questions and shows the probabilities. |
| `serve/main.go` | A static server that sets the headers the multithreaded build needs. |
| `node/run.js` | Runs the same build in Node without a browser. CI uses this test. |
| `node/vlm.js` | The same test for an image model. It makes its own pixels, because Node has no canvas. |
| `node/tools.js` | The same test for tool calling. |
| `node/decide.js` | The same test for typed decisions with `exp/decide`. |
| `node/dawn.js` | Runs `node/run.js` with Dawn, the WebGPU implementation in Chrome, so the WebGPU build runs in Node. |
| `node/wgpu.cjs` | Runs `node/run.js` in Deno, which uses wgpu, the WebGPU implementation in Firefox. |
| `node/bench.js` | Measures the tokens per second of a build in Node. `benchmarks/run.sh --backend wasm` calls it. |

## Build and run

```
# get the WebAssembly build of llama.cpp
make download-llama.cpp-wasm

# build the programs with TinyGo
make wasm-example
make wasm-vlm-example
make wasm-tools-example
make wasm-decide-example

# serve them
make serve-wasm
```

Then <http://localhost:8080> is the chat page,
<http://localhost:8080/vlm.html> is the image page,
<http://localhost:8080/tools.html> is the page where the model calls tools, and
<http://localhost:8080/decide.html> is the page for typed decisions.

Add `?mode=cpu` or `?mode=webgpu` to a page URL to pick the backend yourself.
Add `?gpu=high-performance` or `?gpu=low-power` to pick the GPU on a machine
with two.

`make wasm-example-go` builds the same program with the standard Go toolchain.
The binary is larger, but this helps if TinyGo cannot build a dependency.

## Run it without a browser

```
make wasm-example
make test-wasm
```

`node/run.js` loads a small model, generates tokens with the greedy sampler, and
prints them. The greedy sampler always picks the most probable token, so the
output is the same on every run and a test can compare it.

`make test-wasm-mt` does the same with the multithreaded build. Node provides
`SharedArrayBuffer` without the headers a browser needs, so this tests that
build outside a browser.

`make test-wasm-webgpu` tests the fallback from WebGPU. Node has no WebGPU, so
the loader must pick a CPU build and the program must still generate text.

Two targets run the WebGPU build on a real GPU without a browser. Each one runs
the backend self test first, and fails if the GPU is not usable.

```
make test-wasm-dawn NODE=/path/to/node26
make test-wasm-wgpu
```

| Target | WebGPU | Needs |
| --- | --- | --- |
| `test-wasm-dawn` | Dawn, as in Chrome, from the npm package `webgpu` | Node 25 or later for JSPI. The first run installs the package in `build/node`. |
| `test-wasm-wgpu` | wgpu, as in Firefox | Deno 2.9 or later. |

Set `GPU=high-performance` or `GPU=low-power` to pick the GPU on a machine
with two. `DAWN_FLAGS` passes comma separated flags to Dawn, for example
`backend=vulkan,adapter=NVIDIA`. The default enables `shader-f16` on an NVIDIA
card. `DAWN_FLAGS=backend=opengles` reproduces the OpenGL ES path of Chrome on
Linux.

On an Intel RPL-S and an NVIDIA RTX 4070, both GPUs produce the same text as the
CPU with both targets. wgpu is about five times slower than Dawn. The targets
are not in CI, because CI has no GPU.

## Threads

llama.cpp asks for **four** threads unless the caller changes it. On a machine
with more cores this leaves a lot of speed on the table.

| Tokens a second, in Chrome | Four threads | Every thread |
| --- | --- | --- |
| SmolLM-135M Q2_K | 55.9 | 63.3 |
| Gemma 3 1B Q2_K | 8.7 | 18.5 |
| SmolVLM-256M Q8_0, the answer | 56.3 | 96.9 |

So `ContextDefaultParams` and `MtmdContextParamsDefault` pass
`llamawasm.Threads()`, which the JavaScript glue reads from the machine. The
glue sizes the module's thread pool to match. A pool that is too small is worse
than asking for fewer threads, because llama.cpp then waits for threads that
cannot start while the thread that would start them is busy computing.

The WebGPU build benefits too, even though it has no threads. Asking for one
thread instead of four removes the wait at barriers that expect four threads
when only one arrives. SmolLM-135M went from 47.6 tokens per second to 63.3.

Threads have **no effect on images**. The projector took 30.4 seconds on four
threads and 33.3 seconds on sixteen, and 38.3 seconds against 42.7 seconds in a
browser. So many threads are slightly worse. The GPU makes images fast, not the
CPU.

## Speed

The numbers are in
[benchmarks/webassembly.md](../benchmarks/webassembly.md). Run
`./benchmarks/run.sh --backend wasm` for the Node builds. For a browser, which
WebGPU needs, paste
[benchmarks/browser-bench.js](../benchmarks/browser-bench.js) into the page
console.

The GPU is faster on the larger model. On the smaller model the two results
match, because each operation is too small to be worth the transfer to the GPU.
Test both with `?mode=cpu` and `?mode=webgpu`.

Images are different. A 960 by 720 photo through the SmolVLM-256M Q8_0
projector takes 42.7 seconds on the multithreaded CPU build and 1.6 seconds with
WebGPU.

A projector computes many numbers in parallel, which is what a GPU is for, so
the GPU is 25 times faster. On the CPU the user waits for the image, not the
answer. Threads speed up the answer but not the image. So a page with images
needs WebGPU more than a text only page.

Image size has almost no effect, because the model resizes the image to its own
resolution. 224 by 224 and 448 by 448 took almost the same time and both
produced 148 image tokens.

The CPU builds produce the same text every time. The GPU produces the same first
tokens and then diverges, because the shaders do the calculations in a
different order.

## Images

`vlm.html` and `examples/wasm/vlm` answer a question about an image. The
llama.cpp multimodal library, mtmd, is in every build, so there is nothing more
to install.

**The page decodes the image, not llama.cpp.** It draws the file on a canvas and
sends the pixels to the program.

```js
const bitmap = await createImageBitmap(file);
context.drawImage(bitmap, 0, 0, width, height);
const { data } = context.getImageData(0, 0, width, height); // RGBA
worker.postMessage({ kind: "describe", prompt, width, height, rgba: data.buffer }, [data.buffer]);
```

So any format the browser can read works, and the WebAssembly build needs no
image library. The Go side drops the alpha byte and sends RGB to mtmd.

The calls follow the mtmd package with an Mtmd prefix, because the same package
also holds the llama calls.

```go
mctx, err := llamawasm.MtmdInitFromFile("/models/mmproj.gguf", model, 0, onGPU)
bitmap, err := llamawasm.MtmdBitmapInit(width, height, rgb)
chunks, err := llamawasm.MtmdInputChunksInit()
llamawasm.MtmdTokenize(mctx, chunks, prompt, true, true, []llamawasm.MtmdBitmap{bitmap})
nPast, err := llamawasm.MtmdHelperEvalChunks(mctx, ctx, chunks, 0, 0, nBatch, true)
// then the usual loop of SamplerSample and Decode, as for text
```

The prompt must hold one marker per image. `MtmdMarker` returns the model's
marker. `ChatApplyTemplate` puts the marker and the question into the model's
format.

This needs two downloads, the model and its projector (the mmproj file).

Images only. Audio needs the page to decode and resample the samples, and video
needs ffmpeg in a subprocess, which a browser does not have.

## Why a worker

Each call into llama.cpp is synchronous and one token takes milliseconds. A call
from the main thread would freeze the page. The worker also lets the page show
each token right away, because one `Decode` handles one batch.

## WebGPU

There are three llama.cpp builds. `yzma-loader.js` picks the best one the
browser can run.

| Build | What it needs |
| --- | --- |
| `yzma_wasm_webgpu` | WebGPU with f16 shaders, and JSPI. Chrome and Edge 137 or later. Firefox is much slower than its CPU, so auto mode skips it there. See below. The loader also drops this build if the backend self test fails. |
| `yzma_wasm_mt` | `SharedArrayBuffer`, so a page served with the COOP and COEP headers. |
| `yzma_wasm` | Nothing. It works in all browsers. |

A page can set the choice with `globalThis.yzmaMode`, which accepts `auto` (the
default), `webgpu`, or `cpu`. With `webgpu` the loader still falls back to the
CPU if the browser cannot run that build, because a slow page is better than a
broken one. Auto mode picks the CPU in Firefox, because WebGPU there is much
slower than the CPU. Mode `webgpu` still picks the GPU there, so it is easy to
test a fix.

A machine with an integrated and a discrete GPU gives the browser a choice. A
page picks one with `globalThis.yzmaPowerPreference`, which accepts
`high-performance` or `low-power`. Without a value the browser decides. The
loader tests that GPU and makes llama.cpp request the same one, because
llama.cpp sets no preference of its own. The discrete card is usually faster.
SmolLM-135M Q2_K in Chrome 154 ran at 52.9 tokens per second with
`?gpu=high-performance` on an NVIDIA RTX 4070 and 20.3 with `?gpu=low-power` on
an Intel RPL-S, with the same text.

`llamawasm.Backend()` returns the name of the backend doing the computation and
`llamawasm.GPUDevice()` returns the name of the GPU llama.cpp found. Ask
llama.cpp and not the browser, because a page can have WebGPU while llama.cpp
has no device.

### The backend self test

An adapter that llama.cpp accepts can still compute wrong values. The model's
answer is then random tokens from the vocabulary, and nothing in that text
tells it apart from a weak model, so the loader measures the device instead.

`yzma_backend_check` in the shim runs one small matrix multiply of f16 weights
by f32 activations on the device and the same one on the CPU, then compares the
results by normalized mean squared error. A working device gives about 3e-8.
Noise gives about 1. The limit is 1e-2, far from both, so a slow or unusual but
correct driver raises no false alarm. The test needs no model and takes a few
milliseconds.

`yzma-loader.js` runs the test before it hands off the module. A GPU build that
fails is dropped, the loader picks a CPU build, and `globalThis.yzmaGPUReject`
holds the reason. `llamawasm.BackendOK()` gives a Go program the same answer.

`make test-wasm-loader` covers this choice. It gives the loader fake builds, so
it can test a GPU that computes wrong values, which a test with a real module
cannot reach.

### f16 shaders and NVIDIA

The llama.cpp backend needs `shader-f16` and reports no device without it. In
a browser the backend uses the browser's adapter and sets no options. This has
two consequences.

- An Intel integrated GPU provides f16 and the WebGPU build works.
- A discrete NVIDIA card does **not** provide f16 in a browser. Dawn has the
  `vulkan_enable_f16_on_nvidia` option and llama.cpp sets it outside a browser,
  but not inside one. A page cannot set it either, because it is a browser flag.
  So Chrome uses the CPU on such a machine unless the user starts it with this
  command.

  ```
  google-chrome --enable-dawn-features=vulkan_enable_f16_on_nvidia
  ```

  On Linux this switch alone is not enough. See Vulkan in Chrome on Linux
  below.

The fallback makes the page slow, but it still works.

### Vulkan in Chrome on Linux

Chrome on Linux keeps Vulkan off. WebGPU then uses the ANGLE OpenGL ES backend,
in Dawn compatibility mode. `chrome://gpu` shows this as `Vulkan: Disabled`,
and the first adapter under Dawn Info as an `OpenGLES backend` line ending in
`(Compatibility Mode)`.

This path has two outcomes, and neither is good.

- On many cards the adapter has no `shader-f16`, so llama.cpp reports no
  device and the loader picks the CPU. The page is slow but correct.
- On an Intel Xe with Mesa the adapter does have `shader-f16`, llama.cpp uses
  the device, and the device computes wrong values. This is issue #341. No
  adapter feature tells it apart from a working device, so the loader measures
  it. See the backend self test above.

These three switches give Chrome the Vulkan backend. Close every Chrome window
first, and make sure no Chrome process is left running in the background.

```
google-chrome --enable-unsafe-webgpu --enable-features=Vulkan \
  --enable-dawn-features=vulkan_enable_f16_on_nvidia
```

| Switch | Why |
| --- | --- |
| `--enable-unsafe-webgpu` | Turns off the Dawn adapter blocklist. Without it Chrome hides the Vulkan adapters and only offers the OpenGLES adapter. |
| `--enable-features=Vulkan` | Turns on Vulkan in the Chrome GPU process. |
| `--enable-dawn-features=vulkan_enable_f16_on_nvidia` | Enables `shader-f16` on an NVIDIA card. See above. |

All three are needed. `--use-angle=vulkan` is not. `chrome://version` shows
the command line, so check there first if a switch seems to have no effect.

Tested on Chrome 154 on Ubuntu 24.04 with Mesa 25.2.8, an Intel RPL-S and an
NVIDIA RTX 4070. Without the switches Dawn only offers the Intel OpenGLES
adapter, which has no `shader-f16`, so the page uses the CPU. With the three
switches Dawn offers a Vulkan adapter for each card, both with `shader-f16`, and
SmolLM-135M Q2_K generates the correct text with WebGPU at 44 tokens per second.

Run this in a page console to list the adapters.

```js
for (const p of ["low-power", "high-performance"]) {
  const a = await navigator.gpu.requestAdapter({ powerPreference: p });
  console.log(p, a && a.info.vendor, a && a.info.architecture,
    a && a.features.has("shader-f16"));
}
```

A good result names a real GPU on each line, not `swiftshader`, and prints
`true`.

Issue #341 was tested before `--enable-unsafe-webgpu` was known. On Ubuntu
22.04 with Mesa 23.2.1 and only the other two switches, Dawn still offered the
OpenGLES adapter, so the page computed wrong values. Such a machine can try the
three switches. Without them the self test falls back to the CPU.

### Firefox

Firefox runs the WebGPU build, but WebGPU is not on by default yet. Set both of
these in `about:config` and restart the browser.

| Switch | Why |
| --- | --- |
| `dom.webgpu.enabled` | WebGPU on Linux is still behind this switch. |
| `dom.webgpu.workers.enabled` | llama.cpp loads in the worker, so WebGPU in the page is not enough. |

JSPI arrived in Firefox 153, so 153 or later provides `WebAssembly.Suspending`
and `WebAssembly.promising` without a switch. A test in the page console does
not answer the question, because the worker has its own switch. Test the worker.

```js
const w = new Worker(URL.createObjectURL(new Blob([`
  (async () => {
    const a = self.navigator.gpu && await navigator.gpu.requestAdapter();
    postMessage({
      gpu: !!self.navigator.gpu,
      jspi: typeof WebAssembly.Suspending === "function",
      f16: !!a && a.features.has("shader-f16"),
      info: a && { ...a.info },
    });
  })();
`], { type: "text/javascript" })));
w.onmessage = (e) => console.log(e.data);
```

With `?mode=webgpu` the loader reports which part is missing, so check the
worker console first. Firefox 154 on Linux with the two switches reports `gpu`,
`jspi`, and `f16` in the worker, and the page shows `backend: webgpu (WebGPU)`.

Firefox returns an empty `adapter.info`, so the loader has no name for the card
and `globalThis.yzmaAdapter` just holds `webgpu`. This is not a failure.
`llamawasm.Backend()` still returns the real answer.

Firefox uses wgpu and Chrome uses Dawn, so they do not always report the same
adapter or features on the same machine. On a machine with two GPUs one may
provide f16 and the other not. These variables pick the card before Firefox
starts.

```
__NV_PRIME_RENDER_OFFLOAD=1 __GLX_VENDOR_LIBRARY_NAME=nvidia firefox
MESA_VK_DEVICE_SELECT=<vendor>:<device> firefox
```

Firefox has no subgroups, so llama.cpp uses the plain f16 shaders and the GPU
is slower than the same card in Chrome.

#### Firefox is slow

Firefox 154 computed wrong values. A model stopped at the first token, so a
chat page showed a question and no answer. Firefox 156 produces the correct
text, but very slowly.

SmolLM-135M Q2_K on `index.html`, on an Intel RPL-S and an NVIDIA RTX 4070 with
Linux.

| Browser | Backend | Result | Tokens a second |
| --- | --- | --- | --- |
| Firefox 156 | webgpu, the Intel | a correct answer | 0.66 |
| Firefox 156 | webgpu, the NVIDIA | a correct answer | 0.71 |
| Firefox 156 | cpu-threads | a correct answer | 123 |
| Chrome 154 | webgpu, the Intel | a correct answer | 20.3 |
| Chrome 154 | webgpu, the NVIDIA | a correct answer | 52.9 |

The NVIDIA is no faster than the Intel in Firefox, so the time is not spent on
computation. wgpu alone in Deno runs at 3.5 to 5 tokens per second on the same
cards, so most of the time goes to Firefox itself.

That is why `yzma-loader.js` picks the CPU in Firefox in auto mode. Set
`yzmaMode` to `webgpu`, or add `?mode=webgpu` to a page in this directory, to
test the GPU again when a new Firefox comes out.

## Multithreading

The faster llama.cpp build needs `SharedArrayBuffer`. A browser only provides
it to a page served with these headers.

```
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Embedder-Policy: require-corp
```

`wasm/serve` sets them. If a host does not, `yzma-loader.js` picks the single
thread build, which is slower but works in all browsers.
`llamawasm.Threaded()` reports which one was picked.

A host such as GitHub Pages sends no such headers, so a page there gets them
from a service worker such as `coi-serviceworker`. That worker must not route
the model download through `respondWith`. Firefox stops a service worker that
holds a response open for a long time, and the download then fails with
`TypeError: Error in input stream`. Let the browser make the cross origin
request itself, for example with `event.stopImmediatePropagation()` in a
listener registered before the worker's own.

An isolated page can only fetch a model from another origin if that origin
sends CORS headers. Hugging Face does, so the model in `index.html` downloads
without changes. A model on a host without CORS headers must be copied to the
page's origin.

## Typed decisions

`decide.html` and `examples/wasm/decide` run the `exp/decide` package, which
answers a typed question about a state with a probability for each option. It
takes Jev-Style, JevK5 and decider models, the same as on the host. The page fetches the
model and its config, and asks every question about the state with
`DecideMany`.

`exp/decide` shares one state across several questions with a unified KV
cache. That needs a module with ABI version 10 or later. An older module still
works, with one sequence, and decodes each question separately.

`make test-wasm-decide` asks a Jev-Style model three questions in Node. It
fails when `DecideMany` and `Decide` differ in the exact mode.

## Tool calling

`tools.html` and `examples/wasm/tools` let the model call a function. The
`pkg/template` and `pkg/message` packages are pure Go, so they build for
WebAssembly and the browser gets the same tool calling as the host.

`llamawasm.ModelChatTemplate` returns the template stored in the GGUF.
`template.ApplyWithTools` renders it with the whole conversation and the tool
definitions, so a template with a `tools` branch writes out the tools itself.

```go
tmpl := llamawasm.ModelChatTemplate(model, "")
prompt, err := template.ApplyWithTools(tmpl, messages, tools, true)
```

`message.ParseToolCalls` extracts the calls from the answer. The program runs
them, appends a `message.Tool` and a `message.ToolResponse`, and renders again
for the final answer.

A model must be trained for tool calls to make one. Qwen2.5-0.5B-Instruct is
about the smallest that works. `make test-wasm-tools` uses a smaller model and
only checks the round trip, because that model makes no call. Add
`--expect-tool get_weather` to require one.

`llamawasm.ChatApplyTemplate` only takes one message. Use `pkg/template` for a
multi-turn conversation.

### Stop markers

`message.StopMarkers` cuts the text where a model starts writing the next turn.
The shim has no call for the end of turn token, so the WebAssembly build uses
the text of the end of sequence token and probes a short list of common
markers. The host build reads the token directly.

## Limits

- WebGPU needs an adapter with f16 shaders, and Chrome or Edge 137 or later. All
  other browsers use the CPU with SIMD.
- WebGPU in Firefox is much slower than the CPU with llama.cpp, so auto mode
  picks the CPU there.
- Some drivers expose an adapter that llama.cpp accepts but that computes
  wrong values. The loader checks the device against the CPU and picks a CPU
  build if they do not match.
- A browser does not expose subgroup matrix instructions, which llama.cpp only
  uses outside a browser. So the GPU is slower in a page than the same backend
  on a desktop.
- An operation larger than `maxStorageBufferBindingSize` falls back to the CPU.
- A JavaScript ArrayBuffer holds at most 2 GB, so a larger model must be split.
- `pkg/llamawasm` has the calls for text generation, embeddings, and images, and
  it can save context state in memory. It does not support audio, video, LoRA
  adapters, file based state, or quantization.
- The shim has no end of turn token and no grammar sampler. So
  `message.StopMarkers` is an approximation, and a grammar cannot force a tool
  call as it can on a host.
