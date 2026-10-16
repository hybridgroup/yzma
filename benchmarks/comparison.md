# Engine comparison benchmarks

These benchmarks compare inference performance across 3 engines that support GGUF models:

- yzma
- ollama
- Docker Model Runner

## Summary

| Suite | Result | Evidence |
| --- | --- | --- |
| Embeddings | yzma 2.6 to 2.7 times faster, 1.7 ms against 4.5 ms and 4.6 ms | Five runs, no overlap, each engine at 29 prompt tokens and a vector of 384 |
| Text | yzma 9.2 to 12.1 percent faster, 0.48 to 0.6 of the time to the first token | Five runs, no overlap, each engine at the same count of prompt tokens |
| Images | yzma 1.2 to 1.4 times faster for a request, 0.5 to 0.8 of the time to the first token | Five runs, no overlap, yzma and Docker Model Runner at the same count of prompt tokens, ollama at 5 more |

## Text

<!-- yzma:bench table compare-text -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 118.6 | 12.3 | 134.9 | 1.30.0 | 2026-10-16 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 105.8 | 24.6 | 151.2 | 0.40.5 | 2026-10-16 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 106.7 | 25.7 | 150.0 | v1.2.8 | 2026-10-16 |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 166.9 | 9.0 | 95.9 | 1.30.0 | 2026-10-16 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 152.8 | 15.1 | 104.7 | 0.40.5 | 2026-10-16 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 152.5 | 17.1 | 104.9 | v1.2.8 | 2026-10-16 |
<!-- yzma:bench table end compare-text -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":118.6,"ttft_ms":12.3,"total_ms":134.9,"prompt_tokens":23,"llamacpp":"b11669","engine_version":"1.30.0","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 118.6 tokens a second. 12.3 ms to the first token. 134.9 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 136109556 ns/op	        23.00 prompt_tokens	       118.1 tokens/s	       135.5 total_ms	        12.41 ttft_ms
BenchmarkCompare-32    	      20	 135544608 ns/op	        23.00 prompt_tokens	       118.6 tokens/s	       134.9 total_ms	        12.20 ttft_ms
BenchmarkCompare-32    	      20	 135567397 ns/op	        23.00 prompt_tokens	       118.6 tokens/s	       134.9 total_ms	        12.32 ttft_ms
BenchmarkCompare-32    	      20	 135344246 ns/op	        23.00 prompt_tokens	       118.8 tokens/s	       134.7 total_ms	        12.30 ttft_ms
BenchmarkCompare-32    	      20	 135182599 ns/op	        24.00 prompt_tokens	       118.9 tokens/s	       134.6 total_ms	        12.16 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.452s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":105.8,"ttft_ms":24.59,"total_ms":151.2,"prompt_tokens":23,"engine_version":"0.40.5","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 105.8 tokens a second. 24.6 ms to the first token. 151.2 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 153486881 ns/op	        23.00 prompt_tokens	       104.3 tokens/s	       153.5 total_ms	        25.27 ttft_ms
BenchmarkCompare-32    	      20	 152844415 ns/op	        23.00 prompt_tokens	       104.7 tokens/s	       152.8 total_ms	        24.76 ttft_ms
BenchmarkCompare-32    	      20	 151203022 ns/op	        23.00 prompt_tokens	       105.8 tokens/s	       151.2 total_ms	        24.52 ttft_ms
BenchmarkCompare-32    	      20	 150993552 ns/op	        23.00 prompt_tokens	       106.0 tokens/s	       151.0 total_ms	        24.50 ttft_ms
BenchmarkCompare-32    	      20	 151264380 ns/op	        24.00 prompt_tokens	       105.8 tokens/s	       151.2 total_ms	        24.59 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.236s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":106.7,"ttft_ms":25.74,"total_ms":150,"prompt_tokens":23,"engine_version":"v1.2.8","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 106.7 tokens a second. 25.7 ms to the first token. 150.0 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 149510038 ns/op	        23.00 prompt_tokens	       107.0 tokens/s	       149.5 total_ms	        25.74 ttft_ms
BenchmarkCompare-32    	      20	 149578367 ns/op	        23.00 prompt_tokens	       107.0 tokens/s	       149.5 total_ms	        25.69 ttft_ms
BenchmarkCompare-32    	      20	 152301924 ns/op	        23.00 prompt_tokens	       105.1 tokens/s	       152.3 total_ms	        27.96 ttft_ms
BenchmarkCompare-32    	      20	 150128168 ns/op	        23.00 prompt_tokens	       106.6 tokens/s	       150.1 total_ms	        25.80 ttft_ms
BenchmarkCompare-32    	      20	 150041236 ns/op	        24.00 prompt_tokens	       106.7 tokens/s	       150.0 total_ms	        25.64 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.073s
```

</details>
<!-- yzma:bench end compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":166.9,"ttft_ms":8.964,"total_ms":95.87,"prompt_tokens":22,"llamacpp":"b11669","engine_version":"1.30.0","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 166.9 tokens a second. 9.0 ms to the first token. 95.9 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	  99860357 ns/op	        22.00 prompt_tokens	       166.6 tokens/s	        96.05 total_ms	         9.093 ttft_ms
BenchmarkCompare-32    	      20	  99681575 ns/op	        22.00 prompt_tokens	       166.9 tokens/s	        95.87 total_ms	         8.964 ttft_ms
BenchmarkCompare-32    	      20	  99450817 ns/op	        22.00 prompt_tokens	       167.3 tokens/s	        95.65 total_ms	         8.916 ttft_ms
BenchmarkCompare-32    	      20	  99611383 ns/op	        22.00 prompt_tokens	       167.0 tokens/s	        95.80 total_ms	         8.937 ttft_ms
BenchmarkCompare-32    	      20	  99870632 ns/op	        23.00 prompt_tokens	       166.6 tokens/s	        96.05 total_ms	         8.976 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.770s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":152.8,"ttft_ms":15.09,"total_ms":104.7,"prompt_tokens":22,"engine_version":"0.40.5","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 152.8 tokens a second. 15.1 ms to the first token. 104.7 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 104607792 ns/op	        22.00 prompt_tokens	       153.0 tokens/s	       104.6 total_ms	        14.08 ttft_ms
BenchmarkCompare-32    	      20	 104093140 ns/op	        22.00 prompt_tokens	       153.8 tokens/s	       104.1 total_ms	        14.15 ttft_ms
BenchmarkCompare-32    	      20	 105543328 ns/op	        22.00 prompt_tokens	       151.6 tokens/s	       105.5 total_ms	        15.11 ttft_ms
BenchmarkCompare-32    	      20	 105860293 ns/op	        22.00 prompt_tokens	       151.2 tokens/s	       105.8 total_ms	        15.33 ttft_ms
BenchmarkCompare-32    	      20	 104704800 ns/op	        23.00 prompt_tokens	       152.8 tokens/s	       104.7 total_ms	        15.09 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.524s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":152.5,"ttft_ms":17.13,"total_ms":104.9,"prompt_tokens":22,"engine_version":"v1.2.8","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 152.5 tokens a second. 17.1 ms to the first token. 104.9 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 105683605 ns/op	        22.00 prompt_tokens	       151.5 tokens/s	       105.6 total_ms	        17.89 ttft_ms
BenchmarkCompare-32    	      20	 104462500 ns/op	        22.00 prompt_tokens	       153.2 tokens/s	       104.4 total_ms	        17.13 ttft_ms
BenchmarkCompare-32    	      20	 104956671 ns/op	        22.00 prompt_tokens	       152.5 tokens/s	       104.9 total_ms	        17.12 ttft_ms
BenchmarkCompare-32    	      20	 103985339 ns/op	        22.00 prompt_tokens	       153.9 tokens/s	       104.0 total_ms	        16.59 ttft_ms
BenchmarkCompare-32    	      20	 106025334 ns/op	        23.00 prompt_tokens	       151.0 tokens/s	       106.0 total_ms	        18.58 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.537s
```

