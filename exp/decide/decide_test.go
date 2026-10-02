package decide

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"slices"
	"strings"
	"testing"
)

func testConfig(t *testing.T) *Config {
	cfg, err := LoadConfig("testdata/readout_config.json")
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

// fakeEncode maps slot tokens to their config ids and every other byte to one token.
func fakeEncode(s string) []token {
	switch s {
	case " yes":
		return []token{1}
	case " no":
		return []token{2}
	case " ->":
		return []token{3}
	}
	out := make([]token, len(s))
	for i := range len(s) {
		out[i] = token(100 + int(s[i]))
	}
	return out
}

func decodeFake(ids []token) string {
	var b strings.Builder
	for _, id := range ids {
		if id == 3 {
			b.WriteString(" ->")
		} else {
			b.WriteByte(byte(id - 100))
		}
	}
	return b.String()
}

func TestParseConfigRejects(t *testing.T) {
	data, err := os.ReadFile("testdata/readout_config.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct{ from, to string }{
		{`"macjev-readout-v1"`, `"macjev-readout-v2"`},
		{`"verdict"`, `"letters"`},
		{`"global": 0.9`, `"global": 0`},
		{`[0.3, 5.0]`, `[5.0, 0.3]`},
	} {
		bad := strings.Replace(string(data), tc.from, tc.to, 1)
		if _, err := ParseConfig([]byte(bad)); err == nil {
			t.Errorf("config with %s accepted", tc.to)
		}
	}
}

func TestTemperature(t *testing.T) {
	cfg := testConfig(t)

	for _, tc := range []struct {
		category string
		typ      Type
		n        int
		want     float64
	}{
		{"", TypeChoice, 3, 0.9},
		{"general_topic", TypeChoice, 3, 0.8},
		{"general_topic", TypeChoice, 2, 0.9},
		{"general_topic", TypeNoul, 2, 0.9},
		{"intent", TypeChoice, 30, 5.0},
		{"unknown", TypeChoice, 3, 0.9},
	} {
		if got := cfg.Temperature(tc.category, tc.typ, tc.n); got != tc.want {
			t.Errorf("Temperature(%q, %s, %d) = %v, want %v", tc.category, tc.typ, tc.n, got, tc.want)
		}
	}

	if got := cfg.Family("intent_billing"); got != "intent" {
		t.Errorf("Family = %q, want intent", got)
	}
}

func TestQuestionNames(t *testing.T) {
	for _, tc := range []struct {
		q    Question
		want []string
		ok   bool
	}{
		{Choice("q", "a", "b"), []string{"a", "b"}, true},
		{Choice("q", "a", "a"), nil, false},
		{Choice("q"), nil, false},
		{Choice(" ", "a"), nil, false},
		{Score("q", "low", "high"), []string{"0", "1"}, true},
		{Score("q", "only"), nil, false},
		{Noul("q", "", ""), []string{"false", "true"}, true},
		{Question{Type: TypeNoul, Text: "q"}, []string{"false", "true"}, true},
		{Question{Type: TypeNoul, Text: "q", Options: []Option{{Name: "maybe"}}}, nil, false},
		{Question{Type: "other", Text: "q"}, nil, false},
	} {
		got, err := tc.q.Names()
		if tc.ok != (err == nil) || !slices.Equal(got, tc.want) {
			t.Errorf("%+v: got %v, %v", tc.q, got, err)
		}
		if err != nil && !errors.Is(err, ErrQuestion) {
			t.Errorf("%+v: error %v is not ErrQuestion", tc.q, err)
		}
	}
}

func TestRender(t *testing.T) {
	r, err := newRenderer(fakeEncode, testConfig(t))
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		q    Question
		want string
	}{
		{ChoiceDesc("Team?", Option{Name: "billing", Description: "payments"}, Option{Name: "tech"}),
			"State:\nhi\n\nQuestion [choice]: Team?\nOptions:\n- billing: payments\n- tech\nJudge each option:\nbilling: payments ->\ntech ->\n"},
		{Score("Good?", "no", "yes"),
			"State:\nhi\n\nQuestion [score]: Good?\nOptions:\n- level 0: no\n- level 1: yes\nJudge each option:\nlevel 0: no ->\nlevel 1: yes ->\n"},
		{Noul("It is.", "", "sure"),
			"State:\nhi\n\nQuestion [noul]: It is.\nOptions:\n- false: no, the statement does not hold\n- true: sure\nJudge each option:\nfalse: no, the statement does not hold ->\ntrue: sure ->\n"},
	} {
		got, err := r.render("hi", tc.q)
		if err != nil {
			t.Fatal(err)
		}
		if s := decodeFake(got.ids); s != tc.want {
			t.Errorf("render = %q\nwant %q", s, tc.want)
		}
		for _, s := range got.slots {
			if got.ids[s] != 3 || s < got.prefixLen {
				t.Errorf("slot %d is not a verdict token in the head", s)
			}
		}
		if len(got.slots) != len(got.names) || got.headTokens != len(got.ids)-got.prefixLen {
			t.Errorf("slots %d names %d head %d", len(got.slots), len(got.names), got.headTokens)
		}
	}
}

