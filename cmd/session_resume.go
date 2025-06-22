package cmd

import (
	"fmt"
	"leitner/web"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var sessionResumeID string

var SessionResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the latest or a specific study session",
	Long:  `Resume the most recent study session, a specific one with --id, or interactively select from a list.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir()
		sessionsDir := filepath.Join(homeDir, ".leitner", "__sessions__")
		var sessionFile string
		if sessionResumeID != "" {
			sessionFile = sessionResumeID + ".json"
			if _, err := os.Stat(filepath.Join(sessionsDir, sessionFile)); err != nil {
				fmt.Println("Session file not found:", sessionFile)
				os.Exit(1)
			}
		} else {
			files, err := os.ReadDir(sessionsDir)
			if err != nil || len(files) == 0 {
				fmt.Println("No session files found.")
				os.Exit(1)
			}
			sessionFiles := []string{}
			for _, f := range files {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
					sessionFiles = append(sessionFiles, f.Name())
				}
			}
			if len(sessionFiles) == 0 {
				fmt.Println("No session files found.")
				os.Exit(1)
			}
			sort.Slice(sessionFiles, func(i, j int) bool {
				getTS := func(name string) int64 {
					parts := strings.SplitN(name, "-", 2)
					if len(parts) < 2 {
						return 0
					}
					ts, _ := strconv.ParseInt(parts[0], 10, 64)
					return ts
				}
				return getTS(sessionFiles[i]) > getTS(sessionFiles[j])
			})
			// Build display list for promptui
			displayList := []string{}
			for _, fname := range sessionFiles {
				parts := strings.SplitN(strings.TrimSuffix(fname, ".json"), "-", 3)
				if len(parts) < 3 {
					displayList = append(displayList, fname)
					continue
				}
				ts, _ := strconv.ParseInt(parts[0], 10, 64)
				date := time.Unix(ts, 0).Format("2006-01-02 15:04:05")
				packageName := parts[1]
				deckName := parts[2]
				displayList = append(displayList, fmt.Sprintf("%s | %s | %s", date, packageName, deckName))
			}
			prompt := promptui.Select{
				Label: "Select a session to resume",
				Items: displayList,
				Size:  10,
			}
			idx, _, err := prompt.Run()
			if err != nil {
				fmt.Println("Prompt failed, falling back to default (latest session):", err)
				idx = 0
			}
			sessionFile = sessionFiles[idx]
		}
		// Parse sessionID, package, deck from filename
		parts := strings.SplitN(strings.TrimSuffix(sessionFile, ".json"), "-", 3)
		if len(parts) < 3 {
			fmt.Println("Malformed session file name.")
			os.Exit(1)
		}
		sessionID := parts[0] + "-" + parts[1] + "-" + parts[2]
		packageName := parts[1]
		deckName := parts[2]
		go web.StartStudyServerWithSession(packageName, deckName, sessionID)
		time.Sleep(500 * time.Millisecond)
		if runtime.GOOS == "darwin" {
			exec.Command("open", "-a", "Google Chrome", "http://localhost:8080/study").Start()
		}
		select {}
	},
}

func init() {
	SessionResumeCmd.Flags().StringVar(&sessionResumeID, "id", "", "Session file name (without .json) to resume")
	// Register this command under a new parent 'session' command in root.go
}
