// Package interpolator provides variable interpolation for .env files.
// It resolves references like ${VAR} or $VAR within env entry values.
package interpolator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// AllowMissing silences errors for unresolved references, leaving them as-is.
	AllowMissing bool
}

// DefaultOptions returns sensible default interpolation options.
func DefaultOptions() Options {
	return Options{
		AllowMissing: false,
	}
}

// Interpolate resolves variable references in all entry values within the
// provided EnvFile. It returns a new EnvFile with substituted values and
// does not mutate the original.
func Interpolate(file parser.EnvFile, opts Options) (parser.EnvFile, error) {
	lookup := buildLookup(file)

	resolved := make([]parser.Entry, 0, len(file.Entries))
	for _, entry := range file.Entries {
		val, err := resolve(entry.Value, lookup, opts)
		if err != nil {
			return parser.EnvFile{}, fmt.Errorf("interpolator: key %q: %w", entry.Key, err)
		}
		resolved = append(resolved, parser.Entry{
			Key:     entry.Key,
			Value:   val,
			Comment: entry.Comment,
		})
	}

	return parser.EnvFile{Path: file.Path, Entries: resolved}, nil
}

// buildLookup creates a key→value map from the env file entries.
func buildLookup(file parser.EnvFile) map[string]string {
	m := make(map[string]string, len(file.Entries))
	for _, e := range file.Entries {
		m[e.Key] = e.Value
	}
	return m
}

// resolve performs the actual substitution on a single value string.
func resolve(value string, lookup map[string]string, opts Options) (string, error) {
	var resolveErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		key := extractKey(match)
		if v, ok := lookup[key]; ok {
			return v
		}
		if opts.AllowMissing {
			return match
		}
		resolveErr = fmt.Errorf("unresolved variable reference: %q", key)
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

// extractKey strips the ${ } or $ sigil from a matched token.
func extractKey(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}
