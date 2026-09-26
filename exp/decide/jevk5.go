package decide

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"sort"
	"strconv"

	"github.com/hybridgroup/yzma/pkg/llama"
)

const (
	jevK5Letters = "ABCDEFGHIJKLMNOP"
	jevK5System  = "Apply the supplied criterion to the supplied evidence. Choose exactly one listed option. " +
		"Respond with only its uppercase letter, with no explanation or reasoning."
	jevK5Context = 8192
)

// JevK5Config is the jevk5_config.json of a JevK5 model.
type JevK5Config struct {
	// Temperature calibrates the letter probabilities of every pass.
	Temperature float64 `json:"temperature"`
	// KnockoutTemperature sharpens the combined result of a question with
	// more than 16 options. 0 uses 0.77.
	KnockoutTemperature float64 `json:"knockout_temperature"`
}

// LoadJevK5Config reads and checks a jevk5_config.json file.
func LoadJevK5Config(path string) (*JevK5Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c JevK5Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("jevk5 config: %w", err)
	}
	if c.Temperature <= 0 {
		return nil, errors.New("jevk5 config has no temperature")
	}
	if c.KnockoutTemperature < 0 {
		return nil, fmt.Errorf("jevk5 config has bad knockout_temperature %v", c.KnockoutTemperature)
	}
	if c.KnockoutTemperature == 0 {
		c.KnockoutTemperature = 0.77
	}

	return &c, nil
}

// NewJevK5 loads a JevK5 model and its jevk5_config.json. JevK5 reads the
// probability of each option from the logit of its letter after a chat prompt.
// Call [llama.Load] and [llama.Init] first, and [Decider.Close] when done.
//
// As in the reference runtime, the whole prompt is tokenized with special
// tokens parsed, so control token text in the state is read as a control token.
func NewJevK5(modelPath, configPath string, opts Options) (*Decider, error) {
	if modelPath == "" || configPath == "" {
		return nil, errors.New("decide: model and jevk5 config paths are required")
	}

	cfg, err := LoadJevK5Config(configPath)
	if err != nil {
		return nil, err
	}

	d, err := load(modelPath, jevK5Context, opts)
	if err != nil {
		return nil, err
	}

	k := &jevK5{cfg: cfg, enc: d.tokenizer(true)}
	plain := d.tokenizer(false)
	for _, l := range jevK5Letters {
		ids := plain(string(l))
		if len(ids) != 1 {
			d.Close()
			return nil, fmt.Errorf("tokenizer mismatch: letter %q gives %v, want one token", l, ids)
		}
		k.letters = append(k.letters, ids[0])
	}
	d.family = k

	return d, nil
}

type jevK5 struct {
	cfg     *JevK5Config
	enc     encoder
	letters []llama.Token
}

// jevK5Option is one option as the model sees it, and its place in [Question.Names].
type jevK5Option struct {
	text  string
	index int
}

// options returns the options in the order the model was trained on, with
// "name: description" texts. A noul question lists true before false.
func (k *jevK5) options(q Question) ([]jevK5Option, []string, error) {
	names, err := q.Names()
	if err != nil {
		return nil, nil, err
	}

	var out []jevK5Option
	switch q.Type {
	case TypeNoul:
		desc := map[string]string{}
		for _, o := range q.Options {
			desc[o.Name] = o.Description
		}
		for _, name := range []string{"true", "false"} {
			d := desc[name]
			if d == "" {
				d = "The proposition is " + name + "."
			}
			out = append(out, jevK5Option{text: name + ": " + d, index: slices.Index(names, name)})
		}
	case TypeChoice:
		for i, o := range q.Options {
			d := o.Description
			if d == "" {
				d = o.Name
			}
			out = append(out, jevK5Option{text: o.Name + ": " + d, index: i})
		}
	case TypeScore:
		for i, o := range q.Options {
			out = append(out, jevK5Option{text: strconv.Itoa(i) + ": " + o.Description, index: i})
		}
	}

	return out, names, nil
}

type jevK5Letter struct {
	Letter      string `json:"letter"`
	Description string `json:"description"`
}

type jevK5Payload struct {
	Evidence  any           `json:"evidence"`
	Criterion string        `json:"criterion"`
	Options   []jevK5Letter `json:"options"`
}

