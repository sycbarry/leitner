package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var NewTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Create a new tag",
	Long:  `Create a new empty tag directory in ~/.leitner/__tags__`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}

		tagPath := filepath.Join(homeDir, ".leitner", "__tags__", tagName)

		// Check if tag already exists
		if _, err := os.Stat(tagPath); err == nil {
			fmt.Printf("Tag '%s' already exists.\n", tagName)
			os.Exit(1)
		}

		// Create the tag directory
		if err := os.MkdirAll(tagPath, 0755); err != nil {
			fmt.Println("Error creating tag directory:", err)
			os.Exit(1)
		}

		fmt.Printf("Tag '%s' created successfully.\n", tagName)
	},
}

func init() {
	NewTagCmd.Flags().StringVar(&tagName, "name", "", "Name of the tag")
	NewTagCmd.MarkFlagRequired("name")
}
