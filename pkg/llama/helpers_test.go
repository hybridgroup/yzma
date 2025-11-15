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

func testEncoderModelFileName(t *testing.T) string {
	if os.Getenv("YZMA_TEST_ENCODER_MODEL") == "" {
		t.Skip("no YZMA_TEST_ENCODER_MODEL skipping test")
	}

	return os.Getenv("YZMA_TEST_ENCODER_MODEL")
}

func testMMMModelFileName(t *testing.T) string {
	if os.Getenv("YZMA_TEST_MMMODEL") == "" {
		t.Skip("no YZMA_TEST_MMMODEL skipping test")
	}

	return os.Getenv("YZMA_TEST_MMMODEL")
}

func benchmarkSetup(b *testing.B) {
	testPath := "."
	if err := Load(testPath); err != nil {
		b.Fatal("unable to load library", err.Error())
	}

	LogSet(LogSilent())

	Init()
}

func benchmarkCleanup(_ *testing.B) {
	LogSet(LogNormal)

	BackendFree()
}

func benchmarkModelFileName(b *testing.B) string {
	if os.Getenv("YZMA_TEST_MODEL") == "" {
		b.Skip("no YZMA_TEST_MODEL skipping test")
	}

	return os.Getenv("YZMA_TEST_MODEL")
}
