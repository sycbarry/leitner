# leitner

A simple, command-line based Leitner system for spaced repetition learning.

<img src="leitner.png" alt="Leitner logo" width="300"/>

A video walkthrough is available here: [Watch the video](https://sebbarry-personal.nyc3.digitaloceanspaces.com/videos/leitner-1.mp4)

## Installation

```bash
# Clone the repository
git clone https://github.com/your-username/leitner.git
cd leitner

# Build the binary
go build
```

## Commands

### `init`

Initializes the Leitner system by creating the necessary directories in your home folder (`~/.leitner`).

**Usage:**
```bash
./leitner init
```

---

### `new`

Create new packages, decks, or tags.

**Create a package:**
```bash
./leitner new package --name=<packagename>
```

**Create a deck:**
```bash
./leitner new deck --package=<packagename> --name=<deckname>
```

**Create a new tag:**
```bash
./leitner new tag --name=<tagname>
```
This creates an empty tag directory that you can later populate with content.

---

### `list`

List packages or decks.

**List all packages:**
```bash
./leitner list packages
```

**List decks in a package:**
```bash
./leitner list decks --package=<packagename>
```

**List all packages and decks in a tree:**
```bash
./leitner list all
```

**List all tags and their captured content:**
```bash
./leitner list tags
```

---

### `edit`

Edit a deck's contents using a web-based editor.

**Usage:**
```bash
./leitner edit deck --package=<packagename> --name=<deckname>
```
This will open a web browser where you can add, edit, and delete flashcards.

---

### `delete`

Delete packages or decks.

**Delete a package:**
```bash
./leitner delete package --name=<packagename>
```

**Delete a deck:**
```bash
./leitner delete deck --package=<packagename> --name=<deckname>
```

---

### `study`

Start a study session for a deck.

**Usage:**
```bash
./leitner study deck --package=<packagename> --name=<deckname>
```
This will open a web browser with a flashcard interface to begin your study session.

---

### `session`

Manage and resume study sessions.

**List all sessions:**
```bash
./leitner session list
```
Lists all saved study sessions, ordered by most recent, showing the date, package, and deck.

**Resume a session:**
```bash
./leitner session resume
```
Interactively select a session to resume (arrow keys, enter to select). Defaults to the latest session if you just press enter.

**Resume a session by ID:**
```bash
./leitner session resume --id=<session-file-name-without-.json>
```
Resumes the session with the specified file name (as shown in `session list`).

---

### `tag`

Capture and tag content or manage tagged files.

**Capture from stdin:**
```bash
# Pipe the content of a webpage to a tag
curl -s https://en.wikipedia.org/wiki/Spaced_repetition | ./leitner tag --name=learning
```
This saves the content into a timestamped file within `~/.leitner/__tags__/<tagname>/`.

**Tag a file directly:**
```bash
./leitner tag --name=<tagname> --from-file=<filename>
```
This copies the specified file to the tag directory. If the tag doesn't exist, it will be created automatically.

**Delete a tagged file:**
```bash
./leitner tag delete --name=<tagname> --file=<filename>
```

---

### `config`

Manage LLM provider configuration.

**Set or update your provider:**
```bash
./leitner config set
```
This will launch an interactive prompt to guide you through selecting a provider and entering an API key.

**View your current configuration:**
```bash
./leitner config list
```
This displays the currently configured provider and a masked version of the API key.

---

### `generate`

Generate a new deck of flashcards from your tagged content using a configured LLM.

**Usage:**
```bash
./leitner generate deck --package=<packagename> --name=<deckname> --from-tag=<tagname> --cardcount=15
```
This command reads all the content in the specified tag, sends it to your LLM, and creates a new deck with the generated flashcards.
