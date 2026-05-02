package message

import "strings"

// Format identifies which tool-call grammar a model response uses.
type Format int

const (
	// FormatAuto is returned when no known grammar is detected.
	FormatAuto Format = iota

	// FormatStandard expects bare JSON: {"name":"...","arguments":{...}}
	FormatStandard

	// FormatQwen expects Qwen3-Coder XML-like tags:
	//   <function=name>\n<parameter=key>\nvalue\n</parameter>\n</function>
	FormatQwen

	// FormatGLM expects GLM <arg_key>/<arg_value> tags:
	//   funcname<arg_key>key</arg_key><arg_value>value</arg_value>
	FormatGLM

	// FormatMistral expects Mistral/Devstral bracket markers:
	//   [TOOL_CALLS]funcname[ARGS]{...}
	FormatMistral

	// FormatGemma3 expects Gemma 3's turn-based format:
	//   <start_of_turn>user\n…<end_of_turn>\n
	//   <start_of_turn>model\n…<end_of_turn>\n
	// No system role: system instructions are prepended to the first user turn.
	FormatGemma3

	// FormatGemma expects Gemma 4's call: syntax:
	//   call:funcname{key:<|"|>value<|"|>}
	FormatGemma

	// FormatGPT expects GPT-model tool calls:
	//   .FUNC_NAME <|message|>JSON_ARGS
	FormatGPT

	// FormatPhi is used for Phi-family models (phi-3, phi-4, etc.).
	// They use standard JSON tool calls but have distinct turn-boundary tokens
	// (<|end|>, <|user|>, <|assistant|>, <|system|>) that must be treated as
	// generation stop markers.
	FormatPhi
)

// DetectFormat inspects a tool-call content block and returns the Format that
// matches it, or FormatAuto when no grammar is recognized. It only examines
// structural markers in the content — it does NOT inspect model names. Use
// DetectFormatFromPath to identify the format from a model file path.
func DetectFormat(content string) Format {
	switch {
	case strings.HasPrefix(content, "{\"name\""):
		return FormatStandard
	case strings.HasPrefix(content, "<function="):
		return FormatQwen
	case strings.Contains(content, "<arg_key>"):
		return FormatGLM
	case strings.Contains(content, "[TOOL_CALLS]"):
		return FormatMistral
	case strings.Contains(content, "call:"):
		return FormatGemma
	case strings.Contains(content, "<|message|>"):
		return FormatGPT
	default:
		return FormatAuto
	}
}

// DetectFormatFromPath inspects a model file path and returns the Format for
// that model family based on well-known name substrings (case-insensitive).
// Returns FormatAuto when the path does not match any known family.
func DetectFormatFromPath(path string) Format {
	lower := strings.ToLower(path)
	switch {
	case strings.Contains(lower, "qwen"):
		return FormatQwen
	case strings.Contains(lower, "gemma-3"), strings.Contains(lower, "gemma3"):
		return FormatGemma3
	case strings.Contains(lower, "gemma"):
		return FormatGemma
	case strings.Contains(lower, "mistral"), strings.Contains(lower, "devstral"):
		return FormatMistral
	case strings.Contains(lower, "glm"):
		return FormatGLM
	case strings.Contains(lower, "phi"):
		return FormatPhi
	default:
		return FormatAuto
	}
}
