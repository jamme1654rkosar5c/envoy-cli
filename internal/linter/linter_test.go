package linter_test

import (
	"testing"

	"github.com/envoy-cli/internal/linter"
	"github.com/envoy-cli/internal/parser"
)

func makeEnvFile(entries []parser.Entry) parser.EnvFile {
	return parser.EnvFile{Entries: entries}
}

func TestLint_NoIssues(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "PORT", Value: "8080"},
	})
	issues := linter.Lint(file, linter.DefaultOptions())
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_UppercaseKeys(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "app_name", Value: "myapp"},
		{Key: "Port", Value: "8080"},
	})
	opts := linter.DefaultOptions()
	issues := linter.Lint(file, opts)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	for _, issue := range issues {
		if issue.Rule != linter.RuleUppercaseKeys {
			t.Errorf("expected rule %s, got %s", linter.RuleUppercaseKeys, issue.Rule)
		}
	}
}

func TestLint_TrailingSpaceInValue(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "APP_NAME", Value: "myapp   "},
	})
	opts := linter.DefaultOptions()
	issues := linter.Lint(file, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Rule != linter.RuleNoTrailingSpace {
		t.Errorf("expected rule %s, got %s", linter.RuleNoTrailingSpace, issues[0].Rule)
	}
}

func TestLint_QuotedValues(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "DB_URL", Value: `"postgres://localhost"`},
		{Key: "SECRET", Value: "'abc123'"},
	})
	opts := linter.DefaultOptions()
	issues := linter.Lint(file, opts)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	for _, issue := range issues {
		if issue.Rule != linter.RuleNoQuotedValues {
			t.Errorf("expected rule %s, got %s", linter.RuleNoQuotedValues, issue.Rule)
		}
	}
}

func TestLint_EmptyValueWhenEnabled(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "OPTIONAL_KEY", Value: ""},
	})
	opts := linter.Options{NoEmptyValue: true}
	issues := linter.Lint(file, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Rule != linter.RuleNoEmptyValue {
		t.Errorf("expected rule %s, got %s", linter.RuleNoEmptyValue, issues[0].Rule)
	}
}

func TestLint_IssueString(t *testing.T) {
	issue := linter.Issue{
		Line:    3,
		Key:     "db_host",
		Rule:    linter.RuleUppercaseKeys,
		Message: "key should be uppercase",
	}
	got := issue.String()
	if got == "" {
		t.Error("expected non-empty issue string")
	}
}
