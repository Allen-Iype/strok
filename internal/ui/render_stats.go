package ui

import (
	"fmt"
	"strings"
	"time"

	"strok/internal/domain"

	"github.com/charmbracelet/lipgloss"
)

// renderHeader draws the title and the active keyset.
func renderHeader(t Theme, layoutName string, keyset []rune) string {
	title := t.header.Render("⌨  strok")
	sub := t.statLabel.Render(fmt.Sprintf("· %s · keys: %s", layoutName, spaced(keyset)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, title, " ", sub)
}

// renderStats draws the live statistics bar.
func renderStats(t Theme, s domain.Stats) string {
	parts := []string{
		statCell(t, "WPM", fmt.Sprintf("%.0f", s.WPM)),
		statCell(t, "ACC", fmt.Sprintf("%.0f%%", s.Accuracy*100)),
		statCell(t, "ERR", fmt.Sprintf("%d", s.Errors)),
		statCell(t, "CHARS", fmt.Sprintf("%d", s.Typed)),
		statCell(t, "TIME", formatDuration(s.Elapsed)),
	}
	return strings.Join(parts, "   ")
}

// renderFooter draws the key hints.
func renderFooter(t Theme) string {
	return t.footer.Render("esc/ctrl+c quit · backspace correct · tab restart lesson")
}

func statCell(t Theme, label, value string) string {
	return t.statLabel.Render(label+" ") + t.stat.Render(value)
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
