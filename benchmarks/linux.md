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
| CPU | amd64 | Intel Core i9-13900HX | - | 252.1 | b10964 | 2026-09-16 |
| CPU | arm64 | Jetson Orin Nano Developer Kit 8GB | - | 59.2 | unknown | unknown |
| CPU | arm64 | Raspberry Pi 4 Model B Rev 1.4 8GB | - | 32.7 | unknown | unknown |
| CUDA | amd64 | Intel Core i9-13900HX | CUDA0 | 842.5 | unknown | unknown |
| CUDA | arm64 | Jetson Orin Nano Developer Kit 8GB | CUDA0 | 138.9 | unknown | unknown |
| ROCm | amd64 | AMD EPYC 7443P | ROCm0 | 494.1 | unknown | unknown |
| Vulkan | amd64 | AMD EPYC 7443P | Vulkan0 | 849.0 | unknown | unknown |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan0 | 92.0 | b10964 | 2026-09-16 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan1 | 731.4 | b10964 | 2026-09-16 |
| Vulkan | arm64 | Jetson Orin Nano Developer Kit 8GB | Vulkan0 | 135.2 | unknown | unknown |
<!-- yzma:bench table end text -->

<!-- yzma:bench start text/cpu/amd64/i9-13900hx -->
### CPU, amd64, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"amd64","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":252.1,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 252.1 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=8192 -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	     102	 117707523 ns/op	       254.9 tokens/s
BenchmarkInference-32    	      93	 117272681 ns/op	       255.8 tokens/s
BenchmarkInference-32    	      91	 118984634 ns/op	       252.1 tokens/s
BenchmarkInference-32    	      84	 130386396 ns/op	       230.1 tokens/s
BenchmarkInference-32    	      82	 126130546 ns/op	       237.8 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	62.306s
```

</details>
<!-- yzma:bench end text/cpu/amd64/i9-13900hx -->

<!-- yzma:bench start text/cpu/arm64/jetson-orin-nano -->
### CPU, arm64, Jetson Orin Nano Developer Kit 8GB
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"jetson-orin-nano","label":"Jetson Orin Nano Developer Kit 8GB","cpu":"ARMv8 Processor rev 1 (v8l)","tokens_per_second":59.2} -->

ARMv8 Processor rev 1 (v8l). 59.2 tokens a second.

<details><summary>The device</summary>

```
ron@ubuntu:~/yzma/pkg/mtmd$ vulkaninfo --summary
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.204
...
Devices:
========
GPU0:
        apiVersion         = 4206843 (1.3.251)
        driverVersion      = 2265006080 (0x87014000)
        vendorID           = 0x10de
        deviceID           = 0x97ba03d7
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = NVIDIA Tegra Orin (nvgpu)
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 540.5.0
        conformanceVersion = 1.3.6.0
        deviceUUID         = 1388f9e0-987e-54a0-908f-6a30d8fd5f29
        driverUUID         = ed5ba772-f592-5949-9d1f-236f7ad81bcc
```

</details>

<details><summary>The output of go test</summary>

```
ron@ubuntu:~/yzma/pkg/llama$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="CPU"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkInference-6          43         432432689 ns/op                69.37 tokens/s
BenchmarkInference-6          20         506747397 ns/op                59.20 tokens/s
BenchmarkInference-6          21         514736186 ns/op                58.28 tokens/s
BenchmarkInference-6          27         496646058 ns/op                60.41 tokens/s
BenchmarkInference-6          22         519434233 ns/op                57.76 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   68.009s
```

</details>
<!-- yzma:bench end text/cpu/arm64/jetson-orin-nano -->

<!-- yzma:bench start text/cpu/arm64/raspberry-pi-4 -->
### CPU, arm64, Raspberry Pi 4 Model B Rev 1.4 8GB
<!-- yzma:bench meta {"suite":"text","backend":"cpu","arch":"arm64","machine":"raspberry-pi-4","label":"Raspberry Pi 4 Model B Rev 1.4 8GB","tokens_per_second":32.67} -->

<details><summary>The output of go test</summary>

```
ron@raspberrypi:~/yzma/pkg/llama $ go test -benchtime=10s -count=5 -run=nada -bench .
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
BenchmarkInference-4          15         893788634 ns/op                33.56 tokens/s
BenchmarkInference-4          12         923948131 ns/op                32.47 tokens/s
BenchmarkInference-4          12         918284434 ns/op                32.67 tokens/s
BenchmarkInference-4          12         918693617 ns/op                32.66 tokens/s
BenchmarkInference-4          12         917186754 ns/op                32.71 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   64.583s
```

</details>
<!-- yzma:bench end text/cpu/arm64/raspberry-pi-4 -->

<!-- yzma:bench start text/cuda/amd64/i9-13900hx/cuda0 -->
### CUDA, amd64, Intel Core i9-13900HX, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"amd64","machine":"i9-13900hx","device":"CUDA0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":842.5} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 842.5 tokens a second.

<details><summary>The device</summary>

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

</details>

<details><summary>The output of go test</summary>

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="CUDA0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                332          35746370 ns/op               839.2 tokens/s
BenchmarkInference-32                338          35529926 ns/op               844.4 tokens/s
BenchmarkInference-32                336          35614579 ns/op               842.4 tokens/s
BenchmarkInference-32                336          35609522 ns/op               842.5 tokens/s
BenchmarkInference-32                337          35550352 ns/op               843.9 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   67.491s
```

