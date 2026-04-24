package env

import (
	"fmt"
	"strings"
)

// DiffOptions controls the behaviour of the Diff function.
type DiffOptions struct {
	// MaskSecrets replaces sensitive values with "***" in the output.
	MaskSecrets bool
	// IgnoreKeys is a set of keys to exclude from the diff.
	IgnoreKeys []string
}

// DefaultDiffOptions returns sensible defaults for DiffOptions.
func DefaultDiffOptions() DiffOptions {
	return DiffOptions{
		MaskSecrets: false,
		IgnoreKeys:  nil,
	}
}

// DiffEntry describes a single key-level change between two entry slices.
type DiffEntry struct {
	Key      string
	Status   string // "added", "removed", "changed", "unchanged"
	OldValue string
	NewValue string
}

// Diff compares two slices of EnvEntry and returns a list of DiffEntry values
// describing what changed between base and next.
func Diff(base, next []EnvEntry, opts DiffOptions) ([]DiffEntry, error) {
	ignore := make(map[string]struct{}, len(opts.IgnoreKeys))
	for _, k := range opts.IgnoreKeys {
		ignore[k] = struct{}{}
	}

	baseMap := make(map[string]string, len(base))
	for _, e := range base {
		if _, skip := ignore[e.Key]; skip {
			continue
		}
		baseMap[e.Key] = e.Value
	}

	nextMap := make(map[string]string, len(next))
	for _, e := range next {
		if _, skip := ignore[e.Key]; skip {
			continue
		}
		nextMap[e.Key] = e.Value
	}

	var results []DiffEntry

	for k, oldVal := range baseMap {
		newVal, exists := nextMap[k]
		if !exists {
			results = append(results, DiffEntry{
				Key:      k,
				Status:   "removed",
				OldValue: maybeMask(k, oldVal, opts.MaskSecrets),
				NewValue: "",
			})
		} else if oldVal != newVal {
			results = append(results, DiffEntry{
				Key:      k,
				Status:   "changed",
				OldValue: maybeMask(k, oldVal, opts.MaskSecrets),
				NewValue: maybeMask(k, newVal, opts.MaskSecrets),
			})
		} else {
			results = append(results, DiffEntry{
				Key:      k,
				Status:   "unchanged",
				OldValue: maybeMask(k, oldVal, opts.MaskSecrets),
				NewValue: maybeMask(k, newVal, opts.MaskSecrets),
			})
		}
	}

	for k, newVal := range nextMap {
		if _, existed := baseMap[k]; !existed {
			results = append(results, DiffEntry{
				Key:      k,
				Status:   "added",
				OldValue: "",
				NewValue: maybeMask(k, newVal, opts.MaskSecrets),
			})
		}
	}

	if len(results) == 0 {
		return results, nil
	}

	_ = fmt.Sprintf // keep import
	return results, nil
}

func maybeMask(key, value string, mask bool) string {
	if !mask {
		return value
	}
	lower := strings.ToLower(key)
	for _, kw := range []string{"secret", "password", "passwd", "token", "api_key", "apikey", "private"} {
		if strings.Contains(lower, kw) {
			return "***"
		}
	}
	return value
}
