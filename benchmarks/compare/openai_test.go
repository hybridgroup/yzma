package compare

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const stream = `data: {"choices":[{"delta":{"role":"assistant"}}]}

data: {"choices":[{"delta":{"content":"Yes"}}]}

data: {"choices":[{"delta":{"content":", I am"}}]}

data: {"choices":[],"usage":{"completion_tokens":7}}

data: [DONE]
`

func TestReadStream(t *testing.T) {
	result, err := readStream(strings.NewReader(stream), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if result.Text != "Yes, I am" {
		t.Errorf("text = %q, want %q", result.Text, "Yes, I am")
	}
	// The count of the server wins over the count of the frames.
	if result.Tokens != 7 {
		t.Errorf("tokens = %d, want 7", result.Tokens)
	}
	if result.FirstToken <= 0 || result.FirstToken > result.Total {
		t.Errorf("first token = %v, total = %v", result.FirstToken, result.Total)
	}
}

// A server that sends no usage leaves the count of the frames.
func TestReadStreamWithoutUsage(t *testing.T) {
	raw := `data: {"choices":[{"delta":{"content":"a"}}]}

data: {"choices":[{"delta":{"content":"b"}}]}

data: [DONE]
`
	result, err := readStream(strings.NewReader(raw), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	if result.Tokens != 2 {
		t.Errorf("tokens = %d, want 2", result.Tokens)
	}
}

func TestReadStreamWithoutAnyToken(t *testing.T) {
	if _, err := readStream(strings.NewReader("data: [DONE]\n"), time.Now()); err == nil {
		t.Error("an empty answer must give an error")
	}
}

// Each run needs its own image, or a server answers the second run from the
// cache of the first while yzma does the whole work again.
func TestImageVariantsDiffer(t *testing.T) {
	made, err := ImageVariants("../../images/domestic_llama.jpg", 3)
	if err != nil {
		t.Fatal(err)
	}

	if len(made) != 3 {
		t.Fatalf("got %d images, want 3", len(made))
	}
	for i := range made {
		for j := i + 1; j < len(made); j++ {
			if bytes.Equal(made[i], made[j]) {
				t.Errorf("image %d and image %d are the same", i, j)
			}
		}
	}
}

// The image goes as a data URL of the bytes that yzma also reads.
func TestMessageContentPutsTheImageInTheMessage(t *testing.T) {
	made, err := ImageVariants("../../images/domestic_llama.jpg", 2)
	if err != nil {
		t.Fatal(err)
	}

	content, err := messageContent(Request{Prompt: ImagePrompt, Image: made[0]})
	if err != nil {
		t.Fatal(err)
	}

	parts, ok := content.([]contentPart)
	if !ok || len(parts) != 2 {
		t.Fatalf("content = %#v, want two parts", content)
	}
	if parts[0].Text != ImagePrompt {
		t.Errorf("text = %q, want %q", parts[0].Text, ImagePrompt)
	}
	if !strings.HasPrefix(parts[1].ImageURL.URL, "data:image/jpeg;base64,") {
		t.Errorf("the image is not a JPEG data URL: %.40s", parts[1].ImageURL.URL)
	}
}

func TestMessageContentOfTheTextSuite(t *testing.T) {
	content, err := messageContent(Request{Prompt: TextPrompt})
	if err != nil {
		t.Fatal(err)
	}

	if content != TextPrompt {
		t.Errorf("content = %#v, want %q", content, TextPrompt)
	}
}

// The sampler must be greedy on every engine, or the answers differ.
func TestGenerateAsksForGreedySampling(t *testing.T) {
	var got chatRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"ok\"}}]}\n\ndata: [DONE]\n"))
	}))
	defer srv.Close()

	engine := NewServer("test", srv.URL+"/v1", "a-model")
	if _, err := engine.Generate(Request{Prompt: TextPrompt, MaxTokens: 24, Seed: 1234}); err != nil {
		t.Fatal(err)
	}

	switch {
	case got.Temperature != 0:
		t.Errorf("temperature = %v, want 0", got.Temperature)
	case got.TopP != 1:
		t.Errorf("top_p = %v, want 1", got.TopP)
	case got.Seed != 1234:
		t.Errorf("seed = %v, want 1234", got.Seed)
	case got.MaxTokens != 24:
		t.Errorf("max_tokens = %v, want 24", got.MaxTokens)
	case !got.Stream:
		t.Error("the request must ask for a stream, or there is no time to the first token")
	case !got.StreamOpts.IncludeUsage:
		t.Error("the request must ask for the usage, or the count of the tokens is a guess")
	}
}

func TestGenerateReportsAnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "model not found", http.StatusNotFound)
	}))
	defer srv.Close()

	engine := NewServer("test", srv.URL+"/v1", "a-model")
	_, err := engine.Generate(Request{Prompt: TextPrompt, MaxTokens: 1})
	if err == nil || !strings.Contains(err.Error(), "model not found") {
		t.Errorf("err = %v, want the message of the server", err)
	}
}
