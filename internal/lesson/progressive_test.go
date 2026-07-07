package lesson

import (
	"math/rand"
	"strings"
	"testing"

	"strok/internal/domain"
)

func TestUnlockedForWidens(t *testing.T) {
	if got := unlockedFor(0); len(got) != 2 || got[0] != 'f' || got[1] != 'j' {
		t.Errorf("level 0 = %q, want fj", string(got))
	}
	if got := unlockedFor(1); len(got) != 3 {
		t.Errorf("level 1 len = %d, want 3", len(got))
	}
	if got := unlockedFor(1000); len(got) != len(unlockOrder) {
		t.Errorf("high level len = %d, want %d (capped)", len(got), len(unlockOrder))
	}
}

// TestUnlockOrderLettersBeforePunctuation pins the pedagogical intent of the
// progression: every alphabetic key unlocks before any punctuation, so the
// learner masters the full alphabet before the pinky-operated punctuation keys.
func TestUnlockOrderLettersBeforePunctuation(t *testing.T) {
	seenPunct := false
	letters := map[rune]bool{}
	for _, r := range unlockOrder {
		isLetter := r >= 'a' && r <= 'z'
		if isLetter {
			if seenPunct {
				t.Fatalf("letter %q unlocks after punctuation; letters must all come first", r)
			}
			letters[r] = true
		} else {
			seenPunct = true
		}
	}
	if len(letters) != 26 {
		t.Errorf("unlock order covers %d letters, want all 26", len(letters))
	}
}

func TestNextUsesOnlyUnlockedLetters(t *testing.T) {
	g := NewProgressive(rand.New(rand.NewSource(1)))
	p := domain.NewProfile()
	p.UnlockedLevel = 0
	l := g.Next(p)

	allowed := map[rune]bool{'f': true, 'j': true, ' ': true}
	for _, r := range l.Text {
		if !allowed[r] {
			t.Fatalf("lesson %q contains disallowed rune %q", l.Text, r)
		}
	}
	if strings.TrimSpace(l.Text) == "" {
		t.Fatal("lesson text is empty")
	}
}
