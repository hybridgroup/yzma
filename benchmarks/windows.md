# Windows benchmarks

Benchmarks of yzma on Windows. Each table gives the median of five runs. The output
of each run, and of the device, is below the tables.

To add a machine or to make these numbers again, see
[how to run the benchmarks](README.md).

## Text model benchmarks

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf).
The code is [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go).

<!-- yzma:bench table text -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | amd64 | AMD Ryzen 9 7950X | - | 144.9 | unknown | unknown |
| CUDA | amd64 | AMD Ryzen 9 7950X | CUDA0 | 639.3 | unknown | unknown |
| Vulkan | amd64 | AMD Ryzen 9 7950X | Vulkan0 | 98.5 | unknown | unknown |
| Vulkan | amd64 | AMD Ryzen 9 7950X | Vulkan1 | 739.9 | unknown | unknown |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/ryzen-9-7950x -->
### CPU, amd64, AMD Ryzen 9 7950X
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"ryzen-9-7950x","label":"AMD Ryzen 9 7950X","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":144.9} -->

AMD Ryzen 9 7950X 16-Core Processor. 144.9 tokens a second.

<details><summary>The output of go test</summary>

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

</details>
<!-- yzma:bench end text/cpu/amd64/ryzen-9-7950x -->

<!-- yzma:bench start text/cuda/amd64/ryzen-9-7950x/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"ryzen-9-7950x","device":"CUDA0","label":"AMD Ryzen 9 7950X","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":639.3} -->

AMD Ryzen 9 7950X 16-Core Processor. 639.3 tokens a second.

<details><summary>The device</summary>

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

</details>

<details><summary>The output of go test</summary>

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

</details>
<!-- yzma:bench end text/cuda/amd64/ryzen-9-7950x/cuda0 -->

<!-- yzma:bench start text/vulkan/amd64/ryzen-9-7950x/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"ryzen-9-7950x","device":"Vulkan0","label":"AMD Ryzen 9 7950X","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":98.51} -->

AMD Ryzen 9 7950X 16-Core Processor. 98.5 tokens a second.

<details><summary>The device</summary>

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

</details>

<details><summary>The output of go test</summary>

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
```

</details>
<!-- yzma:bench end text/vulkan/amd64/ryzen-9-7950x/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/ryzen-9-7950x/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X, Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"ryzen-9-7950x","device":"Vulkan1","label":"AMD Ryzen 9 7950X","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":739.9} -->

AMD Ryzen 9 7950X 16-Core Processor. 739.9 tokens a second.

<details><summary>The device</summary>

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

</details>

<details><summary>The output of go test</summary>

```
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

</details>
<!-- yzma:bench end text/vulkan/amd64/ryzen-9-7950x/vulkan1 -->

## Multimodal model benchmarks

The model is
[Qwen3-VL-2B-Instruct.Q4_K_M.gguf](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.Q4_K_M.gguf)
with its
[projector](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf).
The code is [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go).

<!-- yzma:bench table multimodal -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | amd64 | AMD Ryzen 9 7950X | - | 39.5 | unknown | unknown |
| CUDA | amd64 | AMD Ryzen 9 7950X | CUDA0 | 1086.0 | unknown | unknown |
| Vulkan | amd64 | AMD Ryzen 9 7950X | Vulkan0 | 73.1 | unknown | unknown |
| Vulkan | amd64 | AMD Ryzen 9 7950X | Vulkan1 | 1180.0 | unknown | unknown |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/ryzen-9-7950x -->
### CPU, amd64, AMD Ryzen 9 7950X
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"ryzen-9-7950x","label":"AMD Ryzen 9 7950X","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":39.52} -->

AMD Ryzen 9 7950X 16-Core Processor. 39.5 tokens a second.

<details><summary>The output of go test</summary>

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

</details>
<!-- yzma:bench end multimodal/cpu/amd64/ryzen-9-7950x -->

<!-- yzma:bench start multimodal/cuda/amd64/ryzen-9-7950x/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"ryzen-9-7950x","device":"CUDA0","label":"AMD Ryzen 9 7950X","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":1086} -->

AMD Ryzen 9 7950X 16-Core Processor. 1086.0 tokens a second.

<details><summary>The device</summary>

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

</details>

<details><summary>The output of go test</summary>

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

</details>
<!-- yzma:bench end multimodal/cuda/amd64/ryzen-9-7950x/cuda0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/ryzen-9-7950x/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"ryzen-9-7950x","device":"Vulkan0","label":"AMD Ryzen 9 7950X","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":73.08} -->

AMD Ryzen 9 7950X 16-Core Processor. 73.1 tokens a second.

<details><summary>The device</summary>

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

</details>

<details><summary>The output of go test</summary>

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

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/ryzen-9-7950x/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/ryzen-9-7950x/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X, Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"ryzen-9-7950x","device":"Vulkan1","label":"AMD Ryzen 9 7950X","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":1180} -->

AMD Ryzen 9 7950X 16-Core Processor. 1180.0 tokens a second.

<details><summary>The device</summary>

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

</details>

<details><summary>The output of go test</summary>

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

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/ryzen-9-7950x/vulkan1 -->
