//go:build js && wasm

// Tools calls functions with llama.cpp in a browser.
//
// The program has the same structure as examples/wasm/chat. It renders the chat
// template of the model with the template package, offers the tools to it,
// parses the tool calls that come back, runs them, and gives the results to the
// model for a final answer.
//
// Build it with TinyGo.
//
//	tinygo build -target wasm -o build/wasm/yzma-tools.wasm ./examples/wasm/tools
//
// Or build it with the standard toolchain.
//
//	GOOS=js GOARCH=wasm go build -o build/wasm/yzma-tools.wasm ./examples/wasm/tools
//
// See wasm/README.md for the method to serve the result.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/hybridgroup/yzma/pkg/llamawasm"
	"github.com/hybridgroup/yzma/pkg/message"
	"github.com/hybridgroup/yzma/pkg/template"
)

// maxTurns is the number of times that the model can call a tool before it must
// answer. A turn in a browser is slow, thus this is small.
const maxTurns = 3

const modelPath = "/models/model.gguf"

var (
	model   llamawasm.Model
	ctx     llamawasm.Context
	vocab   llamawasm.Vocab
	sampler llamawasm.Sampler

	// format is how this model writes a tool call. It comes from the name of
	// the file of the model.
	format message.Format
)

func main() {
	if err := llamawasm.Load(""); err != nil {
		post("error", err.Error())
		return
	}

	llamawasm.LogSet(llamawasm.LogSilent())
	llamawasm.Init()

	// The page calls these.
	js.Global().Set("yzmaLoadModel", js.FuncOf(loadModel))
	js.Global().Set("yzmaOpenModel", js.FuncOf(openModel))
	js.Global().Set("yzmaAsk", js.FuncOf(ask))

	post("ready", backendReport())

	// Keep the program alive so that the page can call into it.
	<-make(chan struct{})
}

// loadModel(url) gets a model over the network and makes a context for it.
func loadModel(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		post("error", "loadModel needs a URL")
		return nil
	}
	url := args[0].String()

	go func() {
		post("status", "downloading the model")

		err := llamawasm.FetchModelFile(modelPath, url, func(done, total int64) {
			if total > 0 {
				post("progress", fmt.Sprintf("%d%%", done*100/total))
			}
		})
		if err != nil {
			post("error", err.Error())
			return
		}

		open(modelPath, url)
	}()

	return nil
}

// openModel(path) loads a model that is already in the filesystem of the
// llama.cpp module. A test puts the file there itself.
func openModel(this js.Value, args []js.Value) any {
	path := modelPath
	if len(args) > 0 && args[0].Truthy() {
		path = args[0].String()
	}

	// The second argument gives the name of the original file, because the path
	// in the filesystem of the module says nothing about the model.
	name := path
	if len(args) > 1 && args[1].Truthy() {
		name = args[1].String()
	}

	go open(path, name)

	return nil
}

// open loads the model at path and makes a context for it. The name tells which
// format of tool call to expect.
func open(path, name string) {
	post("status", "loading the model")

	params := llamawasm.ModelDefaultParams()

	// A WebGPU build has a device, thus put each layer on it. A CPU build has no
	// device and ignores this value.
	if llamawasm.GPUDevice() != "" {
		params.NGpuLayers = 999
	}

	var err error
	if model, err = llamawasm.ModelLoadFromFile(path, params); err != nil {
		post("error", err.Error())
		return
	}

	ctxParams := llamawasm.ContextDefaultParams()

	// Each turn repeats the whole conversation and the tools, thus this context
	// is larger than the one of the chat example.
	ctxParams.NCtx = 4096
	ctxParams.NBatch = 512

	if ctx, err = llamawasm.InitFromModel(model, ctxParams); err != nil {
		post("error", err.Error())
		return
	}

	vocab = llamawasm.ModelGetVocab(model)
	format = message.DetectFormatFromPath(name)

	post("loaded", llamawasm.ModelDesc(model)+", "+backendReport())
}

// backendReport gives the name of the backend that computes.
func backendReport() string {
	if device := llamawasm.GPUDevice(); device != "" {
		return fmt.Sprintf("backend: %s (%s)", llamawasm.Backend(), device)
	}
	return fmt.Sprintf("backend: %s, %d threads", llamawasm.Backend(), llamawasm.Threads())
}

// ask(question, maxTokens) answers a question and calls tools if it needs them.
func ask(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		post("error", "ask needs a question")
		return nil
	}
	question := args[0].String()

	maxTokens := int32(256)
	if len(args) > 1 && args[1].Truthy() {
		maxTokens = int32(args[1].Int())
	}

	go converse(question, maxTokens)

	return nil
}

