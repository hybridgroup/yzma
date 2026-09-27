package decide

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"
)

// The expected values in these tests come from the decider 0.8b-v1 reference, prompt.py and systemone.py.

func TestParseDeciderConfig(t *testing.T) {
	c, err := ParseDeciderConfig([]byte(`{"temperature": 1.03, "isolated_levels": true, "neutralize_none": false, "max_options": 255}`))
	if err != nil {
		t.Fatal(err)
	}
	if c.Temperature != 1.03 || !c.IsolatedLevels || c.NeutralizeNone || c.TypeTemperature(TypeNoul) != 1.03 {
		t.Errorf("config = %+v", c)
	}

	c, err = ParseDeciderConfig([]byte(`{"temperature_by_type": {"noul": 2.2}, "layout": "plain"}`))
	if err != nil {
		t.Fatal(err)
	}
	if c.Temperature != 1 || !c.NeutralizeNone || c.IsolatedLevels || c.TypeTemperature(TypeNoul) != 2.2 || c.TypeTemperature(TypeChoice) != 1 {
		t.Errorf("defaults = %+v", c)
	}

	for _, s := range []string{
		`{"temperature": 0}`,
		`{"temperature": -1}`,
		`{"temperature_by_type": {"bool": 1.1}}`,
		`{"temperature_by_type": {"score": 0}}`,
		`{"layout": "chat"}`,
		`{"chat_template": true}`,
		`[]`,
	} {
		if _, err := ParseDeciderConfig([]byte(s)); err == nil {
			t.Errorf("%s: no error", s)
		}
	}
}

func TestDeciderRows(t *testing.T) {
	m := &deciderModel{cfg: &DeciderConfig{Temperature: 1, IsolatedLevels: true, NeutralizeNone: true}}

	for _, tc := range []struct {
		q    Question
		rows []deciderRow
	}{
		{ChoiceDesc("Team?", Option{Name: "billing", Description: "Payments"}, Option{Name: "tech"}, Option{Name: "None of the above"}),
			[]deciderRow{{"Team?", []string{"billing: Payments", "tech", "not listed here"}}}},
		{Noul("Refund?", "", ""), []deciderRow{{"Refund?", []string{"no", "yes"}}}},
		{Noul("Refund?", "", "asks for money back"), []deciderRow{{"Refund?", []string{"no", "yes: asks for money back"}}}},
		{Score("Urgent?", "low", "2 : high"), []deciderRow{
			{"Urgent?\nProposed answer: low\nDoes the proposed answer fit?", []string{"no", "yes"}},
			{"Urgent?\nProposed answer: high\nDoes the proposed answer fit?", []string{"no", "yes"}},
		}},
	} {
		_, rows, err := m.rows(tc.q, 255)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.EqualFunc(rows, tc.rows, func(a, b deciderRow) bool { return a.text == b.text && slices.Equal(a.options, b.options) }) {
			t.Errorf("%s: rows = %q, want %q", tc.q.Text, rows, tc.rows)
		}
	}

	m.cfg.IsolatedLevels = false
	m.cfg.NeutralizeNone = false
	_, rows, err := m.rows(Score("Urgent?", "low", "high"), 255)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || !slices.Equal(rows[0].options, []string{"0: low", "1: high"}) {
		t.Errorf("listwise score rows = %q", rows)
	}
	_, rows, _ = m.rows(Choice("x", "none", "b"), 255)
	if rows[0].options[0] != "none" {
		t.Errorf("neutralize off: %q", rows[0].options)
	}

	for _, q := range []Question{Choice("x", "a"), Choice("x", "a", "b", "c"), Noul("", "", "")} {
		if _, _, err := m.rows(q, 2); !errors.Is(err, ErrQuestion) {
			t.Errorf("%+v: got %v, want ErrQuestion", q, err)
		}
	}
}

