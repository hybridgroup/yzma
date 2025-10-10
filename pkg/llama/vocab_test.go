package llama

import (
	"testing"
)

func TestVocabBOS(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	bos := VocabBOS(vocab)
	if bos == TokenNull {
		t.Fatal("VocabBOS returned TokenNull")
	}
}

func TestVocabEOS(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	eos := VocabEOS(vocab)
	if eos == TokenNull {
		t.Fatal("VocabEOS returned TokenNull")
	}
}

func TestVocabIsControl(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabBOS(vocab) // Example token
	isControl := VocabIsControl(vocab, token)
	// Assuming BOS is a control token, adjust as needed
	if !isControl {
		t.Fatal("VocabIsControl incorrectly returned false for BOS token")
	}
}

func TestVocabNTokens(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	nTokens := VocabNTokens(vocab)
	if nTokens <= 0 {
		t.Fatal("VocabNTokens returned an invalid value:", nTokens)
	}
}
