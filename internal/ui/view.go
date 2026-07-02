package ui

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// minHeight is the smallest terminal height the full layout needs. The bordered
// keyboard occupies 15 lines plus the finger legend ~3; the rest is header,
// stats, text, status line, footer and frame.
const minHeight = 29

// View composes the full frame each render. It recomputes layout from the
// current dimensions so terminal resizes are handled gracefully.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	rows := m.deps.Layout.Rows()
	kbWidth := keyboardWidth(rows)

	// Resize guard: if the terminal cannot fit the keyboard, show a hint.
	if m.width > 0 && (m.width < kbWidth || m.height < minHeight) {
		return m.tooSmallView(kbWidth)
	}

	t := m.deps.Theme
	flashing := m.deps.Clock.Now().Before(m.flashTill)
	snap := m.snapshot()

	header := renderHeader(t, m.deps.Layout.Name(), m.state.Keyset())
	statsBar := renderStats(t, snap)

	textWidth := kbWidth
	if m.width > 0 && m.width < textWidth {
		textWidth = m.width
	}
	text := renderText(t, m.deps.Layout, m.state.Entries(), m.state.Cursor(), textWidth)

	kb := renderKeyboard(t, rows, m.state.Expected(), m.state.Feedback(), flashing)
	legend := renderLegend(t, textWidth)
	footer := renderFooter(t)
	status := renderStatus(t, m.justFinished, m.outcome, textWidth)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		statsBar,
		"",
		text,
		status,
		"",
		kb,
		"",
		legend,
		"",
		footer,
	)

	framed := t.box.Render(body)
	if m.width > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, framed)
	}
	return framed
}

func (m Model) tooSmallView(kbWidth int) string {
	msg := "Terminal too small.\nResize to at least " +
		strconv.Itoa(kbWidth+4) + "×" + strconv.Itoa(minHeight+2) + " to play."
	box := m.deps.Theme.box.Render(msg)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
