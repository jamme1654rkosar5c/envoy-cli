package secrets

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeEnvFile(entries []parser.EnvEntry) parser.EnvFile {
	return parser.EnvFile{Path: ".env", Entries: entries}
}

func TestIsSensitive_MatchesSecretKeys(t *testing.T) {
	sensitiveKeys := []string{
		"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN",
		"PRIVATE_KEY", "AWS_ACCESS_KEY", "APP_CREDENTIALS",
	}
	for _, key := range sensitiveKeys {
		if !IsSensitive(key) {
			t.Errorf("expected key %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_AllowsSafeKeys(t *testing.T) {
	safeKeys := []string{"APP_ENV", "PORT", "DEBUG", "LOG_LEVEL", "BASE_URL"}
	for _, key := range safeKeys {
		if IsSensitive(key) {
			t.Errorf("expected key %q to NOT be sensitive", key)
		}
	}
}

func TestRedact_ReplacesSecretValues(t *testing.T) {
	file := makeEnvFile([]parser.EnvEntry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "API_TOKEN", Value: "tok_abc123"},
	})

	result := Redact(file)

	if result.Entries[0].Value != "production" {
		t.Errorf("expected APP_ENV to be unchanged, got %q", result.Entries[0].Value)
	}
	if result.Entries[1].Value != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result.Entries[1].Value)
	}
	if result.Entries[2].Value != "[REDACTED]" {
		t.Errorf("expected API_TOKEN to be redacted, got %q", result.Entries[2].Value)
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	file := makeEnvFile([]parser.EnvEntry{
		{Key: "DB_PASSWORD", Value: "original"},
	})
	Redact(file)
	if file.Entries[0].Value != "original" {
		t.Error("Redact mutated the original EnvFile")
	}
}

func TestRedactValue_MasksLongValues(t *testing.T) {
	result := RedactValue("supersecret")
	if result != "sup********" {
		t.Errorf("unexpected masked value: %q", result)
	}
}

func TestRedactValue_MasksShortValues(t *testing.T) {
	result := RedactValue("ab")
	if result != "**" {
		t.Errorf("unexpected masked value: %q", result)
	}
}
