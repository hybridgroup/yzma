package message

import (
	"encoding/json"
	"strconv"
	"strings"
)

// ParseToolCalls attempts to parse tool calls from a model response,
// detecting the grammar in order and returning the first match.
func ParseToolCalls(response string) []ToolCall {
	// parseStandardToolCalls must be tried first against the full response,
	// because it matches on the outer <tool_call>…</tool_call> envelope tags.
	// DetectFormat (used by parseToolCalls) inspects the unwrapped inner
	// content and will never see those tags, so standard-format responses
	// would fall through undetected without this pre-check.
	if calls := parseStandardToolCalls(response); len(calls) > 0 {
		return calls
	}
	// parseInlineJSONToolCalls handles bare JSON tool call objects emitted by
	// models that don't wrap in <tool_call> tags, e.g. Gemma 4 fine-tunes that
	// output {"name":"func","args":{...}} inline in the generation stream.
	// Only attempt this when there are no <tool_call> wrapper tags at all:
	// if tags are present but malformed (unpaired), treat as a failed call
	// rather than falling back to parsing the JSON content directly.
	if !strings.Contains(response, "<tool_call>") && !strings.Contains(response, "</tool_call>") {
		if calls := parseInlineJSONToolCalls(response); len(calls) > 0 {
			return calls
		}
	}
	return parseToolCalls(response)
}

// parseToolCalls routes tool call content to the appropriate model-specific
// parser based on the format detected by DetectFormat.
func parseToolCalls(content string) []ToolCall {
	if len(content) == 0 {
		return nil
	}

	switch DetectFormat(content) {
	case FormatStandard:
		return parseStandardToolCalls(content)
	case FormatQwen:
		return parseQwenToolCalls(content)
	case FormatGLM:
		return parseGLMToolCalls(content)
	case FormatMistral:
		return parseMistralToolCalls(content)
	case FormatGemma:
		return parseGemmaToolCalls(content)
	case FormatGPT:
		return parseGPTToolCalls(content)
	default:
		return nil
	}
}

// parseStandardToolCalls parses <tool_call>{JSON}</tool_call> blocks.
func parseStandardToolCalls(response string) []ToolCall {
	var calls []ToolCall

	start := strings.Index(response, "<tool_call>")
	end := strings.Index(response, "</tool_call>")

	for start != -1 && end != -1 && start < end {
		content := strings.TrimSpace(response[start+len("<tool_call>") : end])

		var parsed struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal([]byte(content), &parsed); err == nil {
			args := make(map[string]string)
			for k, v := range parsed.Arguments {
				switch val := v.(type) {
				case string:
					args[k] = val
				case float64:
					args[k] = strconv.FormatFloat(val, 'f', -1, 64)
				case int:
					args[k] = strconv.Itoa(val)
				case json.Number:
					args[k] = val.String()
				default:
					if b, err := json.Marshal(v); err == nil {
						args[k] = string(b)
					}
				}
			}
			calls = append(calls, ToolCall{
				Type: "function",
				Function: ToolFunction{
					Name:      parsed.Name,
					Arguments: args,
				},
			})
		}

		response = response[end+len("</tool_call>"):]
		start = strings.Index(response, "<tool_call>")
		end = strings.Index(response, "</tool_call>")
	}

	return calls
}

// unmarshalJSONArgs parses a JSON object and returns all values as strings.
// Handles string, float64, int, and json.Number values; other types are
// JSON-marshalled to their string representation.
func unmarshalJSONArgs(raw json.RawMessage) map[string]string {
	args := make(map[string]string)
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return args
	}
	for k, v := range m {
		switch val := v.(type) {
		case string:
			args[k] = val
		case float64:
			args[k] = strconv.FormatFloat(val, 'f', -1, 64)
		case int:
			args[k] = strconv.Itoa(val)
		case json.Number:
			args[k] = val.String()
		default:
			if b, err := json.Marshal(v); err == nil {
				args[k] = string(b)
			}
		}
	}
	return args
}

