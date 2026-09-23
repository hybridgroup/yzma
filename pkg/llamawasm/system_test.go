//go:build js && wasm

package llamawasm

import "testing"

func TestSystemInfo(t *testing.T) {
	fakeABI9(t)

	if got := PrintSystemInfo(); got != "CPU : WASM_SIMD = 1 | " {
		t.Errorf("PrintSystemInfo gave %q", got)
	}
	if got := FtypeName(FtypeMostlyQ4_K_M); got != "Q4_K - Medium" {
		t.Errorf("FtypeName gave %q", got)
	}
	// This does not fit an int32, thus the shim gives it as a double.
	if got := TimeUs(); got != 5000000000123 {
		t.Errorf("TimeUs gave %d", got)
	}
	if MaxDevices() != 16 || MaxParallelSequences() != 256 || !SupportsGpuOffload() {
		t.Error("the limits of llama.cpp are not correct")
	}
}

func TestSystemOldModule(t *testing.T) {
	fakeOld(t)

	if PrintSystemInfo() != "" || TimeUs() != 0 || MaxDevices() != 0 || SupportsGpuOffload() {
		t.Error("a module of an earlier version must give a zero value")
	}
}
