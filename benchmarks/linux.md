# Linux benchmarks

Benchmarks of yzma on Linux. Each table gives the median of five runs. The output
of each run, and of the device, is below the tables.

To add a machine or to make these numbers again, see
[how to run the benchmarks](README.md).

## Summary

Tokens a second on each machine. The GPU columns give the fastest GPU backend.

| Machine | GPU | Text, CPU | Text, GPU | Multimodal, CPU | Multimodal, GPU |
| --- | --- | --- | --- | --- | --- |
| Intel Core i9-13900HX | RTX 4070 Laptop | 265.1 | 853.4, CUDA | 878.7 | 2277.0, CUDA |
| NVIDIA Jetson Orin Nano Super | Tegra Orin | 87.0 | 191.6, CUDA | 212.0 | 419.1, CUDA |
| Raspberry Pi 4 Model B | none | 35.3 | none | 5.5 | none |
| Arduino UNO Q | none | 32.0 | none | 4.2 | none |

- CUDA on the RTX 4070 gives the fastest results of all Linux machines.
- Vulkan on the same RTX 4070 gives 87 percent of CUDA for text and 94 percent
  for multimodal.
- The integrated Intel GPU of the i9-13900HX is slower than its CPU. On Vulkan0
  it gives 95.5 for text and 476.8 for multimodal.
- On the Jetson Orin Nano, the GPU is 2.2 times faster than the CPU for text and
  2.0 times faster for multimodal. Vulkan gives 93 percent of CUDA for text and
  99 percent for multimodal.
- The Raspberry Pi 4 and the Arduino UNO Q have no GPU backend. They run text at
  more than 30 tokens a second, but multimodal is slow.
- The text suite uses 4 threads on each machine. The multimodal suite uses one
  thread for each performance core, thus its CPU column shows the size of the
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
| CPU | amd64 | Intel Core i9-13900HX | - | 265.1 | b11146 | 2026-09-24 |
| CPU | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | - | 87.0 | b10964 | 2026-09-23 |
| CPU | arm64 | Raspberry Pi 4 Model B Rev 1.4 | - | 35.3 | b10964 | 2026-09-23 |
| CPU | arm64 | Arduino UnoQ | - | 32.0 | b10964 | 2026-09-23 |
| CUDA | amd64 | Intel Core i9-13900HX | CUDA0 | 853.4 | b11146 | 2026-09-24 |
| CUDA | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | CUDA0 | 191.6 | b10964 | 2026-09-23 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan0 | 95.5 | b11146 | 2026-09-24 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan1 | 744.2 | b11146 | 2026-09-24 |
| Vulkan | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | Vulkan0 | 177.5 | b10964 | 2026-09-23 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/i9-13900hx -->
### CPU, amd64, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":265.1,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 265.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192  -threadpool -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	     100	 114444991 ns/op	       262.1 tokens/s
BenchmarkInference-32    	     103	 114085020 ns/op	       263.0 tokens/s
BenchmarkInference-32    	     100	 111509117 ns/op	       269.0 tokens/s
BenchmarkInference-32    	     100	 111645951 ns/op	       268.7 tokens/s
BenchmarkInference-32    	     100	 113144808 ns/op	       265.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.197s
```

</details>
<!-- yzma:bench end text/cpu/amd64/i9-13900hx -->

<!-- yzma:bench start text/cpu/arm64/localhost -->
### CPU, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"localhost","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":87.02,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-6   	      36	 351380582 ns/op	        85.38 tokens/s
BenchmarkInference-6   	      33	 343945806 ns/op	        87.22 tokens/s
BenchmarkInference-6   	      32	 342749008 ns/op	        87.53 tokens/s
BenchmarkInference-6   	      34	 344729854 ns/op	        87.02 tokens/s
BenchmarkInference-6   	      34	 353905994 ns/op	        84.77 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	60.902s
```

</details>
<!-- yzma:bench end text/cpu/arm64/localhost -->

