package llama

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

type LogCallback uintptr // *ffi.Closure

var (
	// LLAMA_API void llama_log_set(ggml_log_callback log_callback, void * user_data);
	logSetFunc ffi.Fun
)

func loadLogFuncs(lib ffi.Lib) error {
	var err error

	if logSetFunc, err = lib.Prep("llama_log_set", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_log_set", err)
	}

	return nil
}

// LogSet sets the logging mode. Pass llama.LogSilent() to turn logging off. Pass nil to use stdout.
func LogSet(cb uintptr) {
	nada := uintptr(0)
	logSetFunc.Call(nil, unsafe.Pointer(&cb), unsafe.Pointer(&nada))
}

// LogSilent is a callback function that you can pass into the LogSet function to turn logging off.
// The equivalent C function signature is:
//
// static void llama_log_callback_null(ggml_log_level level, const char * text, void * user_data) { (void) level; (void) text; (void) user_data; }
func LogSilent() uintptr {
	cb := func(level int32, text, data uintptr) {}

	var callback unsafe.Pointer
	closure := ffi.ClosureAlloc(unsafe.Sizeof(ffi.Closure{}), &callback)

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		arg := unsafe.Slice(args, cif.NArgs)
		level := *(*int32)(arg[0])
		textPtr := *(*uintptr)(arg[1])
		userDataPtr := *(*uintptr)(arg[2])
		cb(level, textPtr, userDataPtr)

		return 0
	})

	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 3, &ffi.TypeVoid, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); status != ffi.OK {
		return uintptr(0)
	}

	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, callback); status != ffi.OK {
			return uintptr(0)
		}
	}

	return uintptr(callback)
}

// LogNormal is a value you can pass into the LogSet function to turn standard logging on.
const LogNormal uintptr = 0
