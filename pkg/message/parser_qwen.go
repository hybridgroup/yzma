package message

import "strings"

// stripQwenFunctionBlocks removes all <function=…>…</function> blocks from s.
func stripQwenFunctionBlocks(s string) string {
	for {
		start := strings.Index(s, "<function=")
		if start == -1 {
			break
		}
		end := strings.Index(s[start:], "</function>")
		if end == -1 {
			s = s[:start]
			break
		}
		s = s[:start] + s[start+end+len("</function>"):]
	}
	return s
}

// parseBareQwenFunctionCalls extracts <function=…>…</function> blocks that
// are NOT already enclosed inside a <tool_call>…</tool_call> wrapper.  Some
// Qwen fine-tunes wrap the first function call in <tool_call> tags but emit
// subsequent calls as bare blocks outside any wrapper; this function collects
// the bare ones so they can be appended to the already-parsed wrapped calls.
func parseBareQwenFunctionCalls(response string) []ToolCall {
	// Build a set of byte ranges covered by <tool_call>…</tool_call> spans.
	type span struct{ start, end int }
	var wrapped []span
	r := response
	offset := 0
	for {
		s := strings.Index(r, "<tool_call>")
		e := strings.Index(r, "</tool_call>")
		if s == -1 || e == -1 || s > e {
			break
		}
		wrapped = append(wrapped, span{offset + s, offset + e + len("</tool_call>")})
		advance := e + len("</tool_call>")
		r = r[advance:]
		offset += advance
	}

	inWrapped := func(pos int) bool {
		for _, sp := range wrapped {
			if pos >= sp.start && pos < sp.end {
				return true
			}
		}
		return false
	}

	// Scan for <function=…> blocks that start outside any wrapped span.
	var calls []ToolCall
	remaining := response
	scanOffset := 0
	for {
		funcStart := strings.Index(remaining, "<function=")
		if funcStart == -1 {
			break
		}
		absStart := scanOffset + funcStart
		if inWrapped(absStart) {
			// Skip past this block (find its </function> end and advance).
			funcEnd := strings.Index(remaining[funcStart:], "</function>")
			if funcEnd == -1 {
				break
			}
			advance := funcStart + funcEnd + len("</function>")
			scanOffset += advance
			remaining = remaining[advance:]
			continue
		}
		// Parse just this one block.
		blockEnd := strings.Index(remaining[funcStart:], "</function>")
		if blockEnd == -1 {
			break
		}
		blockEnd += funcStart + len("</function>")
		block := remaining[funcStart:blockEnd]
		if parsed := parseQwenToolCalls(block); len(parsed) > 0 {
			calls = append(calls, parsed...)
		}
		scanOffset += blockEnd
		remaining = remaining[blockEnd:]
	}
	return calls
}

// parseQwenToolCalls parses Qwen3-Coder style tool calls with XML-like tags.
// Format: <function=get_weather>\n<parameter=location>\nNYC\n</parameter>\n</function>
func parseQwenToolCalls(content string) []ToolCall {
	var calls []ToolCall

	for {
		funcStart := strings.Index(content, "<function=")
		if funcStart == -1 {
			break
		}

		funcEnd := strings.Index(content[funcStart:], ">")
		if funcEnd == -1 {
			break
		}

		name := strings.TrimSpace(content[funcStart+10 : funcStart+funcEnd])

		bodyStart := funcStart + funcEnd + 1
		closeFunc := strings.Index(content[bodyStart:], "</function>")
		if closeFunc == -1 {
			break
		}
		closeFunc += bodyStart

		funcBody := content[bodyStart:closeFunc]
		args := make(map[string]string)

		remaining := funcBody
		for {
			paramStart := strings.Index(remaining, "<parameter=")
			if paramStart == -1 {
				break
			}

			paramNameEnd := strings.Index(remaining[paramStart:], ">")
			if paramNameEnd == -1 {
				break
			}

			paramName := strings.TrimSpace(remaining[paramStart+11 : paramStart+paramNameEnd])

			valueStart := paramStart + paramNameEnd + 1
			paramCloseRel := strings.Index(remaining[valueStart:], "</parameter>")
			if paramCloseRel == -1 {
				break
			}
			paramClose := valueStart + paramCloseRel

			args[paramName] = strings.TrimSpace(remaining[valueStart:paramClose])

			remaining = remaining[paramClose+12:]
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

		content = content[closeFunc+11:]
	}

	return calls
}
