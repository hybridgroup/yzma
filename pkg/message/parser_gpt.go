package message

import (
	"encoding/json"
	"strconv"
	"strings"
)

// stripGPTToolCallBlocks removes all .FUNC_NAME <|message|>{...} blocks from s.
func stripGPTToolCallBlocks(s string) string {
	for {
		dotIdx := strings.Index(s, ".")
		if dotIdx == -1 {
			break
		}

		msgIdx := strings.Index(s[dotIdx:], "<|message|>")
		if msgIdx == -1 {
			break
		}

		jsonStart := dotIdx + msgIdx + len("<|message|>")
		jsonEnd := findJSONObjectEnd(s[jsonStart:])
		if jsonEnd == -1 {
			s = s[:dotIdx]
			break
		}
		s = s[:dotIdx] + s[jsonStart+jsonEnd:]
	}
	return s
}

// parseGPTToolCalls parses GPT-model tool calls.
// Format: .FUNC_NAME <|message|>JSON_ARGS
func parseGPTToolCalls(content string) []ToolCall {
	var calls []ToolCall

	remaining := content
	for {
		dotIdx := strings.Index(remaining, ".")
		if dotIdx == -1 {
			break
		}

		remaining = remaining[dotIdx:]

		msgIdx := strings.Index(remaining, "<|message|>")
		if msgIdx == -1 {
			break
		}

		// Extract function name (between . and space before <|message|>).
		prefix := remaining[:msgIdx]
		parts := strings.SplitN(prefix, " ", 2)
		name := strings.TrimPrefix(parts[0], ".")

		// Move past <|message|> to get the JSON.
		jsonStart := msgIdx + len("<|message|>")
		remaining = remaining[jsonStart:]

		jsonEnd := findJSONObjectEnd(remaining)
		var argsJSON string
		switch {
		case jsonEnd == -1:
			argsJSON = remaining
			remaining = ""
		default:
			argsJSON = remaining[:jsonEnd]
			remaining = remaining[jsonEnd:]
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

		if name != "" {
			calls = append(calls, ToolCall{
				Type: "function",
				Function: ToolFunction{
					Name:      name,
					Arguments: args,
				},
			})
		}

		if remaining == "" {
			break
		}
	}

	return calls
}
