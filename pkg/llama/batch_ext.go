package llama

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/loader"
	"github.com/jupiterrider/ffi"
)

// BatchExt is a handle to a llama.cpp extended batch.
type BatchExt uintptr

// ProcessType selects how [Process] runs a [BatchExt].
type ProcessType int32

const (
	ProcessTypeEncode ProcessType = 0
	ProcessTypeDecode ProcessType = 1
)

// mropeSections is GGML_MROPE_SECTIONS, the most positions llama.cpp reads for one entry.
const mropeSections = 4

// Errors reported by the [BatchExt] functions.
var (
	// ErrBatchExtInvalidToken means the token ID or embedding was rejected.
	ErrBatchExtInvalidToken = errors.New("invalid token or embedding for batch")
	// ErrBatchExtInvalidSeqID means the sequence ID is not below the context n_seq_max.
	ErrBatchExtInvalidSeqID = errors.New("invalid sequence ID for batch")
	// ErrBatchExtRejected means llama.cpp rejected the index or value.
	ErrBatchExtRejected = errors.New("batch rejected the value")

	errInvalidBatchExt = errors.New("invalid extended batch")
)

var (
	// ffiTypeEmbd represents the C struct llama_embd
	ffiTypeEmbd = ffi.NewType(&ffi.TypePointer, &ffiTypeSize, &ffiTypeSize)
)

// embd mirrors the C struct llama_embd.
type embd struct {
	data  *float32
	nRows uint64
	nEmbd uint64
}

var (
	// LLAMA_API struct llama_batch_ext * llama_batch_ext_init (struct llama_context * ctx);
	batchExtInitFunc ffi.Fun

	// LLAMA_API void llama_batch_ext_free (struct llama_batch_ext * batch);
	batchExtFreeFunc ffi.Fun

	// LLAMA_API void llama_batch_ext_clear(struct llama_batch_ext * batch);
	batchExtClearFunc ffi.Fun

	// LLAMA_API int32_t llama_batch_ext_add(struct llama_batch_ext * batch, llama_seq_id seq_id);
	batchExtAddFunc ffi.Fun

	// LLAMA_API int32_t llama_batch_ext_add_token(struct llama_batch_ext * batch, llama_seq_id seq_id, llama_token id);
	batchExtAddTokenFunc ffi.Fun

	// LLAMA_API int32_t llama_batch_ext_add_embd (struct llama_batch_ext * batch, llama_seq_id seq_id, struct llama_embd embd);
	batchExtAddEmbdFunc ffi.Fun

	// LLAMA_API bool llama_batch_ext_add_seq(struct llama_batch_ext * batch, int32_t idx, llama_seq_id seq_id);
	batchExtAddSeqFunc ffi.Fun

	// LLAMA_API bool llama_batch_ext_set_embd_token(struct llama_batch_ext * batch, int32_t idx, struct llama_embd embd);
	batchExtSetEmbdTokenFunc ffi.Fun

	// LLAMA_API bool llama_batch_ext_set_embd_state(struct llama_batch_ext * batch, int32_t idx, struct llama_embd embd);
	batchExtSetEmbdStateFunc ffi.Fun

	// LLAMA_API bool llama_batch_ext_set_output_embd(struct llama_batch_ext * batch, int32_t idx, bool value);
	batchExtSetOutputEmbdFunc ffi.Fun

	// LLAMA_API bool llama_batch_ext_set_output_logits(struct llama_batch_ext * batch, int32_t idx, bool value);
	batchExtSetOutputLogitsFunc ffi.Fun

	// LLAMA_API bool llama_batch_ext_set_pos(struct llama_batch_ext * batch, int32_t idx, const llama_pos * pos);
	batchExtSetPosFunc ffi.Fun

	// LLAMA_API int32_t llama_process(struct llama_context * ctx, enum llama_process_type type, struct llama_batch_ext * batch);
	processFunc ffi.Fun
)

