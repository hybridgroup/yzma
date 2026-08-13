package main

// C constant values: enum members and integer #defines.
//
// RULES 1-3 check the *shape* of every call - widths, classes, offsets - and
// none of them looks at a single value. But yzma also mirrors every llama.cpp
// enum member and every interesting #define as a Go constant, by hand, and
// those values are arguments: LLAMA_POOLING_TYPE_MEAN is the integer 1 because
// of where it sits in its enum, not because anything declares it. Insert a
// member upstream and every later member shifts by one, both sides still
// compile, every width still matches, and yzma starts asking for a different
// pooling type than the caller named. That is the same silent drift as the
// struct-field bug in hybridgroup/yzma#289, one layer up.
//
// So the values have to be read out of the headers and compared. Everything
// here is deliberately conservative: an expression this file cannot evaluate is
// recorded with the reason, and surfaces as a NOT VERIFIED skip if a Go
// constant maps to it, rather than being dropped or guessed at.

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"maps"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// CConst is one C constant yzma may mirror: an enum member or an integer
// #define.
type CConst struct {
	Name string
	Val  int64
	Enum string // enum tag it belongs to, "" for a #define
	File string
	Line int
}

func (c CConst) where() string {
	if c.Enum != "" {
		return fmt.Sprintf("enum %s, %s:%d", c.Enum, filepath.Base(c.File), c.Line)
	}

	return fmt.Sprintf("#define, %s:%d", filepath.Base(c.File), c.Line)
}

var (
	cconsts = map[string]CConst{}

	// cconstBad records, per name, why a constant's value could not be pinned
	// down: an expression this file cannot evaluate, a name defined twice with
	// different values, or an enum member whose predecessor was already
	// unknown. Keeping the reason is what turns "not checked" into a skip a
	// human can act on instead of a hole.
	cconstBad = map[string]string{}

	// cconstByNorm indexes every known name by its normalised form, so a Go
	// constant can be matched without a comment naming its C counterpart. A
	// normalised form shared by several C names is ambiguous and reported as
	// such rather than resolved by preference.
	cconstByNorm = map[string][]string{}

	// cconstBadByNorm is the same index over the unevaluable names, so a Go
	// constant that mirrors one is told why it could not be checked instead of
	// being told nothing matched.
	cconstBadByNorm = map[string][]string{}

	// cconstMirrored records every C constant a Go constant was matched to in
	// this run. It is what the partially-mirrored enum inventory subtracts from:
	// RULE 4 walks the Go side, so a C member nothing mirrors is invisible to it.
	cconstMirrored = map[string]bool{}

	// goEnumTags records, per Go constant type, the C enums its mirrored members
	// belong to. It is the by-product that makes the enum-parameter check
	// possible: RULE 4 already resolves PoolingTypeMean to
	// LLAMA_POOLING_TYPE_MEAN, and that member's enum tag is the one C enum
	// yzma's PoolingType stands for. See goEnumOf and cmpEnumType.
	goEnumTags = map[string]map[string]bool{}
)

func resetCConsts() {
	clear(cconsts)
	clear(cconstBad)
	clear(cconstByNorm)
	clear(cconstBadByNorm)
	clear(cconstMirrored)
	clear(goEnumTags)
}

func markCConstBad(name, why string) {
	cconstBad[name] = why

	for _, k := range constKeys(name) {
		if !contains(cconstBadByNorm[k], name) {
			cconstBadByNorm[k] = append(cconstBadByNorm[k], name)
		}
	}
}