</details>
<!-- yzma:bench end compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

## Embeddings

An embedding returns no tokens, so almost all of the request cost is the
round trip. First token and request time are the same here, because the answer
is a vector and not a stream. Tokens per second counts the prompt tokens the
engine reads.

<!-- yzma:bench table compare-embeddings -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 17123.0 | 1.7 | 1.7 | 1.30.0 | 2026-10-16 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 6393.0 | 4.5 | 4.5 | 0.40.5 | 2026-10-16 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 6360.0 | 4.6 | 4.6 | v1.2.8 | 2026-10-16 |
<!-- yzma:bench table end compare-embeddings -->

<!-- yzma:bench start compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":17123,"ttft_ms":1.694,"total_ms":1.694,"prompt_tokens":29,"llamacpp":"b11669","engine_version":"1.30.0","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 17123.0 tokens a second. 1.7 ms to the first token. 1.7 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	   1725231 ns/op	        29.00 prompt_tokens	     16830 tokens/s	         1.723 total_ms	         1.723 ttft_ms
BenchmarkCompare-32    	      20	   1593110 ns/op	        29.00 prompt_tokens	     18228 tokens/s	         1.591 total_ms	         1.591 ttft_ms
BenchmarkCompare-32    	      20	   1614654 ns/op	        29.00 prompt_tokens	     17991 tokens/s	         1.612 total_ms	         1.612 ttft_ms
BenchmarkCompare-32    	      20	   1703197 ns/op	        29.00 prompt_tokens	     17067 tokens/s	         1.699 total_ms	         1.699 ttft_ms
BenchmarkCompare-32    	      20	   1697607 ns/op	        29.00 prompt_tokens	     17123 tokens/s	         1.694 total_ms	         1.694 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.568s
```

</details>
<!-- yzma:bench end compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":6393,"ttft_ms":4.536,"total_ms":4.536,"prompt_tokens":29,"engine_version":"0.40.5","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 6393.0 tokens a second. 4.5 ms to the first token. 4.5 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	   5263934 ns/op	        29.00 prompt_tokens	      5515 tokens/s	         5.259 total_ms	         5.259 ttft_ms
BenchmarkCompare-32    	      20	   4540531 ns/op	        29.00 prompt_tokens	      6393 tokens/s	         4.536 total_ms	         4.536 ttft_ms
BenchmarkCompare-32    	      20	   4628742 ns/op	        29.00 prompt_tokens	      6272 tokens/s	         4.624 total_ms	         4.624 ttft_ms
BenchmarkCompare-32    	      20	   4295618 ns/op	        29.00 prompt_tokens	      6758 tokens/s	         4.291 total_ms	         4.291 ttft_ms
BenchmarkCompare-32    	      20	   4480446 ns/op	        29.00 prompt_tokens	      6479 tokens/s	         4.476 total_ms	         4.476 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.483s
```

