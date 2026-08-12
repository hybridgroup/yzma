// Package bindings is a fixture for the checker's self-test. It is a
// deliberately miniature imitation of yzma's binding style: package-level
// ffi.Fun vars filled in by lib.Prep, called with unsafe.Pointer avalues.
//
// Eight of the twelve bindings below are wrong, one per rule per direction plus
// one per struct-layout direction plus one transposition, and the other four
// are clean. The mirrored constants at the end carry the ninth plant, for RULE
// 4, among ten clean controls. testdata is invisible to the go tool, so nothing
// here is ever built as part of the checker.
package bindings

import (
	"unsafe"

	"github.com/jupiterrider/ffi"
)

// Mode mirrors a 4-byte enum-like Go type, as yzma's LoadMode is.
type Mode int32

var ffiTypeSize = ffi.TypeUint64

// ffiTypeParams is missing the second uint32 of struct fx_params. The gap is
// reabsorbed by the padding in front of the pointer member, so this descriptor
// is 24 bytes just like the C struct: the RULE 1 struct-layout plant.
var ffiTypeParams = ffi.NewType(
	&ffi.TypeUint32,
	&ffi.TypeUint32,
	&ffi.TypeFloat,
	&ffi.TypePointer)

// ffiTypeGeom matches struct fx_geom exactly. The plant for that binding is on
// the Go side, in Geom.
var ffiTypeGeom = ffi.NewType(
	&ffi.TypeUint32,
	&ffi.TypeUint32,
	&ffi.TypeFloat,
	&ffi.TypePointer)

// Params mirrors struct fx_params the way yzma's ContextParams mirrored
// llama_context_params before yzma#289: the Go struct and the cif descriptor
// are missing the same member, so they agree with each other and only the
// comparison against C can see the drift.
type Params struct {
	NA    uint32
	NC    uint32
	Scale float32
	UD    uintptr
}

// Geom is the RULE 3 struct-layout plant: 24 bytes like the C struct and like
// ffiTypeGeom, but S sits where fx_geom.h belongs.
type Geom struct {
	W  uint32
	S  float32
	H  uint32
	UD uintptr
}

// GeomOK is the clean control: member for member, struct fx_geom.
type GeomOK struct {
	W  uint32
	H  uint32
	S  float32
	UD uintptr
}

// GeomTail is the second clean control, in the shape of llama.Batch: the FFI
// members followed by Go-only bookkeeping past the end of the C struct. libffi
// reads only the 24 bytes the descriptor declares, so the tail is not a
// finding.
type GeomTail struct {
	W  uint32
	H  uint32
	S  float32
	UD uintptr

	capA int32
	capB int32
}

// ffiTypePair matches struct fx_pair, and so does PairSwapped byte for byte:
// the only evidence of the defect is in the names.
var ffiTypePair = ffi.NewType(&ffi.TypeUint32, &ffi.TypeUint32)

// PairSwapped is the transposition plant: alpha and beta are exchanged, which
// no offset, width or ABI class comparison can see.
type PairSwapped struct {
	BetaCount  uint32
	AlphaCount uint32
}

// NamedOK is the clean control for the alias table: n_attn_heads spelled
// AttentionHeads and ctx_size spelled ContextSize are the same members, not a
// transposition.
type NamedOK struct {
	AttentionHeads uint32
	ContextSize    uint32
}

// GeomShort is the other side of that tolerance: 12 bytes behind a 24-byte
// descriptor, so libffi reads 12 bytes of whatever follows it.
type GeomShort struct {
	W uint32
	H uint32
	S float32
}

// FxLevel mirrors enum llama_fx_level and carries the RULE 4 plant:
// FxLevelHigh still holds the value the member had before LLAMA_FX_LEVEL_LOW
// was inserted in front of it, so every call passing FxLevelHigh asks C for LOW
// instead. Both sides compile, and every width still matches.
type FxLevel int32

