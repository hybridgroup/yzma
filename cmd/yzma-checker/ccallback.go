package main

import (
	"fmt"
	"go/token"
	"go/types"
	"strings"
	"unicode"
)

// RULE 5 - the other direction.
//
// Rules 1-3 check calls yzma makes into C. A callback is the reverse: C puts
// arguments on the stack and libffi (or purego) hands them to a Go closure. A
// width error there does not corrupt one call's argument, it makes the closure
// read C stack memory instead of Go memory on every single invocation, and a
// wrong return width makes C read a partly-uninitialised value back. Nothing
// about either is visible to a compiler, and none of these sites goes through
// lib.Prep, so rules 1-3 never see them.
//
// The C side of the comparison is the function-pointer typedef the header
// declares for the callback, parsed here into the same CFunc the rest of the
// tool already compares against.

// cfnptrs maps a function-pointer typedef name to the signature it declares.
var cfnptrs = map[string]CFunc{}

// cfnptrBad records the function-pointer typedefs this parser could not take
// apart, keyed by the raw declaration. They are counted on the RULE 5 summary
// line, and a Go callback that links to one becomes a skip rather than a silent
// pass - dropping it would leave the site looking checked when it is not.
var cfnptrBad = map[string]string{}

func resetCCallbacks() {
	clear(cfnptrs)
	clear(cfnptrBad)
	clear(callbacksByID)
	fnPtrStores, fnPtrStructs = nil, nil
	fnptrChecked, fnptrClean, fnptrTraced = 0, 0, 0
}

// collectFnPtrTypedefs parses "typedef RET (*NAME)(params);" out of a header.
//
// collectTypedefs already sees these declarations, but only far enough to know
// that the name is a pointer; the signature it points at is thrown away there,
// because for every other rule a callback parameter is just 8 bytes.
func collectFnPtrTypedefs(path, src string) {
	for _, idx := range findAll(src, "typedef") {
		if idx > 0 && isIdentChar(src[idx-1]) {
			continue
		}

		if after := idx + len("typedef"); after < len(src) && isIdentChar(src[after]) {
			continue
		}

		// Terminating ';' at brace depth 0, as collectTypedefs does: an enum or
		// struct body may contain its own semicolons.
		depth, end := 0, -1
		for k := idx; k < len(src) && end < 0; k++ {
			switch src[k] {
			case '{':
				depth++
			case '}':
				depth--
			case ';':
				if depth == 0 {
					end = k
				}
			}
		}
		if end < 0 {
			continue
		}

		decl := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(src[idx:end]), "typedef"))
		if !strings.Contains(decl, "(") {
			continue
		}

		name, fn, why := parseFnPtrTypedef(decl)
		switch {
		case name == "" && why == "":
			continue // not a function-pointer typedef at all
		case why != "":
			cfnptrBad[squash(decl)] = why
		default:
			fn.File = path
			fn.Line = 1 + strings.Count(src[:idx], "\n")
			fn.Raw = squash(decl)
			if _, dup := cfnptrs[name]; !dup {
				cfnptrs[name] = fn
			}
		}
	}
}

