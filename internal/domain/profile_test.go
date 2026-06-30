package domain

import (
	"testing"
	"time"
)

func TestProfileApplyAggregates(t *testing.T) {
	p := NewProfile()
	p.Apply(Session{WPM: 40, Accuracy: 0.9, Duration: time.Minute})
	p.Apply(Session{WPM: 60, Accuracy: 1.0, Duration: time.Minute})

	if p.LessonsDone != 2 {
		t.Fatalf("LessonsDone = %d, want 2", p.LessonsDone)
	}
	if p.BestWPM != 60 {
		t.Errorf("BestWPM = %v, want 60", p.BestWPM)
	}
	if p.AvgWPM != 50 {
		t.Errorf("AvgWPM = %v, want 50", p.AvgWPM)
	}
	if p.Accuracy != 0.95 {
		t.Errorf("Accuracy = %v, want 0.95", p.Accuracy)
	}
	if p.PracticeTime != 2*time.Minute {
		t.Errorf("PracticeTime = %v, want 2m", p.PracticeTime)
	}
}

func TestProfileWeakKeys(t *testing.T) {
	p := NewProfile()
	p.Apply(Session{
		KeyHits:   map[rune]int{'a': 10, 'b': 10, 'c': 10},
		KeyErrors: map[rune]int{'a': 5, 'b': 1, 'c': 0},
	})
	got := p.WeakKeys(2)
	if len(got) != 2 || got[0] != 'a' || got[1] != 'b' {
		t.Errorf("WeakKeys = %q, want [a b]", string(got))
	}
}
