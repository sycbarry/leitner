package cmd

import (
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource",
	Long:  `Delete a resource managed by leitner.`,
}
