// Package env provides utilities for manipulating and transforming .env file
// entries in memory.
//
// # Cascader
//
// The Cascade function applies multiple layers of env entries in order,
// simulating environment inheritance (e.g. base → staging → local).
//
// Each subsequent layer can override values from the previous one.
// Behaviour is controlled via CascadeOptions:
//
//   - Overwrite: if true, later layers replace values from earlier ones.
//   - SkipEmpty: if true, entries with empty values in source layers are ignored.
//
// Example:
//
//	result := env.Cascade([][]parser.Entry{base, overlay}, env.DefaultCascadeOptions())
package env
