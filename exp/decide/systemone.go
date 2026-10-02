package decide

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Request is a TypeSafe /v1/systemone request. Questions keep the order of the
// JSON, and so do the options of a choice question.
type Request struct {
	State     any
	Questions []NamedQuestion
}

// NamedQuestion is a question with its id from the request.
type NamedQuestion struct {
	ID string
	Question
}

// Response is a TypeSafe /v1/systemone response. Answers are in request order.
type Response struct {
	Answers []Answer
	Usage   Usage
}

// Usage is the token count of a response. OutputTokens is always 0.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Answer is the answer to one question in the TypeSafe shape.
type Answer struct {
	ID     string
	Result *Result
	q      Question
}

// ParseRequest parses a TypeSafe /v1/systemone request body.
func ParseRequest(data []byte) (*Request, error) {
	var body struct {
		State     json.RawMessage `json:"state"`
		Images    json.RawMessage `json:"images"`
		Questions json.RawMessage `json:"questions"`
	}
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrQuestion, err)
	}
	if isNull(body.State) {
		return nil, fmt.Errorf("%w: state must be provided", ErrQuestion)
	}
	if !isNull(body.Images) && !bytes.Equal(bytes.TrimSpace(body.Images), []byte("[]")) {
		return nil, fmt.Errorf("%w: images are not supported", ErrQuestion)
	}

	req := &Request{State: body.State}
	var s string
	if json.Unmarshal(body.State, &s) == nil {
		req.State = s
	}

	keys, vals, err := orderedObject(body.Questions)
	if err != nil || len(keys) == 0 {
		return nil, fmt.Errorf("%w: questions must be a non-empty object", ErrQuestion)
	}
	for i, id := range keys {
		q, err := parseQuestion(vals[i])
		if err != nil {
			return nil, fmt.Errorf("questions.%s: %w", id, err)
		}
		req.Questions = append(req.Questions, NamedQuestion{ID: id, Question: q})
	}

	return req, nil
}

