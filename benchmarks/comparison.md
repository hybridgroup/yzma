# Engine comparison benchmarks

These benchmarks are to compare inference performance using 3 different engines that support GGUF models:

- yzma
- ollama
- Docker Model Runner

## Summary

| Suite | Result | Evidence |
| --- | --- | --- |
| Embeddings | yzma 3.1 to 3.5 times faster, 1.4 ms against 4.2 ms and 4.9 ms | Five runs, no overlap, each engine at 29 prompt tokens and a vector of 384 |
| Text | yzma 5.4 to 10.7 percent faster, 0.5 to 0.6 of the time to the first token | Five runs, no overlap, each engine at the same count of prompt tokens |
| Images | No numbers. The code runs with `--suite multimodal` | The engines preprocess an image in different ways |

## Text

<!-- yzma:bench table compare-text -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 118.4 | 12.3 | 135.1 | 1.27.0 | 2026-09-23 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 107.0 | 23.6 | 149.6 | 0.34.3 | 2026-09-23 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 107.2 | 25.8 | 149.2 | v1.2.8 | 2026-09-23 |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 168.6 | 8.7 | 94.9 | 1.27.0 | 2026-09-23 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 154.5 | 16.3 | 103.6 | 0.34.3 | 2026-09-23 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 160.0 | 14.1 | 100.0 | v1.2.8 | 2026-09-23 |
<!-- yzma:bench table end compare-text -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":118.4,"ttft_ms":12.32,"total_ms":135.1,"prompt_tokens":23,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 118.4 tokens a second. 12.3 ms to the first token. 135.1 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	 139866453 ns/op	        23.00 prompt_tokens	       114.9 tokens/s	       139.2 total_ms	        12.75 ttft_ms
BenchmarkCompare-32    	      20	 135867741 ns/op	        23.00 prompt_tokens	       118.3 tokens/s	       135.2 total_ms	        12.32 ttft_ms
BenchmarkCompare-32    	      20	 135170985 ns/op	        23.00 prompt_tokens	       118.9 tokens/s	       134.5 total_ms	        12.27 ttft_ms
BenchmarkCompare-32    	      20	 135463893 ns/op	        23.00 prompt_tokens	       118.7 tokens/s	       134.8 total_ms	        12.30 ttft_ms
BenchmarkCompare-32    	      20	 135731652 ns/op	        23.00 prompt_tokens	       118.4 tokens/s	       135.1 total_ms	        12.34 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.542s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":107,"ttft_ms":23.58,"total_ms":149.6,"prompt_tokens":23,"engine_version":"0.34.3","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 107.0 tokens a second. 23.6 ms to the first token. 149.6 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	 152103528 ns/op	        23.00 prompt_tokens	       105.2 tokens/s	       152.1 total_ms	        25.41 ttft_ms
BenchmarkCompare-32    	      20	 149959190 ns/op	        23.00 prompt_tokens	       106.7 tokens/s	       149.9 total_ms	        23.58 ttft_ms
BenchmarkCompare-32    	      20	 149600319 ns/op	        23.00 prompt_tokens	       107.0 tokens/s	       149.6 total_ms	        23.56 ttft_ms
BenchmarkCompare-32    	      20	 149482832 ns/op	        23.00 prompt_tokens	       107.0 tokens/s	       149.5 total_ms	        23.60 ttft_ms
BenchmarkCompare-32    	      20	 148777775 ns/op	        23.00 prompt_tokens	       107.6 tokens/s	       148.8 total_ms	        23.35 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.037s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":107.2,"ttft_ms":25.81,"total_ms":149.2,"prompt_tokens":23,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 107.2 tokens a second. 25.8 ms to the first token. 149.2 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	 148991472 ns/op	        23.00 prompt_tokens	       107.4 tokens/s	       148.9 total_ms	        25.81 ttft_ms
BenchmarkCompare-32    	      20	 148858857 ns/op	        23.00 prompt_tokens	       107.5 tokens/s	       148.8 total_ms	        25.77 ttft_ms
BenchmarkCompare-32    	      20	 149289292 ns/op	        23.00 prompt_tokens	       107.2 tokens/s	       149.2 total_ms	        25.96 ttft_ms
BenchmarkCompare-32    	      20	 149285849 ns/op	        23.00 prompt_tokens	       107.2 tokens/s	       149.2 total_ms	        25.70 ttft_ms
BenchmarkCompare-32    	      20	 150273626 ns/op	        23.00 prompt_tokens	       106.5 tokens/s	       150.2 total_ms	        25.90 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	14.975s
```

</details>
<!-- yzma:bench end compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":168.6,"ttft_ms":8.705,"total_ms":94.89,"prompt_tokens":22,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 168.6 tokens a second. 8.7 ms to the first token. 94.9 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	  99134461 ns/op	        22.00 prompt_tokens	       167.9 tokens/s	        95.32 total_ms	         8.907 ttft_ms
BenchmarkCompare-32    	      20	  98559191 ns/op	        22.00 prompt_tokens	       168.9 tokens/s	        94.75 total_ms	         8.705 ttft_ms
BenchmarkCompare-32    	      20	  98579699 ns/op	        22.00 prompt_tokens	       168.8 tokens/s	        94.77 total_ms	         8.688 ttft_ms
BenchmarkCompare-32    	      20	  98708280 ns/op	        22.00 prompt_tokens	       168.6 tokens/s	        94.89 total_ms	         8.735 ttft_ms
BenchmarkCompare-32    	      20	  98714485 ns/op	        22.00 prompt_tokens	       168.6 tokens/s	        94.90 total_ms	         8.684 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.722s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":154.5,"ttft_ms":16.29,"total_ms":103.6,"prompt_tokens":22,"engine_version":"0.34.3","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 154.5 tokens a second. 16.3 ms to the first token. 103.6 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	 316580736 ns/op	        22.00 prompt_tokens	        50.54 tokens/s	       316.6 total_ms	       227.0 ttft_ms
BenchmarkCompare-32    	      20	 103421565 ns/op	        22.00 prompt_tokens	       154.7 tokens/s	       103.4 total_ms	        16.29 ttft_ms
BenchmarkCompare-32    	      20	 103596516 ns/op	        22.00 prompt_tokens	       154.5 tokens/s	       103.6 total_ms	        16.16 ttft_ms
BenchmarkCompare-32    	      20	 103911467 ns/op	        22.00 prompt_tokens	       154.0 tokens/s	       103.9 total_ms	        16.25 ttft_ms
BenchmarkCompare-32    	      20	 103416782 ns/op	        22.00 prompt_tokens	       154.7 tokens/s	       103.4 total_ms	        16.38 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	14.644s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":160,"ttft_ms":14.14,"total_ms":99.99,"prompt_tokens":22,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 160.0 tokens a second. 14.1 ms to the first token. 100.0 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	 105521778 ns/op	        22.00 prompt_tokens	       151.7 tokens/s	       105.5 total_ms	        18.15 ttft_ms
BenchmarkCompare-32    	      20	 100451949 ns/op	        22.00 prompt_tokens	       159.3 tokens/s	       100.4 total_ms	        14.48 ttft_ms
BenchmarkCompare-32    	      20	  99703999 ns/op	        22.00 prompt_tokens	       160.5 tokens/s	        99.68 total_ms	        14.14 ttft_ms
BenchmarkCompare-32    	      20	  99885274 ns/op	        22.00 prompt_tokens	       160.2 tokens/s	        99.85 total_ms	        13.92 ttft_ms
BenchmarkCompare-32    	      20	 100028430 ns/op	        22.00 prompt_tokens	       160.0 tokens/s	        99.99 total_ms	        13.79 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.141s
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
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 20951.0 | 1.4 | 1.4 | 1.27.0 | 2026-09-23 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 6847.0 | 4.2 | 4.2 | 0.34.3 | 2026-09-23 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 5926.0 | 4.9 | 4.9 | v1.2.8 | 2026-09-23 |
<!-- yzma:bench table end compare-embeddings -->

<!-- yzma:bench start compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":20951,"ttft_ms":1.384,"total_ms":1.384,"prompt_tokens":29,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 20951.0 tokens a second. 1.4 ms to the first token. 1.4 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	   1446580 ns/op	        29.00 prompt_tokens	     20081 tokens/s	         1.444 total_ms	         1.444 ttft_ms
BenchmarkCompare-32    	      20	   1479256 ns/op	        29.00 prompt_tokens	     19651 tokens/s	         1.476 total_ms	         1.476 ttft_ms
BenchmarkCompare-32    	      20	   1386078 ns/op	        29.00 prompt_tokens	     20957 tokens/s	         1.384 total_ms	         1.384 ttft_ms
BenchmarkCompare-32    	      20	   1355019 ns/op	        29.00 prompt_tokens	     21451 tokens/s	         1.352 total_ms	         1.352 ttft_ms
BenchmarkCompare-32    	      20	   1386545 ns/op	        29.00 prompt_tokens	     20951 tokens/s	         1.384 total_ms	         1.384 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.598s
```

