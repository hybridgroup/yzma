package llama

import (
	"testing"
)

func testSetup(t *testing.T) {
	testPath := "."
	if err := Load(testPath); err != nil {
		t.Fatal("unable to load library", err.Error())
	}

	BackendInit()
	GGMLBackendLoadAll()
}

func testCleanup(t *testing.T) {
	BackendFree()
}
