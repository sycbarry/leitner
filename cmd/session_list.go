package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var SessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all study sessions",
	Long:  `List all saved study sessions, ordered by most recent.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir()
		sessionsDir := filepath.Join(homeDir, ".leitner", "__sessions__")
		files, err := os.ReadDir(sessionsDir)
		if err != nil || len(files) == 0 {
			fmt.Println("No session files found.")
			return
		}
		sessionFiles := []string{}
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
				sessionFiles = append(sessionFiles, f.Name())
			}
		}
		if len(sessionFiles) == 0 {
			fmt.Println("No session files found.")
			return
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
		fmt.Printf("%-24s %-20s %-20s\n", "Date", "Package", "Deck")
		fmt.Println(strings.Repeat("-", 64))
		for _, fname := range sessionFiles {
			parts := strings.SplitN(strings.TrimSuffix(fname, ".json"), "-", 3)
			if len(parts) < 3 {
				continue
			}
			ts, _ := strconv.ParseInt(parts[0], 10, 64)
			date := time.Unix(ts, 0).Format("2006-01-02 15:04:05")
			packageName := parts[1]
			deckName := parts[2]
			fmt.Printf("%-24s %-20s %-20s\n", date, packageName, deckName)
		}
	},
}

func init() {
	// Register this command under the 'session' parent in root.go
}
