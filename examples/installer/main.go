package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hybridgroup/yzma/pkg/download"
)

func main() {
	if err := handleFlags(); err != nil {
		showUsage()
		os.Exit(0)
	}

	if !*upgrade {
		if download.AlreadyInstalled(*libPath) {
			fmt.Println("llama.cpp already installed at", *libPath)
			return
		}
	}

	if *processor == "" {
		*processor = "cpu"
		if cudaInstalled, cudaVersion := download.HasCUDA(); cudaInstalled {
			fmt.Printf("CUDA detected (version %s), using CUDA build\n", cudaVersion)
			*processor = "cuda"
		}
	}

	if *version == "" || *version == "latest" {
		fmt.Println("installing latest llama.cpp version to", *libPath)
	} else {
		fmt.Println("installing llama.cpp version", *version, "to", *libPath)
	}
	if err := download.Get(runtime.GOARCH, runtime.GOOS, *processor, *version, *libPath); err != nil {
		fmt.Println("failed to download llama.cpp:", err.Error())
		return
	}

	fmt.Println("done.")
}
