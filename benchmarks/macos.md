# macOS benchmarks

Benchmarks of yzma on macOS. Each table gives the median of five runs. The output
of each run, and of the device, is below the tables.

To add a machine or to make these numbers again, see
[how to run the benchmarks](README.md).

## Summary

Tokens a second on each machine.

| Machine | Text, CPU | Text, Metal | Text, BLAS | Multimodal, CPU | Multimodal, Metal | Multimodal, BLAS |
| --- | --- | --- | --- | --- | --- | --- |
| Apple M4 Pro | 900.7 | 511.8 | 505.2 | 900.8 | 1098.0 | 583.5 |

- For text, the CPU is the fastest backend, 1.8 times faster than Metal. The
  text model is very small, thus the GPU has too little work for each token.
- For multimodal, Metal is the fastest backend, 22 percent faster than the CPU.
- BLAS is the slowest backend for both suites.

## Text model benchmarks

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf).
The code is [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go).
The benchmark uses 4 threads on each machine, see
[the thread count](README.md#run-them).

<!-- yzma:bench table text -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | arm64 | Apple M4 Pro | - | 903.5 | b11146 | 2026-09-24 |
| blas | arm64 | Apple M4 Pro | BLAS | 503.2 | b11146 | 2026-09-24 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 511.5 | b11146 | 2026-09-24 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/arm64/rons-macbook-pro -->
### CPU, arm64, Apple M4 Pro
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":903.5,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Pro. 903.5 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     360	  33203822 ns/op	       903.5 tokens/s
BenchmarkInference-14    	     361	  33175729 ns/op	       904.3 tokens/s
BenchmarkInference-14    	     361	  33179927 ns/op	       904.2 tokens/s
BenchmarkInference-14    	     361	  33295331 ns/op	       901.0 tokens/s
BenchmarkInference-14    	     360	  33227200 ns/op	       902.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.156s
```

</details>
<!-- yzma:bench end text/cpu/arm64/rons-macbook-pro -->

<!-- yzma:bench start text/blas/arm64/rons-macbook-pro/blas -->
### blas, arm64, Apple M4 Pro, BLAS
<!-- yzma:bench meta {"suite":"text","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":503.2,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Pro. 503.2 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     201	  59524000 ns/op	       504.0 tokens/s
BenchmarkInference-14    	     200	  59656747 ns/op	       502.9 tokens/s
BenchmarkInference-14    	     200	  59595159 ns/op	       503.4 tokens/s
BenchmarkInference-14    	     200	  59614975 ns/op	       503.2 tokens/s
BenchmarkInference-14    	     199	  59865650 ns/op	       501.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.695s
```

</details>
<!-- yzma:bench end text/blas/arm64/rons-macbook-pro/blas -->

<!-- yzma:bench start text/mtl/arm64/rons-macbook-pro/mtl0 -->
### mtl, arm64, Apple M4 Pro, MTL0
<!-- yzma:bench meta {"suite":"text","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":511.5,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Pro. 511.5 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     200	  58655590 ns/op	       511.5 tokens/s
BenchmarkInference-14    	     208	  57907047 ns/op	       518.1 tokens/s
BenchmarkInference-14    	     204	  58794825 ns/op	       510.2 tokens/s
BenchmarkInference-14    	     206	  58605772 ns/op	       511.9 tokens/s
BenchmarkInference-14    	     202	  58951936 ns/op	       508.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	76.547s
```

</details>
<!-- yzma:bench end text/mtl/arm64/rons-macbook-pro/mtl0 -->

## Multimodal model benchmarks

The model is
[SmolVLM-256M-Instruct-Q8_0.gguf](https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/SmolVLM-256M-Instruct-Q8_0.gguf)
with its
[projector](https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf).
The code is [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go).

<!-- yzma:bench table multimodal -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | arm64 | Apple M4 Pro | - | 904.3 | b11146 | 2026-09-24 |
| blas | arm64 | Apple M4 Pro | BLAS | 571.0 | b11146 | 2026-09-24 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 1091.0 | b11146 | 2026-09-24 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/arm64/rons-macbook-pro -->
### CPU, arm64, Apple M4 Pro
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":904.3,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Pro. 904.3 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      45	 258986487 ns/op	       905.6 tokens/s
BenchmarkMultimodalInference-14    	      40	 263477632 ns/op	       895.9 tokens/s
BenchmarkMultimodalInference-14    	      40	 254300005 ns/op	       912.0 tokens/s
BenchmarkMultimodalInference-14    	      40	 258140909 ns/op	       904.3 tokens/s
BenchmarkMultimodalInference-14    	      44	 262384770 ns/op	       897.2 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	60.287s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/rons-macbook-pro -->

<!-- yzma:bench start multimodal/blas/arm64/rons-macbook-pro/blas -->
### blas, arm64, Apple M4 Pro, BLAS
<!-- yzma:bench meta {"suite":"multimodal","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":571,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Pro. 571.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      27	 410372920 ns/op	       573.6 tokens/s
BenchmarkMultimodalInference-14    	      25	 409109337 ns/op	       569.3 tokens/s
BenchmarkMultimodalInference-14    	      25	 405644065 ns/op	       571.0 tokens/s
BenchmarkMultimodalInference-14    	      27	 401668620 ns/op	       573.1 tokens/s
BenchmarkMultimodalInference-14    	      25	 413486773 ns/op	       566.3 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	58.886s
```

</details>
<!-- yzma:bench end multimodal/blas/arm64/rons-macbook-pro/blas -->

<!-- yzma:bench start multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
### mtl, arm64, Apple M4 Pro, MTL0
<!-- yzma:bench meta {"suite":"multimodal","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":1091,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Pro. 1091.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      48	 214680949 ns/op	      1086 tokens/s
BenchmarkMultimodalInference-14    	      49	 212384048 ns/op	      1093 tokens/s
BenchmarkMultimodalInference-14    	      52	 213721188 ns/op	      1089 tokens/s
BenchmarkMultimodalInference-14    	      55	 213136656 ns/op	      1091 tokens/s
BenchmarkMultimodalInference-14    	      54	 215383050 ns/op	      1092 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	61.882s
```

</details>
<!-- yzma:bench end multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
