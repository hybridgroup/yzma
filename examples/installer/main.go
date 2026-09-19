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

	// The CUDA version selects the CUDA build, so it is read for a CUDA install that
	// was asked for as well as for one that is found here.
	var cudaVersion string
	if *processor == "" || *processor == download.CUDA.String() {
		cudaInstalled, detected := download.HasCUDA()
		switch {
		case cudaInstalled:
			cudaVersion = detected
			if *processor == "" {
				fmt.Printf("CUDA detected (version %s), using CUDA build\n", cudaVersion)
				*processor = download.CUDA.String()
			}
		case *processor == "":
			*processor = download.CPU.String()
		}
	}

	switch {
	case *version == "" && download.DefaultVersion != "":
		fmt.Println("installing llama.cpp version", download.DefaultTag(), "to", *libPath)
	case *version == "" || *version == "latest":
		fmt.Println("installing latest llama.cpp version to", *libPath)
	default:
		fmt.Println("installing llama.cpp version", *version, "to", *libPath)
	}
	if err := download.Get(runtime.GOARCH, runtime.GOOS, *processor, *version, *libPath,
		download.WithCUDAVersion(cudaVersion)); err != nil {
		fmt.Println("failed to download llama.cpp:", err.Error())
		return
	}

	fmt.Println("done.")
}
