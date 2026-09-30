package llama

import "testing"

func TestCloseUnloadsAllBackends(t *testing.T) {
	testSetup(t)
	defer Init()

	if GGMLBackendRegCount() == 0 {
		t.Skip("no backends registered")
	}

	Close()

	if n := GGMLBackendRegCount(); n != 0 {
		t.Fatalf("GGMLBackendRegCount after Close = %d, want 0", n)
	}
}
