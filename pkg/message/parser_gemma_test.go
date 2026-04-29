package message

import (
	"testing"
)

func TestStripMarkup_TurnMarkerOnly(t *testing.T) {
	s := `<turn>modelHello there.`
	got := StripMarkup(s)
	if got != "Hello there." {
		t.Errorf("got %q, want %q", got, "Hello there.")
	}
}

func TestStripMarkup_BareAndDoubledTurn(t *testing.T) {
	s := `<turn><turn>modelHello there.`
	got := StripMarkup(s)
	if got != "Hello there." {
		t.Errorf("got %q, want %q", got, "Hello there.")
	}
}

func TestStripMarkup_ClosingTurnTag(t *testing.T) {
	s := `</turn><turn>model<turn>model<turn>modelI'm functioning at peak efficiency.`
	got := StripMarkup(s)
	if got != "I'm functioning at peak efficiency." {
		t.Errorf("got %q, want %q", got, "I'm functioning at peak efficiency.")
	}
}

func TestStripMarkup_GemmaCallBlocks(t *testing.T) {
	s := `call:toolmovement{command:<|"|>speak<|"|>}Hello there!`
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_GemmaCallBlocksShortToken(t *testing.T) {
	s := `call:toolmovement{command:<">speak<">}Hello there!`
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_MultipleGemmaCallBlocks(t *testing.T) {
	s := `call:toolmovement{command:<">speak<">}call:toolmovement{command:<">look<">,angle:90}Hello!`
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_ChannelDirective(t *testing.T) {
	s := `<channel>speak:Hello!`
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_TurnAndChannelDirective(t *testing.T) {
	s := `<turn><channel>speak:Hello!`
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_ChannelDirectiveWithToolCall(t *testing.T) {
	s := "<turn><channel>speak:Hello!\ncall:toolmovement{command:<\">slowlook<\">,angle:180}"
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_BareSpeak(t *testing.T) {
	// Model emits speak:content without the <channel> wrapper.
	s := "speak:Hello there!"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_BareWait(t *testing.T) {
	s := "wait:Oh, relax!"
	got := StripMarkup(s)
	if got != "Oh, relax!" {
		t.Errorf("got %q, want %q", got, "Oh, relax!")
	}
}

func TestStripMarkup_BareChannelPreservesURL(t *testing.T) {
	// URL schemes must NOT be stripped.
	s := "https://example.com is a great site."
	got := StripMarkup(s)
	if got != s {
		t.Errorf("got %q, want %q", got, s)
	}
}

func TestParseGemmaToolCalls_SimpleCommand(t *testing.T) {
	response := `call:toolmovement{command:<|"|>look<|"|>,angle:90}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Name != "toolmovement" {
		t.Errorf("name: got %q, want %q", calls[0].Function.Name, "toolmovement")
	}
	if calls[0].Function.Arguments["command"] != "look" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "look")
	}
	if calls[0].Function.Arguments["angle"] != "90" {
		t.Errorf("angle: got %q, want %q", calls[0].Function.Arguments["angle"], "90")
	}
}

func TestParseGemmaToolCalls_PipeQuoteVariant(t *testing.T) {
	response := `call:toolmovement{command:<|"|>slowlook<|"|>,angle:45}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "slowlook" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "slowlook")
	}
	if calls[0].Function.Arguments["angle"] != "45" {
		t.Errorf("angle: got %q, want %q", calls[0].Function.Arguments["angle"], "45")
	}
}

func TestParseGemmaToolCalls_NoArgs(t *testing.T) {
	response := `call:toolmovement{command:<|"|>wait<|"|>}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "wait" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "wait")
	}
}

func TestParseGemmaToolCalls_MultipleSequential(t *testing.T) {
	response := `call:toolmovement{command:<|"|>look<|"|>,angle:90}call:toolmovement{command:<|"|>wait<|"|>}`
	calls := ParseToolCalls(response)
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "look" {
		t.Errorf("call[0] command: got %q", calls[0].Function.Arguments["command"])
	}
	if calls[1].Function.Arguments["command"] != "wait" {
		t.Errorf("call[1] command: got %q", calls[1].Function.Arguments["command"])
	}
}

func TestParseGemmaToolCalls_AutoDetectFallback(t *testing.T) {
	// ParseToolCalls should fall back to Gemma parser when no other grammar matches.
	response := `call:toolmovement{command:<|"|>headshake<|"|>}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call via auto-detect, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "headshake" {
		t.Errorf("command: got %q", calls[0].Function.Arguments["command"])
	}
}

// TestParseGemmaToolCalls_ShortToken tests the <"> alternate quote form that
// some Gemma 4 model variants emit instead of <|"|>.
func TestParseGemmaToolCalls_ShortToken(t *testing.T) {
	response := `call:toolmovement{command:<">slowlook<">,angle:90}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "slowlook" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "slowlook")
	}
	if calls[0].Function.Arguments["angle"] != "90" {
		t.Errorf("angle: got %q, want %q", calls[0].Function.Arguments["angle"], "90")
	}
}

func TestParseGemmaToolCalls_ShortTokenMultiple(t *testing.T) {
	// Replicates the exact format observed in real model output.
	response := `call:toolmovement{command:<">speak<">}call:toolmovement{command:<">look<">,angle:90}`
	calls := ParseToolCalls(response)
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "speak" {
		t.Errorf("call[0] command: got %q", calls[0].Function.Arguments["command"])
	}
	if calls[1].Function.Arguments["command"] != "look" {
		t.Errorf("call[1] command: got %q", calls[1].Function.Arguments["command"])
	}
	if calls[1].Function.Arguments["angle"] != "90" {
		t.Errorf("call[1] angle: got %q", calls[1].Function.Arguments["angle"])
	}
}

func TestParseGemmaToolCalls_NestedArgumentsObject(t *testing.T) {
	// Model wraps extra params inside an arguments:{} object.
	// call:tool_movement{command:<|"|>slowlook<|"|>,arguments:{angle:150}}
	response := "call:tool_movement{command:<|\"|>slowlook<|\"|>,arguments:{angle:150}}"
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "slowlook" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "slowlook")
	}
	if calls[0].Function.Arguments["angle"] != "150" {
		t.Errorf("angle: got %q, want %q", calls[0].Function.Arguments["angle"], "150")
	}
}

func TestParseGemmaToolCalls_NestedArgumentsObjectBareValues(t *testing.T) {
	// Same pattern without Gemma quote tokens (bare values).
	response := "call:tool_movement{command:slowlook,arguments:{angle:150}}"
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(calls))
	}
	if calls[0].Function.Arguments["command"] != "slowlook" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "slowlook")
	}
	if calls[0].Function.Arguments["angle"] != "150" {
		t.Errorf("angle: got %q, want %q", calls[0].Function.Arguments["angle"], "150")
	}
}

func TestParseGemmaToolCalls_EmptyArgsRejected(t *testing.T) {
	// call:func{} with empty braces must not produce a tool call.
	response := "call:tool_movement{}"
	calls := ParseToolCalls(response)
	if len(calls) != 0 {
		t.Fatalf("expected 0 calls for empty Gemma args, got %d: %+v", len(calls), calls)
	}
}

func TestParseGemmaToolCalls_EmptyArgsMixedWithValid(t *testing.T) {
	// A bad call:func{} immediately before a valid call should be filtered out.
	response := `call:tool_movement{}call:tool_movement{command:<|"|>speak<|"|>}`
	calls := ParseToolCalls(response)
	if len(calls) != 1 {
		t.Fatalf("expected 1 call, got %d: %+v", len(calls), calls)
	}
	if calls[0].Function.Arguments["command"] != "speak" {
		t.Errorf("command: got %q, want %q", calls[0].Function.Arguments["command"], "speak")
	}
}
