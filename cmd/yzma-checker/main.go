package main

import (
	"flag"
	"fmt"
	"go/token"
	"go/types"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

var (
	yzmaDir  = flag.String("yzma", "", "yzma tree to audit (default: walk up for a yzma checkout, else `go list -m`)")
	llamaCpp = flag.String("llama", "", "llama.cpp git checkout to read headers from (default: fetch them from upstream)")
	ref      = flag.String("ref", "auto", `llama.cpp git ref to audit against; "auto" resolves the current llama-cpp-builder release`)
	hdrDir   = flag.String("hdrs", "", "use pre-extracted headers from this dir instead of git or the network")
	pkgs     = flag.String("pkgs", yzmaModulePath+"/pkg/llama,"+yzmaModulePath+"/pkg/mtmd",
		"comma-separated package patterns to audit")
	verbose = flag.Bool("v", false, "dump every binding with its C signature and call sites")
)

var skips []string

func noteSkip(format string, a ...any) { skips = append(skips, fmt.Sprintf(format, a...)) }

var structCmp []string

// Struct-by-value layout comparisons, counted separately from the rule the
// comparison belongs to: a layout check that silently stopped resolving would
// otherwise still print "0 violations". See layout.go.
var layoutChecked, layoutClean int

// members renders a leaf count for the STRUCT-BY-VALUE COMPARISONS report,
// distinguishing "no members" from "layout not resolvable".
func members(ls []Leaf, ok bool) string {
	if !ok {
		return "layout?"
	}

	return fmt.Sprintf("%d members", len(ls))
}

// hdrCoverage is how much of one header yzma binds. It is an inventory line,
// never a finding: see the coverage block in analyse.
type hdrCoverage struct {
	Header string
	Decls  int
	Bound  int
}

type violation struct {
	Rule     int
	Fn       string
	GoFile   string
	Detail   string
	Severity string
}

func main() {
	flag.Parse()

	yzmaSrc := "-yzma"
	if *yzmaDir == "" {
		dir, src, err := findYzmaRoot()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		*yzmaDir, yzmaSrc = dir, src
	}

	gitRef, refSrc := *ref, "-ref"
	if gitRef == "auto" {
		r, err := resolveRef()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		gitRef, refSrc = r, "current llama-cpp-builder release"
	}

	headers, hdrSrc, cleanup, err := obtainHeaders(*hdrDir, *llamaCpp, gitRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot obtain headers: %v\n", err)
		os.Exit(2)
	}
	defer cleanup()

	*hdrDir = headers

	fmt.Printf("yzma:   %s (%s)\n", *yzmaDir, yzmaSrc)
	fmt.Printf("ref:    %s (%s)\n", gitRef, refSrc)
	fmt.Printf("hdrs:   %s\n          %s\n\n", *hdrDir, hdrSrc)

	rep, err := analyse(*yzmaDir, *hdrDir, strings.Split(*pkgs, ","))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	printReport(rep)

	if len(rep.Viols) > 0 {
		os.Exit(1)
	}
}

// report is everything one audit produced. The accounting fields matter as
// much as Viols: they are what makes "these are the only ones" a measurement
// rather than an assertion.
type report struct {
	Viols      []violation
	Bindings   []*Binding
	Callbacks  []*Callback
	CFuncs     map[string]CFunc
	NoC        []string
	Unresolved []string
	Skips      []string
	StructCmp  []string

	// Unbound and Coverage are the reverse of NoC: C declarations no binding
	// names. Not defects - yzma binds a subset by design - so they are counted
	// and listed, never turned into violations.
	Unbound  []string
	Coverage []hdrCoverage

	TotalDecls, Unparsed                                  int
	Matched, NCalls, NVariadic                            int
	CheckedR1, CheckedR2, CheckedR3, CheckedR4, CheckedR5 int
	CleanR1, CleanR2, CleanR3, CleanR4, CleanR5           int
	CheckedLayout, CleanLayout                            int

	// CFnPtrs counts the function-pointer typedefs RULE 5 has a signature for,
	// CFnPtrBad the ones this parser could not take apart. As with CConstBad the
	// second number is the rule's coverage limit: a callback linking to one is a
	// skip, never a pass.
	CFnPtrs, CFnPtrBad int

	// CConsts counts the C constants whose value this run pinned down, CConstBad
	// the ones it could not. The second number is the coverage limit of RULE 4:
	// none of those can be compared against, and a Go constant mirroring one
	// becomes a skip.
	CConsts, CConstBad, LocalConsts int
}

// analyse runs the three rules over the packages named by patterns in the
// module rooted at yzmaRoot, against the headers in headerDir.
func analyse(yzmaRoot, headerDir string, patterns []string) (*report, error) {
	// Parser state is package-level and iterated to a fixpoint, so reset it
	// to keep repeated calls within one process independent.
	skips, structCmp = nil, nil
	layoutChecked, layoutClean = 0, 0
	resetCTypes()

	// --- C side ---
	cfuncs := map[string]CFunc{}
	unparsed := 0
	totalDecls := 0
	macros := []string{"LLAMA_API", "GGML_API", "GGML_BACKEND_API", "MTMD_API"}
	var hdrSrcs []struct{ path, src string }
	for _, hf := range headerFiles {
		path := filepath.Join(headerDir, hf.local)
		raw, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SKIP header %s: %v\n", path, err)
			continue
		}
		src := unwrapDeprecated(stripComments(string(raw)))
		collectTypedefs(src)
		collectFnPtrTypedefs(path, src)
		collectStructs(src)
		hdrSrcs = append(hdrSrcs, struct{ path, src string }{path, src})
		var fns []CFunc
		for _, mc := range macros {
			mf, err := parseHeader(path, mc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse %s: %v\n", path, err)
				continue
			}
			fns = append(fns, mf...)
		}
		for _, f := range fns {
			totalDecls++
			if f.Name == "" {
				unparsed++
				fmt.Fprintf(os.Stderr, "UNPARSED %s:%d %s\n", filepath.Base(f.File), f.Line, f.Raw)
				continue
			}
			if _, dup := cfuncs[f.Name]; !dup {
				cfuncs[f.Name] = f
			}
		}
	}
	// Constant values are also iterated, and for the same reason as the structs
	// below: llama.h is parsed first but initialises LLAMA_ROPE_TYPE_NEOX from
	// the GGML_ROPE_TYPE_NEOX defined in ggml.h, so an initialiser can name a
	// constant no earlier pass has seen. The failed attempts are forgotten
	// between passes so that a name resolved late stops being reported as
	// unevaluable.
	for range 3 {
		clear(cconstBad)
		clear(cconstBadByNorm)

		for _, h := range hdrSrcs {
			collectCConsts(h.path, h.src)
		}
	}

	// Second pass: headers are parsed in sequence, so a struct in llama.h can
	// reference a typedef that only appears in ggml-backend.h. Re-classify every
	// field now that every typedef and struct is known, then re-lay-out.
	for pass := 0; pass < 4; pass++ {
		for _, cs := range cstructs {
			for i := range cs.Fields {
				cs.Fields[i].Kind, cs.Fields[i].Size = classify(cs.Fields[i].Norm)
			}
			computeStructSize(cs)
		}
	}
	// re-classify now that all typedefs/structs are known
	for n, f := range cfnptrs {
		f.RetKind, f.RetSize = classify(f.RetRaw)
		for i := range f.Params {
			f.Params[i].Kind, f.Params[i].Size = classify(f.Params[i].Norm)
		}
		cfnptrs[n] = f
	}
	for n, f := range cfuncs {
		f.RetKind, f.RetSize = classify(f.RetRaw)
		for i := range f.Params {
			f.Params[i].Kind, f.Params[i].Size = classify(f.Params[i].Norm)
		}
		cfuncs[n] = f
	}

	// --- Go side ---
	loaded, err := loadPkgs(yzmaRoot, patterns...)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", strings.Join(patterns, ", "), err)
	}
	var all []*Binding
	var callbacks []*Callback
	for _, p := range loaded {
		if len(p.Errors) > 0 {
			for _, e := range p.Errors {
				fmt.Fprintf(os.Stderr, "pkgerr %s: %v\n", p.PkgPath, e)
			}
		}
		if p.Syntax == nil {
			continue
		}
		a := &analyzer{pkg: p, fset: p.Fset}
		a.run()
		for _, k := range a.order {
			all = append(all, a.bindings[k])
		}
		callbacks = append(callbacks, a.callbacks...)
	}
	sort.Slice(callbacks, func(i, j int) bool {
		if callbacks[i].Pos.Filename != callbacks[j].Pos.Filename {
			return callbacks[i].Pos.Filename < callbacks[j].Pos.Filename
		}
		return callbacks[i].Pos.Line < callbacks[j].Pos.Line
	})
	sort.Slice(all, func(i, j int) bool {
		if all[i].PrepPos.Filename != all[j].PrepPos.Filename {
			return all[i].PrepPos.Filename < all[j].PrepPos.Filename
		}
		return all[i].PrepPos.Line < all[j].PrepPos.Line
	})

	var viols []violation
	matched, noC, checkedR1, checkedR2, checkedR3 := 0, 0, 0, 0, 0
	cleanR1, cleanR2, cleanR3 := 0, 0, 0
	var noCList []string
	nCalls := 0
	unresolvedArgs := 0
	var unresolvedList []string
	nVariadic := 0

	for _, b := range all {
		short := shortPos(b.PrepPos)
		if b.Variadi {
			nVariadic++
		}

		cf, ok := cfuncs[b.CName]
		if !ok {
			noC++
			noCList = append(noCList, fmt.Sprintf("%s (%s)", b.CName, short))
		} else {
			matched++
		}

		// ---------- Rule 1: cif vs C prototype ----------
		if ok {
			checkedR1++
			var probs []string

			// Fixed and variadic are two different calling conventions, not two
			// spellings of one. On Apple arm64 a variadic argument is passed on
			// the stack where a fixed one goes in a register, so `nfixed` is not
			// a value the ABI tolerates being off by one: every argument past
			// the fixed ones is read from the wrong place. A variadic binding
			// used to be exempt from arity checking altogether, which meant the
			// one number that decides that split was never looked at.
			switch {
			case b.Variadi && !cf.Vararg:
				probs = append(probs, fmt.Sprintf("variadic: bound with PrepVar (nfixed %d) but C %s declares no \"...\"",
					b.NFixed, cf.Name))
			case !b.Variadi && cf.Vararg:
				probs = append(probs, fmt.Sprintf("variadic: C %s is variadic but the binding is a fixed-arity Prep of %d args",
					cf.Name, len(b.Args)))
			case cf.Vararg:
				// The parser never records the "..." as a parameter, so the
				// length of cf.Params *is* the fixed parameter count.
				if b.NFixed < 0 {
					noteSkip("RULE1 %s (%s): PrepVar nfixed is not a literal - NOT VERIFIED", b.CName, short)
					break
				}
				if b.NFixed != len(cf.Params) {
					probs = append(probs, fmt.Sprintf("nfixed is %d but C %s declares %d parameter(s) before its \"...\"",
						b.NFixed, cf.Name, len(cf.Params)))
				}
				// A PrepVar cif lists the fixed types *and* the concrete
				// variadic types of the call it describes, so it can be longer
				// than nfixed but never shorter: a shorter list means libffi is
				// told fixed arguments exist that it has no type for.
				if len(b.Args) < b.NFixed {
					probs = append(probs, fmt.Sprintf("nfixed is %d but the cif only lists %d arg type(s)",
						b.NFixed, len(b.Args)))
				}
			case len(b.Args) != len(cf.Params):
				probs = append(probs, fmt.Sprintf("arity: cif has %d args, C has %d", len(b.Args), len(cf.Params)))
			}
			// return type
			if p := cmpTypes(b.Ret, cf.RetRaw, cf.RetKind, cf.RetSize, "ret"); p != "" {
				probs = append(probs, p)
			}
			n := len(b.Args)
			if n > len(cf.Params) {
				n = len(cf.Params)
			}
			for i := 0; i < n; i++ {
				if p := cmpTypes(b.Args[i], cf.Params[i].Norm, cf.Params[i].Kind, cf.Params[i].Size, fmt.Sprintf("arg%d", i)); p != "" {
					probs = append(probs, p)
				}
			}
			if len(probs) > 0 {
				viols = append(viols, violation{1, b.CName, short, strings.Join(probs, "; "), sev(probs)})
			} else {
				cleanR1++
			}
		}

		// ---------- Rules 2 & 3: call sites ----------
		for _, cs := range b.Calls {
			nCalls++
			csp := shortPos(cs.Pos)
			// Rule 3: return buffer
			checkedR3++
			cRet := ""
			if ok {
				cRet = cf.RetRaw
			}
			r3 := checkRet(b.Ret, cs, cRet)
			if strings.HasPrefix(r3, "SKIP: ") {
				noteSkip("RULE3 %s (%s): %s", b.CName, csp, strings.TrimPrefix(r3, "SKIP: "))
				r3 = ""
			}
			if r3 != "" {
				viols = append(viols, violation{3, b.CName, csp, r3 + " [in " + cs.Fn + "]", "latent-corruption"})
			} else {
				cleanR3++
			}
			// Rule 2: argument widths
			var probs []string
			na := len(cs.Args)
			if na != len(b.Args) && !b.Variadi {
				probs = append(probs, fmt.Sprintf("call passes %d avalues, cif expects %d (runtime panic)", na, len(b.Args)))
			}
			m := na
			if len(b.Args) < m {
				m = len(b.Args)
			}
			for i := 0; i < m; i++ {
				ca, ct := cs.Args[i], b.Args[i]
				if ca.Pointee == nil {
					unresolvedArgs++
					unresolvedList = append(unresolvedList,
						fmt.Sprintf("%s arg%d %s: %s (%s)", b.CName, i, ca.Expr, ca.Note, csp))
					continue
				}
				checkedR2++
				// Count this slot clean only if nothing below reports on it:
				// "514 / 514 clean" next to a rule2 violation is an accounting
				// bug, and the counts are the whole basis for the coverage claim.
				before := len(probs)
				if ct.Size <= 0 && ct.Kind != KindVoid {
					noteSkip("RULE2 %s arg%d: cif descriptor %s size unknown - NOT VERIFIED", b.CName, i, ct.Name)
					continue
				}
				// For a struct passed by value the C-declared size is the
				// authority: libffi reads exactly that many bytes through the
				// pointer, so a larger Go struct is untouched tail (see
				// cmpStructLayout) and only a short one is a read overrun.
				// Scalars stay exact - a mismatch either way is a defect.
				short := ca.Size != ct.Size
				if ct.Kind == KindStruct && ca.Kind == KindStruct {
					short = ca.Size < ct.Size
				}
				if short {
					probs = append(probs, fmt.Sprintf("arg%d: cif %s wants %dB, Go %s is %dB (%s)",
						i, ct.Name, ct.Size, ca.Pointee.String(), ca.Size, ca.Expr))
					continue
				}
				if !kindCompat(ct.Kind, ca.Kind) {
					probs = append(probs, fmt.Sprintf("arg%d: cif kind %s vs Go kind %s (%s / %s)",
						i, ct.Kind, ca.Kind, ca.Expr, ca.Pointee.String()))
				}
				cArg := ""
				if ok && i < len(cf.Params) {
					cArg = cf.Params[i].Norm
				}
				if p := cmpStructLayout(ct, ca.Pointee, fmt.Sprintf("arg%d", i), ca.Expr, cArg); p != "" {
					probs = append(probs, p)
				}
				if len(probs) == before {
					cleanR2++
				}
			}
			if len(probs) > 0 {
				viols = append(viols, violation{2, b.CName, csp, strings.Join(probs, "; ") + " [in " + cs.Fn + "]", "memory-read-overrun"})
			}
		}
		if *verbose {
			fmt.Printf("BINDING %-45s %s\n  cif ret=%s args=%v\n", b.CName, short, b.Ret, b.Args)
			if ok {
				fmt.Printf("  C   %s:%d %s\n", filepath.Base(cf.File), cf.Line, cf.Sig())
			} else {
				fmt.Printf("  C   <NOT FOUND IN HEADERS>\n")
			}
			for _, cs := range b.Calls {
				fmt.Printf("  call %s ret=%s args=%v\n", shortPos(cs.Pos), cs.RetExpr, argSummary(cs.Args))
			}
		}
	}

	// ---------- Coverage inventory: C declarations nothing binds ----------
	//
	// This is deliberately not a rule. yzma binds a chosen subset of llama.cpp,
	// so an unbound declaration is not a defect and never becomes a violation or
	// affects the exit code. What it is, is the only way the coverage claim is
	// measurable over time: a function yzma should bind appearing upstream, or
	// the neighbours of a bound function changing shape, is invisible unless the
	// unbound set is written down somewhere a diff can see it.
	//
	// Grouped per header because that is the actionable unit - "mtmd-helper.h: 3
	// of 40 bound" is a number a maintainer can decide about, and ~600 names in
	// one list is a number nobody reads.
	bound := make(map[string]bool, len(all))
	for _, b := range all {
		bound[b.CName] = true
	}

	var unbound []string
	perHdr := map[string]*hdrCoverage{}
	for name, cf := range cfuncs {
		hdr := filepath.Base(cf.File)
		hc, ok := perHdr[hdr]
		if !ok {
			hc = &hdrCoverage{Header: hdr}
			perHdr[hdr] = hc
		}

		hc.Decls++
		if bound[name] {
			hc.Bound++
			continue
		}

		unbound = append(unbound, fmt.Sprintf("%s (%s:%d)", name, hdr, cf.Line))
	}

	sort.Strings(unbound)

	var coverage []hdrCoverage
	for _, hdr := range slices.Sorted(maps.Keys(perHdr)) {
		coverage = append(coverage, *perHdr[hdr])
	}

	// ---------- Rule 5: callback descriptors and closure signatures ----------
	checkedR5, cleanR5 := 0, 0
	for _, cb := range callbacks {
		linkCallback(cb)

		probs, checked := checkCallback(cb)
		if !checked {
			continue
		}

		checkedR5++
		if len(probs) == 0 {
			cleanR5++
		} else {
			where := cb.Form + " " + cb.GoID
			if cb.Fn != cb.GoID {
				where += " in " + cb.Fn
			}

			viols = append(viols, violation{5, cb.CTypedef, shortPos(cb.Pos),
				strings.Join(probs, "; ") + " [" + where + "]", sev(probs)})
		}

		if *verbose {
			fmt.Printf("CALLBACK %-32s %s\n  %s %s (linked by %s)\n  C   %s:%d %s\n",
				cb.GoID, shortPos(cb.Pos), cb.Form, cb.CTypedef, cb.Link,
				filepath.Base(cfnptrs[cb.CTypedef].File), cfnptrs[cb.CTypedef].Line, cfnptrs[cb.CTypedef].Sig())
		}
	}

	// ---------- Rule 4: mirrored constant values ----------
	constViols, checkedR4, cleanR4, localConsts := checkConsts(loaded)
	viols = append(viols, constViols...)

	sort.SliceStable(viols, func(i, j int) bool { return viols[i].Rule < viols[j].Rule })

	return &report{
		Viols:      viols,
		Bindings:   all,
		Callbacks:  callbacks,
		CFuncs:     cfuncs,
		NoC:        noCList,
		Unbound:    unbound,
		Coverage:   coverage,
		Unresolved: unresolvedList,
		Skips:      skips,
		StructCmp:  structCmp,
		TotalDecls: totalDecls,
		Unparsed:   unparsed,
		Matched:    matched,
		NCalls:     nCalls,
		NVariadic:  nVariadic,
		CheckedR1:  checkedR1,
		CheckedR2:  checkedR2,
		CheckedR3:  checkedR3,
		CheckedR4:  checkedR4,
		CheckedR5:  checkedR5,
		CleanR1:    cleanR1,
		CleanR2:    cleanR2,
		CleanR3:    cleanR3,
		CleanR4:    cleanR4,
		CleanR5:    cleanR5,

		CFnPtrs:   len(cfnptrs),
		CFnPtrBad: len(cfnptrBad),

		CheckedLayout: layoutChecked,
		CleanLayout:   layoutClean,

		CConsts:     len(cconsts),
		CConstBad:   len(cconstBad),
		LocalConsts: localConsts,
	}, nil
}

