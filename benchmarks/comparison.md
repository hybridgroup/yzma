# Engine comparison benchmarks

These benchmarks are to compare inference performance using 3 different engines that support GGUF models:

- yzma
- ollama
- Docker Model Runner

## Summary

| Suite | Result | Evidence |
| --- | --- | --- |
| Embeddings | yzma 2.6 to 3.1 times faster, 1.4 ms against 4.3 ms and 3.7 ms | Five runs, no overlap, each engine at 29 prompt tokens and a vector of 384 |
| Text | yzma 10.5 to 14.7 percent faster, 0.45 to 0.5 of the time to the first token | Five runs, no overlap, each engine at the same count of prompt tokens |
| Images | yzma 1.2 to 1.4 times faster for a request, 0.5 to 0.8 of the time to the first token | Five runs, no overlap, yzma and Docker Model Runner at the same count of prompt tokens, ollama at 5 more |

## Text

<!-- yzma:bench table compare-text -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 118.5 | 12.3 | 135.0 | 1.29.0-dev | 2026-09-25 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 104.7 | 26.1 | 152.8 | 0.34.4 | 2026-09-25 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 106.7 | 26.0 | 150.0 | v1.2.8 | 2026-09-25 |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 167.6 | 8.9 | 95.5 | 1.29.0-dev | 2026-09-25 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 146.1 | 19.6 | 109.5 | 0.34.4 | 2026-09-25 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 151.7 | 17.7 | 105.5 | v1.2.8 | 2026-09-25 |
<!-- yzma:bench table end compare-text -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":118.5,"ttft_ms":12.29,"total_ms":135,"prompt_tokens":23,"llamacpp":"b11179","engine_version":"1.29.0-dev","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 118.5 tokens a second. 12.3 ms to the first token. 135.0 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 139810811 ns/op	        23.00 prompt_tokens	       115.0 tokens/s	       139.1 total_ms	        12.71 ttft_ms
BenchmarkCompare-32    	      20	 135634360 ns/op	        23.00 prompt_tokens	       118.5 tokens/s	       135.0 total_ms	        12.34 ttft_ms
BenchmarkCompare-32    	      20	 135660450 ns/op	        23.00 prompt_tokens	       118.5 tokens/s	       135.0 total_ms	        12.26 ttft_ms
BenchmarkCompare-32    	      20	 135566740 ns/op	        23.00 prompt_tokens	       118.6 tokens/s	       134.9 total_ms	        12.27 ttft_ms
BenchmarkCompare-32    	      20	 135584263 ns/op	        24.00 prompt_tokens	       118.6 tokens/s	       135.0 total_ms	        12.29 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.682s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":104.7,"ttft_ms":26.1,"total_ms":152.8,"prompt_tokens":23,"engine_version":"0.34.4","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 104.7 tokens a second. 26.1 ms to the first token. 152.8 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 153237147 ns/op	        23.00 prompt_tokens	       104.4 tokens/s	       153.2 total_ms	        26.20 ttft_ms
BenchmarkCompare-32    	      20	 153094785 ns/op	        23.00 prompt_tokens	       104.5 tokens/s	       153.1 total_ms	        26.00 ttft_ms
BenchmarkCompare-32    	      20	 152801651 ns/op	        23.00 prompt_tokens	       104.7 tokens/s	       152.8 total_ms	        26.27 ttft_ms
BenchmarkCompare-32    	      20	 152661019 ns/op	        23.00 prompt_tokens	       104.8 tokens/s	       152.6 total_ms	        25.84 ttft_ms
BenchmarkCompare-32    	      20	 152784116 ns/op	        24.00 prompt_tokens	       104.7 tokens/s	       152.8 total_ms	        26.10 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.328s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":106.7,"ttft_ms":26,"total_ms":150,"prompt_tokens":23,"engine_version":"v1.2.8","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 106.7 tokens a second. 26.0 ms to the first token. 150.0 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 150049430 ns/op	        23.00 prompt_tokens	       106.7 tokens/s	       150.0 total_ms	        26.12 ttft_ms
BenchmarkCompare-32    	      20	 149727690 ns/op	        23.00 prompt_tokens	       106.9 tokens/s	       149.7 total_ms	        25.85 ttft_ms
BenchmarkCompare-32    	      20	 150050261 ns/op	        23.00 prompt_tokens	       106.7 tokens/s	       150.0 total_ms	        25.98 ttft_ms
BenchmarkCompare-32    	      20	 150216936 ns/op	        23.00 prompt_tokens	       106.5 tokens/s	       150.2 total_ms	        26.00 ttft_ms
BenchmarkCompare-32    	      20	 150361056 ns/op	        24.00 prompt_tokens	       106.4 tokens/s	       150.3 total_ms	        26.12 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	15.046s
```

</details>
<!-- yzma:bench end compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":167.6,"ttft_ms":8.938,"total_ms":95.46,"prompt_tokens":22,"llamacpp":"b11179","engine_version":"1.29.0-dev","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 167.6 tokens a second. 8.9 ms to the first token. 95.5 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	  99726605 ns/op	        22.00 prompt_tokens	       166.8 tokens/s	        95.90 total_ms	         9.133 ttft_ms
BenchmarkCompare-32    	      20	  99245455 ns/op	        22.00 prompt_tokens	       167.7 tokens/s	        95.42 total_ms	         8.925 ttft_ms
BenchmarkCompare-32    	      20	  99168447 ns/op	        22.00 prompt_tokens	       167.8 tokens/s	        95.35 total_ms	         8.938 ttft_ms
BenchmarkCompare-32    	      20	  99288677 ns/op	        22.00 prompt_tokens	       167.6 tokens/s	        95.48 total_ms	         8.938 ttft_ms
BenchmarkCompare-32    	      20	  99275950 ns/op	        23.00 prompt_tokens	       167.6 tokens/s	        95.46 total_ms	         8.905 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.777s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":146.1,"ttft_ms":19.63,"total_ms":109.5,"prompt_tokens":22,"engine_version":"0.34.4","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 146.1 tokens a second. 19.6 ms to the first token. 109.5 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 109540730 ns/op	        22.00 prompt_tokens	       146.1 tokens/s	       109.5 total_ms	        19.83 ttft_ms
BenchmarkCompare-32    	      20	 108231342 ns/op	        22.00 prompt_tokens	       147.9 tokens/s	       108.2 total_ms	        19.32 ttft_ms
BenchmarkCompare-32    	      20	 109331570 ns/op	        22.00 prompt_tokens	       146.4 tokens/s	       109.3 total_ms	        19.52 ttft_ms
BenchmarkCompare-32    	      20	 109811138 ns/op	        22.00 prompt_tokens	       145.7 tokens/s	       109.8 total_ms	        19.63 ttft_ms
BenchmarkCompare-32    	      20	 110103446 ns/op	        23.00 prompt_tokens	       145.3 tokens/s	       110.1 total_ms	        19.97 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.981s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":151.7,"ttft_ms":17.71,"total_ms":105.5,"prompt_tokens":22,"engine_version":"v1.2.8","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 151.7 tokens a second. 17.7 ms to the first token. 105.5 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 105508979 ns/op	        22.00 prompt_tokens	       151.7 tokens/s	       105.5 total_ms	        17.79 ttft_ms
BenchmarkCompare-32    	      20	 105593789 ns/op	        22.00 prompt_tokens	       151.6 tokens/s	       105.6 total_ms	        17.71 ttft_ms
BenchmarkCompare-32    	      20	 105620910 ns/op	        22.00 prompt_tokens	       151.6 tokens/s	       105.6 total_ms	        17.79 ttft_ms
BenchmarkCompare-32    	      20	 105255991 ns/op	        22.00 prompt_tokens	       152.1 tokens/s	       105.2 total_ms	        17.34 ttft_ms
BenchmarkCompare-32    	      20	 105380190 ns/op	        23.00 prompt_tokens	       151.9 tokens/s	       105.3 total_ms	        17.34 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	10.582s
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
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 20094.0 | 1.4 | 1.4 | 1.29.0-dev | 2026-09-25 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 6742.0 | 4.3 | 4.3 | 0.34.4 | 2026-09-25 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 7784.0 | 3.7 | 3.7 | v1.2.8 | 2026-09-25 |
<!-- yzma:bench table end compare-embeddings -->

<!-- yzma:bench start compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":20094,"ttft_ms":1.443,"total_ms":1.443,"prompt_tokens":29,"llamacpp":"b11179","engine_version":"1.29.0-dev","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 20094.0 tokens a second. 1.4 ms to the first token. 1.4 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	   1476281 ns/op	        29.00 prompt_tokens	     19668 tokens/s	         1.474 total_ms	         1.474 ttft_ms
BenchmarkCompare-32    	      20	   1362974 ns/op	        29.00 prompt_tokens	     21306 tokens/s	         1.361 total_ms	         1.361 ttft_ms
BenchmarkCompare-32    	      20	   1495376 ns/op	        29.00 prompt_tokens	     19457 tokens/s	         1.490 total_ms	         1.490 ttft_ms
BenchmarkCompare-32    	      20	   1447153 ns/op	        29.00 prompt_tokens	     20094 tokens/s	         1.443 total_ms	         1.443 ttft_ms
BenchmarkCompare-32    	      20	   1347771 ns/op	        29.00 prompt_tokens	     21545 tokens/s	         1.346 total_ms	         1.346 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.599s
```

