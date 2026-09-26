package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/exp/decide"
	"github.com/hybridgroup/yzma/pkg/llama"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := handleFlags(); err != nil {
		showUsage()
		return err
	}

	q, err := buildQuestion()
	if err != nil {
		return err
	}

	if err := llama.Load(*libPath); err != nil {
		return fmt.Errorf("unable to load library: %w", err)
	}

	if !*verbose {
		llama.LogSet(llama.LogSilent())
	}

	llama.Init()
	defer llama.Close()

	newDecider := decide.New
	if *readout == "jevk5" {
		newDecider = decide.NewJevK5
	}
	d, err := newDecider(*modelFile, *configFile, decide.Options{Threads: int32(*threads)})
	if err != nil {
		return err
	}
	defer d.Close()

	// A state that is valid JSON is passed as is, so its key order is kept.
	var st any = *state
	if json.Valid([]byte(*state)) && (*state)[0] != '"' {
		st = json.RawMessage(*state)
	}

	res, err := d.Decide(st, q, *category)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	return nil
}

func buildQuestion() (decide.Question, error) {
	t := decide.Type(*qtype)
	if t == "" {
		t = decide.TypeNoul
		if len(*options) > 0 {
			t = decide.TypeChoice
		}
	}

	var list []string
	var named []decide.Option
	if len(*options) > 0 {
		var err error
		if list, named, err = parseOptions(*options); err != nil {
			return decide.Question{}, err
		}
	}

	switch t {
	case decide.TypeChoice:
		if list != nil {
			return decide.Choice(*question, list...), nil
		}
		return decide.ChoiceDesc(*question, named...), nil
	case decide.TypeScore:
		if list == nil {
			return decide.Question{}, errors.New("score options must be a JSON list of level descriptions")
		}
		return decide.Score(*question, list...), nil
	case decide.TypeNoul:
		if list != nil {
			return decide.Question{}, errors.New(`noul options must be a JSON object {"false":"...","true":"..."}`)
		}
		return decide.Question{Type: decide.TypeNoul, Text: *question, Options: named}, nil
	}

	return decide.Question{}, fmt.Errorf("unknown question type %q", t)
}

// parseOptions reads a JSON list of names or a JSON object of name to
// description. Object keys are kept in the order they are written.
func parseOptions(s string) ([]string, []decide.Option, error) {
	var list []string
	if err := json.Unmarshal([]byte(s), &list); err == nil {
		return list, nil, nil
	}

	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') {
		return nil, nil, errors.New("options must be a JSON list or object")
	}

	var opts []decide.Option
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			return nil, nil, err
		}
		var desc string
		if err := dec.Decode(&desc); err != nil {
			return nil, nil, fmt.Errorf("option %v: description must be a string", key)
		}
		opts = append(opts, decide.Option{Name: key.(string), Description: desc})
	}

	return nil, opts, nil
}
