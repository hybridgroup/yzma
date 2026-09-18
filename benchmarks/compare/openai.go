package compare

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// server talks to one OpenAI compatible model server. ollama and Docker Model
// Runner both give this interface, thus one client covers them.
type server struct {
	name    string
	baseURL string
	model   string
	client  *http.Client
}

// NewServer gives an engine that talks to a model server over HTTP.
func NewServer(name, baseURL, model string) Engine {
	return &server{
		name:    name,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		// A load of a large model can be slow, thus the timeout is generous.
		client: &http.Client{Timeout: 10 * time.Minute},
	}
}

func (s *server) Name() string { return s.name }

func (s *server) Close() {}

// Load asks for one token, which makes the server put the model in memory.
func (s *server) Load(req Request) error {
	if req.Embeddings {
		if _, err := s.Embed(req); err != nil {
			return fmt.Errorf("%s did not answer, is it running and does it have the model %q: %w",
				s.name, s.model, err)
		}

		return nil
	}

	warm := req
	warm.MaxTokens = 1

	if _, err := s.Generate(warm); err != nil {
		return fmt.Errorf("%s did not answer, is it running and does it have the model %q: %w",
			s.name, s.model, err)
	}

	return nil
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
	TopP        float64       `json:"top_p"`
	Seed        uint32        `json:"seed"`
	Stream      bool          `json:"stream"`
	StreamOpts  streamOptions `json:"stream_options"`

	// A model that thinks, such as Gemma 4, answers with its reasoning and
	// leaves the content empty. yzma applies the template of llama.cpp, which
	// does not think, thus the servers must not think either. ollama takes
	// reasoning_effort and Docker Model Runner takes chat_template_kwargs, and
	// each one leaves the field of the other alone.
	ReasoningEffort    string         `json:"reasoning_effort,omitempty"`
	ChatTemplateKwargs map[string]any `json:"chat_template_kwargs,omitempty"`
}

type streamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

type chatChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
			// A model that thinks, such as Gemma 4, puts its tokens here and
			// leaves the content empty. They are tokens all the same, thus the
			// timing and the count must see them. ollama calls the field
			// reasoning and Docker Model Runner calls it reasoning_content.
			Reasoning        string `json:"reasoning"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"delta"`
	} `json:"choices"`
	Usage *struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
	} `json:"usage"`
}

// Generate posts one streaming request. The stream is what gives a true time
// to the first token. A whole answer would only give the time of the request.
func (s *server) Generate(req Request) (Result, error) {
	// The clock starts here, not at the request. An application that uses a
	// server must turn the image into base64 and the whole message into JSON,
	// thus that work belongs to the cost of the REST interface. yzma pays for
	// the decoding of its image inside its own measurement.
	start := time.Now()

	content, err := messageContent(req)
	if err != nil {
		return Result{}, err
	}

	body, err := json.Marshal(chatRequest{
		Model:       s.model,
		Messages:    []chatMessage{{Role: "user", Content: content}},
		MaxTokens:   req.MaxTokens,
		Temperature: 0,
		TopP:        1,
		Seed:        req.Seed,
		Stream:      true,
		StreamOpts:  streamOptions{IncludeUsage: true},

		ReasoningEffort:    "none",
		ChatTemplateKwargs: map[string]any{"enable_thinking": false},
	})
	if err != nil {
		return Result{}, err
	}

	post, err := http.NewRequest(http.MethodPost, s.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	post.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(post)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return Result{}, fmt.Errorf("%s answered %s: %s", s.name, resp.Status, strings.TrimSpace(string(message)))
	}

	return readStream(resp.Body, start)
}

// readStream reads the server sent events of one answer.
func readStream(body io.Reader, start time.Time) (Result, error) {
	var result Result
	chunks := 0

	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}

		var chunk chatChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return result, fmt.Errorf("the stream has a frame that is not JSON: %w", err)
		}

		if chunk.Usage != nil {
			result.Tokens = chunk.Usage.CompletionTokens
			result.PromptTokens = chunk.Usage.PromptTokens
		}
		if len(chunk.Choices) == 0 {
			continue
		}

		piece := chunk.Choices[0].Delta.Content
		if piece == "" {
			piece = chunk.Choices[0].Delta.Reasoning
		}
		if piece == "" {
			piece = chunk.Choices[0].Delta.ReasoningContent
		}
		if piece == "" {
			continue
		}

		if chunks == 0 {
			result.FirstToken = time.Since(start)
		}
		chunks++
		result.Text += piece
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	result.Total = time.Since(start)

	// A server that sends no usage leaves the count of the frames, which is
	// near the count of the tokens but not the same.
	if result.Tokens == 0 {
		result.Tokens = chunks
	}
	if result.Tokens == 0 {
		return result, fmt.Errorf("the answer has no token")
	}

	return result, nil
}

// messageContent gives a plain string for the text suite and the two part form
// for an image. The image goes as a data URL of the same bytes that yzma reads.
func messageContent(req Request) (any, error) {
	if len(req.Image) == 0 {
		return req.Prompt, nil
	}

	url := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(req.Image)

	return []contentPart{
		{Type: "text", Text: req.Prompt},
		{Type: "image_url", ImageURL: &imageURL{URL: url}},
	}, nil
}

type embedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embedResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
	} `json:"usage"`
}

// Embed posts one request to the embeddings interface. The answer is a vector
// and not a stream, thus the time to the first token is the time of the whole
// request. An embedding gives no token back, thus almost all of that time is
// the round trip, which is what this suite measures.
func (s *server) Embed(req Request) (Result, error) {
	start := time.Now()

	body, err := json.Marshal(embedRequest{Model: s.model, Input: req.Prompt})
	if err != nil {
		return Result{}, err
	}

	post, err := http.NewRequest(http.MethodPost, s.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	post.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(post)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return Result{}, fmt.Errorf("%s answered %s: %s", s.name, resp.Status, strings.TrimSpace(string(message)))
	}

	var answer embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		return Result{}, err
	}
	if len(answer.Data) == 0 || len(answer.Data[0].Embedding) == 0 {
		return Result{}, fmt.Errorf("the vector is empty")
	}

	elapsed := time.Since(start)

	return Result{
		PromptTokens: answer.Usage.PromptTokens,
		Dimensions:   len(answer.Data[0].Embedding),
		FirstToken:   elapsed,
		Total:        elapsed,
	}, nil
}