</details>
<!-- yzma:bench end compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":6742,"ttft_ms":4.301,"total_ms":4.301,"prompt_tokens":29,"engine_version":"0.34.4","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 6742.0 tokens a second. 4.3 ms to the first token. 4.3 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	   5473549 ns/op	        29.00 prompt_tokens	      5304 tokens/s	         5.467 total_ms	         5.467 ttft_ms
BenchmarkCompare-32    	      20	   4305789 ns/op	        29.00 prompt_tokens	      6742 tokens/s	         4.301 total_ms	         4.301 ttft_ms
BenchmarkCompare-32    	      20	   4225506 ns/op	        29.00 prompt_tokens	      6870 tokens/s	         4.221 total_ms	         4.221 ttft_ms
BenchmarkCompare-32    	      20	   4219690 ns/op	        29.00 prompt_tokens	      6879 tokens/s	         4.216 total_ms	         4.216 ttft_ms
BenchmarkCompare-32    	      20	   4339701 ns/op	        29.00 prompt_tokens	      6689 tokens/s	         4.335 total_ms	         4.335 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.469s
```

</details>
<!-- yzma:bench end compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":7784,"ttft_ms":3.726,"total_ms":3.726,"prompt_tokens":29,"engine_version":"v1.2.8","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 7784.0 tokens a second. 3.7 ms to the first token. 3.7 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	   3504334 ns/op	        29.00 prompt_tokens	      8285 tokens/s	         3.500 total_ms	         3.500 ttft_ms
BenchmarkCompare-32    	      20	   3730842 ns/op	        29.00 prompt_tokens	      7784 tokens/s	         3.726 total_ms	         3.726 ttft_ms
BenchmarkCompare-32    	      20	   5507510 ns/op	        29.00 prompt_tokens	      5274 tokens/s	         5.499 total_ms	         5.499 ttft_ms
BenchmarkCompare-32    	      20	   4948398 ns/op	        29.00 prompt_tokens	      5869 tokens/s	         4.941 total_ms	         4.941 ttft_ms
BenchmarkCompare-32    	      20	   3114935 ns/op	        29.00 prompt_tokens	      9322 tokens/s	         3.111 total_ms	         3.111 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	0.433s
```

