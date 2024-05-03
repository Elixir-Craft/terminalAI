package main

import (
	"os"
	"testing"
)

func TestGetPrompt(t *testing.T) {
	// Test case 1
	os.Args = []string{"", "Hello"}
	promptText, _ := getPrompt()
	if promptText != "Hello" {
		t.Errorf("Expected 'Hello' but got %s", promptText)
	}

}
