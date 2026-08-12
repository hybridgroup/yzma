package main

import (
	"slices"
	"strings"
	"testing"
)

// The checker's own correctness gate.
//
// It used to be "a run must re-derive these two live defects in yzma", which
// stopped working the moment the defects were fixed: a tool that reports
// nothing is indistinguishable from a tool that checks nothing. Instead the
// gate now runs the full pipeline over testdata/fixture, a miniature binding
// package carrying one planted defect per rule per direction and per callback
// form, plus clean
// controls. Both halves matter — every plant must be found, and the clean
// bindings and constants must never be reported.

func fixtureReport(t *testing.T) *report {
	t.Helper()

	rep, err := analyse("testdata/fixture", "testdata/fixture/hdrs", []string{"yzma-checker-fixture/bindings"})
	if err != nil {
		t.Fatalf("analyse fixture: %v", err)
	}

	return rep
}

func TestFixtureFindsEveryPlantedDefect(t *testing.T) {
	rep := fixtureReport(t)

	tests := []struct {
		name  string
		rule  int
		fn    string
		match string
	}{
		{
			name: "rule1 pointer return bound as void",
			rule: 1, fn: "fx_get_thing",
			match: "cif says TypeVoid",
		},
		{
			name: "rule2 size_t slot fed a 4-byte Go value",
			rule: 2, fn: "fx_desc",
			match: "wants 8B, Go int32 is 4B",
		},
		{
			name: "rule3 float return read through ffi.Arg",
			rule: 3, fn: "fx_score",
			match: "must be float32",
		},
		{
			name: "rule3 integer return into a 4-byte buffer",
			rule: 3, fn: "fx_mode_from_str",
			match: "must be ffi.Arg",
		},
		// The struct-layout plants. Every struct here is the same size on every
		// side, so a size-only comparison passes all of them: this is the
		// yzma#289 gap, in each direction.
		{
			name: "rule1 cif descriptor missing a member absorbed by padding",
			rule: 1, fn: "fx_params_default",
			match: "member 2: cif ffiTypeParams +8 float/4B e2 vs C struct fx_params +8 uint/4B n_c",
		},
		{
			name: "rule2 Go struct shorter than the bytes libffi reads",
			rule: 2, fn: "fx_use_geom_short",
			match: "cif ffiTypeGeom wants 24B, Go yzma-checker-fixture/bindings.GeomShort is 12B",
		},
		{
			// Offsets, widths and ABI classes are all identical here: the members
			// are the same two uint32_t in the other order, so only the names
			// carry the defect - and the C library receives each value in the
			// other's place.
			name: "rule2 same-class members transposed on the Go side",
			rule: 2, fn: "fx_use_pair",
			match: "PairSwapped.BetaCount is C member 1 (struct fx_pair.beta_count)",
		},
		{
			name: "rule3 Go struct members ordered differently from the descriptor",
			rule: 3, fn: "fx_geom_default",
			match: "member 1: Go yzma-checker-fixture/bindings.Geom +4 float/4B S vs cif ffiTypeGeom +4 uint/4B e1",
		},
		{
			// The variadic plant. Every type in the descriptor is right, so
			// only nfixed carries it - and on Apple arm64 a variadic argument
			// travels on the stack where a fixed one travels in a register, so
			// getting the boundary wrong misplaces every variadic argument. A
			// variadic binding used to be exempt from arity checking entirely.
			name: "rule1 PrepVar declares too few fixed arguments",
			rule: 1, fn: "fx_logf",
			match: `nfixed is 1 but C fx_logf declares 2 parameter(s) before its "..."`,
		},
		{
			// The constant plant: LLAMA_FX_LEVEL_LOW was inserted in front of
			// HIGH upstream, so HIGH is 2 in the header and still 1 in Go. No
			// width, offset or class changed, which is why RULES 1-3 cannot see
			// it.
			name: "rule4 enum member value shifted by an insertion upstream",
			rule: 4, fn: "LLAMA_FX_LEVEL_HIGH",
			match: "FxLevelHigh = 1 but C LLAMA_FX_LEVEL_HIGH = 2",
		},
		{
			// The pointer-target plant. The slot is a pointer on both sides and
			// both avalues are 8 bytes, so every rule passes it; the defect is one
			// indirection down, where C writes float bit patterns into memory Go
			// reads as integers.
			name: "rule2 pointer slot aimed at the wrong target type",
			rule: 2, fn: "fx_get_scores",
			match: "arg1: C float * points at float but Go *int32 points at sint",
		},
		{
			// The NUL-termination plant. The slot is a pointer to a char on both
			// sides and the buffer is the *byte every string site in the tree
			// passes, so every width, class and pointer-target comparison passes
			// it: what is missing is the one byte C needs to find the end of the
			// string, and Go never puts it there by itself.
			name: "rule2 char* buffer built from a Go string with no terminator",
			rule: 2, fn: "fx_set_name",
			match: "arg1: C const char * needs a NUL-terminated buffer but unsafe.Pointer(&n) is built from []byte(name) with no terminator",
		},
		{
			// The enum-parameter plant. The slot is 4 bytes of sint on the C side,
			// the descriptor side and the Go side, and there is no pointer target
			// and no member to transpose, so every rule passes it: the value simply
			// belongs to another enumeration than the one C is indexing.
			name: "rule2 enum parameter passed a value of another enumeration",
			rule: 2, fn: "fx_set_level",
			match: "arg1: C takes enum llama_fx_level but Go yzma-checker-fixture/bindings.FxSplitMode mirrors enum llama_fx_split_mode",
		},
		{
			// The callback plants, in the direction where C calls Go. The width
			// is right here and the class is not, so libffi reads the correct
			// four bytes and the closure then reads a float's bit pattern as an
			// integer - on every invocation, for the life of the process.
			name: "rule5 closure descriptor declares a float parameter as an int",
			rule: 5, fn: "llama_fx_report_callback",
			match: "llama_fx_report_callback arg0: C float is float but cif TypeSint32 is sint",
		},
		{
			// A missing parameter shifts every later one, so this is reported as
			// loudly as a width error: text arrives where level belongs.
			name: "rule5 purego closure one parameter short of the C typedef",
			rule: 5, fn: "llama_fx_log_callback",
			match: "closure takes 2 parameter(s) but C llama_fx_log_callback passes 3",
		},
		{
			// The function-pointer struct member, the third way C reaches a Go
			// callback. A Go func value is 8 bytes of pointer class exactly like
			// the code pointer that belongs there, so the layout comparison passes
			// it - and it points at a func descriptor rather than at code.
			name: "rule5 struct member holding a Go func value instead of a code pointer",
			rule: 5, fn: "fx_hooks.cb_progress",
			match: "Go func(progress float32, userData uintptr) bool is a func value",
		},
		{
			// The other half: the member is the uintptr yzma declares, and the
			// code pointer written into it was built for a different typedef. Both
			// sides are 8 bytes of pointer class, both structs lay out identically,
			// and C unpacks its arguments through the wrong descriptor for the life
			// of the process.
			name: "rule5 struct member holding the wrong callback's code pointer",
			rule: 5, fn: "fx_hooks.cb_abort",
			match: "implements llama_fx_report_callback, but C calls this member as llama_fx_abort_callback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range rep.Viols {
				if v.Rule == tt.rule && v.Fn == tt.fn && strings.Contains(v.Detail, tt.match) {
					return
				}
			}

			t.Errorf("planted RULE%d defect in %s was not reported (want detail containing %q)\ngot:\n%s",
				tt.rule, tt.fn, tt.match, dumpViolations(rep))
		})
	}
}