</details>
<!-- yzma:bench end compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->

## Images

Each request has one image and a short question about it. The image gives most
of the prompt tokens. Each run gets an image that no engine has seen, with the
same size and nearly the same pixels, thus a server cannot answer from its
cache. yzma decodes the image inside its own measurement, as a server does.

Each engine scales an image in its own way, thus each model gets an image at
a size that no engine changes, 1280x960 for qwen3-vl-2b and 768x576 for
gemma4-e2b. ollama counts 5 prompt tokens more than the others for an image,
on each model and at each size, while the answers agree.

<!-- yzma:bench table compare-multimodal -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 209 | 73.8 | 91.6 | 216.9 | 1.29.0-dev | 2026-09-25 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 214 | 60.4 | 136.4 | 264.9 | 0.34.4 | 2026-09-25 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 209 | 52.8 | 177.1 | 303.0 | v1.2.8 | 2026-09-25 |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 1216 | 29.8 | 441.2 | 537.6 | 1.29.0-dev | 2026-09-25 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 1221 | 23.7 | 574.2 | 674.0 | 0.34.4 | 2026-09-25 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 1216 | 24.7 | 550.4 | 647.4 | v1.2.8 | 2026-09-25 |
<!-- yzma:bench table end compare-multimodal -->

<!-- yzma:bench start compare-multimodal/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":73.76,"ttft_ms":91.6,"total_ms":216.9,"prompt_tokens":209,"llamacpp":"b11179","engine_version":"1.29.0-dev","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 73.8 tokens a second. 91.6 ms to the first token. 216.9 ms for a request. 209 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 234091691 ns/op	       209.0 prompt_tokens	        75.86 tokens/s	       210.9 total_ms	        87.42 ttft_ms
BenchmarkCompare-32    	      20	 238746642 ns/op	       209.0 prompt_tokens	        74.30 tokens/s	       215.3 total_ms	        91.45 ttft_ms
BenchmarkCompare-32    	      20	 239990194 ns/op	       209.0 prompt_tokens	        73.76 tokens/s	       216.9 total_ms	        92.66 ttft_ms
BenchmarkCompare-32    	      20	 239825230 ns/op	       209.0 prompt_tokens	        73.68 tokens/s	       217.2 total_ms	        92.26 ttft_ms
BenchmarkCompare-32    	      20	 239981426 ns/op	       209.0 prompt_tokens	        73.69 tokens/s	       217.1 total_ms	        91.60 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	26.387s
```

</details>
<!-- yzma:bench end compare-multimodal/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-multimodal/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":60.39,"ttft_ms":136.4,"total_ms":264.9,"prompt_tokens":214,"engine_version":"0.34.4","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 60.4 tokens a second. 136.4 ms to the first token. 264.9 ms for a request. 214 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 283419924 ns/op	       214.0 prompt_tokens	        60.16 tokens/s	       266.0 total_ms	       138.0 ttft_ms
BenchmarkCompare-32    	      20	 279318622 ns/op	       214.0 prompt_tokens	        61.13 tokens/s	       261.7 total_ms	       135.4 ttft_ms
BenchmarkCompare-32    	      20	 280015579 ns/op	       214.0 prompt_tokens	        60.98 tokens/s	       262.4 total_ms	       135.6 ttft_ms
BenchmarkCompare-32    	      20	 283450618 ns/op	       214.0 prompt_tokens	        60.28 tokens/s	       265.4 total_ms	       136.7 ttft_ms
BenchmarkCompare-32    	      20	 281932031 ns/op	       214.0 prompt_tokens	        60.39 tokens/s	       264.9 total_ms	       136.4 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	28.325s
```

