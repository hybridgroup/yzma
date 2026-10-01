//go:build js && wasm

package llamawasm

import (
	"strings"
	"syscall/js"
	"testing"
)

// fakeRejecting installs a module whose calls all return a rejected promise,
// as a WebGPU module does when the GPU device is lost.
func fakeRejecting(t *testing.T) {
	t.Helper()

	js.Global().Call("eval", `globalThis.__yzmaRejecting = new Proxy({ HEAPU8: new Uint8Array(64) }, {
		get: (target, name) => name in target ? target[name] : () => Promise.reject(new Error("device lost")),
	});`)

	previous := mod
	mod = js.Global().Get("__yzmaRejecting")
	t.Cleanup(func() {
		mod = previous
		rejected = nil
	})
}

func TestRejectedPromiseDoesNotPanic(t *testing.T) {
	fakeRejecting(t)

	if got := call("_yzma_model_n_embd", 1); got != errBadHandle {
		t.Errorf("call gave %d, want errBadHandle", got)
	}
	if _, err := callErr("_yzma_decode", 1); err == nil || !strings.Contains(err.Error(), "device lost") {
		t.Errorf("callErr gave %v, want an error with the rejection reason", err)
	}
	if rejected != nil {
		t.Error("shimError did not clear the rejection")
	}

	if got := SamplerSample(Sampler(1), Context(1), -1); got != TokenNull {
		t.Errorf("SamplerSample gave %d, want TokenNull", got)
	}
	if got := SamplerChainInit(SamplerChainDefaultParams()); got != 0 {
		t.Errorf("SamplerChainInit gave %d, want 0", got)
	}
	if got := ModelGetVocab(Model(1)); got != 0 {
		t.Errorf("ModelGetVocab gave %d, want 0", got)
	}
	if got := SamplerGetSeed(Sampler(1)); got != 0xFFFFFFFF {
		t.Errorf("SamplerGetSeed gave %d, want 0xFFFFFFFF", got)
	}
	if got := TimeUs(); got != 0 {
		t.Errorf("TimeUs gave %d, want 0", got)
	}
	if got := ModelSize(Model(1)); got != 0 {
		t.Errorf("ModelSize gave %d, want 0", got)
	}
	if got := StateGetSize(Context(1)); got != 0 {
		t.Errorf("StateGetSize gave %d, want 0", got)
	}
	if got := VocabGetScore(Vocab(1), 1); got != 0 {
		t.Errorf("VocabGetScore gave %v, want 0", got)
	}
	callVoid("_yzma_backend_free")
}
