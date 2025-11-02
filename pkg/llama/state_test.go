package llama

import (
	"os"
	"testing"
)

func TestStateSaveFile(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// tokenize prompt
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	count := Tokenize(vocab, prompt, nil, true, true)
	tokens := make([]Token, count)
	Tokenize(vocab, prompt, tokens, true, true)

	// create batch and decode
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	// Use a temporary file for testing
	tmp, err := os.CreateTemp("", "test_state_save.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile := tmp.Name()
	defer os.Remove(tmpFile)
	tmp.Close()

	ok := StateSaveFile(ctx, tmpFile, tokens)
	if !ok {
		t.Fatal("StateSaveFile failed")
	}

	// Check if file was created
	if _, err := os.Stat(tmpFile); err != nil {
		t.Fatalf("StateSaveFile did not create file: %v", err)
	}

	t.Logf("StateSaveFile created file: %s", tmpFile)
}

func TestStateLoadFile(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// tokenize prompt and decode, then save state
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	count := Tokenize(vocab, prompt, nil, true, true)
	tokens := make([]Token, count)
	Tokenize(vocab, prompt, tokens, true, true)
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	tmp, err := os.CreateTemp("", "test_state_load.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile := tmp.Name()
	defer os.Remove(tmpFile)
	tmp.Close()

	ok := StateSaveFile(ctx, tmpFile, tokens)
	if !ok {
		t.Fatal("StateSaveFile failed")
	}

	// Prepare output buffer for loading
	outTokens := make([]Token, len(tokens))
	var nTokenCountOut uint64

	ok = StateLoadFile(ctx, tmpFile, outTokens, uint64(len(outTokens)), &nTokenCountOut)
	if !ok {
		t.Fatal("StateLoadFile failed")
	}
	if nTokenCountOut == 0 || nTokenCountOut > uint64(len(outTokens)) {
		t.Fatalf("StateLoadFile loaded unexpected number of tokens: %d", nTokenCountOut)
	}

	t.Logf("StateLoadFile loaded %d tokens from %s", nTokenCountOut, tmpFile)
}

func TestStateGetSize(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	size := StateGetSize(ctx)
	t.Logf("StateGetSize returned: %d bytes", size)
	if size == 0 {
		t.Fatal("StateGetSize returned 0, expected non-zero state size")
	}
}

func TestStateGetData(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// Get the required state size
	size := StateGetSize(ctx)
	if size == 0 {
		t.Fatal("StateGetSize returned 0, expected non-zero state size")
	}

	// Allocate a buffer
	buf := make([]byte, size)
	nCopied := StateGetData(ctx, buf)
	t.Logf("StateGetData copied %d bytes", nCopied)
	if nCopied == 0 || nCopied > size {
		t.Fatalf("StateGetData copied an unexpected number of bytes: %d", nCopied)
	}
}

func TestStateSetData(t *testing.T) {
	modelFile := testModelFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)

	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// Save state to buffer
	size := StateGetSize(ctx)
	if size == 0 {
		t.Fatal("StateGetSize returned 0, expected non-zero state size")
	}
	buf := make([]byte, size)
	nCopied := StateGetData(ctx, buf)
	if nCopied == 0 || nCopied > size {
		t.Fatalf("StateGetData copied an unexpected number of bytes: %d", nCopied)
	}

	// Set state from buffer
	nRead := StateSetData(ctx, buf)
	t.Logf("StateSetData read %d bytes", nRead)
	if nRead == 0 || nRead > size {
		t.Fatalf("StateSetData read an unexpected number of bytes: %d", nRead)
	}
}

func TestStateSeqGetSizeAndData(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)
	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// Tokenize and decode
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	count := Tokenize(vocab, prompt, nil, true, true)
	tokens := make([]Token, count)
	Tokenize(vocab, prompt, tokens, true, true)
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	seqId := SeqId(1)
	size := StateSeqGetSize(ctx, seqId)
	if size == 0 {
		t.Fatal("StateSeqGetSize returned 0")
	}
	buf := make([]byte, size)
	nCopied := StateSeqGetData(ctx, buf, seqId)
	if nCopied == 0 || nCopied > size {
		t.Fatalf("StateSeqGetData copied unexpected number of bytes: %d", nCopied)
	}

	nRead := StateSeqSetData(ctx, buf, seqId)
	if nRead == 0 || nRead > size {
		t.Fatalf("StateSeqSetData read unexpected number of bytes: %d", nRead)
	}
}

func TestStateSeqSaveLoadFile(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)
	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	// Tokenize and decode
	prompt := "This is a test"
	vocab := ModelGetVocab(model)
	count := Tokenize(vocab, prompt, nil, true, true)
	tokens := make([]Token, count)
	Tokenize(vocab, prompt, tokens, true, true)
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	seqId := SeqId(1)
	tmp, err := os.CreateTemp("", "test_state_seq_save.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile := tmp.Name()
	defer os.Remove(tmpFile)
	tmp.Close()

	nSaved := StateSeqSaveFile(ctx, tmpFile, seqId, tokens)
	if nSaved == 0 {
		t.Fatal("StateSeqSaveFile failed")
	}

	outTokens := make([]Token, len(tokens))
	var nTokenCountOut uint64
	nLoaded := StateSeqLoadFile(ctx, tmpFile, seqId, outTokens, uint64(len(outTokens)), &nTokenCountOut)
	if nLoaded == 0 {
		t.Fatal("StateSeqLoadFile failed")
	}
	if nTokenCountOut == 0 || nTokenCountOut > uint64(len(outTokens)) {
		t.Fatalf("StateSeqLoadFile loaded unexpected number of tokens: %d", nTokenCountOut)
	}
}

func TestStateSeqGetSizeDataExt(t *testing.T) {
	modelFile := testModelFileName(t)
	testSetup(t)
	defer testCleanup(t)

	model := ModelLoadFromFile(modelFile, ModelDefaultParams())
	defer ModelFree(model)
	ctx := InitFromModel(model, ContextDefaultParams())
	defer Free(ctx)

	seqId := SeqId(1)
	flags := uint32(1) // e.g. LLAMA_STATE_SEQ_FLAGS_SWA_ONLY
	size := StateSeqGetSizeExt(ctx, seqId, flags)
	if size == 0 {
		t.Fatal("StateSeqGetSizeExt returned 0")
	}
	buf := make([]byte, size)
	nCopied := StateSeqGetDataExt(ctx, buf, seqId, flags)
	if nCopied == 0 || nCopied > size {
		t.Fatalf("StateSeqGetDataExt copied unexpected number of bytes: %d", nCopied)
	}

	nRead := StateSeqSetDataExt(ctx, buf, seqId, flags)
	if nRead == 0 || nRead > size {
		t.Fatalf("StateSeqSetDataExt read unexpected number of bytes: %d", nRead)
	}
}
