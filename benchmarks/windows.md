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
| CPU | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | - | 113.2 | b10964 | 2026-09-17 |
| CUDA | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | CUDA0 | 701.3 | b10964 | 2026-09-17 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan0 | 107.6 | b10964 | 2026-09-17 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan1 | 801.1 | b10964 | 2026-09-17 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/desktop-b8a29kd -->
### CPU, amd64, AMD Ryzen 9 7950X 16-Core Processor            
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"desktop-b8a29kd","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":113.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

AMD Ryzen 9 7950X 16-Core Processor. 113.2 tokens a second.

<details><summary>The output of go test</summary>

```
> cd pkg/llama; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkInference-32    	      44	 264061389 ns/op	       113.6 tokens/s
BenchmarkInference-32    	      44	 264901464 ns/op	       113.2 tokens/s
BenchmarkInference-32    	      44	 265594468 ns/op	       113.0 tokens/s
BenchmarkInference-32    	      45	 263602856 ns/op	       113.8 tokens/s
BenchmarkInference-32    	      44	 265617307 ns/op	       112.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	59.553s
```

</details>
<!-- yzma:bench end text/cpu/amd64/desktop-b8a29kd -->

<!-- yzma:bench start text/cuda/amd64/desktop-b8a29kd/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X 16-Core Processor            , CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"desktop-b8a29kd","device":"CUDA0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":701.3,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

AMD Ryzen 9 7950X 16-Core Processor. 701.3 tokens a second.

<details><summary>The device</summary>

```
Thu Sep 17 17:17:52 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 591.86                 Driver Version: 591.86         CUDA Version: 13.1     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   38C    P8              7W /  240W |       0MiB /   8192MiB |      0%      Default |
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
BenchmarkInference-32    	     278	  42802736 ns/op	       700.9 tokens/s
BenchmarkInference-32    	     280	  42734043 ns/op	       702.0 tokens/s
BenchmarkInference-32    	     280	  42704964 ns/op	       702.5 tokens/s
BenchmarkInference-32    	     279	  42779244 ns/op	       701.3 tokens/s
BenchmarkInference-32    	     279	  42787771 ns/op	       701.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.695s
```

</details>
<!-- yzma:bench end text/cuda/amd64/desktop-b8a29kd/cuda0 -->

<!-- yzma:bench start text/vulkan/amd64/desktop-b8a29kd/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":107.6,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

AMD Ryzen 9 7950X 16-Core Processor. 107.6 tokens a second.

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

Instance Layers: count = 7
--------------------------
VK_LAYER_AMD_switchable_graphics AMD switchable graphics layer                 1.4.315  version 1
VK_LAYER_EOS_Overlay             Vulkan overlay layer for Epic Online Services 1.2.136  version 1
VK_LAYER_EOS_Overlay             Vulkan overlay layer for Epic Online Services 1.2.136  version 1
VK_LAYER_NV_optimus              NVIDIA Optimus layer                          1.4.325  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer                     1.4.325  version 1
VK_LAYER_VALVE_steam_fossilize   Steam Pipeline Caching Layer                  1.4.303  version 1
VK_LAYER_VALVE_steam_overlay     Steam Overlay Layer                           1.3.207  version 1

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
BenchmarkInference-32    	      37	 309842411 ns/op	        96.82 tokens/s
BenchmarkInference-32    	      42	 278712176 ns/op	       107.6 tokens/s
BenchmarkInference-32    	      42	 278377079 ns/op	       107.8 tokens/s
BenchmarkInference-32    	      42	 278511669 ns/op	       107.7 tokens/s
BenchmarkInference-32    	      42	 279475457 ns/op	       107.3 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.488s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/desktop-b8a29kd/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan1","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":801.1,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

AMD Ryzen 9 7950X 16-Core Processor. 801.1 tokens a second.

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

Instance Layers: count = 7
--------------------------
VK_LAYER_AMD_switchable_graphics AMD switchable graphics layer                 1.4.315  version 1
VK_LAYER_EOS_Overlay             Vulkan overlay layer for Epic Online Services 1.2.136  version 1
VK_LAYER_EOS_Overlay             Vulkan overlay layer for Epic Online Services 1.2.136  version 1
VK_LAYER_NV_optimus              NVIDIA Optimus layer                          1.4.325  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer                     1.4.325  version 1
VK_LAYER_VALVE_steam_fossilize   Steam Pipeline Caching Layer                  1.4.303  version 1
VK_LAYER_VALVE_steam_overlay     Steam Overlay Layer                           1.3.207  version 1

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
BenchmarkInference-32    	     243	  46302602 ns/op	       647.9 tokens/s
BenchmarkInference-32    	     320	  37450093 ns/op	       801.1 tokens/s
BenchmarkInference-32    	     319	  37483480 ns/op	       800.4 tokens/s
BenchmarkInference-32    	     319	  37439732 ns/op	       801.3 tokens/s
BenchmarkInference-32    	     321	  37285423 ns/op	       804.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	84.888s
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
| CPU | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | - | 410.5 | b10964 | 2026-09-17 |
| CUDA | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | CUDA0 | 1772.0 | b10964 | 2026-09-17 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan0 | 447.9 | b10964 | 2026-09-17 |
| Vulkan | amd64 | AMD Ryzen 9 7950X 16-Core Processor             | Vulkan1 | 2032.0 | b10964 | 2026-09-17 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/desktop-b8a29kd -->
### CPU, amd64, AMD Ryzen 9 7950X 16-Core Processor            
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"desktop-b8a29kd","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":410.5,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

