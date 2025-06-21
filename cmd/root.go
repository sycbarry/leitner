package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasSubCommands}}
  {{.CommandPath}} [command]{{end}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Global Flags:
      --help      Show help for command

Use "{{.CommandPath}} [command] --help" for more information about a command.

Examples:

  # Initialize leitner
  leitner init

  # Create a new package
  leitner new package --name=<packagename>

  # Create a new deck inside a package
  leitner new deck --package=<packagename> --name=<deckname>

  # List all packages
  leitner list packages

  # List all decks in a package
  leitner list decks --package=<packagename>

  # List all packages and decks in a tree
  leitner list all

  # List all tags and their content
  leitner list tags

  # Capture and tag content from stdin
  # curl -s https://example.com | leitner tag --name=<tagname>

  # Edit a deck in a package (opens web editor)
  leitner edit deck --package=<packagename> --name=<deckname>

  # Study a deck
  leitner study deck --package=<packagename> --name=<deckname>

  # Delete a package
  leitner delete package --name=<packagename>

  # Delete a deck from a package
  leitner delete deck --package=<packagename> --name=<deckname>
`

var rootCmd = &cobra.Command{
	Use:   "leitner",
	Short: "testing.",
	Long:  `testing this.`,
}

func Execute() {
	rootCmd.SetUsageTemplate(usageTemplate)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(NewCmd)
	NewCmd.AddCommand(PackageCmd)
	NewCmd.AddCommand(NewDeckCmd)
	rootCmd.AddCommand(ListCmd)
	ListCmd.AddCommand(PackagesCmd)
	ListCmd.AddCommand(ListDecksCmd)
	ListCmd.AddCommand(ListAllCmd)
	ListCmd.AddCommand(ListTagsCmd)
	rootCmd.AddCommand(DeleteCmd)
	DeleteCmd.AddCommand(DeletePackageCmd)
	DeleteCmd.AddCommand(DeleteDeckCmd)
	rootCmd.AddCommand(EditCmd)
	EditCmd.AddCommand(EditDeckCmd)
	rootCmd.AddCommand(StudyCmd)
	StudyCmd.AddCommand(StudyDeckCmd)
	rootCmd.AddCommand(TagCmd)
}
