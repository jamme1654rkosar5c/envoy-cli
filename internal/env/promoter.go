package env

import (
	"fmt"

	"github.com/user/envoy-cli/internal/parser"
)

// PromoteOptions controls the behaviour of the Promote function.
type PromoteOptions struct {
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
	// DryRun returns the result without mutating dst.
	DryRun bool
	// Keys restricts promotion to the listed keys. Empty means all keys.
	Keys []string
}

// DefaultPromoteOptions returns sensible defaults.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		Overwrite: false,
		DryRun:    false,
		Keys:      nil,
	}
}

// Promote copies selected entries from src into dst, simulating a
// promotion workflow (e.g. staging → production).
// It returns the updated dst entries.
func Promote(src, dst []parser.Entry, opts PromoteOptions) ([]parser.Entry, error) {
	allow := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allow[k] = true
	}

	dstMap := make(map[string]int, len(dst))
	result := make([]parser.Entry, len(dst))
	copy(result, dst)
	for i, e := range result {
		dstMap[e.Key] = i
	}

	for _, e := range src {
		if len(allow) > 0 && !allow[e.Key] {
			continue
		}
		if idx, exists := dstMap[e.Key]; exists {
			if !opts.Overwrite {
				return nil, fmt.Errorf("promote: key %q already exists in destination (use Overwrite)", e.Key)
			}
			if !opts.DryRun {
				result[idx].Value = e.Value
			}
		} else {
			if !opts.DryRun {
				result = append(result, e)
				dstMap[e.Key] = len(result) - 1
			}
		}
	}
	return result, nil
}
