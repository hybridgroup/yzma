package template

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hybridgroup/yzma/pkg/message"
)

func testTools() []message.ToolDefinition {
	return []message.ToolDefinition{
		{
			Type: "function",
			Function: message.ToolFunctionDefinition{
				Name:        "get_weather",
				Description: "Get the weather of a place",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{"type": "string"},
					},
					"required": []interface{}{"location"},
				},
			},
		},
	}
}

func qwenTemplate(t *testing.T) string {
	t.Helper()

	tmpl, ok := BuiltinTemplate("qwen2.5-instruct")
	if !ok {
		t.Fatal("no qwen2.5-instruct template")
	}

	return tmpl
}

func TestApplyWithToolsRendersTools(t *testing.T) {
	msgs := []message.Message{message.Chat{Role: "user", Content: "weather in Paris?"}}

	got, err := ApplyWithTools(qwenTemplate(t), msgs, testTools(), true)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"<tools>", "get_weather", "location"} {
		if !strings.Contains(got, want) {
			t.Errorf("no %q in the prompt:\n%s", want, got)
		}
	}
}

// A template must take its plain path when there are no tools.
func TestApplyWithoutToolsHasNoToolsBlock(t *testing.T) {
	msgs := []message.Message{message.Chat{Role: "user", Content: "weather in Paris?"}}
	tmpl := qwenTemplate(t)

	got, err := Apply(tmpl, msgs, true)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(got, "<tools>") {
		t.Errorf("a tools block with no tools:\n%s", got)
	}

	// An empty slice must give the same result as none at all.
	empty, err := ApplyWithTools(tmpl, msgs, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	if empty != got {
		t.Errorf("nil tools changed the prompt:\n%s\n\n%s", got, empty)
	}
}

// The template writes each tool with tojson, thus the result must parse and
// must keep the fields of the definition.
func TestToolsContextKeepsFields(t *testing.T) {
	got, err := toolsContext(testTools())
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("got %d tools, want 1", len(got))
	}

	raw, err := json.Marshal(got[0])
	if err != nil {
		t.Fatal(err)
	}

	var tool struct {
		Type     string `json:"type"`
		Function struct {
			Name        string                 `json:"name"`
			Description string                 `json:"description"`
			Parameters  map[string]interface{} `json:"parameters"`
		} `json:"function"`
	}
	if err := json.Unmarshal(raw, &tool); err != nil {
		t.Fatal(err)
	}

	if tool.Type != "function" {
		t.Errorf("type = %q, want function", tool.Type)
	}
	if tool.Function.Name != "get_weather" {
		t.Errorf("name = %q, want get_weather", tool.Function.Name)
	}
	if tool.Function.Parameters["type"] != "object" {
		t.Errorf("no parameters in %v", tool.Function.Parameters)
	}
}