func TestDeciderPrompt(t *testing.T) {
	got := deciderPrompt(deciderRow{"Team?", []string{"billing: Payments", "tech"}})
	want := "\n\nQuestion: Team?\nOptions:\n(A) billing: Payments\n(B) tech\nAnswer: ("
	if got != want {
		t.Errorf("prompt = %q, want %q", got, want)
	}
}

// labelEncode gives one token to each string of one or two upper case
// letters except "AC", and one token per byte to anything else.
func labelEncode(s string) []token {
	if s != "AC" && len(s) <= 2 && strings.Trim(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") == "" {
		id := token(1000 + int(s[0]))
		if len(s) == 2 {
			id = token(2000 + 26*int(s[0]-'A') + int(s[1]-'A'))
		}
		return []token{id}
	}
	return fakeEncode(s)
}

func TestDeciderLabels(t *testing.T) {
	m, err := newDeciderModel(fakeEncode, labelEncode, &DeciderConfig{Temperature: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(m.labels) != 255 || m.labels[0] != 1000+'A' || m.labels[27] != 2000+1 || m.labels[28] != 2000+3 {
		t.Errorf("labels = %v", m.labels[:30])
	}
	if _, err := newDeciderModel(fakeEncode, fakeEncode, &DeciderConfig{Temperature: 1}); err == nil {
		t.Error("letters of two tokens: no error")
	}
}

func TestDeciderRender(t *testing.T) {
	m, err := newDeciderModel(fakeEncode, labelEncode, &DeciderConfig{Temperature: 1})
	if err != nil {
		t.Fatal(err)
	}
	d := &Decider{nCtx: 4096}
	ctx := fakeEncode("Context:\nhi")

	r, err := m.render(d, ctx, deciderRow{"Q?", []string{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}
	if got := decodeFake(r.ids); got != "Context:\nhi"+deciderPrompt(deciderRow{"Q?", []string{"a", "b"}}) {
		t.Errorf("narrow = %q", got)
	}
	if r.prefixLen != len(ctx) || !slices.Equal(r.slots, []int{len(r.ids) - 1}) || len(r.rows) != 2 {
		t.Errorf("narrow rendered = %+v", r)
	}

	var opts []string
	for i := range 12 {
		opts = append(opts, string(rune('a'+i)))
	}
	r, err = m.render(d, ctx, deciderRow{"Q?", opts})
	if err != nil {
		t.Fatal(err)
	}
	head := len(ctx) + len(fakeEncode("\n\nQuestion: Q?\nOptions:"))
	if r.ids[head+2] != m.labels[0] || r.ids[head+8] != m.labels[1] || len(r.rows) != 12 {
		t.Errorf("wide ids = %v", r.ids[head:])
	}

	d.nCtx = 20
	if _, err := m.render(d, ctx, deciderRow{"Q?", []string{"a", "b"}}); !errors.Is(err, ErrBudget) {
		t.Errorf("got %v, want ErrBudget", err)
	}
}

func TestDeciderState(t *testing.T) {
	for _, tc := range []struct {
		in   any
		want string
	}{
		{"plain text", "plain text"},
		{json.RawMessage(`{"b":[1,2],"a":"<&> "}`), `{"b": [1, 2], "a": "<&>` + " " + `"}`},
		{json.RawMessage(`{"r":[0,1.50,"x",null,true,[1],{"k":"v"},{"_index":"own","z":1}]}`),
			`{"r": [{"_index": 0, "value": 0}, {"_index": 1, "value": 1.50}, {"_index": 2, "value": "x"}, {"_index": 3, "value": null}, ` +
				`{"_index": 4, "value": true}, {"_index": 5, "value": [1]}, {"_index": 6, "k": "v"}, {"_index": "own", "z": 1}]}`},
		{json.RawMessage(`[1,2,3,4,5,6,7]`), `[1, 2, 3, 4, 5, 6, 7]`},
		{map[string]int{"z": 1, "a": 2}, `{"a": 2, "z": 1}`},
	} {
		got, err := deciderState(tc.in)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("state = %s\nwant %s", got, tc.want)
		}
	}
}
