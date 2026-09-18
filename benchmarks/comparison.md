# Engine comparison benchmarks

These benchmarks are to compare inference performance using 3 different engines that support GGUF models:

- yzma
- ollama
- Docker Model Runner

## Summary

| Suite | Result | Evidence |
| --- | --- | --- |
| Embeddings | yzma 3.4 to 3.8 times faster, 1.2 ms against 4.1 ms and 4.6 ms | Ten runs, no overlap, each engine at 29 prompt tokens and a vector of 384 |
| Text | yzma 5.6 to 12.3 percent faster, half the time to the first token | Ten runs, no overlap, each engine at the same count of prompt tokens |
| Images | No numbers. The code runs with `--suite multimodal` | The engines preprocess an image in different ways |

## Text

<!-- yzma:bench table compare-text -->
| Engine | Arch | Machine | Model | Prompt tokens | Tokens a second | First token ms | Request ms | Version | Date |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 118.8 | 12.3 | 134.8 | 1.27.0 | 2026-09-18 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 106.3 | 23.5 | 150.5 | 0.34.2 | 2026-09-18 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | gemma4-e2b | 23 | 105.7 | 26.1 | 151.4 | v1.2.8 | 2026-09-18 |
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 168.6 | 8.8 | 94.9 | 1.27.0 | 2026-09-18 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 154.2 | 16.1 | 103.8 | 0.34.2 | 2026-09-18 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | qwen3-vl-2b | 22 | 159.6 | 14.0 | 100.2 | v1.2.8 | 2026-09-18 |
<!-- yzma:bench table end compare-text -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":118.75,"ttft_ms":12.27,"total_ms":134.75,"prompt_tokens":23,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 118.8 tokens a second. 12.3 ms to the first token. 134.8 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:02:53 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   51C    P8              7W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 136210866 ns/op	        23.00 prompt_tokens	       118.0 tokens/s	       135.6 total_ms	        12.44 ttft_ms
BenchmarkCompare-32    	      20	 134802324 ns/op	        23.00 prompt_tokens	       119.2 tokens/s	       134.2 total_ms	        12.23 ttft_ms
BenchmarkCompare-32    	      20	 135044789 ns/op	        23.00 prompt_tokens	       119.0 tokens/s	       134.4 total_ms	        12.22 ttft_ms
BenchmarkCompare-32    	      20	 135184683 ns/op	        23.00 prompt_tokens	       118.9 tokens/s	       134.6 total_ms	        12.25 ttft_ms
BenchmarkCompare-32    	      20	 135109443 ns/op	        23.00 prompt_tokens	       119.0 tokens/s	       134.5 total_ms	        12.23 ttft_ms
BenchmarkCompare-32    	      20	 135124337 ns/op	        23.00 prompt_tokens	       119.0 tokens/s	       134.5 total_ms	        12.19 ttft_ms
BenchmarkCompare-32    	      20	 135518564 ns/op	        23.00 prompt_tokens	       118.6 tokens/s	       134.9 total_ms	        12.29 ttft_ms
BenchmarkCompare-32    	      20	 135560823 ns/op	        23.00 prompt_tokens	       118.6 tokens/s	       134.9 total_ms	        12.30 ttft_ms
BenchmarkCompare-32    	      20	 136304558 ns/op	        23.00 prompt_tokens	       117.9 tokens/s	       135.7 total_ms	        12.37 ttft_ms
BenchmarkCompare-32    	      20	 136203096 ns/op	        23.00 prompt_tokens	       118.0 tokens/s	       135.6 total_ms	        12.30 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	29.355s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":106.35,"ttft_ms":23.520000000000003,"total_ms":150.5,"prompt_tokens":23,"engine_version":"0.34.2","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 106.3 tokens a second. 23.5 ms to the first token. 150.5 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:02:53 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   51C    P8              7W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 151561427 ns/op	        23.00 prompt_tokens	       105.6 tokens/s	       151.5 total_ms	        25.28 ttft_ms
BenchmarkCompare-32    	      20	 149284508 ns/op	        23.00 prompt_tokens	       107.2 tokens/s	       149.2 total_ms	        23.53 ttft_ms
BenchmarkCompare-32    	      20	 149601235 ns/op	        23.00 prompt_tokens	       107.0 tokens/s	       149.6 total_ms	        23.20 ttft_ms
BenchmarkCompare-32    	      20	 148343410 ns/op	        23.00 prompt_tokens	       107.9 tokens/s	       148.3 total_ms	        23.34 ttft_ms
BenchmarkCompare-32    	      20	 149699759 ns/op	        23.00 prompt_tokens	       106.9 tokens/s	       149.7 total_ms	        23.58 ttft_ms
BenchmarkCompare-32    	      20	 152194000 ns/op	        23.00 prompt_tokens	       105.1 tokens/s	       152.2 total_ms	        24.09 ttft_ms
BenchmarkCompare-32    	      20	 150402015 ns/op	        23.00 prompt_tokens	       106.4 tokens/s	       150.4 total_ms	        23.53 ttft_ms
BenchmarkCompare-32    	      20	 151162112 ns/op	        23.00 prompt_tokens	       105.9 tokens/s	       151.1 total_ms	        23.51 ttft_ms
BenchmarkCompare-32    	      20	 150599504 ns/op	        23.00 prompt_tokens	       106.3 tokens/s	       150.6 total_ms	        23.50 ttft_ms
BenchmarkCompare-32    	      20	 151576805 ns/op	        23.00 prompt_tokens	       105.6 tokens/s	       151.6 total_ms	        23.47 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	30.128s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, gemma4-e2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"gemma4-e2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":105.7,"ttft_ms":26.08,"total_ms":151.35000000000002,"prompt_tokens":23,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 105.7 tokens a second. 26.1 ms to the first token. 151.4 ms for a request. 23 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:02:53 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   51C    P8              7W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 149551810 ns/op	        23.00 prompt_tokens	       107.0 tokens/s	       149.5 total_ms	        25.98 ttft_ms
BenchmarkCompare-32    	      20	 149926427 ns/op	        23.00 prompt_tokens	       106.7 tokens/s	       149.9 total_ms	        25.74 ttft_ms
BenchmarkCompare-32    	      20	 151363114 ns/op	        23.00 prompt_tokens	       105.7 tokens/s	       151.3 total_ms	        26.09 ttft_ms
BenchmarkCompare-32    	      20	 151614164 ns/op	        23.00 prompt_tokens	       105.6 tokens/s	       151.6 total_ms	        26.33 ttft_ms
BenchmarkCompare-32    	      20	 151252162 ns/op	        23.00 prompt_tokens	       105.8 tokens/s	       151.2 total_ms	        26.07 ttft_ms
BenchmarkCompare-32    	      20	 151054090 ns/op	        23.00 prompt_tokens	       106.0 tokens/s	       151.0 total_ms	        25.90 ttft_ms
BenchmarkCompare-32    	      20	 151722760 ns/op	        23.00 prompt_tokens	       105.5 tokens/s	       151.7 total_ms	        26.27 ttft_ms
BenchmarkCompare-32    	      20	 151926633 ns/op	        23.00 prompt_tokens	       105.3 tokens/s	       151.9 total_ms	        26.46 ttft_ms
BenchmarkCompare-32    	      20	 151448487 ns/op	        23.00 prompt_tokens	       105.7 tokens/s	       151.4 total_ms	        26.07 ttft_ms
BenchmarkCompare-32    	      20	 151418766 ns/op	        23.00 prompt_tokens	       105.7 tokens/s	       151.4 total_ms	        26.24 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	30.274s
```

</details>
<!-- yzma:bench end compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/gemma4-e2b -->

<!-- yzma:bench start compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":168.6,"ttft_ms":8.789,"total_ms":94.92,"prompt_tokens":22,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 168.6 tokens a second. 8.8 ms to the first token. 94.9 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:02:53 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   51C    P8              7W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 102422559 ns/op	        22.00 prompt_tokens	       162.5 tokens/s	        98.47 total_ms	         9.229 ttft_ms
BenchmarkCompare-32    	      20	  98708910 ns/op	        22.00 prompt_tokens	       168.6 tokens/s	        94.89 total_ms	         8.792 ttft_ms
BenchmarkCompare-32    	      20	  98731977 ns/op	        22.00 prompt_tokens	       168.6 tokens/s	        94.92 total_ms	         8.794 ttft_ms
BenchmarkCompare-32    	      20	  98841164 ns/op	        22.00 prompt_tokens	       168.4 tokens/s	        95.03 total_ms	         8.815 ttft_ms
BenchmarkCompare-32    	      20	  98735334 ns/op	        22.00 prompt_tokens	       168.6 tokens/s	        94.92 total_ms	         8.768 ttft_ms
BenchmarkCompare-32    	      20	  98723166 ns/op	        22.00 prompt_tokens	       168.6 tokens/s	        94.91 total_ms	         8.770 ttft_ms
BenchmarkCompare-32    	      20	  98625862 ns/op	        22.00 prompt_tokens	       168.7 tokens/s	        94.82 total_ms	         8.747 ttft_ms
BenchmarkCompare-32    	      20	  98736278 ns/op	        22.00 prompt_tokens	       168.6 tokens/s	        94.92 total_ms	         8.816 ttft_ms
BenchmarkCompare-32    	      20	  98672847 ns/op	        22.00 prompt_tokens	       168.7 tokens/s	        94.86 total_ms	         8.786 ttft_ms
BenchmarkCompare-32    	      20	  98748405 ns/op	        22.00 prompt_tokens	       168.5 tokens/s	        94.93 total_ms	         8.780 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	20.747s
```

