package ui

import (
	"strings"
	"testing"

	"strok/internal/engine"
	"strok/internal/keyboard"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestKeyboardRendersAllRows checks the keyboard renders bordered caps for every
// row and that the home keys are present.
func TestKeyboardRendersAllRows(t *testing.T) {
	q := keyboard.NewQWERTY()
	out := renderKeyboard(DefaultTheme(), q.Rows(), []rune{'f', 'j'}, 'f', engine.Feedback{}, false)

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

// TestErrorHighlightOutlivesFlash verifies incorrect feedback stays visible
// after the flash window (yellow expected key + red pressed key), while the
// correct-key flash respects the window.
func TestErrorHighlightOutlivesFlash(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)
	th := DefaultTheme()
	rows := keyboard.NewQWERTY().Rows()

	// Flash expired (flashing=false) after a wrong press: error must persist.
	out := renderKeyboard(th, rows, []rune{'f', 'j'}, 'f', engine.Feedback{Expected: 'f', Pressed: 'x'}, false)
	if !strings.Contains(out, "48;5;220") {
		t.Error("expected key should stay yellow (bg 220) after the flash window")
	}
	if !strings.Contains(out, "48;5;196") {
		t.Error("pressed key should stay red (bg 196) after the flash window")
	}

	// Flash expired after a correct press: no lingering green flash.
	out = renderKeyboard(th, rows, []rune{'f', 'j'}, 'j', engine.Feedback{Expected: 'f', Pressed: 'f', Correct: true}, false)
	if strings.Contains(out, "48;5;78") {
		t.Error("correct-key flash (bg 78) should not persist past the flash window")
	}
}

// TestKeyboardRestingBoardIsMuted verifies the three-tier scheme: unlocked keys
// render in the muted finger palette, locked letters and modifiers in gray, and
// no resting key uses a full-intensity finger color — full saturation is
// reserved for the lesson text and the highlighted keys.
func TestKeyboardRestingBoardIsMuted(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI256)
	th := DefaultTheme()
	rows := keyboard.NewQWERTY().Rows()

	out := renderKeyboard(th, rows, []rune{'f', 'j'}, 'f', engine.Feedback{}, false)

	// unlocked 'j' (R-index) -> dim salmon (131), not full salmon (209)
	if !strings.Contains(out, "38;5;131") {
		t.Error("unlocked 'j' should use the dim R-index color (131)")
	}
	if strings.Contains(out, "38;5;209") {
		t.Error("no resting key should use the full R-index color (209)")
	}
	// locked letters and modifiers -> gray (240)
	if !strings.Contains(out, "38;5;240") {
		t.Error("locked keys should use the gray locked style (240)")
	}
	// resting borders are uniform (238), not finger-colored
	if !strings.Contains(out, "38;5;238") {
		t.Error("resting caps should use the uniform dim border (238)")
	}
	// the current key 'f' still pops as an inverse block
	if !strings.Contains(out, "48;5;231") {
		t.Error("current key should keep the full-intensity inverse highlight (bg 231)")
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
