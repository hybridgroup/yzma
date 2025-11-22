# Benchmarks

These benchmarks all use the [SmolLM-135M-GGUF](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf) model to perform simple text generation.

See https://github.com/hybridgroup/yzma/blob/main/pkg/llama/benchmark_test.go

## Linux

### CPU

```
$ go test -bench=. -benchmem -benchtime=10s -count=5 -run=^$ ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                100         119200742 ns/op           16272 B/op        723 allocs/op
BenchmarkInference-32                 99         120473276 ns/op           16292 B/op        723 allocs/op
BenchmarkInference-32                 99         119045172 ns/op           16286 B/op        723 allocs/op
BenchmarkInference-32                100         116843066 ns/op           16271 B/op        723 allocs/op
BenchmarkInference-32                 99         117814229 ns/op           16286 B/op        723 allocs/op
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   88.048s
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
$ YZMA_BENCHMARK_DEVICE="CUDA0" go test -bench=. -benchmem -benchtime=10s -count=5 -run=^$ ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                345          34503738 ns/op           15328 B/op        695 allocs/op
BenchmarkInference-32                337          34534140 ns/op           15337 B/op        696 allocs/op
BenchmarkInference-32                340          34567047 ns/op           15335 B/op        696 allocs/op
BenchmarkInference-32                339          34642013 ns/op           15336 B/op        696 allocs/op
BenchmarkInference-32                340          34459238 ns/op           15333 B/op        695 allocs/op
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   81.587s
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
$ YZMA_BENCHMARK_DEVICE="VULKAN0" go test -bench=. -benchmem -benchtime=10s -count=5 -run=^$ ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                 16         673672999 ns/op           21560 B/op        877 allocs/op
BenchmarkInference-32                 16         671572114 ns/op           21588 B/op        877 allocs/op
BenchmarkInference-32                 16         674732293 ns/op           21587 B/op        877 allocs/op
BenchmarkInference-32                 14         754711859 ns/op           22555 B/op        904 allocs/op
BenchmarkInference-32                 16         673296851 ns/op           21586 B/op        877 allocs/op
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   91.437s
```

```
$ YZMA_BENCHMARK_DEVICE="VULKAN1" go test -bench=. -benchmem -benchtime=10s -count=5 -run=^$ ./pkg/llama
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                291          37219117 ns/op           15405 B/op        697 allocs/op
BenchmarkInference-32                312          37387668 ns/op           15361 B/op        696 allocs/op
BenchmarkInference-32                314          37462380 ns/op           15361 B/op        696 allocs/op
BenchmarkInference-32                315          36991080 ns/op           15360 B/op        696 allocs/op
BenchmarkInference-32                315          37441630 ns/op           15359 B/op        696 allocs/op
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   72.853s
```

## macOS

### Metal

```
$ go test -bench=. -benchmem -benchtime=10s -count=5 -run=^$ -v ./pkg/llama
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Max
BenchmarkInference
BenchmarkInference-16                207          57066276 ns/op           15490 B/op        701 allocs/op
BenchmarkInference-16                207          57458342 ns/op           15495 B/op        701 allocs/op
BenchmarkInference-16                207          57461009 ns/op           15492 B/op        701 allocs/op
BenchmarkInference-16                208          56645718 ns/op           15490 B/op        701 allocs/op
BenchmarkInference-16                211          56484412 ns/op           15485 B/op        700 allocs/op
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   86.878s
```

## Windows

### CPU

```
C:\Users\ron\yzma>go test -bench=BenchmarkInference -benchmem -benchtime=10s -count=5 -run=nada ./pkg/llama
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor
BenchmarkInference-32                205          55864720 ns/op           15147 B/op        690 allocs/op
BenchmarkInference-32                214          55511057 ns/op           15140 B/op        690 allocs/op
BenchmarkInference-32                214          55102790 ns/op           15141 B/op        690 allocs/op
BenchmarkInference-32                218          56286605 ns/op           15139 B/op        690 allocs/op
BenchmarkInference-32                217          55053613 ns/op           15139 B/op        690 allocs/op
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   81.436s
```

### CUDA

Coming soon...

### Vulkan

Coming soon...
