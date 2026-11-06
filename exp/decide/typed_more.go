package decide

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// The text llama-server gives the lfm2-d1-omni template between the pieces of
// the prompt, and at the start of the pieces that are cut to a token budget.
const (
	d1OmniMarker       = "<<d1omni:"
	d1OmniSep          = "<<d1omni:sep>>"
	d1OmniMarkState    = "<<d1omni:state>>"
	d1OmniMarkQuestion = "<<d1omni:question>>"
	d1OmniMarkOption   = "<<d1omni:option>>"
	// d1OmniMaxTokens is max_length of the model config.
	d1OmniMaxTokens = 16384
)

// h2oMaxStateTokens is the state budget of serve_config.json, a longer state keeps its head and tail.
const h2oMaxStateTokens = 32000

// takesNullState reports whether the model reads a request with a null state.
func (t *typed) takesNullState() bool {
	return t.kind == typeLfm2D1 || t.kind == typeLfm2D1Omni || t.kind == typeVllmD3 || t.isLfm3()
}

// d1Break puts a zero width space after the first character of each special
// token text in s, again while a break makes a new match (d1 prompt.py: plain).
func (t *typed) d1Break(s string) string {
	for changed := true; changed; {
		changed = false
		var b strings.Builder
		for i := 0; i < len(s); {
			hit := t.quotedAt(s, i)
			if hit == "" {
				b.WriteByte(s[i])
				i++
				continue
			}
			_, n := utf8.DecodeRuneInString(hit)
			b.WriteString(hit[:n] + "\u200b" + hit[n:])
			i += len(hit)
			changed = true
		}
		s = b.String()
	}
	return s
}

// h2oQuote writes < and > of each special token text in s as JSON escapes, so they stay text.
func (t *typed) h2oQuote(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		hit := t.quotedAt(s, i)
		if hit == "" {
			b.WriteByte(s[i])
			i++
			continue
		}
		b.WriteString(strings.NewReplacer("<", "\\u003c", ">", "\\u003e").Replace(hit))
		i += len(hit)
	}
	return b.String()
}

// quotedAt returns the longest special token text at s[i:], or "".
func (t *typed) quotedAt(s string, i int) string {
	if !strings.Contains(t.quotedFirst, s[i:i+1]) {
		return ""
	}
	for _, q := range t.quoted {
		if strings.HasPrefix(s[i:], q) {
			return q
		}
	}
	return ""
}

// h2oState returns the state as compact JSON, cut and quoted as the reference shim does.
func (t *typed) h2oState(state any) (string, error) {
	s, ok := state.(string)
	if !ok {
		data, err := stateJSON(state)
		if err != nil {
			return "", err
		}
		if s, err = pyDumps(data); err != nil {
			return "", fmt.Errorf("state: %w", err)
		}
	}
	return t.h2oQuote(t.h2oCut(s)), nil
}

// h2oCut keeps the head and the tail of a state longer than the budget.
func (t *typed) h2oCut(s string) string {
	if len(s) <= h2oMaxStateTokens {
		return s
	}
	ids := t.enc(s)
	if len(ids) <= h2oMaxStateTokens {
		return s
	}
	half := h2oMaxStateTokens / 2
	return t.pieces(ids[:half]) + "\n...\n" + t.pieces(ids[len(ids)-half:])
}

// pieces returns the text of ids, an invalid UTF-8 byte becomes U+FFFD as in a JSON response.
func (t *typed) pieces(ids []token) string {
	var b strings.Builder
	for _, id := range ids {
		b.WriteString(t.piece(id))
	}
	s := b.String()
	if utf8.ValidString(s) {
		return s
	}
	var out strings.Builder
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		out.WriteRune(r)
		i += n
	}
	return out.String()
}

