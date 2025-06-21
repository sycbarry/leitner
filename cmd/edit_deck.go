package cmd

import (
	"fmt"
	"leitner/web"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	editDeckPackageName string
	editDeckName        string
)

var EditDeckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Edit a deck in a package using a web UI",
	Long:  `Spin up a webserver to edit a deck's questions and answers in a browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		if editDeckPackageName == "" || editDeckName == "" {
			fmt.Println("Both --package and --name are required.")
			os.Exit(1)
		}
		go web.StartDeckEditorServer(editDeckPackageName, editDeckName)
		time.Sleep(500 * time.Millisecond) // Give server a moment to start
		if runtime.GOOS == "darwin" {
			err := exec.Command("open", "-a", "Google Chrome", "http://localhost:8080/edit").Start()
			if err != nil {
				fmt.Println("Could not open Chrome automatically. Please open http://localhost:8080/edit manually.")
			}
		} else {
			fmt.Println("Please open http://localhost:8080/edit in your browser.")
		}
		select {} // Keep process alive
	},
}

func init() {
	EditDeckCmd.Flags().StringVar(&editDeckPackageName, "package", "", "Name of the package containing the deck")
	EditDeckCmd.Flags().StringVar(&editDeckName, "name", "", "Name of the deck to edit")
	EditDeckCmd.MarkFlagRequired("package")
	EditDeckCmd.MarkFlagRequired("name")
}
