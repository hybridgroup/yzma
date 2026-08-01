package main

import (
	"fmt"
	"os"
	"strings"
)

// CParam is one parameter of a C function declaration.
type CParam struct {
	Raw  string // original text
	Norm string // normalised type text (name stripped)
	Kind Kind
	Size int
}

// CFunc is one LLAMA_API / GGML_API / MTMD_API declaration.
type CFunc struct {
	Name    string
	File    string
	Line    int
	RetRaw  string
	RetKind Kind
	RetSize int
	Params  []CParam
	Vararg  bool
	Raw     string
}

// Kind classifies a C type for ABI purposes.
type Kind int

const (
	KindUnknown Kind = iota
	KindVoid
	KindSint
	KindUint
	KindFloat  // 4 byte float
	KindDouble // 8 byte float
	KindPointer
	KindStruct
)

func (k Kind) String() string {
	switch k {
	case KindVoid:
		return "void"
	case KindSint:
		return "sint"
	case KindUint:
		return "uint"
	case KindFloat:
		return "float"
	case KindDouble:
		return "double"
	case KindPointer:
		return "ptr"
	case KindStruct:
		return "struct"
	}
	return "?"
}

// stripComments blanks out comments while preserving byte offsets/newlines.
func stripComments(src string) string {
	b := []byte(src)
	out := make([]byte, len(b))
	copy(out, b)
	i := 0
	for i < len(b) {
		switch {
		case b[i] == '/' && i+1 < len(b) && b[i+1] == '/':
			for i < len(b) && b[i] != '\n' {
				out[i] = ' '
				i++
			}
		case b[i] == '/' && i+1 < len(b) && b[i+1] == '*':
			for i < len(b) && !(b[i] == '*' && i+1 < len(b) && b[i+1] == '/') {
				if b[i] != '\n' {
					out[i] = ' '
				}
				i++
			}
			if i < len(b) {
				out[i] = ' '
				i++
			}
			if i < len(b) {
				out[i] = ' '
				i++
			}
		case b[i] == '"':
			i++
			for i < len(b) && b[i] != '"' {
				if b[i] == '\\' {
					i++
				}
				i++
			}
			i++
		default:
			i++
		}
	}
	return string(out)
}

// unwrapDeprecated rewrites DEPRECATED(<decl>, "msg") into <decl> in place,
// blanking the wrapper so byte offsets and line numbers are preserved.
func unwrapDeprecated(src string) string {
	b := []byte(src)
	const tok = "DEPRECATED("
	for _, idx := range findAll(string(b), tok) {
		if idx > 0 && isIdentChar(b[idx-1]) {
			continue
		}
		depth := 0
		close := -1
		for i := idx + len(tok) - 1; i < len(b); i++ {
			if b[i] == '(' {
				depth++
			} else if b[i] == ')' {
				depth--
				if depth == 0 {
					close = i
					break
				}
			}
		}
		if close < 0 {
			continue
		}
		// last top-level comma inside the wrapper
		depth = 0
		comma := -1
		for i := idx + len(tok); i < close; i++ {
			switch b[i] {
			case '(':
				depth++
			case ')':
				depth--
			case ',':
				if depth == 0 {
					comma = i
				}
			}
		}
		for i := idx; i < idx+len(tok); i++ {
			b[i] = ' '
		}
		from := comma
		if from < 0 {
			from = close
		}
		for i := from; i <= close; i++ {
			if b[i] != '\n' {
				b[i] = ' '
			}
		}
	}
	return string(b)
}

// typedefs maps a typedef name -> the underlying type text.
var typedefs = map[string]string{}

// enums holds the set of names that are enum tags.
var enumNames = map[string]bool{}

func collectTypedefs(src string) {
	// typedef int32_t llama_pos;   /  typedef struct x * y;   / typedef enum {...} name;
	s := src
	for i := 0; i < len(s); {
		j := strings.Index(s[i:], "typedef")
		if j < 0 {
			break
		}
		j += i
		// find the terminating ';' at brace depth 0
		depth := 0
		k := j
		for ; k < len(s); k++ {
			switch s[k] {
			case '{':
				depth++
			case '}':
				depth--
			case ';':
				if depth == 0 {
					goto done
				}
			}
		}
	done:
		if k >= len(s) {
			break
		}
		decl := s[j:k]
		i = k + 1
		// remove any {...} body
		if a := strings.Index(decl, "{"); a >= 0 {
			b := strings.LastIndex(decl, "}")
			if b > a {
				decl = decl[:a] + " " + decl[b+1:]
			}
		}
		decl = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(decl), "typedef"))
		decl = strings.Join(strings.Fields(decl), " ")
		if decl == "" || strings.Contains(decl, "(") {
			// function-pointer typedef -> pointer
			// name is inside (*name)
			a := strings.Index(decl, "(*")
			if a >= 0 {
				rest := decl[a+2:]
				b := strings.Index(rest, ")")
				if b > 0 {
					typedefs[strings.TrimSpace(rest[:b])] = "void *"
				}
			}
			continue
		}
		fields := strings.Fields(decl)
		if len(fields) < 2 {
			continue
		}
		name := fields[len(fields)-1]
		under := strings.Join(fields[:len(fields)-1], " ")
		// handle "typedef struct foo * bar" -> name=bar under="struct foo *"
		name = strings.TrimPrefix(name, "*")
		if strings.HasPrefix(fields[len(fields)-1], "*") {
			under += " *"
		}
		if name != "" {
			typedefs[name] = strings.TrimSpace(under)
		}
	}
}

