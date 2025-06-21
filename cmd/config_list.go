package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Re-defining Config struct for this command's scope.
type ListConfig struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

var ConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the current LLM provider configuration",
	Long:  `Displays the currently configured LLM provider and a masked version of the API key.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir()
		configPath := filepath.Join(homeDir, ".leitner", "__config__", "config.json")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println("No configuration found. Please run './leitner config' to set it up.")
			return
		}

		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}

		var config ListConfig
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Println("Error parsing config file:", err)
			return
		}

		maskedKey := "****************"
		if len(config.APIKey) > 4 {
			maskedKey = fmt.Sprintf("...%s", config.APIKey[len(config.APIKey)-4:])
		}

		fmt.Println("Current LLM Configuration:")
		fmt.Printf("  Provider: %s\n", config.Provider)
		fmt.Printf("  API Key:  %s\n", maskedKey)
	},
}
