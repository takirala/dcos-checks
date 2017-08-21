package cmd

import (
	"github.com/dcos/dcos-checks/cmd/checks/components"
	"github.com/dcos/dcos-checks/cmd/checks/executable"
	"github.com/dcos/dcos-checks/cmd/checks/ip"
	"github.com/dcos/dcos-checks/cmd/checks/journald"
	"github.com/dcos/dcos-checks/cmd/checks/mesosmetrics"
	"github.com/dcos/dcos-checks/cmd/checks/time"
	"github.com/dcos/dcos-checks/cmd/checks/version"
)

// this is the place to add each subcommand
func addSubcommands() {
	components.Add(rootCmd)
	executable.Add(rootCmd)
	ip.Add(rootCmd)
	journald.Add(rootCmd)
	mesosmetrics.Add(rootCmd)
	time.Add(rootCmd)
	version.Add(rootCmd)
}
