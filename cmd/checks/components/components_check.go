// Copyright © 2017 Mesosphere Inc. <http://mesosphere.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package components

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dcos/dcos-checks/client"
	"github.com/dcos/dcos-checks/common"
	"github.com/dcos/dcos-checks/constants"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// componentCheck validates that all systemd units are healthy by making a GET request
// to dcos-diagnostics endpoint /system/health/v1 on the localhost.
// In open DC/OS 3dt listens port 1050 on master nodes. On agent nodes, 3dt uses socket activation to bind on
// unix socket. Adminrouter is used to make a reverse proxy.
type componentCheck struct {
	Name string
}

var (
	healthURLPrefix   string
	scheme            string
	port              int
	excludeComponents []string
)

// componentsCmd represents the systemd health check
var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Check DC/OS components",
	Long: `Check DC/OS components health by making a GET request to dcos-3dt service
and validating the health field:

/system/health/v1 is the local endpoint. The response structure is the following
{
  "units": ["unit1", ...]
}
`,
	Run: func(cmd *cobra.Command, args []string) {
		common.RunCheck(context.TODO(),
			&componentCheck{"DC/OS components health check"})
	},
}

// Add adds this command to the root command
func Add(root *cobra.Command) {
	root.AddCommand(componentsCmd)
	componentsCmd.Flags().StringVarP(&healthURLPrefix, "health-url", "u", "/system/health/v1", "Set dcos-diagnostics health url")
	componentsCmd.Flags().StringVarP(&scheme, "scheme", "s", "http", "Set dcos-diagnostics health url scheme")
	componentsCmd.Flags().IntVarP(&port, "port", "p", 1050, "Set TCP port")
	componentsCmd.Flags().StringArrayVarP(&excludeComponents, "exclude", "e", nil, "Exclude components from health check")
}

// Run invokes a systemd check and return error output, exit code and error.
func (c *componentCheck) Run(ctx context.Context, cfg *common.CLIConfigFlags) (string, int, error) {
	httpClient, err := client.NewClient(cfg.IAMConfig, cfg.CACert)
	if err != nil {
		return "", constants.StatusUnknown, errors.Wrap(err, "unable to create HTTP client")
	}

	url, err := c.getHealthURL(httpClient, healthURLPrefix, scheme, port, cfg)
	if err != nil {
		return "", 0, err
	}
	logrus.Debugf("GET %s", url)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", constants.StatusUnknown, errors.Wrap(err, "unable to create a new HTTP request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", constants.StatusUnknown, errors.Wrapf(err, "unable to execute GET %s", healthURLPrefix)
	}
	defer resp.Body.Close()

	var dr diagnosticsResponse
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return "", constants.StatusUnknown, errors.Wrap(err, "unable to unmarshal diagnostics response")
	}

	errorList, retCode := dr.checkHealth(excludeComponents)
	return strings.Join(errorList, "\n"), retCode, nil
}

// ID returns a unique check identifier.
func (c *componentCheck) ID() string {
	return c.Name
}

func (c *componentCheck) getHealthURL(httpClient *http.Client, path, scheme string, port int, cfg *common.CLIConfigFlags) (*url.URL, error) {
	ip, err := cfg.IP(httpClient)
	if err != nil {
		return nil, err
	}

	return &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(ip.String(), strconv.Itoa(port)),
		Path:   path,
	}, nil
}
