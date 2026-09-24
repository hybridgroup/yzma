package llama

import (
	"runtime"
	"testing"
)

func TestThreads(t *testing.T) {
	n := Threads()
	if n < 1 {
		t.Errorf("Threads gave %d, want at least 1", n)
	}
	if n > int32(runtime.NumCPU()) {
		t.Errorf("Threads gave %d, want no more than the %d logical CPUs", n, runtime.NumCPU())
	}
	t.Logf("Threads gave %d of %d logical CPUs", n, runtime.NumCPU())
}

func TestDefaultThreads(t *testing.T) {
	if n := defaultThreads(); n < 1 {
		t.Errorf("defaultThreads gave %d, want at least 1", n)
	}
}

func TestThreadsForSize(t *testing.T) {
	const mib = 1 << 20
	tests := []struct {
		size  uint64
		limit int
		want  int
	}{
		{88 * mib, 10, 4},
		{774 * mib, 10, 9},
		{1056 * mib, 10, 10},
		{10 << 30, 10, 10},
		{88 * mib, 2, 2},
		{0, 0, 1},
	}
	for _, tt := range tests {
		if got := threadsForSize(tt.size, tt.limit); got != tt.want {
			t.Errorf("threadsForSize(%d, %d) gave %d, want %d", tt.size, tt.limit, got, tt.want)
		}
	}
}

func TestMoEActiveBytes(t *testing.T) {
	const size = 18 << 30
	qwen3moe := moeParams{experts: 128, used: 8, ffn: 768, embd: 2048, layers: 48}
	tests := []struct {
		name   string
		moe    moeParams
		params uint64
		want   uint64
	}{
		{"dense", moeParams{}, 1_700_000_000, size},
		{"qwen3-30b-a3b", qwen3moe, 30_532_000_000, size * 3_338_000_000 / 30_532_000_000},
		{"no ffn length", moeParams{experts: 64, used: 8}, 7_000_000_000, size / 8},
	}
	for _, tt := range tests {
		got := tt.moe.activeBytes(size, tt.params)
		if diff := float64(got)/float64(tt.want) - 1; diff > 0.01 || diff < -0.01 {
			t.Errorf("%s: activeBytes gave %d, want near %d", tt.name, got, tt.want)
		}
	}
}
