package mtmd

import (
	"errors"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/loader"
	"github.com/jupiterrider/ffi"
)

// ErrNoResampler is returned when libmtmd does not have the resampler API.
var ErrNoResampler = errors.New("mtmd_helper_resampler not found in library")

// HelperResampler is an opaque handle to a streaming resampler for raw mono F32 PCM.
// It keeps state between calls, so use one per stream.
type HelperResampler uintptr

var (
	// MTMD_API mtmd_helper_resampler * mtmd_helper_resampler_init(int32_t rate_in, int32_t rate_out);
	helperResamplerInitFunc ffi.Fun

	// MTMD_API void mtmd_helper_resampler_free(mtmd_helper_resampler * ctx);
	helperResamplerFreeFunc ffi.Fun

	// MTMD_API int64_t mtmd_helper_resampler_process(mtmd_helper_resampler * ctx, const float * in, size_t n_in, const float ** out);
	helperResamplerProcessFunc ffi.Fun

	hasResampler bool
)

func loadResamplerFuncs(lib loader.Lib) error {
	var err error

	// Older libraries do not have the resampler API, so it is optional.
	hasResampler = false

	if helperResamplerInitFunc, err = lib.Prep("mtmd_helper_resampler_init", &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypeSint32); err != nil {
		return nil
	}

	if helperResamplerFreeFunc, err = lib.Prep("mtmd_helper_resampler_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return nil
	}

	if helperResamplerProcessFunc, err = lib.Prep("mtmd_helper_resampler_process", &ffi.TypeSint64, &ffi.TypePointer, &ffi.TypePointer, &ffiTypeSize, &ffi.TypePointer); err != nil {
		return nil
	}

	hasResampler = true

	return nil
}

// HelperResamplerInit creates a resampler from rateIn to rateOut, for example to GetAudioSampleRate.
// The caller must call HelperResamplerFree when done.
func HelperResamplerInit(rateIn, rateOut int32) (HelperResampler, error) {
	var r HelperResampler
	if !hasResampler {
		return r, ErrNoResampler
	}
	helperResamplerInitFunc.Call(unsafe.Pointer(&r), &rateIn, &rateOut)
	if r == 0 {
		return r, errors.New("mtmd_helper_resampler_init failed")
	}
	return r, nil
}

// HelperResamplerFree frees a resampler.
func HelperResamplerFree(r HelperResampler) {
	if !hasResampler || r == 0 {
		return
	}
	helperResamplerFreeFunc.Call(nil, unsafe.Pointer(&r))
}

// HelperResamplerProcess resamples in and returns the output samples, which can be empty.
// The slice points to native memory owned by the resampler, valid until the next call or free.
func HelperResamplerProcess(r HelperResampler, in []float32) ([]float32, error) {
	if !hasResampler {
		return nil, ErrNoResampler
	}
	if r == 0 {
		return nil, errors.New("invalid resampler handle")
	}

	inPtr := unsafe.SliceData(in)
	nIn := uint64(len(in))
	var outPtr *float32
	outPtrPtr := &outPtr

	var result ffi.Arg
	helperResamplerProcessFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&r), unsafe.Pointer(&inPtr), &nIn, unsafe.Pointer(&outPtrPtr))

	n := int64(result)
	if n < 0 {
		return nil, errors.New("mtmd_helper_resampler_process failed")
	}
	if n == 0 || outPtr == nil {
		return nil, nil
	}
	return unsafe.Slice(outPtr, int(n)), nil
}
