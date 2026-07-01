package ui

import (
	"strok/internal/domain"

	"github.com/charmbracelet/lipgloss"
)

// legendRows groups the fingers into the two displayed rows: the left hand on
// top, the right hand plus thumb below.
var legendRows = [][]domain.Finger{
	{domain.LPinky, domain.LRing, domain.LMiddle, domain.LIndex},
	{domain.RIndex, domain.RMiddle, domain.RRing, domain.RPinky, domain.Thumb},
}

// renderLegend draws the finger-color legend: a swatch in each finger's color
// next to its short label, so the typist can map keyboard colors to fingers.
// Colors come straight from the theme's finger palette, matching the keys.
func renderLegend(t Theme, width int) string {
	var rows []string
	for _, group := range legendRows {
		var items []string
		for _, f := range group {
			swatch := t.fingerStyle(f).Render("●")
			items = append(items, swatch+" "+t.statLabel.Render(f.ShortName()))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, joinWithGap(items, "   ")...))
	}

	legend := lipgloss.JoinVertical(lipgloss.Center, rows...)
	if width > 0 {
		return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(legend)
	}
	return legend
}

// joinWithGap interleaves a gap string between items for JoinHorizontal.
func joinWithGap(items []string, gap string) []string {
	if len(items) == 0 {
		return items
	}
	out := make([]string, 0, len(items)*2-1)
	for i, it := range items {
		if i > 0 {
			out = append(out, gap)
		}
		out = append(out, it)
	}
	return out
}