func loadBatchExtFuncs(lib loader.Lib) error {
	var err error

	if batchExtInitFunc, err = lib.Prep("llama_batch_ext_init", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_batch_ext_init", err)
	}

	if batchExtFreeFunc, err = lib.Prep("llama_batch_ext_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("llama_batch_ext_free", err)
	}

	if batchExtClearFunc, err = lib.Prep("llama_batch_ext_clear", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("llama_batch_ext_clear", err)
	}

	if batchExtAddFunc, err = lib.Prep("llama_batch_ext_add", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("llama_batch_ext_add", err)
	}

	if batchExtAddTokenFunc, err = lib.Prep("llama_batch_ext_add_token", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32); err != nil {
		return loadError("llama_batch_ext_add_token", err)
	}

	if batchExtAddEmbdFunc, err = lib.Prep("llama_batch_ext_add_embd", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffiTypeEmbd); err != nil {
		return loadError("llama_batch_ext_add_embd", err)
	}

	if batchExtAddSeqFunc, err = lib.Prep("llama_batch_ext_add_seq", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32); err != nil {
		return loadError("llama_batch_ext_add_seq", err)
	}

	if batchExtSetEmbdTokenFunc, err = lib.Prep("llama_batch_ext_set_embd_token", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffiTypeEmbd); err != nil {
		return loadError("llama_batch_ext_set_embd_token", err)
	}

	if batchExtSetEmbdStateFunc, err = lib.Prep("llama_batch_ext_set_embd_state", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffiTypeEmbd); err != nil {
		return loadError("llama_batch_ext_set_embd_state", err)
	}

	if batchExtSetOutputEmbdFunc, err = lib.Prep("llama_batch_ext_set_output_embd", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeUint8); err != nil {
		return loadError("llama_batch_ext_set_output_embd", err)
	}

	if batchExtSetOutputLogitsFunc, err = lib.Prep("llama_batch_ext_set_output_logits", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeUint8); err != nil {
		return loadError("llama_batch_ext_set_output_logits", err)
	}

	if batchExtSetPosFunc, err = lib.Prep("llama_batch_ext_set_pos", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_batch_ext_set_pos", err)
	}

	if processFunc, err = lib.Prep("llama_process", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_process", err)
	}

	return nil
}

// BatchExtInit creates an extended batch for the context.
// It holds up to [NBatch] entries and accepts sequence IDs below [NSeqMax].
// The batch has to be freed with [BatchExtFree].
func BatchExtInit(ctx Context) (BatchExt, error) {
	if ctx == 0 {
		return 0, errInvalidContext
	}
	var batch BatchExt
	batchExtInitFunc.Call(unsafe.Pointer(&batch), unsafe.Pointer(&ctx))
	if batch == 0 {
		return 0, errInvalidBatchExt
	}

	return batch, nil
}

// BatchExtFree frees a batch created with [BatchExtInit].
func BatchExtFree(batch BatchExt) error {
	if batch == 0 {
		return errInvalidBatchExt
	}
	batchExtFreeFunc.Call(nil, unsafe.Pointer(&batch))

	return nil
}

// BatchExtClear removes all entries from the batch.
func BatchExtClear(batch BatchExt) error {
	if batch == 0 {
		return errInvalidBatchExt
	}
	batchExtClearFunc.Call(nil, unsafe.Pointer(&batch))

	return nil
}

// BatchExtAdd adds an entry with no token and no embedding, and returns its index.
// Set its position with [BatchExtSetPos] before processing the batch.
func BatchExtAdd(batch BatchExt, seqID SeqId) (int32, error) {
	if batch == 0 {
		return 0, errInvalidBatchExt
	}
	var result ffi.Arg
	batchExtAddFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &seqID)

	return addResult(int32(result))
}

// BatchExtAddToken adds a token and returns its index.
// Set its position with [BatchExtSetPos] before processing the batch.
func BatchExtAddToken(batch BatchExt, seqID SeqId, token Token) (int32, error) {
	if batch == 0 {
		return 0, errInvalidBatchExt
	}
	var result ffi.Arg
	batchExtAddTokenFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &seqID, &token)

	return addResult(int32(result))
}

// BatchExtAddEmbd adds a token embedding made of rows of nEmbd floats, and returns its index.
// llama.cpp copies the data. Set its position with [BatchExtSetPos] before processing the batch.
func BatchExtAddEmbd(batch BatchExt, seqID SeqId, data []float32, nEmbd int) (int32, error) {
	if batch == 0 {
		return 0, errInvalidBatchExt
	}
	e, err := newEmbd(data, nEmbd)
	if err != nil {
		return 0, err
	}
	var result ffi.Arg
	batchExtAddEmbdFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &seqID, unsafe.Pointer(&e))
	runtime.KeepAlive(data)

	return addResult(int32(result))
}

// BatchExtAddSeq adds the entry at idx to another sequence at the same position.
// Call it before the other BatchExtSet functions for that entry.
func BatchExtAddSeq(batch BatchExt, idx int32, seqID SeqId) error {
	if batch == 0 {
		return errInvalidBatchExt
	}
	var result ffi.Arg
	batchExtAddSeqFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &idx, &seqID)

	return setResult(result, "add sequence %d to index %d", seqID, idx)
}