</details>
<!-- yzma:bench end compare-text/yzma/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":154.25,"ttft_ms":16.119999999999997,"total_ms":103.75,"prompt_tokens":22,"engine_version":"0.34.2","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 154.2 tokens a second. 16.1 ms to the first token. 103.8 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:02:53 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   51C    P8              7W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 108674730 ns/op	        22.00 prompt_tokens	       147.3 tokens/s	       108.6 total_ms	        19.67 ttft_ms
BenchmarkCompare-32    	      20	 102711436 ns/op	        22.00 prompt_tokens	       155.8 tokens/s	       102.7 total_ms	        15.84 ttft_ms
BenchmarkCompare-32    	      20	 103596158 ns/op	        22.00 prompt_tokens	       154.5 tokens/s	       103.6 total_ms	        16.11 ttft_ms
BenchmarkCompare-32    	      20	 103381334 ns/op	        22.00 prompt_tokens	       154.8 tokens/s	       103.4 total_ms	        16.03 ttft_ms
BenchmarkCompare-32    	      20	 103608686 ns/op	        22.00 prompt_tokens	       154.5 tokens/s	       103.6 total_ms	        16.02 ttft_ms
BenchmarkCompare-32    	      20	 103946108 ns/op	        22.00 prompt_tokens	       154.0 tokens/s	       103.9 total_ms	        16.14 ttft_ms
BenchmarkCompare-32    	      20	 103254106 ns/op	        22.00 prompt_tokens	       155.0 tokens/s	       103.2 total_ms	        16.04 ttft_ms
BenchmarkCompare-32    	      20	 103930928 ns/op	        22.00 prompt_tokens	       154.0 tokens/s	       103.9 total_ms	        16.32 ttft_ms
BenchmarkCompare-32    	      20	 104016718 ns/op	        22.00 prompt_tokens	       153.9 tokens/s	       104.0 total_ms	        16.13 ttft_ms
BenchmarkCompare-32    	      20	 104364314 ns/op	        22.00 prompt_tokens	       153.3 tokens/s	       104.3 total_ms	        16.13 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	20.866s
```

</details>
<!-- yzma:bench end compare-text/ollama/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->

<!-- yzma:bench start compare-text/dmr/amd64/ron-tuxedo-gemini-gen2/qwen3-vl-2b -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, qwen3-vl-2b
<!-- yzma:bench meta {"suite":"compare-text","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"qwen3-vl-2b","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":159.6,"ttft_ms":13.965,"total_ms":100.25,"prompt_tokens":22,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 159.6 tokens a second. 14.0 ms to the first token. 100.2 ms for a request. 22 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:02:53 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   51C    P8              7W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      20	 106302199 ns/op	        22.00 prompt_tokens	       150.6 tokens/s	       106.3 total_ms	        18.54 ttft_ms
BenchmarkCompare-32    	      20	 100042101 ns/op	        22.00 prompt_tokens	       160.0 tokens/s	       100.0 total_ms	        14.18 ttft_ms
BenchmarkCompare-32    	      20	 100059459 ns/op	        22.00 prompt_tokens	       160.0 tokens/s	       100.0 total_ms	        13.95 ttft_ms
BenchmarkCompare-32    	      20	 100256118 ns/op	        22.00 prompt_tokens	       159.6 tokens/s	       100.2 total_ms	        14.04 ttft_ms
BenchmarkCompare-32    	      20	 100265225 ns/op	        22.00 prompt_tokens	       159.6 tokens/s	       100.2 total_ms	        13.80 ttft_ms
BenchmarkCompare-32    	      20	 100358026 ns/op	        22.00 prompt_tokens	       159.5 tokens/s	       100.3 total_ms	        13.98 ttft_ms
BenchmarkCompare-32    	      20	 100547290 ns/op	        22.00 prompt_tokens	       159.2 tokens/s	       100.5 total_ms	        13.95 ttft_ms
BenchmarkCompare-32    	      20	 100304288 ns/op	        22.00 prompt_tokens	       159.6 tokens/s	       100.3 total_ms	        13.71 ttft_ms
BenchmarkCompare-32    	      20	 100283469 ns/op	        22.00 prompt_tokens	       159.6 tokens/s	       100.2 total_ms	        13.59 ttft_ms
BenchmarkCompare-32    	      20	 100560576 ns/op	        22.00 prompt_tokens	       159.2 tokens/s	       100.5 total_ms	        13.99 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	20.216s
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
| yzma, in process | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 23943.5 | 1.2 | 1.2 | 1.27.0 | 2026-09-18 |
| ollama, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 6363.5 | 4.6 | 4.6 | 0.34.2 | 2026-09-18 |
| Docker Model Runner, REST | amd64 | Intel Core i9-13900HX, RTX 4070 | bge-small | 29 | 7137.5 | 4.1 | 4.1 | v1.2.8 | 2026-09-18 |
<!-- yzma:bench table end compare-embeddings -->

<!-- yzma:bench start compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### yzma, in process, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"yzma","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":23943.5,"ttft_ms":1.2109999999999999,"total_ms":1.2109999999999999,"prompt_tokens":29,"llamacpp":"b10964","engine_version":"1.27.0","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 23943.5 tokens a second. 1.2 ms to the first token. 1.2 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:01:07 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              5W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      50	   1379592 ns/op	        29.00 prompt_tokens	     21048 tokens/s	         1.378 total_ms	         1.378 ttft_ms
BenchmarkCompare-32    	      50	   1347618 ns/op	        29.00 prompt_tokens	     21551 tokens/s	         1.346 total_ms	         1.346 ttft_ms
BenchmarkCompare-32    	      50	   1330413 ns/op	        29.00 prompt_tokens	     21831 tokens/s	         1.328 total_ms	         1.328 ttft_ms
BenchmarkCompare-32    	      50	   1184319 ns/op	        29.00 prompt_tokens	     24518 tokens/s	         1.183 total_ms	         1.183 ttft_ms
BenchmarkCompare-32    	      50	   1230181 ns/op	        29.00 prompt_tokens	     23612 tokens/s	         1.228 total_ms	         1.228 ttft_ms
BenchmarkCompare-32    	      50	   1198717 ns/op	        29.00 prompt_tokens	     24238 tokens/s	         1.196 total_ms	         1.196 ttft_ms
BenchmarkCompare-32    	      50	   1195378 ns/op	        29.00 prompt_tokens	     24293 tokens/s	         1.194 total_ms	         1.194 ttft_ms
BenchmarkCompare-32    	      50	   1190073 ns/op	        29.00 prompt_tokens	     24404 tokens/s	         1.188 total_ms	         1.188 ttft_ms
BenchmarkCompare-32    	      50	   1224237 ns/op	        29.00 prompt_tokens	     23728 tokens/s	         1.222 total_ms	         1.222 ttft_ms
BenchmarkCompare-32    	      50	   1202523 ns/op	        29.00 prompt_tokens	     24159 tokens/s	         1.200 total_ms	         1.200 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	1.124s
```

