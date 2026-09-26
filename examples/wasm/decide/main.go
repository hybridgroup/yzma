//go:build js && wasm

// Decide runs a System One model in a browser with the exp/decide package. The
// model answers a typed question about a state with calibrated probabilities.
//
// Build it with TinyGo.
//
//	tinygo build -target wasm -o build/wasm/yzma-decide.wasm ./examples/wasm/decide
//
// Or build it with the standard toolchain.
//
//	GOOS=js GOARCH=wasm go build -o build/wasm/yzma-decide.wasm ./examples/wasm/decide
//
// See wasm/README.md for the method to serve the result.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"
	"time"

	"github.com/hybridgroup/yzma/exp/decide"
	"github.com/hybridgroup/yzma/pkg/llamawasm"
)

const modelPath = "/models/decide.gguf"

var decider *decide.Decider

func main() {
	if err := llamawasm.Load(""); err != nil {
		post("error", err.Error())
		return
	}

	llamawasm.LogSet(llamawasm.LogSilent())
	llamawasm.Init()

	// The page calls these.
	js.Global().Set("yzmaDecideLoad", js.FuncOf(load))
	js.Global().Set("yzmaDecideOpen", js.FuncOf(openModel))
	js.Global().Set("yzmaDecide", js.FuncOf(decideOne))
	js.Global().Set("yzmaDecideMany", js.FuncOf(decideMany))

	post("ready", backendReport())

	// Keep the program alive so that the page can call into it.
	<-make(chan struct{})
}

// load(modelURL, configURL, readout, manyMode) gets a model and its config
// over the network. readout is "jev" or "jevk5", and manyMode is "exact" or
// "batched".
func load(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		post("error", "yzmaDecideLoad needs a model URL, a config URL and a readout")
		return nil
	}
	modelURL, configURL, readout := args[0].String(), args[1].String(), args[2].String()
	mode := arg(args, 3)

	go func() {
		post("status", "downloading the config")
		config, err := fetchText(configURL)
		if err != nil {
			post("error", err.Error())
			return
		}

		post("status", "downloading the model")
		err = llamawasm.FetchModelFile(modelPath, modelURL, func(done, total int64) {
			if total > 0 {
				post("progress", fmt.Sprintf("%d%%", done*100/total))
			}
		})
		if err != nil {
			post("error", err.Error())
			return
		}

		open(modelPath, config, readout, mode)
	}()

	return nil
}

// openModel(path, configJSON, readout, manyMode) loads a model that is already
// in the file system of the llama.cpp module. A test puts the file there itself.
func openModel(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		post("error", "yzmaDecideOpen needs a model path, the config JSON and a readout")
		return nil
	}
	path, config, readout := args[0].String(), args[1].String(), args[2].String()
	mode := arg(args, 3)

	go open(path, config, readout, mode)

	return nil
}

// open loads the model at path with its config.
func open(path, config, readout, mode string) {
	post("status", "loading the model")

	if decider != nil {
		decider.Close()
		decider = nil
	}

	opts := decide.Options{}
	if mode == "batched" {
		opts.ManyMode = decide.ManyBatched
	}

	var err error
	switch readout {
	case "jev":
		var cfg *decide.Config
		if cfg, err = decide.ParseConfig([]byte(config)); err == nil {
			decider, err = decide.NewFromConfig(path, cfg, opts)
		}
	case "jevk5":
		var cfg *decide.JevK5Config
		if cfg, err = decide.ParseJevK5Config([]byte(config)); err == nil {
			decider, err = decide.NewJevK5FromConfig(path, cfg, opts)
		}
	default:
		err = fmt.Errorf("unknown readout %q, want jev or jevk5", readout)
	}
	if err != nil {
		post("error", err.Error())
		return
	}

	post("loaded", backendReport())
}

// decideOne(state, questionJSON, category) scores one question and posts the result.
func decideOne(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		post("error", "yzmaDecide needs a state and a question")
		return nil
	}
	state, question, category := args[0].String(), args[1].String(), arg(args, 2)

	go func() {
		if decider == nil {
			post("error", "load a model first")
			return
		}
		q, err := parseQuestion([]byte(question))
		if err != nil {
			post("error", err.Error())
			return
		}

		start := time.Now()
		res, err := decider.Decide(parseState(state), q, category)
		if err != nil {
			post("error", err.Error())
			return
		}
		postResult(res, time.Since(start))
	}()

	return nil
}

// decideMany(state, questionsJSON, category) scores a JSON list of questions
// about one state and posts the results.
func decideMany(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		post("error", "yzmaDecideMany needs a state and a list of questions")
		return nil
	}
	state, questions, category := args[0].String(), args[1].String(), arg(args, 2)

	go func() {
		if decider == nil {
			post("error", "load a model first")
			return
		}
		var raw []json.RawMessage
		if err := json.Unmarshal([]byte(questions), &raw); err != nil {
			post("error", "questions must be a JSON list: "+err.Error())
			return
		}
		qs := make([]decide.Question, len(raw))
		for i, r := range raw {
			var err error
			if qs[i], err = parseQuestion(r); err != nil {
				post("error", fmt.Sprintf("question %d: %v", i, err))
				return
			}
		}

		start := time.Now()
		res, err := decider.DecideMany(parseState(state), qs, category)
		if err != nil {
			post("error", err.Error())
			return
		}
		postResult(res, time.Since(start))
	}()

	return nil
}

