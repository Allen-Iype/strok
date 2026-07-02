package mode

import "strok/internal/domain"

// Progressive is the gated progression mode: the learner unlocks the next letter
// only after a lesson that clears the WPM and accuracy thresholds. On a miss
// they stay at the current level and keep practicing the same keys.
type Progressive struct{}

// NewProgressive returns the gated progression mode.
func NewProgressive() Progressive { return Progressive{} }

func (Progressive) Name() string { return "progressive" }

// OnComplete advances the unlock level when the session clears the bar.
func (Progressive) OnComplete(p *domain.Profile, s domain.Session) Outcome {
	if domain.Advanced(s) {
		p.UnlockedLevel++
		return Outcome{Advanced: true, Message: "✓ nice — new key unlocked"}
	}
	return Outcome{Advanced: false, Message: "keep going — need 20 wpm & 90% to advance"}
}
