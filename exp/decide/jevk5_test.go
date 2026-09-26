package decide

import (
	"encoding/json"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// The expected values in these tests come from the jevk5 v0.3.3 reference, prompt.py.

func TestJevK5Prompt(t *testing.T) {
	k := &jevK5{}
	q := ChoiceDesc("Team?", Option{Name: "billing", Description: "Payments"}, Option{Name: "tech"})
	opts, _, err := k.options(q)
	if err != nil {
		t.Fatal(err)
	}
	var texts []string
	for _, o := range opts {
		texts = append(texts, o.text)
	}

	got, err := jevK5Prompt(json.RawMessage(`{"b":[1,"x,y"],"a":"<&> "}`), q.Text, texts)
	if err != nil {
		t.Fatal(err)
	}
	want := "<|im_start|>system\nApply the supplied criterion to the supplied evidence. Choose exactly one listed option. " +
		"Respond with only its uppercase letter, with no explanation or reasoning.<|im_end|>\n<|im_start|>user\n" +
		`{"evidence": {"b": [1, "x,y"], "a": "<&>` + " " + `"}, "criterion": "Team?", "options": [{"letter": "A", "description": "billing: Payments"}, {"letter": "B", "description": "tech: tech"}]}` +
		"<|im_end|>\n<|im_start|>assistant\n<think>\n\n</think>\n\n"
	if got != want {
		t.Errorf("prompt = %q\nwant %q", got, want)
	}
}

func TestJevK5Options(t *testing.T) {
	k := &jevK5{}
	for _, tc := range []struct {
		q     Question
		texts []string
		index []int
	}{
		{Noul("x", "no rocket", ""), []string{"true: The proposition is true.", "false: no rocket"}, []int{1, 0}},
		{Score("x", "low", "high"), []string{"0: low", "1: high"}, []int{0, 1}},
		{Choice("x", "a", "b"), []string{"a: a", "b: b"}, []int{0, 1}},
	} {
		opts, _, err := k.options(tc.q)
		if err != nil {
			t.Fatal(err)
		}
		var texts []string
		var index []int
		for _, o := range opts {
			texts = append(texts, o.text)
			index = append(index, o.index)
		}
		if !slices.Equal(texts, tc.texts) || !slices.Equal(index, tc.index) {
			t.Errorf("%s: got %q %v, want %q %v", tc.q.Type, texts, index, tc.texts, tc.index)
		}
	}
}

func TestKnockout(t *testing.T) {
	// A deterministic reader: softmax of sin(len(text)*1.7 + option number).
	read := func(passes [][]string) ([][]float64, error) {
		out := make([][]float64, len(passes))
		for i, texts := range passes {
			z := make([]float64, len(texts))
			for j, s := range texts {
				n, _ := strconv.Atoi(s[1:strings.Index(s, ":")])
				z[j] = math.Sin(float64(len(s))*1.7 + float64(n))
			}
			out[i], _ = softmax(z, 1)
		}
		return out, nil
	}

	for _, tc := range []struct {
		n    int
		want []float64
	}{
		{20, []float64{0.029576467543508008, 0.07401516775292828, 0.011766780607160455, 0.13073181756285912, 0.01060535280738332}},
		{40, []float64{0.013491139033847839, 0.035712312143171256, 0.005367350678997194, 0.0630780908506991, 0.004823113231802839}},
		{300, []float64{0.0020505515522615246, 0.005131509261658866, 0.0008157968899984412, 0.008300468157334397, 0.000733076895834617}},
	} {
		texts := make([]string, tc.n)
		for i := range texts {
			texts[i] = "o" + strconv.Itoa(i) + ": option " + strings.Repeat("x", i%5)
		}
		p, err := knockout(read, texts)
		if err != nil {
			t.Fatal(err)
		}
		p = sharpen(p, 0.77)
		for i, w := range tc.want {
			if math.Abs(p[i]-w) > 1e-12 {
				t.Errorf("n %d option %d: got %v, want %v", tc.n, i, p[i], w)
			}
		}
	}
}

func TestGroups(t *testing.T) {
	got := groups(40, 3)
	want := [][2]int{{0, 14}, {14, 27}, {27, 40}}
	if !slices.Equal(got, want) {
		t.Errorf("groups(40, 3) = %v, want %v", got, want)
	}
}

func TestLoadJevK5Config(t *testing.T) {
	dir := t.TempDir()
	for _, tc := range []struct {
		data     string
		ok       bool
		knockout float64
	}{
		{`{"temperature": 1.42}`, true, 0.77},
		{`{"temperature": 1.22, "knockout_temperature": 0.93}`, true, 0.93},
		{`{}`, false, 0},
		{`{"temperature": 1, "knockout_temperature": -1}`, false, 0},
	} {
		path := dir + "/jevk5_config.json"
		if err := os.WriteFile(path, []byte(tc.data), 0o600); err != nil {
			t.Fatal(err)
		}
		c, err := LoadJevK5Config(path)
		if tc.ok != (err == nil) {
			t.Errorf("%s: error %v", tc.data, err)
			continue
		}
		if err == nil && c.KnockoutTemperature != tc.knockout {
			t.Errorf("%s: knockout temperature %v, want %v", tc.data, c.KnockoutTemperature, tc.knockout)
		}
	}
}

var jevK5Intents = []string{"play music", "set alarm", "weather", "news", "calendar", "email", "timer",
	"lights on", "lights off", "call contact", "send text", "navigation", "traffic", "recipe",
	"shopping list", "reminder", "volume up", "volume down", "joke", "translate"}
