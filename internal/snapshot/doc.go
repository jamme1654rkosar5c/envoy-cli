// Package snapshot provides functionality for saving and loading point-in-time
// snapshots of parsed .env files. Snapshots are stored as JSON files in a
// configurable directory and can be used to diff environment state over time
// or to audit changes to secrets and configuration values.
//
// Snapshot files are named using the format "<environment>-<timestamp>.json",
// where the timestamp is in RFC3339 format with colons replaced by hyphens for
// filesystem compatibility.
//
// Usage:
//
//	// Save a snapshot
//	err := snapshot.Save(".envoy/snapshots", "production", entries)
//
//	// List existing snapshots
//	paths, err := snapshot.List(".envoy/snapshots")
//
//	// Load a specific snapshot
//	snap, err := snapshot.Load(paths[0])
//
//	// Diff two snapshots
//	diff, err := snapshot.Diff(paths[0], paths[1])
package snapshot
