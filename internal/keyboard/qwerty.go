package keyboard

import "strok/internal/domain"

// fingerOf maps a character rune to its touch-typing finger. Only the runes a
// lesson can emit (lowercase letters, digits, common symbols, space) need an
// entry; non-character keys carry their finger directly in the row tables.
var fingerOf = map[rune]domain.Finger{
	// number row
	'`': domain.LPinky, '1': domain.LPinky, '2': domain.LRing, '3': domain.LMiddle,
	'4': domain.LIndex, '5': domain.LIndex, '6': domain.RIndex, '7': domain.RIndex,
	'8': domain.RMiddle, '9': domain.RRing, '0': domain.RPinky, '-': domain.RPinky, '=': domain.RPinky,
	// top row
	'q': domain.LPinky, 'w': domain.LRing, 'e': domain.LMiddle, 'r': domain.LIndex, 't': domain.LIndex,
	'y': domain.RIndex, 'u': domain.RIndex, 'i': domain.RMiddle, 'o': domain.RRing, 'p': domain.RPinky,
	'[': domain.RPinky, ']': domain.RPinky, '\\': domain.RPinky,
	// home row
	'a': domain.LPinky, 's': domain.LRing, 'd': domain.LMiddle, 'f': domain.LIndex, 'g': domain.LIndex,
	'h': domain.RIndex, 'j': domain.RIndex, 'k': domain.RMiddle, 'l': domain.RRing,
	';': domain.RPinky, '\'': domain.RPinky,
	// bottom row
	'z': domain.LPinky, 'x': domain.LRing, 'c': domain.LMiddle, 'v': domain.LIndex, 'b': domain.LIndex,
	'n': domain.RIndex, 'm': domain.RIndex, ',': domain.RMiddle, '.': domain.RRing, '/': domain.RPinky,
	// thumb
	' ': domain.Thumb,
}

// QWERTY is the US ANSI QWERTY layout.
type QWERTY struct {
	rows  [][]domain.Key
	index map[rune]domain.Key
}

// NewQWERTY builds the QWERTY layout.
func NewQWERTY() *QWERTY {
	q := &QWERTY{index: map[rune]domain.Key{}}
	q.rows = q.build()
	for _, row := range q.rows {
		for _, k := range row {
			if k.IsChar() {
				q.index[k.Rune] = k
			}
		}
	}
	return q
}

func (q *QWERTY) Rows() [][]domain.Key { return q.rows }
func (q *QWERTY) Name() string         { return "QWERTY" }

func (q *QWERTY) Find(r rune) (domain.Key, bool) {
	k, ok := q.index[r]
	return k, ok
}

// charKeys turns a string of characters into a row of single-width Key values,
// looking up each finger from fingerOf.
func charKeys(row int, chars string) []domain.Key {
	keys := make([]domain.Key, 0, len(chars))
	for _, r := range chars {
		keys = append(keys, domain.Key{
			Rune:   r,
			Label:  string(r),
			Finger: fingerOf[r],
			Row:    row,
			Width:  1,
		})
	}
	return keys
}

// mod builds a non-character modifier key (Shift, Tab, ...) of the given width.
func mod(row int, label string, f domain.Finger, width int) domain.Key {
	return domain.Key{Label: label, Finger: f, Row: row, Width: width}
}

func (q *QWERTY) build() [][]domain.Key {
	r0 := append(charKeys(0, "`1234567890-="), mod(0, "⌫", domain.RPinky, 2))

	r1 := append([]domain.Key{mod(1, "⇥", domain.LPinky, 2)}, charKeys(1, "qwertyuiop[]")...)
	r1 = append(r1, mod(1, "\\", domain.RPinky, 2))

	r2 := append([]domain.Key{mod(2, "⇪", domain.LPinky, 2)}, charKeys(2, "asdfghjkl;'")...)
	r2 = append(r2, mod(2, "⏎", domain.RPinky, 3))

	r3 := append([]domain.Key{mod(3, "⇧", domain.LPinky, 3)}, charKeys(3, "zxcvbnm,./")...)
	r3 = append(r3, mod(3, "⇧", domain.RPinky, 3))

	r4 := []domain.Key{
		mod(4, "ctrl", domain.LPinky, 2),
		mod(4, "alt", domain.LIndex, 2),
		{Rune: ' ', Label: "space", Finger: domain.Thumb, Row: 4, Width: 10},
		mod(4, "alt", domain.RIndex, 2),
		mod(4, "ctrl", domain.RPinky, 2),
	}

	return [][]domain.Key{r0, r1, r2, r3, r4}
}
