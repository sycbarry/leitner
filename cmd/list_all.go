package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ListAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all packages, decks, and tags in a tree structure",
	Long:  `Lists all packages, their decks, and tags in ~/.leitner in a tree structure`,
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

		// List packages and decks
		fmt.Println("📦 Packages and Decks:")
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

		// List tags
		fmt.Println("\n🏷️  Tags:")
		tagsRootPath := filepath.Join(leitnerPath, "__tags__")
		tags, err := os.ReadDir(tagsRootPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("  No tags found.")
				return
			}
			fmt.Println("Error reading tags directory:", err)
			return
		}

		if len(tags) == 0 {
			fmt.Println("  No tags found.")
			return
		}

		for _, tag := range tags {
			if tag.IsDir() {
				fmt.Printf("  %s/\n", tag.Name())
				tagPath := filepath.Join(tagsRootPath, tag.Name())
				files, err := os.ReadDir(tagPath)
				if err != nil {
					continue
				}
				for _, file := range files {
					if !file.IsDir() {
						fmt.Printf("    %s\n", file.Name())
					}
				}
			}
		}
	},
}
