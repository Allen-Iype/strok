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

// TestRenderTextFingerColors verifies upcoming chars use their full finger
// color, typed chars dim to the muted green, and the cursor is an inverse
// block matching the keyboard's current-key highlight.
func TestRenderTextFingerColors(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)
	th := DefaultTheme()
	layout := keyboard.NewQWERTY()

	st := engine.New(domain.Lesson{Text: "fjj"})
	st.HandleKey('f') // f now correct, cursor on first j

	out := renderText(th, layout, st.Entries(), st.Cursor(), 0)

	// 'f' typed correctly -> dim green (65).
	if !strings.Contains(out, "38;5;65m") {
		t.Errorf("typed 'f' should be dim green (65); got %q", out)
	}
	// cursor 'j' -> inverse block (white bg 231).
	if !strings.Contains(out, "48;5;231") {
		t.Errorf("cursor 'j' should be an inverse block (bg 231); got %q", out)
	}
	// upcoming second 'j' -> full finger color R-index (209).
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
	if !strings.Contains(out, "48;5;196") {
		t.Errorf("mistyped char should be a red block (bg 196), matching the wrongly pressed key; got %q", out)
	}

	st.HandleKey('f') // correct it
	out = renderText(th, layout, st.Entries(), st.Cursor(), 0)
	if strings.Contains(out, "x") {
		t.Errorf("corrected position should no longer show 'x'; got %q", out)
	}
}
