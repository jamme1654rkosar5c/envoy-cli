// Package env provides utilities for managing environment variable entries.
//
// # Freezer
//
// The freezer module allows individual keys to be "frozen", preventing
// accidental modification via Patch or other mutation operations.
//
// A frozen key is marked with a `frozen` token appended to its inline comment.
// This marker is human-readable in the .env file and survives round-trips
// through the parser.
//
// Usage:
//
//	out, err := env.Freeze(entries, "APP_SECRET", env.DefaultFreezeOptions())
//	if env.IsFrozen(out, "APP_SECRET") {
//	    // key is protected
//	}
//	out, err = env.Unfreeze(out, "APP_SECRET", env.DefaultFreezeOptions())
package env
