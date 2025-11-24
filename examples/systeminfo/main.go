package main

import (
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func main() {
	if err := handleFlags(); err != nil {
		showUsage()
		os.Exit(0)
	}

	llama.Load(*libPath)
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
}
