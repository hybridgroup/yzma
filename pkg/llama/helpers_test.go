package llama

import (
	"os"
	"strings"
	"testing"
)

func testSetup(t *testing.T) {
	if os.Getenv("YZMA_LIB") == "" {
		t.Fatal("no YZMA_LIB set for tests")
	}
	testPath := os.Getenv("YZMA_LIB")
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

func testSplitModelFileNames(t *testing.T) []string {
	if os.Getenv("YZMA_TEST_SPLIT_MODELS") == "" {
		t.Skip("no YZMA_TEST_SPLIT_MODELS skipping test")
	}

	return strings.Split(os.Getenv("YZMA_TEST_SPLIT_MODELS"), ",")
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

func testLoraModelFileName(t *testing.T) string {
	if os.Getenv("YZMA_TEST_LORA_MODEL") == "" {
		t.Skip("no YZMA_TEST_LORA_MODEL skipping test")
	}

	return os.Getenv("YZMA_TEST_LORA_MODEL")
}

func testLoraAdaptorFileName(t *testing.T) string {
	if os.Getenv("YZMA_TEST_LORA_ADAPTER") == "" {
		t.Skip("no YZMA_TEST_LORA_ADAPTER skipping test")
	}

	return os.Getenv("YZMA_TEST_LORA_ADAPTER")
}

func benchmarkSetup(b *testing.B) {
	if os.Getenv("YZMA_LIB") == "" {
		b.Fatal("no YZMA_LIB set for tests")
	}
	testPath := os.Getenv("YZMA_LIB")
	if err := Load(testPath); err != nil {
		b.Fatal("unable to load library", err.Error())
	}

	LogSet(LogSilent())

	Init()
}

func benchmarkCleanup(_ *testing.B) {
	LogSet(LogNormal)

	Close()
}

func benchmarkModelFileName(b *testing.B) string {
	if os.Getenv("YZMA_BENCHMARK_MODEL") == "" {
		b.Skip("no YZMA_BENCHMARK_MODEL skipping test")
	}

	return os.Getenv("YZMA_BENCHMARK_MODEL")
}
