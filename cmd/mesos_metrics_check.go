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

package cmd

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dcos/dcos-checks/client"
	"github.com/dcos/dcos-go/dcos"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	nodeRecovered = 1.0
)

// mesosMetricsCmd represents the mesos-metrics command
var mesosMetricsCmd = &cobra.Command{
	Use:   "mesos-metrics",
	Short: "Get the mesos metrics snapshot",
	Long:  `Metrics snapshot lets us know if the mesos rep logs are synchronized`,
	Run: func(cmd *cobra.Command, args []string) {
		RunCheck(context.TODO(), NewMesosMetricsCheck("DC/OS metrics snapshot check"))
	},
}

// NewMesosMetricsCheck returns an initialized instance of *MesosMetricsCheck.
func NewMesosMetricsCheck(name string) DCOSChecker {
	check := &MesosMetricsCheck{Name: name}
	check.urlFunc = check.getURL
	return check
}

// MesosMetricsCheck checks if mesos replogs are synchronized by checking
// the value of /metrics/snapshot
type MesosMetricsCheck struct {
	Name    string
	urlFunc func(*http.Client, *CLIConfigFlags) (*url.URL, error)
}

// ID returns a unique check identifier.
func (mm *MesosMetricsCheck) ID() string {
	return mm.Name
}

func init() {
	RootCmd.AddCommand(mesosMetricsCmd)
}

// Run invokes a check and return error output, exit code and error.
func (mm *MesosMetricsCheck) Run(ctx context.Context, cfg *CLIConfigFlags) (string, int, error) {
	type masterResponse struct {
		Recovered float64 `json:"registrar/log/recovered"`
	}

	type agentResponse struct {
		Recovered float64 `json:"slave/registered"`
		Output    string
	}

	httpClient, err := client.NewClient(cfg.IAMConfig, cfg.CACert)
	if err != nil {
		return "", statusUnknown, errors.Wrap(err, "Unable to create HTTP client")
	}

	url, err := mm.urlFunc(httpClient, cfg)
	if err != nil {
		return "", statusFailure, errors.Wrap(err, "Unable to get url")
	}

	logrus.Debugf("GET %s", url)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", statusUnknown, errors.Wrap(err, "Unable to create a new HTTP request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", statusUnknown, errors.Wrapf(err, "Unable to execute GET %s", url)
	}
	defer resp.Body.Close()

	if cfg.Role == dcos.RoleMaster {
		var jsonResponse masterResponse
		if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
			return "", statusUnknown, errors.Wrap(err, "Unable to unmarshal response")
		}

		if jsonResponse.Recovered == nodeRecovered {
			return "", statusOK, nil
		}
	} else {
		var jsonResponse agentResponse
		if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
			return "", statusUnknown, errors.Wrap(err, "Unable to unmarshal response")
		}

		if jsonResponse.Recovered == nodeRecovered {
			return "", statusOK, nil
		}
		return "", statusFailure, errors.New("Mesos replog not synchronized")
	}

	return "", statusUnknown, errors.New("Unable to run the check")
}

func (mm *MesosMetricsCheck) getURL(httpClient *http.Client, cfg *CLIConfigFlags) (*url.URL, error) {

	portsMap := map[string]int{
		dcos.RoleMaster:      mesosMasterHTTPPort,
		dcos.RoleAgent:       mesosAgentHTTPPort,
		dcos.RoleAgentPublic: mesosAgentHTTPPort,
	}

	port, ok := portsMap[cfg.Role]
	if !ok {
		return nil, errors.Errorf("invalid role %s", cfg.Role)
	}

	scheme := httpScheme
	if cfg.ForceTLS {
		scheme = httpsScheme
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
