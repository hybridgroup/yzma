package message

import (
	"testing"
)

func TestParseToolCalls_SingleToolCall(t *testing.T) {
	response := `I'll help you with that.
<tool_call>
{"name": "add", "arguments": {"a": 15, "b": 27}}
</tool_call>
Let me calculate that for you.`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	if calls[0].Function.Name != "add" {
		t.Errorf("expected function name 'add', got '%s'", calls[0].Function.Name)
	}

	if calls[0].Function.Arguments["a"] != "15" {
		t.Errorf("expected argument 'a' to be '15', got '%s'", calls[0].Function.Arguments["a"])
	}

	if calls[0].Function.Arguments["b"] != "27" {
		t.Errorf("expected argument 'b' to be '27', got '%s'", calls[0].Function.Arguments["b"])
	}

	if calls[0].Type != "function" {
		t.Errorf("expected type 'function', got '%s'", calls[0].Type)
	}
}

func TestParseToolCalls_MultipleToolCalls(t *testing.T) {
	response := `Let me solve this step by step.
<tool_call>
{"name": "add", "arguments": {"a": 15, "b": 27}}
</tool_call>
Now I'll multiply:
<tool_call>
{"name": "multiply", "arguments": {"a": 42, "b": 3}}
</tool_call>
Done!`

	calls := ParseToolCalls(response)

	if len(calls) != 2 {
		t.Fatalf("expected 2 tool calls, got %d", len(calls))
	}

	if calls[0].Function.Name != "add" {
		t.Errorf("expected first function name 'add', got '%s'", calls[0].Function.Name)
	}

	if calls[1].Function.Name != "multiply" {
		t.Errorf("expected second function name 'multiply', got '%s'", calls[1].Function.Name)
	}

	if calls[1].Function.Arguments["a"] != "42" {
		t.Errorf("expected argument 'a' to be '42', got '%s'", calls[1].Function.Arguments["a"])
	}

	if calls[1].Function.Arguments["b"] != "3" {
		t.Errorf("expected argument 'b' to be '3', got '%s'", calls[1].Function.Arguments["b"])
	}
}

func TestParseToolCalls_NoToolCalls(t *testing.T) {
	response := `The answer is 42. No tools needed for this one.`

	calls := ParseToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 tool calls, got %d", len(calls))
	}
}

func TestParseToolCalls_StringArguments(t *testing.T) {
	response := `<tool_call>
{"name": "search", "arguments": {"query": "weather in Paris", "limit": "10"}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	if calls[0].Function.Name != "search" {
		t.Errorf("expected function name 'search', got '%s'", calls[0].Function.Name)
	}

	if calls[0].Function.Arguments["query"] != "weather in Paris" {
		t.Errorf("expected argument 'query' to be 'weather in Paris', got '%s'", calls[0].Function.Arguments["query"])
	}

	if calls[0].Function.Arguments["limit"] != "10" {
		t.Errorf("expected argument 'limit' to be '10', got '%s'", calls[0].Function.Arguments["limit"])
	}
}

func TestParseToolCalls_MixedArgumentTypes(t *testing.T) {
	response := `<tool_call>
{"name": "process", "arguments": {"name": "test", "count": 5, "ratio": 3.14}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	if calls[0].Function.Arguments["name"] != "test" {
		t.Errorf("expected argument 'name' to be 'test', got '%s'", calls[0].Function.Arguments["name"])
	}

	if calls[0].Function.Arguments["count"] != "5" {
		t.Errorf("expected argument 'count' to be '5', got '%s'", calls[0].Function.Arguments["count"])
	}

	if calls[0].Function.Arguments["ratio"] != "3.14" {
		t.Errorf("expected argument 'ratio' to be '3.14', got '%s'", calls[0].Function.Arguments["ratio"])
	}
}

func TestParseToolCalls_EmptyArguments(t *testing.T) {
	response := `<tool_call>
{"name": "get_time", "arguments": {}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	if calls[0].Function.Name != "get_time" {
		t.Errorf("expected function name 'get_time', got '%s'", calls[0].Function.Name)
	}

	if len(calls[0].Function.Arguments) != 0 {
		t.Errorf("expected 0 arguments, got %d", len(calls[0].Function.Arguments))
	}
}

func TestParseToolCalls_InvalidJSON(t *testing.T) {
	response := `<tool_call>
{"name": "broken", "arguments": {invalid json}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 tool calls for invalid JSON, got %d", len(calls))
	}
}

func TestParseToolCalls_IncompleteToolCall(t *testing.T) {
	response := `<tool_call>
{"name": "add", "arguments": {"a": 1, "b": 2}}
No closing tag here`

	calls := ParseToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 tool calls for incomplete tag, got %d", len(calls))
	}
}

