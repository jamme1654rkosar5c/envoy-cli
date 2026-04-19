// Package env provides utilities for managing, transforming, and inspecting
// collections of environment variable entries.
//
// # Archiver
//
// The Archiver module allows keys to be soft-retired by renaming them with a
// configurable prefix (default: "ARCHIVED_"). This is useful when deprecating
// environment variables without losing their values.
//
// Basic usage:
//
//	out, err := env.Archive(entries, []string{"OLD_API_KEY"}, env.DefaultArchiveOptions())
//
// Options:
//
//	Prefix          — prefix applied to archived key names (default "ARCHIVED_")
//	Timestamp       — when true, appends a unix timestamp to the prefix for uniqueness
//	RemoveOriginal  — when true (default), removes the original key after archiving
package env
