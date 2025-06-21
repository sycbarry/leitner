package cmd

import (
	"leitner/studypack"

	"github.com/spf13/cobra"
)

var packageName string

var PackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Create a new study package",
	Long:  `Create a new study package`,
	Run: func(cmd *cobra.Command, args []string) {
		pack := studypack.StudyPackage{PackageName: packageName}
		pack.InitializePackage()
	},
}

func init() {
	PackageCmd.Flags().StringVarP(&packageName, "name", "n", "", "Package name")
	PackageCmd.MarkFlagRequired("name")
}
