package decide

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"slices"
	"sync"
)

// Options configures a [Decider]. The zero value is valid.
type Options struct {
	// Threads is the number of CPU threads. 0 uses the llama.cpp default.
	Threads int32
	// UBatch is the physical batch size. 0 uses 1024, the size the model was validated with.
	UBatch uint32
	// MaxOptions is the most options one question can have. 0 uses 256.
	MaxOptions uint32
	// ManyMode is how [Decider.DecideMany] shares the state. The default is ManyExact.
	ManyMode ManyMode
	// ContextSize is the context in tokens. 0 uses the model family default,
	// max_len for Jev-Style and 8192 for JevK5 and decider models.
	ContextSize uint32
	// BothOrders also reads each choice question with its options reversed and
	// averages the two, which cancels a preference for the first options.
	BothOrders bool
}

// ManyMode is how [Decider.DecideMany] shares one state across questions.
type ManyMode int

const (
	// ManyExact decodes the full ubatches of the state once and each question
	// separately. Results are identical to [Decider.Decide].
	ManyExact ManyMode = iota
	// ManyBatched decodes the whole state once and up to 16 questions in one
	// decode. It is faster, but probabilities can differ from [Decider.Decide]
	// by a few hundredths and a near tie can change the answer. For a JevK5
	// question with more than 16 options a near tie can also change which
	// options reach the final pass, which moves its probabilities more.
	ManyBatched

	// manySeparate decodes each question separately from an empty memory.
	manySeparate ManyMode = -1
)

// maxSeqs is the number of sequences, the shared state plus up to 16 questions.
const maxSeqs = 17

// Result is the decision for one question.
// Options, Probabilities and Scores are in [Question.Names] order.
// Scores are the raw scores before calibration. They are empty for a JevK5
// question with more than 16 options and for a decider score question with
// isolated levels, which combine several passes.
// Confidence and Expected follow the TypeSafe API. Expected is the mean level
// of a score question and 0 for other questions.
type Result struct {
	Answer               string    `json:"answer"`
	Options              []string  `json:"options"`
	Probabilities        []float64 `json:"probabilities"`
	Scores               []float64 `json:"scores,omitempty"`
	Temperature          float64   `json:"temperature"`
	TopProbability       float64   `json:"top_probability"`
	EntropyConcentration float64   `json:"entropy_concentration"`
	Confidence           float64   `json:"confidence"`
	Expected             float64   `json:"expected"`
	InputTokens          int       `json:"input_tokens"`
	HeadTokens           int       `json:"head_tokens,omitempty"`
}

// Probability returns the probability of the named option, or 0 if there is no such option.
func (r *Result) Probability(name string) float64 {
	for i, n := range r.Options {
		if n == name {
			return r.Probabilities[i]
		}
	}
	return 0
}

// Decider scores questions with a System One model. It is safe for concurrent
// use, but calls run one at a time.
type Decider struct {
	family     family
	b          backend
	maxOptions int
	ubatch     int
	nCtx       int
	nSeqMax    int
	manyMode   ManyMode
	bothOrders bool
	stateless  bool
	cached     []token
	mu         sync.Mutex
}

// family is how one kind of model renders questions and reads its logits.
type family interface {
	decide(d *Decider, state any, qs []Question, category string, many bool) ([]*Result, error)
}

