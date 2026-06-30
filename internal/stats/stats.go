// Package stats computes display statistics from typing counters and elapsed
// time. It is pure: callers pass in the elapsed duration rather than the package
// reading a clock, which keeps the WPM/accuracy math directly testable.
package stats

import (
	"time"

	"strok/internal/domain"
)

// Source is the minimal view of typing state the stats need. The engine's
// TypingState satisfies it.
type Source interface {
	Keystrokes() int
	Errors() int
	Cursor() int
}

// Compute builds a Stats snapshot from a typing source and elapsed time.
//
// WPM uses the standard definition: (correct characters / 5) / minutes, where
// correct characters is the cursor position (how far the user has advanced).
// Accuracy is correct keystrokes / total keystrokes. Both guard against the
// zero-input case to avoid NaN/Inf.
func Compute(s Source, elapsed time.Duration) domain.Stats {
	typed := s.Cursor()
	keystrokes := s.Keystrokes()
	errors := s.Errors()

	accuracy := 1.0
	if keystrokes > 0 {
		accuracy = float64(keystrokes-errors) / float64(keystrokes)
	}

	wpm := 0.0
	if minutes := elapsed.Minutes(); minutes > 0 {
		wpm = (float64(typed) / 5.0) / minutes
	}

	return domain.Stats{
		WPM:      wpm,
		Accuracy: accuracy,
		Errors:   errors,
		Typed:    typed,
		Elapsed:  elapsed,
	}
}
