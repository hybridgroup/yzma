package llama

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const sysfsCPU = "/sys/devices/system/cpu"

// hybridCPUFile lists the performance CPUs on a hybrid machine. Linux creates
// it for the performance counters, and it holds the same set that llama.cpp
// finds with CPUID.
const hybridCPUFile = "/sys/devices/cpu_core/cpus"

// mathCores counts the compute cores on this machine. It is the number of
// physical cores, leaving out efficiency cores on a hybrid machine. It returns
// 0 when sysfs reports nothing, and the caller then uses its own default.
func mathCores() int {
	return len(mathCPUs())
}

// mathCPUs returns one CPU per compute core. SMT siblings are left out,
// because a core runs arithmetic best with one thread. It returns nil when
// sysfs reports nothing.
func mathCPUs() []int32 {
	cpus := performanceCPUs()
	if cpus == nil {
		cpus = allCPUs()
	}

	// One CPU per sibling group, so one per physical core.
	seen := map[string]struct{}{}
	var out []int32
	for _, cpu := range cpus {
		line, err := os.ReadFile(filepath.Join(sysfsCPU, "cpu"+strconv.Itoa(cpu), "topology", "thread_siblings_list"))
		if err != nil {
			continue
		}
		key := strings.TrimSpace(string(line))
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, int32(cpu))
	}
	return out
}

// performanceCPUs returns the CPUs of the performance cores, or nil when this
// machine has only one kind of core.
func performanceCPUs() []int {
	line, err := os.ReadFile(hybridCPUFile)
	if err != nil {
		return nil
	}
	return parseCPUList(string(line))
}

// allCPUs returns every CPU that sysfs lists.
func allCPUs() []int {
	line, err := os.ReadFile(filepath.Join(sysfsCPU, "present"))
	if err != nil {
		return nil
	}
	return parseCPUList(string(line))
}

// parseCPUList parses a CPU list in sysfs format, such as "0-15" or
// "0,2,4-7".
func parseCPUList(s string) []int {
	var cpus []int
	for _, part := range strings.Split(strings.TrimSpace(s), ",") {
		lo, hi, found := strings.Cut(part, "-")
		first, err := strconv.Atoi(lo)
		if err != nil {
			continue
		}
		last := first
		if found {
			if last, err = strconv.Atoi(hi); err != nil {
				continue
			}
		}
		for cpu := first; cpu <= last; cpu++ {
			cpus = append(cpus, cpu)
		}
	}
	return cpus
}
