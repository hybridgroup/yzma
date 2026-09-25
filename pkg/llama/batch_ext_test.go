package llama

import (
	"errors"
	"slices"
	"testing"
)

func testBatchExtContext(t *testing.T, params ContextParams) (Model, Context) {
	t.Helper()
	modelFile := testModelFileName(t)

	testSetup(t)
	t.Cleanup(func() { testCleanup(t) })

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	t.Cleanup(func() { ModelFree(model) })

	ctx, err := InitFromModel(model, params)
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	t.Cleanup(func() { Free(ctx) })

	return model, ctx
}

func TestBatchExtInvalidHandles(t *testing.T) {
	if _, err := BatchExtInit(0); err == nil {
		t.Fatal("BatchExtInit accepted a zero context")
	}
	if err := BatchExtFree(0); err == nil {
		t.Fatal("BatchExtFree accepted a zero batch")
	}
	if _, err := BatchExtAddToken(0, 0, 1); err == nil {
		t.Fatal("BatchExtAddToken accepted a zero batch")
	}
	if _, err := Process(0, ProcessTypeDecode, 1); err == nil {
		t.Fatal("Process accepted a zero context")
	}
}

func TestBatchExtProcessMatchesDecode(t *testing.T) {
	model, ctx := testBatchExtContext(t, ContextDefaultParams())

	vocab := ModelGetVocab(model)
	nVocab := int(VocabNTokens(vocab))
	tokens := Tokenize(vocab, "This is a test", true, true)
	if len(tokens) == 0 {
		t.Fatal("Tokenize returned zero tokens")
	}

	if rc, _ := Decode(ctx, BatchGetOne(tokens)); rc != 0 {
		t.Fatalf("Decode returned %d", rc)
	}
	logits, _ := GetLogitsIth(ctx, -1, nVocab)
	want := slices.Clone(logits)

	mem, err := GetMemory(ctx)
	if err != nil {
		t.Fatalf("GetMemory failed: %v", err)
	}
	MemoryClear(mem, true)

	batch, err := BatchExtInit(ctx)
	if err != nil {
		t.Fatalf("BatchExtInit failed: %v", err)
	}
	defer BatchExtFree(batch)

	for i, tok := range tokens {
		idx, err := BatchExtAddToken(batch, 0, tok)
		if err != nil {
			t.Fatalf("BatchExtAddToken failed: %v", err)
		}
		if idx != int32(i) {
			t.Fatalf("BatchExtAddToken returned index %d, want %d", idx, i)
		}
		if err := BatchExtSetPos(batch, idx, Pos(i)); err != nil {
			t.Fatalf("BatchExtSetPos failed: %v", err)
		}
	}
	if err := BatchExtSetOutputLogits(batch, int32(len(tokens)-1), true); err != nil {
		t.Fatalf("BatchExtSetOutputLogits failed: %v", err)
	}

	rc, err := Process(ctx, ProcessTypeDecode, batch)
	if err != nil || rc != 0 {
		t.Fatalf("Process returned %d, %v", rc, err)
	}

	got, _ := GetLogitsIth(ctx, -1, nVocab)
	if got == nil {
		t.Fatal("GetLogitsIth returned no logits after Process")
	}
	for i := range want {
		if d := got[i] - want[i]; d > 1e-3 || d < -1e-3 {
			t.Fatalf("logit %d is %v after Process, %v after Decode", i, got[i], want[i])
		}
	}

	if err := BatchExtClear(batch); err != nil {
		t.Fatalf("BatchExtClear failed: %v", err)
	}
	if err := BatchExtSetPos(batch, 0, 0); !errors.Is(err, ErrBatchExtRejected) {
		t.Fatalf("BatchExtSetPos on a cleared batch returned %v, want ErrBatchExtRejected", err)
	}
}

func TestBatchExtErrors(t *testing.T) {
	params := ContextDefaultParams()
	params.NBatch = 2
	params.NUbatch = 2
	model, ctx := testBatchExtContext(t, params)

	batch, err := BatchExtInit(ctx)
	if err != nil {
		t.Fatalf("BatchExtInit failed: %v", err)
	}
	defer BatchExtFree(batch)

	if _, err := BatchExtAddToken(batch, SeqId(NSeqMax(ctx)), 1); !errors.Is(err, ErrBatchExtInvalidSeqID) {
		t.Fatalf("out of range sequence ID returned %v, want ErrBatchExtInvalidSeqID", err)
	}
	if _, err := BatchExtAddToken(batch, 0, Token(VocabNTokens(ModelGetVocab(model)))); !errors.Is(err, ErrBatchExtInvalidToken) {
		t.Fatalf("out of range token returned %v, want ErrBatchExtInvalidToken", err)
	}

	// The rejected token above still took a slot, so clear before filling.
	BatchExtClear(batch)

	nEmbd := int(ModelNEmbdInp(model))
	idx, err := BatchExtAddEmbd(batch, 0, make([]float32, nEmbd), nEmbd)
	if err != nil {
		t.Fatalf("BatchExtAddEmbd failed: %v", err)
	}
	if err := BatchExtSetEmbdToken(batch, idx, make([]float32, nEmbd), nEmbd); !errors.Is(err, ErrBatchExtRejected) {
		t.Fatalf("second embedding for one entry returned %v, want ErrBatchExtRejected", err)
	}
	if err := BatchExtSetEmbdState(batch, idx, make([]float32, nEmbd), nEmbd); !errors.Is(err, ErrBatchExtRejected) {
		t.Fatalf("BatchExtSetEmbdState returned %v, want ErrBatchExtRejected", err)
	}
	if _, err := BatchExtAddEmbd(batch, 0, make([]float32, 3), 2); !errors.Is(err, ErrBatchExtInvalidToken) {
		t.Fatalf("ragged embedding returned %v, want ErrBatchExtInvalidToken", err)
	}
	if err := BatchExtSetPos(batch, idx, 0, 0, 0, 0, 0); !errors.Is(err, ErrBatchExtRejected) {
		t.Fatalf("five positions returned %v, want ErrBatchExtRejected", err)
	}
	if err := BatchExtSetPos(batch, 99, 0); !errors.Is(err, ErrBatchExtRejected) {
		t.Fatalf("out of range index returned %v, want ErrBatchExtRejected", err)
	}
	if err := BatchExtAddSeq(batch, idx, 0); err != nil {
		t.Fatalf("BatchExtAddSeq failed: %v", err)
	}
	if err := BatchExtSetOutputEmbd(batch, idx, true); err != nil {
		t.Fatalf("BatchExtSetOutputEmbd failed: %v", err)
	}

	if _, err := BatchExtAdd(batch, 0); err != nil {
		t.Fatalf("BatchExtAdd failed: %v", err)
	}
	if _, err := BatchExtAddToken(batch, 0, 1); !errors.Is(err, ErrBatchFull) {
		t.Fatalf("third entry in a batch of 2 returned %v, want ErrBatchFull", err)
	}
}
