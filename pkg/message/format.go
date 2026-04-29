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

	// FormatGemma expects Gemma 4's call: syntax:
	//   call:funcname{key:<|"|>value<|"|>}
	FormatGemma

	// FormatGPT expects GPT-model tool calls:
	//   .FUNC_NAME <|message|>JSON_ARGS
	FormatGPT
)

// DetectFormat inspects the content of a tool-call block and returns the
// Format that matches it, or FormatAuto when no grammar is recognized.
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
