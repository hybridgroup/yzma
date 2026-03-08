package llama

import "testing"

func TestLogSilent(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	LogSet(LogSilent())

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer ModelFree(model)

	t.Log("Logs should be silent on this test")
}

func TestLogNormal(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	LogSet(LogNormal)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("Failed to load model from file: %v", err)
	}
	defer ModelFree(model)

	t.Log("Logs should be normal on this test")
}

func TestLogGet(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	cb, ud := LogGet()
	t.Logf("LogGet returned callback=%v, userData=%v", cb, ud)
}