</details>
<!-- yzma:bench end compare-multimodal/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-multimodal/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":52.8,"ttft_ms":177.1,"total_ms":303,"prompt_tokens":209,"engine_version":"v1.2.8","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 52.8 tokens a second. 177.1 ms to the first token. 303.0 ms for a request. 209 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 318885998 ns/op	       209.0 prompt_tokens	        52.80 tokens/s	       303.0 total_ms	       177.1 ttft_ms
BenchmarkCompare-32    	      20	 322565968 ns/op	       209.0 prompt_tokens	        52.30 tokens/s	       305.9 total_ms	       180.0 ttft_ms
BenchmarkCompare-32    	      20	 314898521 ns/op	       209.0 prompt_tokens	        53.49 tokens/s	       299.1 total_ms	       173.5 ttft_ms
BenchmarkCompare-32    	      20	 315070481 ns/op	       209.0 prompt_tokens	        53.44 tokens/s	       299.4 total_ms	       174.2 ttft_ms
BenchmarkCompare-32    	      20	 320341785 ns/op	       209.0 prompt_tokens	        52.68 tokens/s	       303.7 total_ms	       178.5 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	32.135s
```

</details>
<!-- yzma:bench end compare-multimodal/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-multimodal/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":29.76,"ttft_ms":441.2,"total_ms":537.6,"prompt_tokens":1216,"llamacpp":"b11179","engine_version":"1.29.0-dev","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 29.8 tokens a second. 441.2 ms to the first token. 537.6 ms for a request. 1216 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 570156105 ns/op	      1216 prompt_tokens	        30.53 tokens/s	       524.1 total_ms	       428.6 ttft_ms
BenchmarkCompare-32    	      20	 601653255 ns/op	      1216 prompt_tokens	        28.83 tokens/s	       555.0 total_ms	       458.8 ttft_ms
BenchmarkCompare-32    	      20	 576061395 ns/op	      1216 prompt_tokens	        30.20 tokens/s	       529.8 total_ms	       433.6 ttft_ms
BenchmarkCompare-32    	      20	 584696748 ns/op	      1216 prompt_tokens	        29.76 tokens/s	       537.6 total_ms	       441.2 ttft_ms
BenchmarkCompare-32    	      20	 588290726 ns/op	      1216 prompt_tokens	        29.57 tokens/s	       541.1 total_ms	       444.4 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	60.191s
```

