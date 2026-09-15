package speculative

import (
	"os"
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func testSetup(t *testing.T) {
	if os.Getenv("YZMA_LIB") == "" {
		t.Fatal("no YZMA_LIB set for tests")
	}
	testPath := os.Getenv("YZMA_LIB")
	if err := llama.Load(testPath); err != nil {
		t.Fatal("unable to load library", err.Error())
	}

	llama.Init()

	if err := Load(testPath); err != nil {
		t.Fatal("unable to load NextN functions", err.Error())
	}
}

func testCleanup(_ *testing.T) {
	llama.BackendFree()
}

func testModelFileName(t *testing.T) string {
	if os.Getenv("YZMA_TEST_MODEL") == "" {
		t.Skip("no YZMA_TEST_MODEL skipping test")
	}

	return os.Getenv("YZMA_TEST_MODEL")
}
