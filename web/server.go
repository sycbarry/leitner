package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Card struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type Deck struct {
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}

func StartDeckEditorServer(packageName, deckName string) {
	homeDir, _ := os.UserHomeDir()
	deckPath := filepath.Join(homeDir, ".leitner", packageName, deckName, "deck.json")
	webDir := "web" // Assuming web files are in a 'web' directory

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "style.css"))
	})
	http.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "app.js"))
	})

	http.HandleFunc("/deck", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			data, err := ioutil.ReadFile(deckPath)
			if err != nil {
				http.Error(w, "Deck not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}
		if r.Method == http.MethodPost {
			var deck Deck
			if err := json.NewDecoder(r.Body).Decode(&deck); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			data, _ := json.MarshalIndent(deck, "", "  ")
			ioutil.WriteFile(deckPath, data, 0644)
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	fmt.Println("Deck editor running at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func StartStudyServer(packageName, deckName string) {
	homeDir, _ := os.UserHomeDir()
	deckPath := filepath.Join(homeDir, ".leitner", packageName, deckName, "deck.json")
	webDir := "web"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "study.html"))
	})
	http.HandleFunc("/study.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "study.css"))
	})
	http.HandleFunc("/study.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(webDir, "study.js"))
	})

	http.HandleFunc("/deck", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(deckPath)
		if err != nil {
			http.Error(w, "Deck not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	fmt.Println("Study session running at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
