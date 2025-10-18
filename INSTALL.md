# Installation

Here is information on how to install the requirements for running an application using `llama.cpp`.

## Linux

You will need to download the `llama.cpp` libraries for Linux. Choose the download option that matches your desired configuration.

### CPU

You can obtain them from https://github.com/ggml-org/llama.cpp/releases

### CUDA

Prebuilt binaries for Ubuntu 24.04 using the NVidia 12.9 drivers for CUDA acceleration are available here:

https://github.com/hybridgroup/llama-cpp-builder/releases

### Vulkan

You can obtain them from https://github.com/ggml-org/llama.cpp/releases

### Installing the prebuilt binaries

Extract the library files into a directory on your local machine.

For Linux, they have the `.so` file extension. For example, `libllama.so`, `libmtmd.so` and so on.

***Important Note***
You currently need to set both the `LD_LIBRARY_PATH` and the `YZMA_LIB` env variable to the directory with your llama.cpp library files. For example:

```shell
export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/home/ron/Development/yzma/lib
export YZMA_LIB=/home/ron/Development/yzma/lib
```

## macOS

You will need to download the `llama.cpp` libraries for macOS. You can obtain them from https://github.com/ggml-org/llama.cpp/releases

Extract the library files into a directory on your local machine.

For macOS, they have the `.dylib` file extension. For example, `libllama.dylib`, `libmtmd.dylib` and so on. You do not need the other downloaded files to use the `llama.cpp` libraries with `yzma`.

***Important Note***
You currently need to set both the `LD_LIBRARY_PATH` and the `YZMA_LIB` env variable to the directory with your `llama.cpp` library files. For example:

```shell
export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/home/ron/Development/yzma/lib
export YZMA_LIB=/home/ron/Development/yzma/lib
```

## Windows

### CPU

You can obtain them from https://github.com/ggml-org/llama.cpp/releases

### CUDA

You can obtain them from https://github.com/ggml-org/llama.cpp/releases

### Vulkan

You can obtain them from https://github.com/ggml-org/llama.cpp/releases

### Installing the prebuilt binaries

Extract the library files into a directory on your local machine.

For Windows, they have the `.dll` file extension. For example, `llama.dll`, `mtmd.dll` and so on.

***Important Note***
You currently need to set both the `PATH` and the `YZMA_LIB` env variable to the directory with your `llama.cpp` library files. For example:

```shell
set PATH=%PATH%;C:\yzma\lib
set YZMA_LIB=C:\yzma\lib
```
