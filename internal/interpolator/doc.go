// Package interpolator resolves variable references embedded in .env file
// values, supporting both the ${VAR} and $VAR syntaxes commonly found in
// shell-style environment files.
//
// # Usage
//
//		file, _ := loader.LoadEnv(".env")
//		resolved, err := interpolator.Interpolate(file, interpolator.DefaultOptions())
//
// By default, Interpolate returns an error if a referenced variable cannot be
// found within the same file. Set Options.AllowMissing to true to leave
// unresolvable references untouched instead.
//
// # Supported syntax
//
//	  $VARNAME          — simple reference
//	  ${VARNAME}        — brace-delimited reference (recommended for clarity)
//
// Circular references are not detected; callers should validate their env
// files before interpolation when this is a concern.
package interpolator
