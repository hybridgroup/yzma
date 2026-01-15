package mtmd

import (
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/jupiterrider/ffi"
)

// EvalCallback is the Go function type for the evaluation callback.
// tensor is a pointer to the ggml_tensor being evaluated.
// ask indicates whether the scheduler wants to know if the user wants to observe this node.
// When ask is true, return true to indicate you want to observe this node.
// When ask is false, return true to continue computation, or false to cancel.
type EvalCallback func(tensor uintptr, ask bool, userData uintptr) bool

var evalCallback unsafe.Pointer
var sizeOfEvalClosure = unsafe.Sizeof(ffi.Closure{})

// NewEvalCallback creates a callback pointer that can be assigned to ContextParamsType.CbEval.
// The callback will be invoked during model evaluation for each tensor node.
func NewEvalCallback(cb EvalCallback) uintptr {
	if cb == nil {
		return 0
	}

	closure := ffi.ClosureAlloc(sizeOfEvalClosure, &evalCallback)

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		if args == nil || ret == nil {
			*(*uint8)(ret) = 0
			return 1
		}

		arg := unsafe.Slice(args, cif.NArgs)
		tensor := *(*uintptr)(arg[0])
		ask := *(*uint8)(arg[1]) != 0
		userDataPtr := *(*uintptr)(arg[2])

		result := cb(tensor, ask, userDataPtr)
		if result {
			*(*uint8)(ret) = 1
		} else {
			*(*uint8)(ret) = 0
		}
		return 0
	})

	var cifCallback ffi.Cif
	if status := ffi.PrepCif(&cifCallback, ffi.DefaultAbi, 3, &ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeUint8, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	if closure != nil {
		if status := ffi.PrepClosureLoc(closure, &cifCallback, fn, nil, evalCallback); status != ffi.OK {
			panic(status)
		}
	}

	return uintptr(evalCallback)
}

// NewEvalCallbackSimple creates a simple callback using purego that logs or monitors evaluation.
// This is a simpler alternative when you don't need the full ffi.Closure approach.
func NewEvalCallbackSimple(cb func(tensor uintptr, ask bool) bool) uintptr {
	return purego.NewCallback(func(tensor uintptr, ask uint8, userData uintptr) uint8 {
		if cb(tensor, ask != 0) {
			return 1
		}
		return 0
	})
}

// SetEvalCallback sets the evaluation callback on the context params.
func (p *ContextParamsType) SetEvalCallback(cb EvalCallback) {
	p.CbEval = NewEvalCallback(cb)
}

// SetEvalCallbackSimple sets a simple evaluation callback on the context params.
func (p *ContextParamsType) SetEvalCallbackSimple(cb func(tensor uintptr, ask bool) bool) {
	p.CbEval = NewEvalCallbackSimple(cb)
}
