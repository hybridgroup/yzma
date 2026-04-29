package message

import (
	"regexp"
	"strings"
)

// gemmaQuoteToken is the canonical Gemma 4 string delimiter.
const gemmaQuoteToken = "<|\"|>"

// gemmaShortQuoteToken is an alternate form emitted by some Gemma 4 variants.
const gemmaShortQuoteToken = "<\">"

// gemmaQuoteTokens lists all recognised Gemma 4 string delimiters, longest first
// so prefix checks prefer the longer form when both could match.
var gemmaQuoteTokens = []string{gemmaQuoteToken, gemmaShortQuoteToken}

// gemma4TurnTagRE matches Gemma 4 turn boundary tokens in all observed forms:
//
//	<|turn>role  – canonical opening token (pipe-delimited)
//	<turn|>      – canonical closing token
//	<turn>role   – decoded form when TokenToPiece strips the pipe characters
//	</turn>      – alternate closing form
var gemma4TurnTagRE = regexp.MustCompile(`<\|turn>(model|user|assistant|system)?|<turn\|>|</?turn>(model|user|assistant|system)?`)

// gemma4ChannelTagRE matches Gemma 4 channel directives of the form
// <channel>channelname: (e.g. <channel>speak:) that prefix channel content.
var gemma4ChannelTagRE = regexp.MustCompile(`<channel>\w+:`)

// gemma4PipeChannelTagRE matches the pipe-delimited channel token forms:
//
//	<|channel>name  – canonical opening token (e.g. <|channel>speak)
//	<channel|>      – canonical closing token
var gemma4PipeChannelTagRE = regexp.MustCompile(`<\|channel>\w*|<channel\|>`)

// gemma4BareChannelRE matches a bare channel/command prefix at the very start
// of the string (e.g. "speak:", "wait:", "look:") emitted by Gemma 4 models
// that omit the <channel> wrapper. URL-scheme detection (http://) is handled
// in stripBareChannelPrefix rather than via a lookahead (unsupported in Go).
var gemma4BareChannelRE = regexp.MustCompile(`^[a-z][a-z0-9_]*:`)

// gemma4OrphanToolCallTagRE matches any remaining <toolcall>/<tool_call> tag
// variants (with or without pipe delimiters, with or without underscore, open
// or close forms) that survive block-level stripping. These are emitted as
// bare wrapper tokens by some Gemma 4 fine-tunes and must be cleaned up as a
// last resort so they never reach spoken/MQTT output.
var gemma4OrphanToolCallTagRE = regexp.MustCompile(`<\|?/?tool_?call\|?>`)

// gemma4OrphanToolResultTagRE matches orphaned <toolresult>, <toolresponse>,
// <toolcode>, <turnend>, and bare <turn> tags (with or without pipe delimiters,
// open or close forms) that survive block-level stripping. Also matches bare
// <turn> without a role suffix, and <turnend> which some Gemma 4 fine-tunes
// emit as a turn-boundary marker.
var gemma4OrphanToolResultTagRE = regexp.MustCompile(`<\|?/?(toolresult|toolresponse|toolcode|turnend)\|?>|<\|?turn\|?>`)

// stripBareChannelPrefix removes one or more leading channel/command prefixes
// from s (e.g. "angle:90wait:text" → "text"). Each prefix is either a bare
// word: or word:number pattern. URL schemes (followed by //) are preserved.
func stripBareChannelPrefix(s string) string {
	for {
		s = strings.TrimSpace(s)
		loc := gemma4BareChannelRE.FindStringIndex(s)
		if loc == nil {
			break
		}
		after := s[loc[1]:]
		if strings.HasPrefix(after, "//") {
			// URL scheme — leave intact.
			break
		}
		// Skip optional number suffix (e.g. angle:90 → skip "90").
		var numLen int
		for numLen < len(after) && after[numLen] >= '0' && after[numLen] <= '9' {
			numLen++
		}
		remaining := strings.TrimSpace(after[numLen:])
		if remaining == s {
			// No progress — avoid infinite loop.
			break
		}
		s = remaining
	}
	return s
}

