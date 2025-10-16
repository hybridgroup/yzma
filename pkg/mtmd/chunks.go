package mtmd

import (
	"unsafe"

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
