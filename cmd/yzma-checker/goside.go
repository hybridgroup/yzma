package main

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// FfiType is a resolved libffi type descriptor.
type FfiType struct {
	Name  string // e.g. TypeUint64, or the Go var name for struct descriptors
	Kind  Kind
	Size  int
	Elems []FfiType // for structs
	Offs  []int     // byte offset of each element, parallel to Elems; nil if Size is -1
}

func (t FfiType) String() string { return fmt.Sprintf("%s{%s/%d}", t.Name, t.Kind, t.Size) }

var ffiScalars = map[string]FfiType{
	"TypeVoid":       {Name: "TypeVoid", Kind: KindVoid, Size: 0},
	"TypeUint8":      {Name: "TypeUint8", Kind: KindUint, Size: 1},
	"TypeSint8":      {Name: "TypeSint8", Kind: KindSint, Size: 1},
	"TypeUint16":     {Name: "TypeUint16", Kind: KindUint, Size: 2},
	"TypeSint16":     {Name: "TypeSint16", Kind: KindSint, Size: 2},
	"TypeUint32":     {Name: "TypeUint32", Kind: KindUint, Size: 4},
	"TypeSint32":     {Name: "TypeSint32", Kind: KindSint, Size: 4},
	"TypeUint64":     {Name: "TypeUint64", Kind: KindUint, Size: 8},
	"TypeSint64":     {Name: "TypeSint64", Kind: KindSint, Size: 8},
	"TypeFloat":      {Name: "TypeFloat", Kind: KindFloat, Size: 4},
	"TypeDouble":     {Name: "TypeDouble", Kind: KindDouble, Size: 8},
	"TypePointer":    {Name: "TypePointer", Kind: KindPointer, Size: 8},
	"TypeLongdouble": {Name: "TypeLongdouble", Kind: KindUnknown, Size: 16},
}

// Binding is one yzma ffi.Fun: its Prep spec plus every Call site.
type Binding struct {
	GoVar   string // Go package-level var name
	CName   string // C symbol
	PrepPos token.Position
	Ret     FfiType
	Args    []FfiType
	NFixed  int  // for PrepVar
	Variadi bool //
	Calls   []CallSite
	Pkg     string
}

// CallSite is one <var>.Call(...) invocation.
type CallSite struct {
	Pos token.Position
	Fn  string // enclosing Go func

	// FnExported and FnDeprecated describe that enclosing func as a consumer sees
	// it: whether it is callable from outside the package at all, and whether its
	// doc comment carries a `Deprecated:` paragraph. Both are here because the
	// enclosing func of a call site *is* the wrapper yzma exposes for the C symbol
	// - see checkDeprecationNote.
	FnExported   bool
	FnDeprecated bool

	RetExpr string
	RetType types.Type // pointee type of the return buffer, nil if untrackable
	RetNil  bool
	Args    []CallArg
}

// CallArg is one avalue slot at a call site.
type CallArg struct {
	Expr    string
	Pointee types.Type // the Go type libffi will read bytes from
	Size    int        // -1 if unknown
	Kind    Kind
	Note    string

	// Str is what could be worked out about the buffer behind this slot, for
	// the C strings among them. See stringBuf and cmpStringTerm.
	Str stringBuf
}

// stringBuf is how a Go buffer handed to a C `char *` was produced, and whether
// anything in that production terminates it.
//
// Width, class and pointer target all match for a Go byte buffer behind a `const
// char *`, so no rule here looks at the one property C actually depends on: that
// the bytes end in a NUL. Nothing in Go produces one - a Go string carries its
// length - so the terminator is always appended by hand, and dropping the
// `+ "\x00"` from `&[]byte(path + "\x00")[0]` leaves every rule passing while C
// reads forward off the end of a Go allocation.
//
// Producer is empty for a buffer this tool cannot reason about, which is not the
// same as an unterminated one: a `*byte` arriving as a function parameter or
// written by C names nothing to trace. Term is empty when the producer *is* one
// of the traced forms and no terminator was found, which is the finding.
type stringBuf struct {
	Producer string // the expression the buffer came from; "" if not one this traces
	Term     string // the evidence it ends in a NUL; "" if there is none
}

