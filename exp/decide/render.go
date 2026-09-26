package decide

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hybridgroup/yzma/pkg/llama"
)

// ErrBudget is returned when a rendered input is over a token budget.
// Inputs are never truncated.
var ErrBudget = errors.New("input over token budget")

type encoder func(text string) []llama.Token

type rendered struct {
	ids        []llama.Token
	slots      []int
	names      []string
	prefixLen  int
	headTokens int
}

type renderer struct {
	enc             encoder
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

	return &rendered{ids: ids, slots: slots, names: names, prefixLen: prefixLen, headTokens: len(head)}, nil
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

// spaceJSON adds a space after each comma and colon outside strings.
func spaceJSON(b []byte) string {
	out := make([]byte, 0, len(b)+len(b)/8)
	inStr, esc := false, false
	for _, c := range b {
		out = append(out, c)
		switch {
		case esc:
			esc = false
		case inStr && c == '\\':
			esc = true
		case c == '"':
			inStr = !inStr
		case !inStr && (c == ',' || c == ':'):
			out = append(out, ' ')
		}
	}
	return string(out)
}
