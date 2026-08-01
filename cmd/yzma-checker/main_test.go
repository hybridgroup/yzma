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
// package carrying one planted defect per rule per direction plus one clean
// binding. Both halves matter — every plant must be found, and the clean
// binding must never be reported.

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

	for _, v := range rep.Viols {
		if v.Fn == "fx_clean" {
			t.Errorf("clean binding reported as a violation: RULE%d %s", v.Rule, v.Detail)
		}
	}

	// Exactly the four plants, with fx_get_thing counted twice: a void return
	// descriptor is both a wrong cif (rule 1) and a return buffer libffi never
	// writes (rule 3), which is how the real ggml_backend_cpu_buffer_type
	// defect presented.
	if got, want := len(rep.Viols), 5; got != want {
		t.Errorf("fixture produced %d violations, want %d:\n%s", got, want, dumpViolations(rep))
	}
}

// TestFixtureAccounting guards the counters the report leans on to claim
// coverage. A rule that silently stops checking would still print "0
// violations", so the checked counts are part of the contract.
func TestFixtureAccounting(t *testing.T) {
	rep := fixtureReport(t)

	if got, want := len(rep.Bindings), 5; got != want {
		t.Errorf("bindings found = %d, want %d", got, want)
	}

	if got, want := rep.Matched, 5; got != want {
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

	if rep.CheckedR1 == 0 || rep.CheckedR2 == 0 || rep.CheckedR3 == 0 {
		t.Errorf("a rule checked nothing: r1=%d r2=%d r3=%d", rep.CheckedR1, rep.CheckedR2, rep.CheckedR3)
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
