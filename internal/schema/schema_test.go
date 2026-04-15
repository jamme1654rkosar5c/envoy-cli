package schema

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeEnvFile(entries []parser.Entry) parser.EnvFile {
	return parser.EnvFile{Entries: entries}
}

func TestEnforce_AllRequiredPresent(t *testing.T) {
	s := Schema{
		Keys: []KeySpec{
			{Key: "APP_ENV", Required: true},
			{Key: "PORT", Required: true},
		},
	}
	file := makeEnvFile([]parser.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	})
	errs := Enforce(s, file)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestEnforce_MissingRequiredKey(t *testing.T) {
	s := Schema{
		Keys: []KeySpec{
			{Key: "DATABASE_URL", Required: true},
		},
	}
	file := makeEnvFile([]parser.Entry{})
	errs := Enforce(s, file)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Key != "DATABASE_URL" {
		t.Errorf("expected key DATABASE_URL, got %s", errs[0].Key)
	}
}

func TestEnforce_OptionalKeyMissing_NoError(t *testing.T) {
	s := Schema{
		Keys: []KeySpec{
			{Key: "LOG_LEVEL", Required: false},
		},
	}
	file := makeEnvFile([]parser.Entry{})
	errs := Enforce(s, file)
	if len(errs) != 0 {
		t.Fatalf("expected no errors for optional missing key, got %v", errs)
	}
}

func TestEnforce_PatternMatch_Valid(t *testing.T) {
	s := Schema{
		Keys: []KeySpec{
			{Key: "API_URL", Required: true, Pattern: "https://*"},
		},
	}
	file := makeEnvFile([]parser.Entry{
		{Key: "API_URL", Value: "https://example.com"},
	})
	errs := Enforce(s, file)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestEnforce_PatternMatch_Invalid(t *testing.T) {
	s := Schema{
		Keys: []KeySpec{
			{Key: "API_URL", Required: true, Pattern: "https://*"},
		},
	}
	file := makeEnvFile([]parser.Entry{
		{Key: "API_URL", Value: "http://example.com"},
	})
	errs := Enforce(s, file)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Key != "API_URL" {
		t.Errorf("unexpected key in error: %s", errs[0].Key)
	}
}

func TestEnforce_MultipleErrors(t *testing.T) {
	s := Schema{
		Keys: []KeySpec{
			{Key: "APP_ENV", Required: true},
			{Key: "SECRET_KEY", Required: true},
		},
	}
	file := makeEnvFile([]parser.Entry{})
	errs := Enforce(s, file)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}
