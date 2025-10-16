package mtmd

import (
	"testing"
)

func TestInputChunksInitAndFree(t *testing.T) {
	chunks := InputChunksInit()
	if chunks == InputChunks(0) {
		t.Fatal("InputChunksInit returned an invalid InputChunks")
	}

	t.Log("InputChunksInit successfully initialized InputChunks")

	InputChunksFree(chunks)
	t.Log("InputChunksFree successfully freed InputChunks")
}

func TestInputChunksSize(t *testing.T) {
	chunks := InputChunksInit()
	defer InputChunksFree(chunks)

	size := InputChunksSize(chunks)
	if size != 0 {
		t.Fatalf("InputChunksSize returned a non-zero size for an empty InputChunks: %d", size)
	}

	t.Logf("InputChunksSize returned: %d", size)
}
