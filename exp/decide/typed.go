package decide

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/ardanlabs/jinja"
)

const (
	typedContext = 8192

	typeOpenJev       = "openjev"
	typeLev           = "lev"
	typeLaya          = "laya"
	typeKev           = "kev"
	typePplxDecider   = "pplx-decider"
	typeLfm2D1        = "lfm2-d1"
	typeClef          = "clef"
	typeVllmDecision2 = "vllm-decision2"
	typeVllmD3        = "vllm-d3"
	typeNimble        = "nimble"
	typeGraniteSwitch = "granite-switch"
	typeLfm2D1Omni    = "lfm2-d1-omni"
	typeLfm3D1        = "lfm3-d1"
	typeLfm3D1Omni    = "lfm3-d1-omni"
	typeH2OLightning  = "h2o-lightning"

	openJevLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// levRatings is the scale of a Lev noul question, from 0 certainly no to 8 certainly yes.
	levRatings = 9
	// layaOptionTokens is the most tokens of one Laya option, as in training.
	layaOptionTokens = 48
)

var errLayout = errors.New("decide: unexpected layout of the decision prompt")

// Open loads a model whose GGUF holds its decision type, prompt template and
// temperatures, as converted for the llama-server /v1/systemone API.
// The openjev, lev, laya, kev, nimble, pplx-decider, lfm2-d1, lfm2-d1-omni,
// lfm3-d1, lfm3-d1-omni, granite-switch, h2o-lightning, clef, vllm-decision2 and
// vllm-d3 types are supported. Other models need [New], [NewJevK5] or
// [NewDeciderModel] and a config file. Clef and vllm-decision2 need ABI version 11
// of the llama.cpp module in the browser, h2o-lightning ABI version 12.
// Load and init llama.cpp first, and call [Decider.Close] when done.
func Open(modelPath string, opts Options) (*Decider, error) {
	if modelPath == "" {
		return nil, errors.New("decide: model path is required")
	}

	d, err := loadModel(modelPath, opts)
	if err != nil {
		return nil, err
	}

	t, err := newTyped(d)
	if err == nil {
		nCtx := typedContext
		if t.kind == typeLfm2D1Omni {
			nCtx = d1OmniMaxTokens
		}
		err = d.newContext(nCtx, opts, t.readout())
	}
	if err != nil {
		d.Close()
		return nil, err
	}
	d.family = t
	d.entropyConfidence = t.kind == typeVllmDecision2 || t.kind == typeVllmD3
	d.medianScore = t.isLfm3()
	// Clef reads the options of a choice question sorted, so their order does not matter.
	if t.kind == typeClef {
		d.bothOrders = false
	}

	return d, nil
}

// typed is a model described by its "<arch>.decision" metadata.
type typed struct {
	kind  string
	tmpl  *jinja.Template
	temps map[string]float64
	// bias is added to the scores of a score question, by its number of levels.
	bias   map[int][]float64
	enc    encoder
	plain  encoder
	piece  func(token) string
	labels []token
	// codes are the label texts the lev template shows.
	codes []string
	// maxOptions is the most options of one question.
	maxOptions int
	// marker is the token an option is read at, for laya and kev.
	marker token
	// sep, markerText and maxHead are the separator, the marker as text and
	// the token budget of the question and its options, for laya.
	sep        token
	markerText string
	maxHead    int
	// quoted are the special token texts that lfm3-d1 breaks and h2o-lightning
	// quotes in user text, longest first. quotedFirst has their first bytes.
	quoted      []string
	quotedFirst string
}

func (t *typed) isLfm3() bool {
	return t.kind == typeLfm3D1 || t.kind == typeLfm3D1Omni
}

// setQuoted keeps texts longest first, as llama-server matches them.
func (t *typed) setQuoted(texts []string) {
	slices.SortStableFunc(texts, func(a, b string) int { return len(b) - len(a) })
	t.quoted = texts
	for _, s := range texts {
		if !strings.Contains(t.quotedFirst, s[:1]) {
			t.quotedFirst += s[:1]
		}
	}
}

