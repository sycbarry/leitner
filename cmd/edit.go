package cmd

import (
	"github.com/spf13/cobra"
)

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit resources",
	Long:  `Edit various resources managed by leitner`,
}
