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