func newTyped(d *Decider) (*typed, error) {
	arch, _ := d.b.meta("general.architecture")
	kind, ok := d.b.meta(arch + ".decision.type")
	if !ok {
		return nil, errors.New("decide: the model has no decision type, load it with New, NewJevK5 or NewDeciderModel and a config file")
	}

	src := d.b.chatTemplate("systemone")
	if src == "" {
		return nil, errors.New("decide: the model has no systemone template")
	}
	tmpl, err := jinja.Compile(src)
	if err != nil {
		return nil, fmt.Errorf("decide: systemone template: %w", err)
	}

	// The values are read as float, as llama-server does.
	plain := d.tokenizer(false)
	t := &typed{kind: kind, tmpl: tmpl, temps: map[string]float64{}, bias: map[int][]float64{}, enc: d.tokenizer(true), plain: plain, piece: d.b.piece}
	for k, v := range d.b.metaPrefix(arch + ".decision.temperature.") {
		f, err := strconv.ParseFloat(v, 32)
		if err != nil || !(f > 0) {
			return nil, fmt.Errorf("decide: invalid decision temperature %s = %s", k, v)
		}
		t.temps[k] = f
	}
	// The keys are "<levels>.<level>".
	for k, v := range d.b.metaPrefix(arch + ".decision.score_bias.") {
		n, level, ok := parseBiasKey(k)
		f, err := strconv.ParseFloat(v, 32)
		if !ok || err != nil {
			return nil, fmt.Errorf("decide: invalid decision score bias %s = %s", k, v)
		}
		if t.bias[n] == nil {
			t.bias[n] = make([]float64, n)
		}
		t.bias[n][level] = f
	}

	switch kind {
	case typeOpenJev:
		for _, l := range openJevLetters {
			ids := plain(string(l))
			if len(ids) != 1 {
				return nil, fmt.Errorf("tokenizer mismatch: letter %q gives %v, want one token", l, ids)
			}
			t.labels = append(t.labels, ids[0])
		}
		t.maxOptions = len(t.labels)
	case typeLev, typePplxDecider, typeVllmD3, typeNimble, typeGraniteSwitch:
		for _, c := range labelCodes() {
			if ids := plain(c); len(ids) == 1 && len(t.labels) < deciderMaxOptions {
				t.labels = append(t.labels, ids[0])
				t.codes = append(t.codes, c)
			}
		}
		if len(t.labels) < levRatings {
			return nil, fmt.Errorf("tokenizer mismatch: %d label tokens, want at least %d", len(t.labels), levRatings)
		}
		t.maxOptions = len(t.labels)
	case typeLaya:
		t.marker, t.sep = d.b.vocabMask(), d.b.vocabSep()
		if t.marker < 0 || t.sep < 0 {
			return nil, errors.New("decide: the model has no mask or sep token")
		}
		t.markerText = d.b.piece(t.marker)
		v, _ := d.b.meta(arch + ".decision.max_head_tokens")
		if t.maxHead, _ = strconv.Atoi(v); t.maxHead <= 0 {
			return nil, errors.New("decide: the model has no valid max_head_tokens")
		}
		t.maxOptions = deciderMaxOptions
	case typeKev:
		ids := t.enc("<|box_end|>")
		if len(ids) != 1 {
			return nil, errors.New("decide: the model has no <|box_end|> token")
		}
		t.marker = ids[0]
		t.maxOptions = deciderMaxOptions
	case typeLfm2D1, typeClef, typeVllmDecision2:
		t.maxOptions = deciderMaxOptions
	case typeLfm3D1, typeLfm3D1Omni:
		t.setQuoted(d.b.specialTexts(false, 2))
		t.maxOptions = deciderMaxOptions
	case typeLfm2D1Omni:
		if t.marker = d.b.vocabMask(); t.marker < 0 {
			return nil, errors.New("decide: the model has no mask token")
		}
		t.maxOptions = deciderMaxOptions
	case typeH2OLightning:
		for i := uint32(0); ; i++ {
			l, err := d.b.decisionLabel(i)
			if err != nil {
				return nil, err
			}
			if l == "" {
				break
			}
			ids := plain(l)
			if len(ids) != 1 || slices.Contains(t.labels, ids[0]) {
				return nil, fmt.Errorf("tokenizer mismatch: decision label %q is not a single distinct token", l)
			}
			t.labels = append(t.labels, ids[0])
			t.codes = append(t.codes, l)
		}
		if len(t.labels) == 0 {
			return nil, errors.New("decide: the model has no decision labels")
		}
		t.setQuoted(d.b.specialTexts(true, 1))
		t.maxOptions = len(t.labels)
	default:
		return nil, fmt.Errorf("decide: decision type %q is not supported", kind)
	}

	return t, nil
}

