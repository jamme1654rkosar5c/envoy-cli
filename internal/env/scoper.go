package env

import "strings"

// ScopeOptions controls how scoping is applied to env entries.
type ScopeOptions struct {
	// Scope is the prefix used to filter or namespace entries (e.g. "APP_").
	Scope string
	// StripPrefix removes the scope prefix from keys after filtering.
	StripPrefix bool
	// AddPrefix prepends the scope prefix to all keys.
	AddPrefix bool
}

// DefaultScopeOptions returns sensible defaults.
func DefaultScopeOptions() ScopeOptions {
	return ScopeOptions{
		Scope:       "",
		StripPrefix: false,
		AddPrefix:   false,
	}
}

// Scope filters entries to those matching the given prefix scope, optionally
// stripping or adding the prefix to the resulting keys.
func Scope(entries []Entry, opts ScopeOptions) []Entry {
	if opts.Scope == "" {
		return entries
	}

	result := make([]Entry, 0, len(entries))

	for _, e := range entries {
		if opts.AddPrefix {
			cloned := e
			cloned.Key = opts.Scope + e.Key
			result = append(result, cloned)
			continue
		}

		if !strings.HasPrefix(e.Key, opts.Scope) {
			continue
		}

		cloned := e
		if opts.StripPrefix {
			cloned.Key = strings.TrimPrefix(e.Key, opts.Scope)
		}
		result = append(result, cloned)
	}

	return result
}
