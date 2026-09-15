package speculative

import (
	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/loader"
	"github.com/jupiterrider/ffi"
)

var (
	// LLAMA_API void llama_set_embeddings_nextn(
	// 				struct llama_context * ctx,
	// 				bool value,
	// 				bool masked);
	setEmbeddingsNextNFunc ffi.Fun

	// LLAMA_API float * llama_get_embeddings_nextn(struct llama_context * ctx);
	getEmbeddingsNextNFunc ffi.Fun

	// LLAMA_API float * llama_get_embeddings_nextn_ith(
	// 				struct llama_context * ctx,
	// 				int32_t i);
	getEmbeddingsNextNIthFunc ffi.Fun

	// LLAMA_API void llama_set_nextn_layer_offset(
	// 				struct llama_context * ctx,
	// 				int32_t offset);
	setNextNLayerOffsetFunc ffi.Fun
)

// prepAny gets the first name that the library exports. It returns the zero
// ffi.Fun when the library has none of them.
func prepAny(lib loader.Lib, names []string, ret *ffi.Type, args ...*ffi.Type) ffi.Fun {
	for _, name := range names {
		if fn, err := lib.Prep(name, ret, args...); err == nil {
			return fn
		}
	}

	return ffi.Fun{}
}

// Load loads the shared llama.cpp library and gets the NextN functions from
// it. An empty path uses the path that [llama.Load] used, then the YZMA_LIB
// env variable.
//
// Load reports an error only when the library does not open. A function that
// the library does not export stays unbound, and [Available] then reports
// false. This keeps older llama.cpp builds usable.
func Load(path string) error {
	if path == "" {
		path = llama.LibPath()
	}

	lib, err := loader.LoadLibrary(path, "llama")
	if err != nil {
		return err
	}

	// yzma-checker does not audit these signatures. llama.cpp does not install
	// src/llama-ext.h, so they are checked by hand against that header.
	setEmbeddingsNextNFunc = prepAny(lib,
		[]string{
			"llama_set_embeddings_nextn",
			"_Z26llama_set_embeddings_nextnP13llama_contextbb",
		},
		&ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeUint8, &ffi.TypeUint8)

	getEmbeddingsNextNFunc = prepAny(lib,
		[]string{
			"llama_get_embeddings_nextn",
			"_Z26llama_get_embeddings_nextnP13llama_context",
		},
		&ffi.TypePointer, &ffi.TypePointer)

	getEmbeddingsNextNIthFunc = prepAny(lib,
		[]string{
			"llama_get_embeddings_nextn_ith",
			"_Z30llama_get_embeddings_nextn_ithP13llama_contexti",
		},
		&ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32)

	setNextNLayerOffsetFunc = prepAny(lib,
		[]string{
			"llama_set_nextn_layer_offset",
			"_Z28llama_set_nextn_layer_offsetP13llama_contexti",
		},
		&ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)

	return nil
}

// Available tells if the loaded llama.cpp library exports all of the NextN
// functions. Callers must use it to select an MTP path.
func Available() bool {
	return setEmbeddingsNextNFunc.Cif != nil &&
		getEmbeddingsNextNFunc.Cif != nil &&
		getEmbeddingsNextNIthFunc.Cif != nil &&
		setNextNLayerOffsetFunc.Cif != nil
}
