package env

import (
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// TrimOptions controls how trimming is applied to env entries.
type TrimOptions struct {
	TrimKeys        bool
	TrimValues      bool
	TrimPrefixes    []string // remove these prefixes from keys
	TrimSuffixes    []string // remove these suffixes from keys
	SkipEmpty       bool
}

// DefaultTrimOptions returns sensible defaults.
func DefaultTrimOptions() TrimOptions {
	return TrimOptions{
		TrimKeys:   true,
		TrimValues: true,
		SkipEmpty:  false,
	}
}

// Trim applies trimming rules to a slice of env entries and returns a new slice.
func Trim(entries []parser.Entry, opts TrimOptions) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if opts.SkipEmpty && e.Value == "" {
			result = append(result, e)
			continue
		}

		if opts.TrimKeys {
			e.Key = strings.TrimSpace(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}

		for _, p := range opts.TrimPrefixes {
			e.Key = strings.TrimPrefix(e.Key, p)
		}
		for _, s := range opts.TrimSuffixes {
			e.Key = strings.TrimSuffix(e.Key, s)
		}

		result = append(result, e)
	}
	return result
}
