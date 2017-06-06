package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	statusOK = iota
	statusWarning
	statusFailure
	statusUnknown

	// exhibitor admin router port
	exhibitorPort = 8181

	// master node has a 3dt instance running on TCP port 1050.
	// ee version has 3dt running via unix socket on both master and agent nodes,
	// depending on security option. Ports 80 or 443 are using accordingly.
	dcosDiagnosticsMasterHTTPPort = 1050
	adminrouterMasterHTTPSPort    = 443

	// agent node runs 3dt via unix socket and is available though the agent
	// adminrouter HTTP TCP port 61001 or HTTPS 61002.
	adminrouterAgentHTTPPort  = 61001
	adminrouterAgentHTTPSPort = 61002

	mesosMasterHTTPPort = 5050
	mesosAgentHTTPPort  = 5051
	mesosDNSPort        = 8123

	httpScheme  = "http"
	httpsScheme = "https"
)

// DCOSChecker defines an interface for a generic DC/OS check.
// ID() returns a check unique ID and RunCheck(...) returns a combined stdout/stderr, exit code and error.
type DCOSChecker interface {
	ID() string
	Run(context.Context, *CLIConfigFlags) (string, int, error)
}

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
