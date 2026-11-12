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

	// LLAMA_API bool llama_set_embeddings_nextn_dev(struct llama_context * ctx, bool value);
	setEmbeddingsNextNDevFunc ffi.Fun

	// LLAMA_API bool llama_get_embeddings_nextn_rows(struct llama_context * ctx, int64_t row0, int64_t n_rows, float * out);
	getEmbeddingsNextNRowsFunc ffi.Fun

	// LLAMA_API bool llama_set_embeddings_nextn_rows(struct llama_context * ctx, int64_t row0, int64_t n_rows, const float * data);
	setEmbeddingsNextNRowsFunc ffi.Fun

	// LLAMA_API bool llama_set_embd_state_src(struct llama_context * ctx, struct llama_context * src);
	setEmbdStateSrcFunc ffi.Fun

	// LLAMA_API bool llama_batch_ext_set_embd_state_row(struct llama_batch_ext * batch, int32_t idx, int32_t row);
	batchExtSetEmbdStateRowFunc ffi.Fun
)

// prepAny binds the first name that the library exports. It returns the zero
// ffi.Fun when the library exports none of them.
func prepAny(lib loader.Lib, names []string, ret *ffi.Type, args ...*ffi.Type) ffi.Fun {
	for _, name := range names {
		if fn, err := lib.Prep(name, ret, args...); err == nil {
			return fn
		}
	}

	return ffi.Fun{}
}

// Load loads the shared llama.cpp library and binds its NextN functions. An
// empty path uses the path that [llama.Load] used, then the YZMA_LIB env
// variable.
//
// Load returns an error only when the library fails to open. A function that
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
			"?llama_set_embeddings_nextn@@YAXPEAUllama_context@@_N1@Z",
		},
		&ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeUint8, &ffi.TypeUint8)

	getEmbeddingsNextNFunc = prepAny(lib,
		[]string{
			"llama_get_embeddings_nextn",
			"_Z26llama_get_embeddings_nextnP13llama_context",
			"?llama_get_embeddings_nextn@@YAPEAMPEAUllama_context@@@Z",
		},
		&ffi.TypePointer, &ffi.TypePointer)

	getEmbeddingsNextNIthFunc = prepAny(lib,
		[]string{
			"llama_get_embeddings_nextn_ith",
			"_Z30llama_get_embeddings_nextn_ithP13llama_contexti",
			"?llama_get_embeddings_nextn_ith@@YAPEAMPEAUllama_context@@H@Z",
		},
		&ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32)

	setNextNLayerOffsetFunc = prepAny(lib,
		[]string{
			"llama_set_nextn_layer_offset",
			"_Z28llama_set_nextn_layer_offsetP13llama_contexti",
			"?llama_set_nextn_layer_offset@@YAXPEAUllama_context@@H@Z",
		},
		&ffi.TypeVoid, &ffi.TypePointer, &ffi.TypeSint32)

	setEmbeddingsNextNDevFunc = prepAny(lib,
		[]string{
			"llama_set_embeddings_nextn_dev",
			"_Z30llama_set_embeddings_nextn_devP13llama_contextb",
			"?llama_set_embeddings_nextn_dev@@YA_NPEAUllama_context@@_N@Z",
		},
		&ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeUint8)

	getEmbeddingsNextNRowsFunc = prepAny(lib,
		[]string{
			"llama_get_embeddings_nextn_rows",
			"_Z31llama_get_embeddings_nextn_rowsP13llama_contextllPf",
			"?llama_get_embeddings_nextn_rows@@YA_NPEAUllama_context@@_J1PEAM@Z",
		},
		&ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint64, &ffi.TypeSint64, &ffi.TypePointer)

	setEmbeddingsNextNRowsFunc = prepAny(lib,
		[]string{
			"llama_set_embeddings_nextn_rows",
			"_Z31llama_set_embeddings_nextn_rowsP13llama_contextllPKf",
			"?llama_set_embeddings_nextn_rows@@YA_NPEAUllama_context@@_J1PEBM@Z",
		},
		&ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint64, &ffi.TypeSint64, &ffi.TypePointer)

	setEmbdStateSrcFunc = prepAny(lib,
		[]string{
			"llama_set_embd_state_src",
			"_Z24llama_set_embd_state_srcP13llama_contextS0_",
			"?llama_set_embd_state_src@@YA_NPEAUllama_context@@0@Z",
		},
		&ffi.TypeUint8, &ffi.TypePointer, &ffi.TypePointer)

	batchExtSetEmbdStateRowFunc = prepAny(lib,
		[]string{
			"llama_batch_ext_set_embd_state_row",
			"_Z34llama_batch_ext_set_embd_state_rowP15llama_batch_extii",
			"?llama_batch_ext_set_embd_state_row@@YA_NPEAUllama_batch_ext@@HH@Z",
		},
		&ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32)

	return nil
}

// Available reports whether the loaded llama.cpp library exports all of the
// NextN functions. Callers must check it before they take an MTP path.
func Available() bool {
	return setEmbeddingsNextNFunc.Cif != nil &&
		getEmbeddingsNextNFunc.Cif != nil &&
		getEmbeddingsNextNIthFunc.Cif != nil &&
		setNextNLayerOffsetFunc.Cif != nil
}

// DeviceRowsAvailable reports whether the loaded llama.cpp library exports the
// device row functions. It is separate from [Available] because older builds
// have the host NextN functions only.
func DeviceRowsAvailable() bool {
	return setEmbeddingsNextNDevFunc.Cif != nil &&
		getEmbeddingsNextNRowsFunc.Cif != nil &&
		setEmbeddingsNextNRowsFunc.Cif != nil &&
		setEmbdStateSrcFunc.Cif != nil &&
		batchExtSetEmbdStateRowFunc.Cif != nil
}
