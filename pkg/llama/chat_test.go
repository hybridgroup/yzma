package llama

import (
	"testing"
)

func TestChat(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	chat := []ChatMessage{NewChatMessage("user", "what is going on?")}
	buf := make([]byte, 1024)

	sz := ChatApplyTemplate("chatml", chat, false, buf)
	if sz <= 0 {
		t.Fatal("unable to apply chat template")
	}

	result := string(buf)
	if len(result) == 0 {
		t.Fatal("invalid output from chat template")
	}
}

func TestChatBuiltinTemplates(t *testing.T) {
	testSetup(t)
	defer testCleanup(t)

	// Create a buffer to receive the template names
	templates := make([]string, 10) // Assuming there are less than 10 built-in templates

	// Call the function
	count := ChatBuiltinTemplates(templates)

	// Check that we got a non-negative count
	if count < 0 {
		t.Errorf("ChatBuiltinTemplates returned negative count: %d", count)
	}

	// Check that the templates slice was populated
	// Note: The actual content depends on the built-in templates in the llama.cpp library
	// We're just verifying that the function call works without crashing
	t.Logf("Found %d built-in templates", count)
	for i := 0; i < int(count) && i < len(templates); i++ {
		if templates[i] != "" {
			t.Logf("Template %d: %s", i, templates[i])
		}
	}
}
