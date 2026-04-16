package env

import (
	"fmt"

	"github.com/envoy-cli/internal/parser"
)

// PatchOp represents a single patch operation.
type PatchOp struct {
	Key    string
	Value  string
	Delete bool
}

// PatchOptions controls patching behaviour.
type PatchOptions struct {
	// ErrorOnMissing causes Set operations to fail if the key does not exist.
	ErrorOnMissing bool
	// AllowNew permits inserting keys that are not already present.
	AllowNew bool
}

// DefaultPatchOptions returns sensible defaults.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{
		ErrorOnMissing: false,
		AllowNew:       true,
	}
}

// Patch applies a slice of PatchOps to a copy of entries and returns the result.
func Patch(entries []parser.Entry, ops []PatchOp, opts PatchOptions) ([]parser.Entry, error) {
	// Work on a shallow copy so we do not mutate the caller's slice.
	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	for _, op := range ops {
		if op.Delete {
			out = deleteKey(out, op.Key)
			continue
		}

		idx := indexOfKey(out, op.Key)
		if idx == -1 {
			if opts.ErrorOnMissing {
				return nil, fmt.Errorf("patch: key %q not found", op.Key)
			}
			if opts.AllowNew {
				out = append(out, parser.Entry{Key: op.Key, Value: op.Value})
			}
			continue
		}
		out[idx].Value = op.Value
	}
	return out, nil
}

func indexOfKey(entries []parser.Entry, key string) int {
	for i, e := range entries {
		if e.Key == key {
			return i
		}
	}
	return -1
}

func deleteKey(entries []parser.Entry, key string) []parser.Entry {
	out := entries[:0:len(entries)]
	for _, e := range entries {
		if e.Key != key {
			out = append(out, e)
		}
	}
	return out
}