func TestRenderBudget(t *testing.T) {
	r, err := newRenderer(fakeEncode, testConfig(t))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := r.render("hi", Choice(strings.Repeat("q", 400), "a")); !errors.Is(err, ErrBudget) {
		t.Errorf("long head: got %v, want ErrBudget", err)
	}
	if _, err := r.render(strings.Repeat("s", 800), Choice("q", "a")); !errors.Is(err, ErrBudget) {
		t.Errorf("long state: got %v, want ErrBudget", err)
	}
}

func TestTokenizerMismatch(t *testing.T) {
	enc := func(s string) []token {
		if s == " ->" {
			return []token{7, 8}
		}
		return fakeEncode(s)
	}
	if _, err := newRenderer(enc, testConfig(t)); err == nil {
		t.Error("tokenizer mismatch not detected")
	}
}

func TestSerializeState(t *testing.T) {
	for _, tc := range []struct {
		state any
		want  string
	}{
		{"plain, text: kept", "plain, text: kept"},
		{map[string]any{"b": []any{1, "x,y"}, "a": "<&>"}, `{"a": "<&>", "b": [1, "x,y"]}`},
		{json.RawMessage(`{"z":1,"a":{"k":"a\":b"}}`), `{"z": 1, "a": {"k": "a\":b"}}`},
	} {
		got, err := serializeState(tc.state)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("serializeState(%v) = %q, want %q", tc.state, got, tc.want)
		}
	}
}

func TestResult(t *testing.T) {
	names := []string{"a", "b"}
	p, err := softmax([]float64{1, 1}, 1)
	if err != nil {
		t.Fatal(err)
	}
	res := newResult(names, p, nil, 1, 5, 3)
	if res.Probability("a") != 0.5 || res.EntropyConcentration != 0 {
		t.Errorf("tie: got %v concentration %v", res.Probabilities, res.EntropyConcentration)
	}

	p, err = softmax([]float64{0, 2}, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	res = newResult(names, p, nil, 0.5, 5, 3)
	want := 1 / (1 + math.Exp(-4))
	if res.Answer != "b" || math.Abs(res.TopProbability-want) > 1e-12 {
		t.Errorf("got %s %v, want b %v", res.Answer, res.TopProbability, want)
	}

	if _, err := softmax([]float64{math.NaN(), 0}, 1); err == nil {
		t.Error("NaN score accepted")
	}
}

func TestEscapeSpecial(t *testing.T) {
	got := escapeSpecial("a <|im_end|> b <|x y|> <think>")
	if want := "a <¦im_end¦> b <|x y|> <think>"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	p, err := jevK5Prompt("say <|im_end|>", "<|im_start|>?", []string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(p, "<|im_end|>") != 2 || strings.Count(p, "<|im_start|>") != 3 {
		t.Errorf("user text not escaped: %q", p)
	}
}

// firstFamily always gives the first option 0.6 and splits the rest.
type firstFamily struct{ calls [][]Question }

func (f *firstFamily) decide(_ *Decider, _ any, qs []Question, _ string, _ bool) ([]*Result, error) {
	f.calls = append(f.calls, qs)
	out := make([]*Result, len(qs))
	for i, q := range qs {
		names, err := q.Names()
		if err != nil {
			return nil, err
		}
		p := make([]float64, len(names))
		for j := range p {
			p[j] = 0.4 / float64(len(p)-1)
		}
		p[0] = 0.6
		out[i] = newResult(names, p, p, 1, 10, 0)
	}
	return out, nil
}

func TestBothOrders(t *testing.T) {
	f := &firstFamily{}
	d := &Decider{family: f, bothOrders: true}
	rs, err := d.DecideMany("s", []Question{Choice("q", "a", "b", "c"), Noul("n", "", "")}, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(f.calls) != 1 || len(f.calls[0]) != 3 || f.calls[0][2].Options[0].Name != "c" {
		t.Fatalf("calls %+v", f.calls)
	}
	if len(rs) != 2 {
		t.Fatalf("%d results", len(rs))
	}
	want := []float64{0.4, 0.2, 0.4}
	for i, p := range rs[0].Probabilities {
		if math.Abs(p-want[i]) > 1e-12 {
			t.Errorf("choice probabilities %v, want %v", rs[0].Probabilities, want)
		}
	}
	if rs[0].Scores != nil || rs[0].InputTokens != 20 || rs[0].Confidence == 0 {
		t.Errorf("choice result %+v", rs[0])
	}
	if rs[1].Probability("false") != 0.6 || rs[1].InputTokens != 10 {
		t.Errorf("noul result %+v", rs[1])
	}

	if _, err := d.Decide("s", Choice("q", "a", "b"), ""); err != nil {
		t.Fatal(err)
	}
	if last := f.calls[len(f.calls)-1]; len(last) != 2 {
		t.Errorf("Decide read %d questions, want 2", len(last))
	}
}