// collectCConsts parses the enum members and integer #defines out of one
// comment-stripped header.
//
// Members and defines are evaluated in source order, because both can be
// initialised from an earlier constant (LLAMA_SESSION_MAGIC is
// LLAMA_FILE_MAGIC_GGSN, LLAMA_ROPE_SCALING_TYPE_MAX_VALUE is
// LLAMA_ROPE_SCALING_TYPE_LONGROPE), and headers are parsed in dependency
// order.
func collectCConsts(path, src string) {
	type item struct {
		off  int
		body func(line int)
	}

	var items []item

	for _, idx := range findAll(src, "enum") {
		if idx > 0 && isIdentChar(src[idx-1]) {
			continue
		}

		if after := idx + len("enum"); after < len(src) && isIdentChar(src[after]) {
			continue
		}

		open := strings.Index(src[idx:], "{")
		if open < 0 {
			continue
		}

		// A `{` further away than the tag name means this is a use of the enum
		// type (a parameter, a field) rather than a definition.
		tag := strings.TrimSpace(src[idx+len("enum") : idx+open])
		if strings.ContainsAny(tag, " \t\n") || strings.Contains(tag, "*") {
			continue
		}

		close := strings.Index(src[idx+open:], "}")
		if close < 0 {
			continue
		}

		body := src[idx+open+1 : idx+open+close]
		items = append(items, item{idx, func(line int) { evalEnumBody(path, line, tag, body) }})
	}

	for _, idx := range findAll(src, "#") {
		nl := strings.IndexByte(src[idx:], '\n')
		if nl < 0 {
			nl = len(src) - idx
		}

		line := strings.TrimSpace(src[idx+1 : idx+nl])
		rest, ok := strings.CutPrefix(line, "define")
		if !ok || (rest != "" && isIdentChar(rest[0])) {
			continue
		}

		// A trailing backslash continues the macro onto the next line, which
		// is never a plain integer.
		if strings.HasSuffix(strings.TrimSpace(rest), "\\") {
			continue
		}

		f := strings.Fields(rest)
		if len(f) < 2 || strings.Contains(f[0], "(") {
			continue
		}

		name, expr := f[0], strings.Join(f[1:], " ")
		items = append(items, item{idx, func(line int) { defineCConst(name, expr, path, line) }})
	}

	sort.SliceStable(items, func(i, j int) bool { return items[i].off < items[j].off })

	for _, it := range items {
		it.body(1 + strings.Count(src[:it.off], "\n"))
	}
}

// evalEnumBody walks one enum body, tracking the implicit counter C uses for
// members without an initialiser.
func evalEnumBody(path string, line int, tag, body string) {
	next := int64(0)
	broken := "" // the member that stopped the implicit counter being knowable
	at := 0      // offset of the current member within body, for its line number

	for _, m := range splitTop(body) {
		// Point at the member name rather than at the comma before it, so the
		// reported line is the member's own.
		off := at + len(m) - len(strings.TrimLeft(m, " \t\r\n"))
		at += len(m) + 1

		name, expr, explicit := strings.Cut(strings.TrimSpace(m), "=")
		name = strings.TrimSpace(name)
		if name == "" || !isIdentChar(name[0]) || strings.ContainsAny(name, " \t") {
			continue
		}

		memberLine := line + strings.Count(body[:off], "\n")

		if explicit {
			v, err := evalCExpr(expr)
			if err != nil {
				markCConstBad(name, fmt.Sprintf("enum %s member %s = %q could not be evaluated: %v", tag, name, squash(expr), err))
				broken = name

				continue
			}

			next, broken = v, ""
		} else if broken != "" {
			markCConstBad(name, fmt.Sprintf("enum %s member %s follows %s, whose value is unknown, so its implicit value is too", tag, name, broken))

			continue
		}

		recordCConst(CConst{Name: name, Val: next, Enum: tag, File: path, Line: memberLine})
		next++
	}
}

// defineCConst records an integer #define.
func defineCConst(name, expr, path string, line int) {
	v, err := evalCExpr(expr)
	if err != nil {
		// Not every #define is an integer: GGML_API is `extern` and
		// GGML_RESTRICT is a keyword. Those are recorded like any other
		// unevaluable name, and only ever surface if a Go constant claims to
		// mirror one.
		if _, seen := cconsts[name]; !seen {
			markCConstBad(name, fmt.Sprintf("#define %s %s is not an integer expression: %v", name, squash(expr), err))
		}

		return
	}

	recordCConst(CConst{Name: name, Val: v, File: path, Line: line})
}

