package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/dcos/dcos-checks/client"
	"github.com/pkg/errors"
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

// URLFields is used to construct the url
type URLFields struct {
	host string
	port int
	path string
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

// HTTPRequest verifies the results of the request
func HTTPRequest(cfg *CLIConfigFlags, urlOptions URLFields) ([]byte, error) {
	httpClient, err := client.NewClient(cfg.IAMConfig, cfg.CACert)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create HTTP client")
	}

	url, err := getURL(httpClient, cfg, urlOptions)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("GET %s", url)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new HTTP request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to execute GET %s", url)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read response body")
	}
	defer resp.Body.Close()

	return responseData, nil
}

func getURL(httpClient *http.Client, cfg *CLIConfigFlags, urlOptions URLFields) (*url.URL, error) {
	scheme := httpScheme
	if cfg.ForceTLS {
		scheme = httpsScheme
	}

	if urlOptions.port == 0 {
		return &url.URL{
			Scheme: scheme,
			Host:   urlOptions.host,
			Path:   urlOptions.path,
		}, nil
	}
	return &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(urlOptions.host, strconv.Itoa(urlOptions.port)),
		Path:   urlOptions.path,
	}, nil
}
