package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hybridgroup/yzma/pkg/download"
	"github.com/urfave/cli/v2"
)

var installCmd = &cli.Command{
	Name:  "install",
	Usage: "Install llama.cpp libraries used by yzma",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "version of llama.cpp to install (leave empty for latest)",
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
			Usage:   "processor to use (cpu, cuda, metal, vulkan)",
			Value:   "cpu",
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
	upgrade := c.Bool("upgrade")

	if libPath == "" {
		return fmt.Errorf("missing lib flag or YZMA_LIB env var")
	}

	if !upgrade {
		if _, err := os.Stat(filepath.Join(libPath, download.LibraryName(runtime.GOOS))); !os.IsNotExist(err) {
			fmt.Println("llama.cpp already installed at", libPath)
			return nil
		}
	}

	if version == "" {
		var err error
		version, err = download.LlamaLatestVersion()
		if err != nil {
			return fmt.Errorf("could not obtain latest version: %w", err)
		}
	}

	quiet := c.Bool("quiet")
	if !quiet {
		fmt.Println("installing llama.cpp version", version, "to", libPath)
	}

	if err := download.Get(runtime.GOARCH, runtime.GOOS, processor, version, libPath); err != nil {
		return fmt.Errorf("failed to download llama.cpp: %w", err)
	}

	if !quiet {
		fmt.Println("done.")
		showInstallRequirements()
	}

	return nil
}

func showInstallRequirements() {
	switch runtime.GOOS {
	case "linux":
		fmt.Println(`
Make sure you have set the LD_LIBRARY_PATH to the directory with your llama.cpp library files. For example:

export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/home/your/location/yzma/lib

You may also want to set the YZMA_LIB environment variable to this path as well:

export YZMA_LIB=/home/your/location/yzma/lib
`)
	case "windows":
		fmt.Println(`
Make sure you have set the PATH to the directory with your llama.cpp library files. For example:

set PATH=%PATH%;C:\your\location\yzma\lib

You may also want to set the YZMA_LIB environment variable to this path as well:

set YZMA_LIB=C:\your\location\yzma\lib
`)
	case "darwin":
		fmt.Println(`
Make sure you have set the LD_LIBRARY_PATH to the directory with your llama.cpp library files. For example:

export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/home/your/location/yzma/lib

You may also want to set the YZMA_LIB environment variable to this path as well:

export YZMA_LIB=/home/your/location/yzma/lib
`)
	}
}
