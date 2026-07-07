package lesson

import (
	"math/rand"
	"strings"

	"strok/internal/domain"
)

// unlockOrder is the sequence in which keys are introduced. It starts on the
// home-row index fingers (f, j) and widens outward row by row — but letters
// only: all punctuation is deferred to a final phase so the learner builds
// confidence with the full alphabet before the awkward pinky-operated
// punctuation keys. Within that final phase, ';' comes first (home row, and
// the most code-relevant), then the bottom-row punctuation.
var unlockOrder = []rune{
	'f', 'j', // home index
	'd', 'k', // home middle
	's', 'l', // home ring
	'a', 'g', // home pinky + index stretch
	'h', 'e', // index stretch + top middle
	'i', 'r', // top middle + top index
	'u', 'w', // top index + top ring
	'o', 'q', // top ring + top pinky
	'p', 't', // top pinky + top index stretch
	'y', 'c', // top index stretch + bottom middle
	'v', 'm', // bottom index
	'x', 'z', // bottom ring + bottom pinky
	'b', 'n', // bottom index stretch
	';', ',', // punctuation: home pinky first…
	'.', '/', // …then the bottom row, once the letters are solid
}

// Progressive is a drill-style generator: it emits random groupings of the
// currently unlocked letters, widening the alphabet as the user levels up.
type Progressive struct {
	rng *rand.Rand
	// Tuning knobs kept as fields so they are easy to adjust or expose later.
	wordsPerLesson int
	minWordLen     int
	maxWordLen     int
}

// NewProgressive returns a generator seeded by the given source.
func NewProgressive(rng *rand.Rand) *Progressive {
	return &Progressive{
		rng:            rng,
		wordsPerLesson: 8,
		minWordLen:     2,
		maxWordLen:     5,
	}
}

// unlockedFor returns the letters available at a given level. Level 0 unlocks
// the first pair; each subsequent level adds one more letter.
func unlockedFor(level int) []rune {
	n := level + 2 // start with at least the first pair (f, j)
	if n > len(unlockOrder) {
		n = len(unlockOrder)
	}
	out := make([]rune, n)
	copy(out, unlockOrder[:n])
	return out
}

// Next builds the next lesson from the profile's unlocked level.
func (g *Progressive) Next(p domain.Profile) domain.Lesson {
	keyset := unlockedFor(p.UnlockedLevel)

	var b strings.Builder
	for i := 0; i < g.wordsPerLesson; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(g.word(keyset))
	}

	return domain.Lesson{Text: b.String(), Keyset: keyset}
}

// word builds one random "word" from the keyset. Very early levels (one or two
// letters) get short words so the first lessons look like "fj jf ff".
func (g *Progressive) word(keyset []rune) string {
	max := g.maxWordLen
	if len(keyset) <= 2 {
		max = min(max, 4)
	}
	n := g.minWordLen
	if max > g.minWordLen {
		n += g.rng.Intn(max - g.minWordLen + 1)
	}

	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteRune(keyset[g.rng.Intn(len(keyset))])
	}
	return b.String()
}
