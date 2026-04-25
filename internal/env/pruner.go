package env

import "github.com/your-org/envoy-cli/internal/parser"

// PruneOptions controls how entries are pruned from a set.
type PruneOptions struct {
	// RemoveEmpty removes entries with empty values.
	RemoveEmpty bool
	// RemoveCommented removes entries whose value starts with a comment marker.
	RemoveCommented bool
	// Keys is an explicit list of keys to remove. If non-empty, only these keys
	// are pruned (other options are still applied in addition).
	Keys []string
	// DryRun returns a copy without mutating the original slice.
	DryRun bool
}

// DefaultPruneOptions returns a sensible default configuration.
func DefaultPruneOptions() PruneOptions {
	return PruneOptions{
		RemoveEmpty:     true,
		RemoveCommented: false,
		DryRun:          false,
	}
}

// Prune removes entries from src according to opts and returns the surviving
// entries together with the list of keys that were pruned.
func Prune(src []parser.Entry, opts PruneOptions) ([]parser.Entry, []string) {
	explicit := toExplicitSet(opts.Keys)

	var kept []parser.Entry
	var pruned []string

	for _, e := range src {
		if shouldPrune(e, opts, explicit) {
			pruned = append(pruned, e.Key)
			continue
		}
		if opts.DryRun {
			copy := e
			kept = append(kept, copy)
		} else {
			kept = append(kept, e)
		}
	}

	return kept, pruned
}

func shouldPrune(e parser.Entry, opts PruneOptions, explicit map[string]struct{}) bool {
	if len(explicit) > 0 {
		if _, ok := explicit[e.Key]; ok {
			return true
		}
	}
	if opts.RemoveEmpty && e.Value == "" {
		return true
	}
	if opts.RemoveCommented && len(e.Value) > 0 && e.Value[0] == '#' {
		return true
	}
	return false
}

func toExplicitSet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}
