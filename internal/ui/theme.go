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

	correct   lipgloss.Style // correct results (status line, key flash)
	typed     lipgloss.Style // already-typed correct text, dimmed so the remainder stays brightest
	incorrect lipgloss.Style // wrongly typed text
	pending   lipgloss.Style // not-yet-typed text
	faint     lipgloss.Style // de-emphasized marks (space dots)
	cursor    lipgloss.Style // current character

	keyCurrent    lipgloss.Style // current expected key highlight
	keyCorrect    lipgloss.Style // brief green flash
	keyExpectErr  lipgloss.Style // expected key after a wrong press (yellow)
	keyWrongPress lipgloss.Style // the wrongly pressed key (red)
	keyLocked     lipgloss.Style // locked letters and unused modifiers

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
	// Muted analog of each finger color: same hue family, lower intensity, so
	// the resting keyboard stays a readable finger-map without competing with
	// the full-saturation lesson text.
	dim := [domain.FingerCount]lipgloss.Color{
		domain.LPinky:  "97",  // muted purple
		domain.LRing:   "67",  // muted blue
		domain.LMiddle: "65",  // muted green
		domain.LIndex:  "137", // muted orange
		domain.RIndex:  "131", // muted salmon
		domain.RMiddle: "71",  // muted green2
		domain.RRing:   "31",  // muted blue2
		domain.RPinky:  "96",  // muted magenta
		domain.Thumb:   "240", // muted gray
	}
	return Theme{
		finger:    fg,
		fingerDim: dim,
		keyBorder: lipgloss.Color("238"),

		correct:   lipgloss.NewStyle().Foreground(lipgloss.Color("78")),
		typed:     lipgloss.NewStyle().Foreground(lipgloss.Color("65")),
		incorrect: lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("196")).Bold(true),
		pending:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		faint:     lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		cursor:    lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("231")).Bold(true),

		keyCurrent:    lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("231")).Bold(true),
		keyCorrect:    lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("78")).Bold(true),
		keyExpectErr:  lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("220")).Bold(true),
		keyWrongPress: lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("196")).Bold(true),
		keyLocked:     lipgloss.NewStyle().Foreground(lipgloss.Color("240")),

		header:    lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Bold(true),
		stat:      lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Bold(true),
		statLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		footer:    lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		box:       lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1),
	}
}

// fingerStyle returns the full-intensity foreground style for a finger, used
// where the finger colors carry the focus (lesson text, key highlights).
func (t Theme) fingerStyle(f domain.Finger) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.finger[f])
}

// fingerDimStyle returns the muted foreground style for a finger, used for the
// resting keyboard so it reads as a background reference.
func (t Theme) fingerDimStyle(f domain.Finger) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.fingerDim[f])
}
