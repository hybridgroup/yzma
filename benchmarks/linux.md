# Linux benchmarks

Benchmarks of yzma on Linux. Each table shows the median of five runs. The output
of each run, and the device details, are below the tables.

To add a machine or to regenerate these numbers, see
[how to run the benchmarks](README.md).

## Summary

Tokens a second on each machine. The GPU columns show the fastest GPU backend.

| Machine | GPU | Text, CPU | Text, GPU | Multimodal, CPU | Multimodal, GPU |
| --- | --- | --- | --- | --- | --- |
| Intel Core i9-13900HX | RTX 4070 Laptop | 251.4 | 846.1, CUDA | 842.6 | 2448.0, CUDA |
| NVIDIA Jetson Orin Nano Super | Tegra Orin | 84.1 | 190.0, CUDA | 208.3 | 427.3, CUDA |
| Raspberry Pi 4 Model B | none | 35.4 | none | 5.6 | none |
| Arduino UNO Q | none | 32.2 | none | 4.1 | none |

- CUDA on the RTX 4070 is the fastest of all the Linux machines.
- Vulkan on the same RTX 4070 reaches 88 percent of CUDA for text and 89 percent
  for multimodal.
- The integrated Intel GPU on the i9-13900HX is slower than its CPU. On Vulkan0
  it reaches 94.9 for text and 395.7 for multimodal.
- On the Jetson Orin Nano, the GPU is 2.3 times faster than the CPU for text and
  2.1 times faster for multimodal. Vulkan reaches 90 percent of CUDA for text and
  99 percent for multimodal.
- The Raspberry Pi 4 and the Arduino UNO Q have no GPU backend. They run text at
  more than 30 tokens a second, but multimodal is slow.
- The text suite uses 4 threads on every machine. The multimodal suite uses one
  thread per performance core, so its CPU column reflects the size of the
  processor.

## Text model benchmarks

The model is
[SmolLM-135M.Q2_K.gguf](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf).
The code is [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go).
The benchmark uses 4 threads on each machine, see
[the thread count](README.md#run-them).

<!-- yzma:bench table text -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | amd64 | AMD EPYC 7443P 24-Core Processor | - | 253.4 | b11146 | 2026-09-25 |
| CPU | amd64 | Intel Core i9-13900HX | - | 251.4 | b11669 | 2026-10-16 |
| CPU | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | - | 80.8 | b11669 | 2026-10-17 |
| CPU | arm64 | Raspberry Pi 4 Model B Rev 1.4 | - | 35.4 | b11146 | 2026-09-24 |
| CPU | arm64 | Arduino UnoQ | - | 32.2 | b11146 | 2026-09-24 |
| CUDA | amd64 | Intel Core i9-13900HX | CUDA0 | 846.1 | b11669 | 2026-10-16 |
| CUDA | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | CUDA0 | 188.3 | b11669 | 2026-10-17 |
| ROCm | amd64 | AMD EPYC 7443P 24-Core Processor | ROCm0 | 479.9 | b11146 | 2026-09-25 |
| ROCm | amd64 | AMD EPYC 7443P 24-Core Processor | ROCm1 | 475.0 | b11146 | 2026-09-25 |
| Vulkan | amd64 | AMD EPYC 7443P 24-Core Processor | Vulkan0 | 806.5 | b11146 | 2026-09-24 |
| Vulkan | amd64 | AMD EPYC 7443P 24-Core Processor | Vulkan1 | 791.5 | b11146 | 2026-09-24 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan0 | 94.9 | b11669 | 2026-10-16 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan1 | 747.5 | b11669 | 2026-10-16 |
| Vulkan | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | Vulkan0 | 171.2 | b11669 | 2026-10-17 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/cookie3 -->
### CPU, amd64, AMD EPYC 7443P 24-Core Processor
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"cookie3","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":253.4,"llamacpp":"b11146","yzma":"1.29.0-dev","date":"2026-09-25"} -->

AMD EPYC 7443P 24-Core Processor. 253.4 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkInference-48    	     135	  89289872 ns/op	       336.0 tokens/s
BenchmarkInference-48    	      96	 118404079 ns/op	       253.4 tokens/s
BenchmarkInference-48    	      90	 112516471 ns/op	       266.6 tokens/s
BenchmarkInference-48    	      98	 119620676 ns/op	       250.8 tokens/s
BenchmarkInference-48    	      91	 122816695 ns/op	       244.3 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.630s
```

</details>
<!-- yzma:bench end text/cpu/amd64/cookie3 -->

<!-- yzma:bench start text/cpu/amd64/i9-13900hx -->
### CPU, amd64, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":251.4,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 251.4 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	      98	 119313769 ns/op	       251.4 tokens/s
BenchmarkInference-32    	      98	 121814723 ns/op	       246.3 tokens/s
BenchmarkInference-32    	      98	 119311267 ns/op	       251.4 tokens/s
BenchmarkInference-32    	      94	 118855140 ns/op	       252.4 tokens/s
BenchmarkInference-32    	     100	 118508397 ns/op	       253.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	64.245s
```

</details>
<!-- yzma:bench end text/cpu/amd64/i9-13900hx -->

<!-- yzma:bench start text/cpu/arm64/localhost -->
### CPU, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"localhost","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":80.77,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-17"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-6   	      46	 371434655 ns/op	        80.77 tokens/s
BenchmarkInference-6   	      33	 367444826 ns/op	        81.64 tokens/s
BenchmarkInference-6   	      33	 368706636 ns/op	        81.37 tokens/s
BenchmarkInference-6   	      31	 388245038 ns/op	        77.27 tokens/s
BenchmarkInference-6   	      33	 379090984 ns/op	        79.14 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	68.222s
```

</details>
<!-- yzma:bench end text/cpu/arm64/localhost -->

<!-- yzma:bench start text/cpu/arm64/raspberrypi -->
### CPU, arm64, Raspberry Pi 4 Model B Rev 1.4
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"raspberrypi","label":"Raspberry Pi 4 Model B Rev 1.4","tokens_per_second":35.4,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-4   	      13	 845508176 ns/op	        35.48 tokens/s
BenchmarkInference-4   	      13	 849931030 ns/op	        35.30 tokens/s
BenchmarkInference-4   	      13	 854401991 ns/op	        35.11 tokens/s
BenchmarkInference-4   	      13	 847558967 ns/op	        35.40 tokens/s
BenchmarkInference-4   	      13	 846519455 ns/op	        35.44 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.117s
```