// readout returns how the context of the model is set up.
func (t *typed) readout() readout {
	switch t.kind {
	case typeLaya:
		return readoutEncoder
	case typeKev:
		return readoutEmbd
	case typeVllmDecision2:
		return readoutHead
	case typeClef:
		return readoutJointHead
	case typeVllmD3, typeLfm3D1, typeLfm3D1Omni:
		return readoutWhole
	case typeLfm2D1Omni:
		return readoutEmbdWhole
	}
	return readoutLogits
}

// parseBiasKey parses "<levels>.<level>".
func parseBiasKey(k string) (n, level int, ok bool) {
	a, b, found := strings.Cut(k, ".")
	n, errN := strconv.Atoi(a)
	level, errL := strconv.Atoi(b)
	if !found || errN != nil || errL != nil || n < 2 || n > deciderMaxOptions || level < 0 || level >= n {
		return 0, 0, false
	}
	return n, level, true
}

// text cleans user text, so it holds no special token text and, for laya, no marker.
// lfm3-d1 breaks the special token texts instead, and h2o-lightning quotes them.
func (t *typed) text(s string) string {
	switch t.kind {
	case typeLfm3D1, typeLfm3D1Omni:
		return t.d1Break(s)
	case typeH2OLightning:
		return t.h2oQuote(s)
	case typeLfm2D1Omni:
		return strings.ReplaceAll(escapeSpecial(s), d1OmniMarker, "<<d1omni ")
	}
	s = escapeSpecial(s)
	if t.markerText != "" {
		s = strings.ReplaceAll(s, t.markerText, " ")
	}
	return s
}

// temperature returns the temperature for q, by its number of options when the model has one.
func (t *typed) temperature(q Question) float64 {
	n := len(q.Options)
	bucket := "11"
	switch {
	case t.kind == typeLev && n <= 8:
		bucket = "small"
	case t.kind == typeLev && n <= 26:
		bucket = "mid"
	case t.kind == typeLev:
		bucket = "large"
	case n <= 2:
		bucket = "2"
	case n <= 5:
		bucket = "3_5"
	case n <= 10:
		bucket = "6_10"
	}

	for _, k := range []string{string(q.Type) + "." + bucket, string(q.Type)} {
		if v, ok := t.temps[k]; ok {
			return v
		}
	}
	return 1
}

// typedOption is one option in the order the model reads it, and its place in [Question.Names].
type typedOption struct {
	key, description string
	index            int
}

// options returns the options of q as the model was trained on them.
// Some models list a noul question true first, clef sorts choice options by name.
func (t *typed) options(q Question) []typedOption {
	var out []typedOption
	switch q.Type {
	case TypeNoul:
		desc := map[string]string{}
		for _, o := range q.Options {
			desc[o.Name] = o.Description
		}
		out = []typedOption{{"false", desc["false"], 0}, {"true", desc["true"], 1}}
		trueFirst := len(q.Options) == 2 && q.Options[0].Name == "true"
		switch {
		case t.kind == typeOpenJev, t.kind == typeLfm2D1, t.kind == typeClef, t.kind == typeGraniteSwitch,
			t.kind == typeH2OLightning, t.isLfm3(), t.kind == typeVllmDecision2 && trueFirst:
			slices.Reverse(out)
		}
	case TypeScore:
		for i, o := range q.Options {
			out = append(out, typedOption{strconv.Itoa(i), o.Description, i})
		}
	default:
		for i, o := range q.Options {
			out = append(out, typedOption{o.Name, o.Description, i})
		}
		if t.kind == typeClef {
			slices.SortFunc(out, func(a, b typedOption) int { return strings.Compare(a.key, b.key) })
		}
	}
	return out
}

