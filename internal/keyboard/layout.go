// Package keyboard describes the physical keyboard layout and the assignment of
// each key to a finger. It is independent of colors (see ui.Theme) so that
// future layouts (Dvorak, Colemak) can be added without touching rendering.
package keyboard

import "strok/internal/domain"

// Layout describes a physical keyboard: its rows for rendering and a lookup from
// a typed rune to the physical key (which carries the finger assignment).
type Layout interface {
	// Rows returns the physical key rows, top to bottom.
	Rows() [][]domain.Key
	// Find maps a rune to its physical key. The second result is false if the
	// rune is not present on the layout.
	Find(r rune) (domain.Key, bool)
	// Name identifies the layout (e.g. "QWERTY").
	Name() string
}
