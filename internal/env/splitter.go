package env

import "github.com/yourorg/envoy-cli/internal/parser"

// SplitOptions controls how entries are split into buckets.
type SplitOptions struct {
	// MaxBuckets limits the number of buckets. 0 means unlimited.
	MaxBuckets int
	// SkipEmpty drops entries with empty values before splitting.
	SkipEmpty bool
}

// DefaultSplitOptions returns sensible defaults.
func DefaultSplitOptions() SplitOptions {
	return SplitOptions{
		MaxBuckets: 0,
		SkipEmpty:  false,
	}
}

// Split divides entries into n roughly equal buckets.
// Useful for distributing env vars across multiple config files or services.
func Split(entries []parser.Entry, n int, opts SplitOptions) [][]parser.Entry {
	if n <= 0 {
		n = 1
	}
	if opts.MaxBuckets > 0 && n > opts.MaxBuckets {
		n = opts.MaxBuckets
	}

	filtered := entries
	if opts.SkipEmpty {
		filtered = make([]parser.Entry, 0, len(entries))
		for _, e := range entries {
			if e.Value != "" {
				filtered = append(filtered, e)
			}
		}
	}

	if len(filtered) == 0 {
		return [][]parser.Entry{}
	}

	buckets := make([][]parser.Entry, n)
	for i := range buckets {
		buckets[i] = []parser.Entry{}
	}

	for i, e := range filtered {
		idx := i % n
		buckets[idx] = append(buckets[idx], e)
	}

	return buckets
}