// parseState keeps a JSON object or list as it is written, so its key order
// stays, and uses any other text as a string.
func parseState(s string) any {
	t := bytes.TrimSpace([]byte(s))
	if len(t) > 0 && (t[0] == '{' || t[0] == '[') && json.Valid(t) {
		return json.RawMessage(t)
	}
	return s
}

// parseQuestion reads {"type": "choice", "question": "...", "options": ...}.
// The options are a list of names, an object of name to description, a list
// of levels for score, or an object of "false" and "true" for noul. The type
// is choice when there are options and noul when there are none.
func parseQuestion(data []byte) (decide.Question, error) {
	var in struct {
		Type     string          `json:"type"`
		Question string          `json:"question"`
		Options  json.RawMessage `json:"options"`
	}
	if err := json.Unmarshal(data, &in); err != nil {
		return decide.Question{}, fmt.Errorf("question: %w", err)
	}

	t := decide.Type(in.Type)
	if t == "" {
		t = decide.TypeNoul
		if len(in.Options) > 0 {
			t = decide.TypeChoice
		}
	}

	var list []string
	var named []decide.Option
	if len(in.Options) > 0 {
		var err error
		if list, named, err = parseOptions(in.Options); err != nil {
			return decide.Question{}, err
		}
	}

	switch t {
	case decide.TypeChoice:
		if list != nil {
			return decide.Choice(in.Question, list...), nil
		}
		return decide.ChoiceDesc(in.Question, named...), nil
	case decide.TypeScore:
		if list == nil {
			return decide.Question{}, errors.New("score options must be a list of level descriptions")
		}
		return decide.Score(in.Question, list...), nil
	case decide.TypeNoul:
		if list != nil {
			return decide.Question{}, errors.New(`noul options must be an object {"false": "...", "true": "..."}`)
		}
		return decide.Question{Type: decide.TypeNoul, Text: in.Question, Options: named}, nil
	}

	return decide.Question{}, fmt.Errorf("unknown question type %q", t)
}

// parseOptions reads a JSON list of names or a JSON object of name to
// description. Object keys are kept in the order they are written.
func parseOptions(data []byte) ([]string, []decide.Option, error) {
	var list []string
	if err := json.Unmarshal(data, &list); err == nil {
		return list, nil, nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
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

// fetchText gets a URL with fetch and gives its body as text.
func fetchText(url string) (string, error) {
	type reply struct {
		text string
		err  error
	}
	done := make(chan reply, 1)

	var onText, onResponse, onError js.Func
	release := func() {
		onText.Release()
		onResponse.Release()
		onError.Release()
	}
	onError = js.FuncOf(func(this js.Value, args []js.Value) any {
		msg := "fetch failed"
		if len(args) > 0 {
			msg = args[0].Call("toString").String()
		}
		done <- reply{err: fmt.Errorf("%s: %s", url, msg)}
		return nil
	})
	onText = js.FuncOf(func(this js.Value, args []js.Value) any {
		done <- reply{text: args[0].String()}
		return nil
	})
	onResponse = js.FuncOf(func(this js.Value, args []js.Value) any {
		r := args[0]
		if !r.Get("ok").Bool() {
			done <- reply{err: fmt.Errorf("%s: status %d", url, r.Get("status").Int())}
			return nil
		}
		r.Call("text").Call("then", onText).Call("catch", onError)
		return nil
	})

	js.Global().Call("fetch", url).Call("then", onResponse).Call("catch", onError)
	r := <-done
	release()

	return r.text, r.err
}

// postResult sends a result, or a list of them, as JSON.
func postResult(res any, elapsed time.Duration) {
	out, err := json.Marshal(map[string]any{"result": res, "ms": elapsed.Milliseconds()})
	if err != nil {
		post("error", err.Error())
		return
	}
	post("result", string(out))
}

// backendReport gives the name of the backend that computes.
func backendReport() string {
	if device := llamawasm.GPUDevice(); device != "" {
		return fmt.Sprintf("backend: %s (%s)", llamawasm.Backend(), device)
	}
	return fmt.Sprintf("backend: %s, %d threads", llamawasm.Backend(), llamawasm.Threads())
}

func arg(args []js.Value, i int) string {
	if len(args) > i && args[i].Truthy() {
		return args[i].String()
	}
	return ""
}

// post sends a message to the container of this module. In a worker that is
// the page, and in Node it is the console.
func post(kind, text string) {
	message := map[string]any{"kind": kind, "text": text}

	if fn := js.Global().Get("postMessage"); fn.Type() == js.TypeFunction {
		fn.Invoke(message)
		return
	}
	if fn := js.Global().Get("yzmaOnMessage"); fn.Type() == js.TypeFunction {
		fn.Invoke(message)
		return
	}
	fmt.Printf("%s: %s\n", kind, text)
}
