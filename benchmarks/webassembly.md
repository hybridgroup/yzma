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

Both kinds read the llama.cpp build from `yzma-install.json` of the build
directory, thus each result says which build made it.

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
| CPU | wasm | Intel Core i9-13900HX | - | 13.8 | b11017 | 2026-09-22 |
| CPU, more threads | wasm | Intel Core i9-13900HX | - | 99.9 | b11017 | 2026-09-22 |
<!-- yzma:bench table end node -->

<!-- yzma:bench start node/cpu/wasm/i9-13900hx -->
### CPU, wasm, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"node","backend":"cpu","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":13.8,"llamacpp":"b11017","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 13.8 tokens a second.

<details><summary>The output of go test</summary>

```
$ node wasm/node/bench.js --dir build/wasm --model SmolLM-135M.Q2_K.gguf --tokens 64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
backend: cpu, 1 thread(s)
llama.cpp: b11017
BenchmarkInference-1	1	4627621114 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4627621114 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4614275415 ns/op	13.9 tokens/s
BenchmarkInference-1	1	4637681159 ns/op	13.8 tokens/s
BenchmarkInference-1	1	4614275415 ns/op	13.9 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	23.141s
```

</details>
<!-- yzma:bench end node/cpu/wasm/i9-13900hx -->

<!-- yzma:bench start node/cpu-threads/wasm/i9-13900hx -->
### CPU, more threads, wasm, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"node","backend":"cpu-threads","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":99.9,"llamacpp":"b11017","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 99.9 tokens a second.

<details><summary>The output of go test</summary>

```
$ node wasm/node/bench.js --dir build/wasm --model SmolLM-135M.Q2_K.gguf --tokens 64 --mt
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
backend: cpu-threads, 16 thread(s)
llama.cpp: b11017
BenchmarkInference-16	1	640897256 ns/op	99.9 tokens/s
BenchmarkInference-16	1	622507538 ns/op	102.8 tokens/s
BenchmarkInference-16	1	601447232 ns/op	106.4 tokens/s
BenchmarkInference-16	1	653461303 ns/op	97.9 tokens/s
BenchmarkInference-16	1	675176706 ns/op	94.8 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	3.218s
```

</details>
<!-- yzma:bench end node/cpu-threads/wasm/i9-13900hx -->

## In a browser

<!-- yzma:bench table browser -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU, more threads | wasm | Intel Core i9-13900HX, Chrome | - | 90.1 | b11017 | 2026-09-23 |
<!-- yzma:bench table end browser -->

<!-- yzma:bench start browser/cpu-threads/wasm/i9-13900hx -->
### CPU, more threads, wasm, Intel Core i9-13900HX, Chrome
<!-- yzma:bench meta {"suite":"browser","backend":"cpu-threads","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX, Chrome","cpu":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36","tokens_per_second":90.1,"llamacpp":"b11017","yzma":"1.27.0","date":"2026-09-23"} -->

Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36. 90.1 tokens a second.

<details><summary>The output of go test</summary>

```
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36
backend: cpu-threads, 16 threads
llama.cpp: b11017
BenchmarkInference-32	1	706245862 ns/op	90.6 tokens/s
BenchmarkInference-32	1	687950124 ns/op	93.0 tokens/s
BenchmarkInference-32	1	725130297 ns/op	88.3 tokens/s
BenchmarkInference-32	1	720720721 ns/op	88.8 tokens/s
BenchmarkInference-32	1	710164225 ns/op	90.1 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	3.572s
```

</details>
<!-- yzma:bench end browser/cpu-threads/wasm/i9-13900hx -->
