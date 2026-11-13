package mtmd

import (
	"errors"
	"math"
	"testing"
)

func TestHelperResamplerZeroHandle(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if !hasResampler {
		if _, err := HelperResamplerInit(16000, 24000); !errors.Is(err, ErrNoResampler) {
			t.Fatalf("HelperResamplerInit without the symbols returned %v", err)
		}
		t.Skip("mtmd_helper_resampler not in this library")
	}

	HelperResamplerFree(0)
	if _, err := HelperResamplerProcess(0, []float32{0}); err == nil {
		t.Fatal("HelperResamplerProcess with zero handle did not fail")
	}
}

func TestHelperResamplerProcess(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	if !hasResampler {
		t.Skip("mtmd_helper_resampler not in this library")
	}

	r, err := HelperResamplerInit(48000, 16000)
	if err != nil {
		t.Fatalf("HelperResamplerInit failed: %v", err)
	}
	defer HelperResamplerFree(r)

	in := make([]float32, 4800)
	for i := range in {
		in[i] = float32(math.Sin(2 * math.Pi * 440 * float64(i) / 48000))
	}

	total := 0
	for range 10 {
		out, err := HelperResamplerProcess(r, in)
		if err != nil {
			t.Fatalf("HelperResamplerProcess failed: %v", err)
		}
		for _, v := range out {
			if math.IsNaN(float64(v)) || math.Abs(float64(v)) > 1.5 {
				t.Fatalf("output sample %f out of range", v)
			}
		}
		total += len(out)
	}

	want := len(in) * 10 / 3
	if total < want*9/10 || total > want {
		t.Fatalf("got %d output samples, expected about %d", total, want)
	}
}
