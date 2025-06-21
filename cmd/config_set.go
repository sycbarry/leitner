package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Config struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

var ConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set or update an LLM provider",
	Long:  `Interactively set or update an LLM provider by selecting from a list and providing an API key.`,
	Run: func(cmd *cobra.Command, args []string) {
		providerPrompt := promptui.Select{
			Label: "Select LLM Provider",
			Items: []string{"openai", "claude"},
		}
		_, provider, err := providerPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		keyPrompt := promptui.Prompt{
			Label: "Enter API Key",
			Mask:  '*',
		}
		apiKey, err := keyPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		config := Config{
			Provider: provider,
			APIKey:   apiKey,
		}

		homeDir, _ := os.UserHomeDir()
		configDir := filepath.Join(homeDir, ".leitner", "__config__")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Println("Error creating config directory:", err)
			os.Exit(1)
		}

		configPath := filepath.Join(configDir, "config.json")
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling config to JSON:", err)
			return
		}

		if err := ioutil.WriteFile(configPath, data, 0600); err != nil {
			fmt.Println("Error writing config file:", err)
			return
		}

		fmt.Printf("Configuration saved for provider: %s\n", provider)
	},
}
