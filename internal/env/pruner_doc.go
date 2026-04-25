// Package env provides utilities for manipulating collections of environment
// variable entries.
//
// # Pruner
//
// The Pruner removes unwanted entries from a slice of [parser.Entry] values
// based on configurable criteria:
//
//   - RemoveEmpty  – drop entries whose value is the empty string.
//   - RemoveCommented – drop entries whose value starts with '#', indicating a
//     commented-out or deprecated assignment.
//   - Keys – an explicit allow-list of key names to remove regardless of value.
//
// All options may be combined. When DryRun is true the original slice is never
// modified; a shallow copy of every surviving entry is returned instead.
//
// Example:
//
//	opts := env.DefaultPruneOptions()
//	opts.Keys = []string{"LEGACY_API_KEY"}
//	kept, pruned := env.Prune(entries, opts)
//	fmt.Printf("removed %d keys: %v\n", len(pruned), pruned)
package env
