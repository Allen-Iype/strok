package ui

import (
	"strings"

	"strok/internal/engine"

	"github.com/charmbracelet/lipgloss"
)

// renderText draws the lesson text with per-character coloring: correct in
// green, incorrect in red, the current character underlined, and the rest
// dimmed. Spaces at the cursor render as a visible underscore marker.
func renderText(t Theme, entries []engine.Entry, cursor, width int) string {
	var b strings.Builder
	for i, e := range entries {
		ch := string(e.Expected)
		switch {
		case i == cursor:
			disp := ch
			if e.Expected == ' ' {
				disp = "_"
			}
			b.WriteString(t.cursor.Render(disp))
		case e.Status == engine.Correct:
			b.WriteString(t.correct.Render(ch))
		case e.Status == engine.Incorrect:
			// Show the wrong character the user typed, in red.
			disp := string(e.Typed)
			if e.Typed == ' ' || e.Typed == 0 {
				disp = "_"
			}
			b.WriteString(t.incorrect.Render(disp))
		default:
			b.WriteString(t.pending.Render(ch))
		}
	}

	text := b.String()
	if width > 0 {
		return lipgloss.NewStyle().Width(width).Render(text)
	}
	return text
}
