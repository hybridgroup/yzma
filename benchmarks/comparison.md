# Engine comparison benchmarks

These benchmarks are to compare inference performance using 3 different engines that support GGUF models:

- yzma
- ollama
- Docker Model Runner

## Summary

| Suite | Result | Evidence |
| --- | --- | --- |
| Embeddings | yzma 3.3 to 3.4 times faster, 1.2 ms against 4.1 ms and 4.2 ms | Ten runs, no overlap, each engine at 29 prompt tokens and a vector of 384 |
| Text | yzma 5.8 to 12.3 percent faster, half the time to the first token | Ten runs, no overlap, each engine at the same count of prompt tokens |
| Images | No numbers. The code runs with `--suite multimodal` | The engines preprocess an image in different ways |

## Text

<!-- yzma:bench table compare-text -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 117.1 | 12.5 | 136.6 | 1.27.0 | 2026-09-22 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 105.2 | 23.8 | 152.2 | 0.34.2 | 2026-09-22 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 104.3 | 26.8 | 153.4 | v1.2.8 | 2026-09-22 |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 166.5 | 9.1 | 96.1 | 1.27.0 | 2026-09-22 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 153.2 | 15.9 | 104.5 | 0.34.2 | 2026-09-22 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 157.4 | 14.1 | 101.7 | v1.2.8 | 2026-09-22 |
<!-- yzma:bench table end compare-text -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":117.1,"ttft_ms":12.524999999999999,"total_ms":136.64999999999998,"prompt_tokens":23,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 117.1 tokens a second. 12.5 ms to the first token. 136.6 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:05:28 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   60C    P8              7W /  115W |     218MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          388569      C   /usr/lib/ollama/llama-server            194MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=10 -run=nada -bench BenchmarkCompare -engine=yzma -suite=text -tokens=16 -nctx=8192 -model=/home/ron/models/gemma-4-E2B-it-Q4_K_M.gguf -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 136403000 ns/op	        23.00 prompt_tokens	       117.8 tokens/s	       135.8 total_ms	        12.48 ttft_ms
BenchmarkCompare-32    	      20	 136039594 ns/op	        23.00 prompt_tokens	       118.2 tokens/s	       135.4 total_ms	        12.41 ttft_ms
BenchmarkCompare-32    	      20	 137284658 ns/op	        23.00 prompt_tokens	       117.1 tokens/s	       136.7 total_ms	        12.54 ttft_ms
BenchmarkCompare-32    	      20	 136798790 ns/op	        23.00 prompt_tokens	       117.5 tokens/s	       136.2 total_ms	        12.50 ttft_ms
BenchmarkCompare-32    	      20	 137082402 ns/op	        23.00 prompt_tokens	       117.3 tokens/s	       136.5 total_ms	        12.50 ttft_ms
BenchmarkCompare-32    	      20	 137232758 ns/op	        23.00 prompt_tokens	       117.1 tokens/s	       136.6 total_ms	        12.51 ttft_ms
BenchmarkCompare-32    	      20	 137512607 ns/op	        23.00 prompt_tokens	       116.9 tokens/s	       136.9 total_ms	        12.59 ttft_ms
BenchmarkCompare-32    	      20	 137779336 ns/op	        23.00 prompt_tokens	       116.7 tokens/s	       137.1 total_ms	        12.57 ttft_ms
BenchmarkCompare-32    	      20	 137426278 ns/op	        23.00 prompt_tokens	       117.0 tokens/s	       136.8 total_ms	        12.56 ttft_ms
BenchmarkCompare-32    	      20	 138141462 ns/op	        23.00 prompt_tokens	       116.4 tokens/s	       137.5 total_ms	        12.62 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	29.606s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":105.15,"ttft_ms":23.805,"total_ms":152.15,"prompt_tokens":23,"engine_version":"0.34.2","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 105.2 tokens a second. 23.8 ms to the first token. 152.2 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:05:28 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   60C    P8              7W /  115W |     218MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          388569      C   /usr/lib/ollama/llama-server            194MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=10 -run=nada -bench BenchmarkCompare -engine=ollama -suite=text -tokens=16 -nctx=8192 -server-model=yzma-bench-gemma4-e2b -ollama-url=http://localhost:11434/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 150795625 ns/op	        23.00 prompt_tokens	       106.1 tokens/s	       150.8 total_ms	        24.92 ttft_ms
BenchmarkCompare-32    	      20	 150224494 ns/op	        23.00 prompt_tokens	       106.5 tokens/s	       150.2 total_ms	        23.30 ttft_ms
BenchmarkCompare-32    	      20	 150692027 ns/op	        23.00 prompt_tokens	       106.2 tokens/s	       150.7 total_ms	        23.35 ttft_ms
BenchmarkCompare-32    	      20	 151926936 ns/op	        23.00 prompt_tokens	       105.3 tokens/s	       151.9 total_ms	        23.54 ttft_ms
BenchmarkCompare-32    	      20	 153076717 ns/op	        23.00 prompt_tokens	       104.5 tokens/s	       153.0 total_ms	        23.95 ttft_ms
BenchmarkCompare-32    	      20	 152940490 ns/op	        23.00 prompt_tokens	       104.6 tokens/s	       152.9 total_ms	        24.07 ttft_ms
BenchmarkCompare-32    	      20	 152449449 ns/op	        23.00 prompt_tokens	       105.0 tokens/s	       152.4 total_ms	        23.66 ttft_ms
BenchmarkCompare-32    	      20	 151403372 ns/op	        23.00 prompt_tokens	       105.7 tokens/s	       151.4 total_ms	        23.59 ttft_ms
BenchmarkCompare-32    	      20	 154838599 ns/op	        23.00 prompt_tokens	       103.3 tokens/s	       154.8 total_ms	        24.29 ttft_ms
BenchmarkCompare-32    	      20	 154551784 ns/op	        23.00 prompt_tokens	       103.5 tokens/s	       154.5 total_ms	        24.07 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	30.495s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":104.3,"ttft_ms":26.785,"total_ms":153.4,"prompt_tokens":23,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 104.3 tokens a second. 26.8 ms to the first token. 153.4 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:05:28 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   60C    P8              7W /  115W |     218MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          388569      C   /usr/lib/ollama/llama-server            194MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=10 -run=nada -bench BenchmarkCompare -engine=dmr -suite=text -tokens=16 -nctx=8192 -server-model=hf.co/unsloth/gemma-4-e2b-it-gguf:q4_k_m -dmr-url=http://localhost:12434/engines/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 150930687 ns/op	        23.00 prompt_tokens	       106.0 tokens/s	       150.9 total_ms	        26.41 ttft_ms
BenchmarkCompare-32    	      20	 152229759 ns/op	        23.00 prompt_tokens	       105.1 tokens/s	       152.2 total_ms	        26.68 ttft_ms
BenchmarkCompare-32    	      20	 152999645 ns/op	        23.00 prompt_tokens	       104.6 tokens/s	       153.0 total_ms	        26.85 ttft_ms
BenchmarkCompare-32    	      20	 153546103 ns/op	        23.00 prompt_tokens	       104.2 tokens/s	       153.5 total_ms	        26.94 ttft_ms
BenchmarkCompare-32    	      20	 153439540 ns/op	        23.00 prompt_tokens	       104.3 tokens/s	       153.4 total_ms	        26.67 ttft_ms
BenchmarkCompare-32    	      20	 154203790 ns/op	        23.00 prompt_tokens	       103.8 tokens/s	       154.2 total_ms	        26.92 ttft_ms
BenchmarkCompare-32    	      20	 153512257 ns/op	        23.00 prompt_tokens	       104.3 tokens/s	       153.5 total_ms	        26.65 ttft_ms
BenchmarkCompare-32    	      20	 153373601 ns/op	        23.00 prompt_tokens	       104.3 tokens/s	       153.3 total_ms	        26.72 ttft_ms
BenchmarkCompare-32    	      20	 153396025 ns/op	        23.00 prompt_tokens	       104.3 tokens/s	       153.4 total_ms	        26.90 ttft_ms
BenchmarkCompare-32    	      20	 153409005 ns/op	        23.00 prompt_tokens	       104.3 tokens/s	       153.4 total_ms	        26.99 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	30.668s
```

</details>
<!-- yzma:bench end compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":166.5,"ttft_ms":9.126,"total_ms":96.095,"prompt_tokens":22,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 166.5 tokens a second. 9.1 ms to the first token. 96.1 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:05:28 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   60C    P8              7W /  115W |     218MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          388569      C   /usr/lib/ollama/llama-server            194MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=10 -run=nada -bench BenchmarkCompare -engine=yzma -suite=text -tokens=16 -nctx=8192 -model=/home/ron/models/Qwen3-VL-2B-Instruct.Q4_K_M.gguf -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	  99370446 ns/op	        22.00 prompt_tokens	       167.4 tokens/s	        95.55 total_ms	         9.094 ttft_ms
BenchmarkCompare-32    	      20	  99050018 ns/op	        22.00 prompt_tokens	       168.0 tokens/s	        95.24 total_ms	         8.883 ttft_ms
BenchmarkCompare-32    	      20	  99418606 ns/op	        22.00 prompt_tokens	       167.4 tokens/s	        95.59 total_ms	         8.914 ttft_ms
BenchmarkCompare-32    	      20	  99564129 ns/op	        22.00 prompt_tokens	       167.1 tokens/s	        95.75 total_ms	         9.009 ttft_ms
BenchmarkCompare-32    	      20	  99765400 ns/op	        22.00 prompt_tokens	       166.7 tokens/s	        95.96 total_ms	         9.063 ttft_ms
BenchmarkCompare-32    	      20	 100033283 ns/op	        22.00 prompt_tokens	       166.3 tokens/s	        96.23 total_ms	         9.158 ttft_ms
BenchmarkCompare-32    	      20	 100249902 ns/op	        22.00 prompt_tokens	       165.9 tokens/s	        96.43 total_ms	         9.227 ttft_ms
BenchmarkCompare-32    	      20	 100172953 ns/op	        22.00 prompt_tokens	       166.1 tokens/s	        96.35 total_ms	         9.164 ttft_ms
BenchmarkCompare-32    	      20	 100394429 ns/op	        22.00 prompt_tokens	       165.7 tokens/s	        96.59 total_ms	         9.170 ttft_ms
BenchmarkCompare-32    	      20	 100531547 ns/op	        22.00 prompt_tokens	       165.4 tokens/s	        96.71 total_ms	         9.186 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	20.895s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":153.2,"ttft_ms":15.905000000000001,"total_ms":104.45,"prompt_tokens":22,"engine_version":"0.34.2","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 153.2 tokens a second. 15.9 ms to the first token. 104.5 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:05:28 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   60C    P8              7W /  115W |     218MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          388569      C   /usr/lib/ollama/llama-server            194MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=10 -run=nada -bench BenchmarkCompare -engine=ollama -suite=text -tokens=16 -nctx=8192 -server-model=yzma-bench-qwen3-vl-2b -ollama-url=http://localhost:11434/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 108708628 ns/op	        22.00 prompt_tokens	       147.2 tokens/s	       108.7 total_ms	        19.28 ttft_ms
BenchmarkCompare-32    	      20	 104519108 ns/op	        22.00 prompt_tokens	       153.1 tokens/s	       104.5 total_ms	        16.40 ttft_ms
BenchmarkCompare-32    	      20	 104306110 ns/op	        22.00 prompt_tokens	       153.4 tokens/s	       104.3 total_ms	        16.06 ttft_ms
BenchmarkCompare-32    	      20	 104029573 ns/op	        22.00 prompt_tokens	       153.8 tokens/s	       104.0 total_ms	        15.83 ttft_ms
BenchmarkCompare-32    	      20	 104313452 ns/op	        22.00 prompt_tokens	       153.4 tokens/s	       104.3 total_ms	        15.80 ttft_ms
BenchmarkCompare-32    	      20	 104681464 ns/op	        22.00 prompt_tokens	       152.9 tokens/s	       104.7 total_ms	        15.80 ttft_ms
BenchmarkCompare-32    	      20	 104156230 ns/op	        22.00 prompt_tokens	       153.7 tokens/s	       104.1 total_ms	        15.90 ttft_ms
BenchmarkCompare-32    	      20	 104412414 ns/op	        22.00 prompt_tokens	       153.3 tokens/s	       104.4 total_ms	        15.91 ttft_ms
BenchmarkCompare-32    	      20	 104775629 ns/op	        22.00 prompt_tokens	       152.7 tokens/s	       104.8 total_ms	        16.10 ttft_ms
BenchmarkCompare-32    	      20	 105266772 ns/op	        22.00 prompt_tokens	       152.0 tokens/s	       105.2 total_ms	        15.87 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	21.021s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":157.35000000000002,"ttft_ms":14.135000000000002,"total_ms":101.7,"prompt_tokens":22,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 157.4 tokens a second. 14.1 ms to the first token. 101.7 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:05:28 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   60C    P8              7W /  115W |     218MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          388569      C   /usr/lib/ollama/llama-server            194MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=20x -count=10 -run=nada -bench BenchmarkCompare -engine=dmr -suite=text -tokens=16 -nctx=8192 -server-model=hf.co/qwen/qwen3-vl-2b-instruct-gguf:q4_k_m -dmr-url=http://localhost:12434/engines/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      20	 105802900 ns/op	        22.00 prompt_tokens	       151.3 tokens/s	       105.8 total_ms	        17.72 ttft_ms
BenchmarkCompare-32    	      20	 101589625 ns/op	        22.00 prompt_tokens	       157.6 tokens/s	       101.6 total_ms	        14.37 ttft_ms
BenchmarkCompare-32    	      20	 101230593 ns/op	        22.00 prompt_tokens	       158.1 tokens/s	       101.2 total_ms	        14.03 ttft_ms
BenchmarkCompare-32    	      20	 101559509 ns/op	        22.00 prompt_tokens	       157.6 tokens/s	       101.5 total_ms	        13.81 ttft_ms
BenchmarkCompare-32    	      20	 101622132 ns/op	        22.00 prompt_tokens	       157.5 tokens/s	       101.6 total_ms	        14.21 ttft_ms
BenchmarkCompare-32    	      20	 101702861 ns/op	        22.00 prompt_tokens	       157.4 tokens/s	       101.7 total_ms	        14.07 ttft_ms
BenchmarkCompare-32    	      20	 101842953 ns/op	        22.00 prompt_tokens	       157.2 tokens/s	       101.8 total_ms	        14.13 ttft_ms
BenchmarkCompare-32    	      20	 101815527 ns/op	        22.00 prompt_tokens	       157.2 tokens/s	       101.8 total_ms	        14.14 ttft_ms
BenchmarkCompare-32    	      20	 101907870 ns/op	        22.00 prompt_tokens	       157.1 tokens/s	       101.9 total_ms	        14.09 ttft_ms
BenchmarkCompare-32    	      20	 101728160 ns/op	        22.00 prompt_tokens	       157.3 tokens/s	       101.7 total_ms	        14.16 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	20.445s
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
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 23653.5 | 1.2 | 1.2 | 1.27.0 | 2026-09-22 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 6953.5 | 4.2 | 4.2 | 0.34.2 | 2026-09-22 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 7079.5 | 4.1 | 4.1 | v1.2.8 | 2026-09-22 |
<!-- yzma:bench table end compare-embeddings -->

<!-- yzma:bench start compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":23653.5,"ttft_ms":1.226,"total_ms":1.226,"prompt_tokens":29,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 23653.5 tokens a second. 1.2 ms to the first token. 1.2 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:08:45 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   84C    P0             71W /  115W |    3000MiB /   8188MiB |     82%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          400827      C   /app/llama-server                      2976MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=50x -count=10 -run=nada -bench BenchmarkCompare -engine=yzma -suite=embeddings -tokens=16 -nctx=8192 -model=/home/ron/models/bge-small-en-v1.5-q8_0.gguf -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      50	   1378795 ns/op	        29.00 prompt_tokens	     21060 tokens/s	         1.377 total_ms	         1.377 ttft_ms
BenchmarkCompare-32    	      50	   1339886 ns/op	        29.00 prompt_tokens	     21682 tokens/s	         1.337 total_ms	         1.337 ttft_ms
BenchmarkCompare-32    	      50	   1338032 ns/op	        29.00 prompt_tokens	     21709 tokens/s	         1.336 total_ms	         1.336 ttft_ms
BenchmarkCompare-32    	      50	   1244858 ns/op	        29.00 prompt_tokens	     23331 tokens/s	         1.243 total_ms	         1.243 ttft_ms
BenchmarkCompare-32    	      50	   1208103 ns/op	        29.00 prompt_tokens	     24044 tokens/s	         1.206 total_ms	         1.206 ttft_ms
BenchmarkCompare-32    	      50	   1222645 ns/op	        29.00 prompt_tokens	     23764 tokens/s	         1.220 total_ms	         1.220 ttft_ms
BenchmarkCompare-32    	      50	   1213908 ns/op	        29.00 prompt_tokens	     23931 tokens/s	         1.212 total_ms	         1.212 ttft_ms
BenchmarkCompare-32    	      50	   1233856 ns/op	        29.00 prompt_tokens	     23543 tokens/s	         1.232 total_ms	         1.232 ttft_ms
BenchmarkCompare-32    	      50	   1202693 ns/op	        29.00 prompt_tokens	     24150 tokens/s	         1.201 total_ms	         1.201 ttft_ms
BenchmarkCompare-32    	      50	   1201775 ns/op	        29.00 prompt_tokens	     24172 tokens/s	         1.200 total_ms	         1.200 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	1.076s
```

