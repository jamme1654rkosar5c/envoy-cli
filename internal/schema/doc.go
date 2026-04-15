// Package schema provides schema enforcement for .env files.
//
// A Schema defines the expected keys for an environment, including
// whether each key is required and an optional pattern the value must match.
//
// Example usage:
//
//	s := schema.Schema{
//		Keys: []schema.KeySpec{
//			{Key: "APP_ENV",      Required: true},
//			{Key: "DATABASE_URL", Required: true, Pattern: "postgres://*"},
//			{Key: "LOG_LEVEL",   Required: false},
//		},
//	}
//
//	errs := schema.Enforce(s, envFile)
//	for _, e := range errs {
//		fmt.Println(e)
//	}
//
// Pattern matching supports simple glob-style expressions:
//   - "prefix_*" matches any value starting with "prefix_"
//   - "*_suffix" matches any value ending with "_suffix"
//   - "*"        matches any non-empty value
//   - exact string for strict equality
package schema