// parseGemmaToolCalls parses Gemma 4's proprietary call: syntax.
// Format: call:funcname{key:<|"|>val<|"|>,key2:bare_val}
// Multiple calls may appear back-to-back or separated by arbitrary text.
func parseGemmaToolCalls(response string) []ToolCall {
	var calls []ToolCall

	remaining := response
	for {
		idx := strings.Index(remaining, "call:")
		if idx == -1 {
			break
		}
		remaining = remaining[idx+5:]

		// Function name ends at the opening brace.
		braceIdx := strings.Index(remaining, "{")
		if braceIdx == -1 {
			break
		}
		name := strings.TrimSpace(remaining[:braceIdx])
		remaining = remaining[braceIdx:]

		// Find the matching closing brace. When the model uses mixed
		// quoting the brace matcher can fail; in that case take
		// everything remaining as the raw args so the call still fires.
		braceEnd := findGemmaBraceEnd(remaining)
		var argsRaw string
		if braceEnd == -1 {
			argsRaw = remaining[1:]
			remaining = ""
		} else {
			argsRaw = remaining[1:braceEnd] // content between { and }
			remaining = remaining[braceEnd+1:]
		}

		args := parseGemmaArgs(argsRaw)
		// Only create a tool call when the function has a name and at least
		// one argument.  An empty brace block call:func{} is a model error —
		// skip it rather than dispatching a call that will always fail.
		if name != "" && len(args) > 0 {
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

// stripGemmaCallBlocks removes all call:funcname{...} blocks from s,
// leaving only the non-tool-call text content.
func stripGemmaCallBlocks(s string) string {
	for {
		callIdx := strings.Index(s, "call:")
		if callIdx == -1 {
			break
		}
		braceStartOff := strings.Index(s[callIdx+5:], "{")
		if braceStartOff == -1 {
			break
		}
		braceStart := callIdx + 5 + braceStartOff
		braceEnd := findGemmaBraceEnd(s[braceStart:])
		if braceEnd == -1 {
			s = s[:callIdx]
			break
		}
		s = s[:callIdx] + s[braceStart+braceEnd+1:]
	}
	return s
}

// findGemmaBraceEnd finds the closing brace that matches the opening brace at
// position 0, accounting for nested braces. Returns the index of the closing
// brace, or -1 if not found. Braces inside quoted strings are ignored so that
// content like `board[move-1] != 0 {` doesn't break the match.
//
// Two quoting modes are supported:
//   - Gemma4 quote tokens (<|"|> or <">): paired as open/close delimiters;
//     everything between them (including standard " and braces) is skipped.
//   - Standard JSON " quotes: used only when no Gemma quote tokens are present
//     in the input, since they contain a literal " that would confuse
//     JSON-style string scanning.
func findGemmaBraceEnd(s string) int {
	if len(s) == 0 || s[0] != '{' {
		return -1
	}

	// When any Gemma quote token is present, use Gemma-style quoting.
	// Standard " characters inside Gemma-delimited values must NOT be
	// treated as JSON string boundaries.
	useJSONQuotes := true
	for _, tok := range gemmaQuoteTokens {
		if strings.Contains(s, tok) {
			useJSONQuotes = false
			break
		}
	}

	depth := 0
	i := 0
	for i < len(s) {
		// Pair any Gemma quote token — skip everything between open and close.
		skipped := false
		for _, tok := range gemmaQuoteTokens {
			if strings.HasPrefix(s[i:], tok) {
				i += len(tok)
				for i < len(s) {
					if strings.HasPrefix(s[i:], tok) {
						i += len(tok)
						break
					}
					i++
				}
				skipped = true
				break
			}
		}
		if skipped {
			continue
		}

		// Skip standard JSON quoted strings only when the model is using
		// pure JSON format (no Gemma quote tokens anywhere in the tool call).
		if useJSONQuotes && s[i] == '"' {
			i++
			for i < len(s) {
				if s[i] == '\\' {
					i += 2 // skip escaped character
					continue
				}
				if s[i] == '"' {
					i++
					break
				}
				i++
			}
			continue
		}

		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
		i++
	}

	return -1
}

// findClosingGemmaQuote finds the position of the closing quote token that
// ends a value. For nested structures containing their own tokens, the correct
// closing token is the one followed by a structural character (comma, closing
// brace/bracket, or double-quote for JSON transition) or end of string.
func findClosingGemmaQuote(s, token string) int {
	searchFrom := 0

	for {
		idx := strings.Index(s[searchFrom:], token)
		if idx == -1 {
			return -1
		}

		pos := searchFrom + idx
		afterQuote := pos + len(token)

		if afterQuote >= len(s) {
			return pos
		}

		// Closing token if followed by a structural character.
		// Accept double-quote as a valid transition character since the
		// model may mix Gemma format with standard JSON mid-output.
		switch s[afterQuote] {
		case ',', '}', ']', '"':
			return pos
		}

		searchFrom = afterQuote
	}
}

// findClosingStandardQuote finds the closing " that ends a JSON-like value.
// The correct closing quote is the one followed by a structural character
// (comma, closing brace) or end of string — not one embedded in content.
func findClosingStandardQuote(s string) int {
	searchFrom := 0

	for {
		idx := strings.Index(s[searchFrom:], "\"")
		if idx == -1 {
			return -1
		}

		pos := searchFrom + idx

		// Skip escaped quotes.
		if pos > 0 && s[pos-1] == '\\' {
			searchFrom = pos + 1
			continue
		}

		afterQuote := pos + 1
		if afterQuote >= len(s) {
			return pos
		}

		next := s[afterQuote]
		if next == ',' || next == '}' || next == ' ' ||
			next == '\n' || next == '\r' || next == '\t' {
			return pos
		}

		searchFrom = afterQuote
	}
}

// parseGemmaArgs parses the key-value pairs inside a Gemma4 tool call argument
// block. Values are delimited by <|"|> or <"> tokens (acting as quotes).
// Format: key1:<|"|>value1<|"|>, key2:<|"|>value2<|"|>
func parseGemmaArgs(raw string) map[string]string {
	args := make(map[string]string)
	remaining := raw

	for len(remaining) > 0 {
		colonIdx := strings.Index(remaining, ":")
		if colonIdx == -1 {
			break
		}

		key := strings.TrimLeft(remaining[:colonIdx], ", \t\n")
		key = strings.Trim(key, "\"")
		remaining = remaining[colonIdx+1:]

		// Strip leading whitespace between the colon and the value. Models that
		// emit JSON-style key: "value" pairs (with a space after the colon) would
		// otherwise fail the quote-prefix checks below.
		remaining = strings.TrimSpace(remaining)

		// Value wrapped in any recognised Gemma quote token.
		var openToken string
		for _, tok := range gemmaQuoteTokens {
			if strings.HasPrefix(remaining, tok) {
				openToken = tok
				break
			}
		}
		if openToken != "" {
			remaining = remaining[len(openToken):]

			endQuote := findClosingGemmaQuote(remaining, openToken)
			if endQuote == -1 {
				args[key] = strings.TrimSpace(remaining)
				break
			}

			args[key] = remaining[:endQuote]
			remaining = remaining[endQuote+len(openToken):]
			continue
		}

		// Value wrapped in standard JSON double quotes.
		if strings.HasPrefix(remaining, "\"") {
			remaining = remaining[1:]

			endQuote := findClosingStandardQuote(remaining)
			if endQuote == -1 {
				args[key] = strings.TrimSpace(remaining)
				break
			}

			args[key] = remaining[:endQuote]
			remaining = remaining[endQuote+1:]
			continue
		}

		// Bare value (number, bool, etc.) — read to next comma or }.
		// If the value begins with { it is a nested object: parse it
		// recursively and merge its key-value pairs into the parent map
		// rather than storing the raw brace text.
		if strings.HasPrefix(remaining, "{") {
			nestedEnd := findGemmaBraceEnd(remaining)
			var nestedRaw string
			if nestedEnd == -1 {
				nestedRaw = remaining[1:]
				remaining = ""
			} else {
				nestedRaw = remaining[1:nestedEnd]
				remaining = remaining[nestedEnd+1:]
			}
			for k, v := range parseGemmaArgs(nestedRaw) {
				args[k] = v
			}
			if remaining == "" {
				break
			}
			continue
		}

		endIdx := strings.IndexAny(remaining, ",}")
		var v string
		if endIdx == -1 {
			v = strings.TrimSpace(remaining)
			remaining = ""
		} else {
			v = strings.TrimSpace(remaining[:endIdx])
			remaining = remaining[endIdx:]
		}
		args[key] = v

		if remaining == "" {
			break
		}
	}

	return args
}
