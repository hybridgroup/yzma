package llama

import (
	"testing"
)

// Note: These tests don't verify the actual functionality since we don't have a loaded model

// TestPerfContextDataString tests the String() method of PerfContextData
func TestPerfContextDataString(t *testing.T) {
	data := PerfContextData{
		TStartMs:      100.5,
		TLoadMs:       200.3,
		TPromptEvalMs: 50.2,
		TEvalMs:       300.1,
		NPEval:        10,
		NEval:         20,
		NReused:       5,
	}

	expected := "PerfContextData{Start: 100.50ms, Load: 200.30ms, Prompt Eval: 50.20ms, Eval: 300.10ms, Prompt Tokens: 10, Gen Tokens: 20, Reused: 5}"
	actual := data.String()

	if actual != expected {
		t.Errorf("PerfContextData.String() = %s; expected %s", actual, expected)
	}
}

// TestPerfSamplerDataString tests the String() method of PerfSamplerData
func TestPerfSamplerDataString(t *testing.T) {
	data := PerfSamplerData{
		TSampleMs: 25.5,
		NSample:   15,
	}

	expected := "PerfSamplerData{Sample Time: 25.50ms, Samples: 15}"
	actual := data.String()

	if actual != expected {
		t.Errorf("PerfSamplerData.String() = %s; expected %s", actual, expected)
	}
}

// TestPerfFunctionsExist tests that the performance functions exist and can be called
func TestPerfFunctionsExist(t *testing.T) {
	// Test that we can call the functions without panicking
	// These should not panic even with invalid contexts/samplers
	ctx := Context(0)
	chain := Sampler(0)

	// These should not panic
	data := PerfContext(ctx)
	if data.TStartMs != 0 || data.TLoadMs != 0 || data.TPromptEvalMs != 0 || data.TEvalMs != 0 {
		t.Errorf("Expected zero values for invalid context")
	}

	PerfContextPrint(ctx)

	samplerData := PerfSampler(chain)
	if samplerData.TSampleMs != 0 || samplerData.NSample != 0 {
		t.Errorf("Expected zero values for invalid sampler")
	}

	PerfSamplerPrint(chain)
	PerfSamplerReset(chain)
	MemoryBreakdownPrint(ctx)
}