</details>
<!-- yzma:bench end text/cuda/amd64/i9-13900hx/cuda0 -->

<!-- yzma:bench start text/cuda/arm64/jetson-orin-nano/cuda0 -->
### CUDA, arm64, Jetson Orin Nano Developer Kit 8GB, CUDA0
<!-- yzma:bench meta {"suite":"text","backend":"cuda","arch":"arm64","machine":"jetson-orin-nano","device":"CUDA0","label":"Jetson Orin Nano Developer Kit 8GB","cpu":"ARMv8 Processor rev 1 (v8l)","tokens_per_second":138.9} -->

ARMv8 Processor rev 1 (v8l). 138.9 tokens a second.

<details><summary>The device</summary>

```
+---------------------------------------------------------------------------------------+
| NVIDIA-SMI 540.5.0                Driver Version: 540.5.0      CUDA Version: 12.6     |
|-----------------------------------------+----------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |         Memory-Usage | GPU-Util  Compute M. |
|                                         |                      |               MIG M. |
|=========================================+======================+======================|
|   0  Orin (nvgpu)                  N/A  | N/A              N/A |                  N/A |
| N/A   N/A  N/A               N/A /  N/A | Not Supported        |     N/A          N/A |
|                                         |                      |                  N/A |
+-----------------------------------------+----------------------+----------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="CUDA0"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkInference-6          51         222138094 ns/op               135.1 tokens/s
BenchmarkInference-6          52         216104925 ns/op               138.8 tokens/s
BenchmarkInference-6          54         215961553 ns/op               138.9 tokens/s
BenchmarkInference-6          52         215498575 ns/op               139.2 tokens/s
BenchmarkInference-6          52         214849130 ns/op               139.6 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   61.014s
```

</details>
<!-- yzma:bench end text/cuda/arm64/jetson-orin-nano/cuda0 -->

<!-- yzma:bench start text/rocm/amd64/epyc-7443p/rocm0 -->
### ROCm, amd64, AMD EPYC 7443P, ROCm0
<!-- yzma:bench meta {"suite":"text","backend":"rocm","arch":"amd64","machine":"epyc-7443p","device":"ROCm0","label":"AMD EPYC 7443P","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":494.1} -->

AMD EPYC 7443P 24-Core Processor. 494.1 tokens a second.

<details><summary>The device</summary>

```
amdgpu_top v0.11.2
┌──────────────────────────────────────────────────────────────────────────────┐
│GPU Name                                | PCI Bus        |    VRAM Usage    | │
│SCLK    MCLK    VDDGFX  Power           | GFX% UMC%Media%|     GTT Usage    | │
│GPU/MEM_T  Fan     Throttle_Status                                          | │
│------------------------------------------------------------------------------│
│#0  [AMD Radeon RX 7900 XTX   ](gfx1100)| 0000:86:00.0   |    26/ 24560 MiB | │
│   0MHz   96MHz   49mV   14/303W        |   0%   0%   0% |    15/128884 MiB | │
│ 40C/ 46C     0RPM []                                                       | │
└──────────────────────────────────────────────────────────────────────────────┘
┌┤ Processes ├─────────────────────────────────────────────────────────────────┐
│┌┤ #0  AMD Radeon RX 7900 XTX ├──────────────────────────────────────────────┐│
││ Name            |  PID  |KFD| VRAM | GTT  |CPU |GFX |COMP|DMA |VCNU|       ││
││ kronk           | 411062|   |    0M|    2M|  1%|  0%|  0%|  0%|  0%|       ││
││ amdgpu_top      | 589729|   |    0M|    2M|  1%|  0%|  0%|  0%|  0%|       ││
│└────────────────────────────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────────────────────┘
```