// variants returns the option orders q is read in. Lev reads a choice
// question a second time in reverse, to cancel its preference for the first label.
func (t *typed) variants(q Question) [][]typedOption {
	opts := t.options(q)
	if t.kind != typeLev || q.Type != TypeChoice || len(opts) < 2 {
		return [][]typedOption{opts}
	}
	rev := slices.Clone(opts)
	slices.Reverse(rev)
	return [][]typedOption{opts, rev}
}

// prompt renders the systemone template. state is already text, or nil for
// no state. labels replace the label codes of the model when given.
// qs are all the questions of the request, which the nimble prompt lists.
func (t *typed) prompt(state any, q Question, opts []typedOption, labels []string, qs []Question) (string, error) {
	if t.kind == typeH2OLightning {
		return t.h2oPrompt(state, q, opts)
	}

	in := map[string]any{
		"id":           q.id,
		"type":         string(q.Type),
		"instructions": t.text(q.Text),
		"state":        state,
		"options":      t.renderOptions(opts, labels),
		"images":       []any{},
	}
	if t.kind == typeNimble {
		var list []any
		for _, o := range qs {
			list = append(list, map[string]any{
				"id":           o.id,
				"type":         string(o.Type),
				"instructions": t.text(o.Text),
				"options":      t.renderOptions(t.options(o), nil),
			})
		}
		in["questions"] = list
	}
	if t.isLfm3() || t.kind == typeLfm2D1Omni {
		in["id"] = t.text(q.id)
	}
	if t.kind == typeLfm2D1Omni {
		in["audio"] = false
		in["sep"] = d1OmniSep
		in["mark_state"] = d1OmniMarkState
		in["mark_question"] = d1OmniMarkQuestion
		in["mark_option"] = d1OmniMarkOption
	}
	return t.tmpl.Render(in)
}

// renderOptions returns the options as the template takes them.
func (t *typed) renderOptions(opts []typedOption, labels []string) []any {
	options := make([]any, len(opts))
	for i, o := range opts {
		opt := map[string]any{"key": t.text(o.key), "description": nil}
		if o.description != "" {
			opt["description"] = t.text(o.description)
		}
		if t.codes != nil {
			opt["label"] = t.codes[i]
		}
		if labels != nil {
			opt["label"] = labels[i]
		}
		options[i] = opt
	}
	return options
}

// stateLen is the number of leading tokens that prompts about state share
// regardless of the question.
func (t *typed) stateLen(state any) (int, error) {
	a, err := t.prompt(state, Noul("a", "", ""), t.options(Noul("a", "", "")), nil, []Question{Noul("a", "", "")})
	if err != nil {
		return 0, err
	}
	b, err := t.prompt(state, Noul("b", "", ""), t.options(Noul("b", "", "")), nil, []Question{Noul("b", "", "")})
	if err != nil {
		return 0, err
	}

	ia, ib := t.enc(a), t.enc(b)
	n := 0
	for n < len(ia) && n < len(ib) && ia[n] == ib[n] {
		n++
	}
	return n, nil
}