</details>
<!-- yzma:bench end compare-multimodal/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-multimodal/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":23.74,"ttft_ms":574.2,"total_ms":674,"prompt_tokens":1221,"engine_version":"0.34.4","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 23.7 tokens a second. 574.2 ms to the first token. 674.0 ms for a request. 1221 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 692042917 ns/op	      1221 prompt_tokens	        24.58 tokens/s	       650.8 total_ms	       552.5 ttft_ms
BenchmarkCompare-32    	      20	 713258600 ns/op	      1221 prompt_tokens	        23.74 tokens/s	       674.0 total_ms	       574.2 ttft_ms
BenchmarkCompare-32    	      20	 710673504 ns/op	      1221 prompt_tokens	        23.83 tokens/s	       671.3 total_ms	       572.6 ttft_ms
BenchmarkCompare-32    	      20	 718406712 ns/op	      1221 prompt_tokens	        23.56 tokens/s	       679.1 total_ms	       579.5 ttft_ms
BenchmarkCompare-32    	      20	 726559024 ns/op	      1221 prompt_tokens	        23.35 tokens/s	       685.1 total_ms	       585.2 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	71.616s
```

</details>
<!-- yzma:bench end compare-multimodal/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-multimodal/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-multimodal","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":24.71,"ttft_ms":550.4,"total_ms":647.4,"prompt_tokens":1216,"engine_version":"v1.2.8","yzma":"1.29.0-dev","date":"2026-09-25"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 24.7 tokens a second. 550.4 ms to the first token. 647.4 ms for a request. 1216 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 25 17:46:48 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   52C    P0            590W /  115W |      15MiB /   8188MiB |     16%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7549      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 682115779 ns/op	      1216 prompt_tokens	        24.86 tokens/s	       643.6 total_ms	       546.5 ttft_ms
BenchmarkCompare-32    	      20	 686154759 ns/op	      1216 prompt_tokens	        24.71 tokens/s	       647.4 total_ms	       550.4 ttft_ms
BenchmarkCompare-32    	      20	 684177649 ns/op	      1216 prompt_tokens	        24.82 tokens/s	       644.6 total_ms	       547.2 ttft_ms
BenchmarkCompare-32    	      20	 690173451 ns/op	      1216 prompt_tokens	        24.56 tokens/s	       651.5 total_ms	       553.4 ttft_ms
BenchmarkCompare-32    	      20	 691251234 ns/op	      1216 prompt_tokens	        24.48 tokens/s	       653.5 total_ms	       555.1 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	69.021s
```

</details>
<!-- yzma:bench end compare-multimodal/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

## Each engine brings its own llama.cpp

The engines do not share one llama.cpp. yzma uses the build of its library
directory. Docker Model Runner pins a build in its image. ollama has a fork of
its own, with a version that does not map to a build of llama.cpp.

| Engine | llama.cpp |
| --- | --- |
| yzma | the build of `lib/`, b11179 of September 2026 here |
| Docker Model Runner | pinned in the image, b9879 of July 2026 in version 1.2 |
| ollama | a fork, version 0.4.1-dev here |

Thus the same GGUF file does not always give the same work. The count of the
prompt tokens shows it, and the script prints that count for each engine.

## What is not measured yet

The cold start is not here. Each table gives the numbers of an engine that has
the model in memory already. yzma also wins the cold start, because it has no
daemon to start, but that number needs a convention that these tables do not
have yet.