// nulHelpers are functions taken on inspection to return a NUL-terminated
// buffer. A closed hand-checked list, for the same reason constAliases and
// callbackTags are: a rule loose enough to guess which helper terminates its
// result is loose enough to accept one that does not, and a new helper should
// cost one reviewed line rather than pass silently. yzma has none today - every
// site appends the terminator inline - so this is empty on purpose.
var nulHelpers = map[string]string{}

// goTargets are the architectures yzma supports, the first of which is the one
// every width in this report is computed under.
//
// pkg/download/arch.go declares exactly these two and MustParseArch panics on
// anything else, and jupiterrider/ffi's build tag narrows it further to the same
// pair plus linux/riscv64 - all 64-bit - so there is no 32-bit target for a
// 4-byte pointer to arrive from. That the two agree used to be an argument in
// the README; goLeaves is run under each of them now so that it is a measurement
// instead. See diffArchLayouts.
var goTargets = []struct {
	Arch  string
	Sizes types.Sizes
}{
	{"arm64", types.SizesFor("gc", "arm64")},
	{"amd64", types.SizesFor("gc", "amd64")},
}

// goSizes is the layout in force. It is a variable rather than a constant
// because the cross-architecture comparison flattens the same struct under each
// target in turn, and goKindOf reads it rather than taking it as a parameter so
// that there stays one struct walker instead of two. See leavesUnder.
var goSizes = goTargets[0].Sizes

func goKindOf(t types.Type) (Kind, int) {
	if t == nil {
		return KindUnknown, -1
	}
	sz := int(goSizes.Sizeof(t))
	u := t.Underlying()
	switch b := u.(type) {
	case *types.Basic:
		switch {
		case b.Info()&types.IsUnsigned != 0:
			if b.Kind() == types.Uintptr || b.Kind() == types.UnsafePointer {
				return KindPointer, sz
			}
			return KindUint, sz
		case b.Info()&types.IsInteger != 0:
			return KindSint, sz
		case b.Kind() == types.Float32:
			return KindFloat, sz
		case b.Kind() == types.Float64:
			return KindDouble, sz
		case b.Kind() == types.Bool:
			return KindUint, sz
		case b.Kind() == types.UnsafePointer:
			return KindPointer, sz
		}
	case *types.Pointer, *types.Signature, *types.Map, *types.Chan:
		return KindPointer, sz
	case *types.Struct:
		return KindStruct, sz
	case *types.Slice, *types.Interface:
		return KindUnknown, sz
	}
	return KindUnknown, sz
}

type analyzer struct {
	pkg       *packages.Package
	fset      *token.FileSet
	ffiAlias  map[string]FfiType // package-level Go vars that alias an ffi.Type
	bindings  map[string]*Binding
	order     []string
	callbacks []*Callback // RULE 5: the sites where C calls back into Go

	// closureCode maps the Go variable holding a libffi closure's *code* address
	// to the cif that describes it, from ffi.PrepClosureLoc(closure, cif, fn,
	// nil, code). That is the one link between the value stored in a
	// function-pointer struct member and the descriptor C will unpack its
	// arguments through.
	closureCode map[string]string

	// stores are the assignments of such a code pointer into a struct field.
	stores []*FnPtrStore
}

func loadPkgs(dir string, patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports,
		Dir: dir,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

// resolveTypeExpr resolves an expression appearing in a Prep type list.
// Accepted forms: &ffi.TypeX, &localVar, ffi.NewType(...), &someStructTypeVar.
func (a *analyzer) resolveTypeExpr(e ast.Expr) FfiType {
	switch x := e.(type) {
	case *ast.UnaryExpr:
		if x.Op == token.AND {
			return a.resolveTypeExpr(x.X)
		}
	case *ast.SelectorExpr:
		if id, ok := x.X.(*ast.Ident); ok && id.Name == "ffi" {
			if t, ok := ffiScalars[x.Sel.Name]; ok {
				return t
			}
			return FfiType{Name: "ffi." + x.Sel.Name, Kind: KindUnknown, Size: -1}
		}
	case *ast.Ident:
		if t, ok := a.ffiAlias[x.Name]; ok {
			return t
		}
		return FfiType{Name: x.Name + "(unresolved)", Kind: KindUnknown, Size: -1}
	case *ast.CallExpr:
		if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
			if id, ok := sel.X.(*ast.Ident); ok && id.Name == "ffi" && sel.Sel.Name == "NewType" {
				return a.buildStruct("anon", x.Args)
			}
		}
	}
	return FfiType{Name: exprStr(e), Kind: KindUnknown, Size: -1}
}

