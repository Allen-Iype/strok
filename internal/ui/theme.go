package ui

import (
	"strok/internal/domain"

	"github.com/charmbracelet/lipgloss"
)

// Theme holds all colors and styles. It is a value injected into the renderers
// so that alternative themes can be added later without changing rendering.
type Theme struct {
	finger [domain.FingerCount]lipgloss.Color
	// fingerDim is the muted variant of each finger color, used for the key
	// caps so the keyboard reads as a background reference while the lesson
	// text — in the bright palette — stays the visual focus.
	fingerDim [domain.FingerCount]lipgloss.Color
	keyBorder lipgloss.Color // uniform border for non-highlighted keys

	correct   lipgloss.Style // correctly typed text
	incorrect lipgloss.Style // wrongly typed text
	pending   lipgloss.Style // not-yet-typed text
	faint     lipgloss.Style // de-emphasized marks (space dots)
	cursor    lipgloss.Style // current character

	keyCurrent    lipgloss.Style // current expected key highlight
	keyCorrect    lipgloss.Style // brief green flash
	keyExpectErr  lipgloss.Style // expected key after a wrong press (yellow)
	keyWrongPress lipgloss.Style // the wrongly pressed key (red)

	header    lipgloss.Style
	stat      lipgloss.Style
	statLabel lipgloss.Style
	footer    lipgloss.Style
	box       lipgloss.Style
}

// DefaultTheme returns the built-in dark-terminal theme. Each finger gets a
// distinct, high-contrast color.
func DefaultTheme() Theme {
	fg := [domain.FingerCount]lipgloss.Color{
		domain.LPinky:  "141", // purple
		domain.LRing:   "75",  // blue
		domain.LMiddle: "78",  // green
		domain.LIndex:  "215", // orange
		domain.RIndex:  "209", // salmon
		domain.RMiddle: "114", // green2
		domain.RRing:   "39",  // blue2
		domain.RPinky:  "176", // magenta
		domain.Thumb:   "245", // gray
	}
	return Theme{
		finger: fg,

		correct:   lipgloss.NewStyle().Foreground(lipgloss.Color("78")),
		incorrect: lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Underline(true),
		pending:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		cursor:    lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Underline(true).Bold(true),

		keyCurrent:    lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("231")).Bold(true),
		keyCorrect:    lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("78")).Bold(true),
		keyExpectErr:  lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("220")).Bold(true),
		keyWrongPress: lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("196")).Bold(true),

		header:    lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Bold(true),
		stat:      lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Bold(true),
		statLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		footer:    lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		box:       lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1),
	}
}

// fingerStyle returns the normal foreground style for a finger.
func (t Theme) fingerStyle(f domain.Finger) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.finger[f])
}