</details>
<!-- yzma:bench end compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":6360,"ttft_ms":4.56,"total_ms":4.56,"prompt_tokens":29,"engine_version":"v1.2.8","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 6360.0 tokens a second. 4.6 ms to the first token. 4.6 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	   4565816 ns/op	        29.00 prompt_tokens	      6360 tokens/s	         4.560 total_ms	         4.560 ttft_ms
BenchmarkCompare-32    	      20	   4687089 ns/op	        29.00 prompt_tokens	      6195 tokens/s	         4.681 total_ms	         4.681 ttft_ms
BenchmarkCompare-32    	      20	   4614523 ns/op	        29.00 prompt_tokens	      6293 tokens/s	         4.609 total_ms	         4.609 ttft_ms
BenchmarkCompare-32    	      20	   3374329 ns/op	        29.00 prompt_tokens	      8605 tokens/s	         3.370 total_ms	         3.370 ttft_ms
BenchmarkCompare-32    	      20	   4064824 ns/op	        29.00 prompt_tokens	      7143 tokens/s	         4.060 total_ms	         4.060 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.441s
```

</details>
<!-- yzma:bench end compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->

## Images

Each request has one image and a short question about it. The image makes up
most of the prompt tokens. Each run gets an image that no engine has seen, with
the same size and nearly the same pixels, so a server can't answer from its
cache. yzma decodes the image inside its own measurement, as a server does.

Each engine scales images its own way, so each model gets an image at a size
that no engine changes, 1280x960 for qwen3-vl-2b and 768x576 for gemma4-e2b.
ollama counts 5 more prompt tokens per image than the others, on every model
and at every size, while the answers match.

<!-- yzma:bench table compare-multimodal -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 209 | 75.8 | 87.7 | 211.1 | 1.30.0 | 2026-10-16 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 214 | 61.3 | 133.3 | 260.9 | 0.40.5 | 2026-10-16 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 209 | 53.7 | 173.5 | 297.8 | v1.2.8 | 2026-10-16 |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 1216 | 30.0 | 437.5 | 534.1 | 1.30.0 | 2026-10-16 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 1221 | 23.9 | 570.7 | 669.6 | 0.40.5 | 2026-10-16 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 1216 | 25.4 | 533.6 | 630.0 | v1.2.8 | 2026-10-16 |
<!-- yzma:bench table end compare-multimodal -->

<!-- yzma:bench start compare-multimodal/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":75.81,"ttft_ms":87.74,"total_ms":211.1,"prompt_tokens":209,"llamacpp":"b11669","engine_version":"1.30.0","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 75.8 tokens a second. 87.7 ms to the first token. 211.1 ms for a request. 209 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=yzma -suite=multimodal -tokens=16 -nctx=8192 -image-size=768x576 -model=/home/ron/models/gemma-4-E2B-it-Q4_K_M.gguf -mmproj=/home/ron/models/mmproj-F16.gguf -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 233743224 ns/op	       209.0 prompt_tokens	        75.92 tokens/s	       210.8 total_ms	        86.26 ttft_ms
BenchmarkCompare-32    	      20	 231682910 ns/op	       209.0 prompt_tokens	        76.65 tokens/s	       208.8 total_ms	        85.95 ttft_ms
BenchmarkCompare-32    	      20	 234243382 ns/op	       209.0 prompt_tokens	        75.81 tokens/s	       211.1 total_ms	        87.74 ttft_ms
BenchmarkCompare-32    	      20	 237080419 ns/op	       209.0 prompt_tokens	        74.79 tokens/s	       213.9 total_ms	        89.95 ttft_ms
BenchmarkCompare-32    	      20	 237248327 ns/op	       209.0 prompt_tokens	        74.58 tokens/s	       214.5 total_ms	        90.79 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	25.836s
```

