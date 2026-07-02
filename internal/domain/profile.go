package domain

import (
	"sort"
	"time"
)

// ProfileVersion is the current on-disk schema version. Bump it when the shape
// of Profile changes so older files can be migrated rather than rejected.
const ProfileVersion = 1

// Profile is the persisted aggregate of a user's progress.
type Profile struct {
	Version       int                `json:"version"`
	BestWPM       float64            `json:"best_wpm"`
	AvgWPM        float64            `json:"avg_wpm"`
	Accuracy      float64            `json:"accuracy"` // lifetime, 0..1
	PracticeTime  time.Duration      `json:"practice_time"`
	LessonsDone   int                `json:"lessons_done"`
	UnlockedLevel int                `json:"unlocked_level"` // drives the generator
	KeyStats      map[string]KeyStat `json:"key_stats"`
}

// KeyStat accumulates per-key attempt counts used to find weak keys.
type KeyStat struct {
	Presses int `json:"presses"`
	Errors  int `json:"errors"`
}

// NewProfile returns a fresh profile at the starting level.
func NewProfile() Profile {
	return Profile{Version: ProfileVersion, KeyStats: map[string]KeyStat{}}
}

// Thresholds a completed lesson must clear to unlock the next letter.
const (
	AdvanceWPM      = 20.0
	AdvanceAccuracy = 0.90
)

// Advanced reports whether a completed session clears the bar to unlock more
// content. It is a pure predicate so progression rules stay testable and
// mode-agnostic.
func Advanced(s Session) bool {
	return s.WPM >= AdvanceWPM && s.Accuracy >= AdvanceAccuracy
}

// Apply folds a completed session into the profile's running aggregates.
func (p *Profile) Apply(s Session) {
	if p.KeyStats == nil {
		p.KeyStats = map[string]KeyStat{}
	}

	// Running averages weighted by lesson count.
	n := float64(p.LessonsDone)
	p.AvgWPM = (p.AvgWPM*n + s.WPM) / (n + 1)
	p.Accuracy = (p.Accuracy*n + s.Accuracy) / (n + 1)

	if s.WPM > p.BestWPM {
		p.BestWPM = s.WPM
	}
	p.PracticeTime += s.Duration
	p.LessonsDone++

	for r, hits := range s.KeyHits {
		st := p.KeyStats[string(r)]
		st.Presses += hits
		st.Errors += s.KeyErrors[r]
		p.KeyStats[string(r)] = st
	}
}

// WeakKeys returns up to n keys with the highest error rate, most error-prone
// first. Keys with no recorded presses are ignored.
func (p *Profile) WeakKeys(n int) []rune {
	type rate struct {
		r    rune
		rate float64
	}
	var rates []rate
	for k, st := range p.KeyStats {
		if st.Presses == 0 {
			continue
		}
		rates = append(rates, rate{[]rune(k)[0], float64(st.Errors) / float64(st.Presses)})
	}
	sort.Slice(rates, func(i, j int) bool {
		if rates[i].rate != rates[j].rate {
			return rates[i].rate > rates[j].rate
		}
		return rates[i].r < rates[j].r
	})
	if n > len(rates) {
		n = len(rates)
	}
	out := make([]rune, n)
	for i := 0; i < n; i++ {
		out[i] = rates[i].r
	}
	return out
}