AMD Ryzen 9 7950X 16-Core Processor. 410.5 tokens a second.

<details><summary>The output of go test</summary>

```
> cd pkg/mtmd; go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: windows
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD Ryzen 9 7950X 16-Core Processor            
BenchmarkMultimodalInference-32    	      18	 566609167 ns/op	       421.7 tokens/s
BenchmarkMultimodalInference-32    	      18	 572250883 ns/op	       410.5 tokens/s
BenchmarkMultimodalInference-32    	      21	 548104448 ns/op	       427.2 tokens/s
BenchmarkMultimodalInference-32    	      27	 601168922 ns/op	       396.1 tokens/s
BenchmarkMultimodalInference-32    	      18	 648956728 ns/op	       373.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	60.803s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/desktop-b8a29kd -->

<!-- yzma:bench start multimodal/cuda/amd64/desktop-b8a29kd/cuda0 -->
### CUDA, amd64, AMD Ryzen 9 7950X 16-Core Processor            , CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"desktop-b8a29kd","device":"CUDA0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":1772,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

AMD Ryzen 9 7950X 16-Core Processor. 1772.0 tokens a second.

<details><summary>The device</summary>

```
Thu Sep 17 17:17:52 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 591.86                 Driver Version: 591.86         CUDA Version: 13.1     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                  Driver-Model | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 3070      WDDM  |   00000000:01:00.0 Off |                  N/A |
|  0%   38C    P8              7W /  240W |       0MiB /   8192MiB |      0%      Default |
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
BenchmarkMultimodalInference-32    	      87	 134776951 ns/op	      1746 tokens/s
BenchmarkMultimodalInference-32    	      87	 132460561 ns/op	      1772 tokens/s
BenchmarkMultimodalInference-32    	      97	 130908414 ns/op	      1785 tokens/s
BenchmarkMultimodalInference-32    	      97	 134073011 ns/op	      1758 tokens/s
BenchmarkMultimodalInference-32    	     100	 132518141 ns/op	      1772 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	63.650s
```

</details>
<!-- yzma:bench end multimodal/cuda/amd64/desktop-b8a29kd/cuda0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/desktop-b8a29kd/vulkan0 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan0","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":447.9,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

AMD Ryzen 9 7950X 16-Core Processor. 447.9 tokens a second.

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

Instance Layers: count = 7
--------------------------
VK_LAYER_AMD_switchable_graphics AMD switchable graphics layer                 1.4.315  version 1
VK_LAYER_EOS_Overlay             Vulkan overlay layer for Epic Online Services 1.2.136  version 1
VK_LAYER_EOS_Overlay             Vulkan overlay layer for Epic Online Services 1.2.136  version 1
VK_LAYER_NV_optimus              NVIDIA Optimus layer                          1.4.325  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer                     1.4.325  version 1
VK_LAYER_VALVE_steam_fossilize   Steam Pipeline Caching Layer                  1.4.303  version 1
VK_LAYER_VALVE_steam_overlay     Steam Overlay Layer                           1.3.207  version 1

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
BenchmarkMultimodalInference-32    	       8	1286016100 ns/op	       181.8 tokens/s
BenchmarkMultimodalInference-32    	      27	 543229837 ns/op	       429.1 tokens/s
BenchmarkMultimodalInference-32    	      20	 508264755 ns/op	       457.6 tokens/s
BenchmarkMultimodalInference-32    	      22	 513704268 ns/op	       447.9 tokens/s
BenchmarkMultimodalInference-32    	      28	 505875479 ns/op	       457.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.603s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/desktop-b8a29kd/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
### Vulkan, amd64, AMD Ryzen 9 7950X 16-Core Processor            , Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"desktop-b8a29kd","device":"Vulkan1","label":"AMD Ryzen 9 7950X 16-Core Processor            ","cpu":"AMD Ryzen 9 7950X 16-Core Processor","tokens_per_second":2032,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

AMD Ryzen 9 7950X 16-Core Processor. 2032.0 tokens a second.

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

Instance Layers: count = 7
--------------------------
VK_LAYER_AMD_switchable_graphics AMD switchable graphics layer                 1.4.315  version 1
VK_LAYER_EOS_Overlay             Vulkan overlay layer for Epic Online Services 1.2.136  version 1
VK_LAYER_EOS_Overlay             Vulkan overlay layer for Epic Online Services 1.2.136  version 1
VK_LAYER_NV_optimus              NVIDIA Optimus layer                          1.4.325  version 1
VK_LAYER_NV_present              NVIDIA Presentation Layer                     1.4.325  version 1
VK_LAYER_VALVE_steam_fossilize   Steam Pipeline Caching Layer                  1.4.303  version 1
VK_LAYER_VALVE_steam_overlay     Steam Overlay Layer                           1.3.207  version 1

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
BenchmarkMultimodalInference-32    	      76	 134298692 ns/op	      1738 tokens/s
BenchmarkMultimodalInference-32    	     103	 114743197 ns/op	      2032 tokens/s
BenchmarkMultimodalInference-32    	      90	 113817927 ns/op	      2042 tokens/s
BenchmarkMultimodalInference-32    	      85	 120897659 ns/op	      1962 tokens/s
BenchmarkMultimodalInference-32    	      93	 114613214 ns/op	      2034 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	61.461s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/desktop-b8a29kd/vulkan1 -->