func (t *typed) decide(d *Decider, state any, qs []Question, _ string, many bool) ([]*Result, error) {
	if t.kind == typeClef || t.kind == typeVllmDecision2 {
		return t.decideHead(d, state, qs, many)
	}

	// The nimble prompt names each question.
	if t.kind == typeNimble {
		qs = slices.Clone(qs)
		for i := range qs {
			if qs[i].id == "" {
				qs[i].id = "q" + strconv.Itoa(i+1)
			}
		}
	}

	s, err := t.state(state)
	if err != nil {
		return nil, err
	}
	stateLen := 0
	if t.kind != typeLaya && t.kind != typeLfm2D1Omni {
		if stateLen, err = t.stateLen(s); err != nil {
			return nil, err
		}
	}

	type plan struct {
		names    []string
		variants [][]typedOption
		first    int
		temp     float64
		// groups is the number of label tokens of each option, for lfm2-d1.
		groups []int
	}
	plans := make([]plan, len(qs))
	var rs []*rendered
	for i, q := range qs {
		// The vllm-d3 and h2o-lightning prompts have a default question.
		names, err := q.names(t.kind == typeVllmD3 || t.kind == typeH2OLightning)
		if err != nil {
			return nil, questionErr(many, i, err)
		}
		if limit := min(d.maxOptions, t.maxOptions); len(q.Options) > limit {
			return nil, questionErr(many, i, fmt.Errorf("%w: %d options, the limit is %d", ErrQuestion, len(q.Options), limit))
		}

		p := plan{names: names, variants: t.variants(q), first: len(rs), temp: t.temperature(q)}
		for _, opts := range p.variants {
			var labels []string
			var rows []token
			if t.kind == typeLfm2D1 || t.isLfm3() {
				var groups [][]token
				if labels, groups, err = t.d1Labels(q, opts); err != nil {
					return nil, questionErr(many, i, err)
				}
				for _, g := range groups {
					rows = append(rows, g...)
					p.groups = append(p.groups, len(g))
				}
			}
			text, err := t.prompt(s, q, opts, labels, qs)
			if err != nil {
				return nil, questionErr(many, i, err)
			}
			var r *rendered
			if t.kind == typeLfm2D1Omni {
				r, err = t.d1OmniRead(text, len(opts))
			} else {
				r, err = t.read(t.enc(text), q, len(opts), stateLen, rows)
			}
			if err != nil {
				return nil, questionErr(many, i, err)
			}
			if len(r.ids) > d.nCtx {
				return nil, questionErr(many, i, fmt.Errorf("%w: input needs %d tokens, the context is %d", ErrBudget, len(r.ids), d.nCtx))
			}
			rs = append(rs, r)
		}
		plans[i] = p
	}

	mode := manySeparate
	if many {
		mode = d.manyMode
	}
	vals, err := d.score(rs, mode)
	if err != nil {
		return nil, err
	}

	out := make([]*Result, len(qs))
	for i, p := range plans {
		probs := make([]float64, len(p.names))
		var scores []float64
		tokens := 0
		for v, opts := range p.variants {
			r := p.first + v
			tokens += len(rs[r].ids)
			logits, err := t.raw(vals[r], qs[i].Type)
			if err != nil {
				return nil, questionErr(many, i, err)
			}
			if p.groups != nil {
				logits = groupMax(logits, p.groups)
			}
			sm, err := softmax(t.biased(logits, qs[i].Type), p.temp)
			if err != nil {
				return nil, questionErr(many, i, err)
			}

			if t.kind == typeLev && qs[i].Type == TypeNoul {
				var yes float64
				for k, pk := range sm {
					yes += pk * float64(k) / (levRatings - 1)
				}
				probs[0], probs[1] = 1-yes, yes
				continue
			}

			for k, o := range opts {
				probs[o.index] += sm[k] / float64(len(p.variants))
			}
			if len(p.variants) == 1 {
				scores = make([]float64, len(p.names))
				for k, o := range opts {
					scores[o.index] = logits[k]
				}
			}
		}
		out[i] = newResult(p.names, probs, scores, p.temp, tokens, 0)
	}

	return out, nil
}

// state returns the state as the template takes it. lfm2-d1, lfm2-d1-omni,
// lfm3-d1 and vllm-d3 take a nil state.
func (t *typed) state(state any) (any, error) {
	if state == nil && t.takesNullState() {
		return nil, nil
	}

	var s string
	var err error
	switch t.kind {
	case typeKev:
		s, err = kevText(state)
	case typeLev:
		s, err = stateText(state, styleDefault.sorted())
	case typeLfm2D1, typeLfm3D1, typeLfm3D1Omni:
		s, err = stateText(state, styleIndent2)
	case typeH2OLightning:
		return t.h2oState(state)
	default:
		s, err = stateText(state, styleDefault)
	}
	if err != nil {
		return nil, err
	}
	return t.text(s), nil
}

