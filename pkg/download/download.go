package download

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	getter "github.com/hashicorp/go-getter"
)

var (
	ErrUnknownArch      = errors.New("unknown architecture")
	ErrUnknownOS        = errors.New("unknown OS")
	ErrUnknownProcessor = errors.New("unknown processor")
	ErrInvalidVersion   = errors.New("invalid version")
)

var (
	// RetryCount is how many times the package will retry to obtain the latest llama.cpp version.
	RetryCount = 10
	// RetryDelay is the delay between retries when obtaining the latest llama.cpp version.
	RetryDelay = 3 * time.Second
)

// LlamaLatestVersion fetches the latest release tag of llama.cpp from the GitHub API.
func LlamaLatestVersion() (string, error) {
	var version string
	var err error
	for range RetryCount {
		version, err = getLatestVersion()
		if err == nil {
			return version, nil
		}
		time.Sleep(RetryDelay)
	}

	return "", errors.New("unable to fetch latest version")
}

func getLatestVersion() (string, error) {
	const apiURL = "https://api.github.com/repos/ggml-org/llama.cpp/releases/latest"

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code %d from GitHub API", resp.StatusCode)
	}

	var result struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.TagName, nil
}

// getDownloadLocationAndFilename returns the download location and filename for the given parameters.
func getDownloadLocationAndFilename(arch Arch, os OS, prcssr Processor, version string, dest string) (location, filename string, err error) {
	location = fmt.Sprintf("https://github.com/ggml-org/llama.cpp/releases/download/%s", version)

	switch os {
	case Linux:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				return "", "", errors.New("precompiled binaries for Linux ARM64 CPU are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-x64.zip//build/bin", version)
		case CUDA:
			location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
			if arch == ARM64 {
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-arm64.zip", version)
			} else {
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-x64.zip", version)
			}
		case Vulkan:
			if arch == ARM64 {
				location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-arm64.zip", version)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-x64.zip//build/bin", version)
		default:
			return "", "", ErrUnknownProcessor
		}

	case Darwin:
		switch prcssr {
		case Metal:
			if arch != ARM64 {
				return "", "", errors.New("precompiled binaries for macOS non-ARM64 CPU/Metal are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-macos-arm64.zip//build/bin", version)
		case CPU:
			if arch == ARM64 {
				filename = fmt.Sprintf("llama-%s-bin-macos-arm64-cpu.zip//build/bin", version)
			} else {
				filename = fmt.Sprintf("llama-%s-bin-macos-x64-cpu.zip//build/bin", version)
			}
		default:
			return "", "", ErrUnknownProcessor
		}

	case Windows:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				filename = fmt.Sprintf("llama-%s-bin-win-cpu-arm64.zip", version)
			} else {
				filename = fmt.Sprintf("llama-%s-bin-win-cpu-x64.zip", version)
			}
		case CUDA:
			if arch == ARM64 {
				return "", "", errors.New("precompiled binaries for Windows ARM64 CUDA are not available")
			}
			// also requires the CUDA RT files
			cudart := "cudart-llama-bin-win-cuda-12.4-x64.zip"
			url := fmt.Sprintf("%s/%s", location, cudart)
			if err := get(url, dest); err != nil {
				return "", "", err
			}
			filename = fmt.Sprintf("llama-%s-bin-win-cuda-12.4-x64.zip", version)
		case Vulkan:
			if arch == ARM64 {
				return "", "", errors.New("precompiled binaries for Windows ARM64 Vulkan are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-win-vulkan-x64.zip", version)
		default:
			return "", "", ErrUnknownProcessor
		}

	default:
		return "", "", ErrUnknownOS
	}

	return location, filename, nil
}

// Get downloads the llama.cpp precompiled binaries for the desired arch/OS/processor.
// arch can be one of the following values: "amd64", "arm64".
// os can be one of the following values: "linux", "darwin", "windows".
// processor can be one of the following values: "cpu", "cuda", "vulkan", "metal".
// version should be the desired `b1234` formatted llama.cpp version. You can use the
// [LlamaLatestVersion] function to obtain the latest release.
// dest in the destination directory for the downloaded binaries.
func Get(architecture string, operatingSystem string, processor string, version string, dest string) error {
	arch, err := ParseArch(architecture)
	if err != nil {
		return ErrUnknownArch
	}

	os, err := ParseOS(operatingSystem)
	if err != nil {
		return ErrUnknownOS
	}

	prcssr, err := ParseProcessor(processor)
	if err != nil {
		return ErrUnknownProcessor
	}

	if err := VersionIsValid(version); err != nil {
		return ErrInvalidVersion
	}

	location, filename, err := getDownloadLocationAndFilename(arch, os, prcssr, version, dest)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", location, filename)
	return get(url, dest)
}

func get(url, dest string) error {
	client := &getter.Client{
		Ctx:  context.Background(),
		Src:  url,
		Dst:  dest,
		Mode: getter.ClientModeAny,
	}

	if err := client.Get(); err != nil {
		return err
	}

	return nil
}

// VersionIsValid checks if the provided version string is valid.
func VersionIsValid(version string) error {
	if !strings.HasPrefix(version, "b") {
		return ErrInvalidVersion
	}

	return nil
}

// LibraryName returns the name for the llama.cpp library for any given OS.
func LibraryName(operatingSystem string) string {
	os, err := ParseOS(operatingSystem)
	if err != nil {
		return "unknown"
	}

	switch os {
	case Linux:
		return "libllama.so"
	case Windows:
		return "llama.dll"
	case Darwin:
		return "libllama.dylib"
	default:
		return "unknown"
	}
}
