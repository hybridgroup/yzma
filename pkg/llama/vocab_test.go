package llama

import (
	"testing"
)

func TestVocabBOS(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabBOS(vocab)
	t.Logf("VocabBOS returned token: %d", token)
}

func TestVocabEOS(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabEOS(vocab)
	t.Logf("VocabEOS returned token: %d", token)

}

func TestVocabEOT(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabEOT(vocab)
	t.Logf("VocabEOT returned token: %d", token)
}

func TestVocabSEP(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabSEP(vocab)
	t.Logf("VocabSEP returned token: %d", token)
}

func TestVocabNL(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabNL(vocab)
	t.Logf("VocabNL returned token: %d", token)
}

func TestVocabPAD(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabPAD(vocab)
	t.Logf("VocabPAD returned token: %d", token)
}

func TestVocabMASK(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabMASK(vocab)
	t.Logf("VocabMASK returned token: %d", token)
}

func TestVocabGetAddBOS(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabFIMPre(vocab)
	t.Logf("VocabFIMPre returned token: %d", token)
}

func TestVocabFIMSuf(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabFIMSuf(vocab)
	t.Logf("VocabFIMSuf returned token: %d", token)
}

func TestVocabFIMMid(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabFIMMid(vocab)
	t.Logf("VocabFIMMid returned token: %d", token)
}

func TestVocabFIMPad(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabFIMPad(vocab)
	t.Logf("VocabFIMPad returned token: %d", token)
}

func TestVocabFIMRep(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabFIMRep(vocab)
	t.Logf("VocabFIMRep returned token: %d", token)
}

func TestVocabFIMSep(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := VocabFIMSep(vocab)
	t.Logf("VocabFIMSep returned token: %d", token)
}

func TestVocabIsEOG(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	// Use a valid token for testing, e.g., BOS token
	token := VocabBOS(vocab)
	if token == TokenNull {
		t.Skip("skipping test, model does not have BOS token")
	}

	isEOG := VocabIsEOG(vocab, token)
	// No specific expected value, just ensure it doesn't fail
	t.Logf("VocabIsEOG returned: %v for token: %d", isEOG, token)
}

func TestVocabIsControl(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	// Use a valid token for testing, e.g., BOS token
	token := VocabBOS(vocab)
	if token == TokenNull {
		t.Skip("skipping test, model does not have BOS token")
	}

	isControl := VocabIsControl(vocab, token)
	// No specific expected value, just ensure it doesn't fail
	t.Logf("VocabIsControl returned: %v for token: %d", isControl, token)
}

func TestTokenToPiece(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	// Use a valid token for testing, e.g., BOS token
	token := VocabBOS(vocab)
	if token == TokenNull {
		t.Skip("skipping test, model does not have BOS token")
	}

	buf := make([]byte, 256)
	piece := TokenToPiece(vocab, token, buf, 0, true)
	if piece == 0 {
		t.Fatalf("TokenToPiece returned an empty string for token: %d", token)
	}

	t.Logf("TokenToPiece returned len: %d for token: %d", piece, token)
}

func TestVocabGetAttr(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
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
	model, err := ModelLoadFromFile(modelFile, params)
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	vocabType := GetVocabType(vocab)
	// No specific expected value, just ensure it doesn't fail
	t.Logf("VocabType returned type: %d", vocabType)
}
