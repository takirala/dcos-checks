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
	"github.com/dcos/dcos-checks/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "checks <check name> [parameters]",
	Short: "DC/OS health checks",
	Long: `DC/OS checks provides an easy interface to check the DC/OS components health

The checks could be executed against a signle node, or a whole cluster.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if common.DCOSConfig.Verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// run the commands parser
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Error parsing subcommands: %s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.checks.yaml)")
	rootCmd.PersistentFlags().BoolVar(&common.DCOSConfig.ForceTLS, "force-tls", false, "use HTTPS for GET/POST requests")
	rootCmd.PersistentFlags().BoolVar(&common.DCOSConfig.Verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVar(&common.DCOSConfig.Role, "role", "", "set DC/OS role. (valid roles: master, agent, public-agent)")
	rootCmd.PersistentFlags().StringVar(&common.DCOSConfig.IAMConfig, "iam-config", "", "a path to identity and access managment config")
	rootCmd.PersistentFlags().StringVar(&common.DCOSConfig.CACert, "ca-cert", "", "a path to certificate authority file")
	rootCmd.PersistentFlags().StringVar(&common.DCOSConfig.DetectIP, "detect-ip", "/opt/mesosphere/bin/detect_ip", "a path to detect ip script")
	rootCmd.PersistentFlags().StringVar(&common.DCOSConfig.NodeIPStr, "node-ip", "", "set node IP address overriding detect_ip output")

	// add the subpackage commands
	addSubcommands()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("dcos-checks-config") // name of config file (without extension)
	viper.AddConfigPath("/opt/mesosphere/etc/")
	viper.AutomaticEnv()

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
		common.DCOSConfig.Role = viper.GetString("role")
		common.DCOSConfig.ForceTLS = viper.GetBool("force-tls")
		common.DCOSConfig.Verbose = viper.GetBool("verbose")
		common.DCOSConfig.IAMConfig = viper.GetString("iam-config")
		common.DCOSConfig.CACert = viper.GetString("ca-cert")
		common.DCOSConfig.DetectIP = viper.GetString("detect-ip")
		common.DCOSConfig.NodeIPStr = viper.GetString("node-ip")
	}
}