</details>
<!-- yzma:bench end compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":6847,"ttft_ms":4.235,"total_ms":4.235,"prompt_tokens":29,"engine_version":"0.34.3","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 6847.0 tokens a second. 4.2 ms to the first token. 4.2 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	   4236347 ns/op	        29.00 prompt_tokens	      6852 tokens/s	         4.232 total_ms	         4.232 ttft_ms
BenchmarkCompare-32    	      20	   4400710 ns/op	        29.00 prompt_tokens	      6596 tokens/s	         4.396 total_ms	         4.396 ttft_ms
BenchmarkCompare-32    	      20	   4128943 ns/op	        29.00 prompt_tokens	      7029 tokens/s	         4.126 total_ms	         4.126 ttft_ms
BenchmarkCompare-32    	      20	   4239413 ns/op	        29.00 prompt_tokens	      6847 tokens/s	         4.235 total_ms	         4.235 ttft_ms
BenchmarkCompare-32    	      20	   4456392 ns/op	        29.00 prompt_tokens	      6513 tokens/s	         4.452 total_ms	         4.452 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.445s
```

</details>
<!-- yzma:bench end compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":5926,"ttft_ms":4.894,"total_ms":4.894,"prompt_tokens":29,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-23"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 5926.0 tokens a second. 4.9 ms to the first token. 4.9 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Wed Sep 23 16:48:06 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   53C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
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
BenchmarkCompare-32    	      20	   4899155 ns/op	        29.00 prompt_tokens	      5926 tokens/s	         4.894 total_ms	         4.894 ttft_ms
BenchmarkCompare-32    	      20	   4373303 ns/op	        29.00 prompt_tokens	      6640 tokens/s	         4.367 total_ms	         4.367 ttft_ms
BenchmarkCompare-32    	      20	   3354847 ns/op	        29.00 prompt_tokens	      8654 tokens/s	         3.351 total_ms	         3.351 ttft_ms
BenchmarkCompare-32    	      20	   4933515 ns/op	        29.00 prompt_tokens	      5885 tokens/s	         4.928 total_ms	         4.928 ttft_ms
BenchmarkCompare-32    	      20	   5322063 ns/op	        29.00 prompt_tokens	      5456 tokens/s	         5.315 total_ms	         5.315 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.473s
```

</details>
<!-- yzma:bench end compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->

## Images, work in progress

This report provides no numbers for images yet. WIP code of the suite is here and
it runs with `--suite multimodal`, but a measurement that compares the engines
fairly is not yet ready.

## Each engine brings its own llama.cpp

The engines do not share one llama.cpp. yzma uses the build of its library
directory. Docker Model Runner pins a build in its image. ollama has a fork of
its own, with a version that does not map to a build of llama.cpp.

| Engine | llama.cpp |
| --- | --- |
| yzma | the build of `lib/`, b10964 of September 2026 here |
| Docker Model Runner | pinned in the image, b9879 of July 2026 in version 1.2 |
| ollama | a fork, version 0.4.1-dev here |

Thus the same GGUF file does not always give the same work. The count of the
prompt tokens shows it, and the script prints that count for each engine.

## What is not measured yet

The cold start is not here. Each table gives the numbers of an engine that has
the model in memory already. yzma also wins the cold start, because it has no
daemon to start, but that number needs a convention that these tables do not
have yet.