// BatchExtSetEmbdToken sets the token embedding for the entry at idx.
// Use it after [BatchExtAddToken] to give an entry both a token ID and an embedding.
func BatchExtSetEmbdToken(batch BatchExt, idx int32, data []float32, nEmbd int) error {
	if batch == 0 {
		return errInvalidBatchExt
	}
	e, err := newEmbd(data, nEmbd)
	if err != nil {
		return err
	}
	var result ffi.Arg
	batchExtSetEmbdTokenFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &idx, unsafe.Pointer(&e))
	runtime.KeepAlive(data)

	return setResult(result, "set token embedding at index %d", idx)
}

// BatchExtSetEmbdState sets the hidden state from a previous stage for the entry at idx.
// llama.cpp does not implement it yet, so it always returns [ErrBatchExtRejected].
func BatchExtSetEmbdState(batch BatchExt, idx int32, data []float32, nEmbd int) error {
	if batch == 0 {
		return errInvalidBatchExt
	}
	e, err := newEmbd(data, nEmbd)
	if err != nil {
		return err
	}
	var result ffi.Arg
	batchExtSetEmbdStateFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &idx, unsafe.Pointer(&e))
	runtime.KeepAlive(data)

	return setResult(result, "set state embedding at index %d", idx)
}

// BatchExtSetOutputEmbd sets whether output embeddings are computed for the entry at idx.
// llama.cpp currently treats this the same as [BatchExtSetOutputLogits].
func BatchExtSetOutputEmbd(batch BatchExt, idx int32, value bool) error {
	if batch == 0 {
		return errInvalidBatchExt
	}
	var result ffi.Arg
	batchExtSetOutputEmbdFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &idx, &value)

	return setResult(result, "set output embedding at index %d", idx)
}

// BatchExtSetOutputLogits sets whether logits are computed for the entry at idx.
// llama.cpp currently treats this the same as [BatchExtSetOutputEmbd].
func BatchExtSetOutputLogits(batch BatchExt, idx int32, value bool) error {
	if batch == 0 {
		return errInvalidBatchExt
	}
	var result ffi.Arg
	batchExtSetOutputLogitsFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &idx, &value)

	return setResult(result, "set output logits at index %d", idx)
}

// BatchExtSetPos sets the position of the entry at idx.
// A token uses one position. An embedding for an M-RoPE model uses one per section, up to 4.
func BatchExtSetPos(batch BatchExt, idx int32, pos ...Pos) error {
	if batch == 0 {
		return errInvalidBatchExt
	}
	if len(pos) == 0 || len(pos) > mropeSections {
		return fmt.Errorf("%w: %d positions, want 1 to %d", ErrBatchExtRejected, len(pos), mropeSections)
	}

	// llama.cpp reads as many positions as the entry needs, so pad to the maximum.
	var p [mropeSections]Pos
	copy(p[:], pos)
	pp := &p[0]
	var result ffi.Arg
	batchExtSetPosFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &idx, unsafe.Pointer(&pp))
	runtime.KeepAlive(&p)

	return setResult(result, "set position at index %d", idx)
}

// Process encodes or decodes the batch. The return values are the same as [Decode].
func Process(ctx Context, typ ProcessType, batch BatchExt) (int32, error) {
	if ctx == 0 {
		return 0, errInvalidContext
	}
	if batch == 0 {
		return 0, errInvalidBatchExt
	}
	var result ffi.Arg
	processFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &typ, unsafe.Pointer(&batch))

	return int32(result), nil
}

func newEmbd(data []float32, nEmbd int) (embd, error) {
	if nEmbd <= 0 || len(data) == 0 || len(data)%nEmbd != 0 {
		return embd{}, fmt.Errorf("%w: %d floats in rows of %d", ErrBatchExtInvalidToken, len(data), nEmbd)
	}

	return embd{data: unsafe.SliceData(data), nRows: uint64(len(data) / nEmbd), nEmbd: uint64(nEmbd)}, nil
}

func addResult(idx int32) (int32, error) {
	switch {
	case idx >= 0:
		return idx, nil
	case idx == -1:
		return idx, ErrBatchFull
	case idx == -2:
		return idx, ErrBatchExtInvalidToken
	case idx == -3:
		return idx, ErrBatchExtInvalidSeqID
	}

	return idx, fmt.Errorf("%w: code %d", ErrBatchExtRejected, idx)
}

func setResult(result ffi.Arg, format string, args ...any) error {
	if result.Bool() {
		return nil
	}

	return fmt.Errorf("%w: "+format, append([]any{ErrBatchExtRejected}, args...)...)
}
