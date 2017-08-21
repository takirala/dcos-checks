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

package version

import (
	"context"
	"encoding/json"

	"github.com/dcos/dcos-checks/common"
	"github.com/dcos/dcos-checks/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// versionCheck struct
type versionCheck struct {
	Name          string
	ClusterLeader string
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Check DC/OS version of the cluster",
	Long: `Check dc/os version on each node in the cluster.
At any point there shouldnt be more than 2 versions that exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.RunCheck(context.TODO(), newVersionCheck("DC/OS version check"))
	},
}

// Add adds this command to the root command
func Add(root *cobra.Command) {
	root.AddCommand(versionCmd)
}

// newVersionCheck returns an initialized instance of *versionCheck.
func newVersionCheck(name string) *versionCheck {
	check := &versionCheck{Name: name}
	check.ClusterLeader = "leader.mesos"
	return check
}

// ID returns a unique check identifier.
func (vc *versionCheck) ID() string {
	return vc.Name
}

// Run is running
func (vc *versionCheck) Run(ctx context.Context, cfg *common.CLIConfigFlags) (string, int, error) {

	// List of masters
	var masterOpt common.URLFields
	masterOpt.Host = vc.ClusterLeader
	masterOpt.Port = constants.MesosDNSPort
	masterOpt.Path = "/v1/hosts/master.mesos"

	masterList, err := vc.ListOfMasters(cfg, masterOpt)
	if err != nil {
		return "", constants.StatusFailure, err
	}

	// List of agents
	var agentOpt common.URLFields
	agentOpt.Host = vc.ClusterLeader
	agentOpt.Port = constants.MesosMasterHTTPPort
	agentOpt.Path = "/slaves"

	agentList, err := vc.ListOfAgents(cfg, agentOpt)
	if err != nil {
		return "", constants.StatusFailure, err
	}

	// Check version endpoint for each endpoint
	var version map[string]bool
	version = make(map[string]bool)
	var versionURL common.URLFields
	versionURL.Path = "/dcos-metadata/dcos-version.json"
	for _, master := range masterList {
		versionURL.Host = master
		versionURL.Port = 0
		if cfg.ForceTLS {
			versionURL.Port = constants.AdminrouterMasterHTTPSPort
		}
		ver, err := vc.GetVersion(cfg, versionURL)
		if err != nil {
			return "", constants.StatusFailure, errors.Wrap(err, "Unable to get version")
		}
		version[ver] = true
	}

	for _, agent := range agentList {
		versionURL.Host = agent
		versionURL.Port = constants.AdminrouterAgentHTTPPort
		if cfg.ForceTLS {
			versionURL.Port = constants.AdminrouterAgentHTTPSPort
		}
		ver, err := vc.GetVersion(cfg, versionURL)
		if err != nil {
			return "", constants.StatusFailure, errors.Wrap(err, "Unable to get version")
		}
		version[ver] = true
	}

	if len(version) <= 2 {
		return "", constants.StatusOK, nil
	}

	if len(version) > 2 {
		return "", constants.StatusWarning, errors.New("More than 2 dc/os versions on the cluster")
	}

	return "", constants.StatusUnknown, nil
}

// ListOfMasters returns the current list of masters in the cluster
func (vc *versionCheck) ListOfMasters(cfg *common.CLIConfigFlags, urlopt common.URLFields) ([]string, error) {
	var masterResponse masterListResponses
	var masterIPs []string
	_, response, err := common.HTTPRequest(cfg, urlopt)
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
func (vc *versionCheck) ListOfAgents(cfg *common.CLIConfigFlags, urlopt common.URLFields) ([]string, error) {
	var agentResponse agentListResponse
	var agentIPs []string
	_, response, err := common.HTTPRequest(cfg, urlopt)
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
func (vc *versionCheck) GetVersion(cfg *common.CLIConfigFlags, urlopt common.URLFields) (string, error) {
	var verResponse versionResponse
	_, response, err := common.HTTPRequest(cfg, urlopt)
	if err != nil {
		return "", errors.Wrap(err, "Unable to get version")
	}

	if err := json.Unmarshal(response, &verResponse); err != nil {
		return "", errors.Wrap(err, "Unable to marshal response")
	}

	return verResponse.Version, nil
}
