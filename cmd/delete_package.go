package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var deletePackageName string

var DeletePackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Delete a study package",
	Long:  `Delete a study package by name.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}

		packagePath := filepath.Join(homeDir, ".leitner", deletePackageName)

		if _, err := os.Stat(packagePath); os.IsNotExist(err) {
			fmt.Printf("Package '%s' not found.\n", deletePackageName)
			os.Exit(1)
		}

		if err := os.RemoveAll(packagePath); err != nil {
			fmt.Printf("Error deleting package '%s': %v\n", deletePackageName, err)
			os.Exit(1)
		}

		fmt.Printf("Package '%s' deleted successfully.\n", deletePackageName)
	},
}

func init() {
	DeletePackageCmd.Flags().StringVarP(&deletePackageName, "name", "n", "", "Name of the package to delete")
	DeletePackageCmd.MarkFlagRequired("name")
}
