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

package mesosmetrics

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dcos/dcos-checks/client"
	"github.com/dcos/dcos-checks/common"
	"github.com/dcos/dcos-checks/constants"
	"github.com/dcos/dcos-go/dcos"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	nodeRecovered = 1.0
)

// mesosMetricsCheck checks if mesos replogs are synchronized by checking
// the value of /metrics/snapshot
type mesosMetricsCheck struct {
	Name    string
	urlFunc func(*http.Client, *common.CLIConfigFlags) (*url.URL, error)
}

// mesosMetricsCmd represents the mesos-metrics command
var mesosMetricsCmd = &cobra.Command{
	Use:   "mesos-metrics",
	Short: "Get the mesos metrics snapshot",
	Long:  `Metrics snapshot lets us know if the mesos rep logs are synchronized`,
	Run: func(cmd *cobra.Command, args []string) {
		common.RunCheck(context.TODO(), newMesosMetricsCheck("DC/OS metrics snapshot check"))
	},
}

// Add adds this command to the root command
func Add(root *cobra.Command) {
	root.AddCommand(mesosMetricsCmd)
}

// newMesosMetricsCheck returns an initialized instance of *mesosMetricsCheck.
func newMesosMetricsCheck(name string) common.DCOSChecker {
	check := &mesosMetricsCheck{Name: name}
	check.urlFunc = check.getURL
	return check
}

// ID returns a unique check identifier.
func (mm *mesosMetricsCheck) ID() string {
	return mm.Name
}

// Run invokes a check and return error output, exit code and error.
func (mm *mesosMetricsCheck) Run(ctx context.Context, cfg *common.CLIConfigFlags) (string, int, error) {
	type masterResponse struct {
		Recovered float64 `json:"registrar/log/recovered"`
	}

	type agentResponse struct {
		Recovered float64 `json:"slave/registered"`
		Output    string
	}

	httpClient, err := client.NewClient(cfg.IAMConfig, cfg.CACert)
	if err != nil {
		return "", constants.StatusUnknown, errors.Wrap(err, "Unable to create HTTP client")
	}

	url, err := mm.urlFunc(httpClient, cfg)
	if err != nil {
		return "", constants.StatusFailure, errors.Wrap(err, "Unable to get url")
	}

	logrus.Debugf("GET %s", url)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", constants.StatusUnknown, errors.Wrap(err, "Unable to create a new HTTP request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", constants.StatusUnknown, errors.Wrapf(err, "Unable to execute GET %s", url)
	}
	defer resp.Body.Close()

	if cfg.Role == dcos.RoleMaster {
		var jsonResponse masterResponse
		if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
			return "", constants.StatusUnknown, errors.Wrap(err, "Unable to unmarshal response")
		}

		if jsonResponse.Recovered == nodeRecovered {
			return "", constants.StatusOK, nil
		}
	} else {
		var jsonResponse agentResponse
		if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
			return "", constants.StatusUnknown, errors.Wrap(err, "Unable to unmarshal response")
		}

		if jsonResponse.Recovered == nodeRecovered {
			return "", constants.StatusOK, nil
		}
		return "", constants.StatusFailure, errors.New("Mesos replog not synchronized")
	}

	return "", constants.StatusUnknown, errors.New("Unable to run the check")
}

func (mm *mesosMetricsCheck) getURL(httpClient *http.Client, cfg *common.CLIConfigFlags) (*url.URL, error) {

	portsMap := map[string]int{
		dcos.RoleMaster:      constants.MesosMasterHTTPPort,
		dcos.RoleAgent:       constants.MesosAgentHTTPPort,
		dcos.RoleAgentPublic: constants.MesosAgentHTTPPort,
	}

	port, ok := portsMap[cfg.Role]
	if !ok {
		return nil, errors.Errorf("invalid role %s", cfg.Role)
	}

	scheme := constants.HTTPScheme
	if cfg.ForceTLS {
		scheme = constants.HTTPSScheme
	}

	ip, err := cfg.IP(httpClient)
	if err != nil {
		return nil, err
	}

	return &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(ip.String(), strconv.Itoa(port)),
		Path:   "/metrics/snapshot",
	}, nil
}
