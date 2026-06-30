package domain

// Lesson is a single unit of practice: the target text the user must type and
// the set of letters this lesson focuses on (used for the header display).
type Lesson struct {
	Text   string
	Keyset []rune
}

// Len returns the number of characters in the lesson text.
func (l Lesson) Len() int { return len([]rune(l.Text)) }
