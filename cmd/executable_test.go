package cmd

import "testing"
import "context"

// TestExecutableExists validates binary exists
func TestExecutableExists(t *testing.T) {
	// negative test case so it passes on mac
	c := &ExecutableCheck{"Test", []string{"curling"}}
	mockCLICfg := &CLIConfigFlags{
		NodeIPStr: "127.0.0.1",
		Role:      "master",
		ForceTLS:  false,
	}

	err := c.executableExists(context.Background(), mockCLICfg)
	if err == nil {
		t.Fatalf("Should have returned an error")
	}
}