func TestParseToolCalls_NestedObjects(t *testing.T) {
	response := `<tool_call>
{"name": "complex", "arguments": {"data": {"nested": "value"}}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	// Nested objects should be marshaled back to JSON string
	if calls[0].Function.Arguments["data"] != `{"nested":"value"}` {
		t.Errorf("expected nested object as JSON string, got '%s'", calls[0].Function.Arguments["data"])
	}
}

func TestParseToolCalls_ArrayArgument(t *testing.T) {
	response := `<tool_call>
{"name": "batch", "arguments": {"items": [1, 2, 3]}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	// Arrays should be marshaled back to JSON string
	if calls[0].Function.Arguments["items"] != `[1,2,3]` {
		t.Errorf("expected array as JSON string, got '%s'", calls[0].Function.Arguments["items"])
	}
}

func TestParseToolCalls_WhitespaceHandling(t *testing.T) {
	response := `<tool_call>

    {"name": "add", "arguments": {"a": 1, "b": 2}}

</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	if calls[0].Function.Name != "add" {
		t.Errorf("expected function name 'add', got '%s'", calls[0].Function.Name)
	}
}

func TestParseToolCalls_FloatPrecision(t *testing.T) {
	response := `<tool_call>
{"name": "calculate", "arguments": {"value": 123.456789}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	if calls[0].Function.Arguments["value"] != "123.456789" {
		t.Errorf("expected argument 'value' to be '123.456789', got '%s'", calls[0].Function.Arguments["value"])
	}
}

// ---- StripMarkup ----

func TestStripMarkup_ToolCallBlock(t *testing.T) {
	s := `<tool_call>{"name": "tool_movement", "arguments": {"command": "speak"}}</tool_call>Hello!`
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_NoMarkup(t *testing.T) {
	s := "Just a normal sentence."
	got := StripMarkup(s)
	if got != s {
		t.Errorf("got %q, want %q", got, s)
	}
}

// ---- TextAfterToolCalls ----

func TestTextAfterToolCalls_Standard(t *testing.T) {
	// Thinking content precedes the tool call; spoken text follows.
	s := "Thinking Process:1.\nAnalyze the request.\n<tool_call>{\"name\":\"tool_movement\",\"arguments\":{\"command\":\"speak\"}}</tool_call>\nGreetings, human!"
	got := TextAfterToolCalls(s)
	if got != "Greetings, human!" {
		t.Errorf("got %q, want %q", got, "Greetings, human!")
	}
}

func TestTextAfterToolCalls_MultipleStandard(t *testing.T) {
	// Two tool calls; text after the last one should be returned.
	s := "<tool_call>{\"name\":\"tool_movement\",\"arguments\":{\"command\":\"speak\"}}</tool_call>Hello!<tool_call>{\"name\":\"tool_movement\",\"arguments\":{\"command\":\"wait\"}}</tool_call>Goodbye!"
	got := TextAfterToolCalls(s)
	if got != "Goodbye!" {
		t.Errorf("got %q, want %q", got, "Goodbye!")
	}
}

func TestTextAfterToolCalls_NoToolCalls(t *testing.T) {
	// No tool calls — whole string returned.
	s := "Just a normal sentence."
	got := TextAfterToolCalls(s)
	if got != s {
		t.Errorf("got %q, want %q", got, s)
	}
}

func TestTextAfterToolCalls_NoTextAfter(t *testing.T) {
	// Tool call at end — result is empty.
	s := "Thinking...\n<tool_call>{\"name\":\"tool_movement\",\"arguments\":{\"command\":\"speak\"}}</tool_call>"
	got := TextAfterToolCalls(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

func TestTextAfterToolCalls_Qwen(t *testing.T) {
	s := "Thinking...\n<function=tool_movement>{\"command\":\"speak\"}</function>\nHello there!"
	got := TextAfterToolCalls(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestParseToolCalls_NegativeNumbers(t *testing.T) {
	response := `<tool_call>
{"name": "subtract", "arguments": {"a": -10, "b": -5.5}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	if calls[0].Function.Arguments["a"] != "-10" {
		t.Errorf("expected argument 'a' to be '-10', got '%s'", calls[0].Function.Arguments["a"])
	}

	if calls[0].Function.Arguments["b"] != "-5.5" {
		t.Errorf("expected argument 'b' to be '-5.5', got '%s'", calls[0].Function.Arguments["b"])
	}
}

func TestParseToolCalls_EmptyResponse(t *testing.T) {
	response := ``

	calls := ParseToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 tool calls for empty response, got %d", len(calls))
	}
}

func TestParseToolCalls_OnlyOpeningTag(t *testing.T) {
	// Phi-4 and similar models emit <tool_call>{JSON} without a closing tag
	// when generation is halted by a stop token. We must parse these anyway.
	response := "<tool_call>\n{\"name\": \"test\", \"arguments\": {}}"

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call for unclosed tag with valid JSON, got %d", len(calls))
	}
	if calls[0].Function.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", calls[0].Function.Name)
	}
}

func TestParseToolCalls_OnlyClosingTag(t *testing.T) {
	response := `{"name": "test", "arguments": {}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 tool calls for missing opening tag, got %d", len(calls))
	}
}

func TestParseToolCalls_BooleanArguments(t *testing.T) {
	response := `<tool_call>
{"name": "toggle", "arguments": {"enabled": true, "verbose": false}}
</tool_call>`

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}

	if calls[0].Function.Arguments["enabled"] != "true" {
		t.Errorf("expected argument 'enabled' to be 'true', got '%s'", calls[0].Function.Arguments["enabled"])
	}

	if calls[0].Function.Arguments["verbose"] != "false" {
		t.Errorf("expected argument 'verbose' to be 'false', got '%s'", calls[0].Function.Arguments["verbose"])
	}
}

// ---- ParseToolCalls: inline JSON format ({"name":"...","args":{...}}) ----

func TestParseToolCalls_InlineJSON_ArgsField(t *testing.T) {
	// Model emits bare JSON with "args" instead of "arguments".
	response := `{"name": "tool_movement", "args": {"command": "speak"}}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Name != "tool_movement" {
		t.Errorf("name: got %q, want %q", calls[0].Function.Name, "tool_movement")
	}
	if calls[0].Function.Arguments["command"] != "speak" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "speak")
	}
}

func TestParseToolCalls_InlineJSON_WithNumericArg(t *testing.T) {
	response := `{"name": "tool_movement", "args": {"command": "slowlook", "angle": 140}}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "slowlook" {
		t.Errorf("command: got %q", calls[0].Function.Arguments["command"])
	}
	if calls[0].Function.Arguments["angle"] != "140" {
		t.Errorf("angle: got %q, want %q", calls[0].Function.Arguments["angle"], "140")
	}
}

func TestParseToolCalls_InlineJSON_EmbeddedInText(t *testing.T) {
	// Model emits JSON tool call followed by toolresult and <turn> separator then speech.
	response := `{"name": "tool_movement", "args": {"command": "speak"}}<toolresult>{"status": "success"}</toolresult><turn>Hello there!`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "speak" {
		t.Errorf("command: got %q", calls[0].Function.Arguments["command"])
	}
}

func TestParseToolCalls_InlineJSON_ArgumentsField(t *testing.T) {
	// Also handle "arguments" (standard field name) in inline JSON format.
	response := `{"name": "tool_movement", "arguments": {"command": "wait"}}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "wait" {
		t.Errorf("command: got %q", calls[0].Function.Arguments["command"])
	}
}

func TestParseToolCalls_InlineJSON_NoToolCallField(t *testing.T) {
	// JSON with "name" but no "args"/"arguments" must NOT be treated as tool call.
	response := `{"name": "Alice", "age": 30}`
	calls := ParseToolCalls(response)
	if len(calls) != 0 {
		t.Fatalf("expected 0 calls, got %d", len(calls))
	}
}

func TestParseToolCalls_InlineJSON_SpacesInName(t *testing.T) {
	// Names with spaces are natural text, not tool calls.
	response := `{"name": "Alice Smith", "args": {"x": "y"}}`
	calls := ParseToolCalls(response)
	if len(calls) != 0 {
		t.Fatalf("expected 0 calls, got %d", len(calls))
	}
}

func TestParseToolCalls_InlineJSON_EmptyArgsRejected(t *testing.T) {
	// {"args":{}} with no keys must NOT produce a tool call (missing required args).
	response := `{"name": "tool_movement", "args": {}}`
	calls := ParseToolCalls(response)
	if len(calls) != 0 {
		t.Fatalf("expected 0 calls for empty args, got %d", len(calls))
	}
}

// TestParseToolCalls_CorruptedQwenOpener covers the pattern where the model
// emits `{"function=name>` instead of `<function=name>` inside a <tool_call>
// block — i.e. the `<` is corrupted to `{"`.
func TestParseToolCalls_CorruptedQwenOpener(t *testing.T) {
	response := "<tool_call>\n{\"function=tool_movement>\n<parameter=command>\nspeak\n</parameter>\n<parameter=angle>\n100\n</parameter>\n</function>\n</tool_call>"

	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Name != "tool_movement" {
		t.Errorf("name: got %q, want %q", calls[0].Function.Name, "tool_movement")
	}
	if calls[0].Function.Arguments["command"] != "speak" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "speak")
	}
	if calls[0].Function.Arguments["angle"] != "100" {
		t.Errorf("angle: got %q, want %q", calls[0].Function.Arguments["angle"], "100")
	}
}

// TestParseToolCalls_TruncatedJSON covers the pattern emitted by some Qwen
// fine-tune models where the closing `}` of the outer JSON object is dropped,
// causing a standard json.Unmarshal to fail. repairJSON should recover it.
func TestParseToolCalls_TruncatedJSON(t *testing.T) {
	// Missing outer closing brace — exactly what was observed in production logs.
	response := "<tool_call>\n{\"name\": \"tool_movement\", \"arguments\": {\"command\": \"speak\"}}\n</tool_call>"
	// Also test the truly truncated variant (outer `}` absent).
	truncated := "<tool_call>\n{\"name\": \"tool_movement\", \"arguments\": {\"command\": \"speak\"}\n</tool_call>"

	for _, tc := range []struct {
		name  string
		input string
	}{
		{"valid", response},
		{"truncated", truncated},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := ParseToolCalls(tc.input)
			if len(calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(calls))
			}
			if calls[0].Function.Name != "tool_movement" {
				t.Errorf("name: got %q, want %q", calls[0].Function.Name, "tool_movement")
			}
			if calls[0].Function.Arguments["command"] != "speak" {
				t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "speak")
			}
		})
	}
}

// TestParseToolCalls_NestedArguments covers the pattern where some Qwen
// fine-tune models wrap the actual arguments inside an extra "arguments" or
// "parameters" sub-object, e.g.:
//
//	{"name":"tool_movement","arguments":{"command":"look","arguments":{"angle":90}}}
//
// flattenNestedArguments should promote the nested values to the top level.
func TestParseToolCalls_NestedArguments(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCmd   string
		wantAngle string
	}{
		{
			name:      "nested_arguments_key",
			input:     "<tool_call>\n{\"name\": \"tool_movement\", \"arguments\": {\"command\": \"look\", \"arguments\": {\"angle\": 90}}}\n</tool_call>",
			wantCmd:   "look",
			wantAngle: "90",
		},
		{
			name:      "nested_parameters_key",
			input:     "<tool_call>\n{\"name\": \"tool_movement\", \"arguments\": {\"command\": \"look\", \"parameters\": {\"angle\": 45}}}\n</tool_call>",
			wantCmd:   "look",
			wantAngle: "45",
		},
		{
			name:    "truncated_nested",
			input:   "<tool_call>\n{\"name\": \"tool_movement\", \"arguments\": {\"command\": \"speak\", \"arguments\": {\"angle\": 90}}\n</tool_call>",
			wantCmd: "speak",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			calls := ParseToolCalls(tc.input)
			if len(calls) != 1 {
				t.Fatalf("expected 1 call, got %d", len(calls))
			}
			if calls[0].Function.Name != "tool_movement" {
				t.Errorf("name: got %q, want %q", calls[0].Function.Name, "tool_movement")
			}
			if calls[0].Function.Arguments["command"] != tc.wantCmd {
				t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], tc.wantCmd)
			}
			if tc.wantAngle != "" && calls[0].Function.Arguments["angle"] != tc.wantAngle {
				t.Errorf("angle: got %q, want %q", calls[0].Function.Arguments["angle"], tc.wantAngle)
			}
		})
	}
}

func TestParseToolCalls_Phi4_PipeToolCallTag_Unclosed(t *testing.T) {
	// Phi-4 uses <|tool_call>{JSON} without a closing tag.
	response := "Hello, human. I am performing optimally today.<|tool_call>{\"name\": \"slowlook\", \"arguments\": {\"angle\": 90}}"

	calls := ParseToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(calls))
	}
	if calls[0].Function.Name != "slowlook" {
		t.Errorf("expected name %q, got %q", "slowlook", calls[0].Function.Name)
	}
	if calls[0].Function.Arguments["angle"] != "90" {
		t.Errorf("expected angle %q, got %q", "90", calls[0].Function.Arguments["angle"])
	}
}