// parseInlineJSONToolCalls scans s for bare JSON objects of the form
// {"name":"funcname","args":{...}} or {"name":"funcname","arguments":{...}}
// — as emitted by some Gemma 4 fine-tunes without any wrapper tags — and
// returns them as ToolCalls. Objects whose "name" value contains whitespace
// are ignored to avoid false positives on natural-language text.
func parseInlineJSONToolCalls(s string) []ToolCall {
	var calls []ToolCall
	remaining := s
	for {
		idx := strings.Index(remaining, `{"name":`)
		if idx == -1 {
			break
		}
		sub := remaining[idx:]
		end := findJSONObjectEnd(sub) // returns one-past the closing `}`
		if end == -1 {
			break
		}
		objStr := sub[:end]
		remaining = sub[end:]

		var parsed struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
			Args      json.RawMessage `json:"args"`
		}
		if err := json.Unmarshal([]byte(objStr), &parsed); err != nil || parsed.Name == "" {
			continue
		}
		// Require a function-name-like value: no whitespace.
		if strings.ContainsAny(parsed.Name, " \t\n") {
			continue
		}
		// Require either "args" or "arguments" field.
		argsJSON := parsed.Arguments
		if len(argsJSON) == 0 {
			argsJSON = parsed.Args
		}
		if len(argsJSON) == 0 {
			continue
		}

		// Require a non-empty args map — an empty {} is not a valid tool call.
		parsedArgs := unmarshalJSONArgs(argsJSON)
		if len(parsedArgs) == 0 {
			continue
		}

		calls = append(calls, ToolCall{
			Type: "function",
			Function: ToolFunction{
				Name:      parsed.Name,
				Arguments: parsedArgs,
			},
		})
	}
	return calls
}

// stripInlineJSONToolCallBlocks removes bare JSON tool call objects of the form
// {"name":"funcname","args":{...}} or {"name":"funcname","arguments":{...}} from
// s, preserving all other text content. JSON objects with a "name" field but
// no "args"/"arguments" field (e.g. serialised data objects) are left intact.
func stripInlineJSONToolCallBlocks(s string) string {
	var b strings.Builder
	remaining := s
	for {
		idx := strings.Index(remaining, `{"name":`)
		if idx == -1 {
			b.WriteString(remaining)
			break
		}
		sub := remaining[idx:]
		end := findJSONObjectEnd(sub) // returns one-past the closing `}`
		if end == -1 {
			b.WriteString(remaining)
			break
		}
		objStr := sub[:end]
		// Only strip if the object has "args" or "arguments" — it is a tool call.
		if strings.Contains(objStr, `"args":`) || strings.Contains(objStr, `"arguments":`) {
			b.WriteString(remaining[:idx])
		} else {
			// Not a tool call object — keep it.
			b.WriteString(remaining[:idx+end])
		}
		remaining = sub[end:]
	}
	return b.String()
}

// stripToolResultEchoBlocks removes bare tool-result echo patterns of the form
// word{...} where the JSON object contains a "status" key. Some Gemma 4
// fine-tunes echo tool responses back into the generation stream prefixed by a
// plain word (e.g. "tool{\"status\":\"SUCCESS\",...}"). These must be stripped
// so they never reach spoken/MQTT output.
func stripToolResultEchoBlocks(s string) string {
	// Walk through the string looking for word{ patterns.
	var b strings.Builder
	i := 0
	for i < len(s) {
		// Find the next `{`.
		braceIdx := strings.Index(s[i:], "{")
		if braceIdx == -1 {
			b.WriteString(s[i:])
			break
		}
		bracePos := i + braceIdx

		// Check whether the `{` is immediately preceded by a bare word (letters,
		// digits, underscore) with no intervening space — e.g. `tool{`.
		wordStart := bracePos
		for wordStart > i && isWordChar(s[wordStart-1]) {
			wordStart--
		}
		hasWordPrefix := wordStart < bracePos

		// Find the matching closing brace.
		end := findJSONObjectEnd(s[bracePos:]) // one-past
		if end == -1 {
			b.WriteString(s[i:])
			break
		}
		objStr := s[bracePos : bracePos+end]

		// Only strip when the object looks like a tool result: contains "status".
		if hasWordPrefix && strings.Contains(objStr, `"status"`) {
			b.WriteString(s[i:wordStart]) // keep text before word prefix
		} else {
			b.WriteString(s[i : bracePos+end])
		}
		i = bracePos + end
	}
	return b.String()
}

