package llama

import (
	"regexp"
	"testing"
)

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

	match := []string{
		`\.ffn_up_exps`,
		`\.ffn_down_exps`,
		`\.ffn_gate_exps`,
		`\.ffn_up_chexps`,
		`\.ffn_down_chexps`,
		`\.ffn_gate_chexps`,
	}
	for _, s := range match {
		if !re.MatchString(s) {
			t.Errorf("pattern should match %q", s)
		}
	}

	noMatch := []string{
		`\.ffn_up_exp`,
		".attn_q",
		`\.ffn_down_chexp`,
	}
	for _, s := range noMatch {
		if re.MatchString(s) {
			t.Errorf("pattern should not match %q", s)
		}
	}
}