const (
	// LLAMA_FX_LEVEL_OFF
	FxLevelOff FxLevel = 0
	// FxLevelLow carries no comment, so it is matched by name normalisation
	// against LLAMA_FX_LEVEL_LOW.
	FxLevelLow FxLevel = 1
	// LLAMA_FX_LEVEL_HIGH
	FxLevelHigh FxLevel = 1
	// LLAMA_FX_LEVEL_MAX
	FxLevelMax FxLevel = 2
)

// FxFlag mirrors enum llama_fx_flag. All four are clean controls, one per
// initialiser form the C parser has to evaluate, and none of them carries a
// comment: they are matched by normalising the names.
type FxFlag int32

const (
	FxFlagAuto FxFlag = -1
	FxFlagNone FxFlag = 0
	FxFlagA    FxFlag = 1 << 2
	FxFlagB    FxFlag = 0x10
)

// The #define controls, all clean. FxMagicAlias is defined from FxMagic on both
// sides.
const (
	FxMagic      = 0x66780001
	FxVersion    = 3
	FxMagicAlias = FxMagic
)

var (
	getThingFunc   ffi.Fun
	descFunc       ffi.Fun
	scoreFunc      ffi.Fun
	modeFromStrFun ffi.Fun
	cleanFunc      ffi.Fun
	paramsDefault  ffi.Fun
	geomDefault    ffi.Fun
	useGeomFunc    ffi.Fun
	useTailFunc    ffi.Fun
	useShortFunc   ffi.Fun
	usePairFunc    ffi.Fun
	useNamedFunc   ffi.Fun
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

	if paramsDefault, err = lib.Prep("fx_params_default", &ffiTypeParams); err != nil {
		return err
	}

	if geomDefault, err = lib.Prep("fx_geom_default", &ffiTypeGeom); err != nil {
		return err
	}

	if useGeomFunc, err = lib.Prep("fx_use_geom", &ffi.TypeVoid, &ffiTypeGeom); err != nil {
		return err
	}

	if useTailFunc, err = lib.Prep("fx_use_geom_tail", &ffi.TypeVoid, &ffiTypeGeom); err != nil {
		return err
	}

	if useShortFunc, err = lib.Prep("fx_use_geom_short", &ffi.TypeVoid, &ffiTypeGeom); err != nil {
		return err
	}

	if usePairFunc, err = lib.Prep("fx_use_pair", &ffi.TypeVoid, &ffiTypePair); err != nil {
		return err
	}

	if useNamedFunc, err = lib.Prep("fx_use_named", &ffi.TypeVoid, &ffiTypePair); err != nil {
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

// ParamsDefault exercises the RULE 1 struct-layout plant: ffiTypeParams and
// Params agree with each other and disagree with struct fx_params.
func ParamsDefault() Params {
	var p Params
	paramsDefault.Call(unsafe.Pointer(&p))

	return p
}

// GeomDefault exercises the RULE 3 struct-layout plant: ffiTypeGeom is right,
// Geom is not.
func GeomDefault() Geom {
	var g Geom
	geomDefault.Call(unsafe.Pointer(&g))

	return g
}

// UseGeom is the clean struct-by-value control and must never be reported.
func UseGeom(g GeomOK) {
	useGeomFunc.Call(nil, unsafe.Pointer(&g))
}

// UseGeomTail passes a struct with a Go-only tail and must never be reported.
func UseGeomTail(g GeomTail) {
	useTailFunc.Call(nil, unsafe.Pointer(&g))
}

// UseGeomShort exercises RULE 2 for struct-by-value: the Go struct is shorter
// than the bytes libffi reads.
func UseGeomShort(g GeomShort) {
	useShortFunc.Call(nil, unsafe.Pointer(&g))
}

// UsePair exercises the transposition check: byte-identical to struct fx_pair,
// with the two members exchanged.
func UsePair(p PairSwapped) {
	usePairFunc.Call(nil, unsafe.Pointer(&p))
}

// UseNamed is the alias-table control and must never be reported.
func UseNamed(n NamedOK) {
	useNamedFunc.Call(nil, unsafe.Pointer(&n))
}

// Clean violates nothing and must never appear in the report.
func Clean(thing uintptr, a int32, n uint64) int32 {
	var result ffi.Arg
	cleanFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&thing), &a, &n)

	return int32(result)
}

var _ = load