</details>

<details><summary>The output of go test</summary>

```
go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="rocm0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD EPYC 7443P 24-Core Processor
BenchmarkInference-48    	     194	  60798061 ns/op	       493.4 tokens/s
BenchmarkInference-48    	     196	  60271732 ns/op	       497.7 tokens/s
BenchmarkInference-48    	     198	  60255594 ns/op	       497.9 tokens/s
BenchmarkInference-48    	     195	  60948909 ns/op	       492.2 tokens/s
BenchmarkInference-48    	     198	  60715718 ns/op	       494.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	60.887s
```

</details>
<!-- yzma:bench end text/rocm/amd64/epyc-7443p/rocm0 -->

<!-- yzma:bench start text/vulkan/amd64/epyc-7443p/vulkan0 -->
### Vulkan, amd64, AMD EPYC 7443P, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"epyc-7443p","device":"Vulkan0","label":"AMD EPYC 7443P","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":849} -->

AMD EPYC 7443P 24-Core Processor. 849.0 tokens a second.

<details><summary>The device</summary>

```
amdgpu_top v0.11.2
┌──────────────────────────────────────────────────────────────────────────────┐
│GPU Name                                | PCI Bus        |    VRAM Usage    | │
│SCLK    MCLK    VDDGFX  Power           | GFX% UMC%Media%|     GTT Usage    | │
│GPU/MEM_T  Fan     Throttle_Status                                          | │
│------------------------------------------------------------------------------│
│#0  [AMD Radeon RX 7900 XTX   ](gfx1100)| 0000:86:00.0   |    26/ 24560 MiB | │
│   0MHz   96MHz   49mV   14/303W        |   0%   0%   0% |    15/128884 MiB | │
│ 40C/ 46C     0RPM []                                                       | │
└──────────────────────────────────────────────────────────────────────────────┘
┌┤ Processes ├─────────────────────────────────────────────────────────────────┐
│┌┤ #0  AMD Radeon RX 7900 XTX ├──────────────────────────────────────────────┐│
││ Name            |  PID  |KFD| VRAM | GTT  |CPU |GFX |COMP|DMA |VCNU|       ││
││ kronk           | 411062|   |    0M|    2M|  1%|  0%|  0%|  0%|  0%|       ││
││ amdgpu_top      | 589729|   |    0M|    2M|  1%|  0%|  0%|  0%|  0%|       ││
│└────────────────────────────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────────────────────┘
```

</details>

<details><summary>The output of go test</summary>

```
go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="vulkan0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: AMD EPYC 7443P 24-Core Processor
BenchmarkInference-48    	     328	  36234037 ns/op	       828.0 tokens/s
BenchmarkInference-48    	     339	  35194859 ns/op	       852.4 tokens/s
BenchmarkInference-48    	     333	  35395438 ns/op	       847.6 tokens/s
BenchmarkInference-48    	     338	  35334138 ns/op	       849.0 tokens/s
BenchmarkInference-48    	     339	  35255138 ns/op	       850.9 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	61.232s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/epyc-7443p/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/i9-13900hx/vulkan0 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":92,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 92.0 tokens a second.

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
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=Vulkan0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	      31	 331359442 ns/op	        90.54 tokens/s
BenchmarkInference-32    	      38	 321097940 ns/op	        93.43 tokens/s
BenchmarkInference-32    	      36	 323849331 ns/op	        92.64 tokens/s
BenchmarkInference-32    	      32	 328564902 ns/op	        91.31 tokens/s
BenchmarkInference-32    	      34	 326092203 ns/op	        92.00 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	65.823s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/i9-13900hx/vulkan0 -->

<!-- yzma:bench start text/vulkan/amd64/i9-13900hx/vulkan1 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan1
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan1","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":731.4,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 731.4 tokens a second.

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
$ cd pkg/llama && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkInference -nctx=32000 -device=Vulkan1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32    	     294	  40291716 ns/op	       744.6 tokens/s
BenchmarkInference-32    	     292	  40873189 ns/op	       734.0 tokens/s
BenchmarkInference-32    	     290	  41069232 ns/op	       730.5 tokens/s
BenchmarkInference-32    	     289	  41194413 ns/op	       728.3 tokens/s
BenchmarkInference-32    	     290	  41018781 ns/op	       731.4 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/llama	78.520s
```

