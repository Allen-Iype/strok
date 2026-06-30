// Package storage persists the user profile. The Store interface isolates the
// rest of the app from the storage mechanism so a future backend (SQLite,
// cloud) can replace the JSON file without other changes.
package storage

import "strok/internal/domain"

// Store loads and saves the user profile.
type Store interface {
	// Load returns the stored profile. A missing profile yields a fresh one and
	// a nil error. A corrupt profile yields a fresh one and a nil error after
	// the corrupt file is preserved as a backup.
	Load() (domain.Profile, error)
	// Save persists the profile.
	Save(domain.Profile) error
}
