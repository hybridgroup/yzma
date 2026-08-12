package main

import (
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
	Pos     token.Position
	Fn      string // enclosing Go func
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
}

var arm64Sizes = types.SizesFor("gc", "arm64")

func goKindOf(t types.Type) (Kind, int) {
	if t == nil {
		return KindUnknown, -1
	}
	sz := int(arm64Sizes.Sizeof(t))
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
	}
	return fmt.Sprintf("%T", e)
}

// pointeeOf figures out the Go type whose bytes libffi will read for an
// argument expression passed to Fun.Call.
func (a *analyzer) pointeeOf(e ast.Expr) (types.Type, string) {
	// unwrap unsafe.Pointer(x) conversions
	for {
		ce, ok := e.(*ast.CallExpr)
		if !ok || len(ce.Args) != 1 {
			break
		}
		if sel, ok := ce.Fun.(*ast.SelectorExpr); ok {
			if id, ok := sel.X.(*ast.Ident); ok && id.Name == "unsafe" && sel.Sel.Name == "Pointer" {
				e = ce.Args[0]
				continue
			}
		}
		break
	}
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
		ast.Inspect(f, func(n ast.Node) bool {
			if fd, ok := n.(*ast.FuncDecl); ok {
				curFn = fd.Name.Name
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
			cs := CallSite{Pos: a.fset.Position(ce.Lparen), Fn: curFn}
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
				cs.Args = append(cs.Args, CallArg{Expr: exprStr(ae), Pointee: pt, Size: sz, Kind: k, Note: note})
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
			}

			return true
		})
	}
}
