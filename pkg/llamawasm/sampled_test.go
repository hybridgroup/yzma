//go:build js && wasm

package llamawasm

import (
	"errors"
	"syscall/js"
	"testing"
)

// The fake module answers the calls of the backend sampler. It has a real
// heap, thus the test covers the way that the package reads the memory.
const fakeSamplerSource = `
globalThis.__yzmaSampler = (function () {
	const heap = new Uint8Array(1 << 20);
	let next = 8;

	let token = 42;
	let taken = 1;
	const counts = { probs: 3, logits: 3, candidates: 3 };

	function floats(out, n, first) {
		const values = new Float32Array(heap.buffer, out, n);
		for (let i = 0; i < n; i++) { values[i] = first + i; }
		return n;
	}

	return {
		module: {
			HEAPU8: heap,
			_malloc: (n) => { const p = next; next += n + (8 - (n % 8)); return p; },
			_free: () => {},
			_yzma_abi_version: () => 7,
			_yzma_last_error: () => 0,
			_yzma_set_sampler: () => taken,
			_yzma_get_sampled_token_ith: () => token,
			_yzma_get_sampled_probs_count_ith: () => counts.probs,
			_yzma_get_sampled_logits_count_ith: () => counts.logits,
			_yzma_get_sampled_candidates_count_ith: () => counts.candidates,
			_yzma_get_sampled_probs_ith: (ctx, i, out, n) => floats(out, n, 0.5),
			_yzma_get_sampled_logits_ith: (ctx, i, out, n) => floats(out, n, 10),
			_yzma_get_sampled_candidates_ith: (ctx, i, out, n) => {
				const tokens = new Int32Array(heap.buffer, out, n);
				for (let j = 0; j < n; j++) { tokens[j] = 100 + j; }
				return n;
			},
		},
		setToken: (t) => { token = t; },
		setTaken: (t) => { taken = t; },
	};
})();
`

func fakeSampler(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeSamplerSource)
	helper := js.Global().Get("__yzmaSampler")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() {
		mod = previous
		outScratch.ptr, outScratch.size = 0, 0
	})

	return helper
}

func TestSetSampler(t *testing.T) {
	helper := fakeSampler(t)

	taken, err := SetSampler(Context(1), 0, Sampler(2))
	if err != nil || !taken {
		t.Errorf("SetSampler gave %v and %v, want true and no error", taken, err)
	}

	// A backend that cannot take the sampler gives 0, which is not an error.
	helper.Call("setTaken", 0)
	taken, err = SetSampler(Context(1), 0, Sampler(2))
	if err != nil || taken {
		t.Errorf("SetSampler gave %v and %v, want false and no error", taken, err)
	}
}

func TestGetSampledTokenIth(t *testing.T) {
	helper := fakeSampler(t)

	got, err := GetSampledTokenIth(Context(1), 0)
	if err != nil || got != Token(42) {
		t.Errorf("GetSampledTokenIth gave %v and %v, want 42 and no error", got, err)
	}

	// A backend that sampled no token gives -1, which is not an error.
	helper.Call("setToken", -1)
	got, err = GetSampledTokenIth(Context(1), 0)
	if err != nil || got != TokenNull {
		t.Errorf("GetSampledTokenIth gave %v and %v, want TokenNull and no error", got, err)
	}

	// Only the value of a bad handle is an error.
	helper.Call("setToken", errBadHandle)
	if _, err = GetSampledTokenIth(Context(9), 0); err == nil {
		t.Error("GetSampledTokenIth gave no error, want the error of a bad handle")
	}
}

func TestSampledCounts(t *testing.T) {
	fakeSampler(t)

	for _, test := range []struct {
		name string
		call func(Context, int32) (int32, error)
	}{
		{"GetSampledProbsCountIth", GetSampledProbsCountIth},
		{"GetSampledLogitsCountIth", GetSampledLogitsCountIth},
		{"GetSampledCandidatesCountIth", GetSampledCandidatesCountIth},
	} {
		got, err := test.call(Context(1), 0)
		if err != nil || got != 3 {
			t.Errorf("%s gave %v and %v, want 3 and no error", test.name, got, err)
		}
	}
}

func TestGetSampledProbsIth(t *testing.T) {
	fakeSampler(t)

	got, err := GetSampledProbsIth(Context(1), 0, 3)
	if err != nil {
		t.Fatalf("GetSampledProbsIth gave %v, want no error", err)
	}
	want := []float32{0.5, 1.5, 2.5}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("value %d is %v, want %v", i, got[i], want[i])
		}
	}
}

func TestGetSampledLogitsIth(t *testing.T) {
	fakeSampler(t)

	got, err := GetSampledLogitsIth(Context(1), 0, 3)
	if err != nil {
		t.Fatalf("GetSampledLogitsIth gave %v, want no error", err)
	}
	if len(got) != 3 || got[0] != 10 || got[2] != 12 {
		t.Errorf("GetSampledLogitsIth gave %v, want [10 11 12]", got)
	}
}

func TestGetSampledCandidatesIth(t *testing.T) {
	fakeSampler(t)

	got, err := GetSampledCandidatesIth(Context(1), 0, 3)
	if err != nil {
		t.Fatalf("GetSampledCandidatesIth gave %v, want no error", err)
	}
	want := []Token{100, 101, 102}
	if len(got) != len(want) {
		t.Fatalf("GetSampledCandidatesIth gave %d tokens, want 3", len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("token %d is %v, want %v", i, got[i], want[i])
		}
	}
}

func TestSampledNoCandidates(t *testing.T) {
	fakeSampler(t)

	got, err := GetSampledCandidatesIth(Context(1), 0, 0)
	if err != nil || got != nil {
		t.Errorf("GetSampledCandidatesIth gave %v and %v, want nil and no error", got, err)
	}
}

func TestSampledOldModule(t *testing.T) {
	// A module before ABI version 7 has none of the calls.
	js.Global().Call("eval", "globalThis.__yzmaOldSampler = { HEAPU8: new Uint8Array(8), _malloc: () => 0 };")

	previous := mod
	mod = js.Global().Get("__yzmaOldSampler")
	t.Cleanup(func() { mod = previous })

	if _, err := SetSampler(Context(1), 0, Sampler(1)); !errors.Is(err, ErrNoBackendSampling) {
		t.Errorf("SetSampler gave %v, want ErrNoBackendSampling", err)
	}
	if _, err := GetSampledTokenIth(Context(1), 0); !errors.Is(err, ErrNoBackendSampling) {
		t.Errorf("GetSampledTokenIth gave %v, want ErrNoBackendSampling", err)
	}
	if _, err := GetSampledProbsCountIth(Context(1), 0); !errors.Is(err, ErrNoBackendSampling) {
		t.Errorf("GetSampledProbsCountIth gave %v, want ErrNoBackendSampling", err)
	}
	if _, err := GetSampledProbsIth(Context(1), 0, 3); !errors.Is(err, ErrNoBackendSampling) {
		t.Errorf("GetSampledProbsIth gave %v, want ErrNoBackendSampling", err)
	}
	if _, err := GetSampledCandidatesIth(Context(1), 0, 3); !errors.Is(err, ErrNoBackendSampling) {
		t.Errorf("GetSampledCandidatesIth gave %v, want ErrNoBackendSampling", err)
	}
}
