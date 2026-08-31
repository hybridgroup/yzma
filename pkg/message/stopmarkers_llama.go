//go:build !(js && wasm)

package message

import (
	"github.com/hybridgroup/yzma/pkg/llama"
)

// StopMarkers returns the markers that should halt generation for a vocabulary
// and a format. See StopMarkersFor.
func StopMarkers(vocab llama.Vocab, format Format) []string {
	return StopMarkersFor(eotMarkers(vocab), format)
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

	return dedupe([]string{piece, llama.VocabGetText(vocab, eot)})
}
