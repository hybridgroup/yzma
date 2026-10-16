# macOS benchmarks

Benchmarks of yzma on macOS. Each table shows the median of five runs. The output
of each run, and the device details, are below the tables.

To add a machine or to regenerate these numbers, see
[how to run the benchmarks](README.md).

## Summary

Tokens a second on each machine.

| Machine | Text, CPU | Text, Metal | Text, BLAS | Multimodal, CPU | Multimodal, Metal | Multimodal, BLAS |
| --- | --- | --- | --- | --- | --- | --- |
| Apple M4 Pro | 903.5 | 511.5 | 503.2 | 904.3 | 1091.0 | 571.0 |

- For text, the CPU is the fastest backend, 1.8 times faster than Metal. The
  text model is very small, so the GPU has too little work per token.
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
| CPU | arm64 | Apple M4 Pro | - | 900.4 | b11669 | 2026-10-16 |
| blas | arm64 | Apple M4 Max | BLAS | 509.4 | b11146 | 2026-09-24 |
| blas | arm64 | Apple M4 Pro | BLAS | 515.6 | b11669 | 2026-10-16 |
| mtl | arm64 | Apple M4 Max | MTL0 | 597.7 | b11146 | 2026-09-24 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 543.1 | b11669 | 2026-10-16 |
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
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":900.4,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

Apple M4 Pro. 900.4 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     357	  33317254 ns/op	       900.4 tokens/s
BenchmarkInference-14    	     362	  35243495 ns/op	       851.2 tokens/s
BenchmarkInference-14    	     354	  33714460 ns/op	       889.8 tokens/s
BenchmarkInference-14    	     355	  32947024 ns/op	       910.6 tokens/s
BenchmarkInference-14    	     361	  32867271 ns/op	       912.8 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.266s
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
<!-- yzma:bench meta {"suite":"text","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":515.6,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

Apple M4 Pro. 515.6 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     204	  58186026 ns/op	       515.6 tokens/s
BenchmarkInference-14    	     205	  58375929 ns/op	       513.9 tokens/s
BenchmarkInference-14    	     206	  58024307 ns/op	       517.0 tokens/s
BenchmarkInference-14    	     206	  58490297 ns/op	       512.9 tokens/s
BenchmarkInference-14    	     205	  58106021 ns/op	       516.3 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.968s
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
<!-- yzma:bench meta {"suite":"text","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":543.1,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

Apple M4 Pro. 543.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Pro
BenchmarkInference-14    	     208	  56127176 ns/op	       534.5 tokens/s
BenchmarkInference-14    	     213	  54991052 ns/op	       545.5 tokens/s
BenchmarkInference-14    	     219	  55236779 ns/op	       543.1 tokens/s
BenchmarkInference-14    	     218	  54986459 ns/op	       545.6 tokens/s
BenchmarkInference-14    	     212	  55931401 ns/op	       536.4 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.110s
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
| CPU | arm64 | Apple M4 Pro | - | 994.1 | b11669 | 2026-10-16 |
| blas | arm64 | Apple M4 Max | BLAS | 690.2 | b11146 | 2026-09-24 |
| blas | arm64 | Apple M4 Pro | BLAS | 587.1 | b11669 | 2026-10-16 |
| mtl | arm64 | Apple M4 Max | MTL0 | 1690.0 | b11146 | 2026-09-24 |
| mtl | arm64 | Apple M4 Pro | MTL0 | 1116.0 | b11669 | 2026-10-16 |
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
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"rons-macbook-pro","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":994.1,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

Apple M4 Pro. 994.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      46	 237599780 ns/op	       988.6 tokens/s
BenchmarkMultimodalInference-14    	      50	 235894192 ns/op	       991.0 tokens/s
BenchmarkMultimodalInference-14    	      44	 232468694 ns/op	       999.6 tokens/s
BenchmarkMultimodalInference-14    	      45	 234607794 ns/op	       994.1 tokens/s
BenchmarkMultimodalInference-14    	      50	 232926051 ns/op	       999.2 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	61.240s
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
<!-- yzma:bench meta {"suite":"multimodal","backend":"blas","arch":"arm64","machine":"rons-macbook-pro","device":"BLAS","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":587.1,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

Apple M4 Pro. 587.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=BLAS
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      25	 405013457 ns/op	       580.0 tokens/s
BenchmarkMultimodalInference-14    	      28	 388242229 ns/op	       591.2 tokens/s
BenchmarkMultimodalInference-14    	      26	 390047570 ns/op	       591.6 tokens/s
BenchmarkMultimodalInference-14    	      27	 391974860 ns/op	       587.1 tokens/s
BenchmarkMultimodalInference-14    	      28	 392618121 ns/op	       586.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	58.904s
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
<!-- yzma:bench meta {"suite":"multimodal","backend":"mtl","arch":"arm64","machine":"rons-macbook-pro","device":"MTL0","label":"Apple M4 Pro","cpu":"Apple M4 Pro","tokens_per_second":1116,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

Apple M4 Pro. 1116.0 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=MTL0
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Pro
BenchmarkMultimodalInference-14    	      52	 208652849 ns/op	      1116 tokens/s
BenchmarkMultimodalInference-14    	      51	 208457348 ns/op	      1116 tokens/s
BenchmarkMultimodalInference-14    	      51	 207257586 ns/op	      1119 tokens/s
BenchmarkMultimodalInference-14    	      49	 208852390 ns/op	      1114 tokens/s
BenchmarkMultimodalInference-14    	      56	 208958072 ns/op	      1113 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	61.077s
```

</details>
<!-- yzma:bench end multimodal/mtl/arm64/rons-macbook-pro/mtl0 -->
