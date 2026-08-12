package llama

import (
	"math"
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

	// Scores are log-probabilities, so they are commonly negative. Bound the
	// magnitude rather than the sign: a float return read through an ffi.Arg
	// is off by orders of magnitude (~1e9 and up), never merely negative.
	for i := range VocabNTokens(vocab) {
		token := Token(i)
		if score := VocabGetScore(vocab, token); math.Abs(float64(score)) > 1e6 {
			t.Fatalf("VocabGetScore returned %g for token %d, want a plausible score", score, token)
		}
	}
}

// TestVocabTokenOutOfRange checks that the token accessors reject an id that
// is not in the vocabulary. llama.cpp resolves them through id_to_token.at(id),
// which throws std::out_of_range; that exception cannot unwind across libffi,
// so an unguarded call aborts the process instead of failing this test.
func TestVocabTokenOutOfRange(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	for _, token := range []Token{-1, TokenNull, Token(VocabNTokens(vocab))} {
		if score := VocabGetScore(vocab, token); score != 0 {
			t.Errorf("VocabGetScore(%d) = %g, want 0", token, score)
		}
		if text := VocabGetText(vocab, token); text != "" {
			t.Errorf("VocabGetText(%d) = %q, want empty", token, text)
		}
		if attr := VocabGetAttr(vocab, token); attr != TokenAttrUnknown {
			t.Errorf("VocabGetAttr(%d) = %v, want TokenAttrUnknown", token, attr)
		}
		if VocabIsEOG(vocab, token) {
			t.Errorf("VocabIsEOG(%d) = true, want false", token)
		}
		if VocabIsControl(vocab, token) {
			t.Errorf("VocabIsControl(%d) = true, want false", token)
		}
		if n := TokenToPiece(vocab, token, make([]byte, 64), 0, false); n != 0 {
			t.Errorf("TokenToPiece(%d) = %d, want 0", token, n)
		}
	}
}

// TestVocabGetScoreNonZero pins the float return type of llama_vocab_get_score.
// The bound in TestVocabGetScore cannot do that on its own: every token of a
// BPE vocab scores exactly 0.0, which is also what a float misread through an
// ffi.Arg yields. The encoder model has a sentencepiece vocab with real
// scores, so reading one through an ffi.Arg lands far outside the bound.
func TestVocabGetScoreNonZero(t *testing.T) {
	modelFile := testEncoderModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	token := Token(TokenNull)
	for i := range VocabNTokens(vocab) {
		if VocabGetScore(vocab, Token(i)) != 0 {
			token = Token(i)
			break
		}
	}
	if token == TokenNull {
		t.Skip("skipping test, model has no token with a non-zero score")
	}

	score := VocabGetScore(vocab, token)
	if math.Abs(float64(score)) > 1e6 {
		t.Fatalf("VocabGetScore returned %g for token %d, want a plausible score", score, token)
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

func TestVocabGetSuppressTokens(t *testing.T) {
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

	// Most models do not specify any tokens to suppress, so an empty result is valid.
	tokens := VocabGetSuppressTokens(vocab)
	t.Logf("VocabGetSuppressTokens returned %d tokens: %v", len(tokens), tokens)
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
