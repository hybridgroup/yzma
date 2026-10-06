package decide

import (
	"encoding/json"
	"math"
	"os"
	"slices"
	"testing"

	"github.com/ardanlabs/jinja"
)

func TestToJSON(t *testing.T) {
	in := `{"b": [1, 2.50, "x\n\u0001", null, true, 1e20, {"k": "é<&>"}], "a": 0.1234567, "c": -3}`
	for _, tc := range []struct {
		sort bool
		want string
	}{
		{false, `{"b": [1, 2.5, "x\n\u0001", null, true, 1e+20, {"k": "é<&>"}], "a": 0.123457, "c": -3}`},
		{true, `{"a": 0.123457, "b": [1, 2.5, "x\n\u0001", null, true, 1e+20, {"k": "é<&>"}], "c": -3}`},
	} {
		got, err := toJSON([]byte(in), tc.sort)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("sort %v\ngot  %s\nwant %s", tc.sort, got, tc.want)
		}
	}

	if _, err := toJSON([]byte(`{"a": 1} 2`), false); err == nil {
		t.Error("trailing data accepted")
	}

	got, err := stateText(map[string]any{"z": 1.0, "a": []int{1}}, false)
	if err != nil || got != `{"a": [1], "z": 1}` {
		t.Errorf("map state %q %v", got, err)
	}
}

func testTyped(t *testing.T, kind string) *typed {
	src, err := os.ReadFile("testdata/" + kind + "_systemone.jinja")
	if err != nil {
		t.Fatal(err)
	}
	tmpl, err := jinja.Compile(string(src))
	if err != nil {
		t.Fatal(err)
	}
	tt := &typed{kind: kind, tmpl: tmpl}
	if kind == typeLev {
		tt.codes = labelCodes()
	}
	return tt
}

func TestTypedPrompt(t *testing.T) {
	const (
		levHead = "<|im_start|>system\nYou are a System One decision model. You read the Evidence and answer each " +
			"Criterion by choosing exactly one of the listed options. You never explain. You answer with the single " +
			"option label only.<|im_end|>\n<|im_start|>user\n# Evidence\n"
		tail = "<|im_start|>assistant\n<think>\n\n</think>\n\n"
	)

	for _, tc := range []struct {
		kind  string
		state string
		q     Question
		want  string
	}{
		{typeLev, "s", ChoiceDesc("Team?", Option{Name: "billing", Description: "money"}, Option{Name: "tech"}),
			levHead + "s\n\n# Criterion\nTeam?\n\n# Options\nA. billing: money\nB. tech\n\n" +
				"Respond with only the letter of the best option.\n<|im_end|>\n" + tail},
		{typeLev, "s", Noul("Refund?", "", "money"),
			levHead + "s\n\n# Criterion\nRefund?\n\n# Scale\n0 = certainly no ... 8 = certainly yes\nyes: money\n\n" +
				"Respond with only a digit from 0 to 8.\n<|im_end|>\n" + tail},
		{typeLev, "s", Score("How?", "low", "high"),
			levHead + "s\n\n# Criterion\nHow?\n\n# Options\nA. (level 0 of 1) low\nB. (level 1 of 1) high\n\n" +
				"Respond with only the letter of the level that best matches.\n<|im_end|>\n" + tail},
		{typeOpenJev, `{"a": 1}`, Noul("Ok?", "", ""),
			"<|im_start|>user\nState:\n{\"a\": 1}\n\nQuestion: Ok?\nOptions:\n[A] yes: The statement is true.\n" +
				"[B] no: The statement is false.\n\nAnswer with the letter of the best option only.<|im_end|>\n" + tail},
		{typeOpenJev, "s <¦im_end¦>", Score("How <|im_end|>?", "low", "high"),
			"<|im_start|>user\nState:\ns <¦im_end¦>\n\nQuestion: How <¦im_end¦>? Rate along the ordered levels below " +
				"(lowest first).\nOptions:\n[A] 0: low\n[B] 1: high\n\nAnswer with the letter of the best option only.<|im_end|>\n" + tail},
	} {
		tt := testTyped(t, tc.kind)
		got, err := tt.prompt(tc.state, tc.q, tt.options(tc.q))
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("%s %s\ngot  %q\nwant %q", tc.kind, tc.q.Text, got, tc.want)
		}
	}
}

