# macOS benchmarks

Benchmarks of yzma on macOS. Each table gives the median of five runs. The output
of each run, and of the device, is below the tables.

To add a machine or to make these numbers again, see
[how to run the benchmarks](README.md).

## Text model benchmarks

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf).
The code is [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go).

<!-- yzma:bench table text -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | arm64 | Apple M4 Pro | - | 774.2 | b10964 | 2026-09-16 |
| Metal | arm64 | Apple M4 Max with 128 GB RAM | - | 577.1 | unknown | unknown |
| blas | arm64 | Apple M4 Pro | BLAS | 505.6 | b10964 | 2026-09-16 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 511.4 | b10964 | 2026-09-16 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/arm64/rons-macbook-pro -->
### CPU, arm64, Apple M4 Pro
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":774.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

Apple M4 Pro. 774.2 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     304	  38703800 ns/op	       775.1 tokens/s
BenchmarkInference-14    	     308	  38747342 ns/op	       774.2 tokens/s
BenchmarkInference-14    	     290	  40030109 ns/op	       749.4 tokens/s
BenchmarkInference-14    	     312	  39050762 ns/op	       768.2 tokens/s
BenchmarkInference-14    	     312	  38369717 ns/op	       781.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	75.340s
```

</details>
<!-- yzma:bench end text/cpu/arm64/rons-macbook-pro -->

<!-- yzma:bench start text/metal/arm64/m4-max -->
### Metal, arm64, Apple M4 Max with 128 GB RAM
<!-- yzma:bench meta {"suite":"text","backend":"metal","arch":"arm64","machine":"m4-max","label":"Apple M4 Max with 128 GB RAM","cpu":"Apple M4 Max","tokens_per_second":577.1} -->

Apple M4 Max. 577.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ go test -run none -benchtime=10s -count=5 -bench BenchmarkInference -nctx=16000
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Max
BenchmarkInference-16	  230		52168178 ns/op	575.1 tokens/s
BenchmarkInference-16	  234		51482815 ns/op	582.7 tokens/s
BenchmarkInference-16	  230		51729562 ns/op	579.9 tokens/s
BenchmarkInference-16	  230		52075140 ns/op	576.1 tokens/s
BenchmarkInference-16	  230		51981549 ns/op	577.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.042s
```

</details>
<!-- yzma:bench end text/metal/arm64/m4-max -->

