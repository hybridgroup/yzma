package message

import (
	"reflect"
	"testing"
)

func TestChatMessage_GetRole(t *testing.T) {
	msg := Chat{Role: "user", Content: "Hello"}
	role := msg.GetRole()
	if role != "user" {
		t.Errorf("GetRole() = %q, want %q", role, "user")
	}
}

func TestChatMessage_GetContent(t *testing.T) {
	msg := Chat{Role: "assistant", Content: "Hi there!"}
	content := msg.GetContent()
	expected := map[string]interface{}{"content": "Hi there!"}
	if !reflect.DeepEqual(content, expected) {
		t.Errorf("GetContent() = %v, want %v", content, expected)
	}
}
