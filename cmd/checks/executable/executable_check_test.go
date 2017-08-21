package executable

import "testing"
import (
	"context"

	"github.com/dcos/dcos-checks/common"
)

// TestExecutableExists validates binary exists
func TestExecutableExists(t *testing.T) {
	// negative test case so it passes on mac
	c := &executableCheck{"Test", []string{"curling"}}
	mockCLICfg := &common.CLIConfigFlags{
		NodeIPStr: "127.0.0.1",
		Role:      "master",
		ForceTLS:  false,
	}

	err := c.executableExists(context.Background(), mockCLICfg)
	if err == nil {
		t.Fatalf("Should have returned an error")
	}
}
