package message

import (
	"reflect"
	"testing"
)

func TestToolDefinition_GetRole(t *testing.T) {
	td := ToolDefinition{
		Type: "function",
		Function: ToolFunctionDefinition{
			Name:        "add",
			Description: "Add two numbers",
			Parameters:  map[string]interface{}{"type": "object"},
		},
	}
	if got := td.GetRole(); got != "tool_definition" {
		t.Errorf("GetRole() = %q, want %q", got, "tool_definition")
	}
}

func TestToolDefinition_GetContent(t *testing.T) {
	params := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"a": map[string]interface{}{"type": "number"},
		},
		"required": []string{"a"},
	}
	td := ToolDefinition{
		Type: "function",
		Function: ToolFunctionDefinition{
			Name:        "add",
			Description: "Add two numbers",
			Parameters:  params,
		},
	}
	got := td.GetContent()
	want := map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "add",
			"description": "Add two numbers",
			"parameters":  params,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetContent() = %v, want %v", got, want)
	}
}

func TestToolMessage_GetRole(t *testing.T) {
	msg := Tool{Role: "tool", ToolCalls: nil}
	if got := msg.GetRole(); got != "tool" {
		t.Errorf("GetRole() = %q, want %q", got, "tool")
	}
}

func TestToolMessage_GetContent(t *testing.T) {
	msg := Tool{
		Role: "tool",
		ToolCalls: []ToolCall{
			{
				Type: "function",
				Function: ToolFunction{
					Name: "add",
					Arguments: map[string]string{
						"a": "1",
						"b": "2",
					},
				},
			},
		},
	}
	got := msg.GetContent()
	want := map[string]interface{}{
		"tool_calls": []map[string]interface{}{
			{
				"type": "function",
				"function": map[string]interface{}{
					"name": "add",
					"arguments": map[string]interface{}{
						"a": "1",
						"b": "2",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetContent() = %v, want %v", got, want)
	}
}

func TestToolResponseMessage_GetRole(t *testing.T) {
	msg := ToolResponse{Role: "tool_response", Name: "result", Content: "42"}
	if got := msg.GetRole(); got != "tool_response" {
		t.Errorf("GetRole() = %q, want %q", got, "tool_response")
	}
}

func TestToolResponseMessage_GetContent(t *testing.T) {
	msg := ToolResponse{Role: "tool_response", Name: "result", Content: "42"}
	got := msg.GetContent()
	want := map[string]interface{}{
		"name":    "result",
		"content": "42",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetContent() = %v, want %v", got, want)
	}
}
