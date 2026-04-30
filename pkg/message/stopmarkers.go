package message

import (
	"github.com/hybridgroup/yzma/pkg/llama"
)

// StopMarkers returns the set of string markers that should halt text generation
// when any of them appears in the accumulated output.
//
// It combines two sources:
//  1. The model's own EOT (end-of-turn) token text, obtained directly from the
//     vocabulary via VocabEOT so it works for any model regardless of format.
//  2. Format-specific role/turn-boundary tokens that signal the model has started
//     simulating the next conversation turn (fabricated Q&A).
//
// The caller should stop generation and discard everything from the first
// matching marker onwards.
func StopMarkers(vocab llama.Vocab, format Format) []string {
	markers := eotMarkers(vocab)

	switch format {
	case FormatGemma:
		markers = append(markers,
			// Turn boundary tokens used by Gemma 4.
			"<|turn>user", "<|turn>model",
			"<turn>user", "<turn>model", // decoded form (| stripped)
			"<|turn|>", "<|turn>",
			"<turn|>",
			// Thought channel — internal reasoning, not spoken text.
			"<|channel>thought", "<channel>thought",
		)
	case FormatQwen:
		markers = append(markers,
			// ChatML role tokens; model simulating next turn.
			"<|im_start|>", "<|im_end|>",
			"<im_start>", "<im_end>", // decoded form
		)
	case FormatPhi:
		markers = append(markers,
			// Phi-3/4 turn-boundary tokens.
			"<|user|>", "<|assistant|>", "<|system|>",
		)
	case FormatStandard, FormatAuto:
		// Include ChatML tokens as a safety net for unknown/auto models.
		markers = append(markers,
			"<|im_start|>", "<|im_end|>",
		)
	}

	// Tool-result tokens are common across all formats: stop if the model
	// starts simulating tool results rather than producing spoken text.
	markers = append(markers,
		"<toolresult", "<|toolresult", "<tool_result",
		"<toolresponse", "<|toolresponse", "<tool_response",
		`tool{"status"`,
		"<turnend>", "<|turnend>",
	)

	return markers
}

// eotMarkers returns a deduplicated list of strings for the model's EOT token.
// It uses TokenToPiece (the decoded form that appears in the output stream) as
// the primary value, with VocabGetText as a fallback, so the returned string
// matches exactly what will appear in accumulated output chunks.
func eotMarkers(vocab llama.Vocab) []string {
	eot := llama.VocabEOT(vocab)
	if eot < 0 {
		return nil
	}

	buf := make([]byte, 64)
	n := llama.TokenToPiece(vocab, eot, buf, 0, true)
	piece := ""
	if n > 0 {
		piece = string(buf[:n])
	}

	text := llama.VocabGetText(vocab, eot)

	seen := map[string]bool{}
	var markers []string
	for _, s := range []string{piece, text} {
		if s != "" && !seen[s] {
			seen[s] = true
			markers = append(markers, s)
		}
	}
	return markers
}
