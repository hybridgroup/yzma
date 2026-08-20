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

func TestGGMLBackendDeviceDescription(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if GGMLBackendDeviceCount() == 0 {
		t.Skip("No backend devices to get description from")
	}

	dev := GGMLBackendDeviceGet(0)
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceGet(0) returned 0")
	}

	desc := GGMLBackendDeviceDescription(dev)
	if desc == "" {
		t.Fatal("GGMLBackendDeviceDescription returned empty string")
	}
	t.Logf("GGMLBackendDeviceDescription returned: %s", desc)
}

func TestGGMLBackendDevType(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// A device found by type must report that same type back. This catches a
	// return buffer that is too narrow and an enum that moved upstream.
	dev := GGMLBackendDeviceByType(GGMLBackendDeviceTypeCPU)
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceByType(GGMLBackendDeviceTypeCPU) returned 0")
	}

	if got := GGMLBackendDevType(dev); got != GGMLBackendDeviceTypeCPU {
		t.Fatalf("GGMLBackendDevType = %v, want %v", got, GGMLBackendDeviceTypeCPU)
	}

	for i := uint64(0); i < GGMLBackendDeviceCount(); i++ {
		d := GGMLBackendDeviceGet(i)
		if d == 0 {
			continue
		}
		t.Logf("device %d (%s) type: %v", i, GGMLBackendDeviceName(d), GGMLBackendDevType(d))
	}
}

func TestGGMLBackendDeviceBackendReg(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if GGMLBackendDeviceCount() == 0 {
		t.Skip("No backend devices to get the registration from")
	}

	dev := GGMLBackendDeviceGet(0)
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceGet(0) returned 0")
	}

	reg := GGMLBackendDeviceBackendReg(dev)
	if reg == 0 {
		t.Fatal("GGMLBackendDeviceBackendReg returned 0")
	}
	t.Logf("GGMLBackendDeviceBackendReg returned: %v", reg)
}

func TestGGMLBackendRegName(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if GGMLBackendDeviceCount() == 0 {
		t.Skip("No backend devices to get the registration name from")
	}

	dev := GGMLBackendDeviceGet(0)
	if dev == 0 {
		t.Fatal("GGMLBackendDeviceGet(0) returned 0")
	}

	reg := GGMLBackendDeviceBackendReg(dev)
	if reg == 0 {
		t.Fatal("GGMLBackendDeviceBackendReg returned 0")
	}

	name := GGMLBackendRegName(reg)
	if name == "" {
		t.Fatal("GGMLBackendRegName returned empty string")
	}
	t.Logf("GGMLBackendRegName returned: %s", name)

	// the name must round-trip through the lookup it came from
	if cpu := GGMLBackendRegByName("CPU"); cpu != 0 {
		if got := GGMLBackendRegName(cpu); got != "CPU" {
			t.Fatalf("GGMLBackendRegName(GGMLBackendRegByName(\"CPU\")) = %q, want \"CPU\"", got)
		}
	}
}

func TestGGMLBackendDeviceAccessorsNilDevice(t *testing.T) {
	if desc := GGMLBackendDeviceDescription(0); desc != "" {
		t.Fatalf("GGMLBackendDeviceDescription(0) = %q, want empty string", desc)
	}

	if reg := GGMLBackendDeviceBackendReg(0); reg != 0 {
		t.Fatalf("GGMLBackendDeviceBackendReg(0) = %v, want 0", reg)
	}

	if name := GGMLBackendRegName(0); name != "" {
		t.Fatalf("GGMLBackendRegName(0) = %q, want empty string", name)
	}
}

func TestGGMLBackendDeviceTypeString(t *testing.T) {
	tests := []struct {
		devType GGMLBackendDeviceType
		want    string
	}{
		{GGMLBackendDeviceTypeCPU, "CPU"},
		{GGMLBackendDeviceTypeGPU, "GPU"},
		{GGMLBackendDeviceTypeIGPU, "IGPU"},
		{GGMLBackendDeviceTypeACCEL, "ACCEL"},
		{GGMLBackendDeviceTypeMETA, "META"},
		{GGMLBackendDeviceType(99), "GGMLBackendDeviceType(99)"},
	}

	for _, tt := range tests {
		if got := tt.devType.String(); got != tt.want {
			t.Errorf("GGMLBackendDeviceType(%d).String() = %q, want %q", int32(tt.devType), got, tt.want)
		}
	}
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

func TestGGMLTypeSizeAndBlockSize(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	tests := []struct {
		typ       GGMLType
		blockSize int64
		typeSize  uint64
	}{
		{GGMLTypeF32, 1, 4},
		{GGMLTypeF16, 1, 2},
		{GGMLTypeQ8_0, 32, 34},
		{GGMLTypeQ4_0, 32, 18},
	}

	for _, tt := range tests {
		if got := GGMLBlockSize(tt.typ); got != tt.blockSize {
			t.Errorf("GGMLBlockSize(%v) = %d, want %d", tt.typ, got, tt.blockSize)
		}

		if got := GGMLTypeSize(tt.typ); got != tt.typeSize {
			t.Errorf("GGMLTypeSize(%v) = %d, want %d", tt.typ, got, tt.typeSize)
		}
	}
}

func TestGGMLRowSize(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	types := []GGMLType{GGMLTypeF32, GGMLTypeF16, GGMLTypeQ8_0, GGMLTypeQ4_0, GGMLTypeQ4_K}

	const ne = 4096

	for _, typ := range types {
		blockSize := GGMLBlockSize(typ)
		if blockSize <= 0 {
			t.Fatalf("GGMLBlockSize(%v) = %d, want a positive block size", typ, blockSize)
		}
		if ne%blockSize != 0 {
			t.Fatalf("GGMLBlockSize(%v) = %d, which does not divide %d", typ, blockSize, ne)
		}

		want := GGMLTypeSize(typ) * uint64(ne) / uint64(blockSize)
		if got := GGMLRowSize(typ, ne); got != want {
			t.Errorf("GGMLRowSize(%v, %d) = %d, want %d", typ, ne, got, want)
		}

		if got := GGMLRowSize(typ, 0); got != 0 {
			t.Errorf("GGMLRowSize(%v, 0) = %d, want 0", typ, got)
		}
	}
}

func TestGGMLTypeName(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if got := GGMLTypeName(GGMLTypeF32); got != "f32" {
		t.Errorf("GGMLTypeName(GGMLTypeF32) = %q, want %q", got, "f32")
	}

	if got := GGMLTypeName(GGMLTypeF16); got != "f16" {
		t.Errorf("GGMLTypeName(GGMLTypeF16) = %q, want %q", got, "f16")
	}

	name := GGMLTypeName(GGMLTypeQ4_K)
	if name == "" {
		t.Error("GGMLTypeName(GGMLTypeQ4_K) returned an empty name")
	}

	if got := GGMLTypeQ4_K.String(); got != name {
		t.Errorf("GGMLTypeQ4_K.String() = %q, want %q", got, name)
	}
}
