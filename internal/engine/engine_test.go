package engine

import (
	"testing"

	"strok/internal/domain"
)

func lesson(text string) domain.Lesson { return domain.Lesson{Text: text} }

func TestCorrectTypingCompletes(t *testing.T) {
	s := New(lesson("fj"))
	s.HandleKey('f')
	s.HandleKey('j')
	if !s.Done() {
		t.Fatal("expected Done after typing all chars")
	}
	if s.Errors() != 0 {
		t.Errorf("Errors = %d, want 0", s.Errors())
	}
}

func TestWrongKeyCountsErrorAndBlocks(t *testing.T) {
	s := New(lesson("fj"))
	s.HandleKey('d') // wrong, expected f
	if s.Cursor() != 0 {
		t.Errorf("cursor advanced on wrong key: %d", s.Cursor())
	}
	if s.Errors() != 1 {
		t.Errorf("Errors = %d, want 1", s.Errors())
	}
	s.HandleKey('f') // correct now
	if s.Cursor() != 1 {
		t.Errorf("cursor = %d, want 1", s.Cursor())
	}
}

func TestBackspaceClearsButKeepsError(t *testing.T) {
	s := New(lesson("fj"))
	s.HandleKey('f') // correct, cursor->1
	s.HandleKey('x') // wrong at pos1
	if s.Errors() != 1 {
		t.Fatalf("Errors = %d, want 1", s.Errors())
	}
	s.Backspace() // clears the incorrect mark at pos1
	if s.entries[1].Status != Pending {
		t.Errorf("pos1 status = %v, want Pending", s.entries[1].Status)
	}
	if s.Errors() != 1 {
		t.Errorf("Errors after backspace = %d, want 1 (permanent)", s.Errors())
	}
}

func TestBackspaceClearsFeedback(t *testing.T) {
	s := New(lesson("fj"))
	s.HandleKey('x') // wrong, feedback set
	if s.Feedback().Pressed == 0 {
		t.Fatal("feedback should be set after a keystroke")
	}
	s.Backspace()
	if s.Feedback().Pressed != 0 {
		t.Error("feedback should clear on backspace (error acknowledged)")
	}
}

func TestKeyTallies(t *testing.T) {
	s := New(lesson("ff"))
	s.HandleKey('d') // wrong on first f
	s.HandleKey('f') // correct
	s.HandleKey('f') // correct
	if s.hits['f'] != 3 {
		t.Errorf("hits[f] = %d, want 3", s.hits['f'])
	}
	if s.keyErr['f'] != 1 {
		t.Errorf("keyErr[f] = %d, want 1", s.keyErr['f'])
	}
}
