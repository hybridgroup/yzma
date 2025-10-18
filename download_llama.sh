#!/usr/bin/env bash

# Check if a parameter is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <windows|macos|linux> <cpu|cuda|vulkan|metal>"
    exit 1
fi

# Get the platform parameter
platform=$1

# Validate the platform parameter
if [[ "$platform" != "windows" && "$platform" != "macos" && "$platform" != "linux" ]]; then
    echo "Invalid platform: $platform"
    echo "Valid options are: windows, macos, linux"
    exit 1
fi

# Get the processor parameter
processor=$2

# Validate the processor parameter
if [[ "$processor" != "cpu" && "$processor" != "cuda" && "$processor" != "vulkan" && "$processor" != "metal" ]]; then
    echo "Invalid platform: $processor"
    echo "Valid options are: cpu, cuda, vulkan, metal"
    exit 1
fi

latest_llama=$(curl -s https://api.github.com/repos/ggml-org/llama.cpp/releases/latest | jq -r '.tag_name')
echo "latest llama: $latest_llama"

download_files() {
    local platform=$1
    local processor=$2
    local base_url="https://github.com/ggml-org/llama.cpp/releases/download/$latest_llama"

    # Determine the file to download based on platform and processor
    case "$platform" in
        windows)
            case "$processor" in
                cpu) file="llama-$latest_llama-bin-win-cpu-x64.zip" ;;
                cuda) file="llama-$latest_llama-bin-win-cuda-12.4-x64.zip" ;;
                vulkan) file="llama-$latest_llama-bin-win-vulkan-x64.zip" ;;
                metal) echo "Metal is not supported on Windows"; exit 1 ;;
            esac
            ;;
        macos)
            case "$processor" in
                cpu) file="llama-$latest_llama-bin-macos-arm64.zip" ;;
                cuda) echo "CUDA is not supported on macOS"; exit 1 ;;
                vulkan) echo "Vulkan is not supported on macOS"; exit 1 ;;
                metal) file="llama-$latest_llama-bin-macos-arm64.zip" ;;
            esac
            ;;
        linux)
            case "$processor" in
                cpu) file="llama-$latest_llama-bin-ubuntu-x64.zip" ;;
                cuda) base_url="https://github.com/hybridgroup/llama-cpp-builder/releases/download/$latest_llama"; file="llama-$latest_llama-bin-ubuntu-cuda-x64.zip" ;;
                cuda) file="llama-$latest_llama-bin-vulkan-x64.zip" ;;
                metal) echo "Metal is not supported on Linux"; exit 1 ;;
            esac
            ;;
    esac

    # Download the file
    echo "Downloading $file for $platform with $processor..."
    curl -LO "$base_url/$file"

    # Verify the download
    if [ $? -eq 0 ]; then
        echo "Download completed: $file"
    else
        echo "Failed to download $file"
        exit 1
    fi
}

# Call the download function
download_files "$platform" "$processor"
