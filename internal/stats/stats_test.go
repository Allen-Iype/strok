package stats

import (
	"testing"
	"time"
)

type fake struct{ keystrokes, errors, cursor int }

func (f fake) Keystrokes() int { return f.keystrokes }
func (f fake) Errors() int     { return f.errors }
func (f fake) Cursor() int     { return f.cursor }

func TestComputeZeroInput(t *testing.T) {
	got := Compute(fake{}, 0)
	if got.WPM != 0 {
		t.Errorf("WPM = %v, want 0", got.WPM)
	}
	if got.Accuracy != 1.0 {
		t.Errorf("Accuracy = %v, want 1.0", got.Accuracy)
	}
}

func TestComputeWPMAndAccuracy(t *testing.T) {
	// 25 correct chars in 30s = 5 words / 0.5min = 10 WPM.
	got := Compute(fake{keystrokes: 30, errors: 5, cursor: 25}, 30*time.Second)
	if got.WPM != 10 {
		t.Errorf("WPM = %v, want 10", got.WPM)
	}
	// 25 correct / 30 total = 0.8333...
	if got.Accuracy < 0.83 || got.Accuracy > 0.834 {
		t.Errorf("Accuracy = %v, want ~0.833", got.Accuracy)
	}
}
