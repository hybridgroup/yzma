# Windows benchmarks

Benchmarks of yzma on Windows. Each table shows the median of five runs. The output
of each run, and the device details, are below the tables.

To add a machine or to regenerate these numbers, see
[how to run the benchmarks](README.md).

## Summary

Tokens a second on each machine. The GPU columns show the fastest GPU backend.

| Machine | GPU | Text, CPU | Text, GPU | Multimodal, CPU | Multimodal, GPU |
| --- | --- | --- | --- | --- | --- |
| AMD Ryzen 9 7950X | RTX 3070 | 112.9 | 802.3, Vulkan | 291.9 | 2011.0, Vulkan |

- On the RTX 3070, Vulkan is faster than CUDA. It is 15 percent faster for text
  and 13 percent faster for multimodal.
- The integrated AMD Radeon GPU on Vulkan0 reaches 115.6 for text, which is
  almost the same as the CPU. For multimodal it reaches 448.9, which is faster
  than the CPU.
- The CPU results vary a lot from run to run, because Windows can't pin a
  thread to a core. The multimodal runs on the CPU go from 286.7 to
  339.9.

## Text model benchmarks

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf).
The code is [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go).
The benchmark uses 4 threads on each machine, see
[the thread count](README.md#run-them).

<!-- yzma:bench table text -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | - | 110.7 | b11669 | 2026-10-19 |
| CUDA | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | CUDA0 | 703.3 | b11669 | 2026-10-19 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan0 | 114.3 | b11669 | 2026-10-19 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan1 | 820.5 | b11669 | 2026-10-19 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/desktop-b8a29kd -->
### CPU, amd64, AMD Ryzen 9 7950X 16-Core Processor            
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"desktop-b8a29kd","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":110.7,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-19"} -->

AMD Ryzen 9 7950X 16-Core Processor. 110.7 tokens a second.

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	      42	 275651569 ns/op	       108.8 tokens/s
BenchmarkInference-32    	      43	 267559830 ns/op	       112.1 tokens/s
BenchmarkInference-32    	      43	 270985609 ns/op	       110.7 tokens/s
BenchmarkInference-32    	      42	 274921905 ns/op	       109.1 tokens/s
BenchmarkInference-32    	      43	 270224542 ns/op	       111.0 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	59.806s
```

</details>
<!-- yzma:bench end text/cpu/amd64/desktop-b8a29kd -->

<!-- yzma:bench start text/cuda/amd64/desktop-b8a29kd/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X 16-Core Processor            , CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"desktop-b8a29kd","device":"CUDA0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":703.3,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-19"} -->

AMD Ryzen 9 7950X 16-Core Processor. 703.3 tokens a second.

<details><summary>The device</summary>

```
Mon Oct 19 12:47:21 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 610.88                 KMD Version: 610.88        CUDA UMD Version: 13.3     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   27C    P0              6W /  240W |       0MiB /   8192MiB |     19%      Default |
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
BenchmarkInference-32    	     278	  42679360 ns/op	       702.9 tokens/s
BenchmarkInference-32    	     282	  42425318 ns/op	       707.1 tokens/s
BenchmarkInference-32    	     282	  42572143 ns/op	       704.7 tokens/s
BenchmarkInference-32    	     280	  42676884 ns/op	       703.0 tokens/s
BenchmarkInference-32    	     280	  42656456 ns/op	       703.3 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.934s
```

</details>
<!-- yzma:bench end text/cuda/amd64/desktop-b8a29kd/cuda0 -->

<!-- yzma:bench start text/vulkan/amd64/desktop-b8a29kd/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":114.3,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-19"} -->

AMD Ryzen 9 7950X 16-Core Processor. 114.3 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.341


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
VK_LAYER_NV_optimus              NVIDIA Optimus layer          1.4.341  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer     1.4.341  version 1
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
	driverInfo         = 2026.04.17 (AMD proprietary shader compiler)
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-0c00-0000-0000-000000000000
	driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
	apiVersion         = 1.4.341
	driverVersion      = 610.88.0.0
	vendorID           = 0x10de
	deviceID           = 0x2488
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 3070
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 610.88
	conformanceVersion = 1.4.3.3
	deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
	driverUUID         = 8690f1c8-0a3f-5499-9bf6-ea2aee515602
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=Vulkan0
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	      38	 293648632 ns/op	       102.2 tokens/s
BenchmarkInference-32    	      44	 262366432 ns/op	       114.3 tokens/s
BenchmarkInference-32    	      44	 263832305 ns/op	       113.7 tokens/s
BenchmarkInference-32    	      44	 261808439 ns/op	       114.6 tokens/s
BenchmarkInference-32    	      45	 260732080 ns/op	       115.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	60.770s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/desktop-b8a29kd/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan1","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":820.5,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-19"} -->

AMD Ryzen 9 7950X 16-Core Processor. 820.5 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.341


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
VK_LAYER_NV_optimus              NVIDIA Optimus layer          1.4.341  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer     1.4.341  version 1
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
	driverInfo         = 2026.04.17 (AMD proprietary shader compiler)
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-0c00-0000-0000-000000000000
	driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
	apiVersion         = 1.4.341
	driverVersion      = 610.88.0.0
	vendorID           = 0x10de
	deviceID           = 0x2488
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 3070
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 610.88
	conformanceVersion = 1.4.3.3
	deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
	driverUUID         = 8690f1c8-0a3f-5499-9bf6-ea2aee515602
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=Vulkan1
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	     244	  45706907 ns/op	       656.4 tokens/s
BenchmarkInference-32    	     327	  36632594 ns/op	       818.9 tokens/s
BenchmarkInference-32    	     326	  36562383 ns/op	       820.5 tokens/s
BenchmarkInference-32    	     331	  36183650 ns/op	       829.1 tokens/s
BenchmarkInference-32    	     331	  36271350 ns/op	       827.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	85.415s
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
| CPU | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | - | 460.4 | b11669 | 2026-10-19 |
| CUDA | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | CUDA0 | 1848.0 | b11669 | 2026-10-19 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan0 | 446.7 | b11669 | 2026-10-19 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan1 | 2093.0 | b11669 | 2026-10-19 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/desktop-b8a29kd -->
### CPU, amd64, AMD Ryzen 9 7950X 16-Core Processor            
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"desktop-b8a29kd","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":460.4,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-19"} -->

AMD Ryzen 9 7950X 16-Core Processor. 460.4 tokens a second.

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	      22	 530597977 ns/op	       450.1 tokens/s
BenchmarkMultimodalInference-32    	      21	 515711986 ns/op	       460.4 tokens/s
BenchmarkMultimodalInference-32    	      19	 536271721 ns/op	       448.5 tokens/s
BenchmarkMultimodalInference-32    	      25	 509296628 ns/op	       465.4 tokens/s
BenchmarkMultimodalInference-32    	      27	 498370763 ns/op	       472.8 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	60.006s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/desktop-b8a29kd -->

<!-- yzma:bench start multimodal/cuda/amd64/desktop-b8a29kd/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X 16-Core Processor            , CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"desktop-b8a29kd","device":"CUDA0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":1848,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-19"} -->

AMD Ryzen 9 7950X 16-Core Processor. 1848.0 tokens a second.

<details><summary>The device</summary>

```
Mon Oct 19 12:47:21 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 610.88                 KMD Version: 610.88        CUDA UMD Version: 13.3     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   27C    P0              6W /  240W |       0MiB /   8192MiB |     19%      Default |
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
BenchmarkMultimodalInference-32    	      91	 128573067 ns/op	      1830 tokens/s
BenchmarkMultimodalInference-32    	      82	 127622978 ns/op	      1844 tokens/s
BenchmarkMultimodalInference-32    	      88	 127211835 ns/op	      1848 tokens/s
BenchmarkMultimodalInference-32    	      93	 126958076 ns/op	      1851 tokens/s
BenchmarkMultimodalInference-32    	     100	 124871678 ns/op	      1870 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	59.553s
```

</details>
<!-- yzma:bench end multimodal/cuda/amd64/desktop-b8a29kd/cuda0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/desktop-b8a29kd/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":446.7,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-19"} -->

AMD Ryzen 9 7950X 16-Core Processor. 446.7 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.341


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
VK_LAYER_NV_optimus              NVIDIA Optimus layer          1.4.341  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer     1.4.341  version 1
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
	driverInfo         = 2026.04.17 (AMD proprietary shader compiler)
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-0c00-0000-0000-000000000000
	driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
	apiVersion         = 1.4.341
	driverVersion      = 610.88.0.0
	vendorID           = 0x10de
	deviceID           = 0x2488
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 3070
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 610.88
	conformanceVersion = 1.4.3.3
	deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
	driverUUID         = 8690f1c8-0a3f-5499-9bf6-ea2aee515602
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=Vulkan0
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	      10	1025074650 ns/op	       224.7 tokens/s
BenchmarkMultimodalInference-32    	      31	 458625958 ns/op	       499.2 tokens/s
BenchmarkMultimodalInference-32    	      26	 504358562 ns/op	       460.7 tokens/s
BenchmarkMultimodalInference-32    	      26	 518228762 ns/op	       444.6 tokens/s
BenchmarkMultimodalInference-32    	      22	 516416677 ns/op	       446.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	64.895s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/desktop-b8a29kd/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan1","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":2093,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-19"} -->

AMD Ryzen 9 7950X 16-Core Processor. 2093.0 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.341


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
VK_LAYER_NV_optimus              NVIDIA Optimus layer          1.4.341  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer     1.4.341  version 1
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
	driverInfo         = 2026.04.17 (AMD proprietary shader compiler)
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-0c00-0000-0000-000000000000
	driverUUID         = 414d442d-5749-4e2d-4452-560000000000
GPU1:
	apiVersion         = 1.4.341
	driverVersion      = 610.88.0.0
	vendorID           = 0x10de
	deviceID           = 0x2488
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 3070
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 610.88
	conformanceVersion = 1.4.3.3
	deviceUUID         = 91c0b9f4-e340-3c73-1422-c227930ae260
	driverUUID         = 8690f1c8-0a3f-5499-9bf6-ea2aee515602
```

</details>

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=Vulkan1
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	      80	 135170000 ns/op	      1719 tokens/s
BenchmarkMultimodalInference-32    	     100	 112391238 ns/op	      2089 tokens/s
BenchmarkMultimodalInference-32    	      97	 111984727 ns/op	      2093 tokens/s
BenchmarkMultimodalInference-32    	     100	 110464707 ns/op	      2112 tokens/s
BenchmarkMultimodalInference-32    	     100	 110957981 ns/op	      2101 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	63.871s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
