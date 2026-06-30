package keyboard

import (
	"testing"

	"strok/internal/domain"
)

func TestFindAndFingers(t *testing.T) {
	q := NewQWERTY()
	cases := map[rune]domain.Finger{
		'f': domain.LIndex, 'j': domain.RIndex, 'a': domain.LPinky,
		';': domain.RPinky, ' ': domain.Thumb, 'e': domain.LMiddle,
	}
	for r, want := range cases {
		k, ok := q.Find(r)
		if !ok {
			t.Errorf("Find(%q) not found", r)
			continue
		}
		if k.Finger != want {
			t.Errorf("Find(%q).Finger = %v, want %v", r, k.Finger, want)
		}
	}
	if _, ok := q.Find('€'); ok {
		t.Error("Find(€) should not be found")
	}
}

func TestEveryCharKeyHasFingerAssigned(t *testing.T) {
	q := NewQWERTY()
	for _, row := range q.Rows() {
		for _, k := range row {
			if k.Width <= 0 {
				t.Errorf("key %q has non-positive width", k.Label)
			}
		}
	}
}
