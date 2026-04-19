// Package env provides utilities for manipulating .env file entries.
//
// # Promoter
//
// The Promote function copies entries from a source environment into a
// destination environment, modelling a typical promotion workflow such
// as copying variables from staging into production.
//
// Usage:
//
//	opts := env.DefaultPromoteOptions()
//	opts.Overwrite = true          // replace existing keys
//	opts.Keys = []string{"DB_URL"} // promote only specific keys
//
//	updated, err := env.Promote(staging, production, opts)
//
// DryRun mode returns the result slice without mutating the destination,
// which is useful for previewing changes before committing them.
package env
