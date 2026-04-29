// Package message provides types and functions for creating, manipulating,
// and processing LLM messages in various formats.
//
// # Message types
//
// The [Message] interface is the core abstraction. Three concrete types
// implement it:
//
//   - [Chat] - a standard user/assistant/system turn.
//   - [Tool] - an assistant turn that contains one or more tool calls.
//   - [ToolResponse] - a tool turn carrying the result of a tool call.
//
// # Tool calls
//
// [ParseToolCalls] parses tool calls from a raw model response string,
// automatically detecting the grammar used. All argument values are
// normalised to strings in the returned [ToolCall] slice.
//
// [DetectFormat] identifies which grammar a piece of content uses and
// returns the corresponding [Format] constant. The following formats are
// supported:
//
//   - [FormatStandard] - bare JSON: {"name":"…","arguments":{…}}
//   - [FormatQwen] - Qwen3-Coder XML tags: <function=name>…</function>
//   - [FormatGLM] - GLM key/value tags: name<arg_key>k</arg_key><arg_value>v</arg_value>
//   - [FormatMistral] - Mistral/Devstral markers: [TOOL_CALLS]name[ARGS]{…}
//   - [FormatGemma] - Gemma 4 call syntax: call:name{key:<|"|>val<|"|>}
//   - [FormatGPT] - GPT tool call syntax: .name <|message|>{…}
//
// Standard-format responses wrap bare JSON inside <tool_call>…</tool_call>
// envelope tags; all other formats are detected from the raw content.
//
// # Markup stripping
//
// [StripMarkup] removes all tool-call blocks and model-specific markers
// (turn tags, channel directives, bare channel prefixes) from a string,
// leaving only the plain text content. It handles every format listed above.
package message
