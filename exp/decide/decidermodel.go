package decide

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	// deciderNarrow is how many options are written with the letters A to J.
	deciderNarrow     = 10
	deciderMaxOptions = 255
	deciderContext    = 8192
	deciderNeutral    = "not listed here"
)

// DeciderConfig is the decider_config.json of a decider model.
type DeciderConfig struct {
	// Temperature calibrates every answer. It is 1 when the file has none.
	Temperature float64 `json:"temperature"`
	// TemperatureByType replaces Temperature for the question types it names.
	TemperatureByType map[Type]float64 `json:"temperature_by_type"`
	// IsolatedLevels judges each level of a score question in its own yes or no question.
	IsolatedLevels bool `json:"isolated_levels"`
	// NeutralizeNone rewrites "none of the above" options. It is true when the file has none.
	NeutralizeNone bool `json:"neutralize_none"`
	// Layout is the prompt layout. Only the plain layout is supported.
	Layout       string `json:"layout"`
	ChatTemplate bool   `json:"chat_template"`
}

// LoadDeciderConfig reads and checks a decider_config.json file.
func LoadDeciderConfig(path string) (*DeciderConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseDeciderConfig(data)
}

// ParseDeciderConfig parses and checks the contents of a decider_config.json file.
func ParseDeciderConfig(data []byte) (*DeciderConfig, error) {
	c := DeciderConfig{Temperature: 1, NeutralizeNone: true}
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("decider config: %w", err)
	}

	if !positive(c.Temperature) {
		return nil, fmt.Errorf("decider config has bad temperature %v", c.Temperature)
	}
	for t, v := range c.TemperatureByType {
		if t != TypeChoice && t != TypeScore && t != TypeNoul {
			return nil, fmt.Errorf("decider config has temperature for unknown type %q", t)
		}
		if !positive(v) {
			return nil, fmt.Errorf("decider config has bad %s temperature %v", t, v)
		}
	}
	layout := c.Layout
	if layout == "" && c.ChatTemplate {
		layout = "chat"
	}
	if layout != "" && layout != "plain" {
		return nil, fmt.Errorf("decider config uses layout %q, only plain is supported", layout)
	}

	return &c, nil
}

func positive(v float64) bool {
	return v > 0 && !math.IsInf(v, 0)
}

// TypeTemperature returns the calibration temperature of a question type.
func (c *DeciderConfig) TypeTemperature(t Type) float64 {
	if v, ok := c.TemperatureByType[t]; ok {
		return v
	}
	return c.Temperature
}

// NewDeciderModel loads a decider model and its decider_config.json.
// Load and init llama.cpp first, and call [Decider.Close] when done.
func NewDeciderModel(modelPath, configPath string, opts Options) (*Decider, error) {
	if modelPath == "" || configPath == "" {
		return nil, errors.New("decide: model and decider config paths are required")
	}

	cfg, err := LoadDeciderConfig(configPath)
	if err != nil {
		return nil, err
	}
	return NewDeciderModelFromConfig(modelPath, cfg, opts)
}

// NewDeciderModelFromConfig loads a decider model with a config that is
// already parsed, for example with [ParseDeciderConfig] where there are no files.
func NewDeciderModelFromConfig(modelPath string, cfg *DeciderConfig, opts Options) (*Decider, error) {
	if modelPath == "" || cfg == nil {
		return nil, errors.New("decide: model path and decider config are required")
	}

	d, err := load(modelPath, deciderContext, opts)
	if err != nil {
		return nil, err
	}

	m, err := newDeciderModel(d.tokenizer(true), d.tokenizer(false), cfg)
	if err != nil {
		d.Close()
		return nil, err
	}
	d.family = m

	return d, nil
}

type deciderModel struct {
	cfg    *DeciderConfig
	enc    encoder
	labels []token
	open   []token
}

// newDeciderModel finds the label tokens, A to Z and then the two letter
// strings that are one token, as the reference runtime does.
func newDeciderModel(enc, plain encoder, cfg *DeciderConfig) (*deciderModel, error) {
	m := &deciderModel{cfg: cfg, enc: enc, open: enc("\n(")}

	var names []string
	for a := 'A'; a <= 'Z'; a++ {
		names = append(names, string(a))
	}
	for a := 'A'; a <= 'Z'; a++ {
		for b := 'A'; b <= 'Z'; b++ {
			names = append(names, string([]rune{a, b}))
		}
	}
	for i, n := range names {
		ids := plain(n)
		if len(ids) == 1 {
			m.labels = append(m.labels, ids[0])
		} else if i < deciderNarrow {
			return nil, fmt.Errorf("tokenizer mismatch: letter %q gives %v, want one token", n, ids)
		}
		if len(m.labels) == deciderMaxOptions {
			return m, nil
		}
	}

	return nil, fmt.Errorf("tokenizer mismatch: %d label tokens, want %d", len(m.labels), deciderMaxOptions)
}

