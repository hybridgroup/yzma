# Engine comparison benchmarks

These benchmarks are to compare inference performance using 3 different engines that support GGUF models:

- yzma
- ollama
- Docker Model Runner

## Summary

| Suite | Result | Evidence |
| --- | --- | --- |
| Embeddings | yzma 3.0 to 3.3 times faster, 1.4 ms against 4.2 ms and 4.6 ms | Five runs, no overlap, each engine at 29 prompt tokens and a vector of 384 |
| Text | yzma 5.4 to 11.8 percent faster, 0.5 to 0.6 of the time to the first token | Five runs, no overlap, each engine at the same count of prompt tokens |
| Images | No numbers. The code runs with `--suite multimodal` | The engines preprocess an image in different ways |

## Text

<!-- yzma:bench table compare-text -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 119.0 | 12.2 | 134.5 | 1.28.0 | 2026-09-24 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 106.4 | 23.5 | 150.4 | 0.34.4 | 2026-09-24 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 107.4 | 25.6 | 149.0 | v1.2.8 | 2026-09-24 |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 168.7 | 8.7 | 94.8 | 1.28.0 | 2026-09-24 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 153.5 | 16.3 | 104.2 | 0.34.4 | 2026-09-24 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 160.1 | 14.2 | 100.0 | v1.2.8 | 2026-09-24 |
<!-- yzma:bench table end compare-text -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":119,"ttft_ms":12.17,"total_ms":134.5,"prompt_tokens":23,"llamacpp":"b11146","engine_version":"1.28.0","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 119.0 tokens a second. 12.2 ms to the first token. 134.5 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=yzma -suite=text -tokens=16 -nctx=8192 -model=/home/ron/models/gemma-4-E2B-it-Q4_K_M.gguf -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 136237099 ns/op	        23.00 prompt_tokens	       118.0 tokens/s	       135.6 total_ms	        12.35 ttft_ms
BenchmarkCompare-32    	      20	 134667062 ns/op	        23.00 prompt_tokens	       119.4 tokens/s	       134.0 total_ms	        12.14 ttft_ms
BenchmarkCompare-32    	      20	 135116531 ns/op	        23.00 prompt_tokens	       119.0 tokens/s	       134.5 total_ms	        12.23 ttft_ms
BenchmarkCompare-32    	      20	 134859591 ns/op	        23.00 prompt_tokens	       119.2 tokens/s	       134.2 total_ms	        12.17 ttft_ms
BenchmarkCompare-32    	      20	 135314549 ns/op	        23.00 prompt_tokens	       118.8 tokens/s	       134.7 total_ms	        12.14 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.480s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":106.4,"ttft_ms":23.52,"total_ms":150.4,"prompt_tokens":23,"engine_version":"0.34.4","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 106.4 tokens a second. 23.5 ms to the first token. 150.4 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=ollama -suite=text -tokens=16 -nctx=8192 -server-model=yzma-bench-gemma4-e2b -ollama-url=http://localhost:11434/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 150461333 ns/op	        23.00 prompt_tokens	       106.4 tokens/s	       150.4 total_ms	        25.25 ttft_ms
BenchmarkCompare-32    	      20	 150550513 ns/op	        23.00 prompt_tokens	       106.3 tokens/s	       150.5 total_ms	        23.75 ttft_ms
BenchmarkCompare-32    	      20	 149871668 ns/op	        23.00 prompt_tokens	       106.8 tokens/s	       149.8 total_ms	        23.52 ttft_ms
BenchmarkCompare-32    	      20	 150439988 ns/op	        23.00 prompt_tokens	       106.4 tokens/s	       150.4 total_ms	        23.51 ttft_ms
BenchmarkCompare-32    	      20	 149971711 ns/op	        23.00 prompt_tokens	       106.7 tokens/s	       149.9 total_ms	        23.21 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.062s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":107.4,"ttft_ms":25.6,"total_ms":149,"prompt_tokens":23,"engine_version":"v1.2.8","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 107.4 tokens a second. 25.6 ms to the first token. 149.0 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=dmr -suite=text -tokens=16 -nctx=8192 -server-model=hf.co/unsloth/gemma-4-e2b-it-gguf:q4_k_m -dmr-url=http://localhost:12434/engines/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 149661989 ns/op	        23.00 prompt_tokens	       106.9 tokens/s	       149.6 total_ms	        25.91 ttft_ms
BenchmarkCompare-32    	      20	 148773933 ns/op	        23.00 prompt_tokens	       107.6 tokens/s	       148.7 total_ms	        25.60 ttft_ms
BenchmarkCompare-32    	      20	 149095372 ns/op	        23.00 prompt_tokens	       107.4 tokens/s	       149.0 total_ms	        25.80 ttft_ms
BenchmarkCompare-32    	      20	 148830665 ns/op	        23.00 prompt_tokens	       107.5 tokens/s	       148.8 total_ms	        25.44 ttft_ms
BenchmarkCompare-32    	      20	 149064338 ns/op	        23.00 prompt_tokens	       107.4 tokens/s	       149.0 total_ms	        25.51 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	14.950s
```

</details>
<!-- yzma:bench end compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":168.7,"ttft_ms":8.743,"total_ms":94.85,"prompt_tokens":22,"llamacpp":"b11146","engine_version":"1.28.0","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 168.7 tokens a second. 8.7 ms to the first token. 94.8 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=yzma -suite=text -tokens=16 -nctx=8192 -model=/home/ron/models/Qwen3-VL-2B-Instruct.Q4_K_M.gguf -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	  99552956 ns/op	        22.00 prompt_tokens	       167.1 tokens/s	        95.75 total_ms	         9.040 ttft_ms
BenchmarkCompare-32    	      20	  98588059 ns/op	        22.00 prompt_tokens	       168.8 tokens/s	        94.78 total_ms	         8.729 ttft_ms
BenchmarkCompare-32    	      20	  98495852 ns/op	        22.00 prompt_tokens	       169.0 tokens/s	        94.69 total_ms	         8.743 ttft_ms
BenchmarkCompare-32    	      20	  98711652 ns/op	        22.00 prompt_tokens	       168.6 tokens/s	        94.90 total_ms	         8.780 ttft_ms
BenchmarkCompare-32    	      20	  98653903 ns/op	        22.00 prompt_tokens	       168.7 tokens/s	        94.85 total_ms	         8.739 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.725s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":153.5,"ttft_ms":16.31,"total_ms":104.2,"prompt_tokens":22,"engine_version":"0.34.4","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 153.5 tokens a second. 16.3 ms to the first token. 104.2 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=ollama -suite=text -tokens=16 -nctx=8192 -server-model=yzma-bench-qwen3-vl-2b -ollama-url=http://localhost:11434/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 294998410 ns/op	        22.00 prompt_tokens	        54.24 tokens/s	       295.0 total_ms	       205.0 ttft_ms
BenchmarkCompare-32    	      20	 104440472 ns/op	        22.00 prompt_tokens	       153.2 tokens/s	       104.4 total_ms	        16.39 ttft_ms
BenchmarkCompare-32    	      20	 104240556 ns/op	        22.00 prompt_tokens	       153.5 tokens/s	       104.2 total_ms	        16.31 ttft_ms
BenchmarkCompare-32    	      20	 103672698 ns/op	        22.00 prompt_tokens	       154.4 tokens/s	       103.7 total_ms	        16.28 ttft_ms
BenchmarkCompare-32    	      20	 104239396 ns/op	        22.00 prompt_tokens	       153.5 tokens/s	       104.2 total_ms	        16.30 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	14.257s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":160.1,"ttft_ms":14.16,"total_ms":99.96,"prompt_tokens":22,"engine_version":"v1.2.8","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 160.1 tokens a second. 14.2 ms to the first token. 100.0 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=dmr -suite=text -tokens=16 -nctx=8192 -server-model=hf.co/qwen/qwen3-vl-2b-instruct-gguf:q4_k_m -dmr-url=http://localhost:12434/engines/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 105134221 ns/op	        22.00 prompt_tokens	       152.3 tokens/s	       105.1 total_ms	        17.70 ttft_ms
BenchmarkCompare-32    	      20	 100557542 ns/op	        22.00 prompt_tokens	       159.2 tokens/s	       100.5 total_ms	        14.71 ttft_ms
BenchmarkCompare-32    	      20	  99989887 ns/op	        22.00 prompt_tokens	       160.1 tokens/s	        99.96 total_ms	        14.16 ttft_ms
BenchmarkCompare-32    	      20	  99999749 ns/op	        22.00 prompt_tokens	       160.1 tokens/s	        99.96 total_ms	        14.13 ttft_ms
BenchmarkCompare-32    	      20	  99693409 ns/op	        22.00 prompt_tokens	       160.6 tokens/s	        99.66 total_ms	        13.90 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.137s
```