func TestFixtureDoesNotReportTheCleanBinding(t *testing.T) {
	rep := fixtureReport(t)

	// fx_clean is the scalar control and fx_use_geom the struct-by-value one: a
	// struct whose C, cif and Go members all agree must survive the layout
	// comparison silently, or the check is just noise. fx_use_geom_tail is the
	// llama.Batch shape - Go-only bookkeeping appended past the end of the C
	// struct, which libffi never reads - and pins that it stays not-a-finding.
	// fx_use_named is the alias table's control: members whose Go and C spellings
	// differ by more than case are the same members, not a transposition.
	// The RULE 4 controls are the constants around the plant: the three other
	// members of the same enum, every initialiser form of llama_fx_flag, and the
	// three #defines. A rule that reported those would be worse than useless,
	// because every mirrored constant in yzma looks like them.
	// The RULE 5 controls are one per callback form: a descriptor that models its
	// typedef exactly, and a purego closure whose parameters and word-sized
	// result match. Both are written the way the four real yzma sites are, so a
	// rule that reported them would report all of those too.
	clean := map[string]bool{
		"fx_clean": true, "fx_use_geom": true, "fx_use_geom_tail": true, "fx_use_named": true,
		// fx_printf is the variadic control - correct nfixed, so nothing about
		// it is a finding - and fx_unbound is the coverage-inventory control: a
		// C declaration nothing binds is not a defect, and the assertion that
		// matters for it is this one, not its presence in the inventory.
		"fx_printf": true, "fx_unbound": true,
		// The pointer-target controls: fx_get_logits points at the float32 its
		// header declares, and fx_get_token at the int32 that can hold its -1
		// sentinel. fx_get_count and fx_decode carry the two signedness plants
		// and are on this list on purpose: signed against unsigned is the same
		// register and the same bytes, so it must be reported as its own class
		// and never as an ABI violation. fx_decode_ok is its control.
		"fx_get_logits": true, "fx_get_token": true, "fx_get_count": true,
		"fx_decode": true, "fx_decode_ok": true,
		// The NUL-termination controls, one per idiom yzma actually uses. This is
		// the half that matters: both are written exactly as the six real string
		// sites in the tree are, so a rule that reported either would report all
		// of them.
		"fx_set_path": true, "fx_set_text": true,
		// The enum-parameter controls: fx_set_flag is fed the yzma type that
		// mirrors its own C enum, which is how every enum-taking call in the tree is
		// written, and fx_set_mode is fed a plain int32, which names no enum and so
		// must be neither a finding nor a skip.
		"fx_set_flag": true, "fx_set_mode": true,
		"LLAMA_FX_LEVEL_OFF": true, "LLAMA_FX_LEVEL_LOW": true, "LLAMA_FX_LEVEL_MAX": true,
		"LLAMA_FX_FLAG_AUTO": true, "LLAMA_FX_FLAG_NONE": true,
		"LLAMA_FX_FLAG_A": true, "LLAMA_FX_FLAG_B": true,
		"LLAMA_FX_MAGIC": true, "LLAMA_FX_VERSION": true, "LLAMA_FX_MAGIC_ALIAS": true,
		"llama_fx_progress_callback": true, "llama_fx_abort_callback": true,
		// The partially-mirrored enum: the member yzma mirrors is a clean
		// comparison, and the two it does not mirror are an inventory entry. This
		// is the load-bearing half of that inventory - a C member with no Go
		// constant must never be reported as a defect, because mirroring a subset
		// is deliberate.
		"LLAMA_FX_SPLIT_MODE_NONE": true, "LLAMA_FX_SPLIT_MODE_LAYER": true,
		"LLAMA_FX_SPLIT_MODE_TENSOR": true,
	}

	for _, v := range rep.Viols {
		if clean[v.Fn] {
			t.Errorf("clean binding %s reported as a violation: RULE%d %s", v.Fn, v.Rule, v.Detail)
		}
	}

	// Exactly the fourteen plants, with fx_get_thing counted twice: a void return
	// descriptor is both a wrong cif (rule 1) and a return buffer libffi never
	// writes (rule 3), which is how the real ggml_backend_cpu_buffer_type
	// defect presented.
	// HooksClean is the control for the function-pointer members, and the half
	// that matters: both of its members hold the code pointer of the callback
	// that implements them, written exactly as yzma's two live SetProgressCallback
	// methods write theirs. It cannot go on the map above, because the member
	// violations are named for the C member and all three fixture structs mirror
	// the same struct fx_hooks - so it is pinned by the Go struct instead.
	for _, v := range rep.Viols {
		if strings.Contains(v.Detail, "HooksClean") {
			t.Errorf("the clean fn-ptr member control was reported: RULE%d %s %s", v.Rule, v.Fn, v.Detail)
		}
	}

	if got, want := len(rep.Viols), 19; got != want {
		t.Errorf("fixture produced %d violations, want %d:\n%s", got, want, dumpViolations(rep))
	}
}

