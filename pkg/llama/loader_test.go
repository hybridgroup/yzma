package llama

import (
	"os"
	"os/exec"
	"testing"
)

// Close cannot be undone when backends are linked into libggml, as on macOS.
// So the check runs in a child process and leaves the other tests alone.
func TestCloseUnloadsAllBackends(t *testing.T) {
	if os.Getenv("YZMA_TEST_CLOSE_CHILD") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestCloseUnloadsAllBackends$", "-test.v")
		cmd.Env = append(os.Environ(), "YZMA_TEST_CLOSE_CHILD=1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("child process failed: %v\n%s", err, out)
		}
		return
	}

	testSetup(t)

	if GGMLBackendRegCount() == 0 {
		t.Skip("no backends registered")
	}

	Close()

	if n := GGMLBackendRegCount(); n != 0 {
		t.Fatalf("GGMLBackendRegCount after Close = %d, want 0", n)
	}
}
