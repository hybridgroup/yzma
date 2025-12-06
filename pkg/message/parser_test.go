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
	response := `<tool_call>
{"name": "test", "arguments": {}}`

	calls := ParseToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 tool calls for missing closing tag, got %d", len(calls))
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
