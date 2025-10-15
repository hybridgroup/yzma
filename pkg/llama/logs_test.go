package llama

import "testing"

func TestLogSilent(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	LogSet(LogSilent(), uintptr(0))

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	t.Log("Logs should be silent on this test")
}
