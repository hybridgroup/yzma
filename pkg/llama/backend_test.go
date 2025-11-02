package llama

import (
	"testing"
)

func TestMaxDevices(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	maxDevices := MaxDevices()
	if maxDevices == 0 {
		t.Fatal("MaxDevices returned 0, which is invalid")
	}
	t.Logf("MaxDevices returned: %d", maxDevices)
}

func TestMaxParallelSequences(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	maxParallelSequences := MaxParallelSequences()
	if maxParallelSequences == 0 {
		t.Fatal("MaxParallelSequences returned 0, which is invalid")
	}
	t.Logf("MaxParallelSequences returned: %d", maxParallelSequences)
}

func TestSupportsMmap(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	supportsMmap := SupportsMmap()
	t.Logf("SupportsMmap returned: %v", supportsMmap)
}

func TestSupportsMlock(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	supportsMlock := SupportsMlock()
	t.Logf("SupportsMlock returned: %v", supportsMlock)
}

func TestSupportsGpuOffload(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	supportsGpuOffload := SupportsGpuOffload()
	t.Logf("SupportsGpuOffload returned: %v", supportsGpuOffload)
}

func TestSupportsRpc(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	supportsRpc := SupportsRpc()
	t.Logf("SupportsRpc returned: %v", supportsRpc)
}

func TestTimeUs(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	time := TimeUs()
	if time <= 0 {
		t.Fatal("TimeUs returned an invalid value:", time)
	}
	t.Logf("TimeUs returned: %d microseconds", time)
}

func TestFlashAttnTypeName(t *testing.T) {
	var flashAttnType FlashAttentionType = FlashAttentionTypeAuto

	name := FlashAttnTypeName(flashAttnType)
	if name == "" {
		t.Fatal("FlashAttnTypeName returned empty string")
	}
	t.Logf("FlashAttnTypeName returned: %s", name)
}

func TestNumaInit(t *testing.T) {
	var strategy NumaStrategy = NumaStrategyDisabled

	// Should not panic or error
	NumaInit(strategy)
	t.Logf("NumaInit called with strategy: %d", strategy)
}
