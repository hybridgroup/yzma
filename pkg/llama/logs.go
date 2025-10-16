package llama

import (
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

type LogCallback uintptr // *ffi.Closure

var (
	// LLAMA_API void llama_log_set(ggml_log_callback log_callback, void * user_data);
	logSetFunc ffi.Fun

	// static void llama_log_callback_null(ggml_log_level level, const char * text, void * user_data) { (void) level; (void) text; (void) user_data; }
	// logSilent uintptr
)

func loadLogFuncs(lib ffi.Lib) error {
	var err error

	if logSetFunc, err = lib.Prep("llama_log_set", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_log_set", err)
	}

	return nil
}

// LogSet sets the logging mode. Pass [LogSilent()] to turn logging off. Pass nil to use stdout.
// Note that if you turn logging off when using the [mtmd] package, you must also set Verbosity = llama.LogLevelContinue.
func LogSet(cb uintptr) {
	nada := uintptr(0)
	logSetFunc.Call(nil, unsafe.Pointer(&cb), unsafe.Pointer(&nada))
}

// LogSilent is a callback function that you can pass into the [LogSet] function to turn logging off.
func LogSilent() uintptr {
	return purego.NewCallback(func(level int32, text, data uintptr) uintptr {
		return 0
	})
}
