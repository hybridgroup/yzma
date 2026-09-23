# macOS benchmarks

Benchmarks of yzma on macOS. Each table gives the median of five runs. The output
of each run, and of the device, is below the tables.

To add a machine or to make these numbers again, see
[how to run the benchmarks](README.md).

## Text model benchmarks

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf).
The code is [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go).
The benchmark uses 4 threads on each machine, see
[the thread count](README.md#run-them).

<!-- yzma:bench table text -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | arm64 | Apple M4 Pro | - | 900.7 | b10964 | 2026-09-23 |
| blas | arm64 | Apple M4 Pro | BLAS | 505.2 | b10964 | 2026-09-23 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 511.8 | b10964 | 2026-09-23 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/arm64/rons-macbook-pro -->
### CPU, arm64, Apple M4 Pro
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":900.7,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

Apple M4 Pro. 900.7 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     358	  33306418 ns/op	       900.7 tokens/s
BenchmarkInference-14    	     360	  33389123 ns/op	       898.5 tokens/s
BenchmarkInference-14    	     358	  33323882 ns/op	       900.3 tokens/s
BenchmarkInference-14    	     360	  33278179 ns/op	       901.5 tokens/s
BenchmarkInference-14    	     360	  33273386 ns/op	       901.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.849s
```

</details>
<!-- yzma:bench end text/cpu/arm64/rons-macbook-pro -->

<!-- yzma:bench start text/blas/arm64/rons-macbook-pro/blas -->
### blas, arm64, Apple M4 Pro, BLAS
<!-- yzma:bench meta {"suite":"text","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":505.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

Apple M4 Pro. 505.2 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     201	  59311649 ns/op	       505.8 tokens/s
BenchmarkInference-14    	     201	  59377787 ns/op	       505.2 tokens/s
BenchmarkInference-14    	     200	  59656495 ns/op	       502.9 tokens/s
BenchmarkInference-14    	     201	  59373107 ns/op	       505.3 tokens/s
BenchmarkInference-14    	     201	  59422324 ns/op	       504.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.515s
```

</details>
<!-- yzma:bench end text/blas/arm64/rons-macbook-pro/blas -->

<!-- yzma:bench start text/mtl/arm64/rons-macbook-pro/mtl0 -->
### mtl, arm64, Apple M4 Pro, MTL0
<!-- yzma:bench meta {"suite":"text","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":511.8,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

Apple M4 Pro. 511.8 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     204	  58683391 ns/op	       511.2 tokens/s
BenchmarkInference-14    	     204	  58614407 ns/op	       511.8 tokens/s
BenchmarkInference-14    	     204	  58466797 ns/op	       513.1 tokens/s
BenchmarkInference-14    	     204	  58566325 ns/op	       512.2 tokens/s
BenchmarkInference-14    	     204	  58636609 ns/op	       511.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.328s
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
| CPU | arm64 | Apple M4 Pro | - | 900.8 | b10964 | 2026-09-23 |
| blas | arm64 | Apple M4 Pro | BLAS | 583.5 | b10964 | 2026-09-23 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 1098.0 | b10964 | 2026-09-23 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/arm64/rons-macbook-pro -->
### CPU, arm64, Apple M4 Pro
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":900.8,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

Apple M4 Pro. 900.8 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      43	 252746628 ns/op	       916.9 tokens/s
BenchmarkMultimodalInference-14    	      43	 322871652 ns/op	       798.6 tokens/s
BenchmarkMultimodalInference-14    	      46	 255015899 ns/op	       911.1 tokens/s
BenchmarkMultimodalInference-14    	      40	 260702073 ns/op	       900.8 tokens/s
BenchmarkMultimodalInference-14    	      42	 267103955 ns/op	       888.4 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	64.169s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/rons-macbook-pro -->

<!-- yzma:bench start multimodal/blas/arm64/rons-macbook-pro/blas -->
### blas, arm64, Apple M4 Pro, BLAS
<!-- yzma:bench meta {"suite":"multimodal","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":583.5,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

Apple M4 Pro. 583.5 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      28	 400229002 ns/op	       579.7 tokens/s
BenchmarkMultimodalInference-14    	      28	 402351662 ns/op	       578.6 tokens/s
BenchmarkMultimodalInference-14    	      28	 396971705 ns/op	       583.5 tokens/s
BenchmarkMultimodalInference-14    	      27	 400842864 ns/op	       584.0 tokens/s
BenchmarkMultimodalInference-14    	      28	 396421071 ns/op	       585.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	61.780s
```

</details>
<!-- yzma:bench end multimodal/blas/arm64/rons-macbook-pro/blas -->

<!-- yzma:bench start multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
### mtl, arm64, Apple M4 Pro, MTL0
<!-- yzma:bench meta {"suite":"multimodal","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":1098,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

Apple M4 Pro. 1098.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      55	 210081728 ns/op	      1103 tokens/s
BenchmarkMultimodalInference-14    	      55	 213161283 ns/op	      1093 tokens/s
BenchmarkMultimodalInference-14    	      49	 211227783 ns/op	      1098 tokens/s
BenchmarkMultimodalInference-14    	      51	 210659331 ns/op	      1099 tokens/s
BenchmarkMultimodalInference-14    	      54	 216634327 ns/op	      1080 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.709s
```

</details>
<!-- yzma:bench end multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
