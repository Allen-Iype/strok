package lesson

import (
	"math/rand"
	"strings"

	"strok/internal/domain"
)

// Curriculum composes the texture generators and owns the graduation policy:
// each lesson is delegated to the richest texture the current keyset can
// support. Delegation is capability-driven (what the keyset contains), not
// level-driven, so reordering the unlock sequence can never break it. Future
// textures (real vocabulary, programming lessons) slot in as additional
// delegates with their own capability rule.
type Curriculum struct {
	drill  *Progressive
	pseudo *PseudoWord
}

// NewCurriculum returns the standard texture progression seeded by the given
// source: random drills while the keyset is consonant-only, pronounceable
// pseudo-words from the first unlocked vowel on.
func NewCurriculum(rng *rand.Rand) *Curriculum {
	return &Curriculum{
		drill:  NewProgressive(rng),
		pseudo: NewPseudoWord(rng),
	}
}

// Next delegates to the richest texture the profile's keyset supports.
func (c *Curriculum) Next(p domain.Profile) domain.Lesson {
	keyset := unlockedFor(p.UnlockedLevel)
	if strings.ContainsAny(string(keyset), vowels) {
		return c.pseudo.Next(p)
	}
	return c.drill.Next(p)
}
