package download

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetModel(t *testing.T) {
	url := "https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/stories260K.gguf"
	dest := filepath.Join(t.TempDir(), "stories260K.gguf")

	ProgressTracker = nil
	err := GetModel(url, dest)
	if err != nil {
		t.Fatalf("GetModel() failed: %v", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatalf("Downloaded model file not found: %s", dest)
	}

	t.Logf("GetModel() successfully downloaded the model to: %s", dest)
}

func TestGetModelProgress(t *testing.T) {
	url := "https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/stories260K.gguf"
	dest := filepath.Join(t.TempDir(), "stories260K.gguf")

	ProgressTracker = DefaultProgressTracker()
	err := GetModel(url, dest)
	if err != nil {
		t.Fatalf("GetModel() failed: %v", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatalf("Downloaded model file not found: %s", dest)
	}

	t.Logf("GetModel() successfully downloaded the model to: %s", dest)
}

func TestDefaultModelsDir(t *testing.T) {
	dir := DefaultModelsDir()
	if dir == "" {
		t.Fatal("DefaultModelsDir should not return an empty string")
	}
}
