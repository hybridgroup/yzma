package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	help      *string
	libPath   *string
	version   *string
	processor *string
	upgrade   *bool
)

func showUsage() {
	fmt.Println(`
Usage:
installer -version [version] -lib [llama.cpp .so file path] -processor [cpu, cuda, metal, vulkan]`)
}

func handleFlags() error {
	help = flag.String("help", "", "show help")
	version = flag.String("version", "", "version of llama.cpp to install (leave empty for latest)")
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files (leave empty to use YZMA_LIB env var)")
	processor = flag.String("processor", "cpu", "processor to use (cpu, cuda, metal, vulkan)")
	upgrade = flag.Bool("upgrade", false, "upgrade existing installation")

	flag.Parse()

	if *help != "" {
		return errors.New("showing help")
	}

	if os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	return nil
}