// parseFnPtrTypedef takes apart the body of one typedef, with "typedef" already
// stripped. It returns an empty name and an empty reason for a declaration that
// is not a function-pointer typedef, and a reason for one that is but could not
// be read.
func parseFnPtrTypedef(decl string) (string, CFunc, string) {
	var fn CFunc

	open := strings.Index(decl, "(")
	if open < 0 {
		return "", fn, ""
	}

	rest := decl[open+1:]
	close := strings.Index(rest, ")")
	if close < 0 {
		return "", fn, ""
	}

	// The declarator is "(*name)" or, as mtmd.h spells it, "(* name)".
	inner := strings.TrimSpace(rest[:close])
	if !strings.HasPrefix(inner, "*") {
		return "", fn, "" // a plain typedef whose type happens to be parenthesised
	}

	name := strings.TrimSpace(strings.TrimLeft(inner, "* "))
	if name == "" || !isIdentChar(name[0]) || strings.ContainsAny(name, " *") {
		return name, fn, "declarator " + squash(inner) + " is not a plain name"
	}

	// The parameter list follows immediately after the declarator.
	tail := strings.TrimSpace(rest[close+1:])
	if !strings.HasPrefix(tail, "(") {
		return name, fn, "no parameter list after the declarator"
	}

	depth, pclose := 0, -1
	for i := 0; i < len(tail) && pclose < 0; i++ {
		switch tail[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				pclose = i
			}
		}
	}
	if pclose < 0 {
		return name, fn, "unterminated parameter list"
	}

	fn.Name = name
	fn.RetRaw = squash(decl[:open])
	if fn.RetRaw == "" {
		return name, fn, "no return type"
	}

	fn.RetKind, fn.RetSize = classify(fn.RetRaw)

	body := strings.TrimSpace(tail[1:pclose])
	if body == "" || body == "void" {
		return name, fn, ""
	}

	for _, p := range splitTop(body) {
		ps := strings.TrimSpace(p)
		if ps == "..." {
			fn.Vararg = true
			continue
		}

		if ps == "" {
			return name, fn, "empty parameter in " + squash(body)
		}

		norm := stripDeclName(ps)
		k, sz := classify(norm)
		fn.Params = append(fn.Params, CParam{Raw: squash(ps), Norm: squash(norm), Kind: k, Size: sz})
	}

	return name, fn, ""
}

// Callback is one site where C calls back into Go: either a libffi closure
// described by ffi.PrepCif, or a purego.NewCallback closure whose Go signature
// *is* the descriptor.
type Callback struct {
	Form string // "ffi.PrepCif" or "purego.NewCallback"
	GoID string // cif variable, or the enclosing func for the purego form
	Fn   string // enclosing Go func
	Pkg  string
	Pos  token.Position

	// ffi.PrepCif form.
	Ret    FfiType
	Args   []FfiType
	NFixed int

	// purego.NewCallback form.
	Sig *types.Signature

	// CTypedef is the function-pointer typedef this site must match, and Link
	// how that was established. CTypedef is empty when no link could be made,
	// which is a skip rather than a pass.
	CTypedef string
	Link     string
}

// callbackTags names the C typedef for a callback site whose Go identifier does
// not normalise onto it. Like memberAliases and constEnumTags it is a closed,
// hand-checked list on purpose: a matcher loose enough to guess a name is loose
// enough to validate a closure against the wrong signature, so a new site costs
// one reviewed line here rather than passing silently. Anything neither listed
// nor normalising onto exactly one typedef is reported as a skip.
var callbackTags = map[string]string{
	// LogSilent is named for what it does, not for what it implements.
	"llama.LogSilent": "ggml_log_callback",
}

// cPrefixes are the library prefixes a yzma identifier drops.
var cPrefixes = []string{"llama_", "ggml_", "mtmd_"}

// linkCallback resolves the C typedef a callback site must match.
//
// The tag table decides first. Otherwise the site's own identifier is
// normalised - progressCallbackCif -> progresscallback, newAbortCallback ->
// abortcallback - and matched against every function-pointer typedef name with
// its library prefix removed. Exactly one match is a link; several are
// disambiguated by the prefix belonging to the Go package (pkg/mtmd owns
// mtmd_progress_callback, pkg/llama owns llama_progress_callback) and are
// otherwise left unresolved.
func linkCallback(cb *Callback) {
	pkg := cb.Pkg
	if i := strings.LastIndex(pkg, "/"); i >= 0 {
		pkg = pkg[i+1:]
	}

	if td, ok := callbackTags[pkg+"."+cb.GoID]; ok {
		cb.CTypedef, cb.Link = td, "callbackTags"
		return
	}

	want := normCallbackName(cb.GoID)

	var hits []string
	for name := range cfnptrs {
		if normCallbackName(trimCPrefix(name)) == want {
			hits = append(hits, name)
		}
	}
	for raw, why := range cfnptrBad {
		// An unparseable typedef still has a name, and a site normalising onto
		// it must not fall through to some other typedef that happens to match.
		if n, _, _ := parseFnPtrTypedef(raw); n != "" && normCallbackName(trimCPrefix(n)) == want {
			cb.Link = fmt.Sprintf("C typedef %s could not be parsed: %s", n, why)
			return
		}
	}

	switch len(hits) {
	case 0:
		cb.Link = fmt.Sprintf("identifier %s normalises to %q, which names no function-pointer typedef", cb.GoID, want)
	case 1:
		cb.CTypedef, cb.Link = hits[0], "identifier "+cb.GoID
	default:
		var owned []string
		for _, h := range hits {
			if strings.HasPrefix(h, pkg+"_") {
				owned = append(owned, h)
			}
		}
		if len(owned) == 1 {
			cb.CTypedef, cb.Link = owned[0], "identifier "+cb.GoID+" in package "+pkg
			return
		}

		cb.Link = fmt.Sprintf("identifier %s normalises to %q, which is ambiguous between %s",
			cb.GoID, want, strings.Join(hits, ", "))
	}
}

