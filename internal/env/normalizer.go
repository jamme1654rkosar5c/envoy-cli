package env

import (
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// NormalizeOptions controls how normalization is applied.
type NormalizeOptions struct {
	// UppercaseKeys converts all key names to uppercase.
	UppercaseKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// StripQuotes removes surrounding single or double quotes from values.
	StripQuotes bool
	// RemoveEmpty drops entries with empty values.
	RemoveEmpty bool
	// CollapseWhitespace replaces runs of whitespace in values with a single space.
	CollapseWhitespace bool
}

// DefaultNormalizeOptions returns sensible defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		UppercaseKeys:      true,
		TrimValues:         true,
		StripQuotes:        false,
		RemoveEmpty:        false,
		CollapseWhitespace: false,
	}
}

// Normalize applies normalization rules to a slice of EnvEntry values and
// returns a new slice. The original slice is never mutated.
func Normalize(entries []parser.EnvEntry, opts NormalizeOptions) []parser.EnvEntry {
	result := make([]parser.EnvEntry, 0, len(entries))
	for _, e := range entries {
		key := e.Key
		val := e.Value

		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}
		if opts.TrimValues {
			val = strings.TrimSpace(val)
		}
		if opts.StripQuotes {
			val = stripNormQuotes(val)
		}
		if opts.CollapseWhitespace {
			val = strings.Join(strings.Fields(val), " ")
		}
		if opts.RemoveEmpty && val == "" {
			continue
		}

		result = append(result, parser.EnvEntry{
			Key:     key,
			Value:   val,
			Comment: e.Comment,
		})
	}
	return result
}

func stripNormQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