// buildStruct computes the libffi/C aggregate size of an ffi.NewType(...) descriptor.
func (a *analyzer) buildStruct(name string, args []ast.Expr) FfiType {
	st := FfiType{Name: name, Kind: KindStruct}
	maxAlign := 1
	off := 0
	ok := true
	for _, e := range args {
		el := a.resolveTypeExpr(e)
		st.Elems = append(st.Elems, el)
		if el.Size <= 0 && el.Kind != KindVoid {
			ok = false
			continue
		}
		al := el.Size
		if el.Kind == KindStruct {
			al = structAlign(el)
		}
		if al > 8 {
			al = 16
		}
		if al > maxAlign {
			maxAlign = al
		}
		if off%al != 0 {
			off += al - off%al
		}
		st.Offs = append(st.Offs, off)
		off += el.Size
	}
	if !ok {
		st.Size = -1
		st.Offs = nil
		return st
	}
	if off%maxAlign != 0 {
		off += maxAlign - off%maxAlign
	}
	st.Size = off
	return st
}

func structAlign(t FfiType) int {
	if t.Kind != KindStruct {
		if t.Size > 8 {
			return 16
		}
		return t.Size
	}
	m := 1
	for _, e := range t.Elems {
		a := structAlign(e)
		if a > m {
			m = a
		}
	}
	return m
}

func exprStr(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return exprStr(x.X) + "." + x.Sel.Name
	case *ast.UnaryExpr:
		return x.Op.String() + exprStr(x.X)
	case *ast.CallExpr:
		var as []string
		for _, a := range x.Args {
			as = append(as, exprStr(a))
		}
		return exprStr(x.Fun) + "(" + strings.Join(as, ", ") + ")"
	case *ast.IndexExpr:
		return exprStr(x.X) + "[" + exprStr(x.Index) + "]"
	case *ast.BasicLit:
		return x.Value
	case *ast.StarExpr:
		return "*" + exprStr(x.X)
	case *ast.ParenExpr:
		return "(" + exprStr(x.X) + ")"
	case *ast.CompositeLit:
		return "composite{...}"
	case *ast.SliceExpr:
		return exprStr(x.X) + "[:]"
	case *ast.ArrayType:
		if x.Len == nil {
			return "[]" + exprStr(x.Elt)
		}
		return "[" + exprStr(x.Len) + "]" + exprStr(x.Elt)
	}
	return fmt.Sprintf("%T", e)
}

// pointeeOf figures out the Go type whose bytes libffi will read for an
// argument expression passed to Fun.Call.
func (a *analyzer) pointeeOf(e ast.Expr) (types.Type, string) {
	e = unwrapUnsafePointer(e)
	if id, ok := e.(*ast.Ident); ok && id.Name == "nil" {
		return nil, "nil"
	}
	t := a.pkg.TypesInfo.TypeOf(e)
	if t == nil {
		return nil, "no type info"
	}
	if p, ok := t.Underlying().(*types.Pointer); ok {
		return p.Elem(), ""
	}
	if b, ok := t.Underlying().(*types.Basic); ok && b.Kind() == types.UnsafePointer {
		return nil, "opaque unsafe.Pointer"
	}
	if b, ok := t.Underlying().(*types.Basic); ok && b.Kind() == types.Uintptr {
		return nil, "uintptr (not addressable)"
	}
	return nil, "non-pointer " + t.String()
}

// stringBufOf traces one avalue expression back to the buffer it points at.
//
// The avalue is the address of the pointer C receives, so `unsafe.Pointer(&file)`
// says nothing by itself: what matters is every value `file` is assigned inside
// the same function, which is where the terminator is either appended or
// forgotten.
func stringBufOf(e ast.Expr, body *ast.BlockStmt) stringBuf {
	u, ok := unparen(unwrapUnsafePointer(e)).(*ast.UnaryExpr)
	if !ok || u.Op != token.AND || body == nil {
		return stringBuf{}
	}

	id, ok := unparen(u.X).(*ast.Ident)
	if !ok {
		return stringBuf{}
	}

	return charBuf(id, body, 0)
}

