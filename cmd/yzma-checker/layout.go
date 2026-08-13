package main

import (
	"fmt"
	"go/types"
	"slices"
	"strings"
	"unicode"
)

// Struct-by-value slots used to be compared by total size alone, which is the
// weakest check the tool makes and the one that failed in the field
// (hybridgroup/yzma#289): llama.cpp b10375 inserted a uint32_t into
// llama_context_params, the Go and cif sides did not gain it, and the four
// missing bytes were reabsorbed by the alignment padding in front of the first
// pointer member. Both sides really were 160 bytes, while every field between
// the insertion point and that padding was displaced by four bytes -- silently,
// so embeddings simply stopped working.
//
// Any struct mixing 4-byte and 8-byte members has interior padding, so this is
// not a contrived shape: it is the shape of llama_context_params,
// llama_model_params and llama_batch, the three structs upstream churns most.
//
// The fix is to compare field layouts rather than total sizes. All three
// representations of a struct-by-value slot get flattened to the same
// primitive form -- a Leaf list -- and compared element by element:
//
//	C struct        cstructs[name]  -> cLeaves
//	cif descriptor  ffi.NewType(..) -> ffiLeaves
//	Go struct       go/types        -> goLeaves
//
// Flattening to leaves rather than to immediate fields is what makes the three
// comparable at all: a nested struct is one field on the C side and one nested
// ffi.NewType on the cif side, and arrays are one field on both but N slots to
// libffi. Padding is never emitted, so an equal leaf list means equal offsets
// for every byte either side will actually read.

// Leaf is one primitive (non-struct, non-array) member of a flattened struct,
// at its byte offset from the start of the outermost struct.
type Leaf struct {
	Off  int
	Size int
	Kind Kind
	Path string // dotted field path, for a diff a human can act on

	// GoType is the member's Go type, set on the Go side only. Offset, width and
	// ABI class are all a layout diff needs, and they are also exactly what a Go
	// `func` field and a `uintptr` field have in common - which is why a member C
	// will *call* through needs the type itself. See checkFnPtrMembers.
	GoType types.Type
}

func (l Leaf) String() string { return fmt.Sprintf("+%d %s/%dB %s", l.Off, l.Kind, l.Size, l.Path) }

// maxLeaves bounds flattening: a struct holding a huge array would otherwise
// expand to a leaf per element. Exceeding it fails resolution, which surfaces
// as a NOT VERIFIED skip rather than as a silent pass.
const maxLeaves = 4096

// maxLeafDepth bounds nesting, so a malformed self-referential struct cannot
// recurse forever.
const maxLeafDepth = 8

func elemPath(name string, cnt, i int) string {
	if cnt == 1 {
		return name
	}

	return fmt.Sprintf("%s[%d]", name, i)
}

// cLeaves flattens a parsed C struct.
func cLeaves(cs *CStruct) ([]Leaf, bool) {
	var out []Leaf
	if !appendCLeaves(&out, cs, 0, "", 0) {
		return nil, false
	}

	return out, true
}

func appendCLeaves(out *[]Leaf, cs *CStruct, base int, prefix string, depth int) bool {
	if cs == nil || cs.Size < 0 || len(cs.Offs) != len(cs.Fields) || depth > maxLeafDepth {
		return false
	}

	for i, f := range cs.Fields {
		off, name := base+cs.Offs[i], prefix+f.Name
		if f.Cnt <= 0 {
			return false
		}

		if f.Kind == KindStruct {
			sub := cstructOf(f.Norm)
			if sub == nil || sub.Size <= 0 {
				return false
			}

			for e := range f.Cnt {
				if !appendCLeaves(out, sub, off+e*sub.Size, elemPath(name, f.Cnt, e)+".", depth+1) {
					return false
				}
			}

			continue
		}

		if f.Size <= 0 {
			return false
		}

		for e := range f.Cnt {
			if len(*out) >= maxLeaves {
				return false
			}

			*out = append(*out, Leaf{Off: off + e*f.Size, Size: f.Size, Kind: f.Kind, Path: elemPath(name, f.Cnt, e)})
		}
	}

	return true
}

// ffiLeaves flattens an ffi.NewType(...) descriptor. Element paths are
// positional because a libffi descriptor carries no member names.
func ffiLeaves(t FfiType) ([]Leaf, bool) {
	var out []Leaf
	if t.Kind != KindStruct {
		return nil, false
	}

	if !appendFfiLeaves(&out, t, 0, "", 0) {
		return nil, false
	}

	return out, true
}