func trimCPrefix(s string) string {
	for _, p := range cPrefixes {
		if strings.HasPrefix(s, p) {
			return s[len(p):]
		}
	}

	return s
}

// normCallbackName reduces a callback identifier to a form comparable across
// the two naming conventions. The affixes dropped are the ones yzma adds around
// the C name: a "Cif" suffix for the descriptor variable and a "new" prefix for
// the constructor func.
func normCallbackName(s string) string {
	s = strings.TrimSuffix(s, "Cif")
	s = strings.TrimPrefix(s, "new")

	var b strings.Builder
	for _, r := range s {
		if r == '_' {
			continue
		}

		b.WriteRune(unicode.ToLower(r))
	}

	return b.String()
}

// checkCallback implements RULE 5 for one site. It returns the problems found,
// or nil for a clean site; an unresolved link is reported as a skip and checks
// nothing.
func checkCallback(cb *Callback) (probs []string, checked bool) {
	if cb.CTypedef == "" {
		noteSkip("RULE5 %s (%s): %s - NOT VERIFIED", cb.GoID, shortPos(cb.Pos), cb.Link)
		return nil, false
	}

	cf, ok := cfnptrs[cb.CTypedef]
	if !ok {
		noteSkip("RULE5 %s (%s): C typedef %s not found - NOT VERIFIED", cb.GoID, shortPos(cb.Pos), cb.CTypedef)
		return nil, false
	}

	if cb.Form == "purego.NewCallback" {
		return checkPuregoCallback(cb, cf), true
	}

	return checkCifCallback(cb, cf), true
}

// checkCifCallback compares an ffi.PrepCif descriptor against the typedef, the
// same comparison RULE 1 makes for a call in the other direction: libffi reads
// the C arguments through this descriptor, so a wrong width here reads the
// wrong bytes off the C stack on every invocation.
func checkCifCallback(cb *Callback, cf CFunc) []string {
	var probs []string

	// cmpTypes labels a skip with the slot alone, so the typedef name goes in
	// front of it: otherwise a NOT VERIFIED line would not say which callback it
	// belongs to.
	slot := func(s string) string { return cf.Name + " " + s }

	// nfixed is what libffi will believe about the argument count, so it is
	// checked separately from the descriptor list it indexes into.
	if cb.NFixed >= 0 && cb.NFixed != len(cf.Params) {
		probs = append(probs, fmt.Sprintf("nfixed is %d but C %s takes %d parameter(s)",
			cb.NFixed, cf.Name, len(cf.Params)))
	}

	if len(cb.Args) != len(cf.Params) {
		probs = append(probs, fmt.Sprintf("descriptor has %d arg type(s), C %s takes %d",
			len(cb.Args), cf.Name, len(cf.Params)))
	}

	if p := cmpTypes(cb.Ret, cf.RetRaw, cf.RetKind, cf.RetSize, slot("ret")); p != "" {
		probs = append(probs, p)
	}

	for i := range min(len(cb.Args), len(cf.Params)) {
		p := cf.Params[i]
		if d := cmpTypes(cb.Args[i], p.Norm, p.Kind, p.Size, slot(fmt.Sprintf("arg%d", i))); d != "" {
			probs = append(probs, d)
		}
	}

	return probs
}

