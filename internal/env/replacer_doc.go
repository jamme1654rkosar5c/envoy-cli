// Package env provides utilities for manipulating collections of environment
// variable entries.
//
// # Replacer
//
// The Replace function performs bulk value substitution across a slice of
// [parser.Entry] values. It is useful when migrating infrastructure (e.g.
// swapping a hostname across all entries) or sanitising values before export.
//
// Basic usage:
//
//	result, count, err := env.Replace(entries, env.ReplaceOptions{
//		OldValue: "localhost",
//		NewValue: "prod.example.com",
//	})
//
// Key filtering restricts substitutions to entries whose key contains a given
// substring:
//
//	opts.KeyFilter = "DB_"
//
// Exact matching requires the entire value to equal OldValue:
//
//	opts.ExactMatch = true
//
// Dry-run mode returns a mutated copy without altering the original slice:
//
//	opts.DryRun = true
package env
