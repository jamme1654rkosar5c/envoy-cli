package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// AliasOptions controls behaviour of the Alias operation.
type AliasOptions struct {
	// Overwrite allows the alias key to replace an existing entry.
	Overwrite bool
	// KeepOriginal retains the source key after aliasing.
	KeepOriginal bool
	// DryRun skips mutation and returns what would change.
	DryRun bool
}

// DefaultAliasOptions returns sensible defaults.
func DefaultAliasOptions() AliasOptions {
	return AliasOptions{
		Overwrite:    false,
		KeepOriginal: true,
		DryRun:       false,
	}
}

// Alias creates a new key (alias) that copies the value of an existing key.
// The comment on the alias entry records its origin.
func Alias(entries []parser.Entry, srcKey, aliasKey string, opts AliasOptions) ([]parser.Entry, error) {
	var src *parser.Entry
	for i := range entries {
		if entries[i].Key == srcKey {
			src = &entries[i]
			break
		}
	}
	if src == nil {
		return nil, fmt.Errorf("alias: source key %q not found", srcKey)
	}

	for _, e := range entries {
		if e.Key == aliasKey {
			if !opts.Overwrite {
				return nil, fmt.Errorf("alias: key %q already exists; use Overwrite to replace", aliasKey)
			}
		}
	}

	if opts.DryRun {
		return entries, nil
	}

	aliasEntry := parser.Entry{
		Key:     aliasKey,
		Value:   src.Value,
		Comment: fmt.Sprintf("alias of %s", srcKey),
	}

	result := make([]parser.Entry, 0, len(entries)+1)
	for _, e := range entries {
		if e.Key == aliasKey && opts.Overwrite {
			continue
		}
		if e.Key == srcKey && !opts.KeepOriginal {
			result = append(result, aliasEntry)
			continue
		}
		result = append(result, e)
	}

	if opts.KeepOriginal || !strings.Contains(fmt.Sprint(result), aliasKey) {
		result = append(result, aliasEntry)
	}

	return result, nil
}
