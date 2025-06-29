package cmd

import (
	"github.com/spf13/cobra"
)

var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new resources",
	Long:  `Create new resources managed by leitner`,
}

func init() {
	NewCmd.AddCommand(PackageCmd)
	NewCmd.AddCommand(NewDeckCmd)
	NewCmd.AddCommand(NewTagCmd)
}
