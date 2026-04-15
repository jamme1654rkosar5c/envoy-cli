package linter

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// Rule represents a linting rule with a name and description.
type Rule string

const (
	RuleNoTrailingSpace  Rule = "no-trailing-space"
	RuleUppercaseKeys    Rule = "uppercase-keys"
	RuleNoEmptyValue     Rule = "no-empty-value"
	RuleNoQuotedValues   Rule = "no-quoted-values"
)

// Issue represents a single linting violation.
type Issue struct {
	Line    int
	Key     string
	Rule    Rule
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("line %d [%s] %s: %s", i.Line, i.Rule, i.Key, i.Message)
}

// Options configures which rules are active.
type Options struct {
	NoTrailingSpace bool
	UppercaseKeys   bool
	NoEmptyValue    bool
	NoQuotedValues  bool
}

// DefaultOptions returns an Options with all rules enabled.
func DefaultOptions() Options {
	return Options{
		NoTrailingSpace: true,
		UppercaseKeys:   true,
		NoEmptyValue:    false,
		NoQuotedValues:  true,
	}
}

// Lint checks the entries in an EnvFile against the configured rules.
func Lint(file parser.EnvFile, opts Options) []Issue {
	var issues []Issue

	for i, entry := range file.Entries {
		lineNum := i + 1

		if opts.NoTrailingSpace {
			if strings.TrimRight(entry.Key, " ") != entry.Key {
				issues = append(issues, Issue{
					Line:    lineNum,
					Key:     entry.Key,
					Rule:    RuleNoTrailingSpace,
					Message: "key contains trailing whitespace",
				})
			}
			if strings.TrimRight(entry.Value, " ") != entry.Value {
				issues = append(issues, Issue{
					Line:    lineNum,
					Key:     entry.Key,
					Rule:    RuleNoTrailingSpace,
					Message: "value contains trailing whitespace",
				})
			}
		}

		if opts.UppercaseKeys {
			if entry.Key != strings.ToUpper(entry.Key) {
				issues = append(issues, Issue{
					Line:    lineNum,
					Key:     entry.Key,
					Rule:    RuleUppercaseKeys,
					Message: "key should be uppercase",
				})
			}
		}

		if opts.NoEmptyValue {
			if strings.TrimSpace(entry.Value) == "" {
				issues = append(issues, Issue{
					Line:    lineNum,
					Key:     entry.Key,
					Rule:    RuleNoEmptyValue,
					Message: "value is empty",
				})
			}
		}

		if opts.NoQuotedValues {
			v := entry.Value
			if (strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`)) ||
				(strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'")) {
				issues = append(issues, Issue{
					Line:    lineNum,
					Key:     entry.Key,
					Rule:    RuleNoQuotedValues,
					Message: "value should not be wrapped in quotes",
				})
			}
		}
	}

	return issues
}
