// Package profile manages named environment profiles, allowing users to
// define and switch between sets of .env files for different contexts
// (e.g., "local", "staging", "production").
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Profile represents a named collection of environment file paths.
type Profile struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

// profileStore is the on-disk structure for persisting profiles.
type profileStore struct {
	Profiles []Profile `json:"profiles"`
}

// Save writes a profile to the given directory. If a profile with the same
// name already exists it is replaced.
func Save(dir string, p Profile) error {
	if p.Name == "" {
		return errors.New("profile name must not be empty")
	}
	store, err := load(dir)
	if err != nil {
		return err
	}
	for i, existing := range store.Profiles {
		if existing.Name == p.Name {
			store.Profiles[i] = p
			return persist(dir, store)
		}
	}
	store.Profiles = append(store.Profiles, p)
	return persist(dir, store)
}

// Get returns the profile with the given name from dir.
func Get(dir, name string) (Profile, error) {
	store, err := load(dir)
	if err != nil {
		return Profile{}, err
	}
	for _, p := range store.Profiles {
		if p.Name == name {
			return p, nil
		}
	}
	return Profile{}, fmt.Errorf("profile %q not found", name)
}

// List returns all profiles stored in dir.
func List(dir string) ([]Profile, error) {
	store, err := load(dir)
	if err != nil {
		return nil, err
	}
	return store.Profiles, nil
}

// Delete removes the named profile from dir. It returns an error if the
// profile does not exist.
func Delete(dir, name string) error {
	store, err := load(dir)
	if err != nil {
		return err
	}
	for i, p := range store.Profiles {
		if p.Name == name {
			store.Profiles = append(store.Profiles[:i], store.Profiles[i+1:]...)
			return persist(dir, store)
		}
	}
	return fmt.Errorf("profile %q not found", name)
}

func storePath(dir string) string {
	return filepath.Join(dir, "profiles.json")
}

func load(dir string) (profileStore, error) {
	path := storePath(dir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return profileStore{}, nil
	}
	if err != nil {
		return profileStore{}, fmt.Errorf("reading profile store: %w", err)
	}
	var store profileStore
	if err := json.Unmarshal(data, &store); err != nil {
		return profileStore{}, fmt.Errorf("parsing profile store: %w", err)
	}
	return store, nil
}

func persist(dir string, store profileStore) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating profile dir: %w", err)
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding profile store: %w", err)
	}
	return os.WriteFile(storePath(dir), data, 0o644)
}