// biased returns the scores of a score question with the score bias of the model added.
func (t *typed) biased(scores []float64, typ Type) []float64 {
	b, ok := t.bias[len(scores)]
	if typ != TypeScore || !ok {
		return scores
	}
	out := make([]float64, len(scores))
	for i, s := range scores {
		out[i] = s + b[i]
	}
	return out
}

// groupMax returns the highest score of each group of n[i] scores.
func groupMax(scores []float64, n []int) []float64 {
	out := make([]float64, len(n))
	for i, k := range n {
		out[i] = slices.Max(scores[:k])
		scores = scores[k:]
	}
	return out
}

// d1Labels returns the label of each lfm2-d1 and lfm3-d1 option and the tokens it is read at,
// as prompt.py of the model repo picks them.
func (t *typed) d1Labels(q Question, opts []typedOption) ([]string, [][]token, error) {
	single := func(forms ...string) []token {
		var out []token
		for _, f := range forms {
			if ids := t.plain(f); len(ids) == 1 && !slices.Contains(out, ids[0]) {
				out = append(out, ids[0])
			}
		}
		return out
	}

	var labels []string
	var groups [][]token
	if q.Type != TypeChoice {
		for _, o := range opts {
			var g []token
			switch {
			case q.Type == TypeScore:
				g = single(o.key)
			case o.key == "true":
				g = single("yes", "Yes", "YES")
			default:
				g = single("no", "No", "NO")
			}
			if len(g) == 0 {
				return nil, nil, fmt.Errorf("tokenizer mismatch: label %q is not one token", o.key)
			}
			labels = append(labels, o.key)
			groups = append(groups, g)
		}
		return labels, groups, nil
	}

	letters := true
	for _, o := range opts {
		letters = letters && len(o.key) == 1 && ('a' <= o.key[0] && o.key[0] <= 'z' || 'A' <= o.key[0] && o.key[0] <= 'Z')
	}
	var used []token
	take := func(code string) bool {
		ids := t.plain(code)
		if len(ids) != 1 || slices.Contains(used, ids[0]) {
			return false
		}
		used = append(used, ids[0])
		g := []token{ids[0]}
		for _, id := range single(" " + code) {
			if id != ids[0] {
				g = append(g, id)
			}
		}
		labels = append(labels, code)
		groups = append(groups, g)
		return true
	}

	pool := d1LabelPool()
	for i, o := range opts {
		code := fmt.Sprintf("%02d", i)
		switch {
		case letters:
			code = o.key
		case len(opts) <= 26, t.isLfm3() && i < 26:
			code = string(rune('A' + i))
		case t.isLfm3():
			code = string([]rune{rune('A' + (i-26)/26), rune('A' + (i-26)%26)})
		}
		if take(code) {
			continue
		}
		found := false
		for _, c := range pool {
			if found = take(c); found {
				break
			}
		}
		if !found {
			return nil, nil, fmt.Errorf("%w: no single token label left for %d options", ErrQuestion, len(opts))
		}
	}
	return labels, groups, nil
}

// d1LabelPool returns the labels lfm2-d1 falls back to, A to Z, 00 to 99, a to z, #0 to #199 and AA to ZZ.
func d1LabelPool() []string {
	var pool []string
	for c := 'A'; c <= 'Z'; c++ {
		pool = append(pool, string(c))
	}
	for i := range 100 {
		pool = append(pool, fmt.Sprintf("%02d", i))
	}
	for c := 'a'; c <= 'z'; c++ {
		pool = append(pool, string(c))
	}
	for i := range 200 {
		pool = append(pool, "#"+strconv.Itoa(i))
	}
	return append(pool, labelCodes()[26:]...)
}

