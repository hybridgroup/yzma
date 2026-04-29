package message

// ToolDefinition represents a single tool available to the model.
type ToolDefinition struct {
	Type     string                 `json:"type"`
	Function ToolFunctionDefinition `json:"function"`
}

// GetRole returns the role of the tool definition message.
func (td ToolDefinition) GetRole() string {
	return "tool_definition"
}

// GetContent returns the tool definition as a map.
func (td ToolDefinition) GetContent() map[string]interface{} {
	return map[string]interface{}{
		"type": td.Type,
		"function": map[string]interface{}{
			"name":        td.Function.Name,
			"description": td.Function.Description,
			"parameters":  td.Function.Parameters,
		},
	}
}

// ToolFunctionDefinition represents a function definition
type ToolFunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// Tool represents a message that contains tool calls.
// Content may optionally hold any spoken text that was generated alongside the
// tool calls so that conversation history preserves both.
type Tool struct {
	Role      string
	Content   string
	ToolCalls []ToolCall
}

// ToolCall represents a call to a tool function within a tool message.
type ToolCall struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents a function within a tool call.
type ToolFunction struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}

// GetRole returns the role of the tool message.
func (tm Tool) GetRole() string {
	return tm.Role
}

// GetContent returns the content of the tool message as a map.
func (tm Tool) GetContent() map[string]interface{} {
	calls := make([]map[string]interface{}, len(tm.ToolCalls))
	for i, call := range tm.ToolCalls {
		args := make(map[string]interface{})
		for k, v := range call.Function.Arguments {
			args[k] = v
		}
		calls[i] = map[string]interface{}{
			"type": call.Type,
			"function": map[string]interface{}{
				"name":      call.Function.Name,
				"arguments": args,
			},
		}
	}
	m := map[string]interface{}{
		"tool_calls": calls,
	}
	if tm.Content != "" {
		m["content"] = tm.Content
	}
	return m
}

// ToolResponse represents a message that contains a tool response.
type ToolResponse struct {
	Role    string
	Name    string
	Content string
}

// GetRole returns the role of the tool response message.
func (trm ToolResponse) GetRole() string {
	return trm.Role
}

// GetContent returns the content of the tool response message as a map.
func (trm ToolResponse) GetContent() map[string]interface{} {
	return map[string]interface{}{
		"name":    trm.Name,
		"content": trm.Content,
	}
}
