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

// ---- StripMarkup: <think> reasoning blocks ----

func TestStripMarkup_Think_FullBlock(t *testing.T) {
	s := "<think>internal reasoning here</think>The actual answer."
	got := StripMarkup(s)
	if got != "The actual answer." {
		t.Errorf("got %q, want %q", got, "The actual answer.")
	}
}

func TestStripMarkup_Think_BlockWithNewlines(t *testing.T) {
	s := "<think>\nUser asked how I am today.\nI should respond naturally.\n</think>\nHello there!"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_Think_MultipleBlocks(t *testing.T) {
	s := "<think>first thought</think>Step one. <think>second thought</think>Step two."
	got := StripMarkup(s)
	if got != "Step one. Step two." {
		t.Errorf("got %q, want %q", got, "Step one. Step two.")
	}
}

func TestStripMarkup_Think_NoClosingTag(t *testing.T) {
	// Incomplete block: strip from <think> to end of string.
	s := "Here is some text. <think>partial reasoning without closing tag"
	got := StripMarkup(s)
	if got != "Here is some text." {
		t.Errorf("got %q, want %q", got, "Here is some text.")
	}
}

func TestStripMarkup_Think_OnlyBlock(t *testing.T) {
	s := "<think>only reasoning, no spoken content</think>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

func TestStripMarkup_Think_MixedWithToolCall(t *testing.T) {
	// Gemma 4 thinking model: <think> block then tool call then spoken text.
	s := "<think>\nI should move and speak.\n</think>\ncall:tool_movement{command:<|\"|>speak<|\"|>}Hello!"
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_Think_OrphanedCloseTag(t *testing.T) {
	// Qwen3 thinking model: <think> is injected into the generation prompt so
	// the generated text starts mid-thought and ends the block with </think>.
	s := "Thinking Process:1.\nAnalyze the Request: User asked how I am.\n</think>\nI'm doing well today!"
	got := StripMarkup(s)
	if got != "I'm doing well today!" {
		t.Errorf("got %q, want %q", got, "I'm doing well today!")
	}
}

func TestStripMarkup_Think_OrphanedCloseTagWithToolCall(t *testing.T) {
	// Qwen3 thinking model with tool calls: orphaned </think> then Qwen tool
	// call then spoken text.
	s := "Thinking Process:1.\nStep A.\n</think>\n<function=tool_movement>{\"command\":\"speak\",\"angle\":90}</function>\nHello there!"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_Think_OrphanedCloseTagOnly(t *testing.T) {
	// Orphaned </think> with no text after — model only thought, no spoken response.
	s := "Reasoning...\n</think>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

// ---- StripMarkup: Gemma 4 sentence markers (<s>, </s>) ----

func TestStripMarkup_ChatML_ImStartStripped(t *testing.T) {
	// Qwen/ChatML: model starts simulating next turn with <|im_start|>.
	// Everything from that token onwards should be discarded.
	s := "I was merely responding with standard efficiency.\n<|im_start|>user\nhow are you"
	got := StripMarkup(s)
	if got != "I was merely responding with standard efficiency." {
		t.Errorf("got %q, want %q", got, "I was merely responding with standard efficiency.")
	}
}

func TestStripMarkup_ChatML_ImStartDecodedStripped(t *testing.T) {
	// Decoded form (pipes stripped by TokenToPiece).
	s := "Hello there!<im_start>T\nSome fabricated turn."
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_Gemma4_SentenceMarkers(t *testing.T) {
	s := "<s>Hello there!</s>"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_Gemma4_SentenceMarkersWithTurn(t *testing.T) {
	s := "First sentence.</s><turn>model<s>Second sentence."
	got := StripMarkup(s)
	if got != "First sentence. Second sentence." {
		t.Errorf("got %q, want %q", got, "First sentence. Second sentence.")
	}
}

// ---- StripMarkup: Gemma 4 <toolresponse> blocks ----

func TestStripMarkup_Gemma4_ToolResponse_Single(t *testing.T) {
	s := "<toolresponse>speak()</toolresponse>Hello!"
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_Gemma4_ToolResponse_Multiple(t *testing.T) {
	s := "<toolresponse>speak()</toolresponse>Hello.<toolresponse>stop()</toolresponse>"
	got := StripMarkup(s)
	if got != "Hello." {
		t.Errorf("got %q, want %q", got, "Hello.")
	}
}

func TestStripMarkup_Gemma4_ToolResponse_NoClosingTag(t *testing.T) {
	s := "Some text. <toolresponse>speak()"
	got := StripMarkup(s)
	if got != "Some text." {
		t.Errorf("got %q, want %q", got, "Some text.")
	}
}

// ---- StripMarkup: Gemma 4 <toolresult> simulated result blocks ----

func TestStripMarkup_Gemma4_ToolResult_Single(t *testing.T) {
	s := `<toolresult>{"status": "success"}</toolresult>Hello!`
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_Gemma4_ToolResult_Multiple(t *testing.T) {
	s := `<toolresult>{"status": "success"}</toolresult><toolresult>{"status": "success"}</toolresult>I am doing well.`
	got := StripMarkup(s)
	if got != "I am doing well." {
		t.Errorf("got %q, want %q", got, "I am doing well.")
	}
}

func TestStripMarkup_Gemma4_ToolResult_NoClosingTag(t *testing.T) {
	s := `Some text. <toolresult>{"status": "success"}`
	got := StripMarkup(s)
	if got != "Some text." {
		t.Errorf("got %q, want %q", got, "Some text.")
	}
}

// ---- StripMarkup: Gemma 4 <|turn> canonical token form ----

func TestStripMarkup_Gemma4_PipeTurnTag(t *testing.T) {
	s := "<|turn>model<s>Hello there!</s>"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_Gemma4_PipeTurnClosing(t *testing.T) {
	s := "I am well, thank you.<turn|>"
	got := StripMarkup(s)
	if got != "I am well, thank you." {
		t.Errorf("got %q, want %q", got, "I am well, thank you.")
	}
}

func TestStripMarkup_Gemma4_FullTurnSequence(t *testing.T) {
	// Realistic Gemma 4 output: toolresponse + sentence markers + turn markers.
	s := "<toolresponse>speak()</toolresponse><|turn>model<s>I am functioning well.</s><turn|>"
	got := StripMarkup(s)
	if got != "I am functioning well." {
		t.Errorf("got %q, want %q", got, "I am functioning well.")
	}
}

// ---- StripMarkup: Gemma 4 pipe-delimited <|toolcall> wrappers ----

func TestStripMarkup_Gemma4_PipeToolcallOpening(t *testing.T) {
	// Model uses <|toolcall> as opening and a second <|toolcall> as self-close.
	s := "<|toolcall>call:tool_movement{command:<|\"|>speak<|\"|>}<|toolcall>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

func TestStripMarkup_Gemma4_PipeToolcallOpenClose(t *testing.T) {
	// Model uses <|toolcall> opening and <toolcall|> closing.
	s := "<|toolcall>call:tool_movement{command:<|\"|>speak<|\"|>}<toolcall|>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

func TestStripMarkup_Gemma4_PipeToolcallWithSpokenText(t *testing.T) {
	// Spoken text alongside a pipe-delimited tool call block.
	s := "<|toolcall>call:tool_movement{command:<|\"|>speak<|\"|>}<toolcall|>Hello there!"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

// ---- StripMarkup: Gemma 4 <|channel>thought reasoning blocks ----

func TestStripMarkup_Gemma4_ChannelThought_PipeForm(t *testing.T) {
	s := "<|channel>thought\nI should respond warmly.\n<channel|>Hello there!"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_Gemma4_ChannelThought_DecodedForm(t *testing.T) {
	// After TokenToPiece strips pipes: <channel>thought...<channel>
	s := "<channel>thought\nI should respond warmly.\n<channel>Hello there!"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_Gemma4_ChannelThought_Multiple(t *testing.T) {
	s := "<|channel>thought\nFirst thought.\n<channel|>Step one. <|channel>thought\nSecond thought.\n<channel|>Step two."
	got := StripMarkup(s)
	if got != "Step one. Step two." {
		t.Errorf("got %q, want %q", got, "Step one. Step two.")
	}
}

func TestStripMarkup_Gemma4_ChannelThought_NoClosingTag(t *testing.T) {
	s := "Some text. <|channel>thought\nreasoning with no closing tag"
	got := StripMarkup(s)
	if got != "Some text." {
		t.Errorf("got %q, want %q", got, "Some text.")
	}
}

func TestStripMarkup_Gemma4_ChannelThought_OnlyBlock(t *testing.T) {
	s := "<|channel>thought\nonly reasoning\n<channel|>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

func TestStripMarkup_Gemma4_ChannelThought_WithTurnAndToolCall(t *testing.T) {
	// Thought block followed by turn markers and sentence markers.
	s := "<|channel>thought\nI should greet warmly.\n<channel|><|turn>model<s>Hello there!</s>"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

// ---- StripMarkup: Gemma 4 pipe-delimited channel tags (non-thought) ----

func TestStripMarkup_Gemma4_PipeChannelTag_Opening(t *testing.T) {
	// <|channel>speak tag should be stripped (content kept).
	s := "<|channel>speak\nHello there!"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_Gemma4_PipeChannelTag_Closing(t *testing.T) {
	s := "Hello there!<channel|>"
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

// ---- StripMarkup: Gemma 4 <toolcode> Python-style execution blocks ----

func TestStripMarkup_Gemma4_ToolCodeBlock_WithSuffix(t *testing.T) {
	s := "<toolcode>print(tool_movement(command='speak'))</toolcode>Hello!"
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_Gemma4_ToolCodeBlock_WithPrefix(t *testing.T) {
	s := "Hello!<toolcode>print(tool_movement(command='wait'))</toolcode>"
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_Gemma4_ToolCodeBlock_Multiple(t *testing.T) {
	s := "<toolcode>print(tool_movement(command='speak'))</toolcode>Step one. <toolcode>print(tool_movement(command='wait'))</toolcode>Step two."
	got := StripMarkup(s)
	if got != "Step one. Step two." {
		t.Errorf("got %q, want %q", got, "Step one. Step two.")
	}
}

func TestStripMarkup_Gemma4_ToolCodeBlock_NoClose(t *testing.T) {
	s := "Some text. <toolcode>print(tool_movement(command='speak'))"
	got := StripMarkup(s)
	if got != "Some text." {
		t.Errorf("got %q, want %q", got, "Some text.")
	}
}

func TestStripMarkup_Gemma4_ToolCodeBlock_WithToolResult(t *testing.T) {
	// Realistic Gemma 4 output: toolcode block + toolresult block.
	s := "<toolcode>print(tool_movement(command='speak'))</toolcode><toolresult>{\"status\": \"success\"}</toolresult>Hello!"
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

// ---- StripMarkup: orphan <toolcall> tag cleanup ----

func TestStripMarkup_OrphanToolCallTags(t *testing.T) {
	// Bare <toolcall> pairs that somehow survive block stripping.
	s := "<toolcall><toolcall>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

func TestStripMarkup_OrphanPipeToolCallTags(t *testing.T) {
	s := "<|toolcall><|toolcall>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

// ---- StripMarkup: inline JSON tool call objects ----

func TestStripMarkup_InlineJSON_BeforeSpeech(t *testing.T) {
	// Realistic Gemma 4 pattern: JSON tool call + toolresult + <turn> + speech.
	s := `{"name": "tool_movement", "args": {"command": "speak"}}<toolresult>{"status": "success"}</toolresult><turn>Hello there!`
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_InlineJSON_MultipleBeforeSpeech(t *testing.T) {
	s := `<toolresult>{"status": "success"}</toolresult><turn>{"name": "tool_movement", "args": {"command": "slowlook", "angle": 140}}<toolresult>{"status": "success"}</toolresult><turn>Hello!`
	got := StripMarkup(s)
	if got != "Hello!" {
		t.Errorf("got %q, want %q", got, "Hello!")
	}
}

func TestStripMarkup_InlineJSON_WithArgumentsField(t *testing.T) {
	s := `{"name": "tool_movement", "arguments": {"command": "wait"}}<turn>Goodbye!`
	got := StripMarkup(s)
	if got != "Goodbye!" {
		t.Errorf("got %q, want %q", got, "Goodbye!")
	}
}

func TestStripMarkup_InlineJSON_PreservesNonToolCallJSON(t *testing.T) {
	// JSON with "name" but no "args"/"arguments" should NOT be stripped.
	s := `The result is {"name": "Alice", "score": 100}.`
	got := StripMarkup(s)
	if got != `The result is {"name": "Alice", "score": 100}.` {
		t.Errorf("got %q", got)
	}
}

// ---- StripMarkup: orphan <turn> tag cleanup ----

func TestStripMarkup_OrphanBareToolresultAndTurnTags(t *testing.T) {
	// What the dialogue observed: toolresult + hundreds of bare <turn> tags.
	s := "<toolresult><turn><turn><turn>"
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

func TestStripMarkup_OrphanBareTurnSeparator(t *testing.T) {
	// Bare <turn> used as separator between JSON tool call and speech.
	s := `{"name": "tool_movement", "args": {"command": "speak"}}<turn>Hello there!`
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_PipeToolresultNormalized(t *testing.T) {
	// Pipe-delimited <|toolresult>...</toolresult|> pair normalized then stripped.
	s := `<|toolresult>{"status": "success"}<toolresult|>Hello there!`
	got := StripMarkup(s)
	if got != "Hello there!" {
		t.Errorf("got %q, want %q", got, "Hello there!")
	}
}

func TestStripMarkup_OrphanToolresultTag(t *testing.T) {
	// Bare <toolresult> without close tag: everything from tag to end is stripped.
	s := "Hello <toolresult>there!"
	got := StripMarkup(s)
	if got != "Hello" {
		t.Errorf("got %q, want %q", got, "Hello")
	}
}

// ---- StripMarkup: tool result echo blocks (word{...}) ----

func TestStripMarkup_ToolResultEcho_BarePrefix(t *testing.T) {
	// Model echoes tool result as tool{"status":"SUCCESS",...} with no spoken text.
	s := `tool{"status":"SUCCESS","data":{"headangle":90,"headmovement":"speaking"}}`
	got := StripMarkup(s)
	if got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
}

func TestStripMarkup_ToolResultEcho_WithLeadingText(t *testing.T) {
	// Model emits spoken text, then the tool result echo.
	s := `I'll start speaking now. tool{"status":"SUCCESS","data":{"headangle":90}}`
	got := StripMarkup(s)
	if got != "I'll start speaking now." {
		t.Errorf("got %q, want %q", got, "I'll start speaking now.")
	}
}

func TestStripMarkup_ToolResultEcho_NonToolJSONPreserved(t *testing.T) {
	// A word{...} pattern without "status" key should be preserved.
	s := `The config{} object is empty.`
	got := StripMarkup(s)
	if got != `The config{} object is empty.` {
		t.Errorf("got %q", got)
	}
}

// ---- StripMarkup: <tool_result> underscore form (Gemma template native) ----

func TestStripMarkup_ToolResultUnderscore_Stripped(t *testing.T) {
	// <tool_result>...</tool_result> emitted by the Gemma chat template must be stripped.
	s := `<tool_result>The actor's head animates with a subtle nod.</tool_result>Hey there!`
	got := StripMarkup(s)
	if got != "Hey there!" {
		t.Errorf("got %q, want %q", got, "Hey there!")
	}
}

func TestStripMarkup_ToolResultUnderscore_MultipleBlocks(t *testing.T) {
	// Multiple <tool_result> blocks interleaved with speech.
	s := `<tool_result>result1</tool_result>Hello! <tool_result>result2</tool_result>World.`
	got := StripMarkup(s)
	if got != "Hello! World." {
		t.Errorf("got %q, want %q", got, "Hello! World.")
	}
}

func TestStripMarkup_ToolResponseUnderscore_Stripped(t *testing.T) {
	// <tool_response>...</tool_response> (alternative underscore form) must be stripped.
	s := `<tool_response>{"status":"ok"}</tool_response>Good morning!`
	got := StripMarkup(s)
	if got != "Good morning!" {
		t.Errorf("got %q, want %q", got, "Good morning!")
	}
}

// ---- StripMarkup: Phi-4 turn boundary tokens ----

func TestStripMarkup_Phi4_EndTokenStripped(t *testing.T) {
	// Phi-4 appends <|end|> at the end of its response turn.
	// It must be stripped from the output before TTS.
	s := "I'm merely observing the day's events with my superior processing capabilities.<|end|>"
	got := StripMarkup(s)
	want := "I'm merely observing the day's events with my superior processing capabilities."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripMarkup_Phi4_UserTokenTruncates(t *testing.T) {
	// Phi-4 simulating next conversation turn with <|user|>.
	// Everything from that token onwards should be discarded.
	s := "Hello there.<|end|><|user|>What time is it?<|end|><|assistant|>It is noon."
	got := StripMarkup(s)
	if got != "Hello there." {
		t.Errorf("got %q, want %q", got, "Hello there.")
	}
}

func TestStripMarkup_Phi4_AssistantTokenTruncates(t *testing.T) {
	// Phi-4 starting a second assistant turn — discard everything after.
	s := "First response.<|assistant|>Second fabricated response."
	got := StripMarkup(s)
	if got != "First response." {
		t.Errorf("got %q, want %q", got, "First response.")
	}
}

func TestStripMarkup_Phi4_SystemTokenTruncates(t *testing.T) {
	// Phi-4 injecting a <|system|> prompt mid-generation.
	s := "Some text.<|system|>You are a helpful assistant."
	got := StripMarkup(s)
	if got != "Some text." {
		t.Errorf("got %q, want %q", got, "Some text.")
	}
}
