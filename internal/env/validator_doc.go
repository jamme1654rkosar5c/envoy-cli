// Package env provides utilities for manipulating and inspecting .env file
// entries in memory.
//
// # ValidateEntries
//
// ValidateEntries performs extended validation on a slice of parsed entries
// beyond what the core parser checks. It supports the following rules:
//
//   - RequireValues: every entry must have a non-empty value.
//   - ForbiddenKeys: specific keys that must not appear in the file.
//   - MaxValueLength: values must not exceed a given character count.
//   - AllowedPrefixes: every key must start with one of the listed prefixes.
//
// Example:
//
//	opts := env.DefaultValidateOptions()
//	opts.RequireValues = true
//	opts.AllowedPrefixes = []string{"APP_", "DB_"}
//	issues := env.ValidateEntries(entries, opts)
//	for _, issue := range issues {
//		fmt.Println(issue)
//	}
package env
