package mtmd

import (
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

// enum mtmd_input_chunk_type
type InputChunkType int32

const (
	InputChunkTypeText InputChunkType = iota
	InputChunkTypeImage
	InputChunkTypeAudio
)

var (
	// MTMD_API mtmd_input_chunks *      mtmd_input_chunks_init(void);
	inputChunksInitFunc ffi.Fun

	// MTMD_API void mtmd_input_chunks_free(mtmd_input_chunks * chunks);
	inputChunksFreeFunc ffi.Fun

	// MTMD_API size_t mtmd_input_chunks_size(const mtmd_input_chunks * chunks);
	inputChunksSizeFunc ffi.Fun

	// MTMD_API const mtmd_input_chunk * mtmd_input_chunks_get(const mtmd_input_chunks * chunks, size_t idx);
	inputChunksGetFunc ffi.Fun

	// MTMD_API enum mtmd_input_chunk_type mtmd_input_chunk_get_type(const mtmd_input_chunk * chunk);
	inputChunkGetTypeFunc ffi.Fun

	// MTMD_API const llama_token * mtmd_input_chunk_get_tokens_text(const mtmd_input_chunk * chunk, size_t * n_tokens_output);
	inputChunkGetTokensTextFunc ffi.Fun

	// MTMD_API const mtmd_image_tokens * mtmd_input_chunk_get_tokens_image(const mtmd_input_chunk * chunk);
	inputChunkGetTokensImageFunc ffi.Fun

	// MTMD_API size_t mtmd_input_chunk_get_n_tokens(const mtmd_input_chunk * chunk);
	inputChunkGetNTokensFunc ffi.Fun

	// MTMD_API const char * mtmd_input_chunk_get_id(const mtmd_input_chunk * chunk);
	inputChunkGetIdFunc ffi.Fun

	// MTMD_API llama_pos mtmd_input_chunk_get_n_pos(const mtmd_input_chunk * chunk);
	inputChunkGetNPosFunc ffi.Fun

	// MTMD_API mtmd_input_chunk * mtmd_input_chunk_copy(const mtmd_input_chunk * chunk);
	inputChunkCopyFunc ffi.Fun

	// MTMD_API void mtmd_input_chunk_free(mtmd_input_chunk * chunk);
	inputChunkFreeFunc ffi.Fun
)

func loadChunkFuncs(lib ffi.Lib) error {
	var err error

	if inputChunksInitFunc, err = lib.Prep("mtmd_input_chunks_init", &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunks_init", err)
	}

	if inputChunksFreeFunc, err = lib.Prep("mtmd_input_chunks_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunks_free", err)
	}

	if inputChunksSizeFunc, err = lib.Prep("mtmd_input_chunks_size", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunks_size", err)
	}

	if inputChunksGetFunc, err = lib.Prep("mtmd_input_chunks_get", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return loadError("mtmd_input_chunks_get", err)
	}

	if inputChunkGetTypeFunc, err = lib.Prep("mtmd_input_chunk_get_type", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunk_get_type", err)
	}

	if inputChunkGetTokensTextFunc, err = lib.Prep("mtmd_input_chunk_get_tokens_text", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunk_get_tokens_text", err)
	}

	if inputChunkGetNTokensFunc, err = lib.Prep("mtmd_input_chunk_get_n_tokens", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunk_get_n_tokens", err)
	}

	if inputChunkGetIdFunc, err = lib.Prep("mtmd_input_chunk_get_id", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunk_get_id", err)
	}

	if inputChunkGetNPosFunc, err = lib.Prep("mtmd_input_chunk_get_n_pos", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunk_get_n_pos", err)
	}

	if inputChunkCopyFunc, err = lib.Prep("mtmd_input_chunk_copy", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunk_copy", err)
	}

	if inputChunkFreeFunc, err = lib.Prep("mtmd_input_chunk_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("mtmd_input_chunk_free", err)
	}

	return nil
}

// InputChunksInit initializes a list of InputChunk.
// It can only be populated via Tokenize().
func InputChunksInit() InputChunks {
	var chunks InputChunks
	inputChunksInitFunc.Call(unsafe.Pointer(&chunks))

	return chunks
}

// InputChunksFree frees the InputChunks.
func InputChunksFree(chunks InputChunks) {
	inputChunksFreeFunc.Call(nil, unsafe.Pointer(&chunks))
}

// InputChunksSize returns the number of InputChunk in the list.
func InputChunksSize(chunks InputChunks) uint32 {
	var result ffi.Arg
	inputChunksSizeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&chunks))

	return uint32(result)
}

// InputChunksGet retrieves the input chunk at the specified index.
func InputChunksGet(chunks InputChunks, idx uint32) InputChunk {
	var chunk InputChunk
	inputChunksGetFunc.Call(unsafe.Pointer(&chunk), unsafe.Pointer(&chunks), unsafe.Pointer(&idx))
	return chunk
}

// InputChunkGetType retrieves the type of the input chunk.
func InputChunkGetType(chunk InputChunk) InputChunkType {
	var result ffi.Arg
	inputChunkGetTypeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&chunk))
	return InputChunkType(result)
}

// InputChunkGetTokensText retrieves the text tokens of the input chunk.
func InputChunkGetTokensText(chunk InputChunk) []llama.Token {
	var tokensPtr *llama.Token
	var nTokens uint32
	inputChunkGetTokensTextFunc.Call(unsafe.Pointer(&tokensPtr), unsafe.Pointer(&chunk), unsafe.Pointer(&nTokens))

	if tokensPtr == nil || nTokens == 0 {
		return nil
	}

	return unsafe.Slice(tokensPtr, int(nTokens))
}

// InputChunkGetNTokens retrieves the number of tokens in the input chunk.
func InputChunkGetNTokens(chunk InputChunk) uint32 {
	var result ffi.Arg
	inputChunkGetNTokensFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&chunk))
	return uint32(result)
}

// InputChunkGetId retrieves the ID of the input chunk.
func InputChunkGetId(chunk InputChunk) string {
	var idPtr *byte
	inputChunkGetIdFunc.Call(unsafe.Pointer(&idPtr), unsafe.Pointer(&chunk))

	if idPtr == nil {
		return ""
	}

	return utils.BytePtrToString(idPtr)
}

// InputChunkGetNPos retrieves the number of temporal positions in the input chunk.
func InputChunkGetNPos(chunk InputChunk) llama.Pos {
	var result ffi.Arg
	inputChunkGetNPosFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&chunk))
	return llama.Pos(result)
}

// InputChunkCopy creates a copy of the input chunk.
func InputChunkCopy(chunk InputChunk) InputChunk {
	var copy InputChunk
	inputChunkCopyFunc.Call(unsafe.Pointer(&copy), unsafe.Pointer(&chunk))
	return copy
}

// InputChunkFree frees the input chunk.
func InputChunkFree(chunk InputChunk) {
	inputChunkFreeFunc.Call(nil, unsafe.Pointer(&chunk))
}
