package message

import (
	"encoding/json"
	"strings"
)

const (
	lfmToolCallStart = "<|tool_call_start|>"
	lfmToolCallEnd   = "<|tool_call_end|>"
)

// lfmToolCallBlocks returns the content of each <|tool_call_start|>…<|tool_call_end|> block.
// An unclosed block runs to the end of s.
func lfmToolCallBlocks(s string) []string {
	var blocks []string
	for {
		start := strings.Index(s, lfmToolCallStart)
		if start == -1 {
			return blocks
		}
		s = s[start+len(lfmToolCallStart):]
		end := strings.Index(s, lfmToolCallEnd)
		if end == -1 {
			return append(blocks, strings.TrimSpace(s))
		}
		blocks = append(blocks, strings.TrimSpace(s[:end]))
		s = s[end+len(lfmToolCallEnd):]
	}
}

// stripLFMToolCallBlocks removes all <|tool_call_start|>…<|tool_call_end|> blocks from s.
func stripLFMToolCallBlocks(s string) string {
	for {
		start := strings.Index(s, lfmToolCallStart)
		if start == -1 {
			return s
		}
		end := strings.Index(s[start:], lfmToolCallEnd)
		if end == -1 {
			return s[:start]
		}
		s = s[:start] + s[start+end+len(lfmToolCallEnd):]
	}
}

// parseLFMToolCalls parses LFM2 tool calls in the default Python list form or in JSON.
// <|tool_call_start|>[get_weather(location="NYC")]<|tool_call_end|>
func parseLFMToolCalls(content string) []ToolCall {
	var calls []ToolCall
	for _, block := range lfmToolCallBlocks(content) {
		if strings.HasPrefix(block, "{") || strings.HasPrefix(block, "[{") {
			calls = append(calls, parseLFMJSONToolCalls(block)...)
			continue
		}
		calls = append(calls, parseLFMPythonToolCalls(block)...)
	}
	return calls
}

// parseLFMJSONToolCalls parses the JSON form LFM2 uses when asked for JSON tool calls.
func parseLFMJSONToolCalls(block string) []ToolCall {
	objects := []json.RawMessage{json.RawMessage(block)}
	if strings.HasPrefix(block, "[") {
		objects = nil
		if err := json.Unmarshal([]byte(block), &objects); err != nil {
			return nil
		}
	}
	var calls []ToolCall
	for _, obj := range objects {
		if call := parseOneStandardToolCall(string(obj)); call != nil {
			calls = append(calls, *call)
		}
	}
	return calls
}

// parseLFMPythonToolCalls parses a list such as [a(x="1", y=2), b()].
func parseLFMPythonToolCalls(block string) []ToolCall {
	s := strings.TrimSpace(block)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	var calls []ToolCall
	for {
		s = strings.TrimLeft(s, " \t\r\n,")
		open := strings.Index(s, "(")
		if open == -1 {
			return calls
		}
		name := strings.TrimSpace(s[:open])
		args, rest := parsePythonArgs(s[open+1:])
		if name != "" {
			calls = append(calls, ToolCall{
				Type: "function",
				Function: ToolFunction{
					Name:      name,
					Arguments: args,
				},
			})
		}
		s = rest
	}
}

// parsePythonArgs parses key=value pairs up to the closing paren and returns the text after it.
func parsePythonArgs(s string) (map[string]string, string) {
	args := make(map[string]string)
	for {
		s = strings.TrimLeft(s, " \t\r\n,")
		if s == "" {
			return args, ""
		}
		if s[0] == ')' {
			return args, s[1:]
		}
		eq := strings.Index(s, "=")
		if eq == -1 {
			return args, ""
		}
		key := strings.TrimSpace(s[:eq])
		value, n := parsePythonValue(strings.TrimLeft(s[eq+1:], " \t"))
		args[key] = value
		s = strings.TrimLeft(s[eq+1:], " \t")[n:]
	}
}

// parsePythonValue reads one Python literal from the start of s and returns it as a string
// along with the number of bytes it used.
func parsePythonValue(s string) (string, int) {
	if s == "" {
		return "", 0
	}
	switch s[0] {
	case '"', '\'':
		return parsePythonString(s)
	case '[', '{':
		depth := 0
		var quote byte
		for i := 0; i < len(s); i++ {
			c := s[i]
			if quote != 0 {
				if c == '\\' {
					i++
				} else if c == quote {
					quote = 0
				}
				continue
			}
			switch c {
			case '"', '\'':
				quote = c
			case '[', '{':
				depth++
			case ']', '}':
				depth--
			}
			if depth == 0 {
				return s[:i+1], i + 1
			}
		}
		return s, len(s)
	}
	n := strings.IndexAny(s, ",)")
	if n == -1 {
		n = len(s)
	}
	v := strings.TrimSpace(s[:n])
	switch v {
	case "True":
		v = "true"
	case "False":
		v = "false"
	case "None":
		v = "null"
	}
	return v, n
}

// parsePythonString reads a quoted Python string from the start of s.
func parsePythonString(s string) (string, int) {
	quote := s[0]
	var b strings.Builder
	for i := 1; i < len(s); i++ {
		c := s[i]
		switch {
		case c == quote:
			return b.String(), i + 1
		case c == '\\' && i+1 < len(s):
			i++
			switch s[i] {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			default:
				b.WriteByte(s[i])
			}
		default:
			b.WriteByte(c)
		}
	}
	return b.String(), len(s)
}
