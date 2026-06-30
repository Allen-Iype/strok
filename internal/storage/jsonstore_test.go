package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveThenLoadRoundTrips(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile.json")
	s := NewJSONStore(path)

	p, _ := s.Load()
	p.BestWPM = 77
	p.UnlockedLevel = 3
	if err := s.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.BestWPM != 77 || got.UnlockedLevel != 3 {
		t.Errorf("round trip mismatch: %+v", got)
	}
}

func TestLoadMissingReturnsFresh(t *testing.T) {
	s := NewJSONStore(filepath.Join(t.TempDir(), "nope.json"))
	p, err := s.Load()
	if err != nil {
		t.Fatalf("Load missing: %v", err)
	}
	if p.Version == 0 || p.KeyStats == nil {
		t.Errorf("missing load not a fresh profile: %+v", p)
	}
}

func TestLoadCorruptBacksUpAndRecovers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := NewJSONStore(path)
	p, err := s.Load()
	if err != nil {
		t.Fatalf("Load corrupt: %v", err)
	}
	if p.Version == 0 {
		t.Error("corrupt load should return fresh profile")
	}
	if _, err := os.Stat(path + ".corrupt"); err != nil {
		t.Error("corrupt file was not backed up")
	}
}
