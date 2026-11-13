package speculative

import (
	"runtime"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/jupiterrider/ffi"
)

// SetEmbeddingsNextNDev sets whether the context writes its unmasked NextN
// rows to a device buffer instead of host memory. Call [SetEmbeddingsNextN]
// with masked false first.
//
// The buffer has n_batch + n_seq_max rows of [llama.ModelNEmbdOut] floats. Row
// r below n_batch is batch position r of the last decode. Row n_batch + s is
// free for the caller, for example for sequence s.
//
// It returns false when the device buffer cannot be used or the library does
// not export the function. The rows then go to host memory as before.
func SetEmbeddingsNextNDev(ctx llama.Context, value bool) bool {
	if ctx == 0 || setEmbeddingsNextNDevFunc.Cif == nil {
		return false
	}

	var result ffi.Arg
	setEmbeddingsNextNDevFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &value)

	return result.Bool()
}

// GetEmbeddingsNextNRows copies device rows starting at row0 into out. nEmbd
// is the row width, [llama.ModelNEmbdOut]. The number of rows is len(out)/nEmbd.
//
// The copy waits for the last decode. It returns false when out is not a whole
// number of rows, the rows are out of range, no device buffer is in use, or the
// library does not export the function.
func GetEmbeddingsNextNRows(ctx llama.Context, row0 int64, nEmbd int, out []float32) bool {
	if ctx == 0 || getEmbeddingsNextNRowsFunc.Cif == nil {
		return false
	}

	nRows, ok := rowCount(len(out), nEmbd)
	if !ok {
		return false
	}

	ptr := unsafe.SliceData(out)
	var result ffi.Arg
	getEmbeddingsNextNRowsFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &row0, &nRows, unsafe.Pointer(&ptr))
	runtime.KeepAlive(out)

	return result.Bool()
}

// SetEmbeddingsNextNRows copies data into the device rows starting at row0.
// nEmbd is the row width, [llama.ModelNEmbdOut]. The number of rows is
// len(data)/nEmbd. llama.cpp does not keep data after the call.
//
// It returns false in the same cases as [GetEmbeddingsNextNRows].
func SetEmbeddingsNextNRows(ctx llama.Context, row0 int64, nEmbd int, data []float32) bool {
	if ctx == 0 || setEmbeddingsNextNRowsFunc.Cif == nil {
		return false
	}

	nRows, ok := rowCount(len(data), nEmbd)
	if !ok {
		return false
	}

	ptr := unsafe.SliceData(data)
	var result ffi.Arg
	setEmbeddingsNextNRowsFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), &row0, &nRows, unsafe.Pointer(&ptr))
	runtime.KeepAlive(data)

	return result.Bool()
}

// SetCtxStateSrc makes the MTP draft ctx read batch state rows from the NextN
// device buffer of ctxSrc. It returns whether decode accepts batches with state rows.
//
// A ctxSrc of 0 detaches. It returns false then because no source is left,
// which is not a failed detach. Keep ctxSrc alive while ctx uses it, and detach before
// freeing either context.
func SetCtxStateSrc(ctx, ctxSrc llama.Context) bool {
	if ctx == 0 || setCtxStateSrcFunc.Cif == nil {
		return false
	}

	var result ffi.Arg
	setCtxStateSrcFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&ctx), unsafe.Pointer(&ctxSrc))

	return result.Bool()
}

// BatchExtSetEmbdStateRow sets the state of the entry at idx to a row of the
// device buffer set with [SetCtxStateSrc]. Every entry of a batch must have a
// state row, or none.
//
// It returns false when llama.cpp rejects the index or row, or the library
// does not export the function.
func BatchExtSetEmbdStateRow(batch llama.BatchExt, idx, row int32) bool {
	if batch == 0 || batchExtSetEmbdStateRowFunc.Cif == nil {
		return false
	}

	var result ffi.Arg
	batchExtSetEmbdStateRowFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&batch), &idx, &row)

	return result.Bool()
}

func rowCount(n, nEmbd int) (int64, bool) {
	if nEmbd <= 0 || n == 0 || n%nEmbd != 0 {
		return 0, false
	}

	return int64(n / nEmbd), true
}
