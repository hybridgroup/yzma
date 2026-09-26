package decide

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"sync"

	"github.com/hybridgroup/yzma/pkg/llama"
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
	// max_len for Jev-Style and 8192 for JevK5.
	ContextSize uint32
}

// ManyMode is how [Decider.DecideMany] shares one state across questions.
type ManyMode int

const (
	// ManyExact decodes the whole ubatches of the state once and each question
	// on its own. Results are identical to [Decider.Decide].
	ManyExact ManyMode = iota
	// ManyBatched decodes the whole state once and up to 16 questions in one
	// decode. It is faster, but probabilities can differ from [Decider.Decide]
	// by a few hundredths and a near tie can change the answer. For a JevK5
	// question with more than 16 options a near tie can also change which
	// options reach the final pass, which moves its probabilities more.
	ManyBatched

	// manySeparate decodes each question on its own from an empty memory.
	manySeparate ManyMode = -1
)

// maxSeqs is the number of sequences, the shared state plus up to 16 questions.
const maxSeqs = 17

// Result is the decision for one question.
// Options, Probabilities and Scores are in [Question.Names] order.
// Scores are the raw scores before calibration. They are empty for a JevK5
// question with more than 16 options, which combines several passes.
type Result struct {
	Answer               string    `json:"answer"`
	Options              []string  `json:"options"`
	Probabilities        []float64 `json:"probabilities"`
	Scores               []float64 `json:"scores,omitempty"`
	Temperature          float64   `json:"temperature"`
	TopProbability       float64   `json:"top_probability"`
	EntropyConcentration float64   `json:"entropy_concentration"`
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
	model      llama.Model
	vocab      llama.Vocab
	ctx        llama.Context
	mem        llama.Memory
	nVocab     int
	maxOptions int
	ubatch     int
	nCtx       int
	nSeqMax    int
	manyMode   ManyMode
	cached     []llama.Token
	mu         sync.Mutex
}

// family is how one kind of model renders questions and reads its logits.
type family interface {
	decide(d *Decider, state any, qs []Question, category string, many bool) ([]*Result, error)
}

