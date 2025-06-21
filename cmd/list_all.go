package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ListAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all packages and decks in a tree structure",
	Long:  `Lists all packages and their decks in ~/.leitner in a tree structure`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}
		leitnerPath := filepath.Join(homeDir, ".leitner")
		packages, err := os.ReadDir(leitnerPath)
		if err != nil {
			fmt.Println("Error reading .leitner directory:", err)
			os.Exit(1)
		}
		for _, pkg := range packages {
			if pkg.IsDir() && pkg.Name() != "__config__" && pkg.Name() != "__tags__" {
				fmt.Printf("%s/\n", pkg.Name())
				packagePath := filepath.Join(leitnerPath, pkg.Name())
				entries, err := os.ReadDir(packagePath)
				if err != nil {
					continue
				}
				for _, entry := range entries {
					if entry.IsDir() {
						fmt.Printf("  %s/\n", entry.Name())
					}
				}
			}
		}
	},
}