func printReport(r *report) {
	fmt.Println("================ VIOLATIONS ================")
	for _, v := range r.Viols {
		fmt.Printf("RULE%d %-14s %-42s %s\n        %s\n", v.Rule, v.Severity, v.Fn, v.GoFile, v.Detail)
	}
	fmt.Println("================ UNRESOLVED ARG EXPRS ================")
	for _, u := range r.Unresolved {
		fmt.Println("  ", u)
	}
	fmt.Println("================ BINDINGS WITH NO C DECL ================")
	for _, u := range r.NoC {
		fmt.Println("  ", u)
	}
	fmt.Println("================ C DECLS WITH NO BINDING (inventory, not defects) ================")
	for _, c := range r.Coverage {
		fmt.Printf("   %-18s %3d of %3d bound (%d unbound)\n", c.Header+":", c.Bound, c.Decls, c.Decls-c.Bound)
	}
	if *verbose {
		for _, u := range r.Unbound {
			fmt.Println("     ", u)
		}
	}
	hint := ""
	if !*verbose && len(r.Unbound) > 0 {
		hint = "; -v lists them"
	}
	fmt.Printf("(%d unbound C declarations%s)\n", len(r.Unbound), hint)
	fmt.Println("================ NOT VERIFIED (skips) ================")
	for _, k := range r.Skips {
		fmt.Println("  ", k)
	}
	fmt.Printf("(%d skips)\n", len(r.Skips))
	fmt.Println("================ STRUCT-BY-VALUE COMPARISONS ================")
	for _, k := range r.StructCmp {
		fmt.Println("  ", k)
	}
	fmt.Printf("(%d lines; %d struct layouts compared, %d clean)\n", len(r.StructCmp), r.CheckedLayout, r.CleanLayout)
	fmt.Println("================ SUMMARY ================")
	fmt.Printf("C declarations found:      %d (unparsed: %d)\n", r.TotalDecls, r.Unparsed)
	fmt.Printf("distinct C functions:      %d\n", len(r.CFuncs))
	fmt.Printf("yzma bindings (Prep):      %d\n", len(r.Bindings))
	fmt.Printf("  matched to a C decl:     %d\n", r.Matched)
	fmt.Printf("  no C decl in headers:    %d\n", len(r.NoC))
	fmt.Printf("  of them variadic:        %d (PrepVar)\n", r.NVariadic)
	fmt.Printf("C decls bound / unbound:   %d / %d (unbound is an inventory, not a defect)\n",
		len(r.CFuncs)-len(r.Unbound), len(r.Unbound))
	fmt.Printf("Call sites analysed:       %d\n", r.NCalls)
	fmt.Printf("callback sites (C->Go):    %d\n", len(r.Callbacks))
	fmt.Printf("Rule1 checked/clean:       %d / %d\n", r.CheckedR1, r.CleanR1)
	fmt.Printf("Rule2 arg slots checked:   %d / %d clean (unresolvable exprs: %d)\n", r.CheckedR2, r.CleanR2, len(r.Unresolved))
	fmt.Printf("Rule3 return bufs checked: %d / %d clean\n", r.CheckedR3, r.CleanR3)
	fmt.Printf("Rule4 constants checked:   %d / %d clean (C constants parsed: %d, unevaluable: %d; yzma-local: %d)\n",
		r.CheckedR4, r.CleanR4, r.CConsts, r.CConstBad, r.LocalConsts)
	fmt.Printf("Rule5 callbacks checked:   %d / %d clean (C function-pointer typedefs parsed: %d, unparseable: %d)\n",
		r.CheckedR5, r.CleanR5, r.CFnPtrs, r.CFnPtrBad)
	fmt.Printf("struct layouts compared:   %d / %d clean\n", r.CheckedLayout, r.CleanLayout)
	nr := map[int]int{}
	for _, v := range r.Viols {
		nr[v.Rule]++
	}
	fmt.Printf("violations: rule1=%d rule2=%d rule3=%d rule4=%d rule5=%d\n", nr[1], nr[2], nr[3], nr[4], nr[5])
}

