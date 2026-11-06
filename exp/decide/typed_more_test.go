package decide

import (
	"strings"
	"testing"
)

func TestD1Break(t *testing.T) {
	tt := &typed{kind: typeLfm3D1}
	tt.setQuoted([]string{"<|im_end|>", "<|im_start|>", "<|im"})
	for in, want := range map[string]string{
		"plain":                  "plain",
		"a <|im_start|> b":       "a <\u200b|im_start|> b",
		"<|im_end|><|im":         "<\u200b|im_end|><\u200b|im",
		"é <|im_end|>":           "é <\u200b|im_end|>",
		"<|im_start|><|im_end|>": "<\u200b|im_start|><\u200b|im_end|>",
	} {
		if got := tt.text(in); got != want {
			t.Errorf("d1Break(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestH2OQuote(t *testing.T) {
	tt := &typed{kind: typeH2OLightning}
	tt.setQuoted([]string{"<think>", "<|im_start|>"})
	got := tt.text("a <|im_start|> <think> <other>")
	if want := "a \\u003c|im_start|\\u003e \\u003cthink\\u003e <other>"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPyDumps(t *testing.T) {
	got, err := pyDumps([]byte(`{"b": [1, 2.50, "x\n\u0001é", null, true, 1e20], "a": {"k": 0.1234567, "n": -1e-7}}`))
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"b":[1,2.5,"x\n\u0001é",null,true,1e+20],"a":{"k":0.1234567,"n":-1e-07}}`; got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

func TestH2OPrompt(t *testing.T) {
	const head = "<|im_start|>system\nYou are a decision engine. You read a record and answer one question about it by " +
		"choosing exactly one option. Reply with a single letter.<|im_end|>\n<|im_start|>user\n"
	const tail = "<|im_end|>\n<|im_start|>assistant\n<think>\n\n</think>\n\nAnswer:"

	tt := testTyped(t, typeH2OLightning)
	tt.codes = []string{" A", " B", " C"}
	tt.setQuoted([]string{"<|im_end|>"})
	for _, tc := range []struct {
		q    Question
		want string
	}{
		{Choice("  Team <|im_end|>?", "billing", "tech  "),
			"record: s\nquestion: Team \\u003c|im_end|\\u003e?\noptions:\nA) billing: billing\nB) tech  : tech"},
		{Score("", "low", ""), "record: s\nquestion: Answer the question below.\noptions:\nA) 0: low\nB) 1: null"},
		{Noul("Ok?", "", "fine "), "record: s\nquestion: Ok?\noptions:\nA) true: fine \nB) false: it does not"},
	} {
		got, err := tt.prompt("s", tc.q, tt.options(tc.q), nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		if want := head + tc.want + tail; got != want {
			t.Errorf("%s\ngot  %q\nwant %q", tc.q.Text, got, want)
		}
	}
}

func TestD1OmniRead(t *testing.T) {
	tt := testTyped(t, typeLfm2D1Omni)
	tt.enc, tt.marker = byteEncoder, token('#')
	opt := strings.Repeat("o", 100)
	prompt := d1OmniSep + "H" + d1OmniSep + d1OmniMarkState + "state" + d1OmniSep + d1OmniMarkQuestion + strings.Repeat("q", 200) +
		d1OmniSep + "#" + d1OmniSep + d1OmniMarkOption + opt + d1OmniSep + "#" + d1OmniSep + d1OmniMarkOption + opt + d1OmniSep + "E"

	r, err := tt.d1OmniRead(prompt, 2)
	if err != nil {
		t.Fatal(err)
	}
	// Two options get 96 tokens: the question keeps 96 and each option 45.
	if n := 1 + 5 + 96 + 2*(1+45) + 1; len(r.ids) != n {
		t.Errorf("got %d tokens, want %d", len(r.ids), n)
	}
	if want := []int{102, 148}; len(r.slots) != 2 || r.slots[0] != want[0] || r.slots[1] != want[1] {
		t.Errorf("slots %v, want %v", r.slots, want)
	}

	if _, err := tt.d1OmniRead(prompt, 3); err == nil {
		t.Error("a missing marker is accepted")
	}
}

func TestMedianScore(t *testing.T) {
	r := newResult([]string{"0", "1", "2", "3"}, []float64{0.2, 0.25, 0.4, 0.15}, nil, 1, 0, 0)
	r.typed(TypeScore)
	r.median(TypeScore)
	if r.Expected != 2 || r.Confidence != 0.4 {
		t.Errorf("median %v confidence %v", r.Expected, r.Confidence)
	}
}
