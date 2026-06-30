package engine

// Status is the correctness state of a single typed position.
type Status int

const (
	Pending   Status = iota // not yet reached
	Correct                 // typed correctly
	Incorrect               // currently shows a wrong character
)

// Entry is one position in the lesson: the expected rune, what the user most
// recently typed there, and whether any wrong key was ever pressed for it.
type Entry struct {
	Expected rune
	Typed    rune
	Status   Status
	Errored  bool // a wrong key was pressed here at least once (permanent)
}
