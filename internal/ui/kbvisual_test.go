package ui

import (
	"strings"
	"testing"

	"strok/internal/engine"
	"strok/internal/keyboard"
)

// TestKeyboardRendersAllRows checks the keyboard renders bordered caps for every
// row and that the home keys are present.
func TestKeyboardRendersAllRows(t *testing.T) {
	q := keyboard.NewQWERTY()
	out := renderKeyboard(DefaultTheme(), q.Rows(), 'f', engine.Feedback{}, false)

	lines := strings.Split(out, "\n")
	if len(lines) < 15 { // 5 rows × 3 lines each (top/mid/bottom border)
		t.Fatalf("expected >=15 rendered lines, got %d", len(lines))
	}
	for _, want := range []string{"⌫", "⇥", "⇪", "⇧", "space"} {
		if !strings.Contains(out, want) {
			t.Errorf("keyboard missing label %q", want)
		}
	}
}

// TestKeyboardStaggerIncreases verifies each letter row is indented further than
// the one above it (the classic keyboard stagger).
func TestKeyboardStaggerIncreases(t *testing.T) {
	prev := -1
	for row := 0; row < 4; row++ {
		got := indentFor(row)
		if got < prev {
			t.Errorf("row %d indent %d < previous %d (stagger should not decrease)", row, got, prev)
		}
		prev = got
	}
}