func appendFfiLeaves(out *[]Leaf, t FfiType, base int, prefix string, depth int) bool {
	if t.Size < 0 || len(t.Offs) != len(t.Elems) || depth > maxLeafDepth {
		return false
	}

	for i, el := range t.Elems {
		off, name := base+t.Offs[i], fmt.Sprintf("%se%d", prefix, i)
		if el.Kind == KindStruct {
			if !appendFfiLeaves(out, el, off, name+".", depth+1) {
				return false
			}

			continue
		}

		if el.Size <= 0 || el.Kind == KindUnknown {
			return false
		}

		if len(*out) >= maxLeaves {
			return false
		}

		*out = append(*out, Leaf{Off: off, Size: el.Size, Kind: el.Kind, Path: name})
	}

	return true
}

// Go structs flattened under every architecture yzma supports, counted for the
// same reason the pointer targets are: every width in this report is computed
// under goTargets[0], so a struct that lays out differently under the other
// target is being audited by numbers that are not its own - and that would be
// invisible in a report full of clean comparisons.
var archChecked, archClean int

// archSeen is the structs already compared. The same struct is a slot in
// several bindings, and one finding per slot would be the same defect printed
// five times.
var archSeen []types.Type

// diffArchLayouts flattens one Go struct under each architecture yzma supports
// and reports how the member lists disagree.
//
// The README used to *argue* that amd64 agrees with arm64 member for member for
// the types that cross this boundary, which left every number in the report
// resting on an unmeasured claim: the tool prints arm64 widths whatever GOARCH
// it runs on, so a divergence would mean one of the two supported architectures
// is audited by another's offsets, silently, with every comparison green.
// Running the same walker twice is the whole of what it takes to know instead.
func diffArchLayouts(t types.Type) string {
	if t == nil {
		return ""
	}

	if slices.ContainsFunc(archSeen, func(s types.Type) bool { return types.Identical(s, t) }) {
		return ""
	}

	archSeen = append(archSeen, t)

	home, ok := leavesUnder(t, goTargets[0].Sizes)
	if !ok {
		return "" // cmpStructLayout owns that skip; a second copy of it says nothing new
	}

	var probs []string
	for _, tgt := range goTargets[1:] {
		other, ok := leavesUnder(t, tgt.Sizes)
		if !ok {
			noteSkip("Go struct %s member layout not resolvable under %s - NOT VERIFIED across architectures", t.String(), tgt.Arch)
			continue
		}

		archChecked++

		d := diffLeaves(home, other, goTargets[0].Arch, tgt.Arch)
		if len(d) == 0 {
			archClean++
			continue
		}

		probs = append(probs, fmt.Sprintf("Go struct %s lays out differently under %s than under %s, so this report's widths are not that target's: %s",
			t.String(), tgt.Arch, goTargets[0].Arch, strings.Join(d, "; ")))
	}

	return strings.Join(probs, "; ")
}

// leavesUnder flattens a Go struct under another target's layout.
//
// The active target is swapped around the walk rather than threaded through
// goLeaves, appendGoLeaf and goKindOf as a parameter, so that the
// cross-architecture comparison uses the one struct walker every other check
// uses instead of a second one that could drift from it.
func leavesUnder(t types.Type, s types.Sizes) ([]Leaf, bool) {
	prev := goSizes
	goSizes = s

	defer func() { goSizes = prev }()

	return goLeaves(t)
}

// goLeaves flattens a Go struct type using the gc layout of the target in force.
func goLeaves(t types.Type) ([]Leaf, bool) {
	var out []Leaf
	if t == nil {
		return nil, false
	}

	st, ok := t.Underlying().(*types.Struct)
	if !ok {
		return nil, false
	}

	if !appendGoLeaves(&out, st, 0, "", 0) {
		return nil, false
	}

	return out, true
}

func appendGoLeaves(out *[]Leaf, st *types.Struct, base int, prefix string, depth int) bool {
	if depth > maxLeafDepth {
		return false
	}

	fields := make([]*types.Var, st.NumFields())
	for i := range fields {
		fields[i] = st.Field(i)
	}

	offs := goSizes.Offsetsof(fields)

	for i, f := range fields {
		if !appendGoLeaf(out, f.Type(), base+int(offs[i]), prefix+f.Name(), depth) {
			return false
		}
	}

	return true
}

func appendGoLeaf(out *[]Leaf, t types.Type, off int, name string, depth int) bool {
	switch u := t.Underlying().(type) {
	case *types.Struct:
		return appendGoLeaves(out, u, off, name+".", depth+1)
	case *types.Array:
		n := int(u.Len())
		if n < 0 {
			return false
		}

		stride := int(goSizes.Sizeof(u.Elem()))
		if stride <= 0 {
			return false
		}

		for e := range n {
			if !appendGoLeaf(out, u.Elem(), off+e*stride, elemPath(name, n, e), depth+1) {
				return false
			}
		}

		return true
	}

	k, sz := goKindOf(t)
	if sz <= 0 || k == KindUnknown {
		return false
	}

	if len(*out) >= maxLeaves {
		return false
	}

	*out = append(*out, Leaf{Off: off, Size: sz, Kind: k, Path: name, GoType: t})

	return true
}

