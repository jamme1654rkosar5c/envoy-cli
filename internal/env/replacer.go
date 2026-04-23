package env

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// ReplaceOptions controls the behaviour of Replace.
type ReplaceOptions struct {
	// OldValue is the substring (or exact value) to search for.
	OldValue string

	// NewValue is the replacement string.
	NewValue string

	// KeyFilter, when non-empty, restricts replacements to keys that contain
	// this substring.
	KeyFilter string

	// ExactMatch requires the entire value to equal OldValue rather than just
	// containing it as a substring.
	ExactMatch bool

	// DryRun returns a mutated copy without modifying the original slice.
	DryRun bool
}

// DefaultReplaceOptions returns sensible defaults.
func DefaultReplaceOptions() ReplaceOptions {
	return ReplaceOptions{}
}

// Replace scans entries and substitutes occurrences of OldValue with NewValue
// in entry values. It returns the (possibly cloned) entries and a count of
// how many substitutions were made.
func Replace(entries []parser.Entry, opts ReplaceOptions) ([]parser.Entry, int, error) {
	if opts.OldValue == "" {
		return nil, 0, fmt.Errorf("replace: OldValue must not be empty")
	}

	working := entries
	if opts.DryRun {
		working = make([]parser.Entry, len(entries))
		copy(working, entries)
	}

	count := 0
	for i, e := range working {
		if opts.KeyFilter != "" && !strings.Contains(e.Key, opts.KeyFilter) {
			continue
		}

		var newVal string
		if opts.ExactMatch {
			if e.Value != opts.OldValue {
				continue
			}
			newVal = opts.NewValue
		} else {
			if !strings.Contains(e.Value, opts.OldValue) {
				continue
			}
			newVal = strings.ReplaceAll(e.Value, opts.OldValue, opts.NewValue)
		}

		working[i].Value = newVal
		count++
	}

	return working, count, nil
}
