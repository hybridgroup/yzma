package llama

import (
	"encoding/binary"
	"unsafe"

	"golang.org/x/sys/windows"
)

var procGetLogicalProcessorInformationEx = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetLogicalProcessorInformationEx")

// relationProcessorCore requests one record per physical core.
const relationProcessorCore = 0

// mathCores counts the compute cores on this machine. It is the number of
// physical cores, leaving out efficiency cores on a hybrid machine. It returns
// 0 when the system reports nothing, and the caller then uses its own default.
func mathCores() int {
	buf := processorCores()
	if buf == nil {
		return 0
	}
	return countPerformanceCores(buf)
}

// mathCPUs returns nil on Windows. This package does not pin threads to CPUs
// there.
func mathCPUs() []int32 { return nil }

// processorCores returns the SYSTEM_LOGICAL_PROCESSOR_INFORMATION_EX records
// for the physical cores, or nil when the call fails.
func processorCores() []byte {
	if procGetLogicalProcessorInformationEx.Find() != nil {
		return nil
	}
	var size uint32
	procGetLogicalProcessorInformationEx.Call(relationProcessorCore, 0, uintptr(unsafe.Pointer(&size)))
	if size == 0 {
		return nil
	}
	buf := make([]byte, size)
	r, _, _ := procGetLogicalProcessorInformationEx.Call(relationProcessorCore,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r == 0 {
		return nil
	}
	return buf[:size]
}

// countPerformanceCores counts the cores with the highest EfficiencyClass. A
// machine with one kind of core reports 0 for every core, so all of them count.
func countPerformanceCores(buf []byte) int {
	counts := map[byte]int{}
	best := -1
	for len(buf) >= 10 {
		relationship := binary.LittleEndian.Uint32(buf[0:])
		size := binary.LittleEndian.Uint32(buf[4:])
		if size < 10 || int(size) > len(buf) {
			break
		}
		if relationship == relationProcessorCore {
			class := buf[9]
			counts[class]++
			best = max(best, int(class))
		}
		buf = buf[size:]
	}
	if best < 0 {
		return 0
	}
	return counts[byte(best)]
}
