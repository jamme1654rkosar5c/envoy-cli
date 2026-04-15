// Package rename provides functionality for renaming keys across .env files.
package rename

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// Result holds the outcome of a rename operation.
type Result struct {
	OldKey  string
	NewKey  string
	File    string
	Renamed bool
}

// Options configures the rename behaviour.
type Options struct {
	// DryRun reports what would change without modifying the EnvFile.
	DryRun bool
	// ErrorIfNotFound returns an error when the old key is absent.
	ErrorIfNotFound bool
}

// DefaultOptions returns sensible rename defaults.
func DefaultOptions() Options {
	return Options{
		DryRun:          false,
		ErrorIfNotFound: true,
	}
}

// RenameKey renames oldKey to newKey inside the provided EnvFile.
// It returns a Result describing what happened and an error if the
// operation could not be completed.
func RenameKey(file *parser.EnvFile, oldKey, newKey string, opts Options) (Result, error) {
	result := Result{
		OldKey: oldKey,
		NewKey: newKey,
		File:   file.Path,
	}

	if strings.TrimSpace(newKey) == "" {
		return result, fmt.Errorf("rename: new key must not be empty")
	}

	// Check for duplicate target key.
	for _, e := range file.Entries {
		if e.Key == newKey {
			return result, fmt.Errorf("rename: key %q already exists in %s", newKey, file.Path)
		}
	}

	for i, e := range file.Entries {
		if e.Key == oldKey {
			if !opts.DryRun {
				file.Entries[i].Key = newKey
			}
			result.Renamed = true
			return result, nil
		}
	}

	if opts.ErrorIfNotFound {
		return result, fmt.Errorf("rename: key %q not found in %s", oldKey, file.Path)
	}
	return result, nil
}

// RenameKeyInAll applies RenameKey to every file in the slice.
// All results are collected; the first hard error aborts the run.
func RenameKeyInAll(files []*parser.EnvFile, oldKey, newKey string, opts Options) ([]Result, error) {
	results := make([]Result, 0, len(files))
	for _, f := range files {
		r, err := RenameKey(f, oldKey, newKey, opts)
		if err != nil {
			return results, err
		}
		results = append(results, r)
	}
	return results, nil
}
