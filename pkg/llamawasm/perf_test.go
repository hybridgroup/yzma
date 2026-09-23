//go:build js && wasm

package llamawasm

import (
	"errors"
	"testing"
)

func TestPerfContext(t *testing.T) {
	fakeABI9(t)

	got := PerfContext(Context(1))
	want := PerfContextData{TStartMs: 10.5, TLoadMs: 20.25, TPromptEvalMs: 30, TEvalMs: 40, NPEval: 7, NEval: 3, NReused: 2}
	if got != want {
		t.Errorf("PerfContext gave %v, want %v", got, want)
	}
	if got := PerfContext(Context(2)); got != (PerfContextData{}) {
		t.Errorf("PerfContext of a bad handle gave %v, want zero values", got)
	}

	if err := PerfContextReset(Context(1)); err != nil {
		t.Errorf("PerfContextReset gave %v, want no error", err)
	}
	if err := PerfContextReset(Context(2)); err == nil {
		t.Error("PerfContextReset of a bad handle gave no error")
	}
}

func TestPerfSampler(t *testing.T) {
	fakeABI9(t)

	if got := PerfSampler(Sampler(1)); got != (PerfSamplerData{TSampleMs: 1.5, NSample: 9}) {
		t.Errorf("PerfSampler gave %v", got)
	}
	// The shim refuses a sampler that is not a chain.
	if got := PerfSampler(Sampler(3)); got != (PerfSamplerData{}) {
		t.Errorf("PerfSampler of a sampler that is not a chain gave %v", got)
	}
	PerfSamplerReset(Sampler(1))
}

func TestPerfOldModule(t *testing.T) {
	fakeOld(t)

	if got := PerfContext(Context(1)); got != (PerfContextData{}) {
		t.Errorf("PerfContext gave %v, want zero values", got)
	}
	if err := PerfContextReset(Context(1)); !errors.Is(err, ErrNoPerf) {
		t.Errorf("PerfContextReset gave %v, want ErrNoPerf", err)
	}
}

func TestContextNoPerf(t *testing.T) {
	helper := fakeABI9(t)

	params := ContextDefaultParams()
	if params.NoPerf != 1 {
		t.Errorf("the default NoPerf is %d, want 1 as in llama.cpp", params.NoPerf)
	}

	params.NoPerf = 0
	if _, err := InitFromModel(Model(1), params); err != nil {
		t.Fatalf("InitFromModel gave %v", err)
	}
	if got := helper.Call("lastNoPerf").Int(); got != 0 {
		t.Errorf("the shim got NoPerf %d, want 0", got)
	}
}
