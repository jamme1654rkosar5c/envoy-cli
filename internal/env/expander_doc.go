// Package env provides utilities for manipulating collections of environment
// variable entries represented as []parser.Entry.
//
// # Expander
//
// Expand resolves shell-style variable references within entry values using
// other entries in the same slice as the lookup source.
//
// Supported reference syntaxes:
//
//	${VAR}   — brace-delimited reference
//	$VAR     — bare reference (identifier characters only)
//	$$       — escaped literal dollar sign, produces a single '$'
//
// Entries are processed in declaration order, so a key defined later can
// reference a key defined earlier, but not vice-versa.
//
// Example:
//
//	HOST=db.local
//	DSN=postgres://${HOST}/app   → postgres://db.local/app
//
// Options:
//
//	AllowMissing — when true, unresolvable references are left as-is instead
//	               of returning an error.
//	MaxDepth     — maximum recursive expansion depth (default 10). Prevents
//	               runaway expansion from self-referential values.
package env
