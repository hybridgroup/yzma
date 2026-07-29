// Installs the llama.cpp libraries using a custom Resolver.
//
// A Resolver decides which release assets a Target needs, so an application can
// install builds the built-in table does not name: an internal mirror, a local
// file:// path on an air-gapped machine, its own llama.cpp build, or a different
// CUDA major version. Targets it does not care about fall through to
// download.DefaultResolver.
//
// Run with -mirror to point the CUDA build at somewhere else:
//
//	go run ./examples/resolver -lib ./lib -mirror https://mirror.example.com/llama
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/hybridgroup/yzma/pkg/download"
)

func main() {
	libPath := flag.String("lib", os.Getenv("YZMA_LIB"), "directory to install the libraries into")
	version := flag.String("version", "latest", "llama.cpp release tag, or latest")
	mirror := flag.String("mirror", "", "base URL to fetch the CUDA build from instead")
	flag.Parse()

	if *libPath == "" {
		fmt.Println("missing -lib flag or YZMA_LIB env var")
		os.Exit(1)
	}

	arch, err := download.ParseArch(runtime.GOARCH)
	if err != nil {
		fmt.Println("unsupported architecture:", runtime.GOARCH)
		os.Exit(1)
	}
	operatingSystem, err := download.ParseOS(runtime.GOOS)
	if err != nil {
		fmt.Println("unsupported OS:", runtime.GOOS)
		os.Exit(1)
	}

	processor := download.CPU
	if installed, cudaVersion := download.HasCUDA(); installed {
		fmt.Printf("CUDA detected (version %s), using CUDA build\n", cudaVersion)
		processor = download.CUDA
	}

	// Redirect only the target we have a mirror for; everything else keeps whatever
	// the built-in table says.
	resolver := download.ResolverFunc(func(target download.Target) ([]string, error) {
		if *mirror == "" || target.Processor != download.CUDA {
			return download.DefaultResolver.Resolve(target)
		}
		url := fmt.Sprintf("%s/llama-%s-bin-%s-cuda-%s.tar.gz", *mirror, target.Version, target.OS, target.Arch)
		fmt.Println("resolving CUDA build to", url)
		return []string{url}, nil
	})

	target := download.Target{Arch: arch, OS: operatingSystem, Processor: processor, Version: *version}
	fmt.Println("installing llama.cpp", *version, "to", *libPath)
	if err := download.Install(context.Background(), target, *libPath, download.ProgressTracker, resolver); err != nil {
		fmt.Println("failed to download llama.cpp:", err)
		os.Exit(1)
	}
	fmt.Println("done.")
}