// recordCConst files a resolved constant, rejecting a name whose definitions
// disagree.
func recordCConst(c CConst) {
	if _, bad := cconstBad[c.Name]; bad {
		return
	}

	if prev, seen := cconsts[c.Name]; seen {
		if prev.Val == c.Val {
			return
		}

		// GGML_MEM_ALIGN and friends are defined per platform inside #if
		// branches. The checker does not evaluate the preprocessor, so it
		// cannot say which value this build compiles with.
		markCConstBad(c.Name, fmt.Sprintf("%s is defined more than once with different values (%d at %s, %d at %s)",
			c.Name, prev.Val, prev.where(), c.Val, c.where()))
		delete(cconsts, c.Name)

		return
	}

	cconsts[c.Name] = c
	indexCConst(c.Name)
}

func indexCConst(name string) {
	for _, k := range constKeys(name) {
		if !contains(cconstByNorm[k], name) {
			cconstByNorm[k] = append(cconstByNorm[k], name)
		}
	}
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}

	return false
}

// cLibPrefixes are the three library prefixes every C name here carries and
// yzma's Go names sometimes drop: LLAMA_POOLING_TYPE_MEAN is PoolingTypeMean,
// while GGML_TYPE_F32 keeps its prefix as GGMLTypeF32. Both spellings are
// indexed, so neither convention needs guessing at.
var cLibPrefixes = []string{"llama", "ggml", "mtmd"}

// constKeys returns every normalised form a C name may be matched by.
func constKeys(name string) []string {
	n := normConst(name)
	keys := []string{n}

	for _, p := range cLibPrefixes {
		if rest, ok := strings.CutPrefix(n, p); ok && rest != "" {
			keys = append(keys, rest)
			break
		}
	}

	return keys
}

// constAliases reconciles the words yzma spells out where llama.cpp abbreviates
// them, or the other way round. Each pair is rewritten to its second, shorter
// form on both sides so the two conventions meet in the middle — the same
// device, and the same discipline, as memberAliases in layout.go: a closed
// hand-checked list, because an alias loose enough to guess a name is loose
// enough to validate a constant against the wrong C member.
var constAliases = [][2]string{
	{"attention", "attn"},      // FlashAttentionTypeAuto / LLAMA_FLASH_ATTN_TYPE_AUTO
	{"userdefined", "userdef"}, // TokenAttrUserDef / LLAMA_TOKEN_ATTR_USER_DEFINED
	{"continue", "cont"},       // LogLevelContinue / GGML_LOG_LEVEL_CONT
	{"probability", "prob"},    // ...SamplingXTCProb / ..._SAMPLING_XTC_PROBABILITY
	{"threshold", "thold"},     // ...SamplingXTCThold / ..._SAMPLING_XTC_THRESHOLD
}

// normConst reduces a constant name to a form comparable across the two naming
// conventions, the way normName does for struct members: GGML_TYPE_Q4_0 and
// GGMLTypeQ4_0 both become "ggmltypeq40".
//
// Unlike normName it strips no counting prefix. Member names are matched
// positionally against a struct whose shape is already known, so a loose match
// there costs little; here the match is the whole check, and two constants that
// normalise together would silently validate each other's value.
func normConst(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '_' {
			continue
		}

		b.WriteRune(unicode.ToLower(r))
	}

	s = b.String()

	for _, a := range constAliases {
		s = strings.ReplaceAll(s, a[0], a[1])
	}

	return s
}

// evalCExpr evaluates a C integer constant expression: decimal and hex
// literals with the usual suffixes, references to constants already parsed,
// and the operators llama.h and ggml.h actually use in enum initialisers
// (1 << 4, (1 << 8), -1).
func evalCExpr(s string) (int64, error) {
	p := &cexpr{toks: tokenizeC(s)}

	v, err := p.parse(0)
	if err != nil {
		return 0, err
	}

	if p.i != len(p.toks) {
		return 0, fmt.Errorf("trailing %q", strings.Join(p.toks[p.i:], " "))
	}

	return v, nil
}

func tokenizeC(s string) []string {
	var out []string
	for i := 0; i < len(s); {
		switch c := s[i]; {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case isIdentChar(c):
			j := i
			for j < len(s) && isIdentChar(s[j]) {
				j++
			}

			out = append(out, s[i:j])
			i = j
		case (c == '<' || c == '>') && i+1 < len(s) && s[i+1] == c:
			out = append(out, s[i:i+2])
			i += 2
		default:
			out = append(out, s[i:i+1])
			i++
		}
	}

	return out
}

