package env

import (
	"fmt"

	"github.com/user/envoy-cli/internal/parser"
)

// CloneOptions controls how entries are cloned.
type CloneOptions struct {
	// Prefix to prepend to all cloned keys (e.g. "STAGING_").
	Prefix string
	// Suffix to append to all cloned keys (e.g. "_BACKUP").
	Suffix string
	// SkipKeys is a set of keys to exclude from the clone.
	SkipKeys map[string]bool
	// OverwriteExisting allows cloned keys to replace existing entries.
	OverwriteExisting bool
}

// DefaultCloneOptions returns a CloneOptions with sensible defaults.
func DefaultCloneOptions() CloneOptions {
	return CloneOptions{
		SkipKeys: map[string]bool{},
	}
}

// Clone duplicates entries from src into dst, applying prefix/suffix
// transformations and respecting skip rules. Returns the merged slice.
func Clone(dst, src []parser.Entry, opts CloneOptions) ([]parser.Entry, error) {
	existing := make(map[string]int, len(dst))
	for i, e := range dst {
		existing[e.Key] = i
	}

	result := make([]parser.Entry, len(dst))
	copy(result, dst)

	for _, e := range src {
		if opts.SkipKeys[e.Key] {
			continue
		}
		newKey := opts.Prefix + e.Key + opts.Suffix
		cloned := parser.Entry{
			Key:     newKey,
			Value:   e.Value,
			Comment: e.Comment,
		}
		if idx, exists := existing[newKey]; exists {
			if !opts.OverwriteExisting {
				return nil, fmt.Errorf("clone: key %q already exists in destination", newKey)
			}
			result[idx] = cloned
		} else {
			existing[newKey] = len(result)
			result = append(result, cloned)
		}
	}
	return result, nil
}
