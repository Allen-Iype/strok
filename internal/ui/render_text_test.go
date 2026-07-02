package ui

import (
	"strings"
	"testing"

	"strok/internal/domain"
	"strok/internal/engine"
	"strok/internal/keyboard"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestRenderTextFingerColors verifies upcoming chars use their finger color
// while typed chars keep their status color.
func TestRenderTextFingerColors(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)
	th := DefaultTheme()
	layout := keyboard.NewQWERTY()

	st := engine.New(domain.Lesson{Text: "fj"})
	st.HandleKey('f') // f now correct, cursor on j

	out := renderText(th, layout, st.Entries(), st.Cursor(), 0)

	// 'f' typed correctly -> green (78); 'j' upcoming/cursor -> R-index (209).
	if !strings.Contains(out, "38;5;78m") {
		t.Errorf("typed 'f' should be green (78); got %q", out)
	}
	if !strings.Contains(out, "38;5;209") {
		t.Errorf("upcoming 'j' should be finger-colored R-index (209); got %q", out)
	}
}

// TestRenderTextShowsError verifies a wrong keystroke is visible at the cursor:
// the wrongly typed rune renders in the error color instead of the expected
// character, until it is corrected.
func TestRenderTextShowsError(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)
	th := DefaultTheme()
	layout := keyboard.NewQWERTY()

	st := engine.New(domain.Lesson{Text: "fj"})
	st.HandleKey('x') // wrong: expected 'f'

	out := renderText(th, layout, st.Entries(), st.Cursor(), 0)

	if !strings.Contains(out, "x") {
		t.Errorf("mistyped 'x' should be shown at the cursor; got %q", out)
	}
	if !strings.Contains(out, "38;5;203") {
		t.Errorf("mistyped char should use the error color (203); got %q", out)
	}

	st.HandleKey('f') // correct it
	out = renderText(th, layout, st.Entries(), st.Cursor(), 0)
	if strings.Contains(out, "x") {
		t.Errorf("corrected position should no longer show 'x'; got %q", out)
	}
}
