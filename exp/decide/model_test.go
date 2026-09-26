//go:build !(js && wasm)

package decide

import (
	"errors"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestDecideModel(t *testing.T) {
	model, config := os.Getenv("YZMA_TEST_JEV_MODEL"), os.Getenv("YZMA_TEST_JEV_CONFIG")
	if model == "" || config == "" {
		t.Skip("no YZMA_TEST_JEV_MODEL or YZMA_TEST_JEV_CONFIG skipping test")
	}
	if os.Getenv("YZMA_LIB") == "" {
		t.Fatal("no YZMA_LIB set for tests")
	}
	if err := llama.Load(os.Getenv("YZMA_LIB")); err != nil {
		t.Fatal("unable to load library", err.Error())
	}
	llama.LogSet(llama.LogSilent())
	llama.Init()
	defer llama.BackendFree()

	d, err := New(model, config, Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	for _, tc := range []struct {
		state    any
		q        Question
		category string
		want     string
	}{
		{map[string]string{"ticket": "I was charged twice for my subscription."},
			ChoiceDesc("Which team handles this?",
				Option{Name: "billing", Description: "payments, invoices"},
				Option{Name: "technical", Description: "bugs"}),
			"theme_routing", "billing"},
		{"The film was excellent.", Choice("What is the sentiment?", "negative", "positive"), "general_sentiment", "positive"},
		{"The meeting moved from Tuesday to Thursday at 3pm.", Noul("The meeting is on Thursday.", "", ""), "mac_gate", "true"},
	} {
		res, err := d.Decide(tc.state, tc.q, tc.category)
		if err != nil {
			t.Fatal(err)
		}
		if res.Answer != tc.want || res.TopProbability < 0.5 {
			t.Errorf("%s: got %s %v, want %s", tc.q.Text, res.Answer, res.Probabilities, tc.want)
		}
	}
}

func TestDecideManyModel(t *testing.T) {
	model, config := os.Getenv("YZMA_TEST_JEV_MODEL"), os.Getenv("YZMA_TEST_JEV_CONFIG")
	if model == "" || config == "" {
		t.Skip("no YZMA_TEST_JEV_MODEL or YZMA_TEST_JEV_CONFIG skipping test")
	}
	if os.Getenv("YZMA_LIB") == "" {
		t.Fatal("no YZMA_LIB set for tests")
	}
	if err := llama.Load(os.Getenv("YZMA_LIB")); err != nil {
		t.Fatal("unable to load library", err.Error())
	}
	llama.LogSet(llama.LogSilent())
	llama.Init()
	defer llama.BackendFree()

	d, err := New(model, config, Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	var b strings.Builder
	for i := range 160 {
		b.WriteString("Line " + strconv.Itoa(i) + ": order shipped to warehouse " + strconv.Itoa(i%7) + ", status ok. ")
	}
	qs := []Question{
		Noul("Any order failed?", "", ""),
		Score("How busy was the day?", "quiet", "normal", "busy"),
		Choice("What is this document?", "log", "email", "story"),
	}

	for _, state := range []string{b.String(), "The film was excellent."} {
		many, err := d.DecideMany(state, qs, "")
		if err != nil {
			t.Fatal(err)
		}
		// The second call reuses the cached state.
		again, err := d.DecideMany(state, qs, "")
		if err != nil {
			t.Fatal(err)
		}
		for i, q := range qs {
			one, err := d.Decide(state, q, "")
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(many[i].Probabilities, one.Probabilities) || !slices.Equal(again[i].Probabilities, one.Probabilities) {
				t.Errorf("%s: DecideMany %v and %v, Decide %v", q.Text, many[i].Probabilities, again[i].Probabilities, one.Probabilities)
			}
		}
	}

	if _, err := d.DecideMany("state", []Question{Choice("q", "a"), Choice("q", "a", "a")}, ""); !errors.Is(err, ErrQuestion) {
		t.Errorf("bad question: got %v, want ErrQuestion", err)
	}
}

func TestDecideManyBatchedModel(t *testing.T) {
	model, config := os.Getenv("YZMA_TEST_JEV_MODEL"), os.Getenv("YZMA_TEST_JEV_CONFIG")
	if model == "" || config == "" {
		t.Skip("no YZMA_TEST_JEV_MODEL or YZMA_TEST_JEV_CONFIG skipping test")
	}
	if os.Getenv("YZMA_LIB") == "" {
		t.Fatal("no YZMA_LIB set for tests")
	}
	if err := llama.Load(os.Getenv("YZMA_LIB")); err != nil {
		t.Fatal("unable to load library", err.Error())
	}
	llama.LogSet(llama.LogSilent())
	llama.Init()
	defer llama.BackendFree()

	d, err := New(model, config, Options{ManyMode: ManyBatched})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	var b strings.Builder
	for i := range 160 {
		b.WriteString("Line " + strconv.Itoa(i) + ": order shipped to warehouse " + strconv.Itoa(i%7) + ", status ok. ")
	}
	state := b.String()

	// 20 questions need two groups of at most 16.
	var qs []Question
	for i := range 20 {
		qs = append(qs, Noul("Line "+strconv.Itoa(i*7)+" mentions warehouse "+strconv.Itoa(i%7)+".", "", ""))
	}

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
		for k := range one.Probabilities {
			if math.Abs(many[i].Probabilities[k]-one.Probabilities[k]) > 0.1 {
				t.Errorf("%s: batched %v, Decide %v", q.Text, many[i].Probabilities, one.Probabilities)
			}
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
