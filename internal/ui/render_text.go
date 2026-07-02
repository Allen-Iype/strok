package ui

import (
	"strings"

	"strok/internal/engine"
	"strok/internal/keyboard"

	"github.com/charmbracelet/lipgloss"
)

// renderText draws the lesson text with per-character coloring. A mistyped
// character shows the wrongly typed rune in the error style right at the
// cursor, so a typo is visible where the eye is focused until it is corrected.
// Correctly typed characters are green. Not-yet-typed characters — including
// the current one — are colored by their finger so the typist can see which
// finger to use before pressing; the current character is additionally
// underlined and bold. Spaces show as a middle dot, or an underscore at the
// cursor.
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
			b.WriteString(fingerStyleFor(t, layout, e.Expected).Underline(true).Bold(true).Render(disp))
		case e.Status == engine.Correct:
			b.WriteString(t.correct.Render(ch))
		default:
			b.WriteString(fingerStyleFor(t, layout, e.Expected).Render(ch))
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
