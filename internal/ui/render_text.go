package ui

import (
	"strings"

	"strok/internal/engine"
	"strok/internal/keyboard"

	"github.com/charmbracelet/lipgloss"
)

// renderText draws the lesson text with per-character coloring. The line is
// the brightest element in the frame: not-yet-typed characters carry the
// full-saturation finger colors, already-typed characters dim to a muted
// green, and word-separator dots are faint — so what's ahead stays brightest
// and the bright/dim boundary tracks the typist's position. The cursor is an
// inverse block matching the keyboard's current-key highlight, and a mistyped
// character shows the wrongly typed rune as a red block right at the cursor
// (matching the wrongly pressed key) until it is corrected. Spaces show as a
// middle dot, or an underscore at the cursor.
func renderText(t Theme, layout keyboard.Layout, entries []engine.Entry, cursor, width int) string {
	var b strings.Builder
	for i, e := range entries {
		ch := string(e.Expected)
		if e.Expected == ' ' {
			ch = "·"
		}
		switch {
		case e.Status == engine.Incorrect:
			// The engine holds the cursor on an error, so this is the cursor
			// cell: show what was actually typed until it is fixed.
			disp := string(e.Typed)
			if e.Typed == ' ' || e.Typed == 0 {
				disp = "_"
			}
			b.WriteString(t.incorrect.Render(disp))
		case i == cursor:
			disp := ch
			if e.Expected == ' ' {
				disp = "_"
			}
			b.WriteString(t.cursor.Render(disp))
		case e.Status == engine.Correct:
			if e.Expected == ' ' {
				b.WriteString(t.faint.Render(ch))
			} else {
				b.WriteString(t.typed.Render(ch))
			}
		default:
			if e.Expected == ' ' {
				b.WriteString(t.faint.Render(ch))
			} else {
				b.WriteString(fingerStyleFor(t, layout, e.Expected).Render(ch))
			}
		}
	}

	text := b.String()
	if width > 0 {
		return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(text)
	}
	return text
}

// fingerStyleFor returns the finger-colored style for the key that types r,
// falling back to the dim pending style for runes not on the layout.
func fingerStyleFor(t Theme, layout keyboard.Layout, r rune) lipgloss.Style {
	if k, ok := layout.Find(r); ok {
		return t.fingerStyle(k.Finger)
	}
	return t.pending
}
