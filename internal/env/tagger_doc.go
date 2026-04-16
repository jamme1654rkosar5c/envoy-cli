// Package env provides utilities for sorting, grouping, filtering,
// deduplicating, transforming, patching, pinning, and tagging environment
// variable entries.
//
// # Tagger
//
// The tagger module allows attaching arbitrary string tags to individual
// entries via structured inline comments. Tags are stored as comment
// annotations using a configurable prefix (default: "@tag").
//
// Example usage:
//
//	opts := env.DefaultTagOptions()
//	out, err := env.Tag(entries, "APP_ENV", "stable", opts)
//	value := env.GetTag(out, "APP_ENV", opts) // "stable"
//	out = env.Untag(out, "APP_ENV", opts)
//
// Tags survive round-trip serialisation via the exporter as long as
// the dot-env format is used, since comments are preserved.
package env
