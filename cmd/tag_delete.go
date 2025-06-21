package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	deleteTagName     string
	deleteTagFilename string
)

var TagDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific file from a tag",
	Long:  `Deletes a single content file from a specified tag directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir()
		filePath := filepath.Join(homeDir, ".leitner", "__tags__", deleteTagName, deleteTagFilename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' not found in tag '%s'.\n", deleteTagFilename, deleteTagName)
			return
		}

		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Error deleting file: %v\n", err)
			return
		}

		fmt.Printf("Successfully deleted '%s' from tag '%s'.\n", deleteTagFilename, deleteTagName)
	},
}

func init() {
	TagDeleteCmd.Flags().StringVar(&deleteTagName, "name", "", "Name of the tag")
	TagDeleteCmd.Flags().StringVar(&deleteTagFilename, "file", "", "Filename to delete from the tag")
	TagDeleteCmd.MarkFlagRequired("name")
	TagDeleteCmd.MarkFlagRequired("file")
}