// TestFixtureSignedness pins the signedness class, both halves of it.
//
// kindCompat merges signed and unsigned on purpose, because on arm64 they are
// the same register and the same bytes, so this is deliberately not a rule
// violation: the assertion that the two plants stay out of rep.Viols lives in
// TestFixtureDoesNotReportTheCleanBinding, and the assertion that they are still
// *reported* lives here. The controls matter as much: a rule that flagged
// fx_get_token or fx_decode_ok would flag most of yzma.
func TestFixtureSignedness(t *testing.T) {
	rep := fixtureReport(t)

	want := []string{
		// The pointer target: C writes int32 and Go reads uint32, so a -1 comes
		// back as 4294967295.
		"fx_get_count",
		// The return buffer: fx_decode's negative error code read unsigned.
		"fx_decode",
	}

	for _, w := range want {
		if !slices.ContainsFunc(rep.Signs, func(s string) bool { return strings.HasPrefix(s, w+" (") }) {
			t.Errorf("signedness plant %s was not reported:\n%v", w, rep.Signs)
		}
	}

	if got := len(rep.Signs); got != len(want) {
		t.Errorf("signedness findings = %d, want %d:\n%v", got, len(want), rep.Signs)
	}

	// The 1-byte exemption. Plain char is signed in the C table and every string
	// argument in the tree is a Go *byte, so a rule without this would report
	// fx_mode_from_str and nothing useful.
	for _, s := range rep.Signs {
		if strings.Contains(s, "char") {
			t.Errorf("a 1-byte char target was reported as a signedness finding: %s", s)
		}
	}
}

