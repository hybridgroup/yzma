# Linux benchmarks

Benchmarks of yzma on Linux. Each table gives the median of five runs. The output
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
| CPU | amd64 | Intel Core i9-13900HX | - | 270.1 | b10964 | 2026-09-22 |
| CPU | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | - | 84.2 | b10964 | 2026-09-19 |
| CPU | arm64 | Raspberry Pi 4 Model B Rev 1.4 | - | 28.8 | b10964 | 2026-09-17 |
| CPU | arm64 | Arduino UnoQ | - | 32.0 | b10964 | 2026-09-19 |
| CUDA | amd64 | Intel Core i9-13900HX | CUDA0 | 852.6 | b10964 | 2026-09-22 |
| CUDA | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | CUDA0 | 190.5 | b10964 | 2026-09-19 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan0 | 95.9 | b10964 | 2026-09-22 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan1 | 746.2 | b10964 | 2026-09-22 |
| Vulkan | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | Vulkan0 | 183.2 | b10964 | 2026-09-19 |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/i9-13900hx -->
### CPU, amd64, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":270.1,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 270.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192  -threadpool -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	      87	 116222107 ns/op	       258.1 tokens/s
BenchmarkInference-32    	     100	 111062497 ns/op	       270.1 tokens/s
BenchmarkInference-32    	     100	 109065889 ns/op	       275.1 tokens/s
BenchmarkInference-32    	     100	 108390243 ns/op	       276.8 tokens/s
BenchmarkInference-32    	     100	 114162275 ns/op	       262.8 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	60.564s
```

</details>
<!-- yzma:bench end text/cpu/amd64/i9-13900hx -->

<!-- yzma:bench start text/cpu/arm64/localhost -->
### CPU, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"localhost","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":84.17,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-19"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-6   	      36	 347242000 ns/op	        86.40 tokens/s
BenchmarkInference-6   	      33	 351954953 ns/op	        85.24 tokens/s
BenchmarkInference-6   	      33	 356435241 ns/op	        84.17 tokens/s
BenchmarkInference-6   	      32	 359732485 ns/op	        83.40 tokens/s
BenchmarkInference-6   	      33	 357596293 ns/op	        83.89 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.214s
```

</details>
<!-- yzma:bench end text/cpu/arm64/localhost -->

<!-- yzma:bench start text/cpu/arm64/raspberrypi -->
### CPU, arm64, Raspberry Pi 4 Model B Rev 1.4
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"raspberrypi","label":"Raspberry Pi 4 Model B Rev 1.4","tokens_per_second":28.8,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-4   	      14	 917452696 ns/op	        32.70 tokens/s
BenchmarkInference-4   	      12	 996564566 ns/op	        30.10 tokens/s
BenchmarkInference-4   	      10	1041497747 ns/op	        28.80 tokens/s
BenchmarkInference-4   	      10	1055651177 ns/op	        28.42 tokens/s
BenchmarkInference-4   	      10	1071348666 ns/op	        28.00 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.556s
```

</details>
<!-- yzma:bench end text/cpu/arm64/raspberrypi -->

<!-- yzma:bench start text/cpu/arm64/yzma -->
### CPU, arm64, Arduino UnoQ
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"yzma","label":"Arduino UnoQ","tokens_per_second":32.01,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-19"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-4   	      12	 941140525 ns/op	        31.88 tokens/s
BenchmarkInference-4   	      12	 933477923 ns/op	        32.14 tokens/s
BenchmarkInference-4   	      12	 937238417 ns/op	        32.01 tokens/s
BenchmarkInference-4   	      12	 936656218 ns/op	        32.03 tokens/s
BenchmarkInference-4   	      12	 941273494 ns/op	        31.87 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	59.124s
```

</details>
<!-- yzma:bench end text/cpu/arm64/yzma -->