// charBuf classifies an expression that produces a `char *`.
func charBuf(e ast.Expr, body *ast.BlockStmt, depth int) stringBuf {
	if depth > 4 {
		return stringBuf{}
	}

	switch x := unparen(unwrapUnsafePointer(e)).(type) {
	case *ast.Ident:
		// Merge every value the identifier is given: a producer with no
		// terminator anywhere on that list is the finding, so it wins over one
		// that has it.
		var found stringBuf
		for _, v := range assignedValues(x.Name, body) {
			got := charBuf(v, body, depth+1)
			if got.Producer == "" {
				continue
			}
			if got.Term == "" {
				return got
			}
			found = got
		}

		return found

	case *ast.UnaryExpr:
		// &buf[0], the idiom every yzma site uses.
		if x.Op != token.AND {
			return stringBuf{}
		}
		ix, ok := unparen(x.X).(*ast.IndexExpr)
		if !ok || exprStr(ix.Index) != "0" {
			return stringBuf{}
		}

		return byteSliceBuf(ix.X, body, depth+1)

	case *ast.CallExpr:
		switch callee := exprStr(x.Fun); {
		case callee == "unsafe.SliceData" && len(x.Args) == 1:
			return byteSliceBuf(x.Args[0], body, depth+1)
		case callee == "unsafe.StringData" && len(x.Args) == 1:
			return stringBuf{Producer: exprStr(x), Term: nulTerm(x.Args[0], body, depth+1)}
		default:
			if why, ok := nulHelpers[callee]; ok {
				return stringBuf{Producer: exprStr(x), Term: why}
			}
		}
	}

	return stringBuf{}
}

// byteSliceBuf classifies the []byte a char * was taken the address of.
func byteSliceBuf(e ast.Expr, body *ast.BlockStmt, depth int) stringBuf {
	if depth > 4 {
		return stringBuf{}
	}

	switch x := unparen(e).(type) {
	case *ast.Ident:
		var found stringBuf
		for _, v := range assignedValues(x.Name, body) {
			got := byteSliceBuf(v, body, depth+1)
			if got.Producer == "" {
				continue
			}
			if got.Term == "" {
				return got
			}
			found = got
		}

		return found

	case *ast.CallExpr:
		// []byte(s): a Go string copied into a byte buffer, terminator and all
		// or terminator and none.
		if at, ok := x.Fun.(*ast.ArrayType); ok && at.Len == nil && len(x.Args) == 1 {
			if id, ok := at.Elt.(*ast.Ident); ok && (id.Name == "byte" || id.Name == "uint8") {
				return stringBuf{Producer: exprStr(x), Term: nulTerm(x.Args[0], body, depth+1)}
			}
		}

	case *ast.CompositeLit:
		at, ok := x.Type.(*ast.ArrayType)
		if !ok || at.Len != nil || len(x.Elts) == 0 {
			return stringBuf{}
		}
		if id, ok := at.Elt.(*ast.Ident); ok && (id.Name == "byte" || id.Name == "uint8") {
			buf := stringBuf{Producer: exprStr(x)}
			if last, ok := x.Elts[len(x.Elts)-1].(*ast.BasicLit); ok && (last.Value == "0" || last.Value == `'\x00'` || last.Value == "0x00") {
				buf.Term = "[]byte composite ends in 0"
			}

			return buf
		}
	}

	return stringBuf{}
}

// nulTerm looks for the evidence that a Go *string* expression ends in a NUL.
//
// Only the last operand of a concatenation can terminate it, so `path + "\x00"`
// counts and `"\x00" + path` does not.
func nulTerm(e ast.Expr, body *ast.BlockStmt, depth int) string {
	if depth > 6 {
		return ""
	}

	switch x := unparen(e).(type) {
	case *ast.BasicLit:
		if x.Kind == token.STRING {
			if s, err := strconv.Unquote(x.Value); err == nil && strings.HasSuffix(s, "\x00") {
				return `string literal ending in "\x00"`
			}
		}

	case *ast.BinaryExpr:
		if x.Op == token.ADD {
			return nulTerm(x.Y, body, depth+1)
		}

	case *ast.Ident:
		for _, v := range assignedValues(x.Name, body) {
			if t := nulTerm(v, body, depth+1); t != "" {
				return t
			}
		}

	case *ast.CallExpr:
		if why, ok := nulHelpers[exprStr(x.Fun)]; ok {
			return why
		}
	}

	return ""
}

