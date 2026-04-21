package env

import (
	"fmt"

	"github.com/envoy-cli/internal/parser"
)

// StashEntry holds a named stash of env entries.
type StashEntry struct {
	Name    string
	Entries []parser.Entry
}

// StashOptions configures Stash and Pop behaviour.
type StashOptions struct {
	// AllowOverwrite replaces an existing stash with the same name.
	AllowOverwrite bool
	// RestoreOnPop merges stashed entries back into dst on Pop.
	RestoreOnPop bool
}

// DefaultStashOptions returns sensible defaults.
func DefaultStashOptions() StashOptions {
	return StashOptions{
		AllowOverwrite: false,
		RestoreOnPop:   true,
	}
}

// Stash saves a named copy of entries into the provided stash map.
// Returns an error if a stash with the same name already exists and
// AllowOverwrite is false.
func Stash(name string, entries []parser.Entry, store map[string]StashEntry, opts StashOptions) error {
	if _, exists := store[name]; exists && !opts.AllowOverwrite {
		return fmt.Errorf("stash %q already exists; use AllowOverwrite to replace it", name)
	}
	copied := make([]parser.Entry, len(entries))
	copy(copied, entries)
	store[name] = StashEntry{Name: name, Entries: copied}
	return nil
}

// Pop retrieves and removes the named stash from the store.
// If RestoreOnPop is true, stashed entries are merged into dst (existing keys
// are not overwritten).
// Returns the stashed entries and an error if the name is not found.
func Pop(name string, dst []parser.Entry, store map[string]StashEntry, opts StashOptions) ([]parser.Entry, error) {
	se, ok := store[name]
	if !ok {
		return nil, fmt.Errorf("stash %q not found", name)
	}
	delete(store, name)

	if !opts.RestoreOnPop {
		return se.Entries, nil
	}

	existing := make(map[string]struct{}, len(dst))
	for _, e := range dst {
		existing[e.Key] = struct{}{}
	}

	result := make([]parser.Entry, len(dst))
	copy(result, dst)
	for _, e := range se.Entries {
		if _, found := existing[e.Key]; !found {
			result = append(result, e)
		}
	}
	return result, nil
}

// ListStashes returns the names of all stashes currently in the store.
func ListStashes(store map[string]StashEntry) []string {
	names := make([]string, 0, len(store))
	for k := range store {
		names = append(names, k)
	}
	return names
}
