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
| Metal | arm64 | Apple M4 Max with 128 GB RAM | - | 577.1 | unknown | unknown |
<!-- yzma:bench table end text -->

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

## Multimodal model benchmarks

The model is
[Qwen3-VL-2B-Instruct.Q4_K_M.gguf](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.Q4_K_M.gguf)
with its
[projector](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf).
The code is [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go).

<!-- yzma:bench table multimodal -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| Metal | arm64 | Apple M4 Max with 128 GB RAM | - | 771.9 | unknown | unknown |
<!-- yzma:bench table end multimodal -->

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
