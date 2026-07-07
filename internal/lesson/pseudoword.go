package lesson

import (
	"math/rand"
	"strings"

	"strok/internal/domain"
)

// vowels are the rune classes pseudo-word phonotactics are built from. 'y' is
// treated as a consonant to keep the rules simple.
const vowels = "aeiou"

// PseudoWord is a drill generator that emits pronounceable non-words built
// from the unlocked keyset: consonants and vowels alternate keybr-style, so
// the text reads as word-shaped chunks instead of random noise. It requires
// at least one unlocked vowel; the Curriculum only delegates here once the
// keyset can support it.
type PseudoWord struct {
	rng *rand.Rand
	// Tuning knobs, mirroring Progressive so lesson shape and stats stay
	// comparable across textures.
	wordsPerLesson int
	minWordLen     int
	maxWordLen     int
	// doubleChance is the probability of doubling a letter mid-word ("dall");
	// punctChance is the probability a word ends with an unlocked punctuation
	// rune ("fask;"), keeping punctuation practice at natural word boundaries.
	doubleChance float64
	punctChance  float64
}

// NewPseudoWord returns a pseudo-word generator seeded by the given source.
func NewPseudoWord(rng *rand.Rand) *PseudoWord {
	return &PseudoWord{
		rng:            rng,
		wordsPerLesson: 8,
		minWordLen:     2,
		maxWordLen:     6,
		doubleChance:   0.15,
		punctChance:    0.3,
	}
}

// Next builds a lesson of pronounceable pseudo-words from the profile's
// unlocked keyset.
func (g *PseudoWord) Next(p domain.Profile) domain.Lesson {
	keyset := unlockedFor(p.UnlockedLevel)
	vs, cs, punct := classify(keyset)

	var b strings.Builder
	for i := 0; i < g.wordsPerLesson; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(g.word(vs, cs))
		if len(punct) > 0 && g.rng.Float64() < g.punctChance {
			b.WriteRune(punct[g.rng.Intn(len(punct))])
		}
	}

	return domain.Lesson{Text: b.String(), Keyset: keyset}
}

// word builds one pronounceable pseudo-word by alternating consonant and
// vowel positions, starting from either class, with an occasional doubled
// letter. Alternation guarantees a vowel in every word of length >= 2. If
// either class is empty the word degrades to the other class alone, so the
// generator never panics even if handed a keyset the Curriculum would not
// normally send here.
func (g *PseudoWord) word(vs, cs []rune) string {
	n := g.minWordLen + g.rng.Intn(g.maxWordLen-g.minWordLen+1)

	var b strings.Builder
	isVowel := g.rng.Intn(2) == 0
	hasVowel := false
	var prev rune
	for b.Len() < n {
		useVowel := (isVowel && len(vs) > 0) || len(cs) == 0
		class := cs
		if useVowel {
			class = vs
		}
		r := class[g.rng.Intn(len(class))]
		b.WriteRune(r)
		hasVowel = hasVowel || useVowel
		// Occasionally double a letter ("dall", "soo") — but never before the
		// word has its vowel, or a short consonant-start word could fill up
		// as all consonants ("ff") and lose pronounceability.
		if b.Len() < n && hasVowel && r != prev && g.rng.Float64() < g.doubleChance {
			b.WriteRune(r)
		}
		prev = r
		isVowel = !isVowel
	}
	return b.String()
}

// classify splits a keyset into vowels, consonant letters, and punctuation.
func classify(keyset []rune) (vs, cs, punct []rune) {
	for _, r := range keyset {
		switch {
		case strings.ContainsRune(vowels, r):
			vs = append(vs, r)
		case r >= 'a' && r <= 'z':
			cs = append(cs, r)
		default:
			punct = append(punct, r)
		}
	}
	return vs, cs, punct
}
