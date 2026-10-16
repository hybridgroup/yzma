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
| Intel Core i9-13900HX | 13.9 | 111.7 | 93.3 | 19.6 | 66.2 |

- The multithreaded build is 8.0 times faster than the single threaded build.
- In Chrome, the multithreaded build reaches 84 percent of its speed in Node.
- On this small model the multithreaded CPU build is faster than WebGPU. Each
  operation is too small to pay for the trip to the GPU. The GPU pays off more
  with a larger model and with images.
- The RTX 4070 is 3.4 times faster than the Intel RPL-S. Chrome on Linux needs
  three switches for WebGPU with Vulkan, see
  [wasm/README.md](../wasm/README.md).
- Every number uses b11669.

## In Node

<!-- yzma:bench table node -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | wasm | Intel Core i9-13900HX | - | 13.9 | b11669 | 2026-10-16 |
| CPU, more threads | wasm | Intel Core i9-13900HX | - | 111.7 | b11669 | 2026-10-16 |
<!-- yzma:bench table end node -->

<!-- yzma:bench start node/cpu/wasm/i9-13900hx -->
### CPU, wasm, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"node","backend":"cpu","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":13.9,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 13.9 tokens a second.

<details><summary>The output of go test</summary>

```
$ node wasm/node/bench.js --dir build/wasm --model SmolLM-135M.Q2_K.gguf --tokens 64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
backend: cpu, 1 thread(s)
llama.cpp: b11669
BenchmarkInference-1	1	4620938628 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4614275415 ns/op	13.9 tokens/s
BenchmarkInference-1	1	4617604618 ns/op	13.9 tokens/s
BenchmarkInference-1	1	4610951009 ns/op	13.9 tokens/s
BenchmarkInference-1	1	4610951009 ns/op	13.9 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	23.091s
```

</details>
<!-- yzma:bench end node/cpu/wasm/i9-13900hx -->

<!-- yzma:bench start node/cpu-threads/wasm/i9-13900hx -->
### CPU, more threads, wasm, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"node","backend":"cpu-threads","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":111.7,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 111.7 tokens a second.

<details><summary>The output of go test</summary>

```
$ node wasm/node/bench.js --dir build/wasm --model SmolLM-135M.Q2_K.gguf --tokens 64 --mt
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
backend: cpu-threads, 16 thread(s)
llama.cpp: b11669
BenchmarkInference-16	1	638722555 ns/op	100.2 tokens/s
BenchmarkInference-16	1	574351611 ns/op	111.4 tokens/s
BenchmarkInference-16	1	573168547 ns/op	111.7 tokens/s
BenchmarkInference-16	1	560175055 ns/op	114.3 tokens/s
BenchmarkInference-16	1	542097239 ns/op	118.1 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	2.908s
```

</details>
<!-- yzma:bench end node/cpu-threads/wasm/i9-13900hx -->

## In a browser

<!-- yzma:bench table browser -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU, more threads | wasm | Intel Core i9-13900HX, Chrome | - | 93.3 | b11669 | 2026-10-16 |
| WebGPU | wasm | Intel Core i9-13900HX, Chrome | Intel-RPL-S | 19.6 | b11669 | 2026-10-16 |
| WebGPU | wasm | Intel Core i9-13900HX, Chrome | RTX-4070 | 66.2 | b11669 | 2026-10-16 |
<!-- yzma:bench table end browser -->

<!-- yzma:bench start browser/cpu-threads/wasm/i9-13900hx -->
### CPU, more threads, wasm, Intel Core i9-13900HX, Chrome
<!-- yzma:bench meta {"suite":"browser","backend":"cpu-threads","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX, Chrome","cpu":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36","tokens_per_second":93.3,"llamacpp":"b11669","date":"2026-10-16"} -->

Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36. 93.3 tokens a second.

<details><summary>The output of go test</summary>

```
$ browser-bench.js mode=cpu tokens=64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36
backend: cpu-threads, 16 threads
llama.cpp: b11669
BenchmarkInference-32    1    685591859 ns/op    93.3 tokens/s
BenchmarkInference-32    1    678109769 ns/op    94.4 tokens/s
BenchmarkInference-32    1    685959271 ns/op    93.3 tokens/s
BenchmarkInference-32    1    684053014 ns/op    93.6 tokens/s
BenchmarkInference-32    1    690026954 ns/op    92.8 tokens/s
PASS
ok    github.com/hybridgroup/yzma/examples/wasm/chat    3.443s
```

</details>
<!-- yzma:bench end browser/cpu-threads/wasm/i9-13900hx -->

<!-- yzma:bench start browser/webgpu/wasm/i9-13900hx/intel-rpl-s -->
### WebGPU, wasm, Intel Core i9-13900HX, Chrome, Intel-RPL-S
<!-- yzma:bench meta {"suite":"browser","backend":"webgpu","arch":"wasm","machine":"i9-13900hx","device":"Intel-RPL-S","label":"Intel Core i9-13900HX, Chrome","cpu":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36","tokens_per_second":19.6,"llamacpp":"b11669","date":"2026-10-16"} -->

Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36. 19.6 tokens a second.

<details><summary>The output of go test</summary>

```
$ browser-bench.js mode=webgpu gpu=low-power tokens=64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36
backend: webgpu (WebGPU)
llama.cpp: b11669
BenchmarkInference-32    1    3158933860 ns/op    20.3 tokens/s
BenchmarkInference-32    1    3266972945 ns/op    19.6 tokens/s
BenchmarkInference-32    1    3340292276 ns/op    19.2 tokens/s
BenchmarkInference-32    1    3324675325 ns/op    19.3 tokens/s
BenchmarkInference-32    1    3258655804 ns/op    19.6 tokens/s
PASS
ok    github.com/hybridgroup/yzma/examples/wasm/chat    16.358s
```

</details>
<!-- yzma:bench end browser/webgpu/wasm/i9-13900hx/intel-rpl-s -->

<!-- yzma:bench start browser/webgpu/wasm/i9-13900hx/rtx-4070 -->
### WebGPU, wasm, Intel Core i9-13900HX, Chrome, RTX-4070
<!-- yzma:bench meta {"suite":"browser","backend":"webgpu","arch":"wasm","machine":"i9-13900hx","device":"RTX-4070","label":"Intel Core i9-13900HX, Chrome","cpu":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36","tokens_per_second":66.2,"llamacpp":"b11669","date":"2026-10-16"} -->

Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36. 66.2 tokens a second.

<details><summary>The output of go test</summary>

```
$ browser-bench.js mode=webgpu gpu=high-performance tokens=64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/155.0.0.0 Safari/537.36
backend: webgpu (WebGPU)
llama.cpp: b11669
BenchmarkInference-32    1    873005047 ns/op    73.3 tokens/s
BenchmarkInference-32    1    987959247 ns/op    64.8 tokens/s
BenchmarkInference-32    1    967498110 ns/op    66.2 tokens/s
BenchmarkInference-32    1    846784864 ns/op    75.6 tokens/s
BenchmarkInference-32    1    985676883 ns/op    64.9 tokens/s
PASS
ok    github.com/hybridgroup/yzma/examples/wasm/chat    4.665s
```

</details>
<!-- yzma:bench end browser/webgpu/wasm/i9-13900hx/rtx-4070 -->
