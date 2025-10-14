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
		t.Skip("skipping test, model does not have BOS token")
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
		t.Skip("skipping test, model does not have EOS token")
	}
}

func TestVocabEOT(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	eot := VocabEOT(vocab)
	if eot == TokenNull {
		t.Skip("skipping test, model does not have EOT token")
	}
}

func TestVocabSEP(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	sep := VocabSEP(vocab)
	if sep == TokenNull {
		t.Skip("skipping test, model does not have SEP token")
	}
}

func TestVocabNL(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	nl := VocabNL(vocab)
	if nl == TokenNull {
		t.Fatal("VocabNL returned TokenNull")
	}
}

func TestVocabPAD(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	pad := VocabPAD(vocab)
	if pad == TokenNull {
		t.Skip("skipping test, model does not have PAD token")
	}
}

func TestVocabMASK(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	mask := VocabMASK(vocab)
	if mask == TokenNull {
		t.Skip("skipping test, model does not have MASK token")
	}
}

func TestVocabGetAddBOS(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	addBOS := VocabGetAddBOS(vocab)
	// No specific expected value, just ensure it doesn't fail
	_ = addBOS
}

func TestVocabGetAddEOS(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	addEOS := VocabGetAddEOS(vocab)
	// No specific expected value, just ensure it doesn't fail
	_ = addEOS
}

func TestVocabGetAddSEP(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	sep := VocabSEP(vocab)
	if sep == TokenNull {
		t.Skip("skipping test, model does not have SEP token")
	}

	addSEP := VocabGetAddSEP(vocab)
	// No specific expected value, just ensure it doesn't fail
	_ = addSEP
}

func TestVocabFIMPre(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimPre := VocabFIMPre(vocab)
	if fimPre == TokenNull {
		t.Skip("skipping test, model does not have FIMPre token")
	}
}

func TestVocabFIMSuf(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimSuf := VocabFIMSuf(vocab)
	if fimSuf == TokenNull {
		t.Skip("skipping test, model does not have FIMSuf token")
	}
}

func TestVocabFIMMid(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimMid := VocabFIMMid(vocab)
	if fimMid == TokenNull {
		t.Skip("skipping test, model does not have FIMMid token")
	}
}

func TestVocabFIMPad(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimPad := VocabFIMPad(vocab)
	if fimPad == TokenNull {
		t.Skip("skipping test, model does not have FIMPad token")
	}
}

func TestVocabFIMRep(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimRep := VocabFIMRep(vocab)
	if fimRep == TokenNull {
		t.Skip("skipping test, model does not have FIMRep token")
	}
}

func TestVocabFIMSep(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimSep := VocabFIMSep(vocab)
	if fimSep == TokenNull {
		t.Skip("skipping test, model does not have FIMSep token")
	}
}

func TestVocabGetAttr(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	// Use a valid token for testing, e.g., BOS token
	token := VocabBOS(vocab)
	if token == TokenNull {
		t.Skip("skipping test, model does not have BOS token")
	}

	attr := VocabGetAttr(vocab, token)
	// No specific expected value, just ensure it doesn't fail
	t.Logf("VocabGetAttr returned attribute: %d for token: %d", attr, token)
}

func TestVocabGetScore(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	// Use a valid token for testing, e.g., BOS token
	token := VocabBOS(vocab)
	if token == TokenNull {
		t.Skip("skipping test, model does not have BOS token")
	}

	score := VocabGetScore(vocab, token)
	if score < 0 {
		t.Fatalf("VocabGetScore returned an invalid score: %f for token: %d", score, token)
	}

	t.Logf("VocabGetScore returned score: %f for token: %d", score, token)
}

func TestVocabGetText(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	// Use a valid token for testing, e.g., BOS token
	token := VocabBOS(vocab)
	if token == TokenNull {
		t.Skip("skipping test, model does not have BOS token")
	}

	text := VocabGetText(vocab, token)
	if text == "" {
		t.Fatalf("VocabGetText returned an empty string for token: %d", token)
	}

	t.Logf("VocabGetText returned text: %s for token: %d", text, token)
}

func TestGetVocabType(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	vocabType := GetVocabType(vocab)
	// No specific expected value, just ensure it doesn't fail
	t.Logf("VocabType returned type: %d", vocabType)
}
