package llama

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

var (
	// LLAMA_API void llama_memory_clear(
	// 				llama_memory_t mem,
	// 				bool data);
	memoryClearFunc ffi.Fun

	// LLAMA_API bool llama_memory_seq_rm(
	//         		llama_memory_t mem,
	//           	llama_seq_id seq_id,
	//              llama_pos p0,
	//              llama_pos p1);
	memorySeqRmFunc ffi.Fun
)

func loadMemoryFuncs(lib ffi.Lib) error {
	var err error

	if memoryClearFunc, err = lib.Prep("llama_memory_clear", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeUint8); err != nil {
		return loadError("llama_memory_clear", err)
	}
	if memorySeqRmFunc, err = lib.Prep("llama_memory_seq_rm", &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32, &ffi.TypeSint32); err != nil {
		return loadError("llama_memory_seq_rm", err)
	}

	return nil
}

// MemoryClear clears the memory contents.
// If data == true, the data buffers will also be cleared together with the metadata.
func MemoryClear(mem Memory, data bool) {
	memoryClearFunc.Call(nil, unsafe.Pointer(&mem), unsafe.Pointer(&data))
}

// MemorySeqRm removes all tokens that belong to the specified sequence and have positions in [p0, p1).
// Returns false if a partial sequence cannot be removed. Removing a whole sequence never fails.
// seqID < 0 : match any sequence
// p0 < 0     : [0,  p1]
// p1 < 0     : [p0, inf)
func MemorySeqRm(mem Memory, seqID SeqId, p0, p1 Pos) bool {
	var result ffi.Arg
	memorySeqRmFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&mem), &seqID, &p0, &p1)

	return result.Bool()
}