<!-- yzma:bench start text/cpu/arm64/raspberrypi -->
### CPU, arm64, Raspberry Pi 4 Model B Rev 1.4
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"raspberrypi","label":"Raspberry Pi 4 Model B Rev 1.4","tokens_per_second":35.26,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-4   	      13	 850633907 ns/op	        35.27 tokens/s
BenchmarkInference-4   	      13	 860719468 ns/op	        34.85 tokens/s
BenchmarkInference-4   	      13	 855645047 ns/op	        35.06 tokens/s
BenchmarkInference-4   	      13	 849870406 ns/op	        35.30 tokens/s
BenchmarkInference-4   	      13	 850919028 ns/op	        35.26 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.566s
```

</details>
<!-- yzma:bench end text/cpu/arm64/raspberrypi -->

<!-- yzma:bench start text/cpu/arm64/yzma -->
### CPU, arm64, Arduino UnoQ
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"yzma","label":"Arduino UnoQ","tokens_per_second":31.96,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-4   	      12	 966081054 ns/op	        31.05 tokens/s
BenchmarkInference-4   	      12	 932465305 ns/op	        32.17 tokens/s
BenchmarkInference-4   	      12	 941589886 ns/op	        31.86 tokens/s
BenchmarkInference-4   	      12	 931413098 ns/op	        32.21 tokens/s
BenchmarkInference-4   	      12	 938804481 ns/op	        31.96 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	58.437s
```

</details>
<!-- yzma:bench end text/cpu/arm64/yzma -->

<!-- yzma:bench start text/cuda/amd64/i9-13900hx/cuda0 -->
### CUDA, amd64, Intel Core i9-13900HX, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"i9-13900hx","device":"CUDA0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":853.4,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 853.4 tokens a second.

<details><summary>The device</summary>

```
Thu Sep 24 08:32:51 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   51C    P0             18W /  115W |      15MiB /   8188MiB |      1%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkInference-32    	     333	  35387764 ns/op	       847.8 tokens/s
BenchmarkInference-32    	     339	  35300312 ns/op	       849.9 tokens/s
BenchmarkInference-32    	     339	  35153057 ns/op	       853.4 tokens/s
BenchmarkInference-32    	     340	  35067881 ns/op	       855.5 tokens/s
BenchmarkInference-32    	     340	  35153941 ns/op	       853.4 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	65.310s
```

</details>
<!-- yzma:bench end text/cuda/amd64/i9-13900hx/cuda0 -->

<!-- yzma:bench start text/cuda/arm64/localhost/cuda0 -->
### CUDA, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"arm64","machine":"localhost","device":"CUDA0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":191.6,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

<details><summary>The device</summary>

```
Wed Sep 23 10:46:52 2026       
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
BenchmarkInference-6   	      74	 158071388 ns/op	       189.8 tokens/s
BenchmarkInference-6   	      72	 156535194 ns/op	       191.7 tokens/s
BenchmarkInference-6   	      74	 156677079 ns/op	       191.5 tokens/s
BenchmarkInference-6   	      74	 156541354 ns/op	       191.6 tokens/s
BenchmarkInference-6   	      74	 156442212 ns/op	       191.8 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.084s
```

</details>
<!-- yzma:bench end text/cuda/arm64/localhost/cuda0 -->

<!-- yzma:bench start text/vulkan/amd64/i9-13900hx/vulkan0 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":95.52,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 95.5 tokens a second.

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
	driverVersion      = 595.84.0.0
	vendorID           = 0x10de
	deviceID           = 0x2860
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.84
	conformanceVersion = 1.4.3.3
	deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
	driverUUID         = 027cbdd1-c478-513e-af69-ebc9eaa7de87
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
BenchmarkInference-32    	      33	 314484019 ns/op	        95.39 tokens/s
BenchmarkInference-32    	      37	 314078421 ns/op	        95.52 tokens/s
BenchmarkInference-32    	      36	 326216341 ns/op	        91.96 tokens/s
BenchmarkInference-32    	      37	 313294349 ns/op	        95.76 tokens/s
BenchmarkInference-32    	      37	 310940713 ns/op	        96.48 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	65.493s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/i9-13900hx/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/i9-13900hx/vulkan1 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan1","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":744.2,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 744.2 tokens a second.

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
	driverVersion      = 595.84.0.0
	vendorID           = 0x10de
	deviceID           = 0x2860
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.84
	conformanceVersion = 1.4.3.3
	deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
	driverUUID         = 027cbdd1-c478-513e-af69-ebc9eaa7de87
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
BenchmarkInference-32    	     294	  40435435 ns/op	       741.9 tokens/s
BenchmarkInference-32    	     297	  40231602 ns/op	       745.7 tokens/s
BenchmarkInference-32    	     298	  40003017 ns/op	       749.9 tokens/s
BenchmarkInference-32    	     295	  40601893 ns/op	       738.9 tokens/s
BenchmarkInference-32    	     296	  40311326 ns/op	       744.2 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	79.250s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/i9-13900hx/vulkan1 -->

