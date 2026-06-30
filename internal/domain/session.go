package domain

import "time"

// Session is the result of completing one lesson. It is folded into a Profile
// via Profile.Apply, which owns all aggregation rules.
type Session struct {
	WPM       float64
	Accuracy  float64
	Errors    int
	Duration  time.Duration
	Keyset    []rune
	KeyErrors map[rune]int // wrong attempts per expected key
	KeyHits   map[rune]int // total attempts per expected key
}
