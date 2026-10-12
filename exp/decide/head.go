package decide

import (
	"bytes"
	"fmt"
	"strings"
)

// The text llama-server gives the clef and vllm-decision2 templates between
// the pieces of the prompt and at the start of a span.
const (
	headMarker   = "<<clef:"
	headSep      = "<<clef:sep>>"
	headQuestion = "<<clef:question>>"
	headOption   = "<<clef:option>>"
)

// The decision order of a token, the same values as llama_decision_order.
const (
	orderNone           int32 = 0
	orderQuestionNoul   int32 = 1
	orderQuestionChoice int32 = 2
	orderQuestionScore  int32 = 3
	orderOption         int32 = 4
)

// The tojson filter of the Go templates escapes these characters and the one
// of llama.cpp does not, so they are hidden from it.
var (
	jsonHide   = strings.NewReplacer("<", "\U0010FFF0", ">", "\U0010FFF1", "&", "\U0010FFF2", " ", "\U0010FFF3", " ", "\U0010FFF4")
	jsonShow   = strings.NewReplacer("\U0010FFF0", "<", "\U0010FFF1", ">", "\U0010FFF2", "&", "\U0010FFF3", " ", "\U0010FFF4", " ")
	jsonHidden = "\U0010FFF0\U0010FFF1\U0010FFF2\U0010FFF3\U0010FFF4"
)

// headText cleans user text for a decision head prompt.
func (t *typed) headText(s string) (string, error) {
	if strings.ContainsAny(s, jsonHidden) {
		return "", fmt.Errorf("%w: text has characters from U+10FFF0 to U+10FFF4", ErrQuestion)
	}
	return strings.ReplaceAll(t.text(s), headMarker, "<<clef "), nil
}

// headState returns the state of a decision head prompt, a JSON value has its keys sorted.
// vllm-decision2 only takes text, an object or an array.
func (t *typed) headState(state any) (string, error) {
	if _, ok := state.(string); !ok && t.kind == typeVllmDecision2 {
		data, err := stateJSON(state)
		if err != nil {
			return "", err
		}
		if d := bytes.TrimSpace(data); len(d) == 0 || (d[0] != '{' && d[0] != '[') {
			return "", fmt.Errorf("%w: state must be text, an object or an array", ErrQuestion)
		}
	}
	s, err := stateText(state, styleCompact.sorted())
	if err != nil {
		return "", err
	}
	return t.headText(s)
}

// headPrompt renders qs about state in one prompt and marks the spans the head reads.
// It also returns the options of each question in the order of the scores.
func (t *typed) headPrompt(state string, qs []Question) (*rendered, [][]typedOption, error) {
	questions := make([]any, len(qs))
	opts := make([][]typedOption, len(qs))
	total := 0
	for i, q := range qs {
		text, err := t.headText(q.Text)
		if err != nil {
			return nil, nil, err
		}
		opts[i] = t.options(q)
		options := make([]any, len(opts[i]))
		for j, o := range opts[i] {
			key, err := t.headText(o.key)
			if err != nil {
				return nil, nil, err
			}
			opt := map[string]any{"key": jsonHide.Replace(key), "description": nil}
			if o.description != "" {
				desc, err := t.headText(o.description)
				if err != nil {
					return nil, nil, err
				}
				opt["description"] = jsonHide.Replace(desc)
			}
			options[j] = opt
		}
		questions[i] = map[string]any{"id": q.id, "type": string(q.Type), "instructions": text, "options": options}
		total += len(options)
	}

	prompt, err := t.tmpl.Render(map[string]any{
		"state":         state,
		"questions":     questions,
		"images":        []any{},
		"sep":           headSep,
		"mark_question": headQuestion,
		"mark_option":   headOption,
	})
	if err != nil {
		return nil, nil, err
	}
	prompt = jsonShow.Replace(prompt)

	// The model was trained with the pieces tokenized one by one.
	r := &rendered{prefixLen: -1}
	next := 0
	for _, piece := range strings.Split(prompt, headSep) {
		order := orderNone
		if rest, ok := strings.CutPrefix(piece, headQuestion); ok {
			if next >= len(qs) {
				return nil, nil, errLayout
			}
			piece, order = rest, questionOrder(qs[next].Type)
			next++
		} else if rest, ok := strings.CutPrefix(piece, headOption); ok {
			piece, order = rest, orderOption
			r.nScores++
		}

		ids := t.enc(piece)
		if order != orderNone && len(ids) == 0 {
			return nil, nil, fmt.Errorf("%w: the instructions and the options of a question must not be empty", ErrQuestion)
		}
		if order != orderNone && r.prefixLen < 0 {
			r.prefixLen = len(r.ids)
		}
		r.ids = append(r.ids, ids...)
		for range ids {
			r.order = append(r.order, order)
		}
	}
	if next != len(qs) || r.nScores != total || r.prefixLen < 0 {
		return nil, nil, errLayout
	}

	return r, opts, nil
}

