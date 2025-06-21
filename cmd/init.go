package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize leitner on the local system",
	Long:  `Initialize leitner on the local system by creating the necessary directories.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}

		leitnerPath := filepath.Join(homeDir, ".leitner")
		configPath := filepath.Join(leitnerPath, "__config__")
		tagsPath := filepath.Join(leitnerPath, "__tags__")

		paths := []string{leitnerPath, configPath, tagsPath}

		for _, path := range paths {
			if err := os.MkdirAll(path, 0755); err != nil {
				fmt.Printf("Error creating directory %s: %v\n", path, err)
				os.Exit(1)
			}
		}
		fmt.Println("Leitner initialized successfully.")
	},
}
