package env

import (
	"fmt"

	"github.com/user/envoy-cli/internal/parser"
)

// RotateOptions controls key rotation behaviour.
type RotateOptions struct {
	// DryRun prevents mutations when true.
	DryRun bool
	// ErrorOnMissing returns an error if a key to rotate is absent.
	ErrorOnMissing bool
}

// DefaultRotateOptions returns sensible defaults.
func DefaultRotateOptions() RotateOptions {
	return RotateOptions{
		DryRun:         false,
		ErrorOnMissing: true,
	}
}

// RotateEntry describes a single key rotation: OldKey -> NewKey with NewValue.
type RotateEntry struct {
	OldKey   string
	NewKey   string
	NewValue string
}

// Rotate replaces keys (and optionally values) in entries according to the
// supplied rotation map. The original key is removed and the new key is
// appended unless DryRun is set.
func Rotate(entries []parser.Entry, rotations []RotateEntry, opts RotateOptions) ([]parser.Entry, error) {
	if opts.DryRun {
		entries = cloneEntries(entries)
	}

	for _, r := range rotations {
		idx := indexOfEntry(entries, r.OldKey)
		if idx == -1 {
			if opts.ErrorOnMissing {
				return nil, fmt.Errorf("rotate: key %q not found", r.OldKey)
			}
			continue
		}

		newValue := r.NewValue
		if newValue == "" {
			newValue = entries[idx].Value
		}

		newEntry := parser.Entry{
			Key:     r.NewKey,
			Value:   newValue,
			Comment: entries[idx].Comment,
		}

		// Replace in-place to preserve ordering.
		entries[idx] = newEntry
	}

	return entries, nil
}

func indexOfEntry(entries []parser.Entry, key string) int {
	for i, e := range entries {
		if e.Key == key {
			return i
		}
	}
	return -1
}

func cloneEntries(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)
	return out
}
