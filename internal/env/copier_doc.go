// Package env provides utilities for manipulating collections of environment
// variable entries parsed from .env files.
//
// # Copier
//
// The Copy function duplicates the value of one key into another within the
// same entry slice. It supports the following behaviours via CopyOptions:
//
//   - Overwrite: allow replacing an existing destination key.
//   - DryRun: validate the operation without modifying the slice.
//   - KeepSource: when false, the source key is removed after copying,
//     effectively performing a rename / move.
//
// Example:
//
//	opts := env.DefaultCopyOptions()
//	opts.KeepSource = false // move semantics
//	updated, err := env.Copy(entries, "OLD_KEY", "NEW_KEY", opts)
package env
