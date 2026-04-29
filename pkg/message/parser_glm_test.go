package message

import (
	"testing"
)

func TestParseGLMToolCalls_Single(t *testing.T) {
	response := "get_weather<arg_key>location</arg_key><arg_value>NYC</arg_value>"

	calls := parseGLMToolCalls(response)

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

func TestParseGLMToolCalls_MultipleArguments(t *testing.T) {
	response := "search<arg_key>query</arg_key><arg_value>Go language</arg_value><arg_key>limit</arg_key><arg_value>10</arg_value>"

	calls := parseGLMToolCalls(response)

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

func TestParseGLMToolCalls_MultipleCalls(t *testing.T) {
	response := "get_weather<arg_key>location</arg_key><arg_value>NYC</arg_value>\n" +
		"get_time<arg_key>timezone</arg_key><arg_value>UTC</arg_value>"

	calls := parseGLMToolCalls(response)

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

func TestParseGLMToolCalls_NoMatch(t *testing.T) {
	response := "The answer is 42."

	calls := parseGLMToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 calls, got %d", len(calls))
	}
}

func TestParseGLMToolCalls_SkipsEmptyLines(t *testing.T) {
	response := "\nget_weather<arg_key>location</arg_key><arg_value>Paris</arg_value>\n\n"

	calls := parseGLMToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["location"] != "Paris" {
		t.Errorf("location: got %q, want %q", calls[0].Function.Arguments["location"], "Paris")
	}
}

func TestParseGLMToolCalls_ViaParseToolCall(t *testing.T) {
	// Verify parseToolCalls routes correctly to the GLM parser.
	content := "add<arg_key>a</arg_key><arg_value>15</arg_value><arg_key>b</arg_key><arg_value>27</arg_value>"

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
