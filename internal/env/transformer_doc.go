// Package env provides utilities for sorting, grouping, filtering,
// deduplicating, and transforming environment variable entries.
//
// # Transformer
//
// The Transform function applies a TransformFunc to each entry's value,
// producing a new slice without mutating the original.
//
// Built-in transform functions:
//
//   - BuiltinUppercase  — converts values to UPPER CASE
//   - BuiltinLowercase  — converts values to lower case
//   - BuiltinTrimSpace  — strips leading/trailing whitespace
//
// Use TransformOptions to restrict which keys are affected:
//
//	opts := env.DefaultTransformOptions()
//	opts.OnlyKeys = []string{"APP_ENV"}
//	out := env.Transform(entries, env.BuiltinTrimSpace, opts)
package env
