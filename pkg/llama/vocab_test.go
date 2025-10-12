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
		t.Fatal("VocabEOT returned TokenNull")
	}
}

func TestVocabSEP(t *testing.T) {
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	sep := VocabSEP(vocab)
	if sep == TokenNull {
		t.Fatal("VocabSEP returned TokenNull")
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
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	pad := VocabPAD(vocab)
	if pad == TokenNull {
		t.Fatal("VocabPAD returned TokenNull")
	}
}

func TestVocabMASK(t *testing.T) {
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	mask := VocabMASK(vocab)
	if mask == TokenNull {
		t.Fatal("VocabMASK returned TokenNull")
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
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	addSEP := VocabGetAddSEP(vocab)
	// No specific expected value, just ensure it doesn't fail
	_ = addSEP
}

func TestVocabFIMPre(t *testing.T) {
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimPre := VocabFIMPre(vocab)
	if fimPre == TokenNull {
		t.Fatal("VocabFIMPre returned TokenNull")
	}
}

func TestVocabFIMSuf(t *testing.T) {
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimSuf := VocabFIMSuf(vocab)
	if fimSuf == TokenNull {
		t.Fatal("VocabFIMSuf returned TokenNull")
	}
}

func TestVocabFIMMid(t *testing.T) {
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimMid := VocabFIMMid(vocab)
	if fimMid == TokenNull {
		t.Fatal("VocabFIMMid returned TokenNull")
	}
}

func TestVocabFIMPad(t *testing.T) {
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimPad := VocabFIMPad(vocab)
	if fimPad == TokenNull {
		t.Fatal("VocabFIMPad returned TokenNull")
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
		t.Fatal("VocabFIMRep returned TokenNull")
	}
}

func TestVocabFIMSep(t *testing.T) {
	t.Skip("TODO: test this function for model with vocab that has this")

	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	params := ModelDefaultParams()
	model := ModelLoadFromFile(modelFile, params)
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	fimSep := VocabFIMSep(vocab)
	if fimSep == TokenNull {
		t.Fatal("VocabFIMSep returned TokenNull")
	}
}
