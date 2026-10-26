package message

import (
	"testing"
)

func TestParseLFMToolCalls_Single(t *testing.T) {
	response := `<|tool_call_start|>[get_weather(location="NYC")]<|tool_call_end|>`

	calls := parseLFMToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Type != "function" {
		t.Errorf("type: got %q, want %q", calls[0].Type, "function")
	}
	if calls[0].Function.Name != "get_weather" {
		t.Errorf("name: got %q, want %q", calls[0].Function.Name, "get_weather")
	}
	if calls[0].Function.Arguments["location"] != "NYC" {
		t.Errorf("location: got %q, want %q", calls[0].Function.Arguments["location"], "NYC")
	}
}

func TestParseLFMToolCalls_ArgumentTypes(t *testing.T) {
	response := `<|tool_call_start|>[move(command='look', angle=135, fast=True, slow=False, note=None, ` +
		`text='it\'s\na test', tags=["a", "b"], opts={"x": 1})]<|tool_call_end|>`

	calls := parseLFMToolCalls(response)

	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	want := map[string]string{
		"command": "look",
		"angle":   "135",
		"fast":    "true",
		"slow":    "false",
		"note":    "null",
		"text":    "it's\na test",
		"tags":    `["a", "b"]`,
		"opts":    `{"x": 1}`,
	}
	args := calls[0].Function.Arguments
	if len(args) != len(want) {
		t.Errorf("got %d arguments, want %d: %v", len(args), len(want), args)
	}
	for k, v := range want {
		if args[k] != v {
			t.Errorf("%s: got %q, want %q", k, args[k], v)
		}
	}
}

func TestParseLFMToolCalls_MultipleCalls(t *testing.T) {
	response := `<|tool_call_start|>[a(x="1, 2)"), b(), c(y=2)]<|tool_call_end|>`

	calls := parseLFMToolCalls(response)

	if len(calls) != 3 {
		t.Fatalf("expected 3 calls, got %d: %v", len(calls), calls)
	}
	for i, name := range []string{"a", "b", "c"} {
		if calls[i].Function.Name != name {
			t.Errorf("call %d name: got %q, want %q", i, calls[i].Function.Name, name)
		}
	}
	if calls[0].Function.Arguments["x"] != "1, 2)" {
		t.Errorf("x: got %q", calls[0].Function.Arguments["x"])
	}
	if len(calls[1].Function.Arguments) != 0 {
		t.Errorf("b args: got %v, want none", calls[1].Function.Arguments)
	}
	if calls[2].Function.Arguments["y"] != "2" {
		t.Errorf("y: got %q", calls[2].Function.Arguments["y"])
	}
}

func TestParseLFMToolCalls_MultipleBlocks(t *testing.T) {
	response := `<|tool_call_start|>[a()]<|tool_call_end|>Hello.<|tool_call_start|>[b()]<|tool_call_end|>`

	calls := parseLFMToolCalls(response)

	if len(calls) != 2 || calls[0].Function.Name != "a" || calls[1].Function.Name != "b" {
		t.Fatalf("got %v, want calls a and b", calls)
	}
}

func TestParseLFMToolCalls_Unclosed(t *testing.T) {
	response := `<|tool_call_start|>[get_weather(location="NYC")]`

	calls := parseLFMToolCalls(response)

	if len(calls) != 1 || calls[0].Function.Arguments["location"] != "NYC" {
		t.Fatalf("got %v, want one get_weather call", calls)
	}
}

func TestParseLFMToolCalls_JSON(t *testing.T) {
	tests := []string{
		`<|tool_call_start|>[{"name": "get_weather", "arguments": {"location": "NYC"}}]<|tool_call_end|>`,
		`<|tool_call_start|>{"name": "get_weather", "arguments": {"location": "NYC"}}<|tool_call_end|>`,
	}
	for _, response := range tests {
		calls := parseLFMToolCalls(response)
		if len(calls) != 1 {
			t.Fatalf("%s: expected 1 call, got %d", response, len(calls))
		}
		if calls[0].Function.Name != "get_weather" || calls[0].Function.Arguments["location"] != "NYC" {
			t.Errorf("%s: got %v", response, calls[0])
		}
	}
}

func TestParseLFMToolCalls_NoMatch(t *testing.T) {
	if calls := parseLFMToolCalls("Just some text."); len(calls) != 0 {
		t.Errorf("expected no calls, got %v", calls)
	}
}

func TestParseLFMToolCalls_ViaParseToolCalls(t *testing.T) {
	response := `<|tool_call_start|>[get_weather(location="NYC")]<|tool_call_end|>Checking the weather.`

	calls := ParseToolCalls(response)

	if len(calls) != 1 || calls[0].Function.Name != "get_weather" {
		t.Fatalf("got %v, want one get_weather call", calls)
	}
}

func TestLFMStripMarkup(t *testing.T) {
	response := `<|tool_call_start|>[get_weather(location="NYC")]<|tool_call_end|>Checking the weather.`

	if got := StripMarkup(response); got != "Checking the weather." {
		t.Errorf("StripMarkup: got %q", got)
	}
	if got := TextAfterToolCalls(response); got != "Checking the weather." {
		t.Errorf("TextAfterToolCalls: got %q", got)
	}
}

func TestLFMDetectFormat(t *testing.T) {
	if got := DetectFormat(`<|tool_call_start|>[a()]<|tool_call_end|>`); got != FormatLFM {
		t.Errorf("DetectFormat: got %v, want %v", got, FormatLFM)
	}
	if got := DetectFormatFromPath("models/LFM2.5-350M-Q4_K_M.gguf"); got != FormatLFM {
		t.Errorf("DetectFormatFromPath: got %v, want %v", got, FormatLFM)
	}
}

func TestLFMStopMarkers(t *testing.T) {
	markers := StopMarkersFor(nil, FormatLFM)
	for _, want := range []string{"<|im_start|>", "<|im_end|>", "<|tool_response_start|>"} {
		found := false
		for _, m := range markers {
			if m == want {
				found = true
			}
		}
		if !found {
			t.Errorf("missing stop marker %q", want)
		}
	}
}
