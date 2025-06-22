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

var ListAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all packages, decks, and tags in a tree structure",
	Long:  `Lists all packages, their decks, and tags in ~/.leitner in a tree structure`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}
		leitnerPath := filepath.Join(homeDir, ".leitner")
		packages, err := os.ReadDir(leitnerPath)
		if err != nil {
			fmt.Println("Error reading .leitner directory:", err)
			os.Exit(1)
		}

		// Count items for summary
		packageCount := 0
		totalDeckCount := 0
		tagCount := 0
		totalFileCount := 0

		// List packages and decks (but count first for summary)
		for _, pkg := range packages {
			if pkg.IsDir() && pkg.Name() != "__config__" && pkg.Name() != "__tags__" && pkg.Name() != "__sessions__" {
				packageCount++
				packagePath := filepath.Join(leitnerPath, pkg.Name())
				entries, err := os.ReadDir(packagePath)
				if err != nil {
					continue
				}
				for _, entry := range entries {
					if entry.IsDir() {
						totalDeckCount++
					}
				}
			}
		}
		// Count tags and files
		tagsRootPath := filepath.Join(leitnerPath, "__tags__")
		tags, err := os.ReadDir(tagsRootPath)
		if err == nil {
			for _, tag := range tags {
				if tag.IsDir() {
					tagCount++
					tagPath := filepath.Join(tagsRootPath, tag.Name())
					files, err := os.ReadDir(tagPath)
					if err == nil {
						for _, file := range files {
							if !file.IsDir() {
								totalFileCount++
							}
						}
					}
				}
			}
		}

		// Print summary box at the top
		fmt.Println("┌──────────────────────────────────────────────┐")
		fmt.Printf("| 📦 Packages: %-6d | 📚 Decks: %-8d     |\n", packageCount, totalDeckCount)
		fmt.Printf("| 🏷️  Tags:    %-6d | 📄 Files: %-8d     |\n", tagCount, totalFileCount)
		fmt.Println("└──────────────────────────────────────────────┘")

		// List packages and decks
		fmt.Println("\n📦 Packages and Decks")
		fmt.Println("═══════════════════════")

		hasPackages := false
		for _, pkg := range packages {
			if pkg.IsDir() && pkg.Name() != "__config__" && pkg.Name() != "__tags__" && pkg.Name() != "__sessions__" {
				hasPackages = true
				connector := "├"
				fmt.Printf("%s── 📁 %s/\n", connector, pkg.Name())
				packagePath := filepath.Join(leitnerPath, pkg.Name())
				entries, err := os.ReadDir(packagePath)
				if err != nil {
					continue
				}

				decks := []os.DirEntry{}
				for _, entry := range entries {
					if entry.IsDir() {
						decks = append(decks, entry)
					}
				}
				for di, deck := range decks {
					isLastDeck := di == len(decks)-1
					branch := "├"
					if isLastDeck {
						branch = "└"
					}
					fmt.Printf("%s   %s── 📚 %s/\n", connector, branch, deck.Name())
				}
				if len(decks) == 0 {
					fmt.Printf("%s   └── (no decks)\n", connector)
				}
			}
		}

		if !hasPackages {
			fmt.Println("   (no packages found)")
		}

		// List tags
		fmt.Println("\n🏷️  Tags")
		fmt.Println("═══════════")

		if len(tags) == 0 {
			fmt.Println("   (no tags found)")
		} else {
			for ti, tag := range tags {
				if tag.IsDir() {
					tagCount++
					tagConnector := "├"
					if ti == len(tags)-1 {
						tagConnector = "└"
					}
					fmt.Printf("%s── 📂 %s/\n", tagConnector, tag.Name())
					tagPath := filepath.Join(tagsRootPath, tag.Name())
					files, err := os.ReadDir(tagPath)
					if err != nil {
						continue
					}

					fileList := []os.DirEntry{}
					for _, file := range files {
						if !file.IsDir() {
							fileList = append(fileList, file)
						}
					}
					for fi, file := range fileList {
						totalFileCount++
						fileBranch := "├"
						if fi == len(fileList)-1 {
							fileBranch = "└"
						}
						fmt.Printf("%s   %s── 📄 %s\n", tagConnector, fileBranch, file.Name())
					}
					if len(fileList) == 0 {
						fmt.Printf("%s   └── (no files)\n", tagConnector)
					}
				}
			}
		}

		// List sessions
		fmt.Println("\n🕒 Sessions")
		fmt.Println("════════════")
		sessionsDir := filepath.Join(leitnerPath, "__sessions__")
		sessionFiles, err := os.ReadDir(sessionsDir)
		if err != nil || len(sessionFiles) == 0 {
			fmt.Println("   (no sessions found)")
		} else {
			sessions := []string{}
			for _, f := range sessionFiles {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
					sessions = append(sessions, f.Name())
				}
			}
			if len(sessions) == 0 {
				fmt.Println("   (no sessions found)")
			} else {
				sort.Slice(sessions, func(i, j int) bool {
					getTS := func(name string) int64 {
						parts := strings.SplitN(name, "-", 2)
						if len(parts) < 2 {
							return 0
						}
						ts, _ := strconv.ParseInt(parts[0], 10, 64)
						return ts
					}
					return getTS(sessions[i]) > getTS(sessions[j])
				})
				for si, fname := range sessions {
					parts := strings.SplitN(strings.TrimSuffix(fname, ".json"), "-", 3)
					if len(parts) < 3 {
						continue
					}
					ts, _ := strconv.ParseInt(parts[0], 10, 64)
					date := time.Unix(ts, 0).Format("2006-01-02 15:04:05")
					packageName := parts[1]
					deckName := parts[2]
					branch := "├"
					if si == len(sessions)-1 {
						branch = "└"
					}
					fmt.Printf("%s── 🗂️  %s | %s | %s\n", branch, date, packageName, deckName)
				}
			}
		}
	},
}