// assignedValues collects every value assigned to name in body.
//
// Flow-insensitive on purpose: `x += "\x00"` after the buffer is built and
// `if len(name) > 0 { n = ... }` are both the shape yzma writes, and the
// question asked of the result - is there positive evidence of a terminator -
// does not need an order.
func assignedValues(name string, body *ast.BlockStmt) []ast.Expr {
	var out []ast.Expr

	ast.Inspect(body, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.AssignStmt:
			if len(s.Lhs) != len(s.Rhs) {
				return true
			}
			for i, l := range s.Lhs {
				if id, ok := l.(*ast.Ident); ok && id.Name == name {
					out = append(out, s.Rhs[i])
				}
			}
		case *ast.ValueSpec:
			for i, id := range s.Names {
				if id.Name == name && i < len(s.Values) {
					out = append(out, s.Values[i])
				}
			}
		}

		return true
	})

	return out
}

func unparen(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}

// unwrapUnsafePointer strips unsafe.Pointer(x) conversions, which carry no
// information about x.
func unwrapUnsafePointer(e ast.Expr) ast.Expr {
	for {
		ce, ok := unparen(e).(*ast.CallExpr)
		if !ok || len(ce.Args) != 1 || exprStr(ce.Fun) != "unsafe.Pointer" {
			return e
		}
		e = ce.Args[0]
	}
}

// hasDeprecatedNote reports whether a doc comment carries the deprecation
// convention: a paragraph *beginning* with "Deprecated: ".
//
// The exact form is what matters rather than the sentiment, because the form is
// what has tooling behind it - gopls surfaces it and staticcheck's SA1019 flags
// every consumer call site. Prose saying a function is deprecated somewhere in
// the middle of a comment reads the same to a human and is invisible to both, so
// it deliberately does not count: mtmd.Encode's "Note: this function is marked as
// deprecated upstream" is exactly the case this must not accept.
func hasDeprecatedNote(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}

	for para := range strings.SplitSeq(doc.Text(), "\n\n") {
		if strings.HasPrefix(para, "Deprecated: ") {
			return true
		}
	}

	return false
}

