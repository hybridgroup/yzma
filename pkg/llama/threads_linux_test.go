package llama

import (
	"reflect"
	"testing"
)

func TestParseCPUList(t *testing.T) {
	tests := []struct {
		in   string
		want []int
	}{
		{"0-15\n", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
		{"16", []int{16}},
		{"0,2,4-6", []int{0, 2, 4, 5, 6}},
		{"", nil},
		{"bad", nil},
	}
	for _, tt := range tests {
		if got := parseCPUList(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseCPUList(%q) gave %v, want %v", tt.in, got, tt.want)
		}
	}
}