// cexpr is a precedence-climbing evaluator over the token list.
type cexpr struct {
	toks []string
	i    int
}

// cPrec is the subset of C's precedence table these headers need. Anything not
// listed makes the expression unevaluable, which is a skip rather than a guess.
var cPrec = map[string]int{
	"|": 1, "^": 2, "&": 3,
	"<<": 4, ">>": 4,
	"+": 5, "-": 5,
	"*": 6, "/": 6,
}

func (p *cexpr) parse(minPrec int) (int64, error) {
	lhs, err := p.unary()
	if err != nil {
		return 0, err
	}

	for p.i < len(p.toks) {
		prec, ok := cPrec[p.toks[p.i]]
		if !ok || prec < minPrec {
			break
		}

		op := p.toks[p.i]
		p.i++

		rhs, err := p.parse(prec + 1)
		if err != nil {
			return 0, err
		}

		switch op {
		case "|":
			lhs |= rhs
		case "^":
			lhs ^= rhs
		case "&":
			lhs &= rhs
		case "<<":
			if rhs < 0 || rhs > 62 {
				return 0, fmt.Errorf("shift by %d", rhs)
			}

			lhs <<= uint(rhs)
		case ">>":
			if rhs < 0 || rhs > 62 {
				return 0, fmt.Errorf("shift by %d", rhs)
			}

			lhs >>= uint(rhs)
		case "+":
			lhs += rhs
		case "-":
			lhs -= rhs
		case "*":
			lhs *= rhs
		case "/":
			if rhs == 0 {
				return 0, fmt.Errorf("division by zero")
			}

			lhs /= rhs
		}
	}

	return lhs, nil
}

func (p *cexpr) unary() (int64, error) {
	if p.i >= len(p.toks) {
		return 0, fmt.Errorf("empty expression")
	}

	switch t := p.toks[p.i]; t {
	case "-":
		p.i++
		v, err := p.unary()

		return -v, err
	case "+":
		p.i++

		return p.unary()
	case "~":
		p.i++
		v, err := p.unary()

		return ^v, err
	case "(":
		p.i++
		v, err := p.parse(0)
		if err != nil {
			return 0, err
		}

		if p.i >= len(p.toks) || p.toks[p.i] != ")" {
			return 0, fmt.Errorf("unclosed (")
		}

		p.i++

		return v, nil
	}

	t := p.toks[p.i]
	p.i++

	if t[0] >= '0' && t[0] <= '9' {
		return parseCInt(t)
	}

	if !isIdentChar(t[0]) {
		return 0, fmt.Errorf("unexpected %q", t)
	}

	if c, ok := cconsts[t]; ok {
		return c.Val, nil
	}

	if why, bad := cconstBad[t]; bad {
		return 0, fmt.Errorf("refers to %s, itself unknown (%s)", t, why)
	}

	return 0, fmt.Errorf("refers to unknown identifier %s", t)
}

// parseCInt handles the u/U/l/L literal suffixes, which appear on the
// LLAMA_FILE_MAGIC_* defines.
func parseCInt(t string) (int64, error) {
	s := strings.TrimRight(t, "uUlL")

	v, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		return v, nil
	}

	// 0xFFFFFFFFFFFFFFFF and friends do not fit a signed 64-bit value but are
	// still exact bit patterns.
	if u, uerr := strconv.ParseUint(s, 0, 64); uerr == nil {
		return int64(u), nil
	}

	return 0, fmt.Errorf("literal %q: %w", t, err)
}

// GoConst is one package-level integer constant in an audited package.
type GoConst struct {
	Name  string
	Type  string // the Go named type, "" for an untyped constant
	Pos   token.Position
	Val   int64
	CName string // the C name from a `// GGML_TYPE_F32` doc line, "" if absent
}

