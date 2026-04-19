package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// ValidateOptions controls validation behaviour.
type ValidateOptions struct {
	RequireValues  bool
	AllowedPrefixes []string
	ForbiddenKeys  []string
	MaxValueLength int
}

// DefaultValidateOptions returns sensible defaults.
func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		RequireValues:  false,
		MaxValueLength: 0,
	}
}

// ValidationIssue describes a single validation problem.
type ValidationIssue struct {
	Key     string
	Message string
}

func (v ValidationIssue) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// ValidateEntries runs entry-level checks beyond the basic parser validation.
func ValidateEntries(entries []parser.Entry, opts ValidateOptions) []ValidationIssue {
	var issues []ValidationIssue

	forbidden := make(map[string]bool, len(opts.ForbiddenKeys))
	for _, k := range opts.ForbiddenKeys {
		forbidden[k] = true
	}

	for _, e := range entries {
		if forbidden[e.Key] {
			issues = append(issues, ValidationIssue{Key: e.Key, Message: "key is forbidden"})
		}

		if opts.RequireValues && strings.TrimSpace(e.Value) == "" {
			issues = append(issues, ValidationIssue{Key: e.Key, Message: "value is required but empty"})
		}

		if opts.MaxValueLength > 0 && len(e.Value) > opts.MaxValueLength {
			issues = append(issues, ValidationIssue{
				Key:     e.Key,
				Message: fmt.Sprintf("value exceeds max length %d", opts.MaxValueLength),
			})
		}

		if len(opts.AllowedPrefixes) > 0 {
			matched := false
			for _, p := range opts.AllowedPrefixes {
				if strings.HasPrefix(e.Key, p) {
					matched = true
					break
				}
			}
			if !matched {
				issues = append(issues, ValidationIssue{
					Key:     e.Key,
					Message: "key does not match any allowed prefix",
				})
			}
		}
	}
	return issues
}
