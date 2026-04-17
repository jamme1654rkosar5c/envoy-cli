package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError describes a problem found in an env file.
type ValidationError struct {
	Line    int
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d [%s]: %s", e.Line, e.Key, e.Message)
	}
	return fmt.Sprintf("[%s]: %s", e.Key, e.Message)
}

var validKey = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Validate checks an EnvFile for common issues and returns a list of errors.
func Validate(env *EnvFile) []ValidationError {
	var errs []ValidationError
	seen := make(map[string]int)

	for _, entry := range env.Entries {
		// Check key naming convention.
		if !validKey.MatchString(entry.Key) {
			errs = append(errs, ValidationError{
				Line:    entry.Line,
				Key:     entry.Key,
				Message: "key must be uppercase letters, digits, or underscores and start with a letter",
			})
		}

		// Check for empty values.
		if strings.TrimSpace(entry.Value) == "" {
			errs = append(errs, ValidationError{
				Line:    entry.Line,
				Key:     entry.Key,
				Message: "value is empty",
			})
		}

		// Check for duplicate keys.
		if prev, ok := seen[entry.Key]; ok {
			errs = append(errs, ValidationError{
				Line:    entry.Line,
				Key:     entry.Key,
				Message: fmt.Sprintf("duplicate key, first defined on line %d", prev),
			})
		} else {
			seen[entry.Key] = entry.Line
		}
	}

	return errs
}
