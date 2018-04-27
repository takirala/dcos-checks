package executable

import "testing"
import (
	"context"

	"github.com/dcos/dcos-checks/common"
)

func checkExecutable(e string) error {
	c := &executableCheck{"Test", []string{e}}
	mockCLICfg := &common.CLIConfigFlags{
		NodeIPStr: "127.0.0.1",
		Role:      "master",
		ForceTLS:  false,
	}
	return c.executableExists(context.Background(), mockCLICfg)
}

// TestExecutableExists validates binary exists
func TestExecutableExists(t *testing.T) {
	// negative test case
	executable := "nonexistent_executable"
	if err := checkExecutable(executable); err == nil {
		t.Fatalf("unexpectedly found executable '%s'", executable)
	}

	// positive test case
	executable = "bash"
	if err := checkExecutable(executable); err != nil {
		t.Fatalf("executable '%s' not found", executable)
	}
}
