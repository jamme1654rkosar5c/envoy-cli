// Package env provides utilities for manipulating .env file entries.
//
// # Trimmer
//
// The Trim function applies whitespace and affix trimming to a slice of
// [parser.Entry] values. It is useful for normalising keys and values that
// may have been manually edited or imported from external sources.
//
// Basic usage:
//
//	opts := env.DefaultTrimOptions()
//	opts.TrimPrefixes = []string{"APP_"}
//	cleaned := env.Trim(entries, opts)
//
// Options:
//   - TrimKeys   – strip surrounding whitespace from each key
//   - TrimValues – strip surrounding whitespace from each value
//   - TrimPrefixes – remove a list of string prefixes from keys
//   - TrimSuffixes – remove a list of string suffixes from keys
//   - SkipEmpty  – leave entries with an empty value completely untouched
package env
