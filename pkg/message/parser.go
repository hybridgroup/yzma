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

// StripMarkup removes all tool call blocks and model-specific markers from s,
// returning only the plain text content. It handles all known formats:
// Standard (<tool_call>), Qwen (<function=…>), GLM (<arg_key>), Mistral
// ([TOOL_CALLS]), and Gemma 4 (call:), as well as Gemma 4 turn/channel markers.
func StripMarkup(s string) string {
	// Normalise <toolcall> (no underscore) to the canonical form so the
	// Standard block-removal below handles both spellings.
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

	// Remove all <turn> markers with optional role suffix.
	s = gemma4TurnTagRE.ReplaceAllString(s, "")

	// Remove Gemma 4 <channel>name: directives (e.g. <channel>speak:).
	s = gemma4ChannelTagRE.ReplaceAllString(s, "")

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
