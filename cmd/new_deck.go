package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	deckPackageName string
	deckName        string
)

type Deck struct {
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}

type Card struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

var NewDeckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Create a new deck within a package",
	Long:  `Create a new deck as a subfolder within a specified study package.`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting user home directory:", err)
			os.Exit(1)
		}

		packagePath := filepath.Join(homeDir, ".leitner", deckPackageName)
		if _, err := os.Stat(packagePath); os.IsNotExist(err) {
			fmt.Printf("Error: Package '%s' not found. Please create it first.\n", deckPackageName)
			os.Exit(1)
		}

		deckPath := filepath.Join(packagePath, deckName)
		if err := os.Mkdir(deckPath, 0755); err != nil {
			if os.IsExist(err) {
				fmt.Printf("Error: Deck '%s' already exists in package '%s'.\n", deckName, deckPackageName)
			} else {
				fmt.Printf("Error creating deck directory: %v\n", err)
			}
			os.Exit(1)
		}

		deck := Deck{
			Name:  deckName,
			Cards: []Card{},
		}
		jsonContent, err := json.MarshalIndent(deck, "", "  ")
		if err != nil {
			fmt.Printf("Error creating deck.json: %v\n", err)
			os.Exit(1)
		}

		jsonFilePath := filepath.Join(deckPath, "deck.json")
		if err := os.WriteFile(jsonFilePath, jsonContent, 0644); err != nil {
			fmt.Printf("Error writing deck.json: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deck '%s' and deck.json created successfully in package '%s'.\n", deckName, deckPackageName)
	},
}

func init() {
	NewDeckCmd.Flags().StringVar(&deckPackageName, "package", "", "Name of the package to add the deck to")
	NewDeckCmd.Flags().StringVar(&deckName, "name", "", "Name of the new deck")
	NewDeckCmd.MarkFlagRequired("package")
	NewDeckCmd.MarkFlagRequired("name")
}