// TestFixtureAccounting guards the counters the report leans on to claim
// coverage. A rule that silently stops checking would still print "0
// violations", so the checked counts are part of the contract.
func TestFixtureAccounting(t *testing.T) {
	rep := fixtureReport(t)

	if got, want := len(rep.Bindings), 29; got != want {
		t.Errorf("bindings found = %d, want %d", got, want)
	}

	if got, want := rep.Matched, 29; got != want {
		t.Errorf("bindings matched to a C decl = %d, want %d", got, want)
	}

	if len(rep.NoC) != 0 {
		t.Errorf("fixture bindings with no C decl: %v", rep.NoC)
	}

	// The variadic pair: the plant and its control. Nothing checked nfixed
	// before, and a variadic binding was exempt from arity checking altogether,
	// so a count of zero here would mean the check has nothing to run on.
	if got, want := rep.NVariadic, 2; got != want {
		t.Errorf("variadic bindings found = %d, want %d", got, want)
	}

	// The coverage inventory. fx_unbound is the only declaration in the fixture
	// header nothing binds, and it must be an inventory entry rather than a
	// violation - which is what TestFixtureDoesNotReportTheCleanBinding pins.
	if got, want := len(rep.Unbound), 1; got != want {
		t.Errorf("unbound C declarations = %d, want %d: %v", got, want, rep.Unbound)
	}

	if len(rep.Unbound) == 1 && !strings.HasPrefix(rep.Unbound[0], "fx_unbound (llama.h:") {
		t.Errorf("unbound inventory entry = %q, want fx_unbound from llama.h", rep.Unbound[0])
	}

	// Grouping is per header, and the fixture's declarations all live in
	// llama.h, so there is exactly one line and its two numbers must add up
	// against the bindings above.
	if got, want := len(rep.Coverage), 1; got != want {
		t.Fatalf("coverage lines = %d, want %d: %+v", got, want, rep.Coverage)
	}

	if c := rep.Coverage[0]; c.Header != "llama.h" || c.Bound != 29 || c.Decls != 30 {
		t.Errorf("coverage line = %+v, want llama.h 29 of 30 bound", c)
	}

	if len(rep.Unresolved) != 0 {
		t.Errorf("unresolvable arg exprs in fixture: %v", rep.Unresolved)
	}

	if rep.Unparsed != 0 {
		t.Errorf("unparsed C declarations = %d, want 0", rep.Unparsed)
	}

	if len(rep.Skips) != 0 {
		t.Errorf("fixture produced %d skip(s), so its coverage is not total: %v", len(rep.Skips), rep.Skips)
	}

	if rep.CheckedR1 == 0 || rep.CheckedR2 == 0 || rep.CheckedR3 == 0 || rep.CheckedR4 == 0 || rep.CheckedR5 == 0 {
		t.Errorf("a rule checked nothing: r1=%d r2=%d r3=%d r4=%d r5=%d",
			rep.CheckedR1, rep.CheckedR2, rep.CheckedR3, rep.CheckedR4, rep.CheckedR5)
	}

	// The four callback sites: two plants and one clean control per form. An
	// unresolved link to a C typedef is a skip rather than a pass, so a
	// normalisation that stopped working would drop the checked count here long
	// before it dropped a violation.
	if got, want := rep.CheckedR5, 4; got != want {
		t.Errorf("callbacks checked = %d, want %d (skips: %v)", got, want, rep.Skips)
	}

	if got, want := rep.CleanR5, 2; got != want {
		t.Errorf("clean callbacks = %d, want %d", got, want)
	}

	if got, want := len(rep.Callbacks), 4; got != want {
		t.Errorf("callback sites found = %d, want %d", got, want)
	}

	if got, want := rep.CFnPtrs, 4; got != want {
		t.Errorf("C function-pointer typedefs parsed = %d, want %d", got, want)
	}

	if rep.CFnPtrBad != 0 {
		t.Errorf("C function-pointer typedefs that could not be parsed = %d, want 0", rep.CFnPtrBad)
	}

	// Every constant in the fixture bindings is mirrored from the header, so all
	// eleven are compared and only the plant differs. An exact count is what
	// pins the mapping: a normalisation that stopped matching would show up here
	// as a lower checked count long before it showed up as a missing violation.
	if got, want := rep.CheckedR4, 12; got != want {
		t.Errorf("constants checked = %d, want %d (skips: %v)", got, want, rep.Skips)
	}

	if got, want := rep.CleanR4, 11; got != want {
		t.Errorf("clean constants = %d, want %d", got, want)
	}

	// The partially-mirrored enum inventory. llama_fx_split_mode is the only enum
	// in the fixture header with some members mirrored and some not:
	// llama_fx_level and llama_fx_flag are mirrored in full, so listing them
	// would be noise. The two unmirrored members are inventory entries, which
	// TestFixtureDoesNotReportTheCleanBinding pins as never being violations.
	if got, want := len(rep.PartialEnums), 1; got != want {
		t.Fatalf("partially mirrored enums = %d, want %d: %+v", got, want, rep.PartialEnums)
	}

	if ec := rep.PartialEnums[0]; ec.Enum != "llama_fx_split_mode" || ec.Members != 3 || ec.Mirrored != 1 || len(ec.Missing) != 2 {
		t.Errorf("partial enum line = %+v, want llama_fx_split_mode 1 of 3 mirrored", ec)
	}

	if got, want := rep.EnumMissing, 2; got != want {
		t.Errorf("unmirrored enum members = %d, want %d", got, want)
	}

	if !strings.HasPrefix(rep.PartialEnums[0].Missing[0], "LLAMA_FX_SPLIT_MODE_LAYER = 1 (llama.h:") {
		t.Errorf("first unmirrored member = %q, want LLAMA_FX_SPLIT_MODE_LAYER from llama.h", rep.PartialEnums[0].Missing[0])
	}

	if rep.CConstBad != 0 {
		t.Errorf("C constants whose value could not be evaluated = %d, want 0", rep.CConstBad)
	}

	if rep.LocalConsts != 0 {
		t.Errorf("fixture constants excluded as yzma-local = %d, want 0", rep.LocalConsts)
	}

	// Struct layouts are compared per struct-by-value slot, once against the C
	// struct and once against the Go struct the bytes come from: ten slots,
	// nineteen comparisons (fx_use_geom_short never reaches its Go-side
	// comparison, the size gate catches it first), four of which are the plants.
	// The three fx_hooks slots are clean on every side, which is the point of
	// them: their plants are invisible to a layout comparison.
	if got, want := rep.CheckedLayout, 19; got != want {
		t.Errorf("struct layouts compared = %d, want %d", got, want)
	}

	if got, want := rep.CleanLayout, 15; got != want {
		t.Errorf("clean struct layouts = %d, want %d", got, want)
	}

	// The function-pointer struct members: two per fx_hooks slot, deduplicated to
	// one comparison per pair of representations rather than one per slot. Four are
	// clean and two are the plants. Traced counts the narrower claim - a member
	// whose stored code pointer could be followed back to a callback site - so a
	// dataflow that stopped resolving would show up here rather than as silence.
	if got, want := rep.CheckedFnPtr, 6; got != want {
		t.Errorf("fn-ptr members compared = %d, want %d (skips: %v)", got, want, rep.Skips)
	}

	if got, want := rep.CleanFnPtr, 4; got != want {
		t.Errorf("clean fn-ptr members = %d, want %d", got, want)
	}

	if got, want := rep.TracedFnPtr, 4; got != want {
		t.Errorf("fn-ptr members with a traced code pointer = %d, want %d", got, want)
	}

	// The pointer-target comparisons. Eight of the fixture's pointer slots name a
	// concrete target on both sides; the other pointer slots aim at an opaque
	// struct fx_thing through a Go uintptr, which names nothing to compare, and
	// are out of scope rather than skips. Nine are clean - the three char *
	// buffers among them included, since a *byte is the right target for a char -
	// one is the plant and one is the signedness finding.
	if got, want := rep.CheckedPtr, 11; got != want {
		t.Errorf("pointer targets compared = %d, want %d", got, want)
	}

	if got, want := rep.CleanPtr, 9; got != want {
		t.Errorf("clean pointer targets = %d, want %d", got, want)
	}

	// The NUL-termination comparisons. Three of the fixture's char * slots are fed
	// a buffer this can trace back to a Go string: the plant and its two controls.
	// The other char * slots - fx_desc's output buffer from make, and
	// fx_mode_from_str's *byte parameter - are produced somewhere this cannot see,
	// so they are out of scope rather than skips, exactly as a void * is for the
	// pointer targets. A count that drifted down would mean the check quietly
	// stopped looking.
	if got, want := rep.CheckedNUL, 3; got != want {
		t.Errorf("C string buffers checked = %d, want %d", got, want)
	}

	if got, want := rep.CleanNUL, 2; got != want {
		t.Errorf("NUL-terminated C string buffers = %d, want %d", got, want)
	}

	// The enum-parameter comparisons. Two of the fixture's three enum slots are fed
	// a Go type whose mirrored members say which C enum it stands for: the plant and
	// its control. fx_set_mode's plain int32 mirrors nothing, so it is out of scope
	// rather than a skip, exactly as a void * is for the pointer targets - and
	// pinning the count at two is what keeps that narrowing visible.
	if got, want := rep.CheckedEnum, 2; got != want {
		t.Errorf("enum params compared = %d, want %d", got, want)
	}

	if got, want := rep.CleanEnum, 1; got != want {
		t.Errorf("clean enum params = %d, want %d", got, want)
	}

	// Ten slot lines plus the one Go-only-tail note for fx_use_geom_tail.
	if got, want := len(rep.StructCmp), 11; got != want {
		t.Errorf("struct-by-value slots compared = %d, want %d", got, want)
	}
}

// TestAnalyseIsRepeatable pins the parser-state reset. The C typedef, struct
// and enum tables are package-level and iterated to a fixpoint, so a leak
// between runs would make results depend on call order.
func TestAnalyseIsRepeatable(t *testing.T) {
	first := dumpViolations(fixtureReport(t))
	second := dumpViolations(fixtureReport(t))

	if first != second {
		t.Errorf("two identical runs disagreed:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func dumpViolations(r *report) string {
	if len(r.Viols) == 0 {
		return "  (none)"
	}

	var b strings.Builder
	for _, v := range r.Viols {
		b.WriteString("  RULE")
		b.WriteString(string(rune('0' + v.Rule)))
		b.WriteString(" ")
		b.WriteString(v.Fn)
		b.WriteString(": ")
		b.WriteString(v.Detail)
		b.WriteString("\n")
	}

	return b.String()
}
