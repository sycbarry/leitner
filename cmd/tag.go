package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var tagName string

var TagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag content from stdin or manage tags",
	Long: `Tag content from stdin or manage tags with subcommands.
	
When run without a subcommand, it reads from stdin.`,
	Run: func(cmd *cobra.Command, args []string) {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Println("This command is intended to be used with piped data.")
			os.Exit(1)
		}

		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			os.Exit(1)
		}

		homeDir, _ := os.UserHomeDir()
		tagPath := filepath.Join(homeDir, ".leitner", "__tags__", tagName)
		if err := os.MkdirAll(tagPath, 0755); err != nil {
			fmt.Println("Error creating tag directory:", err)
			os.Exit(1)
		}

		filename := fmt.Sprintf("%s-%d.html", tagName, time.Now().Unix())
		filePath := filepath.Join(tagPath, filename)
		if err := ioutil.WriteFile(filePath, input, 0644); err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}

		fmt.Printf("Content saved to %s\n", filePath)
	},
}

func init() {
	TagCmd.Flags().StringVar(&tagName, "name", "", "Name of the tag")
	TagCmd.MarkFlagRequired("name")
	TagCmd.AddCommand(TagDeleteCmd)
}
