package decide

import (
	"encoding/json"
	"errors"
	"math"
	"slices"
	"testing"
)

// Values from the /v1/systemone example in the llama-server README.
func TestConfidence(t *testing.T) {
	route := &Result{Probabilities: []float64{0.9998, 0.0001, 0.0001}}
	route.typed(TypeChoice)
	if math.Abs(route.Confidence-0.9997) > 1e-4 {
		t.Errorf("choice confidence %v, want 0.9997", route.Confidence)
	}

	urgency := &Result{Probabilities: []float64{0.0023, 0.116, 0.6753, 0.2064}}
	urgency.typed(TypeScore)
	if math.Abs(urgency.Confidence-0.673) > 1e-3 {
		t.Errorf("score confidence %v, want 0.673", urgency.Confidence)
	}
	if math.Abs(urgency.Expected-2.0858) > 1e-3 {
		t.Errorf("expected level %v, want 2.0858", urgency.Expected)
	}

	tie := &Result{Probabilities: []float64{0.5, 0.5}}
	tie.typed(TypeNoul)
	if tie.Confidence != 0 || tie.Expected != 0 {
		t.Errorf("tie: confidence %v expected %v", tie.Confidence, tie.Expected)
	}
}

func TestParseRequest(t *testing.T) {
	req, err := ParseRequest([]byte(`{
		"state": {"z": 1, "a": "x"},
		"questions": {
			"team": {"type": "choice", "instructions": "Which team?",
				"criteria": {"shipping": null, "billing": "money", "other": {"k": 1}}},
			"angry": {"type": "noul", "instructions": {"q": "angry?"}},
			"urgency": {"type": "score", "instructions": "How urgent?", "criteria": ["later", "now"]}
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}

	if s, ok := req.State.(json.RawMessage); !ok || string(s) != `{"z": 1, "a": "x"}` {
		t.Errorf("state %#v", req.State)
	}
	ids := []string{req.Questions[0].ID, req.Questions[1].ID, req.Questions[2].ID}
	if !slices.Equal(ids, []string{"team", "angry", "urgency"}) {
		t.Errorf("question order %v", ids)
	}

	team := req.Questions[0].Question
	want := []Option{{"shipping", ""}, {"billing", "money"}, {"other", `{"k":1}`}}
	if team.Type != TypeChoice || !slices.Equal(team.Options, want) {
		t.Errorf("team %+v", team)
	}
	if angry := req.Questions[1].Question; angry.Type != TypeNoul || angry.Text != `{"q":"angry?"}` {
		t.Errorf("angry %+v", angry)
	}
	if urgency := req.Questions[2].Question; urgency.Type != TypeScore || urgency.Options[1].Description != "now" {
		t.Errorf("urgency %+v", urgency)
	}

	req, err = ParseRequest([]byte(`{"state": "text", "questions": {"q": {"type": "noul", "instructions": "ok?"}}}`))
	if err != nil {
		t.Fatal(err)
	}
	if req.State != "text" {
		t.Errorf("string state %#v", req.State)
	}
}

func TestParseRequestErrors(t *testing.T) {
	for _, body := range []string{
		`{"questions": {"q": {"type": "noul", "instructions": "ok?"}}}`,
		`{"state": "s", "questions": {}}`,
		`{"state": "s", "questions": []}`,
		`{"state": "s", "images": ["data:image/png;base64,AA"], "questions": {"q": {"type": "noul", "instructions": "ok?"}}}`,
		`{"state": "s", "questions": {"q": {"type": "noul"}}}`,
		`{"state": "s", "questions": {"q": {"type": "maybe", "instructions": "ok?"}}}`,
		`{"state": "s", "questions": {"q": {"type": "choice", "instructions": "ok?", "criteria": ["a"]}}}`,
		`{"state": "s", "questions": {"q": {"type": "score", "instructions": "ok?", "criteria": ["a"]}}}`,
		`{"state": "s", "questions": {"q": {"type": "noul", "instructions": "ok?", "criteria": [1]}}}`,
	} {
		if _, err := ParseRequest([]byte(body)); !errors.Is(err, ErrQuestion) {
			t.Errorf("%s: got %v, want ErrQuestion", body, err)
		}
	}
}

func TestResponseJSON(t *testing.T) {
	choice := ChoiceDesc("Which?", Option{Name: "b"}, Option{Name: "a"})
	score := Score("How?", "low", "high")
	noul := Noul("Ok?", "", "")

	rc := &Result{Answer: "b", Options: []string{"b", "a"}, Probabilities: []float64{0.75, 0.25}, InputTokens: 3}
	rc.typed(TypeChoice)
	rs := &Result{Answer: "1", Options: []string{"0", "1"}, Probabilities: []float64{0.25, 0.75}, InputTokens: 4}
	rs.typed(TypeScore)
	rn := &Result{Answer: "true", Options: []string{"false", "true"}, Probabilities: []float64{0.25, 0.75}, InputTokens: 5}
	rn.typed(TypeNoul)

	resp := &Response{
		Answers: []Answer{{ID: "z", Result: rc, q: choice}, {ID: "s", Result: rs, q: score}, {ID: "n", Result: rn, q: noul}},
		Usage:   Usage{InputTokens: 12},
	}
	got, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"answers":{` +
		`"z":{"type":"choice","choice":"b","probabilities":{"b":0.75,"a":0.25},"confidence":0.5},` +
		`"s":{"type":"score","score":0.75,"legend":{"0":"low","1":"high"},"probabilities":{"0":0.25,"1":0.75},"confidence":0.5},` +
		`"n":{"type":"noul","noul":0.75}},` +
		`"usage":{"input_tokens":12,"output_tokens":0}}`
	if string(got) != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
}
