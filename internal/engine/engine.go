// Package engine is the typing state machine. It is pure: it holds no rendering
// or I/O concerns and never reads the clock itself, so it can be unit-tested
// headlessly. The UI drives it with HandleKey/Backspace and reads its state.
package engine

import (
	"time"

	"strok/internal/domain"
)

// Feedback describes the result of the most recent keystroke, used by the UI to
// briefly flash keys. Zero value means "no recent feedback".
type Feedback struct {
	Expected rune // the key that was expected
	Pressed  rune // the key the user actually pressed
	Correct  bool // whether the press matched
}

// TypingState is the live state of typing one lesson.
type TypingState struct {
	entries []Entry
	cursor  int

	// Aggregate keystroke counters (a keystroke = one HandleKey call that
	// advances or attempts a position). Backspace does not count as a keystroke.
	keystrokes int
	errors     int

	// Per-expected-key tallies for weak-key tracking.
	hits   map[rune]int
	keyErr map[rune]int

	keyset   []rune
	feedback Feedback
	started  bool
}

// New builds a fresh TypingState for the given lesson.
func New(l domain.Lesson) *TypingState {
	runes := []rune(l.Text)
	entries := make([]Entry, len(runes))
	for i, r := range runes {
		entries[i] = Entry{Expected: r, Status: Pending}
	}
	return &TypingState{
		entries: entries,
		hits:    map[rune]int{},
		keyErr:  map[rune]int{},
		keyset:  l.Keyset,
	}
}

// HandleKey processes a typed rune. It records correctness at the cursor,
// advances on a correct key, and (Keybr-style) permanently counts any wrong key
// as an error even though the user may later backspace and retype.
func (s *TypingState) HandleKey(r rune) {
	if s.Done() {
		return
	}
	s.started = true
	s.keystrokes++

	e := &s.entries[s.cursor]
	exp := e.Expected
	s.hits[exp]++

	if r == exp {
		e.Typed = r
		e.Status = Correct
		s.cursor++
		s.feedback = Feedback{Expected: exp, Pressed: r, Correct: true}
		return
	}

	// Wrong key: record a permanent error but do not advance.
	e.Typed = r
	e.Status = Incorrect
	e.Errored = true
	s.errors++
	s.keyErr[exp]++
	s.feedback = Feedback{Expected: exp, Pressed: r, Correct: false}
}

// Backspace moves the cursor back one position and clears the visible character
// there. It does not erase the permanent error record (errors stay counted).
func (s *TypingState) Backspace() {
	if s.cursor == 0 {
		// Clear an in-place incorrect mark at the very first position.
		if s.entries[0].Status == Incorrect {
			s.entries[0].Status = Pending
			s.entries[0].Typed = 0
		}
		return
	}
	// If the current position shows an incorrect char, clear it in place first.
	if s.cursor < len(s.entries) && s.entries[s.cursor].Status == Incorrect {
		s.entries[s.cursor].Status = Pending
		s.entries[s.cursor].Typed = 0
		return
	}
	s.cursor--
	s.entries[s.cursor].Status = Pending
	s.entries[s.cursor].Typed = 0
}

// Done reports whether every position has been typed correctly.
func (s *TypingState) Done() bool { return s.cursor >= len(s.entries) }

// Entries exposes the per-position entries for rendering.
func (s *TypingState) Entries() []Entry { return s.entries }

// Cursor returns the current position index.
func (s *TypingState) Cursor() int { return s.cursor }

// Started reports whether any key has been typed (used to start the timer).
func (s *TypingState) Started() bool { return s.started }

// Keyset returns the lesson's focus letters.
func (s *TypingState) Keyset() []rune { return s.keyset }

// Feedback returns the most recent keystroke feedback.
func (s *TypingState) Feedback() Feedback { return s.feedback }

// Keystrokes returns the total number of typed keys (excludes backspace).
func (s *TypingState) Keystrokes() int { return s.keystrokes }

// Errors returns the total number of wrong keystrokes.
func (s *TypingState) Errors() int { return s.errors }

// Expected returns the rune the user is currently expected to type, or 0 if the
// lesson is complete.
func (s *TypingState) Expected() rune {
	if s.Done() {
		return 0
	}
	return s.entries[s.cursor].Expected
}

// Session builds the completed-lesson result. The caller supplies the snapshot
// stats (WPM, accuracy, duration) since the engine has no clock; the engine
// contributes the per-key tallies it has accumulated.
func (s *TypingState) Session(stats domain.Stats, dur time.Duration) domain.Session {
	return domain.Session{
		WPM:       stats.WPM,
		Accuracy:  stats.Accuracy,
		Errors:    s.errors,
		Duration:  dur,
		Keyset:    s.keyset,
		KeyErrors: s.keyErr,
		KeyHits:   s.hits,
	}
}
