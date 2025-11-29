# Installation

Here is information on how to install the requirements for running an application using `llama.cpp`.

## `yzma` command

You can install or update the `llama.cpp` libraries on your local machine using the `yzma` command line tool.

Install it like this:

```
go install github.com/hybridgroup/yzma/cmd/yzma@latest
```

Decide where you want put the files for your local installation, then run the `yzma install` command:

```
yzma install --lib /path/to/lib
```

If you want to install one of the hardware accelerated versions for your configuration such as CUDA:

```
yzma install --lib /path/to/lib --processor cuda
```

For more info, see the [`yzma` command documentation](./cmd/yzma/README.md).

## Programmatic Installation

Want to use Go code to install the `llama.cpp` precompiled binaries from within your own application? We have the `download` package for that!

Check out the [installer example code](./examples/installer/).

## Installation via bash script

We also have a helpful bash script you can use to download the latest binaries.

For example:

```
./download_llama.sh linux cuda
```

Once you have downloaded it, unzip the downloaded file to the correct location for your configuration.

## Linux

You will need to download the `llama.cpp` libraries for Linux. Choose the download option that matches your desired configuration.

### CPU

You can use the `yzma` installer program or obtain prebuilt `llama.cpp` libraries from here:

https://github.com/ggml-org/llama.cpp/releases

### CUDA

You will need to install the CUDA drivers for your Nvidia GPU.

See https://docs.nvidia.com/cuda/cuda-installation-guide-linux/

Once your drivers are installed, you can use the `yzma` installer program or obtain prebuilt `llama.cpp` binaries for Ubuntu 24.04 using the NVidia 12.9 drivers for CUDA acceleration from here:

https://github.com/hybridgroup/llama-cpp-builder/releases

### Vulkan

You will need to install Vulkan on your system. For example:

```shell
sudo apt install -y mesa-vulkan-drivers vulkan-tools
```

Once your drivers are installed, you can use the `yzma` installer program or obtain prebuilt `llama.cpp` binaries from here:

https://github.com/ggml-org/llama.cpp/releases

### Installing the prebuilt binaries (manual)

If you do not use the `yzma` installer, you must extract the library files into a directory on your local machine.

For Linux, they have the `.so` file extension. For example, `libllama.so`, `libmtmd.so` and so on.

***Important Note***
You currently need to set both the `LD_LIBRARY_PATH` and the `YZMA_LIB` env variable to the directory with your llama.cpp library files. For example:

```shell
export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/home/ron/Development/yzma/lib
export YZMA_LIB=/home/ron/Development/yzma/lib
```

## macOS

You can use the `yzma` installer program or else download the `llama.cpp` libraries for macOS.

You can obtain them from https://github.com/ggml-org/llama.cpp/releases

If you do not use the `yzma` installer, extract the library files into a directory on your local machine.

For macOS, they have the `.dylib` file extension. For example, `libllama.dylib`, `libmtmd.dylib` and so on. You do not need the other downloaded files to use the `llama.cpp` libraries with `yzma`.

***Important Note***
You currently need to set both the `LD_LIBRARY_PATH` and the `YZMA_LIB` env variable to the directory with your `llama.cpp` library files. For example:

```shell
export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/home/ron/Development/yzma/lib
export YZMA_LIB=/home/ron/Development/yzma/lib
```

## Windows

### CPU

You can use the `yzma` installer program or obtain prebuilt `llama.cpp` libraries from here:

https://github.com/ggml-org/llama.cpp/releases

### CUDA

You will need to install the CUDA drivers for your Nvidia GPU.

See https://docs.nvidia.com/cuda/cuda-installation-guide-microsoft-windows/

Once your drivers are installed, you can use the `yzma` installer program or obtain prebuilt `llama.cpp` binaries for Ubuntu 24.04 using the NVidia 12.9 drivers for CUDA acceleration from here:

https://github.com/ggml-org/llama.cpp/releases

You will also need to download the `cudart` files from the same location as the other `llama.cpp` libraries when using CUDA on Windows.

### Vulkan

You can use the `yzma` installer program or obtain prebuilt `llama.cpp` libraries from here:

https://github.com/ggml-org/llama.cpp/releases

### Installing the prebuilt binaries (manual install)

Extract the library files into a directory on your local machine.

For Windows, they have the `.dll` file extension. For example, `llama.dll`, `mtmd.dll` and so on.

***Important Note***
You currently need to set both the `PATH` and the `YZMA_LIB` env variable to the directory with your `llama.cpp` library files. For example:

```shell
set PATH=%PATH%;C:\yzma\lib
set YZMA_LIB=C:\yzma\lib
```
