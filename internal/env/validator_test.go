package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeValidateEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestValidateEntries_NoIssues(t *testing.T) {
	entries := makeValidateEntries("APP_NAME", "envoy", "APP_PORT", "8080")
	issues := ValidateEntries(entries, DefaultValidateOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestValidateEntries_RequireValues_Fails(t *testing.T) {
	entries := makeValidateEntries("APP_NAME", "", "APP_PORT", "8080")
	opts := DefaultValidateOptions()
	opts.RequireValues = true
	issues := ValidateEntries(entries, opts)
	if len(issues) != 1 || issues[0].Key != "APP_NAME" {
		t.Fatalf("expected 1 issue for APP_NAME, got %v", issues)
	}
}

func TestValidateEntries_ForbiddenKey(t *testing.T) {
	entries := makeValidateEntries("SECRET", "abc", "SAFE", "xyz")
	opts := DefaultValidateOptions()
	opts.ForbiddenKeys = []string{"SECRET"}
	issues := ValidateEntries(entries, opts)
	if len(issues) != 1 || issues[0].Key != "SECRET" {
		t.Fatalf("expected 1 forbidden-key issue, got %v", issues)
	}
}

func TestValidateEntries_MaxValueLength(t *testing.T) {
	entries := makeValidateEntries("KEY", "this-is-a-very-long-value")
	opts := DefaultValidateOptions()
	opts.MaxValueLength = 5
	issues := ValidateEntries(entries, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 length issue, got %v", issues)
	}
}

func TestValidateEntries_AllowedPrefixes_Pass(t *testing.T) {
	entries := makeValidateEntries("APP_NAME", "envoy", "APP_PORT", "8080")
	opts := DefaultValidateOptions()
	opts.AllowedPrefixes = []string{"APP_"}
	issues := ValidateEntries(entries, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestValidateEntries_AllowedPrefixes_Fail(t *testing.T) {
	entries := makeValidateEntries("APP_NAME", "envoy", "DB_HOST", "localhost")
	opts := DefaultValidateOptions()
	opts.AllowedPrefixes = []string{"APP_"}
	issues := ValidateEntries(entries, opts)
	if len(issues) != 1 || issues[0].Key != "DB_HOST" {
		t.Fatalf("expected 1 prefix issue for DB_HOST, got %v", issues)
	}
}