func questionOrder(t Type) int32 {
	switch t {
	case TypeNoul:
		return orderQuestionNoul
	case TypeChoice:
		return orderQuestionChoice
	}
	return orderQuestionScore
}

// decideHead scores qs with a decision head in the llama.cpp graph. Clef reads
// all the questions in one prompt, vllm-decision2 reads each in its own.
func (t *typed) decideHead(d *Decider, state any, qs []Question, many bool) ([]*Result, error) {
	s, err := t.headState(state)
	if err != nil {
		return nil, err
	}

	names := make([][]string, len(qs))
	for i, q := range qs {
		if names[i], err = q.Names(); err != nil {
			return nil, questionErr(many, i, err)
		}
		if limit := min(d.maxOptions, t.maxOptions); len(q.Options) > limit {
			return nil, questionErr(many, i, fmt.Errorf("%w: %d options, the limit is %d", ErrQuestion, len(q.Options), limit))
		}
		if t.kind == typeVllmDecision2 && q.Type == TypeChoice && len(q.Options) < 2 {
			return nil, questionErr(many, i, fmt.Errorf("%w: choice needs at least two options", ErrQuestion))
		}
	}

	groups := [][]int{}
	if t.kind == typeClef {
		all := make([]int, len(qs))
		for i := range all {
			all[i] = i
		}
		groups = append(groups, all)
	} else {
		for i := range qs {
			groups = append(groups, []int{i})
		}
	}

	type plan struct {
		input, first int
		opts         []typedOption
	}
	plans := make([]plan, len(qs))
	rs := make([]*rendered, len(groups))
	for g, idx := range groups {
		group := make([]Question, len(idx))
		for k, i := range idx {
			group[k] = qs[i]
		}
		r, opts, err := t.headPrompt(s, group)
		if err == nil && len(r.ids) > d.nCtx {
			err = fmt.Errorf("%w: input needs %d tokens, the context is %d", ErrBudget, len(r.ids), d.nCtx)
		}
		if err != nil {
			if len(idx) == 1 {
				return nil, questionErr(many, idx[0], err)
			}
			return nil, err
		}
		rs[g] = r
		first := 0
		for k, i := range idx {
			plans[i] = plan{input: g, first: first, opts: opts[k]}
			first += len(opts[k])
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
	counted := make([]bool, len(rs))
	for i, p := range plans {
		logits := make([]float64, len(p.opts))
		for k := range logits {
			logits[k] = vals[p.input][p.first+k][0]
		}
		temp := t.temperature(qs[i])
		sm, err := softmax(t.biased(logits, qs[i].Type), temp)
		if err != nil {
			return nil, questionErr(many, i, err)
		}

		probs := make([]float64, len(names[i]))
		scores := make([]float64, len(names[i]))
		for k, o := range p.opts {
			probs[o.index] = sm[k]
			scores[o.index] = logits[k]
		}
		// The tokens of a prompt with several questions count once.
		tokens := 0
		if !counted[p.input] {
			tokens, counted[p.input] = len(rs[p.input].ids), true
		}
		out[i] = newResult(names[i], probs, scores, temp, tokens, 0)
	}

	return out, nil
}
