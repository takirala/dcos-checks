// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewVersionCheck returns an initialized instance of *ComponentCheck.
func NewVersionCheck(name string) DCOSChecker {
	check := &VersionCheck{Name: name}
	check.ClusterLeader = "leader.mesos"
	return check
}

// MasterListResponse response for leader.mesos/master.mesos
type MasterListResponse []struct {
	Host string `json:"host"`
	IP   string `json:"ip"`
}

// AgentListResponse response for /slaves
type AgentListResponse struct {
	Slaves []struct {
		ID         string `json:"id"`
		Hostname   string `json:"hostname"`
		Port       int    `json:"port"`
		Attributes struct {
			PublicIP string `json:"public_ip"`
		} `json:"attributes"`
	} `json:"slaves"`
	RecoveredSlaves []interface{} `json:"recovered_slaves"`
}

// VersionResponse responses /dcos-metadata/dcos-version.json
type VersionResponse struct {
	Version         string `json:"version"`
	DcosImageCommit string `json:"dcos-image-commit"`
	BootstrapID     string `json:"bootstrap-id"`
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Check DC/OS version of the cluster",
	Long: `Check dc/os version on each node in the cluster.
At any point there shouldnt be more than 2 versions that exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		RunCheck(context.TODO(), NewVersionCheck("DC/OS version check"))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

// VersionCheck struct
type VersionCheck struct {
	Name          string
	ClusterLeader string
}

// ID returns a unique check identifier.
func (vc *VersionCheck) ID() string {
	return vc.Name
}

// Run is running
func (vc *VersionCheck) Run(ctx context.Context, cfg *CLIConfigFlags) (string, int, error) {

	// List of masters
	var masterOpt URLFields
	masterOpt.host = vc.ClusterLeader
	masterOpt.port = mesosDNSPort
	masterOpt.path = "/v1/hosts/master.mesos"

	masterList, err := vc.ListOfMasters(cfg, masterOpt)
	if err != nil {
		return "", statusFailure, err
	}

	// List of agents
	var agentOpt URLFields
	agentOpt.host = vc.ClusterLeader
	agentOpt.port = mesosMasterHTTPPort
	agentOpt.path = "/slaves"

	agentList, err := vc.ListOfAgents(cfg, agentOpt)
	if err != nil {
		return "", statusFailure, err
	}

	// Check version endpoint for each endpoint
	var version map[string]bool
	version = make(map[string]bool)
	var versionURL URLFields
	versionURL.path = "/dcos-metadata/dcos-version.json"
	for _, master := range masterList {
		versionURL.host = master
		versionURL.port = 0
		if cfg.ForceTLS {
			versionURL.port = adminrouterMasterHTTPSPort
		}
		ver, err := vc.GetVersion(cfg, versionURL)
		if err != nil {
			return "", statusFailure, errors.Wrap(err, "Unable to get version")
		}
		version[ver] = true
	}

	for _, agent := range agentList {
		versionURL.host = agent
		versionURL.port = adminrouterAgentHTTPPort
		if cfg.ForceTLS {
			versionURL.port = adminrouterAgentHTTPSPort
		}
		ver, err := vc.GetVersion(cfg, versionURL)
		if err != nil {
			return "", statusFailure, errors.Wrap(err, "Unable to get version")
		}
		version[ver] = true
	}

	if len(version) <= 2 {
		return "", statusOK, nil
	}

	if len(version) > 2 {
		return "", statusWarning, errors.New("More than 2 dc/os versions on the cluster")
	}

	return "", statusUnknown, nil
}

// ListOfMasters returns the current list of masters in the cluster
func (vc *VersionCheck) ListOfMasters(cfg *CLIConfigFlags, urlopt URLFields) ([]string, error) {
	var masterResponse MasterListResponse
	var masterIPs []string
	_, response, err := HTTPRequest(cfg, urlopt)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch list of masters")
	}

	if err := json.Unmarshal(response, &masterResponse); err != nil {
		return nil, errors.Wrap(err, "Unable to unmarshal response")
	}

	for _, addr := range masterResponse {
		masterIPs = append(masterIPs, addr.IP)
	}
	return masterIPs, nil
}

// ListOfAgents returns the current list of agents in the cluster
func (vc *VersionCheck) ListOfAgents(cfg *CLIConfigFlags, urlopt URLFields) ([]string, error) {
	var agentResponse AgentListResponse
	var agentIPs []string
	_, response, err := HTTPRequest(cfg, urlopt)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch list of agents")
	}

	if err := json.Unmarshal(response, &agentResponse); err != nil {
		return nil, errors.Wrap(err, "Unable to unmarshal response")
	}

	for _, hosts := range agentResponse.Slaves {
		agentIPs = append(agentIPs, hosts.Hostname)
	}
	return agentIPs, nil
}

// GetVersion returns the dc/os version of a node
func (vc *VersionCheck) GetVersion(cfg *CLIConfigFlags, urlopt URLFields) (string, error) {
	var verResponse VersionResponse
	_, response, err := HTTPRequest(cfg, urlopt)
	if err != nil {
		return "", errors.Wrap(err, "Unable to get version")
	}

	if err := json.Unmarshal(response, &verResponse); err != nil {
		return "", errors.Wrap(err, "Unable to marshal response")
	}

	return verResponse.Version, nil
}
