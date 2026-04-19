// Package env provides utilities for manipulating collections of environment
// variable entries.
//
// # Flattener
//
// The Flatten function normalises keys that use a multi-character separator
// (default "__") into single-underscore-joined keys. This is useful when
// consuming environment variables exported from systems that use double-
// underscore notation to represent hierarchical configuration, such as
// Docker Compose override files or Kubernetes ConfigMaps.
//
// Example
//
//	APP__DATABASE__HOST=db.internal  →  APP_DATABASE_HOST=db.internal
//
// Duplicate keys that resolve to the same flat key are deduplicated; the
// last entry wins, mirroring typical shell variable override semantics.
//
// Use FlattenOptions.Prefix to restrict flattening to a namespace, and
// FlattenOptions.StripPrefix to remove that prefix from the output keys.
package env
