// Package env provides utilities for manipulating collections of env entries.
//
// # Stasher
//
// The stasher module provides a lightweight named-stash mechanism for
// env entry slices. It is analogous to `git stash` — you can save the
// current state of an entry list under a name, continue working, and
// later pop the stash back to restore those entries.
//
// Usage:
//
//	store := map[string]env.StashEntry{}
//	opts  := env.DefaultStashOptions()
//
//	// Save current entries
//	env.Stash("pre-deploy", entries, store, opts)
//
//	// ... mutate entries ...
//
//	// Restore stashed entries (new keys only)
//	restored, err := env.Pop("pre-deploy", entries, store, opts)
//
// Stash names must be unique unless AllowOverwrite is set.
// Pop removes the stash from the store after retrieval.
// When RestoreOnPop is true, only keys absent from dst are merged in.
package env
