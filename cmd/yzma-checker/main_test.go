package main

import (
	"strings"
	"testing"
)

// The checker's own correctness gate.
//
// It used to be "a run must re-derive these two live defects in yzma", which
// stopped working the moment the defects were fixed: a tool that reports
// nothing is indistinguishable from a tool that checks nothing. Instead the
// gate now runs the full pipeline over testdata/fixture, a miniature binding
// package carrying one planted defect per rule per direction plus clean
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
			// The constant plant: LLAMA_FX_LEVEL_LOW was inserted in front of
			// HIGH upstream, so HIGH is 2 in the header and still 1 in Go. No
			// width, offset or class changed, which is why RULES 1-3 cannot see
			// it.
			name: "rule4 enum member value shifted by an insertion upstream",
			rule: 4, fn: "LLAMA_FX_LEVEL_HIGH",
			match: "FxLevelHigh = 1 but C LLAMA_FX_LEVEL_HIGH = 2",
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
	clean := map[string]bool{
		"fx_clean": true, "fx_use_geom": true, "fx_use_geom_tail": true, "fx_use_named": true,
		"LLAMA_FX_LEVEL_OFF": true, "LLAMA_FX_LEVEL_LOW": true, "LLAMA_FX_LEVEL_MAX": true,
		"LLAMA_FX_FLAG_AUTO": true, "LLAMA_FX_FLAG_NONE": true,
		"LLAMA_FX_FLAG_A": true, "LLAMA_FX_FLAG_B": true,
		"LLAMA_FX_MAGIC": true, "LLAMA_FX_VERSION": true, "LLAMA_FX_MAGIC_ALIAS": true,
	}

	for _, v := range rep.Viols {
		if clean[v.Fn] {
			t.Errorf("clean binding %s reported as a violation: RULE%d %s", v.Fn, v.Rule, v.Detail)
		}
	}

	// Exactly the nine plants, with fx_get_thing counted twice: a void return
	// descriptor is both a wrong cif (rule 1) and a return buffer libffi never
	// writes (rule 3), which is how the real ggml_backend_cpu_buffer_type
	// defect presented.
	if got, want := len(rep.Viols), 11; got != want {
		t.Errorf("fixture produced %d violations, want %d:\n%s", got, want, dumpViolations(rep))
	}
}

// TestFixtureAccounting guards the counters the report leans on to claim
// coverage. A rule that silently stops checking would still print "0
// violations", so the checked counts are part of the contract.
func TestFixtureAccounting(t *testing.T) {
	rep := fixtureReport(t)

	if got, want := len(rep.Bindings), 12; got != want {
		t.Errorf("bindings found = %d, want %d", got, want)
	}

	if got, want := rep.Matched, 12; got != want {
		t.Errorf("bindings matched to a C decl = %d, want %d", got, want)
	}

	if len(rep.NoC) != 0 {
		t.Errorf("fixture bindings with no C decl: %v", rep.NoC)
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

	if rep.CheckedR1 == 0 || rep.CheckedR2 == 0 || rep.CheckedR3 == 0 || rep.CheckedR4 == 0 {
		t.Errorf("a rule checked nothing: r1=%d r2=%d r3=%d r4=%d",
			rep.CheckedR1, rep.CheckedR2, rep.CheckedR3, rep.CheckedR4)
	}

	// Every constant in the fixture bindings is mirrored from the header, so all
	// eleven are compared and only the plant differs. An exact count is what
	// pins the mapping: a normalisation that stopped matching would show up here
	// as a lower checked count long before it showed up as a missing violation.
	if got, want := rep.CheckedR4, 11; got != want {
		t.Errorf("constants checked = %d, want %d (skips: %v)", got, want, rep.Skips)
	}

	if got, want := rep.CleanR4, 10; got != want {
		t.Errorf("clean constants = %d, want %d", got, want)
	}

	if rep.CConstBad != 0 {
		t.Errorf("C constants whose value could not be evaluated = %d, want 0", rep.CConstBad)
	}

	if rep.LocalConsts != 0 {
		t.Errorf("fixture constants excluded as yzma-local = %d, want 0", rep.LocalConsts)
	}

	// Struct layouts are compared per struct-by-value slot, once against the C
	// struct and once against the Go struct the bytes come from: seven slots,
	// thirteen comparisons (fx_use_geom_short never reaches its Go-side
	// comparison, the size gate catches it first), four of which are the plants.
	if got, want := rep.CheckedLayout, 13; got != want {
		t.Errorf("struct layouts compared = %d, want %d", got, want)
	}

	if got, want := rep.CleanLayout, 9; got != want {
		t.Errorf("clean struct layouts = %d, want %d", got, want)
	}

	// Seven slot lines plus the one Go-only-tail note for fx_use_geom_tail.
	if got, want := len(rep.StructCmp), 8; got != want {
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
