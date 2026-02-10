# Installation

Here is information on how to install `yzma`.

## Prerequisites

### Linux - CPU

No extra installation required to use the CPU on any Linux computer.

### Linux - CUDA

If you want to use a GPU with CUDA on a Linux machine, you will need to install the CUDA drivers.

See https://docs.nvidia.com/cuda/cuda-installation-guide-linux/

### Linux - Vulkan

To use Vulkan on your Linux system, your will also need to install the Vulkan drivers. For example:

```shell
sudo apt install -y mesa-vulkan-drivers vulkan-tools
```

### NVIDIA Jetson Orin

To use CUDA with the GPU on your [NVIDIA Jetson Orin](https://www.nvidia.com/en-us/autonomous-machines/embedded-systems/jetson-orin/nano-super-developer-kit/) you should install the latest version of the Jetpack software for your device.

To use Vulkan with the GPU on your Jetson Orin, you will also need to also update the GLIBC shared libraries:

```shell
sudo add-apt-repository ppa:ubuntu-toolchain-r/test
sudo apt-get update
sudo apt-get install --only-upgrade libstdc++6
```

### macOS

No extra installation required to use the GPU on any M-series computer.

### Windows - CPU

No extra installation required to use the CPU on any Windows computer.

### Windows - CUDA

If you want to use a GPU on your Windows machine, you will need to install the CUDA drivers.

See https://docs.nvidia.com/cuda/cuda-installation-guide-microsoft-windows/

### Windows - Vulkan

To use Vulkan, you will need to install the Vulkan SDK.

https://vulkan.lunarg.com/doc/sdk/latest/windows/getting_started.html

## Quick Start

### Install `yzma` command

Install the `yzma` command line tool.

```
go install github.com/hybridgroup/yzma/cmd/yzma@latest
```

For more info, see the [`yzma` command documentation](./cmd/yzma/README.md).

### Install `llama.cpp` libraries

Use the `yzma` command to install the `llama.cpp` libraries on your local machine.

Decide where you want put the files for your local installation. Do you have a GPU with either CUDA or Vulkan installed? You can use the `-processor` flag to download the hardware accelerated version for your system configuration.

For CPU only installation, run the following command:

```
yzma install --lib /path/to/lib
```

To install one of the hardware accelerated versions for your configuration such as CUDA, instead run the following command:

```
yzma install --lib /path/to/lib --processor cuda
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### Next steps

Now the installation is complete. Try running one of the example programs!

## Manual installation

If you prefer a manual installation, you can obtain most of the prebuilt `llama.cpp` binaries from here:

https://github.com/ggml-org/llama.cpp/releases

We also have binaries available for Ubuntu CUDA and Vulkan for arm64 located here:

https://github.com/hybridgroup/llama-cpp-builder/releases

### Installing the prebuilt binaries (manual)

If you do not use the `yzma` installer, you must download and extract the library files into a directory on your local machine.

#### Linux

For Linux, they have the `.so` file extension. For example, `libllama.so`, `libmtmd.so` and so on.

***Important Note***
You currently need to set the `YZMA_LIB` env variable to the directory with your llama.cpp library files. For example:

```shell
export YZMA_LIB=/home/ron/Development/yzma/lib
```

#### macOS

For macOS, the `llama.cpp` binaries have a `.dylib` file extension. For example, `libllama.dylib`, `libmtmd.dylib` and so on. You do not need the other downloaded files to use the `llama.cpp` libraries with `yzma`.

***Important Note***
You currently need to set the `YZMA_LIB` env variable to the directory with your `llama.cpp` library files. For example:

```shell
export YZMA_LIB=/home/ron/Development/yzma/lib
```

#### Windows

On Windows, the `llama.cpp` binaries have the `.dll` file extension. For example, `llama.dll`, `mtmd.dll` and so on.

You will also need to download the `cudart` files from the same location as the other `llama.cpp` libraries when using CUDA on Windows.

***Important Note***
You currently need to set the `YZMA_LIB` env variable to the directory with your `llama.cpp` library files. For example:

```shell
set YZMA_LIB=C:\yzma\lib
```
## Programmatic Installation

Want to use Go code to install the `llama.cpp` precompiled binaries from within your own application? We have the `download` package for that!

Check out the [installer example code](./examples/installer/).
