// Package lesson generates practice text. The Generator interface is the seam
// where a future adaptive (Keybr-style) algorithm can be swapped in without
// touching the rest of the application.
package lesson

import "strok/internal/domain"

// Generator produces the next lesson given the user's current progress.
type Generator interface {
	Next(p domain.Profile) domain.Lesson
}
