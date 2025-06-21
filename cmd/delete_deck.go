package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	deleteDeckPackageName string
	deleteDeckName        string
)

var DeleteDeckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Delete a deck from a package",
	Long:  `Delete a deck and its contents from a specified study package.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}

		packagePath := filepath.Join(homeDir, ".leitner", deleteDeckPackageName)
		if _, err := os.Stat(packagePath); os.IsNotExist(err) {
			fmt.Printf("Error: Package '%s' not found.\n", deleteDeckPackageName)
			os.Exit(1)
		}

		deckPath := filepath.Join(packagePath, deleteDeckName)
		if _, err := os.Stat(deckPath); os.IsNotExist(err) {
			fmt.Printf("Error: Deck '%s' not found in package '%s'.\n", deleteDeckName, deleteDeckPackageName)
			os.Exit(1)
		}

		if err := os.RemoveAll(deckPath); err != nil {
			fmt.Printf("Error deleting deck: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deck '%s' deleted successfully from package '%s'.\n", deleteDeckName, deleteDeckPackageName)
	},
}

func init() {
	DeleteDeckCmd.Flags().StringVar(&deleteDeckPackageName, "package", "", "Name of the package containing the deck")
	DeleteDeckCmd.Flags().StringVar(&deleteDeckName, "name", "", "Name of the deck to delete")
	DeleteDeckCmd.MarkFlagRequired("package")
	DeleteDeckCmd.MarkFlagRequired("name")
}
