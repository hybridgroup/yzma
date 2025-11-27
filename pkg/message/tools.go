package message

// ToolFunction represents a function within a tool call.
type ToolFunction struct {
	Name      string
	Arguments map[string]string
}

// ToolCall represents a call to a tool function within a tool message.
type ToolCall struct {
	Type     string
	Function ToolFunction
}

// Tool represents a message that contains tool calls.
type Tool struct {
	Role      string
	ToolCalls []ToolCall
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
	return map[string]interface{}{
		"tool_calls": calls,
	}
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
