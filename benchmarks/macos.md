# macOS benchmarks

Benchmarks of yzma on macOS. Each table gives the median of five runs. The output
of each run, and of the device, is below the tables.

To add a machine or to make these numbers again, see
[how to run the benchmarks](README.md).

## Summary

Tokens a second on each machine.

| Machine | Text, CPU | Text, Metal | Text, BLAS | Multimodal, CPU | Multimodal, Metal | Multimodal, BLAS |
| --- | --- | --- | --- | --- | --- | --- |
| Apple M4 Pro | 903.5 | 511.5 | 503.2 | 904.3 | 1091.0 | 571.0 |

- For text, the CPU is the fastest backend, 1.8 times faster than Metal. The
  text model is very small, thus the GPU has too little work for each token.
- For multimodal, Metal is the fastest backend, 21 percent faster than the CPU.
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
| CPU | arm64 | Apple M4 Max | - | 959.7 | b11146 | 2026-09-24 |
| CPU | arm64 | Apple M4 Pro | - | 903.5 | b11146 | 2026-09-24 |
| blas | arm64 | Apple M4 Max | BLAS | 509.4 | b11146 | 2026-09-24 |
| blas | arm64 | Apple M4 Pro | BLAS | 503.2 | b11146 | 2026-09-24 |
| mtl | arm64 | Apple M4 Max | MTL0 | 597.7 | b11146 | 2026-09-24 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 511.5 | b11146 | 2026-09-24 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/arm64/monzi14 -->
### CPU, arm64, Apple M4 Max
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"monzi14","label":"Apple M4 Max","cpu":"Apple M4 Max","tokens_per_second":959.7,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Max. 959.7 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Max
BenchmarkInference-16    	     384	  31171824 ns/op	       962.4 tokens/s
BenchmarkInference-16    	     376	  31622733 ns/op	       948.7 tokens/s
BenchmarkInference-16    	     381	  31258531 ns/op	       959.7 tokens/s
BenchmarkInference-16    	     380	  31386417 ns/op	       955.8 tokens/s
BenchmarkInference-16    	     387	  31182761 ns/op	       962.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.768s
```

</details>
<!-- yzma:bench end text/cpu/arm64/monzi14 -->

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

<!-- yzma:bench start text/blas/arm64/monzi14/blas -->
### blas, arm64, Apple M4 Max, BLAS
<!-- yzma:bench meta {"suite":"text","backend":"blas","arch":"arm64","machine":"monzi14","device":"BLAS","label":"Apple M4 Max","cpu":"Apple M4 Max","tokens_per_second":509.4,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Max. 509.4 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Max
BenchmarkInference-16    	     202	  58916992 ns/op	       509.2 tokens/s
BenchmarkInference-16    	     204	  58797734 ns/op	       510.2 tokens/s
BenchmarkInference-16    	     202	  58959088 ns/op	       508.8 tokens/s
BenchmarkInference-16    	     204	  58895263 ns/op	       509.4 tokens/s
BenchmarkInference-16    	     204	  58729855 ns/op	       510.8 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.377s
```

</details>
<!-- yzma:bench end text/blas/arm64/monzi14/blas -->

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

<!-- yzma:bench start text/mtl/arm64/monzi14/mtl0 -->
### mtl, arm64, Apple M4 Max, MTL0
<!-- yzma:bench meta {"suite":"text","backend":"mtl","arch":"arm64","machine":"monzi14","device":"MTL0","label":"Apple M4 Max","cpu":"Apple M4 Max","tokens_per_second":597.7,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Max. 597.7 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Max
BenchmarkInference-16    	     240	  50191677 ns/op	       597.7 tokens/s
BenchmarkInference-16    	     238	  50224946 ns/op	       597.3 tokens/s
BenchmarkInference-16    	     236	  50655662 ns/op	       592.2 tokens/s
BenchmarkInference-16    	     238	  49764861 ns/op	       602.8 tokens/s
BenchmarkInference-16    	     246	  48400489 ns/op	       619.8 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	77.084s
```

</details>
<!-- yzma:bench end text/mtl/arm64/monzi14/mtl0 -->

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
| CPU | arm64 | Apple M4 Max | - | 1143.0 | b11146 | 2026-09-24 |
| CPU | arm64 | Apple M4 Pro | - | 904.3 | b11146 | 2026-09-24 |
| blas | arm64 | Apple M4 Max | BLAS | 690.2 | b11146 | 2026-09-24 |
| blas | arm64 | Apple M4 Pro | BLAS | 571.0 | b11146 | 2026-09-24 |
| mtl | arm64 | Apple M4 Max | MTL0 | 1690.0 | b11146 | 2026-09-24 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 1091.0 | b11146 | 2026-09-24 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/arm64/monzi14 -->
### CPU, arm64, Apple M4 Max
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"monzi14","label":"Apple M4 Max","cpu":"Apple M4 Max","tokens_per_second":1143,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Max. 1143.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Max
BenchmarkMultimodalInference-16    	      62	 246963139 ns/op	      1009 tokens/s
BenchmarkMultimodalInference-16    	      54	 205005366 ns/op	      1143 tokens/s
BenchmarkMultimodalInference-16    	      50	 209100112 ns/op	      1123 tokens/s
BenchmarkMultimodalInference-16    	      51	 203381570 ns/op	      1143 tokens/s
BenchmarkMultimodalInference-16    	      57	 203550999 ns/op	      1152 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	78.464s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/monzi14 -->

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

<!-- yzma:bench start multimodal/blas/arm64/monzi14/blas -->
### blas, arm64, Apple M4 Max, BLAS
<!-- yzma:bench meta {"suite":"multimodal","backend":"blas","arch":"arm64","machine":"monzi14","device":"BLAS","label":"Apple M4 Max","cpu":"Apple M4 Max","tokens_per_second":690.2,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Max. 690.2 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Max
BenchmarkMultimodalInference-16    	      33	 336156308 ns/op	       690.2 tokens/s
BenchmarkMultimodalInference-16    	      33	 329855066 ns/op	       699.0 tokens/s
BenchmarkMultimodalInference-16    	      33	 333096062 ns/op	       692.3 tokens/s
BenchmarkMultimodalInference-16    	      33	 343581660 ns/op	       683.5 tokens/s
BenchmarkMultimodalInference-16    	      33	 345061419 ns/op	       670.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	61.712s
```

</details>
<!-- yzma:bench end multimodal/blas/arm64/monzi14/blas -->

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

<!-- yzma:bench start multimodal/mtl/arm64/monzi14/mtl0 -->
### mtl, arm64, Apple M4 Max, MTL0
<!-- yzma:bench meta {"suite":"multimodal","backend":"mtl","arch":"arm64","machine":"monzi14","device":"MTL0","label":"Apple M4 Max","cpu":"Apple M4 Max","tokens_per_second":1690,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

Apple M4 Max. 1690.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Max
BenchmarkMultimodalInference-16    	      82	 142364081 ns/op	      1641 tokens/s
BenchmarkMultimodalInference-16    	      81	 140955126 ns/op	      1660 tokens/s
BenchmarkMultimodalInference-16    	      88	 135926264 ns/op	      1702 tokens/s
BenchmarkMultimodalInference-16    	      79	 137054854 ns/op	      1690 tokens/s
BenchmarkMultimodalInference-16    	      87	 135459486 ns/op	      1711 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	64.799s
```

</details>
<!-- yzma:bench end multimodal/mtl/arm64/monzi14/mtl0 -->

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
