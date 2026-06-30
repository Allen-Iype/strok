package domain

// Key describes a single physical key on the keyboard.
//
// Rune is the primary lowercase character produced by the key (e.g. 'a'); it is
// zero for non-character keys such as Shift or Tab. Label is what gets drawn in
// the cell. Width is a relative cell width used to render wide keys (Space,
// Enter, ...) — a normal letter key has Width 1.
type Key struct {
	Rune   rune
	Label  string
	Finger Finger
	Row    int
	Width  int
}

// IsChar reports whether the key produces a typeable character.
func (k Key) IsChar() bool { return k.Rune != 0 }
