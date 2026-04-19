// Package env provides utilities for manipulating collections of environment
// file entries.
//
// # Rotator
//
// The Rotate function handles key rotation: renaming one or more keys while
// optionally replacing their values in a single pass. This is useful when
// migrating environment variable naming conventions or rotating secrets.
//
// Example usage:
//
//	rotations := []env.RotateEntry{
//		{OldKey: "DB_PASS", NewKey: "DB_PASSWORD", NewValue: "newSecret123"},
//	}
//	updated, err := env.Rotate(entries, rotations, env.DefaultRotateOptions())
//
// Options:
//   - DryRun: simulate the rotation without modifying the original slice.
//   - ErrorOnMissing: return an error when a source key cannot be found.
package env
