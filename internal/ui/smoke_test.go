package ui

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"strok/internal/domain"
	"strok/internal/keyboard"
	"strok/internal/lesson"

	tea "github.com/charmbracelet/bubbletea"
)

type memStore struct{ p domain.Profile }

func (m *memStore) Load() (domain.Profile, error) { return m.p, nil }
func (m *memStore) Save(p domain.Profile) error   { m.p = p; return nil }

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

func newTestModel() Model {
	deps := Deps{
		Layout:    keyboard.NewQWERTY(),
		Generator: lesson.NewProgressive(rand.New(rand.NewSource(1))),
		Store:     &memStore{p: domain.NewProfile()},
		Clock:     fixedClock{t: time.Unix(1000, 0)},
		Theme:     DefaultTheme(),
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
