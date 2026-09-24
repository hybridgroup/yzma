package llama

import (
	"encoding/binary"
	"unsafe"

	"golang.org/x/sys/windows"
)

var procGetLogicalProcessorInformationEx = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetLogicalProcessorInformationEx")

// relationProcessorCore asks for one record for each physical core.
const relationProcessorCore = 0

// mathCores counts the cores of this machine that do the arithmetic well. It
// is the number of physical cores, and on a machine with performance cores and
// efficiency cores it leaves the efficiency cores out. The count is 0 when
// the system says nothing, and then the caller uses its own default.
func mathCores() int {
	buf := processorCores()
	if buf == nil {
		return 0
	}
	return countPerformanceCores(buf)
}

// mathCPUs gives nothing on Windows. This package does not hold threads to a
// CPU there.
func mathCPUs() []int32 { return nil }

// processorCores gives the SYSTEM_LOGICAL_PROCESSOR_INFORMATION_EX records of
// the physical cores, or nil when the call fails.
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
// machine with one kind of core gives 0 to each core, thus all of them count.
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
