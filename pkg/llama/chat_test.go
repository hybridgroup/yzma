package llama

import (
	"slices"
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

	templates := ChatBuiltinTemplates()

	t.Logf("Found %d built-in templates\n", len(templates))
	if len(templates) == 0 {
		t.Fatal("no built-in templates found")
	}

	t.Logf("Template %v\n", templates)

	existingTemplates := []string{"deepseek3", "gemma", "gpt-oss", "llama3", "llama4"}
	for _, existingTemplate := range existingTemplates {
		if !slices.Contains(templates, existingTemplate) {
			t.Fatalf("missing expected template: %s\n", existingTemplate)
		}
	}
}
