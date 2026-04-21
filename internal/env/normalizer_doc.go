// Package env provides utilities for manipulating collections of environment
// variable entries represented as []parser.EnvEntry.
//
// # Normalizer
//
// The Normalize function applies a configurable set of transformations to a
// slice of EnvEntry values, returning a new, cleaned slice without mutating
// the original input.
//
// Supported transformations:
//
//   - UppercaseKeys  – converts every key to UPPER_CASE
//   - TrimValues     – strips leading and trailing whitespace from values
//   - StripQuotes    – removes surrounding single or double quotes from values
//   - RemoveEmpty    – drops entries whose value is empty (after other transforms)
//   - CollapseWhitespace – collapses runs of whitespace inside values to a
//     single space
//
// Example:
//
//	opts := env.DefaultNormalizeOptions()
//	opts.StripQuotes = true
//	cleaned := env.Normalize(entries, opts)
package env
