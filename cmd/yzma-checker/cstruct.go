package main

import (
	"strconv"
	"strings"
)

// CField is one member of a C struct.
type CField struct {
	Name string
	Norm string
	Kind Kind
	Size int
	Cnt  int // array count (1 for scalars)
}

// CStruct is a parsed C struct definition.
type CStruct struct {
	Name   string
	Fields []CField
	Size   int // -1 if not computable
	Align  int
}

var cstructs = map[string]*CStruct{}

// resetCTypes clears the C parser's accumulated state so that more than one
// audit can run in a single process (the self-test does exactly that).
func resetCTypes() {
	clear(cstructs)
	clear(typedefs)
	clear(enumNames)
}

// collectStructs parses "struct NAME { ... };" definitions out of a header.
func collectStructs(src string) {
	for _, idx := range findAll(src, "struct ") {
		if idx > 0 && isIdentChar(src[idx-1]) {
			continue
		}
		rest := src[idx+len("struct "):]
		f := strings.Fields(rest)
		if len(f) < 2 {
			continue
		}
		name := f[0]
		if !strings.HasPrefix(strings.TrimSpace(rest[len(name):]), "{") {
			continue
		}
		open := strings.Index(rest, "{")
		depth := 0
		end := -1
		for i := open; i < len(rest); i++ {
			if rest[i] == '{' {
				depth++
			} else if rest[i] == '}' {
				depth--
				if depth == 0 {
					end = i
					break
				}
			}
		}
		if end < 0 {
			continue
		}
		body := rest[open+1 : end]
		cs := &CStruct{Name: name}
		for _, stmt := range strings.Split(body, ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if strings.Contains(stmt, "(") && !strings.Contains(stmt, "(*") {
				continue // method-ish / macro
			}
			// possible comma-separated declarators: "int32_t a, b;"
			parts := splitTop(stmt)
			base := ""
			for pi, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				decl := p
				if pi > 0 && base != "" {
					decl = base + " " + p
				}
				fname, norm, cnt := splitField(decl)
				if pi == 0 {
					base = fieldBaseType(decl)
				}
				k, sz := classify(norm)
				cs.Fields = append(cs.Fields, CField{Name: fname, Norm: squash(norm), Kind: k, Size: sz, Cnt: cnt})
			}
		}
		computeStructSize(cs)
		if _, exists := cstructs[name]; !exists {
			cstructs[name] = cs
		}
	}
}

func fieldBaseType(decl string) string {
	f := strings.Fields(decl)
	if len(f) < 2 {
		return ""
	}
	return strings.Join(f[:len(f)-1], " ")
}

// splitField returns (name, normalisedType, arrayCount).
func splitField(decl string) (string, string, int) {
	decl = strings.TrimSpace(decl)
	cnt := 1
	if i := strings.Index(decl, "["); i >= 0 {
		j := strings.Index(decl, "]")
		if j > i {
			n, err := strconv.Atoi(strings.TrimSpace(decl[i+1 : j]))
			if err == nil {
				cnt = n
			} else {
				cnt = -1
			}
		}
		decl = decl[:i]
	}
	f := strings.Fields(decl)
	if len(f) == 0 {
		return "", "", cnt
	}
	last := f[len(f)-1]
	name := strings.TrimLeft(last, "*")
	stars := len(last) - len(name)
	norm := strings.Join(f[:len(f)-1], " ")
	if stars > 0 {
		norm += " " + strings.Repeat("*", stars)
	}
	if norm == "" {
		norm = decl
		name = ""
	}
	return name, norm, cnt
}

func computeStructSize(cs *CStruct) {
	off := 0
	maxAlign := 1
	for _, f := range cs.Fields {
		if f.Size <= 0 || f.Cnt <= 0 {
			if f.Kind == KindStruct {
				if sub, ok := cstructs[strings.TrimPrefix(f.Norm, "struct ")]; ok && sub.Size > 0 {
					al := sub.Align
					if off%al != 0 {
						off += al - off%al
					}
					off += sub.Size * f.Cnt
					if al > maxAlign {
						maxAlign = al
					}
					continue
				}
			}
			cs.Size = -1
			cs.Align = -1
			return
		}
		al := f.Size
		if al > 8 {
			al = 8
		}
		if off%al != 0 {
			off += al - off%al
		}
		off += f.Size * f.Cnt
		if al > maxAlign {
			maxAlign = al
		}
	}
	if off%maxAlign != 0 {
		off += maxAlign - off%maxAlign
	}
	cs.Size = off
	cs.Align = maxAlign
}

// cTypeSize resolves a possibly-struct C type text to a byte size (-1 unknown).
func cTypeSize(norm string) int {
	k, sz := classify(norm)
	if k != KindStruct {
		return sz
	}
	n := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.ReplaceAll(norm, "const", "")), "struct"))
	if cs, ok := cstructs[strings.TrimSpace(n)]; ok {
		return cs.Size
	}
	// maybe a typedef to a struct
	if u, ok := typedefs[strings.TrimSpace(n)]; ok {
		un := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(u), "struct"))
		if cs, ok := cstructs[un]; ok {
			return cs.Size
		}
	}
	return -1
}
