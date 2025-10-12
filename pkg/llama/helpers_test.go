package llama

import (
	"os"
	"testing"
)

func testSetup(t *testing.T) {
	testPath := "."
	if err := Load(testPath); err != nil {
		t.Fatal("unable to load library", err.Error())
	}

	Init()
}

func testCleanup(_ *testing.T) {
	BackendFree()
}

func testModelFileName(t *testing.T) string {
	if os.Getenv("YZMA_TEST_MODEL") == "" {
		t.Skip("no YZMA_TEST_MODEL skipping test")
	}

	return os.Getenv("YZMA_TEST_MODEL")
}
