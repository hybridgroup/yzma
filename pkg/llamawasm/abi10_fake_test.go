//go:build js && wasm

package llamawasm

import (
	"errors"
	"syscall/js"
	"testing"
)

// The fake module of ABI version 10 records the arguments of the new context call.
const fakeABI10Source = `
globalThis.__yzmaABI10 = (function () {
	let last = null;
	return {
		module: {
			HEAPU8: new Uint8Array(64),
			_malloc: () => 8,
			_free: () => {},
			_yzma_abi_version: () => 10,
			_yzma_last_error: () => 0,
			_yzma_context_new_seq: () => 3,
			_yzma_context_new_kv: (...args) => { last = args; return 5; },
		},
		last: () => last,
	};
})();
`

func fakeABI10(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeABI10Source)
	helper := js.Global().Get("__yzmaABI10")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() { mod = previous })

	return helper
}

func TestContextKVUnified(t *testing.T) {
	helper := fakeABI10(t)

	params := ContextDefaultParams()
	params.NSeqMax = 17
	params.KVUnified = 1
	params.NOutputsMax = 256
	ctx, err := InitFromModel(Model(1), params)
	if err != nil {
		t.Fatalf("InitFromModel gave %v", err)
	}
	if ctx != 5 {
		t.Errorf("InitFromModel gave context %d, want 5", ctx)
	}

	last := helper.Call("last")
	if n := last.Length(); n != 11 {
		t.Fatalf("the shim got %d arguments, want 11", n)
	}
	for i, want := range map[int]int{7: 17, 9: 1, 10: 256} {
		if got := last.Index(i).Int(); got != want {
			t.Errorf("argument %d is %d, want %d", i, got, want)
		}
	}
}

func TestContextKVUnifiedOld(t *testing.T) {
	fakeABI9(t)

	params := ContextDefaultParams()
	params.KVUnified = 1
	if _, err := InitFromModel(Model(1), params); !errors.Is(err, ErrNoKVUnified) {
		t.Errorf("InitFromModel gave %v, want ErrNoKVUnified", err)
	}

	// Without KVUnified the old call still makes the context.
	params.KVUnified = 0
	if _, err := InitFromModel(Model(1), params); err != nil {
		t.Errorf("InitFromModel gave %v", err)
	}
}
