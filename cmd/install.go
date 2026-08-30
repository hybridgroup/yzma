package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hybridgroup/yzma/pkg/download"
	"github.com/urfave/cli/v2"
)

var InstallCmd = &cli.Command{
	Name:  "install",
	Usage: "Install llama.cpp libraries used by yzma",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version of llama.cpp to install (leave empty for the version this yzma release uses)",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "lib",
			Aliases: []string{"l"},
			Usage:   "path to llama.cpp compiled library files",
			EnvVars: []string{"YZMA_LIB"},
		},
		&cli.StringFlag{
			Name:    "processor",
			Aliases: []string{"p"},
			Usage:   "processor to use (cpu, cuda, metal, rocm, vulkan)",
			Value:   "",
		},
		&cli.StringFlag{
			Name:  "os",
			Usage: "target to use (linux, windows, darwin, bookworm, trixie, wasm)",
			Value: runtime.GOOS,
		},
		&cli.BoolFlag{
			Name:    "upgrade",
			Aliases: []string{"u"},
			Usage:   "upgrade existing installation",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
			Usage:   "suppress output during installation",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {
		return runInstall(c)
	},
}

func runInstall(c *cli.Context) error {
	libPath := c.String("lib")
	version := c.String("version")
	processor := c.String("processor")
	osInstall := c.String("os")
	upgrade := c.Bool("upgrade")

	if libPath == "" {
		return fmt.Errorf("missing lib flag or YZMA_LIB env var")
	}

	// A wasm install puts the WebAssembly build of llama.cpp in place, which has
	// its own file name, so the check for an existing install follows the target
	// and not the machine that runs the command.
	if !upgrade {
		if _, err := os.Stat(filepath.Join(libPath, download.LibraryName(osInstall))); !os.IsNotExist(err) {
			fmt.Println("llama.cpp already installed at", libPath)
			return nil
		}
	}

	quiet := c.Bool("quiet")
	if !quiet {
		switch {
		case version == "" && download.DefaultVersion != "":
			fmt.Println("installing llama.cpp version", download.DefaultVersion, "to", libPath)
		case version == "" || version == "latest":
			fmt.Println("installing latest llama.cpp version to", libPath)
		default:
			fmt.Println("installing llama.cpp version", version, "to", libPath)
		}
	} else {
		download.ProgressTracker = nil
	}

	if osInstall == download.Wasm.String() {
		// A WebAssembly build has no GPU.
		processor = download.CPU.String()
	}

	if processor == "" {
		processor = "cpu"

		if cudaInstalled, cudaVersion := download.HasCUDA(); cudaInstalled {
			if !quiet {
				fmt.Printf("CUDA detected (version %s), using CUDA build\n", cudaVersion)
			}
			processor = "cuda"
		} else if rocmInstalled, rocmVersion := download.HasROCm(); rocmInstalled {
			if !quiet {
				fmt.Printf("ROCm detected (version %s), using ROCm build\n", rocmVersion)
			}
			processor = "rocm"
		}
	}

	if err := download.Get(runtime.GOARCH, osInstall, processor, version, libPath); err != nil {
		return fmt.Errorf("failed to download llama.cpp: %w", err)
	}

	if !quiet {
		fmt.Println("done.")
		if osInstall == download.Wasm.String() {
			showWasmRequirements(libPath)
		} else {
			showInstallRequirements(libPath)
		}
	}

	return nil
}

// showWasmRequirements says what to do with a WebAssembly install. YZMA_LIB has
// no part in it, because a browser loads the files over HTTP.
func showWasmRequirements(libPath string) {
	fmt.Println(`
The WebAssembly build of llama.cpp is in ` + libPath + `

Put the JavaScript glue and the program of your page in the same directory, then
serve it with the headers that a build with more than one thread needs:

    make wasm-example
    make serve-wasm`)
}

func showInstallRequirements(libPath string) {
	if os.Getenv("YZMA_LIB") == libPath {
		return
	}

	switch runtime.GOOS {
	case "linux":
		fmt.Println(`
You may want to set the YZMA_LIB environment variable to the directory with your llama.cpp library files. For example:

    export YZMA_LIB=` + libPath)
	case "windows":
		fmt.Println(`
You may want to set the YZMA_LIB environment variable to the directory with your llama.cpp library files. For example:

    set YZMA_LIB=` + libPath)
	case "darwin":
		fmt.Println(`
You may want to set the YZMA_LIB environment variable to the directory with your llama.cpp library files. For example:

    export YZMA_LIB=` + libPath)
	}
}
