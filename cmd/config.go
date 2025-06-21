package cmd

import (
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage LLM provider configuration",
	Long:  `Manage LLM provider configuration. Use subcommands to set or list configurations.`,
}

func init() {
	ConfigCmd.AddCommand(ConfigSetCmd)
	ConfigCmd.AddCommand(ConfigListCmd)
}
