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
