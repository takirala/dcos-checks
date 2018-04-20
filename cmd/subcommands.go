package cmd

import (
	"github.com/dcos/dcos-checks/cmd/checks/components"
	"github.com/dcos/dcos-checks/cmd/checks/executable"
	"github.com/dcos/dcos-checks/cmd/checks/ip"
	"github.com/dcos/dcos-checks/cmd/checks/journald"
	"github.com/dcos/dcos-checks/cmd/checks/mesosmetrics"
	"github.com/dcos/dcos-checks/cmd/checks/time"
	"github.com/dcos/dcos-checks/cmd/checks/version"
	"github.com/spf13/cobra"
)

// RegisterFunc represents a function that adds a subcommand to the rootCmd.
type RegisterFunc func(rootCmd *cobra.Command)

// RegisterSubcommand calls the given RegisterFunc, passing in the root command.
func RegisterSubcommand(fn RegisterFunc) {
	fn(rootCmd)
}

func addSubcommands() {
	RegisterSubcommand(components.Register)
	RegisterSubcommand(executable.Register)
	RegisterSubcommand(ip.Register)
	RegisterSubcommand(journald.Register)
	RegisterSubcommand(mesosmetrics.Register)
	RegisterSubcommand(time.Register)
	RegisterSubcommand(version.Register)
}