</details>
<!-- yzma:bench end compare-multimodal/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-multimodal/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":61.33,"ttft_ms":133.3,"total_ms":260.9,"prompt_tokens":214,"engine_version":"0.40.5","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 61.3 tokens a second. 133.3 ms to the first token. 260.9 ms for a request. 214 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=ollama -suite=multimodal -tokens=16 -nctx=8192 -image-size=768x576 -server-model=yzma-bench-gemma4-e2b -ollama-url=http://localhost:11434/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 283615732 ns/op	       214.0 prompt_tokens	        60.42 tokens/s	       264.8 total_ms	       136.2 ttft_ms
BenchmarkCompare-32    	      20	 276605656 ns/op	       214.0 prompt_tokens	        62.04 tokens/s	       257.9 total_ms	       131.5 ttft_ms
BenchmarkCompare-32    	      20	 277712781 ns/op	       214.0 prompt_tokens	        61.78 tokens/s	       259.0 total_ms	       132.3 ttft_ms
BenchmarkCompare-32    	      20	 279737218 ns/op	       214.0 prompt_tokens	        61.33 tokens/s	       260.9 total_ms	       133.3 ttft_ms
BenchmarkCompare-32    	      20	 282988305 ns/op	       214.0 prompt_tokens	        60.65 tokens/s	       263.8 total_ms	       135.9 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	28.172s
```

</details>
<!-- yzma:bench end compare-multimodal/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-multimodal/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":53.73,"ttft_ms":173.5,"total_ms":297.8,"prompt_tokens":209,"engine_version":"v1.2.8","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 53.7 tokens a second. 173.5 ms to the first token. 297.8 ms for a request. 209 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=dmr -suite=multimodal -tokens=16 -nctx=8192 -image-size=768x576 -server-model=hf.co/unsloth/gemma-4-e2b-it-gguf:q4_k_m -dmr-url=http://localhost:12434/engines/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 321697096 ns/op	       209.0 prompt_tokens	        52.33 tokens/s	       305.7 total_ms	       181.0 ttft_ms
BenchmarkCompare-32    	      20	 317255603 ns/op	       209.0 prompt_tokens	        53.15 tokens/s	       301.1 total_ms	       176.1 ttft_ms
BenchmarkCompare-32    	      20	 312184390 ns/op	       209.0 prompt_tokens	        54.00 tokens/s	       296.3 total_ms	       171.7 ttft_ms
BenchmarkCompare-32    	      20	 313295786 ns/op	       209.0 prompt_tokens	        53.94 tokens/s	       296.6 total_ms	       172.2 ttft_ms
BenchmarkCompare-32    	      20	 313919391 ns/op	       209.0 prompt_tokens	        53.73 tokens/s	       297.8 total_ms	       173.5 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	31.833s
```