<!-- yzma:bench start text/vulkan/arm64/localhost/vulkan0 -->
### Vulkan, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"arm64","machine":"localhost","device":"Vulkan0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":177.5,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

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
BenchmarkInference-6   	      67	 174882641 ns/op	       171.5 tokens/s
BenchmarkInference-6   	      61	 169256059 ns/op	       177.2 tokens/s
BenchmarkInference-6   	      62	 166685163 ns/op	       180.0 tokens/s
BenchmarkInference-6   	      62	 169041983 ns/op	       177.5 tokens/s
BenchmarkInference-6   	      72	 166200892 ns/op	       180.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	69.999s
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
| CPU | amd64 | Intel Core i9-13900HX | - | 878.7 | b11146 | 2026-09-24 |
| CPU | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | - | 212.0 | b10964 | 2026-09-23 |
| CPU | arm64 | Raspberry Pi 4 Model B Rev 1.4 | - | 5.5 | b10964 | 2026-09-23 |
| CPU | arm64 | Arduino UnoQ | - | 4.2 | b10964 | 2026-09-23 |
| CUDA | amd64 | Intel Core i9-13900HX | CUDA0 | 2277.0 | b11146 | 2026-09-24 |
| CUDA | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | CUDA0 | 419.1 | b10964 | 2026-09-23 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan0 | 476.8 | b11146 | 2026-09-24 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan1 | 2133.0 | b11146 | 2026-09-24 |
| Vulkan | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | Vulkan0 | 416.1 | b10964 | 2026-09-23 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/i9-13900hx -->
### CPU, amd64, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":878.7,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 878.7 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192  -threadpool -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	      48	 268641515 ns/op	       876.9 tokens/s
BenchmarkMultimodalInference-32    	      44	 266624950 ns/op	       878.7 tokens/s
BenchmarkMultimodalInference-32    	      37	 280456836 ns/op	       853.4 tokens/s
BenchmarkMultimodalInference-32    	      56	 255300441 ns/op	       917.9 tokens/s
BenchmarkMultimodalInference-32    	      45	 258057383 ns/op	       910.0 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	64.335s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/i9-13900hx -->

<!-- yzma:bench start multimodal/cpu/arm64/localhost -->
### CPU, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"localhost","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":212,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-6   	      10	1141532842 ns/op	       206.7 tokens/s
BenchmarkMultimodalInference-6   	      10	1086205780 ns/op	       215.0 tokens/s
BenchmarkMultimodalInference-6   	      10	1100411231 ns/op	       212.0 tokens/s
BenchmarkMultimodalInference-6   	      10	1085658571 ns/op	       214.5 tokens/s
BenchmarkMultimodalInference-6   	       9	1165162499 ns/op	       203.0 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	55.891s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/localhost -->

<!-- yzma:bench start multimodal/cpu/arm64/raspberrypi -->
### CPU, arm64, Raspberry Pi 4 Model B Rev 1.4
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"raspberrypi","label":"Raspberry Pi 4 Model B Rev 1.4","tokens_per_second":5.522,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-4   	       1	38009244634 ns/op	         6.393 tokens/s
BenchmarkMultimodalInference-4   	       1	40374723120 ns/op	         5.647 tokens/s
BenchmarkMultimodalInference-4   	       1	46903363233 ns/op	         5.522 tokens/s
BenchmarkMultimodalInference-4   	       1	48361727646 ns/op	         4.611 tokens/s
BenchmarkMultimodalInference-4   	       1	51685423342 ns/op	         4.450 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	233.261s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/raspberrypi -->

<!-- yzma:bench start multimodal/cpu/arm64/yzma -->
### CPU, arm64, Arduino UnoQ
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"yzma","label":"Arduino UnoQ","tokens_per_second":4.185,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192   -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-4   	       1	56260051346 ns/op	         4.017 tokens/s
BenchmarkMultimodalInference-4   	       1	58066872101 ns/op	         4.546 tokens/s
BenchmarkMultimodalInference-4   	       1	57046166609 ns/op	         4.260 tokens/s
BenchmarkMultimodalInference-4   	       1	56626371282 ns/op	         4.185 tokens/s
BenchmarkMultimodalInference-4   	       1	56248529264 ns/op	         4.018 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	285.758s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/yzma -->

<!-- yzma:bench start multimodal/cuda/amd64/i9-13900hx/cuda0 -->
### CUDA, amd64, Intel Core i9-13900HX, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"i9-13900hx","device":"CUDA0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":2277,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 2277.0 tokens a second.

<details><summary>The device</summary>

