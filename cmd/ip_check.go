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
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/dcos/dcos-go/exec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const defaultDetectIP = "/opt/mesosphere/bin/detect_ip"

var detectIP string

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Validate `detect_ip` output",
	Long:  `detect_ip is used to determine the node IP address.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		RunCheck(ctx, NewDetectIPCheck(detectIP))
	},
}

// NewDetectIPCheck returns a new instance of DetectIPCheck.
func NewDetectIPCheck(path string) DCOSChecker {
	return &DetectIPCheck{path}
}

// DetectIPCheck is a structure to accommodate detect_ip check.
type DetectIPCheck struct {
	Path string
}

// ID returns check ID.
func (d *DetectIPCheck) ID() string {
	return "detect_ip check " + d.Path
}

// Run executes the check.
func (d *DetectIPCheck) Run(ctx context.Context, cfg *CLIConfigFlags) (string, int, error) {
	if d.Path == "" {
		return "", 0, errors.New("path must be set")
	}

	stdout, stderr, code, err := exec.Output(ctx, d.Path)
	if err != nil {
		return "", 0, err
	}

	if code != 0 {
		return "", 0, errors.Wrapf(err, "return code non zero: %d", code)
	}

	if len(stderr) > 0 {
		return "", 0, errors.Errorf("detect_ip returned stderr: %s", string(stderr))
	}

	trimmedIP := bytes.TrimSpace(stdout)

	ip := net.ParseIP(string(trimmedIP))
	if ip == nil {
		return "", 0, errors.Errorf("invalid IP address %s", stdout)
	}

	return fmt.Sprintf("%s is a valid IPV4 address", ip), 0, nil
}

func init() {
	RootCmd.AddCommand(ipCmd)
	ipCmd.Flags().StringVarP(&detectIP, "detect-ip", "d", defaultDetectIP, "Set path to detect_ip script")
}