// collectGoConsts gathers every package-level integer constant of a package,
// with the C name from the comment line immediately above it when there is one.
//
// Values come from go/types rather than from the source text, so `1 << 6`,
// `iota` and `RopeScalingTypeLongROPE` all arrive as the integer the compiler
// computed.
func collectGoConsts(p *packages.Package) []GoConst {
	var out []GoConst

	for _, f := range p.Syntax {
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}

			for _, s := range gd.Specs {
				vs := s.(*ast.ValueSpec)

				doc := vs.Doc
				if doc == nil && len(gd.Specs) == 1 {
					doc = gd.Doc
				}

				for _, n := range vs.Names {
					c, ok := p.TypesInfo.Defs[n].(*types.Const)
					if !ok || n.Name == "_" {
						continue
					}

					v := constant.ToInt(c.Val())
					if v.Kind() != constant.Int {
						continue
					}

					iv, exact := constant.Int64Val(v)
					if !exact {
						continue
					}

					gc := GoConst{Name: n.Name, Pos: p.Fset.Position(n.Pos()), Val: iv, CName: docCName(doc)}
					if named, ok := c.Type().(*types.Named); ok {
						gc.Type = named.Obj().Name()
					}

					out = append(out, gc)
				}
			}
		}
	}

	return out
}

// docCName reports the C constant a doc comment names, which yzma writes as a
// bare `// GGML_TYPE_F32` line directly above the constant.
//
// Only the last line counts. The comment block above GGMLTypeQ5_0 also lists
// GGML_TYPE_Q4_2 and GGML_TYPE_Q4_3 as removed types, and a scan of the whole
// block would have three candidates to choose between.
func docCName(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}

	last := strings.TrimSpace(strings.TrimPrefix(doc.List[len(doc.List)-1].Text, "//"))
	if last == "" || !isIdentChar(last[0]) || strings.ContainsAny(last, " \t") {
		return ""
	}

	for _, r := range last {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) && r != '_' {
			return ""
		}
	}

	return last
}

// constEnumTags names the C enum a Go constant type mirrors, for the cases
// where the name alone is ambiguous. It only needs an entry where two C enums
// declare members that normalise to the same string.
var constEnumTags = map[string]string{
	// llama_ftype and ggml_ftype are parallel quantisation enumerations with
	// overlapping member names and different values (LLAMA_FTYPE_MOSTLY_Q8_0 is
	// 7, GGML_FTYPE_MOSTLY_Q8_0 is 7 today but the two lists have drifted
	// apart before). yzma's Ftype is the llama.h one: it is what
	// llama_model_quantize_params.ftype takes.
	"Ftype": "llama_ftype",

	// llama_rope_type initialises four of its members straight from the
	// GGML_ROPE_TYPE_* defines, so both spellings exist with the same value.
	// yzma's RoPEType is the llama.h enum, which is what llama_model_rope_type
	// returns.
	"RoPEType": "llama_rope_type",
}

// goOnlyConsts are the audited packages' own constants: names yzma defines for
// its own API rather than to mirror a C constant, so there is nothing in the
// headers to compare them against.
//
// Like memberAliases in layout.go this is a closed, hand-maintained list rather
// than a pattern, and for the same reason: a rule loose enough to recognise
// "this one is ours" would also excuse a real mirror that stopped matching. A
// new yzma-local constant therefore shows up as a skip needing one line added
// here, never as a silent pass. The names are matched exactly.
var goOnlyConsts = map[string]bool{
	// yzma's own GPU-backend enumeration, used to pick which llama.cpp shared
	// library to load. llama.cpp has no equivalent enum.
	"GpuBackendNone": true, "GpuBackendCPU": true, "GpuBackendCUDA": true,
	"GpuBackendMetal": true, "GpuBackendHIP": true, "GpuBackendVulkan": true,
	"GpuBackendOpenCL": true, "GpuBackendSYCL": true,

	// Sampler kinds are yzma's own dispatch tags: llama.cpp exposes one
	// llama_sampler_init_* function per sampler and no enum over them.
	"SamplerTypeNone": true, "SamplerTypeDry": true, "SamplerTypeTopK": true,
	"SamplerTypeTopP": true, "SamplerTypeMinP": true, "SamplerTypeTypicalP": true,
	"SamplerTypeTemperature": true, "SamplerTypeXTC": true, "SamplerTypeInfill": true,
	"SamplerTypePenalties": true, "SamplerTypeTopNSigma": true,
	"SamplerTypeAdaptiveP": true, "SamplerTypeLogitBias": true,

	// mtmd_batch_add_chunk returns bare int32 codes documented in a comment on
	// the declaration; BatchAddResult is yzma naming them.
	"BatchAddSuccess": true, "BatchAddError": true,
	"BatchAddTooLarge": true, "BatchAddIncompatible": true,

	// Conveniences yzma adds on top of the C API.
	"MaxToken":  true, // the largest llama_token, not declared in llama.h
	"LogNormal": true, // the "no user data" sentinel for the log callback
}

