package llama

import (
	"runtime"
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type LogCallback *ffi.Closure

var (
	// LLAMA_API void llama_log_set(ggml_log_callback log_callback, void * user_data);
	logSet     func(cb LogCallback, data uintptr)
	logSetFunc ffi.Fun

	// static void llama_log_callback_null(ggml_log_level level, const char * text, void * user_data) { (void) level; (void) text; (void) user_data; }
	logSilent *ffi.Closure
)

func loadLogFuncs(lib ffi.Lib) error {
	var err error

	if logSetFunc, err = lib.Prep("llama_log_set", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_log_set", err)
	}
	logSet = func(cb LogCallback, data uintptr) {
		logSetFunc.Call(nil, unsafe.Pointer(&cb), unsafe.Pointer(&data))
	}

	var callback unsafe.Pointer
	logSilent = ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)

	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 0, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		return 0
	})

	if logSilent != nil {
		if status := ffi.PrepClosureLoc(logSilent, &cifCallback, fn, nil, callback); status != ffi.OK {
			panic(status)
		}
	}

	return nil
}

// LogSet sets the logging mode. Pass [LogSilent()] to turn logging off. Pass nil to use stdout.
// Note that if you turn logging off when using the [mtmd] package, you must also set Verbosity > 10.
func LogSet(cb LogCallback, data uintptr) {
	// calling LogSet on macOS currently results in a SIGBUS error.
	// TODO: solve this. For now, you cannot disable logging on macOS.
	if runtime.GOOS == "darwin" {
		return
	}
	logSet(cb, data)
}

// LogSilent is a callback function that you can pass into the [LogSet] function to turn logging off.
func LogSilent() *ffi.Closure {
	return logSilent
}
