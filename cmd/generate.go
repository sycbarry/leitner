package cmd

import (
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate decks or other resources using an LLM",
	Long:  `Use a configured LLM provider to generate learning resources.`,
}

func init() {
	GenerateCmd.AddCommand(GenerateDeckCmd)
}
