package ui

import (
	"strok/internal/domain"
	"strok/internal/engine"

	"github.com/charmbracelet/lipgloss"
)

// keyHighlight describes how a key should be drawn this frame, beyond its
// normal finger color.
type keyHighlight int

const (
	hlNone keyHighlight = iota
	hlCurrent
	hlCorrect
	hlExpectErr
	hlWrongPress
)

// Keyboard geometry. unit is the horizontal columns one key-width occupies,
// including its borders, so cells in different rows line up on a fixed grid.
const (
	keyUnit = 5 // columns per Width unit; larger = bigger keys that fill the screen
	// keyPadY is the vertical padding inside each cap. 0 keeps keys single-line
	// (wide but not tall) so the board fits comfortably in most terminals.
	keyPadY = 0
	// keyboardLeftPad is a small left margin before the whole board.
	keyboardLeftPad = 1
)

// rowIndent is the leading-space offset per row that produces the classic
// staggered keyboard look (number row flush, each row below shifts right).
// Offsets scale with keyUnit so the stagger stays proportional to key size.
var rowIndent = []int{0, 3, 5, 7, 0}

// renderKeyboard draws the full keyboard as bordered key caps with the correct
// staggered offsets. The current expected key is always highlighted; recent
// feedback briefly colors the correct/incorrect keys.
func renderKeyboard(t Theme, rows [][]domain.Key, expected rune, fb engine.Feedback, flashing bool) string {
	var lines []string
	for ri, row := range rows {
		var cells []string
		for _, k := range row {
			cells = append(cells, renderKey(t, k, expected, fb, flashing))
		}
		joined := lipgloss.JoinHorizontal(lipgloss.Top, cells...)
		indent := keyboardLeftPad + indentFor(ri)
		lines = append(lines, lipgloss.NewStyle().MarginLeft(indent).Render(joined))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func indentFor(row int) int {
	if row < len(rowIndent) {
		return rowIndent[row]
	}
	return 0
}

// renderKey draws a single bordered key cap. Its total width is a multiple of
// keyUnit so keys align across rows regardless of the row stagger.
func renderKey(t Theme, k domain.Key, expected rune, fb engine.Feedback, flashing bool) string {
	hl := highlightFor(k, expected, fb, flashing)
	style := keyStyle(t, k, hl)

	// Inner width so the bordered cap spans Width*keyUnit columns.
	// A bordered box adds 2 cols (left+right border); cells share no borders so
	// each cap is independent. inner = Width*keyUnit - 2, floored at 1.
	inner := max(k.Width*keyUnit-2, 1)

	cap := style.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor(t, k, hl)).
		Padding(keyPadY, 0).
		Width(inner).
		Align(lipgloss.Center).
		Render(k.Label)
	return cap
}

func highlightFor(k domain.Key, expected rune, fb engine.Feedback, flashing bool) keyHighlight {
	if k.IsChar() && flashing && fb.Pressed != 0 {
		if !fb.Correct {
			if k.Rune == fb.Expected {
				return hlExpectErr
			}
			if k.Rune == fb.Pressed {
				return hlWrongPress
			}
		} else if k.Rune == fb.Expected {
			return hlCorrect
		}
	}
	if k.IsChar() && k.Rune == expected {
		return hlCurrent
	}
	return hlNone
}

func keyStyle(t Theme, k domain.Key, hl keyHighlight) lipgloss.Style {
	switch hl {
	case hlCurrent:
		return t.keyCurrent
	case hlCorrect:
		return t.keyCorrect
	case hlExpectErr:
		return t.keyExpectErr
	case hlWrongPress:
		return t.keyWrongPress
	default:
		return t.fingerStyle(k.Finger)
	}
}

// borderColor keeps the cap border subtle for normal keys but matches the
// highlight color for the current/feedback keys so they pop.
func borderColor(t Theme, k domain.Key, hl keyHighlight) lipgloss.Color {
	switch hl {
	case hlCurrent:
		return lipgloss.Color("231")
	case hlCorrect:
		return lipgloss.Color("78")
	case hlExpectErr:
		return lipgloss.Color("220")
	case hlWrongPress:
		return lipgloss.Color("196")
	default:
		return t.finger[k.Finger]
	}
}

// keyboardWidth returns the rendered width of the widest keyboard row including
// its stagger, used by the resize guard.
func keyboardWidth(rows [][]domain.Key) int {
	maxw := 0
	for ri, row := range rows {
		w := keyboardLeftPad + indentFor(ri)
		for _, k := range row {
			w += max(k.Width*keyUnit, 3) // each bordered cap spans Width*unit cols
		}
		maxw = max(maxw, w)
	}
	return maxw
}
