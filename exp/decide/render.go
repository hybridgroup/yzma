package decide

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/hybridgroup/yzma/pkg/llama"
)

// ErrBudget is returned when a rendered input is over a token budget.
// Inputs are never truncated.
var ErrBudget = errors.New("input over token budget")

type encoder func(text string) []llama.Token

// rendered is one model input. The logits of rows are read at each slot.
// Inputs about one state can share their first prefixLen tokens.
type rendered struct {
	ids        []llama.Token
	slots      []int
	rows       []llama.Token
	names      []string
	prefixLen  int
	headTokens int
}

type renderer struct {
	enc             encoder
	yesNo           []llama.Token
	arrow, nl, dash []llama.Token
	judge           []llama.Token
	maxLen, headMax int
}

func newRenderer(enc encoder, cfg *Config) (*renderer, error) {
	for _, st := range []SlotToken{cfg.SlotTokens.Yes, cfg.SlotTokens.No, cfg.SlotTokens.VerdictSlot} {
		got := enc(st.Text)
		if len(got) != 1 || int32(got[0]) != st.ID {
			return nil, fmt.Errorf("tokenizer mismatch: %q gives %v, readout config expects [%d]", st.Text, got, st.ID)
		}
	}

	return &renderer{
		enc:     enc,
		yesNo:   []llama.Token{llama.Token(cfg.SlotTokens.Yes.ID), llama.Token(cfg.SlotTokens.No.ID)},
		arrow:   []llama.Token{llama.Token(cfg.SlotTokens.VerdictSlot.ID)},
		nl:      enc("\n"),
		dash:    enc("- "),
		judge:   enc("Judge each option:\n"),
		maxLen:  cfg.Budgets.MaxLen,
		headMax: cfg.Budgets.HeadMax,
	}, nil
}

func (r *renderer) render(state string, q Question) (*rendered, error) {
	names, err := q.Names()
	if err != nil {
		return nil, err
	}

	opts := q.renderOptions()
	optIDs := make([][]llama.Token, len(opts))
	for i, o := range opts {
		optIDs[i] = r.enc(o)
	}

	head := r.enc(fmt.Sprintf("Question [%s]: %s\nOptions:\n", q.Type, q.Text))
	for _, o := range optIDs {
		head = append(head, r.dash...)
		head = append(head, o...)
		head = append(head, r.nl...)
	}
	head = append(head, r.judge...)

	rel := make([]int, len(optIDs))
	for i, o := range optIDs {
		head = append(head, o...)
		head = append(head, r.arrow...)
		rel[i] = len(head) - 1
		head = append(head, r.nl...)
	}
	if len(head) > r.headMax {
		return nil, fmt.Errorf("%w: question and options need %d tokens, the head budget is %d", ErrBudget, len(head), r.headMax)
	}

	ids := r.enc("State:\n")
	ids = append(ids, r.enc(state)...)
	ids = append(ids, r.enc("\n\n")...)
	prefixLen := len(ids)
	ids = append(ids, head...)
	if len(ids) > r.maxLen {
		return nil, fmt.Errorf("%w: input needs %d tokens (state %d, head %d), the limit is %d",
			ErrBudget, len(ids), prefixLen, len(head), r.maxLen)
	}

	slots := make([]int, len(rel))
	for i, s := range rel {
		slots[i] = prefixLen + s
	}

	return &rendered{ids: ids, slots: slots, rows: r.yesNo, names: names, prefixLen: prefixLen, headTokens: len(head)}, nil
}

// jevStyle is the Jev-Style family. The score of option k is
// logit(" yes") minus logit(" no") at its slot.
type jevStyle struct {
	cfg    *Config
	render *renderer
}

func (j *jevStyle) decide(d *Decider, state any, qs []Question, category string, many bool) ([]*Result, error) {
	s, err := serializeState(state)
	if err != nil {
		return nil, err
	}

	rs := make([]*rendered, len(qs))
	for i, q := range qs {
		if rs[i], err = j.render.render(s, q); err == nil && len(rs[i].slots) > d.maxOptions {
			err = fmt.Errorf("%w: %d options, the limit is %d", ErrQuestion, len(rs[i].slots), d.maxOptions)
		}
		if err != nil {
			return nil, questionErr(many, i, err)
		}
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
	for i, q := range qs {
		scores := make([]float64, len(vals[i]))
		for k, v := range vals[i] {
			scores[k] = v[0] - v[1]
		}
		t := j.cfg.Temperature(category, q.Type, len(rs[i].names))
		p, err := softmax(scores, t)
		if err != nil {
			return nil, questionErr(many, i, err)
		}
		out[i] = newResult(rs[i].names, p, scores, t, len(rs[i].ids), rs[i].headTokens)
	}

	return out, nil
}

// questionErr adds the question index to err when there are several questions.
func questionErr(many bool, i int, err error) error {
	if many {
		return fmt.Errorf("question %d: %w", i, err)
	}
	return err
}

// serializeState passes a string through and encodes any other value as JSON
// with the ", " and ": " separators the model was trained on.
func serializeState(state any) (string, error) {
	if s, ok := state.(string); ok {
		return s, nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(state); err != nil {
		return "", fmt.Errorf("state: %w", err)
	}

	return spaceJSON(bytes.TrimRight(buf.Bytes(), "\n")), nil
}

// spaceJSON adds a space after each comma and colon outside strings, and
// writes U+2028 and U+2029 raw, as Python's json.dumps does.
func spaceJSON(b []byte) string {
	out := make([]byte, 0, len(b)+len(b)/8)
	inStr := false
	for i := 0; i < len(b); i++ {
		c := b[i]
		switch {
		case inStr && c == '\\':
			if i+5 < len(b) && b[i+1] == 'u' && string(b[i+2:i+5]) == "202" && (b[i+5] == '8' || b[i+5] == '9') {
				out = utf8.AppendRune(out, rune(0x2020+int(b[i+5]-'0')))
				i += 5
				continue
			}
			out = append(out, c, b[i+1])
			i++
			continue
		case c == '"':
			inStr = !inStr
		case !inStr && (c == ',' || c == ':'):
			out = append(out, c, ' ')
			continue
		}
		out = append(out, c)
	}
	return string(out)
}
