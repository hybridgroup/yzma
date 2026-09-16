//go:build js && wasm

package llamawasm

// NThreads gives the number of threads that the context uses for one token.
func NThreads(ctx Context) int32 { return contextInt("_yzma_context_n_threads", ctx) }

// NThreadsBatch gives the number of threads that the context uses for a batch
// of more than one token.
func NThreadsBatch(ctx Context) int32 { return contextInt("_yzma_context_n_threads_batch", ctx) }

// NRsSeq gives the number of sequences of the recurrent state of the context.
func NRsSeq(ctx Context) int32 { return contextInt("_yzma_context_n_rs_seq", ctx) }

// SetNThreads changes the threads of a context after it is made.
// ContextParams sets the same thing when the context is made.
func SetNThreads(ctx Context, nThreads, nThreadsBatch int32) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	if !has("_yzma_set_n_threads") {
		return ErrNoContextFlags
	}
	_, err := callErr("_yzma_set_n_threads", int(ctx), int(nThreads), int(nThreadsBatch))
	return err
}

// contextInt reads a whole number of a context. It gives 0 when the module has
// no such call and when the call fails.
func contextInt(name string, ctx Context) int32 {
	if !has(name) {
		return 0
	}
	n := call(name, int(ctx))
	if n < 0 {
		return 0
	}
	return n
}