</details>
<!-- yzma:bench end compare-multimodal/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-multimodal/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":29.96,"ttft_ms":437.5,"total_ms":534.1,"prompt_tokens":1216,"llamacpp":"b11669","engine_version":"1.30.0","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 30.0 tokens a second. 437.5 ms to the first token. 534.1 ms for a request. 1216 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=yzma -suite=multimodal -tokens=16 -nctx=8192 -image-size=1280x960 -model=/home/ron/models/Qwen3-VL-2B-Instruct.Q4_K_M.gguf -mmproj=/home/ron/models/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 569337460 ns/op	      1216 prompt_tokens	        30.60 tokens/s	       522.9 total_ms	       427.1 ttft_ms
BenchmarkCompare-32    	      20	 594721491 ns/op	      1216 prompt_tokens	        29.18 tokens/s	       548.4 total_ms	       451.7 ttft_ms
BenchmarkCompare-32    	      20	 575438039 ns/op	      1216 prompt_tokens	        30.25 tokens/s	       528.9 total_ms	       432.4 ttft_ms
BenchmarkCompare-32    	      20	 580441120 ns/op	      1216 prompt_tokens	        29.96 tokens/s	       534.1 total_ms	       437.5 ttft_ms
BenchmarkCompare-32    	      20	 597960021 ns/op	      1216 prompt_tokens	        29.01 tokens/s	       551.6 total_ms	       454.7 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	59.985s
```

</details>
<!-- yzma:bench end compare-multimodal/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-multimodal/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":23.89,"ttft_ms":570.7,"total_ms":669.6,"prompt_tokens":1221,"engine_version":"0.40.5","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 23.9 tokens a second. 570.7 ms to the first token. 669.6 ms for a request. 1221 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=ollama -suite=multimodal -tokens=16 -nctx=8192 -image-size=1280x960 -server-model=yzma-bench-qwen3-vl-2b -ollama-url=http://localhost:11434/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 691431614 ns/op	      1221 prompt_tokens	        24.58 tokens/s	       650.9 total_ms	       551.4 ttft_ms
BenchmarkCompare-32    	      20	 710315813 ns/op	      1221 prompt_tokens	        23.89 tokens/s	       669.6 total_ms	       570.7 ttft_ms
BenchmarkCompare-32    	      20	 707129374 ns/op	      1221 prompt_tokens	        23.98 tokens/s	       667.2 total_ms	       567.7 ttft_ms
BenchmarkCompare-32    	      20	 717456218 ns/op	      1221 prompt_tokens	        23.62 tokens/s	       677.4 total_ms	       578.1 ttft_ms
BenchmarkCompare-32    	      20	 716560904 ns/op	      1221 prompt_tokens	        23.63 tokens/s	       677.0 total_ms	       577.3 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	71.238s
```

</details>
<!-- yzma:bench end compare-multimodal/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-multimodal/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":25.4,"ttft_ms":533.6,"total_ms":630,"prompt_tokens":1216,"engine_version":"v1.2.8","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 25.4 tokens a second. 533.6 ms to the first token. 630.0 ms for a request. 1216 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Oct 16 14:49:11 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   56C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=5 -run=nada -bench BenchmarkCompare -engine=dmr -suite=multimodal -tokens=16 -nctx=8192 -image-size=1280x960 -server-model=hf.co/qwen/qwen3-vl-2b-instruct-gguf:q4_k_m -dmr-url=http://localhost:12434/engines/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 666118364 ns/op	      1216 prompt_tokens	        25.58 tokens/s	       625.5 total_ms	       529.1 ttft_ms
BenchmarkCompare-32    	      20	 667992194 ns/op	      1216 prompt_tokens	        25.40 tokens/s	       630.0 total_ms	       533.6 ttft_ms
BenchmarkCompare-32    	      20	 667191009 ns/op	      1216 prompt_tokens	        25.45 tokens/s	       628.7 total_ms	       531.8 ttft_ms
BenchmarkCompare-32    	      20	 675189270 ns/op	      1216 prompt_tokens	        25.10 tokens/s	       637.5 total_ms	       540.4 ttft_ms
BenchmarkCompare-32    	      20	 672890172 ns/op	      1216 prompt_tokens	        25.21 tokens/s	       634.6 total_ms	       537.8 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	67.321s
```

</details>
<!-- yzma:bench end compare-multimodal/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

## Each engine ships its own llama.cpp

The engines don't share one llama.cpp. yzma uses the build in its library
directory. Docker Model Runner pins a build in its image. ollama has its own
fork, with a version that doesn't map to a llama.cpp build.

| Engine | llama.cpp |
| --- | --- |
| yzma | the build of `lib/`, b11669 (v0.7.0) of October 2026 here |
| Docker Model Runner | pinned in the image, b9879 of July 2026 in version 1.2 |
| ollama | a fork, with no llama.cpp build number |

So the same GGUF file doesn't always mean the same work. The prompt token count
shows it, and the script prints that count for each engine.

## What is not measured yet

Cold start is not measured. Each table shows an engine that already has the
model in memory. yzma also wins on cold start, because it has no daemon to
start, but measuring it needs a convention these tables don't have yet.
