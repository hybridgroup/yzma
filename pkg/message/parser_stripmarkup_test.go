package message

import (
	"testing"
)

// ---- StripMarkup: Standard format (<tool_call>) ----

func TestStripMarkup_Standard_ToolCallAtStart(t *testing.T) {
	s := `<tool_call>{"name": "add", "arguments": {"a": 1}}</tool_call>The answer is 42.`
	got := StripMarkup(s)
	if got != "The answer is 42." {
		t.Errorf("got %q, want %q", got, "The answer is 42.")
	}
}

func TestStripMarkup_Standard_ToolCallAtEnd(t *testing.T) {
	s := `Here you go.<tool_call>{"name": "add", "arguments": {"a": 1}}</tool_call>`
	got := StripMarkup(s)
	if got != "Here you go." {
		t.Errorf("got %q, want %q", got, "Here you go.")
	}
}

func TestStripMarkup_Standard_MultipleToolCalls(t *testing.T) {
	s := `Step one.
<tool_call>{"name": "add", "arguments": {}}</tool_call>
Step two.
<tool_call>{"name": "multiply", "arguments": {}}</tool_call>
Done.`
	got := StripMarkup(s)
	want := "Step one.\n\nStep two.\n\nDone."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripMarkup_Standard_ToolcallNoUnderscore(t *testing.T) {
	// <toolcall> (no underscore) should be normalised and removed.
	s := `<toolcall>{"name": "add", "arguments": {}}</toolcall>Hello!`
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_Standard_OnlyToolCall(t *testing.T) {
	s := `<tool_call>{"name": "add", "arguments": {}}</tool_call>`
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

// ---- StripMarkup: Qwen format (<function=…>) ----

func TestStripMarkup_Qwen_SingleBlock(t *testing.T) {
	s := "<function=get_weather>\n<parameter=location>\nNYC\n</parameter>\n</function>It's sunny."
	got := StripMarkup(s)
	if got != "It's sunny." {
		t.Errorf("got %q, want %q", got, "It's sunny.")
	}
}

func TestStripMarkup_Qwen_MultipleBlocks(t *testing.T) {
	s := "<function=get_weather>\n<parameter=location>\nNYC\n</parameter>\n</function>" +
		"<function=get_time>\n<parameter=tz>\nUTC\n</parameter>\n</function>" +
		"Done."
	got := StripMarkup(s)
	if got != "Done." {
		t.Errorf("got %q, want %q", got, "Done.")
	}
}

func TestStripMarkup_Qwen_TextAroundBlock(t *testing.T) {
	s := "Let me check. <function=get_weather>\n<parameter=location>\nNYC\n</parameter>\n</function> All done."
	got := StripMarkup(s)
	if got != "Let me check.  All done." {
		t.Errorf("got %q, want %q", got, "Let me check.  All done.")
	}
}

// ---- StripMarkup: GLM format (<arg_key>) ----

func TestStripMarkup_GLM_SingleLine(t *testing.T) {
	s := "get_weather<arg_key>location</arg_key><arg_value>NYC</arg_value>\nHere is your answer."
	got := StripMarkup(s)
	if got != "Here is your answer." {
		t.Errorf("got %q, want %q", got, "Here is your answer.")
	}
}

func TestStripMarkup_GLM_MultipleLines(t *testing.T) {
	s := "get_weather<arg_key>location</arg_key><arg_value>NYC</arg_value>\n" +
		"get_time<arg_key>tz</arg_key><arg_value>UTC</arg_value>\n" +
		"Here is your answer."
	got := StripMarkup(s)
	if got != "Here is your answer." {
		t.Errorf("got %q, want %q", got, "Here is your answer.")
	}
}

func TestStripMarkup_GLM_OnlyToolCall(t *testing.T) {
	s := "get_weather<arg_key>location</arg_key><arg_value>NYC</arg_value>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

// ---- StripMarkup: Mistral format ([TOOL_CALLS]) ----

func TestStripMarkup_Mistral_Single(t *testing.T) {
	s := `[TOOL_CALLS]get_weather[ARGS]{"location": "NYC"}The weather is nice.`
	got := StripMarkup(s)
	if got != "The weather is nice." {
		t.Errorf("got %q, want %q", got, "The weather is nice.")
	}
}

func TestStripMarkup_Mistral_Multiple(t *testing.T) {
	s := `[TOOL_CALLS]get_weather[ARGS]{"location": "NYC"}[TOOL_CALLS]get_time[ARGS]{"tz": "UTC"}Done.`
	got := StripMarkup(s)
	if got != "Done." {
		t.Errorf("got %q, want %q", got, "Done.")
	}
}

func TestStripMarkup_Mistral_OnlyToolCall(t *testing.T) {
	s := `[TOOL_CALLS]get_weather[ARGS]{"location": "NYC"}`
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

// ---- StripMarkup: plain text passthrough ----

func TestStripMarkup_PlainText(t *testing.T) {
	s := "No tool calls here. Just a normal sentence."
	got := StripMarkup(s)
	if got != s {
		t.Errorf("got %q, want %q", got, s)
	}
}

func TestStripMarkup_Empty(t *testing.T) {
	got := StripMarkup("")
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}
