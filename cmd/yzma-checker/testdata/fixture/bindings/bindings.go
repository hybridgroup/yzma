// Package bindings is a fixture for the checker's self-test. It is a
// deliberately miniature imitation of yzma's binding style: package-level
// ffi.Fun vars filled in by lib.Prep, called with unsafe.Pointer avalues.
//
// Thirteen of the twenty-six bindings below are wrong, one per rule per
// direction plus one per struct-layout direction plus one transposition plus one
// variadic nfixed plus one pointer aimed at the wrong target plus one C string
// with no terminator plus two function-pointer struct members, and eleven are
// clean. Two
// more carry the signedness plants, which are deliberately not violations: same
// width, same register, other meaning. The mirrored constants in the middle
// carry the RULE 4 plant among ten clean controls, and the four callback sites
// at the end carry the two RULE 5 plants, with one clean control per form. The
// three struct fx_hooks mirrors carry the other two RULE 5 plants, in what a
// function-pointer member *holds* rather than in a signature.
// testdata is invisible to the go tool, so nothing here is ever built as part of
// the checker.
package bindings

import (
	"unsafe"

	"github.com/ebitengine/purego"
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

// ffiTypeHooks matches struct fx_hooks exactly, and so do all three Go structs
// below: every plant here is in what the function-pointer members *hold*, which
// no offset, width or ABI class comparison can see.
var ffiTypeHooks = ffi.NewType(
	&ffi.TypeUint32,
	&ffi.TypePointer,
	&ffi.TypePointer,
	&ffi.TypePointer)

// HooksFunc is the func-value plant: cb_progress is a function pointer C jumps
// through, and a Go func value is a pointer to a func *descriptor*. It is 8
// bytes of pointer class, so the layout comparison passes it and C jumps into
// the descriptor's first word.
type HooksFunc struct {
	NSlots     uint32
	CbProgress func(progress float32, userData uintptr) bool
	CbAbort    uintptr
	UserData   uintptr
}

// Hooks is the wrong-code-pointer plant, in the shape yzma declares these
// members: uintptr on both, so nothing in the Go types says which callback
// belongs where. SetHooks fills them in.
type Hooks struct {
	NSlots     uint32
	CbProgress uintptr
	CbAbort    uintptr
	UserData   uintptr
}

// HooksClean is the control, and the half that matters: both members hold the
// code pointer of the callback that actually implements them, written exactly as
// yzma's ModelParams.SetProgressCallback writes its own. A rule that reported
// this would report that too.
type HooksClean struct {
	NSlots     uint32
	CbProgress uintptr
	CbAbort    uintptr
	UserData   uintptr
}

// SetHooks carries one store of each kind. CbProgress gets the closure built for
// llama_fx_progress_callback, which is what C calls that member as. CbAbort gets
// the one built for llama_fx_report_callback, so C unpacks a
// llama_fx_abort_callback's single pointer argument through a descriptor
// declaring a float and a pointer - and every width, offset and class still
// matches.
func (h *Hooks) SetHooks() {
	h.CbProgress = uintptr(fxProgressCallbackCode)
	h.CbAbort = uintptr(fxReportCallbackCode)
}

// SetHooksClean is the control for both callback forms at once: a descriptor
// closure's code address for the descriptor-form member, and the purego
// constructor's result for the purego-form one.
func (h *HooksClean) SetHooksClean() {
	h.CbProgress = uintptr(fxProgressCallbackCode)
	h.CbAbort = newFxAbortCallback()
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

// FxSplitMode mirrors one member of enum llama_fx_split_mode and leaves LAYER
// and TENSOR alone, the way pkg/llama mirrors llama_split_mode. The mirrored
// member is a clean RULE 4 comparison; the two it omits are the inventory's
// control, and must stay an inventory entry rather than a finding.
type FxSplitMode int32

const (
	// LLAMA_FX_SPLIT_MODE_NONE
	FxSplitModeNone FxSplitMode = 0
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
	useHooksFunc   ffi.Fun
	useHooksOKFunc ffi.Fun
	useHooksCleanF ffi.Fun
	logfFunc       ffi.Fun
	printfFunc     ffi.Fun
	getScoresFunc  ffi.Fun
	getLogitsFunc  ffi.Fun
	getCountFunc   ffi.Fun
	getTokenFunc   ffi.Fun
	decodeFunc     ffi.Fun
	decodeOKFunc   ffi.Fun
	setNameFunc    ffi.Fun
	setPathFunc    ffi.Fun
	setTextFunc    ffi.Fun
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

	if useHooksFunc, err = lib.Prep("fx_use_hooks", &ffi.TypeVoid, &ffiTypeHooks); err != nil {
		return err
	}

	if useHooksOKFunc, err = lib.Prep("fx_use_hooks_ok", &ffi.TypeVoid, &ffiTypeHooks); err != nil {
		return err
	}

	if useHooksCleanF, err = lib.Prep("fx_use_hooks_clean", &ffi.TypeVoid, &ffiTypeHooks); err != nil {
		return err
	}

	// The variadic plant: fx_logf declares two parameters before its "...", and
	// this says one. Every type in the list is right, so nothing but nfixed
	// carries the defect - and on Apple arm64 it decides register versus stack
	// for fmt and everything after it.
	if logfFunc, err = lib.PrepVar("fx_logf", 1, &ffi.TypeVoid,
		&ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return err
	}

	// The variadic control: same C shape, correct nfixed, one concrete variadic
	// type after the two fixed ones. It must never be reported.
	if printfFunc, err = lib.PrepVar("fx_printf", 2, &ffi.TypeSint32,
		&ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32); err != nil {
		return err
	}

	// The pointer-target group. Every descriptor here is &ffi.TypePointer and
	// every avalue is 8 bytes, so nothing in the cif or in the widths can carry
	// the defect: it is one indirection down, in what the pointer aims at.
	if getScoresFunc, err = lib.Prep("fx_get_scores", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if getLogitsFunc, err = lib.Prep("fx_get_logits", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if getCountFunc, err = lib.Prep("fx_get_count", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if getTokenFunc, err = lib.Prep("fx_get_token", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if decodeFunc, err = lib.Prep("fx_decode", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return err
	}

	if decodeOKFunc, err = lib.Prep("fx_decode_ok", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return err
	}

	// The C-string group. Every descriptor is &ffi.TypePointer and every buffer
	// is a Go *byte, so the only thing that separates the plant from its two
	// controls is whether the bytes end in a NUL.
	if setNameFunc, err = lib.Prep("fx_set_name", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if setPathFunc, err = lib.Prep("fx_set_path", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return err
	}

	if setTextFunc, err = lib.Prep("fx_set_text", &ffi.TypeVoid, &ffi.TypePointer, &ffi.TypePointer); err != nil {
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

// UseHooksFunc hands over the struct whose cb_progress is a Go func value.
func UseHooksFunc(h HooksFunc) {
	useHooksFunc.Call(nil, unsafe.Pointer(&h))
}

// UseHooks hands over the struct SetHooks filled in, one member of it with the
// wrong closure's code pointer.
func UseHooks(h Hooks) {
	useHooksOKFunc.Call(nil, unsafe.Pointer(&h))
}

// UseHooksClean is the control and must never be reported.
func UseHooksClean(h HooksClean) {
	useHooksCleanF.Call(nil, unsafe.Pointer(&h))
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

// Logf exercises the nfixed plant. Every avalue is the width the cif claims, so
// RULES 2 and 3 have nothing to say about it.
func Logf(thing uintptr, fmtStr *byte, n int32) {
	logfFunc.Call(nil, unsafe.Pointer(&thing), unsafe.Pointer(&fmtStr), unsafe.Pointer(&n))
}

// Printf is the variadic control and must never be reported.
func Printf(thing uintptr, fmtStr *byte, n int32) int32 {
	var result ffi.Arg
	printfFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&thing), unsafe.Pointer(&fmtStr), unsafe.Pointer(&n))

	return int32(result)
}

// Clean violates nothing and must never appear in the report.
func Clean(thing uintptr, a int32, n uint64) int32 {
	var result ffi.Arg
	cleanFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&thing), &a, &n)

	return int32(result)
}

// GetScores is the pointer-target plant: fx_get_scores writes float32s and this
// hands it a *int32, so every value Go reads back is a float's bit pattern read
// as an integer. The slot is a pointer on both sides and both avalues are 8
// bytes, so no width, class or layout comparison can see it.
func GetScores(thing uintptr, out *int32) {
	getScoresFunc.Call(nil, unsafe.Pointer(&thing), unsafe.Pointer(&out))
}

// GetLogits is the control for it and must never be reported: the same slot,
// pointing at the float32 the header declares.
func GetLogits(thing uintptr, out *float32) {
	getLogitsFunc.Call(nil, unsafe.Pointer(&thing), unsafe.Pointer(&out))
}

// GetCount is the pointer-target signedness plant. Signed and unsigned 4-byte
// integers are the same bytes in the same register, so this is not an ABI break
// and must not be a violation - but the negative values C writes come back as
// very large positive ones.
func GetCount(thing uintptr, out *uint32) {
	getCountFunc.Call(nil, unsafe.Pointer(&thing), unsafe.Pointer(&out))
}

// GetToken is the control for that: a *int32 that can hold the -1 sentinel.
func GetToken(thing uintptr, out *int32) {
	getTokenFunc.Call(nil, unsafe.Pointer(&thing), unsafe.Pointer(&out))
}

// Decode is the return-buffer signedness plant. The buffer is a full 8 bytes, so
// RULE 3 is satisfied, and fx_decode's negative error code is read out of it as
// 4294967295.
func Decode(thing uintptr) uint64 {
	var result uint64
	decodeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&thing))

	return result
}

// DecodeOK is the control: the same declaration, the same width, read through a
// signed buffer.
func DecodeOK(thing uintptr) int64 {
	var result int64
	decodeOKFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&thing))

	return result
}

// SetName is the NUL-termination plant: the buffer is a Go string copied into a
// []byte and nothing appends a terminator, so C keeps reading past the end of
// that allocation until it happens to find a zero byte. The slot is a pointer of
// the same width on both sides, aimed at a char through a *byte, so every other
// rule passes it.
func SetName(thing uintptr, name string) {
	n := &[]byte(name)[0]
	setNameFunc.Call(nil, unsafe.Pointer(&thing), unsafe.Pointer(&n))
}

// SetPath is the control for it, written the way every yzma string site is, and
// must never be reported.
func SetPath(thing uintptr, path string) {
	p := &[]byte(path + "\x00")[0]
	setPathFunc.Call(nil, unsafe.Pointer(&thing), unsafe.Pointer(&p))
}

// SetText is the second control, in the other idiom yzma uses: the string is
// terminated in place and handed over with unsafe.StringData.
func SetText(thing uintptr, text string) {
	text += "\x00"
	p := unsafe.StringData(text)
	setTextFunc.Call(nil, unsafe.Pointer(&thing), unsafe.Pointer(&p))
}

var _ = load

// The RULE 5 sites: C calling back into Go. None of them goes through lib.Prep,
// so none of them is one of the bindings above.

var (
	fxProgressCallbackCode unsafe.Pointer
	fxProgressCallbackCif  *ffi.Cif
	fxReportCallbackCode   unsafe.Pointer
	fxReportCallbackCif    *ffi.Cif
	sizeOfClosure          = unsafe.Sizeof(ffi.Closure{})
)

// setFxProgressCallback is the clean control for the descriptor form: uint8
// return for bool, float for float, pointer for void *, and nfixed matching the
// two parameters llama_fx_progress_callback declares.
func setFxProgressCallback() uintptr {
	closure := ffi.ClosureAlloc(sizeOfClosure, &fxProgressCallbackCode)

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		arg := unsafe.Slice(args, cif.NArgs)
		*(*uint8)(ret) = uint8(*(*float32)(arg[0]))

		return 0
	})

	fxProgressCallbackCif = new(ffi.Cif)
	if status := ffi.PrepCif(fxProgressCallbackCif, ffi.DefaultAbi, 2, &ffi.TypeUint8, &ffi.TypeFloat, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	if status := ffi.PrepClosureLoc(closure, fxProgressCallbackCif, fn, nil, fxProgressCallbackCode); status != ffi.OK {
		panic(status)
	}

	return uintptr(fxProgressCallbackCode)
}

// setFxReportCallback is the descriptor plant: llama_fx_report_callback passes a
// float, and the descriptor claims a 4-byte int. The width matches, so libffi
// reads the right bytes and the closure then reads a float's bit pattern as an
// integer on every invocation.
func setFxReportCallback() uintptr {
	closure := ffi.ClosureAlloc(sizeOfClosure, &fxReportCallbackCode)

	fn := ffi.NewCallback(func(cif *ffi.Cif, ret unsafe.Pointer, args *unsafe.Pointer, userData unsafe.Pointer) uintptr {
		arg := unsafe.Slice(args, cif.NArgs)
		*(*uint8)(ret) = uint8(*(*int32)(arg[0]))

		return 0
	})

	fxReportCallbackCif = new(ffi.Cif)
	if status := ffi.PrepCif(fxReportCallbackCif, ffi.DefaultAbi, 2, &ffi.TypeUint8, &ffi.TypeSint32, &ffi.TypePointer); status != ffi.OK {
		panic(status)
	}

	if status := ffi.PrepClosureLoc(closure, fxReportCallbackCif, fn, nil, fxReportCallbackCode); status != ffi.OK {
		panic(status)
	}

	return uintptr(fxReportCallbackCode)
}

// newFxLogCallback is the purego plant: llama_fx_log_callback passes three
// arguments and this closure declares two, so text arrives where level belongs
// and data is never read at all.
func newFxLogCallback() uintptr {
	return purego.NewCallback(func(level int32, text uintptr) uintptr {
		return 0
	})
}

// newFxAbortCallback is the clean control for that form: one pointer parameter
// for llama_fx_abort_callback's void *, and a word-sized result C reads the low
// byte of as a bool.
func newFxAbortCallback() uintptr {
	return purego.NewCallback(func(data uintptr) uintptr {
		return 1
	})
}