// maxLeafDiffs bounds how much of a diff is reported. One displaced field
// shifts every field after it, so printing them all buries the insertion point
// that is the actual finding.
const maxLeafDiffs = 3

// diffLeaves reports how two flattened layouts disagree. aName/bName label the
// sides. An empty result means every member sits at the same offset with the
// same width and ABI class.
func diffLeaves(a, b []Leaf, aName, bName string) []string {
	var out []string

	for i := 0; i < len(a) && i < len(b); i++ {
		if len(out) >= maxLeafDiffs {
			break
		}

		x, y := a[i], b[i]
		if x.Off == y.Off && x.Size == y.Size && kindCompat(x.Kind, y.Kind) {
			continue
		}

		out = append(out, fmt.Sprintf("member %d: %s %s vs %s %s", i, aName, x, bName, y))
	}

	if len(a) != len(b) {
		out = append(out, fmt.Sprintf("%s has %d members, %s has %d (%s)",
			aName, len(a), bName, len(b), extraMembers(a, b, aName, bName)))
	}

	return out
}

// diffLeafPrefix compares a Go struct's layout against a cif descriptor's,
// where the Go struct is allowed to carry members the descriptor does not have
// — provided they sit past the end of the C struct.
//
// libffi reads (and for a struct return, writes) exactly ffi_type.size bytes
// through the buffer it is handed, so the C-declared size is the authority on
// how far it reaches. A Go struct appending its own bookkeeping after the FFI
// members, as llama.Batch does with capTokens/capSeq, is bytes libffi never
// touches: not a finding. A Go struct that is *short*, or whose extra members
// land inside the region libffi does reach, is.
func diffLeafPrefix(gl, fl []Leaf, goName, cifName string, cifSize int) []string {
	out := diffLeaves(gl[:min(len(gl), len(fl))], fl[:min(len(gl), len(fl))], goName, cifName)

	if len(gl) < len(fl) {
		out = append(out, fmt.Sprintf("%s has only %d members, %s reads %d (%s)",
			goName, len(gl), cifName, len(fl), extraMembers(gl, fl, goName, cifName)))

		return out
	}

	for _, l := range gl[len(fl):] {
		if l.Off < cifSize {
			out = append(out, fmt.Sprintf("%s member %s is not in %s but lies inside the %dB libffi reads",
				goName, l.String(), cifName, cifSize))
		}
	}

	return out
}

// memberAliases reconciles the places where yzma's Go field names diverge from
// the C member names by more than case and underscores. Each pair is rewritten
// to its second form on both sides, so the two conventions meet in the middle.
//
// It is deliberately a short, closed list rather than a fuzzy matcher: an alias
// loose enough to guess would also be loose enough to accept a genuine
// mismatch. Anything not listed here stays unmatched and is reported as NOT
// VERIFIED, so a new divergence shows up as a skip that needs one line added
// rather than as a silent hole.
var memberAliases = [][2]string{
	{"attention", "attn"},          // FlashAttentionType / flash_attn_type
	{"context", "ctx"},             // VideoContext / video_ctx
	{"prompteval", "peval"},        // TPromptEvalMs / t_p_eval_ms
	{"tensortypes", "ttoverrides"}, // TensorTypes / tt_overrides
}

// normName reduces a member name to a form comparable across the two naming
// conventions: Go's NThreadsBatch and C's n_threads_batch become the same
// string. The counting prefix is kept, so this is the exact form.
func normName(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '_' {
			continue
		}

		b.WriteRune(unicode.ToLower(r))
	}

	s = b.String()

	for _, a := range memberAliases {
		s = strings.ReplaceAll(s, a[0], a[1])
	}

	return s
}

// normNameLoose is normName with the C counting prefix dropped, for the second
// pass of the two-pass match in diffSwappedMembers.
//
// A leading n_ is a prefix yzma sometimes keeps (NThreads for n_threads) and
// sometimes drops (mtmd's Threads for n_threads), so it has to be optional
// somewhere. It cannot be optional in the exact form, though, because dropping
// it unconditionally is lossy in a way that matters: struct llama_batch has both
// n_seq_id and seq_id, so a single normalisation that treats the prefix as
// noise maps two real, adjacent members onto one name and can no longer tell
// which of them a Go member is. Matching exactly first and only falling back to
// this keeps NSeqId on n_seq_id and SeqId on seq_id.
//
// The prefix has to be recognised on the name as written: by the time normName
// has removed the underscores there is no prefix left to see, only a leading n,
// and trimming that would make `name` and `ame` the same member. The Go spelling
// of the same prefix is an N before another capital, which is what separates
// NThreads from Name.
func normNameLoose(s string) string {
	if !hasCountPrefix(s) {
		return normName(s)
	}

	return strings.TrimPrefix(normName(s), "n")
}

