package env

import (
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// TransformFunc is a function that transforms an entry's value.
type TransformFunc func(key, value string) string

// TransformOptions controls how transformations are applied.
type TransformOptions struct {
	// Keys to include; empty means all keys.
	OnlyKeys []string
	// Skip keys matching these prefixes.
	SkipPrefixes []string
}

// DefaultTransformOptions returns sensible defaults.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{}
}

// Transform applies fn to each entry in file, returning a new slice.
func Transform(entries []parser.Entry, fn TransformFunc, opts TransformOptions) []parser.Entry {
	only := make(map[string]bool, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		only[k] = true
	}

	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if len(only) > 0 && !only[e.Key] {
			result = append(result, e)
			continue
		}
		skipped := false
		for _, pfx := range opts.SkipPrefixes {
			if strings.HasPrefix(e.Key, pfx) {
				skipped = true
				break
			}
		}
		if skipped {
			result = append(result, e)
			continue
		}
		result = append(result, parser.Entry{
			Key:     e.Key,
			Value:   fn(e.Key, e.Value),
			Comment: e.Comment,
		})
	}
	return result
}

// BuiltinUppercase transforms all values to uppercase.
func BuiltinUppercase(_, v string) string { return strings.ToUpper(v) }

// BuiltinLowercase transforms all values to lowercase.
func BuiltinLowercase(_, v string) string { return strings.ToLower(v) }

// BuiltinTrimSpace trims leading/trailing whitespace from values.
func BuiltinTrimSpace(_, v string) string { return strings.TrimSpace(v) }
