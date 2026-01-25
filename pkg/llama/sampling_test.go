package llama

import (
	"testing"
	"unsafe"
)

func TestSamplerChainDefaultParams(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := SamplerChainDefaultParams()
	if params == (SamplerChainParams{}) {
		t.Fatal("SamplerChainDefaultParams returned empty parameters")
	}
}

func TestSamplerChainInit(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	params := SamplerChainDefaultParams()
	chain := SamplerChainInit(params)
	if chain == (Sampler(0)) {
		t.Fatal("SamplerChainInit failed to initialize a sampler chain")
	}

	SamplerFree(chain)
}

func TestSamplerInitGreedy(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitGreedy()
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitGreedy failed to initialize a greedy sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitDist(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitDist(12345)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitDist failed to initialize a distribution sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitTopK(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitTopK(40)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitTopK failed to initialize a top-k sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitTopP(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitTopP(0.95, 0)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitTopP failed to initialize a top-p sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitMinP(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitMinP(0.05, 0)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitMinP failed to initialize a min-p sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitTypical(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitTypical(1.0, 0)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitTypical failed to initialize a typical sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitPenalties(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitPenalties(64, 1.0, 0.0, 0.0)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitPenalties failed to initialize a penalties sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitDry(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	modelFile := testModelFileName(t)
	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	sampler := SamplerInitDry(vocab, 1024, 0.2, 0.2, 1.0, 2, []string{"the", "and", "of"})
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitDry failed to initialize a dry sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitDry2(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	modelFile := testModelFileName(t)
	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	sampler := SamplerInitDry(vocab, 1024, 0.2, 0.2, 1.0, 2, nil)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitDry failed to initialize a dry sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitLogitBias(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	biases := []LogitBias{
		{Token: 10, Bias: -1.0},
		{Token: 20, Bias: 2.0},
	}
	sampler := SamplerInitLogitBias(100, int32(len(biases)), unsafe.SliceData(biases))
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitLogitBias failed to initialize a logit bias sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitTopNSigma(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitTopNSigma(1.0)
	if sampler == 0 {
		t.Fatal("SamplerInitTopNSigma failed to initialize a top-n-sigma sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitXTC(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitXTC(0.5, 1.0, 0, 42)
	if sampler == 0 {
		t.Fatal("SamplerInitXTC failed to initialize an XTC sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitTempExt(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitTempExt(0.8, 0, 1.0)
	if sampler == 0 {
		t.Fatal("SamplerInitTempExt failed to initialize a temp-ext sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerInitGrammar(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	modelFile := testModelFileName(t)
	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	vocab := ModelGetVocab(model)

	// Use a simple grammar string and root for testing
	grammar := "root ::= \"hello\""
	root := "root"

	sampler := SamplerInitGrammar(vocab, grammar, root)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitGrammar failed to initialize a grammar sampler")
	}

	SamplerFree(sampler)
}

func TestSamplerSample(t *testing.T) {
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

	// Use a simple sampler (e.g., greedy)
	sampler := SamplerInitGreedy()
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitGreedy failed to initialize a greedy sampler")
	}
	defer SamplerFree(sampler)

	// Tokenize a prompt and decode to produce logits
	prompt := "Hello world"
	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, prompt, true, true)
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	// Sample a token (index 0)
	token := SamplerSample(sampler, ctx, -1)
	if token == TokenNull {
		t.Fatal("SamplerSample returned TokenNull")
	}
	t.Logf("SamplerSample returned token: %d", token)
}

func TestSamplerAccept(t *testing.T) {
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

	// Use a simple sampler (e.g., greedy)
	sampler := SamplerInitGreedy()
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitGreedy failed to initialize a greedy sampler")
	}
	defer SamplerFree(sampler)

	// Tokenize a prompt and decode to produce logits
	prompt := "Hello world"
	vocab := ModelGetVocab(model)
	tokens := Tokenize(vocab, prompt, true, true)
	batch := BatchGetOne(tokens)
	Decode(ctx, batch)

	// Sample a token (index 0)
	token := SamplerSample(sampler, ctx, -1)
	if token == TokenNull {
		t.Fatal("SamplerSample returned TokenNull")
	}

	// Accept the sampled token (should not panic or error)
	SamplerAccept(sampler, token)
	t.Logf("SamplerAccept succeeded for token: %d", token)
}

func TestNewSampler(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	modelFile := testModelFileName(t)
	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	samplers := []SamplerType{
		SamplerTypeTopK,
		SamplerTypeTopP,
	}

	sampler := NewSampler(model, samplers, DefaultSamplerParams())
	if sampler == (Sampler(0)) {
		t.Fatal("NewSampler failed to create a new sampler chain")
	}

	SamplerFree(sampler)
}

func TestSamplerReset(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	sampler := SamplerInitDist(12345)
	if sampler == (Sampler(0)) {
		t.Fatal("SamplerInitDist failed to initialize a distribution sampler")
	}

	// Reset the sampler
	SamplerReset(sampler)

	SamplerFree(sampler)
}
