//go:build js && wasm

package llamawasm

import "syscall/js"

// The calls here save and restore context state in memory. Each one follows
// the same call in pkg/llama. A module before ABI version 9 has none of them,
// so each returns 0.
//
// A program can keep the state in IndexedDB or OPFS, so a page that reloads
// does not have to process the same prompt again.

// StateGetSize returns the size of the context state in bytes.
func StateGetSize(ctx Context) uint64 {
	return stateUint64("_yzma_state_get_size", int(ctx))
}

// StateGetData copies the context state into dst and returns the number of
// bytes copied.
func StateGetData(ctx Context, dst []byte) uint64 {
	return stateGet("_yzma_state_get_data", dst, int(ctx))
}

// StateSetData restores the context state from src and returns the number of
// bytes read.
func StateSetData(ctx Context, src []byte) uint64 {
	return stateSet("_yzma_state_set_data", src, int(ctx))
}

// StateSeqGetSize returns the size of one sequence state in bytes.
func StateSeqGetSize(ctx Context, seqId SeqId) uint64 {
	return StateSeqGetSizeExt(ctx, seqId, 0)
}

// StateSeqGetData copies the state of one sequence into dst and returns the
// number of bytes copied.
func StateSeqGetData(ctx Context, dst []byte, seqId SeqId) uint64 {
	return StateSeqGetDataExt(ctx, dst, seqId, 0)
}

// StateSeqSetData restores a sequence state from src into destSeqId and
// returns the number of bytes read.
func StateSeqSetData(ctx Context, src []byte, destSeqId SeqId) uint64 {
	return StateSeqSetDataExt(ctx, src, destSeqId, 0)
}

// StateSeqGetSizeExt is StateSeqGetSize with llama.cpp flags.
func StateSeqGetSizeExt(ctx Context, seqId SeqId, flags uint32) uint64 {
	return stateUint64("_yzma_state_seq_get_size_ext", int(ctx), int(seqId), int(flags))
}

// StateSeqGetDataExt is StateSeqGetData with llama.cpp flags.
func StateSeqGetDataExt(ctx Context, dst []byte, seqId SeqId, flags uint32) uint64 {
	return stateGet("_yzma_state_seq_get_data_ext", dst, int(ctx), int(seqId), int(flags))
}

// StateSeqSetDataExt is StateSeqSetData with llama.cpp flags.
func StateSeqSetDataExt(ctx Context, src []byte, destSeqId SeqId, flags uint32) uint64 {
	return stateSet("_yzma_state_seq_set_data_ext", src, int(ctx), int(destSeqId), int(flags))
}

// stateUint64 runs a state call. The shim returns a double, because a state
// can be larger than an int32.
func stateUint64(name string, args ...any) uint64 {
	if !has(name) {
		return 0
	}
	v, ok := callNumber(name, args...)
	if !ok || v < 0 {
		return 0
	}
	return uint64(v)
}

// stateGet runs a call that copies a state into a module buffer, then copies
// the buffer into dst. The context comes first and the other args come after
// the buffer and its size.
func stateGet(name string, dst []byte, ctx int, rest ...any) uint64 {
	if len(dst) == 0 || !has(name) {
		return 0
	}

	ptr, err := malloc(len(dst))
	if err != nil {
		return 0
	}
	defer free(ptr)

	n := stateUint64(name, append([]any{ctx, ptr, len(dst)}, rest...)...)
	if n > uint64(len(dst)) {
		return 0
	}
	js.CopyBytesToGo(dst[:n], view(ptr, int(n)))
	return n
}

// stateSet copies src into a module buffer and runs a call that restores a
// state from it.
func stateSet(name string, src []byte, ctx int, rest ...any) uint64 {
	if len(src) == 0 || !has(name) {
		return 0
	}

	ptr, err := malloc(len(src))
	if err != nil {
		return 0
	}
	defer free(ptr)

	writeBytes(ptr, src)
	return stateUint64(name, append([]any{ctx, ptr, len(src)}, rest...)...)
}
