package llama

import (
	"testing"
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
