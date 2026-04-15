// Package snapshot provides functionality for saving and loading point-in-time
// snapshots of parsed .env files. Snapshots are stored as JSON files in a
// configurable directory and can be used to diff environment state over time
// or to audit changes to secrets and configuration values.
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
package snapshot
