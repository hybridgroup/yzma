package message

import (
	"slices"
	"testing"

	"github.com/hybridgroup/yzma/pkg/llama"
)

// StopMarkers is tested with a zero (nil) vocab so VocabEOT returns -1 and
// eotMarkers returns nil. This exercises the format-dispatch logic cleanly
// without needing a real model file.

func TestStopMarkers_Gemma_ContainsGemmaTurnTokens(t *testing.T) {
	markers := StopMarkers(llama.Vocab(0), FormatGemma)
	for _, want := range []string{
		"<|turn>user", "<|turn>model",
		"<turn>user", "<turn>model",
		"<|turn|>", "<|turn>",
		"<turn|>",
		"<|channel>thought", "<channel>thought",
	} {
		if !slices.Contains(markers, want) {
			t.Errorf("StopMarkers(FormatGemma) missing %q", want)
		}
	}
}

func TestStopMarkers_Gemma_DoesNotContainChatMLOrPhiTokens(t *testing.T) {
	markers := StopMarkers(llama.Vocab(0), FormatGemma)
	for _, unwanted := range []string{"<|im_start|>", "<|user|>", "<|assistant|>"} {
		if slices.Contains(markers, unwanted) {
			t.Errorf("StopMarkers(FormatGemma) unexpectedly contains %q", unwanted)
		}
	}
}

func TestStopMarkers_Qwen_ContainsChatMLTokens(t *testing.T) {
	markers := StopMarkers(llama.Vocab(0), FormatQwen)
	for _, want := range []string{"<|im_start|>", "<|im_end|>", "<im_start>", "<im_end>"} {
		if !slices.Contains(markers, want) {
			t.Errorf("StopMarkers(FormatQwen) missing %q", want)
		}
	}
}

func TestStopMarkers_Qwen_DoesNotContainGemmaOrPhiTokens(t *testing.T) {
	markers := StopMarkers(llama.Vocab(0), FormatQwen)
	for _, unwanted := range []string{"<|turn>user", "<|user|>", "<|assistant|>"} {
		if slices.Contains(markers, unwanted) {
			t.Errorf("StopMarkers(FormatQwen) unexpectedly contains %q", unwanted)
		}
	}
}

func TestStopMarkers_Phi_ContainsPhiTurnTokens(t *testing.T) {
	markers := StopMarkers(llama.Vocab(0), FormatPhi)
	for _, want := range []string{"<|user|>", "<|assistant|>", "<|system|>"} {
		if !slices.Contains(markers, want) {
			t.Errorf("StopMarkers(FormatPhi) missing %q", want)
		}
	}
}

func TestStopMarkers_Phi_DoesNotContainGemmaOrChatMLTokens(t *testing.T) {
	markers := StopMarkers(llama.Vocab(0), FormatPhi)
	for _, unwanted := range []string{"<|turn>user", "<|im_start|>"} {
		if slices.Contains(markers, unwanted) {
			t.Errorf("StopMarkers(FormatPhi) unexpectedly contains %q", unwanted)
		}
	}
}

func TestStopMarkers_AllFormats_ContainToolResultMarkers(t *testing.T) {
	// Tool-result stop tokens should be present regardless of format.
	toolMarkers := []string{
		"<toolresult", "<tool_result",
		"<toolresponse", "<tool_response",
		`tool{"status"`,
		"<turnend>", "<|turnend>",
	}
	for _, format := range []Format{FormatGemma, FormatGemma3, FormatQwen, FormatPhi, FormatStandard, FormatAuto} {
		markers := StopMarkers(llama.Vocab(0), format)
		for _, want := range toolMarkers {
			if !slices.Contains(markers, want) {
				t.Errorf("StopMarkers(format=%v) missing tool marker %q", format, want)
			}
		}
	}
}

func TestStopMarkers_Gemma3_ContainsTurnTokens(t *testing.T) {
	markers := StopMarkers(llama.Vocab(0), FormatGemma3)
	for _, want := range []string{
		"<start_of_turn>user", "<start_of_turn>model",
	} {
		if !slices.Contains(markers, want) {
			t.Errorf("StopMarkers(FormatGemma3) missing %q", want)
		}
	}
}

func TestStopMarkers_Gemma3_DoesNotContainGemma4Tokens(t *testing.T) {
	markers := StopMarkers(llama.Vocab(0), FormatGemma3)
	for _, unwanted := range []string{"<|turn>user", "<|turn>model", "<|channel>thought"} {
		if slices.Contains(markers, unwanted) {
			t.Errorf("StopMarkers(FormatGemma3) unexpectedly contains %q", unwanted)
		}
	}
}

func TestDetectFormatFromPath_PhiReturnsPhi(t *testing.T) {
	tests := []string{
		"models/phi-4-mini-instruct-abliterated-Q4_K_M.gguf",
		"models/Phi-3-mini-4k-instruct-q4.gguf",
		"models/phi3.5-mini-instruct.gguf",
	}
	for _, path := range tests {
		got := DetectFormatFromPath(path)
		if got != FormatPhi {
			t.Errorf("DetectFormatFromPath(%q) = %v, want FormatPhi", path, got)
		}
	}
}
