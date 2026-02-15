package vlm

import (
	"testing"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
)

func TestVLM_Init_Close(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	vlm := NewVLM(modelFile, mmprojFile)
	vlm.ModelParams = llama.ModelDefaultParams()
	vlm.ContextParams = llama.ContextDefaultParams()
	vlm.ProjectorParams = mtmd.ContextParamsDefault()
	if err := vlm.Init(); err != nil {
		t.Fatalf("VLM.Init failed: %v", err)
	}
	vlm.Close()
}

func TestVLM_ChatTemplate(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	vlm := NewVLM(modelFile, mmprojFile)
	vlm.ModelParams = llama.ModelDefaultParams()
	vlm.ContextParams = llama.ContextDefaultParams()
	vlm.ProjectorParams = mtmd.ContextParamsDefault()
	if err := vlm.Init(); err != nil {
		t.Fatalf("VLM.Init failed: %v", err)
	}
	defer vlm.Close()

	messages := []llama.ChatMessage{llama.NewChatMessage("user", "Hello")}
	out := vlm.ChatTemplate(messages, true)
	if out == "" {
		t.Error("ChatTemplate returned empty string")
	}
}

func TestVLM_Clear(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	vlm := NewVLM(modelFile, mmprojFile)
	vlm.ModelParams = llama.ModelDefaultParams()
	vlm.ContextParams = llama.ContextDefaultParams()
	vlm.ProjectorParams = mtmd.ContextParamsDefault()
	if err := vlm.Init(); err != nil {
		t.Fatalf("VLM.Init failed: %v", err)
	}
	defer vlm.Close()

	vlm.Clear() // Should not panic or error
}

func TestVLM_Tokenize(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	vlm := NewVLM(modelFile, mmprojFile)
	if err := vlm.Init(); err != nil {
		t.Fatalf("VLM.Init failed: %v", err)
	}
	defer vlm.Close()

	chunks := mtmd.InputChunksInit()
	defer mtmd.InputChunksFree(chunks)

	text := mtmd.NewInputText(mtmd.DefaultMarker()+"what is in this image?", true, true)

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open image file")
	}

	bitmap := mtmd.BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer mtmd.BitmapFree(bitmap)

	if err := vlm.Tokenize(text, []mtmd.Bitmap{bitmap}, chunks); err != nil {
		// Accept both nil and not-nil error, as it may depend on model/projector
		t.Fatalf("VLM.Tokenize failed: %v", err)
	}
}

func TestVLM_Results(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	vlm := NewVLM(modelFile, mmprojFile)
	if err := vlm.Init(); err != nil {
		t.Fatalf("VLM.Init failed: %v", err)
	}
	defer vlm.Close()

	chunks := mtmd.InputChunksInit()
	defer mtmd.InputChunksFree(chunks)

	text := mtmd.NewInputText(mtmd.DefaultMarker()+"what is in this image?", true, true)

	data, x, y, err := openImageFile("../../images/domestic_llama.jpg")
	if err != nil {
		t.Fatal("could not open image file")
	}

	bitmap := mtmd.BitmapInit(x, y, uintptr(unsafe.Pointer(&data[0])))
	defer mtmd.BitmapFree(bitmap)

	if err := vlm.Tokenize(text, []mtmd.Bitmap{bitmap}, chunks); err != nil {
		// Accept both nil and not-nil error, as it may depend on model/projector
		t.Fatalf("VLM.Tokenize failed: %v", err)
	}

	// This is a minimal call; actual results depend on model/projector and input
	if _, err := vlm.Results(chunks); err != nil {
		t.Fatalf("Results returned error): %v", err)
	}
}
