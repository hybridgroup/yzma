package message

import "testing"

func TestDetectFormatFromPath_Phi(t *testing.T) {
	tests := []struct {
		path string
		want Format
	}{
		{"models/phi-4-mini-instruct-abliterated-Q4_K_M.gguf", FormatPhi},
		{"models/Phi-3-mini-4k-instruct-q4.gguf", FormatPhi},
		{"models/phi3.5-mini-instruct.gguf", FormatPhi},
	}
	for _, tt := range tests {
		got := DetectFormatFromPath(tt.path)
		if got != tt.want {
			t.Errorf("DetectFormatFromPath(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestDetectFormatFromPath_KnownFamilies(t *testing.T) {
	tests := []struct {
		path string
		want Format
	}{
		{"models/qwen3.5-4B-instruct.gguf", FormatQwen},
		{"models/gemma-4-E4B-it-Q4_K_M.gguf", FormatGemma},
		{"models/mistral-7b-instruct-v0.2.Q4_K_M.gguf", FormatMistral},
		{"models/devstral-small.gguf", FormatMistral},
		{"models/glm-4-9b-chat-q4.gguf", FormatGLM},
		{"models/llama-3.2-3B-instruct.gguf", FormatAuto},
	}
	for _, tt := range tests {
		got := DetectFormatFromPath(tt.path)
		if got != tt.want {
			t.Errorf("DetectFormatFromPath(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}
