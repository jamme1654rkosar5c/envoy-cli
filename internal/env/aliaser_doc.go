// Package env provides utilities for manipulating collections of environment
// variable entries parsed from .env files.
//
// # Aliaser
//
// The Alias function creates a new key whose value mirrors an existing source
// key. This is useful when a dependency expects a different variable name than
// the one already defined in your .env file.
//
// Example:
//
//	out, err := env.Alias(entries, "DB_HOST", "DATABASE_HOST", env.DefaultAliasOptions())
//
// By default the original source key is retained. Set KeepOriginal = false to
// remove it after aliasing. Use Overwrite = true to replace an existing key
// with the same name as the alias. DryRun = true returns the unchanged slice
// without performing any mutation.
package env
