// Package compare provides cross-environment comparison for .env files.
//
// It allows users to compare the same keys across multiple named environments
// (e.g., dev, staging, prod) and identify:
//
//   - Keys present in some environments but missing in others
//   - Keys whose values differ between environments
//   - A full matrix view of all keys and their per-environment values
//
// # Usage
//
//	envs := compare.EnvMap{
//	    "dev":  devEntries,
//	    "prod": prodEntries,
//	}
//	report := compare.CrossCompare(envs)
//	fmt.Print(compare.FormatTable(report))
//	fmt.Print(compare.FormatMissing(report))
//
// The Report type exposes structured data suitable for further processing
// or rendering in CLI output.
package compare
