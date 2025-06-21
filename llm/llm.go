package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	anthropic "github.com/liushuangls/go-anthropic"
	"github.com/sashabaranov/go-openai"
)

type Card struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type Deck struct {
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}

type Config struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

func loadConfig() (Config, error) {
	var config Config
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".leitner", "__config__", "config.json")

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, errors.New("configuration not found, please run 'leitner config set'")
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

func chunkText(text string, chunkSize int) []string {
	var chunks []string
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}

func GenerateFlashcards(deckName, context, tag string, cardCount int) ([]byte, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	chunks := chunkText(context, 100000) // ~25k tokens, well within model limits
	if len(chunks) == 0 {
		return nil, errors.New("no content found in the tag to generate from")
	}

	cardsPerChunk := (cardCount + len(chunks) - 1) / len(chunks) // Ceiling division
	allCards := make([]Card, 0, cardCount)
	fmt.Printf("Content divided into %d chunks. Requesting ~%d cards per chunk.\n", len(chunks), cardsPerChunk)

	for i, chunk := range chunks {
		fmt.Printf("Processing chunk %d of %d...\n", i+1, len(chunks))
		prompt := fmt.Sprintf(`Based on the following text about "%s", generate %d flashcards. Each flashcard should be a distinct question-answer pair.

Provide the output as a valid JSON array of objects, like this: [{"front": "question", "back": "answer"}]. Do not include any text or explanation outside of the JSON array.

Context:
---
%s
---
`, tag, cardsPerChunk, chunk)

		var responseJSON string
		switch config.Provider {
		case "openai":
			responseJSON, err = generateWithOpenAI(config.APIKey, prompt)
		case "claude":
			responseJSON, err = generateWithClaude(config.APIKey, prompt)
		default:
			return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
		}

		if err != nil {
			fmt.Printf("Warning: Could not process chunk %d: %v\n", i+1, err)
			continue // Skip to the next chunk on error
		}

		// Find the start and end of the JSON array
		start := strings.Index(responseJSON, "[")
		end := strings.LastIndex(responseJSON, "]")

		if start == -1 || end == -1 || start > end {
			fmt.Printf("Warning: Could not find a valid JSON array in the response for chunk %d, skipping.\n", i+1)
			continue
		}

		// Extract the JSON array string
		jsonArrayStr := responseJSON[start : end+1]

		var newCards []Card
		if err := json.Unmarshal([]byte(jsonArrayStr), &newCards); err != nil {
			fmt.Printf("Warning: Failed to parse JSON for chunk %d, skipping.\n", i+1)
			continue
		}
		allCards = append(allCards, newCards...)
	}

	finalDeck := Deck{
		Name:  deckName,
		Cards: allCards,
	}

	if len(finalDeck.Cards) == 0 {
		return nil, errors.New("failed to generate any cards from the provided content")
	}

	fmt.Printf("Generated a total of %d cards.\n", len(finalDeck.Cards))
	return json.MarshalIndent(finalDeck, "", "  ")
}

func generateWithOpenAI(apiKey, prompt string) (string, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func generateWithClaude(apiKey, prompt string) (string, error) {
	client := anthropic.NewClient(apiKey)
	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Sonnet20240229,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(prompt),
		},
		MaxTokens: 2048,
	})

	if err != nil {
		return "", err
	}

	if len(resp.Content) == 0 || resp.Content[0].Text == "" {
		return "", errors.New("claude returned an empty or invalid response")
	}

	return resp.Content[0].Text, nil
}