</details>
<!-- yzma:bench end compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":6953.5,"ttft_ms":4.170999999999999,"total_ms":4.170999999999999,"prompt_tokens":29,"engine_version":"0.34.2","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 6953.5 tokens a second. 4.2 ms to the first token. 4.2 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:08:45 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   84C    P0             71W /  115W |    3000MiB /   8188MiB |     82%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          400827      C   /app/llama-server                      2976MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=50x -count=10 -run=nada -bench BenchmarkCompare -engine=ollama -suite=embeddings -tokens=16 -nctx=8192 -server-model=yzma-bench-bge-small -ollama-url=http://localhost:11434/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      50	   4209450 ns/op	        29.00 prompt_tokens	      6895 tokens/s	         4.206 total_ms	         4.206 ttft_ms
BenchmarkCompare-32    	      50	   4092344 ns/op	        29.00 prompt_tokens	      7092 tokens/s	         4.089 total_ms	         4.089 ttft_ms
BenchmarkCompare-32    	      50	   4076307 ns/op	        29.00 prompt_tokens	      7121 tokens/s	         4.072 total_ms	         4.072 ttft_ms
BenchmarkCompare-32    	      50	   4313272 ns/op	        29.00 prompt_tokens	      6730 tokens/s	         4.309 total_ms	         4.309 ttft_ms
BenchmarkCompare-32    	      50	   4177851 ns/op	        29.00 prompt_tokens	      6947 tokens/s	         4.175 total_ms	         4.175 ttft_ms
BenchmarkCompare-32    	      50	   4042427 ns/op	        29.00 prompt_tokens	      7181 tokens/s	         4.039 total_ms	         4.039 ttft_ms
BenchmarkCompare-32    	      50	   4178664 ns/op	        29.00 prompt_tokens	      6947 tokens/s	         4.175 total_ms	         4.175 ttft_ms
BenchmarkCompare-32    	      50	   4170376 ns/op	        29.00 prompt_tokens	      6960 tokens/s	         4.167 total_ms	         4.167 ttft_ms
BenchmarkCompare-32    	      50	   4139670 ns/op	        29.00 prompt_tokens	      7011 tokens/s	         4.136 total_ms	         4.136 ttft_ms
BenchmarkCompare-32    	      50	   4332859 ns/op	        29.00 prompt_tokens	      6699 tokens/s	         4.329 total_ms	         4.329 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	2.107s
```

</details>
<!-- yzma:bench end compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":7079.5,"ttft_ms":4.0965,"total_ms":4.0965,"prompt_tokens":29,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 7079.5 tokens a second. 4.1 ms to the first token. 4.1 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Tue Sep 22 22:08:45 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   84C    P0             71W /  115W |    3000MiB /   8188MiB |     82%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
|    0   N/A  N/A          400827      C   /app/llama-server                      2976MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd benchmarks/compare && go test -benchtime=50x -count=10 -run=nada -bench BenchmarkCompare -engine=dmr -suite=embeddings -tokens=16 -nctx=8192 -server-model=hf.co/ggml-org/bge-small-en-v1.5-q8_0-gguf:q8_0 -dmr-url=http://localhost:12434/engines/v1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/benchmarks/compare
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkCompare-32    	      50	   4068277 ns/op	        29.00 prompt_tokens	      7136 tokens/s	         4.064 total_ms	         4.064 ttft_ms
BenchmarkCompare-32    	      50	   4134543 ns/op	        29.00 prompt_tokens	      7023 tokens/s	         4.129 total_ms	         4.129 ttft_ms
BenchmarkCompare-32    	      50	   4448712 ns/op	        29.00 prompt_tokens	      6527 tokens/s	         4.443 total_ms	         4.443 ttft_ms
BenchmarkCompare-32    	      50	   4511909 ns/op	        29.00 prompt_tokens	      6436 tokens/s	         4.506 total_ms	         4.506 ttft_ms
BenchmarkCompare-32    	      50	   4260254 ns/op	        29.00 prompt_tokens	      6816 tokens/s	         4.255 total_ms	         4.255 ttft_ms
BenchmarkCompare-32    	      50	   3493618 ns/op	        29.00 prompt_tokens	      8311 tokens/s	         3.489 total_ms	         3.489 ttft_ms
BenchmarkCompare-32    	      50	   3904358 ns/op	        29.00 prompt_tokens	      7436 tokens/s	         3.900 total_ms	         3.900 ttft_ms
BenchmarkCompare-32    	      50	   3759771 ns/op	        29.00 prompt_tokens	      7723 tokens/s	         3.755 total_ms	         3.755 ttft_ms
BenchmarkCompare-32    	      50	   4275770 ns/op	        29.00 prompt_tokens	      6791 tokens/s	         4.270 total_ms	         4.270 ttft_ms
BenchmarkCompare-32    	      50	   3688088 ns/op	        29.00 prompt_tokens	      7873 tokens/s	         3.684 total_ms	         3.684 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	2.051s
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