<!-- yzma:bench start text/cuda/amd64/i9-13900hx/cuda0 -->
### CUDA, amd64, Intel Core i9-13900HX, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"i9-13900hx","device":"CUDA0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":852.6,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 852.6 tokens a second.

<details><summary>The device</summary>

```
Tue Sep 22 21:42:10 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   60C    P0             23W /  115W |      16MiB /   8188MiB |     17%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkInference-32    	     333	  35480149 ns/op	       845.5 tokens/s
BenchmarkInference-32    	     340	  35225276 ns/op	       851.7 tokens/s
BenchmarkInference-32    	     340	  35185188 ns/op	       852.6 tokens/s
BenchmarkInference-32    	     339	  35169759 ns/op	       853.0 tokens/s
BenchmarkInference-32    	     340	  35185066 ns/op	       852.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	65.442s
```

</details>
<!-- yzma:bench end text/cuda/amd64/i9-13900hx/cuda0 -->

<!-- yzma:bench start text/cuda/arm64/localhost/cuda0 -->
### CUDA, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"arm64","machine":"localhost","device":"CUDA0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":190.5,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-19"} -->

<details><summary>The device</summary>

```
Sat Sep 19 08:01:34 2026       
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
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=CUDA0
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-6   	      74	 158433298 ns/op	       189.4 tokens/s
BenchmarkInference-6   	      73	 157411886 ns/op	       190.6 tokens/s
BenchmarkInference-6   	      74	 157470586 ns/op	       190.5 tokens/s
BenchmarkInference-6   	      74	 157282683 ns/op	       190.7 tokens/s
BenchmarkInference-6   	      73	 157831218 ns/op	       190.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	63.445s
```

</details>
<!-- yzma:bench end text/cuda/arm64/localhost/cuda0 -->

<!-- yzma:bench start text/vulkan/amd64/i9-13900hx/vulkan0 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":95.91,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 95.9 tokens a second.

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
BenchmarkInference-32    	      32	 314944244 ns/op	        95.25 tokens/s
BenchmarkInference-32    	      37	 313001122 ns/op	        95.85 tokens/s
BenchmarkInference-32    	      37	 311199789 ns/op	        96.40 tokens/s
BenchmarkInference-32    	      38	 312504542 ns/op	        96.00 tokens/s
BenchmarkInference-32    	      37	 312805713 ns/op	        95.91 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	65.276s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/i9-13900hx/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/i9-13900hx/vulkan1 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan1","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":746.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 746.2 tokens a second.

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
BenchmarkInference-32    	     292	  40563991 ns/op	       739.6 tokens/s
BenchmarkInference-32    	     302	  39778344 ns/op	       754.2 tokens/s
BenchmarkInference-32    	     294	  40198758 ns/op	       746.3 tokens/s
BenchmarkInference-32    	     298	  40366280 ns/op	       743.2 tokens/s
BenchmarkInference-32    	     294	  40202694 ns/op	       746.2 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	79.021s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/i9-13900hx/vulkan1 -->