</details>
<!-- yzma:bench end text/cpu/arm64/raspberrypi -->

<!-- yzma:bench start text/cpu/arm64/yzma -->
### CPU, arm64, Arduino UnoQ
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"yzma","label":"Arduino UnoQ","tokens_per_second":32.21,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-4   	      12	 944310001 ns/op	        31.77 tokens/s
BenchmarkInference-4   	      12	 931419997 ns/op	        32.21 tokens/s
BenchmarkInference-4   	      12	 928250065 ns/op	        32.32 tokens/s
BenchmarkInference-4   	      12	 929622606 ns/op	        32.27 tokens/s
BenchmarkInference-4   	      12	 939298061 ns/op	        31.94 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	58.410s
```

</details>
<!-- yzma:bench end text/cpu/arm64/yzma -->

<!-- yzma:bench start text/cuda/amd64/i9-13900hx/cuda0 -->
### CUDA, amd64, Intel Core i9-13900HX, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"i9-13900hx","device":"CUDA0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":846.1,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 846.1 tokens a second.

<details><summary>The device</summary>

```
Fri Oct 16 13:56:21 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   48C    P8              7W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	     333	  35556203 ns/op	       843.7 tokens/s
BenchmarkInference-32    	     337	  35457226 ns/op	       846.1 tokens/s
BenchmarkInference-32    	     339	  35370692 ns/op	       848.2 tokens/s
BenchmarkInference-32    	     338	  35409305 ns/op	       847.2 tokens/s
BenchmarkInference-32    	     337	  35478668 ns/op	       845.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	65.480s
```

</details>
<!-- yzma:bench end text/cuda/amd64/i9-13900hx/cuda0 -->

<!-- yzma:bench start text/cuda/arm64/localhost/cuda0 -->
### CUDA, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"arm64","machine":"localhost","device":"CUDA0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":188.3,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-17"} -->

<details><summary>The device</summary>

```
Sat Oct 17 09:09:45 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.78                 Driver Version: 595.78         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  Orin (nvgpu)                  N/A  |   N/A              N/A |                  N/A |
| N/A   N/A  N/A             N/A  /  N/A  | Not Supported          |     N/A          N/A |
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
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=CUDA0
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-6   	      73	 160523505 ns/op	       186.9 tokens/s
BenchmarkInference-6   	      73	 159344113 ns/op	       188.3 tokens/s
BenchmarkInference-6   	      73	 156831722 ns/op	       191.3 tokens/s
BenchmarkInference-6   	      74	 159531321 ns/op	       188.1 tokens/s
BenchmarkInference-6   	      73	 158979317 ns/op	       188.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.593s
```

</details>
<!-- yzma:bench end text/cuda/arm64/localhost/cuda0 -->

<!-- yzma:bench start text/rocm/amd64/cookie3/rocm0 -->
### ROCm, amd64, AMD EPYC 7443P 24-Core Processor, ROCm0
<!-- yzma:bench meta {"suite":"text","backend":"rocm","arch":"amd64","machine":"cookie3","device":"ROCm0","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":479.9,"llamacpp":"b11146","yzma":"1.29.0-dev","date":"2026-09-25"} -->

AMD EPYC 7443P 24-Core Processor. 479.9 tokens a second.

<details><summary>The device</summary>

```


