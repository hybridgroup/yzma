package download

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetModel(t *testing.T) {
	url := "https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/stories260K.gguf"
	dest := filepath.Join(t.TempDir(), "stories260K.gguf")

	err := GetModel(url, dest, false)
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

	err := GetModel(url, dest, true)
	if err != nil {
		t.Fatalf("GetModel() failed: %v", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatalf("Downloaded model file not found: %s", dest)
	}

	t.Logf("GetModel() successfully downloaded the model to: %s", dest)
}
