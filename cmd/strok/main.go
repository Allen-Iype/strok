// Command strok is a terminal-based typing app.
//
// Run it with:
//
//	go run ./cmd/strok
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"strok/internal/keyboard"
	"strok/internal/lesson"
	"strok/internal/mode"
	"strok/internal/storage"
	"strok/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "strok:", err)
		os.Exit(1)
	}
}

func run() error {
	dataPath := flag.String("data", "", "path to the profile JSON file (default: OS config dir)")
	flag.Parse()

	path := *dataPath
	if path == "" {
		p, err := storage.DefaultPath()
		if err != nil {
			return err
		}
		path = p
	}

	store := storage.NewJSONStore(path)
	profile, err := store.Load()
	if err != nil {
		return fmt.Errorf("load profile: %w", err)
	}

	deps := ui.Deps{
		Layout:    keyboard.NewQWERTY(),
		Generator: lesson.NewCurriculum(rand.New(rand.NewSource(time.Now().UnixNano()))),
		Store:     store,
		Clock:     ui.RealClock{},
		Theme:     ui.DefaultTheme(),
		Mode:      mode.NewProgressive(),
	}

	model := ui.New(deps, profile)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("run program: %w", err)
	}

	printSummary(path, store)
	return nil
}

// printSummary prints a one-line progress summary after the TUI exits.
func printSummary(path string, store *storage.JSONStore) {
	p, err := store.Load()
	if err != nil {
		return
	}
	fmt.Printf("strok · lessons: %d · best WPM: %.0f · avg WPM: %.0f · accuracy: %.0f%%\n",
		p.LessonsDone, p.BestWPM, p.AvgWPM, p.Accuracy*100)
	fmt.Printf("progress saved to %s\n", path)
}
