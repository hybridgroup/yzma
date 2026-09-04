//go:build js && wasm

package llamawasm

import "errors"

// ErrNoMemorySeq says that the module is from a release before the calls for
// the memory of one sequence, which are in ABI version 6 and later.
var ErrNoMemorySeq = errors.New("llamawasm: this llama.cpp module has no calls for the memory of a sequence, install a newer build")

// The calls of this file take a context and not a llama.Memory, because the
// shim holds the memory itself. This is the same as [MemoryClear].

// MemoryClear removes everything from the memory of the context, which starts
// the generation again from an empty state.
//
// On a native platform this takes a llama.Memory that comes from
// llama.GetMemory. Here it takes the context, because the shim holds the
// memory itself.
func MemoryClear(ctx Context, data bool) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	_, err := callErr("_yzma_memory_clear", int(ctx), boolInt(data))
	return err
}

// MemorySeqRm removes the tokens of a sequence between the positions p0 and
// p1. A p0 below zero is the start of the sequence and a p1 below zero is the
// end of it. A seqID below zero matches every sequence.
//
// The result is false if the memory cannot remove the tokens, which happens
// when it keeps whole sequences only. The sequence stays as it was then.
func MemorySeqRm(ctx Context, seqID SeqId, p0, p1 Pos) (bool, error) {
	rc, err := memoryCall("_yzma_memory_seq_rm", int(ctx), int(seqID), int(p0), int(p1))
	return rc == 1, err
}

// MemorySeqCp copies the tokens of one sequence between the positions p0 and
// p1 into another sequence.
func MemorySeqCp(ctx Context, seqIDSrc, seqIDDst SeqId, p0, p1 Pos) error {
	_, err := memoryCall("_yzma_memory_seq_cp", int(ctx), int(seqIDSrc), int(seqIDDst), int(p0), int(p1))
	return err
}

// MemorySeqKeep removes every sequence of the memory except one.
func MemorySeqKeep(ctx Context, seqID SeqId) error {
	_, err := memoryCall("_yzma_memory_seq_keep", int(ctx), int(seqID))
	return err
}

// MemorySeqAdd moves the tokens of a sequence between the positions p0 and p1
// by delta positions. A program shifts a context that is full this way.
//
// Examine [MemoryCanShift] first, because not every memory accepts this.
func MemorySeqAdd(ctx Context, seqID SeqId, p0, p1, delta Pos) error {
	_, err := memoryCall("_yzma_memory_seq_add", int(ctx), int(seqID), int(p0), int(p1), int(delta))
	return err
}

// MemorySeqDiv divides the positions of the tokens of a sequence between p0
// and p1 by d.
func MemorySeqDiv(ctx Context, seqID SeqId, p0, p1 Pos, d int) error {
	_, err := memoryCall("_yzma_memory_seq_div", int(ctx), int(seqID), int(p0), int(p1), d)
	return err
}

// MemorySeqPosMin gives the smallest position of a sequence, and -1 if the
// sequence holds no tokens.
func MemorySeqPosMin(ctx Context, seqID SeqId) (Pos, error) {
	return memoryPos("_yzma_memory_seq_pos_min", ctx, seqID)
}

// MemorySeqPosMax gives the largest position of a sequence, and -1 if the
// sequence holds no tokens.
func MemorySeqPosMax(ctx Context, seqID SeqId) (Pos, error) {
	return memoryPos("_yzma_memory_seq_pos_max", ctx, seqID)
}

// MemoryCanShift tells if the memory accepts [MemorySeqAdd].
func MemoryCanShift(ctx Context) (bool, error) {
	rc, err := memoryCall("_yzma_memory_can_shift", int(ctx))
	return rc == 1, err
}

// memoryCall runs a call of the memory and reports a module that is too old as
// a clear error.
func memoryCall(name string, args ...any) (int32, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if !has(name) {
		return 0, ErrNoMemorySeq
	}
	return callErr(name, args...)
}

// memoryPos reads a position of a sequence. A position of -1 is normal, thus
// the shim reports a bad handle with errBadHandle and not with a small
// negative value.
func memoryPos(name string, ctx Context, seqID SeqId) (Pos, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if !has(name) {
		return 0, ErrNoMemorySeq
	}

	rc := call(name, int(ctx), int(seqID))
	if rc <= errBadHandle {
		return 0, shimError(name, rc)
	}
	return Pos(rc), nil
}

// boolInt gives the value that the shim takes for a bool.
func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
