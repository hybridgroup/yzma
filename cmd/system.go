package cmd

import (
	"fmt"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/urfave/cli/v2"
)

var SystemCmd = &cli.Command{
	Name:  "system",
	Usage: "Show llama.cpp system information",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "lib",
			Aliases: []string{"l"},
			Usage:   "path to llama.cpp compiled library files",
			EnvVars: []string{"YZMA_LIB"},
		},
	},
	Action: func(c *cli.Context) error {
		return runSystemInfo(c)
	},
}

func runSystemInfo(c *cli.Context) error {
	return showSystemInfo(c)
}

func showSystemInfo(c *cli.Context) error {
	libPath := c.String("lib")
	if libPath == "" {
		return fmt.Errorf("missing lib flag or YZMA_LIB env var")
	}

	llama.Load(libPath)
	llama.LogSet(llama.LogSilent())

	llama.Init()
	defer llama.Close()

	fmt.Println("-- Devices --")

	for i := uint64(0); i < llama.GGMLBackendDeviceCount(); i++ {
		device := llama.GGMLBackendDeviceGet(i)
		if device == 0 {
			continue
		}

		fmt.Printf("Device %d: %s\n", i, llama.GGMLBackendDeviceName(device))
		fmt.Printf("  Type:        %s\n", llama.GGMLBackendDevType(device))
		fmt.Printf("  Backend:     %s\n", llama.GGMLBackendRegName(llama.GGMLBackendDeviceBackendReg(device)))

		if desc := llama.GGMLBackendDeviceDescription(device); desc != "" {
			fmt.Printf("  Description: %s\n", desc)
		}

		if free, total := llama.GGMLBackendDeviceMemory(device); total > 0 {
			fmt.Printf("  Memory:      free %d MiB / total %d MiB\n", free/(1024*1024), total/(1024*1024))
		}
	}

	fmt.Println()

	sysInfo := llama.PrintSystemInfo()
	fmt.Println("-- llama.cpp System Information --")
	fmt.Println(sysInfo)

	return nil
}