// checkConsts implements RULE 4: every mirrored Go constant must still hold the
// value its C counterpart has in the headers this run parsed.
func checkConsts(loaded []*packages.Package) (viols []violation, checked, clean, local int) {
	for _, p := range loaded {
		if p.Syntax == nil {
			continue
		}

		for _, g := range collectGoConsts(p) {
			if goOnlyConsts[g.Name] {
				local++
				continue
			}

			c, why := matchCConst(g)
			if why != "" {
				noteSkip("RULE4 %s (%s): %s", g.Name, shortPos(g.Pos), why)
				continue
			}

			checked++
			cconstMirrored[c.Name] = true

			// Which C enum a Go enum type stands for is not stated anywhere in
			// either language; it is only observable through the members, and
			// this is where they have just been matched one to one.
			if g.Type != "" && c.Enum != "" {
				if goEnumTags[g.Type] == nil {
					goEnumTags[g.Type] = map[string]bool{}
				}

				goEnumTags[g.Type][c.Enum] = true
			}

			if *verbose {
				fmt.Printf("CONST %-34s = %-12d %s (%s)\n", g.Name, g.Val, c.Name, c.where())
			}

			if c.Val == g.Val {
				clean++
				continue
			}

			viols = append(viols, violation{4, c.Name, shortPos(g.Pos),
				fmt.Sprintf("%s = %d but C %s = %d (%s): every value passed through this constant is off by %d",
					g.Name, g.Val, c.Name, c.Val, c.where(), c.Val-g.Val),
				"wrong-value"})
		}
	}

	return viols, checked, clean, local
}

// EnumCoverage is how much of one C enum yzma mirrors: an inventory line, never
// a finding, in the same spirit as hdrCoverage.
type EnumCoverage struct {
	Enum     string
	Members  int
	Mirrored int
	Missing  []string // member name, value and header line, in declaration order
}

// partialEnums inventories the members of every partially mirrored C enum: the
// ones with no Go constant at all.
//
// RULE 4 walks the Go side, so a C member yzma never transcribed cannot be
// compared and is invisible to the rule. Usually that is deliberate - yzma
// mirrors a chosen subset of llama.cpp, exactly as it binds a chosen subset of
// its functions - so this never becomes a violation and never affects the exit
// code. What it is, is the signal for the one event RULE 4 exists to catch: an
// enum gaining a member upstream is both a new unmirrored name here and a value
// shift in every mirrored member after it, and only the first of those two is
// visible before someone renumbers.
//
// Only enums with at least one mirrored member are listed. An enum with none is
// not a partial mirror, it is simply unused, and the several hundred members of
// the ggml enums yzma never touches would drown the report.
func partialEnums() []EnumCoverage {
	byEnum := map[string][]CConst{}
	for _, c := range cconsts {
		if c.Enum == "" {
			continue
		}

		byEnum[c.Enum] = append(byEnum[c.Enum], c)
	}

	var out []EnumCoverage
	for _, tag := range slices.Sorted(maps.Keys(byEnum)) {
		ms := byEnum[tag]
		slices.SortFunc(ms, func(a, b CConst) int { return a.Line - b.Line })

		ec := EnumCoverage{Enum: tag, Members: len(ms)}
		for _, m := range ms {
			if cconstMirrored[m.Name] {
				ec.Mirrored++
				continue
			}

			ec.Missing = append(ec.Missing, fmt.Sprintf("%s = %d (%s:%d)", m.Name, m.Val, filepath.Base(m.File), m.Line))
		}

		if ec.Mirrored == 0 || len(ec.Missing) == 0 {
			continue
		}

		out = append(out, ec)
	}

	return out
}