func hasCountPrefix(s string) bool {
	if rest, ok := strings.CutPrefix(s, "n_"); ok {
		return rest != ""
	}

	rest, ok := strings.CutPrefix(s, "N")

	return ok && rest != "" && unicode.IsUpper(rune(rest[0]))
}

// indexMember records the C member at index i under key, marking a key two
// members share as -1 so it can never answer "found at another index".
func indexMember(m map[string]int, key string, i int) {
	if _, dup := m[key]; dup {
		m[key] = -1

		return
	}

	m[key] = i
}

// diffSwappedMembers finds members the Go struct has transposed relative to the
// C struct.
//
// Offsets cannot see this: swapping two int32_t members leaves every offset,
// width and ABI class identical, so the call succeeds and the C library simply
// receives each value in the other's parameter. Names can see it, and a
// *permutation* is the evidence to key on — a Go name that matches a C member
// at a different index is a swap, whereas a name matching nothing at all is
// just a binding that renamed a field, so that is reported as unverified rather
// than as a defect.
func diffSwappedMembers(gl, cl []Leaf, goName, cName string) (probs, unmatched []string) {
	// Two indexes, because the counting prefix has to be optional without being
	// discarded: the exact one answers first, and the loose one only sees a Go
	// member whose name is not in it. A duplicate in either would make "found at
	// another index" meaningless, so those are marked -1 and dropped from
	// consideration. Dropping them silently would be indistinguishable from
	// checking them, though, so a Go member that lands on one is reported as
	// unverified below rather than passing quietly.
	at, loose := make(map[string]int, len(cl)), make(map[string]int, len(cl))
	for i, l := range cl {
		indexMember(at, normName(l.Path), i)
		indexMember(loose, normNameLoose(l.Path), i)
	}

	for i, l := range gl {
		if i >= len(cl) {
			break // Go-only tail, handled by diffLeafPrefix
		}

		j, ok := at[normName(l.Path)]
		if !ok || j == -1 {
			// Only now is the prefix allowed to be noise: mtmd spells C's
			// n_threads as Threads, which nothing in the exact index answers.
			if lj, lok := loose[normNameLoose(l.Path)]; lok {
				j, ok = lj, true
			}
		}

		switch {
		case !ok:
			unmatched = append(unmatched, fmt.Sprintf("%s.%s has no counterpart in %s", goName, l.Path, cName))
		case j == -1:
			unmatched = append(unmatched, fmt.Sprintf("%s.%s matches more than one member of %s once names are normalised, so its position is not checked", goName, l.Path, cName))
		case j == i:
		default:
			if len(probs) < maxLeafDiffs {
				probs = append(probs, fmt.Sprintf("member %d: %s.%s is C member %d (%s.%s), and member %d holds %s.%s: transposed, so each receives the other's value",
					i, goName, l.Path, j, cName, cl[j].Path, i, cName, cl[i].Path))
			}
		}
	}

	return probs, unmatched
}

// goTail describes Go-only members past the end of the C struct, for the
// STRUCT-BY-VALUE COMPARISONS report. They are legitimate, but they are also
// the reason a size comparison alone cannot be exact, so they stay visible.
func goTail(gl, fl []Leaf, goSize, cifSize int) string {
	if len(gl) <= len(fl) || goSize <= cifSize {
		return ""
	}

	names := make([]string, 0, len(gl)-len(fl))
	for _, l := range gl[len(fl):] {
		names = append(names, l.Path)
	}

	return fmt.Sprintf(" (+%dB Go-only tail past the C struct: %s)", goSize-cifSize, strings.Join(names, ", "))
}

func extraMembers(a, b []Leaf, aName, bName string) string {
	long, longName := a, aName
	if len(b) > len(a) {
		long, longName = b, bName
	}

	var s []string
	for _, l := range long[min(len(a), len(b)):] {
		if len(s) == maxLeafDiffs {
			s = append(s, "...")
			break
		}

		s = append(s, l.String())
	}

	return "only in " + longName + ": " + strings.Join(s, ", ")
}

func leafList(ls []Leaf) string {
	s := make([]string, 0, len(ls))
	for _, l := range ls {
		s = append(s, l.String())
	}

	return strings.Join(s, " ")
}
