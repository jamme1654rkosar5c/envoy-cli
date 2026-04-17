package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

const freezeTag = "frozen"

// FreezeOptions controls Freeze/Unfreeze behaviour.
type FreezeOptions struct {
	AllowMissing bool
}

// DefaultFreezeOptions returns sensible defaults.
func DefaultFreezeOptions() FreezeOptions {
	return FreezeOptions{AllowMissing: false}
}

// Freeze marks a key as frozen by appending a # frozen comment.
// Frozen keys are protected from modification by Patch.
func Freeze(entries []parser.Entry, key string, opts FreezeOptions) ([]parser.Entry, error) {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if e.Key == key {
			out[i].Comment = setFreezeComment(e.Comment, true)
			return out, nil
		}
	}
	if !opts.AllowMissing {
		return nil, fmt.Errorf("freeze: key %q not found", key)
	}
	return out, nil
}

// Unfreeze removes the frozen marker from a key's comment.
func Unfreeze(entries []parser.Entry, key string, opts FreezeOptions) ([]parser.Entry, error) {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if e.Key == key {
			out[i].Comment = setFreezeComment(e.Comment, false)
			return out, nil
		}
	}
	if !opts.AllowMissing {
		return nil, fmt.Errorf("unfreeze: key %q not found", key)
	}
	return out, nil
}

// IsFrozen reports whether a key is frozen.
func IsFrozen(entries []parser.Entry, key string) bool {
	for _, e := range entries {
		if e.Key == key {
			return strings.Contains(e.Comment, freezeTag)
		}
	}
	return false
}

func setFreezeComment(comment string, add bool) string {
	clean := strings.TrimSpace(strings.ReplaceAll(comment, freezeTag, ""))
	clean = strings.Trim(clean, "# ")
	if add {
		if clean == "" {
			return freezeTag
		}
		return clean + " " + freezeTag
	}
	return clean
}
