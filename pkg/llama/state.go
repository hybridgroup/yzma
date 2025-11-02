package llama

import (
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

var (
	// LLAMA_API bool llama_state_save_file(
	//     struct llama_context * ctx,
	//     const char * path_session,
	//     const llama_token * tokens,
	//     size_t   n_token_count);
	stateSaveFileFunc ffi.Fun

	// LLAMA_API bool llama_state_load_file(
	//     struct llama_context * ctx,
	//               const char * path_session,
	//              llama_token * tokens_out,
	//                   size_t   n_token_capacity,
	//                   size_t * n_token_count_out);
	stateLoadFileFunc ffi.Fun

	// Returns the *actual* size in bytes of the state
	// (logits, embedding and memory)
	// Only use when saving the state, not when restoring it, otherwise the size may be too small.
	// LLAMA_API size_t llama_state_get_size(struct llama_context * ctx);
	stateGetSizeFunc ffi.Fun

	// Copies the state to the specified destination address.
	// Destination needs to have allocated enough memory.
	// Returns the number of bytes copied
	// LLAMA_API size_t llama_state_get_data(
	//         struct llama_context * ctx,
	//                      uint8_t * dst,
	//                       size_t   size);
	stateGetDataFunc ffi.Fun

	// Set the state reading from the specified address
	// Returns the number of bytes read
	// LLAMA_API size_t llama_state_set_data(
	//         struct llama_context * ctx,
	//                const uint8_t * src,
	//                       size_t   size);
	stateSetDataFunc ffi.Fun
)

func loadStateFuncs(lib ffi.Lib) error {
	var err error

	if stateSaveFileFunc, err = lib.Prep("llama_state_save_file", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_state_save_file", err)
	}

	if stateLoadFileFunc, err = lib.Prep("llama_state_load_file", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("llama_state_load_file", err)
	}

	if stateGetSizeFunc, err = lib.Prep("llama_state_get_size", &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("llama_state_get_size", err)
	}

	if stateGetDataFunc, err = lib.Prep("llama_state_get_data", &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_state_get_data", err)
	}

	if stateSetDataFunc, err = lib.Prep("llama_state_set_data", &ffi.TypeUint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_state_set_data", err)
	}

	return err
}

// StateSaveFile saves the state to a file and returns true on success.
func StateSaveFile(ctx Context, path string, tokens []Token) bool {
	pathPtr, _ := utils.BytePtrFromString(path)
	var toks *Token
	if len(tokens) > 0 {
		toks = unsafe.SliceData(tokens)
	}
	tlen := int64(len(tokens))

	var result ffi.Arg
	stateSaveFileFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), unsafe.Pointer(&pathPtr), unsafe.Pointer(&toks), &tlen)
	return result.Bool()
}

// StateLoadFile loads the state from a file and returns true on success.
// tokensOut should be a slice with capacity nTokenCapacity. nTokenCountOut will be set to the number of tokens loaded.
func StateLoadFile(ctx Context, path string, tokensOut []Token, nTokenCapacity uint64, nTokenCountOut *uint64) bool {
	pathPtr, _ := utils.BytePtrFromString(path)
	var toks *Token
	if len(tokensOut) > 0 {
		toks = unsafe.SliceData(tokensOut)
	}
	var result ffi.Arg
	stateLoadFileFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&pathPtr),
		unsafe.Pointer(&toks),
		unsafe.Pointer(&nTokenCapacity),
		unsafe.Pointer(&nTokenCountOut),
	)
	return result.Bool()
}

// StateGetSize returns the actual size in bytes of the state (logits, embedding and memory).
// Only use when saving the state, not when restoring it.
func StateGetSize(ctx Context) uint64 {
	var result ffi.Arg
	stateGetSizeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx))
	return uint64(result)
}

// StateGetData copies the state to the specified destination address.
// Returns the number of bytes copied.
func StateGetData(ctx Context, dst []byte) uint64 {
	var result ffi.Arg
	var size int64 = int64(len(dst))
	var dstPtr *byte
	if len(dst) > 0 {
		dstPtr = &dst[0]
	}
	stateGetDataFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), unsafe.Pointer(&dstPtr), &size)
	return uint64(result)
}

// StateSetData sets the state by reading from the specified address.
// Returns the number of bytes read.
func StateSetData(ctx Context, src []byte) uint64 {
	var result ffi.Arg
	var size int64 = int64(len(src))
	var srcPtr *byte
	if len(src) > 0 {
		srcPtr = &src[0]
	}
	stateSetDataFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), unsafe.Pointer(&srcPtr), &size)
	return uint64(result)
}
