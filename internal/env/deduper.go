package env

import "github.com/envoy-cli/internal/parser"

// DedupeStrategy controls how duplicate keys are resolved.
type DedupeStrategy int

const (
	// KeepFirst retains the first occurrence of a duplicate key.
	KeepFirst DedupeStrategy = iota
	// KeepLast retains the last occurrence of a duplicate key.
	KeepLast
)

// DedupeOptions configures deduplication behaviour.
type DedupeOptions struct {
	Strategy DedupeStrategy
}

// DefaultDedupeOptions returns sensible defaults.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{Strategy: KeepFirst}
}

// Dedupe removes duplicate keys from entries according to the given options.
// It returns the deduplicated slice and a list of keys that were removed.
func Dedupe(entries []parser.Entry, opts DedupeOptions) ([]parser.Entry, []string) {
	seen := make(map[string]int) // key -> index in result
	result := make([]parser.Entry, 0, len(entries))
	var removed []string

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			switch opts.Strategy {
			case KeepLast:
				removed = append(removed, result[idx].Key)
				result[idx] = e
			default: // KeepFirst
				removed = append(removed, e.Key)
			}
		} else {
			seen[e.Key] = len(result)
			result = append(result, e)
		}
	}

	return result, removed
}
