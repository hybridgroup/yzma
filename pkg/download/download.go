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
	ErrUnknownOS        = errors.New("unknown OS")
	ErrUnknownProcessor = errors.New("unknown processor")
	ErrInvalidVersion   = errors.New("invalid version")
)

// RetryCount is how many times the package will retry to obtain the latest llama.cpp version.
var RetryCount = 3

// LlamaLatestVersion fetches the latest release tag of llama.cpp from the GitHub API.
func LlamaLatestVersion() (string, error) {
	var version string
	var err error
	for range RetryCount {
		version, err = getLatestVersion()
		if err == nil {
			return version, nil
		}
		time.Sleep(3 * time.Second)
	}

	return "", err
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

// Get downloads the llama.cpp precompiled binaries for the desired OS/processor.
// os can be one of the following values: "linux", "darwin", "windows".
// processor can be one of the following values: "cpu", "cuda", "vulkan", "metal".
// version should be the desired `b1234` formatted llama.cpp version. You can use the
// [LlamaLatestVersion] function to obtain the latest release.
// dest in the destination directory for the downloaded binaries.
func Get(operatingSystem string, processor string, version string, dest string) error {
	os, err := ParseOS(operatingSystem)
	if err != nil {
		return ErrUnknownOS
	}

	prcssr, err := ParseProcessor(processor)
	if err != nil {
		return ErrUnknownProcessor
	}

	if err := VersionIsValid(version); err != nil {
		return err
	}

	var location, filename string
	location = fmt.Sprintf("https://github.com/ggml-org/llama.cpp/releases/download/%s", version)

	switch os {
	case Linux:
		switch prcssr {
		case CPU:
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-x64.zip//build/bin", version)
		case CUDA:
			location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-x64.zip", version)
		case Vulkan:
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-x64.zip//build/bin", version)
		default:
			return ErrUnknownProcessor
		}

	case Darwin:
		switch prcssr {
		case CPU, Metal:
			filename = fmt.Sprintf("llama-%s-bin-macos-arm64.zip//build/bin", version)
		default:
			return ErrUnknownProcessor
		}

	case Windows:
		switch prcssr {
		case CPU:
			filename = fmt.Sprintf("llama-%s-bin-win-cpu-x64.zip", version)
		case CUDA:
			// also requires the CUDA RT files
			cudart := "cudart-llama-bin-win-cuda-12.4-x64.zip"
			url := fmt.Sprintf("%s/%s", location, cudart)
			if err := get(url, dest); err != nil {
				return err
			}

			filename = fmt.Sprintf("llama-%s-bin-win-cuda-12.4-x64.zip", version)
		case Vulkan:
			filename = fmt.Sprintf("llama-%s-bin-win-vulkan-x64.zip", version)
		default:
			return ErrUnknownProcessor
		}

	default:
		return ErrUnknownOS
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