// baseTypes maps a fully-resolved scalar C type spelling to kind+size (arm64/LP64).
var baseTypes = map[string]struct {
	k Kind
	s int
}{
	"void":               {KindVoid, 0},
	"bool":               {KindUint, 1},
	"_Bool":              {KindUint, 1},
	"char":               {KindSint, 1},
	"signed char":        {KindSint, 1},
	"unsigned char":      {KindUint, 1},
	"short":              {KindSint, 2},
	"unsigned short":     {KindUint, 2},
	"int":                {KindSint, 4},
	"unsigned":           {KindUint, 4},
	"unsigned int":       {KindUint, 4},
	"long":               {KindSint, 8},
	"unsigned long":      {KindUint, 8},
	"long long":          {KindSint, 8},
	"unsigned long long": {KindUint, 8},
	"float":              {KindFloat, 4},
	"double":             {KindDouble, 8},
	"int8_t":             {KindSint, 1},
	"uint8_t":            {KindUint, 1},
	"int16_t":            {KindSint, 2},
	"uint16_t":           {KindUint, 2},
	"int32_t":            {KindSint, 4},
	"uint32_t":           {KindUint, 4},
	"int64_t":            {KindSint, 8},
	"uint64_t":           {KindUint, 8},
	"size_t":             {KindUint, 8},
	"ssize_t":            {KindSint, 8},
	"ptrdiff_t":          {KindSint, 8},
	"intptr_t":           {KindSint, 8},
	"uintptr_t":          {KindUint, 8},
	"time_t":             {KindSint, 8},
	"FILE":               {KindStruct, -1},
}

// classify resolves a C type text (no declarator name) to kind+size.
func classify(t string) (Kind, int) {
	t = strings.TrimSpace(t)
	t = strings.ReplaceAll(t, "const", " ")
	t = strings.ReplaceAll(t, "volatile", " ")
	t = strings.ReplaceAll(t, "restrict", " ")
	t = strings.Join(strings.Fields(t), " ")
	if strings.Contains(t, "*") || strings.Contains(t, "[") {
		return KindPointer, 8
	}
	if t == "" {
		return KindUnknown, -1
	}
	// enum X -> 4 bytes (C enums in llama.h/ggml.h all fit in int)
	if strings.HasPrefix(t, "enum ") {
		return KindSint, 4
	}
	if strings.HasPrefix(t, "struct ") || strings.HasPrefix(t, "union ") {
		return KindStruct, -1
	}
	if bt, ok := baseTypes[t]; ok {
		return bt.k, bt.s
	}
	// resolve typedefs (bounded depth)
	seen := map[string]bool{}
	cur := t
	for i := 0; i < 12; i++ {
		u, ok := typedefs[cur]
		if !ok || seen[cur] {
			break
		}
		seen[cur] = true
		u = strings.Join(strings.Fields(strings.ReplaceAll(u, "const", " ")), " ")
		if strings.Contains(u, "*") {
			return KindPointer, 8
		}
		if strings.HasPrefix(u, "enum") {
			return KindSint, 4
		}
		if strings.HasPrefix(u, "struct") || strings.HasPrefix(u, "union") {
			return KindStruct, -1
		}
		if bt, ok := baseTypes[u]; ok {
			return bt.k, bt.s
		}
		cur = u
	}
	if enumNames[t] {
		return KindSint, 4
	}
	return KindUnknown, -1
}

// splitTop splits s on top-level commas (respecting parens).
func splitTop(s string) []string {
	var out []string
	depth := 0
	last := 0
	for i, c := range s {
		switch c {
		case '(', '[':
			depth++
		case ')', ']':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, s[last:i])
				last = i + 1
			}
		}
	}
	out = append(out, s[last:])
	return out
}