</details>
<!-- yzma:bench end text/vulkan/amd64/i9-13900hx/vulkan1 -->

<!-- yzma:bench start text/vulkan/arm64/jetson-orin-nano/vulkan0 -->
### Vulkan, arm64, Jetson Orin Nano Developer Kit 8GB, Vulkan0
<!-- yzma:bench meta {"suite":"text","backend":"vulkan","arch":"arm64","machine":"jetson-orin-nano","device":"Vulkan0","label":"Jetson Orin Nano Developer Kit 8GB","cpu":"ARMv8 Processor rev 1 (v8l)","tokens_per_second":135.2} -->

ARMv8 Processor rev 1 (v8l). 135.2 tokens a second.

<details><summary>The device</summary>

```
ron@ubuntu:~/yzma/pkg/mtmd$ vulkaninfo --summary
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.204
...
Devices:
========
GPU0:
        apiVersion         = 4206843 (1.3.251)
        driverVersion      = 2265006080 (0x87014000)
        vendorID           = 0x10de
        deviceID           = 0x97ba03d7
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = NVIDIA Tegra Orin (nvgpu)
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 540.5.0
        conformanceVersion = 1.3.6.0
        deviceUUID         = 1388f9e0-987e-54a0-908f-6a30d8fd5f29
        driverUUID         = ed5ba772-f592-5949-9d1f-236f7ad81bcc
```

</details>

<details><summary>The output of go test</summary>

```
ron@ubuntu:~/yzma/pkg/llama$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="VULKAN0"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkInference-6          52         222098600 ns/op               135.1 tokens/s
BenchmarkInference-6          52         222072877 ns/op               135.1 tokens/s
BenchmarkInference-6          54         219825013 ns/op               136.5 tokens/s
BenchmarkInference-6          52         220919304 ns/op               135.8 tokens/s
BenchmarkInference-6          54         221925680 ns/op               135.2 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   63.318s
```

</details>
<!-- yzma:bench end text/vulkan/arm64/jetson-orin-nano/vulkan0 -->

## Multimodal model benchmarks

The model is
[Qwen3-VL-2B-Instruct.Q4_K_M.gguf](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.Q4_K_M.gguf)
with its
[projector](https://huggingface.co/mradermacher/Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf).
The code is [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go).

<!-- yzma:bench table multimodal -->
| Backend | Arch | Machine | Device | Tokens a second | llama.cpp | Date |
| --- | --- | --- | --- | --- | --- | --- |
| CPU | amd64 | Intel Core i9-13900HX | - | 67.5 | b10964 | 2026-09-16 |
| CPU | arm64 | Jetson Orin Nano Developer Kit 8GB | - | 15.4 | unknown | unknown |
| CPU | arm64 | Raspberry Pi 4 Model B Rev 1.4 8GB | - | 5.9 | unknown | unknown |
| CUDA | amd64 | Intel Core i9-13900HX | CUDA0 | 1114.0 | unknown | unknown |
| CUDA | arm64 | Jetson Orin Nano Developer Kit 8GB | CUDA0 | 127.6 | unknown | unknown |
| ROCm | amd64 | AMD EPYC 7443P | ROCm0 | 961.5 | unknown | unknown |
| Vulkan | amd64 | AMD EPYC 7443P | Vulkan0 | 1179.0 | unknown | unknown |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan0 | 62.5 | b10964 | 2026-09-16 |
| Vulkan | amd64 | Intel Core i9-13900HX | Vulkan1 | 754.9 | b10964 | 2026-09-16 |
| Vulkan | arm64 | Jetson Orin Nano Developer Kit 8GB | Vulkan0 | 82.4 | unknown | unknown |
<!-- yzma:bench table end multimodal -->

<!-- yzma:bench start multimodal/cpu/amd64/i9-13900hx -->
### CPU, amd64, Intel Core i9-13900HX
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"amd64","machine":"i9-13900hx","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":67.46,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 67.5 tokens a second.

<details><summary>The output of go test</summary>

```
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=8192 -device=CPU
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	       1	15340603536 ns/op	        71.97 tokens/s
BenchmarkMultimodalInference-32    	       1	16647937836 ns/op	        67.46 tokens/s
BenchmarkMultimodalInference-32    	       1	20618641014 ns/op	        60.43 tokens/s
BenchmarkMultimodalInference-32    	       1	16358869871 ns/op	        67.98 tokens/s
BenchmarkMultimodalInference-32    	       1	20431101072 ns/op	        59.86 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	92.132s
```

</details>
<!-- yzma:bench end multimodal/cpu/amd64/i9-13900hx -->

<!-- yzma:bench start multimodal/cpu/arm64/jetson-orin-nano -->
### CPU, arm64, Jetson Orin Nano Developer Kit 8GB
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"jetson-orin-nano","label":"Jetson Orin Nano Developer Kit 8GB","cpu":"ARMv8 Processor rev 1 (v8l)","tokens_per_second":15.37} -->