func argSummary(as []CallArg) string {
	var s []string
	for _, a := range as {
		if a.Pointee == nil {
			s = append(s, a.Expr+"=?"+a.Note)
		} else {
			s = append(s, fmt.Sprintf("%s:%s/%dB", a.Expr, a.Pointee.String(), a.Size))
		}
	}
	return strings.Join(s, ", ")
}

func sev(p []string) string {
	for _, s := range p {
		// A wrong argument count breaks the ABI in either direction: on a call
		// libffi reads the wrong number of avalues, and on a callback every
		// parameter after the missing one arrives shifted.
		if strings.Contains(s, "arity") || strings.Contains(s, "nfixed") ||
			strings.Contains(s, "variadic") ||
			strings.Contains(s, "parameter(s)") || strings.Contains(s, "arg type(s)") {
			return "ABI-BREAK"
		}
	}
	return "type-mismatch"
}

// cmpTypes compares a cif descriptor against the C type it must model.
func cmpTypes(t FfiType, cRaw string, cKind Kind, cSize int, slot string) string {
	if cKind == KindStruct {
		cSize = cTypeSize(cRaw)
		if cSize < 0 {
			noteSkip("%s: C struct %q size unknown - NOT VERIFIED", slot, squash(cRaw))
			return ""
		}
		if t.Kind != KindStruct {
			return fmt.Sprintf("%s: C passes struct %s (%dB) by value but cif says %s", slot, squash(cRaw), cSize, t.Name)
		}

		var probs []string
		if t.Size >= 0 && t.Size != cSize {
			probs = append(probs, fmt.Sprintf("%s: struct descriptor %s is %dB, C %s is %dB", slot, t.Name, t.Size, squash(cRaw), cSize))
		}

		// Equal sizes are not equal layouts: see the comment at the top of
		// layout.go. Compare member offsets too, which is the check that
		// catches a field inserted upstream in front of interior padding.
		cl, cOK := cLeaves(cstructOf(cRaw))
		fl, fOK := ffiLeaves(t)
		switch {
		case !cOK:
			noteSkip("%s: C struct %q member layout not resolvable - LAYOUT NOT VERIFIED", slot, squash(cRaw))
		case !fOK:
			noteSkip("%s: cif descriptor %s member layout not resolvable - LAYOUT NOT VERIFIED", slot, t.Name)
		default:
			layoutChecked++
			if d := diffLeaves(fl, cl, "cif "+t.Name, "C "+squash(cRaw)); len(d) > 0 {
				probs = append(probs, slot+": struct layout differs from C: "+strings.Join(d, "; "))
			} else {
				layoutClean++
			}
		}

		structCmp = append(structCmp, fmt.Sprintf("%s structcmp: cif %s=%dB/%s vs C %s=%dB/%s",
			slot, t.Name, t.Size, members(fl, fOK), squash(cRaw), cSize, members(cl, cOK)))
		if *verbose {
			structCmp = append(structCmp,
				"    cif members: "+leafList(fl),
				"    C   members: "+leafList(cl))
		}

		return strings.Join(probs, "; ")
	}
	if cKind == KindUnknown || cSize < 0 {
		noteSkip("%s: C type %q unclassifiable - NOT VERIFIED", slot, squash(cRaw))
		return ""
	}
	if t.Size < 0 {
		noteSkip("%s: cif descriptor %q size unknown - NOT VERIFIED", slot, t.Name)
		return ""
	}
	if cKind == KindVoid {
		if t.Kind != KindVoid {
			return fmt.Sprintf("%s: C returns void but cif says %s", slot, t.Name)
		}
		return ""
	}
	if t.Kind == KindVoid {
		return fmt.Sprintf("%s: C is %s (%s/%dB) but cif says TypeVoid", slot, squash(cRaw), cKind, cSize)
	}
	if t.Size != cSize {
		return fmt.Sprintf("%s: C %s is %dB but cif %s is %dB", slot, squash(cRaw), cSize, t.Name, t.Size)
	}
	if !kindCompat(t.Kind, cKind) {
		return fmt.Sprintf("%s: C %s is %s but cif %s is %s", slot, squash(cRaw), cKind, t.Name, t.Kind)
	}
	return ""
}

