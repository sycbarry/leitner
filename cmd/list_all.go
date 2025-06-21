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

		// Count items for summary
		packageCount := 0
		totalDeckCount := 0
		tagCount := 0
		totalFileCount := 0

		// List packages and decks
		fmt.Println("📦 Packages and Decks")
		fmt.Println("═══════════════════════")

		hasPackages := false
		for _, pkg := range packages {
			if pkg.IsDir() && pkg.Name() != "__config__" && pkg.Name() != "__tags__" {
				hasPackages = true
				packageCount++
				fmt.Printf("📁 %s/\n", pkg.Name())
				packagePath := filepath.Join(leitnerPath, pkg.Name())
				entries, err := os.ReadDir(packagePath)
				if err != nil {
					continue
				}

				deckCount := 0
				for _, entry := range entries {
					if entry.IsDir() {
						deckCount++
						totalDeckCount++
						fmt.Printf("   📚 %s/\n", entry.Name())
					}
				}

				if deckCount == 0 {
					fmt.Printf("   └── (no decks)\n")
				}
			}
		}

		if !hasPackages {
			fmt.Println("   (no packages found)")
		}

		// List tags
		fmt.Println("\n🏷️  Tags")
		fmt.Println("═══════════")

		tagsRootPath := filepath.Join(leitnerPath, "__tags__")
		tags, err := os.ReadDir(tagsRootPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("   (no tags found)")
			} else {
				fmt.Println("Error reading tags directory:", err)
				return
			}
		} else {
			if len(tags) == 0 {
				fmt.Println("   (no tags found)")
			} else {
				for _, tag := range tags {
					if tag.IsDir() {
						tagCount++
						fmt.Printf("📂 %s/\n", tag.Name())
						tagPath := filepath.Join(tagsRootPath, tag.Name())
						files, err := os.ReadDir(tagPath)
						if err != nil {
							continue
						}

						fileCount := 0
						for _, file := range files {
							if !file.IsDir() {
								fileCount++
								totalFileCount++
								fmt.Printf("   📄 %s\n", file.Name())
							}
						}

						if fileCount == 0 {
							fmt.Printf("   └── (no files)\n")
						}
					}
				}
			}
		}

		// Summary
		fmt.Println("\n📊 Summary")
		fmt.Println("═══════════")
		fmt.Printf("📦 Packages: %d\n", packageCount)
		fmt.Printf("📚 Total Decks: %d\n", totalDeckCount)
		fmt.Printf("🏷️  Tags: %d\n", tagCount)
		fmt.Printf("📄 Total Files: %d\n", totalFileCount)
	},
}