// New loads a Jev-Style model and its readout_config.json.
// Call [llama.Load] and [llama.Init] first, and [Decider.Close] when done.
func New(modelPath, configPath string, opts Options) (*Decider, error) {
	if modelPath == "" || configPath == "" {
		return nil, errors.New("decide: model and readout config paths are required")
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
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
	model, err := llama.ModelLoadFromFile(modelPath, llama.ModelDefaultParams())
	if err != nil {
		return nil, fmt.Errorf("decide: load model: %w", err)
	}
	if model == 0 {
		return nil, fmt.Errorf("decide: unable to load model %s", modelPath)
	}

	d := &Decider{model: model, maxOptions: 256, manyMode: opts.ManyMode}
	if opts.MaxOptions > 0 {
		d.maxOptions = int(opts.MaxOptions)
	}
	d.vocab = llama.ModelGetVocab(model)
	d.nVocab = int(llama.VocabNTokens(d.vocab))

	// One decode holds the whole input, so the batch is as big as the context.
	params := llama.ContextDefaultParams()
	params.NCtx = uint32(defCtx)
	if opts.ContextSize > 0 {
		params.NCtx = opts.ContextSize
	}
	params.NBatch = params.NCtx
	params.NUbatch = min(1024, params.NCtx)
	if opts.UBatch > 0 {
		params.NUbatch = min(opts.UBatch, params.NCtx)
	}
	params.NSeqMax = maxSeqs
	// llama.cpp needs room for at least one output per sequence.
	params.NOutputsMax = uint32(max(d.maxOptions, maxSeqs))
	params.KVUnified = 1
	params.NoPerf = 1
	if opts.Threads > 0 {
		params.NThreads = opts.Threads
		params.NThreadsBatch = opts.Threads
	}

	d.ubatch = int(params.NUbatch)
	d.nCtx = int(params.NCtx)
	d.nSeqMax = maxSeqs

	d.ctx, err = llama.InitFromModel(model, params)
	if err != nil || d.ctx == 0 {
		d.Close()
		return nil, fmt.Errorf("decide: unable to create context: %v", err)
	}
	if d.mem, err = llama.GetMemory(d.ctx); err != nil {
		d.Close()
		return nil, err
	}

	return d, nil
}

func (d *Decider) tokenizer(parseSpecial bool) encoder {
	return func(s string) []llama.Token {
		return llama.Tokenize(d.vocab, s, false, parseSpecial)
	}
}

// Close frees the model and context.
func (d *Decider) Close() {
	if d.ctx != 0 {
		llama.Free(d.ctx)
		d.ctx = 0
	}
	if d.model != 0 {
		llama.ModelFree(d.model)
		d.model = 0
	}
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
// the global temperature. JevK5 ignores it.
func (d *Decider) Decide(state any, q Question, category string) (*Result, error) {
	rs, err := d.family.decide(d, state, []Question{q}, category, false)
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
	return d.family.decide(d, state, qs, category, true)
}

// score decodes the rendered inputs, which share one state, and returns the
// logits at each slot for each row token. mode manySeparate decodes each on its own.
func (d *Decider) score(rs []*rendered, mode ManyMode) ([][][]float64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.ctx == 0 {
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
		v, err := d.decode(part{ids: r.ids, slots: r.slots, rows: r.rows})
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

// exact decodes the whole ubatches of the shared state once into sequence 0
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
		llama.MemorySeqRm(d.mem, 1, -1, -1)
		llama.MemorySeqCp(d.mem, 0, 1, -1, -1)
		v, err := d.decode(part{ids: r.ids[shared:], pos0: shared, seq: 1, slots: r.slots, rows: r.rows})
		llama.MemorySeqRm(d.mem, 1, -1, -1)
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
			seq := llama.SeqId(len(parts) + 1)
			parts = append(parts, part{ids: rs[i].ids[n:], pos0: n, seq: seq, slots: rs[i].slots, rows: rs[i].rows})
			used += l
			outs += ns
			i++
		}

		for _, p := range parts {
			llama.MemorySeqRm(d.mem, p.seq, -1, -1)
			llama.MemorySeqCp(d.mem, 0, p.seq, -1, -1)
		}
		v, err := d.decode(parts...)
		for _, p := range parts {
			llama.MemorySeqRm(d.mem, p.seq, -1, -1)
		}
		if err != nil {
			return nil, err
		}
		out = append(out, v...)
	}

	return out, nil
}

// keepPrefix makes sequence 0 hold exactly prefix, reusing it when it is already there.
func (d *Decider) keepPrefix(prefix []llama.Token) error {
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
	llama.MemoryClear(d.mem, true)
	d.cached = nil
}

// part is a run of tokens at positions pos0 onward in seq. slots are absolute
// positions, and rows are the tokens whose logits are read at each slot.
type part struct {
	ids   []llama.Token
	pos0  int
	seq   llama.SeqId
	slots []int
	rows  []llama.Token
}

// decode runs all parts in one batch and returns, per part, the logits of its
// rows at each of its slots.
func (d *Decider) decode(parts ...part) ([][][]float64, error) {
	total := 0
	for _, p := range parts {
		total += len(p.ids)
	}

	batch := llama.BatchInit(int32(total), 0, 1)
	defer llama.BatchFree(batch)

	idx := make([][]int32, len(parts))
	for k, p := range parts {
		isSlot := make(map[int]bool, len(p.slots))
		for _, s := range p.slots {
			isSlot[s] = true
		}
		seqs := []llama.SeqId{p.seq}
		for i, tok := range p.ids {
			if isSlot[p.pos0+i] {
				idx[k] = append(idx[k], batch.NTokens)
			}
			if err := batch.Add(tok, llama.Pos(p.pos0+i), seqs, isSlot[p.pos0+i]); err != nil {
				return nil, err
			}
		}
	}

	ret, err := llama.Decode(d.ctx, batch)
	if err != nil {
		return nil, err
	}
	if ret != 0 {
		return nil, fmt.Errorf("decide: decode returned %d", ret)
	}

	out := make([][][]float64, len(parts))
	for k, p := range parts {
		out[k] = make([][]float64, len(idx[k]))
		for j, i := range idx[k] {
			logits, err := llama.GetLogitsIth(d.ctx, i, d.nVocab)
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