ARMv8 Processor rev 1 (v8l). 15.4 tokens a second.

<details><summary>The device</summary>

```
ron@ubuntu:~/yzma/pkg/mtmd$ vulkaninfo --summary
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.204
...
Devices:
========
GPU0:
        apiVersion         = 4206843 (1.3.251)
        driverVersion      = 2265006080 (0x87014000)
        vendorID           = 0x10de
        deviceID           = 0x97ba03d7
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = NVIDIA Tegra Orin (nvgpu)
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 540.5.0
        conformanceVersion = 1.3.6.0
        deviceUUID         = 1388f9e0-987e-54a0-908f-6a30d8fd5f29
        driverUUID         = ed5ba772-f592-5949-9d1f-236f7ad81bcc
```

</details>

<details><summary>The output of go test</summary>

```
ron@ubuntu:~/yzma/pkg/mtmd$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="CPU"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkMultimodalInference-6                 1        72233629960 ns/op               15.03 tokens/s
BenchmarkMultimodalInference-6                 1        75555489707 ns/op               15.37 tokens/s
BenchmarkMultimodalInference-6                 1        87238792057 ns/op               14.65 tokens/s
BenchmarkMultimodalInference-6                 1        71406835155 ns/op               15.70 tokens/s
BenchmarkMultimodalInference-6                 1        70659234723 ns/op               15.74 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    383.358s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/jetson-orin-nano -->

<!-- yzma:bench start multimodal/cpu/arm64/raspberry-pi-4 -->
### CPU, arm64, Raspberry Pi 4 Model B Rev 1.4 8GB
<!-- yzma:bench meta {"suite":"multimodal","backend":"cpu","arch":"arm64","machine":"raspberry-pi-4","label":"Raspberry Pi 4 Model B Rev 1.4 8GB","tokens_per_second":5.917} -->

NOTE: this device has less memory, thus these numbers use the [SmolVLM2-500M-Video-Instruct-Q8_0](https://huggingface.co/ggml-org/SmolVLM2-500M-Video-Instruct-GGUF) model and its projector, not the model of the table above.

<details><summary>The output of go test</summary>

```
ron@raspberrypi:~/yzma/pkg/mtmd $ export YZMA_BENCHMARK_MMMODEL=/home/ron/models/SmolVLM2-500M-Video-Instruct-Q8_0.gguf
ron@raspberrypi:~/yzma/pkg/mtmd $ export YZMA_BENCHMARK_MMPROJ=/home/ron/models/mmproj-SmolVLM2-500M-Video-Instruct-Q8_0.gguf
ron@raspberrypi:~/yzma/pkg/mtmd $ go test -benchtime=10s -count=5 -run=nada -bench .
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
BenchmarkMultimodalInference-4                 1        50239133481 ns/op                6.748 tokens/s
BenchmarkMultimodalInference-4                 1        49358181828 ns/op                6.341 tokens/s
BenchmarkMultimodalInference-4                 1        48164506831 ns/op                5.917 tokens/s
BenchmarkMultimodalInference-4                 1        40171997080 ns/op                5.551 tokens/s
BenchmarkMultimodalInference-4                 1        41428165840 ns/op                5.504 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    243.876s
```

</details>
<!-- yzma:bench end multimodal/cpu/arm64/raspberry-pi-4 -->

<!-- yzma:bench start multimodal/cuda/amd64/i9-13900hx/cuda0 -->
### CUDA, amd64, Intel Core i9-13900HX, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"amd64","machine":"i9-13900hx","device":"CUDA0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":1114} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 1114.0 tokens a second.

<details><summary>The device</summary>

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

</details>

<details><summary>The output of go test</summary>

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="CUDA0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32               21         921205057 ns/op              1240 tokens/s
BenchmarkMultimodalInference-32               15        1043496530 ns/op              1114 tokens/s
BenchmarkMultimodalInference-32               18         939373857 ns/op              1219 tokens/s
BenchmarkMultimodalInference-32               14        1118362797 ns/op              1047 tokens/s
BenchmarkMultimodalInference-32                8        1336574088 ns/op               900.2 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    82.619s
```

