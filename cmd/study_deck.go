package cmd

import (
	"fmt"
	"leitner/web"
	"os"
	"os/exec"
	"path/filepath"
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
		// Generate session ID
		timestamp := time.Now().Unix()
		sessionID := fmt.Sprintf("%d-%s-%s", timestamp, studyDeckPackageName, studyDeckName)

		// Create session file path
		homeDir, _ := os.UserHomeDir()
		sessionsDir := filepath.Join(homeDir, ".leitner", "__sessions__")
		os.MkdirAll(sessionsDir, 0755)
		sessionFile := filepath.Join(sessionsDir, sessionID+".json")
		// Create empty session file (or touch it)
		os.WriteFile(sessionFile, []byte("{}"), 0644)

		// Pass sessionID to the web server
		go web.StartStudyServerWithSession(studyDeckPackageName, studyDeckName, sessionID)
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
