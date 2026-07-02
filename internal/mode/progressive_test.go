package mode

import (
	"testing"

	"strok/internal/domain"
)

func TestProgressiveAdvancesOnPass(t *testing.T) {
	p := domain.NewProfile()
	m := NewProgressive()

	out := m.OnComplete(&p, domain.Session{WPM: 40, Accuracy: 0.95})
	if !out.Advanced {
		t.Error("passing session should report Advanced")
	}
	if p.UnlockedLevel != 1 {
		t.Errorf("UnlockedLevel = %d, want 1 after a passing lesson", p.UnlockedLevel)
	}
}

func TestProgressiveHoldsOnFail(t *testing.T) {
	p := domain.NewProfile()
	p.UnlockedLevel = 3
	m := NewProgressive()

	out := m.OnComplete(&p, domain.Session{WPM: 10, Accuracy: 0.80})
	if out.Advanced {
		t.Error("failing session should not report Advanced")
	}
	if p.UnlockedLevel != 3 {
		t.Errorf("UnlockedLevel = %d, want 3 (unchanged) after a failing lesson", p.UnlockedLevel)
	}
	if out.Message == "" {
		t.Error("outcome should carry a message")
	}
}
