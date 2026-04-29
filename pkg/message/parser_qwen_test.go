package message

import (
	"testing"
)

func TestParseQwenToolCalls_Single(t *testing.T) {
	response := "<function=get_weather>\n<parameter=location>\nNYC\n</parameter>\n</function>"

	calls := parseQwenToolCalls(response)

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

func TestParseQwenToolCalls_MultipleParameters(t *testing.T) {
	response := "<function=search>\n<parameter=query>\nGo language\n</parameter>\n<parameter=limit>\n10\n</parameter>\n</function>"

	calls := parseQwenToolCalls(response)

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

func TestParseQwenToolCalls_MultipleCalls(t *testing.T) {
	response := "<function=get_weather>\n<parameter=location>\nNYC\n</parameter>\n</function>" +
		"<function=get_time>\n<parameter=timezone>\nUTC\n</parameter>\n</function>"

	calls := parseQwenToolCalls(response)

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

func TestParseQwenToolCalls_NoParameters(t *testing.T) {
	response := "<function=get_time>\n</function>"

	calls := parseQwenToolCalls(response)

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

func TestParseQwenToolCalls_NoMatch(t *testing.T) {
	response := "The answer is 42."

	calls := parseQwenToolCalls(response)

	if len(calls) != 0 {
		t.Fatalf("expected 0 calls, got %d", len(calls))
	}
}

func TestParseQwenToolCalls_ViaParseTool(t *testing.T) {
	// Verify parseToolCalls routes correctly to the Qwen parser.
	content := "<function=add>\n<parameter=a>\n15\n</parameter>\n<parameter=b>\n27\n</parameter>\n</function>"

	calls := parseToolCalls(content)

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
