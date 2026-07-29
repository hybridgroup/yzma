package download

import (
	"context"
	"errors"
	"fmt"

	getter "github.com/hashicorp/go-getter"
)

// Target identifies the machine a llama.cpp build is for.
type Target struct {
	Arch      Arch
	OS        OS
	Processor Processor

	// Version is the llama.cpp release tag, e.g. "b7974". "" or "latest" resolves to
	// the newest release.
	Version string
}

// Resolver reports the release assets to install for a Target, as URLs downloaded in
// the order returned. Implement it to reach builds the built-in table does not name.
// Implementations must not download.
type Resolver interface {
	Resolve(target Target) (urls []string, err error)
}

// ResolverFunc adapts an ordinary function to [Resolver].
type ResolverFunc func(target Target) ([]string, error)

// Resolve calls f.
func (f ResolverFunc) Resolve(target Target) ([]string, error) { return f(target) }

// DefaultResolver resolves the assets published on the llama.cpp and llama-cpp-builder
// release pages. [Install] uses it when no resolver is given.
var DefaultResolver Resolver = ResolverFunc(defaultResolve)

// defaultResolve is the built-in platform table.
func defaultResolve(target Target) ([]string, error) {
	arch, os, prcssr, version := target.Arch, target.OS, target.Processor, target.Version

	var extra []string
	var location, filename string
	location = fmt.Sprintf("https://github.com/ggml-org/llama.cpp/releases/download/%s", version)

	location = fmt.Sprintf("https://github.com/ggml-org/llama.cpp/releases/download/%s", version)

	switch os {
	case Linux:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cpu-arm64.tar.gz", version)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-x64.tar.gz", version)
		case CUDA:
			location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
			if arch == ARM64 {
				// defaults to CUDA 12 assuming that is running Jetson Orin.
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-arm64.tar.gz", version)
			} else {
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-13-x64.tar.gz", version)
			}
		case Vulkan:
			if arch == ARM64 {
				location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-arm64.tar.gz", version)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-x64.tar.gz", version)
		case ROCm:
			if arch != AMD64 {
				return nil, errors.New("precompiled binaries for Linux ARM64 ROCm are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-rocm-7.2-x64.tar.gz", version)
		default:
			return nil, ErrUnknownProcessor
		}

	case Bookworm:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cpu-arm64.tar.gz", version)
				break
			}

			// no AMD64 for bookworm
			return nil, ErrUnknownProcessor
		case CUDA:
			location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
			if arch == ARM64 {
				// Jetson Orin.
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-arm64.tar.gz", version)
				break
			}

			// no AMD64 for bookworm
			return nil, ErrUnknownProcessor
		case Vulkan:
			if arch == ARM64 {
				location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-arm64.tar.gz", version)
				break
			}

			// no AMD64 for bookworm
			return nil, ErrUnknownProcessor
		default:
			return nil, ErrUnknownProcessor
		}

	case Trixie:
		switch prcssr {
		case CPU:
			if arch == ARM64 {
				location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-trixie-cpu-arm64.tar.gz", version)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-x64.tar.gz", version)
		case CUDA:
			location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
			if arch == ARM64 {
				// not yet
				return nil, ErrUnknownProcessor
			} else {
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-cuda-13-x64.tar.gz", version)
			}
		case Vulkan:
			if arch == ARM64 {
				location = fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s", version)
				filename = fmt.Sprintf("llama-%s-bin-ubuntu-trixie-vulkan-arm64.tar.gz", version)
				break
			}
			filename = fmt.Sprintf("llama-%s-bin-ubuntu-vulkan-x64.tar.gz", version)
		default:
			return nil, ErrUnknownProcessor
		}

	case Darwin:
		switch prcssr {
		case Metal:
			if arch != ARM64 {
				return nil, errors.New("precompiled binaries for macOS non-ARM64 CPU/Metal are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-macos-arm64.tar.gz", version)
		case CPU:
			if arch == ARM64 {
				filename = fmt.Sprintf("llama-%s-bin-macos-arm64.tar.gz", version)
			} else {
				filename = fmt.Sprintf("llama-%s-bin-macos-x64.tar.gz", version)
			}
		default:
			return nil, ErrUnknownProcessor
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
				return nil, errors.New("precompiled binaries for Windows ARM64 CUDA are not available")
			}
			// also requires the CUDA RT files
			cudart := "cudart-llama-bin-win-cuda-13.3-x64.zip"
			extra = append(extra, fmt.Sprintf("%s/%s", location, cudart))
			filename = fmt.Sprintf("llama-%s-bin-win-cuda-13.3-x64.zip", version)
		case Vulkan:
			if arch == ARM64 {
				return nil, errors.New("precompiled binaries for Windows ARM64 Vulkan are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-win-vulkan-x64.zip", version)
		case ROCm:
			if arch != AMD64 {
				return nil, errors.New("precompiled binaries for Windows ARM64 ROCm are not available")
			}
			filename = fmt.Sprintf("llama-%s-bin-win-hip-radeon-x64.zip", version)
		default:
			return nil, ErrUnknownProcessor
		}

	default:
		return nil, ErrUnknownOS
	}

	return append(extra, fmt.Sprintf("%s/%s", location, filename)), nil
}

// Install downloads the llama.cpp binaries for target into dest. A nil resolver means
// [DefaultResolver].
func Install(ctx context.Context, target Target, dest string, progress getter.ProgressTracker, resolver Resolver) error {
	if resolver == nil {
		resolver = DefaultResolver
	}

	autoVersion := target.Version == "" || target.Version == "latest"
	if autoVersion {
		latest, err := LlamaLatestVersion()
		if err != nil {
			return err
		}
		target.Version = latest
	}
	if err := VersionIsValid(target.Version); err != nil {
		return ErrInvalidVersion
	}

	err := installAssets(ctx, target, dest, progress, resolver)

	// The newest release may still be building for this platform.
	if err != nil && autoVersion && errors.Is(err, ErrFileNotFound) {
		previous, prevErr := LlamaPreviousVersion()
		if prevErr != nil {
			return err
		}
		target.Version = previous
		return installAssets(ctx, target, dest, progress, resolver)
	}

	return err
}

func installAssets(ctx context.Context, target Target, dest string, progress getter.ProgressTracker, resolver Resolver) error {
	urls, err := resolver.Resolve(target)
	if err != nil {
		return err
	}
	for _, url := range urls {
		if err := getFunc(ctx, url, dest, progress); err != nil {
			return err
		}
	}
	return nil
}