func (a *analyzer) run() {
	a.ffiAlias = map[string]FfiType{}
	// pass 1: package-level vars that alias ffi types
	for _, f := range a.pkg.Syntax {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			for _, s := range gd.Specs {
				vs := s.(*ast.ValueSpec)
				for i, n := range vs.Names {
					if i >= len(vs.Values) {
						continue
					}
					v := vs.Values[i]
					switch x := v.(type) {
					case *ast.SelectorExpr:
						if id, ok := x.X.(*ast.Ident); ok && id.Name == "ffi" {
							if t, ok := ffiScalars[x.Sel.Name]; ok {
								a.ffiAlias[n.Name] = t
							}
						}
					case *ast.CallExpr:
						if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
							if id, ok := sel.X.(*ast.Ident); ok && id.Name == "ffi" && sel.Sel.Name == "NewType" {
								a.ffiAlias[n.Name] = a.buildStruct(n.Name, x.Args)
							}
						}
					}
				}
			}
		}
	}
	// second resolution pass so aliases referring to other aliases settle
	for i := 0; i < 3; i++ {
		for _, f := range a.pkg.Syntax {
			for _, d := range f.Decls {
				gd, ok := d.(*ast.GenDecl)
				if !ok || gd.Tok != token.VAR {
					continue
				}
				for _, s := range gd.Specs {
					vs := s.(*ast.ValueSpec)
					for i, n := range vs.Names {
						if i >= len(vs.Values) {
							continue
						}
						if ce, ok := vs.Values[i].(*ast.CallExpr); ok {
							if sel, ok := ce.Fun.(*ast.SelectorExpr); ok {
								if id, ok := sel.X.(*ast.Ident); ok && id.Name == "ffi" && sel.Sel.Name == "NewType" {
									a.ffiAlias[n.Name] = a.buildStruct(n.Name, ce.Args)
								}
							}
						}
					}
				}
			}
		}
	}

	a.bindings = map[string]*Binding{}

	// pass 2: Prep calls
	for _, f := range a.pkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			var lhs []ast.Expr
			var rhs ast.Expr
			switch x := n.(type) {
			case *ast.AssignStmt:
				if len(x.Rhs) == 1 {
					lhs, rhs = x.Lhs, x.Rhs[0]
				}
			default:
				return true
			}
			ce, ok := rhs.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := ce.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			variadic := false
			nfixed := -1
			var typeArgs []ast.Expr
			var nameArg ast.Expr
			switch sel.Sel.Name {
			case "Prep", "MustPrep":
				if len(ce.Args) < 2 {
					return true
				}
				nameArg = ce.Args[0]
				typeArgs = ce.Args[1:]
			case "PrepVar", "MustPrepVar":
				if len(ce.Args) < 3 {
					return true
				}
				variadic = true
				nameArg = ce.Args[0]
				if bl, ok := ce.Args[1].(*ast.BasicLit); ok {
					nfixed, _ = strconv.Atoi(bl.Value)
				}
				typeArgs = ce.Args[2:]
			default:
				return true
			}
			// receiver must be ffi.Lib
			rt := a.pkg.TypesInfo.TypeOf(sel.X)
			if rt == nil || !strings.HasSuffix(rt.String(), "ffi.Lib") {
				return true
			}
			bl, ok := nameArg.(*ast.BasicLit)
			if !ok {
				return true
			}
			cname, _ := strconv.Unquote(bl.Value)
			// LHS var name (first, skipping err)
			govar := ""
			for _, l := range lhs {
				if id, ok := l.(*ast.Ident); ok && id.Name != "err" && id.Name != "_" {
					govar = id.Name
					break
				}
			}
			if govar == "" {
				govar = cname
			}
			b := &Binding{
				GoVar:   govar,
				CName:   cname,
				PrepPos: a.fset.Position(ce.Lparen),
				Ret:     a.resolveTypeExpr(typeArgs[0]),
				NFixed:  nfixed,
				Variadi: variadic,
				Pkg:     a.pkg.PkgPath,
			}
			for _, ta := range typeArgs[1:] {
				b.Args = append(b.Args, a.resolveTypeExpr(ta))
			}
			if _, dup := a.bindings[govar]; !dup {
				a.order = append(a.order, govar)
			}
			a.bindings[govar] = b
			return true
		})
	}

	// pass 3: Call sites
	for _, f := range a.pkg.Syntax {
		var curFn string
		// The enclosing body is where a C string buffer is built, one or more
		// statements before the Call that hands it over.
		var curBody *ast.BlockStmt
		curExported, curDeprecated := false, false
		ast.Inspect(f, func(n ast.Node) bool {
			if fd, ok := n.(*ast.FuncDecl); ok {
				curFn, curBody = fd.Name.Name, fd.Body
				curExported, curDeprecated = fd.Name.IsExported(), hasDeprecatedNote(fd.Doc)
			}
			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := ce.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Call" {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}
			b, ok := a.bindings[id.Name]
			if !ok {
				return true
			}
			cs := CallSite{
				Pos: a.fset.Position(ce.Lparen), Fn: curFn,
				FnExported: curExported, FnDeprecated: curDeprecated,
			}
			if len(ce.Args) > 0 {
				cs.RetExpr = exprStr(ce.Args[0])
				pt, note := a.pointeeOf(ce.Args[0])
				cs.RetType = pt
				if note == "nil" {
					cs.RetNil = true
				}
			}
			for _, ae := range ce.Args[1:] {
				pt, note := a.pointeeOf(ae)
				k, sz := goKindOf(pt)
				if pt == nil {
					sz = -1
					k = KindUnknown
				}
				cs.Args = append(cs.Args, CallArg{
					Expr: exprStr(ae), Pointee: pt, Size: sz, Kind: k, Note: note,
					Str: stringBufOf(ae, curBody),
				})
			}
			b.Calls = append(b.Calls, cs)
			return true
		})
	}

	a.collectCallbacks()
}