// isWordChar reports whether c is a letter, digit, or underscore.
func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// StripMarkup removes all tool call blocks and model-specific markers from s,
// returning only the plain text content. It handles all known formats:
// Standard (<tool_call>), Qwen (<function=…>), GLM (<arg_key>), Mistral
// ([TOOL_CALLS]), and Gemma 4 (call:), as well as Gemma 4 turn/channel markers
// and thinking/reasoning <think>…</think> and <|channel>thought…<channel|> blocks.
func StripMarkup(s string) string {
	// Remove <think>…</think> reasoning blocks before any other processing.
	// If <think> has no matching </think>, everything from <think> to end is stripped.
	s = stripThinkBlocks(s)

	// Remove Gemma 4 pipe-delimited <|channel>thought…<channel|> reasoning blocks.
	// These are used by some fine-tunes in place of <think> for internal reasoning.
	s = stripChannelThoughtBlocks(s)

	// Remove <toolresponse>…</toolresponse>, <toolresult>…</toolresult>, and
	// <toolcode>…</toolcode> blocks. Some Gemma 4 fine-tunes use these as
	// tool-invocation or simulated-result markers / Python-style execution blocks.
	// Normalise pipe-delimited variants (<|toolresult|>, <|toolresult>, etc.) and
	// underscore variants (<tool_result>, <tool_response>) before block-stripping
	// so all forms are caught.
	s = strings.ReplaceAll(s, "<|toolresult|>", "<toolresult>")
	s = strings.ReplaceAll(s, "<|toolresult>", "<toolresult>")
	s = strings.ReplaceAll(s, "<toolresult|>", "</toolresult>")
	s = strings.ReplaceAll(s, "<tool_result>", "<toolresult>")
	s = strings.ReplaceAll(s, "</tool_result>", "</toolresult>")
	s = strings.ReplaceAll(s, "<|toolresponse|>", "<toolresponse>")
	s = strings.ReplaceAll(s, "<|toolresponse>", "<toolresponse>")
	s = strings.ReplaceAll(s, "<toolresponse|>", "</toolresponse>")
	s = strings.ReplaceAll(s, "<tool_response>", "<toolresponse>")
	s = strings.ReplaceAll(s, "</tool_response>", "</toolresponse>")
	s = stripToolResponseBlocks(s)
	s = stripToolResultBlocks(s)
	s = stripToolCodeBlocks(s)

	// Remove Gemma 4 sentence boundary markers.
	s = strings.ReplaceAll(s, "<s>", "")
	s = strings.ReplaceAll(s, "</s>", " ")

	// Normalise <toolcall> (no underscore) to the canonical form so the
	// Standard block-removal below handles both spellings.
	// Also normalise the Gemma 4 pipe-delimited forms: <|toolcall> (opening)
	// and <toolcall|> (closing), mapping both to <tool_call> so the block
	// stripper below removes the whole block including its call: content.
	s = strings.ReplaceAll(s, "<|toolcall>", "<tool_call>")
	s = strings.ReplaceAll(s, "<toolcall|>", "<tool_call>")
	s = strings.ReplaceAll(s, "<toolcall>", "<tool_call>")
	s = strings.ReplaceAll(s, "</toolcall>", "</tool_call>")

	// Remove Standard <tool_call>…</tool_call> blocks.
	s = stripStandardToolCallBlocks(s)

	// Remove Qwen <function=…>…</function> blocks.
	s = stripQwenFunctionBlocks(s)

	// Remove GLM lines that contain <arg_key>…</arg_key> tags.
	s = stripGLMToolCallLines(s)

	// Remove Mistral [TOOL_CALLS]…[ARGS]{…} blocks.
	s = stripMistralToolCallBlocks(s)

	// Remove GPT .FUNC_NAME <|message|>{…} blocks.
	s = stripGPTToolCallBlocks(s)

	// Remove Gemma 4 call:funcname{...} blocks.
	s = stripGemmaCallBlocks(s)

	// Remove bare inline JSON tool call objects: {"name":"func","args":{...}}.
	// Some Gemma 4 fine-tunes emit tool calls as raw JSON without wrapper tags,
	// interleaved with spoken text separated by bare <turn> markers.
	s = stripInlineJSONToolCallBlocks(s)

	// Remove bare tool result echo blocks: tool{...} or word{...} where the JSON
	// object contains "status" — emitted by some Gemma 4 fine-tunes that echo
	// the tool response back into the generation stream.
	s = stripToolResultEchoBlocks(s)

	// Remove all <turn> markers with optional role suffix.
	s = gemma4TurnTagRE.ReplaceAllString(s, "")

	// Remove Gemma 4 channel tags in all forms:
	//   <|channel>name  – canonical pipe-delimited opening token
	//   <channel|>      – canonical pipe-delimited closing token
	//   <channel>name:  – decoded form (pipes stripped)
	s = gemma4PipeChannelTagRE.ReplaceAllString(s, "")
	s = gemma4ChannelTagRE.ReplaceAllString(s, "")

	// Last-resort cleanup: remove any orphaned <toolcall> / <tool_call> tags
	// (with or without pipe delimiters) that survived block-level stripping.
	s = gemma4OrphanToolCallTagRE.ReplaceAllString(s, "")

	// Last-resort cleanup: remove any orphaned <toolresult>, <toolresponse>,
	// <toolcode>, or bare <turn> tags that survived block-level stripping.
	// These are emitted as unpaired markers by some Gemma 4 fine-tunes.
	s = gemma4OrphanToolResultTagRE.ReplaceAllString(s, "")

	s = strings.TrimSpace(s)

	// Strip a bare channel/command prefix (e.g. "speak:", "wait:", "look:")
	// emitted by models that omit the <channel> wrapper. URL schemes such as
	// "http://" are preserved. Must happen after trimming so the start-of-string
	// anchor in gemma4BareChannelRE is reliable.
	s = stripBareChannelPrefix(s)

	return s
}