</details>
<!-- yzma:bench end multimodal/cuda/amd64/i9-13900hx/cuda0 -->

<!-- yzma:bench start multimodal/cuda/arm64/jetson-orin-nano/cuda0 -->
### CUDA, arm64, Jetson Orin Nano Developer Kit 8GB, CUDA0
<!-- yzma:bench meta {"suite":"multimodal","backend":"cuda","arch":"arm64","machine":"jetson-orin-nano","device":"CUDA0","label":"Jetson Orin Nano Developer Kit 8GB","cpu":"ARMv8 Processor rev 1 (v8l)","tokens_per_second":127.6} -->

ARMv8 Processor rev 1 (v8l). 127.6 tokens a second.

<details><summary>The device</summary>

```
+---------------------------------------------------------------------------------------+
| NVIDIA-SMI 540.5.0                Driver Version: 540.5.0      CUDA Version: 12.6     |
|-----------------------------------------+----------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |         Memory-Usage | GPU-Util  Compute M. |
|                                         |                      |               MIG M. |
|=========================================+======================+======================|
|   0  Orin (nvgpu)                  N/A  | N/A              N/A |                  N/A |
| N/A   N/A  N/A               N/A /  N/A | Not Supported        |     N/A          N/A |
|                                         |                      |                  N/A |
+-----------------------------------------+----------------------+----------------------+
```

</details>

<details><summary>The output of go test</summary>

```
$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="CUDA0"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkMultimodalInference-6                 2        7077293280 ns/op               166.9 tokens/s
BenchmarkMultimodalInference-6                 2        8106794026 ns/op               150.8 tokens/s
BenchmarkMultimodalInference-6                 1        10837943077 ns/op              120.7 tokens/s
BenchmarkMultimodalInference-6                 1        12015033493 ns/op              112.1 tokens/s
BenchmarkMultimodalInference-6                 1        10055887615 ns/op              127.6 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    69.733s
```

</details>
<!-- yzma:bench end multimodal/cuda/arm64/jetson-orin-nano/cuda0 -->

<!-- yzma:bench start multimodal/rocm/amd64/epyc-7443p/rocm0 -->
### ROCm, amd64, AMD EPYC 7443P, ROCm0
<!-- yzma:bench meta {"suite":"multimodal","backend":"rocm","arch":"amd64","machine":"epyc-7443p","device":"ROCm0","label":"AMD EPYC 7443P","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":961.5} -->

AMD EPYC 7443P 24-Core Processor. 961.5 tokens a second.

<details><summary>The device</summary>

```
amdgpu_top v0.11.2
┌──────────────────────────────────────────────────────────────────────────────┐
│GPU Name                                | PCI Bus        |    VRAM Usage    | │
│SCLK    MCLK    VDDGFX  Power           | GFX% UMC%Media%|     GTT Usage    | │
│GPU/MEM_T  Fan     Throttle_Status                                          | │
│------------------------------------------------------------------------------│
│#0  [AMD Radeon RX 7900 XTX   ](gfx1100)| 0000:86:00.0   |    26/ 24560 MiB | │
│   0MHz   96MHz   49mV   14/303W        |   0%   0%   0% |    15/128884 MiB | │
│ 40C/ 46C     0RPM []                                                       | │
└──────────────────────────────────────────────────────────────────────────────┘
┌┤ Processes ├─────────────────────────────────────────────────────────────────┐
│┌┤ #0  AMD Radeon RX 7900 XTX ├──────────────────────────────────────────────┐│
││ Name            |  PID  |KFD| VRAM | GTT  |CPU |GFX |COMP|DMA |VCNU|       ││
││ kronk           | 411062|   |    0M|    2M|  1%|  0%|  0%|  0%|  0%|       ││
││ amdgpu_top      | 589729|   |    0M|    2M|  1%|  0%|  0%|  0%|  0%|       ││
│└────────────────────────────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────────────────────┘
```

