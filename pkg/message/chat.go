package message

// Chat represents a standard chat message with a role and content.
type Chat struct {
	Role    string
	Content string
}

// GetRole returns the role of the chat message.
func (cm Chat) GetRole() string {
	return cm.Role
}

// GetContent returns the content of the chat message as a map.
func (cm Chat) GetContent() map[string]interface{} {
	return map[string]interface{}{
		"content": cm.Content,
	}
}
