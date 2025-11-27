package message

type Message interface {
	GetRole() string
	GetContent() map[string]interface{}
}
