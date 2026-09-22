package main

import (
	"fmt"
	"os"
	"runtime"

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

	showThreads()

	sysInfo := llama.PrintSystemInfo()
	fmt.Println("-- llama.cpp System Information --")
	fmt.Println(sysInfo)
}

// showThreads prints what yzma makes of the cores of this machine. Call it
// after llama.Init, because the register of the CPU backend exists only then.
func showThreads() {
	fmt.Println("-- CPU Threads --")
	fmt.Printf("Logical CPUs:      %d\n", runtime.NumCPU())
	fmt.Printf("Inference threads: %d\n", llama.Threads())

	cpus := llama.PerformanceCPUs()
	if len(cpus) == 0 {
		fmt.Println("Performance CPUs:  the system does not say")
	} else {
		fmt.Printf("Performance CPUs:  %v\n", cpus)
	}

	tp, err := llama.NewPerformanceThreadpool()
	if err != nil {
		fmt.Printf("Thread pool:       none, %v\n", err)
	} else {
		fmt.Printf("Thread pool:       %d threads, each held to a CPU\n", len(cpus))
		llama.ThreadpoolFree(tp)
	}

	fmt.Println()
}
