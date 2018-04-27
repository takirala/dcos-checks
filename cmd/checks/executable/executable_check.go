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

package executable

import (
	"context"
	"fmt"

	"github.com/dcos/dcos-checks/common"
	"github.com/dcos/dcos-checks/constants"
	"github.com/dcos/dcos-go/exec"
	"github.com/spf13/cobra"
)

// executableCmd represents the executable command
var executableCmd = &cobra.Command{
	Use:   "executable",
	Short: "Check for the availability of an executable",
	Long:  "Check for the availability of an executable",
	Run: func(cmd *cobra.Command, args []string) {
		common.RunCheck(context.TODO(),
			newExecutableCheck("check availability of executable", args))
	},
}

// Register adds this command to the root command
func Register(root *cobra.Command) {
	root.AddCommand(executableCmd)
}

// newExecutableCheck returns an intialized instance of *executableCheck
func newExecutableCheck(name string, args []string) *executableCheck {
	return &executableCheck{
		Name: name,
		Args: args,
	}
}

// executableCheck validates we have the required executable to install/run DC/OS
type executableCheck struct {
	Name string
	Args []string
}

// ID returns a unique check identifier.
func (c *executableCheck) ID() string {
	return c.Name
}

// Run the binary check
func (c *executableCheck) Run(ctx context.Context, cfg *common.CLIConfigFlags) (string, int, error) {
	err := c.executableExists(ctx, cfg)
	if err != nil {
		return "", constants.StatusFailure, err
	}
	return "", constants.StatusOK, nil
}

func (c *executableCheck) executableExists(ctx context.Context, cfg *common.CLIConfigFlags) error {
	var args = c.Args

	if len(args) == 0 {
		return fmt.Errorf("No executable to check")
	}

	if len(args) > 1 {
		return fmt.Errorf("Only one executable allowed at a time")
	}

	_, _, exitCode, err := exec.FullOutput(exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("command -v %s", args[0])))
	if err != nil {
		return fmt.Errorf("ERROR: Unable to determine whether %s is available", args[0])
	}
	if exitCode != 0 {
		return fmt.Errorf("%s not available", args[0])
	}

	return nil
}
