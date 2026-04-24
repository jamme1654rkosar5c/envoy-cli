// Package env provides utilities for managing environment variable entries.
//
// # Differ
//
// The Diff function compares two slices of EnvEntry and produces a list of
// DiffEntry values that describe what changed between a base set and a next
// (updated) set of entries.
//
// Supported change statuses:
//
//   - added    — key exists in next but not in base
//   - removed  — key exists in base but not in next
//   - changed  — key exists in both but the value differs
//   - unchanged — key exists in both with the same value
//
// # Options
//
// DiffOptions.MaskSecrets will redact values for keys whose names contain
// common sensitive keywords (e.g. "secret", "password", "token").
//
// DiffOptions.IgnoreKeys accepts a list of key names to exclude entirely from
// the comparison.
//
// # Formatting
//
// FormatDiff renders a tabular summary of all diff entries.
// BuildDiffSummaries returns a compact slice of annotated strings.
package env