</details>
<!-- yzma:bench end compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

## Embeddings

An embedding returns no tokens, thus almost all of the cost of a request is
the round trip. The first token and the request contain the same number here,
because the answer is a vector and not a stream. Tokens per second counts the
tokens of the prompt that the engine read.

<!-- yzma:bench table compare-embeddings -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 20954.0 | 1.4 | 1.4 | 1.28.0 | 2026-09-24 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 6968.0 | 4.2 | 4.2 | 0.34.4 | 2026-09-24 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 6366.0 | 4.6 | 4.6 | v1.2.8 | 2026-09-24 |
<!-- yzma:bench table end compare-embeddings -->

<!-- yzma:bench start compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":20954,"ttft_ms":1.384,"total_ms":1.384,"prompt_tokens":29,"llamacpp":"b11146","engine_version":"1.28.0","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 20954.0 tokens a second. 1.4 ms to the first token. 1.4 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=yzma -suite=embeddings -tokens=16 -nctx=8192 -model=/home/ron/models/bge-small-en-v1.5-q8_0.gguf -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	   1433154 ns/op	        29.00 prompt_tokens	     20260 tokens/s	         1.431 total_ms	         1.431 ttft_ms
BenchmarkCompare-32    	      20	   1336158 ns/op	        29.00 prompt_tokens	     21733 tokens/s	         1.334 total_ms	         1.334 ttft_ms
BenchmarkCompare-32    	      20	   1386452 ns/op	        29.00 prompt_tokens	     20954 tokens/s	         1.384 total_ms	         1.384 ttft_ms
BenchmarkCompare-32    	      20	   1385680 ns/op	        29.00 prompt_tokens	     20983 tokens/s	         1.382 total_ms	         1.382 ttft_ms
BenchmarkCompare-32    	      20	   1412500 ns/op	        29.00 prompt_tokens	     20579 tokens/s	         1.409 total_ms	         1.409 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.589s
```

</details>
<!-- yzma:bench end compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":6968,"ttft_ms":4.162,"total_ms":4.162,"prompt_tokens":29,"engine_version":"0.34.4","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 6968.0 tokens a second. 4.2 ms to the first token. 4.2 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=ollama -suite=embeddings -tokens=16 -nctx=8192 -server-model=yzma-bench-bge-small -ollama-url=http://localhost:11434/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	   4196134 ns/op	        29.00 prompt_tokens	      6918 tokens/s	         4.192 total_ms	         4.192 ttft_ms
BenchmarkCompare-32    	      20	   4165955 ns/op	        29.00 prompt_tokens	      6968 tokens/s	         4.162 total_ms	         4.162 ttft_ms
BenchmarkCompare-32    	      20	   4217454 ns/op	        29.00 prompt_tokens	      6883 tokens/s	         4.213 total_ms	         4.213 ttft_ms
BenchmarkCompare-32    	      20	   3977332 ns/op	        29.00 prompt_tokens	      7299 tokens/s	         3.973 total_ms	         3.973 ttft_ms
BenchmarkCompare-32    	      20	   4027555 ns/op	        29.00 prompt_tokens	      7207 tokens/s	         4.024 total_ms	         4.024 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.426s
```