<!-- yzma:bench start text/vulkan/arm64/localhost/vulkan0 -->
### Vulkan, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"arm64","machine":"localhost","device":"Vulkan0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":183.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-19"} -->

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
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=Vulkan0
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-6   	      14	 720641293 ns/op	        41.63 tokens/s
BenchmarkInference-6   	      66	 163608069 ns/op	       183.4 tokens/s
BenchmarkInference-6   	      64	 163069410 ns/op	       184.0 tokens/s
BenchmarkInference-6   	      67	 164635755 ns/op	       182.2 tokens/s
BenchmarkInference-6   	      73	 163729836 ns/op	       183.2 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	66.882s
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
| CPU | amd64 | Intel Core i9-13900HX | - | 856.1 | b10964 | 2026-09-22 |
| CPU | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | - | 138.2 | b10964 | 2026-09-19 |
| CPU | arm64 | Raspberry Pi 4 Model B Rev 1.4 | - | 3.5 | b10964 | 2026-09-17 |
| CPU | arm64 | Arduino UnoQ | - | 4.1 | b10964 | 2026-09-19 |
| CUDA | amd64 | Intel Core i9-13900HX | CUDA0 | 2297.0 | b10964 | 2026-09-22 |
| CUDA | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | CUDA0 | 423.0 | b10964 | 2026-09-19 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan0 | 449.2 | b10964 | 2026-09-22 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan1 | 2138.0 | b10964 | 2026-09-22 |
| Vulkan | arm64 | NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super | Vulkan0 | 427.2 | b10964 | 2026-09-19 |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/i9-13900hx -->
### CPU, amd64, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":856.1,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 856.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192  -threadpool -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	      39	 267396046 ns/op	       881.8 tokens/s
BenchmarkMultimodalInference-32    	      61	 281662694 ns/op	       845.8 tokens/s
BenchmarkMultimodalInference-32    	      55	 277539137 ns/op	       856.1 tokens/s
BenchmarkMultimodalInference-32    	      39	 266912013 ns/op	       881.7 tokens/s
BenchmarkMultimodalInference-32    	      46	 293252861 ns/op	       813.7 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	70.284s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/i9-13900hx -->

<!-- yzma:bench start multimodal/cpu/arm64/localhost -->
### CPU, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"localhost","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":138.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-19"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-6   	       7	1625863099 ns/op	       142.3 tokens/s
BenchmarkMultimodalInference-6   	       7	1689402565 ns/op	       138.2 tokens/s
BenchmarkMultimodalInference-6   	       7	1724905268 ns/op	       136.7 tokens/s
BenchmarkMultimodalInference-6   	       7	1660245148 ns/op	       139.2 tokens/s
BenchmarkMultimodalInference-6   	       6	1768562379 ns/op	       135.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	58.545s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/localhost -->

<!-- yzma:bench start multimodal/cpu/arm64/raspberrypi -->
### CPU, arm64, Raspberry Pi 4 Model B Rev 1.4
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"raspberrypi","label":"Raspberry Pi 4 Model B Rev 1.4","tokens_per_second":3.508,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-17"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-4   	       1	52235071588 ns/op	         4.403 tokens/s
BenchmarkMultimodalInference-4   	       1	60595582227 ns/op	         3.779 tokens/s
BenchmarkMultimodalInference-4   	       1	64996824736 ns/op	         3.508 tokens/s
BenchmarkMultimodalInference-4   	       1	68645379347 ns/op	         3.394 tokens/s
BenchmarkMultimodalInference-4   	       1	71142212836 ns/op	         3.233 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	325.622s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/raspberrypi -->

<!-- yzma:bench start multimodal/cpu/arm64/yzma -->
### CPU, arm64, Arduino UnoQ
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"yzma","label":"Arduino UnoQ","tokens_per_second":4.127,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-19"} -->

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-4   	       1	56697056172 ns/op	         4.127 tokens/s
BenchmarkMultimodalInference-4   	       1	56879974757 ns/op	         4.237 tokens/s
BenchmarkMultimodalInference-4   	       1	56321539181 ns/op	         4.030 tokens/s
BenchmarkMultimodalInference-4   	       1	57043481210 ns/op	         4.277 tokens/s
BenchmarkMultimodalInference-4   	       1	56457938553 ns/op	         4.092 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	286.836s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/yzma -->

<!-- yzma:bench start multimodal/cuda/amd64/i9-13900hx/cuda0 -->
### CUDA, amd64, Intel Core i9-13900HX, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"i9-13900hx","device":"CUDA0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":2297,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 2297.0 tokens a second.

<details><summary>The device</summary>