// checkPuregoCallback compares a purego closure's Go signature against the
// typedef. purego hands every C argument over in a register, so each parameter
// must be a type it can carry there, and each must agree with the C parameter
// in width and ABI class: a Go int32 where C passes a pointer reads half a
// pointer, and one missing parameter shifts every later one.
func checkPuregoCallback(cb *Callback, cf CFunc) []string {
	var probs []string

	if cb.Sig == nil {
		noteSkip("RULE5 %s (%s): closure signature not resolvable - NOT VERIFIED", cb.GoID, shortPos(cb.Pos))
		return nil
	}

	ps := cb.Sig.Params()
	if ps.Len() != len(cf.Params) {
		probs = append(probs, fmt.Sprintf("closure takes %d parameter(s) but C %s passes %d: every parameter after the first missing one receives the wrong value",
			ps.Len(), cf.Name, len(cf.Params)))
	}

	for i := range min(ps.Len(), len(cf.Params)) {
		cp := cf.Params[i]
		gt := ps.At(i).Type()
		gk, gsz := goKindOf(gt)

		if cp.Kind == KindUnknown || cp.Size < 0 {
			noteSkip("RULE5 %s arg%d: C type %q unclassifiable - NOT VERIFIED", cf.Name, i, squash(cp.Norm))
			continue
		}

		switch {
		case gk == KindUnknown || gsz < 0:
			probs = append(probs, fmt.Sprintf("arg%d: Go %s is not a type purego can receive in a register",
				i, gt.String()))
		case gsz > 8 || gk == KindStruct:
			probs = append(probs, fmt.Sprintf("arg%d: Go %s is %dB; purego passes every argument in a register, so it cannot carry it",
				i, gt.String(), gsz))
		case gsz != cp.Size:
			probs = append(probs, fmt.Sprintf("arg%d: C %s is %dB but Go %s is %dB",
				i, squash(cp.Norm), cp.Size, gt.String(), gsz))
		case !kindCompat(gk, cp.Kind):
			probs = append(probs, fmt.Sprintf("arg%d: C %s is %s but Go %s is %s",
				i, squash(cp.Norm), cp.Kind, gt.String(), gk))
		}
	}

	return append(probs, puregoRet(cb.Sig, cf)...)
}

// puregoRet checks the closure's result against the typedef's return type.
//
// A Go result wider than C reads is harmless - C reads the low bytes of the
// register the closure filled - and so is a result C ignores entirely, which is
// how yzma's void-returning log callback is written. A result that is *narrower*
// than what C reads, or of the wrong ABI class, is not: C then reads bytes the
// closure never wrote, or reads an integer register for a float.
func puregoRet(sig *types.Signature, cf CFunc) []string {
	if cf.RetKind == KindVoid {
		return nil
	}

	rs := sig.Results()
	if rs.Len() == 0 {
		return []string{fmt.Sprintf("closure returns nothing but C %s returns %s: C reads an uninitialised register",
			cf.Name, squash(cf.RetRaw))}
	}

	if rs.Len() != 1 {
		return []string{fmt.Sprintf("closure returns %d values; a C callback returns one", rs.Len())}
	}

	rt := rs.At(0).Type()
	gk, gsz := goKindOf(rt)

	switch {
	case cf.RetKind == KindUnknown || cf.RetSize < 0:
		noteSkip("RULE5 %s ret: C type %q unclassifiable - NOT VERIFIED", cf.Name, squash(cf.RetRaw))
	case !kindCompat(gk, cf.RetKind):
		return []string{fmt.Sprintf("ret: C %s is %s but closure returns %s (%s)",
			squash(cf.RetRaw), cf.RetKind, rt.String(), gk)}
	case gsz < cf.RetSize:
		return []string{fmt.Sprintf("ret: C reads %dB of %s but closure returns %s (%dB), leaving %dB uninitialised",
			cf.RetSize, squash(cf.RetRaw), rt.String(), gsz, cf.RetSize-gsz)}
	}

	return nil
}

