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
| CPU | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | - | 112.9 | b11146 | 2026-09-24 |
| CUDA | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | CUDA0 | 700.5 | b11146 | 2026-09-24 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan0 | 115.6 | b11146 | 2026-09-24 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan1 | 802.3 | b11146 | 2026-09-24 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/desktop-b8a29kd -->
### CPU, amd64, AMD Ryzen 9 7950X 16-Core Processor            
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"desktop-b8a29kd","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":112.9,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD Ryzen 9 7950X 16-Core Processor. 112.9 tokens a second.

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	      44	 266535414 ns/op	       112.6 tokens/s
BenchmarkInference-32    	      44	 265157861 ns/op	       113.1 tokens/s
BenchmarkInference-32    	      44	 265621375 ns/op	       112.9 tokens/s
BenchmarkInference-32    	      45	 268107333 ns/op	       111.9 tokens/s
BenchmarkInference-32    	      44	 261941895 ns/op	       114.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	59.764s
```

</details>
<!-- yzma:bench end text/cpu/amd64/desktop-b8a29kd -->

<!-- yzma:bench start text/cuda/amd64/desktop-b8a29kd/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X 16-Core Processor            , CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"desktop-b8a29kd","device":"CUDA0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":700.5,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD Ryzen 9 7950X 16-Core Processor. 700.5 tokens a second.

<details><summary>The device</summary>

```
Thu Sep 24 10:17:50 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 591.86                 Driver Version: 591.86         CUDA Version: 13.1     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   38C    P0              4W /  240W |       0MiB /   8192MiB |      0%      Default |
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
BenchmarkInference-32    	     273	  43194863 ns/op	       694.5 tokens/s
BenchmarkInference-32    	     280	  42824033 ns/op	       700.5 tokens/s
BenchmarkInference-32    	     279	  42895518 ns/op	       699.4 tokens/s
BenchmarkInference-32    	     282	  42513730 ns/op	       705.7 tokens/s
BenchmarkInference-32    	     280	  42755701 ns/op	       701.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.823s
```

</details>
<!-- yzma:bench end text/cuda/amd64/desktop-b8a29kd/cuda0 -->

<!-- yzma:bench start text/vulkan/amd64/desktop-b8a29kd/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":115.6,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD Ryzen 9 7950X 16-Core Processor. 115.6 tokens a second.

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
BenchmarkInference-32    	      39	 288354731 ns/op	       104.0 tokens/s
BenchmarkInference-32    	      45	 259580356 ns/op	       115.6 tokens/s
BenchmarkInference-32    	      45	 259779556 ns/op	       115.5 tokens/s
BenchmarkInference-32    	      45	 259522553 ns/op	       115.6 tokens/s
BenchmarkInference-32    	      45	 259512469 ns/op	       115.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.030s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/desktop-b8a29kd/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan1","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":802.3,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD Ryzen 9 7950X 16-Core Processor. 802.3 tokens a second.

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
BenchmarkInference-32    	     307	  37952148 ns/op	       790.5 tokens/s
BenchmarkInference-32    	     319	  37391985 ns/op	       802.3 tokens/s
BenchmarkInference-32    	     321	  37283221 ns/op	       804.7 tokens/s
BenchmarkInference-32    	     321	  37368861 ns/op	       802.8 tokens/s
BenchmarkInference-32    	     320	  37430941 ns/op	       801.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	86.223s
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
| CPU | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | - | 291.9 | b11146 | 2026-09-24 |
| CUDA | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | CUDA0 | 1782.0 | b11146 | 2026-09-24 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan0 | 448.9 | b11146 | 2026-09-24 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan1 | 2011.0 | b11146 | 2026-09-24 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/desktop-b8a29kd -->
### CPU, amd64, AMD Ryzen 9 7950X 16-Core Processor            
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"desktop-b8a29kd","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":291.9,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD Ryzen 9 7950X 16-Core Processor. 291.9 tokens a second.

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	      16	 683603094 ns/op	       339.9 tokens/s
BenchmarkMultimodalInference-32    	      13	 831616277 ns/op	       286.7 tokens/s
BenchmarkMultimodalInference-32    	      24	 812582217 ns/op	       291.9 tokens/s
BenchmarkMultimodalInference-32    	      16	 825620562 ns/op	       288.0 tokens/s
BenchmarkMultimodalInference-32    	      13	 774315915 ns/op	       304.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	65.331s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/desktop-b8a29kd -->

<!-- yzma:bench start multimodal/cuda/amd64/desktop-b8a29kd/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X 16-Core Processor            , CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"desktop-b8a29kd","device":"CUDA0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":1782,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD Ryzen 9 7950X 16-Core Processor. 1782.0 tokens a second.

<details><summary>The device</summary>

```
Thu Sep 24 10:17:50 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 591.86                 Driver Version: 591.86         CUDA Version: 13.1     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   38C    P0              4W /  240W |       0MiB /   8192MiB |      0%      Default |
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
BenchmarkMultimodalInference-32    	      90	 133344401 ns/op	      1763 tokens/s
BenchmarkMultimodalInference-32    	      93	 131132820 ns/op	      1790 tokens/s
BenchmarkMultimodalInference-32    	      87	 132017324 ns/op	      1782 tokens/s
BenchmarkMultimodalInference-32    	      87	 132366326 ns/op	      1776 tokens/s
BenchmarkMultimodalInference-32    	      94	 131995491 ns/op	      1782 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	61.563s
```

</details>
<!-- yzma:bench end multimodal/cuda/amd64/desktop-b8a29kd/cuda0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/desktop-b8a29kd/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":448.9,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD Ryzen 9 7950X 16-Core Processor. 448.9 tokens a second.

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
BenchmarkMultimodalInference-32    	       8	1290016338 ns/op	       179.2 tokens/s
BenchmarkMultimodalInference-32    	      25	 528152336 ns/op	       440.9 tokens/s
BenchmarkMultimodalInference-32    	      24	 513449204 ns/op	       453.6 tokens/s
BenchmarkMultimodalInference-32    	      24	 497500900 ns/op	       465.2 tokens/s
BenchmarkMultimodalInference-32    	      20	 519748675 ns/op	       448.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	60.078s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/desktop-b8a29kd/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan1","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":2011,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD Ryzen 9 7950X 16-Core Processor. 2011.0 tokens a second.

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
BenchmarkMultimodalInference-32    	      85	 123081724 ns/op	      1910 tokens/s
BenchmarkMultimodalInference-32    	     100	 116428688 ns/op	      2012 tokens/s
BenchmarkMultimodalInference-32    	     100	 118117676 ns/op	      1992 tokens/s
BenchmarkMultimodalInference-32    	      99	 116481655 ns/op	      2011 tokens/s
BenchmarkMultimodalInference-32    	      97	 113441393 ns/op	      2046 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	65.135s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
