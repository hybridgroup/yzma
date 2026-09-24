package llama

import (
	"encoding/binary"
	"testing"
)

// coreRecord makes one record of GetLogicalProcessorInformationEx for a core.
func coreRecord(class byte) []byte {
	b := make([]byte, 48)
	binary.LittleEndian.PutUint32(b[0:], relationProcessorCore)
	binary.LittleEndian.PutUint32(b[4:], uint32(len(b)))
	b[9] = class
	return b
}

func TestCountPerformanceCores(t *testing.T) {
	var hybrid, flat []byte
	for range 8 {
		hybrid = append(hybrid, coreRecord(1)...)
		flat = append(flat, coreRecord(0)...)
	}
	for range 16 {
		hybrid = append(hybrid, coreRecord(0)...)
	}
	if n := countPerformanceCores(hybrid); n != 8 {
		t.Errorf("hybrid machine gave %d, want 8", n)
	}
	if n := countPerformanceCores(flat); n != 8 {
		t.Errorf("machine with one kind of core gave %d, want 8", n)
	}
	if n := countPerformanceCores(nil); n != 0 {
		t.Errorf("no records gave %d, want 0", n)
	}
}