// RULE 5, third form - the function-pointer struct member.
//
// The two forms above are where a callback is *built*. A C struct member of
// function-pointer type is where one is *installed*, and it is the form nothing
// looked at: llama_model_params.progress_callback, llama_context_params.cb_eval
// and .abort_callback, mtmd_context_params' two, are all function pointers C
// will jump through, and to layout.go's Leaf each is 8 bytes of pointer class
// like any other member. So two things were silent.
//
// A Go `func` field is also 8 bytes of pointer class, and it is a pointer to a
// func *descriptor* rather than to code, so C jumping through it executes
// whatever the descriptor's first word happens to be. Every offset, width and
// class comparison passes it.
//
// And when the field is the uintptr yzma actually declares, nothing says which
// callback belongs in it. Storing the log callback's code pointer into cb_eval
// type-checks, lays out identically, and makes C unpack a ggml_backend_sched
// _eval_callback's arguments through a ggml_log_callback descriptor on every
// graph node. So the code pointer is traced back to the site that produced it
// and the typedef that site was linked to must be the member's own - which,
// with the descriptor comparison above, is what makes the signature C calls
// through the member checked rather than assumed.

// FnPtrStore is one assignment of a callback's code pointer into a struct field.
type FnPtrStore struct {
	GoStruct types.Type // the struct the field belongs to
	Field    string     // Go field name
	Expr     string
	Pkg      string // the package the assignment is in
	Site     string // the callback site's GoID: a cif var, or a constructor func
	Pos      token.Position
}

var fnPtrStores []*FnPtrStore

// callbacksByID indexes the linked callback sites by the identifier a store
// resolves to, so a stored code pointer can be traced to the typedef RULE 5
// already compared its descriptor against.
//
// Keyed by package as well as identifier, because pkg/llama and pkg/mtmd give
// their progress-callback closure the same variable names while the typedefs
// behind them are llama_progress_callback and mtmd_progress_callback: a bare
// identifier key resolves one package's store to the other package's callback,
// which reports two correct sites as defects.
var callbacksByID = map[string]*Callback{}

func callbackID(pkg, id string) string { return pkg + "." + id }

// fnPtrStruct is one struct-by-value slot's pair of representations, recorded
// where both are known. Only the pair matters, not the slot: the same struct
// appears in many slots, and one finding per slot would be the same defect
// printed five times.
type fnPtrStruct struct {
	GoType types.Type
	CS     *CStruct
	CRaw   string
	Pos    string
}

var fnPtrStructs []fnPtrStruct

// Function-pointer members compared, counted separately for the same reason the
// pointer targets are: the claim is only made where the C member names a typedef
// this parser took apart and the Go struct has a member at that offset, and
// Traced is the narrower subset whose stored code pointer could be followed back
// to a callback site. A narrowing of either would otherwise leave RULE 5
// reporting zero findings without reporting that it stopped looking.
var fnptrChecked, fnptrClean, fnptrTraced int

// noteFnPtrStruct records a struct-by-value slot for the member scan, once per
// pair of representations.
func noteFnPtrStruct(goType types.Type, cRaw, pos string) {
	cs := cstructOf(cRaw)
	if goType == nil || cs == nil {
		return
	}

	// Keyed on the resolved struct rather than the parameter text, which spells
	// the same struct "struct mtmd_context_params" in one slot and "const struct
	// mtmd_context_params" in the next.
	for _, s := range fnPtrStructs {
		if s.CS == cs && types.Identical(s.GoType, goType) {
			return
		}
	}

	fnPtrStructs = append(fnPtrStructs, fnPtrStruct{GoType: goType, CS: cs, CRaw: cRaw, Pos: pos})
}

