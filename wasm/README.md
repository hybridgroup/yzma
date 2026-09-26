# yzma in WebAssembly

This directory holds the parts that a browser needs to run yzma. These are the
JavaScript glue that selects the correct WebAssembly build of llama.cpp, a Web
Worker that holds the program, a page, a static server, and a test for Node.

The Go code is in [`pkg/llamawasm`](../pkg/llamawasm) and the example is in
[`examples/wasm/chat`](../examples/wasm/chat).

This page is for the people who work on this code. If you want to use yzma in a
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
| `tools.html` | A page where the model calls tools. |
| `decide.html` | A page that asks a System One model typed questions and shows the probabilities. |
| `serve/main.go` | A static server that sets the headers for a build with more than one thread. |
| `node/run.js` | Runs the same build in Node with no browser. CI uses this test. |
| `node/vlm.js` | The same test for an image model. It makes its own pixels, because Node has no canvas. |
| `node/tools.js` | The same test for tool calling. |
| `node/decide.js` | The same test for typed decisions with `exp/decide`. |
| `node/dawn.js` | Runs `node/run.js` with Dawn, the WebGPU of Chrome, thus the WebGPU build runs in Node. |
| `node/wgpu.cjs` | Runs `node/run.js` in Deno, which has wgpu, the WebGPU of Firefox. |
| `node/bench.js` | Measures the tokens a second of a build in Node. `benchmarks/run.sh --backend wasm` calls it. |

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
<http://localhost:8080/vlm.html> is the page that takes an image, and
<http://localhost:8080/tools.html> is the page where the model calls tools, and
<http://localhost:8080/decide.html> is the page for typed decisions.

Add `?mode=cpu` or `?mode=webgpu` to the URL of a page to select the backend
yourself. Add `?gpu=high-performance` or `?gpu=low-power` to select the GPU on a
machine with two.

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
the loader must select a CPU build and the program must make text.

Two targets run the WebGPU build on a real GPU with no browser. Each one runs
the self test of the backend first, and fails if the GPU is not usable.

```
make test-wasm-dawn NODE=/path/to/node26
make test-wasm-wgpu
```

| Target | WebGPU | Needs |
| --- | --- | --- |
| `test-wasm-dawn` | Dawn, as in Chrome, from the npm package `webgpu` | Node 25 or later for JSPI. The first run installs the package in `build/node`. |
| `test-wasm-wgpu` | wgpu, as in Firefox | Deno 2.9 or later. |

Set `GPU=high-performance` or `GPU=low-power` to select the GPU on a machine
with two. `DAWN_FLAGS` gives flags to Dawn, separated by commas, for example
`backend=vulkan,adapter=NVIDIA`. The default gives `shader-f16` on an NVIDIA
card. `DAWN_FLAGS=backend=opengles` copies the OpenGL ES path of Chrome on
Linux.

On an Intel RPL-S and an NVIDIA RTX 4070 both GPUs give the same text as the
CPU with both targets. wgpu is about five times slower than Dawn. The targets
are not in CI, because CI has no GPU.

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

The numbers are in
[benchmarks/webassembly.md](../benchmarks/webassembly.md). Run
`./benchmarks/run.sh --backend wasm` for the builds in Node, and paste
[benchmarks/browser-bench.js](../benchmarks/browser-bench.js) in the console of
the page for a browser, which WebGPU needs.

The GPU is faster on the larger model. On the smaller model the two results
agree, because each operation is too small to justify the transfer to the GPU.
Test both with `?mode=cpu` and `?mode=webgpu`.

An image gives a different result. A photo of 960 by 720 through the projector
of SmolVLM-256M Q8_0 takes 42.7 seconds on the CPU with more threads and 1.6
seconds with WebGPU.

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
| `yzma_wasm_webgpu` | WebGPU with f16 shaders, and JSPI. Chrome and Edge 137 or later. Firefox gives wrong values, thus auto mode does not use it there. See below. The loader also drops this build if the self test of the backend fails. |
| `yzma_wasm_mt` | `SharedArrayBuffer`, thus a page with the COOP and COEP headers. |
| `yzma_wasm` | Nothing. It operates in all browsers. |

A page can set the choice with `globalThis.yzmaMode`, which accepts `auto` (the
default), `webgpu`, or `cpu`. With `webgpu` the loader still falls back to the
CPU if the browser cannot run that build, because a slow page is better than a
page that does not operate. Auto mode takes the CPU in Firefox, because the
WebGPU of that browser gives wrong values. Mode `webgpu` still selects the GPU
there, thus a test of a repair is easy.

A machine with an integrated and a discrete GPU gives the browser a choice. A
page selects one with `globalThis.yzmaPowerPreference`, which accepts
`high-performance` or `low-power`. With no value the browser selects. The loader
tests that GPU and makes llama.cpp ask for the same one, because llama.cpp sets
no preference of its own. The discrete card is usually faster. SmolLM-135M Q2_K
in Chrome 154 gave 52.9 tokens a second with `?gpu=high-performance` on an NVIDIA
RTX 4070 and 20.3 with `?gpu=low-power` on an Intel RPL-S, with the same text.

