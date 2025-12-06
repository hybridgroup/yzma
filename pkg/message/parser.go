package message

import (
	"encoding/json"
	"strconv"
	"strings"
)

// ParseToolCalls attempts to parse tool calls from a model response
func ParseToolCalls(response string) []ToolCall {
	var calls []ToolCall

	// Look for <tool_call> tags
	start := strings.Index(response, "<tool_call>")
	end := strings.Index(response, "</tool_call>")

	for start != -1 && end != -1 && start < end {
		content := response[start+len("<tool_call>") : end]
		content = strings.TrimSpace(content)

		// Parse JSON content using interface{} to handle any type
		var parsed struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}

		if err := json.Unmarshal([]byte(content), &parsed); err == nil {
			// Convert arguments to map[string]string
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
					// For other types, try to marshal back to JSON string
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

		// Look for more tool calls
		response = response[end+len("</tool_call>"):]
		start = strings.Index(response, "<tool_call>")
		end = strings.Index(response, "</tool_call>")
	}

	return calls
}
