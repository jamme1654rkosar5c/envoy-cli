package env

import (
	"fmt"

	"github.com/envoy-cli/internal/parser"
)

// CopyOptions controls how key copying behaves.
type CopyOptions struct {
	Overwrite  bool
	DryRun     bool
	KeepSource bool
}

// DefaultCopyOptions returns sensible defaults.
func DefaultCopyOptions() CopyOptions {
	return CopyOptions{
		Overwrite:  false,
		DryRun:     false,
		KeepSource: true,
	}
}

// Copy duplicates the value of srcKey into dstKey within the same entries slice.
// If KeepSource is false, the source key is removed (move semantics).
func Copy(entries []parser.Entry, srcKey, dstKey string, opts CopyOptions) ([]parser.Entry, error) {
	srcIdx := -1
	dstIdx := -1

	for i, e := range entries {
		if e.Key == srcKey {
			srcIdx = i
		}
		if e.Key == dstKey {
			dstIdx = i
		}
	}

	if srcIdx == -1 {
		return nil, fmt.Errorf("copy: source key %q not found", srcKey)
	}

	if dstIdx != -1 && !opts.Overwrite {
		return nil, fmt.Errorf("copy: destination key %q already exists (use Overwrite)", dstKey)
	}

	if opts.DryRun {
		return entries, nil
	}

	result := make([]parser.Entry, len(entries))
	copy(result, entries)

	newEntry := parser.Entry{
		Key:   dstKey,
		Value: result[srcIdx].Value,
	}

	if dstIdx != -1 {
		result[dstIdx] = newEntry
	} else {
		result = append(result, newEntry)
	}

	if !opts.KeepSource {
		result = deleteKey(result, srcKey)
	}

	return result, nil
}
