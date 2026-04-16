// Package env provides utilities for working with collections of environment
// variable entries, including sorting and grouping operations.
//
// # Sorting
//
// Sort reorders entries alphabetically, reverse-alphabetically, or by key
// length using [SortOptions].
//
// # Grouping
//
// Group partitions entries by their key prefix (the segment before the first
// underscore). For example DB_HOST and DB_PORT both belong to the "DB" group.
// Entries with no underscore are placed in the catch-all "_" group.
//
// Use [GroupNames] to retrieve an optionally sorted list of group names from
// the resulting map.
package env
