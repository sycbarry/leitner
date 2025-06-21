package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ListTagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags and their content",
	Long:  `Lists all tags and the files within them from ~/.leitner/__tags__ in a tree structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}
		tagsRootPath := filepath.Join(homeDir, ".leitner", "__tags__")

		tags, err := os.ReadDir(tagsRootPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No tags found. The '__tags__' directory doesn't exist.")
				return
			}
			fmt.Println("Error reading tags directory:", err)
			os.Exit(1)
		}

		if len(tags) == 0 {
			fmt.Println("No tags found.")
			return
		}

		for _, tag := range tags {
			if tag.IsDir() {
				fmt.Printf("%s/\n", tag.Name())
				tagPath := filepath.Join(tagsRootPath, tag.Name())
				files, err := os.ReadDir(tagPath)
				if err != nil {
					continue // Skip directories that can't be read
				}
				for _, file := range files {
					if !file.IsDir() {
						fmt.Printf("  %s\n", file.Name())
					}
				}
			}
		}
	},
}
