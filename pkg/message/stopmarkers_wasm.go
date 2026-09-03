//go:build js && wasm

package message

import (
	"github.com/hybridgroup/yzma/pkg/llamawasm"
)

// eotCandidates are the end of turn tokens of the usual instruct models. A
// module from a release before ABI version 5 has no call that gives this token,
// thus the vocabulary decides which of these are real.
var eotCandidates = []string{
	"<|im_end|>",
	"<end_of_turn>",
	"<|eot_id|>",
	"<|end|>",
	"<|endoftext|>",
}

// StopMarkers returns the markers that should halt generation for a vocabulary
// and a format. See StopMarkersFor.
func StopMarkers(vocab llamawasm.Vocab, format Format) []string {
	return StopMarkersFor(eotMarkers(vocab), format)
}

// eotMarkers gives the end of turn text of the model.
//
// It takes the token of the vocabulary, the same as the version for a host. A
// module from a release before ABI version 5, or a model with no such token,
// goes to guessEOTMarkers.
func eotMarkers(vocab llamawasm.Vocab) []string {
	eot := llamawasm.VocabEOT(vocab)
	if eot < 0 {
		return guessEOTMarkers(vocab)
	}

	buf := make([]byte, 64)
	piece := ""
	if n := llamawasm.TokenToPiece(vocab, eot, buf, 0, true); n > 0 {
		piece = string(buf[:n])
	}

	markers := dedupe([]string{piece, llamawasm.VocabGetText(vocab, eot)})
	if len(markers) == 0 {
		return guessEOTMarkers(vocab)
	}
	return markers
}

// guessEOTMarkers takes the text of the end of sequence token, which is the
// same token for most instruct models, and adds each candidate that the
// vocabulary holds as one token.
func guessEOTMarkers(vocab llamawasm.Vocab) []string {
	buf := make([]byte, 64)

	var markers []string
	if eos := llamawasm.VocabEOS(vocab); eos >= 0 {
		if n := llamawasm.TokenToPiece(vocab, eos, buf, 0, true); n > 0 {
			markers = append(markers, string(buf[:n]))
		}
	}

	for _, s := range eotCandidates {
		if len(llamawasm.Tokenize(vocab, s, false, true)) == 1 {
			markers = append(markers, s)
		}
	}

	return dedupe(markers)
}
