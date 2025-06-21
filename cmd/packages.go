package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var PackagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "List all packages (folders) in ~/.leitner",
	Long:  `Lists all folders in the user's ~/.leitner directory`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}
		leitnerPath := filepath.Join(homeDir, ".leitner")
		entries, err := os.ReadDir(leitnerPath)
		if err != nil {
			fmt.Println("Error reading .leitner directory:", err)
			os.Exit(1)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				fmt.Println(entry.Name())
			}
		}
	},
}
