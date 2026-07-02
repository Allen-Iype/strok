package ui

import (
	"time"

	"strok/internal/engine"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages: window resize, ticks, and key input.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tickMsg:
		// Periodic redraw keeps the timer live; re-arm the tick.
		return m, tick()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m.quit()

	case tea.KeyBackspace:
		m.state.Backspace()
		return m, nil

	case tea.KeyTab:
		return m.restartLesson(), nil

	case tea.KeyRunes, tea.KeySpace:
		return m.typeRunes(msg)
	}
	return m, nil
}

// typeRunes feeds printable runes into the engine, starting the timer on the
// first keystroke and advancing to the next lesson on completion.
func (m Model) typeRunes(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	runes := msg.Runes
	if msg.Type == tea.KeySpace {
		runes = []rune{' '}
	}
	for _, r := range runes {
		if !m.state.Started() {
			m.startedAt = m.deps.Clock.Now()
		}
		m.state.HandleKey(r)
		m.flashTill = m.deps.Clock.Now().Add(flashDuration)
		m.justFinished = false

		if m.state.Done() {
			m.completeLesson()
			break
		}
	}
	return m, nil
}

// completeLesson records the session, persists progress, advances the level and
// generates the next lesson.
func (m *Model) completeLesson() {
	snap := m.snapshot()
	session := m.state.Session(snap, m.elapsed())

	m.profile.Apply(session)
	m.outcome = m.deps.Mode.OnComplete(&m.profile, session)
	_ = m.deps.Store.Save(m.profile)

	next := m.deps.Generator.Next(m.profile)
	m.state = engine.New(next)
	m.startedAt = time.Time{}
	m.flashTill = time.Time{}
	m.justFinished = true
}

func (m Model) restartLesson() Model {
	// Regenerate at the current level for a fresh attempt.
	next := m.deps.Generator.Next(m.profile)
	m.state = engine.New(next)
	m.startedAt = time.Time{}
	m.flashTill = time.Time{}
	m.justFinished = false
	return m
}

func (m Model) quit() (tea.Model, tea.Cmd) {
	_ = m.deps.Store.Save(m.profile)
	m.quitting = true
	return m, tea.Quit
}
