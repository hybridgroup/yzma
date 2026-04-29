package message

// Message represents a single message in a conversation, which can be from the user, assistant, or system.
type Message interface {
	GetRole() string
	GetContent() map[string]interface{}
}
