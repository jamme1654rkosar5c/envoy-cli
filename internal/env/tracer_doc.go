// Package env provides utilities for managing and transforming environment
// variable entries.
//
// # Tracer
//
// The Tracer resolves the full reference chain for each key in a set of
// environment entries. Given entries where values may reference other keys
// via $KEY or ${KEY} syntax, Trace walks each reference chain and records:
//
//   - The ordered chain of keys visited (e.g. A -> B -> C)
//   - The depth of the chain (number of hops)
//   - Whether a cycle was detected
//
// Example:
//
//	results, err := env.Trace(entries, env.DefaultTraceOptions())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(env.FormatTrace(results, keys))
//
// Cycle detection is enabled by default and returns an error when a cycle is
// found. Set AllowCycles: true in TraceOptions to record cycles without
// failing.
package env
