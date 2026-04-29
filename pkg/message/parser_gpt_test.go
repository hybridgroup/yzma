package message

import (
	"testing"
)

func TestParseGPTToolCalls_Single(t *testing.T) {
	content := `.get_weather <|message|>{"location": "NYC"}`

	calls := parseGPTToolCalls(content)

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

func TestParseGPTToolCalls_MultipleArguments(t *testing.T) {
	content := `.search <|message|>{"query": "Go language", "limit": 10}`

	calls := parseGPTToolCalls(content)

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

func TestParseGPTToolCalls_MultipleCalls(t *testing.T) {
	content := `.get_weather <|message|>{"location": "NYC"}.get_time <|message|>{"timezone": "UTC"}`

	calls := parseGPTToolCalls(content)

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

func TestParseGPTToolCalls_BoolArgument(t *testing.T) {
	content := `.set_flag <|message|>{"enabled": true}`

	calls := parseGPTToolCalls(content)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["enabled"] != "true" {
		t.Errorf("enabled: got %q, want %q", calls[0].Function.Arguments["enabled"], "true")
	}
}

func TestParseGPTToolCalls_EmptyArgs(t *testing.T) {
	content := `.get_time <|message|>{}`

	calls := parseGPTToolCalls(content)

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

func TestParseGPTToolCalls_NoMatch(t *testing.T) {
	content := "The answer is 42."

	calls := parseGPTToolCalls(content)

	if len(calls) != 0 {
		t.Fatalf("expected 0 calls, got %d", len(calls))
	}
}

func TestParseGPTToolCalls_ViaParseToolCalls(t *testing.T) {
	content := `.add <|message|>{"a": "15", "b": "27"}`

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

func TestStripMarkup_GPT_Single(t *testing.T) {
	s := `.get_weather <|message|>{"location": "NYC"}The weather is nice.`
	got := StripMarkup(s)
	if got != "The weather is nice." {
		t.Errorf("got %q, want %q", got, "The weather is nice.")
	}
}

func TestStripMarkup_GPT_Multiple(t *testing.T) {
	s := `.get_weather <|message|>{"location": "NYC"}.get_time <|message|>{"tz": "UTC"}Done.`
	got := StripMarkup(s)
	if got != "Done." {
		t.Errorf("got %q, want %q", got, "Done.")
	}
}

func TestStripMarkup_GPT_OnlyToolCall(t *testing.T) {
	s := `.get_weather <|message|>{"location": "NYC"}`
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}
