package ui

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// minHeight is the smallest terminal height the full layout needs: the framed
// body — bordered keyboard (15 lines), one-line legend, header, stats, the
// lesson/status block in its whitespace pocket, footer — plus the frame border.
const minHeight = 30

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

	textWidth := kbWidth
	if m.width > 0 && m.width < textWidth {
		textWidth = m.width
	}
	// Everything in the play loop shares one centered axis; the header alone
	// spans the width as a balanced top bar.
	center := lipgloss.NewStyle().Width(textWidth).Align(lipgloss.Center)

	header := renderHeader(t, m.deps.Layout.Name(), m.state.Keyset(), textWidth)
	statsBar := center.Render(renderStats(t, snap))
	text := renderText(t, m.deps.Layout, m.state.Entries(), m.state.Cursor(), textWidth)

	kb := renderKeyboard(t, rows, m.state.Keyset(), m.state.Expected(), m.state.Feedback(), flashing)
	legend := renderLegend(t, textWidth)
	footer := center.Render(renderFooter(t))
	status := renderStatus(t, m.justFinished, m.outcome, m.lastResult, textWidth)

	// The lesson/status block sits in a double-blank pocket of whitespace —
	// the terminal's only "font size" — so the eye lands on it first.
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		statsBar,
		"",
		"",
		text,
		status,
		"",
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
		strconv.Itoa(kbWidth+4) + "×" + strconv.Itoa(minHeight+2) +
		" to play.\nnow " + strconv.Itoa(m.width) + "×" + strconv.Itoa(m.height) + "."
	box := m.deps.Theme.box.Render(msg)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
