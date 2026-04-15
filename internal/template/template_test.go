package template

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeEnvFile(entries []parser.Entry) parser.EnvFile {
	return parser.EnvFile{Entries: entries}
}

func TestGenerate_DefaultPlaceholders(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret123"},
	})

	result := Generate(file, nil)

	if !strings.Contains(result, "APP_NAME=<APP_NAME>") {
		t.Errorf("expected placeholder for APP_NAME, got:\n%s", result)
	}
	if !strings.Contains(result, "DB_PASSWORD=<DB_PASSWORD>") {
		t.Errorf("expected placeholder for DB_PASSWORD, got:\n%s", result)
	}
}

func TestGenerate_IncludeDefaults_NonSensitive(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "API_SECRET", Value: "topsecret"},
	})

	opts := &Options{IncludeDefaults: true, PlaceholderFormat: "<%s>"}
	result := Generate(file, opts)

	if !strings.Contains(result, "APP_ENV=production") {
		t.Errorf("expected default value for APP_ENV, got:\n%s", result)
	}
	if !strings.Contains(result, "API_SECRET=<API_SECRET>") {
		t.Errorf("sensitive key should still be a placeholder, got:\n%s", result)
	}
}

func TestGenerate_CustomPlaceholderFormat(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "PORT", Value: "8080"},
	})

	opts := &Options{PlaceholderFormat: "YOUR_%s_HERE"}
	result := Generate(file, opts)

	if !strings.Contains(result, "PORT=YOUR_PORT_HERE") {
		t.Errorf("expected custom placeholder format, got:\n%s", result)
	}
}

func TestGenerate_SensitiveKeyVariants(t *testing.T) {
	sensitiveKeys := []string{"DB_PASSWORD", "AUTH_TOKEN", "AWS_SECRET", "PRIVATE_KEY", "API_KEY"}

	for _, key := range sensitiveKeys {
		if !isSensitiveKey(key) {
			t.Errorf("expected %q to be detected as sensitive", key)
		}
	}
}

func TestGenerate_EmptyFile(t *testing.T) {
	file := makeEnvFile([]parser.Entry{})
	result := Generate(file, nil)

	if result != "" {
		t.Errorf("expected empty output for empty file, got: %q", result)
	}
}

func TestGenerate_EntryWithComment(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "LOG_LEVEL", Value: "info", Comment: "logging verbosity"},
	})

	result := Generate(file, nil)

	if !strings.Contains(result, "# logging verbosity") {
		t.Errorf("expected comment in output, got:\n%s", result)
	}
	if !strings.Contains(result, "LOG_LEVEL=<LOG_LEVEL>") {
		t.Errorf("expected placeholder for LOG_LEVEL, got:\n%s", result)
	}
}
