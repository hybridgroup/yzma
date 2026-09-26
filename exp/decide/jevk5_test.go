package decide

import (
	"encoding/json"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
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

func jevK5TestDecider(t *testing.T, mode ManyMode) *Decider {
	model, config := os.Getenv("YZMA_TEST_JEVK5_MODEL"), os.Getenv("YZMA_TEST_JEVK5_CONFIG")
	if model == "" || config == "" {
		t.Skip("no YZMA_TEST_JEVK5_MODEL or YZMA_TEST_JEVK5_CONFIG skipping test")
	}
	if os.Getenv("YZMA_LIB") == "" {
		t.Fatal("no YZMA_LIB set for tests")
	}
	if err := llama.Load(os.Getenv("YZMA_LIB")); err != nil {
		t.Fatal("unable to load library", err.Error())
	}
	llama.LogSet(llama.LogSilent())
	llama.Init()
	t.Cleanup(llama.BackendFree)

	d, err := NewJevK5(model, config, Options{ManyMode: mode})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(d.Close)
	return d
}

var jevK5Intents = []string{"play music", "set alarm", "weather", "news", "calendar", "email", "timer",
	"lights on", "lights off", "call contact", "send text", "navigation", "traffic", "recipe",
	"shopping list", "reminder", "volume up", "volume down", "joke", "translate"}

func TestJevK5Model(t *testing.T) {
	d := jevK5TestDecider(t, ManyExact)

	for _, tc := range []struct {
		state string
		q     Question
		want  string
	}{
		{"I was billed twice for order #4411. Please refund the duplicate charge today.",
			ChoiceDesc("Which team should handle this?",
				Option{Name: "billing", Description: "Payments and refunds"},
				Option{Name: "tech", Description: "Bugs"},
				Option{Name: "sales", Description: "New purchases"}), "billing"},
		{"Refunds need a receipt and a purchase within 30 days. The customer bought 12 days ago and has no receipt.",
			Noul("Is a refund permitted under the policy?", "", ""), "false"},
		{"Turn the living room lights off.", Choice("What does the user want?", jevK5Intents...), "lights off"},
	} {
		res, err := d.Decide(tc.state, tc.q, "")
		if err != nil {
			t.Fatal(err)
		}
		var sum float64
		for _, p := range res.Probabilities {
			sum += p
		}
		if res.Answer != tc.want || math.Abs(sum-1) > 1e-9 {
			t.Errorf("%s: got %s %v, want %s", tc.q.Text, res.Answer, res.Probabilities, tc.want)
		}
	}
}

func TestJevK5ManyModel(t *testing.T) {
	var b strings.Builder
	for i := range 150 {
		b.WriteString("Ticket " + strconv.Itoa(i) + ": the customer reports a double charge on the invoice. ")
	}
	state := b.String()
	qs := []Question{
		Noul("The customer was charged twice.", "", ""),
		Score("How urgent is this?", "low", "medium", "high"),
		Choice("What does the user want?", jevK5Intents...),
	}

	for _, mode := range []ManyMode{ManyExact, ManyBatched} {
		t.Run(strconv.Itoa(int(mode)), func(t *testing.T) {
			d := jevK5TestDecider(t, mode)
			many, err := d.DecideMany(state, qs, "")
			if err != nil {
				t.Fatal(err)
			}
			again, err := d.DecideMany(state, qs, "")
			if err != nil {
				t.Fatal(err)
			}
			for i, q := range qs {
				one, err := d.Decide(state, q, "")
				if err != nil {
					t.Fatal(err)
				}
				if !slices.Equal(many[i].Probabilities, again[i].Probabilities) {
					t.Errorf("%s: cached call %v, first call %v", q.Text, again[i].Probabilities, many[i].Probabilities)
				}
				// A near tie can change the finalists of a batched knockout,
				// so only its sum is checked there.
				if mode == ManyBatched && len(q.Options) > 16 {
					var sum float64
					for _, p := range many[i].Probabilities {
						sum += p
					}
					if math.Abs(sum-1) > 1e-9 {
						t.Errorf("%s: probabilities sum to %v", q.Text, sum)
					}
					continue
				}
				for k := range one.Probabilities {
					d := math.Abs(many[i].Probabilities[k] - one.Probabilities[k])
					if (mode == ManyExact && d != 0) || d > 0.1 {
						t.Errorf("%s: DecideMany %v, Decide %v", q.Text, many[i].Probabilities, one.Probabilities)
						break
					}
				}
			}
		})
	}
}
