package ui

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"strok/internal/domain"
	"strok/internal/keyboard"
	"strok/internal/lesson"
	"strok/internal/mode"

	tea "github.com/charmbracelet/bubbletea"
)

type memStore struct{ p domain.Profile }

func (m *memStore) Load() (domain.Profile, error) { return m.p, nil }
func (m *memStore) Save(p domain.Profile) error   { m.p = p; return nil }

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

// stepClock returns a time that the test can advance, so elapsed-time-dependent
// behavior (WPM) can be controlled.
type stepClock struct{ t *time.Time }

func (c stepClock) Now() time.Time { return *c.t }

// typeLesson types the current lesson to completion exactly once, feeding each
// expected key. It captures the lesson length up front so that completing the
// lesson (which generates a fresh one) does not extend the loop.
func typeLesson(m Model) Model {
	for range m.state.Entries() {
		exp := m.state.Expected()
		mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{exp}})
		m = mm.(Model)
	}
	return m
}

// typeLessonWithClock is typeLesson but advances the clock by step per keystroke
// so elapsed-time-dependent stats (WPM) can be controlled.
func typeLessonWithClock(m Model, now *time.Time, step time.Duration) Model {
	for range m.state.Entries() {
		exp := m.state.Expected()
		mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{exp}})
		m = mm.(Model)
		*now = now.Add(step)
	}
	return m
}

func newTestModel() Model {
	deps := Deps{
		Layout:    keyboard.NewQWERTY(),
		Generator: lesson.NewProgressive(rand.New(rand.NewSource(1))),
		Store:     &memStore{p: domain.NewProfile()},
		Clock:     fixedClock{t: time.Unix(1000, 0)},
		Theme:     DefaultTheme(),
		Mode:      mode.NewProgressive(),
	}
	return New(deps, domain.NewProfile())
}

func TestSmokeRenderAndType(t *testing.T) {
	m := newTestModel()
	mm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = mm.(Model)

	out := m.View()
	if !strings.Contains(out, "strok") {
		t.Fatal("header not rendered")
	}
	if !strings.Contains(out, "WPM") {
		t.Fatal("stats bar not rendered")
	}
	if !strings.Contains(out, "L-index") {
		t.Fatal("finger legend not rendered")
	}
	if !strings.Contains(out, "●") {
		t.Fatal("legend swatch not rendered")
	}

	exp := m.state.Expected()
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{exp}})
	m = mm.(Model)
	if m.state.Cursor() != 1 {
		t.Fatalf("cursor = %d after correct key, want 1", m.state.Cursor())
	}
	_ = m.View() // must not panic
}

func TestSmokeTooSmall(t *testing.T) {
	m := newTestModel()
	mm, _ := m.Update(tea.WindowSizeMsg{Width: 10, Height: 10})
	m = mm.(Model)
	if !strings.Contains(m.View(), "too small") {
		t.Fatal("expected too-small guard")
	}
}

// TestGatedProgressionHoldsOnSlowLesson types a whole lesson with a constant
// clock (elapsed 0 → WPM 0), which fails the gate, so the unlock level must not
// advance and the "keep going" note must show.
func TestGatedProgressionHoldsOnSlowLesson(t *testing.T) {
	m := newTestModel()
	mm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = mm.(Model)

	startLevel := m.profile.UnlockedLevel
	m = typeLesson(m)

	if m.profile.UnlockedLevel != startLevel {
		t.Errorf("UnlockedLevel advanced on a 0-WPM lesson: %d -> %d", startLevel, m.profile.UnlockedLevel)
	}
	if m.outcome.Advanced {
		t.Error("outcome should not report Advanced for a failing lesson")
	}
	if !strings.Contains(m.View(), "keep going") {
		t.Error("expected the keep-going note after a failing lesson")
	}
}

// TestGatedProgressionAdvancesOnGoodLesson types a whole lesson correctly with a
// tiny elapsed time so WPM clears the threshold, and asserts the unlock level
// advances and the "unlocked" note shows.
func TestGatedProgressionAdvancesOnGoodLesson(t *testing.T) {
	now := time.Unix(1000, 0)
	deps := Deps{
		Layout:    keyboard.NewQWERTY(),
		Generator: lesson.NewProgressive(rand.New(rand.NewSource(1))),
		Store:     &memStore{p: domain.NewProfile()},
		Clock:     stepClock{&now},
		Theme:     DefaultTheme(),
		Mode:      mode.NewProgressive(),
	}
	m := New(deps, domain.NewProfile())
	mm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = mm.(Model)

	startLevel := m.profile.UnlockedLevel
	m = typeLessonWithClock(m, &now, 50*time.Millisecond)

	if m.profile.UnlockedLevel != startLevel+1 {
		t.Errorf("UnlockedLevel = %d, want %d after a clean fast lesson", m.profile.UnlockedLevel, startLevel+1)
	}
	if !m.outcome.Advanced {
		t.Error("outcome should report Advanced for a passing lesson")
	}
}
