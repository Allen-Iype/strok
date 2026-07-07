package lesson

import (
	"math/rand"
	"strings"
	"testing"

	"strok/internal/domain"
)

// pseudoAt generates one lesson at the given unlock level with a fixed seed.
func pseudoAt(t *testing.T, level int, seed int64) domain.Lesson {
	t.Helper()
	g := NewPseudoWord(rand.New(rand.NewSource(seed)))
	p := domain.NewProfile()
	p.UnlockedLevel = level
	return g.Next(p)
}

// TestPseudoWordUsesOnlyUnlockedKeys verifies every generated rune is either
// an unlocked key or the word separator, across a spread of levels and seeds.
func TestPseudoWordUsesOnlyUnlockedKeys(t *testing.T) {
	for _, level := range []int{6, 10, 16, 24, 28, 1000} {
		for seed := int64(0); seed < 20; seed++ {
			l := pseudoAt(t, level, seed)
			allowed := map[rune]bool{' ': true}
			for _, r := range l.Keyset {
				allowed[r] = true
			}
			for _, r := range l.Text {
				if !allowed[r] {
					t.Fatalf("level %d seed %d: lesson %q contains locked rune %q",
						level, seed, l.Text, r)
				}
			}
		}
	}
}

// TestPseudoWordEveryWordHasAVowel pins the pronounceability property: once a
// vowel is unlocked, every word contains at least one.
func TestPseudoWordEveryWordHasAVowel(t *testing.T) {
	for _, level := range []int{6, 12, 24, 1000} {
		for seed := int64(0); seed < 20; seed++ {
			l := pseudoAt(t, level, seed)
			for _, w := range strings.Fields(l.Text) {
				if !strings.ContainsAny(w, vowels) {
					t.Fatalf("level %d seed %d: word %q has no vowel in lesson %q",
						level, seed, w, l.Text)
				}
			}
		}
	}
}

// TestPseudoWordShape verifies the lesson keeps the drill generator's shape:
// 8 words within the length bounds (plus at most one trailing punctuation
// rune per word once punctuation is unlocked).
func TestPseudoWordShape(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		l := pseudoAt(t, 1000, seed) // full keyset, punctuation included
		words := strings.Fields(l.Text)
		if len(words) != 8 {
			t.Fatalf("seed %d: got %d words, want 8 (%q)", seed, len(words), l.Text)
		}
		for _, w := range words {
			core := strings.TrimRight(w, ";,./")
			if len(w)-len(core) > 1 {
				t.Errorf("seed %d: word %q has more than one trailing punctuation rune", seed, w)
			}
			if len(core) < 2 || len(core) > 6 {
				t.Errorf("seed %d: word core %q length %d outside [2,6]", seed, core, len(core))
			}
			if strings.ContainsAny(core, ";,./") {
				t.Errorf("seed %d: word %q has punctuation before the end", seed, w)
			}
		}
	}
}

// TestPseudoWordDeterministic verifies identical seeds produce identical
// lessons, keeping generation reproducible under test.
func TestPseudoWordDeterministic(t *testing.T) {
	a := pseudoAt(t, 10, 7)
	b := pseudoAt(t, 10, 7)
	if a.Text != b.Text {
		t.Errorf("same seed produced different lessons:\n%q\n%q", a.Text, b.Text)
	}
}

// TestPseudoWordSurvivesVowellessKeyset verifies the documented degradation:
// a keyset with no vowels (early levels) must not panic and must still emit
// keyset-only text.
func TestPseudoWordSurvivesVowellessKeyset(t *testing.T) {
	l := pseudoAt(t, 0, 1) // keyset {f, j}
	for _, r := range l.Text {
		if r != 'f' && r != 'j' && r != ' ' {
			t.Fatalf("lesson %q contains rune %q outside keyset", l.Text, r)
		}
	}
}
