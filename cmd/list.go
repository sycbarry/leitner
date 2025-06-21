package cmd

import (
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Long:  `List various resources managed by leitner`,
}