// stripStandardToolCallBlocks removes all <tool_call>…</tool_call> blocks,
// handling the self-closing style where a second <tool_call> acts as the
// closing tag.
func stripStandardToolCallBlocks(s string) string {
	for {
		start := strings.Index(s, "<tool_call>")
		if start < 0 {
			break
		}
		after := s[start+len("<tool_call>"):]
		closeIdx := strings.Index(after, "</tool_call>")
		selfCloseIdx := strings.Index(after, "<tool_call>")

		var endInAfter, closeTagLen int
		switch {
		case closeIdx < 0 && selfCloseIdx < 0:
			s = s[:start]
			break
		case closeIdx < 0 || (selfCloseIdx >= 0 && selfCloseIdx < closeIdx):
			endInAfter, closeTagLen = selfCloseIdx, len("<tool_call>")
		default:
			endInAfter, closeTagLen = closeIdx, len("</tool_call>")
		}
		if closeIdx < 0 && selfCloseIdx < 0 {
			break
		}
		s = s[:start] + s[start+len("<tool_call>")+endInAfter+closeTagLen:]
	}
	return s
}

// stripThinkBlocks removes <think>…</think> reasoning/thinking blocks.
// If a <think> tag has no matching </think>, everything from that tag to the
// end of the string is stripped (the block is treated as incomplete/trailing).
func stripThinkBlocks(s string) string {
	for {
		start := strings.Index(s, "<think>")
		if start < 0 {
			break
		}
		end := strings.Index(s[start:], "</think>")
		if end < 0 {
			// No closing tag; strip from <think> to end of string.
			s = strings.TrimSpace(s[:start])
			break
		}
		s = s[:start] + s[start+end+len("</think>"):]
	}
	return s
}

// stripChannelThoughtBlocks removes <|channel>thought…<channel|> reasoning
// blocks used by some Gemma 4 fine-tunes for internal reasoning in place of
// <think> tags. Both the canonical pipe-delimited token forms and the decoded
// forms (pipes stripped by TokenToPiece) are handled.
// If an opening tag has no matching closing tag, everything to end of string
// is stripped.
func stripChannelThoughtBlocks(s string) string {
	// openings and their corresponding closings, longest first so the more
	// specific pipe-delimited form is tried before the decoded form.
	type pair struct{ open, close string }
	variants := []pair{
		{"<|channel>thought", "<channel|>"},
		{"<channel>thought", "<channel>"},
	}
	for _, v := range variants {
		for {
			start := strings.Index(s, v.open)
			if start < 0 {
				break
			}
			searchFrom := start + len(v.open)
			end := strings.Index(s[searchFrom:], v.close)
			if end < 0 {
				s = strings.TrimSpace(s[:start])
				break
			}
			s = s[:start] + s[searchFrom+end+len(v.close):]
		}
	}
	return s
}

// stripToolResponseBlocks removes <toolresponse>…</toolresponse> blocks emitted
// by Gemma 4 fine-tunes as their tool-invocation syntax (e.g.
// <toolresponse>speak()</toolresponse>). If no closing tag is found, everything
// from the opening tag to end of string is stripped.
func stripToolResponseBlocks(s string) string {
	const open = "<toolresponse>"
	const close = "</toolresponse>"
	for {
		start := strings.Index(s, open)
		if start < 0 {
			break
		}
		end := strings.Index(s[start:], close)
		if end < 0 {
			s = strings.TrimSpace(s[:start])
			break
		}
		s = s[:start] + s[start+end+len(close):]
	}
	return s
}

// stripToolResultBlocks removes <toolresult>…</toolresult> blocks that some
// Gemma 4 fine-tunes generate inline to simulate tool execution results.
func stripToolResultBlocks(s string) string {
	const open = "<toolresult>"
	const close = "</toolresult>"
	for {
		start := strings.Index(s, open)
		if start < 0 {
			break
		}
		end := strings.Index(s[start:], close)
		if end < 0 {
			s = strings.TrimSpace(s[:start])
			break
		}
		s = s[:start] + s[start+end+len(close):]
	}
	return s
}

// stripToolCodeBlocks removes <toolcode>…</toolcode> blocks that some Gemma 4
// fine-tunes emit to represent Python-style simulated tool execution.
// Example: <toolcode>print(tool_movement(command='speak'))</toolcode>
func stripToolCodeBlocks(s string) string {
	const open = "<toolcode>"
	const close = "</toolcode>"
	for {
		start := strings.Index(s, open)
		if start < 0 {
			break
		}
		end := strings.Index(s[start:], close)
		if end < 0 {
			s = strings.TrimSpace(s[:start])
			break
		}
		s = s[:start] + s[start+end+len(close):]
	}
	return s
}
