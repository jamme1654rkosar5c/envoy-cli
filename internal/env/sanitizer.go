package env

import (
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// SanitizeOptions controls sanitization behavior.
type SanitizeOptions struct {
	TrimKeys        bool
	TrimValues      bool
	RemoveEmpty     bool
	NormalizeKeys   bool // uppercase keys
	StripQuotes     bool
}

// DefaultSanitizeOptions returns sensible defaults.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		TrimKeys:      true,
		TrimValues:    true,
		RemoveEmpty:   false,
		NormalizeKeys: false,
		StripQuotes:   false,
	}
}

// Sanitize cleans entries according to the provided options.
// It returns a new slice and does not mutate the input.
func Sanitize(entries []parser.Entry, opts SanitizeOptions) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		key := e.Key
		val := e.Value

		if opts.TrimKeys {
			key = strings.TrimSpace(key)
		}
		if opts.NormalizeKeys {
			key = strings.ToUpper(key)
		}
		if opts.TrimValues {
			val = strings.TrimSpace(val)
		}
		if opts.StripQuotes {
			val = stripQuotes(val)
		}
		if opts.RemoveEmpty && val == "" {
			continue
		}
		result = append(result, parser.Entry{
			Key:     key,
			Value:   val,
			Comment: e.Comment,
		})
	}
	return result
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
