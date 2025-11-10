package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hybridgroup/yzma/pkg/download"
)

func main() {
	if err := handleFlags(); err != nil {
		showUsage()
		os.Exit(0)
	}

	if !*upgrade {
		if _, err := os.Stat(filepath.Join(*libPath, download.LibraryName(runtime.GOOS))); !os.IsNotExist(err) {
			fmt.Println("llama.cpp already installed at", libPath)
			return
		}
	}

	if *version == "" {
		var err error
		*version, err = download.LlamaLatestVersion()
		if err != nil {
			fmt.Println("could not obtain latest version:", err.Error())
			return
		}
	}

	fmt.Println("installing llama.cpp version", *version, "to", *libPath)
	if err := download.Get(runtime.GOOS, *processor, *version, *libPath); err != nil {
		fmt.Println("failed to download llama.cpp:", err.Error())
		return
	}

	fmt.Println("done.")
}
