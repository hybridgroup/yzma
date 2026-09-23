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
| CPU | arm64 | Apple M4 Pro | - | 779.0 | b10964 | 2026-09-17 |
| blas | arm64 | Apple M4 Pro | BLAS | 505.0 | b10964 | 2026-09-17 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 513.3 | b10964 | 2026-09-17 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/arm64/rons-macbook-pro -->
### CPU, arm64, Apple M4 Pro
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":779,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

Apple M4 Pro. 779.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     310	  38499938 ns/op	       779.2 tokens/s
BenchmarkInference-14    	     310	  38483170 ns/op	       779.6 tokens/s
BenchmarkInference-14    	     309	  38537063 ns/op	       778.5 tokens/s
BenchmarkInference-14    	     307	  38666246 ns/op	       775.9 tokens/s
BenchmarkInference-14    	     310	  38509962 ns/op	       779.0 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.601s
```

</details>
<!-- yzma:bench end text/cpu/arm64/rons-macbook-pro -->

<!-- yzma:bench start text/blas/arm64/rons-macbook-pro/blas -->
### blas, arm64, Apple M4 Pro, BLAS
<!-- yzma:bench meta {"suite":"text","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":505,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

Apple M4 Pro. 505.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     201	  59410262 ns/op	       505.0 tokens/s
BenchmarkInference-14    	     201	  59402231 ns/op	       505.0 tokens/s
BenchmarkInference-14    	     201	  59387807 ns/op	       505.2 tokens/s
BenchmarkInference-14    	     200	  59606867 ns/op	       503.3 tokens/s
BenchmarkInference-14    	     201	  59414071 ns/op	       504.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.428s
```

</details>
<!-- yzma:bench end text/blas/arm64/rons-macbook-pro/blas -->

<!-- yzma:bench start text/mtl/arm64/rons-macbook-pro/mtl0 -->
### mtl, arm64, Apple M4 Pro, MTL0
<!-- yzma:bench meta {"suite":"text","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":513.3,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

Apple M4 Pro. 513.3 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     204	  58431832 ns/op	       513.4 tokens/s
BenchmarkInference-14    	     205	  58444535 ns/op	       513.3 tokens/s
BenchmarkInference-14    	     205	  58494539 ns/op	       512.9 tokens/s
BenchmarkInference-14    	     205	  58190734 ns/op	       515.5 tokens/s
BenchmarkInference-14    	     204	  58708515 ns/op	       511.0 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	75.834s
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
| CPU | arm64 | Apple M4 Pro | - | 799.4 | b10964 | 2026-09-17 |
| blas | arm64 | Apple M4 Pro | BLAS | 701.1 | b10964 | 2026-09-17 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 1085.0 | b10964 | 2026-09-17 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/arm64/rons-macbook-pro -->
### CPU, arm64, Apple M4 Pro
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":799.4,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

Apple M4 Pro. 799.4 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      40	 291283828 ns/op	       802.0 tokens/s
BenchmarkMultimodalInference-14    	      37	 291033359 ns/op	       800.2 tokens/s
BenchmarkMultimodalInference-14    	      38	 292351040 ns/op	       799.4 tokens/s
BenchmarkMultimodalInference-14    	      37	 292761112 ns/op	       797.2 tokens/s
BenchmarkMultimodalInference-14    	      40	 292947904 ns/op	       798.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.001s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/rons-macbook-pro -->

<!-- yzma:bench start multimodal/blas/arm64/rons-macbook-pro/blas -->
### blas, arm64, Apple M4 Pro, BLAS
<!-- yzma:bench meta {"suite":"multimodal","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":701.1,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

Apple M4 Pro. 701.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      34	 331062126 ns/op	       701.1 tokens/s
BenchmarkMultimodalInference-14    	      32	 329820315 ns/op	       702.0 tokens/s
BenchmarkMultimodalInference-14    	      32	 322404771 ns/op	       708.9 tokens/s
BenchmarkMultimodalInference-14    	      31	 331453765 ns/op	       701.0 tokens/s
BenchmarkMultimodalInference-14    	      32	 332634116 ns/op	       700.8 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	59.361s
```

</details>
<!-- yzma:bench end multimodal/blas/arm64/rons-macbook-pro/blas -->

<!-- yzma:bench start multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
### mtl, arm64, Apple M4 Pro, MTL0
<!-- yzma:bench meta {"suite":"multimodal","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":1085,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

Apple M4 Pro. 1085.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      55	 216544326 ns/op	      1077 tokens/s
BenchmarkMultimodalInference-14    	      49	 210986187 ns/op	      1099 tokens/s
BenchmarkMultimodalInference-14    	      48	 215437464 ns/op	      1085 tokens/s
BenchmarkMultimodalInference-14    	      49	 210890504 ns/op	      1100 tokens/s
BenchmarkMultimodalInference-14    	      49	 217819764 ns/op	      1077 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	59.926s
```

</details>
<!-- yzma:bench end multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
