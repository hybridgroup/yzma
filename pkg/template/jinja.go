package template

import (
	"embed"
	"strings"

	"github.com/ardanlabs/jinja"
	"github.com/hybridgroup/yzma/pkg/message"
)

//go:embed prompts/*.jinja
var builtinTemplates embed.FS

// BuiltinTemplate returns the raw Jinja template string for a named built-in
// template (e.g. "gemma3", "chatml", "qwen2.5-instruct"). The second return
// value is false when no built-in with that name exists.
func BuiltinTemplate(name string) (string, bool) {
	data, err := builtinTemplates.ReadFile("prompts/" + strings.TrimSuffix(name, ".jinja") + ".jinja")
	if err != nil {
		return "", false
	}
	return string(data), true
}

// Options controls optional template rendering behaviour.
type Options struct {
	// EnableThinking controls the enable_thinking template variable passed to
	// models that support a thinking/reasoning mode (e.g. Qwen3). When false
	// the model is instructed to skip its internal chain-of-thought and respond
	// directly. Models whose templates do not use this variable ignore it.
	// Defaults to true (thinking enabled) so existing callers are unaffected.
	EnableThinking bool
}

// DefaultOptions returns Options with all fields set to their defaults.
func DefaultOptions() Options {
	return Options{
		EnableThinking: true,
	}
}

// Apply applies a jinja chat template to a slice of [message.Message].
// Set addAssistantPrompt to true to generate the assistant prompt.
// Thinking mode is enabled by default; use ApplyWithOptions to disable it.
func Apply(tmpl string, messages []message.Message, addAssistantPrompt bool) (string, error) {
	return ApplyWithOptions(tmpl, messages, addAssistantPrompt, DefaultOptions())
}

// ApplyWithOptions is like Apply but accepts an Options struct to control
// template rendering behaviour such as thinking mode.
func ApplyWithOptions(tmpl string, messages []message.Message, addAssistantPrompt bool, opts Options) (string, error) {
	t, err := jinja.Compile(tmpl)
	if err != nil {
		return "", err
	}

	msgs := make([]any, len(messages))
	for i, m := range messages {
		msg := map[string]any{
			"role": m.GetRole(),
		}
		for k, v := range m.GetContent() {
			msg[k] = v
		}
		msgs[i] = msg
	}

	return t.Render(map[string]any{
		"messages":              msgs,
		"add_generation_prompt": addAssistantPrompt,
		"enable_thinking":       opts.EnableThinking,
	})
}
