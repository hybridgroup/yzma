# Benchmarks

- [Text model benchmarks](#text-model-benchmarks)
- [Multimodal model benchmarks](#multimodal-model-benchmarks)

## Text model benchmarks

These benchmarks all use the [SmolLM-135M-GGUF](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf) model to perform simple text generation.

```shell
yzma model get -u https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf
export YZMA_BENCHMARK_MODEL=~/models/SmolLM-135M.Q2_K.gguf
```

See https://github.com/hybridgroup/yzma/blob/main/pkg/llama/benchmark_test.go

### Linux

#### CPU

##### amd64

```
$ go test -benchtime=10s -count=5 -run=nada -bench .
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                 99         110913774 ns/op               270.5 tokens/s
BenchmarkInference-32                100         111035054 ns/op               270.2 tokens/s
BenchmarkInference-32                100         110369390 ns/op               271.8 tokens/s
BenchmarkInference-32                100         112705133 ns/op               266.2 tokens/s
BenchmarkInference-32                100         111892770 ns/op               268.1 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   61.199s
```

##### arm64

Raspberry Pi 4 Model B Rev 1.4 8GB

```
ron@raspberrypi:~/yzma/pkg/llama $ go test -benchtime=10s -count=5 -run=nada -bench .
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-4          15         893788634 ns/op                33.56 tokens/s
BenchmarkInference-4          12         923948131 ns/op                32.47 tokens/s
BenchmarkInference-4          12         918284434 ns/op                32.67 tokens/s
BenchmarkInference-4          12         918693617 ns/op                32.66 tokens/s
BenchmarkInference-4          12         917186754 ns/op                32.71 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   64.583s
```

#### CUDA

##### amd64

```
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 580.95.05              Driver Version: 580.95.05      CUDA Version: 13.0     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   38C    P0            590W /  115W |      15MiB /   8188MiB |     17%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+
```

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="CUDA0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                332          35746370 ns/op               839.2 tokens/s
BenchmarkInference-32                338          35529926 ns/op               844.4 tokens/s
BenchmarkInference-32                336          35614579 ns/op               842.4 tokens/s
BenchmarkInference-32                336          35609522 ns/op               842.5 tokens/s
BenchmarkInference-32                337          35550352 ns/op               843.9 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   67.491s
```

##### arm64

Jetson Orin Nano Developer Kit - 8GB

```
+---------------------------------------------------------------------------------------+
| NVIDIA-SMI 540.5.0                Driver Version: 540.5.0      CUDA Version: 12.6     |
|-----------------------------------------+----------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |         Memory-Usage | GPU-Util  Compute M. |
|                                         |                      |               MIG M. |
|=========================================+======================+======================|
|   0  Orin (nvgpu)                  N/A  | N/A              N/A |                  N/A |
| N/A   N/A  N/A               N/A /  N/A | Not Supported        |     N/A          N/A |
|                                         |                      |                  N/A |
+-----------------------------------------+----------------------+----------------------+
```

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="CUDA0"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkInference-6          51         222138094 ns/op               135.1 tokens/s
BenchmarkInference-6          52         216104925 ns/op               138.8 tokens/s
BenchmarkInference-6          54         215961553 ns/op               138.9 tokens/s
BenchmarkInference-6          52         215498575 ns/op               139.2 tokens/s
BenchmarkInference-6          52         214849130 ns/op               139.6 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   61.014s
```

#### Vulkan

##### amd64

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.275

Devices:
========
GPU0:
        apiVersion         = 1.4.318
        driverVersion      = 25.2.8
        vendorID           = 0x8086
        deviceID           = 0xa788
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = Intel(R) Graphics (RPL-S)
        driverID           = DRIVER_ID_INTEL_OPEN_SOURCE_MESA
        driverName         = Intel open-source Mesa driver
        driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.1
        conformanceVersion = 1.4.0.0
        deviceUUID         = 868088a7-0400-0000-0002-000000000000
        driverUUID         = 032fbbbb-ddee-3516-8477-c17071969177
GPU1:
        apiVersion         = 1.4.312
        driverVersion      = 580.95.5.0
        vendorID           = 0x10de
        deviceID           = 0x2860
        deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
        deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 580.95.05
        conformanceVersion = 1.4.1.3
        deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
        driverUUID         = b92269a1-b525-5615-ab8a-e2095ee37192
```

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="VULKAN0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                 31         354329548 ns/op                84.67 tokens/s
BenchmarkInference-32                 34         351859490 ns/op                85.26 tokens/s
BenchmarkInference-32                 32         353665267 ns/op                84.83 tokens/s
BenchmarkInference-32                 33         349151210 ns/op                85.92 tokens/s
BenchmarkInference-32                 33         348216889 ns/op                86.15 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   70.757s
```

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="VULKAN1"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                328          36362981 ns/op               825.0 tokens/s
BenchmarkInference-32                330          36353223 ns/op               825.2 tokens/s
BenchmarkInference-32                327          36207519 ns/op               828.6 tokens/s
BenchmarkInference-32                331          36366451 ns/op               824.9 tokens/s
BenchmarkInference-32                330          36262953 ns/op               827.3 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   83.142s
```

##### arm64

Jetson Orin Nano Developer Kit - 8GB

```
ron@ubuntu:~/yzma/pkg/mtmd$ vulkaninfo --summary
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.204
...
Devices:
========
GPU0:
        apiVersion         = 4206843 (1.3.251)
        driverVersion      = 2265006080 (0x87014000)
        vendorID           = 0x10de
        deviceID           = 0x97ba03d7
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = NVIDIA Tegra Orin (nvgpu)
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 540.5.0
        conformanceVersion = 1.3.6.0
        deviceUUID         = 1388f9e0-987e-54a0-908f-6a30d8fd5f29
        driverUUID         = ed5ba772-f592-5949-9d1f-236f7ad81bcc
```

```
ron@ubuntu:~/yzma/pkg/llama$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="CPU"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkInference-6          43         432432689 ns/op                69.37 tokens/s
BenchmarkInference-6          20         506747397 ns/op                59.20 tokens/s
BenchmarkInference-6          21         514736186 ns/op                58.28 tokens/s
BenchmarkInference-6          27         496646058 ns/op                60.41 tokens/s
BenchmarkInference-6          22         519434233 ns/op                57.76 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   68.009s
```

```
ron@ubuntu:~/yzma/pkg/llama$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="VULKAN0"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkInference-6          52         222098600 ns/op               135.1 tokens/s
BenchmarkInference-6          52         222072877 ns/op               135.1 tokens/s
BenchmarkInference-6          54         219825013 ns/op               136.5 tokens/s
BenchmarkInference-6          52         220919304 ns/op               135.8 tokens/s
BenchmarkInference-6          54         221925680 ns/op               135.2 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   63.318s
```

### macOS

#### Metal

Apple M4 Max with 128 GB RAM

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

### Windows

#### CPU

```
C:\Users\limbo\ron\yzma\pkg\llama>go test -benchtime=10s -count=5 -run=nada -bench . -nctx=8192
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                 51         214577557 ns/op               139.8 tokens/s
BenchmarkInference-32                 56         210247484 ns/op               142.7 tokens/s
BenchmarkInference-32                 52         206580071 ns/op               145.2 tokens/s
BenchmarkInference-32                 57         206447956 ns/op               145.3 tokens/s
BenchmarkInference-32                 57         207005089 ns/op               144.9 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   58.254s
```

#### CUDA

```
C:\Users\ron>nvidia-smi
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 581.57                 Driver Version: 581.57         CUDA Version: 13.0     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   42C    P8              6W /  240W |      22MiB /   8192MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+
```

```
C:\Users\limbo\ron\yzma\pkg\llama>go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="CUDA0"
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                254          46914384 ns/op               639.5 tokens/s
BenchmarkInference-32                258          46820920 ns/op               640.7 tokens/s
BenchmarkInference-32                255          46929827 ns/op               639.3 tokens/s
BenchmarkInference-32                255          46958283 ns/op               638.9 tokens/s
BenchmarkInference-32                250          47880058 ns/op               626.6 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   62.888s
```


#### Vulkan

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.309


Devices:
========
GPU0:
        apiVersion         = 1.3.270
        driverVersion      = 2.0.294
        vendorID           = 0x1002
        deviceID           = 0x164e
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = AMD Radeon(TM) Graphics
        driverID           = DRIVER_ID_AMD_PROPRIETARY
        driverName         = AMD proprietary driver
        driverInfo         = 23.40.02 (AMD proprietary shader compiler)
        conformanceVersion = 1.3.3.1
        deviceUUID         = 00000000-0c00-0000-0000-000000000000
        driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
        apiVersion         = 1.4.312
        driverVersion      = 581.57.0.0
        vendorID           = 0x10de
        deviceID           = 0x2488
        deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
        deviceName         = NVIDIA GeForce RTX 3070
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 581.57
        conformanceVersion = 1.4.1.3
        deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
        driverUUID         = 08a6deb5-2838-56d3-b7da-f79802447960
```

```
C:\Users\limbo\ron\yzma\pkg\llama>go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="VULKAN0"
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                 34         329955426 ns/op                90.92 tokens/s
BenchmarkInference-32                 39         302329823 ns/op                99.23 tokens/s
BenchmarkInference-32                 39         302524487 ns/op                99.17 tokens/s
BenchmarkInference-32                 39         304700162 ns/op                98.46 tokens/s
BenchmarkInference-32                 39         304536574 ns/op                98.51 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   61.326s

C:\Users\limbo\ron\yzma\pkg\llama>go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="VULKAN1"
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                294          40543699 ns/op               739.9 tokens/s
BenchmarkInference-32                295          40568015 ns/op               739.5 tokens/s
BenchmarkInference-32                295          40579471 ns/op               739.3 tokens/s
BenchmarkInference-32                297          40277643 ns/op               744.8 tokens/s
BenchmarkInference-32                296          40319531 ns/op               744.1 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   84.981s
```

## Multimodal model benchmarks

These benchmarks all use the [Qwen3-VL-2B-Instruct.Q4_K_M.gguf](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.Q4_K_M.gguf) model and [projector](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf) to provide a description for an image.

```shell
yzma model get -u https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.Q4_K_M.gguf
yzma model get -u https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf
export YZMA_BENCHMARK_MMMODEL=~/models/Qwen3-VL-2B-Instruct.Q4_K_M.gguf
export YZMA_BENCHMARK_MMPROJ=~/models/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf
```

See https://github.com/hybridgroup/yzma/blob/main/pkg/mtmd/benchmark_test.go

### Linux

#### CPU

##### amd64

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=8192
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32                1        47402263232 ns/op               26.16 tokens/s
BenchmarkMultimodalInference-32                1        42673907034 ns/op               26.08 tokens/s
BenchmarkMultimodalInference-32                1        42432080672 ns/op               25.81 tokens/s
BenchmarkMultimodalInference-32                1        46803510445 ns/op               26.15 tokens/s
BenchmarkMultimodalInference-32                1        45700830384 ns/op               25.91 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    226.685s
```

##### arm64

Raspberry Pi 4 Model B Rev 1.4 8GB

NOTE: Due to less available memory, the benchmarks on this device used the [SmolVLM2-500M-Video-Instruct-Q8_0](https://huggingface.co/ggml-org/SmolVLM2-500M-Video-Instruct-GGUF) model.


```
ron@raspberrypi:~/yzma/pkg/mtmd $ export YZMA_BENCHMARK_MMMODEL=/home/ron/models/SmolVLM2-500M-Video-Instruct-Q8_0.gguf
ron@raspberrypi:~/yzma/pkg/mtmd $ export YZMA_BENCHMARK_MMPROJ=/home/ron/models/mmproj-SmolVLM2-500M-Video-Instruct-Q8_0.gguf
ron@raspberrypi:~/yzma/pkg/mtmd $ go test -benchtime=10s -count=5 -run=nada -bench .
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-4                 1        50239133481 ns/op                6.748 tokens/s
BenchmarkMultimodalInference-4                 1        49358181828 ns/op                6.341 tokens/s
BenchmarkMultimodalInference-4                 1        48164506831 ns/op                5.917 tokens/s
BenchmarkMultimodalInference-4                 1        40171997080 ns/op                5.551 tokens/s
BenchmarkMultimodalInference-4                 1        41428165840 ns/op                5.504 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    243.876s
```

#### CUDA

##### amd64

```
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 580.95.05              Driver Version: 580.95.05      CUDA Version: 13.0     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   38C    P0            590W /  115W |      15MiB /   8188MiB |     17%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+
```

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="CUDA0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32               21         921205057 ns/op              1240 tokens/s
BenchmarkMultimodalInference-32               15        1043496530 ns/op              1114 tokens/s
BenchmarkMultimodalInference-32               18         939373857 ns/op              1219 tokens/s
BenchmarkMultimodalInference-32               14        1118362797 ns/op              1047 tokens/s
BenchmarkMultimodalInference-32                8        1336574088 ns/op               900.2 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    82.619s
```

##### arm64

Jetson Orin Nano Developer Kit - 8GB

```
+---------------------------------------------------------------------------------------+
| NVIDIA-SMI 540.5.0                Driver Version: 540.5.0      CUDA Version: 12.6     |
|-----------------------------------------+----------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |         Memory-Usage | GPU-Util  Compute M. |
|                                         |                      |               MIG M. |
|=========================================+======================+======================|
|   0  Orin (nvgpu)                  N/A  | N/A              N/A |                  N/A |
| N/A   N/A  N/A               N/A /  N/A | Not Supported        |     N/A          N/A |
|                                         |                      |                  N/A |
+-----------------------------------------+----------------------+----------------------+
```

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="CUDA0"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkMultimodalInference-6                 2        7077293280 ns/op               166.9 tokens/s
BenchmarkMultimodalInference-6                 2        8106794026 ns/op               150.8 tokens/s
BenchmarkMultimodalInference-6                 1        10837943077 ns/op              120.7 tokens/s
BenchmarkMultimodalInference-6                 1        12015033493 ns/op              112.1 tokens/s
BenchmarkMultimodalInference-6                 1        10055887615 ns/op              127.6 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    69.733s
```

#### Vulkan

##### amd64

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.275

Devices:
========
GPU0:
        apiVersion         = 1.4.318
        driverVersion      = 25.2.8
        vendorID           = 0x8086
        deviceID           = 0xa788
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = Intel(R) Graphics (RPL-S)
        driverID           = DRIVER_ID_INTEL_OPEN_SOURCE_MESA
        driverName         = Intel open-source Mesa driver
        driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.1
        conformanceVersion = 1.4.0.0
        deviceUUID         = 868088a7-0400-0000-0002-000000000000
        driverUUID         = 032fbbbb-ddee-3516-8477-c17071969177
GPU1:
        apiVersion         = 1.4.312
        driverVersion      = 580.95.5.0
        vendorID           = 0x10de
        deviceID           = 0x2860
        deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
        deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 580.95.05
        conformanceVersion = 1.4.1.3
        deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
        driverUUID         = b92269a1-b525-5615-ab8a-e2095ee37192
```

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="VULKAN0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32                1        14578268628 ns/op               78.34 tokens/s
BenchmarkMultimodalInference-32                1        22073783877 ns/op               55.59 tokens/s
BenchmarkMultimodalInference-32                1        11278156188 ns/op               97.62 tokens/s
BenchmarkMultimodalInference-32                1        14723860691 ns/op               77.43 tokens/s
BenchmarkMultimodalInference-32                1        11996066619 ns/op               92.45 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    79.922s
```

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="VULKAN1"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32                8        1339951138 ns/op               891.1 tokens/s
BenchmarkMultimodalInference-32               10        1172385505 ns/op               997.5 tokens/s
BenchmarkMultimodalInference-32               13        1276183643 ns/op               929.1 tokens/s
BenchmarkMultimodalInference-32               18        1122849292 ns/op              1035 tokens/s
BenchmarkMultimodalInference-32                7        1471154871 ns/op               825.9 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    76.276s
```

##### arm64

Jetson Orin Nano Developer Kit - 8GB

```
ron@ubuntu:~/yzma/pkg/mtmd$ vulkaninfo --summary
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.204
...
Devices:
========
GPU0:
        apiVersion         = 4206843 (1.3.251)
        driverVersion      = 2265006080 (0x87014000)
        vendorID           = 0x10de
        deviceID           = 0x97ba03d7
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = NVIDIA Tegra Orin (nvgpu)
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 540.5.0
        conformanceVersion = 1.3.6.0
        deviceUUID         = 1388f9e0-987e-54a0-908f-6a30d8fd5f29
        driverUUID         = ed5ba772-f592-5949-9d1f-236f7ad81bcc
```

```
ron@ubuntu:~/yzma/pkg/mtmd$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="CPU"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkMultimodalInference-6                 1        72233629960 ns/op               15.03 tokens/s
BenchmarkMultimodalInference-6                 1        75555489707 ns/op               15.37 tokens/s
BenchmarkMultimodalInference-6                 1        87238792057 ns/op               14.65 tokens/s
BenchmarkMultimodalInference-6                 1        71406835155 ns/op               15.70 tokens/s
BenchmarkMultimodalInference-6                 1        70659234723 ns/op               15.74 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    383.358s
```

```
ron@ubuntu:~/yzma/pkg/mtmd$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="VULKAN0"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkMultimodalInference-6                 1        13718208893 ns/op               81.13 tokens/s
BenchmarkMultimodalInference-6                 1        16724822437 ns/op               71.39 tokens/s
BenchmarkMultimodalInference-6                 1        13133369170 ns/op               84.14 tokens/s
BenchmarkMultimodalInference-6                 1        13515072899 ns/op               82.43 tokens/s
BenchmarkMultimodalInference-6                 1        12471954537 ns/op               87.24 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    76.766s
```

### macOS

#### Metal

Apple M4 Max with 128 GB RAM

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

### Windows

#### CPU

```
C:\Users\limbo\ron\yzma\pkg\mtmd>go test -benchtime=10s -count=5 -run=nada -bench . -nctx=8192
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkMultimodalInference-32                1        26850046400 ns/op               43.17 tokens/s
BenchmarkMultimodalInference-32                1        48420966900 ns/op               35.44 tokens/s
BenchmarkMultimodalInference-32                1        34259612500 ns/op               39.52 tokens/s
BenchmarkMultimodalInference-32                1        24749935100 ns/op               44.44 tokens/s
BenchmarkMultimodalInference-32                1        36232681200 ns/op               38.75 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    171.920s
```

#### CUDA

```
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 581.57                 Driver Version: 581.57         CUDA Version: 13.0     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   42C    P8              6W /  240W |      22MiB /   8192MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+
```

```
C:\Users\limbo\ron\yzma\pkg\mtmd>go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="CUDA0"
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkMultimodalInference-32               14         975072514 ns/op              1212 tokens/s
BenchmarkMultimodalInference-32                9        1124768556 ns/op              1080 tokens/s
BenchmarkMultimodalInference-32                9        1138583744 ns/op              1071 tokens/s
BenchmarkMultimodalInference-32               10        1099877300 ns/op              1099 tokens/s
BenchmarkMultimodalInference-32               10        1116220610 ns/op              1086 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    57.908s
```

#### Vulkan

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.309


Devices:
========
GPU0:
        apiVersion         = 1.3.270
        driverVersion      = 2.0.294
        vendorID           = 0x1002
        deviceID           = 0x164e
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = AMD Radeon(TM) Graphics
        driverID           = DRIVER_ID_AMD_PROPRIETARY
        driverName         = AMD proprietary driver
        driverInfo         = 23.40.02 (AMD proprietary shader compiler)
        conformanceVersion = 1.3.3.1
        deviceUUID         = 00000000-0c00-0000-0000-000000000000
        driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
        apiVersion         = 1.4.312
        driverVersion      = 581.57.0.0
        vendorID           = 0x10de
        deviceID           = 0x2488
        deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
        deviceName         = NVIDIA GeForce RTX 3070
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 581.57
        conformanceVersion = 1.4.1.3
        deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
        driverUUID         = 08a6deb5-2838-56d3-b7da-f79802447960
```

```
C:\Users\limbo\ron\yzma\pkg\mtmd>go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="VULKAN0"
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkMultimodalInference-32                1        14997592100 ns/op               73.08 tokens/s
BenchmarkMultimodalInference-32                1        14469341200 ns/op               76.71 tokens/s
BenchmarkMultimodalInference-32                1        24988773000 ns/op               49.22 tokens/s
BenchmarkMultimodalInference-32                1        24924637400 ns/op               49.35 tokens/s
BenchmarkMultimodalInference-32                1        14559276800 ns/op               76.31 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    96.114s
```

```
C:\Users\limbo\ron\yzma\pkg\mtmd>go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="VULKAN1"
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkMultimodalInference-32               16         937497038 ns/op              1262 tokens/s
BenchmarkMultimodalInference-32               20        1079753220 ns/op              1126 tokens/s
BenchmarkMultimodalInference-32               19        1003840647 ns/op              1194 tokens/s
BenchmarkMultimodalInference-32                9        1535556511 ns/op               856.7 tokens/s
BenchmarkMultimodalInference-32               12        1018743817 ns/op              1180 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    90.525s
```
