package message

import (
	"testing"
)

func TestParseMistralToolCalls_Single(t *testing.T) {
	response := `[TOOL_CALLS]get_weather[ARGS]{"location": "NYC"}`

	calls := parseMistralToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Name != "get_weather" {
		t.Errorf("name: got %q, want %q", calls[0].Function.Name, "get_weather")
	}
	if calls[0].Function.Arguments["location"] != "NYC" {
		t.Errorf("location: got %q, want %q", calls[0].Function.Arguments["location"], "NYC")
	}
	if calls[0].Type != "function" {
		t.Errorf("type: got %q, want %q", calls[0].Type, "function")
	}
}

func TestParseMistralToolCalls_MultipleArguments(t *testing.T) {
	response := `[TOOL_CALLS]search[ARGS]{"query": "Go language", "limit": 10}`

	calls := parseMistralToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Name != "search" {
		t.Errorf("name: got %q, want %q", calls[0].Function.Name, "search")
	}
	if calls[0].Function.Arguments["query"] != "Go language" {
		t.Errorf("query: got %q, want %q", calls[0].Function.Arguments["query"], "Go language")
	}
	if calls[0].Function.Arguments["limit"] != "10" {
		t.Errorf("limit: got %q, want %q", calls[0].Function.Arguments["limit"], "10")
	}
}

func TestParseMistralToolCalls_MultipleCalls(t *testing.T) {
	response := `[TOOL_CALLS]get_weather[ARGS]{"location": "NYC"}[TOOL_CALLS]get_time[ARGS]{"timezone": "UTC"}`

	calls := parseMistralToolCalls(response)

	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
	if calls[0].Function.Name != "get_weather" {
		t.Errorf("call[0] name: got %q, want %q", calls[0].Function.Name, "get_weather")
	}
	if calls[1].Function.Name != "get_time" {
		t.Errorf("call[1] name: got %q, want %q", calls[1].Function.Name, "get_time")
	}
	if calls[1].Function.Arguments["timezone"] != "UTC" {
		t.Errorf("call[1] timezone: got %q, want %q", calls[1].Function.Arguments["timezone"], "UTC")
	}
}

func TestParseMistralToolCalls_NumericArgument(t *testing.T) {
	response := `[TOOL_CALLS]add[ARGS]{"a": 15, "b": 27}`

	calls := parseMistralToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["a"] != "15" {
		t.Errorf("a: got %q, want %q", calls[0].Function.Arguments["a"], "15")
	}
	if calls[0].Function.Arguments["b"] != "27" {
		t.Errorf("b: got %q, want %q", calls[0].Function.Arguments["b"], "27")
	}
}

func TestParseMistralToolCalls_BoolArgument(t *testing.T) {
	response := `[TOOL_CALLS]set_flag[ARGS]{"enabled": true}`

	calls := parseMistralToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["enabled"] != "true" {
		t.Errorf("enabled: got %q, want %q", calls[0].Function.Arguments["enabled"], "true")
	}
}

func TestParseMistralToolCalls_EmptyArgs(t *testing.T) {
	response := `[TOOL_CALLS]get_time[ARGS]{}`

	calls := parseMistralToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Name != "get_time" {
		t.Errorf("name: got %q, want %q", calls[0].Function.Name, "get_time")
	}
	if len(calls[0].Function.Arguments) != 0 {
		t.Errorf("expected 0 arguments, got %d", len(calls[0].Function.Arguments))
	}
}

func TestParseMistralToolCalls_NoMatch(t *testing.T) {
	response := "The answer is 42."

	calls := parseMistralToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 calls, got %d", len(calls))
	}
}

func TestParseMistralToolCalls_ViaParseToolCall(t *testing.T) {
	// Verify parseToolCalls routes correctly to the Mistral parser.
	content := `[TOOL_CALLS]add[ARGS]{"a": "15", "b": "27"}`

	calls := parseToolCalls(content)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Name != "add" {
		t.Errorf("name: got %q, want %q", calls[0].Function.Name, "add")
	}
	if calls[0].Function.Arguments["a"] != "15" {
		t.Errorf("a: got %q, want %q", calls[0].Function.Arguments["a"], "15")
	}
	if calls[0].Function.Arguments["b"] != "27" {
		t.Errorf("b: got %q, want %q", calls[0].Function.Arguments["b"], "27")
	}
}

func TestFindJSONObjectEnd_Simple(t *testing.T) {
	s := `{"key": "value"} extra`
	got := findJSONObjectEnd(s)
	if got != 16 {
		t.Errorf("got %d, want %d", got, 16)
	}
}

func TestFindJSONObjectEnd_Nested(t *testing.T) {
	s := `{"outer": {"inner": 1}} extra`
	got := findJSONObjectEnd(s)
	if got != 23 {
		t.Errorf("got %d, want %d", got, 23)
	}
}

func TestFindJSONObjectEnd_NoObject(t *testing.T) {
	s := `no object here`
	got := findJSONObjectEnd(s)
	if got != -1 {
		t.Errorf("got %d, want -1", got)
	}
}

func TestFindJSONObjectEnd_EscapedBraceInString(t *testing.T) {
	s := `{"key": "val{ue"} extra`
	got := findJSONObjectEnd(s)
	if got != 17 {
		t.Errorf("got %d, want %d", got, 17)
	}
}