```
Thu Sep 24 08:32:51 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   51C    P0             18W /  115W |      15MiB /   8188MiB |      1%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7732      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkMultimodalInference-32    	     116	 103334961 ns/op	      2270 tokens/s
BenchmarkMultimodalInference-32    	     100	 102473333 ns/op	      2288 tokens/s
BenchmarkMultimodalInference-32    	      99	 103212343 ns/op	      2277 tokens/s
BenchmarkMultimodalInference-32    	     100	 102320336 ns/op	      2294 tokens/s
BenchmarkMultimodalInference-32    	     100	 103537540 ns/op	      2267 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	55.205s
```

</details>
<!-- yzma:bench end multimodal/cuda/amd64/i9-13900hx/cuda0 -->

<!-- yzma:bench start multimodal/cuda/arm64/localhost/cuda0 -->
### CUDA, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"arm64","machine":"localhost","device":"CUDA0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":419.1,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

<details><summary>The device</summary>

```
Wed Sep 23 10:46:52 2026       
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
BenchmarkMultimodalInference-6   	      20	 567053354 ns/op	       417.1 tokens/s
BenchmarkMultimodalInference-6   	      22	 800435411 ns/op	       357.0 tokens/s
BenchmarkMultimodalInference-6   	      21	 545664945 ns/op	       428.4 tokens/s
BenchmarkMultimodalInference-6   	      22	 569710673 ns/op	       419.1 tokens/s
BenchmarkMultimodalInference-6   	      21	 550007790 ns/op	       431.4 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	66.513s
```

</details>
<!-- yzma:bench end multimodal/cuda/arm64/localhost/cuda0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/i9-13900hx/vulkan0 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":476.8,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 476.8 tokens a second.

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
	driverVersion      = 595.84.0.0
	vendorID           = 0x10de
	deviceID           = 0x2860
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.84
	conformanceVersion = 1.4.3.3
	deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
	driverUUID         = 027cbdd1-c478-513e-af69-ebc9eaa7de87
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
BenchmarkMultimodalInference-32    	      26	 458706669 ns/op	       505.4 tokens/s
BenchmarkMultimodalInference-32    	      22	 500817466 ns/op	       471.0 tokens/s
BenchmarkMultimodalInference-32    	      33	 493483073 ns/op	       476.8 tokens/s
BenchmarkMultimodalInference-32    	      25	 493515625 ns/op	       477.1 tokens/s
BenchmarkMultimodalInference-32    	      24	 510538972 ns/op	       463.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	70.448s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/i9-13900hx/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/i9-13900hx/vulkan1 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan1","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":2133,"llamacpp":"b11146","yzma":"1.28.0","date":"2026-09-24"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 2133.0 tokens a second.

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
	driverVersion      = 595.84.0.0
	vendorID           = 0x10de
	deviceID           = 0x2860
	deviceType         = PHYSICAL_DEVICE_TYPE_DISCRETE_GPU
	deviceName         = NVIDIA GeForce RTX 4070 Laptop GPU
	driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
	driverName         = NVIDIA
	driverInfo         = 595.84
	conformanceVersion = 1.4.3.3
	deviceUUID         = 7e611089-1272-699d-8985-ab84fef4311e
	driverUUID         = 027cbdd1-c478-513e-af69-ebc9eaa7de87
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
BenchmarkMultimodalInference-32    	     108	 109258887 ns/op	      2135 tokens/s
BenchmarkMultimodalInference-32    	     100	 110990663 ns/op	      2122 tokens/s
BenchmarkMultimodalInference-32    	     100	 110075571 ns/op	      2133 tokens/s
BenchmarkMultimodalInference-32    	     100	 108413462 ns/op	      2155 tokens/s
BenchmarkMultimodalInference-32    	     100	 110280568 ns/op	      2127 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	63.073s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/i9-13900hx/vulkan1 -->

<!-- yzma:bench start multimodal/vulkan/arm64/localhost/vulkan0 -->
### Vulkan, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"arm64","machine":"localhost","device":"Vulkan0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":416.1,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-23"} -->

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
BenchmarkMultimodalInference-6   	      18	 563675996 ns/op	       416.1 tokens/s
BenchmarkMultimodalInference-6   	      22	 538901458 ns/op	       430.9 tokens/s
BenchmarkMultimodalInference-6   	      18	 571077486 ns/op	       413.6 tokens/s
BenchmarkMultimodalInference-6   	      18	 580751599 ns/op	       409.1 tokens/s
BenchmarkMultimodalInference-6   	      22	 540216471 ns/op	       431.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	60.272s
```

</details>
<!-- yzma:bench end multimodal/vulkan/arm64/localhost/vulkan0 -->
