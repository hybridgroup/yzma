package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
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
	CFuncs     map[string]CFunc
	NoC        []string
	Unresolved []string
	Skips      []string
	StructCmp  []string

	TotalDecls, Unparsed            int
	Matched, NCalls                 int
	CheckedR1, CheckedR2, CheckedR3 int
	CleanR1, CleanR2, CleanR3       int
}

// analyse runs the three rules over the packages named by patterns in the
// module rooted at yzmaRoot, against the headers in headerDir.
func analyse(yzmaRoot, headerDir string, patterns []string) (*report, error) {
	// Parser state is package-level and iterated to a fixpoint, so reset it
	// to keep repeated calls within one process independent.
	skips, structCmp = nil, nil
	resetCTypes()

	// --- C side ---
	cfuncs := map[string]CFunc{}
	unparsed := 0
	totalDecls := 0
	macros := []string{"LLAMA_API", "GGML_API", "GGML_BACKEND_API", "MTMD_API"}
	for _, hf := range headerFiles {
		path := filepath.Join(headerDir, hf.local)
		raw, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SKIP header %s: %v\n", path, err)
			continue
		}
		src := unwrapDeprecated(stripComments(string(raw)))
		collectTypedefs(src)
		collectStructs(src)
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
	}
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

	for _, b := range all {
		short := shortPos(b.PrepPos)
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
			if !cf.Vararg && len(b.Args) != len(cf.Params) {
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
			r3 := checkRet(b.Ret, cs)
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
				if ct.Size <= 0 && ct.Kind != KindVoid {
					noteSkip("RULE2 %s arg%d: cif descriptor %s size unknown - NOT VERIFIED", b.CName, i, ct.Name)
					continue
				}
				if ca.Size != ct.Size {
					probs = append(probs, fmt.Sprintf("arg%d: cif %s wants %dB, Go %s is %dB (%s)",
						i, ct.Name, ct.Size, ca.Pointee.String(), ca.Size, ca.Expr))
					continue
				}
				if !kindCompat(ct.Kind, ca.Kind) {
					probs = append(probs, fmt.Sprintf("arg%d: cif kind %s vs Go kind %s (%s / %s)",
						i, ct.Kind, ca.Kind, ca.Expr, ca.Pointee.String()))
				}
				cleanR2++
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

	sort.SliceStable(viols, func(i, j int) bool { return viols[i].Rule < viols[j].Rule })

	return &report{
		Viols:      viols,
		Bindings:   all,
		CFuncs:     cfuncs,
		NoC:        noCList,
		Unresolved: unresolvedList,
		Skips:      skips,
		StructCmp:  structCmp,
		TotalDecls: totalDecls,
		Unparsed:   unparsed,
		Matched:    matched,
		NCalls:     nCalls,
		CheckedR1:  checkedR1,
		CheckedR2:  checkedR2,
		CheckedR3:  checkedR3,
		CleanR1:    cleanR1,
		CleanR2:    cleanR2,
		CleanR3:    cleanR3,
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
	fmt.Println("================ NOT VERIFIED (skips) ================")
	for _, k := range r.Skips {
		fmt.Println("  ", k)
	}
	fmt.Printf("(%d skips)\n", len(r.Skips))
	fmt.Println("================ STRUCT-BY-VALUE COMPARISONS ================")
	for _, k := range r.StructCmp {
		fmt.Println("  ", k)
	}
	fmt.Printf("(%d struct-by-value slots compared)\n", len(r.StructCmp))
	fmt.Println("================ SUMMARY ================")
	fmt.Printf("C declarations found:      %d (unparsed: %d)\n", r.TotalDecls, r.Unparsed)
	fmt.Printf("distinct C functions:      %d\n", len(r.CFuncs))
	fmt.Printf("yzma bindings (Prep):      %d\n", len(r.Bindings))
	fmt.Printf("  matched to a C decl:     %d\n", r.Matched)
	fmt.Printf("  no C decl in headers:    %d\n", len(r.NoC))
	fmt.Printf("Call sites analysed:       %d\n", r.NCalls)
	fmt.Printf("Rule1 checked/clean:       %d / %d\n", r.CheckedR1, r.CleanR1)
	fmt.Printf("Rule2 arg slots checked:   %d / %d clean (unresolvable exprs: %d)\n", r.CheckedR2, r.CleanR2, len(r.Unresolved))
	fmt.Printf("Rule3 return bufs checked: %d / %d clean\n", r.CheckedR3, r.CleanR3)
	nr := map[int]int{}
	for _, v := range r.Viols {
		nr[v.Rule]++
	}
	fmt.Printf("violations: rule1=%d rule2=%d rule3=%d\n", nr[1], nr[2], nr[3])
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
		if strings.Contains(s, "arity") {
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
		structCmp = append(structCmp, fmt.Sprintf("%s %s: cif %s=%dB vs C %s=%dB", slot, "structcmp", t.Name, t.Size, squash(cRaw), cSize))
		if t.Size >= 0 && t.Size != cSize {
			return fmt.Sprintf("%s: struct descriptor %s is %dB, C %s is %dB", slot, t.Name, t.Size, squash(cRaw), cSize)
		}
		return ""
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
func checkRet(ret FfiType, cs CallSite) string {
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
		_, sz := goKindOf(cs.RetType)
		if sz != ret.Size {
			return fmt.Sprintf("struct return buffer %s is %s (%dB) but cif struct descriptor %s is %dB", cs.RetExpr, cs.RetType.String(), sz, ret.Name, ret.Size)
		}
		return ""
	}
	return ""
}

func shortPos(p token.Position) string {
	d := filepath.Base(filepath.Dir(p.Filename))
	return fmt.Sprintf("%s/%s:%d", d, filepath.Base(p.Filename), p.Line)
}
