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

## Summary

Tokens a second on each machine.

| Machine | Node, CPU | Node, more threads | Chrome, more threads |
| --- | --- | --- | --- |
| Intel Core i9-13900HX | 13.8 | 105.7 | 91.8 |

- The build with more threads is 7.7 times faster than the build with one
  thread.
- In Chrome, the build with more threads gives 87 percent of the speed in Node.
- There is no WebGPU result yet.

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
| CPU, more threads | wasm | Intel Core i9-13900HX, Chrome | - | 91.8 | b11146 | 2026-09-24 |
<!-- yzma:bench table end browser -->

<!-- yzma:bench start browser/cpu-threads/wasm/i9-13900hx -->
### CPU, more threads, wasm, Intel Core i9-13900HX, Chrome
<!-- yzma:bench meta {"suite":"browser","backend":"cpu-threads","arch":"wasm","machine":"i9-13900hx","label":"Intel Core i9-13900HX, Chrome","cpu":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36","tokens_per_second":91.8,"llamacpp":"b11146","date":"2026-09-24"} -->

Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36. 91.8 tokens a second.

<details><summary>The output of go test</summary>

```
$ browser-bench.js mode=auto tokens=64
goos: js
goarch: wasm
pkg: github.com/hybridgroup/yzma/examples/wasm/chat
cpu: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36
backend: cpu-threads, 16 threads
llama.cpp: b11146
BenchmarkInference-32	1	690995465 ns/op	92.6 tokens/s
BenchmarkInference-32	1	700755502 ns/op	91.3 tokens/s
BenchmarkInference-32	1	696864111 ns/op	91.8 tokens/s
BenchmarkInference-32	1	693316000 ns/op	92.3 tokens/s
BenchmarkInference-32	1	722266110 ns/op	88.6 tokens/s
PASS
ok	github.com/hybridgroup/yzma/examples/wasm/chat	3.527s
```

</details>
<!-- yzma:bench end browser/cpu-threads/wasm/i9-13900hx -->
