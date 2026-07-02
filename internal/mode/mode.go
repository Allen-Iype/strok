// Package mode defines the progression policy for a play mode: what happens when
// a lesson is completed. Different modes (gated progression now; a free practice
// mode or a timed test later) plug in through the Mode interface, keeping
// mode-specific rules out of the UI. Completion detection currently stays
// "type to end" (engine.Done); a future timed mode can extend this interface
// with a completion trigger without disturbing existing modes.
package mode

import "strok/internal/domain"

// Outcome is a mode's decision after a completed lesson, surfaced by the UI.
type Outcome struct {
	Advanced bool   // did the learner unlock new content?
	Message  string // short note to display after the lesson
}

// Mode owns the progression policy. OnComplete is the single place where
// mode-specific progression lives.
type Mode interface {
	// OnComplete applies the mode's rule to profile (e.g. advancing the unlock
	// level) and reports what happened.
	OnComplete(profile *domain.Profile, s domain.Session) Outcome
	// Name identifies the mode (for the header and future mode switching).
	Name() string
}
