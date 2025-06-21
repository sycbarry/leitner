package cmd

import (
	"leitner/web"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var (
	studyDeckPackageName string
	studyDeckName        string
)

var StudyDeckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Study a deck using the web UI",
	Long:  `Spin up a webserver to study a deck's flashcards in a browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		go web.StartStudyServer(studyDeckPackageName, studyDeckName)
		time.Sleep(500 * time.Millisecond) // Give server a moment to start
		if runtime.GOOS == "darwin" {
			exec.Command("open", "-a", "Google Chrome", "http://localhost:8080/study").Start()
		}
		select {} // Keep process alive
	},
}

func init() {
	StudyDeckCmd.Flags().StringVar(&studyDeckPackageName, "package", "", "Name of the package containing the deck")
	StudyDeckCmd.Flags().StringVar(&studyDeckName, "name", "", "Name of the deck to study")
	StudyDeckCmd.MarkFlagRequired("package")
	StudyDeckCmd.MarkFlagRequired("name")
}
