package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dcos/dcos-checks/common"
	"github.com/dcos/dcos-checks/constants"
	"github.com/dcos/dcos-go/dcos"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	localhost string = "127.0.0.1"
)

type arguments map[string]bool

var validArgs = arguments{
	"dcos-registry":   	true,
	"dcos-registry-dns":true,
}

func (arguments) list() []string {
	valid := make([]string, 0)
	for k := range validArgs {
		valid = append(valid, k)
	}
	return valid
}

// httpCmd represents the HTTP proxy command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Checks for component HTTP servers running correctly",
	Long: `DC/OS checks to confirm that DC/OS component HTTP servers are up and responding correctly.
Usage:
http dcos-registry
http dcos-registry-dns
`,
	Run: func(cmd *cobra.Command, args []string) {
		common.RunCheck(context.TODO(),
			newHTTPCheck("DC/OS checks for component HTTP servers", args))
	},
}

// Register adds this command to the root command
func Register(root *cobra.Command) {
	root.AddCommand(httpCmd)
}

type httpCheck struct {
	Name     string
	Args     []string
}

// newHTTPCheck returns an initialized instance of *httpCheck
func newHTTPCheck(name string, args []string) *httpCheck {
	return &httpCheck{
		Name:     name,
		Args:     args,
	}
}

// ID returns a unique check identifier.
func (c *httpCheck) ID() string {
	return c.Name
}

// Run runs the specified checks
func (c *httpCheck) Run(ctx context.Context, cfg *common.CLIConfigFlags) (string, int, error) {
	if len(c.Args) != 1 {
		return "", constants.StatusFailure, fmt.Errorf("Provide one argument only, valid args %v", validArgs.list())
	}
	if _, ok := validArgs[c.Args[0]]; !ok {
		return "", constants.StatusFailure, fmt.Errorf("Option not supported, valid args %v", validArgs.list())
	}
	if cfg.Role != dcos.RoleMaster {
		return "", constants.StatusFailure, errors.New("Check can be run only on masters")
	}
	switch c.Args[0] {
	case "dcos-registry":
		req, err := http.NewRequest("GET", portURL(localhost, 5001, "/", false).String(), nil)
		if err != nil {
			return "", constants.StatusUnknown, errors.Wrap(err, "Unable to create HTTP request")
		}
		return c.checkHTTP(cfg, insecureHttpClient(), req, http.StatusOK)

	case "dcos-registry-dns":
		req, err := http.NewRequest("GET", portURL("registry.component.thisdcos.directory", 443, "/", true).String(), nil)
		if err != nil {
			return "", constants.StatusUnknown, errors.Wrap(err, "Unable to create HTTP request")
		}
		return c.checkHTTP(cfg, insecureHttpClient(), req, http.StatusOK)

	default:
	}
	return "", constants.StatusUnknown, errors.Errorf("Unable not find matching check for arg %s", c.Args[0])
}

func insecureHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func portURL(host string, port int, path string, tls bool) *url.URL {
	scheme := "http"
	if tls {
		scheme = "https"
	}
	return &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   path,
	}
}


func (c *httpCheck) checkHTTP(cfg *common.CLIConfigFlags, httpClient *http.Client, req *http.Request, expectedStatus int) (string, int, error) {
	componentStatus, err := httpClient.Do(req)
	if err != nil {
		return "", constants.StatusUnknown, errors.Wrapf(err, "Unable to fetch %s status", req.URL.String())
	}
	httpResponse := http.StatusText(componentStatus.StatusCode)
	if componentStatus.StatusCode != expectedStatus {
		output := fmt.Sprintf("HTTP Server: %d %s", componentStatus.StatusCode, httpResponse)
		if cfg.Verbose {
			output = fmt.Sprintf("%s\n%d Expected status\n%d %s", output, expectedStatus, componentStatus.StatusCode, req.URL.String())
		}
		return output, constants.StatusFailure, nil
	}
	return httpResponse, constants.StatusOK, nil
}