# Benchmarks

- [Text model benchmarks](#text-model-benchmarks)
- [Multimodal model benchmarks](#multimodal-model-benchmarks)

## Text model benchmarks

These benchmarks all use the [SmolLM-135M-GGUF](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf) model to perform simple text generation.

See https://github.com/hybridgroup/yzma/blob/main/pkg/llama/benchmark_test.go

### Linux

#### CPU

```
$ go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                 96         111641439 ns/op               268.7 tokens/s
BenchmarkInference-32                100         109954726 ns/op               272.8 tokens/s
BenchmarkInference-32                 99         111312200 ns/op               269.5 tokens/s
BenchmarkInference-32                100         113785241 ns/op               263.7 tokens/s
BenchmarkInference-32                 93         112232836 ns/op               267.3 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   57.943s
```

#### CUDA

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
$ YZMA_BENCHMARK_DEVICE="CUDA0" go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                330          35788282 ns/op               838.3 tokens/s
BenchmarkInference-32                337          35562432 ns/op               843.6 tokens/s
BenchmarkInference-32                336          35605583 ns/op               842.6 tokens/s
BenchmarkInference-32                337          35610519 ns/op               842.4 tokens/s
BenchmarkInference-32                337          35509592 ns/op               844.8 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   62.938s
```

#### Vulkan

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
$ YZMA_BENCHMARK_DEVICE="VULKAN0" go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                 32         352943116 ns/op                85.00 tokens/s
BenchmarkInference-32                 32         368368027 ns/op                81.44 tokens/s
BenchmarkInference-32                 32         348202981 ns/op                86.16 tokens/s
BenchmarkInference-32                 32         345188429 ns/op                86.91 tokens/s
BenchmarkInference-32                 32         346846063 ns/op                86.49 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   62.292s
```

```
$ YZMA_BENCHMARK_DEVICE="VULKAN1" go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                327          35299822 ns/op               849.9 tokens/s
BenchmarkInference-32                342          34959925 ns/op               858.1 tokens/s
BenchmarkInference-32                343          34859873 ns/op               860.6 tokens/s
BenchmarkInference-32                340          35097909 ns/op               854.8 tokens/s
BenchmarkInference-32                340          35057110 ns/op               855.7 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   63.269s
```

### macOS

#### Metal

```
$ go test -bench=BenchmarkInference -benchtime=10s -count=5 -v -run=nada ./pkg/llama
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Max
BenchmarkInference
BenchmarkInference-16    	     212	  56221789 ns/op	       533.6 tokens/s
BenchmarkInference-16    	     212	  56651795 ns/op	       529.6 tokens/s
BenchmarkInference-16    	     213	  56220516 ns/op	       533.6 tokens/s
BenchmarkInference-16    	     213	  56204004 ns/op	       533.8 tokens/s
BenchmarkInference-16    	     208	  57035355 ns/op	       526.0 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	60.415s
```

### Windows

#### CPU

```
C:\Users\limbo\ron\yzma>go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                 51         211158431 ns/op               142.1 tokens/s
BenchmarkInference-32                 54         211732807 ns/op               141.7 tokens/s
BenchmarkInference-32                 55         208719991 ns/op               143.7 tokens/s
BenchmarkInference-32                 55         209386684 ns/op               143.3 tokens/s
BenchmarkInference-32                 54         207956502 ns/op               144.3 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   57.616s
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
C:\Users\ron\yzma>set YZMA_BENCHMARK_DEVICE=CUDA0

C:\Users\limbo\ron\yzma>go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                254          47057470 ns/op               637.5 tokens/s
BenchmarkInference-32                252          47270689 ns/op               634.6 tokens/s
BenchmarkInference-32                254          47281096 ns/op               634.5 tokens/s
BenchmarkInference-32                252          47607048 ns/op               630.2 tokens/s
BenchmarkInference-32                249          47880971 ns/op               626.6 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   62.562s
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
C:\Users\ron\yzma>set YZMA_BENCHMARK_DEVICE=VULKAN0

C:\Users\limbo\ron\yzma>go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                 34         330604594 ns/op                90.74 tokens/s
BenchmarkInference-32                 39         300735477 ns/op                99.76 tokens/s
BenchmarkInference-32                 38         300771805 ns/op                99.74 tokens/s
BenchmarkInference-32                 39         300599485 ns/op                99.80 tokens/s
BenchmarkInference-32                 39         300543574 ns/op                99.82 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   59.380s

C:\Users\ron\yzma>set YZMA_BENCHMARK_DEVICE=VULKAN1

C:\Users\limbo\ron\yzma>go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                252          45340691 ns/op               661.7 tokens/s
BenchmarkInference-32                296          40006768 ns/op               749.9 tokens/s
BenchmarkInference-32                294          40285934 ns/op               744.7 tokens/s
BenchmarkInference-32                294          40255082 ns/op               745.2 tokens/s
BenchmarkInference-32                295          40230172 ns/op               745.7 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   67.275ss
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

```
$ go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=5 -run=nada ./pkg/mtmd/
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32                1        42969171941 ns/op               26.16 tokens/s
BenchmarkMultimodalInference-32                1        47926256808 ns/op               26.00 tokens/s
BenchmarkMultimodalInference-32                1        47091090108 ns/op               26.40 tokens/s
BenchmarkMultimodalInference-32                1        43416576116 ns/op               25.89 tokens/s
BenchmarkMultimodalInference-32                1        46271244033 ns/op               27.36 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    233.600s
```

#### CUDA

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
$ YZMA_BENCHMARK_DEVICE="CUDA0" go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=5 -run=nada ./pkg/mtmd/
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32               12        1165425512 ns/op              1013 tokens/s
BenchmarkMultimodalInference-32               10        1142957281 ns/op              1032 tokens/s
BenchmarkMultimodalInference-32               10        1004667780 ns/op              1152 tokens/s
BenchmarkMultimodalInference-32               18        1151246613 ns/op              1025 tokens/s
BenchmarkMultimodalInference-32               13         978776993 ns/op              1176 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    76.957s
```

#### Vulkan

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
$ YZMA_BENCHMARK_DEVICE="VULKAN0" go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=5 -run=nada ./pkg/mtmd/
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32                1        20030722296 ns/op               60.16 tokens/s
BenchmarkMultimodalInference-32                1        18770055628 ns/op               63.24 tokens/s
BenchmarkMultimodalInference-32                1        35520108840 ns/op               38.71 tokens/s
BenchmarkMultimodalInference-32                1        24745078194 ns/op               50.76 tokens/s
BenchmarkMultimodalInference-32                1        24393051539 ns/op               51.29 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    136.107s
```

```
$ YZMA_BENCHMARK_DEVICE="VULKAN1" go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=5 -run=nada ./pkg/mtmd/
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32                9        1587911257 ns/op               775.7 tokens/s
BenchmarkMultimodalInference-32                8        1438853680 ns/op               841.6 tokens/s
BenchmarkMultimodalInference-32               13        1348601428 ns/op               888.4 tokens/s
BenchmarkMultimodalInference-32                9        1277507362 ns/op               929.3 tokens/s
BenchmarkMultimodalInference-32               15        1420791519 ns/op               851.4 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    85.301s
```

### macOS

#### Metal

```
$ go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=1 -run=nada -v ./pkg/mtmd/
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Max
BenchmarkMultimodalInference
BenchmarkMultimodalInference-16 183 64912527 ns/op 1448 tokens/s
...
```

At present, this benchmark can only be run once.

### Windows

#### CPU

```
C:\Users\limbo\ron\yzma>go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=5 -run=nada ./pkg/mtmd/
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkMultimodalInference-32                1        31575326100 ns/op               40.66 tokens/s
BenchmarkMultimodalInference-32                1        27726097200 ns/op               42.63 tokens/s
BenchmarkMultimodalInference-32                1        27337048100 ns/op               42.76 tokens/s
BenchmarkMultimodalInference-32                1        31019634700 ns/op               39.94 tokens/s
BenchmarkMultimodalInference-32                1        25894544200 ns/op               43.33 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    148.749s
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
C:\Users\ron\yzma>set YZMA_BENCHMARK_DEVICE=CUDA0

C:\Users\limbo\ron\yzma>go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=5 -run=nada ./pkg/mtmd/
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkMultimodalInference-32               20        1021667905 ns/op              1167 tokens/s
BenchmarkMultimodalInference-32               19         983015732 ns/op              1204 tokens/s
BenchmarkMultimodalInference-32               16        1110895319 ns/op              1091 tokens/s
BenchmarkMultimodalInference-32               15         908261600 ns/op              1286 tokens/s
BenchmarkMultimodalInference-32               16         941364706 ns/op              1248 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    91.365s
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
C:\Users\limbo\ron\yzma>set YZMA_BENCHMARK_DEVICE=VULKAN0

C:\Users\limbo\ron\yzma>go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=5 -run=nada ./pkg/mtmd/
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkMultimodalInference-32                1        31921097800 ns/op               40.38 tokens/s
BenchmarkMultimodalInference-32                1        29552171600 ns/op               43.41 tokens/s
BenchmarkMultimodalInference-32                1        15080235900 ns/op               74.00 tokens/s
BenchmarkMultimodalInference-32                1        22364990200 ns/op               53.66 tokens/s
BenchmarkMultimodalInference-32                1        15404597400 ns/op               72.71 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    122.115s
```

```
C:\Users\limbo\ron\yzma>set YZMA_BENCHMARK_DEVICE=VULKAN1

C:\Users\limbo\ron\yzma>go test -bench=BenchmarkMultimodalInference -benchtime=10s -count=5 -run=nada ./pkg/mtmd/
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkMultimodalInference-32               15         890444833 ns/op              1314 tokens/s
BenchmarkMultimodalInference-32               16         952690756 ns/op              1243 tokens/s
BenchmarkMultimodalInference-32                9        1128096289 ns/op              1080 tokens/s
BenchmarkMultimodalInference-32                9        1247414711 ns/op               996.7 tokens/s
BenchmarkMultimodalInference-32               12         848426400 ns/op              1361 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    76.636s
```