</details>
<!-- yzma:bench end compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":6366,"ttft_ms":4.555,"total_ms":4.555,"prompt_tokens":29,"engine_version":"v1.2.8","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 6366.0 tokens a second. 4.6 ms to the first token. 4.6 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Thu Sep 24 08:51:00 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              2W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=dmr -suite=embeddings -tokens=16 -nctx=8192 -server-model=hf.co/ggml-org/bge-small-en-v1.5-q8_0-gguf:q8_0 -dmr-url=http://localhost:12434/engines/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	   4154278 ns/op	        29.00 prompt_tokens	      6990 tokens/s	         4.149 total_ms	         4.149 ttft_ms
BenchmarkCompare-32    	      20	   4858284 ns/op	        29.00 prompt_tokens	      5977 tokens/s	         4.852 total_ms	         4.852 ttft_ms
BenchmarkCompare-32    	      20	   4851916 ns/op	        29.00 prompt_tokens	      5986 tokens/s	         4.845 total_ms	         4.845 ttft_ms
BenchmarkCompare-32    	      20	   4561526 ns/op	        29.00 prompt_tokens	      6366 tokens/s	         4.555 total_ms	         4.555 ttft_ms
BenchmarkCompare-32    	      20	   4519028 ns/op	        29.00 prompt_tokens	      6427 tokens/s	         4.512 total_ms	         4.512 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.477s
```

</details>
<!-- yzma:bench end compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->

## Images

Each request has one image and a short question about it. The image gives most
of the prompt tokens. Each run gets an image that no engine has seen, with the
same size and nearly the same pixels, thus a server cannot answer from its
cache. yzma decodes the image inside its own measurement, as a server does.

<!-- yzma:bench table compare-multimodal -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
<!-- yzma:bench table end compare-multimodal -->

## Each engine brings its own llama.cpp

The engines do not share one llama.cpp. yzma uses the build of its library
directory. Docker Model Runner pins a build in its image. ollama has a fork of
its own, with a version that does not map to a build of llama.cpp.

| Engine | llama.cpp |
| --- | --- |
| yzma | the build of `lib/`, b11146 (v0.5.0) of September 2026 here |
| Docker Model Runner | pinned in the image, b9879 of July 2026 in version 1.2 |
| ollama | a fork, version 0.4.1-dev here |

Thus the same GGUF file does not always give the same work. The count of the
prompt tokens shows it, and the script prints that count for each engine.

## What is not measured yet

The cold start is not here. Each table gives the numbers of an engine that has
the model in memory already. yzma also wins the cold start, because it has no
daemon to start, but that number needs a convention that these tables do not
have yet.
