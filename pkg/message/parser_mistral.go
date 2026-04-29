package message

import (
	"encoding/json"
	"strconv"
	"strings"
)

// stripMistralToolCallBlocks removes all [TOOL_CALLS]…[ARGS]{…} blocks from s.
func stripMistralToolCallBlocks(s string) string {
	for {
		start := strings.Index(s, "[TOOL_CALLS]")
		if start == -1 {
			break
		}
		argsIdx := strings.Index(s[start:], "[ARGS]")
		if argsIdx == -1 {
			s = s[:start]
			break
		}
		jsonStart := start + argsIdx + len("[ARGS]")
		endIdx := findJSONObjectEnd(s[jsonStart:])
		if endIdx == -1 {
			s = s[:start]
			break
		}
		s = s[:start] + s[jsonStart+endIdx:]
	}
	return s
}

// parseMistralToolCalls parses Mistral/Devstral style tool calls.
// Format: [TOOL_CALLS]get_weather[ARGS]{"location": "NYC"}
func parseMistralToolCalls(content string) []ToolCall {
	var calls []ToolCall

	remaining := content
	for {
		callStart := strings.Index(remaining, "[TOOL_CALLS]")
		if callStart == -1 {
			break
		}

		argsStart := strings.Index(remaining[callStart:], "[ARGS]")
		if argsStart == -1 {
			break
		}

		name := strings.TrimSpace(remaining[callStart+12 : callStart+argsStart])
		argsContent := remaining[callStart+argsStart+6:]

		endIdx := findJSONObjectEnd(argsContent)
		var argsJSON string
		switch {
		case endIdx == -1:
			argsJSON = argsContent
			remaining = ""
		default:
			argsJSON = argsContent[:endIdx]
			remaining = argsContent[endIdx:]
		}

		var rawArgs map[string]any
		args := make(map[string]string)
		if err := json.Unmarshal([]byte(argsJSON), &rawArgs); err == nil {
			for k, v := range rawArgs {
				switch val := v.(type) {
				case string:
					args[k] = val
				case float64:
					args[k] = strconv.FormatFloat(val, 'f', -1, 64)
				case bool:
					args[k] = strconv.FormatBool(val)
				default:
					if b, err := json.Marshal(v); err == nil {
						args[k] = string(b)
					}
				}
			}
		}

		calls = append(calls, ToolCall{
			Type: "function",
			Function: ToolFunction{
				Name:      name,
				Arguments: args,
			},
		})
	}

	return calls
}

// findJSONObjectEnd returns the index one past the closing brace of the first
// complete JSON object in s, or -1 if no complete object is found.
func findJSONObjectEnd(s string) int {
	depth := 0
	inString := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inString {
			if c == '\\' {
				i++ // skip escaped character
			} else if c == '"' {
				inString = false
			}
			continue
		}
		switch c {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}
