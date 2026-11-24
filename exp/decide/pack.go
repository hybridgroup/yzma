package decide

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// lfm3-d1 reads up to packSize questions after their state in one pass, as
// llama-server and pack_set of the model repo do.
const (
	packSize   = 16
	packHead   = "<|startoftext|><|im_start|>user\n"
	packHeader = "QUESTION:\n"
)

// packSuffix returns the prompt of q with a null state, without its head.
func (t *typed) packSuffix(q Question, opts []typedOption, labels []string, qs []Question) (string, error) {
	lone, err := t.prompt(nil, q, opts, labels, qs)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(lone, packHead) {
		return "", errLayout
	}
	return lone[len(packHead):], nil
}

// packed reads the questions in sets of packSize. rs[i] is the input of
// question i alone and suffixes[i] its packSuffix. It returns the scores of each
// input and the tokens of each set, counted at its first question.
func (t *typed) packed(d *Decider, rs []*rendered, suffixes []string, nullState bool) ([][][]float64, []int, error) {
	out := make([][][]float64, len(rs))
	tokens := make([]int, len(rs))
	for i := 0; i < len(rs); i += packSize {
		end := min(len(rs), i+packSize)

		// The first question is its own input, the others follow it.
		ids := slices.Clone(rs[i].ids)
		n := len(t.enc(suffixes[i]))
		if n >= len(ids) {
			return nil, nil, errLayout
		}
		group := make([]int32, len(ids))
		for j := len(ids) - n; j < len(ids); j++ {
			group[j] = 1
		}
		slots := []int{len(ids) - 1}
		rows := [][]token{rs[i].rows}
		for k := i + 1; k < end; k++ {
			text := suffixes[k]
			if !nullState {
				text = packHeader + text
			}
			for _, tok := range t.enc(text) {
				ids = append(ids, tok)
				group = append(group, int32(k-i+1))
			}
			slots = append(slots, len(ids)-1)
			rows = append(rows, rs[k].rows)
		}
		if limit := min(d.ubatch, d.nCtx); len(ids) > limit {
			return nil, nil, fmt.Errorf("%w: input needs %d tokens, the context is %d", ErrBudget, len(ids), limit)
		}

		v, err := d.decodeGrouped(ids, group, slots, rows)
		if err != nil {
			return nil, nil, err
		}
		for k := i; k < end; k++ {
			out[k] = [][]float64{v[k-i]}
		}
		tokens[i] = len(ids)
	}
	return out, tokens, nil
}

// decodeGrouped decodes one packed input and returns the logits of rows[j] at slots[j].
func (d *Decider) decodeGrouped(ids []token, group []int32, slots []int, rows [][]token) ([][]float64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.b.open() {
		return nil, errors.New("decide: decider is closed")
	}
	if err := d.b.decodeGrouped(ids, group, slots); err != nil {
		return nil, err
	}

	out := make([][]float64, len(slots))
	for j, s := range slots {
		logits, err := d.b.logits(int32(s))
		if err != nil {
			return nil, err
		}
		if logits == nil {
			return nil, fmt.Errorf("decide: no logits at batch index %d", s)
		}
		out[j] = make([]float64, len(rows[j]))
		for r, t := range rows[j] {
			out[j][r] = float64(logits[t])
		}
	}
	return out, nil
}