// deciderRow is one question as the model reads it.
type deciderRow struct {
	text    string
	options []string
}

var levelNumber = regexp.MustCompile(`^[\t\n\v\f\r\x1c-\x1f\x85\p{Z}]*-?\p{Nd}+[\t\n\v\f\r\x1c-\x1f\x85\p{Z}]*:[\t\n\v\f\r\x1c-\x1f\x85\p{Z}]*`)

// rows returns the option names of q and its rows. An isolated score question
// has one yes or no row per level, and every other question has one row.
func (m *deciderModel) rows(q Question, limit int) ([]string, []deciderRow, error) {
	names, err := q.Names()
	if err != nil {
		return nil, nil, err
	}

	var opts []string
	switch q.Type {
	case TypeChoice:
		if len(q.Options) < 2 {
			return nil, nil, fmt.Errorf("%w: choice needs at least two options", ErrQuestion)
		}
		for _, o := range q.Options {
			if o.Description == "" {
				opts = append(opts, o.Name)
			} else {
				opts = append(opts, o.Name+": "+o.Description)
			}
		}
	case TypeScore:
		if m.cfg.IsolatedLevels {
			rows := make([]deciderRow, len(q.Options))
			for i, o := range q.Options {
				level := levelNumber.ReplaceAllString(o.Description, "")
				rows[i] = deciderRow{text: q.Text + "\nProposed answer: " + level + "\nDoes the proposed answer fit?", options: []string{"no", "yes"}}
			}
			return names, rows, nil
		}
		for i, o := range q.Options {
			opts = append(opts, strconv.Itoa(i)+": "+o.Description)
		}
	case TypeNoul:
		opts = []string{"no", "yes"}
		for _, o := range q.Options {
			if o.Description == "" {
				continue
			}
			if o.Name == "false" {
				opts[0] = "no: " + o.Description
			} else {
				opts[1] = "yes: " + o.Description
			}
		}
	}

	if len(opts) > limit {
		return nil, nil, fmt.Errorf("%w: %d options, the limit is %d", ErrQuestion, len(opts), limit)
	}
	if m.cfg.NeutralizeNone {
		for i, o := range opts {
			if isNoneOption(o) {
				opts[i] = deciderNeutral
			}
		}
	}

	return names, []deciderRow{{text: q.Text, options: opts}}, nil
}

// isNoneOption reports whether an option reads like "none of the above",
// which older models learned as a signal to abstain.
func isNoneOption(o string) bool {
	k := strings.ToLower(strings.TrimSpace(o))
	return strings.HasPrefix(k, "none of the above") || k == "none" || k == "n/a" || k == "none of these"
}

// deciderPrompt returns the text after the state of a row of at most 10 options.
func deciderPrompt(row deciderRow) string {
	var b strings.Builder
	b.WriteString("\n\nQuestion: " + row.text + "\nOptions:")
	for j, o := range row.options {
		b.WriteString("\n(" + string(rune('A'+j)) + ") " + o)
	}
	b.WriteString("\nAnswer: (")
	return b.String()
}

// render tokenizes one row after the state tokens ctx. The labels are read at the last token.
func (m *deciderModel) render(d *Decider, ctx []token, row deciderRow) (*rendered, error) {
	ids := slices.Clip(ctx)
	if len(row.options) <= deciderNarrow {
		ids = append(ids, m.enc(deciderPrompt(row))...)
	} else {
		ids = append(ids, m.enc("\n\nQuestion: "+row.text+"\nOptions:")...)
		for j, o := range row.options {
			ids = append(ids, m.open...)
			ids = append(ids, m.labels[j])
			ids = append(ids, m.enc(") "+o)...)
		}
		ids = append(ids, m.enc("\nAnswer: (")...)
	}
	if len(ids) > d.nCtx {
		return nil, fmt.Errorf("%w: input needs %d tokens, the context is %d", ErrBudget, len(ids), d.nCtx)
	}

	return &rendered{ids: ids, slots: []int{len(ids) - 1}, rows: m.labels[:len(row.options)], prefixLen: len(ctx)}, nil
}

