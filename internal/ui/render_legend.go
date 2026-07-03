package ui

import (
	"strings"

	"strok/internal/domain"

	"github.com/charmbracelet/lipgloss"
)

// legendGroups orders the fingers as the hands sit on the board — left hand,
// right hand, thumb — so the legend mirrors the keyboard's left-to-right color
// order. Within a group the labels alone are unambiguous, so the L-/R-
// prefixes are dropped to keep the line compact.
var legendGroups = [][]domain.Finger{
	{domain.LPinky, domain.LRing, domain.LMiddle, domain.LIndex},
	{domain.RIndex, domain.RMiddle, domain.RRing, domain.RPinky},
	{domain.Thumb},
}

// renderLegend draws the finger-color legend as a single centered line: a
// swatch in each finger's color next to its label, groups separated by a faint
// dot. Swatches use the muted palette so the legend matches the resting board
// it describes.
func renderLegend(t Theme, width int) string {
	var groups []string
	for _, group := range legendGroups {
		var items []string
		for _, f := range group {
			label := strings.TrimPrefix(strings.TrimPrefix(f.ShortName(), "L-"), "R-")
			swatch := t.fingerDimStyle(f).Render("●")
			items = append(items, swatch+" "+t.statLabel.Render(label))
		}
		groups = append(groups, strings.Join(items, "  "))
	}

	legend := strings.Join(groups, t.faint.Render("   ·   "))
	if width > 0 {
		return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(legend)
	}
	return legend
}