// cmpStructLayout compares the Go struct libffi will read bytes from against
// the cif struct descriptor that says how to read them, member by member.
//
// Equal total size is not enough: the Go struct and the descriptor are written
// and maintained by hand in two separate places, so one can gain a field the
// other did not. That is silent when the difference is absorbed by interior
// padding, which is exactly how hybridgroup/yzma#289 presented.
func cmpStructLayout(ct FfiType, goType types.Type, slot, expr, cRaw string) string {
	if ct.Kind != KindStruct || goType == nil {
		return ""
	}

	if _, ok := goType.Underlying().(*types.Struct); !ok {
		return ""
	}

	gl, gOK := goLeaves(goType)
	if !gOK {
		noteSkip("%s %s: Go struct %s member layout not resolvable - LAYOUT NOT VERIFIED", slot, expr, goType.String())
		return ""
	}

	fl, fOK := ffiLeaves(ct)
	if !fOK {
		noteSkip("%s %s: cif descriptor %s member layout not resolvable - LAYOUT NOT VERIFIED", slot, expr, ct.Name)
		return ""
	}

	layoutChecked++

	goSize := int(arm64Sizes.Sizeof(goType))
	if tail := goTail(gl, fl, goSize, ct.Size); tail != "" {
		structCmp = append(structCmp, fmt.Sprintf("%s Go %s=%dB vs cif %s=%dB%s",
			slot, goType.String(), goSize, ct.Name, ct.Size, tail))
	}

	d := diffLeafPrefix(gl, fl, "Go "+goType.String(), "cif "+ct.Name, ct.Size)

	// Two members of the same width and ABI class can be swapped without moving
	// a single offset, and the C library then simply receives each value in the
	// other's place. Only the member names can see that, and only the header has
	// names to compare against: a cif descriptor carries none.
	if cl, ok := cLeaves(cstructOf(cRaw)); ok {
		sw, unmatched := diffSwappedMembers(gl, cl, goType.String(), squash(cRaw))
		d = append(d, sw...)

		for _, u := range unmatched {
			noteSkip("%s %s: %s - NOT VERIFIED against a transposition", slot, expr, u)
		}
	} else if cRaw != "" {
		noteSkip("%s %s: C struct %q member names not resolvable - NOT VERIFIED against a transposition", slot, expr, squash(cRaw))
	}

	if len(d) == 0 {
		layoutClean++
		return ""
	}

	return fmt.Sprintf("%s: Go struct %s layout differs from cif descriptor %s (%s): %s",
		slot, goType.String(), ct.Name, expr, strings.Join(d, "; "))
}