// New loads a Jev-Style model and its readout_config.json.
// Load and init llama.cpp first, and call [Decider.Close] when done.
func New(modelPath, configPath string, opts Options) (*Decider, error) {
	if modelPath == "" || configPath == "" {
		return nil, errors.New("decide: model and readout config paths are required")
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return NewFromConfig(modelPath, cfg, opts)
}

// NewFromConfig loads a Jev-Style model with a readout config that is already
// parsed, for example with [ParseConfig] where there are no files.
func NewFromConfig(modelPath string, cfg *Config, opts Options) (*Decider, error) {
	if modelPath == "" || cfg == nil {
		return nil, errors.New("decide: model path and readout config are required")
	}

	d, err := load(modelPath, cfg.Budgets.MaxLen, opts)
	if err != nil {
		return nil, err
	}

	r, err := newRenderer(d.tokenizer(false), cfg)
	if err != nil {
		d.Close()
		return nil, err
	}
	r.maxLen = min(r.maxLen, d.nCtx)
	d.family = &jevStyle{cfg: cfg, render: r}

	return d, nil
}

// load loads the model and creates a context of opts.ContextSize tokens, or defCtx when it is 0.
func load(modelPath string, defCtx int, opts Options) (*Decider, error) {
	d, err := loadModel(modelPath, opts)
	if err != nil {
		return nil, err
	}
	if err := d.newContext(defCtx, opts, readoutLogits); err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}

// readout is what a model outputs and how its context is set up.
type readout int

const (
	// readoutLogits reads the logits of label tokens.
	readoutLogits readout = iota
	// readoutEmbd reads the embeddings output of each token, with no pooling.
	readoutEmbd
	// readoutEncoder is readoutEmbd for a model with no memory, so each
	// input is decoded alone in one ubatch.
	readoutEncoder
)

func loadModel(modelPath string, opts Options) (*Decider, error) {
	d := &Decider{maxOptions: 256, manyMode: opts.ManyMode, bothOrders: opts.BothOrders}
	if opts.MaxOptions > 0 {
		d.maxOptions = int(opts.MaxOptions)
	}
	if err := d.b.load(modelPath); err != nil {
		return nil, err
	}
	return d, nil
}

// newContext creates a context of opts.ContextSize tokens, or defCtx when it is 0.
func (d *Decider) newContext(defCtx int, opts Options, r readout) error {
	// One decode holds the whole input, so the batch is as big as the context.
	p := ctxParams{nCtx: uint32(defCtx), nSeqMax: maxSeqs, threads: opts.Threads, embeddings: r != readoutLogits}
	if opts.ContextSize > 0 {
		p.nCtx = opts.ContextSize
	}
	p.nBatch = p.nCtx
	p.nUbatch = min(1024, p.nCtx)
	if opts.UBatch > 0 {
		p.nUbatch = min(opts.UBatch, p.nCtx)
	}
	if r == readoutEncoder {
		p.nUbatch = p.nCtx
	}
	// llama.cpp needs room for at least one output per sequence.
	p.nOutputsMax = uint32(max(d.maxOptions, maxSeqs))

	n, err := d.b.newContext(p)
	if err != nil {
		return err
	}
	d.ubatch = int(p.nUbatch)
	d.nCtx = int(p.nCtx)
	d.nSeqMax = int(n)
	d.stateless = r == readoutEncoder
	// A context of one sequence cannot share the state, and an encoder keeps none.
	if d.nSeqMax < 2 || d.stateless {
		d.manyMode = manySeparate
	}

	return nil
}

func (d *Decider) tokenizer(parseSpecial bool) encoder {
	return func(s string) []token {
		return d.b.tokenize(s, parseSpecial)
	}
}

var specialText = regexp.MustCompile(`<\|([A-Za-z0-9_]+)\|>`)

// escapeSpecial changes <|name|> in user text to <¦name¦>, as llama-server does,
// so a prompt tokenized with special tokens parsed keeps it as text.
func escapeSpecial(s string) string {
	return specialText.ReplaceAllString(s, "<\u00a6$1\u00a6>")
}

// Close frees the model and context.
func (d *Decider) Close() {
	d.b.close()
}

// Config returns the readout config of a Jev-Style model, or nil for other models.
func (d *Decider) Config() *Config {
	if j, ok := d.family.(*jevStyle); ok {
		return j.cfg
	}
	return nil
}

// Decide scores one question about state. A string state is used as is, and
// any other value is encoded as JSON. Go sorts map keys, so pass a string or a
// json.RawMessage to keep a key order.
// category picks a Jev-Style calibration temperature. An empty category uses
// the global temperature. JevK5 and decider models ignore it.
func (d *Decider) Decide(state any, q Question, category string) (*Result, error) {
	rs, err := d.run(state, []Question{q}, category, false)
	if err != nil {
		return nil, err
	}
	return rs[0], nil
}

// DecideMany scores several questions about one state. The state is shared
// as set by [Options.ManyMode] and kept, so a later call with the same state
// reuses it. All questions are checked before any scoring.
func (d *Decider) DecideMany(state any, qs []Question, category string) ([]*Result, error) {
	if len(qs) == 0 {
		return nil, nil
	}
	return d.run(state, qs, category, true)
}

// run scores qs and, with [Options.BothOrders], reads each choice question
// a second time with its options reversed.
func (d *Decider) run(state any, qs []Question, category string, many bool) ([]*Result, error) {
	all := qs
	var flipped []int
	if d.bothOrders {
		for i, q := range qs {
			if q.Type != TypeChoice || len(q.Options) < 2 {
				continue
			}
			q.Options = slices.Clone(q.Options)
			slices.Reverse(q.Options)
			all = append(slices.Clip(all), q)
			flipped = append(flipped, i)
		}
	}

	rs, err := d.family.decide(d, state, all, category, many)
	if err != nil {
		return nil, err
	}
	for j, i := range flipped {
		rs[i] = bothOrders(rs[i], rs[len(qs)+j])
	}
	rs = rs[:len(qs)]
	for i, r := range rs {
		r.typed(qs[i].Type)
	}
	return rs, nil
}

// bothOrders averages the probabilities of a question read in two orders.
func bothOrders(a, b *Result) *Result {
	p := make([]float64, len(a.Options))
	for i, n := range a.Options {
		p[i] = (a.Probabilities[i] + b.Probability(n)) / 2
	}
	return newResult(a.Options, p, nil, a.Temperature, a.InputTokens+b.InputTokens, a.HeadTokens+b.HeadTokens)
}

// score decodes the rendered inputs, which share one state, and returns the
// logits at each slot for each row token. mode manySeparate decodes each separately.
func (d *Decider) score(rs []*rendered, mode ManyMode) ([][][]float64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.b.open() {
		return nil, errors.New("decide: decider is closed")
	}

	var (
		out [][][]float64
		err error
	)
	switch mode {
	case manySeparate:
		out, err = d.separate(rs)
	case ManyBatched:
		out, err = d.batched(rs)
	default:
		out, err = d.exact(rs)
	}
	if err != nil {
		d.clear()
		return nil, err
	}

	return out, nil
}

// separate decodes each input from an empty memory.
func (d *Decider) separate(rs []*rendered) ([][][]float64, error) {
	out := make([][][]float64, len(rs))
	for i, r := range rs {
		d.clear()
		v, err := d.decode(part{ids: r.ids, slots: r.slots, rows: r.rows, embd: r.embd})
		if err != nil {
			return nil, err
		}
		out[i] = v[0]
	}
	d.clear()
	return out, nil
}

// sharedLen is the number of leading tokens all inputs have in common, at
// most each input's prefixLen, so every input keeps its slots.
func sharedLen(rs []*rendered) int {
	n := rs[0].prefixLen
	for _, r := range rs[1:] {
		n = min(n, r.prefixLen)
		for i := range n {
			if r.ids[i] != rs[0].ids[i] {
				n = i
				break
			}
		}
	}
	return n
}

// exact decodes the full ubatches of the shared state once into sequence 0
// and each input on a copy of it in sequence 1. The ubatch splits are the same
// as in separate, so the results are identical.
func (d *Decider) exact(rs []*rendered) ([][][]float64, error) {
	shared := (sharedLen(rs) / d.ubatch) * d.ubatch
	if shared == 0 {
		return d.separate(rs)
	}

	if err := d.keepPrefix(rs[0].ids[:shared]); err != nil {
		return nil, err
	}

	out := make([][][]float64, len(rs))
	for i, r := range rs {
		d.b.memSeqRm(1)
		d.b.memSeqCp(0, 1)
		v, err := d.decode(part{ids: r.ids[shared:], pos0: shared, seq: 1, slots: r.slots, rows: r.rows, embd: r.embd})
		d.b.memSeqRm(1)
		if err != nil {
			return nil, err
		}
		out[i] = v[0]
	}

	return out, nil
}

// batched decodes the whole state once into sequence 0, then groups of
// questions in one decode, each on its own copy of the state.
func (d *Decider) batched(rs []*rendered) ([][][]float64, error) {
	n := sharedLen(rs)
	if err := d.keepPrefix(rs[0].ids[:n]); err != nil {
		return nil, err
	}

	out := make([][][]float64, 0, len(rs))
	for i := 0; i < len(rs); {
		var parts []part
		used, outs := 0, 0
		for i < len(rs) && len(parts) < d.nSeqMax-1 {
			l, ns := len(rs[i].ids)-n, len(rs[i].slots)
			if len(parts) > 0 && (n+used+l > d.nCtx || outs+ns > d.maxOptions) {
				break
			}
			seq := seqID(len(parts) + 1)
			parts = append(parts, part{ids: rs[i].ids[n:], pos0: n, seq: seq, slots: rs[i].slots, rows: rs[i].rows, embd: rs[i].embd})
			used += l
			outs += ns
			i++
		}

		for _, p := range parts {
			d.b.memSeqRm(p.seq)
			d.b.memSeqCp(0, p.seq)
		}
		v, err := d.decode(parts...)
		for _, p := range parts {
			d.b.memSeqRm(p.seq)
		}
		if err != nil {
			return nil, err
		}
		out = append(out, v...)
	}

	return out, nil
}

// keepPrefix makes sequence 0 hold exactly prefix, reusing it when it is already there.
func (d *Decider) keepPrefix(prefix []token) error {
	if slices.Equal(d.cached, prefix) {
		return nil
	}

	d.clear()
	if _, err := d.decode(part{ids: prefix}); err != nil {
		return err
	}
	d.cached = slices.Clone(prefix)

	return nil
}

// clear empties the memory and forgets the cached state.
func (d *Decider) clear() {
	if !d.stateless {
		d.b.memClear()
	}
	d.cached = nil
}

// part is a run of tokens at positions pos0 onward in seq. slots are absolute
// positions, and rows are the tokens whose logits are read at each slot.
// With embd, the whole embeddings output is read at each slot instead.
type part struct {
	ids   []token
	pos0  int
	seq   seqID
	slots []int
	rows  []token
	embd  bool
}

// decode runs all parts in one batch and returns, per part, the logits of its
// rows or the embeddings at each of its slots.
func (d *Decider) decode(parts ...part) ([][][]float64, error) {
	total := 0
	for _, p := range parts {
		total += len(p.ids)
	}

	bt := newBatch(total)
	defer bt.free()

	idx := make([][]int32, len(parts))
	for k, p := range parts {
		isSlot := make(map[int]bool, len(p.slots))
		for _, s := range p.slots {
			isSlot[s] = true
		}
		for i, tok := range p.ids {
			if isSlot[p.pos0+i] {
				idx[k] = append(idx[k], bt.len())
			}
			if err := bt.add(tok, p.pos0+i, p.seq, isSlot[p.pos0+i]); err != nil {
				return nil, err
			}
		}
	}

	if err := d.b.decode(bt); err != nil {
		return nil, err
	}

	out := make([][][]float64, len(parts))
	for k, p := range parts {
		out[k] = make([][]float64, len(idx[k]))
		for j, i := range idx[k] {
			if p.embd {
				e, err := d.b.embeddings(i)
				if err != nil {
					return nil, err
				}
				if e == nil {
					return nil, fmt.Errorf("decide: no embeddings at batch index %d", i)
				}
				out[k][j] = make([]float64, len(e))
				for r, v := range e {
					out[k][j][r] = float64(v)
				}
				continue
			}
			logits, err := d.b.logits(i)
			if err != nil {
				return nil, err
			}
			if logits == nil {
				return nil, fmt.Errorf("decide: no logits at batch index %d", i)
			}
			out[k][j] = make([]float64, len(p.rows))
			for r, t := range p.rows {
				out[k][j][r] = float64(logits[t])
			}
		}
	}

	return out, nil
}

// softmax returns softmax(scores / temp).
func softmax(scores []float64, temp float64) ([]float64, error) {
	z := make([]float64, len(scores))
	zmax := math.Inf(-1)
	for i, s := range scores {
		z[i] = s / temp
		if math.IsNaN(z[i]) || math.IsInf(z[i], 0) {
			return nil, errors.New("decide: non-finite decision scores")
		}
		zmax = max(zmax, z[i])
	}

	var sum float64
	p := make([]float64, len(z))
	for i := range z {
		p[i] = math.Exp(z[i] - zmax)
		sum += p[i]
	}
	for i := range p {
		p[i] /= sum
	}
	return p, nil
}

// newResult builds a Result from probabilities in names order.
func newResult(names []string, p, scores []float64, temp float64, inputTokens, headTokens int) *Result {
	top := 0
	for i := range p {
		if p[i] > p[top] {
			top = i
		}
	}

	return &Result{
		Answer:               names[top],
		Options:              names,
		Probabilities:        p,
		Scores:               scores,
		Temperature:          temp,
		TopProbability:       p[top],
		EntropyConcentration: concentration(p),
		InputTokens:          inputTokens,
		HeadTokens:           headTokens,
	}
}

// typed sets the fields that depend on the question type.
// Score questions use the score confidence, the others the choice confidence.
func (r *Result) typed(t Type) {
	if t != TypeScore {
		r.Confidence = confidenceChoice(r.Probabilities)
		return
	}
	r.Confidence = confidenceScore(r.Probabilities)
	for i, v := range r.Probabilities {
		r.Expected += float64(i) * v
	}
}

// confidenceChoice is 0 when all options are equally likely and 1 when one is certain.
func confidenceChoice(p []float64) float64 {
	if len(p) < 2 {
		return 1
	}
	u := 1 / float64(len(p))
	return math.Max(0, (slices.Max(p)-u)/(1-u))
}

// confidenceScore compares the mean distance to the mode with the one of a uniform distribution.
func confidenceScore(p []float64) float64 {
	n := len(p)
	if n < 2 {
		return 1
	}
	mode := 0
	for i := range p {
		if p[i] > p[mode] {
			mode = i
		}
	}
	var dist, uniform float64
	for i, v := range p {
		dist += v * math.Abs(float64(i-mode))
		uniform += math.Abs(float64(i)-float64(n-1)/2) / float64(n)
	}
	return math.Max(0, 1-dist/uniform)
}

// concentration is 1 minus the entropy of p over its maximum, in [0, 1].
func concentration(p []float64) float64 {
	if len(p) < 2 {
		return 1
	}
	var ent float64
	for _, v := range p {
		ent -= v * math.Log(math.Min(1, math.Max(1e-12, v)))
	}
	return math.Min(1, math.Max(0, 1-ent/math.Log(float64(len(p)))))
}