// converse asks the model, runs the tools that it calls, and asks again with
// the results until the model gives an answer.
func converse(question string, maxTokens int32) {
	if model == 0 || ctx == 0 {
		post("error", "load a model first")
		return
	}

	tmpl, tools, messages := start(question)
	began := time.Now()

	var total int32

	for turn := 0; turn < maxTurns; turn++ {
		prompt, err := template.ApplyWithTools(tmpl, messages, tools, true)
		if err != nil {
			post("error", err.Error())
			return
		}

		response, count, err := generate(prompt, maxTokens)
		total += count
		if err != nil {
			post("error", err.Error())
			return
		}

		calls := message.ParseToolCalls(response)
		if len(calls) == 0 {
			post("answer", message.StripMarkup(response))
			break
		}

		messages = append(messages, message.Tool{
			Role:      "assistant",
			Content:   message.StripMarkup(response),
			ToolCalls: calls,
		})

		for _, call := range calls {
			messages = append(messages, run(call))
		}
	}

	elapsed := time.Since(began).Seconds()
	if elapsed > 0 {
		post("done", fmt.Sprintf("%d tokens, %.2f tokens/s", total, float64(total)/elapsed))
		return
	}
	post("done", fmt.Sprintf("%d tokens", total))
}

// start gives the template, the tools, and the first messages of a
// conversation.
func start(question string) (string, []message.ToolDefinition, []message.Message) {
	tools := toolDefinitions()

	tmpl := llamawasm.ModelChatTemplate(model, "")
	if tmpl != "" {
		return tmpl, tools, []message.Message{
			message.Chat{Role: "user", Content: question},
		}
	}

	// A model with no template of its own takes chatml, which has no tools
	// branch. Thus the tools go in a system message instead.
	tmpl, _ = template.BuiltinTemplate("chatml")

	return tmpl, tools, []message.Message{
		message.Chat{Role: "system", Content: toolsPrompt(tools)},
		message.Chat{Role: "user", Content: question},
	}
}

// toolsPrompt describes the tools for a template that cannot show them itself.
func toolsPrompt(tools []message.ToolDefinition) string {
	list, err := json.Marshal(tools)
	if err != nil {
		return ""
	}

	return "You can call these tools:\n" + string(list) +
		"\n\nTo call one, answer with <tool_call>{\"name\": \"...\", \"arguments\": {...}}</tool_call>." +
		" After you get the result, answer the question."
}

// run calls one tool and makes the message that holds its result.
func run(call message.ToolCall) message.Message {
	args, _ := json.Marshal(call.Function.Arguments)
	post("tool", call.Function.Name+" "+string(args))

	result, err := executeToolCall(call)
	if err != nil {
		result = "error: " + err.Error()
	}
	post("result", result)

	return message.ToolResponse{
		Role:    "tool",
		Name:    call.Function.Name,
		Content: result,
	}
}

// generate makes text for a prompt and sends each piece to the page. It stops
// at the end of the turn of the model.
func generate(prompt string, maxTokens int32) (string, int32, error) {
	// Each turn starts from an empty state, because the prompt holds the whole
	// conversation.
	if err := llamawasm.MemoryClear(ctx, true); err != nil {
		return "", 0, err
	}

	if sampler != 0 {
		llamawasm.SamplerFree(sampler)
	}
	sampler = llamawasm.SamplerChainInit(llamawasm.SamplerChainDefaultParams())
	llamawasm.SamplerChainAdd(sampler, llamawasm.SamplerInitGreedy())

	tokens := llamawasm.Tokenize(vocab, prompt, true, true)
	if len(tokens) == 0 {
		return "", 0, fmt.Errorf("the prompt has no tokens")
	}

	markers := message.StopMarkers(vocab, format)
	batch := llamawasm.BatchGetOne(tokens)
	buf := make([]byte, 64)

	var text strings.Builder
	var count int32

	for pos := int32(0); pos < maxTokens; pos += batch.NTokens {
		if _, err := llamawasm.Decode(ctx, batch); err != nil {
			return text.String(), count, err
		}

		token := llamawasm.SamplerSample(sampler, ctx, -1)
		if llamawasm.VocabIsEOG(vocab, token) {
			break
		}
		llamawasm.SamplerAccept(sampler, token)

		if n := llamawasm.TokenToPiece(vocab, token, buf, 0, true); n > 0 {
			piece := string(buf[:n])
			text.WriteString(piece)
			post("token", piece)
		}

		count++

		// A small model can write the next turn itself. Cut the text at the
		// first marker and stop.
		if cut, found := trim(text.String(), markers); found {
			return cut, count, nil
		}

		batch = llamawasm.BatchGetOne([]llamawasm.Token{token})
	}

	return text.String(), count, nil
}

// trim cuts the text at the first marker that it holds.
func trim(text string, markers []string) (string, bool) {
	cut := -1
	for _, marker := range markers {
		if i := strings.Index(text, marker); i >= 0 && (cut < 0 || i < cut) {
			cut = i
		}
	}

	if cut < 0 {
		return text, false
	}

	return text[:cut], true
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
