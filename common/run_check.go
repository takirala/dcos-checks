package common

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// RunCheck is a helper function to run the check and emit the result.
func RunCheck(ctx context.Context, check DCOSChecker) {
	output, retCode, err := check.Run(ctx, DCOSConfig)
	if err != nil {
		logrus.Fatalf("Error executing %s: %s", check.ID(), err)
	}

	if output != "" {
		fmt.Println(output)
	}

	os.Exit(retCode)
}