```
Tue Sep 22 21:42:10 2026       
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 595.84                 Driver Version: 595.84         CUDA Version: 13.2     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   60C    P0             23W /  115W |      16MiB /   8188MiB |     17%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+

+-----------------------------------------------------------------------------------------+
| Processes:                                                                              |
|  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
|        ID   ID                                                               Usage      |
|=========================================================================================|
|    0   N/A  N/A            7793      G   /usr/lib/xorg/Xorg                        4MiB |
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
BenchmarkMultimodalInference-32    	     118	 101893125 ns/op	      2296 tokens/s
BenchmarkMultimodalInference-32    	     100	 102045264 ns/op	      2297 tokens/s
BenchmarkMultimodalInference-32    	     100	 102900707 ns/op	      2279 tokens/s
BenchmarkMultimodalInference-32    	     120	  99002053 ns/op	      2344 tokens/s
BenchmarkMultimodalInference-32    	     100	 101842857 ns/op	      2297 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	56.797s
```

</details>
<!-- yzma:bench end multimodal/cuda/amd64/i9-13900hx/cuda0 -->

<!-- yzma:bench start multimodal/cuda/arm64/localhost/cuda0 -->
### CUDA, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"arm64","machine":"localhost","device":"CUDA0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":423,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-19"} -->

<details><summary>The device</summary>

```
Sat Sep 19 08:01:34 2026       
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
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=CUDA0
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-6   	      20	 564174666 ns/op	       416.1 tokens/s
BenchmarkMultimodalInference-6   	      21	 563107157 ns/op	       423.0 tokens/s
BenchmarkMultimodalInference-6   	      20	 564359298 ns/op	       420.5 tokens/s
BenchmarkMultimodalInference-6   	      21	 556031042 ns/op	       423.2 tokens/s
BenchmarkMultimodalInference-6   	      24	 544903652 ns/op	       429.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	61.198s
```

</details>
<!-- yzma:bench end multimodal/cuda/arm64/localhost/cuda0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/i9-13900hx/vulkan0 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":449.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 449.2 tokens a second.

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
BenchmarkMultimodalInference-32    	      20	 522260533 ns/op	       449.2 tokens/s
BenchmarkMultimodalInference-32    	      30	 533139125 ns/op	       447.3 tokens/s
BenchmarkMultimodalInference-32    	      24	 525356835 ns/op	       453.5 tokens/s
BenchmarkMultimodalInference-32    	      27	 467869487 ns/op	       497.8 tokens/s
BenchmarkMultimodalInference-32    	      20	 543390832 ns/op	       440.5 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	68.780s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/i9-13900hx/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/i9-13900hx/vulkan1 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan1","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":2138,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-22"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 2138.0 tokens a second.

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
BenchmarkMultimodalInference-32    	     108	 110992694 ns/op	      2116 tokens/s
BenchmarkMultimodalInference-32    	      94	 110149820 ns/op	      2128 tokens/s
BenchmarkMultimodalInference-32    	     100	 109019839 ns/op	      2138 tokens/s
BenchmarkMultimodalInference-32    	     100	 107798457 ns/op	      2164 tokens/s
BenchmarkMultimodalInference-32    	     100	 109163266 ns/op	      2148 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	62.195s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/i9-13900hx/vulkan1 -->

<!-- yzma:bench start multimodal/vulkan/arm64/localhost/vulkan0 -->
### Vulkan, arm64, NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"arm64","machine":"localhost","device":"Vulkan0","label":"NVIDIA Jetson Orin Nano Engineering Reference Developer Kit Super","tokens_per_second":427.2,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-19"} -->

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
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=Vulkan0
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-6   	       1	23310196573 ns/op	         9.610 tokens/s
BenchmarkMultimodalInference-6   	      21	 535231492 ns/op	       434.4 tokens/s
BenchmarkMultimodalInference-6   	      21	 552260375 ns/op	       427.2 tokens/s
BenchmarkMultimodalInference-6   	      20	 555362272 ns/op	       425.2 tokens/s
BenchmarkMultimodalInference-6   	      22	 544643297 ns/op	       430.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	73.711s
```

</details>
<!-- yzma:bench end multimodal/vulkan/arm64/localhost/vulkan0 -->