func (m *deciderModel) decide(d *Decider, state any, qs []Question, _ string, many bool) ([]*Result, error) {
	s, err := deciderState(state)
	if err != nil {
		return nil, err
	}
	ctx := m.enc("Context:\n" + escapeSpecial(s))

	type plan struct {
		names       []string
		first, rows int
	}
	plans := make([]plan, len(qs))
	var rs []*rendered
	for i, q := range qs {
		names, rows, err := m.rows(q, min(d.maxOptions, deciderMaxOptions))
		if err != nil {
			return nil, questionErr(many, i, err)
		}
		for j := range rows {
			rows[j].text = escapeSpecial(rows[j].text)
			for k, o := range rows[j].options {
				rows[j].options[k] = escapeSpecial(o)
			}
		}
		plans[i] = plan{names: names, first: len(rs), rows: len(rows)}
		for _, row := range rows {
			r, err := m.render(d, ctx, row)
			if err != nil {
				return nil, questionErr(many, i, err)
			}
			rs = append(rs, r)
		}
	}

	// The rows of an isolated score question share the state even in Decide.
	mode := manySeparate
	if many {
		mode = d.manyMode
	} else if len(rs) > 1 && d.nSeqMax > 1 {
		mode = ManyExact
	}
	vals, err := d.score(rs, mode)
	if err != nil {
		return nil, err
	}

	out := make([]*Result, len(qs))
	for i, q := range qs {
		p := plans[i]
		t := m.cfg.TypeTemperature(q.Type)
		if p.rows == 1 {
			logits := vals[p.first][0]
			probs, err := softmax(logits, t)
			if err != nil {
				return nil, questionErr(many, i, err)
			}
			out[i] = newResult(p.names, probs, logits, t, len(rs[p.first].ids), 0)
			continue
		}

		// Each level's chance to fit, normalized over the levels.
		fit := make([]float64, p.rows)
		var total float64
		tokens := 0
		for j := range fit {
			probs, err := softmax(vals[p.first+j][0], t)
			if err != nil {
				return nil, questionErr(many, i, err)
			}
			fit[j] = probs[1]
			total += fit[j]
			tokens += len(rs[p.first+j].ids)
		}
		if total == 0 {
			total = 1e-9
		}
		for j := range fit {
			fit[j] /= total
		}
		out[i] = newResult(p.names, fit, nil, t, tokens, 0)
	}

	return out, nil
}

// deciderState passes a string through and encodes any other value as JSON,
// with each element of an array of 8 or more given an "_index" field.
func deciderState(state any) (string, error) {
	if s, ok := state.(string); ok {
		return s, nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(state); err != nil {
		return "", fmt.Errorf("state: %w", err)
	}

	dec := json.NewDecoder(&buf)
	dec.UseNumber()
	n, err := readJSON(dec)
	if err != nil {
		return "", fmt.Errorf("state: %w", err)
	}

	var out bytes.Buffer
	n.annotate().write(&out)
	return spaceJSON(out.Bytes()), nil
}

// jsonNode is a JSON value that keeps its key order.
type jsonNode struct {
	kind json.Delim
	keys []string
	vals []*jsonNode
	raw  []byte
}

func readJSON(dec *json.Decoder) (*jsonNode, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}

	switch tok {
	case json.Delim('{'), json.Delim('['):
		n := &jsonNode{kind: tok.(json.Delim)}
		for dec.More() {
			if n.kind == '{' {
				key, err := dec.Token()
				if err != nil {
					return nil, err
				}
				n.keys = append(n.keys, key.(string))
			}
			v, err := readJSON(dec)
			if err != nil {
				return nil, err
			}
			n.vals = append(n.vals, v)
		}
		_, err := dec.Token()
		return n, err
	}

	raw, err := marshalJSON(tok)
	return &jsonNode{raw: raw}, err
}

func marshalJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// annotate returns n with the positions written into arrays of 8 or more.
// An object keeps its keys after "_index", and any other value moves to "value".
func (n *jsonNode) annotate() *jsonNode {
	if n.kind == 0 {
		return n
	}

	out := &jsonNode{kind: n.kind, keys: n.keys, vals: make([]*jsonNode, len(n.vals))}
	for i, v := range n.vals {
		out.vals[i] = v.annotate()
	}
	if n.kind != '[' || len(n.vals) < 8 {
		return out
	}

	for i, v := range out.vals {
		idx := &jsonNode{raw: []byte(strconv.Itoa(i))}
		if v.kind != '{' {
			out.vals[i] = &jsonNode{kind: '{', keys: []string{"_index", "value"}, vals: []*jsonNode{idx, v}}
			continue
		}
		// An element's own "_index" keeps the first place, as a Python dict merge does.
		obj := &jsonNode{kind: '{', keys: []string{"_index"}, vals: []*jsonNode{idx}}
		for k, key := range v.keys {
			if key == "_index" {
				obj.vals[0] = v.vals[k]
				continue
			}
			obj.keys = append(obj.keys, key)
			obj.vals = append(obj.vals, v.vals[k])
		}
		out.vals[i] = obj
	}

	return out
}

// write writes n as compact JSON.
func (n *jsonNode) write(b *bytes.Buffer) {
	if n.kind == 0 {
		b.Write(n.raw)
		return
	}

	b.WriteByte(byte(n.kind))
	for i, v := range n.vals {
		if i > 0 {
			b.WriteByte(',')
		}
		if n.kind == '{' {
			key, _ := marshalJSON(n.keys[i])
			b.Write(key)
			b.WriteByte(':')
		}
		v.write(b)
	}
	if n.kind == '{' {
		b.WriteByte('}')
	} else {
		b.WriteByte(']')
	}
}