</details>
<!-- yzma:bench end compare-embeddings/yzma/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### ollama, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"ollama","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":6363.5,"ttft_ms":4.5585,"total_ms":4.5585,"prompt_tokens":29,"engine_version":"0.34.2","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 6363.5 tokens a second. 4.6 ms to the first token. 4.6 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:01:07 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              5W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      50	   4881724 ns/op	        29.00 prompt_tokens	      5945 tokens/s	         4.878 total_ms	         4.878 ttft_ms
BenchmarkCompare-32    	      50	   4876982 ns/op	        29.00 prompt_tokens	      5951 tokens/s	         4.873 total_ms	         4.873 ttft_ms
BenchmarkCompare-32    	      50	   4786748 ns/op	        29.00 prompt_tokens	      6063 tokens/s	         4.783 total_ms	         4.783 ttft_ms
BenchmarkCompare-32    	      50	   4644642 ns/op	        29.00 prompt_tokens	      6250 tokens/s	         4.640 total_ms	         4.640 ttft_ms
BenchmarkCompare-32    	      50	   4468155 ns/op	        29.00 prompt_tokens	      6497 tokens/s	         4.464 total_ms	         4.464 ttft_ms
BenchmarkCompare-32    	      50	   4819264 ns/op	        29.00 prompt_tokens	      6023 tokens/s	         4.815 total_ms	         4.815 ttft_ms
BenchmarkCompare-32    	      50	   4336000 ns/op	        29.00 prompt_tokens	      6694 tokens/s	         4.332 total_ms	         4.332 ttft_ms
BenchmarkCompare-32    	      50	   4480974 ns/op	        29.00 prompt_tokens	      6477 tokens/s	         4.477 total_ms	         4.477 ttft_ms
BenchmarkCompare-32    	      50	   4342832 ns/op	        29.00 prompt_tokens	      6683 tokens/s	         4.339 total_ms	         4.339 ttft_ms
BenchmarkCompare-32    	      50	   4233120 ns/op	        29.00 prompt_tokens	      6856 tokens/s	         4.230 total_ms	         4.230 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	2.315s
```

</details>
<!-- yzma:bench end compare-embeddings/ollama/amd64/ron-tuxedo-gemini-gen2/bge-small -->

<!-- yzma:bench start compare-embeddings/dmr/amd64/ron-tuxedo-gemini-gen2/bge-small -->
### Docker Model Runner, REST, amd64, Intel Core i9-13900HX, RTX 4070, bge-small
<!-- yzma:bench meta {"suite":"compare-embeddings","backend":"dmr","arch":"amd64","machine":"ron-tuxedo-gemini-gen2","model":"bge-small","label":"Intel Core i9-13900HX, RTX 4070","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":7137.5,"ttft_ms":4.063000000000001,"total_ms":4.063000000000001,"prompt_tokens":29,"engine_version":"v1.2.8","yzma":"1.27.0","date":"2026-09-18"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 7137.5 tokens a second. 4.1 ms to the first token. 4.1 ms for a request. 29 prompt tokens.

16 tokens, greedy sampling, one request at a time.

<details><summary>The device</summary>

```
Fri Sep 18 17:01:07 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P8              5W /  115W |      16MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkCompare-32    	      50	   4062419 ns/op	        29.00 prompt_tokens	      7148 tokens/s	         4.057 total_ms	         4.057 ttft_ms
BenchmarkCompare-32    	      50	   4039859 ns/op	        29.00 prompt_tokens	      7187 tokens/s	         4.035 total_ms	         4.035 ttft_ms
BenchmarkCompare-32    	      50	   3599542 ns/op	        29.00 prompt_tokens	      8065 tokens/s	         3.596 total_ms	         3.596 ttft_ms
BenchmarkCompare-32    	      50	   4439740 ns/op	        29.00 prompt_tokens	      6541 tokens/s	         4.433 total_ms	         4.433 ttft_ms
BenchmarkCompare-32    	      50	   4074073 ns/op	        29.00 prompt_tokens	      7127 tokens/s	         4.069 total_ms	         4.069 ttft_ms
BenchmarkCompare-32    	      50	   4425678 ns/op	        29.00 prompt_tokens	      6561 tokens/s	         4.420 total_ms	         4.420 ttft_ms
BenchmarkCompare-32    	      50	   3849324 ns/op	        29.00 prompt_tokens	      7546 tokens/s	         3.843 total_ms	         3.843 ttft_ms
BenchmarkCompare-32    	      50	   4249413 ns/op	        29.00 prompt_tokens	      6837 tokens/s	         4.242 total_ms	         4.242 ttft_ms
BenchmarkCompare-32    	      50	   4557897 ns/op	        29.00 prompt_tokens	      6371 tokens/s	         4.552 total_ms	         4.552 ttft_ms
BenchmarkCompare-32    	      50	   3414754 ns/op	        29.00 prompt_tokens	      8507 tokens/s	         3.409 total_ms	         3.409 ttft_ms
PASS
ok  	github.com/hybridgroup/yzma/benchmarks/compare	2.060s
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
