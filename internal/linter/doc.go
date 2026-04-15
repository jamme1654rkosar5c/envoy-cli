// Package linter provides style and consistency linting for .env files.
//
// It checks entries in a parsed EnvFile against configurable rules such as:
//   - no-trailing-space: keys and values must not have trailing whitespace
//   - uppercase-keys: all keys must be fully uppercase
//   - no-empty-value: values must not be blank (opt-in)
//   - no-quoted-values: values must not be wrapped in single or double quotes
//
// Usage:
//
//	file, _ := parser.ParseFile("path/to/.env")
//	opts := linter.DefaultOptions()
//	issues := linter.Lint(file, opts)
//	for _, issue := range issues {
//		fmt.Println(issue)
//	}
package linter
