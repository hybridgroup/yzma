# WebAssembly benchmarks

Benchmarks of yzma in WebAssembly. Each table shows the median of five runs. The
output of each run is below the tables.

There are two kinds of run.

- **Node**, which `benchmarks/run.sh --backend wasm` runs with
  [wasm/node/bench.js](../wasm/node/bench.js). It covers the single threaded
  build and the multithreaded build.
- **A browser**, which no script here can drive. Serve the example with
  `make serve-wasm`, paste [browser-bench.js](browser-bench.js) into the page
  console, and pass the result to `yzma-bench update`. WebGPU needs this,
  because `bench.js` has no WebGPU. On a machine with two GPUs, set `gpu` in
  the script to pick one.

Both kinds read the llama.cpp build from `yzma-install.json` in the build
directory, so each result records which build produced it.

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf)
and the tokens a second come from the generation loop in
[examples/wasm/chat](../examples/wasm/chat), which is the loop the page uses.
These numbers are not comparable with the native tables, because the prompt
and the token count are different.

## Summary

Tokens a second on each machine.

| Machine | Node, CPU | Node, more threads | Chrome, more threads | Chrome, WebGPU Intel | Chrome, WebGPU NVIDIA |
| --- | --- | --- | --- | --- | --- |
| Intel Core i9-13900HX | 13.8 | 105.7 | 92.8 | 19.9 | 69.8 |

- The multithreaded build is 7.7 times faster than the single threaded build.
- In Chrome, the multithreaded build reaches 88 percent of its speed in Node.
- On this small model the multithreaded CPU build is faster than WebGPU. Each
  operation is too small to pay for the trip to the GPU. The GPU pays off more
  with a larger model and with images.
- The RTX 4070 is 3.5 times faster than the Intel RPL-S. Chrome on Linux needs
  three switches for WebGPU with Vulkan, see
  [wasm/README.md](../wasm/README.md).
- The Node numbers use b11146 and the Chrome numbers use b11202.

## In Node

<!-- yzma:bench table node -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | wasm | Intel Core i9-13900HX | - | 13.8 | b11146 | 2026-09-24 |
| CPU, more threads | wasm | Intel Core i9-13900HX | - | 105.7 | b11146 | 2026-09-24 |
<!-- yzma:bench table end node -->

<!-- yzma:bench start node/cpu/wasm/i9-13900hx -->
### CPU, wasm, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"node","backend":"cpu","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":13.8,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 13.8 tokens a second.

<details><summary>The output of go test</summary>

```
$ node wasm/node/bench.js --dir build/wasm --model SmolLM-135M.Q2_K.gguf --tokens 64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
backend: cpu, 1 thread(s)
llama.cpp: b11146
BenchmarkInference-1	1	4624277457 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4637681159 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4627621114 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4620938628 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4630969609 ns/op	13.8 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	23.162s
```

</details>
<!-- yzma:bench end node/cpu/wasm/i9-13900hx -->

<!-- yzma:bench start node/cpu-threads/wasm/i9-13900hx -->
### CPU, more threads, wasm, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"node","backend":"cpu-threads","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":105.7,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 105.7 tokens a second.

<details><summary>The output of go test</summary>

```
$ node wasm/node/bench.js --dir build/wasm --model SmolLM-135M.Q2_K.gguf --tokens 64 --mt
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
backend: cpu-threads, 16 thread(s)
llama.cpp: b11146
BenchmarkInference-16	1	640192058 ns/op	100.0 tokens/s
BenchmarkInference-16	1	565720852 ns/op	113.1 tokens/s
BenchmarkInference-16	1	607902736 ns/op	105.3 tokens/s
BenchmarkInference-16	1	605258180 ns/op	105.7 tokens/s
BenchmarkInference-16	1	594077787 ns/op	107.7 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	3.033s
```

</details>
<!-- yzma:bench end node/cpu-threads/wasm/i9-13900hx -->

## In a browser

<!-- yzma:bench table browser -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU, more threads | wasm | Intel Core i9-13900HX, Chrome | - | 92.8 | b11202 | 2026-09-27 |
| WebGPU | wasm | Intel Core i9-13900HX, Chrome | Intel-RPL-S | 19.9 | b11202 | 2026-09-27 |
| WebGPU | wasm | Intel Core i9-13900HX, Chrome | RTX-4070 | 69.8 | b11202 | 2026-09-27 |
<!-- yzma:bench table end browser -->