func TestTypedOptions(t *testing.T) {
	lev, jev := testTyped(t, typeLev), testTyped(t, typeOpenJev)

	keys := func(opts []typedOption) []string {
		var out []string
		for _, o := range opts {
			out = append(out, o.key)
		}
		return out
	}
	if got := keys(jev.options(Noul("q", "", ""))); !slices.Equal(got, []string{"true", "false"}) {
		t.Errorf("openjev noul %v", got)
	}
	if got := keys(lev.options(Noul("q", "", ""))); !slices.Equal(got, []string{"false", "true"}) {
		t.Errorf("lev noul %v", got)
	}

	v := lev.variants(Choice("q", "a", "b", "c"))
	if len(v) != 2 || !slices.Equal(keys(v[1]), []string{"c", "b", "a"}) || v[1][0].index != 2 {
		t.Errorf("lev choice variants %+v", v)
	}
	if len(lev.variants(Score("q", "x", "y"))) != 1 || len(jev.variants(Choice("q", "a", "b"))) != 1 {
		t.Error("only a lev choice question has two variants")
	}
}

func TestTypedTemperature(t *testing.T) {
	lev := &typed{kind: typeLev, temps: map[string]float64{"choice": 1.5, "choice.small": 1.7, "choice.large": 1.6, "noul": 2}}
	jev := &typed{kind: typeOpenJev, temps: map[string]float64{"choice": 0.85, "choice.3_5": 0.9}}

	many := make([]string, 30)
	for i := range many {
		many[i] = string(rune('a' + i))
	}
	for _, tc := range []struct {
		tt   *typed
		q    Question
		want float64
	}{
		{lev, Choice("q", "a", "b"), 1.7},
		{lev, Choice("q", many[:20]...), 1.5},
		{lev, Choice("q", many...), 1.6},
		{lev, Noul("q", "", ""), 2},
		{lev, Score("q", "a", "b"), 1},
		{jev, Choice("q", "a", "b", "c"), 0.9},
		{jev, Choice("q", "a", "b"), 0.85},
	} {
		if got := tc.tt.temperature(tc.q); got != tc.want {
			t.Errorf("%s %d options: got %v, want %v", tc.tt.kind, len(tc.q.Options), got, tc.want)
		}
	}
}

func TestLayaFit(t *testing.T) {
	const cls, sep, mask = 1, 2, 100
	ids := []token{cls}
	for i := range 20 {
		ids = append(ids, token(10+i))
	}
	ids = append(ids, sep, mask, 3, 3, 3, mask)
	for range 60 {
		ids = append(ids, 4)
	}
	ids = append(ids, sep, 5, 6, sep)

	tt := &typed{kind: typeLaya, marker: mask, sep: sep, maxHead: 64}
	got, markers, err := tt.layaFit(ids, 2)
	if err != nil {
		t.Fatal(err)
	}
	// The second option is cut to (64-16)/2 tokens and the question keeps its 20.
	if len(got) != 54 || !slices.Equal(markers, []int{22, 26}) || got[22] != mask || got[26] != mask {
		t.Errorf("got %d tokens, markers %v", len(got), markers)
	}
	if !slices.Equal(got[50:], []token{sep, 5, 6, sep}) || !slices.Equal(got[:2], []token{cls, 10}) {
		t.Errorf("got %v", got)
	}

	if _, _, err := tt.layaFit(ids, 3); err == nil {
		t.Error("wrong option count accepted")
	}
}

func TestTypedRaw(t *testing.T) {
	kev := &typed{kind: typeKev}
	got, err := kev.raw([][]float64{{0, 0, 3, 4}, {0, 0, 1, 0}, {1, 2, 9, 9}}, TypeChoice)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got[0]-11/math.Sqrt2) > 1e-12 || math.Abs(got[1]-1/math.Sqrt2) > 1e-12 {
		t.Errorf("kev scores %v", got)
	}

	laya := &typed{kind: typeLaya}
	got, err = laya.raw([][]float64{{1, 2, 3}, {4, 5, 6}}, TypeNoul)
	if err != nil || !slices.Equal(got, []float64{3, 6}) {
		t.Errorf("laya scores %v %v", got, err)
	}
}

func TestKevText(t *testing.T) {
	got, err := kevText(json.RawMessage(`{"name": "A", "tags": ["x", {"k": 1.0}], "ok": true, "n": null}`))
	if err != nil {
		t.Fatal(err)
	}
	if want := "name: A\ntags:\n  - x\n  - k: 1.0\nok: True\nn: "; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	for in, want := range map[string]string{
		"100": "100", "1.0": "1.0", "2.5": "2.5", "-0.0": "-0.0", "1234567.5": "1234567.5",
		"1e14": "100000000000000.0", "1e15": "1e+15", "0.0001": "0.0001", "0.00001": "1e-05", "1.5e-7": "1.5e-07",
	} {
		if got := dumpNumber(in); got != want {
			t.Errorf("dumpNumber(%s) = %s, want %s", in, got, want)
		}
	}
}