// cEnumOf names the C enum a parameter type is, "" for anything else.
//
// classify() reduces every enum to KindSint/4, which is all the ABI needs and
// exactly what makes one enum indistinguishable from another - see cmpEnumType.
// The tag is the only thing that tells them apart, so it is recovered here, out
// of the `enum <tag>` spelling llama.h uses and through a typedef chain for the
// spellings it does not.
func cEnumOf(cRaw string) string {
	t := squash(strings.NewReplacer("const", " ", "volatile", " ", "restrict", " ").Replace(cRaw))
	if strings.ContainsAny(t, "*[") {
		return ""
	}

	for range 12 {
		if tag, ok := strings.CutPrefix(t, "enum "); ok {
			return strings.TrimSpace(tag)
		}

		if enumNames[t] {
			return t
		}

		u, ok := typedefs[t]
		if !ok {
			return ""
		}

		t = squash(u)
	}

	return ""
}

// goEnumOf names the C enum a Go type mirrors, "" where nothing keys on it.
//
// The evidence is RULE 4's: a Go type whose mirrored members all belong to one C
// enum is that enum, and constEnumTags already records the two types whose names
// alone are ambiguous. A type with no mirrored member at all - and a type whose
// members straddle two C enums, which no yzma type does today - names nothing to
// compare against, so the caller makes no claim rather than a guess.
func goEnumOf(goType string) string {
	if tag, ok := constEnumTags[goType]; ok {
		return tag
	}

	if tags := goEnumTags[goType]; len(tags) == 1 {
		for tag := range tags {
			return tag
		}
	}

	return ""
}

// matchCConst resolves the C constant a Go constant mirrors, preferring the
// name yzma wrote in the comment above it over any normalisation.
//
// It never picks between candidates: a name matching nothing, or matching more
// than one, is returned as a reason to skip. Guessing here would mean checking
// a value against the wrong constant, which is worse than not checking it.
func matchCConst(g GoConst) (CConst, string) {
	if g.CName != "" {
		if c, ok := cconsts[g.CName]; ok {
			return c, ""
		}

		if why, bad := cconstBad[g.CName]; bad {
			return CConst{}, why
		}

		return CConst{}, fmt.Sprintf("comment names %s, which is not an enum member or #define in the headers", g.CName)
	}

	// A name recordCConst later dropped - because a second #define disagreed
	// with the first - stays in cconstByNorm, since the index is written when
	// the first definition lands and never unwound. Such a name is not a
	// candidate: cconsts no longer holds it, so taking it would return a
	// zero-valued CConst with no reason attached, count towards checked, and
	// compare the Go value against 0. Filtering it here is also what lets the
	// case 0 arm below find the real reason in cconstBadByNorm.
	var cands []string
	for _, k := range constKeys(g.Name) {
		for _, n := range cconstByNorm[k] {
			if _, bad := cconstBad[n]; bad {
				continue
			}

			if !contains(cands, n) {
				cands = append(cands, n)
			}
		}
	}

	// ggml.h and llama.h each declare a parallel quantisation enum, so
	// FtypeMostlyQ4_0 normalises onto both GGML_FTYPE_MOSTLY_Q4_0 and
	// LLAMA_FTYPE_MOSTLY_Q4_0 — which do not hold the same values. The Go type
	// says which enum is being mirrored, so use it to narrow the candidates.
	if tag, ok := constEnumTags[g.Type]; ok && len(cands) > 1 {
		var keep []string
		for _, n := range cands {
			if cconsts[n].Enum == tag {
				keep = append(keep, n)
			}
		}

		cands = keep
	}

	switch len(cands) {
	case 1:
		return cconsts[cands[0]], ""
	case 0:
		for _, k := range constKeys(g.Name) {
			if bad := cconstBadByNorm[k]; len(bad) > 0 {
				return CConst{}, fmt.Sprintf("mirrors %s: %s", bad[0], cconstBad[bad[0]])
			}
		}

		return CConst{}, "no C enum member or #define matches this name - add a `// C_NAME` comment above it, or list it in goOnlyConsts if it is yzma's own"
	default:
		sort.Strings(cands)

		return CConst{}, fmt.Sprintf("matches %d C names (%s) - add a `// C_NAME` comment above it to say which, or map its Go type in constEnumTags",
			len(cands), strings.Join(cands, ", "))
	}
}
