package template

import (
	"github.com/ardanlabs/jinja"
	"github.com/hybridgroup/yzma/pkg/message"
)

// Apply applies a jinja chat template to a slice of [message.Message], Set addAssistantPrompt to true to generate the
// assistant prompt, for example on the first message.
func Apply(tmpl string, messages []message.Message, addAssistantPrompt bool) (string, error) {
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
	})
}