// read returns the input of the tokens ids of a question with n options, and where to read it.
// rows replace the label tokens of the model when given.
func (t *typed) read(ids []token, q Question, n, stateLen int, rows []token) (*rendered, error) {
	last := len(ids) - 1
	switch t.kind {
	case typeLaya:
		ids, markers, err := t.layaFit(ids, n)
		if err != nil {
			return nil, err
		}
		return &rendered{ids: ids, slots: markers, embd: true}, nil
	case typeKev:
		var slots []int
		for i, id := range ids {
			if id == t.marker {
				slots = append(slots, i)
			}
		}
		if len(slots) != n {
			return nil, errLayout
		}
		// The question is read at the last token.
		return &rendered{ids: ids, slots: append(slots, last), embd: true, prefixLen: min(slots[0], stateLen)}, nil
	}

	switch {
	case rows != nil:
	case t.kind == typeLev && q.Type == TypeNoul:
		rows = t.labels[:levRatings]
	default:
		rows = t.labels[:n]
	}
	return &rendered{ids: ids, slots: []int{last}, rows: rows, prefixLen: min(last, stateLen)}, nil
}

// raw returns the raw score of each option from the outputs at the slots of one input.
func (t *typed) raw(out [][]float64, typ Type) ([]float64, error) {
	switch t.kind {
	case typeLaya, typeLfm2D1Omni:
		// The output has one score per question type.
		col := slices.Index([]Type{TypeChoice, TypeScore, TypeNoul}, typ)
		s := make([]float64, len(out))
		for i, e := range out {
			if col >= len(e) {
				return nil, errors.New("decide: the model output has no score for this question type")
			}
			s[i] = e[col]
		}
		return s, nil
	case typeKev:
		// Each output is [q | k], and an option scores q of the last token dot k of its marker.
		q := out[len(out)-1]
		h := len(q) / 2
		s := make([]float64, len(out)-1)
		for i, e := range out[:len(out)-1] {
			var dot float64
			for j := range h {
				dot += q[j] * e[h+j]
			}
			s[i] = dot / math.Sqrt(float64(h))
		}
		return s, nil
	}
	return out[0], nil
}

// layaFit cuts the question and options of a Laya prompt to the head budget,
// as the model was trained, and returns the tokens and the marker positions.
// The prompt is [cls] question [sep] ([mask] option)* [sep] state [sep].
func (t *typed) layaFit(ids []token, n int) ([]token, []int, error) {
	var markers []int
	for i, id := range ids {
		if id == t.marker {
			markers = append(markers, i)
		}
	}
	if len(markers) != n || markers[0] < 2 || ids[markers[0]-1] != t.sep || ids[len(ids)-1] != t.sep {
		return nil, nil, errLayout
	}
	headEnd := markers[0] - 1
	optsEnd := len(ids)
	if k := slices.Index(ids[markers[n-1]:], t.sep); k >= 0 {
		optsEnd = markers[n-1] + k
	}
	if optsEnd+1 >= len(ids) {
		return nil, nil, errLayout
	}

	opts := make([][]token, n)
	for i := range opts {
		end := optsEnd
		if i+1 < n {
			end = markers[i+1]
		}
		opts[i] = ids[markers[i]:end]
	}
	total := 0
	setMax := func(m int) {
		total = 0
		for i := range opts {
			opts[i] = opts[i][:min(len(opts[i]), m)]
			total += len(opts[i])
		}
	}
	setMax(layaOptionTokens + 1)
	if total+16 > t.maxHead {
		setMax(max(4, (t.maxHead-min(t.maxHead, 16))/n))
	}
	question := max(8, t.maxHead-min(t.maxHead, total))

	out := append([]token{ids[0]}, ids[1:min(headEnd, 1+question)]...)
	out = append(out, t.sep)
	pos := make([]int, n)
	for i, o := range opts {
		pos[i] = len(out)
		out = append(out, o...)
	}
	return append(out, ids[optsEnd:]...), pos, nil
}