// stripDeclName removes the parameter name from a declaration like "const char * path".
func stripDeclName(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return p
	}
	// array declarator: "int32_t buf[]" -> pointer
	if i := strings.Index(p, "["); i >= 0 {
		p = p[:i] + "*"
	}
	if strings.Contains(p, "(*") {
		return "void *"
	}
	f := strings.Fields(p)
	if len(f) == 0 {
		return p
	}
	last := f[len(f)-1]
	if last == "*" || last == "**" || strings.HasSuffix(last, "*") && strings.Trim(last, "*") == "" {
		return p // no name, e.g. "const char *"
	}
	// last token is either the name, or "*name"
	if strings.HasPrefix(last, "*") {
		stars := len(last) - len(strings.TrimLeft(last, "*"))
		return strings.Join(f[:len(f)-1], " ") + " " + strings.Repeat("*", stars)
	}
	// if it is a known single-word type with no other tokens (e.g. "void", "int"), keep it
	if len(f) == 1 {
		return p
	}
	// "unsigned int" style with no name? handled by baseTypes lookups after dropping last
	cand := strings.Join(f[:len(f)-1], " ")
	if _, ok := baseTypes[strings.Join(f, " ")]; ok {
		return strings.Join(f, " ")
	}
	if strings.HasSuffix(cand, "*") {
		return cand
	}
	return cand
}

func parseHeader(path, apiMacro string) ([]CFunc, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	src := unwrapDeprecated(stripComments(string(raw)))
	collectTypedefs(src)
	// collect enum tag names: "enum llama_foo {"
	for _, m := range findAll(src, "enum ") {
		rest := src[m+5:]
		f := strings.Fields(rest)
		if len(f) > 0 {
			enumNames[strings.TrimSuffix(f[0], "{")] = true
		}
	}

	var out []CFunc
	for _, idx := range findAll(src, apiMacro) {
		// must be a token boundary
		if idx > 0 && isIdentChar(src[idx-1]) {
			continue
		}
		after := idx + len(apiMacro)
		if after < len(src) && isIdentChar(src[after]) {
			continue
		}
		// find ';' at paren depth 0
		depth := 0
		end := -1
		for k := after; k < len(src); k++ {
			switch src[k] {
			case '(':
				depth++
			case ')':
				depth--
			case ';':
				if depth == 0 {
					end = k
				}
			case '{':
				if depth == 0 {
					end = -2
				}
			}
			if end != -1 {
				break
			}
		}
		if end < 0 {
			continue
		}
		decl := src[after:end]
		line := 1 + strings.Count(src[:idx], "\n")
		fn, ok := parseDecl(decl)
		if !ok {
			out = append(out, CFunc{Name: "", File: path, Line: line, Raw: squash(decl)})
			continue
		}
		fn.File = path
		fn.Line = line
		fn.Raw = squash(decl)
		out = append(out, fn)
	}
	return out, nil
}

func squash(s string) string { return strings.Join(strings.Fields(s), " ") }

func isIdentChar(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func findAll(s, sub string) []int {
	var out []int
	for i := 0; ; {
		j := strings.Index(s[i:], sub)
		if j < 0 {
			return out
		}
		out = append(out, i+j)
		i += j + len(sub)
	}
}

func parseDecl(decl string) (CFunc, bool) {
	var fn CFunc
	// find first top-level '('
	open := -1
	for i := 0; i < len(decl); i++ {
		if decl[i] == '(' {
			open = i
			break
		}
	}
	if open < 0 {
		return fn, false
	}
	// matching close
	depth := 0
	close := -1
	for i := open; i < len(decl); i++ {
		if decl[i] == '(' {
			depth++
		} else if decl[i] == ')' {
			depth--
			if depth == 0 {
				close = i
				break
			}
		}
	}
	if close < 0 {
		return fn, false
	}
	head := strings.TrimSpace(decl[:open])
	f := strings.Fields(head)
	if len(f) < 2 {
		return fn, false
	}
	name := f[len(f)-1]
	if strings.HasPrefix(name, "*") {
		// "const char *llama_x" style
		fn.Name = strings.TrimLeft(name, "*")
		fn.RetRaw = strings.Join(f[:len(f)-1], " ") + " " + strings.Repeat("*", len(name)-len(strings.TrimLeft(name, "*")))
	} else {
		fn.Name = name
		fn.RetRaw = strings.Join(f[:len(f)-1], " ")
	}
	if fn.Name == "" || !isIdentChar(fn.Name[0]) {
		return fn, false
	}
	fn.RetKind, fn.RetSize = classify(fn.RetRaw)
	body := strings.TrimSpace(decl[open+1 : close])
	if body != "" && body != "void" {
		for _, p := range splitTop(body) {
			ps := strings.TrimSpace(p)
			if ps == "..." {
				fn.Vararg = true
				continue
			}
			norm := stripDeclName(ps)
			k, sz := classify(norm)
			fn.Params = append(fn.Params, CParam{Raw: squash(ps), Norm: squash(norm), Kind: k, Size: sz})
		}
	}
	return fn, true
}

func (f CFunc) Sig() string {
	var ps []string
	for _, p := range f.Params {
		ps = append(ps, fmt.Sprintf("%s{%s/%d}", p.Norm, p.Kind, p.Size))
	}
	return fmt.Sprintf("%s{%s/%d} %s(%s)", f.RetRaw, f.RetKind, f.RetSize, f.Name, strings.Join(ps, ", "))
}
