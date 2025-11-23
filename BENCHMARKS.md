# Benchmarks

These benchmarks all use the [SmolLM-135M-GGUF](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf) model to perform simple text generation.

See https://github.com/hybridgroup/yzma/blob/main/pkg/llama/benchmark_test.go

## Linux

### CPU

```
$ go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                100         111468211 ns/op               269.1 tokens/s
BenchmarkInference-32                108         109933297 ns/op               272.9 tokens/s
BenchmarkInference-32                 99         108908937 ns/op               275.5 tokens/s
BenchmarkInference-32                 96         109742158 ns/op               273.4 tokens/s
BenchmarkInference-32                100         108993386 ns/op               275.2 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   58.614s
```

### CUDA

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
BenchmarkInference-32                349          34024582 ns/op               881.7 tokens/s
BenchmarkInference-32                349          34442142 ns/op               871.0 tokens/s
BenchmarkInference-32                350          34039915 ns/op               881.3 tokens/s
BenchmarkInference-32                352          34228134 ns/op               876.5 tokens/s
BenchmarkInference-32                351          34005275 ns/op               882.2 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   61.194s
```

### Vulkan

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.275

Devices:
========
GPU0:
        apiVersion         = 1.4.305
        driverVersion      = 25.0.7
        vendorID           = 0x8086
        deviceID           = 0xa788
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = Intel(R) Graphics (RPL-S)
        driverID           = DRIVER_ID_INTEL_OPEN_SOURCE_MESA
        driverName         = Intel open-source Mesa driver
        driverInfo         = Mesa 25.0.7-0ubuntu0.24.04.2
        conformanceVersion = 1.4.0.0
        deviceUUID         = 868088a7-0400-0000-0002-000000000000
        driverUUID         = 802b0057-40c2-aed9-e538-d78b797f04f4
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
BenchmarkInference-32                 16         684923482 ns/op                43.80 tokens/s
BenchmarkInference-32                 16         683619165 ns/op                43.88 tokens/s
BenchmarkInference-32                 16         684708734 ns/op                43.81 tokens/s
BenchmarkInference-32                 16         684328523 ns/op                43.84 tokens/s
BenchmarkInference-32                 16         683029370 ns/op                43.92 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   58.073s
```

```
$ YZMA_BENCHMARK_DEVICE="VULKAN1" go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                334          35314783 ns/op               849.5 tokens/s
BenchmarkInference-32                338          35190386 ns/op               852.5 tokens/s
BenchmarkInference-32                339          35332049 ns/op               849.1 tokens/s
BenchmarkInference-32                338          35536944 ns/op               844.2 tokens/s
BenchmarkInference-32                337          35432645 ns/op               846.7 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   63.392s
```

## macOS

### Metal

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

## Windows

### CPU

```
C:\Users\ron\yzma>go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                219          54324292 ns/op               552.2 tokens/s
BenchmarkInference-32                218          54365003 ns/op               551.8 tokens/s
BenchmarkInference-32                219          54175797 ns/op               553.8 tokens/s
BenchmarkInference-32                223          53792174 ns/op               557.7 tokens/s
BenchmarkInference-32                216          55443764 ns/op               541.1 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   62.396s
```

### CUDA

```
C:\Users\ron>nvidia-smi
Sat Nov 22 16:58:06 2025
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

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A           19020    C+G   ...s\Win64\EpicGamesLauncher.exe      N/A      |
|    0   N/A  N/A           56512      C   ...989226070\b001\llama.test.exe      N/A      |
+-----------------------------------------------------------------------------------------+
```

```
C:\Users\ron\yzma>set YZMA_BENCHMARK_DEVICE=CUDA0

C:\Users\ron\yzma>go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                260          45704590 ns/op               656.4 tokens/s
BenchmarkInference-32                260          45915642 ns/op               653.4 tokens/s
BenchmarkInference-32                262          45523990 ns/op               659.0 tokens/s
BenchmarkInference-32                258          46101993 ns/op               650.7 tokens/s
BenchmarkInference-32                261          45832454 ns/op               654.6 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   64.296s
```


### Vulkan

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

C:\Users\ron\yzma>go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                 30         370791227 ns/op                80.91 tokens/s
BenchmarkInference-32                 36         326416361 ns/op                91.91 tokens/s
BenchmarkInference-32                 36         325644942 ns/op                92.12 tokens/s
BenchmarkInference-32                 36         325254353 ns/op                92.24 tokens/s
BenchmarkInference-32                 36         324711861 ns/op                92.39 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   59.649s

C:\Users\ron\yzma>set YZMA_BENCHMARK_DEVICE=VULKAN1

C:\Users\ron\yzma>go test -bench=BenchmarkInference -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                272          43940693 ns/op               682.7 tokens/s
BenchmarkInference-32                272          44012208 ns/op               681.6 tokens/s
BenchmarkInference-32                271          44027857 ns/op               681.4 tokens/s
BenchmarkInference-32                271          43851773 ns/op               684.1 tokens/s
BenchmarkInference-32                274          43559242 ns/op               688.7 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   65.798s
```
