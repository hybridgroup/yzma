// Package bindings is a fixture for the checker's self-test. It is a
// deliberately miniature imitation of yzma's binding style: package-level
// ffi.Fun vars filled in by lib.Prep, called with unsafe.Pointer avalues.
//
// Four of the five bindings below are wrong, one per rule per direction, and
// the fifth is clean. testdata is invisible to the go tool, so nothing here is
// ever built as part of the checker.
package bindings

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

// Mode mirrors a 4-byte enum-like Go type, as yzma's LoadMode is.
type Mode int32

var ffiTypeSize = ffi.TypeUint64

var (
	getThingFunc   ffi.Fun
	descFunc       ffi.Fun
	scoreFunc      ffi.Fun
	modeFromStrFun ffi.Fun
	cleanFunc      ffi.Fun
)

func load(lib ffi.Lib) error {
	var err error

	// RULE 1: C returns a pointer, the cif says void.
	if getThingFunc, err = lib.Prep("fx_get_thing", &ffi.TypeVoid); err != nil {
		return err
	}

	if descFunc, err = lib.Prep("fx_desc", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffiTypeSize); err != nil {
		return err
	}

	if scoreFunc, err = lib.Prep("fx_score", &ffi.TypeFloat, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return err
	}

	if modeFromStrFun, err = lib.Prep("fx_mode_from_str", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return err
	}

	if cleanFunc, err = lib.Prep("fx_clean", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffiTypeSize); err != nil {
		return err
	}

	return nil
}

// GetThing exercises RULE 1 only: the binding is void-returning, so libffi
// writes nothing into ret.
func GetThing() uintptr {
	var ret uintptr
	getThingFunc.Call(unsafe.Pointer(&ret))

	return ret
}

// Desc exercises RULE 2: bLen is 4 bytes behind an 8-byte size_t slot.
func Desc(thing uintptr) string {
	buf := make([]byte, 128)
	b := unsafe.SliceData(buf)
	bLen := int32(len(buf))

	var result ffi.Arg
	descFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&thing), unsafe.Pointer(&b), &bLen)

	return string(buf[:int32(result)])
}

// Score exercises RULE 3 for floats: an ffi.Arg is the wrong return buffer.
func Score(thing uintptr, token int32) float32 {
	var score ffi.Arg
	scoreFunc.Call(unsafe.Pointer(&score), unsafe.Pointer(&thing), unsafe.Pointer(&token))

	return float32(score)
}

// ModeFromStr exercises RULE 3 for integers: a 4-byte return buffer is
// written past, because libffi always stores a full 8-byte ffi_arg.
func ModeFromStr(s *byte) Mode {
	var result Mode
	modeFromStrFun.Call(unsafe.Pointer(&result), unsafe.Pointer(&s))

	return result
}

// Clean violates nothing and must never appear in the report.
func Clean(thing uintptr, a int32, n uint64) int32 {
	var result ffi.Arg
	cleanFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&thing), &a, &n)

	return int32(result)
}

var _ = load
