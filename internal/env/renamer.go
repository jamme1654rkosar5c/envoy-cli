package env

import (
	"fmt"

	"github.com/envoy-cli/internal/parser"
)

// RenameOptions controls behaviour of the Rename function.
type RenameOptions struct {
	// FailIfNotFound returns an error if the old key does not exist.
	FailIfNotFound bool
	// FailIfDestExists returns an error if the new key already exists.
	FailIfDestExists bool
	// DryRun returns a copy without mutating the original slice.
	DryRun bool
}

// DefaultRenameOptions returns sensible defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		FailIfNotFound:  true,
		FailIfDestExists: true,
		DryRun:          false,
	}
}

// Rename renames oldKey to newKey within entries.
// It preserves the original comment and order of the entry.
func Rename(entries []parser.EnvEntry, oldKey, newKey string, opts RenameOptions) ([]parser.EnvEntry, error) {
	result := make([]parser.EnvEntry, len(entries))
	copy(result, entries)

	oldIdx := -1
	newIdx := -1
	for i, e := range result {
		if e.Key == oldKey {
			oldIdx = i
		}
		if e.Key == newKey {
			newIdx = i
		}
	}

	if oldIdx == -1 {
		if opts.FailIfNotFound {
			return nil, fmt.Errorf("key %q not found", oldKey)
		}
		return result, nil
	}

	if newIdx != -1 {
		if opts.FailIfDestExists {
			return nil, fmt.Errorf("destination key %q already exists", newKey)
		}
	}

	if opts.DryRun {
		return result, nil
	}

	result[oldIdx].Key = newKey
	return result, nil
}