// checkFnPtrMembers compares every function-pointer member of one C struct
// against the Go member that occupies the same offset.
//
// The offset correspondence comes from the same flattening the layout diff uses,
// so there is one struct walker rather than two: a C member's offset is
// CStruct.Offs, and the Go leaf at that offset is the field libffi will copy the
// bytes of.
func checkFnPtrMembers(s fnPtrStruct) []violation {
	cs := s.CS
	if cs == nil || len(cs.Offs) != len(cs.Fields) {
		return nil
	}

	gl, ok := goLeaves(s.GoType)
	if !ok {
		return nil // already a LAYOUT NOT VERIFIED skip from cmpStructLayout
	}

	byOff := make(map[int]Leaf, len(gl))
	for _, l := range gl {
		byOff[l.Off] = l
	}

	var viols []violation
	for i, f := range cs.Fields {
		td := strings.TrimSpace(strings.ReplaceAll(f.Norm, "const", " "))

		cf, known := cfnptrs[td]
		if !known {
			if why, bad := cfnptrBadOf(td); bad {
				noteSkip("RULE5 %s.%s: C typedef %s could not be parsed: %s - NOT VERIFIED",
					cs.Name, f.Name, td, why)
			}

			continue
		}

		if f.Cnt != 1 {
			noteSkip("RULE5 %s.%s: an array of %d function pointers is not one member - NOT VERIFIED",
				cs.Name, f.Name, f.Cnt)
			continue
		}

		gf, at := byOff[cs.Offs[i]]
		if !at {
			// No Go member at that offset at all. cmpStructLayout owns that
			// finding, and repeating it here would bury the insertion point.
			continue
		}

		fnptrChecked++

		probs := cmpFnPtrMember(cf, s.GoType, gf, cs.Name, f.Name)
		if len(probs) == 0 {
			fnptrClean++
			continue
		}

		// Named for the member rather than for the typedef, which is also the
		// name of the callback *site* RULE 5 already reports on: two subjects
		// under one name would make a report line ambiguous about which of them
		// is wrong.
		viols = append(viols, violation{5, cs.Name + "." + f.Name, s.Pos,
			strings.Join(probs, "; ") + " [C calls it as " + cf.Name + "]", sev(probs)})
	}

	return viols
}

// cmpFnPtrMember checks one member: that the Go field can hold a code pointer at
// all, and that the code pointer stored in it came from the callback site
// implementing this member's typedef.
func cmpFnPtrMember(cf CFunc, goType types.Type, gf Leaf, cName, cField string) []string {
	var probs []string

	where := fmt.Sprintf("Go %s.%s / C %s.%s", goType.String(), gf.Path, cName, cField)

	if gf.GoType != nil {
		if _, isFunc := gf.GoType.Underlying().(*types.Signature); isFunc {
			probs = append(probs, fmt.Sprintf("%s: C calls %s through this member, but Go %s is a func value - 8 bytes of pointer class like a code pointer, and a pointer to a func descriptor rather than to code",
				where, cf.Name, gf.GoType.String()))
		}
	}

	// Only a top-level field can be named by an assignment this tool traced; a
	// nested one is one selector further away than collectFnPtrStores looks, and
	// no bound struct has one.
	if strings.Contains(gf.Path, ".") {
		return probs
	}

	for _, st := range fnPtrStores {
		if st.Field != gf.Path || !types.Identical(st.GoStruct, goType) {
			continue
		}

		cb, isCallback := callbacksByID[callbackID(st.Pkg, st.Site)]
		if !isCallback {
			// Not a callback code pointer, so there is no signature behind it to
			// compare - the same way a void * names no pointer target. Outside
			// the claim rather than a skip, which is what fnptrTraced counts.
			continue
		}

		fnptrTraced++

		switch {
		case cb.CTypedef == "":
			noteSkip("RULE5 %s: code pointer from %s, whose own typedef is unresolved (%s) - NOT VERIFIED",
				where, st.Site, cb.Link)
		case cb.CTypedef != cf.Name:
			probs = append(probs, fmt.Sprintf("%s: the code pointer stored here (%s) implements %s, but C calls this member as %s: every argument is unpacked through the wrong descriptor",
				where, st.Expr, cb.CTypedef, cf.Name))
		}
	}

	return probs
}

// cfnptrBadOf reports whether a type name is a function-pointer typedef this
// parser could not take apart. cfnptrBad is keyed by the raw declaration,
// because an unparseable one may have no usable name, so the name is recovered
// the way linkCallback recovers it. Without this a member whose typedef failed
// to parse would look like a member that is not a function pointer at all.
func cfnptrBadOf(name string) (string, bool) {
	if name == "" {
		return "", false
	}

	for raw, why := range cfnptrBad {
		if n, _, _ := parseFnPtrTypedef(raw); n == name {
			return why, true
		}
	}

	return "", false
}
