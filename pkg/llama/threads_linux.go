package llama

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const sysfsCPU = "/sys/devices/system/cpu"

// hybridCPUFile lists the performance CPUs of a machine that has two kinds of
// core. Linux makes it for the performance counters, and it holds the same set
// that llama.cpp finds with CPUID.
const hybridCPUFile = "/sys/devices/cpu_core/cpus"

// mathCores counts the cores of this machine that do the arithmetic well. It
// is the number of physical cores, and on a machine with performance cores and
// efficiency cores it leaves the efficiency cores out. The count is 0 when
// sysfs says nothing, and then the caller uses its own default.
func mathCores() int {
	return len(mathCPUs())
}

// mathCPUs gives one CPU of each core that does the arithmetic well. Two CPUs
// of one core do not both appear, because a core gives its best with one
// thread of arithmetic. It gives nothing when sysfs says nothing.
func mathCPUs() []int32 {
	cpus := performanceCPUs()
	if cpus == nil {
		cpus = allCPUs()
	}

	// One CPU for each group that shares a core, thus one for each physical
	// core.
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

// performanceCPUs gives the CPUs of the performance cores, or nil when this
// machine has one kind of core only.
func performanceCPUs() []int {
	line, err := os.ReadFile(hybridCPUFile)
	if err != nil {
		return nil
	}
	return parseCPUList(string(line))
}

// allCPUs gives every CPU that sysfs knows.
func allCPUs() []int {
	line, err := os.ReadFile(filepath.Join(sysfsCPU, "present"))
	if err != nil {
		return nil
	}
	return parseCPUList(string(line))
}

// parseCPUList reads a list in the form that sysfs uses, such as "0-15" or
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
