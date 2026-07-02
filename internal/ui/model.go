// Package ui is the Bubble Tea presentation layer. The Model orchestrates the
// engine, stats, lesson generator and store; it holds no typing logic itself.
package ui

import (
	"time"

	"strok/internal/domain"
	"strok/internal/engine"
	"strok/internal/keyboard"
	"strok/internal/lesson"
	"strok/internal/mode"
	"strok/internal/stats"
	"strok/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

// Clock abstracts time so stats and animation are deterministic in tests.
type Clock interface{ Now() time.Time }

// RealClock is the production clock.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// flashDuration is how long a keystroke flash stays lit.
const flashDuration = 130 * time.Millisecond

// tickMsg drives time-based redraws (timer + flash expiry).
type tickMsg time.Time

// Deps bundles the injected collaborators so wiring stays in cmd/strok.
type Deps struct {
	Layout    keyboard.Layout
	Generator lesson.Generator
	Store     storage.Store
	Clock     Clock
	Theme     Theme
	Mode      mode.Mode
}

// Model is the root Bubble Tea model.
type Model struct {
	deps Deps

	profile domain.Profile
	state   *engine.TypingState

	width, height int

	startedAt time.Time // first keystroke time of the current lesson
	flashTill time.Time // when the current key flash expires

	justFinished bool         // shows the completion note this frame
	outcome      mode.Outcome // result of the last completed lesson
	quitting     bool
}

// New constructs the initial model from injected dependencies and a loaded
// profile.
func New(deps Deps, profile domain.Profile) Model {
	first := deps.Generator.Next(profile)
	return Model{
		deps:    deps,
		profile: profile,
		state:   engine.New(first),
	}
}

// Init starts the periodic tick.
func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// elapsed returns the time spent typing the current lesson so far.
func (m Model) elapsed() time.Duration {
	if !m.state.Started() {
		return 0
	}
	return m.deps.Clock.Now().Sub(m.startedAt)
}

// snapshot computes the current stats snapshot.
func (m Model) snapshot() domain.Stats {
	return stats.Compute(m.state, m.elapsed())
}