// prompt returns the chat prompt for one pass over at most 16 option texts.
func jevK5Prompt(state any, criterion string, texts []string) (string, error) {
	p := jevK5Payload{Evidence: state, Criterion: criterion, Options: make([]jevK5Letter, len(texts))}
	for i, t := range texts {
		p.Options[i] = jevK5Letter{Letter: jevK5Letters[i : i+1], Description: t}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(p); err != nil {
		return "", fmt.Errorf("state: %w", err)
	}
	user := spaceJSON(bytes.TrimRight(buf.Bytes(), "\n"))

	return "<|im_start|>system\n" + jevK5System + "<|im_end|>\n" +
		"<|im_start|>user\n" + user + "<|im_end|>\n" +
		"<|im_start|>assistant\n<think>\n\n</think>\n\n", nil
}

// stateLen is the number of leading tokens that prompts about state share no
// matter the criterion. Passes share at most these, so the kept state is the
// same for every pass.
func (k *jevK5) stateLen(state any) (int, error) {
	a, err := jevK5Prompt(state, "a", []string{"a"})
	if err != nil {
		return 0, err
	}
	b, err := jevK5Prompt(state, "b", []string{"b"})
	if err != nil {
		return 0, err
	}

	ia, ib := k.enc(a), k.enc(b)
	n := 0
	for n < len(ia) && n < len(ib) && ia[n] == ib[n] {
		n++
	}
	return n, nil
}

// render tokenizes one pass. The letters are read at the last token.
func (k *jevK5) render(d *Decider, state any, criterion string, texts []string, stateLen int) (*rendered, error) {
	prompt, err := jevK5Prompt(state, criterion, texts)
	if err != nil {
		return nil, err
	}

	ids := k.enc(prompt)
	if len(ids) > d.nCtx {
		return nil, fmt.Errorf("%w: input needs %d tokens, the context is %d", ErrBudget, len(ids), d.nCtx)
	}

	last := len(ids) - 1
	return &rendered{ids: ids, slots: []int{last}, rows: k.letters[:len(texts)], prefixLen: min(last, stateLen)}, nil
}

// jevK5Pass is one pass, a criterion and at most 16 option texts.
type jevK5Pass struct {
	criterion string
	texts     []string
}

// read scores passes about one state and returns calibrated letter probabilities,
// the raw letter logits and the input tokens of each pass.
func (k *jevK5) read(d *Decider, state any, stateLen int, passes []jevK5Pass, mode ManyMode) ([][]float64, [][]float64, []int, error) {
	rs := make([]*rendered, len(passes))
	for i, p := range passes {
		var err error
		if rs[i], err = k.render(d, state, p.criterion, p.texts, stateLen); err != nil {
			return nil, nil, nil, err
		}
	}

	vals, err := d.score(rs, mode)
	if err != nil {
		return nil, nil, nil, err
	}

	probs := make([][]float64, len(rs))
	logits := make([][]float64, len(rs))
	tokens := make([]int, len(rs))
	for i := range rs {
		logits[i] = vals[i][0]
		if probs[i], err = softmax(logits[i], k.cfg.Temperature); err != nil {
			return nil, nil, nil, err
		}
		tokens[i] = len(rs[i].ids)
	}

	return probs, logits, tokens, nil
}

func (k *jevK5) decide(d *Decider, state any, qs []Question, _ string, many bool) ([]*Result, error) {
	stateLen, err := k.stateLen(state)
	if err != nil {
		return nil, err
	}

	opts := make([][]jevK5Option, len(qs))
	names := make([][]string, len(qs))
	texts := make([][]string, len(qs))
	for i, q := range qs {
		var err error
		if opts[i], names[i], err = k.options(q); err == nil && len(opts[i]) > d.maxOptions {
			err = fmt.Errorf("%w: %d options, the limit is %d", ErrQuestion, len(opts[i]), d.maxOptions)
		}
		if err == nil {
			for _, o := range opts[i] {
				texts[i] = append(texts[i], o.text)
			}
			_, err = k.render(d, state, q.Text, texts[i][:min(len(texts[i]), len(jevK5Letters))], stateLen)
		}
		if err != nil {
			return nil, questionErr(many, i, err)
		}
	}

	mode := manySeparate
	if many {
		mode = d.manyMode
	}

	// Questions of up to 16 options take one pass each and are scored together.
	var single []int
	var passes []jevK5Pass
	for i, q := range qs {
		if len(texts[i]) <= len(jevK5Letters) {
			single = append(single, i)
			passes = append(passes, jevK5Pass{q.Text, texts[i]})
		}
	}

	out := make([]*Result, len(qs))
	if len(passes) > 0 {
		probs, logits, tokens, err := k.read(d, state, stateLen, passes, mode)
		if err != nil {
			return nil, err
		}
		for n, i := range single {
			out[i] = k.result(names[i], opts[i], probs[n], logits[n], tokens[n])
		}
	}

	// Each other question takes its own knockout passes, which share the state
	// as in ManyExact unless ManyBatched was asked for.
	for i, q := range qs {
		if out[i] != nil {
			continue
		}
		m := mode
		if m == manySeparate {
			m = ManyExact
		}
		inputTokens := 0
		read := func(sets [][]string) ([][]float64, error) {
			ps := make([]jevK5Pass, len(sets))
			for j, t := range sets {
				ps[j] = jevK5Pass{q.Text, t}
			}
			probs, _, tokens, err := k.read(d, state, stateLen, ps, m)
			for _, t := range tokens {
				inputTokens += t
			}
			return probs, err
		}
		p, err := knockout(read, texts[i])
		if err != nil {
			return nil, questionErr(many, i, err)
		}
		out[i] = k.result(names[i], opts[i], sharpen(p, k.cfg.KnockoutTemperature), nil, inputTokens)
	}

	return out, nil
}

// result puts probabilities in model order into a Result in names order.
func (k *jevK5) result(names []string, opts []jevK5Option, p, logits []float64, tokens int) *Result {
	np := make([]float64, len(names))
	var ns []float64
	if logits != nil {
		ns = make([]float64, len(names))
	}
	for j, o := range opts {
		np[o.index] = p[j]
		if logits != nil {
			ns[o.index] = logits[j]
		}
	}
	return newResult(names, np, ns, k.cfg.Temperature, tokens, 0)
}

// passReader reads passes of at most 16 option texts and returns their probabilities.
type passReader func(passes [][]string) ([][]float64, error)

// knockout gives a probability to each of more than 16 options, as the JevK5
// runtime does. The options are split in order into near equal groups of at
// most 16 and each group is read. A final pass then reads the top options of
// every group. Finalists keep their final share, times the chance the answer
// is a finalist, and every other option gets its group's share of the final
// times its probability in the group.
func knockout(read passReader, texts []string) ([]float64, error) {
	if len(texts) <= len(jevK5Letters) {
		p, err := read([][]string{texts})
		if err != nil {
			return nil, err
		}
		return p[0], nil
	}

	runs := groups(len(texts), (len(texts)+len(jevK5Letters)-1)/len(jevK5Letters))
	passes := make([][]string, len(runs))
	for g, run := range runs {
		passes[g] = texts[run[0]:run[1]]
	}
	inner, err := read(passes)
	if err != nil {
		return nil, err
	}
	for _, p := range inner {
		normalize(p)
	}

	keep := max(1, len(jevK5Letters)/len(runs))
	type pick struct{ g, j int }
	chosen := map[pick]bool{}
	var rest []pick
	for g, p := range inner {
		order := make([]int, len(p))
		for j := range order {
			order[j] = j
		}
		// Ties go to the earlier option.
		sort.SliceStable(order, func(a, b int) bool { return p[order[a]] > p[order[b]] })
		for r, j := range order {
			if r < keep {
				chosen[pick{g, j}] = true
			} else {
				rest = append(rest, pick{g, j})
			}
		}
	}
	sort.SliceStable(rest, func(a, b int) bool { return inner[rest[a].g][rest[a].j] > inner[rest[b].g][rest[b].j] })
	for _, pk := range rest[:max(0, min(len(rest), len(jevK5Letters)-len(chosen)))] {
		chosen[pk] = true
	}

	tops := make([][]int, len(runs))
	var finalTexts []string
	for g, run := range runs {
		for j := 0; j < run[1]-run[0]; j++ {
			if chosen[pick{g, j}] {
				tops[g] = append(tops[g], j)
				finalTexts = append(finalTexts, texts[run[0]+j])
			}
		}
	}

	final, err := knockout(read, finalTexts)
	if err != nil {
		return nil, err
	}

	// Sums run in the reference's order, so the floats match.
	shares := make([]map[int]float64, len(runs))
	masses := make([]float64, len(runs))
	at := 0
	for g, top := range tops {
		shares[g] = map[int]float64{}
		for _, j := range top {
			shares[g][j] = final[at]
			masses[g] += final[at]
			at++
		}
	}

	var inFinal float64
	for g, p := range inner {
		var in float64
		for _, j := range tops[g] {
			in += p[j]
		}
		inFinal += masses[g] * in
	}

	var weights []float64
	for g, p := range inner {
		for j, q := range p {
			if f, ok := shares[g][j]; ok {
				weights = append(weights, f*inFinal)
			} else {
				weights = append(weights, masses[g]*q)
			}
		}
	}
	normalize(weights)

	return weights, nil
}

// groups splits n items into count contiguous [start, stop) runs whose sizes differ by at most one.
func groups(n, count int) [][2]int {
	base, extra := n/count, n%count
	runs := make([][2]int, count)
	start := 0
	for g := range runs {
		stop := start + base
		if g < extra {
			stop++
		}
		runs[g] = [2]int{start, stop}
		start = stop
	}
	return runs
}

// sharpen raises p to 1/temp and normalizes it.
func sharpen(p []float64, temp float64) []float64 {
	if temp == 1 {
		return p
	}
	out := make([]float64, len(p))
	for i, q := range p {
		out[i] = math.Pow(q, 1/temp)
	}
	normalize(out)
	return out
}

func normalize(p []float64) {
	var sum float64
	for _, q := range p {
		sum += q
	}
	for i := range p {
		p[i] /= sum
	}
}