// h2oPrompt renders the template with the texts the reference shim gives it.
// state is the text from h2oState.
func (t *typed) h2oPrompt(state any, q Question, opts []typedOption) (string, error) {
	options := make([]any, len(opts))
	for i, o := range opts {
		key := t.text(o.key)
		var desc any
		switch {
		case o.description != "":
			desc = t.text(o.description)
		case q.Type == TypeScore:
			desc = "null"
		}
		// The template trims ASCII space only, Python strips the end of the user turn.
		if i == len(opts)-1 {
			text, ok := desc.(string)
			if desc == nil && q.Type == TypeChoice {
				text, ok = key, true
			}
			if s := strings.TrimRightFunc(text, pyIsSpace); ok && s != text {
				desc = s
				if s == "" {
					desc = " "
				}
			}
		}
		options[i] = map[string]any{"label": t.codes[i], "key": key, "description": desc}
	}

	return t.tmpl.Render(map[string]any{
		"type":         string(q.Type),
		"state":        state,
		"instructions": t.text(strings.TrimLeftFunc(q.Text, pyIsSpace)),
		"options":      options,
		"images":       []any{},
	})
}

// pyIsSpace is str.isspace of Python.
func pyIsSpace(r rune) bool {
	return unicode.IsSpace(r) || (r >= 0x1c && r <= 0x1f)
}

// d1OmniRead cuts each marked piece of an lfm2-d1-omni prompt to its token
// budget, as d1-omni prompt.py encode does, and finds the marker of each option.
func (t *typed) d1OmniRead(prompt string, n int) (*rendered, error) {
	pieces := strings.Split(prompt, d1OmniSep)
	if len(pieces) == 0 || pieces[0] != "" {
		return nil, errLayout
	}

	budget := max(96, min(24*n+32, d1OmniMaxTokens/2))
	optionMax := max(2, (budget-3*n)/n)
	questionMax := max(16, budget)

	var head, body, tail []token
	hasState := false
	for _, p := range pieces[1:] {
		limit := -1
		isState := false
		switch {
		case strings.HasPrefix(p, d1OmniMarkState):
			p, isState = p[len(d1OmniMarkState):], true
		case strings.HasPrefix(p, d1OmniMarkQuestion):
			p, limit = p[len(d1OmniMarkQuestion):], questionMax
		case strings.HasPrefix(p, d1OmniMarkOption):
			p, limit = p[len(d1OmniMarkOption):], optionMax
		}
		ids := t.enc(p)
		if limit >= 0 && len(ids) > limit {
			ids = ids[:limit]
		}
		switch {
		case isState && hasState:
			return nil, errLayout
		case isState:
			body, hasState = ids, true
		case hasState:
			tail = append(tail, ids...)
		default:
			head = append(head, ids...)
		}
	}
	if !hasState {
		return nil, errLayout
	}

	room := d1OmniMaxTokens - min(d1OmniMaxTokens, len(head)+len(tail))
	ids := append(head, body[:min(len(body), room)]...)
	ids = append(ids, tail...)
	ids = ids[:min(len(ids), d1OmniMaxTokens)]

	var slots []int
	for i, id := range ids {
		if id == t.marker {
			slots = append(slots, i)
		}
	}
	if len(slots) != n {
		return nil, fmt.Errorf("%w: the options do not fit in the context", ErrBudget)
	}
	return &rendered{ids: ids, slots: slots, embd: true}, nil
}

// pyDumps writes a JSON value as json.dumps of Python does with compact
// separators and ensure_ascii off.
func pyDumps(data []byte) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	n, err := readJSON(dec)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	writePy(&b, n)
	return b.String(), nil
}

func writePy(b *strings.Builder, n *jsonNode) {
	switch n.kind {
	case '[':
		b.WriteByte('[')
		for i, v := range n.vals {
			if i > 0 {
				b.WriteByte(',')
			}
			writePy(b, v)
		}
		b.WriteByte(']')
	case '{':
		b.WriteByte('{')
		for i, k := range n.keys {
			if i > 0 {
				b.WriteByte(',')
			}
			writeJSONString(b, k)
			b.WriteByte(':')
			writePy(b, n.vals[i])
		}
		b.WriteByte('}')
	default:
		var s string
		raw := string(n.raw)
		switch {
		case strings.HasPrefix(raw, `"`) && json.Unmarshal(n.raw, &s) == nil:
			writeJSONString(b, s)
		case raw == "true", raw == "false", raw == "null", !strings.ContainsAny(raw, ".eE"):
			b.WriteString(raw)
		default:
			b.WriteString(pyFloat(raw))
		}
	}
}
