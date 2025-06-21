package cmd

import (
	"github.com/spf13/cobra"
)

var StudyCmd = &cobra.Command{
	Use:   "study",
	Short: "Study a deck of flashcards",
	Long:  `Start a study session for a specified deck.`,
}