</details>

<details><summary>The output of go test</summary>

```
go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="rocm0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD EPYC 7443P 24-Core Processor
BenchmarkMultimodalInference-48    	       9	1182597512 ns/op	       987.1 tokens/s
BenchmarkMultimodalInference-48    	      10	1241401135 ns/op	       961.6 tokens/s
BenchmarkMultimodalInference-48    	       8	1323004757 ns/op	       912.9 tokens/s
BenchmarkMultimodalInference-48    	      12	1241431410 ns/op	       961.5 tokens/s
BenchmarkMultimodalInference-48    	       8	1715075982 ns/op	       755.4 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	63.492s
```

</details>
<!-- yzma:bench end multimodal/rocm/amd64/epyc-7443p/rocm0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/epyc-7443p/vulkan0 -->
### Vulkan, amd64, AMD EPYC 7443P, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"epyc-7443p","device":"Vulkan0","label":"AMD EPYC 7443P","cpu":"AMD EPYC 7443P 24-Core Processor","tokens_per_second":1179} -->

AMD EPYC 7443P 24-Core Processor. 1179.0 tokens a second.

<details><summary>The device</summary>

```
amdgpu_top v0.11.2
┌──────────────────────────────────────────────────────────────────────────────┐
│GPU Name                                | PCI Bus        |    VRAM Usage    | │
│SCLK    MCLK    VDDGFX  Power           | GFX% UMC%Media%|     GTT Usage    | │
│GPU/MEM_T  Fan     Throttle_Status                                          | │
│------------------------------------------------------------------------------│
│#0  [AMD Radeon RX 7900 XTX   ](gfx1100)| 0000:86:00.0   |    26/ 24560 MiB | │
│   0MHz   96MHz   49mV   14/303W        |   0%   0%   0% |    15/128884 MiB | │
│ 40C/ 46C     0RPM []                                                       | │
└──────────────────────────────────────────────────────────────────────────────┘
┌┤ Processes ├─────────────────────────────────────────────────────────────────┐
│┌┤ #0  AMD Radeon RX 7900 XTX ├──────────────────────────────────────────────┐│
││ Name            |  PID  |KFD| VRAM | GTT  |CPU |GFX |COMP|DMA |VCNU|       ││
││ kronk           | 411062|   |    0M|    2M|  1%|  0%|  0%|  0%|  0%|       ││
││ amdgpu_top      | 589729|   |    0M|    2M|  1%|  0%|  0%|  0%|  0%|       ││
│└────────────────────────────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────────────────────┘
```

</details>

<details><summary>The output of go test</summary>

