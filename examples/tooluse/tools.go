package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hybridgroup/yzma/pkg/message"
)

// Tool represents a tool definition for the LLM
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents a function definition
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

func getToolDefinitions() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "add",
				Description: "Add two numbers together",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"a": map[string]interface{}{
							"type":        "number",
							"description": "First number",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "Second number",
						},
					},
					"required": []string{"a", "b"},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "multiply",
				Description: "Multiply two numbers together",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"a": map[string]interface{}{
							"type":        "number",
							"description": "First number",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "Second number",
						},
					},
					"required": []string{"a", "b"},
				},
			},
		},
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "subtract",
				Description: "Subtract the second number from the first",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"a": map[string]interface{}{
							"type":        "number",
							"description": "First number",
						},
						"b": map[string]interface{}{
							"type":        "number",
							"description": "Second number to subtract",
						},
					},
					"required": []string{"a", "b"},
				},
			},
		},
	}
}

// executeToolCall executes a tool call and returns the result
func executeToolCall(call message.ToolCall) (string, error) {
	switch call.Function.Name {
	case "add":
		a, err := strconv.ParseFloat(call.Function.Arguments["a"], 64)
		if err != nil {
			return "", fmt.Errorf("invalid argument 'a': %v", err)
		}
		b, err := strconv.ParseFloat(call.Function.Arguments["b"], 64)
		if err != nil {
			return "", fmt.Errorf("invalid argument 'b': %v", err)
		}
		result := a + b
		return fmt.Sprintf("%.2f", result), nil

	case "multiply":
		a, err := strconv.ParseFloat(call.Function.Arguments["a"], 64)
		if err != nil {
			return "", fmt.Errorf("invalid argument 'a': %v", err)
		}
		b, err := strconv.ParseFloat(call.Function.Arguments["b"], 64)
		if err != nil {
			return "", fmt.Errorf("invalid argument 'b': %v", err)
		}
		result := a * b
		return fmt.Sprintf("%.2f", result), nil

	case "subtract":
		a, err := strconv.ParseFloat(call.Function.Arguments["a"], 64)
		if err != nil {
			return "", fmt.Errorf("invalid argument 'a': %v", err)
		}
		b, err := strconv.ParseFloat(call.Function.Arguments["b"], 64)
		if err != nil {
			return "", fmt.Errorf("invalid argument 'b': %v", err)
		}
		result := a - b
		return fmt.Sprintf("%.2f", result), nil

	default:
		return "", fmt.Errorf("unknown function: %s", call.Function.Name)
	}
}

// parseToolCalls attempts to parse tool calls from the model's response
// This is a simplified parser - real implementations would be more robust
func parseToolCalls(response string) []message.ToolCall {
	var calls []message.ToolCall

	// Look for <tool_call> tags
	start := strings.Index(response, "<tool_call>")
	end := strings.Index(response, "</tool_call>")

	for start != -1 && end != -1 && start < end {
		content := response[start+len("<tool_call>") : end]
		content = strings.TrimSpace(content)

		// Parse JSON content
		var parsed struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}

		if err := json.Unmarshal([]byte(content), &parsed); err == nil {
			calls = append(calls, message.ToolCall{
				Type: "function",
				Function: message.ToolFunction{
					Name:      parsed.Name,
					Arguments: parsed.Arguments,
				},
			})
		}

		// Look for more tool calls
		response = response[end+len("</tool_call>"):]
		start = strings.Index(response, "<tool_call>")
		end = strings.Index(response, "</tool_call>")
	}

	return calls
}
