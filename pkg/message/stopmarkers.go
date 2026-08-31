package message

// StopMarkersFor returns the set of string markers that should halt text
// generation when any of them appears in the accumulated output.
//
// The eot argument holds the end of turn text of the model, which a caller gets
// from the vocabulary. StopMarkers gives it for a vocabulary of the backend of
// the build.
//
// The caller should stop generation and discard everything from the first
// matching marker onwards.
func StopMarkersFor(eot []string, format Format) []string {
	markers := append([]string(nil), eot...)

	switch format {
	case FormatGemma3:
		markers = append(markers,
			// Turn boundary tokens for Gemma 3. Stop if the model starts
			// simulating the next conversation turn.
			"<start_of_turn>user", "<start_of_turn>model",
		)
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
			// Phi-3/4 EOT token. Added explicitly because VocabEOT may return -1
			// for some fine-tuned/abliterated variants that don't register a
			// distinct EOT token — eotMarkers() would then contribute nothing.
			"<|end|>",
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

// dedupe removes the empty strings and the repeated values, and keeps the order.
func dedupe(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range values {
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
