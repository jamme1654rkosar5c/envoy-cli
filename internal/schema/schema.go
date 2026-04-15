// Package schema provides functionality for defining and enforcing
// a required set of keys for a given environment configuration.
package schema

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// KeySpec describes the requirements for a single environment key.
type KeySpec struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern for value validation
}

// Schema holds the full set of key specifications for an environment.
type Schema struct {
	Keys []KeySpec
}

// ValidationError represents a single schema violation.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// Enforce checks the provided EnvFile against the schema and returns
// a list of validation errors. An empty slice means the file is compliant.
func Enforce(s Schema, file parser.EnvFile) []ValidationError {
	present := make(map[string]string, len(file.Entries))
	for _, e := range file.Entries {
		present[e.Key] = e.Value
	}

	var errs []ValidationError

	for _, spec := range s.Keys {
		val, found := present[spec.Key]
		if !found {
			if spec.Required {
				errs = append(errs, ValidationError{
					Key:     spec.Key,
					Message: "required key is missing",
				})
			}
			continue
		}
		if spec.Pattern != "" {
			if !matchesPattern(spec.Pattern, val) {
				errs = append(errs, ValidationError{
					Key:     spec.Key,
					Message: fmt.Sprintf("value does not match pattern %q", spec.Pattern),
				})
			}
		}
	}

	return errs
}

// matchesPattern performs a simple glob-style prefix/suffix match.
// Supports patterns like "prefix_*" or "*_suffix" or exact strings.
func matchesPattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(value, strings.TrimPrefix(pattern, "*"))
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(value, strings.TrimSuffix(pattern, "*"))
	}
	return value == pattern
}
