package llama

import (
	"regexp"
	"testing"
)

func TestGGMLBackendCpuBufferType(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// Binding this pointer-returning function with a void return type makes
	// libffi write nothing, so the result is always NULL and every
	// TensorBuftOverride built from it aborts llama.cpp during model load.
	buft := GGMLBackendCpuBufferType()
	if buft == 0 {
		t.Fatal("GGMLBackendCpuBufferType returned NULL")
	}
	t.Logf("GGMLBackendCpuBufferType returned: %#x", uintptr(buft))

	if o := NewTensorBuftAllFFNExprsOverride(); o.Type == 0 {
		t.Fatal("NewTensorBuftAllFFNExprsOverride produced a NULL buffer type")
	}
}

func TestGGMLBackendDevCount(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	count := GGMLBackendDeviceCount()
	t.Logf("GGMLBackendDeviceCount returned: %d", count)
	if count == 0 {
		t.Skip("No backend devices found")
	}
}

func TestGGMLBackendDeviceGet(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	count := GGMLBackendDeviceCount()
	if count == 0 {
		t.Skip("No backend devices to get")
	}
	dev := GGMLBackendDeviceGet(0)
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceGet returned 0 for index 0")
	}
	t.Logf("GGMLBackendDeviceGet(0) returned: %v", dev)
}

func TestGGMLBackendDeviceByType(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// Try CPU device type
	dev := GGMLBackendDeviceByType(GGMLBackendDeviceTypeCPU)
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceByType(GGMLBackendDeviceTypeCPU) returned 0")
	} else {
		t.Logf("GGMLBackendDeviceByType(GGMLBackendDeviceTypeCPU) returned: %v", dev)
	}
}

func TestGGMLBackendDeviceByName(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	count := GGMLBackendDeviceCount()
	if count == 0 {
		t.Skip("No backend devices to get name from")
	}

	dev := GGMLBackendDeviceByName("CPU")
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceByName(\"CPU\") returned 0")
	} else {
		t.Logf("GGMLBackendDeviceByName(\"CPU\") returned: %v", dev)
	}
}

func TestGGMLBackendDeviceName(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	count := GGMLBackendDeviceCount()
	if count == 0 {
		t.Skip("No backend devices to get name from")
	}

	dev := GGMLBackendDeviceGet(0)
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceGet(0) returned 0")
	}

	name := GGMLBackendDeviceName(dev)
	if name == "" {
		t.Fatal("GGMLBackendDeviceName returned empty string")
	} else {
		t.Logf("GGMLBackendDeviceName returned: %s", name)
	}
}

func TestGGMLBackendRegCount(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	count := GGMLBackendRegCount()
	t.Logf("GGMLBackendRegCount returned: %d", count)
	if count == 0 {
		t.Skip("No backend registrations found")
	}
}

func TestGGMLBackendRegGet(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	count := GGMLBackendRegCount()
	if count == 0 {
		t.Skip("No backend registrations to get")
	}
	reg := GGMLBackendRegGet(0)
	if reg == 0 {
		t.Fatal("GGMLBackendRegGet returned 0 for index 0")
	}
	t.Logf("GGMLBackendRegGet(0) returned: %v", reg)
}

func TestGGMLBackendRegByName(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	reg := GGMLBackendRegByName("CPU")
	if reg == 0 {
		t.Log("GGMLBackendRegByName(\"CPU\") returned 0 (may be expected if no CPU backend)")
	} else {
		t.Logf("GGMLBackendRegByName(\"CPU\") returned: %v", reg)
	}
}

func TestGGMLBackendUnload(t *testing.T) {
	t.Skip("skipped: unloading a backend invalidates function pointers for subsequent tests")

	testSetup(t)
	defer testCleanup(t)

	count := GGMLBackendRegCount()
	if count == 0 {
		t.Skip("No backend registrations to unload")
	}
	reg := GGMLBackendRegGet(0)
	if reg == 0 {
		t.Skip("GGMLBackendRegGet returned 0 for index 0")
	}
	// Should not panic or error
	GGMLBackendUnload(reg)
	t.Logf("GGMLBackendUnload succeeded for reg: %v", reg)
}

func TestGGMLBackendDeviceMemoryNilDevice(t *testing.T) {
	free, total := GGMLBackendDeviceMemory(0)
	if free != 0 || total != 0 {
		t.Fatalf("GGMLBackendDeviceMemory(0) = (%d, %d), want (0, 0)", free, total)
	}
}

func TestGGMLBackendDeviceMemory(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if GGMLBackendDeviceCount() == 0 {
		t.Skip("No backend devices available")
	}

	dev := GGMLBackendDeviceGet(0)
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceGet(0) returned 0")
	}

	free, total := GGMLBackendDeviceMemory(dev)
	if total > 0 && free > total {
		t.Fatalf("free (%d) > total (%d)", free, total)
	}
	t.Logf("Device memory: free=%d, total=%d", free, total)
}

func TestMoEExpertTensorPattern(t *testing.T) {
	re := regexp.MustCompile(MoEExpertTensorPattern)

	// Match against real llama.cpp tensor names, not against strings that
	// contain a literal backslash.
	match := []string{
		"blk.0.ffn_up_exps.weight",
		"blk.12.ffn_down_exps.weight",
		"blk.3.ffn_gate_exps.weight",
		"blk.0.ffn_up_chexps.weight",
		"blk.7.ffn_down_chexps.weight",
		"blk.7.ffn_gate_chexps.weight",
	}
	for _, s := range match {
		if !re.MatchString(s) {
			t.Errorf("pattern should match %q", s)
		}
	}

	noMatch := []string{
		"blk.0.ffn_up_exp.weight",
		"blk.0.attn_q.weight",
		"blk.0.ffn_down_chexp.weight",
	}
	for _, s := range noMatch {
		if re.MatchString(s) {
			t.Errorf("pattern should not match %q", s)
		}
	}

	// the per-block pattern must keep the block prefix escaped correctly
	blockRe := regexp.MustCompile(ffnExprBlockRegex(3))
	if !blockRe.MatchString("blk.3.ffn_up_exps.weight") {
		t.Error("block pattern should match its own block")
	}
	if blockRe.MatchString("blk.4.ffn_up_exps.weight") {
		t.Error("block pattern should not match another block")
	}
}
