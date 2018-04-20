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

var validExecutables = map[string]bool{
	"curl":  true,
	"wget":  true,
	"tar":   true,
	"git":   true,
	"xz":    true,
	"unzip": true,
}

// executableCmd represents the executable command
var executableCmd = &cobra.Command{
	Use:   "executable",
	Short: "Check for executable/executables required to install DC/OS",
	Long: `Check for existence of the following executable:
curl
wget
tar
git
xz
unzip
`,
	Run: func(cmd *cobra.Command, args []string) {
		common.RunCheck(context.TODO(),
			newExecutableCheck("DC/OS verify existence of executables", args))
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

	if !validExecutables[args[0]] {
		var keys []string
		for key := range validExecutables {
			keys = append(keys, key)
		}
		return fmt.Errorf("Choose from valid list of executables %v", keys)
	}

	if _, _, _, err := exec.FullOutput(exec.CommandContext(ctx, args[0])); err != nil {
		return fmt.Errorf("%s not installed", args[0])
	}
	return nil
}
