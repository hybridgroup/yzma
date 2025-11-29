package main

import (
	"fmt"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/urfave/cli/v2"
)

var systemCmd = &cli.Command{
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
		deviceName := llama.GGMLBackendDeviceName(device)

		fmt.Printf("Device %d: %s\n", i, deviceName)
	}

	fmt.Println()

	sysInfo := llama.PrintSystemInfo()
	fmt.Println("-- llama.cpp System Information --")
	fmt.Println(sysInfo)

	return nil
}