```
go test -benchtime=10s -count=5 -run=nada -bench . -nctx=32000 -device="vulkan0"
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: AMD EPYC 7443P 24-Core Processor
BenchmarkMultimodalInference-48    	       9	1147394253 ns/op	      1053 tokens/s
BenchmarkMultimodalInference-48    	      15	 941516811 ns/op	      1245 tokens/s
BenchmarkMultimodalInference-48    	      13	 924097033 ns/op	      1265 tokens/s
BenchmarkMultimodalInference-48    	      18	1018284301 ns/op	      1179 tokens/s
BenchmarkMultimodalInference-48    	      15	1022548971 ns/op	      1172 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	71.331s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/epyc-7443p/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/i9-13900hx/vulkan0 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan0","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":62.49,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 62.5 tokens a second.

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
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=Vulkan0
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	       1	29325408518 ns/op	        45.59 tokens/s
BenchmarkMultimodalInference-32    	       1	19091432501 ns/op	        65.47 tokens/s
BenchmarkMultimodalInference-32    	       1	20340714964 ns/op	        62.49 tokens/s
BenchmarkMultimodalInference-32    	       1	19056191976 ns/op	        65.54 tokens/s
BenchmarkMultimodalInference-32    	       1	22343760148 ns/op	        58.27 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	114.592s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/i9-13900hx/vulkan0 -->

<!-- yzma:bench start multimodal/vulkan/amd64/i9-13900hx/vulkan1 -->
### Vulkan, amd64, Intel Core i9-13900HX, Vulkan1
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"amd64","machine":"i9-13900hx","device":"Vulkan1","label":"Intel Core i9-13900HX","cpu":"13th Gen Intel(R) Core(TM) i9-13900HX","tokens_per_second":754.9,"llamacpp":"b10964","yzma":"1.27.0","date":"2026-09-16"} -->

13th Gen Intel(R) Core(TM) i9-13900HX. 754.9 tokens a second.

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
$ cd pkg/mtmd && go test -benchtime=10s -count=5 -run=nada -bench BenchmarkMultimodalInference -nctx=32000 -device=Vulkan1
goos: linux
goarch: amd64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkMultimodalInference-32    	       1	16223066476 ns/op	        68.85 tokens/s
BenchmarkMultimodalInference-32    	      16	1565100628 ns/op	       784.4 tokens/s
BenchmarkMultimodalInference-32    	       7	1630751597 ns/op	       754.9 tokens/s
BenchmarkMultimodalInference-32    	       7	1558894683 ns/op	       783.6 tokens/s
BenchmarkMultimodalInference-32    	       6	2056770062 ns/op	       629.6 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	79.713s
```

</details>
<!-- yzma:bench end multimodal/vulkan/amd64/i9-13900hx/vulkan1 -->

<!-- yzma:bench start multimodal/vulkan/arm64/jetson-orin-nano/vulkan0 -->
### Vulkan, arm64, Jetson Orin Nano Developer Kit 8GB, Vulkan0
<!-- yzma:bench meta {"suite":"multimodal","backend":"vulkan","arch":"arm64","machine":"jetson-orin-nano","device":"Vulkan0","label":"Jetson Orin Nano Developer Kit 8GB","cpu":"ARMv8 Processor rev 1 (v8l)","tokens_per_second":82.43} -->

ARMv8 Processor rev 1 (v8l). 82.4 tokens a second.

<details><summary>The device</summary>

```
ron@ubuntu:~/yzma/pkg/mtmd$ vulkaninfo --summary
==========
VULKANINFO
==========

Vulkan Instance Version: 1.3.204
...
Devices:
========
GPU0:
        apiVersion         = 4206843 (1.3.251)
        driverVersion      = 2265006080 (0x87014000)
        vendorID           = 0x10de
        deviceID           = 0x97ba03d7
        deviceType         = PHYSICAL_DEVICE_TYPE_INTEGRATED_GPU
        deviceName         = NVIDIA Tegra Orin (nvgpu)
        driverID           = DRIVER_ID_NVIDIA_PROPRIETARY
        driverName         = NVIDIA
        driverInfo         = 540.5.0
        conformanceVersion = 1.3.6.0
        deviceUUID         = 1388f9e0-987e-54a0-908f-6a30d8fd5f29
        driverUUID         = ed5ba772-f592-5949-9d1f-236f7ad81bcc
```

</details>

<details><summary>The output of go test</summary>

```
ron@ubuntu:~/yzma/pkg/mtmd$ go test -benchtime=10s -count=5 -run=nada -bench . -nctx=16000 -device="VULKAN0"
goos: linux
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: ARMv8 Processor rev 1 (v8l)
BenchmarkMultimodalInference-6                 1        13718208893 ns/op               81.13 tokens/s
BenchmarkMultimodalInference-6                 1        16724822437 ns/op               71.39 tokens/s
BenchmarkMultimodalInference-6                 1        13133369170 ns/op               84.14 tokens/s
BenchmarkMultimodalInference-6                 1        13515072899 ns/op               82.43 tokens/s
BenchmarkMultimodalInference-6                 1        12471954537 ns/op               87.24 tokens/s
PASS
ok      github.com/hybridgroup/yzma/pkg/mtmd    76.766s
```

</details>
<!-- yzma:bench end multimodal/vulkan/arm64/jetson-orin-nano/vulkan0 -->
