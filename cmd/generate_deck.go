package cmd

import (
	"fmt"
	"io/ioutil"
	"leitner/llm"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	generateDeckPackageName string
	generateDeckName        string
	fromTag                 string
	cardCount               int
)

var GenerateDeckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Generate a new deck from a tag using an LLM",
	Long:  `Generate a new deck of flashcards from the content of a tag using your configured LLM provider.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reading content from tag:", fromTag)
		homeDir, _ := os.UserHomeDir()
		tagPath := filepath.Join(homeDir, ".leitner", "__tags__", fromTag)
		files, err := ioutil.ReadDir(tagPath)
		if err != nil {
			fmt.Printf("Error reading tag directory '%s': %v\n", fromTag, err)
			return
		}

		var contentBuilder strings.Builder
		for _, file := range files {
			data, err := ioutil.ReadFile(filepath.Join(tagPath, file.Name()))
			if err != nil {
				continue
			}
			contentBuilder.Write(data)
			contentBuilder.WriteString("\n\n")
		}

		fmt.Println("Generating flashcards with LLM...")
		deckJSON, err := llm.GenerateFlashcards(generateDeckName, contentBuilder.String(), fromTag, cardCount)
		if err != nil {
			fmt.Println("Error generating flashcards:", err)
			return
		}

		fmt.Println("Creating new deck...")
		deckDir := filepath.Join(homeDir, ".leitner", generateDeckPackageName, generateDeckName)
		if err := os.MkdirAll(deckDir, 0755); err != nil {
			fmt.Printf("Error creating deck directory: %v\n", err)
			return
		}

		deckFilePath := filepath.Join(deckDir, "deck.json")
		if err := ioutil.WriteFile(deckFilePath, deckJSON, 0644); err != nil {
			fmt.Printf("Error writing deck.json: %v\n", err)
			return
		}

		fmt.Printf("Successfully generated deck '%s' in package '%s'.\n", generateDeckName, generateDeckPackageName)
	},
}

func init() {
	GenerateDeckCmd.Flags().StringVar(&generateDeckPackageName, "package", "", "Name of the package for the new deck")
	GenerateDeckCmd.Flags().StringVar(&generateDeckName, "name", "", "Name of the new deck to generate")
	GenerateDeckCmd.Flags().StringVar(&fromTag, "from-tag", "", "Tag to use as context for generation")
	GenerateDeckCmd.Flags().IntVar(&cardCount, "cardcount", 10, "Number of flashcards to generate")

	GenerateDeckCmd.MarkFlagRequired("package")
	GenerateDeckCmd.MarkFlagRequired("name")
	GenerateDeckCmd.MarkFlagRequired("from-tag")
}
