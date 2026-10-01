package message

import "strings"

const (
	glmKeyOpen    = "<arg_key>"
	glmKeyClose   = "</arg_key>"
	glmValueOpen  = "<arg_value>"
	glmValueClose = "</arg_value>"
)

// stripGLMToolCallLines removes lines that contain GLM-style <arg_key> tags.
func stripGLMToolCallLines(s string) string {
	lines := strings.Split(s, "\n")
	filtered := lines[:0]
	for _, line := range lines {
		if !strings.Contains(line, "<arg_key>") {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, "\n")
}

// parseGLMToolCalls parses GLM-style tool calls with <arg_key>/<arg_value> tags.
// Format: get_weather<arg_key>location</arg_key><arg_value>NYC</arg_value>
func parseGLMToolCalls(content string) []ToolCall {
	var calls []ToolCall

	for _, call := range strings.Split(content, "\n") {
		if call == "" {
			continue
		}

		// Find the function name (everything before the first <arg_key>)
		argKeyIdx := strings.Index(call, "<arg_key>")
		if argKeyIdx == -1 {
			continue
		}

		name := strings.TrimSpace(call[:argKeyIdx])
		args := make(map[string]string)

		// Parse all <arg_key>...</arg_key><arg_value>...</arg_value> pairs.
		// Each tag is searched after the one before it, so out of order tags cannot panic.
		remaining := call[argKeyIdx:]
		for {
			key, rest, ok := cutBetween(remaining, glmKeyOpen, glmKeyClose)
			if !ok {
				break
			}
			value, rest, ok := cutBetween(rest, glmValueOpen, glmValueClose)
			if !ok {
				break
			}
			args[key] = value
			remaining = rest
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
	}

	return calls
}

// cutBetween returns the text between the first start tag in s and the end tag after it,
// and the rest of s after the end tag.
func cutBetween(s, start, end string) (string, string, bool) {
	_, after, ok := strings.Cut(s, start)
	if !ok {
		return "", "", false
	}
	return strings.Cut(after, end)
}
