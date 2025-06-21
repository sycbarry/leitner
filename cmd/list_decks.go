package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var listDecksPackageName string

var ListDecksCmd = &cobra.Command{
	Use:   "decks",
	Short: "List all decks in a package",
	Long:  `Lists all decks (subfolders) in the specified package in ~/.leitner`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}
		packagePath := filepath.Join(homeDir, ".leitner", listDecksPackageName)
		entries, err := os.ReadDir(packagePath)
		if err != nil {
			fmt.Printf("Error reading package '%s': %v\n", listDecksPackageName, err)
			os.Exit(1)
		}
		found := false
		for _, entry := range entries {
			if entry.IsDir() {
				fmt.Println(entry.Name())
				found = true
			}
		}
		if !found {
			fmt.Println("No decks found in this package.")
		}
	},
}

func init() {
	ListDecksCmd.Flags().StringVar(&listDecksPackageName, "package", "", "Name of the package to list decks from")
	ListDecksCmd.MarkFlagRequired("package")
}
