package llama

import "testing"

func TestDraftGenerateGuards(t *testing.T) {
	nPast := Pos(7)
	outTokens := []Token{111, 222}
	outDists := [][]DraftCandidate{{{Tok: 3, Prob: 0.5}}}

	tests := []struct {
		name    string
		ctx     Context
		sampler Sampler
		nDraft  int
	}{
		{"nil_context", 0, Sampler(1), 2},
		{"nil_sampler", Context(1), 0, 2},
		{"zero_nDraft", Context(1), Sampler(1), 0},
		{"negative_nDraft", Context(1), Sampler(1), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch := Batch{}
			drafted, finalPast := DraftGenerate(
				tt.ctx, &batch, Vocab(1), tt.sampler,
				Token(1), nPast, []SeqId{0}, tt.nDraft,
				true, outTokens, outDists,
			)
			if drafted != 0 {
				t.Fatalf("drafted = %d, want 0", drafted)
			}
			if finalPast != nPast {
				t.Fatalf("finalPast = %d, want %d", finalPast, nPast)
			}
		})
	}

	if outTokens[0] != 111 || outTokens[1] != 222 {
		t.Fatal("outTokens was unexpectedly modified")
	}
	if len(outDists[0]) != 1 || outDists[0][0].Tok != 3 {
		t.Fatal("outDists was unexpectedly modified")
	}
}

func TestDraftGenerateGreedy(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	modelFile := testModelFileName(t)
	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, "Hello world", true, true)
	if len(tokens) < 2 {
		t.Skip("prompt tokenized to fewer than 2 tokens")
	}

	prefix := tokens[:len(tokens)-1]
	Decode(ctx, BatchGetOne(prefix))

	batch := BatchInit(1, 0, 1)
	defer BatchFree(batch)

	lastToken := tokens[len(tokens)-1]
	nPast := Pos(len(prefix))
	sampler := SamplerInitGreedy()
	if sampler == 0 {
		t.Fatal("SamplerInitGreedy failed")
	}
	defer SamplerFree(sampler)

	const nDraft = 2
	outTokens := make([]Token, nDraft)

	drafted, finalPast := DraftGenerate(
		ctx, &batch, vocab, sampler,
		lastToken, nPast, []SeqId{0}, nDraft,
		true, outTokens, nil,
	)

	if drafted < 0 || drafted > nDraft {
		t.Fatalf("drafted = %d, want 0..%d", drafted, nDraft)
	}

	// finalPast advances by drafted (all non-EOG) or drafted+1 (EOG terminated the loop).
	delta := int(finalPast - nPast)
	if delta < drafted || delta > drafted+1 {
		t.Fatalf("finalPast advanced by %d with drafted=%d", delta, drafted)
	}

	for i := range drafted {
		if VocabIsEOG(vocab, outTokens[i]) {
			t.Fatalf("outTokens[%d] is EOG, should not be in output", i)
		}
	}
}

func TestDraftGenerateNonGreedy(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	modelFile := testModelFileName(t)
	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, "Hello world", true, true)
	if len(tokens) < 2 {
		t.Skip("prompt tokenized to fewer than 2 tokens")
	}

	prefix := tokens[:len(tokens)-1]
	Decode(ctx, BatchGetOne(prefix))

	batch := BatchInit(1, 0, 1)
	defer BatchFree(batch)

	lastToken := tokens[len(tokens)-1]
	nPast := Pos(len(prefix))

	params := SamplerChainDefaultParams()
	chain := SamplerChainInit(params)
	if chain == 0 {
		t.Fatal("SamplerChainInit failed")
	}
	defer SamplerFree(chain)
	SamplerChainAdd(chain, SamplerInitDist(42))

	const nDraft = 2
	outTokens := make([]Token, nDraft)
	outDists := make([][]DraftCandidate, nDraft)
	for i := range outDists {
		outDists[i] = make([]DraftCandidate, 0, 16)
	}

	drafted, finalPast := DraftGenerate(
		ctx, &batch, vocab, chain,
		lastToken, nPast, []SeqId{0}, nDraft,
		false, outTokens, outDists,
	)

	if drafted < 0 || drafted > nDraft {
		t.Fatalf("drafted = %d, want 0..%d", drafted, nDraft)
	}

	delta := int(finalPast - nPast)
	if delta < drafted || delta > drafted+1 {
		t.Fatalf("finalPast advanced by %d with drafted=%d", delta, drafted)
	}

	for i := range drafted {
		if VocabIsEOG(vocab, outTokens[i]) {
			t.Fatalf("outTokens[%d] is EOG, should not be in output", i)
		}
	}

	t.Logf("non-greedy drafted %d tokens", drafted)
	for i := range drafted {
		t.Logf("  dist[%d]: %d candidates", i, len(outDists[i]))
	}
}
