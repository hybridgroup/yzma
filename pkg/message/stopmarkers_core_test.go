package message

import (
	"slices"
	"testing"
)

// allFormats holds each format, so a test covers the whole switch.
var allFormats = []Format{
	FormatAuto, FormatStandard, FormatQwen, FormatGLM,
	FormatMistral, FormatGemma3, FormatGemma, FormatGPT, FormatPhi,
}

func TestStopMarkersForKeepsEOT(t *testing.T) {
	eot := []string{"<|im_end|>", "<end_of_turn>"}

	for _, format := range allFormats {
		markers := StopMarkersFor(eot, format)
		for _, want := range eot {
			if !slices.Contains(markers, want) {
				t.Errorf("format %d: no marker %q in %v", format, want, markers)
			}
		}
	}
}

func TestStopMarkersForAlwaysHasToolResults(t *testing.T) {
	for _, format := range allFormats {
		markers := StopMarkersFor(nil, format)
		if !slices.Contains(markers, "<tool_response") {
			t.Errorf("format %d: no tool result marker in %v", format, markers)
		}
	}
}

func TestStopMarkersForFormatMarkers(t *testing.T) {
	tests := []struct {
		format Format
		want   string
	}{
		{FormatGemma3, "<start_of_turn>user"},
		{FormatGemma, "<|turn>user"},
		{FormatQwen, "<|im_start|>"},
		{FormatPhi, "<|end|>"},
		{FormatStandard, "<|im_end|>"},
		{FormatAuto, "<|im_end|>"},
	}

	for _, tt := range tests {
		markers := StopMarkersFor(nil, tt.format)
		if !slices.Contains(markers, tt.want) {
			t.Errorf("format %d: no marker %q in %v", tt.format, tt.want, markers)
		}
	}
}

// StopMarkersFor must not write into the slice of the caller.
func TestStopMarkersForCopiesEOT(t *testing.T) {
	eot := make([]string, 1, 8)
	eot[0] = "<|im_end|>"

	StopMarkersFor(eot, FormatQwen)

	if len(eot) != 1 || eot[0] != "<|im_end|>" {
		t.Errorf("the eot slice changed: %v", eot)
	}
}

func TestDedupe(t *testing.T) {
	got := dedupe([]string{"a", "", "b", "a", "", "c"})
	want := []string{"a", "b", "c"}

	if !slices.Equal(got, want) {
		t.Errorf("dedupe = %v, want %v", got, want)
	}
}