======================================= ROCm System Management Interface =======================================
================================================= Concise Info =================================================
Device  Node  IDs              Temp    Power  Partitions          SCLK  MCLK    Fan  Perf  PwrCap  VRAM%  GPU%  
[3m              (DID,     GUID)  (Edge)  (Avg)  (Mem, Compute, ID)                                                [0m
================================================================================================================
0       1     0x744c,   16626  31.0°C  59.0W  N/A, N/A, 0         0Mhz  456Mhz  0%   auto  0.0W    0%     0%    
1       2     0x744c,   19095  27.0°C  61.0W  N/A, N/A, 0         0Mhz  456Mhz  0%   auto  0.0W    0%     0%    
================================================================================================================
============================================= End of ROCm SMI Log ==============================================
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=ROCm0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkInference-48    	     189	  62374303 ns/op	       481.0 tokens/s
BenchmarkInference-48    	     192	  62524617 ns/op	       479.8 tokens/s
BenchmarkInference-48    	     192	  62219906 ns/op	       482.2 tokens/s
BenchmarkInference-48    	     190	  62511263 ns/op	       479.9 tokens/s
BenchmarkInference-48    	     190	  62777837 ns/op	       477.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.381s
```

</details>
<!-- yzma:bench end text/rocm/amd64/cookie3/rocm0 -->

<!-- yzma:bench start text/rocm/amd64/cookie3/rocm1 -->
### ROCm, amd64, AMD EPYC 7443P 24-Core Processor, ROCm1
<!-- yzma:bench meta {"suite":"text","backend":"rocm","arch":"amd64","machine":"cookie3","device":"ROCm1","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":475,"llamacpp":"b11146","yzma":"1.29.0-dev","date":"2026-09-25"} -->

AMD EPYC 7443P 24-Core Processor. 475.0 tokens a second.

<details><summary>The device</summary>

```


========================================= ROCm System Management Interface =========================================
=================================================== Concise Info ===================================================
Device  Node  IDs              Temp    Power   Partitions          SCLK  MCLK   Fan     Perf  PwrCap  VRAM%  GPU%  
[3m              (DID,     GUID)  (Edge)  (Avg)   (Mem, Compute, ID)                                                  [0m
====================================================================================================================
0       1     0x744c,   16626  46.0°C  139.0W  N/A, N/A, 0         0Mhz  96Mhz  28.63%  auto  0.0W    0%     0%    
1       2     0x744c,   19095  26.0°C  12.0W   N/A, N/A, 0         0Mhz  96Mhz  0%      auto  0.0W    0%     0%    
====================================================================================================================
=============================================== End of ROCm SMI Log ================================================
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=ROCm1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkInference-48    	     186	  63754700 ns/op	       470.6 tokens/s
BenchmarkInference-48    	     189	  63221721 ns/op	       474.5 tokens/s
BenchmarkInference-48    	     189	  62945370 ns/op	       476.6 tokens/s
BenchmarkInference-48    	     189	  63152136 ns/op	       475.0 tokens/s
BenchmarkInference-48    	     189	  62817984 ns/op	       477.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	60.890s
```

</details>
<!-- yzma:bench end text/rocm/amd64/cookie3/rocm1 -->

<!-- yzma:bench start text/vulkan/amd64/cookie3/vulkan0 -->
### Vulkan, amd64, AMD EPYC 7443P 24-Core Processor, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"cookie3","device":"Vulkan0","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":806.5,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD EPYC 7443P 24-Core Processor. 806.5 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.341


Instance Extensions: count = 26
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
VK_EXT_headless_surface                : extension revision 1
VK_EXT_layer_settings                  : extension revision 2
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
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_INTEL_nullhw       INTEL NULL HW                                                1.1.73   version 1
VK_LAYER_MESA_anti_lag      Open-source implementation of the VK_AMD_anti_lag extension. 1.4.303  version 1
VK_LAYER_MESA_device_select Linux device selection layer                                 1.4.303  version 1
VK_LAYER_MESA_overlay       Mesa Overlay layer                                           1.4.303  version 1
VK_LAYER_MESA_screenshot    Mesa Screenshot layer                                        1.4.303  version 1

Devices:
========
GPU0:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x1002
	deviceID           = 0x744c
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = AMD Radeon RX 7900 XTX (RADV NAVI31)
	driverID           = DRIVER_ID_MESA_RADV
	driverName         = radv
	driverInfo         = Mesa 26.0.8-1ubuntu0.3
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-8300-0000-0000-000000000000
	driverUUID         = 414d442d-4d45-5341-2d44-525600000000
GPU1:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x1002
	deviceID           = 0x744c
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = AMD Radeon RX 7900 XTX (RADV NAVI31)
	driverID           = DRIVER_ID_MESA_RADV
	driverName         = radv
	driverInfo         = Mesa 26.0.8-1ubuntu0.3
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-8600-0000-0000-000000000000
	driverUUID         = 414d442d-4d45-5341-2d44-525600000000
GPU2:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x10005
	deviceID           = 0x0000
	deviceType         = PHYSICAL_DEVICE_TYPE_CPU
	deviceName         = llvmpipe (LLVM 21.1.8, 256 bits)
	driverID           = DRIVER_ID_MESA_LLVMPIPE
	driverName         = llvmpipe
	driverInfo         = Mesa 26.0.8-1ubuntu0.3 (LLVM 21.1.8)
	conformanceVersion = 1.3.1.1
	deviceUUID         = 6d657361-3236-2e30-2e38-2d3175627500
	driverUUID         = 6c6c766d-7069-7065-5555-494400000000
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=Vulkan0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkInference-48    	     319	  37200317 ns/op	       806.4 tokens/s
BenchmarkInference-48    	     322	  37022935 ns/op	       810.3 tokens/s
BenchmarkInference-48    	     320	  37031336 ns/op	       810.1 tokens/s
BenchmarkInference-48    	     322	  37196656 ns/op	       806.5 tokens/s
BenchmarkInference-48    	     316	  37438856 ns/op	       801.3 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.763s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/cookie3/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/cookie3/vulkan1 -->
### Vulkan, amd64, AMD EPYC 7443P 24-Core Processor, Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"cookie3","device":"Vulkan1","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":791.5,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD EPYC 7443P 24-Core Processor. 791.5 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.341


Instance Extensions: count = 26
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
VK_EXT_headless_surface                : extension revision 1
VK_EXT_layer_settings                  : extension revision 2
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
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_INTEL_nullhw       INTEL NULL HW                                                1.1.73   version 1
VK_LAYER_MESA_anti_lag      Open-source implementation of the VK_AMD_anti_lag extension. 1.4.303  version 1
VK_LAYER_MESA_device_select Linux device selection layer                                 1.4.303  version 1
VK_LAYER_MESA_overlay       Mesa Overlay layer                                           1.4.303  version 1
VK_LAYER_MESA_screenshot    Mesa Screenshot layer                                        1.4.303  version 1

Devices:
========
GPU0:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x1002
	deviceID           = 0x744c
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = AMD Radeon RX 7900 XTX (RADV NAVI31)
	driverID           = DRIVER_ID_MESA_RADV
	driverName         = radv
	driverInfo         = Mesa 26.0.8-1ubuntu0.3
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-8300-0000-0000-000000000000
	driverUUID         = 414d442d-4d45-5341-2d44-525600000000
GPU1:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x1002
	deviceID           = 0x744c
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = AMD Radeon RX 7900 XTX (RADV NAVI31)
	driverID           = DRIVER_ID_MESA_RADV
	driverName         = radv
	driverInfo         = Mesa 26.0.8-1ubuntu0.3
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-8600-0000-0000-000000000000
	driverUUID         = 414d442d-4d45-5341-2d44-525600000000
GPU2:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x10005
	deviceID           = 0x0000
	deviceType         = PHYSICAL_DEVICE_TYPE_CPU
	deviceName         = llvmpipe (LLVM 21.1.8, 256 bits)
	driverID           = DRIVER_ID_MESA_LLVMPIPE
	driverName         = llvmpipe
	driverInfo         = Mesa 26.0.8-1ubuntu0.3 (LLVM 21.1.8)
	conformanceVersion = 1.3.1.1
	deviceUUID         = 6d657361-3236-2e30-2e38-2d3175627500
	driverUUID         = 6c6c766d-7069-7065-5555-494400000000
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=Vulkan1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkInference-48    	     314	  38202881 ns/op	       785.3 tokens/s
BenchmarkInference-48    	     314	  37782531 ns/op	       794.0 tokens/s
BenchmarkInference-48    	     314	  37904955 ns/op	       791.5 tokens/s
BenchmarkInference-48    	     318	  37672420 ns/op	       796.3 tokens/s
BenchmarkInference-48    	     313	  37924995 ns/op	       791.0 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.036s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/cookie3/vulkan1 -->

<!-- yzma:bench start text/vulkan/amd64/i9-13900hx/vulkan0 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":94.93,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 94.9 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.275


Instance Extensions: count = 24
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
VK_EXT_headless_surface                : extension revision 1
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
VK_KHR_surface_protected_capabilities  : extension revision 1
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_INTEL_nullhw       INTEL NULL HW                1.1.73   version 1
VK_LAYER_MESA_device_select Linux device selection layer 1.4.303  version 1
VK_LAYER_MESA_overlay       Mesa Overlay layer           1.4.303  version 1
VK_LAYER_NV_optimus         NVIDIA Optimus layer         1.4.329  version 1
VK_LAYER_NV_present         NVIDIA Presentation Layer    1.4.329  version 1

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
	driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.2
	conformanceVersion = 1.4.0.0
	deviceUUID         = 868088a7-0400-0000-0002-000000000000
	driverUUID         = ee99561e-45e1-e718-c612-1d36d8345582
GPU1:
	apiVersion         = 1.4.329
	driverVersion      = 595.99.2.0
	vendorID           = 0x10de
	deviceID           = 0x2860
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.99.02
	conformanceVersion = 1.4.3.3
	deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
	driverUUID         = c9bf27cc-ad5e-57db-8930-96d25aad2e91
GPU2:
	apiVersion         = 1.4.318
	driverVersion      = 25.2.8
	vendorID           = 0x10005
	deviceID           = 0x0000
	deviceType         = PHYSICAL_DEVICE_TYPE_CPU
	deviceName         = llvmpipe (LLVM 20.1.2, 256 bits)
	driverID           = DRIVER_ID_MESA_LLVMPIPE
	driverName         = llvmpipe
	driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.2 (LLVM 20.1.2)
	conformanceVersion = 1.3.1.1
	deviceUUID         = 6d657361-3235-2e32-2e38-2d3075627500
	driverUUID         = 6c6c766d-7069-7065-5555-494400000000
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=Vulkan0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	      33	 331250522 ns/op	        90.57 tokens/s
BenchmarkInference-32    	      37	 315408816 ns/op	        95.11 tokens/s
BenchmarkInference-32    	      36	 313697435 ns/op	        95.63 tokens/s
BenchmarkInference-32    	      37	 316404824 ns/op	        94.82 tokens/s
BenchmarkInference-32    	      36	 316029018 ns/op	        94.93 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	67.359s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/i9-13900hx/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/i9-13900hx/vulkan1 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan1","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":747.5,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 747.5 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.275


Instance Extensions: count = 24
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
VK_EXT_headless_surface                : extension revision 1
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
VK_KHR_surface_protected_capabilities  : extension revision 1
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_INTEL_nullhw       INTEL NULL HW                1.1.73   version 1
VK_LAYER_MESA_device_select Linux device selection layer 1.4.303  version 1
VK_LAYER_MESA_overlay       Mesa Overlay layer           1.4.303  version 1
VK_LAYER_NV_optimus         NVIDIA Optimus layer         1.4.329  version 1
VK_LAYER_NV_present         NVIDIA Presentation Layer    1.4.329  version 1

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
	driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.2
	conformanceVersion = 1.4.0.0
	deviceUUID         = 868088a7-0400-0000-0002-000000000000
	driverUUID         = ee99561e-45e1-e718-c612-1d36d8345582
GPU1:
	apiVersion         = 1.4.329
	driverVersion      = 595.99.2.0
	vendorID           = 0x10de
	deviceID           = 0x2860
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.99.02
	conformanceVersion = 1.4.3.3
	deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
	driverUUID         = c9bf27cc-ad5e-57db-8930-96d25aad2e91
GPU2:
	apiVersion         = 1.4.318
	driverVersion      = 25.2.8
	vendorID           = 0x10005
	deviceID           = 0x0000
	deviceType         = PHYSICAL_DEVICE_TYPE_CPU
	deviceName         = llvmpipe (LLVM 20.1.2, 256 bits)
	driverID           = DRIVER_ID_MESA_LLVMPIPE
	driverName         = llvmpipe
	driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.2 (LLVM 20.1.2)
	conformanceVersion = 1.3.1.1
	deviceUUID         = 6d657361-3235-2e32-2e38-2d3075627500
	driverUUID         = 6c6c766d-7069-7065-5555-494400000000
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=Vulkan1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	     212	  51789844 ns/op	       579.3 tokens/s
BenchmarkInference-32    	     296	  40018127 ns/op	       749.7 tokens/s
BenchmarkInference-32    	     300	  40134672 ns/op	       747.5 tokens/s
BenchmarkInference-32    	     296	  40571890 ns/op	       739.4 tokens/s
BenchmarkInference-32    	     302	  40049532 ns/op	       749.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	77.406s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/i9-13900hx/vulkan1 -->

<!-- yzma:bench start text/vulkan/arm64/localhost/vulkan0 -->
### Vulkan, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"arm64","machine":"localhost","device":"Vulkan0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":171.2,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-17"} -->

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.321


Instance Extensions: count = 25
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
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
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1
VK_NV_display_stereo                   : extension revision 1

Instance Layers:
----------------

Devices:
========
GPU0:
	apiVersion         = 1.4.329
	driverVersion      = 595.78.0.0
	vendorID           = 0x10de
	deviceID           = 0x97ba03d7
	deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
	deviceName         = NVIDIA Tegra Orin (nvgpu)
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.78
	conformanceVersion = 1.4.3.3
	deviceUUID         = 1388f9e0-987e-54a0-908f-6a30d8fd5f29
	driverUUID         = 5feccbaa-e166-5f7f-8cd4-26d96ad572dc
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000   -device=Vulkan0
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-6   	      33	 310347413 ns/op	        96.67 tokens/s
BenchmarkInference-6   	      62	 171340736 ns/op	       175.1 tokens/s
BenchmarkInference-6   	      62	 172843447 ns/op	       173.6 tokens/s
BenchmarkInference-6   	      63	 175256766 ns/op	       171.2 tokens/s
BenchmarkInference-6   	      66	 176780428 ns/op	       169.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	68.210s
```

</details>
<!-- yzma:bench end text/vulkan/arm64/localhost/vulkan0 -->

## Multimodal model benchmarks

The model is
[SmolVLM-256M-Instruct-Q8_0.gguf](https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/SmolVLM-256M-Instruct-Q8_0.gguf)
with its
[projector](https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf).
The code is [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go).

<!-- yzma:bench table multimodal -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | amd64 | AMD EPYC 7443P 24-Core Processor | - | 531.1 | b11146 | 2026-09-25 |
| CPU | amd64 | Intel Core i9-13900HX | - | 842.6 | b11669 | 2026-10-16 |
| CPU | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | - | 213.5 | b11669 | 2026-10-17 |
| CPU | arm64 | Raspberry Pi 4 Model B Rev 1.4 | - | 5.6 | b11146 | 2026-09-24 |
| CPU | arm64 | Arduino UnoQ | - | 4.1 | b11146 | 2026-09-24 |
| CUDA | amd64 | Intel Core i9-13900HX | CUDA0 | 2448.0 | b11669 | 2026-10-16 |
| CUDA | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | CUDA0 | 451.1 | b11669 | 2026-10-17 |
| ROCm | amd64 | AMD EPYC 7443P 24-Core Processor | ROCm0 | 1461.0 | b11146 | 2026-09-25 |
| ROCm | amd64 | AMD EPYC 7443P 24-Core Processor | ROCm1 | 1203.0 | b11146 | 2026-09-25 |
| Vulkan | amd64 | AMD EPYC 7443P 24-Core Processor | Vulkan0 | 1350.0 | b11146 | 2026-09-24 |
| Vulkan | amd64 | AMD EPYC 7443P 24-Core Processor | Vulkan1 | 1108.0 | b11146 | 2026-09-24 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan0 | 395.7 | b11669 | 2026-10-16 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan1 | 2178.0 | b11669 | 2026-10-16 |
| Vulkan | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | Vulkan0 | 395.6 | b11669 | 2026-10-17 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/cookie3 -->
### CPU, amd64, AMD EPYC 7443P 24-Core Processor
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"cookie3","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":531.1,"llamacpp":"b11146","yzma":"1.29.0-dev","date":"2026-09-25"} -->

AMD EPYC 7443P 24-Core Processor. 531.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkMultimodalInference-48    	      28	 416150767 ns/op	       557.1 tokens/s
BenchmarkMultimodalInference-48    	      26	 435851928 ns/op	       531.1 tokens/s
BenchmarkMultimodalInference-48    	      31	 434341982 ns/op	       533.3 tokens/s
BenchmarkMultimodalInference-48    	      25	 466286117 ns/op	       504.9 tokens/s
BenchmarkMultimodalInference-48    	      28	 481531058 ns/op	       493.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	63.979s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/cookie3 -->

<!-- yzma:bench start multimodal/cpu/amd64/i9-13900hx -->
### CPU, amd64, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":842.6,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 842.6 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	      38	 303118702 ns/op	       791.1 tokens/s
BenchmarkMultimodalInference-32    	      44	 263852449 ns/op	       882.5 tokens/s
BenchmarkMultimodalInference-32    	      44	 264812321 ns/op	       881.8 tokens/s
BenchmarkMultimodalInference-32    	      44	 280906572 ns/op	       842.6 tokens/s
BenchmarkMultimodalInference-32    	      58	 281323129 ns/op	       838.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	66.512s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/i9-13900hx -->

<!-- yzma:bench start multimodal/cpu/arm64/localhost -->
### CPU, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"localhost","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":213.5,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-17"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-6   	      10	1150553015 ns/op	       207.7 tokens/s
BenchmarkMultimodalInference-6   	      12	1011646476 ns/op	       227.7 tokens/s
BenchmarkMultimodalInference-6   	       9	1149530144 ns/op	       207.7 tokens/s
BenchmarkMultimodalInference-6   	      10	1064425774 ns/op	       220.0 tokens/s
BenchmarkMultimodalInference-6   	      10	1105388193 ns/op	       213.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	57.076s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/localhost -->

<!-- yzma:bench start multimodal/cpu/arm64/raspberrypi -->
### CPU, arm64, Raspberry Pi 4 Model B Rev 1.4
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"raspberrypi","label":"Raspberry Pi 4 Model B Rev 1.4","tokens_per_second":5.596,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-4   	       1	34320762778 ns/op	         6.876 tokens/s
BenchmarkMultimodalInference-4   	       1	37906633685 ns/op	         6.358 tokens/s
BenchmarkMultimodalInference-4   	       1	40207740487 ns/op	         5.596 tokens/s
BenchmarkMultimodalInference-4   	       1	43692589003 ns/op	         5.333 tokens/s
BenchmarkMultimodalInference-4   	       1	44633338087 ns/op	         5.041 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	208.605s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/raspberrypi -->

<!-- yzma:bench start multimodal/cpu/arm64/yzma -->
### CPU, arm64, Arduino UnoQ
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"yzma","label":"Arduino UnoQ","tokens_per_second":4.056,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-4   	       1	55802770374 ns/op	         4.032 tokens/s
BenchmarkMultimodalInference-4   	       1	56029573780 ns/op	         4.051 tokens/s
BenchmarkMultimodalInference-4   	       1	55965817409 ns/op	         4.056 tokens/s
BenchmarkMultimodalInference-4   	       1	55928726514 ns/op	         4.059 tokens/s
BenchmarkMultimodalInference-4   	       1	56101565335 ns/op	         4.100 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	282.332s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/yzma -->

<!-- yzma:bench start multimodal/cuda/amd64/i9-13900hx/cuda0 -->
### CUDA, amd64, Intel Core i9-13900HX, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"i9-13900hx","device":"CUDA0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":2448,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 2448.0 tokens a second.

<details><summary>The device</summary>

```
Fri Oct 16 13:56:21 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.99.02              Driver Version: 595.99.02      CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   48C    P8              7W /  115W |      15MiB /   8188MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            6259      G   /usr/lib/xorg/Xorg                        4MiB |
+-----------------------------------------------------------------------------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=CUDA0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	     128	  92516598 ns/op	      2532 tokens/s
BenchmarkMultimodalInference-32    	     123	  95804675 ns/op	      2446 tokens/s
BenchmarkMultimodalInference-32    	     129	  92276039 ns/op	      2534 tokens/s
BenchmarkMultimodalInference-32    	     126	  95593974 ns/op	      2448 tokens/s
BenchmarkMultimodalInference-32    	     121	  98204289 ns/op	      2403 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.224s
```

</details>
<!-- yzma:bench end multimodal/cuda/amd64/i9-13900hx/cuda0 -->

<!-- yzma:bench start multimodal/cuda/arm64/localhost/cuda0 -->
### CUDA, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"arm64","machine":"localhost","device":"CUDA0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":451.1,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-17"} -->

<details><summary>The device</summary>

```
Sat Oct 17 09:09:45 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.78                 Driver Version: 595.78         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  Orin (nvgpu)                  N/A  |   N/A              N/A |                  N/A |
| N/A   N/A  N/A             N/A  /  N/A  | Not Supported          |     N/A          N/A |
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
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=CUDA0
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-6   	      20	 553187494 ns/op	       432.7 tokens/s
BenchmarkMultimodalInference-6   	      21	 548401368 ns/op	       441.4 tokens/s
BenchmarkMultimodalInference-6   	      22	 526099399 ns/op	       451.1 tokens/s
BenchmarkMultimodalInference-6   	      22	 522303247 ns/op	       454.9 tokens/s
BenchmarkMultimodalInference-6   	      25	 512168925 ns/op	       458.4 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	60.574s
```

</details>
<!-- yzma:bench end multimodal/cuda/arm64/localhost/cuda0 -->

<!-- yzma:bench start multimodal/rocm/amd64/cookie3/rocm0 -->
### ROCm, amd64, AMD EPYC 7443P 24-Core Processor, ROCm0
<!-- yzma:bench meta {"suite":"multimodal","backend":"rocm","arch":"amd64","machine":"cookie3","device":"ROCm0","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":1461,"llamacpp":"b11146","yzma":"1.29.0-dev","date":"2026-09-25"} -->

AMD EPYC 7443P 24-Core Processor. 1461.0 tokens a second.

<details><summary>The device</summary>

```


======================================= ROCm System Management Interface =======================================
================================================= Concise Info =================================================
Device  Node  IDs              Temp    Power  Partitions          SCLK  MCLK    Fan  Perf  PwrCap  VRAM%  GPU%  
[3m              (DID,     GUID)  (Edge)  (Avg)  (Mem, Compute, ID)                                                [0m
================================================================================================================
0       1     0x744c,   16626  31.0°C  59.0W  N/A, N/A, 0         0Mhz  456Mhz  0%   auto  0.0W    0%     0%    
1       2     0x744c,   19095  27.0°C  61.0W  N/A, N/A, 0         0Mhz  456Mhz  0%   auto  0.0W    0%     0%    
================================================================================================================
============================================= End of ROCm SMI Log ==============================================
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=ROCm0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkMultimodalInference-48    	      66	 169443589 ns/op	      1390 tokens/s
BenchmarkMultimodalInference-48    	      74	 159135187 ns/op	      1479 tokens/s
BenchmarkMultimodalInference-48    	      90	 164132082 ns/op	      1443 tokens/s
BenchmarkMultimodalInference-48    	      74	 157266751 ns/op	      1490 tokens/s
BenchmarkMultimodalInference-48    	      74	 161120145 ns/op	      1461 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.499s
```

</details>
<!-- yzma:bench end multimodal/rocm/amd64/cookie3/rocm0 -->

<!-- yzma:bench start multimodal/rocm/amd64/cookie3/rocm1 -->
### ROCm, amd64, AMD EPYC 7443P 24-Core Processor, ROCm1
<!-- yzma:bench meta {"suite":"multimodal","backend":"rocm","arch":"amd64","machine":"cookie3","device":"ROCm1","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":1203,"llamacpp":"b11146","yzma":"1.29.0-dev","date":"2026-09-25"} -->

AMD EPYC 7443P 24-Core Processor. 1203.0 tokens a second.

<details><summary>The device</summary>

```


========================================= ROCm System Management Interface =========================================
=================================================== Concise Info ===================================================
Device  Node  IDs              Temp    Power   Partitions          SCLK  MCLK   Fan     Perf  PwrCap  VRAM%  GPU%  
[3m              (DID,     GUID)  (Edge)  (Avg)   (Mem, Compute, ID)                                                  [0m
====================================================================================================================
0       1     0x744c,   16626  46.0°C  139.0W  N/A, N/A, 0         0Mhz  96Mhz  28.63%  auto  0.0W    0%     0%    
1       2     0x744c,   19095  26.0°C  12.0W   N/A, N/A, 0         0Mhz  96Mhz  0%      auto  0.0W    0%     0%    
====================================================================================================================
=============================================== End of ROCm SMI Log ================================================
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=ROCm1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkMultimodalInference-48    	      55	 202216190 ns/op	      1161 tokens/s
BenchmarkMultimodalInference-48    	      54	 200107492 ns/op	      1178 tokens/s
BenchmarkMultimodalInference-48    	      54	 192329472 ns/op	      1217 tokens/s
BenchmarkMultimodalInference-48    	      67	 192463774 ns/op	      1219 tokens/s
BenchmarkMultimodalInference-48    	      57	 195761670 ns/op	      1203 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	57.350s
```

</details>
<!-- yzma:bench end multimodal/rocm/amd64/cookie3/rocm1 -->

<!-- yzma:bench start multimodal/vulkan/amd64/cookie3/vulkan0 -->
### Vulkan, amd64, AMD EPYC 7443P 24-Core Processor, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"cookie3","device":"Vulkan0","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":1350,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD EPYC 7443P 24-Core Processor. 1350.0 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.341


Instance Extensions: count = 26
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
VK_EXT_headless_surface                : extension revision 1
VK_EXT_layer_settings                  : extension revision 2
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
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_INTEL_nullhw       INTEL NULL HW                                                1.1.73   version 1
VK_LAYER_MESA_anti_lag      Open-source implementation of the VK_AMD_anti_lag extension. 1.4.303  version 1
VK_LAYER_MESA_device_select Linux device selection layer                                 1.4.303  version 1
VK_LAYER_MESA_overlay       Mesa Overlay layer                                           1.4.303  version 1
VK_LAYER_MESA_screenshot    Mesa Screenshot layer                                        1.4.303  version 1

Devices:
========
GPU0:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x1002
	deviceID           = 0x744c
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = AMD Radeon RX 7900 XTX (RADV NAVI31)
	driverID           = DRIVER_ID_MESA_RADV
	driverName         = radv
	driverInfo         = Mesa 26.0.8-1ubuntu0.3
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-8300-0000-0000-000000000000
	driverUUID         = 414d442d-4d45-5341-2d44-525600000000
GPU1:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x1002
	deviceID           = 0x744c
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = AMD Radeon RX 7900 XTX (RADV NAVI31)
	driverID           = DRIVER_ID_MESA_RADV
	driverName         = radv
	driverInfo         = Mesa 26.0.8-1ubuntu0.3
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-8600-0000-0000-000000000000
	driverUUID         = 414d442d-4d45-5341-2d44-525600000000
GPU2:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x10005
	deviceID           = 0x0000
	deviceType         = PHYSICAL_DEVICE_TYPE_CPU
	deviceName         = llvmpipe (LLVM 21.1.8, 256 bits)
	driverID           = DRIVER_ID_MESA_LLVMPIPE
	driverName         = llvmpipe
	driverInfo         = Mesa 26.0.8-1ubuntu0.3 (LLVM 21.1.8)
	conformanceVersion = 1.3.1.1
	deviceUUID         = 6d657361-3236-2e30-2e38-2d3175627500
	driverUUID         = 6c6c766d-7069-7065-5555-494400000000
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=Vulkan0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkMultimodalInference-48    	      62	 171931417 ns/op	      1350 tokens/s
BenchmarkMultimodalInference-48    	      57	 176524904 ns/op	      1323 tokens/s
BenchmarkMultimodalInference-48    	      70	 175311360 ns/op	      1333 tokens/s
BenchmarkMultimodalInference-48    	      61	 171654851 ns/op	      1354 tokens/s
BenchmarkMultimodalInference-48    	      60	 169699917 ns/op	      1367 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	54.959s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/cookie3/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/cookie3/vulkan1 -->
### Vulkan, amd64, AMD EPYC 7443P 24-Core Processor, Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"cookie3","device":"Vulkan1","label":"AMD EPYC 7443P 24-Core Processor","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":1108,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

AMD EPYC 7443P 24-Core Processor. 1108.0 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.341


Instance Extensions: count = 26
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
VK_EXT_headless_surface                : extension revision 1
VK_EXT_layer_settings                  : extension revision 2
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
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_INTEL_nullhw       INTEL NULL HW                                                1.1.73   version 1
VK_LAYER_MESA_anti_lag      Open-source implementation of the VK_AMD_anti_lag extension. 1.4.303  version 1
VK_LAYER_MESA_device_select Linux device selection layer                                 1.4.303  version 1
VK_LAYER_MESA_overlay       Mesa Overlay layer                                           1.4.303  version 1
VK_LAYER_MESA_screenshot    Mesa Screenshot layer                                        1.4.303  version 1

Devices:
========
GPU0:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x1002
	deviceID           = 0x744c
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = AMD Radeon RX 7900 XTX (RADV NAVI31)
	driverID           = DRIVER_ID_MESA_RADV
	driverName         = radv
	driverInfo         = Mesa 26.0.8-1ubuntu0.3
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-8300-0000-0000-000000000000
	driverUUID         = 414d442d-4d45-5341-2d44-525600000000
GPU1:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x1002
	deviceID           = 0x744c
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = AMD Radeon RX 7900 XTX (RADV NAVI31)
	driverID           = DRIVER_ID_MESA_RADV
	driverName         = radv
	driverInfo         = Mesa 26.0.8-1ubuntu0.3
	conformanceVersion = 1.4.0.0
	deviceUUID         = 00000000-8600-0000-0000-000000000000
	driverUUID         = 414d442d-4d45-5341-2d44-525600000000
GPU2:
	apiVersion         = 1.4.335
	driverVersion      = 26.0.8
	vendorID           = 0x10005
	deviceID           = 0x0000
	deviceType         = PHYSICAL_DEVICE_TYPE_CPU
	deviceName         = llvmpipe (LLVM 21.1.8, 256 bits)
	driverID           = DRIVER_ID_MESA_LLVMPIPE
	driverName         = llvmpipe
	driverInfo         = Mesa 26.0.8-1ubuntu0.3 (LLVM 21.1.8)
	conformanceVersion = 1.3.1.1
	deviceUUID         = 6d657361-3236-2e30-2e38-2d3175627500
	driverUUID         = 6c6c766d-7069-7065-5555-494400000000
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=Vulkan1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD EPYC 7443P 24-Core Processor               
BenchmarkMultimodalInference-48    	      52	 208638597 ns/op	      1121 tokens/s
BenchmarkMultimodalInference-48    	      54	 215018444 ns/op	      1089 tokens/s
BenchmarkMultimodalInference-48    	      49	 210643485 ns/op	      1108 tokens/s
BenchmarkMultimodalInference-48    	      55	 214275826 ns/op	      1087 tokens/s
BenchmarkMultimodalInference-48    	      56	 208981753 ns/op	      1108 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	56.983s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/cookie3/vulkan1 -->

<!-- yzma:bench start multimodal/vulkan/amd64/i9-13900hx/vulkan0 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":395.7,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 395.7 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.275


Instance Extensions: count = 24
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
VK_EXT_headless_surface                : extension revision 1
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
VK_KHR_surface_protected_capabilities  : extension revision 1
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_INTEL_nullhw       INTEL NULL HW                1.1.73   version 1
VK_LAYER_MESA_device_select Linux device selection layer 1.4.303  version 1
VK_LAYER_MESA_overlay       Mesa Overlay layer           1.4.303  version 1
VK_LAYER_NV_optimus         NVIDIA Optimus layer         1.4.329  version 1
VK_LAYER_NV_present         NVIDIA Presentation Layer    1.4.329  version 1

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
	driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.2
	conformanceVersion = 1.4.0.0
	deviceUUID         = 868088a7-0400-0000-0002-000000000000
	driverUUID         = ee99561e-45e1-e718-c612-1d36d8345582
GPU1:
	apiVersion         = 1.4.329
	driverVersion      = 595.99.2.0
	vendorID           = 0x10de
	deviceID           = 0x2860
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.99.02
	conformanceVersion = 1.4.3.3
	deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
	driverUUID         = c9bf27cc-ad5e-57db-8930-96d25aad2e91
GPU2:
	apiVersion         = 1.4.318
	driverVersion      = 25.2.8
	vendorID           = 0x10005
	deviceID           = 0x0000
	deviceType         = PHYSICAL_DEVICE_TYPE_CPU
	deviceName         = llvmpipe (LLVM 20.1.2, 256 bits)
	driverID           = DRIVER_ID_MESA_LLVMPIPE
	driverName         = llvmpipe
	driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.2 (LLVM 20.1.2)
	conformanceVersion = 1.3.1.1
	deviceUUID         = 6d657361-3235-2e32-2e38-2d3075627500
	driverUUID         = 6c6c766d-7069-7065-5555-494400000000
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=Vulkan0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	       7	1438115517 ns/op	       158.9 tokens/s
BenchmarkMultimodalInference-32    	      20	 596850604 ns/op	       389.7 tokens/s
BenchmarkMultimodalInference-32    	      22	 579556373 ns/op	       400.8 tokens/s
BenchmarkMultimodalInference-32    	      20	 585445522 ns/op	       395.7 tokens/s
BenchmarkMultimodalInference-32    	      20	 580367606 ns/op	       399.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.778s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/i9-13900hx/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/i9-13900hx/vulkan1 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan1","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":2178,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 2178.0 tokens a second.

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.275


Instance Extensions: count = 24
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
VK_EXT_headless_surface                : extension revision 1
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
VK_KHR_surface_protected_capabilities  : extension revision 1
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1

Instance Layers: count = 5
--------------------------
VK_LAYER_INTEL_nullhw       INTEL NULL HW                1.1.73   version 1
VK_LAYER_MESA_device_select Linux device selection layer 1.4.303  version 1
VK_LAYER_MESA_overlay       Mesa Overlay layer           1.4.303  version 1
VK_LAYER_NV_optimus         NVIDIA Optimus layer         1.4.329  version 1
VK_LAYER_NV_present         NVIDIA Presentation Layer    1.4.329  version 1

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
	driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.2
	conformanceVersion = 1.4.0.0
	deviceUUID         = 868088a7-0400-0000-0002-000000000000
	driverUUID         = ee99561e-45e1-e718-c612-1d36d8345582
GPU1:
	apiVersion         = 1.4.329
	driverVersion      = 595.99.2.0
	vendorID           = 0x10de
	deviceID           = 0x2860
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.99.02
	conformanceVersion = 1.4.3.3
	deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
	driverUUID         = c9bf27cc-ad5e-57db-8930-96d25aad2e91
GPU2:
	apiVersion         = 1.4.318
	driverVersion      = 25.2.8
	vendorID           = 0x10005
	deviceID           = 0x0000
	deviceType         = PHYSICAL_DEVICE_TYPE_CPU
	deviceName         = llvmpipe (LLVM 20.1.2, 256 bits)
	driverID           = DRIVER_ID_MESA_LLVMPIPE
	driverName         = llvmpipe
	driverInfo         = Mesa 25.2.8-0ubuntu0.24.04.2 (LLVM 20.1.2)
	conformanceVersion = 1.3.1.1
	deviceUUID         = 6d657361-3235-2e32-2e38-2d3075627500
	driverUUID         = 6c6c766d-7069-7065-5555-494400000000
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=Vulkan1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	      75	 136997259 ns/op	      1703 tokens/s
BenchmarkMultimodalInference-32    	     100	 107425562 ns/op	      2178 tokens/s
BenchmarkMultimodalInference-32    	     100	 106421518 ns/op	      2190 tokens/s
BenchmarkMultimodalInference-32    	     100	 106116474 ns/op	      2195 tokens/s
BenchmarkMultimodalInference-32    	     110	 107742168 ns/op	      2166 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	60.887s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/i9-13900hx/vulkan1 -->

<!-- yzma:bench start multimodal/vulkan/arm64/localhost/vulkan0 -->
### Vulkan, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"arm64","machine":"localhost","device":"Vulkan0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":395.6,"llamacpp":"b11669","yzma":"1.30.0","date":"2026-10-17"} -->

<details><summary>The device</summary>

```
==========
VULKANINFO
==========

Vulkan Instance Version: 1.4.321


Instance Extensions: count = 25
-------------------------------
VK_EXT_acquire_drm_display             : extension revision 1
VK_EXT_acquire_xlib_display            : extension revision 1
VK_EXT_debug_report                    : extension revision 10
VK_EXT_debug_utils                     : extension revision 2
VK_EXT_direct_mode_display             : extension revision 1
VK_EXT_display_surface_counter         : extension revision 1
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
VK_KHR_wayland_surface                 : extension revision 6
VK_KHR_xcb_surface                     : extension revision 6
VK_KHR_xlib_surface                    : extension revision 6
VK_LUNARG_direct_driver_loading        : extension revision 1
VK_NV_display_stereo                   : extension revision 1

Instance Layers:
----------------

Devices:
========
GPU0:
	apiVersion         = 1.4.329
	driverVersion      = 595.78.0.0
	vendorID           = 0x10de
	deviceID           = 0x97ba03d7
	deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
	deviceName         = NVIDIA Tegra Orin (nvgpu)
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.78
	conformanceVersion = 1.4.3.3
	deviceUUID         = 1388f9e0-987e-54a0-908f-6a30d8fd5f29
	driverUUID         = 5feccbaa-e166-5f7f-8cd4-26d96ad572dc
```

</details>

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000   -device=Vulkan0
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-6   	      10	1014784873 ns/op	       231.3 tokens/s
BenchmarkMultimodalInference-6   	      20	 598311893 ns/op	       395.6 tokens/s
BenchmarkMultimodalInference-6   	      21	 548857285 ns/op	       422.0 tokens/s
BenchmarkMultimodalInference-6   	      20	 570344220 ns/op	       410.5 tokens/s
BenchmarkMultimodalInference-6   	      19	 619280811 ns/op	       387.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.660s
```

</details>
<!-- yzma:bench end multimodal/vulkan/arm64/localhost/vulkan0 -->
