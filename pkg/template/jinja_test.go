package template

import (
	"os"
	"strings"
	"testing"

	"github.com/hybridgroup/yzma/pkg/message"
)

func TestChatMLTemplate(t *testing.T) {
	// Load the chatml.jinja template
	tmplPath := "./prompts/chatml.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare chat messages
	messages := []message.Message{
		message.Chat{
			Role:    "user",
			Content: "Hello, how are you?",
		},
		message.Chat{
			Role:    "assistant",
			Content: "I'm fine, thank you!",
		},
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)
	// Basic checks
	if !strings.Contains(result, "<|im_start|>user") || !strings.Contains(result, "<|im_start|>assistant") {
		t.Error("Output does not contain expected role markers")
	}
	if !strings.Contains(result, "Hello, how are you?") || !strings.Contains(result, "I'm fine, thank you!") {
		t.Error("Output does not contain expected message content")
	}
}

func TestQwen25InstructTemplateWithToolCall(t *testing.T) {
	// Load the qwen2.5-instruct.jinja template
	tmplPath := "./prompts/qwen2.5-instruct.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare messages with a tool call
	messages := []message.Message{
		message.Chat{
			Role:    "user",
			Content: "What is 2 + 3?",
		},
		message.Tool{
			Role: "assistant",
			ToolCalls: []message.ToolCall{
				{
					Type: "function",
					Function: message.ToolFunction{
						Name: "add",
						Arguments: map[string]string{
							"a": "2",
							"b": "3",
						},
					},
				},
			},
		},
		message.ToolResponse{
			Role:    "tool",
			Name:    "add",
			Content: "5",
		},
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// Basic checks
	if !strings.Contains(result, "<|im_start|>user") {
		t.Error("Output does not contain expected user role marker")
	}
	if !strings.Contains(result, "What is 2 + 3?") {
		t.Error("Output does not contain expected user message content")
	}
	if !strings.Contains(result, "<tool_call>") {
		t.Error("Output does not contain expected tool_call marker")
	}
	if !strings.Contains(result, "add") {
		t.Error("Output does not contain expected function name 'add'")
	}
	if !strings.Contains(result, "<tool_response>") {
		t.Error("Output does not contain expected tool_response marker")
	}
	if !strings.Contains(result, "5") {
		t.Error("Output does not contain expected tool response content '5'")
	}
}

func TestApplyJinjaTemplateWithToolMessage(t *testing.T) {
	// Load the qwen2.5-instruct.jinja template (supports tool calls)
	tmplPath := "./prompts/qwen2.5-instruct.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare messages with ToolMessage
	messages := []message.Message{
		message.Chat{
			Role:    "user",
			Content: "Call the calculator function",
		},
		message.Tool{
			Role: "assistant",
			ToolCalls: []message.ToolCall{
				{
					Type: "function",
					Function: message.ToolFunction{
						Name: "calculator",
						Arguments: map[string]string{
							"operation": "add",
							"x":         "10",
							"y":         "20",
						},
					},
				},
			},
		},
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// Basic checks
	if !strings.Contains(result, "<|im_start|>user") {
		t.Error("Output does not contain expected user role marker")
	}
	if !strings.Contains(result, "Call the calculator function") {
		t.Error("Output does not contain expected user message content")
	}
	if !strings.Contains(result, "<|im_start|>assistant") {
		t.Error("Output does not contain expected assistant role marker")
	}
	if !strings.Contains(result, "<tool_call>") {
		t.Error("Output does not contain expected tool_call marker")
	}
	if !strings.Contains(result, "calculator") {
		t.Error("Output does not contain expected function name 'calculator'")
	}
}

func TestApplyJinjaTemplateWithToolResponseMessage(t *testing.T) {
	// Load the qwen2.5-instruct.jinja template (supports tool calls)
	tmplPath := "./prompts/qwen2.5-instruct.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare messages with ToolResponseMessage
	messages := []message.Message{
		message.Chat{
			Role:    "user",
			Content: "What is the result?",
		},
		message.Tool{
			Role: "assistant",
			ToolCalls: []message.ToolCall{
				{
					Type: "function",
					Function: message.ToolFunction{
						Name: "calculator",
						Arguments: map[string]string{
							"x": "10",
							"y": "20",
						},
					},
				},
			},
		},
		message.ToolResponse{
			Role:    "tool",
			Name:    "calculator",
			Content: "30",
		},
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// Basic checks
	if !strings.Contains(result, "<|im_start|>user") {
		t.Error("Output does not contain expected user role marker")
	}
	if !strings.Contains(result, "What is the result?") {
		t.Error("Output does not contain expected user message content")
	}
	if !strings.Contains(result, "<tool_response>") {
		t.Error("Output does not contain expected tool_response marker")
	}
	if !strings.Contains(result, "30") {
		t.Error("Output does not contain expected tool response content '30'")
	}
}

func TestApplyJinjaTemplateWithMultipleToolCalls(t *testing.T) {
	// Load the qwen2.5-instruct.jinja template
	tmplPath := "./prompts/qwen2.5-instruct.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	// Prepare messages with multiple tool calls
	messages := []message.Message{
		message.Chat{
			Role:    "user",
			Content: "Calculate 2+3 and 5*7",
		},
		message.Tool{
			Role: "assistant",
			ToolCalls: []message.ToolCall{
				{
					Type: "function",
					Function: message.ToolFunction{
						Name: "add",
						Arguments: map[string]string{
							"a": "2",
							"b": "3",
						},
					},
				},
				{
					Type: "function",
					Function: message.ToolFunction{
						Name: "multiply",
						Arguments: map[string]string{
							"a": "5",
							"b": "7",
						},
					},
				},
			},
		},
		message.ToolResponse{
			Role:    "tool",
			Name:    "add",
			Content: "5",
		},
		message.ToolResponse{
			Role:    "tool",
			Name:    "multiply",
			Content: "35",
		},
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// Basic checks
	if !strings.Contains(result, "add") {
		t.Error("Output does not contain expected function name 'add'")
	}
	if !strings.Contains(result, "multiply") {
		t.Error("Output does not contain expected function name 'multiply'")
	}
	if !strings.Contains(result, "5") {
		t.Error("Output does not contain expected tool response '5'")
	}
	if !strings.Contains(result, "35") {
		t.Error("Output does not contain expected tool response '35'")
	}
}

func TestGemma3Template_BasicConversation(t *testing.T) {
	tmplPath := "./prompts/gemma3.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	messages := []message.Message{
		message.Chat{Role: "user", Content: "knock knock"},
		message.Chat{Role: "assistant", Content: "who is there"},
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	if !strings.Contains(result, "<start_of_turn>user") {
		t.Error("Output does not contain <start_of_turn>user")
	}
	if !strings.Contains(result, "<start_of_turn>model") {
		t.Error("Output does not contain <start_of_turn>model")
	}
	if !strings.Contains(result, "<end_of_turn>") {
		t.Error("Output does not contain <end_of_turn>")
	}
	if !strings.Contains(result, "knock knock") {
		t.Error("Output does not contain user message content")
	}
	if !strings.Contains(result, "who is there") {
		t.Error("Output does not contain assistant message content")
	}
	// Generation prompt should be appended
	if !strings.HasSuffix(strings.TrimRight(result, "\n"), "<start_of_turn>model") {
		t.Errorf("Output does not end with <start_of_turn>model, got: %q", result)
	}
}

func TestGemma3Template_SystemMessageMergedIntoFirstUserTurn(t *testing.T) {
	tmplPath := "./prompts/gemma3.jinja"
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}
	tmpl := string(tmplBytes)

	messages := []message.Message{
		message.Chat{Role: "system", Content: "Only reply like a pirate."},
		message.Chat{Role: "user", Content: "What is the answer to life the universe and everything?"},
	}

	result, err := Apply(tmpl, messages, true)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	t.Logf("Apply output:\n%s", result)

	// System content must appear inside the first user turn, not as its own turn.
	if strings.Contains(result, "<start_of_turn>system") {
		t.Error("Output must not contain a <start_of_turn>system turn")
	}
	if !strings.Contains(result, "Only reply like a pirate.") {
		t.Error("Output does not contain system instruction text")
	}
	if !strings.Contains(result, "What is the answer to life the universe and everything?") {
		t.Error("Output does not contain user message content")
	}
	// Both system and user text should be inside the same user turn.
	userTurnIdx := strings.Index(result, "<start_of_turn>user")
	systemIdx := strings.Index(result, "Only reply like a pirate.")
	userContentIdx := strings.Index(result, "What is the answer to life the universe and everything?")
	if systemIdx < userTurnIdx {
		t.Error("System instruction appears before the user turn marker")
	}
	if userContentIdx < systemIdx {
		t.Error("User content appears before system instruction within the user turn")
	}
}
