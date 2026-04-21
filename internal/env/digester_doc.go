// Package env provides the Digest function for computing a deterministic
// cryptographic hash over a set of environment entries.
//
// # Overview
//
// Digest is useful for detecting whether an .env file has changed between
// deployments or CI runs without exposing the raw values. Two sets of entries
// with identical keys and values (in any order when SortKeys is true) will
// always produce the same hash.
//
// # Algorithms
//
// Both MD5 and SHA-256 are supported via the DigestAlgorithm type:
//
//	DigestMD5    – fast, 32-char hex output
//	DigestSHA256 – stronger, 64-char hex output (default)
//
// # Filtering
//
// Keys can be scoped using IncludeKeys (allowlist) or ExcludeKeys (denylist)
// inside DigestOptions. This is handy for ignoring volatile keys such as
// timestamps or request IDs when comparing snapshots.
//
// # Example
//
//	opts := env.DefaultDigestOptions()
//	opts.ExcludeKeys = []string{"BUILD_TIMESTAMP"}
//	hash, err := env.Digest(entries, opts)
package env
