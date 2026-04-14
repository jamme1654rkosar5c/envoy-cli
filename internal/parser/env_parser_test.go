package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDATABASE_URL=postgres://localhost/db\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(env.Entries))
	}
	if env.Entries[0].Key != "APP_ENV" || env.Entries[0].Value != "production" {
		t.Errorf("unexpected first entry: %+v", env.Entries[0])
	}
}

func TestParseFile_SkipsComments(t *testing.T) {
	path := writeTempEnv(t, "# This is a comment\nKEY=value\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(env.Entries))
	}
}

func TestParseFile_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # default port\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Entries[0].Value != "8080" {
		t.Errorf("expected value '8080', got %q", env.Entries[0].Value)
	}
	if env.Entries[0].Comment != "default port" {
		t.Errorf("expected comment 'default port', got %q", env.Entries[0].Comment)
	}
}

func TestParseFile_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")

	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestToMap(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := env.ToMap()
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("unexpected map contents: %v", m)
	}
}
