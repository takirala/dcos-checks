// +build linux
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

package time

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/dcos/dcos-checks/common"
	"github.com/dcos/dcos-checks/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	// STA_UNSYNC is taken from https://github.com/torvalds/linux/blob/master/include/uapi/linux/timex.h#L137
	staUnsync = 0x0040

	// 100 millisecond
	maxEstErrorUs = int64(time.Microsecond * 100000)
)

// timeCheck is a time check structure.
type timeCheck struct {
	Name string

	runAdjtimex func(*syscall.Timex) (int, error)
}

// timeCmd represents the time command
var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Verify time is synced",
	Long:  `This check uses a system call adjtimex to validate time is synced.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.RunCheck(context.TODO(), newTimeCheck("Check clock synchronization"))
	},
}

// Add adds this command to the root command
func Add(root *cobra.Command) {
	root.AddCommand(timeCmd)
}

// newTimeCheck returns a new initialized instance of timeCheck.
func newTimeCheck(name string) common.DCOSChecker {
	return &timeCheck{
		Name:        name,
		runAdjtimex: syscall.Adjtimex,
	}
}

// ID returns a check ID.
func (t *timeCheck) ID() string {
	return t.Name
}

// Run executes the check.
func (t *timeCheck) Run(ctx context.Context, cfg *common.CLIConfigFlags) (string, int, error) {
	tBuf := syscall.Timex{}

	// intentionally ignore status. If err != nil, status != 0
	_, err := t.runAdjtimex(&tBuf)
	if err != nil {
		return "", constants.StatusUnknown, errors.Wrap(err, "unable to make a system call adjtimex")
	}

	// This is to check if NTP thinks the clock is unstable
	if diff := tBuf.Esterror - maxEstErrorUs; diff > 0 {
		return fmt.Sprintf("Clock is less stable than allowed. Max estimated error exceeded by: %s", time.Duration(diff)*time.Microsecond), constants.StatusFailure, nil
	}

	// If NTP is down for ~16000 seconds, the clock will go unsync, based on
	// modern kernels. Unfortunately, even though there are a bunch of other
	// heuristics in the timex struct, it doesn't make a ton of sense to look
	// at them. Maybe in the future we can do something smarter.
	if (tBuf.Status & staUnsync) > 0 {
		return "Clock is out of sync / in unsync state. Must be synchronized for proper operation.", constants.StatusFailure, nil
	}

	return "Clock is synced", constants.StatusOK, nil
}
