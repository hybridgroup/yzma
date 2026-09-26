package decide

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Type is the kind of answer a question expects.
type Type string

const (
	// TypeChoice picks one option from a list.
	TypeChoice Type = "choice"
	// TypeScore picks one of 2 to 10 ordered levels.
	TypeScore Type = "score"
	// TypeNoul decides if a statement is false or true.
	TypeNoul Type = "noul"
)

// ErrQuestion is returned for a malformed question.
var ErrQuestion = errors.New("invalid question")

// Option is one answer a question can have.
// Description is optional for choice and noul questions.
type Option struct {
	Name        string
	Description string
}

// Question is a typed question about a state. Options are judged in order.
type Question struct {
	Type    Type
	Text    string
	Options []Option
}

// Choice returns a choice question with options that have no description.
func Choice(text string, names ...string) Question {
	opts := make([]Option, len(names))
	for i, n := range names {
		opts[i] = Option{Name: n}
	}
	return Question{Type: TypeChoice, Text: text, Options: opts}
}

// ChoiceDesc returns a choice question with described options.
func ChoiceDesc(text string, opts ...Option) Question {
	return Question{Type: TypeChoice, Text: text, Options: opts}
}

// Score returns a score question. levels describe level 0 first.
func Score(text string, levels ...string) Question {
	opts := make([]Option, len(levels))
	for i, l := range levels {
		opts[i] = Option{Name: strconv.Itoa(i), Description: l}
	}
	return Question{Type: TypeScore, Text: text, Options: opts}
}

// Noul returns a true or false question. Empty descriptions use the defaults.
func Noul(text, falseDesc, trueDesc string) Question {
	return Question{Type: TypeNoul, Text: text, Options: []Option{
		{Name: "false", Description: falseDesc},
		{Name: "true", Description: trueDesc},
	}}
}

// Names returns the option names in the order the probabilities are returned.
func (q Question) Names() ([]string, error) {
	if strings.TrimSpace(q.Text) == "" {
		return nil, fmt.Errorf("%w: question text must not be empty", ErrQuestion)
	}

	switch q.Type {
	case TypeChoice:
		if len(q.Options) == 0 {
			return nil, fmt.Errorf("%w: choice needs at least one option", ErrQuestion)
		}
		names := make([]string, len(q.Options))
		seen := make(map[string]bool, len(q.Options))
		for i, o := range q.Options {
			if o.Name == "" {
				return nil, fmt.Errorf("%w: option %d has no name", ErrQuestion, i)
			}
			if seen[o.Name] {
				return nil, fmt.Errorf("%w: duplicate option name %q", ErrQuestion, o.Name)
			}
			seen[o.Name] = true
			names[i] = o.Name
		}
		return names, nil
	case TypeScore:
		if len(q.Options) < 2 || len(q.Options) > 10 {
			return nil, fmt.Errorf("%w: score needs 2 to 10 levels, got %d", ErrQuestion, len(q.Options))
		}
		names := make([]string, len(q.Options))
		for i := range q.Options {
			names[i] = strconv.Itoa(i)
		}
		return names, nil
	case TypeNoul:
		for _, o := range q.Options {
			if o.Name != "false" && o.Name != "true" {
				return nil, fmt.Errorf("%w: noul options must be named false or true, got %q", ErrQuestion, o.Name)
			}
		}
		return []string{"false", "true"}, nil
	}

	return nil, fmt.Errorf("%w: unknown question type %q", ErrQuestion, q.Type)
}

// renderOptions returns the option lines as the model was trained on them.
func (q Question) renderOptions() []string {
	switch q.Type {
	case TypeChoice:
		out := make([]string, len(q.Options))
		for i, o := range q.Options {
			if o.Description == "" {
				out[i] = o.Name
			} else {
				out[i] = o.Name + ": " + o.Description
			}
		}
		return out
	case TypeScore:
		out := make([]string, len(q.Options))
		for i, o := range q.Options {
			out[i] = "level " + strconv.Itoa(i) + ": " + o.Description
		}
		return out
	}

	falseDesc, trueDesc := "no, the statement does not hold", "yes, the statement holds"
	for _, o := range q.Options {
		if o.Description == "" {
			continue
		}
		if o.Name == "false" {
			falseDesc = o.Description
		} else {
			trueDesc = o.Description
		}
	}
	return []string{"false: " + falseDesc, "true: " + trueDesc}
}
