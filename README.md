# leitner

A simple, command-line based Leitner system for spaced repetition learning.

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

Create new packages or decks.

**Create a package:**
```bash
./leitner new package --name=<packagename>
```

**Create a deck:**
```bash
./leitner new deck --package=<packagename> --name=<deckname>
```

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

### `tag`

Capture and tag content from the web or other command-line tools.

**Usage:**
```bash
# Pipe the content of a webpage to a tag
curl -s https://en.wikipedia.org/wiki/Spaced_repetition | ./leitner tag --name=learning

# Pipe from a local file
cat my_notes.txt | ./leitner tag --name=personal
```
This saves the content into a timestamped file within `~/.leitner/__tags__/<tagname>/`.
