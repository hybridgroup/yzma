//go:build js && wasm

package llamawasm

import (
	"errors"
	"syscall/js"
	"testing"
)

const fakeContextSource = `
globalThis.__yzmaContext = (function () {
	let pooling = 2;
	const seen = {};

	return {
		module: {
			HEAPU8: new Uint8Array(64),
			_malloc: () => 0,
			_free: () => {},
			_yzma_abi_version: () => 7,
			_yzma_last_error: () => 0,
			_yzma_context_pooling_type: () => pooling,
			_yzma_set_embeddings: (ctx, on) => { seen.embeddings = on; return 0; },
			_yzma_set_causal_attn: (ctx, on) => { seen.causal = on; return 0; },
			_yzma_synchronize: () => { seen.sync = true; return 0; },
		},
		setPooling: (p) => { pooling = p; },
		seen: () => seen,
	};
})();
`

func fakeContext(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeContextSource)
	helper := js.Global().Get("__yzmaContext")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() { mod = previous })

	return helper
}

func TestGetPoolingType(t *testing.T) {
	helper := fakeContext(t)

	if got := GetPoolingType(Context(1)); got != PoolingTypeCLS {
		t.Errorf("GetPoolingType gave %v, want PoolingTypeCLS", got)
	}

	// A pooling type of -1 is a value and not a failure.
	helper.Call("setPooling", -1)
	if got := GetPoolingType(Context(1)); got != PoolingTypeUnspecified {
		t.Errorf("GetPoolingType gave %v, want PoolingTypeUnspecified", got)
	}

	// Only the value of a bad handle means that the call failed.
	helper.Call("setPooling", errBadHandle)
	if got := GetPoolingType(Context(9)); got != PoolingTypeUnspecified {
		t.Errorf("GetPoolingType gave %v, want PoolingTypeUnspecified", got)
	}
}

func TestSetEmbeddingsAndCausalAttn(t *testing.T) {
	helper := fakeContext(t)

	if err := SetEmbeddings(Context(1), true); err != nil {
		t.Fatalf("SetEmbeddings gave %v, want no error", err)
	}
	if err := SetCausalAttn(Context(1), false); err != nil {
		t.Fatalf("SetCausalAttn gave %v, want no error", err)
	}

	seen := helper.Call("seen")
	if seen.Get("embeddings").Int() != 1 {
		t.Error("the shim did not get a true value for the embeddings")
	}
	if seen.Get("causal").Int() != 0 {
		t.Error("the shim did not get a false value for the attention")
	}
}

func TestSynchronize(t *testing.T) {
	helper := fakeContext(t)

	if err := Synchronize(Context(1)); err != nil {
		t.Fatalf("Synchronize gave %v, want no error", err)
	}
	if !helper.Call("seen").Get("sync").Bool() {
		t.Error("Synchronize did not reach the shim")
	}
}

func TestContextCallsOldModule(t *testing.T) {
	// A module before ABI version 7 has none of the calls.
	js.Global().Call("eval", "globalThis.__yzmaOldContext = { HEAPU8: new Uint8Array(8), _malloc: () => 0 };")

	previous := mod
	mod = js.Global().Get("__yzmaOldContext")
	t.Cleanup(func() { mod = previous })

	if got := GetPoolingType(Context(1)); got != PoolingTypeUnspecified {
		t.Errorf("GetPoolingType gave %v, want PoolingTypeUnspecified", got)
	}
	if err := SetEmbeddings(Context(1), true); !errors.Is(err, ErrNoContextFlags) {
		t.Errorf("SetEmbeddings gave %v, want ErrNoContextFlags", err)
	}
	if err := SetCausalAttn(Context(1), true); !errors.Is(err, ErrNoContextFlags) {
		t.Errorf("SetCausalAttn gave %v, want ErrNoContextFlags", err)
	}
	if err := Synchronize(Context(1)); !errors.Is(err, ErrNoContextFlags) {
		t.Errorf("Synchronize gave %v, want ErrNoContextFlags", err)
	}
}
