package download

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	getter "github.com/hashicorp/go-getter/v2"
)

var (
	errUnknownOS        = errors.New("unknown OS")
	errUnknownProcessor = errors.New("unknown processor")
	errInvalidVersion   = errors.New("invalid version")
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
func Get(os string, processor string, version string, dest string) error {
	if err := VersionIsValid(version); err != nil {
		return err
	}

	var location, filename string
	location = fmt.Sprintf("https://github.com/ggml-org/llama.cpp/releases/download/%s", version)

	switch os {
	case "linux":
		switch processor {
		case "cpu":
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-x64.zip", version)
		case "cuda":
			location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-x64.zip", version)
		case "vulkan":
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-x64.zip", version)
		default:
			return errUnknownProcessor
		}
	case "darwin":
		switch processor {
		case "cpu", "metal":
			filename = fmt.Sprintf("llama-%s-bin-macos-arm64.zip", version)
		default:
			return errUnknownProcessor
		}

	case "windows":
		switch processor {
		case "cpu":
			filename = fmt.Sprintf("llama-%s-bin-win-cpu-x64.zip", version)
		case "cuda":
			filename = fmt.Sprintf("llama-%s-bin-win-cuda-12.4-x64.zip", version)
		case "vulkan":
			filename = fmt.Sprintf("llama-%s-bin-win-vulkan-x64.zip", version)
		default:
			return errUnknownProcessor
		}

	default:
		return errUnknownOS
	}

	url := fmt.Sprintf("%s/%s", location, filename)
	return get(url, dest)
}

func get(url, dest string) error {
	req := &getter.Request{
		Src:     url,
		Dst:     dest,
		GetMode: getter.ModeAny,
	}

	client := &getter.Client{}
	if _, err := client.Get(context.Background(), req); err != nil {
		return err
	}

	return nil
}

// VersionIsValid checks if the provided version string is valid.
func VersionIsValid(version string) error {
	if !strings.HasPrefix(version, "b") {
		return errInvalidVersion
	}

	return nil
}

// LibraryName returns the name for the llama.cpp library for any given OS.
func LibraryName(os string) string {
	switch os {
	case "linux", "freebsd":
		return "libllama.so"
	case "windows":
		return "llama.dll"
	case "darwin":
		return "libllama.dylib"
	default:
		return "unknown"
	}
}
