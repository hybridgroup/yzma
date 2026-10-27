package mtmd

import (
	"errors"
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func TestHelperCapLayout(t *testing.T) {
	if size := unsafe.Sizeof(HelperCap{}); size != 8 {
		t.Fatalf("HelperCap size %d does not match C size 8", size)
	}
}

func TestHelperGetCapabilities(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	if !hasHelperGetCapabilities {
		if _, err := HelperGetCapabilities(0, 0); !errors.Is(err, ErrNoCapabilities) {
			t.Fatalf("HelperGetCapabilities without the symbol returned %v", err)
		}
		t.Skip("mtmd_helper_get_capabilities not in this library")
	}

	if caps, err := HelperGetCapabilities(0, 0); err != nil || caps != (HelperCap{}) {
		t.Fatalf("HelperGetCapabilities with zero model returned %+v, %v", caps, err)
	}

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer llama.ModelFree(model)

	textCaps, err := HelperGetCapabilities(model, 0)
	if err != nil {
		t.Fatalf("HelperGetCapabilities failed: %v", err)
	}
	t.Logf("text only capabilities: %+v", textCaps)
	if !textCaps.Chat || !textCaps.TextGen {
		t.Fatalf("text only model expected chat and text_gen, got %+v", textCaps)
	}
	if textCaps.Vision || textCaps.Audio || textCaps.Video {
		t.Fatalf("text only model expected no media input, got %+v", textCaps)
	}

	ctx, err := InitFromFile(mmprojFile, model, ContextParamsDefault())
	if err != nil {
		t.Fatalf("InitFromFile failed: %v", err)
	}
	defer Free(ctx)

	caps, err := HelperGetCapabilities(model, ctx)
	if err != nil {
		t.Fatalf("HelperGetCapabilities failed: %v", err)
	}
	t.Logf("projector capabilities: %+v", caps)
	if caps.Vision != SupportVision(ctx) {
		t.Fatalf("Vision %v does not match SupportVision %v", caps.Vision, SupportVision(ctx))
	}
	if caps.Audio != SupportAudio(ctx) {
		t.Fatalf("Audio %v does not match SupportAudio %v", caps.Audio, SupportAudio(ctx))
	}
	if caps.Video != SupportVideo(ctx) {
		t.Fatalf("Video %v does not match SupportVideo %v", caps.Video, SupportVideo(ctx))
	}
	if !caps.Chat || !caps.TextGen {
		t.Fatalf("expected chat and text_gen, got %+v", caps)
	}
}
