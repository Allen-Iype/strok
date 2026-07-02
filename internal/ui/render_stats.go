package ui

import (
	"fmt"
	"strings"
	"time"

	"strok/internal/domain"
	"strok/internal/mode"

	"github.com/charmbracelet/lipgloss"
)

// renderHeader draws the title and the active keyset.
func renderHeader(t Theme, layoutName string, keyset []rune) string {
	title := t.header.Render("⌨  strok")
	sub := t.statLabel.Render(fmt.Sprintf("· %s · keys: %s", layoutName, spaced(keyset)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, title, " ", sub)
}

// renderStats draws the live statistics bar. Every value sits in a fixed-width
// slot so changing digit counts never shift the bar mid-lesson, and the WPM and
// ACC values turn green once they clear the advance thresholds, showing at a
// glance whether the lesson is on pace to unlock the next key.
func renderStats(t Theme, s domain.Stats) string {
	wpm := s.WPM
	if wpm > 999 {
		wpm = 999 // clamp the first-keystroke spike so its slot never overflows
	}

	passing := t.correct.Bold(true)
	wpmStyle, accStyle := t.stat, t.stat
	if s.WPM >= domain.AdvanceWPM {
		wpmStyle = passing
	}
	if s.Accuracy >= domain.AdvanceAccuracy {
		accStyle = passing
	}

	parts := []string{
		statCell(t, "WPM", fmt.Sprintf("%.0f", wpm), 3, wpmStyle),
		statCell(t, "ACC", fmt.Sprintf("%.0f%%", s.Accuracy*100), 4, accStyle),
		statCell(t, "ERR", fmt.Sprintf("%d", s.Errors), 3, t.stat),
		statCell(t, "CHARS", fmt.Sprintf("%d", s.Typed), 3, t.stat),
		statCell(t, "TIME", formatDuration(s.Elapsed), 5, t.stat),
	}
	return strings.Join(parts, "   ")
}

// renderStatus draws the one-line status area under the lesson text, centered
// to the same width. It always occupies exactly one line — blank when there is
// nothing to report — so a message appearing or clearing never resizes the
// frame. After a lesson it leads with the measured result so the learner sees
// what they scored against the advance gate.
func renderStatus(t Theme, show bool, o mode.Outcome, result domain.Stats, width int) string {
	line := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	if !show {
		return line.Render("")
	}
	style := t.pending
	if o.Advanced {
		style = t.correct
	}
	res := t.stat.Render(fmt.Sprintf("%.0f wpm · %.0f%%", result.WPM, result.Accuracy*100))
	return line.Render(res + "  " + style.Render(o.Message))
}

// renderFooter draws the key hints.
func renderFooter(t Theme) string {
	return t.footer.Render("esc/ctrl+c quit · backspace correct · tab restart lesson")
}

// statCell renders one label + value pair, the value left-aligned in a
// fixed-width slot so the bar's geometry is independent of the value.
func statCell(t Theme, label, value string, width int, style lipgloss.Style) string {
	return t.statLabel.Render(label+" ") + style.Render(fmt.Sprintf("%-*s", width, value))
}

func spaced(rs []rune) string {
	if len(rs) == 0 {
		return "-"
	}
	parts := make([]string, len(rs))
	for i, r := range rs {
		parts[i] = string(r)
	}
	return strings.Join(parts, " ")
}

func formatDuration(d time.Duration) string {
	total := int(d.Seconds())
	return fmt.Sprintf("%02d:%02d", total/60, total%60)
}