// kindCompat: same-width integer signedness is ABI-identical on arm64, so only
// flag cross-class confusion (int vs float vs pointer vs struct).
func kindCompat(a, b Kind) bool {
	cls := func(k Kind) int {
		switch k {
		case KindSint, KindUint:
			return 1
		case KindPointer:
			return 1 // pointer and 8-byte int are ABI-identical in integer regs
		case KindFloat, KindDouble:
			return 2
		case KindStruct:
			return 3
		case KindVoid:
			return 4
		}
		return 0
	}
	ca, cb := cls(a), cls(b)
	if ca == 0 || cb == 0 {
		return true
	}
	return ca == cb
}

// checkRet implements Rule 3.
func checkRet(ret FfiType, cs CallSite, cRaw string) string {
	switch ret.Kind {
	case KindVoid:
		if !cs.RetNil && cs.RetExpr != "" {
			return fmt.Sprintf("cif return type is TypeVoid but call passes a return buffer %q: libffi writes nothing, the Go variable keeps its zero value", cs.RetExpr)
		}
		return ""
	case KindSint, KindUint, KindPointer:
		if cs.RetNil || cs.RetExpr == "" {
			return "" // discarding a value is safe
		}
		if cs.RetType == nil {
			return "SKIP: integer return buffer " + cs.RetExpr + " type not resolvable"
		}
		_, sz := goKindOf(cs.RetType)
		if sz < 8 {
			return fmt.Sprintf("integer return buffer %s is %s (%dB); libffi always stores a full 8-byte ffi_arg, overwriting %d adjacent byte(s) - must be ffi.Arg",
				cs.RetExpr, cs.RetType.String(), sz, 8-sz)
		}
		return ""
	case KindFloat, KindDouble:
		if cs.RetNil || cs.RetType == nil {
			return ""
		}
		k, sz := goKindOf(cs.RetType)
		if ret.Kind == KindFloat && (k != KindFloat || sz != 4) {
			return fmt.Sprintf("float return buffer %s is %s (%dB); must be float32 (ffi.Arg is WRONG for floats)", cs.RetExpr, cs.RetType.String(), sz)
		}
		if ret.Kind == KindDouble && (k != KindDouble || sz != 8) {
			return fmt.Sprintf("double return buffer %s is %s; must be float64", cs.RetExpr, cs.RetType.String())
		}
		return ""
	case KindStruct:
		if cs.RetNil || cs.RetType == nil || ret.Size < 0 {
			return ""
		}
		// As for a struct argument, the cif size is how far libffi writes, so
		// a larger Go buffer is headroom and only a short one is overrun.
		_, sz := goKindOf(cs.RetType)
		if sz < ret.Size {
			return fmt.Sprintf("struct return buffer %s is %s (only %dB) but libffi writes the %dB of cif struct descriptor %s", cs.RetExpr, cs.RetType.String(), sz, ret.Size, ret.Name)
		}
		// A struct return is where the padding-absorbed field drift of
		// hybridgroup/yzma#289 landed: llama_context_default_params returns
		// llama_context_params by value into a Go ContextParams.
		return cmpStructLayout(ret, cs.RetType, "ret", cs.RetExpr, cRaw)
	}
	return ""
}

func shortPos(p token.Position) string {
	d := filepath.Base(filepath.Dir(p.Filename))
	return fmt.Sprintf("%s/%s:%d", d, filepath.Base(p.Filename), p.Line)
}
