// Command demorec is a throwaway harness for recording the README demo GIF.
// It builds the real strok UI with a FIXED-seed generator so the lesson text
// is deterministic, then (with -print) emits that exact text so the vhs tape
// can type it back for an all-correct, all-green recording.
//
// This file is NOT part of the shipped app and is deleted after recording.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"strok/internal/domain"
	"strok/internal/keyboard"
	"strok/internal/lesson"
	"strok/internal/mode"
	"strok/internal/storage"
	"strok/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const seed = 42

func main() {
	printText := flag.Bool("print", false, "print the deterministic first lesson text and exit")
	flag.Parse()

	if *printText {
		gen := lesson.NewCurriculum(rand.New(rand.NewSource(seed)))
		fmt.Print(gen.Next(domain.NewProfile()).Text)
		return
	}

	// Record against a throwaway in-memory-ish profile so the recording always
	// starts from level 0 (the f/j home-row drill) regardless of real progress.
	tmp, _ := os.CreateTemp("", "strok-demo-*.json")
	store := storage.NewJSONStore(tmp.Name())
	profile := domain.NewProfile()

	deps := ui.Deps{
		Layout:    keyboard.NewQWERTY(),
		Generator: lesson.NewCurriculum(rand.New(rand.NewSource(seed))),
		Store:     store,
		Clock:     ui.RealClock{},
		Theme:     ui.DefaultTheme(),
		Mode:      mode.NewProgressive(),
	}

	program := tea.NewProgram(ui.New(deps, profile), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "demorec:", err)
		os.Exit(1)
	}
}
