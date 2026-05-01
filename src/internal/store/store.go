// Package store provides named configuration entry storage (v1 legacy).
package store

// Store represents the v1 configuration store for named entries.
// This minimal type exists for backward compatibility with session-based operations.
type Store struct {
	rootDir string
}

// NewStore creates a new Store instance for v1 named configuration storage.
func NewStore(rootDir string) *Store {
	return &Store{
		rootDir: rootDir,
	}
}

// GetRootDir returns the root directory path.
func (s *Store) GetRootDir() string {
	return s.rootDir
}