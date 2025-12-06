package main

import (
	"fmt"
	"strconv"

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
