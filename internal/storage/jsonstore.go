package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"strok/internal/domain"
)

// JSONStore persists the profile as a JSON file on disk.
type JSONStore struct {
	path string
}

// NewJSONStore returns a store writing to the given file path.
func NewJSONStore(path string) *JSONStore { return &JSONStore{path: path} }

// DefaultPath returns the standard per-user profile location, e.g.
// ~/.config/strok/profile.json on Linux/macOS or %AppData%\strok on Windows.
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locate config dir: %w", err)
	}
	return filepath.Join(dir, "strok", "profile.json"), nil
}

// Load reads the profile. A missing file returns a fresh profile. A corrupt
// file is backed up and a fresh profile is returned so the user is never blocked
// by bad data.
func (s *JSONStore) Load() (domain.Profile, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return domain.NewProfile(), nil
	}
	if err != nil {
		return domain.NewProfile(), fmt.Errorf("read profile: %w", err)
	}

	var p domain.Profile
	if err := json.Unmarshal(data, &p); err != nil {
		// Preserve the corrupt file and start fresh rather than crash.
		_ = os.Rename(s.path, s.path+".corrupt")
		return domain.NewProfile(), nil
	}
	if p.KeyStats == nil {
		p.KeyStats = map[string]domain.KeyStat{}
	}
	return p, nil
}

// Save writes the profile atomically: it writes to a temp file in the same
// directory and renames it into place so a crash mid-write cannot corrupt the
// existing profile.
func (s *JSONStore) Save(p domain.Profile) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create profile dir: %w", err)
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profile: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(s.path), "profile-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op if the rename succeeded

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpName, s.path); err != nil {
		return fmt.Errorf("replace profile: %w", err)
	}
	return nil
}
