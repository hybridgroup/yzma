//go:build js && wasm

package llamawasm

import "syscall/js"

// A Go function cannot be the abort callback of llama.cpp in a browser. So the
// shim holds one flag, and a context with SetAbortCallback stops its compute
// when the flag is set.
//
// The CPU backend reads the flag after each node of the graph. The WebGPU
// backend does not read it, so on the GPU a batch runs to the end.
//
// Decode is synchronous, so the worker cannot run Go code while a batch runs.
// In the build with threads the module memory is a SharedArrayBuffer, so the
// page can set the flag itself. Post [AbortFlagView] to the page, and the page
// calls Atomics.store(view, 0, 1). Decode then returns 2.

// AbortSupported reports whether the module has the abort flag, which arrived
// in ABI version 13.
func AbortSupported() bool {
	return has("_yzma_set_abort_callback")
}

// SetAbortCallback connects the abort flag to a context when on is true, and
// disconnects it when on is false.
func SetAbortCallback(ctx Context, on bool) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	if !AbortSupported() {
		return ErrNoAbort
	}
	_, err := callErr("_yzma_set_abort_callback", int(ctx), boolToInt(on))
	return err
}

// Abort sets the abort flag. A connected context stops at the next node it
// computes on the CPU, and the call that runs it returns 2.
func Abort() {
	if AbortSupported() {
		callVoid("_yzma_abort_set", 1)
	}
}

// AbortReset clears the abort flag. Clear it before the next request,
// because the flag stays set after a stop.
func AbortReset() {
	if AbortSupported() {
		callVoid("_yzma_abort_set", 0)
	}
}

// Aborted reports whether the abort flag is set.
func Aborted() bool {
	return AbortSupported() && call("_yzma_abort_get") != 0
}

// AbortFlagView returns an Int32Array of length 1 over the abort flag. Only
// the build with threads shares its memory, so it returns undefined for the
// other builds.
func AbortFlagView() js.Value {
	if !AbortSupported() || !threaded {
		return js.Undefined()
	}
	ptr := call("_yzma_abort_flag")
	if ptr <= 0 {
		return js.Undefined()
	}
	return js.Global().Get("Int32Array").New(heapU8().Get("buffer"), int(ptr), 1)
}
