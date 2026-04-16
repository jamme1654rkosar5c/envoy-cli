// Package env provides utilities for sorting, grouping, filtering,
// and deduplicating environment variable entries.
//
// # Deduplication
//
// Dedupe removes entries with duplicate keys from a slice of parser.Entry
// values. The caller controls which occurrence is retained via DedupeStrategy:
//
//   - KeepFirst (default) – the first occurrence is kept; subsequent
//     duplicates are discarded.
//   - KeepLast – the last occurrence overwrites earlier ones in-place,
//     preserving original ordering of unique keys.
//
// In both cases Dedupe returns the cleaned slice together with a list of
// the duplicate key names that were removed, which can be surfaced to the
// user as warnings or audit records.
package env
