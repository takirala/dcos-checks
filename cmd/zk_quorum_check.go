// Copyright © 2017 Mesosphere Inc. <www.mesosphere.com>
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

const zkResponseServing = 3

// NewZkQuorumCheck returns an initialized instance of *ComponentCheck.
func NewZkQuorumCheck(name string) DCOSChecker {
	check := &ZkQuorumCheck{Name: name}
	check.urlFunc = check.getURL
	return check
}

// zkQuorumCmd represents the zk-quorum command
var zkQuorumCmd = &cobra.Command{
	Use:   "zk-quorum",
	Short: "DC/OS check zk in quorum",
	Long:  `DC/OS check to verify that the node exhibitor is up and serving`,
	Run: func(cmd *cobra.Command, args []string) {
		RunCheck(context.TODO(), NewZkQuorumCheck("DC/OS zk quorum check"))
	},
}

func init() {
	RootCmd.AddCommand(zkQuorumCmd)
}

// ZkResponse struct response for localhost:8181/exhibitor/v1/cluster/state/<host>
type ZkResponse struct {
	Response struct {
		Switches struct {
			Restarts bool `json:"restarts"`
			Cleanup  bool `json:"cleanup"`
			Backups  bool `json:"backups"`
		} `json:"switches"`
		State       int    `json:"state"`
		Description string `json:"description"`
		IsLeader    bool   `json:"isLeader"`
	} `json:"response"`
	ErrorMessage string `json:"errorMessage"`
	Success      bool   `json:"success"`
}

// ZkQuorumCheck struct
type ZkQuorumCheck struct {
	Name    string
	urlFunc func(*http.Client, *CLIConfigFlags) (*url.URL, error)
}

// Run invokes a zkquorum check and return error output, exit code and error.
func (zk *ZkQuorumCheck) Run(ctx context.Context, cfg *CLIConfigFlags) (string, int, error) {
	if cfg.Role != dcos.RoleMaster {
		return "", statusFailure, errors.New("Check can be run only on masters")
	}

	httpClient, err := client.NewClient(cfg.IAMConfig, cfg.CACert)
	if err != nil {
		return "", statusUnknown, errors.Wrap(err, "unable to create HTTP client")
	}

	url, err := zk.urlFunc(httpClient, cfg)
	if err != nil {
		return "", statusFailure, errors.Wrap(err, "Unable to get the zk status endpoint")
	}
	logrus.Debugf("GET %s", url)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", statusUnknown, errors.Wrap(err, "unable to create a new HTTP request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", statusUnknown, errors.Wrapf(err, "unable to execute GET %s", url)
	}
	defer resp.Body.Close()

	var zr ZkResponse
	if err := json.NewDecoder(resp.Body).Decode(&zr); err != nil {
		return "", statusUnknown, errors.Wrap(err, "unable to unmarshal exhibitor response")
	}

	if zr.Response.State == zkResponseServing {
		// the host is serving
		return zr.Response.Description, statusOK, errors.Wrap(nil, zr.ErrorMessage)
	}

	return zr.Response.Description, statusFailure, errors.Wrap(nil, zr.ErrorMessage)
}

// ID returns a unique check identifier.
func (zk *ZkQuorumCheck) ID() string {
	return zk.Name
}

func (zk *ZkQuorumCheck) getURL(httpClient *http.Client, cfg *CLIConfigFlags) (*url.URL, error) {
	ip, err := cfg.IP(httpClient)
	if err != nil {
		return nil, err
	}

	scheme := httpScheme
	if cfg.ForceTLS {
		scheme = httpsScheme
	}

	path := "exhibitor/v1/cluster/state/" + ip.String()

	return &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(ip.String(), strconv.Itoa(exhibitorPort)),
		Path:   path,
	}, nil

}
