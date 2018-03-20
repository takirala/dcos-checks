package common

import (
	"context"
	"fmt"
	"os"
)

// RunCheck is a helper function to run the check and emit the result.
func RunCheck(ctx context.Context, check DCOSChecker) {
	output, retCode, err := check.Run(ctx, DCOSConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing %s: %s\n", check.ID(), err)
	}

	if output != "" {
		fmt.Println(output)
	}

	os.Exit(retCode)
}
