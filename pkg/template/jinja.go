package template

import (
	"errors"
	"io"

	"github.com/hybridgroup/yzma/pkg/message"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

// Apply applies a jinja chat template to a slice of [message.Message], Set addAssistantPrompt to true to generate the
// assistant prompt, for example on the first message.
func Apply(tmpl string, messages []message.Message, addAssistantPrompt bool) (string, error) {
	// prevent filesystem access
	gonja.DefaultLoader = &NoFSLoader{}

	t, err := gonja.FromString(tmpl)
	if err != nil {
		return "", err
	}

	msgs := make([]map[string]interface{}, len(messages))
	for i, m := range messages {
		msgs[i] = map[string]interface{}{
			"role": m.GetRole(),
		}
		for k, v := range m.GetContent() {
			msgs[i][k] = v
		}
	}

	data := exec.NewContext(map[string]interface{}{
		"messages":              msgs,
		"add_generation_prompt": addAssistantPrompt,
	})

	return t.ExecuteToString(data)
}

// NoFSLoader is a template loader that provides no filesystem access.
// This prevents template injection attacks like {% include "/etc/passwd" %}.
type NoFSLoader struct{}

func (nl *NoFSLoader) Read(path string) (io.Reader, error) {
	return nil, errors.New("filesystem access disabled")
}

// Resolve always returns an error to prevent filesystem access.
func (nl *NoFSLoader) Resolve(path string) (string, error) {
	return "", errors.New("filesystem access disabled")
}

func (nl *NoFSLoader) Inherit(from string) (loaders.Loader, error) {
	return nil, errors.New("filesystem access disabled")
}
