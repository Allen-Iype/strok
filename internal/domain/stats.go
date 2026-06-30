package domain

import "time"

// Stats is an immutable snapshot of typing performance for display. It is
// recomputed each frame from the engine state and the elapsed time.
type Stats struct {
	WPM      float64       // words per minute (correct chars / 5 / minutes)
	Accuracy float64       // 0..1, correct keystrokes / total keystrokes
	Errors   int           // total wrong keystrokes
	Typed    int           // characters advanced through the lesson
	Elapsed  time.Duration // time since the first keystroke
}