<!-- yzma:bench start text/blas/arm64/rons-macbook-pro/blas -->
### blas, arm64, Apple M4 Pro, BLAS
<!-- yzma:bench meta {"suite":"text","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":505.6,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

Apple M4 Pro. 505.6 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     201	  59307228 ns/op	       505.8 tokens/s
BenchmarkInference-14    	     201	  59297460 ns/op	       505.9 tokens/s
BenchmarkInference-14    	     200	  59531441 ns/op	       503.9 tokens/s
BenchmarkInference-14    	     201	  59513729 ns/op	       504.1 tokens/s
BenchmarkInference-14    	     201	  59336192 ns/op	       505.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.289s
```

</details>
<!-- yzma:bench end text/blas/arm64/rons-macbook-pro/blas -->

<!-- yzma:bench start text/mtl/arm64/rons-macbook-pro/mtl0 -->
### mtl, arm64, Apple M4 Pro, MTL0
<!-- yzma:bench meta {"suite":"text","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":511.4,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

Apple M4 Pro. 511.4 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     195	  59304356 ns/op	       505.9 tokens/s
BenchmarkInference-14    	     204	  58668166 ns/op	       511.4 tokens/s
BenchmarkInference-14    	     204	  58475187 ns/op	       513.0 tokens/s
BenchmarkInference-14    	     205	  58465505 ns/op	       513.1 tokens/s
BenchmarkInference-14    	     204	  58811355 ns/op	       510.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	75.849s
```

</details>
<!-- yzma:bench end text/mtl/arm64/rons-macbook-pro/mtl0 -->

## Multimodal model benchmarks

The model is
[Qwen3-VL-2B-Instruct.Q4_K_M.gguf](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.Q4_K_M.gguf)
with its
[projector](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf).
The code is [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go).

<!-- yzma:bench table multimodal -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | arm64 | Apple M4 Pro | - | 102.4 | b10964 | 2026-09-16 |
| Metal | arm64 | Apple M4 Max with 128 GB RAM | - | 771.9 | unknown | unknown |
| blas | arm64 | Apple M4 Pro | BLAS | 159.9 | b10964 | 2026-09-16 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 472.5 | b10964 | 2026-09-16 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/arm64/rons-macbook-pro -->
### CPU, arm64, Apple M4 Pro
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":102.4,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

Apple M4 Pro. 102.4 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	       2	10843545958 ns/op	       106.5 tokens/s
BenchmarkMultimodalInference-14    	       1	11664156459 ns/op	       103.9 tokens/s
BenchmarkMultimodalInference-14    	       1	12752652167 ns/op	        99.82 tokens/s
BenchmarkMultimodalInference-14    	       1	12370603583 ns/op	       101.1 tokens/s
BenchmarkMultimodalInference-14    	       1	12070906000 ns/op	       102.4 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	76.389s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/rons-macbook-pro -->

<!-- yzma:bench start multimodal/metal/arm64/m4-max -->
### Metal, arm64, Apple M4 Max with 128 GB RAM
<!-- yzma:bench meta {"suite":"multimodal","backend":"metal","arch":"arm64","machine":"m4-max","label":"Apple M4 Max with 128 GB RAM","cpu":"Apple M4 Max","tokens_per_second":771.9} -->

Apple M4 Max. 771.9 tokens a second.

<details><summary>The output of go test</summary>

```
$ go test -run none -benchtime=10s -count=5 -bench BenchmarkMultimodalInference -nctx=16000
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Max
BenchmarkMultimodalInference-16		10		1577948683 ns/op	788.9 tokens/s
BenchmarkMultimodalInference-16		12		1243692014 ns/op	910.8 tokens/s
BenchmarkMultimodalInference-16		 7		1654741804 ns/op	737.2 tokens/s
BenchmarkMultimodalInference-16		 7		1568106947 ns/op	771.9 tokens/s
BenchmarkMultimodalInference-16		10		1704669371 ns/op	706.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	76.644s
```

</details>
<!-- yzma:bench end multimodal/metal/arm64/m4-max -->

<!-- yzma:bench start multimodal/blas/arm64/rons-macbook-pro/blas -->
### blas, arm64, Apple M4 Pro, BLAS
<!-- yzma:bench meta {"suite":"multimodal","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":159.9,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

Apple M4 Pro. 159.9 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	       2	8038776146 ns/op	       149.6 tokens/s
BenchmarkMultimodalInference-14    	       2	8820003708 ns/op	       138.8 tokens/s
BenchmarkMultimodalInference-14    	       2	7393613646 ns/op	       159.9 tokens/s
BenchmarkMultimodalInference-14    	       2	7075656958 ns/op	       165.5 tokens/s
BenchmarkMultimodalInference-14    	       2	6736491521 ns/op	       172.0 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	82.167s
```

</details>
<!-- yzma:bench end multimodal/blas/arm64/rons-macbook-pro/blas -->

<!-- yzma:bench start multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
### mtl, arm64, Apple M4 Pro, MTL0
<!-- yzma:bench meta {"suite":"multimodal","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":472.5,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

Apple M4 Pro. 472.5 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	       4	2853429073 ns/op	       431.6 tokens/s
BenchmarkMultimodalInference-14    	       4	2597057125 ns/op	       466.4 tokens/s
BenchmarkMultimodalInference-14    	       6	2253171528 ns/op	       516.7 tokens/s
BenchmarkMultimodalInference-14    	       6	1960910243 ns/op	       572.8 tokens/s
BenchmarkMultimodalInference-14    	       4	2549525906 ns/op	       472.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	63.945s
```

</details>
<!-- yzma:bench end multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
