//go:build js && wasm

package llamawasm

import (
	"errors"
	"syscall/js"
	"testing"
)

// The fake module has a real heap, thus the test covers the way that the
// package writes and reads the memory of the module.
const fakeModuleSource = `
globalThis.__yzmaModule = (function () {
	const heap = new Uint8Array(1 << 20);
	let next = 8;
	const calls = [];

	// values holds what the shim writes. A test sets it.
	let values = [];
	let result = null;

	function put(out, n) {
		const floats = new Float32Array(heap.buffer, out, n);
		for (let i = 0; i < n; i++) {
			floats[i] = values[i % values.length];
		}
		return result === null ? n : result;
	}

	return {
		module: {
			HEAPU8: heap,
			_malloc: (n) => { const p = next; next += n + (8 - (n % 8)); return p; },
			_free: () => {},
			_yzma_abi_version: () => 7,
			_yzma_last_error: () => 0,
			_yzma_get_logits_ith: (ctx, i, out, n) => {
				calls.push(["ith", ctx, i, n]);
				return put(out, n);
			},
			_yzma_get_logits: (ctx, out, n) => {
				calls.push(["all", ctx, n]);
				return put(out, n);
			},
			_yzma_get_embeddings_ith: (ctx, i, out, n) => {
				calls.push(["embd ith", ctx, i, n]);
				return put(out, n);
			},
			_yzma_get_embeddings: (ctx, out, n) => {
				calls.push(["embd all", ctx, n]);
				return put(out, n);
			},
		},
		setValues: (v) => { values = v; },
		setResult: (r) => { result = r; },
		calls: () => calls,
	};
})();
`

// fakeModule puts a module of the test in place and gives the helper back.
func fakeModule(t *testing.T, values []any) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeModuleSource)
	helper := js.Global().Get("__yzmaModule")
	helper.Call("setValues", values)

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() {
		mod = previous
		outScratch.ptr, outScratch.size = 0, 0
	})

	return helper
}

func TestGetLogitsIth(t *testing.T) {
	helper := fakeModule(t, []any{1.5, -2.25, 0.0, 4.75})

	got, err := GetLogitsIth(Context(3), -1, 4)
	if err != nil {
		t.Fatalf("GetLogitsIth gave %v, want no error", err)
	}

	want := []float32{1.5, -2.25, 0, 4.75}
	if len(got) != len(want) {
		t.Fatalf("GetLogitsIth gave %d values, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("value %d is %v, want %v", i, got[i], want[i])
		}
	}

	call := helper.Call("calls").Index(0)
	if ctx, i := call.Index(1).Int(), call.Index(2).Int(); ctx != 3 || i != -1 {
		t.Errorf("the shim got ctx %d and i %d, want 3 and -1", ctx, i)
	}
}

func TestGetLogits(t *testing.T) {
	helper := fakeModule(t, []any{0.5})

	// Two tokens of a vocabulary of three give six values.
	got, err := GetLogits(Context(1), 2, 3)
	if err != nil {
		t.Fatalf("GetLogits gave %v, want no error", err)
	}
	if len(got) != 6 {
		t.Fatalf("GetLogits gave %d values, want 6", len(got))
	}

	if n := helper.Call("calls").Index(0).Index(2).Int(); n != 6 {
		t.Errorf("the shim got a count of %d, want 6", n)
	}
}

func TestGetLogitsNoTokens(t *testing.T) {
	fakeModule(t, []any{1.0})

	got, err := GetLogits(Context(1), 0, 100)
	if err != nil || got != nil {
		t.Errorf("GetLogits gave %v and %v, want nil and no error", got, err)
	}
}

func TestGetLogitsShimFails(t *testing.T) {
	helper := fakeModule(t, []any{1.0})
	helper.Call("setResult", errHandle)

	if _, err := GetLogitsIth(Context(9), 0, 4); err == nil {
		t.Error("GetLogitsIth gave no error, want the error of the shim")
	}
}

func TestGetEmbeddingsIth(t *testing.T) {
	helper := fakeModule(t, []any{0.25, -0.5})

	got, err := GetEmbeddingsIth(Context(2), -1, 2)
	if err != nil {
		t.Fatalf("GetEmbeddingsIth gave %v, want no error", err)
	}
	if len(got) != 2 || got[0] != 0.25 || got[1] != -0.5 {
		t.Errorf("GetEmbeddingsIth gave %v, want [0.25 -0.5]", got)
	}

	call := helper.Call("calls").Index(0)
	if ctx, i := call.Index(1).Int(), call.Index(2).Int(); ctx != 2 || i != -1 {
		t.Errorf("the shim got ctx %d and i %d, want 2 and -1", ctx, i)
	}
}

func TestGetEmbeddings(t *testing.T) {
	helper := fakeModule(t, []any{1.0})

	// Three outputs of an embedding of four values give twelve values.
	got, err := GetEmbeddings(Context(1), 3, 4)
	if err != nil {
		t.Fatalf("GetEmbeddings gave %v, want no error", err)
	}
	if len(got) != 12 {
		t.Fatalf("GetEmbeddings gave %d values, want 12", len(got))
	}

	if n := helper.Call("calls").Index(0).Index(2).Int(); n != 12 {
		t.Errorf("the shim got a count of %d, want 12", n)
	}
}

func TestGetEmbeddingsNoOutputs(t *testing.T) {
	fakeModule(t, []any{1.0})

	got, err := GetEmbeddings(Context(1), 0, 100)
	if err != nil || got != nil {
		t.Errorf("GetEmbeddings gave %v and %v, want nil and no error", got, err)
	}
}

func TestOutputsOldModule(t *testing.T) {
	// A module before ABI version 7 has neither call.
	js.Global().Call("eval", "globalThis.__yzmaOld = { HEAPU8: new Uint8Array(8), _malloc: () => 0 };")

	previous := mod
	mod = js.Global().Get("__yzmaOld")
	t.Cleanup(func() { mod = previous })

	if _, err := GetLogitsIth(Context(1), -1, 4); !errors.Is(err, ErrNoOutputs) {
		t.Errorf("GetLogitsIth gave %v, want ErrNoOutputs", err)
	}
	if _, err := GetLogits(Context(1), 1, 4); !errors.Is(err, ErrNoOutputs) {
		t.Errorf("GetLogits gave %v, want ErrNoOutputs", err)
	}
	if _, err := GetEmbeddingsIth(Context(1), -1, 4); !errors.Is(err, ErrNoOutputs) {
		t.Errorf("GetEmbeddingsIth gave %v, want ErrNoOutputs", err)
	}
	if _, err := GetEmbeddings(Context(1), 1, 4); !errors.Is(err, ErrNoOutputs) {
		t.Errorf("GetEmbeddings gave %v, want ErrNoOutputs", err)
	}
}
