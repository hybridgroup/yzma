# Installation

Here is information on how to install `yzma`.

## Install `yzma` command

First install the `yzma` command line tool. You can then use it to install the `llama.cpp` libraries for your platform.

```
go install github.com/hybridgroup/yzma/cmd/yzma@latest
```

For more info, see the [`yzma` command documentation](./cmd/yzma/README.md).

## Install `llama.cpp` libraries

Now, using the `yzma` command, you can install the `llama.cpp` libraries. Follow the instructions for your system:

- [Linux - CPU](#linux-cpu)
- [Linux - CUDA](#linux-cuda)
- [Linux - Vulkan](#linux-vulkan)
- [NVIDIA Jetson Orin](#nvidia-jetson-orin)
- [Raspberry Pi 4/5](#raspberry-pi)
- [macOS](#macos)
- [Windows - CPU](#windows-cpu)
- [Windows - CUDA](#windows-cuda)
- [Windows - Vulkan](#windows-vulkan)

### Linux CPU

Decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### Linux CUDA

If you want to use a GPU with CUDA on a Linux machine, you will need to install the CUDA drivers.

See https://docs.nvidia.com/cuda/cuda-installation-guide-linux/

Once that is complete, decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor cuda
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### Linux Vulkan

To use Vulkan on your Linux system, your will also need to install the Vulkan drivers. For example:

```shell
sudo apt install -y mesa-vulkan-drivers vulkan-tools
```

Once that is complete, decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor vulkan
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### NVIDIA Jetson Orin

To the GPU on your [NVIDIA Jetson Orin](https://www.nvidia.com/en-us/autonomous-machines/embedded-systems/jetson-orin/nano-super-developer-kit/) you should install the latest version of the Jetpack software for your device.

#### CUDA

Decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor cuda
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

#### Vulkan

To use Vulkan with the GPU on your Jetson Orin, you will also need to also update the GLIBC shared libraries:

```shell
sudo add-apt-repository ppa:ubuntu-toolchain-r/test
sudo apt-get update
sudo apt-get install --only-upgrade libstdc++6
```

Once that is complete, decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor vulkan
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### Raspberry Pi

You can run `yzma` on a Raspberry Pi 4 or 5.

#### Raspberry Pi OS (64-bit)

If you are running the latest version of the Raspberry Pi OS, decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor cpu --os trixie
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

#### Raspberry Pi OS (Legacy, 64-bit)

If you are running an older version of the Raspberry Pi OS, decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor cpu --os bookworm
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### macOS

Decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### Windows CPU

Decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor cuda
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### Windows CUDA

If you want to use a GPU on your Windows machine, you will need to install the CUDA drivers.

See https://docs.nvidia.com/cuda/cuda-installation-guide-microsoft-windows/

Decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor cuda
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

### Windows Vulkan

To use Vulkan, you will need to install the Vulkan SDK.

https://vulkan.lunarg.com/doc/sdk/latest/windows/getting_started.html

Decide where you want put the files for your local installation, then run the following command:

```
yzma install --lib /path/to/lib --processor vulkan
```

To complete your installation, follow any specific instructions for your operating system displayed by the results of the `yzma install` command.

## Next steps

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
