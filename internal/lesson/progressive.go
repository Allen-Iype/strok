package lesson

import (
	"math/rand"
	"strings"

	"strok/internal/domain"
)

// unlockOrder is the sequence in which letters are introduced. It starts on the
// home-row index fingers (f, j) and widens outward, matching the spec's
// progression (f, j, then more letters).
var unlockOrder = []rune{
	'f', 'j', // home index
	'd', 'k', // home middle
	's', 'l', // home ring
	'a', ';', // home pinky
	'g', 'h', // index stretch
	'e', 'i', // top middle
	'r', 'u', // top index
	'w', 'o', // top ring
	'q', 'p', // top pinky
	't', 'y', // top index stretch
	'c', ',', // bottom middle
	'v', 'm', // bottom index
	'x', '.', // bottom ring
	'z', '/', // bottom pinky
	'b', 'n', // bottom index stretch
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
