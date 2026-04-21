package env

import (
	"fmt"

	"github.com/user/envoy-cli/internal/parser"
)

// ReorderOptions controls how entries are reordered.
type ReorderOptions struct {
	// Keys defines the desired key order. Keys not listed are appended at the end.
	Keys []string
	// PushUnknownToEnd places unlisted keys after listed ones (default true).
	PushUnknownToEnd bool
	// ErrorOnMissing returns an error if a key in Keys is not found in entries.
	ErrorOnMissing bool
}

// DefaultReorderOptions returns sensible defaults.
func DefaultReorderOptions() ReorderOptions {
	return ReorderOptions{
		PushUnknownToEnd: true,
		ErrorOnMissing:   false,
	}
}

// Reorder returns a new slice of entries arranged according to opts.Keys.
// Entries not mentioned in opts.Keys are appended at the end when
// PushUnknownToEnd is true, or prepended when false.
func Reorder(entries []parser.Entry, opts ReorderOptions) ([]parser.Entry, error) {
	index := make(map[string]parser.Entry, len(entries))
	for _, e := range entries {
		index[e.Key] = e
	}

	if opts.ErrorOnMissing {
		for _, k := range opts.Keys {
			if _, ok := index[k]; !ok {
				return nil, fmt.Errorf("reorder: key %q not found in entries", k)
			}
		}
	}

	seen := make(map[string]bool, len(opts.Keys))
	result := make([]parser.Entry, 0, len(entries))

	for _, k := range opts.Keys {
		if e, ok := index[k]; ok {
			result = append(result, e)
			seen[k] = true
		}
	}

	var rest []parser.Entry
	for _, e := range entries {
		if !seen[e.Key] {
			rest = append(rest, e)
		}
	}

	if opts.PushUnknownToEnd {
		result = append(result, rest...)
	} else {
		result = append(rest, result...)
	}

	return result, nil
}