func parseQuestion(data json.RawMessage) (Question, error) {
	var raw struct {
		Type         Type            `json:"type"`
		Instructions json.RawMessage `json:"instructions"`
		Criteria     json.RawMessage `json:"criteria"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Question{}, fmt.Errorf("%w: must be an object", ErrQuestion)
	}
	if isNull(raw.Instructions) {
		return Question{}, fmt.Errorf("%w: instructions must be provided", ErrQuestion)
	}
	text, err := jsonText(raw.Instructions)
	if err != nil {
		return Question{}, err
	}

	switch raw.Type {
	case TypeChoice:
		keys, vals, err := orderedObject(raw.Criteria)
		if err != nil || len(keys) == 0 {
			return Question{}, fmt.Errorf("%w: criteria must be a non-empty object", ErrQuestion)
		}
		opts := make([]Option, len(keys))
		for i, k := range keys {
			if opts[i].Description, err = jsonText(vals[i]); err != nil {
				return Question{}, err
			}
			opts[i].Name = k
		}
		return ChoiceDesc(text, opts...), nil
	case TypeScore:
		var levels []json.RawMessage
		if err := json.Unmarshal(raw.Criteria, &levels); err != nil || len(levels) < 2 || len(levels) > 10 {
			return Question{}, fmt.Errorf("%w: criteria must be an array of 2 to 10 levels", ErrQuestion)
		}
		descs := make([]string, len(levels))
		for i, l := range levels {
			if descs[i], err = jsonText(l); err != nil {
				return Question{}, err
			}
		}
		return Score(text, descs...), nil
	case TypeNoul:
		var crit map[string]json.RawMessage
		if !isNull(raw.Criteria) {
			if err := json.Unmarshal(raw.Criteria, &crit); err != nil {
				return Question{}, fmt.Errorf("%w: criteria must be an object", ErrQuestion)
			}
		}
		f, err := jsonText(crit["false"])
		if err != nil {
			return Question{}, err
		}
		t, err := jsonText(crit["true"])
		if err != nil {
			return Question{}, err
		}
		return Noul(text, f, t), nil
	}

	return Question{}, fmt.Errorf("%w: type must be one of choice, score, noul", ErrQuestion)
}

// Answer scores every question of req about its state.
func (d *Decider) Answer(req *Request) (*Response, error) {
	if req == nil || len(req.Questions) == 0 {
		return nil, fmt.Errorf("%w: no questions", ErrQuestion)
	}

	qs := make([]Question, len(req.Questions))
	for i, nq := range req.Questions {
		qs[i] = nq.Question
	}
	rs, err := d.DecideMany(req.State, qs, "")
	if err != nil {
		return nil, err
	}

	resp := &Response{Answers: make([]Answer, len(rs))}
	for i, r := range rs {
		resp.Answers[i] = Answer{ID: req.Questions[i].ID, Result: r, q: qs[i]}
		resp.Usage.InputTokens += r.InputTokens
	}
	return resp, nil
}

// MarshalJSON writes answers as an object keyed by question id, in order.
func (r *Response) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(`{"answers":{`)
	for i, a := range r.Answers {
		if i > 0 {
			b.WriteByte(',')
		}
		writeKey(&b, a.ID)
		v, err := a.MarshalJSON()
		if err != nil {
			return nil, err
		}
		b.Write(v)
	}
	b.WriteString(`},"usage":`)
	u, err := json.Marshal(r.Usage)
	if err != nil {
		return nil, err
	}
	b.Write(u)
	b.WriteByte('}')
	return b.Bytes(), nil
}

// MarshalJSON writes the answer fields of its question type.
func (a Answer) MarshalJSON() ([]byte, error) {
	r := a.Result
	if r == nil {
		return nil, errors.New("decide: answer has no result")
	}

	var b bytes.Buffer
	b.WriteString(`{"type":`)
	writeValue(&b, a.q.Type)
	switch a.q.Type {
	case TypeNoul:
		b.WriteString(`,"noul":`)
		writeValue(&b, r.Probability("true"))
	case TypeChoice:
		b.WriteString(`,"choice":`)
		writeValue(&b, r.Answer)
		writeProbabilities(&b, r)
	case TypeScore:
		b.WriteString(`,"score":`)
		writeValue(&b, r.Expected)
		b.WriteString(`,"legend":{`)
		for i, o := range a.q.Options {
			if i > 0 {
				b.WriteByte(',')
			}
			writeKey(&b, strconv.Itoa(i))
			writeValue(&b, o.Description)
		}
		b.WriteByte('}')
		writeProbabilities(&b, r)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func writeProbabilities(b *bytes.Buffer, r *Result) {
	b.WriteString(`,"probabilities":{`)
	for i, n := range r.Options {
		if i > 0 {
			b.WriteByte(',')
		}
		writeKey(b, n)
		writeValue(b, r.Probabilities[i])
	}
	b.WriteString(`},"confidence":`)
	writeValue(b, r.Confidence)
}

func writeKey(b *bytes.Buffer, k string) {
	writeValue(b, k)
	b.WriteByte(':')
}

// writeValue writes a string or a finite number, which json.Marshal cannot fail on.
func writeValue(b *bytes.Buffer, v any) {
	data, _ := json.Marshal(v)
	b.Write(data)
}

// orderedObject returns the keys and values of a JSON object in order.
func orderedObject(data json.RawMessage) ([]string, []json.RawMessage, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	if t, err := dec.Token(); err != nil || t != json.Delim('{') {
		return nil, nil, errors.New("not an object")
	}

	var keys []string
	var vals []json.RawMessage
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, nil, err
		}
		var v json.RawMessage
		if err := dec.Decode(&v); err != nil {
			return nil, nil, err
		}
		keys = append(keys, t.(string))
		vals = append(vals, v)
	}
	return keys, vals, nil
}

// jsonText returns a JSON string as is, null as empty, and any other value as compact JSON.
func jsonText(data json.RawMessage) (string, error) {
	if isNull(data) {
		return "", nil
	}
	var s string
	if json.Unmarshal(data, &s) == nil {
		return s, nil
	}
	var b bytes.Buffer
	if err := json.Compact(&b, data); err != nil {
		return "", fmt.Errorf("%w: %v", ErrQuestion, err)
	}
	return b.String(), nil
}

func isNull(data json.RawMessage) bool {
	d := bytes.TrimSpace(data)
	return len(d) == 0 || bytes.Equal(d, []byte("null"))
}
