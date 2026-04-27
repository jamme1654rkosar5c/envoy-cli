// Package env provides utilities for manipulating and transforming collections
// of environment variable entries.
//
// # Caster
//
// The Caster module normalises the string representation of env values by
// casting them to a specified target type. This is useful when env files
// contain values that are logically typed (integers, booleans, floats) but
// stored as raw strings with inconsistent formatting.
//
// Supported target types:
//
//   - CastString  – no-op; trims surrounding whitespace only.
//   - CastInt     – parses as a base-10 integer and re-serialises.
//   - CastFloat   – parses as a 64-bit float and re-serialises without
//     trailing zeros (e.g. "3.14000" → "3.14").
//   - CastBool    – accepts Go's strconv.ParseBool inputs ("1", "t", "TRUE",
//     "false", "0", etc.) and normalises to "true" or "false".
//
// Example usage:
//
//	opts := env.DefaultCastOptions()
//	opts.TargetType = env.CastBool
//	opts.Keys = []string{"DEBUG", "VERBOSE"}
//	out, results, err := env.Cast(entries, opts)
package env
