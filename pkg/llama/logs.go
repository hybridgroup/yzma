package llama

import (
	"runtime"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type LogCallback *ffi.Closure

var (
	// LLAMA_API void llama_log_set(ggml_log_callback log_callback, void * user_data);
	logSetFunc ffi.Fun

	// static void llama_log_callback_null(ggml_log_level level, const char * text, void * user_data) { (void) level; (void) text; (void) user_data; }
	logSilent *ffi.Closure

	callback    unsafe.Pointer
	cifCallback ffi.Cif
	cbFn        uintptr
)

func loadLogFuncs(lib ffi.Lib) error {
	var err error

	if logSetFunc, err = lib.Prep("llama_log_set", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_log_set", err)
	}

	return nil
}

// LogSet sets the logging mode. Pass [LogSilent()] to turn logging off. Pass nil to use stdout.
// Note that if you turn logging off when using the [mtmd] package, you must also set Verbosity > 5.
func LogSet(cb LogCallback, data uintptr) {
	if runtime.GOOS == "darwin" {
		// setting log callback currently causes a SIGBUS: bus error on macOS. so don't do it.
		return
	}
	logSet(cb, data)
}

func logSet(cb LogCallback, data uintptr) {
	logSetFunc.Call(nil, unsafe.Pointer(&cb), unsafe.Pointer(&data))
}

// LogSilent is a callback function that you can pass into the [LogSet] function to turn logging off.
func LogSilent() *ffi.Closure {
	logSilent = ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)

	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 3, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	cbFn = ffi.NewCallback(silentLogCallbackFunc)

	if logSilent != nil {
		if status := ffi.PrepClosureLoc(logSilent, &cifCallback, cbFn, nil, callback); status != ffi.OK {
			panic(status)
		}
	}

	return logSilent
}

func silentLogCallbackFunc(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
	return 0
}
