# WebAssembly benchmarks

Benchmarks of yzma in WebAssembly. Each table gives the median of five runs. The
output of each run is below the tables.

There are two kinds of run.

- **Node**, which `benchmarks/run.sh --backend wasm` makes with
  [wasm/node/bench.js](../wasm/node/bench.js). It covers the build with one
  thread and the build with more threads.
- **A browser**, which no script here can drive. Serve the example with
  `make serve-wasm`, paste [browser-bench.js](browser-bench.js) in the console
  of the page, and give the result to `yzma-bench update`. WebGPU needs this,
  because Node has no WebGPU.

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf)
and the tokens a second come from the generation loop of
[examples/wasm/chat](../examples/wasm/chat), which is the loop that the page
uses. These numbers are not comparable with the native tables, because the
prompt and the count of tokens are different.

## In Node

<!-- yzma:bench table node -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | wasm | Intel Core i9-13900HX | - | 13.8 | unknown | 2026-09-16 |
| CPU, more threads | wasm | Intel Core i9-13900HX | - | 107.3 | unknown | 2026-09-16 |
<!-- yzma:bench table end node -->

<!-- yzma:bench start node/cpu/wasm/i9-13900hx -->
### CPU, wasm, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"node","backend":"cpu","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":13.8,"yzma":"1.27.0","date":"2026-09-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 13.8 tokens a second.

<details><summary>The output of go test</summary>

```
$ node wasm/node/bench.js --dir build/wasm --model SmolLM-135M.Q2_K.gguf --tokens 64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
backend: cpu, 1 thread(s)
BenchmarkInference-1	1	4627621114 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4674945215 ns/op	13.7 tokens/s
BenchmarkInference-1	1	4620938628 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4630969609 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4624277457 ns/op	13.8 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	23.204s
```

</details>
<!-- yzma:bench end node/cpu/wasm/i9-13900hx -->

<!-- yzma:bench start node/cpu-threads/wasm/i9-13900hx -->
### CPU, more threads, wasm, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"node","backend":"cpu-threads","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":107.3,"yzma":"1.27.0","date":"2026-09-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 107.3 tokens a second.

<details><summary>The output of go test</summary>

```
$ node wasm/node/bench.js --dir build/wasm --model SmolLM-135M.Q2_K.gguf --tokens 64 --mt
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
backend: cpu-threads, 16 thread(s)
BenchmarkInference-16	1	596680962 ns/op	107.3 tokens/s
BenchmarkInference-16	1	539128970 ns/op	118.7 tokens/s
BenchmarkInference-16	1	597460792 ns/op	107.1 tokens/s
BenchmarkInference-16	1	614793468 ns/op	104.1 tokens/s
BenchmarkInference-16	1	587857077 ns/op	108.9 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	2.962s
```

</details>
<!-- yzma:bench end node/cpu-threads/wasm/i9-13900hx -->

## In a browser

<!-- yzma:bench table browser -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
<!-- yzma:bench table end browser -->

## Earlier measurements, by hand

These came from Chrome on one machine, an RTX 4070 with an Intel integrated GPU,
with the greedy sampler. The llama.cpp build is not known. They stay here until
a run replaces them.

| Model | Backend | Tokens a second |
| --- | --- | --- |
| SmolLM-135M Q2_K | one thread | 10.8 |
| SmolLM-135M Q2_K | more threads | 63.3 |
| SmolLM-135M Q2_K | WebGPU | 63.3 |
| Gemma 3 1B Q2_K | more threads | 18.5 |
| Gemma 3 1B Q2_K | WebGPU | 38.7 |

The GPU is faster on the larger model. On the smaller model the two results
agree, because each operation is too small to justify the transfer to the GPU.

An image gives a different result. This is a photo of 960 by 720 through the
projector of SmolVLM-256M Q8_0, and then 32 tokens of answer.

| Backend | Time for the image | Tokens a second |
| --- | --- | --- |
| more threads, in Chrome | 42.7 s | 96.9 |
| WebGPU, in Chrome | 1.6 s | 64.4 |
| one thread, in Node | 80 s | 17.4 |

A projector computes many numbers at the same time, which is the function of a
GPU. Thus the GPU is much faster for an image.
