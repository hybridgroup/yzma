package llama

import "testing"

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