<!-- yzma:bench start browser/cpu-threads/wasm/i9-13900hx -->
### CPU, more threads, wasm, Intel Core i9-13900HX, Chrome
<!-- yzma:bench meta {"suite":"browser","backend":"cpu-threads","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX, Chrome","cpu":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36","tokens_per_second":92.8,"llamacpp":"b11202","date":"2026-09-27"} -->

Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36. 92.8 tokens a second.

<details><summary>The output of go test</summary>

```
$ browser-bench.js mode=cpu tokens=64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36
backend: cpu-threads, 16 threads
llama.cpp: b11202
BenchmarkInference-32	1	689506572 ns/op	92.8 tokens/s
BenchmarkInference-32	1	674252002 ns/op	94.9 tokens/s
BenchmarkInference-32	1	700065631 ns/op	91.4 tokens/s
BenchmarkInference-32	1	683833743 ns/op	93.6 tokens/s
BenchmarkInference-32	1	703528636 ns/op	91.0 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	3.471s
```

</details>
<!-- yzma:bench end browser/cpu-threads/wasm/i9-13900hx -->

<!-- yzma:bench start browser/webgpu/wasm/i9-13900hx/intel-rpl-s -->
### WebGPU, wasm, Intel Core i9-13900HX, Chrome, Intel-RPL-S
<!-- yzma:bench meta {"suite":"browser","backend":"webgpu","arch":"wasm","machine":"i9-13900hx","device":"Intel-RPL-S","label":"Intel Core i9-13900HX, Chrome","cpu":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36","tokens_per_second":19.9,"llamacpp":"b11202","date":"2026-09-27"} -->

Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36. 19.9 tokens a second.

Chrome 154 with --enable-unsafe-webgpu --enable-features=Vulkan --enable-dawn-features=vulkan_enable_f16_on_nvidia, and ?gpu=low-power.

<details><summary>The output of go test</summary>

```
$ browser-bench.js mode=webgpu gpu=low-power tokens=64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36
backend: webgpu (WebGPU)
llama.cpp: b11202
BenchmarkInference-32	1	3247082699 ns/op	19.7 tokens/s
BenchmarkInference-32	1	3261977574 ns/op	19.6 tokens/s
BenchmarkInference-32	1	3219315895 ns/op	19.9 tokens/s
BenchmarkInference-32	1	2996254682 ns/op	21.4 tokens/s
BenchmarkInference-32	1	2752688172 ns/op	23.3 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	15.483s
```

</details>
<!-- yzma:bench end browser/webgpu/wasm/i9-13900hx/intel-rpl-s -->

<!-- yzma:bench start browser/webgpu/wasm/i9-13900hx/rtx-4070 -->
### WebGPU, wasm, Intel Core i9-13900HX, Chrome, RTX-4070
<!-- yzma:bench meta {"suite":"browser","backend":"webgpu","arch":"wasm","machine":"i9-13900hx","device":"RTX-4070","label":"Intel Core i9-13900HX, Chrome","cpu":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36","tokens_per_second":69.8,"llamacpp":"b11202","date":"2026-09-27"} -->

Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36. 69.8 tokens a second.

Chrome 154 with --enable-unsafe-webgpu --enable-features=Vulkan --enable-dawn-features=vulkan_enable_f16_on_nvidia, and ?gpu=high-performance.

<details><summary>The output of go test</summary>

```
$ browser-bench.js mode=webgpu gpu=high-performance tokens=64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36
backend: webgpu (WebGPU)
llama.cpp: b11202
BenchmarkInference-32	1	856875084 ns/op	74.7 tokens/s
BenchmarkInference-32	1	925791986 ns/op	69.1 tokens/s
BenchmarkInference-32	1	1102118133 ns/op	58.1 tokens/s
BenchmarkInference-32	1	884955752 ns/op	72.3 tokens/s
BenchmarkInference-32	1	916511528 ns/op	69.8 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	4.691s
```

</details>
<!-- yzma:bench end browser/webgpu/wasm/i9-13900hx/rtx-4070 -->
