package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed *.html *.css *.js
var webFiles embed.FS

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

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/edit", http.StatusFound)
	})

	mux.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		// Serve index.html for the edit endpoint
		data, err := webFiles.ReadFile("index.html")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})
	mux.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		data, err := webFiles.ReadFile("style.css")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/css")
		w.Write(data)
	})
	mux.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		data, err := webFiles.ReadFile("app.js")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(data)
	})

	mux.HandleFunc("/deck", func(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Deck editor running at http://localhost:8080/edit")
	http.ListenAndServe(":8080", mux)
}

func StartStudyServer(packageName, deckName string) {
	homeDir, _ := os.UserHomeDir()
	deckPath := filepath.Join(homeDir, ".leitner", packageName, deckName, "deck.json")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/study", http.StatusFound)
	})

	mux.HandleFunc("/study", func(w http.ResponseWriter, r *http.Request) {
		// Serve study.html for the study endpoint
		data, err := webFiles.ReadFile("study.html")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})
	mux.HandleFunc("/study.css", func(w http.ResponseWriter, r *http.Request) {
		data, err := webFiles.ReadFile("study.css")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/css")
		w.Write(data)
	})
	mux.HandleFunc("/study.js", func(w http.ResponseWriter, r *http.Request) {
		data, err := webFiles.ReadFile("study.js")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(data)
	})

	mux.HandleFunc("/deck", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(deckPath)
		if err != nil {
			http.Error(w, "Deck not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	fmt.Println("Study session running at http://localhost:8080/study")
	http.ListenAndServe(":8080", mux)
}

func StartStudyServerWithSession(packageName, deckName, sessionID string) {
	homeDir, _ := os.UserHomeDir()
	deckPath := filepath.Join(homeDir, ".leitner", packageName, deckName, "deck.json")
	sessionsDir := filepath.Join(homeDir, ".leitner", "__sessions__")
	sessionFile := filepath.Join(sessionsDir, sessionID+".json")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "leitner_session_id",
			Value:  sessionID,
			Path:   "/",
			MaxAge: 60 * 60 * 24 * 7, // 1 week
		})
		http.Redirect(w, r, "/study", http.StatusFound)
	})

	mux.HandleFunc("/study", func(w http.ResponseWriter, r *http.Request) {
		data, err := webFiles.ReadFile("study.html")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})
	mux.HandleFunc("/study.css", func(w http.ResponseWriter, r *http.Request) {
		data, err := webFiles.ReadFile("study.css")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/css")
		w.Write(data)
	})
	mux.HandleFunc("/study.js", func(w http.ResponseWriter, r *http.Request) {
		data, err := webFiles.ReadFile("study.js")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(data)
	})

	mux.HandleFunc("/deck", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(deckPath)
		if err != nil {
			http.Error(w, "Deck not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	// Session state endpoints
	mux.HandleFunc("/api/session/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		ioutil.WriteFile(sessionFile, body, 0644)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/api/session/load", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(sessionFile)
		if err != nil {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	fmt.Println("Study session with session tracking running at http://localhost:8080/study")
	http.ListenAndServe(":8080", mux)
}
