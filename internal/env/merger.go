package env

import "github.com/your-org/envoy-cli/internal/parser"

// OverlapOptions controls how key overlaps are handled during a merge.
type OverlapOptions struct {
	Overwrite bool // if true, source values overwrite destination
	SkipEmpty bool // if true, empty source values are skipped
}

// DefaultOverlapOptions returns sensible defaults.
func DefaultOverlapOptions() OverlapOptions {
	return OverlapOptions{
		Overwrite: false,
		SkipEmpty: true,
	}
}

// Overlap merges entries from src into dst according to opts.
// It returns a new slice without mutating the inputs.
func Overlap(dst, src []parser.Entry, opts OverlapOptions) []parser.Entry {
	result := make([]parser.Entry, len(dst))
	copy(result, dst)

	index := make(map[string]int, len(result))
	for i, e := range result {
		index[e.Key] = i
	}

	for _, e := range src {
		if opts.SkipEmpty && e.Value == "" {
			continue
		}
		if i, exists := index[e.Key]; exists {
			if opts.Overwrite {
				result[i].Value = e.Value
			}
		} else {
			result = append(result, e)
			index[e.Key] = len(result) - 1
		}
	}

	return result
}