`llamawasm.Backend()` gives the name of the part that computes and
`llamawasm.GPUDevice()` gives the name of the GPU that llama.cpp found. Ask
llama.cpp and not the browser, because a page can have WebGPU while llama.cpp
has no device.

### The self test of the backend

An adapter that llama.cpp accepts can still compute wrong values. The answer of
a model is then random tokens of the vocabulary, and nothing in that text tells
it apart from a weak model, so the loader measures the device instead.

`yzma_backend_check` in the shim runs one small matrix multiply of f16 weights
by f32 activations on the device and the same one on the CPU, then compares the
two results with the normalized mean squared error. A device that operates
gives about 3e-8. Noise gives about 1. The limit is 1e-2, which is far from
both, so a slow or an unusual but correct driver raises no false alarm. The
test needs no model and costs a few milliseconds.

`yzma-loader.js` runs the test before it gives the module away. A GPU build
that fails goes away, the loader takes a CPU build, and
`globalThis.yzmaGPUReject` holds the reason. `llamawasm.BackendOK()` gives the
same answer to a Go program.

`make test-wasm-loader` covers this choice. It gives the loader false builds,
so it can test a GPU that computes wrong values, which a test with a real
module cannot reach.

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

  On Linux this switch alone is not sufficient. See Vulkan in Chrome on Linux
  below.

The fallback makes the page slow, but the page operates.

### Vulkan in Chrome on Linux

Chrome on Linux keeps Vulkan off. WebGPU then uses the OpenGL ES backend of
ANGLE, in the compatibility mode of Dawn. `chrome://gpu` shows this as
`Vulkan: Disabled` and the first adapter of Dawn Info as an `OpenGLES backend`
line with `(Compatibility Mode)` at the end.

This path gives two results, and neither is good.

- On many cards the adapter has no `shader-f16`, thus llama.cpp reports no
  device and the loader takes the CPU. The page is slow but correct.
- On an Intel Xe with Mesa the adapter does have `shader-f16`, llama.cpp takes
  the device, and the device computes wrong values. Issue #341 is this case. No
  feature of the adapter tells it apart from a device that operates, thus the
  loader measures it. See the self test of the backend above.

These three switches give Chrome the Vulkan backend. Close every window of
Chrome first, and make sure that no Chrome process stays in the background.

```
google-chrome --enable-unsafe-webgpu --enable-features=Vulkan \
  --enable-dawn-features=vulkan_enable_f16_on_nvidia
```

| Switch | Why |
| --- | --- |
| `--enable-unsafe-webgpu` | Turns off the list of blocked adapters in Dawn. Without it Chrome keeps the Vulkan adapters hidden and gives only the OpenGLES adapter. |
| `--enable-features=Vulkan` | Turns on Vulkan in the GPU process of Chrome. |
| `--enable-dawn-features=vulkan_enable_f16_on_nvidia` | Gives `shader-f16` on an NVIDIA card. See above. |

All three are necessary. `--use-angle=vulkan` is not. `chrome://version` shows
the command line, thus look there first if a switch seems to have no effect.

Tested on Chrome 154 on Ubuntu 24.04 with Mesa 25.2.8, an Intel RPL-S and an
NVIDIA RTX 4070. Without the switches Dawn gives only the OpenGLES adapter of
the Intel, which has no `shader-f16`, thus the page takes the CPU. With the
three switches Dawn gives a Vulkan adapter for each card, both with
`shader-f16`, and SmolLM-135M Q2_K makes the correct text with WebGPU at 44
tokens a second.

This test in the console of a page shows the adapters.

```js
for (const p of ["low-power", "high-performance"]) {
  const a = await navigator.gpu.requestAdapter({ powerPreference: p });
  console.log(p, a && a.info.vendor, a && a.info.architecture,
    a && a.features.has("shader-f16"));
}
```

A good result names a GPU for each line, not `swiftshader`, and gives `true`.

Issue #341 was tested before `--enable-unsafe-webgpu` was known. On Ubuntu
22.04 with Mesa 23.2.1 and only the other two switches, Dawn still gave the
OpenGLES adapter, thus the page computed wrong values. Such a machine can try
the three switches. Without them the self test takes the CPU.

### Firefox

Firefox runs the WebGPU build, but WebGPU is not yet on by default. Set both of
these in `about:config` and restart the browser.

| Switch | Why |
| --- | --- |
| `dom.webgpu.enabled` | WebGPU on Linux is still behind this switch. |
| `dom.webgpu.workers.enabled` | llama.cpp loads in the worker, thus WebGPU in a page is not sufficient. |

JSPI came in Firefox 153, thus 153 or later gives `WebAssembly.Suspending` and
`WebAssembly.promising` with no switch. A test in the console of the page does
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

With `?mode=webgpu` the loader says which of the parts is missing, thus read the
console of the worker first. Firefox 154 on Linux with the two switches gives
`gpu`, `jspi`, and `f16` in the worker, and the page reports
`backend: webgpu (WebGPU)`.

Firefox gives an empty `adapter.info`, thus the loader has no name for the card
and `globalThis.yzmaAdapter` holds the plain word `webgpu`. This is not a
failure. `llamawasm.Backend()` still gives the true answer.

