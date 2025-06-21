package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	tagName  string
	fromFile string
)

var TagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag content from stdin, file, or manage tags",
	Long: `Tag content from stdin, file, or manage tags with subcommands.
	
When run without a subcommand, it reads from stdin or from a specified file.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir()
		tagPath := filepath.Join(homeDir, ".leitner", "__tags__", tagName)
		if err := os.MkdirAll(tagPath, 0755); err != nil {
			fmt.Println("Error creating tag directory:", err)
			os.Exit(1)
		}

		var content []byte
		var err error
		var filename string

		if fromFile != "" {
			// Read from specified file
			content, err = ioutil.ReadFile(fromFile)
			if err != nil {
				fmt.Printf("Error reading file '%s': %v\n", fromFile, err)
				os.Exit(1)
			}

			// Use original filename with timestamp to avoid conflicts
			baseName := filepath.Base(fromFile)
			ext := filepath.Ext(baseName)
			nameWithoutExt := baseName[:len(baseName)-len(ext)]
			filename = fmt.Sprintf("%s-%d%s", nameWithoutExt, time.Now().Unix(), ext)
		} else {
			// Read from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				fmt.Println("This command is intended to be used with piped data or --from-file flag.")
				os.Exit(1)
			}

			content, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("Error reading from stdin:", err)
				os.Exit(1)
			}

			filename = fmt.Sprintf("%s-%d.html", tagName, time.Now().Unix())
		}

		filePath := filepath.Join(tagPath, filename)
		if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}

		fmt.Printf("Content saved to %s\n", filePath)
	},
}

func init() {
	TagCmd.Flags().StringVar(&tagName, "name", "", "Name of the tag")
	TagCmd.Flags().StringVar(&fromFile, "from-file", "", "File to tag (optional)")
	TagCmd.MarkFlagRequired("name")
	TagCmd.AddCommand(TagDeleteCmd)
}
