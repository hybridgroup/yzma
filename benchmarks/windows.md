# Windows benchmarks

Benchmarks of yzma on Windows. Each table gives the median of five runs. The output
of each run, and of the device, is below the tables.

To add a machine or to make these numbers again, see
[how to run the benchmarks](README.md).

## Summary

Tokens a second on each machine. The GPU columns give the fastest GPU backend.

| Machine | GPU | Text, CPU | Text, GPU | Multimodal, CPU | Multimodal, GPU |
| --- | --- | --- | --- | --- | --- |
| AMD Ryzen 9 7950X | RTX 3070 | 115.0 | 803.2, Vulkan | 289.5 | 2018.0, Vulkan |

- On the RTX 3070, Vulkan is faster than CUDA. It is 14 percent faster for text
  and 13 percent faster for multimodal.
- The integrated AMD Radeon GPU on Vulkan0 gives 114.3 for text, which is
  almost the same as the CPU. For multimodal it gives 479.2, which is faster
  than the CPU.
- The CPU runs change much from one run to the next, because Windows cannot
  hold a thread to a core. The multimodal runs on the CPU go from 262.5 to
  315.2.

## Text model benchmarks

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf).
The code is [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go).
The benchmark uses 4 threads on each machine, see
[the thread count](README.md#run-them).

<!-- yzma:bench table text -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | - | 115.0 | b10964 | 2026-09-23 |
| CUDA | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | CUDA0 | 707.4 | b10964 | 2026-09-23 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan0 | 114.3 | b10964 | 2026-09-23 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan1 | 803.2 | b10964 | 2026-09-23 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/desktop-b8a29kd -->
### CPU, amd64, AMD Ryzen 9 7950X 16-Core Processor            
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"desktop-b8a29kd","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":115,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

AMD Ryzen 9 7950X 16-Core Processor. 115.0 tokens a second.

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	      45	 258961422 ns/op	       115.8 tokens/s
BenchmarkInference-32    	      45	 262115822 ns/op	       114.5 tokens/s
BenchmarkInference-32    	      45	 256640851 ns/op	       116.9 tokens/s
BenchmarkInference-32    	      42	 260806229 ns/op	       115.0 tokens/s
BenchmarkInference-32    	      44	 266618600 ns/op	       112.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	58.786s
```

</details>
<!-- yzma:bench end text/cpu/amd64/desktop-b8a29kd -->

<!-- yzma:bench start text/cuda/amd64/desktop-b8a29kd/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X 16-Core Processor            , CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"desktop-b8a29kd","device":"CUDA0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":707.4,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

AMD Ryzen 9 7950X 16-Core Processor. 707.4 tokens a second.

<details><summary>The device</summary>

```
Wed Sep 23 19:24:55 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 591.86                 Driver Version: 591.86         CUDA Version: 13.1     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   39C    P0              9W /  240W |       0MiB /   8192MiB |     10%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|  No running processes found                                                             |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=CUDA0
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	     279	  42565482 ns/op	       704.8 tokens/s
BenchmarkInference-32    	     282	  42406014 ns/op	       707.4 tokens/s
BenchmarkInference-32    	     282	  42423854 ns/op	       707.1 tokens/s
BenchmarkInference-32    	     282	  42400097 ns/op	       707.5 tokens/s
BenchmarkInference-32    	     284	  42116328 ns/op	       712.3 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.834s
```

</details>
<!-- yzma:bench end text/cuda/amd64/desktop-b8a29kd/cuda0 -->

<!-- yzma:bench start text/vulkan/amd64/desktop-b8a29kd/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":114.3,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

AMD Ryzen 9 7950X 16-Core Processor. 114.3 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.321


Instance Extensions: count = 20
-------------------------------
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_surface_maintenance1            : extension revision 1
VK_EXT_swapchain_colorspace            : extension revision 5
VK_KHR_device_group_creation           : extension revision 1
VK_KHR_display                         : extension revision 23
VK_KHR_external_fence_capabilities     : extension revision 1
VK_KHR_external_memory_capabilities    : extension revision 1
VK_KHR_external_semaphore_capabilities : extension revision 1
VK_KHR_get_display_properties2         : extension revision 1
VK_KHR_get_physical_device_properties2 : extension revision 2
VK_KHR_get_surface_capabilities2       : extension revision 1
VK_KHR_portability_enumeration         : extension revision 1
VK_KHR_surface                         : extension revision 25
VK_KHR_surface_maintenance1            : extension revision 1
VK_KHR_surface_protected_capabilities  : extension revision 1
VK_KHR_win32_surface                   : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1
VK_NV_external_memory_capabilities     : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_AMD_switchable_graphics AMD switchable graphics layer 1.4.315  version 1
VK_LAYER_NV_optimus              NVIDIA Optimus layer          1.4.325  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer     1.4.325  version 1
VK_LAYER_VALVE_steam_fossilize   Steam Pipeline Caching Layer  1.4.303  version 1
VK_LAYER_VALVE_steam_overlay     Steam Overlay Layer           1.3.207  version 1

Devices:
========
GPU0:
	apiVersion         = 1.4.315
	driverVersion      = 2.0.353
	vendorID           = 0x1002
	deviceID           = 0x164e
	deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
	deviceName         = AMD Radeon(TM) Graphics
	driverID           = DRIVER_ID_AMD_PROPRIETARY
	driverName         = AMD proprietary driver
	driverInfo         = 26.3.1 (AMD proprietary shader compiler)
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-0c00-0000-0000-000000000000
	driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
	apiVersion         = 1.4.325
	driverVersion      = 591.86.0.0
	vendorID           = 0x10de
	deviceID           = 0x2488
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 3070
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 591.86
	conformanceVersion = 1.4.3.0
	deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
	driverUUID         = ab24b0bb-7bd9-59fc-a310-4c2bdaba9b72
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=Vulkan0
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	      38	 293363245 ns/op	       102.3 tokens/s
BenchmarkInference-32    	      44	 262383202 ns/op	       114.3 tokens/s
BenchmarkInference-32    	      45	 262079087 ns/op	       114.5 tokens/s
BenchmarkInference-32    	      45	 262501896 ns/op	       114.3 tokens/s
BenchmarkInference-32    	      45	 262116369 ns/op	       114.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.115s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/desktop-b8a29kd/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan1","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":803.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

AMD Ryzen 9 7950X 16-Core Processor. 803.2 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.321


Instance Extensions: count = 20
-------------------------------
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_surface_maintenance1            : extension revision 1
VK_EXT_swapchain_colorspace            : extension revision 5
VK_KHR_device_group_creation           : extension revision 1
VK_KHR_display                         : extension revision 23
VK_KHR_external_fence_capabilities     : extension revision 1
VK_KHR_external_memory_capabilities    : extension revision 1
VK_KHR_external_semaphore_capabilities : extension revision 1
VK_KHR_get_display_properties2         : extension revision 1
VK_KHR_get_physical_device_properties2 : extension revision 2
VK_KHR_get_surface_capabilities2       : extension revision 1
VK_KHR_portability_enumeration         : extension revision 1
VK_KHR_surface                         : extension revision 25
VK_KHR_surface_maintenance1            : extension revision 1
VK_KHR_surface_protected_capabilities  : extension revision 1
VK_KHR_win32_surface                   : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1
VK_NV_external_memory_capabilities     : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_AMD_switchable_graphics AMD switchable graphics layer 1.4.315  version 1
VK_LAYER_NV_optimus              NVIDIA Optimus layer          1.4.325  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer     1.4.325  version 1
VK_LAYER_VALVE_steam_fossilize   Steam Pipeline Caching Layer  1.4.303  version 1
VK_LAYER_VALVE_steam_overlay     Steam Overlay Layer           1.3.207  version 1

Devices:
========
GPU0:
	apiVersion         = 1.4.315
	driverVersion      = 2.0.353
	vendorID           = 0x1002
	deviceID           = 0x164e
	deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
	deviceName         = AMD Radeon(TM) Graphics
	driverID           = DRIVER_ID_AMD_PROPRIETARY
	driverName         = AMD proprietary driver
	driverInfo         = 26.3.1 (AMD proprietary shader compiler)
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-0c00-0000-0000-000000000000
	driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
	apiVersion         = 1.4.325
	driverVersion      = 591.86.0.0
	vendorID           = 0x10de
	deviceID           = 0x2488
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 3070
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 591.86
	conformanceVersion = 1.4.3.0
	deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
	driverUUID         = ab24b0bb-7bd9-59fc-a310-4c2bdaba9b72
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=Vulkan1
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	     318	  37424038 ns/op	       801.6 tokens/s
BenchmarkInference-32    	     321	  37217502 ns/op	       806.1 tokens/s
BenchmarkInference-32    	     321	  37291330 ns/op	       804.5 tokens/s
BenchmarkInference-32    	     320	  37351894 ns/op	       803.2 tokens/s
BenchmarkInference-32    	     320	  37421694 ns/op	       801.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	86.633s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/desktop-b8a29kd/vulkan1 -->

## Multimodal model benchmarks

The model is
[SmolVLM-256M-Instruct-Q8_0.gguf](https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/SmolVLM-256M-Instruct-Q8_0.gguf)
with its
[projector](https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf).
The code is [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go).

<!-- yzma:bench table multimodal -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | - | 289.5 | b10964 | 2026-09-23 |
| CUDA | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | CUDA0 | 1787.0 | b10964 | 2026-09-23 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan0 | 479.2 | b10964 | 2026-09-23 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan1 | 2018.0 | b10964 | 2026-09-23 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/desktop-b8a29kd -->
### CPU, amd64, AMD Ryzen 9 7950X 16-Core Processor            
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"desktop-b8a29kd","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":289.5,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

AMD Ryzen 9 7950X 16-Core Processor. 289.5 tokens a second.

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	      16	 745796738 ns/op	       315.2 tokens/s
BenchmarkMultimodalInference-32    	      22	 909492400 ns/op	       265.5 tokens/s
BenchmarkMultimodalInference-32    	      18	 923854639 ns/op	       262.5 tokens/s
BenchmarkMultimodalInference-32    	      14	 789683586 ns/op	       299.5 tokens/s
BenchmarkMultimodalInference-32    	      15	 821971773 ns/op	       289.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	72.818s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/desktop-b8a29kd -->

<!-- yzma:bench start multimodal/cuda/amd64/desktop-b8a29kd/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X 16-Core Processor            , CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"desktop-b8a29kd","device":"CUDA0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":1787,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

AMD Ryzen 9 7950X 16-Core Processor. 1787.0 tokens a second.

<details><summary>The device</summary>

```
Wed Sep 23 19:24:55 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 591.86                 Driver Version: 591.86         CUDA Version: 13.1     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   39C    P0              9W /  240W |       0MiB /   8192MiB |     10%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|  No running processes found                                                             |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=CUDA0
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	      90	 133313998 ns/op	      1767 tokens/s
BenchmarkMultimodalInference-32    	      90	 131427786 ns/op	      1788 tokens/s
BenchmarkMultimodalInference-32    	      94	 132823227 ns/op	      1779 tokens/s
BenchmarkMultimodalInference-32    	      94	 131777089 ns/op	      1787 tokens/s
BenchmarkMultimodalInference-32    	      97	 129325838 ns/op	      1807 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.993s
```

</details>
<!-- yzma:bench end multimodal/cuda/amd64/desktop-b8a29kd/cuda0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/desktop-b8a29kd/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":479.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

AMD Ryzen 9 7950X 16-Core Processor. 479.2 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.321


Instance Extensions: count = 20
-------------------------------
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_surface_maintenance1            : extension revision 1
VK_EXT_swapchain_colorspace            : extension revision 5
VK_KHR_device_group_creation           : extension revision 1
VK_KHR_display                         : extension revision 23
VK_KHR_external_fence_capabilities     : extension revision 1
VK_KHR_external_memory_capabilities    : extension revision 1
VK_KHR_external_semaphore_capabilities : extension revision 1
VK_KHR_get_display_properties2         : extension revision 1
VK_KHR_get_physical_device_properties2 : extension revision 2
VK_KHR_get_surface_capabilities2       : extension revision 1
VK_KHR_portability_enumeration         : extension revision 1
VK_KHR_surface                         : extension revision 25
VK_KHR_surface_maintenance1            : extension revision 1
VK_KHR_surface_protected_capabilities  : extension revision 1
VK_KHR_win32_surface                   : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1
VK_NV_external_memory_capabilities     : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_AMD_switchable_graphics AMD switchable graphics layer 1.4.315  version 1
VK_LAYER_NV_optimus              NVIDIA Optimus layer          1.4.325  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer     1.4.325  version 1
VK_LAYER_VALVE_steam_fossilize   Steam Pipeline Caching Layer  1.4.303  version 1
VK_LAYER_VALVE_steam_overlay     Steam Overlay Layer           1.3.207  version 1

Devices:
========
GPU0:
	apiVersion         = 1.4.315
	driverVersion      = 2.0.353
	vendorID           = 0x1002
	deviceID           = 0x164e
	deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
	deviceName         = AMD Radeon(TM) Graphics
	driverID           = DRIVER_ID_AMD_PROPRIETARY
	driverName         = AMD proprietary driver
	driverInfo         = 26.3.1 (AMD proprietary shader compiler)
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-0c00-0000-0000-000000000000
	driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
	apiVersion         = 1.4.325
	driverVersion      = 591.86.0.0
	vendorID           = 0x10de
	deviceID           = 0x2488
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 3070
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 591.86
	conformanceVersion = 1.4.3.0
	deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
	driverUUID         = ab24b0bb-7bd9-59fc-a310-4c2bdaba9b72
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=Vulkan0
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	      24	 478238429 ns/op	       484.8 tokens/s
BenchmarkMultimodalInference-32    	      26	 465214100 ns/op	       497.9 tokens/s
BenchmarkMultimodalInference-32    	      31	 480689284 ns/op	       479.2 tokens/s
BenchmarkMultimodalInference-32    	      22	 506849014 ns/op	       454.9 tokens/s
BenchmarkMultimodalInference-32    	      22	 504452005 ns/op	       460.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.883s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/desktop-b8a29kd/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan1","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":2018,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

AMD Ryzen 9 7950X 16-Core Processor. 2018.0 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.321


Instance Extensions: count = 20
-------------------------------
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_surface_maintenance1            : extension revision 1
VK_EXT_swapchain_colorspace            : extension revision 5
VK_KHR_device_group_creation           : extension revision 1
VK_KHR_display                         : extension revision 23
VK_KHR_external_fence_capabilities     : extension revision 1
VK_KHR_external_memory_capabilities    : extension revision 1
VK_KHR_external_semaphore_capabilities : extension revision 1
VK_KHR_get_display_properties2         : extension revision 1
VK_KHR_get_physical_device_properties2 : extension revision 2
VK_KHR_get_surface_capabilities2       : extension revision 1
VK_KHR_portability_enumeration         : extension revision 1
VK_KHR_surface                         : extension revision 25
VK_KHR_surface_maintenance1            : extension revision 1
VK_KHR_surface_protected_capabilities  : extension revision 1
VK_KHR_win32_surface                   : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1
VK_NV_external_memory_capabilities     : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_AMD_switchable_graphics AMD switchable graphics layer 1.4.315  version 1
VK_LAYER_NV_optimus              NVIDIA Optimus layer          1.4.325  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer     1.4.325  version 1
VK_LAYER_VALVE_steam_fossilize   Steam Pipeline Caching Layer  1.4.303  version 1
VK_LAYER_VALVE_steam_overlay     Steam Overlay Layer           1.3.207  version 1

Devices:
========
GPU0:
	apiVersion         = 1.4.315
	driverVersion      = 2.0.353
	vendorID           = 0x1002
	deviceID           = 0x164e
	deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
	deviceName         = AMD Radeon(TM) Graphics
	driverID           = DRIVER_ID_AMD_PROPRIETARY
	driverName         = AMD proprietary driver
	driverInfo         = 26.3.1 (AMD proprietary shader compiler)
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-0c00-0000-0000-000000000000
	driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
	apiVersion         = 1.4.325
	driverVersion      = 591.86.0.0
	vendorID           = 0x10de
	deviceID           = 0x2488
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 3070
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 591.86
	conformanceVersion = 1.4.3.0
	deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
	driverUUID         = ab24b0bb-7bd9-59fc-a310-4c2bdaba9b72
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=Vulkan1
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	     102	 116508174 ns/op	      2009 tokens/s
BenchmarkMultimodalInference-32    	     100	 115528723 ns/op	      2022 tokens/s
BenchmarkMultimodalInference-32    	     100	 115872453 ns/op	      2018 tokens/s
BenchmarkMultimodalInference-32    	     100	 116910499 ns/op	      2007 tokens/s
BenchmarkMultimodalInference-32    	      99	 114065233 ns/op	      2041 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	67.089s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
