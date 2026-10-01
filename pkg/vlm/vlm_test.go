package vlm

import (
	"runtime"
	"strings"
	"testing"

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

	// A history longer than the first buffer must not panic or be cut short.
	long := strings.Repeat("a", 40000)
	out = vlm.ChatTemplate([]llama.ChatMessage{llama.NewChatMessage("user", long)}, true)
	if !strings.Contains(out, long) {
		t.Errorf("ChatTemplate returned %d bytes, want the whole %d byte message", len(out), len(long))
	}
}

func TestTokenPieceGrowsBuffer(t *testing.T) {
	modelFile := testModelFileName(t)
	mmprojFile := testMMProjFileName(t)

	testSetup(t)
	defer testCleanup(t)

	vlm := NewVLM(modelFile, mmprojFile)
	if err := vlm.Init(); err != nil {
		t.Fatalf("VLM.Init failed: %v", err)
	}
	defer vlm.Close()

	vocab := llama.ModelGetVocab(vlm.Model)
	tokens := llama.Tokenize(vocab, "information", false, false)
	if len(tokens) == 0 {
		t.Fatal("Tokenize returned no tokens")
	}

	want, _ := tokenPiece(vocab, tokens[0], make([]byte, 128))
	if len(want) < 2 {
		t.Skipf("piece %q is too short to test a small buffer", want)
	}
	got, buf := tokenPiece(vocab, tokens[0], make([]byte, 1))
	if got != want {
		t.Errorf("tokenPiece with a 1 byte buffer gave %q, want %q", got, want)
	}
	if len(buf) < len(want) {
		t.Errorf("tokenPiece returned a %d byte buffer, want at least %d", len(buf), len(want))
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

	bitmap := mtmd.BitmapInit(x, y, data)
	defer mtmd.BitmapFree(bitmap)

	if err := vlm.Tokenize(text, []mtmd.Bitmap{bitmap}, chunks); err != nil {
		// Accept both nil and not-nil error, as it may depend on model/projector
		t.Fatalf("VLM.Tokenize failed: %v", err)
	}
}

func TestVLM_Results(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Results test times out on macOS, skipping for now")
	}

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

	bitmap := mtmd.BitmapInit(x, y, data)
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
