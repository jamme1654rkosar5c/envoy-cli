package env

import (
	"fmt"

	"github.com/envoy-cli/internal/parser"
)

// PinOptions controls how keys are pinned or unpinned.
type PinOptions struct {
	// FailIfMissing returns an error if a key to pin is not found.
	FailIfMissing bool
}

// DefaultPinOptions returns sensible defaults.
func DefaultPinOptions() PinOptions {
	return PinOptions{
		FailIfMissing: true,
	}
}

// Pin marks the given keys as pinned by appending a "# pinned" comment.
// Pinned keys are preserved during merge and patch operations.
func Pin(entries []parser.Entry, keys []string, opts PinOptions) ([]parser.Entry, error) {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	matched := make(map[string]bool)
	result := make([]parser.Entry, len(entries))
	copy(result, entries)

	for i, e := range result {
		if keySet[e.Key] {
			if e.Comment != "pinned" {
				result[i].Comment = "pinned"
			}
			matched[e.Key] = true
		}
	}

	if opts.FailIfMissing {
		for _, k := range keys {
			if !matched[k] {
				return nil, fmt.Errorf("pin: key %q not found", k)
			}
		}
	}

	return result, nil
}

// Unpin removes the "# pinned" comment from the given keys.
func Unpin(entries []parser.Entry, keys []string, opts PinOptions) ([]parser.Entry, error) {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	matched := make(map[string]bool)
	result := make([]parser.Entry, len(entries))
	copy(result, entries)

	for i, e := range result {
		if keySet[e.Key] {
			if e.Comment == "pinned" {
				result[i].Comment = ""
			}
			matched[e.Key] = true
		}
	}

	if opts.FailIfMissing {
		for _, k := range keys {
			if !matched[k] {
				return nil, fmt.Errorf("unpin: key %q not found", k)
			}
		}
	}

	return result, nil
}

// IsPinned reports whether an entry is pinned.
func IsPinned(e parser.Entry) bool {
	return e.Comment == "pinned"
}