// collectCallbacks finds the sites where C calls back into Go (RULE 5).
//
// Neither form goes through lib.Prep, so neither is one of the bindings above:
// ffi.PrepCif builds a descriptor libffi will use to unpack the C stack for a
// closure, and purego.NewCallback uses the Go func literal's own signature as
// the descriptor.
func (a *analyzer) collectCallbacks() {
	a.closureCode = map[string]string{}

	for _, f := range a.pkg.Syntax {
		var curFn string

		ast.Inspect(f, func(n ast.Node) bool {
			if fd, ok := n.(*ast.FuncDecl); ok {
				curFn = fd.Name.Name
			}

			ce, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := ce.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			id, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}

			switch {
			case id.Name == "ffi" && sel.Sel.Name == "PrepCif":
				// PrepCif(cif, abi, nfixed, ret, args...)
				if len(ce.Args) < 4 {
					return true
				}

				cb := &Callback{
					Form:   "ffi.PrepCif",
					GoID:   exprStr(ce.Args[0]),
					Fn:     curFn,
					Pkg:    a.pkg.PkgPath,
					Pos:    a.fset.Position(ce.Lparen),
					NFixed: -1,
					Ret:    a.resolveTypeExpr(ce.Args[3]),
				}

				if bl, ok := ce.Args[2].(*ast.BasicLit); ok {
					cb.NFixed, _ = strconv.Atoi(bl.Value)
				}

				for _, ta := range ce.Args[4:] {
					cb.Args = append(cb.Args, a.resolveTypeExpr(ta))
				}

				a.callbacks = append(a.callbacks, cb)

			case id.Name == "purego" && sel.Sel.Name == "NewCallback":
				if len(ce.Args) != 1 {
					return true
				}

				cb := &Callback{
					Form: "purego.NewCallback",
					GoID: curFn,
					Fn:   curFn,
					Pkg:  a.pkg.PkgPath,
					Pos:  a.fset.Position(ce.Lparen),
				}

				if sig, ok := a.pkg.TypesInfo.TypeOf(ce.Args[0]).(*types.Signature); ok {
					cb.Sig = sig
				}

				a.callbacks = append(a.callbacks, cb)

			case id.Name == "ffi" && sel.Sel.Name == "PrepClosureLoc":
				// PrepClosureLoc(closure, cif, fn, userData, code): the code
				// address is what gets stored in a function-pointer struct
				// member, and this is the only place it is tied to a cif.
				if len(ce.Args) == 5 {
					a.closureCode[exprStr(ce.Args[4])] = exprStr(ce.Args[1])
				}
			}

			return true
		})
	}

	a.collectFnPtrStores()
}

// collectFnPtrStores finds the assignments that install a callback's code
// pointer into a struct field - `p.ProgressCallback = uintptr(progressCallbackCode)`.
//
// That field is a function pointer C will *jump through*, but its Go type is a
// uintptr, so nothing in it says which callback belongs there: storing the log
// callback's code pointer into cb_eval compiles, lays out identically and is
// found by no other rule.
func (a *analyzer) collectFnPtrStores() {
	for _, f := range a.pkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			as, ok := n.(*ast.AssignStmt)
			if !ok || len(as.Lhs) != len(as.Rhs) {
				return true
			}

			for i, l := range as.Lhs {
				sel, ok := unparen(l).(*ast.SelectorExpr)
				if !ok {
					continue
				}

				st := a.pkg.TypesInfo.TypeOf(sel.X)
				if st == nil {
					continue
				}

				if p, ok := st.Underlying().(*types.Pointer); ok {
					st = p.Elem()
				}

				if _, ok := st.Underlying().(*types.Struct); !ok {
					continue
				}

				src := codeSource(as.Rhs[i])
				if src == "" {
					continue
				}

				a.stores = append(a.stores, &FnPtrStore{
					GoStruct: st,
					Field:    sel.Sel.Name,
					Expr:     exprStr(as.Rhs[i]),
					Pkg:      a.pkg.PkgPath,
					Site:     cmp.Or(a.closureCode[src], src),
					Pos:      a.fset.Position(as.TokPos),
				})
			}

			return true
		})
	}
}

// codeSource reduces the right-hand side of such an assignment to the
// identifier the pointer came from, looking through the uintptr and
// unsafe.Pointer conversions that carry no information of their own:
// `uintptr(progressCallbackCode)` is the code variable, and
// `uintptr(newAbortCallback(fn))` is the constructor that built the closure.
//
// It returns "" for anything else, including `uintptr(0)` - a cleared field is
// not a wrong code pointer, so there is nothing there to decide.
func codeSource(e ast.Expr) string {
	for range 4 {
		e = unparen(unwrapUnsafePointer(e))

		ce, ok := e.(*ast.CallExpr)
		if !ok {
			break
		}

		if id, ok := ce.Fun.(*ast.Ident); ok && id.Name == "uintptr" && len(ce.Args) == 1 {
			e = ce.Args[0]
			continue
		}

		return exprStr(ce.Fun)
	}

	if id, ok := e.(*ast.Ident); ok {
		return id.Name
	}

	return ""
}
