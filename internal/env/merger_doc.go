// Package env provides utilities for manipulating collections of environment
// variable entries.
//
// # Overlap
//
// Overlap merges two slices of entries (dst and src) into a new slice.
// Keys present in src but absent in dst are appended.
// Keys present in both are handled according to OverlapOptions:
//
//   - Overwrite: when true, the src value replaces the dst value.
//   - SkipEmpty: when true, src entries with empty values are ignored entirely,
//     preventing accidental erasure of existing values.
//
// The original dst and src slices are never mutated.
package env
