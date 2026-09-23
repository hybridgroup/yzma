//go:build js && wasm

package llamawasm

import "fmt"

// PerfContextData holds the performance counters of a context. The times stay
// 0 unless the context has NoPerf 0 in its [ContextParams].
type PerfContextData struct {
	TStartMs      float64 // absolute start time
	TLoadMs       float64 // time needed for loading the model
	TPromptEvalMs float64 // time needed for processing the prompt
	TEvalMs       float64 // time needed for generating tokens

	NPEval  int32 // number of prompt tokens
	NEval   int32 // number of generated tokens
	NReused int32 // number of times a ggml compute graph had been reused
}

// String gives the counters as text.
func (p PerfContextData) String() string {
	return fmt.Sprintf("PerfContextData{Start: %.2fms, Load: %.2fms, Prompt Eval: %.2fms, Eval: %.2fms, Prompt Tokens: %d, Gen Tokens: %d, Reused: %d}",
		p.TStartMs, p.TLoadMs, p.TPromptEvalMs, p.TEvalMs, p.NPEval, p.NEval, p.NReused)
}

// PerfSamplerData holds the performance counters of a chain of samplers.
type PerfSamplerData struct {
	TSampleMs float64 // time needed for sampling in ms

	NSample int32 // number of sampled tokens
}

// String gives the counters as text.
func (p PerfSamplerData) String() string {
	return fmt.Sprintf("PerfSamplerData{Sample Time: %.2fms, Samples: %d}", p.TSampleMs, p.NSample)
}

// PerfContext gives the performance counters of the context. A module before
// ABI version 9 gives zero values.
func PerfContext(ctx Context) PerfContextData {
	v, err := perfValues("_yzma_perf_context", 7, int(ctx))
	if err != nil {
		return PerfContextData{}
	}
	return PerfContextData{
		TStartMs:      v[0],
		TLoadMs:       v[1],
		TPromptEvalMs: v[2],
		TEvalMs:       v[3],
		NPEval:        int32(v[4]),
		NEval:         int32(v[5]),
		NReused:       int32(v[6]),
	}
}

// PerfContextReset sets the performance counters of the context to zero.
func PerfContextReset(ctx Context) error {
	return perfReset("_yzma_perf_context_reset", int(ctx))
}

// PerfSampler gives the performance counters of a chain of samplers. A sampler
// that is not a chain gives zero values.
func PerfSampler(chain Sampler) PerfSamplerData {
	v, err := perfValues("_yzma_perf_sampler", 2, int(chain))
	if err != nil {
		return PerfSamplerData{}
	}
	return PerfSamplerData{TSampleMs: v[0], NSample: int32(v[1])}
}

// PerfSamplerReset sets the performance counters of a chain of samplers to
// zero.
func PerfSamplerReset(chain Sampler) {
	perfReset("_yzma_perf_sampler_reset", int(chain))
}

// perfValues runs a call that writes n doubles and reads them.
func perfValues(name string, n int, handle int) ([]float64, error) {
	if !Loaded() {
		return nil, ErrNotLoaded
	}
	if !has(name) {
		return nil, ErrNoPerf
	}

	ptr, err := outScratch.reserve(n * 8)
	if err != nil {
		return nil, err
	}
	if _, err := callErr(name, handle, ptr); err != nil {
		return nil, err
	}
	return readFloat64s(ptr, n), nil
}

func perfReset(name string, handle int) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	if !has(name) {
		return ErrNoPerf
	}
	_, err := callErr(name, handle)
	return err
}