Firefox uses wgpu and Chrome uses Dawn, so the two do not always give the same
adapter or the same features on the same machine. A machine with two GPUs can
give f16 with one and not with the other, and these variables select the card
before Firefox starts.

```
__NV_PRIME_RENDER_OFFLOAD=1 __GLX_VENDOR_LIBRARY_NAME=nvidia firefox
MESA_VK_DEVICE_SELECT=<vendor>:<device> firefox
```

Firefox gives no subgroups, thus llama.cpp takes the plain f16 shaders and the
GPU is slower than the same card in Chrome.

#### Firefox gives wrong values

The WebGPU build operates in Firefox, but the numbers that come back are wrong.
A model stops at the first token, thus a chat page shows a question and no
answer. `Decode` reports no failure, thus the logits reach the sampler with
values that make the first token an end of generation.

The same page, the same model, and the same WebAssembly build give this.

| Browser | Backend | Result |
| --- | --- | --- |
| Firefox 154 | webgpu, which is wgpu | no tokens |
| Firefox 154 with `?mode=cpu` | cpu | a correct answer |
| Chrome 152 | webgpu, which is Dawn | a correct answer |

Thus `yzma-loader.js` takes the CPU in Firefox in auto mode. Set `yzmaMode` to
`webgpu`, or add `?mode=webgpu` to a page of this directory, to test the GPU
again when a new Firefox comes.

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

A host such as GitHub Pages sends no headers, thus a page there gets them from a
service worker such as `coi-serviceworker`. Such a worker must not send the
download of the model through `respondWith`. Firefox stops a service worker that
holds a response open for a long time, and the download then fails with
`TypeError: Error in input stream`. Let the browser make the cross origin
request, for example with `event.stopImmediatePropagation()` in a listener
before the one of the worker.

An isolated page can get a model from another origin only if that origin sends
the CORS headers. Hugging Face sends them, thus the model in `index.html` comes
down without a change. A model on a host with no CORS headers needs a copy on
the origin of the page.

## Typed decisions

`decide.html` and `examples/wasm/decide` run the `exp/decide` package, which
answers a typed question about a state with a probability for each option. It
takes Jev-Style and JevK5 models, the same as on the host. The page fetches the
model and its config, and asks every question about the state with
`DecideMany`.

`exp/decide` shares one state across several questions with a unified KV
cache. That needs a module of ABI version 10 or later. An older module still
runs it, with one sequence, and then decodes each question on its own.

`make test-wasm-decide` asks a Jev-Style model three questions in Node. It
fails when `DecideMany` and `Decide` differ in the exact mode.

## Tool calling

`tools.html` and `examples/wasm/tools` let the model call a function. The
`pkg/template` and `pkg/message` packages are pure Go, thus they build for
WebAssembly and the browser gets the same tool calling as the host.

`llamawasm.ModelChatTemplate` gives the template that the GGUF holds.
`template.ApplyWithTools` renders it with the whole conversation and the tool
definitions, so a template that has a `tools` branch writes the tools itself.

```go
tmpl := llamawasm.ModelChatTemplate(model, "")
prompt, err := template.ApplyWithTools(tmpl, messages, tools, true)
```

`message.ParseToolCalls` reads the calls out of the answer. The program runs
them, appends a `message.Tool` and a `message.ToolResponse`, and renders again
for the final answer.

A model must be trained for tool calls to make one. Qwen2.5-0.5B-Instruct is
about the smallest that works. `make test-wasm-tools` uses a smaller model and
checks the round trip only, because that model makes no call. Add
`--expect-tool get_weather` to require one.

`llamawasm.ChatApplyTemplate` takes one message only. Use `pkg/template` for a
conversation with turns.

### Stop markers

`message.StopMarkers` cuts the text where a model starts to write the next turn.
The shim has no call for the end of turn token, thus the WebAssembly build takes
the text of the end of sequence token and probes a short list of the usual
markers. The host build reads the token itself.

## Limits

- WebGPU needs an adapter with f16 shaders, and Chrome or Edge 137 or later. All
  other browsers use the CPU with SIMD.
- The WebGPU of Firefox gives wrong values to llama.cpp, thus auto mode takes
  the CPU there.
- Some drivers give an adapter that llama.cpp accepts and that then computes
  wrong values. The loader measures the device against the CPU and takes a CPU
  build if the two do not agree.
- A browser does not give the matrix instructions of a subgroup, which llama.cpp
  uses only outside a browser. Thus the GPU is slower in a page than the same
  backend on a desktop.
- An operation larger than `maxStorageBufferBindingSize` returns to the CPU.
- One JavaScript ArrayBuffer holds a maximum of 2 GB, thus a larger model must be
  in splits.
- `pkg/llamawasm` has the calls for text generation, embeddings, and images, and
  it saves the state of a context in memory. It does not have audio, video, LoRA
  adapters, state in a file, or quantization.
- The shim gives no end of turn token and no grammar sampler. Thus
  `message.StopMarkers` approximates, and a tool call cannot be forced by a
  grammar as it can on a host.
