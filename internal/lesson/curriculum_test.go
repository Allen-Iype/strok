package lesson

import (
	"math/rand"
	"strings"
	"testing"

	"strok/internal/domain"
)

// curriculumAt generates one lesson at the given unlock level.
func curriculumAt(level int, seed int64) domain.Lesson {
	g := NewCurriculum(rand.New(rand.NewSource(seed)))
	p := domain.NewProfile()
	p.UnlockedLevel = level
	return g.Next(p)
}

// TestCurriculumDrillsBeforeFirstVowel verifies the earliest levels — whose
// keyset has no vowel — fall back to the drill texture and stay within the
// keyset.
func TestCurriculumDrillsBeforeFirstVowel(t *testing.T) {
	l := curriculumAt(0, 1) // keyset {f, j}
	if strings.TrimSpace(l.Text) == "" {
		t.Fatal("lesson text is empty")
	}
	for _, r := range l.Text {
		if r != 'f' && r != 'j' && r != ' ' {
			t.Fatalf("lesson %q contains rune %q outside keyset", l.Text, r)
		}
	}
}

// TestCurriculumGraduatesToPseudoWordsAtFirstVowel pins the graduation rule:
// from the level that unlocks 'a' onward, every word is pronounceable (has a
// vowel) — the pseudo-word texture's defining property, which the uniform
// random drill cannot guarantee.
func TestCurriculumGraduatesToPseudoWordsAtFirstVowel(t *testing.T) {
	for _, level := range []int{5, 12, 1000} { // level 5 unlocks 'a' (7th key)
		for seed := int64(0); seed < 20; seed++ {
			l := curriculumAt(level, seed)
			for _, w := range strings.Fields(l.Text) {
				if !strings.ContainsAny(w, vowels) {
					t.Fatalf("level %d seed %d: word %q has no vowel — expected pseudo-word texture (%q)",
						level, seed, w, l.Text)
				}
			}
		}
	}
}
